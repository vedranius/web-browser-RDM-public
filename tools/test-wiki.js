// Tests of the wiki generator (tools/wiki/build-wiki.js): every help article has a page in
// both languages, every wiki link points to a page that exists, README sections are found
// and their relative links made absolute.  node --test tools/test-wiki.js
'use strict';
const test = require('node:test');
const assert = require('node:assert');
const fs = require('fs');
const os = require('os');
const path = require('path');
const {build, pageName, absLinks} = require('./wiki/build-wiki');
const {readArticles} = require('./wiki/articles');

const helpDir = path.join(__dirname, '..', 'remote-manager', 'help');

test('page names are clean', () => {
  assert.equal(pageName('SSH terminal, snippets and broadcast'), 'SSH-terminal-snippets-and-broadcast');
  assert.equal(pageName('Storages (S3)'), 'Storages-S3');
  assert.equal(pageName("What's new"), 'Whats-new');
});

test('README links become absolute', () => {
  const md = absLinks('[a](docs/X.md) [b](#quick-start) [c](https://example.com) <img src="docs/y.png">');
  assert.match(md, /\(https:\/\/github\.com\/[^)]+\/blob\/main\/docs\/X\.md\)/);
  assert.match(md, /README\.md#quick-start/);
  assert.match(md, /\(https:\/\/example\.com\)/);
  assert.match(md, /src="https:\/\/github\.com\/[^"]+\/raw\/main\/docs\/y\.png"/);
});

test('wiki: every article in both languages and no broken links', () => {
  const out = fs.mkdtempSync(path.join(os.tmpdir(), 'wrm-wiki-'));
  const img = fs.mkdtempSync(path.join(os.tmpdir(), 'wrm-wiki-img-'));
  try {
    fs.writeFileSync(path.join(img, 'en-terminal.jpg'), 'x');
    fs.writeFileSync(path.join(img, 'evil.jpg'), 'x');
    build(out, img);
    const pages = new Set(fs.readdirSync(out).filter(f => f.endsWith('.md')).map(f => f.slice(0, -3)));
    const en = readArticles(helpDir, 'en');
    const hr = readArticles(helpDir, 'hr');
    assert.deepEqual(en.map(a => a.id).sort(), hr.map(a => a.id).sort(), 'every article in en and hr');
    assert.equal(pages.size, en.length + hr.length + 9, 'articles + 2 homes, 2 indexes, 3 README pages, sidebar, footer');
    for (const p of ['Home', 'Početna', 'Keyword-index', 'Indeks-pojmova', 'Whats-new', 'How-WRM-compares', 'Quick-start', '_Sidebar', '_Footer']) assert.ok(pages.has(p), p);
    for (const f of pages) {
      const md = fs.readFileSync(path.join(out, f + '.md'), 'utf8');
      for (const m of md.matchAll(/\]\(([^)\s]+)\)/g)) {
        const target = m[1];
        if (/^(https?:|mailto:|#)/.test(target)) continue;
        if (target.startsWith('images/')) { assert.ok(fs.existsSync(path.join(out, target)), `${f}: ${target}`); continue; }
        assert.ok(pages.has(decodeURI(target)), `${f}: broken link ${target}`);
      }
    }
    const term = fs.readFileSync(path.join(out, 'SSH-terminal-snippets-and-broadcast.md'), 'utf8');
    assert.match(term, /images\/en-terminal\.jpg/, 'screenshot on the article page');
    assert.match(term, /\[Hrvatski\]/, 'link to the Croatian page');
    assert.ok(!fs.existsSync(path.join(out, 'images', 'evil.jpg')), 'only lang-id images are copied');
    assert.match(fs.readFileSync(path.join(out, 'How-WRM-compares.md'), 'utf8'), /\|/, 'comparison table from the README');
    assert.match(fs.readFileSync(path.join(out, 'Whats-new.md'), 'utf8'), /^## v\d+\.\d+\.\d+/m);
  } finally {
    fs.rmSync(out, {recursive: true, force: true});
    fs.rmSync(img, {recursive: true, force: true});
  }
});
