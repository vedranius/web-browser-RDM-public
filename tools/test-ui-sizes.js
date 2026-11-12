#!/usr/bin/env node
'use strict';

// Drives a running WRM in Chromium at phone, tablet, laptop and desktop sizes and at
// browser zoom levels, and checks that the top bar never overlaps or overflows, that
// buttons keep their labels where there is room, and that the page has no console errors.
//
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node tools/test-ui-sizes.js
//
// Skips when WRM_URL is not set (CI runs it only where a test server is started).
const assert = require('node:assert/strict');
const {test} = require('node:test');

const url = process.env.WRM_URL;
const sizes = [
  // name, width, height, mobile, deviceScaleFactor (browser zoom = CSS width shrinks, DSF grows)
  ['phone-small', 320, 640, true, 2],
  ['phone', 390, 844, true, 3],
  ['phone-landscape', 844, 390, true, 3],
  ['tablet-portrait', 768, 1024, true, 2],
  ['tablet', 820, 1180, true, 2],
  ['tablet-landscape', 1180, 820, true, 2],
  ['laptop', 1280, 800, false, 1],
  ['desktop-1100', 1100, 800, false, 1],
  ['desktop', 1920, 1080, false, 1],
  ['zoom-200', 960, 540, false, 2],
  ['zoom-150', 1280, 720, false, 1.5],
  ['zoom-67', 2880, 1620, false, 0.67],
];

test('top bar fits at every screen size and zoom', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const browser = await chromium.launch();
  try {
    for (const [name, width, height, mobile, dsf] of sizes) {
      const ctx = await browser.newContext({viewport: {width, height}, isMobile: mobile, hasTouch: mobile, deviceScaleFactor: dsf});
      const page = await ctx.newPage();
      const errors = [];
      page.on('pageerror', e => errors.push(e.message));
      await page.goto(url);
      await page.evaluate(async ([u, p]) => {
        await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
      }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '']);
      await page.goto(url);
      await page.waitForSelector('#top-bar', {state: 'visible'});
      await page.waitForTimeout(800);
      const r = await page.evaluate(() => {
        const tb = document.getElementById('top-bar');
        const vis = el => { const s = getComputedStyle(el); const b = el.getBoundingClientRect(); return s.display !== 'none' && s.visibility !== 'hidden' && b.width > 0; };
        const items = [...tb.querySelectorAll('#top-bar > *, #auth-info > *, #tb-canvas-grp > *')].filter(e => e.id !== 'auth-info' && e.id !== 'tb-canvas-grp' && vis(e));
        const boxes = items.map(e => ({id: e.id || e.className, r: e.getBoundingClientRect()}));
        const overlaps = [];
        for (let i = 0; i < boxes.length; i++) for (let j = i + 1; j < boxes.length; j++) {
          const a = boxes[i].r, b = boxes[j].r;
          if (a.left < b.right - 1 && b.left < a.right - 1 && a.top < b.bottom - 1 && b.top < a.bottom - 1) overlaps.push(boxes[i].id + ' × ' + boxes[j].id);
        }
        const outside = boxes.filter(b => b.r.right > window.innerWidth + 1 || b.r.left < -1).map(b => b.id);
        const labelled = [...tb.querySelectorAll('.btn-top .lbl')].filter(vis).length;
        return {overlaps, outside, labelled, hscroll: document.documentElement.scrollWidth > window.innerWidth + 1};
      });
      assert.deepEqual(r.overlaps, [], `${name}: top bar items overlap`);
      assert.deepEqual(r.outside, [], `${name}: top bar items outside the screen`);
      assert.equal(r.hscroll, false, `${name}: page scrolls sideways`);
      if (width >= 1280) assert.ok(r.labelled >= 3, `${name}: top bar buttons lost their labels (${r.labelled} shown)`);
      assert.deepEqual(errors, [], `${name}: page errors`);
      await ctx.close();
    }
  } finally {
    await browser.close();
  }
});

test('dialogs minimise to a chip and toasts are grouped', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const browser = await chromium.launch();
  try {
    const page = await browser.newPage({viewport: {width: 1400, height: 900}});
    const errors = [];
    page.on('pageerror', e => errors.push(e.message));
    await page.goto(url);
    await page.evaluate(async ([u, p]) => {
      await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
    }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '']);
    await page.goto(url);
    await page.waitForSelector('#top-bar', {state: 'visible'});
    await page.evaluate(() => openSettings());
    await page.click('.wm-min-btn');
    const count = () => page.evaluate(() => [document.querySelectorAll('.modal-overlay.open').length, document.querySelectorAll('#wm-tray .wm-chip').length]);
    assert.deepEqual(await count(), [0, 1], 'minimised dialog leaves a chip');
    await page.click('.wm-chip-open');
    await page.waitForTimeout(100);
    assert.deepEqual(await count(), [1, 0], 'the chip opens the dialog again');
    await page.click('.wm-min-btn');
    await page.evaluate(() => openSettings());
    await page.waitForTimeout(100);
    assert.deepEqual(await count(), [1, 0], 'opening the dialog from the menu removes its chip');
    const toasts = await page.evaluate(() => {
      for (let i = 0; i < 3; i++) showToast('same', 'error', 20000);
      for (let i = 1; i <= 6; i++) showToast('server ' + i, 'error', 20000);
      const box = document.getElementById('toast-container');
      return {n: box.querySelectorAll('.toast').length, group: !!box.querySelector('.toast-group'), same: box.textContent.includes('×3')};
    });
    assert.deepEqual(toasts, {n: 5, group: true, same: true});
    assert.deepEqual(errors, []);
  } finally {
    await browser.close();
  }
});
