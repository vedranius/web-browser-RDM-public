#!/usr/bin/env node
'use strict';
// Browser test of Settings → Data → "Import from another tool": the picker detects the format,
// shows a preview without writing, imports on "Import", and a second import is skipped as
// duplicates. Croatian texts and phone width are checked too.
//
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node tools/test-import-picker.js
//
// Skips when WRM_URL is not set. Use a scratch database: the test creates connections.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const {test} = require('node:test');

const url = process.env.WRM_URL;

const RUN = Date.now().toString(36); // unique names: the test can run again on the same database
const RDG = `<?xml version="1.0" encoding="utf-8"?>
<RDCMan programVersion="2.90" schemaVersion="3"><file><properties><name>demo</name></properties>
<group><properties><name>Demo DC ${RUN}</name></properties>
<server><properties><name>demo-web.example.com</name><displayName>demo-web</displayName></properties></server>
<server><properties><name>demo-api.example.com</name><displayName>demo-api</displayName></properties></server>
</group></file></RDCMan>`;

async function login(page) {
  await page.goto(url);
  await page.evaluate(async ([u, p]) => {
    await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
  }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '']);
  await page.goto(url);
  await page.waitForFunction(() => typeof openImportPicker === 'function' && typeof currentUser !== 'undefined' && !!currentUser);
}

test('import picker previews, imports and skips duplicates', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const file = path.join(fs.mkdtempSync(path.join(os.tmpdir(), 'wrm-imp-')), 'demo.rdg');
  fs.writeFileSync(file, RDG);
  const browser = await chromium.launch();
  try {
    const page = await browser.newPage({viewport: {width: 1280, height: 900}});
    await login(page);
    const errors = []; // after the sign-in: the first page load answers 401 before it
    page.on('pageerror', e => errors.push(e.message));
    page.on('console', m => { if (m.type() === 'error') errors.push(m.text()); });
    const count = () => page.evaluate(async () => (await (await fetch('/api/connections')).json()).length);
    const before = await count();

    await page.evaluate(() => openImportPicker());
    await page.setInputFiles('#imp-file', file);
    await page.waitForSelector('#imp-out table');
    const out = await page.textContent('#imp-out');
    assert.match(out, /Remote Desktop Connection Manager/);
    assert.match(out, /demo-web/);
    assert.match(out, /Demo DC/);
    assert.equal(await count(), before, 'the preview must not write');
    if (process.env.SHOT) await page.screenshot({path: process.env.SHOT + '-en.png'});

    await page.click('#imp-go');
    await page.waitForSelector('#imp-modal', {state: 'detached'});
    await page.waitForFunction(n => fetch('/api/connections').then(r => r.json()).then(l => l.length === n), before + 2);
    await page.keyboard.press('Escape');

    // Same file again: the preview lists both as skipped and nothing is importable.
    await page.evaluate(() => openImportPicker());
    await page.setInputFiles('#imp-file', file);
    await page.waitForFunction(() => /demo-web/.test(document.querySelector('#imp-out')?.textContent || ''));
    assert.equal(await page.isDisabled('#imp-go'), true);
    await page.click('#imp-cancel');

    // Croatian and phone width.
    await page.evaluate(() => { changeLang('hr'); });
    await page.setViewportSize({width: 390, height: 800});
    await page.evaluate(() => openImportPicker());
    assert.match(await page.textContent('#imp-modal h3'), /Uvoz veza/);
    const box = await page.locator('#imp-modal .modal').boundingBox();
    assert.ok(box.width <= 390, 'dialog fits the phone width');
    if (process.env.SHOT) { await page.waitForTimeout(500); await page.screenshot({path: process.env.SHOT + '-hr-phone.png'}); }
    assert.deepEqual(errors, []);
  } finally {
    await browser.close();
  }
});
