// Generates the GitHub wiki from the in-app help articles (remote-manager/help/en|hr),
// the README (what's new, comparison, quick start) and the screenshots of screenshots.js.
//   node tools/wiki/build-wiki.js <out dir> [images dir]
// The out dir becomes the wiki repository's content (Home, _Sidebar, _Footer, one page per
// article and language, keyword indexes). Images are copied to <out>/images.
'use strict';
const fs = require('fs');
const path = require('path');
const {readArticles} = require('./articles');

const root = path.join(__dirname, '..', '..');
const REPO = 'https://github.com/vedranius/web-browser-RDM-public';
const BLOB = `${REPO}/blob/main/`;

const L = {
  en: {home: 'Home', index: 'Keyword index', news: "What's new", compare: 'How WRM compares', quick: 'Quick start',
    other: 'Hrvatski', shot: 'Screenshot', keys: 'Search terms', admin: 'Administrators only', prev: 'Previous', next: 'Next',
    sections: 'Topics', search: 'Search: use the **Search wiki** box at the top of this page (it searches every page), or open the',
    inApp: 'The same article is in WRM: **?** in the top bar (or **F1**) opens Help, where you can search it and ask the help AI.',
    phone: 'On a phone', mainShot: 'The main screen with invented demo servers',
    intro: '**WRM PRO (Web Remote Manager)** is a remote server manager that runs in any web browser: SSH terminals, SFTP/FTP files, RDP/VNC desktops, tunnels, keys and credentials, monitoring, logs, Git deployment, sharing, an AI assistant and much more, in one Go binary.'},
  hr: {home: 'Početna', index: 'Indeks pojmova', news: 'Novosti', compare: 'Usporedba', quick: 'Brzi početak',
    other: 'English', shot: 'Snimka zaslona', keys: 'Pojmovi za pretragu', admin: 'Samo administratori', prev: 'Prethodno', next: 'Sljedeće',
    sections: 'Teme', search: 'Pretraga: upišite pojam u okvir **Search wiki** na vrhu stranice (pretražuje sve stranice) ili otvorite',
    inApp: 'Isti članak je i u WRM-u: **?** u gornjoj traci (ili **F1**) otvara Pomoć, gdje ga možete pretražiti i pitati AI pomoći.',
    phone: 'Na mobitelu', mainShot: 'Glavni zaslon s izmišljenim demo serverima',
    intro: '**WRM PRO (Web Remote Manager)** je upravitelj udaljenim serverima koji radi u bilo kojem web pregledniku: SSH terminali, SFTP/FTP datoteke, RDP/VNC radne površine, tuneli, ključevi i vjerodajnice, nadzor, logovi, Git deploy, dijeljenje, AI asistent i još puno toga, u jednoj Go binarnoj datoteci.'},
};

// pageName turns a title into a wiki page name (the file name without .md).
function pageName(title) {
  return title.replace(/[\\/:*?"<>|#%,()'’]/g, '').replace(/\s+/g, '-').replace(/-+/g, '-');
}
const link = (text, page) => `[${text}](${encodeURI(page)})`;

// readmeSection returns the body of a "## heading" section of the README, links made absolute.
function readmeSection(readme, heading) {
  const lines = readme.split('\n');
  const i = lines.findIndex(l => l.trim() === `## ${heading}`);
  if (i < 0) return '';
  let j = i + 1;
  while (j < lines.length && !/^## /.test(lines[j])) j++;
  return absLinks(lines.slice(i + 1, j).join('\n').trim());
}
function readmeSections(readme, re) {
  return readme.split('\n').filter(l => re.test(l)).map(l => l.slice(3).trim());
}
function absLinks(md) {
  return md
    .replace(/\]\(#([^)]+)\)/g, `](${BLOB}README.md#$1)`)
    .replace(/\]\((?!https?:|mailto:|#)([^)\s]+)\)/g, `](${BLOB}$1)`)
    .replace(/(src|srcset|href)="(?!https?:|#)([^"]+)"/g, `$1="${REPO}/raw/main/$2"`);
}

function build(out, images) {
  const helpDir = path.join(root, 'remote-manager', 'help');
  const readme = fs.readFileSync(path.join(root, 'README.md'), 'utf8');
  const version = (/\*\*Current version: (v[\d.]+)\*\*/.exec(readme) || [])[1] || '';
  fs.rmSync(out, {recursive: true, force: true});
  fs.mkdirSync(path.join(out, 'images'), {recursive: true});
  const have = new Set();
  if (images && fs.existsSync(images)) {
    for (const f of fs.readdirSync(images)) {
      if (/^[a-z]{2}-[a-z0-9-]+\.(jpg|png)$/.test(f)) { fs.copyFileSync(path.join(images, f), path.join(out, 'images', f)); have.add(f); }
    }
  }
  const img = (lang, id, alt) => {
    for (const ext of ['jpg', 'png']) if (have.has(`${lang}-${id}.${ext}`)) return `![${alt}](images/${lang}-${id}.${ext})`;
    return '';
  };
  const files = {};
  const put = (name, md) => { files[name] = md.replace(/\n{3,}/g, '\n\n').trim() + '\n'; };

  const arts = {en: readArticles(helpDir, 'en'), hr: readArticles(helpDir, 'hr')};
  const names = {en: {}, hr: {}};
  for (const a of arts.en) names.en[a.id] = pageName(a.title);
  for (const a of arts.hr) {
    let n = pageName(a.title);
    if (Object.values(names.en).includes(n)) n += '-HR';
    names.hr[a.id] = n;
  }
  const special = {
    en: {home: 'Home', index: 'Keyword-index', news: 'Whats-new', compare: 'How-WRM-compares', quick: 'Quick-start'},
    hr: {home: 'Početna', index: 'Indeks-pojmova'},
  };
  const sectionsOf = (list) => {
    const s = [];
    for (const a of list) {
      let g = s.find(x => x.name === a.section);
      if (!g) s.push(g = {name: a.section, items: []});
      g.items.push(a);
    }
    return s;
  };

  for (const lang of ['en', 'hr']) {
    const t = L[lang];
    const other = lang === 'en' ? 'hr' : 'en';
    const list = arts[lang];
    // one page per article
    list.forEach((a, i) => {
      const peer = names[other][a.id];
      const nav = [];
      if (i > 0) nav.push(`← ${t.prev}: ${link(list[i - 1].title, names[lang][list[i - 1].id])}`);
      nav.push(link(t.home, special[lang].home));
      if (i < list.length - 1) nav.push(`${t.next}: ${link(list[i + 1].title, names[lang][list[i + 1].id])} →`);
      put(names[lang][a.id], [
        `*${a.section}* · ${peer ? link(t.other, peer) : ''}${a.admin ? ` · **${t.admin}**` : ''}`,
        `> ${a.summary}`,
        img(lang, a.id, `${t.shot}: ${a.title}`),
        a.body,
        '---',
        `**${t.keys}:** ${a.keywords.join(', ')}`,
        t.inApp,
        nav.join(' · '),
      ].join('\n\n'));
    });
    // keyword index
    const kw = new Map();
    for (const a of list) for (const k of [a.title, ...a.keywords]) {
      const key = k.toLowerCase();
      if (!kw.has(key)) kw.set(key, {k, ids: new Set()});
      kw.get(key).ids.add(a.id);
    }
    const sorted = [...kw.values()].sort((x, y) => x.k.localeCompare(y.k, lang));
    let letter = '';
    const idx = [];
    for (const e of sorted) {
      const first = e.k[0].toUpperCase();
      if (first !== letter) { letter = first; idx.push(`\n### ${letter}\n`); }
      idx.push(`- **${e.k}**: ${[...e.ids].map(id => link(list.find(a => a.id === id).title, names[lang][id])).join(', ')}`);
    }
    put(special[lang].index, idx.join('\n'));
    // home
    const home = [
      `${lang === 'en' ? link('Hrvatski', 'Početna') : link('English', 'Home')}`,
      t.intro,
      version ? `**${lang === 'en' ? 'Current version' : 'Trenutna verzija'}: ${version}** · [Download](${REPO}/releases/latest) · [Changelog](${BLOB}CHANGELOG.md)` : '',
      img(lang, 'main', t.mainShot),
      `${t.search} ${link(t.index, special[lang].index)}.`,
      lang === 'en' ? `${link(t.quick, special.en.quick)} · ${link(t.news, special.en.news)} · ${link(t.compare, special.en.compare)}` : `${link('Brzi početak (engleski)', special.en.quick)} · ${link('Novosti (engleski)', special.en.news)} · ${link('Usporedba (engleski)', special.en.compare)}`,
      `## ${t.sections}`,
    ];
    for (const g of sectionsOf(list)) {
      home.push(`### ${g.name}`);
      home.push(`| | |\n|---|---|\n` + g.items.map(a => `| ${link(a.title, names[lang][a.id])}${a.admin ? ' 🛡' : ''} | ${a.summary} |`).join('\n'));
    }
    const phone = img(lang, 'phone', t.phone);
    if (phone) home.push(`### ${t.phone}`, phone.replace(/\)$/, ' "390 px")'));
    put(special[lang].home, home.join('\n\n'));
  }

  // pages taken from the README (English)
  const news = readmeSections(readme, /^## What's new in v/).map(h => `## ${h.replace("What's new in ", '')}\n\n${readmeSection(readme, h)}`);
  put(special.en.news, `${news.join('\n\n')}\n\nOlder versions: [docs/WHATS-NEW.md](${BLOB}docs/WHATS-NEW.md) · [CHANGELOG.md](${BLOB}CHANGELOG.md)`);
  put(special.en.compare, readmeSection(readme, 'How WRM compares'));
  put(special.en.quick, readmeSection(readme, 'Quick start'));

  // sidebar and footer
  const side = [`**${link('Home', 'Home')} · ${link('Početna', 'Početna')}**`, `${link('Quick start', special.en.quick)} · ${link("What's new", special.en.news)} · ${link('Comparison', special.en.compare)}`];
  for (const lang of ['en', 'hr']) {
    side.push(`### ${lang === 'en' ? 'English' : 'Hrvatski'}`, link(L[lang].index, special[lang].index));
    for (const g of sectionsOf(arts[lang])) {
      side.push(`<details${lang === 'en' ? ' open' : ''}><summary><b>${g.name}</b></summary>\n\n${g.items.map(a => `- ${link(a.title, names[lang][a.id])}`).join('\n')}\n\n</details>`);
    }
  }
  put('_Sidebar', side.join('\n\n'));
  put('_Footer', `WRM PRO ${version} · generated from the [in-app help](${REPO}/tree/main/remote-manager/help) on every release · screenshots use invented demo data · [Repository](${REPO}) · [Issues](${REPO}/issues)`);

  for (const [name, md] of Object.entries(files)) fs.writeFileSync(path.join(out, `${name}.md`), md);
  return Object.keys(files).length;
}

if (require.main === module) {
  const out = process.argv[2];
  if (!out) { console.error('usage: node tools/wiki/build-wiki.js <out dir> [images dir]'); process.exit(2); }
  console.log(`${build(out, process.argv[3])} wiki pages in ${out}`);
}
module.exports = {build, pageName, readmeSection, absLinks};
