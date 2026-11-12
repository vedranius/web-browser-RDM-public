#!/usr/bin/env node
'use strict';
// Unit tests for the terminal clipboard history (pins, limit, old format, paste safety).
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const {test} = require('node:test');
const html = fs.readFileSync(path.join(__dirname, '..', 'remote-manager', 'static', 'index.html'), 'utf8');
const slice = (from, to) => html.slice(html.indexOf(from), html.indexOf(to));

function makeCtx(extra) {
  const store = {};
  const ctx = vm.createContext(Object.assign({
    CLIP_KEY: 'wrm_clipboard_v9', clipData: {}, currentSessionId: null, currentUser: {id: 7}, winObjects: [],
    localStorage: {getItem: k => (k in store ? store[k] : null), setItem: (k, v) => { store[k] = String(v); }, removeItem: k => { delete store[k]; }},
    document: {getElementById: () => null}, t: k => k, esc: s => String(s), jsq: s => JSON.stringify(s),
    winConn: () => null, copied: [], copyText(s) { ctx.copied.push(s); }, toasts: [], showToast(m) { ctx.toasts.push(m); },
    snippetTargetWin: () => null, isDangerousCommand: l => /rm -rf/.test(l), uiConfirm: async () => { ctx.confirms++; return ctx.confirmAnswer; }, confirms: 0, confirmAnswer: true,
    sent: [], termSend(wo, d) { ctx.sent.push(d); return true; }, opened: null, openSnippetModal(id, p) { ctx.opened = p; },
  }, extra || {}));
  vm.runInContext(slice('function sessionClipKey(', '// ══ WINDOWS ═'), ctx);
  ctx.store = store;
  return ctx;
}
const texts = ctx => JSON.parse(vm.runInContext('JSON.stringify((clipData.global || []).map(e => e.text))', ctx));

test('history keeps the newest entries, moves duplicates up and stays within the limit', () => {
  const ctx = makeCtx();
  for (let i = 0; i < 25; i++) ctx.addClip('cmd ' + i, '');
  ctx.addClip('cmd 10', '');
  const l = texts(ctx);
  assert.equal(l.length, 20);
  assert.equal(l[0], 'cmd 10');
  assert.equal(l.filter(x => x === 'cmd 10').length, 1);
  assert.ok(ctx.store['wrm_clipboard_v9:u7']);
});

test('pinned entries stay on top, survive the limit and Clear', () => {
  const ctx = makeCtx();
  ctx.addClip('keep me', 'web-1');
  ctx.clipTogglePin('global', 0);
  for (let i = 0; i < 30; i++) ctx.addClip('x' + i, '');
  let l = texts(ctx);
  assert.equal(l[0], 'keep me');
  assert.equal(l.length, 21);
  ctx.clearClipboard();
  assert.deepEqual(texts(ctx), ['keep me']);
  ctx.clearClipboard();
  assert.deepEqual(texts(ctx), []);
});

test('old plain-string history is converted and an expired history is dropped', () => {
  const ctx = makeCtx();
  ctx.store['wrm_clipboard_v9:u7'] = JSON.stringify({global: ['old one', 5, null]});
  ctx.loadClipboard();
  assert.deepEqual(texts(ctx), ['old one']);
  ctx.store['wrm_clipboard_v9:u7:at'] = String(Date.now() - 2 * 3600 * 1000);
  ctx.loadClipboard();
  assert.equal(vm.runInContext('Object.keys(clipData).length', ctx), 0);
});

test('paste asks before running several lines without bracketed paste or a dangerous command', async () => {
  const sock = {readyState: 1};
  const wo = {el: {classList: {contains: () => true}}, socket: sock, term: {modes: {bracketedPasteMode: false}, focus() {}}};
  const ctx = makeCtx({winObjects: [wo]});
  ctx.addClip('echo a\necho b', '');
  ctx.confirmAnswer = false;
  await ctx.pasteClipItem('global', 0);
  assert.equal(ctx.confirms, 1);
  assert.equal(ctx.sent.length, 0);
  wo.term.modes.bracketedPasteMode = true;
  await ctx.pasteClipItem('global', 0);
  assert.equal(ctx.confirms, 1);
  assert.deepEqual(ctx.sent, ['\x1b[200~echo a\recho b\x1b[201~']);
  ctx.addClip('rm -rf /tmp/x', '');
  await ctx.pasteClipItem('global', 0);
  assert.equal(ctx.confirms, 2);
  assert.equal(ctx.sent.length, 1);
});

test('paste without a terminal copies instead; save as snippet prefills the dialog', async () => {
  const ctx = makeCtx();
  ctx.addClip('systemctl status nginx', '');
  await ctx.pasteClipItem('global', 0);
  assert.deepEqual(ctx.copied, ['systemctl status nginx']);
  ctx.clipToSnippet('global', 0);
  assert.equal(ctx.opened.command, 'systemctl status nginx');
  assert.equal(ctx.opened.group, 'clip_snip_group');
});
