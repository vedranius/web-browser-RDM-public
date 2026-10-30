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

const fileActionsBegin = html.indexOf('async function createRemoteFolder(');
const fileActionsEnd = html.indexOf('\nfunction fmtBytes(', fileActionsBegin);
if (fileActionsBegin < 0 || fileActionsEnd < 0) throw new Error('remote file action functions not found');
const fileActions = html.slice(fileActionsBegin, fileActionsEnd);

function fileFixture(fetch) {
  const messages = [];
  let reloads = 0;
  const context = vm.createContext({
    fetch,
    winObjects: [{id: 7, connId: 3, path: '/srv', _searchResults: null}],
    guardWin: () => true,
    uiPrompt: async () => 'reports',
    uiConfirm: async () => true,
    showToast: (...args) => messages.push(args),
    apiError: async () => 'server refused the request',
    joinPath: (base, name) => base + '/' + name,
    connQS: () => 'id=3',
    loadFiles: async () => { reloads++; },
    endFmSearch: () => {},
    t: key => key,
  });
  vm.runInContext(fileActions, context);
  return {context, messages, reloads: () => reloads};
}

test('remote folder creation reports network failure without claiming success', async () => {
  const f = fileFixture(async () => { throw new Error('offline'); });
  await f.context.createRemoteFolder(7);
  assert.match(f.messages[0][0], /network_error.*offline/);
  assert.equal(f.messages[0][1], 'error');
  assert.equal(f.messages.some(m => m[1] === 'success'), false);
  assert.equal(f.reloads(), 0);
});

test('remote item deletion reports network failure without claiming success', async () => {
  const f = fileFixture(async () => { throw new Error('offline'); });
  await f.context.deleteRemoteItem(7, '/srv/report.txt', false, 'report.txt');
  assert.match(f.messages[0][0], /network_error.*offline/);
  assert.equal(f.messages[0][1], 'error');
  assert.equal(f.messages.some(m => m[1] === 'success'), false);
  assert.equal(f.reloads(), 0);
});

const appearance = fs.readFileSync(path.join(__dirname, '..', 'remote-manager', 'static', 'appearance.js'), 'utf8');
const syncBegin = appearance.indexOf('  let saveTimer = null;');
const syncEnd = appearance.indexOf('\n\n  async function load()', syncBegin);
if (syncBegin < 0 || syncEnd < 0) throw new Error('appearance sync functions not found');
const syncSource = appearance.slice(syncBegin, syncEnd);

function appearanceFixture(fetch) {
  const messages = [];
  const timers = new Map();
  let nextTimer = 1;
  let onOnline;
  const context = vm.createContext({
    prefs: {theme: 'graphite'},
    apply: () => {},
    fetch,
    navigator: {onLine: false},
    window: {addEventListener: (name, fn) => { if (name === 'online') onOnline = fn; }},
    setTimeout: fn => { const id = nextTimer++; timers.set(id, fn); return id; },
    clearTimeout: id => timers.delete(id),
    showToast: (...args) => messages.push(args),
    t: key => key,
  });
  vm.runInContext(syncSource, context);
  return {
    context,
    messages,
    async flushSave() {
      const [id, fn] = [...timers.entries()].at(-1);
      timers.delete(id);
      await fn();
    },
    online: () => onOnline(),
  };
}

test('appearance sync reports network loss and retries when the browser is online', async () => {
  let requests = 0;
  const f = appearanceFixture(async () => {
    requests++;
    if (requests === 1) throw new Error('offline');
    return {ok: true, status: 200};
  });
  f.context.save({theme: 'fjord'});
  await f.flushSave();
  assert.equal(requests, 1);
  assert.equal(f.messages[0][0], 'ap_sync_failed');
  assert.equal(f.messages[0][1], 'error');

  f.context.navigator.onLine = true;
  await f.online();
  assert.equal(requests, 2);
  assert.equal(f.messages.length, 1);
});

const shareBootBegin = html.indexOf('function shareTokenFromLocation()');
const shareBootEnd = html.indexOf('\n\nwindow.addEventListener(\'load\'', shareBootBegin);
if (shareBootBegin < 0 || shareBootEnd < 0) throw new Error('share URL cleanup helper not found');
const shareBootSource = html.slice(shareBootBegin, shareBootEnd);

function shareBootFixture(pathname, failStorage = false) {
  let currentPath = pathname;
  let currentState = null;
  const values = new Map();
  const replaced = [];
  const context = vm.createContext({
    location: {get pathname() { return currentPath; }},
    sessionStorage: {
      getItem: key => values.get(key) || null,
      setItem: (key, value) => { if (failStorage) throw new Error('storage disabled'); values.set(key, value); },
      removeItem: key => values.delete(key),
    },
    history: {
      get state() { return currentState; },
      replaceState: (state, _title, url) => { replaced.push(url); currentPath = url; currentState = state; },
    },
  });
  vm.runInContext(shareBootSource, context);
  return {
    context,
    values,
    replaced,
    get state() { return currentState; },
    navigate(path, state = null) { currentPath = path; currentState = state; },
  };
}

test('share tokens stay associated with their history entries across rooms and back navigation', () => {
  const first = shareBootFixture('/share/secret-capability/');
  assert.equal(first.context.shareTokenFromLocation(), 'secret-capability');
  assert.equal(first.values.has('wrm.share.entries'), true);
  assert.deepEqual(first.replaced, ['/share']);
  const firstState = first.state;

  first.navigate('/share/another-capability');
  assert.equal(first.context.shareTokenFromLocation(), 'another-capability');
  const secondState = first.state;
  assert.notEqual(firstState.wrmShareEntry, secondState.wrmShareEntry);

  first.navigate('/share', firstState);
  assert.equal(first.context.shareTokenFromLocation(), 'secret-capability');
  first.navigate('/', null);
  assert.equal(first.context.shareTokenFromLocation(), '');
  first.navigate('/share', firstState);
  assert.equal(first.context.shareTokenFromLocation(), 'secret-capability');
});

test('share token URL remains intact if session storage cannot persist it', () => {
  const f = shareBootFixture('/share/secret-capability', true);
  assert.equal(f.context.shareTokenFromLocation(), 'secret-capability');
  assert.deepEqual(f.replaced, []);
});

const alertsBegin = html.indexOf('var wrmSoundCtx = null, wrmAlertSeen = {}, wrmDesktopAlerts = {};');
const alertsEnd = html.indexOf('\nasync function checkStatusNow(', alertsBegin);
if (alertsBegin < 0 || alertsEnd < 0) throw new Error('AI notification functions not found');
const alertsSource = html.slice(alertsBegin, alertsEnd);

test('approval notifications are persistent per approval and click opens its card', () => {
  const opened = [], notifications = [];
  class FakeNotification {
    static permission = 'granted';
    constructor(title, options) { this.title = title; this.options = options; notifications.push(this); }
    close() { this.closed = true; }
  }
  const context = vm.createContext({
    window: {Notification: FakeNotification, focus() {}, wrmFocusApproval: id => opened.push(id)},
    Notification: FakeNotification,
    document: {hidden: true, addEventListener() {}},
    localStorage: {getItem: () => '0'},
    showToast: () => null,
    t: key => key,
    Date,
  });
  vm.runInContext(alertsSource, context);
  context.wrmAiAlert('approval', 'secret command text', 'approval-a');
  context.wrmAiAlert('approval', 'another secret command', 'approval-b');
  assert.equal(notifications.length, 2);
  assert.notEqual(notifications[0].options.tag, notifications[1].options.tag);
  assert.equal(notifications[0].options.requireInteraction, true);
  assert.equal(notifications[0].options.body, 'ai_approval_notify');
  assert.equal(notifications[0].options.body.includes('secret'), false);
  notifications[0].onclick();
  assert.deepEqual(opened, ['a']);
  context.wrmAiAlertDone('approval-a');
  assert.equal(notifications[0].closed, true);
});
