// Takes the screenshots for the GitHub wiki: the main screen and, for every help article,
// the feature it describes (the article's "open" target) or the help panel showing it.
// English and Croatian. Needs a running WRM with an admin account:
//   WRM_URL=http://localhost:8780 WRM_USER=admin WRM_PASS=... node tools/wiki/screenshots.js <out dir>
// An empty WRM is filled with invented demo folders and connections first (no real hosts).
// CHROMIUM_PATH points at a browser binary when the Playwright one is not installed.
'use strict';
const fs = require('fs');
const path = require('path');
const {readArticles} = require('./articles');

const url = (process.env.WRM_URL || 'http://localhost:8780').replace(/\/$/, '');
const out = process.argv[2] || 'wiki-images';
const helpDir = path.join(__dirname, '..', '..', 'remote-manager', 'help');

function playwright() {
  try { return require('playwright'); } catch (_) { return require('playwright-core'); }
}

const DEMO = {
  folders: ['Production', 'Staging', 'Network'],
  connections: [
    ['demo-web-01', 'ssh', '192.0.2.10', 'Production', ['web', 'nginx']],
    ['demo-web-02', 'ssh', '192.0.2.11', 'Production', ['web', 'nginx']],
    ['demo-api', 'ssh', '192.0.2.20', 'Production', ['api']],
    ['demo-db', 'ssh', '192.0.2.30', 'Production', ['postgres']],
    ['demo-win-app', 'rdp', '192.0.2.40', 'Production', ['windows']],
    ['demo-stage-web', 'ssh', '198.51.100.10', 'Staging', ['web']],
    ['demo-stage-db', 'ssh', '198.51.100.30', 'Staging', ['postgres']],
    ['demo-switch', 'telnet', '203.0.113.2', 'Network', ['switch']],
    ['demo-router-ui', 'https', '203.0.113.1', 'Network', ['router']],
    ['demo-console', 'vnc', '203.0.113.5', 'Network', []],
  ],
};

async function login(page, lang) {
  await page.goto(url);
  await page.evaluate(async ([u, p, l]) => {
    localStorage.setItem('wrm_lang', l);
    const r = await fetch('/api/auth/login', {method: 'POST', headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'}, body: JSON.stringify({username: u, password: p})});
    if (!r.ok) throw new Error('login failed: ' + r.status);
  }, [process.env.WRM_USER || 'admin', process.env.WRM_PASS || '', lang]);
}

async function seed(page) {
  await page.evaluate(async (demo) => {
    const h = {'Content-Type': 'application/json', 'X-WRM-Request': '1'};
    const have = await (await fetch('/api/connections')).json();
    if (Array.isArray(have) && have.length) return;
    const ids = {};
    for (const name of demo.folders) {
      const r = await fetch('/api/folders', {method: 'POST', headers: h, body: JSON.stringify({name})});
      ids[name] = (await r.json()).id;
    }
    for (const [name, protocol, host, folder, tags] of demo.connections) {
      await fetch('/api/connections', {method: 'POST', headers: h, body: JSON.stringify({
        name, protocol, host, username: 'admin', auth_method: 'password', folder_id: ids[folder], tags, monitor: false,
        notes: name === 'demo-api' ? '## Runbook\nRestart: `systemctl restart demo-api`' : undefined})});
    }
  }, DEMO);
}

async function fresh(page) {
  await page.goto(url);
  await page.waitForSelector('#top-bar', {state: 'visible'});
  await page.waitForTimeout(700);
}

async function main() {
  fs.mkdirSync(out, {recursive: true});
  const {chromium} = playwright();
  const browser = await chromium.launch(process.env.CHROMIUM_PATH ? {executablePath: process.env.CHROMIUM_PATH} : {});
  const errors = [];
  let count = 0;
  try {
    for (const lang of ['en', 'hr']) {
      const locale = lang === 'hr' ? 'hr-HR' : 'en-GB';
      const ctx = await browser.newContext({viewport: {width: 1440, height: 900}, deviceScaleFactor: 1, locale});
      const page = await ctx.newPage();
      let at = 'main';
      page.on('pageerror', e => errors.push(`${lang}/${at}: ${e.message}`));
      await login(page, lang);
      await seed(page);
      await fresh(page);
      const shot = async (name) => {
        await page.screenshot({path: path.join(out, `${lang}-${name}.jpg`), type: 'jpeg', quality: 82});
        count++;
      };
      // main screen with the demo folders expanded
      for (const row of await page.locator('.folder-row').all()) await row.click().catch(() => {});
      await page.waitForTimeout(400);
      await shot('main');
      for (const a of readArticles(helpDir, lang)) {
        at = a.id;
        try {
          await fresh(page);
          await page.evaluate((id) => wrmHelp.open(id), a.id);
          await page.waitForSelector('#help-panel .hp-md', {timeout: 5000});
          const feature = a.open && !a.open.startsWith('help:');
          if (feature) {
            await page.evaluate(() => { wrmHelp.openFeature(); wrmHelp.close(); });
            await page.waitForTimeout(1200);
          }
          await shot(a.id);
          if (a.open === 'lock') await login(page, lang); // the lock screen signs out of this tab
        } catch (e) {
          errors.push(`${lang}/${a.id}: ${e.message.split('\n')[0]}`);
        }
      }
      // the same screen at phone width
      await ctx.close();
      const phone = await browser.newContext({viewport: {width: 390, height: 844}, deviceScaleFactor: 2, isMobile: true, hasTouch: true, locale});
      const pp = await phone.newPage();
      await login(pp, lang);
      await pp.goto(url);
      await pp.waitForSelector('#top-bar', {state: 'visible'});
      await pp.waitForTimeout(700);
      await pp.click('.btn-menu').catch(() => {}); // the connection drawer with the demo tree
      for (const row of await pp.locator('.folder-row').all()) await row.click().catch(() => {});
      await pp.waitForTimeout(500);
      await pp.screenshot({path: path.join(out, `${lang}-phone.jpg`), type: 'jpeg', quality: 82});
      count++;
      await phone.close();
    }
  } finally {
    await browser.close();
  }
  console.log(`${count} screenshots in ${out}`);
  if (errors.length) {
    console.error(errors.join('\n'));
    process.exitCode = 1;
  }
}

main().catch(e => { console.error(e); process.exit(1); });
