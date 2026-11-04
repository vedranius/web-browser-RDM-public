#!/usr/bin/env node
'use strict';

// Exercise the bundled addon helper's secure MessageChannel handshake and its
// compatibility fallback when an older WRM host does not acknowledge the port.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const {MessageChannel} = require('node:worker_threads');
const {test} = require('node:test');

const source = fs.readFileSync(path.join(__dirname, '..', 'remote-manager', 'addon_templates', '_shared', 'wrm-bridge.js'), 'utf8');

function fixture() {
  const listeners = {};
  const messages = [];
  const timers = [];
  const parent = {
    postMessage(data, _origin, ports = []) { messages.push({data, ports}); },
  };
  const window = {
    parent,
    addEventListener(type, listener) { listeners[type] = listener; },
  };
  const context = vm.createContext({
    window,
    location: {hash: '#wrm-bridge=one-time-capability'},
    document: {documentElement: {style: {setProperty() {}}}},
    URLSearchParams,
    MessageChannel,
    setTimeout(fn) { timers.push(fn); return timers.length; },
    clearTimeout() {},
    setInterval() { return 1; },
  });
  vm.runInContext(source, context, {filename: 'wrm-bridge.js'});
  return {listeners, messages, timers, wrm: window.wrm, parent};
}

test('uses a one-time capability and the private port after host acknowledgement', async () => {
  const f = fixture();
  f.wrm.ready();
  assert.equal(f.messages.length, 1);
  assert.equal(f.messages[0].data.bridge_token, 'one-time-capability');
  assert.equal(f.messages[0].ports.length, 1);

  const hostPort = f.messages[0].ports[0];
  hostPort.postMessage({v: 2, type: 'wrm:bridge:ready'});
  await new Promise(resolve => setImmediate(resolve));
  f.wrm.notify('private message');
  const received = await new Promise(resolve => hostPort.once('message', resolve));
  assert.deepEqual(received, {v: 2, type: 'wrm:notify', text: 'private message', level: 'info'});
  assert.equal(f.messages.length, 1, 'postMessage must not carry traffic after acknowledgement');
  hostPort.close();
});

test('falls back to window messages when an older host does not acknowledge', async () => {
  const f = fixture();
  f.wrm.ready();
  f.timers[0]();
  assert.equal(f.messages.length, 2);
  assert.equal(f.messages[1].data.type, 'wrm:ready');
  assert.equal('bridge_token' in f.messages[1].data, false);

  const pending = f.wrm.storeSet('layout', {left: 12});
  const request = f.messages[2].data;
  assert.equal(request.type, 'wrm:store.set');
  f.listeners.message({source: f.parent, data: {v: 2, id: request.id, ok: true}});
  await assert.doesNotReject(pending);
});
