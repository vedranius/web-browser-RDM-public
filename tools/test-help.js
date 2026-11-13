// Browser test of the help panel (static/help.js): ? button, topics, search in Croatian,
// articles, ⌖ Pick, the AI setup checklist, the ? on error toasts and the parsing of the
// help AI's action block. Needs a running WRM:
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node --test tools/test-help.js
// WRM_SHOTS=<dir> saves screenshots.
'use strict';
const test = require('node:test');
const assert = require('node:assert');
const path = require('path');

const url = process.env.WRM_URL;
const shots = process.env.WRM_SHOTS;

async function signedIn(browser, opts, lang) {
  const ctx = await browser.newContext(opts);
  const page = await ctx.newPage();
  const errors = [];
  page.on('pageerror', e => errors.push(e.message));
  page.on('console', m => { if (m.type() === 'error' && !/favicon|status of 40[13]/.test(m.text())) errors.push(m.text()); });
  await page.goto(url);
  await page.evaluate(async ([u, p, l]) => {
    localStorage.setItem('wrm_lang', l);
    await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
  }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '', lang]);
  await page.goto(url);
  await page.waitForSelector('#top-bar', {state: 'visible'});
  await page.waitForTimeout(600);
  return {ctx, page, errors};
}

test('help panel: topics, search, articles, pick and AI setup', {skip: !url && 'WRM_URL not set'}, async () => {
  const {chromium} = require('playwright');
  const browser = await chromium.launch();
  try {
    const {ctx, page, errors} = await signedIn(browser, {viewport: {width: 1440, height: 900}}, 'hr');
    await page.click('#btn-help');
    await page.waitForSelector('#help-panel .hp-item');
    assert.ok(await page.locator('#help-panel .hp-item').count() >= 30, 'topics list every article');
    assert.equal(await page.textContent('#hp-ttl'), 'Pomoć');
    assert.ok((await page.textContent('#hp-body')).includes('Prvi koraci'), 'Croatian article titles');
    if (shots) await page.screenshot({path: path.join(shots, 'help-topics.png')});

    // search without diacritics finds Croatian text
    await page.fill('#hp-q', 'kljuc');
    await page.waitForFunction(() => document.querySelectorAll('#hp-body .hp-item').length > 0 && /ključ/i.test(document.getElementById('hp-body').textContent));
    await page.fill('#hp-q', '');

    // an article
    await page.evaluate(() => wrmHelp.open('ai-setup'));
    await page.waitForSelector('#help-panel .hp-md h3, #help-panel .hp-md h4');
    assert.ok((await page.textContent('#hp-body')).includes('Postavi AI'));
    if (shots) await page.screenshot({path: path.join(shots, 'help-article.png')});

    // pick a part of the screen
    await page.click('#hp-pick');
    await page.waitForSelector('#help-pick');
    const b = await page.locator('#btn-settings').boundingBox();
    await page.mouse.move(b.x + b.width / 2, b.y + b.height / 2);
    await page.mouse.down();
    await page.mouse.up();
    await page.waitForSelector('#help-pick', {state: 'detached'});
    const picked = await page.textContent('#hp-body');
    assert.ok(picked.includes('Odabrano') && picked.includes('Vaš račun i postavke'), 'the pick shows the articles about it: ' + picked.slice(0, 200));

    // Esc cancels a pick
    await page.click('#hp-pick');
    await page.keyboard.press('Escape');
    await page.waitForSelector('#help-pick', {state: 'detached'});

    // the help AI's action block: only known actions survive
    const acts = await page.evaluate(() => wrmHelp.parseActions('x\n```wrm-actions\n[{"type":"open","target":"settings:st-ai","label":"a"},{"type":"open","target":"javascript:alert(1)"},{"type":"eval","code":"1"},{"type":"connection","fields":{"host":"10.0.0.5","evil":"x"}}]\n```'));
    assert.deepEqual(acts.map(a => a.type), ['open', 'connection']);
    assert.equal(acts[1].fields.evil, undefined);
    const html = await page.evaluate(() => wrmHelp.mdHTML('**b** <img src=x onerror=alert(1)> [x](javascript:alert(1)) [ok](https://example.com)'));
    assert.ok(!html.includes('<img') && !html.includes('href="javascript'), 'markdown is escaped: ' + html);

    // the AI tab without a service shows the setup button; the checklist opens
    await page.click('#hp-tab-ask');
    await page.waitForSelector('#hp-body button');
    await page.evaluate(() => openHelpAISetup());
    await page.waitForSelector('#ai-setup .ais-step');
    assert.equal(await page.locator('#ai-setup .ais-step').count(), 4);
    if (shots) await page.screenshot({path: path.join(shots, 'help-ai-setup.png')});
    await page.evaluate(() => wrmHelp.closeSetup());

    // F1 closes and opens it
    await page.keyboard.press('F1');
    assert.equal(await page.evaluate(() => document.body.classList.contains('help-open')), false);
    await page.keyboard.press('F1');
    assert.equal(await page.evaluate(() => document.body.classList.contains('help-open')), true);

    // an error toast has a ? that opens Help with the error
    await page.evaluate(() => showToast('Connection refused (demo)', 'error', 8000));
    await page.click('.toast.error .toast-help');
    await page.waitForSelector('#help-panel .hp-ctx, #help-panel .hint-box');
    assert.deepEqual(errors, [], 'no page errors');
    await ctx.close();

    // phone: the panel fills the screen, ? is in the account menu
    const ph = await signedIn(browser, {viewport: {width: 390, height: 844}, isMobile: true, hasTouch: true}, 'hr');
    await ph.page.evaluate(() => openHelp());
    await ph.page.waitForSelector('#help-panel .hp-item');
    const w = await ph.page.evaluate(() => document.getElementById('help-panel').getBoundingClientRect().width);
    assert.ok(w >= 389, 'full width on a phone: ' + w);
    if (shots) await ph.page.screenshot({path: path.join(shots, 'help-phone.png')});
    assert.deepEqual(ph.errors, [], 'no page errors on the phone');
    await ph.ctx.close();
  } finally {
    await browser.close();
  }
});
