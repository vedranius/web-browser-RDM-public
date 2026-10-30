#!/usr/bin/env node
'use strict';

// Check that transient tunnel API failures are visible and the start control recovers.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const {test} = require('node:test');

const html = fs.readFileSync(path.join(__dirname, '..', 'remote-manager', 'static', 'index.html'), 'utf8');
const begin = html.indexOf('async function tunnelStart(');
const end = html.indexOf('async function startConnTunnels(', begin);
if (begin < 0 || end < 0) throw new Error('tunnel action functions not found');
const source = html.slice(begin, end);

function fixture(fetch) {
  const messages = [];
  let refreshes = 0;
  const context = vm.createContext({
    fetch,
    showToast: (...args) => messages.push(args),
    apiError: async () => 'server refused the request',
    loadTunnels: async () => { refreshes++; },
    t: key => key,
    listenAddress: () => '127.0.0.1:5000',
  });
  vm.runInContext(source, context);
  return {context, messages, refreshes: () => refreshes};
}

test('tunnel start reports network failure and restores its button', async () => {
  const f = fixture(async () => { throw new Error('offline'); });
  const button = {disabled: false, textContent: 'Start'};
  await f.context.tunnelStart('local:test', button);
  assert.equal(button.disabled, false);
  assert.equal(f.messages[0][0], 'offline');
  assert.equal(f.messages[0][1], 'error');
  assert.equal(f.refreshes(), 1);
});

test('tunnel stop reports HTTP failure and refreshes state', async () => {
  const f = fixture(async () => ({ok: false}));
  await f.context.tunnelStop('local:test');
  assert.equal(f.messages[0][0], 'server refused the request');
  assert.equal(f.messages[0][1], 'error');
  assert.equal(f.refreshes(), 1);
});
