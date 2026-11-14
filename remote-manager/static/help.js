// Help centre: the ? button in the top bar (and F1) opens Help on the side. It searches the
// articles (GET /api/help…), opens the article about what is on the screen, lets the user pick
// a part of the screen (⌖ Pick, or drag a frame) and asks the help AI (POST /api/help/ask,
// server-sent events). The help AI never runs anything: its answer may carry a ```wrm-actions
// block; each action becomes a button that says what it does, and nothing happens until the
// user presses it. A setting change is shown again in a confirmation before it is sent.
// Also: the one-page AI setup checklist (openHelpAISetup) and a ? on error toasts.
'use strict';

Object.assign(LANGS.en, {
  help_lbl: 'Help', help_t: 'Help (F1): search, articles, pick a part of the screen, ask the help AI',
  help_search_ph: 'Search the help…', help_topics: 'Topics', help_ask_tab: 'Ask AI', help_pick: '⌖ Pick', help_pick_t: 'Point at a part of WRM and click it, or drag a frame around an area. Esc cancels.',
  help_pick_hint: 'Click a part of WRM or drag a frame around an area · Esc cancels', help_picked: 'Picked: {t}', help_about: 'About what you picked',
  help_on_screen: 'About what is on the screen', help_no_hits: 'Nothing found. Try other words, or ask the help AI.', help_github: 'Open on GitHub',
  help_back: '← Back', help_open_feature: 'Open it', help_admin_only: 'administrators',
  help_ask_ph: 'Ask in your own words, e.g. how do I reach a server behind a bastion?', help_send: 'Ask', help_stop: 'Stop', help_new_chat: 'New chat',
  help_ai_intro: 'The help AI answers from these articles. It never runs anything: when it can do something for you it shows a button that says exactly what it does.',
  help_ai_off: 'The help AI is not available: {r}', help_ai_setup: '🧭 Set up AI', help_model: 'AI service',
  help_ctx: 'Context: {t}', help_ctx_clear: 'remove', help_ask_about: 'Ask about it', help_err_ask: 'Explain this error',
  help_act_title: 'Actions offered by the help AI: nothing happens until you press one.',
  help_act_open: 'Opens {x}', help_act_conn: 'Opens the New connection dialog filled in ({x}); you check it, add the password or key and save.',
  help_act_connect: 'Opens the saved connection {x}', help_act_assistant: 'Opens {x} with the AI assistant and puts this in its prompt box (you send it): {p}',
  help_act_setting: 'Changes the WRM setting {k} to {v} (administrators; asks again before saving)', help_act_unknown: 'Unknown action, ignored',
  help_act_confirm: 'Change the setting {k} from "{o}" to "{v}"?', help_act_done: 'Done: {x}', help_act_noconn: 'No connection named {x}',
  help_act_ok: 'Setting saved.', help_err_btn_t: 'Ask Help about this error',
  // AI setup checklist
  ais_title: 'Set up AI', ais_intro: 'Every step with a light: green when done, with a button for what is missing. Nothing is sent to an AI service until a service is added and somebody asks.',
  ais_s1: 'AI is turned on', ais_s1_d: 'Who may use the AI assistant (policy ai_assistant): {v}.', ais_s1_admin: 'Administrators', ais_s1_all: 'Everybody',
  ais_s2: 'An AI service is added', ais_s2_d: '{o} for everybody, {u} of your own. The wizard connects to the service and adds every model it offers.', ais_s2_btn: 'Add an AI service', ais_s2_own: 'Add my own key',
  ais_s3: 'Help AI', ais_s3_d: 'Policy ai_help: {v}.', ais_s4: 'AI apps (MCP), optional', ais_s4_d: 'Claude Desktop, ChatGPT, Claude Code, Cursor and other MCP apps. {n} connected.',
  ais_s4_on: 'Turn on AI connections', ais_s4_btn: 'Connect an AI app', ais_ask_admin: 'Ask an administrator.', ais_done: 'Done', ais_todo: 'To do', ais_off: 'Off',
  ais_reason: 'Why not: {r}', ais_kill: 'The AI kill switch is on (Admin → AI assistant).', ais_read_more: 'Read: Set up AI', ais_close: 'Close',
});
Object.assign(LANGS.hr, {
  help_lbl: 'Pomoć', help_t: 'Pomoć (F1): pretraga, članci, odabir dijela zaslona, pitaj AI pomoć',
  help_search_ph: 'Pretraži pomoć…', help_topics: 'Teme', help_ask_tab: 'Pitaj AI', help_pick: '⌖ Odaberi', help_pick_t: 'Pokažite dio WRM-a i kliknite ga ili povucite okvir oko područja. Esc odustaje.',
  help_pick_hint: 'Kliknite dio WRM-a ili povucite okvir oko područja · Esc odustaje', help_picked: 'Odabrano: {t}', help_about: 'O onome što ste odabrali',
  help_on_screen: 'O onome što je na zaslonu', help_no_hits: 'Ništa nije pronađeno. Probajte drugim riječima ili pitajte AI pomoć.', help_github: 'Otvori na GitHubu',
  help_back: '← Natrag', help_open_feature: 'Otvori', help_admin_only: 'administratori',
  help_ask_ph: 'Pitajte svojim riječima, npr. kako doći do servera iza bastiona?', help_send: 'Pitaj', help_stop: 'Zaustavi', help_new_chat: 'Novi razgovor',
  help_ai_intro: 'AI pomoć odgovara iz ovih članaka. Nikad ništa ne pokreće: kad nešto može napraviti za vas, pokaže gumb koji točno kaže što radi.',
  help_ai_off: 'AI pomoć nije dostupna: {r}', help_ai_setup: '🧭 Postavi AI', help_model: 'AI servis',
  help_ctx: 'Kontekst: {t}', help_ctx_clear: 'ukloni', help_ask_about: 'Pitaj o tome', help_err_ask: 'Objasni ovu grešku',
  help_act_title: 'Radnje koje nudi AI pomoć: ništa se ne događa dok ne pritisnete jednu.',
  help_act_open: 'Otvara {x}', help_act_conn: 'Otvara dijalog Nova veza s popunjenim poljima ({x}); vi provjerite, dodate lozinku ili ključ i spremite.',
  help_act_connect: 'Otvara spremljenu vezu {x}', help_act_assistant: 'Otvara {x} s AI asistentom i stavlja ovo u polje za upit (vi šaljete): {p}',
  help_act_setting: 'Mijenja WRM postavku {k} u {v} (administratori; prije spremanja pita još jednom)', help_act_unknown: 'Nepoznata radnja, zanemarena',
  help_act_confirm: 'Promijeniti postavku {k} iz "{o}" u "{v}"?', help_act_done: 'Gotovo: {x}', help_act_noconn: 'Nema veze s imenom {x}',
  help_act_ok: 'Postavka je spremljena.', help_err_btn_t: 'Pitaj Pomoć o ovoj grešci',
  ais_title: 'Postavi AI', ais_intro: 'Svaki korak sa svjetlom: zeleno kad je gotov, s gumbom za ono što nedostaje. AI servisu se ništa ne šalje dok se servis ne doda i netko ne pita.',
  ais_s1: 'AI je uključen', ais_s1_d: 'Tko smije koristiti AI asistenta (pravilo ai_assistant): {v}.', ais_s1_admin: 'Administratori', ais_s1_all: 'Svi',
  ais_s2: 'Dodan je AI servis', ais_s2_d: '{o} za sve, {u} vaših. Čarobnjak se spaja na servis i dodaje sve modele koje nudi.', ais_s2_btn: 'Dodaj AI servis', ais_s2_own: 'Dodaj svoj ključ',
  ais_s3: 'AI pomoć', ais_s3_d: 'Pravilo ai_help: {v}.', ais_s4: 'AI aplikacije (MCP), neobavezno', ais_s4_d: 'Claude Desktop, ChatGPT, Claude Code, Cursor i druge MCP aplikacije. Spojeno: {n}.',
  ais_s4_on: 'Uključi AI veze', ais_s4_btn: 'Spoji AI aplikaciju', ais_ask_admin: 'Pitajte administratora.', ais_done: 'Gotovo', ais_todo: 'Za napraviti', ais_off: 'Isključeno',
  ais_reason: 'Zašto ne: {r}', ais_kill: 'AI prekidač za isključivanje je uključen (Admin → AI asistent).', ais_read_more: 'Pročitaj: Postavi AI', ais_close: 'Zatvori',
});

const wrmHelp = (() => {
  const $ = id => document.getElementById(id);
  const tf = (k, v) => t(k).replace(/\{(\w+)\}/g, (m, x) => (v && v[x] != null ? v[x] : m));
  const lang = () => (typeof currentLang === 'string' && currentLang === 'hr' ? 'hr' : 'en');
  const css = `
#help-panel { position: fixed; top: var(--topbar-h, 44px); right: 0; bottom: 0; width: min(420px, 100vw); z-index: 9600; display: none; flex-direction: column;
  background: var(--bg2, #1e1e2e); border-left: 1px solid var(--border, #444); box-shadow: -6px 0 24px rgba(0,0,0,.35); font-size: 13px; }
body.help-open #help-panel { display: flex; }
#help-panel .hp-head { display: flex; align-items: center; gap: 8px; padding: 8px 10px; border-bottom: 1px solid var(--border, #444); }
#help-panel .hp-search { padding: 8px 10px 4px; } #help-panel .hp-search input { width: 100%; }
#help-panel .hp-tabs { display: flex; gap: 4px; padding: 4px 10px; } #help-panel .hp-tabs button { flex: 1; background: transparent; color: var(--text2); border: 1px solid var(--border, #444); padding: 4px; }
#help-panel .hp-tabs button.on { background: var(--accent, #4a7); color: var(--accent-fg, #fff); border-color: transparent; }
#help-panel .hp-body { flex: 1; overflow: auto; padding: 8px 12px; }
#help-panel .hp-foot { padding: 8px 10px; border-top: 1px solid var(--border, #444); display: flex; flex-direction: column; gap: 6px; }
#help-panel .hp-foot:empty { display: none; } #help-panel .hp-foot textarea { width: 100%; resize: vertical; }
#help-panel .hp-row { display: flex; gap: 6px; align-items: center; } #help-panel .hp-row select { flex: 1; min-width: 0; }
#help-panel .hp-sec { margin: 12px 0 4px; font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: .04em; color: var(--text2); }
#help-panel .hp-item { display: block; width: 100%; text-align: left; background: transparent; color: var(--text); border: 0; border-radius: 6px; padding: 6px 8px; cursor: pointer; }
#help-panel .hp-item:hover, #help-panel .hp-item:focus-visible { background: var(--bg3, rgba(255,255,255,.06)); }
#help-panel .hp-item span { display: block; font-size: 11.5px; color: var(--text2); margin-top: 2px; }
#help-panel .hp-h { margin: 10px 0 6px; } #help-panel .hp-md h3, #help-panel .hp-md h4 { margin: 12px 0 4px; font-size: 13px; }
#help-panel .hp-md p { margin: 4px 0; line-height: 1.5; } #help-panel .hp-md ul, #help-panel .hp-md ol { margin: 4px 0; padding-left: 20px; line-height: 1.5; }
#help-panel pre { white-space: pre-wrap; font: 12px var(--font-mono); background: var(--bg, #111); padding: 6px; border-radius: 6px; }
#help-panel .hp-gh { margin: 12px 0 4px; display: flex; gap: 8px; align-items: center; flex-wrap: wrap; font-size: 12px; }
#help-panel .hp-ctx { font-size: 11.5px; color: var(--text2); margin: 4px 0 8px; }
#help-panel .ai-see, #ai-setup .ai-see { font-size: 11.5px; color: var(--text2); margin: 4px 0; line-height: 1.45; }
#help-panel .hp-msg { margin: 8px 0; padding: 8px 10px; border-radius: 8px; line-height: 1.5; }
#help-panel .hp-msg.me { background: var(--bg3, rgba(255,255,255,.06)); white-space: pre-wrap; } #help-panel .hp-msg.ai { border: 1px solid var(--border, #444); }
#help-panel .hp-msg p { margin: 3px 0; } #help-panel .hp-meta { font-size: 10.5px; color: var(--text2); margin-top: 6px; }
#help-panel .hp-acts { margin-top: 8px; border-top: 1px dashed var(--border, #444); padding-top: 6px; }
#help-panel .hp-act { display: flex; gap: 8px; align-items: flex-start; margin: 6px 0; } #help-panel .hp-act span { font-size: 11px; color: var(--text2); }
#help-pick { position: fixed; inset: 0; z-index: 9900; cursor: crosshair; }
#help-pick .hp-pick-box { position: fixed; display: none; border: 2px solid var(--accent, #4a7); background: rgba(80,160,255,.12); pointer-events: none; border-radius: 4px; }
#help-pick .hp-pick-hint { position: fixed; left: 50%; top: 12px; transform: translateX(-50%); background: var(--bg2, #222); color: var(--text); padding: 6px 12px; border-radius: 6px; box-shadow: 0 2px 10px rgba(0,0,0,.4); font-size: 12px; }
#ai-setup .ais-step { display: flex; gap: 10px; align-items: flex-start; padding: 8px 0; border-top: 1px solid var(--border, #444); }
#ai-setup .ais-light { font-size: 18px; line-height: 1; } #ai-setup .ais-light.ok { color: var(--green-t, #4c8); } #ai-setup .ais-light.todo { color: var(--yellow-t, #db4); }
#ai-setup .ais-btns { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 4px; font-size: 12px; }
.toast .toast-help { margin-left: 8px; background: transparent; border: 1px solid currentColor; color: inherit; border-radius: 50%; width: 20px; height: 20px; padding: 0; font-size: 11px; cursor: pointer; flex: none; }
@media (max-width: 600px) { #help-panel { width: 100vw; border-left: 0; } }
`;
  if (!$('help-style')) { const st = document.createElement('style'); st.id = 'help-style'; st.textContent = css; document.head.appendChild(st); }
  let index = null, indexLang = '', view = 'topics', article = null, ctx = null, chat = [], status = null, busy = null, picking = null;

  const OPEN = {
    new_connection: () => openConnModal(null),
    quick_connect: () => openQuickModal(),
    folder_new: () => openFolderModal(null),
    import: () => openImportPicker(),
    inventory: () => openInventory(),
    vault: () => openVault(), 'vault:keys': () => openVault('keys'), 'vault:creds': () => openVault('creds'), 'vault:proxies': () => openVault('proxies'),
    tunnels: () => openTunnelsPanel(), tools: () => openTools(), broadcast: () => openBroadcastModal(),
    monitoring: () => openMonitoring(), 'monitoring:logs': () => openMonitoring('logs'),
    git: () => openGit(), storages: () => openStorages(), chains: () => openChains(), shares: () => openManageShares(),
    panels: () => window.wrmLayout && wrmLayout.openPanels($('btn-panels')), workboard: () => window.wrmWorkboard && wrmWorkboard.menu($('btn-workboard')),
    snippets: () => openSnippetPicker(), bookmarks: () => openBookmarksManager(), history: () => openSettings('st-history'),
    ai_setup: () => openHelpAISetup(), mcp_setup: () => openSettings('st-mcp'), addons: () => openSettings('st-addons'),
    appearance: () => openSettings('st-appearance'), lock: () => window.wrmLock && wrmLock.lock(), updates: () => openUpdateModal(),
  };
  const SETTINGS_TABS = ['st-general', 'st-appearance', 'st-terminal', 'st-security', 'st-notify', 'st-audio', 'st-history', 'st-ai', 'st-mcp', 'st-addons', 'st-data'];
  const ADMIN_TABS = ['overview', 'security', 'users', 'teams', 'sso', 'shares', 'policies', 'ai', 'mcp', 'voice', 'notify', 'mail', 'hosts', 'audit', 'sessions', 'transfers', 'edits'];

  // openTarget runs one "open" target from the whitelist; anything else is refused.
  function openTarget(target) {
    target = String(target || '');
    if (OPEN[target]) { OPEN[target](); return true; }
    const [kind, arg] = target.split(/:(.*)/s);
    if (kind === 'settings' && SETTINGS_TABS.includes(arg)) { openSettings(arg); return true; }
    if (kind === 'admin' && ADMIN_TABS.includes(arg) && currentUser && currentUser.is_admin) { openAdminPanel(); showAdminTab(arg); return true; }
    if (kind === 'help' && /^[a-z0-9-]{1,60}$/.test(arg)) { open(arg); return true; }
    return false;
  }
  function validTarget(target) {
    target = String(target || '');
    if (OPEN[target]) return true;
    const [kind, arg] = target.split(/:(.*)/s);
    return (kind === 'settings' && SETTINGS_TABS.includes(arg)) || (kind === 'admin' && ADMIN_TABS.includes(arg)) || (kind === 'help' && /^[a-z0-9-]{1,60}$/.test(arg || ''));
  }

  // ── small Markdown renderer (articles and answers); everything is escaped first ──
  function inline(s) {
    return esc(s).replace(/`([^`\n]+)`/g, '<code>$1</code>').replace(/\*\*([^*\n]+)\*\*/g, '<b>$1</b>').replace(/(^|[\s(])\*([^*\n]+)\*/g, '$1<i>$2</i>')
      .replace(/\[([^\]\n]+)\]\((https:\/\/[^)\s"']+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>');
  }
  function mdHTML(src) {
    const out = [];
    let list = '';
    const close = () => { if (list) { out.push(`</${list}>`); list = ''; } };
    const parts = String(src || '').split(/```/);
    parts.forEach((p, i) => {
      if (i % 2) { close(); if (!/^wrm-actions/.test(p)) out.push('<pre>' + esc(p.replace(/^[a-z0-9_+-]*\n/i, '')) + '</pre>'); return; }
      for (const line of p.split('\n')) {
        let m;
        if ((m = line.match(/^(#{1,4})\s+(.*)/))) { close(); out.push(`<h${Math.min(m[1].length + 2, 6)}>${inline(m[2])}</h${Math.min(m[1].length + 2, 6)}>`); }
        else if ((m = line.match(/^\s*[-*]\s+(.*)/))) { if (list !== 'ul') { close(); out.push('<ul>'); list = 'ul'; } out.push(`<li>${inline(m[1])}</li>`); }
        else if ((m = line.match(/^\s*\d+[.)]\s+(.*)/))) { if (list !== 'ol') { close(); out.push('<ol>'); list = 'ol'; } out.push(`<li>${inline(m[1])}</li>`); }
        else if (!line.trim()) close();
        else { close(); out.push(`<p>${inline(line)}</p>`); }
      }
    });
    close();
    return out.join('');
  }

  // ── panel ──
  function ensurePanel() {
    let p = $('help-panel');
    if (p) return p;
    p = document.createElement('aside');
    p.id = 'help-panel';
    p.setAttribute('role', 'complementary');
    p.innerHTML = `<div class="hp-head"><b>❓ <span id="hp-ttl"></span></b><span style="flex:1"></span>
      <button class="btn-sec btn-sm" id="hp-pick" onclick="wrmHelp.pick()"></button><button class="icon-btn" onclick="wrmHelp.close()" aria-label="Close">✕</button></div>
      <div class="hp-search"><input type="search" id="hp-q" autocomplete="off" oninput="wrmHelp.search(this.value)" onkeydown="if(event.key==='Escape'){this.value='';wrmHelp.search('')}"></div>
      <div class="hp-tabs"><button id="hp-tab-topics" onclick="wrmHelp.show('topics')"></button><button id="hp-tab-ask" onclick="wrmHelp.show('ask')"></button></div>
      <div class="hp-body" id="hp-body"></div><div class="hp-foot" id="hp-foot"></div>`;
    document.body.appendChild(p);
    return p;
  }
  function labels() {
    $('hp-ttl').textContent = t('help_lbl');
    $('hp-pick').textContent = t('help_pick'); $('hp-pick').title = t('help_pick_t');
    $('hp-q').placeholder = t('help_search_ph');
    $('hp-tab-topics').textContent = t('help_topics'); $('hp-tab-ask').textContent = '✦ ' + t('help_ask_tab');
  }
  async function loadIndex() {
    if (index && indexLang === lang()) return index;
    const r = await fetch('/api/help?lang=' + lang());
    if (!r.ok) throw new Error(await apiError(r));
    index = await r.json(); indexLang = lang();
    return index;
  }
  // articles whose UI selectors match an element (or what is open on the screen)
  function articlesFor(el) {
    if (!index) return [];
    const hits = [];
    for (const a of index.articles) {
      for (const sel of a.ui || []) {
        try {
          if (el ? el.closest(sel) || (el.querySelector && el.querySelector(sel)) : visible(document.querySelector(sel))) { hits.push(a); break; }
        } catch (e) { /* not a valid selector */ }
      }
    }
    return hits;
  }
  function visible(el) {
    if (!el || el.closest('#help-panel')) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0 && getComputedStyle(el).visibility !== 'hidden' && !!el.closest('.modal-overlay.open, .open, #top-bar, body > :not(.modal-overlay)');
  }
  function onScreen() {
    const open = [...document.querySelectorAll('.modal-overlay.open')].filter(m => !m.closest('#help-panel'));
    const top = open[open.length - 1];
    return top ? articlesFor(top) : [];
  }

  async function open(id) {
    ensurePanel(); labels();
    document.body.classList.add('help-open');
    try { await loadIndex(); } catch (e) { $('hp-body').innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    if (id) return showArticle(id);
    if (view === 'article' && article) return showArticle(article.id);
    show(view === 'ask' ? 'ask' : 'topics');
  }
  function close() {
    document.body.classList.remove('help-open');
    if (picking) stopPick();
    if (busy) busy.abort();
  }
  function toggle() { document.body.classList.contains('help-open') ? close() : open(); }

  function show(v) {
    view = v;
    ensurePanel();
    $('hp-tab-topics').classList.toggle('on', v !== 'ask');
    $('hp-tab-ask').classList.toggle('on', v === 'ask');
    if (v === 'ask') return renderAsk();
    $('hp-foot').innerHTML = '';
    const q = $('hp-q').value.trim();
    if (q) return search(q);
    const here = onScreen();
    let html = '';
    if (here.length) html += `<div class="hp-sec">${esc(t('help_on_screen'))}</div>` + here.map(itemHTML).join('');
    let sec = '';
    for (const a of index.articles) {
      if (a.section !== sec) { sec = a.section; html += `<div class="hp-sec">${esc(sec)}</div>`; }
      html += itemHTML(a);
    }
    html += `<div class="hp-gh"><a href="${esc(index.github)}" target="_blank" rel="noopener noreferrer">${esc(t('help_github'))} ↗</a></div>`;
    $('hp-body').innerHTML = html;
  }
  function itemHTML(a, snippet) {
    return `<button class="hp-item" onclick="wrmHelp.open('${esc(a.id)}')"><b>${esc(a.title)}</b>${a.admin ? ` <span class="pill">${esc(t('help_admin_only'))}</span>` : ''}<span>${esc(snippet || a.summary)}</span></button>`;
  }
  let searchSeq = 0;
  async function search(q) {
    q = String(q || '').trim();
    if (!q) { if (view !== 'ask') show('topics'); return; }
    if (view === 'ask' || view === 'article') { view = 'topics'; $('hp-tab-topics').classList.add('on'); $('hp-tab-ask').classList.remove('on'); $('hp-foot').innerHTML = ''; }
    const seq = ++searchSeq;
    try {
      const r = await fetch(`/api/help/search?lang=${lang()}&q=${encodeURIComponent(q)}`);
      if (!r.ok) throw new Error(await apiError(r));
      const hits = await r.json();
      if (seq !== searchSeq) return;
      $('hp-body').innerHTML = hits.length ? hits.map(h => itemHTML(h, h.snippet)).join('') : `<div class="chat-empty">${esc(t('help_no_hits'))}</div>`;
    } catch (e) { $('hp-body').innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; }
  }
  async function showArticle(id) {
    view = 'article';
    try {
      const r = await fetch(`/api/help/articles/${encodeURIComponent(id)}?lang=${lang()}`);
      if (!r.ok) throw new Error(await apiError(r));
      article = await r.json();
    } catch (e) { $('hp-body').innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    const a = article;
    $('hp-body').innerHTML = `<button class="btn-sec btn-sm" onclick="wrmHelp.show('topics')">${esc(t('help_back'))}</button>
      <h3 class="hp-h">${esc(a.title)}</h3><div class="hp-md">${mdHTML(a.markdown)}</div>
      <div class="hp-gh">${a.open && validTarget(a.open) && !a.open.startsWith('help:') ? `<button class="btn-sm" onclick="wrmHelp.openFeature()">${esc(t('help_open_feature'))}</button> ` : ''}
      <button class="btn-sec btn-sm" onclick="wrmHelp.askAbout()">✦ ${esc(t('help_ask_about'))}</button> <a href="${esc(a.github)}" target="_blank" rel="noopener noreferrer">${esc(t('help_github'))} ↗</a></div>`;
    $('hp-body').scrollTop = 0;
    $('hp-foot').innerHTML = '';
  }
  function openFeature() { if (article && article.open) openTarget(article.open); }
  function askAbout() { show('ask'); const q = $('hp-ask-q'); if (q) q.focus(); }

  // ── pick a part of the screen ──
  function pick() {
    if (picking) { stopPick(); return; }
    const ov = document.createElement('div');
    ov.id = 'help-pick';
    ov.innerHTML = `<div class="hp-pick-box"></div><div class="hp-pick-hint">${esc(t('help_pick_hint'))}</div>`;
    document.body.appendChild(ov);
    const box = ov.querySelector('.hp-pick-box');
    let start = null, hover = null;
    const under = (x, y) => { ov.style.pointerEvents = 'none'; const el = document.elementFromPoint(x, y); ov.style.pointerEvents = ''; return el && !el.closest('#help-panel') ? el : null; };
    const frame = r => Object.assign(box.style, {left: r.left + 'px', top: r.top + 'px', width: r.width + 'px', height: r.height + 'px', display: 'block'});
    const move = e => {
      if (start) { frame({left: Math.min(start.x, e.clientX), top: Math.min(start.y, e.clientY), width: Math.abs(e.clientX - start.x), height: Math.abs(e.clientY - start.y)}); return; }
      hover = under(e.clientX, e.clientY);
      if (hover) frame(hover.getBoundingClientRect()); else box.style.display = 'none';
    };
    const down = e => { e.preventDefault(); start = {x: e.clientX, y: e.clientY}; };
    const up = e => {
      e.preventDefault();
      const dragged = start && (Math.abs(e.clientX - start.x) > 8 || Math.abs(e.clientY - start.y) > 8);
      let el = hover;
      if (dragged) el = commonAncestor(start, {x: e.clientX, y: e.clientY}, under);
      stopPick();
      if (el) picked(el);
    };
    const key = e => { if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); stopPick(); } };
    ov.addEventListener('pointermove', move); ov.addEventListener('pointerdown', down); ov.addEventListener('pointerup', up);
    document.addEventListener('keydown', key, true);
    picking = {ov, key};
  }
  function stopPick() {
    if (!picking) return;
    picking.ov.remove();
    document.removeEventListener('keydown', picking.key, true);
    picking = null;
  }
  // the smallest element that holds both corners of the frame
  function commonAncestor(a, b, under) {
    const x = under(a.x, a.y), y = under(b.x, b.y);
    if (!x || !y) return x || y;
    for (let el = x; el; el = el.parentElement) if (el.contains(y)) return el === document.body ? null : el;
    return null;
  }
  function describe(el) {
    const path = [];
    for (let e = el; e && e !== document.body && path.length < 6; e = e.parentElement) {
      if (e.id) { path.unshift('#' + e.id); if (/^(modal-|st-|admin-)/.test(e.id)) break; }
    }
    const dlg = el.closest('.modal, .modal-overlay, [role=dialog]');
    const h = dlg && dlg.querySelector('h3, h2');
    const own = el.getAttribute('aria-label') || el.title || (el.textContent || '').trim().slice(0, 60);
    const title = [h && h.textContent.trim().slice(0, 80), own].filter(Boolean).filter((v, i, l) => l.indexOf(v) === i).join(' › ');
    // form values are never sent: only visible text and labels
    const text = (el.innerText || el.value === undefined && el.textContent || '').replace(/\s+\n/g, '\n').slice(0, 3000);
    return {kind: 'pick', title, path: path.join(' > '), text};
  }
  function picked(el) {
    ctx = describe(el);
    const list = articlesFor(el);
    view = 'topics';
    $('hp-tab-topics').classList.add('on'); $('hp-tab-ask').classList.remove('on');
    $('hp-body').innerHTML = `<div class="hp-ctx">${esc(tf('help_picked', {t: ctx.title || ctx.path || '…'}))}</div>
      <button class="btn-sm" onclick="wrmHelp.show('ask')">✦ ${esc(t('help_ask_about'))}</button>
      ${list.length ? `<div class="hp-sec">${esc(t('help_about'))}</div>` + list.map(a => itemHTML(a)).join('') : ''}`;
    $('hp-foot').innerHTML = '';
  }

  // ── ask the help AI ──
  async function loadStatus() {
    const r = await fetch('/api/help/ai');
    if (!r.ok) throw new Error(await apiError(r));
    status = await r.json();
    return status;
  }
  async function renderAsk() {
    $('hp-foot').innerHTML = '';
    try { await loadStatus(); } catch (e) { $('hp-body').innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    if (!status.help.allowed || !status.providers.length) {
      $('hp-body').innerHTML = `<div class="hint-box warn">${esc(tf('help_ai_off', {r: status.help.allowed ? t('ais_s2') + ': ' + t('ais_todo') : status.help.reason}))}</div>
        <button class="btn-sm" onclick="openHelpAISetup()">${esc(t('help_ai_setup'))}</button>`;
      return;
    }
    renderChat();
    const provs = status.providers;
    const sel = provs.map(p => p.models.map(m => `<option value="${p.id}|${esc(m)}"${m === p.default_model && p === provs[0] ? ' selected' : ''}>${esc(p.name)} · ${esc(m)}</option>`).join('')).join('');
    $('hp-foot').innerHTML = `${ctx ? `<div class="hp-ctx">${esc(tf('help_ctx', {t: ctx.title || ctx.path || ctx.kind}))} <a href="#" onclick="wrmHelp.clearCtx();return false;">${esc(t('help_ctx_clear'))}</a></div>` : ''}
      <textarea id="hp-ask-q" rows="3" placeholder="${esc(t('help_ask_ph'))}" onkeydown="if(event.key==='Enter'&&!event.shiftKey){event.preventDefault();wrmHelp.send();}"></textarea>
      <div class="hp-row"><select id="hp-model" title="${esc(t('help_model'))}">${sel}</select>
      <button class="btn-sec btn-sm" onclick="wrmHelp.newChat()">${esc(t('help_new_chat'))}</button><button class="btn-sm" id="hp-send" onclick="wrmHelp.send()">${esc(t('help_send'))}</button></div>`;
  }
  function renderChat() {
    const body = $('hp-body');
    body.innerHTML = `<div class="ai-see">${esc(t('help_ai_intro'))}</div>` + chat.map((m, i) => m.role === 'user'
      ? `<div class="hp-msg me">${esc(m.text)}</div>`
      : `<div class="hp-msg ai">${m.error ? `<div class="fm-error">${esc(m.error)}</div>` : mdHTML(m.text)}${actionsHTML(m, i)}${m.meta ? `<div class="hp-meta">${esc(m.meta)}</div>` : ''}</div>`).join('');
    body.scrollTop = body.scrollHeight;
  }
  function clearCtx() { ctx = null; renderAsk(); }
  function newChat() { chat = []; ctx = null; renderAsk(); }

  // parseActions reads the ```wrm-actions block(s) of an answer; malformed or unknown entries are dropped.
  function parseActions(text) {
    const out = [];
    const re = /```wrm-actions\s*\n([\s\S]*?)```/g;
    let m;
    while ((m = re.exec(text)) && out.length < 4) {
      let arr;
      try { arr = JSON.parse(m[1]); } catch (e) { continue; }
      for (const a of Array.isArray(arr) ? arr : []) {
        if (out.length >= 4 || !a || typeof a !== 'object') break;
        const label = String(a.label || '').slice(0, 120);
        if (a.type === 'open' && validTarget(a.target)) out.push({type: 'open', target: String(a.target), label});
        else if (a.type === 'connection' && a.fields && typeof a.fields === 'object') {
          const f = {};
          for (const k of ['name', 'protocol', 'host', 'port', 'username', 'folder']) if (a.fields[k] != null) f[k] = String(a.fields[k]).slice(0, 200);
          out.push({type: 'connection', fields: f, label});
        } else if (a.type === 'connect' && a.name) out.push({type: 'connect', name: String(a.name).slice(0, 200), label});
        else if (a.type === 'assistant' && a.name) out.push({type: 'assistant', name: String(a.name).slice(0, 200), prompt: String(a.prompt || '').slice(0, 2000), label});
        else if (a.type === 'setting' && /^[a-z0-9_]{1,60}$/.test(String(a.key || '')) && currentUser && currentUser.is_admin) out.push({type: 'setting', key: String(a.key), value: String(a.value == null ? '' : a.value).slice(0, 2000), label});
      }
    }
    return out;
  }
  function actionWhat(a) {
    if (a.type === 'open') return tf('help_act_open', {x: a.target});
    if (a.type === 'connection') return tf('help_act_conn', {x: Object.entries(a.fields).map(([k, v]) => k + ': ' + v).join(', ')});
    if (a.type === 'connect') return tf('help_act_connect', {x: a.name});
    if (a.type === 'assistant') return tf('help_act_assistant', {x: a.name, p: a.prompt});
    if (a.type === 'setting') return tf('help_act_setting', {k: a.key, v: a.value});
    return t('help_act_unknown');
  }
  function actionsHTML(m, i) {
    if (!m.actions || !m.actions.length) return '';
    return `<div class="hp-acts"><div class="ai-see">${esc(t('help_act_title'))}</div>${m.actions.map((a, j) =>
      `<div class="hp-act"><button class="btn-sm" onclick="wrmHelp.run(${i},${j})">${esc(a.label || actionWhat(a))}</button><span>${esc(actionWhat(a))}</span></div>`).join('')}</div>`;
  }
  function findConn(name) {
    const n = String(name).toLowerCase();
    return (typeof conns !== 'undefined' ? conns : []).find(c => (c.name || '').toLowerCase() === n) || null;
  }
  function openConnWindow(c) {
    if (typeof isWebConn === 'function' && isWebConn(c)) { openWebConn(c); return null; }
    const wo = addWindow();
    wo.connId = c.id;
    setWinMode(wo.id, defaultWinMode(c));
    return wo;
  }
  // run is called only by the user's click on an action button.
  async function run(i, j) {
    const a = chat[i] && chat[i].actions && chat[i].actions[j];
    if (!a) return;
    try {
      if (a.type === 'open') { if (!openTarget(a.target)) throw new Error(t('help_act_unknown')); }
      else if (a.type === 'connection') {
        openConnModal(null);
        const set = (id, v) => { const el = $(id); if (el && v != null && v !== '') { el.value = v; el.dispatchEvent(new Event('input', {bubbles: true})); } };
        const f = a.fields;
        if (f.protocol && [...$('m-proto').options].some(o => o.value === f.protocol.toUpperCase())) { $('m-proto').value = f.protocol.toUpperCase(); onProtoChange(); }
        set('m-name', f.name); set('m-host', f.host); set('m-user', f.username); set('m-port', f.port);
        if (f.folder && typeof folders !== 'undefined') { const fo = folders.find(x => (x.name || '').toLowerCase() === f.folder.toLowerCase()); if (fo) set('m-folder', fo.id); }
      } else if (a.type === 'connect' || a.type === 'assistant') {
        const c = findConn(a.name);
        if (!c) throw new Error(tf('help_act_noconn', {x: a.name}));
        const wo = openConnWindow(c);
        if (a.type === 'assistant' && wo) {
          await loadAI();
          setTimeout(() => {
            if (typeof aiAvailable === 'function' && !aiAvailable(wo, wo.mode)) { openSettings('st-ai'); return; }
            wrmAI.toggle(wo.id);
            setTimeout(() => { const box = $('ai-in-' + wo.id); if (box) { box.value = a.prompt; box.focus(); } }, 400);
          }, 600);
        }
      } else if (a.type === 'setting') {
        const r = await fetch('/api/admin/settings');
        if (!r.ok) throw new Error(await apiError(r));
        const all = await r.json();
        const cur = Array.isArray(all) ? all.find(x => x.key === a.key) : null;
        if (!cur) throw new Error(t('help_act_unknown') + ': ' + a.key);
        const old = cur.value;
        if (!confirm(tf('help_act_confirm', {k: a.key, o: old == null ? '' : old, v: a.value}))) return;
        const w = await fetch('/api/admin/settings', {method: 'PUT', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({[a.key]: a.value})});
        if (!w.ok) throw new Error(await apiError(w));
        showToast(t('help_act_ok'), 'success');
        return;
      }
      showToast(tf('help_act_done', {x: a.label || actionWhat(a)}), 'success', 2000);
    } catch (e) { showToast(e.message, 'error'); }
  }

  async function send() {
    const q = $('hp-ask-q');
    if (busy) { busy.abort(); return; }
    const question = (q && q.value || '').trim();
    if (!question) return;
    const [pid, model] = ($('hp-model').value || '').split(/\|(.*)/s);
    const history = chat.filter(m => !m.error).slice(-16).map(m => ({role: m.role, text: m.text}));
    chat.push({role: 'user', text: question});
    const ans = {role: 'assistant', text: ''};
    chat.push(ans);
    q.value = '';
    renderChat();
    const ctrl = new AbortController();
    busy = ctrl;
    $('hp-send').textContent = t('help_stop');
    try {
      const res = await fetch('/api/help/ask', {method: 'POST', headers: {'Content-Type': 'application/json'}, signal: ctrl.signal,
        body: JSON.stringify({question, lang: lang(), article: article ? article.id : '', context: ctx || {}, history, provider_id: Number(pid) || 0, model: model || ''})});
      if (!res.ok) throw new Error(await apiError(res));
      const reader = res.body.getReader(), dec = new TextDecoder();
      let buf = '';
      for (;;) {
        const {value, done} = await reader.read();
        if (done) break;
        buf += dec.decode(value, {stream: true});
        let k;
        while ((k = buf.indexOf('\n\n')) >= 0) {
          const block = buf.slice(0, k); buf = buf.slice(k + 2);
          const ev = (block.match(/^event: (.*)$/m) || [])[1], data = (block.match(/^data: (.*)$/m) || [])[1];
          let d = {};
          try { d = JSON.parse(data || '{}'); } catch (e) { continue; }
          if (ev === 'text') ans.text += d.t || '';
          else if (ev === 'done') { ans.text = d.text || ans.text; ans.meta = `${d.tokens_in || 0} → ${d.tokens_out || 0} tokens`; }
          else if (ev === 'meta') ans.metaPre = `${d.provider} · ${d.model}`;
          else if (ev === 'error') ans.error = d.error;
          renderChat();
        }
      }
    } catch (e) { if (e.name !== 'AbortError') ans.error = e.message; }
    busy = null;
    ans.actions = ans.error ? [] : parseActions(ans.text);
    if (ans.metaPre) ans.meta = ans.metaPre + (ans.meta ? ' · ' + ans.meta : '');
    if ($('hp-send')) $('hp-send').textContent = t('help_send');
    renderChat();
  }

  // an error toast's ? opens Help with the error as context
  function askError(msg) {
    ctx = {kind: 'error', title: t('help_err_ask'), path: '', text: String(msg || '').slice(0, 2000)};
    open().then(() => { show('ask'); setTimeout(() => { const q = $('hp-ask-q'); if (q && !q.value) { q.value = t('help_err_ask'); q.focus(); } }, 300); });
  }

  // ── one-page AI setup checklist ──
  async function aiSetup() {
    let s;
    try { s = await loadStatus(); } catch (e) { showToast(e.message, 'error'); return; }
    const admin = s.is_admin;
    const light = (ok, off) => `<span class="ais-light ${ok ? 'ok' : off ? 'off' : 'todo'}">${ok ? '●' : '○'}</span>`;
    const pol = v => ({off: t('ais_off'), admins: t('ais_s1_admin'), all: t('ais_s1_all'), assistant: t('adm_ai')}[v] || v);
    const step = (ok, title, desc, btns, why) => `<div class="ais-step">${light(ok)}<div><b>${esc(title)}</b> <span class="pill">${esc(ok ? t('ais_done') : t('ais_todo'))}</span>
      <div class="ai-see">${esc(desc)}</div>${!ok && why ? `<div class="ai-see">${esc(tf('ais_reason', {r: why}))}</div>` : ''}<div class="ais-btns">${btns || ''}</div></div></div>`;
    const onAI = s.assistant.policy !== 'off' && !s.assistant.kill_switch;
    const html = `<div class="modal" style="width:min(620px,calc(100vw - 24px))" id="ai-setup"><div style="display:flex;align-items:center;gap:10px;"><h3 style="margin:0;flex:1;">🧭 ${esc(t('ais_title'))}</h3>
      <button class="icon-btn" onclick="wrmHelp.closeSetup()">✕</button></div><div class="ai-see" style="margin:6px 0 10px;">${esc(t('ais_intro'))}</div>
      ${s.assistant.kill_switch ? `<div class="hint-box warn">${esc(t('ais_kill'))}</div>` : ''}
      ${step(onAI, t('ais_s1'), tf('ais_s1_d', {v: pol(s.assistant.policy)}), admin && !s.assistant.locked
        ? `<button class="btn-sec btn-sm" onclick="wrmHelp.setupSet('ai_assistant','admins')">${esc(t('ais_s1_admin'))}</button><button class="btn-sec btn-sm" onclick="wrmHelp.setupSet('ai_assistant','all')">${esc(t('ais_s1_all'))}</button>`
        : (!onAI ? esc(t('ais_ask_admin')) : ''), s.assistant.allowed ? '' : s.assistant.reason)}
      ${step(s.providers.length > 0, t('ais_s2'), tf('ais_s2_d', {o: s.org_providers, u: s.own_providers}),
        (admin ? `<button class="btn-sm" onclick="wrmHelp.setupWizard('org')">${esc(t('ais_s2_btn'))}</button>` : '') + (s.personal_keys ? ` <button class="btn-sec btn-sm" onclick="wrmHelp.setupWizard('user')">${esc(t('ais_s2_own'))}</button>` : '') + (!admin && !s.personal_keys ? esc(t('ais_ask_admin')) : ''))}
      ${step(s.help.allowed, t('ais_s3'), tf('ais_s3_d', {v: pol(s.help.policy)}), admin && !s.help.locked && !s.help.allowed ? `<button class="btn-sec btn-sm" onclick="wrmHelp.setupSet('ai_help','assistant')">${esc(t('adm_ai'))}</button>` : '', s.help.reason)}
      ${step(s.mcp.enabled && s.mcp.connections > 0, t('ais_s4'), tf('ais_s4_d', {n: s.mcp.connections}),
        (admin && !s.mcp.enabled && !s.mcp.locked ? `<button class="btn-sec btn-sm" onclick="wrmHelp.setupSet('ai_mcp_enabled','1')">${esc(t('ais_s4_on'))}</button> ` : '') + (s.mcp.enabled && s.mcp.allowed ? `<button class="btn-sec btn-sm" onclick="wrmHelp.setupClose(()=>openSettings('st-mcp'))">${esc(t('ais_s4_btn'))}</button>` : ''))}
      <div class="modal-footer"><button class="btn-sec" onclick="wrmHelp.setupClose(()=>wrmHelp.open('ai-setup'))">${esc(t('ais_read_more'))}</button><button onclick="wrmHelp.closeSetup()">${esc(t('ais_close'))}</button></div></div>`;
    let ov = $('ai-setup-ov');
    if (!ov) { ov = document.createElement('div'); ov.id = 'ai-setup-ov'; ov.className = 'modal-overlay open'; ov.style.zIndex = 9700; document.body.appendChild(ov); }
    ov.innerHTML = html;
  }
  function closeSetup() { const ov = $('ai-setup-ov'); if (ov) ov.remove(); }
  function setupClose(fn) { closeSetup(); fn(); }
  async function setupSet(key, value) {
    try {
      const r = await fetch('/api/admin/settings', {method: 'PUT', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({[key]: value})});
      if (!r.ok) throw new Error(await apiError(r));
      if (typeof refreshUserPolicy === 'function') await refreshUserPolicy();
    } catch (e) { showToast(e.message, 'error'); }
    aiSetup();
  }
  function setupWizard(scope) {
    closeSetup();
    loadAI().then(() => wrmAI.awOpen(scope), () => showToast(t('ai_load_failed'), 'error'));
  }

  document.addEventListener('keydown', e => {
    if (e.key === 'F1' && !e.ctrlKey && !e.altKey && !e.metaKey) {
      // a terminal keeps F1 when it has the focus
      if (document.activeElement && document.activeElement.closest('.xterm')) return;
      e.preventDefault(); toggle();
    }
  });

  return {open, close, toggle, show, search, pick, send, run, newChat, clearCtx, openFeature, askAbout, askError, aiSetup, closeSetup, setupClose, setupSet, setupWizard, parseActions, mdHTML};
})();
window.wrmHelp = wrmHelp;
function openHelp(id) { wrmHelp.open(id); }
function openHelpAISetup() { wrmHelp.aiSetup(); }
