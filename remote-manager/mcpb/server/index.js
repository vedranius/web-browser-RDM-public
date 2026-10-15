#!/usr/bin/env node
// WRM PRO stdio bridge for Claude Desktop (MCP bundle, .mcpb).
//
// Claude Desktop starts local MCP servers only (stdio). This file is the server of the
// "WRM PRO" desktop extension that WRM generates in Settings → AI connections: it reads
// newline-delimited JSON-RPC messages from stdin, POSTs each one to WRM's /mcp endpoint
// with the bearer token, and writes every message of the answer (JSON or an event stream,
// including WRM's approval questions) to stdout, one per line. It mirrors `wrm -mcp-stdio`
// (mcp_stdio.go) so the user needs no WRM program on the computer: Claude Desktop ships
// the Node.js runtime this file needs. No dependencies; only protocol messages go to
// stdout, diagnostics go to stderr.
//
// Environment: WRM_MCP_URL (the /mcp address), WRM_MCP_TOKEN (wrm_pat_…), optional
// WRM_MCP_CA_FILE (PEM file to trust a private CA or a self-signed WRM certificate).
'use strict';

const http = require('http');
const https = require('https');
const fs = require('fs');
const readline = require('readline');

const VERSION = '0.0.0'; // replaced with the WRM version when the bundle is generated
const MAX_BODY = 32 << 20;

let url = (process.env.WRM_MCP_URL || '').trim();
const token = (process.env.WRM_MCP_TOKEN || '').trim();
let caFile = (process.env.WRM_MCP_CA_FILE || '').trim();
if (caFile.startsWith('${')) caFile = ''; // an optional setting the app left unfilled

function log(text) { process.stderr.write('wrm-mcp-bridge: ' + text + '\n'); }

if (!url || !token) {
  log('WRM_MCP_URL and WRM_MCP_TOKEN are required (set them in the extension settings).');
  process.exit(2);
}
url = url.replace(/\/+$/, '');
if (!/\/mcp$/.test(url)) url += '/mcp';
let target;
try { target = new URL(url); } catch (e) { log('invalid WRM address: ' + url); process.exit(2); }
const lib = target.protocol === 'https:' ? https : http;
let ca;
if (caFile) {
  try { ca = fs.readFileSync(caFile); } catch (e) { log(e.message); process.exit(2); }
}

let session = '';
let protocolVersion = '';

function write(obj) { process.stdout.write((typeof obj === 'string' ? obj.trim() : JSON.stringify(obj)) + '\n'); }
function rpcError(id, message) { return {jsonrpc: '2.0', id: id === undefined ? null : id, error: {code: -32000, message}}; }

function request(method, body) {
  return new Promise((resolve, reject) => {
    const headers = {
      'Authorization': 'Bearer ' + token,
      'Accept': 'application/json, text/event-stream',
      'User-Agent': 'wrm-mcpb/' + VERSION,
    };
    if (body) { headers['Content-Type'] = 'application/json'; headers['Content-Length'] = Buffer.byteLength(body); }
    if (session) headers['Mcp-Session-Id'] = session;
    if (protocolVersion) headers['MCP-Protocol-Version'] = protocolVersion;
    const opts = {method, headers};
    if (ca) opts.ca = ca;
    const req = lib.request(target, opts, resolve);
    req.on('error', reject);
    if (body) req.write(body);
    req.end();
  });
}

function readAll(res, limit) {
  return new Promise((resolve, reject) => {
    const chunks = [];
    let n = 0;
    res.on('data', c => { n += c.length; if (n <= limit) chunks.push(c); });
    res.on('end', () => resolve(Buffer.concat(chunks).toString('utf8')));
    res.on('error', reject);
  });
}

// copyEvents writes the data of every server-sent event to stdout as it arrives (an
// approval question must reach the app while the tool call still waits).
function copyEvents(res) {
  return new Promise(resolve => {
    res.setEncoding('utf8');
    let buf = '';
    let data = [];
    const flush = () => { if (data.length) { write(data.join('\n')); data = []; } };
    const line = l => {
      if (l.endsWith('\r')) l = l.slice(0, -1);
      if (l === '') flush();
      else if (l.startsWith('data:')) data.push(l.slice(5).replace(/^ /, ''));
    };
    res.on('data', chunk => {
      buf += chunk;
      let i;
      while ((i = buf.indexOf('\n')) >= 0) { line(buf.slice(0, i)); buf = buf.slice(i + 1); }
    });
    const end = () => { if (buf) line(buf); flush(); resolve(); };
    res.on('end', end);
    res.on('error', end);
  });
}

async function forward(raw, msg) {
  const isRequest = msg.id !== undefined && msg.id !== null && typeof msg.method === 'string';
  const fail = text => { log(text); if (isRequest) write(rpcError(msg.id, text)); };
  let res;
  try { res = await request('POST', raw); } catch (e) { fail('WRM is not reachable at ' + url + ': ' + e.message); return; }
  if (msg.method === 'initialize' && res.statusCode === 200) session = res.headers['mcp-session-id'] || '';
  const ct = String(res.headers['content-type'] || '');
  if (res.statusCode === 202) { res.resume(); return; }
  if (ct.startsWith('text/event-stream')) { await copyEvents(res); return; }
  const body = await readAll(res, ct.startsWith('application/json') ? MAX_BODY : 4096);
  if (res.statusCode !== 200 || !ct.startsWith('application/json')) {
    let m = body.trim();
    try { m = JSON.parse(body).error.message || m; } catch (e) { /* plain text */ }
    if (res.statusCode === 401) m += ' (the token is wrong, expired or revoked: create a new one in WRM → Settings → AI connections)';
    fail('WRM answered HTTP ' + res.statusCode + ': ' + m);
    return;
  }
  if (msg.method === 'initialize') {
    try { protocolVersion = JSON.parse(body).result.protocolVersion || ''; } catch (e) { /* keep the default */ }
  }
  write(body);
}

const pending = new Set();
let chain = Promise.resolve();
const rl = readline.createInterface({input: process.stdin, crlfDelay: Infinity});
rl.on('line', line => {
  const raw = line.trim();
  if (!raw) return;
  let msg;
  try { msg = JSON.parse(raw); } catch (e) { write({jsonrpc: '2.0', id: null, error: {code: -32700, message: 'parse error'}}); return; }
  if (msg.method === 'initialize') {
    // nothing else may run before the session exists
    chain = Promise.all([...pending]).then(() => forward(raw, msg));
    track(chain);
    return;
  }
  track(chain.then(() => forward(raw, msg)));
});
function track(p) {
  const q = p.catch(e => log(String(e && e.message || e))).finally(() => pending.delete(q));
  pending.add(q);
}
rl.on('close', async () => {
  await Promise.all([...pending]);
  if (session) {
    try { const res = await request('DELETE'); res.resume(); } catch (e) { /* ending anyway */ }
  }
  process.exit(0);
});
