#!/usr/bin/env node
'use strict';

// Drives a running WRM in Chromium and checks that the ▤ Workboard menu manages panels and
// addons in one place: the separate ▦ Panels button is hidden, the menu lists every panel
// with its place, a panel hidden there can be shown again from "Add to the Workboard",
// locking the layout marks the Workboard button, and the texts are translated.
//
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node tools/test-workboard-panels.js
//
// Skips when WRM_URL is not set. Restores the panel layout it changes.
const assert = require('node:assert/strict');
const {test} = require('node:test');

const url = process.env.WRM_URL;

test('Workboard menu manages panels and addons together', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const browser = await chromium.launch();
  try {
    const page = await (await browser.newContext({viewport: {width: 1600, height: 900}})).newPage();
    const errors = [];
    page.on('pageerror', e => errors.push(e.message));
    // The first page load runs before the sign-in, so its 401 is expected.
    page.on('console', m => { if (m.type() === 'error' && !m.text().includes('401')) errors.push(m.text()); });
    await page.goto(url);
    await page.evaluate(async ([u, p]) => {
      await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
    }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '']);
    await page.goto(url);
    await page.waitForSelector('#top-bar', {state: 'visible'});
    await page.waitForFunction(() => window.wrmWorkboard && wrmWorkboard.ready() && wrmLayout.ready());
    await page.evaluate(() => wrmLayout.show('snip'));

    // One button: ▦ Panels is hidden while the Workboard is available.
    assert.equal(await page.evaluate(() => getComputedStyle(document.getElementById('btn-panels')).display), 'none');

    const menuItems = () => page.evaluate(() => [...document.querySelectorAll('#dk-menu .dk-mi, #dk-menu .dk-mh')].map(e => e.textContent.trim()));
    const openMenu = () => page.evaluate(() => wrmWorkboard.menu(document.getElementById('btn-more')));
    await openMenu();
    let items = await menuItems();
    for (const want of ['Add to the Workboard…', 'Connections', 'Sessions', 'Snippets', 'Clipboard', 'Snap to the grid', 'Lock the layout', 'Reset the layout', 'Save this layout…']) {
      assert.ok(items.some(x => x.includes(want)), `menu lacks "${want}": ${items.join(' | ')}`);
    }
    assert.ok(!items.some(x => x === 'Panels…'), 'the menu should not send to a second Panels menu');

    // Hide Snippets from the menu; the menu stays open and shows it unticked.
    await page.click('#dk-menu .dk-mi:has-text("Snippets")');
    assert.equal(await page.evaluate(() => wrmLayout.where('snip')), 'hidden');
    assert.ok(await page.isVisible('#dk-menu .dk-mi:has-text("Snippets") .ico:text("☐")'), 'Snippets should be unticked');
    await page.keyboard.press('Escape');

    // The panel list of a panel's ⋯ menu ("Panels…") opens the same Workboard menu.
    await page.evaluate(() => wrmLayout.openPanels(document.getElementById('btn-more')));
    items = await menuItems();
    assert.ok(items.some(x => x.includes('Add to the Workboard')), 'Panels… should open the Workboard menu');
    await page.keyboard.press('Escape');

    // "Add to the Workboard" lists the hidden panel with Show, and Show brings it back.
    await page.evaluate(() => wrmWorkboard.catalog());
    await page.waitForSelector('#modal-wb-catalog.open [data-p="snip"]:not([disabled])');
    assert.ok(await page.isDisabled('#modal-wb-catalog [data-p="conn"]'), 'a shown panel cannot be added twice');
    await page.click('#modal-wb-catalog [data-p="snip"]');
    assert.notEqual(await page.evaluate(() => wrmLayout.where('snip')), 'hidden');
    assert.equal(await page.isVisible('#modal-wb-catalog'), false);

    // Locking the layout marks the Workboard button with 🔒.
    await page.evaluate(() => wrmLayout.setLocked(true));
    assert.equal(await page.evaluate(() => document.querySelector('#btn-workboard span').textContent), '🔒');
    await page.evaluate(() => wrmLayout.setLocked(false));
    assert.equal(await page.evaluate(() => document.querySelector('#btn-workboard span').textContent), '▤');

    // Croatian texts.
    await page.evaluate(() => { changeLang('hr'); });
    await openMenu();
    items = await menuItems();
    assert.ok(items.some(x => x.includes('Dodaj na Workboard')) && items.some(x => x.includes('Snippeti') || x.includes('Veze')), `Croatian menu: ${items.join(' | ')}`);
    await page.keyboard.press('Escape');
    await page.evaluate(() => { changeLang('en'); });

    assert.deepEqual(errors, [], 'page errors');
  } finally {
    await browser.close();
  }
});
