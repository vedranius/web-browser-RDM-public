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

// An unreachable WRM must not send the user through a new approval while the saved access is
// still valid: the bridge says so, reconnects by itself, and wrm_connect does not start a new
// sign-in. Only a refused token (401) asks for wrm_connect.
test('Claude bridge keeps a valid saved access while WRM is unreachable', async t => {
  const fs = require('node:fs');
  const os = require('node:os');
  let tokenOK = true;
  const server = http.createServer(async (req, res) => {
    const chunks = [];
    for await (const chunk of req) chunks.push(chunk);
    const body = chunks.length ? JSON.parse(Buffer.concat(chunks).toString('utf8')) : {};
    if (!tokenOK || req.headers.authorization !== 'Bearer wrm_oat_saved') {
      res.writeHead(401, {'Content-Type': 'application/json'});
      res.end(JSON.stringify({jsonrpc: '2.0', id: null, error: {code: -32001, message: 'the token has expired'}}));
      return;
    }
    if (body.method === 'initialize') {
      res.writeHead(200, {'Content-Type': 'application/json', 'Mcp-Session-Id': 's1'});
      res.end(JSON.stringify({jsonrpc: '2.0', id: body.id, result: {protocolVersion: '2025-06-18', capabilities: {tools: {}}, serverInfo: {name: 'wrm', version: '1'}}}));
      return;
    }
    if (body.method === 'notifications/initialized') { res.writeHead(202); res.end(); return; }
    res.writeHead(200, {'Content-Type': 'application/json'});
    res.end(JSON.stringify({jsonrpc: '2.0', id: body.id, result: {tools: [{name: 'list_connections'}]}}));
  });
  // reserve a port, then leave it closed so WRM is "down" when the bridge starts
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  const port = server.address().port;
  await new Promise(resolve => server.close(resolve));

  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'wrm-mcpb-'));
  t.after(() => fs.rmSync(dir, {recursive: true, force: true}));
  const cfg = path.join(dir, 'mcp.json');
  fs.writeFileSync(cfg, JSON.stringify({url: `http://127.0.0.1:${port}/mcp`, token: 'wrm_oat_saved', expires_at: new Date(Date.now() + 7 * 864e5).toISOString()}));
  const env = {...process.env, WRM_MCP_CONFIG: cfg, WRM_MCP_RECONNECT_MS: '1', WRM_MCP_NO_BROWSER: '1', WRM_MCP_CONNECT_WAIT_MS: '200'};
  delete env.WRM_MCP_URL;
  delete env.WRM_MCP_TOKEN;
  const child = spawn(process.execPath, [bridge], {env, stdio: ['pipe', 'pipe', 'pipe']});
  t.after(() => { if (!child.killed) child.kill(); });
  child.stderr.resume();
  const pending = new Map();
  const notes = [];
  readline.createInterface({input: child.stdout}).on('line', line => {
    let m;
    try { m = JSON.parse(line); } catch (e) { return; }
    if (m.id !== undefined && pending.has(String(m.id))) { pending.get(String(m.id))(m); pending.delete(String(m.id)); } else if (m.method) notes.push(m.method);
  });
  const send = m => new Promise((resolve, reject) => {
    pending.set(String(m.id), resolve);
    child.stdin.write(JSON.stringify(m) + '\n');
    setTimeout(() => { if (pending.delete(String(m.id))) reject(new Error('timed out waiting for ' + m.method)); }, 5000).unref();
  });

  const init = await send({jsonrpc: '2.0', id: 1, method: 'initialize', params: {protocolVersion: '2025-06-18', capabilities: {}, clientInfo: {name: 'down-test'}}});
  assert.match(init.result.instructions, /do NOT call wrm_connect/);
  const call = await send({jsonrpc: '2.0', id: 2, method: 'tools/call', params: {name: 'list_connections', arguments: {}}});
  assert.equal(call.result.isError, true);
  assert.match(call.result.content[0].text, /saved sign-in is still valid/);
  const connect = await send({jsonrpc: '2.0', id: 3, method: 'tools/call', params: {name: 'wrm_connect', arguments: {}}});
  assert.match(connect.result.content[0].text, /do NOT call wrm_connect/, 'no new sign-in while the saved access is valid');

  // WRM is back: the next request reconnects with the saved token
  server.listen(port, '127.0.0.1');
  await new Promise(resolve => server.once('listening', resolve));
  t.after(() => new Promise(resolve => server.close(resolve)));
  const tools = await send({jsonrpc: '2.0', id: 4, method: 'tools/list'});
  assert.ok(tools.result.tools.some(x => x.name === 'list_connections'));
  assert.ok(notes.includes('notifications/tools/list_changed'));
  const again = await send({jsonrpc: '2.0', id: 5, method: 'tools/call', params: {name: 'wrm_connect', arguments: {}}});
  assert.match(again.result.content[0].text, /Already connected/);

  // a refused token is the one case that needs a new sign-in
  tokenOK = false;
  const refused = await send({jsonrpc: '2.0', id: 6, method: 'tools/list'});
  assert.match(refused.error.message, /call the wrm_connect tool/);
  const signin = await send({jsonrpc: '2.0', id: 7, method: 'tools/call', params: {name: 'wrm_connect', arguments: {}}});
  assert.doesNotMatch(signin.result.content[0].text, /Already connected/);
});

// An approval request from WRM becomes one desktop notification on this computer (Linux:
// notify-send), also when no WRM tab is open; WRM_MCP_NOTIFY=0 turns it off.
test('Claude bridge shows a desktop notification for an approval request', {skip: process.platform !== 'linux'}, async t => {
  const fs = require('node:fs');
  const os = require('node:os');
  const dir = fs.mkdtempSync(path.join(os.tmpdir(), 'wrm-mcpb-notify-'));
  t.after(() => fs.rmSync(dir, {recursive: true, force: true}));
  const out = path.join(dir, 'notified.txt');
  fs.writeFileSync(path.join(dir, 'notify-send'), '#!/bin/sh\nprintf "%s|" "$@" >> "' + out + '"\necho >> "' + out + '"\n', {mode: 0o755});
  const note = {jsonrpc: '2.0', method: 'notifications/message', params: {level: 'warning', logger: 'wrm.approval',
    data: {type: 'approval_required', approval_id: 'ab12', server: 'demo-api', message: 'systemctl restart nginx'}}};
  const server = http.createServer(async (req, res) => {
    const chunks = [];
    for await (const chunk of req) chunks.push(chunk);
    const body = chunks.length ? JSON.parse(Buffer.concat(chunks).toString('utf8')) : {};
    if (body.method === 'initialize') {
      res.writeHead(200, {'Content-Type': 'application/json', 'Mcp-Session-Id': 's1'});
      res.end(JSON.stringify({jsonrpc: '2.0', id: body.id, result: {protocolVersion: '2025-06-18', capabilities: {tools: {}}, serverInfo: {name: 'wrm', version: '1'}}}));
      return;
    }
    if (body.method === 'notifications/initialized') { res.writeHead(202); res.end(); return; }
    res.writeHead(200, {'Content-Type': 'text/event-stream'});
    const ev = m => res.write('event: message\ndata: ' + JSON.stringify(m) + '\n\n');
    ev(note);
    ev(note); // the same request twice: one notification
    ev({jsonrpc: '2.0', id: body.id, result: {content: [{type: 'text', text: 'done'}]}});
    res.end();
  });
  await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
  t.after(() => new Promise(resolve => server.close(resolve)));
  const child = spawn(process.execPath, [bridge], {
    env: {...process.env, PATH: dir + ':' + process.env.PATH, WRM_MCP_URL: `http://127.0.0.1:${server.address().port}/mcp`, WRM_MCP_TOKEN: 'wrm_pat_test', WRM_MCP_CONFIG: path.join(dir, 'none.json')},
    stdio: ['pipe', 'pipe', 'pipe'],
  });
  t.after(() => { if (!child.killed) child.kill(); });
  child.stderr.resume();
  const got = [];
  const waitFor = id => new Promise((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error('timed out')), 5000);
    readline.createInterface({input: child.stdout}).on('line', l => { const m = JSON.parse(l); got.push(m); if (m.id === id) { clearTimeout(timer); resolve(m); } });
  });
  const done = waitFor(2);
  child.stdin.write(JSON.stringify({jsonrpc: '2.0', id: 1, method: 'initialize', params: {protocolVersion: '2025-06-18', capabilities: {}, clientInfo: {name: 'n'}}}) + '\n');
  child.stdin.write(JSON.stringify({jsonrpc: '2.0', id: 2, method: 'tools/call', params: {name: 'run_command', arguments: {}}}) + '\n');
  const res = await done;
  assert.equal(res.result.content[0].text, 'done');
  assert.ok(got.some(m => m.method === 'notifications/message'), 'the log message still reaches the app');
  for (let i = 0; i < 40 && !fs.existsSync(out); i++) await new Promise(r => setTimeout(r, 50));
  await new Promise(r => setTimeout(r, 200));
  const lines = fs.readFileSync(out, 'utf8').trim().split('\n');
  assert.equal(lines.length, 1, 'one notification per approval');
  assert.match(lines[0], /WRM PRO: approval needed\|demo-api: systemctl restart nginx/);
});
