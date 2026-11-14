// Reads the in-app help articles (remote-manager/help/<lang>/*.md): front matter + body.
// The same format help.go parses; used by the wiki generator and the screenshot script.
'use strict';
const fs = require('fs');
const path = require('path');

function parseArticle(raw, file) {
  const m = /^---\r?\n([\s\S]*?)\r?\n---\r?\n([\s\S]*)$/.exec(raw);
  if (!m) throw new Error(`${file}: no front matter`);
  const a = {file, body: m[2].trim(), order: 0, keywords: [], admin: false};
  for (const line of m[1].split(/\r?\n/)) {
    const i = line.indexOf(':');
    if (i < 0) continue;
    const k = line.slice(0, i).trim();
    const v = line.slice(i + 1).trim();
    if (k === 'order') a.order = parseInt(v, 10) || 0;
    else if (k === 'keywords') a.keywords = v.split(',').map(s => s.trim()).filter(Boolean);
    else if (k === 'admin') a.admin = v === 'true';
    else a[k] = v;
  }
  if (!a.id) a.id = path.basename(file, '.md');
  if (!a.title || !a.section || !a.summary || !a.body) throw new Error(`${file}: title, section, summary and a body are required`);
  return a;
}

// readArticles returns one language's articles sorted by order.
function readArticles(helpDir, lang) {
  const dir = path.join(helpDir, lang);
  return fs.readdirSync(dir).filter(f => f.endsWith('.md')).sort()
    .map(f => parseArticle(fs.readFileSync(path.join(dir, f), 'utf8'), `${lang}/${f}`))
    .sort((x, y) => x.order - y.order || x.id.localeCompare(y.id));
}

module.exports = {parseArticle, readArticles};
