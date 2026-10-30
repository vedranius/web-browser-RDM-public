#!/usr/bin/env node
'use strict';

// Exercise the host-side capability check and iframe navigation revocation with
// the same bundled helper and MessageChannel used by addon panels.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const {MessageChannel} = require('node:worker_threads');
const {test} = require('node:test');

const root = path.join(__dirname, '..');
const host = fs.readFileSync(path.join(root, 'remote-manager', 'static', 'addons.js'), 'utf8');
const helper = fs.readFileSync(path.join(root, 'remote-manager', 'addon_templates', '_shared', 'wrm-bridge.js'), 'utf8');
const start = host.indexOf('function wrmAddonAcceptBridgeMessage(');
const end = host.indexOf('\nObject.assign(LANGS.en', start);
if (start < 0 || end < 0) throw new Error('host bridge functions not found');

function hostFixture() {
  const handled = [], diagnostics = [], banners = [];
  const context = vm.createContext({
    addDiag: (_st, d) => diagnostics.push(d),
    banner: (_st, kind) => banners.push(kind),
    MessageChannel,
  });
  vm.runInContext(host.slice(start, end), context, {filename: 'addons.js bridge'});
  return {context, handled, diagnostics, banners, onMessage: (_st, d) => handled.push(d)};
}

function panelFixture() {
  const listeners = {}, messages = [];
  const window = {parent: {postMessage(data, _origin, ports = []) { messages.push({data, ports}); }}, addEventListener(type, fn) { listeners[type] = fn; }};
  const context = vm.createContext({window, location: {hash: '#wrm-bridge=host-capability'}, document: {documentElement: {style: {setProperty() {}}}}, URLSearchParams, MessageChannel,
    setTimeout, clearTimeout, setInterval: () => 1});
  vm.runInContext(helper, context, {filename: 'wrm-bridge.js'});
  return {window, listeners, messages};
}

test('host rejects a forged handshake, then accepts the one-time token and ignores window traffic', async () => {
  const f = hostFixture();
  const st = {bridgeToken: 'host-capability', secureBridge: false, revoked: false, dropped: 0};
  const forged = new MessageChannel();
  assert.equal(f.context.wrmAddonAcceptBridgeMessage(st, {v: 2, type: 'wrm:ready', bridge_token: 'wrong'}, [forged.port1], f.onMessage, (_s, d) => f.diagnostics.push(d)), false);
  assert.equal(st.secureBridge, false);
  assert.equal(st.bridgeToken, 'host-capability');
  forged.port1.close(); forged.port2.close();

  const panel = panelFixture();
  panel.window.wrm.ready();
  const init = panel.messages[0];
  assert.equal(f.context.wrmAddonAcceptBridgeMessage(st, init.data, init.ports, f.onMessage, (_s, d) => f.diagnostics.push(d)), true);
  assert.equal(st.secureBridge, true);
  assert.equal(st.bridgeToken, '');
  await new Promise(resolve => setImmediate(resolve));
  assert.equal(f.handled[0].type, 'wrm:ready');

  panel.window.wrm.notify('through private port');
  await new Promise(resolve => setImmediate(resolve));
  assert.equal(f.handled.at(-1).type, 'wrm:notify');
  assert.equal(f.handled.at(-1).text, 'through private port');
  assert.equal(f.context.wrmAddonAcceptBridgeMessage(st, {v: 2, type: 'wrm:http', id: 'forged'}, [], f.onMessage, () => {}), false);
  assert.equal(f.handled.length, 2, 'legacy window messages are ignored after port setup');
  st.bridgePort.close();
});

test('a second load revokes and removes the navigated iframe and closes its port', () => {
  const f = hostFixture();
  let removed = 0, closed = 0;
  const frame = {remove() { removed++; }};
  const st = {frame, bridgePort: {close() { closed++; }}, revoked: false, ready: true, wo: {id: 10}, diag: []};
  let loads = 0;
  loads = f.context.wrmAddonFrameLoaded(st, frame, loads);
  assert.equal(st.revoked, false, 'the initial load keeps the bridge available');
  loads = f.context.wrmAddonFrameLoaded(st, frame, loads);
  assert.equal(st.revoked, true);
  assert.equal(st.ready, false);
  assert.equal(st.frame, null);
  assert.equal(st.bridgePort, null);
  assert.equal(removed, 1);
  assert.equal(closed, 1);
  assert.equal(f.diagnostics[0].kind, 'E_BRIDGE_NAVIGATION');
  assert.deepEqual(f.banners, ['start']);
  f.context.wrmAddonFrameLoaded(st, frame, loads);
  assert.equal(removed, 1, 'a stale load event must not revoke a replacement frame');
});
