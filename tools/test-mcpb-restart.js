#!/usr/bin/env node
'use strict';

// A stale MCP session after WRM restarts must reconnect with the saved token instead of
// sending the user through wrm_connect. Only WRM's explicit "session not found" 404 is retried.
const assert = require('node:assert/strict');
const {spawn} = require('node:child_process');
const http = require('node:http');
const path = require('node:path');
const readline = require('node:readline');
const {test} = require('node:test');

const bridge = path.join(__dirname, '..', 'remote-manager', 'mcpb', 'server', 'index.js');

test('Claude bridge reinitializes after WRM restart without reauth or server reselection', async t => {
  let activeSession = '';
  let initCount = 0;
  let listCount = 0;
  const server = http.createServer(async (req, res) => {
    const chunks = [];
    for await (const chunk of req) chunks.push(chunk);
    const body = chunks.length ? JSON.parse(Buffer.concat(chunks).toString('utf8')) : {};
    if (req.headers.authorization !== 'Bearer wrm_pat_test') {
      res.writeHead(401, {'Content-Type': 'application/json'});
      res.end(JSON.stringify({error: 'invalid bearer'}));
      return;
    }
    if (body.method === 'initialize') {
      activeSession = 'session-' + (++initCount);
      res.writeHead(200, {'Content-Type': 'application/json', 'Mcp-Session-Id': activeSession});
      res.end(JSON.stringify({jsonrpc: '2.0', id: body.id, result: {protocolVersion: '2025-06-18', capabilities: {tools: {}}, serverInfo: {name: 'test', version: '1'}}}));
      return;
    }
    if (body.method === 'notifications/initialized') { res.writeHead(202); res.end(); return; }
    if (body.method === 'tools/list') {
      listCount++;
      if (req.headers['mcp-session-id'] !== activeSession) {
        res.writeHead(404, {'Content-Type': 'application/json'});
        res.end(JSON.stringify({error: {message: 'session not found: initialize again'}}));
        return;
      }
      if (listCount === 1) {
        activeSession = ''; // simulate the process restart invalidating its in-memory sessions
        res.writeHead(404, {'Content-Type': 'application/json'});
        res.end(JSON.stringify({error: {message: 'session not found: initialize again'}}));
        return;
      }
      res.writeHead(200, {'Content-Type': 'application/json'});
      res.end(JSON.stringify({jsonrpc: '2.0', id: body.id, result: {tools: [{name: 'list_connections'}]}}));
      return;
    }
    res.writeHead(404, {'Content-Type': 'text/plain'});
    res.end('Not Found');
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  t.after(() => new Promise(resolve => server.close(resolve)));

  const addr = server.address();
  const child = spawn(process.execPath, [bridge], {
    env: {...process.env, WRM_MCP_URL: `http://127.0.0.1:${addr.port}/mcp`, WRM_MCP_TOKEN: 'wrm_pat_test', WRM_MCP_CONFIG: path.join(__dirname, '.missing-mcp-test-config')},
    stdio: ['pipe', 'pipe', 'pipe'],
  });
  t.after(() => { if (!child.killed) child.kill(); });
  let stderr = '';
  child.stderr.setEncoding('utf8');
  child.stderr.on('data', chunk => { stderr += chunk; });

  const pending = new Map();
  const lines = readline.createInterface({input: child.stdout});
  lines.on('line', line => {
    let message;
    try { message = JSON.parse(line); } catch (e) { return; }
    if (message.id !== undefined && pending.has(String(message.id))) {
      const resolve = pending.get(String(message.id));
      pending.delete(String(message.id));
      resolve(message);
    }
  });
  function send(message) {
    return new Promise((resolve, reject) => {
      pending.set(String(message.id), resolve);
      child.stdin.write(JSON.stringify(message) + '\n', err => { if (err) reject(err); });
      setTimeout(() => {
        if (pending.has(String(message.id))) {
          pending.delete(String(message.id));
          reject(new Error(`timed out waiting for ${message.method}; stderr: ${stderr}`));
        }
      }, 5000).unref();
    });
  }

  const initialized = await send({jsonrpc: '2.0', id: 1, method: 'initialize', params: {protocolVersion: '2025-06-18', capabilities: {}, clientInfo: {name: 'restart-test'}}});
  assert.equal(initialized.result.serverInfo.name, 'test');
  child.stdin.write(JSON.stringify({jsonrpc: '2.0', method: 'notifications/initialized'}) + '\n');

  const tools = await send({jsonrpc: '2.0', id: 2, method: 'tools/list'});
  assert.equal(tools.result.tools[0].name, 'list_connections');
  assert.equal(initCount, 2, 'the bridge should make one internal initialize request after the stale-session response');
  assert.equal(listCount, 2, 'the rejected list request should be retried once with the new session');

  const missing = await send({jsonrpc: '2.0', id: 3, method: 'resources/list'});
  assert.equal(missing.error.code, -32000);
  assert.equal(initCount, 2, 'an unrelated 404 must not trigger reinitialization');
});
