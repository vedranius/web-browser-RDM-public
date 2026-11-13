#!/usr/bin/env node
'use strict';

// Drives a running WRM in Chromium: drags a connection in the tree above another one and
// checks that the user's order is kept after a reload, that "Move connection down" in the
// context menu works, and that "Sort connections A–Z" returns to the order by name.
//
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node tools/test-conn-order.js
//
// Skips when WRM_URL is not set. Creates three connections demo-order-a/b/c and deletes them.
const assert = require('node:assert/strict');
const {test} = require('node:test');

const url = process.env.WRM_URL;

test('connections can be dragged into the user order', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const browser = await chromium.launch();
  try {
    for (const [width, height] of [[1280, 800], [390, 844]]) {
      const page = await (await browser.newContext({viewport: {width, height}})).newPage();
      const errors = [];
      page.on('pageerror', e => errors.push(e.message));
      page.on('console', m => { if (m.type() === 'error') errors.push(m.text()); });
      await page.goto(url);
      const api = (path, body) => page.evaluate(async ([p, b]) => {
        const r = await fetch(p, {method: b ? 'POST' : 'GET', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: b ? JSON.stringify(b) : undefined});
        return r.ok ? r.json() : null;
      }, [path, body]);
      await api('/api/auth/login', {username: process.env.WRM_USER || 'admin', password: process.env.WRM_PASS || ''});
      const ids = [];
      for (const n of ['demo-order-a', 'demo-order-b', 'demo-order-c']) {
        const c = await api('/api/connections', {name: n, protocol: 'SSH', host: n + '.example.com', username: 'root', auth_method: 'PASSWORD'});
        assert.ok(c && c.id, 'connection created');
        ids.push(c.id);
      }
      await page.goto(url);
      errors.length = 0; // the first load before the login asks /api/auth/me (401)
      await page.evaluate(() => { if (typeof changeLang === 'function') changeLang('hr'); });
      if (width < 800) await page.evaluate(() => { const b = document.querySelector('#sidebar-toggle, .sidebar-toggle, [onclick*="toggleSidebar"]'); if (b) b.click(); });
      const order = () => page.evaluate(() => [...document.querySelectorAll('#tree .conn-item .conn-name, .conn-item .conn-name')].map(e => e.textContent).filter(n => n.startsWith('demo-order-')));
      await page.waitForFunction(() => document.querySelectorAll('.conn-item').length >= 3);
      assert.deepEqual(await order(), ['demo-order-a', 'demo-order-b', 'demo-order-c']);
      if (width >= 800) {
        // drag c onto the upper half of a → c, a, b
        const src = page.locator('.conn-item', {hasText: 'demo-order-c'}).first();
        const dst = page.locator('.conn-item', {hasText: 'demo-order-a'}).first();
        await src.dragTo(dst, {targetPosition: {x: 20, y: 3}});
        await page.waitForFunction(() => [...document.querySelectorAll('.conn-item .conn-name')].map(e => e.textContent).filter(n => n.startsWith('demo-order-'))[0] === 'demo-order-c');
        await page.reload();
        await page.waitForFunction(() => document.querySelectorAll('.conn-item').length >= 3);
        assert.deepEqual(await order(), ['demo-order-c', 'demo-order-a', 'demo-order-b'], 'order kept after reload');
        await page.screenshot({path: process.env.SHOT_DIR ? process.env.SHOT_DIR + '/conn-order.png' : '/dev/null'});
        // context menu: move c down → a, c, b (Croatian label)
        await page.locator('.conn-item', {hasText: 'demo-order-c'}).first().click({button: 'right'});
        await page.locator('.ctx-item', {hasText: 'Pomakni vezu dolje'}).click();
        await page.waitForFunction(() => [...document.querySelectorAll('.conn-item .conn-name')].map(e => e.textContent).filter(n => n.startsWith('demo-order-'))[1] === 'demo-order-c');
        // back to A–Z for the top level
        await page.evaluate(() => resetConnOrder(null));
        await page.waitForFunction(() => [...document.querySelectorAll('.conn-item .conn-name')].map(e => e.textContent).filter(n => n.startsWith('demo-order-')).join() === 'demo-order-a,demo-order-b,demo-order-c');
      }
      await api('/api/connections/bulk', {action: 'delete', ids});
      assert.deepEqual(errors, [], 'no console errors at ' + width);
      await page.context().close();
    }
  } finally {
    await browser.close();
  }
});
