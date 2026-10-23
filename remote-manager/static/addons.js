// Addons of WRM PRO (v13.0): Settings → Addons. An addon is a folder under
// <DB>.addons with addon.json that adds a panel and/or a colour theme. A panel
// runs in a sandboxed iframe (sandbox="allow-scripts", no same-origin), so it
// cannot read the page, cookies or other panels; it only gets the small context
// sent through postMessage. Themes are a stylesheet the administrator enabled.
// The server decides everything (manifest, capabilities, who may install): this
// file only shows it.
'use strict';

Object.assign(LANGS.en, {
  tab_addons: 'Addons',
  stin_addons: 'Addons extend WRM from inside itself: an extra panel or a colour theme, installed as a ZIP. A panel runs sandboxed and cannot touch the rest of the interface. Only an administrator installs, enables or removes an addon.',
  add_intro: 'An addon is a ZIP with addon.json. It can add a panel (a small sandboxed window) and a colour theme. Install, enable and remove are administrator actions.',
  add_disclaimer: 'Addons are contributed by the community and are used at your own risk. WRM PRO is a hobby project without a warranty: check what an addon does before you enable it, and prefer addons you trust.',
  add_none: 'No addons are installed.',
  add_install: 'Install an addon (.zip)',
  add_install_d: 'The archive is checked before anything is written: a manifest, a safe layout and declared capabilities. Server code runs only as WebAssembly inside a sandbox.',
  add_pick: 'Choose the .zip file',
  add_go: 'Install',
  add_need_file: 'Choose the .zip file first',
  add_too_big: 'The archive is larger than 20 MB',
  add_installed: 'Addon {name} is installed',
  add_open: 'Open panel',
  add_enable: 'Enable',
  add_disable: 'Disable',
  add_remove: 'Remove',
  add_remove_q: 'Remove the addon {name}? Its files are deleted.',
  add_enabled_toast: 'Addon {name} is enabled',
  add_disabled_toast: 'Addon {name} is disabled',
  add_removed: 'Addon {name} is removed',
  add_by: 'by {author}',
  add_caps: 'Capabilities',
  add_run: 'Run',
  add_run_d: 'Server code runs in a WebAssembly sandbox: it has no files, no network and no environment, and a call is stopped after five seconds. Its output is shown below.',
  add_run_in: 'Input (JSON)',
  add_run_out: 'Result',
  add_run_running: 'Running…',
  add_run_empty: '(no output)',
  add_ai: 'AI addon',
  add_ai_d: 'Describe the addon you want to an AI client: WRM builds the prompt with the manifest, the capabilities this build accepts and the sandbox rules.',
  add_ai_q: 'Create an addon with AI',
  add_ai_help: 'Copy this prompt into the AI client you use (Claude, Codex, OpenCode…) or open the WRM assistant. With the addon tools the assistant can write a draft directly into WRM; otherwise zip the files it returns and install the ZIP above. A draft is never served or executed and never enables itself: you review the diff and install it here. Only what the manifest declares is executed, and server code only as WebAssembly in the sandbox.',
  add_ai_copy: 'Copy prompt',
  add_ai_copied: 'The prompt is copied',
  add_ai_copy_fail: 'The clipboard is not available; the text is selected, copy it with Ctrl+C.',
  add_ai_open: 'Open assistant',
  add_ai_close: 'Close',
  add_community: 'Community library',
  add_community_d: 'Addons published in the public repository "WRM Pro - community public addons". The badge is the trust level recorded in the index; every download must match its SHA-256 before anything is installed.',
  add_community_none: 'The community index is empty or cannot be reached.',
  add_community_install: 'Install from community',
  add_trust_official: 'Official',
  add_trust_reviewed: 'Reviewed',
  add_trust_community: 'Community',
  add_trust_unverified: 'Unverified',
  add_drafts: 'AI drafts',
  add_drafts_d: 'Addons the assistant wrote, kept apart from the installed ones. Review the diff, then install; a new draft is installed disabled and you enable it yourself.',
  add_drafts_none: 'No drafts.',
  add_draft_files: '{n} files',
  add_draft_installed: 'installed {version}',
  add_draft_not_installed: 'not installed',
  add_draft_diff: 'Diff',
  add_draft_install: 'Install draft',
  add_draft_discard: 'Discard',
  add_draft_install_q: 'Install the draft {name}? It goes through the same checks as a ZIP; a new addon stays disabled until you enable it.',
  add_draft_installed_toast: 'Draft {name} is installed; enable it when you are ready',
  add_draft_discarded: 'Draft {name} is discarded',
  add_draft_diff_title: 'Draft {name} compared with the installed addon',
  add_draft_invalid: 'This draft is not valid: {error}',
  add_draft_empty: '(no changes)',
});

Object.assign(LANGS.hr, {
  tab_addons: 'Dodaci',
  stin_addons: 'Dodaci proširuju WRM iz njega samoga: dodatni panel ili tema boja, instalirani kao ZIP. Panel radi u sandboxu i ne može dirati ostatak sučelja. Samo administrator instalira, uključuje ili uklanja dodatak.',
  add_intro: 'Dodatak je ZIP s addon.json. Može dodati panel (mali sandbox prozor) i temu boja. Instalacija, uključivanje i uklanjanje su administratorske radnje.',
  add_disclaimer: 'Dodatke piše zajednica i koristite ih na vlastitu odgovornost. WRM PRO je hobi projekt bez jamstva: provjerite što dodatak radi prije uključivanja i koristite dodatke kojima vjerujete.',
  add_none: 'Nema instaliranih dodataka.',
  add_install: 'Instaliraj dodatak (.zip)',
  add_install_d: 'Arhiva se provjerava prije nego se išta zapiše: manifest, siguran raspored i deklarirane mogućnosti. Kôd na serveru izvodi se samo kao WebAssembly u sandboxu.',
  add_pick: 'Odaberite .zip datoteku',
  add_go: 'Instaliraj',
  add_need_file: 'Prvo odaberite .zip datoteku',
  add_too_big: 'Arhiva je veća od 20 MB',
  add_installed: 'Dodatak {name} je instaliran',
  add_open: 'Otvori panel',
  add_enable: 'Uključi',
  add_disable: 'Isključi',
  add_remove: 'Ukloni',
  add_remove_q: 'Ukloniti dodatak {name}? Njegove datoteke se brišu.',
  add_enabled_toast: 'Dodatak {name} je uključen',
  add_disabled_toast: 'Dodatak {name} je isključen',
  add_removed: 'Dodatak {name} je uklonjen',
  add_by: 'autor {author}',
  add_caps: 'Mogućnosti',
  add_run: 'Pokreni',
  add_run_d: 'Kôd na serveru radi u WebAssembly sandboxu: nema datoteka, mreže ni okoline, a poziv se prekida nakon pet sekundi. Rezultat se prikazuje ispod.',
  add_run_in: 'Ulaz (JSON)',
  add_run_out: 'Rezultat',
  add_run_running: 'Radi…',
  add_run_empty: '(nema izlaza)',
  add_ai: 'AI dodatak',
  add_ai_d: 'Opišite željeni dodatak AI klijentu: WRM sastavlja upit s manifestom, mogućnostima koje ova verzija prihvaća i pravilima sandboxa.',
  add_ai_q: 'Napravi dodatak uz AI',
  add_ai_help: 'Kopirajte ovaj upit u AI klijent koji koristite (Claude, Codex, OpenCode…) ili otvorite WRM asistenta. S addon alatima asistent može izravno napisati skicu u WRM; inače zapakirajte datoteke koje AI vrati u ZIP i instalirajte ga gore. Skica se nikad ne poslužuje ni izvršava i nikad se sama ne uključuje: pregledate razlike i instalirate je ovdje. Izvršava se samo ono što manifest deklarira, a kôd na serveru samo kao WebAssembly u sandboxu.',
  add_ai_copy: 'Kopiraj upit',
  add_ai_copied: 'Upit je kopiran',
  add_ai_copy_fail: 'Međuspremnik nije dostupan; tekst je označen, kopirajte ga s Ctrl+C.',
  add_ai_open: 'Otvori asistenta',
  add_ai_close: 'Zatvori',
  add_community: 'Zajednica',
  add_community_d: 'Dodaci objavljeni u javnom repozitoriju "WRM Pro - community public addons". Značka je razina povjerenja iz indeksa; svako preuzimanje mora odgovarati svom SHA-256 prije nego se išta instalira.',
  add_community_none: 'Indeks zajednice je prazan ili nedostupan.',
  add_community_install: 'Instaliraj iz zajednice',
  add_trust_official: 'Službeno',
  add_trust_reviewed: 'Pregledano',
  add_trust_community: 'Zajednica',
  add_trust_unverified: 'Neprovjereno',
  add_drafts: 'AI skice',
  add_drafts_d: 'Dodaci koje je napisao asistent, odvojeni od instaliranih. Pregledajte razlike i instalirajte; nova skica instalira se isključena i sami je uključujete.',
  add_drafts_none: 'Nema skica.',
  add_draft_files: '{n} datoteka',
  add_draft_installed: 'instalirano {version}',
  add_draft_not_installed: 'nije instalirano',
  add_draft_diff: 'Razlike',
  add_draft_install: 'Instaliraj skicu',
  add_draft_discard: 'Odbaci',
  add_draft_install_q: 'Instalirati skicu {name}? Prolazi iste provjere kao ZIP; novi dodatak ostaje isključen dok ga sami ne uključite.',
  add_draft_installed_toast: 'Skica {name} je instalirana; uključite je kad ste spremni',
  add_draft_discarded: 'Skica {name} je odbačena',
  add_draft_diff_title: 'Skica {name} u odnosu na instalirani dodatak',
  add_draft_invalid: 'Skica nije valjana: {error}',
  add_draft_empty: '(nema promjena)',
});

(function () {
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(String(b)); return s; };
  let list = [], caps = {}, loaded = false, loading = null;
  let community = [], trust = {}, communityError = '';
  let drafts = [], draftsError = '';
  const openWins = {};

  async function api(method, url, body) {
    const opts = {method, headers: {}};
    if (body !== undefined) { opts.headers['Content-Type'] = 'application/json'; opts.body = JSON.stringify(body); }
    const res = await fetch(url, opts);
    if (!res.ok) throw new Error(await apiError(res));
    return res.json().catch(() => ({}));
  }

  async function refresh() {
    if (loading) return loading;
    loading = api('GET', '/api/addons').then(d => {
      list = d.addons || [];
      caps = d.capabilities || {};
      loaded = true;
    }).finally(() => { loading = null; });
    return loading;
  }

  // ── community library ────────────────────────────────
  async function loadCommunity() {
    try {
      const d = await api('GET', '/api/addons/community');
      community = d.addons || [];
      trust = d.trust || {};
      communityError = '';
    } catch (e) {
      community = [];
      communityError = e.message;
    }
  }

  // ── AI drafts ────────────────────────────────────────
  // The assistant writes files into <DB>.addons/.drafts/<id>/. They are not
  // installed, served or executed until the administrator installs one here.
  async function loadDrafts() {
    const admin = !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
    if (!admin) { drafts = []; draftsError = ''; return; }
    try {
      const d = await api('GET', '/api/addons/drafts');
      drafts = d.drafts || [];
      draftsError = '';
    } catch (e) {
      drafts = [];
      draftsError = e.message;
    }
  }

  function draftsHTML() {
    const admin = !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
    if (!admin) return '';
    const rows = drafts.map(d => {
      const installed = d.installed ? tf('add_draft_installed', {version: d.installed_version || ''}) + (d.installed_enabled ? '' : ' · ' + t('add_disable')) : t('add_draft_not_installed');
      const invalid = d.valid ? '' : `<div class="hint-box warn" style="margin-top:6px;">${esc(tf('add_draft_invalid', {error: d.error || ''}))}</div>`;
      return `<div style="border:1px solid var(--border);border-radius:8px;padding:10px;margin:8px 0;${d.valid ? '' : 'opacity:.75;'}">
        <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
          <b>${esc(d.name || d.id)}</b> ${d.version ? pill(d.version) : ''}
          ${(d.capabilities || []).map(c => pill(c)).join(' ')}
          <span style="flex:1;"></span>
          <button class="btn-sec btn-sm" onclick="wrmAddons.draftDiff('${esc(d.id)}')">${esc(t('add_draft_diff'))}</button>
          <button class="btn-sec btn-sm" ${d.valid ? '' : 'disabled'} onclick="wrmAddons.draftInstall('${esc(d.id)}')">${esc(t('add_draft_install'))}</button>
          <button class="btn-danger btn-sm" onclick="wrmAddons.draftDiscard('${esc(d.id)}')">${esc(t('add_draft_discard'))}</button>
        </div>
        ${d.description ? `<div class="sub" style="font-size:11.5px;color:var(--text2);margin-top:4px;">${esc(d.description)}</div>` : ''}
        <div class="sub" style="font-size:10.5px;color:var(--text3);margin-top:2px;">${esc(d.id)} · ${esc(tf('add_draft_files', {n: (d.files || []).length}))} · ${esc(installed)}${d.updated_at ? ' · ' + esc(fmtTime(d.updated_at)) : ''}</div>
        ${invalid}
      </div>`;
    }).join('');
    return `<div class="adm-h">${esc(t('add_drafts'))}</div>
      <div class="hint-box">${esc(t('add_drafts_d'))}</div>
      ${rows || `<div class="sub" style="color:var(--text3);font-size:11.5px;">${esc(draftsError || t('add_drafts_none'))}</div>`}`;
  }

  async function draftInstall(id) {
    const d = drafts.find(x => x.id === id);
    if (!d || !d.valid) return;
    if (!(await uiConfirm(tf('add_draft_install_q', {name: d.name || d.id}), {okText: t('add_draft_install')}))) return;
    try {
      await api('POST', '/api/addons/drafts/' + encodeURIComponent(id) + '/install');
      await after(tf('add_draft_installed_toast', {name: d.name || id}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    }
  }

  async function draftDiff(id) {
    let d;
    try {
      d = await api('GET', '/api/addons/drafts/' + encodeURIComponent(id) + '/diff');
    } catch (e) {
      showToast(e.message, 'error', 8000);
      return;
    }
    const draft = drafts.find(x => x.id === id);
    let m = $('addon-draft-diff');
    if (!m) {
      m = document.createElement('div');
      m.id = 'addon-draft-diff';
      m.className = 'modal-overlay';
      m.style.zIndex = 9660;
      document.body.appendChild(m);
    }
    m.innerHTML = `<div class="modal" style="width:900px;max-width:calc(100vw - 32px);max-height:calc(100vh - 32px);overflow:auto;">
      <h3>${esc(tf('add_draft_diff_title', {name: (draft && (draft.name || draft.id)) || id}))}</h3>
      <pre style="margin:0;padding:8px;background:var(--bg2);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;white-space:pre-wrap;max-height:60vh;overflow:auto;color:var(--text2);">${esc(d.diff || t('add_draft_empty'))}</pre>
      <div style="display:flex;justify-content:flex-end;margin-top:10px;"><button onclick="wrmAddons.closeDraftDiff()">${esc(t('add_ai_close'))}</button></div>
    </div>`;
    m.classList.add('open');
  }

  function closeDraftDiff() {
    const m = $('addon-draft-diff');
    if (m) m.classList.remove('open');
  }

  async function draftDiscard(id) {
    const d = drafts.find(x => x.id === id);
    if (!(await uiConfirm(tf('add_draft_discard', {name: (d && (d.name || d.id)) || id}), {okText: t('add_draft_discard'), danger: true}))) return;
    try {
      await api('DELETE', '/api/addons/drafts/' + encodeURIComponent(id));
      await after(tf('add_draft_discarded', {name: (d && (d.name || d.id)) || id}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    }
  }

  function communityHTML() {
    const admin = !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
    return `<div class="adm-h">${esc(t('add_community'))}</div>
      <div class="hint-box warn">${esc(t('add_community_d'))}</div>
      ${community.length ? community.map(a => `<div style="border:1px solid var(--border);border-radius:8px;padding:10px;margin:8px 0;">
        <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
          <b>${esc(a.name)}</b> ${pill(a.version)}
          <span class="pill ${a.trust === 'official' ? 'ok' : a.trust === 'unverified' ? 'warn' : ''}">${esc(t('add_trust_' + a.trust))}</span>
          ${(a.capabilities || []).map(c => pill(c)).join(' ')}
          <span style="flex:1;"></span>
          ${admin ? `<button class="btn-sec btn-sm" onclick="wrmAddons.communityInstall('${esc(a.id)}')">${esc(t('add_community_install'))}</button>` : ''}
        </div>
        ${a.description ? `<div class="sub" style="font-size:11.5px;color:var(--text2);margin-top:4px;">${esc(a.description)}</div>` : ''}
        <div class="sub" style="font-size:10.5px;color:var(--text3);margin-top:2px;">${esc(a.id)} · ${esc(a.license || '')} · ${esc(trust[a.trust] || '')}</div>
      </div>`).join('') : `<div class="sub" style="color:var(--text3);font-size:11.5px;">${esc(communityError || t('add_community_none'))}</div>`}`;
  }

  async function communityInstall(id) {
    const a = community.find(x => x.id === id);
    try {
      await api('POST', '/api/addons/community/install', {id});
      await loadCommunity();
      await after(tf('add_installed', {name: a ? a.name : id}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    }
  }

  // ── themes ───────────────────────────────────────────
  // An enabled theme is one stylesheet link. Removing or disabling the addon
  // takes the link out again, so the look always matches the enabled list.
  function applyThemes() {
    const want = {};
    list.filter(a => a.enabled && a.ui && a.ui.theme).forEach(a => { want['addon-theme-' + a.id] = assetURL(a, a.ui.theme); });
    document.querySelectorAll('link[data-addon-theme]').forEach(l => {
      if (!want[l.id]) l.remove();
    });
    Object.entries(want).forEach(([id, href]) => {
      let l = $(id);
      if (!l) {
        l = document.createElement('link');
        l.id = id; l.rel = 'stylesheet'; l.dataset.addonTheme = '1';
        document.head.appendChild(l);
      }
      if (l.getAttribute('href') !== href) l.href = href;
    });
  }

  function assetURL(a, file) {
    return '/api/addons/' + encodeURIComponent(a.id) + '/asset/' + String(file).split('/').map(encodeURIComponent).join('/');
  }

  // ── the panel window ─────────────────────────────────
  function themeTokens() {
    const cs = getComputedStyle(document.documentElement), out = {};
    ['--bg', '--bg2', '--bg3', '--bg4', '--bg5', '--border', '--border2', '--text', '--text2', '--text3',
      '--accent', '--accent-h', '--green', '--red'].forEach(k => { out[k] = cs.getPropertyValue(k).trim(); });
    return out;
  }

  function sendContext(frame, a) {
    let conn = [];
    try {
      conn = (typeof conns !== 'undefined' ? conns : []).map(c => ({id: c.id, name: c.name, host: c.host, folder: c.folder_id || null}));
    } catch (_) {}
    const me = (typeof currentUser !== 'undefined' && currentUser) ? currentUser : {};
    const ctx = {type: 'wrm:context', addon: a.id, user: me.username || '', lang: typeof currentLang !== 'undefined' ? currentLang : 'en', theme: themeTokens(), connections: conn};
    try { frame.contentWindow.postMessage(ctx, '*'); } catch (_) {}
  }

  function openPanel(id) {
    const a = list.find(x => x.id === id);
    if (!a || !a.enabled || !a.ui || !a.ui.panel) return;
    const have = openWins[id];
    if (have && typeof winObjects !== 'undefined' && winObjects.some(w => w.id === have.id && w.el)) {
      if (have.el.style.display === 'none') have.el.style.display = 'flex';
      focusWindow(have.id);
      return;
    }
    if (typeof addWindow !== 'function') return;
    const wo = addWindow({mode: 'addon-panel'});
    wo.addonPanel = id;
    wo.addonTitle = a.ui.title || a.name;
    openWins[id] = wo;
    const tt = $('wt-' + wo.id);
    if (tt) tt.textContent = '▣ ' + wo.addonTitle;
    const wc = $('wc-' + wo.id);
    wc.innerHTML = '';
    wc.style.padding = '0';
    wc.style.overflow = 'hidden';
    const frame = document.createElement('iframe');
    frame.className = 'addon-frame';
    frame.setAttribute('sandbox', 'allow-scripts');
    frame.setAttribute('referrerpolicy', 'no-referrer');
    frame.style.cssText = 'width:100%;height:100%;border:0;display:block;background:#fff';
    // The context goes out once, to the panel page WRM loaded. If the frame later
    // navigates (a link in the panel), the new page may be another site: it gets
    // nothing unless it is an addon page that asks (see the message listener).
    frame.addEventListener('load', () => { if (frame.dataset.ctxSent) return; frame.dataset.ctxSent = '1'; sendContext(frame, a); });
    frame.src = assetURL(a, a.ui.panel);
    wc.appendChild(frame);
    if (typeof renderTabs === 'function') renderTabs();
  }

  // The only words a panel gets back into WRM: a toast, or a request for the
  // context again. Anything else is ignored.
  window.addEventListener('message', ev => {
    const d = ev.data;
    if (!d || typeof d !== 'object' || typeof d.type !== 'string') return;
    // Addon pages are served with the sandbox CSP, so their origin is opaque
    // ("null"); a page from any real origin in the frame is not an addon page.
    if (ev.origin !== 'null') return;
    const frame = Array.from(document.querySelectorAll('iframe.addon-frame')).find(f => f.contentWindow === ev.source);
    if (!frame) return;
    const a = list.find(x => { const wo = openWins[x.id]; return wo && wo.el && wo.el.contains(frame); });
    if (d.type === 'wrm:notify' && typeof d.text === 'string') {
      showToast(String(d.text).slice(0, 500), d.kind === 'error' ? 'error' : 'success', 6000);
    } else if (d.type === 'wrm:context' && a) {
      sendContext(frame, a);
    }
  });

  // ── settings pane ────────────────────────────────────
  function pill(text) { return `<span class="pill">${esc(text)}</span>`; }

  function rowHTML(a) {
    const admin = !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
    const panel = a.ui && a.ui.panel;
    const btns = [];
    if (panel && a.enabled) btns.push(`<button class="btn-sec btn-sm" onclick="wrmAddons.openPanel('${esc(a.id)}')">${esc(t('add_open'))}</button>`);
    if (admin) {
      btns.push(`<label class="mcp-check" style="margin:0;display:flex;align-items:center;gap:6px;"><input type="checkbox" ${a.enabled ? 'checked' : ''} onchange="wrmAddons.toggle('${esc(a.id)}', this.checked)"><span>${esc(a.enabled ? t('add_disable') : t('add_enable'))}</span></label>`);
      btns.push(`<button class="btn-danger btn-sm" onclick="wrmAddons.remove('${esc(a.id)}')">${esc(t('add_remove'))}</button>`);
    }
    const meta = [a.author ? tf('add_by', {author: a.author}) : '', a.license || '', a.homepage || ''].filter(Boolean).map(esc).join(' · ');
    const exports = (a.server && a.server.exports) || [];
    const runner = admin && a.enabled && exports.length ? `<details style="margin-top:6px;">
        <summary style="cursor:pointer;font-size:11.5px;color:var(--text2);">${esc(t('add_run'))} · ${esc(a.server.module)}</summary>
        <div class="setting-desc" style="margin:4px 0;">${esc(t('add_run_d'))}</div>
        <div style="display:flex;gap:6px;align-items:flex-start;flex-wrap:wrap;">
          <select id="addon-export-${esc(a.id)}" style="background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font-size:11.5px;padding:4px 6px;">${exports.map(x => `<option value="${esc(x)}">${esc(x)}</option>`).join('')}</select>
          <textarea id="addon-in-${esc(a.id)}" placeholder="${esc(t('add_run_in'))}" style="flex:1;min-width:180px;min-height:44px;background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;padding:5px 7px;"></textarea>
          <button class="btn-sec btn-sm" onclick="wrmAddons.run('${esc(a.id)}')">${esc(t('add_run'))}</button>
        </div>
        <pre id="addon-out-${esc(a.id)}" style="margin:6px 0 0;padding:6px 8px;background:var(--bg2);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;white-space:pre-wrap;max-height:180px;overflow:auto;color:var(--text2);">${esc(t('add_run_out'))}</pre>
      </details>` : '';
    return `<div style="border:1px solid var(--border);border-radius:8px;padding:10px;margin:8px 0;${a.enabled ? '' : 'opacity:.65;'}">
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        <b>${esc(a.name)}</b> ${pill(a.version)}
        ${(a.capabilities || []).map(c => pill(c)).join(' ')}
        <span style="flex:1;"></span>${btns.join(' ')}
      </div>
      ${a.description ? `<div class="sub" style="font-size:11.5px;color:var(--text2);margin-top:4px;">${esc(a.description)}</div>` : ''}
      ${meta ? `<div class="sub" style="font-size:11px;color:var(--text3);margin-top:2px;">${meta}</div>` : ''}
      <div class="sub" style="font-size:10.5px;color:var(--text3);margin-top:2px;">${esc(a.id)} · ${esc(a.checksum || '')}${a.installed_at ? ' · ' + esc(fmtTime(a.installed_at)) : ''}</div>
      ${runner}
    </div>`;
  }

  function render() {
    const el = $('addons-pane');
    if (!el) return;
    const admin = !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
    const install = admin ? `<div class="setting-row">
      <div><div class="setting-label">${esc(t('add_install'))}</div><div class="setting-desc">${esc(t('add_install_d'))}</div></div>
      <div class="setting-control" style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        <input type="file" id="addon-file" accept=".zip,application/zip" style="max-width:240px;">
        <button class="btn-sec btn-sm" id="addon-install-btn" onclick="wrmAddons.install()">${esc(t('add_go'))}</button>
      </div></div>` : '';
    const aiBox = admin ? `<div class="setting-row">
      <div><div class="setting-label">✨ ${esc(t('add_ai_q'))}</div><div class="setting-desc">${esc(t('add_ai_d'))}</div></div>
      <div class="setting-control"><button class="btn-sec btn-sm" onclick="wrmAddons.aiPrompt()">${esc(t('add_ai'))}</button></div></div>` : '';
    el.innerHTML = `<div class="hint-box">${esc(t('add_intro'))}</div>
      <div class="hint-box warn">${esc(t('add_disclaimer'))}</div>
      ${install}
      ${aiBox}
      ${draftsHTML()}
      <div class="adm-h">${esc(t('tab_addons'))}</div>
      ${list.length ? list.map(rowHTML).join('') : `<div class="sub" style="color:var(--text3);font-size:11.5px;">${esc(t('add_none'))}</div>`}
      ${communityHTML()}`;
  }

  async function pane() {
    const el = $('addons-pane');
    if (!el) return;
    el.innerHTML = `<div class="sub" style="color:var(--text3);">…</div>`;
    try { await refresh(); } catch (e) { el.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    await loadCommunity();
    await loadDrafts();
    render();
  }

  async function after(msg) {
    if (msg) showToast(msg, 'success');
    try { await refresh(); } catch (_) {}
    await loadDrafts();
    applyThemes();
    if ($('addons-pane') && $('st-addons') && $('st-addons').classList.contains('active')) render();
  }

  async function install() {
    const inp = $('addon-file'), btn = $('addon-install-btn');
    const file = inp && inp.files && inp.files[0];
    if (!file) { showToast(t('add_need_file'), 'error'); return; }
    if (file.size > 20 * 1024 * 1024) { showToast(t('add_too_big'), 'error'); return; }
    if (btn) btn.disabled = true;
    try {
      const res = await fetch('/api/addons/install', {method: 'POST', headers: {'Content-Type': 'application/zip'}, body: file});
      if (!res.ok) throw new Error(await apiError(res));
      const d = await res.json().catch(() => ({}));
      await after(tf('add_installed', {name: d.name || d.id || file.name}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    } finally {
      if (btn) btn.disabled = false;
    }
  }

  async function toggle(id, on) {
    const a = list.find(x => x.id === id);
    try {
      await api('POST', '/api/addons/' + encodeURIComponent(id) + '/' + (on ? 'enable' : 'disable'));
      await after(tf(on ? 'add_enabled_toast' : 'add_disabled_toast', {name: a ? a.name : id}));
      if (!on) {
        const wo = openWins[id];
        if (wo) { delete openWins[id]; try { if (wo.el) closeWindow(wo.id); } catch (_) {} }
      }
    } catch (e) {
      showToast(e.message, 'error', 8000);
      render();
    }
  }

  async function remove(id) {
    const a = list.find(x => x.id === id);
    if (!(await uiConfirm(tf('add_remove_q', {name: a ? a.name : id}), {okText: t('add_remove'), danger: true}))) return;
    try {
      await api('DELETE', '/api/addons/' + encodeURIComponent(id));
      const wo = openWins[id];
      if (wo) { delete openWins[id]; try { if (wo.el) closeWindow(wo.id); } catch (_) {} }
      await after(tf('add_removed', {name: a ? a.name : id}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    }
  }

  // ── server code ──────────────────────────────────────
  // The input box is JSON text sent as-is; the result is shown in the card.
  async function run(id) {
    const a = list.find(x => x.id === id);
    if (!a || !a.server) return;
    const sel = $('addon-export-' + id), inp = $('addon-in-' + id), out = $('addon-out-' + id);
    const ex = (sel && sel.value) || (a.server.exports || [])[0];
    const body = (inp && inp.value.trim()) || '{}';
    if (out) out.textContent = t('add_run_running');
    try {
      const res = await fetch('/api/addons/' + encodeURIComponent(id) + '/call/' + encodeURIComponent(ex), {
        method: 'POST', headers: {'Content-Type': 'application/json'}, body
      });
      const d = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(d.error || ('HTTP ' + res.status));
      if (out) out.textContent = d.output ? d.output : t('add_run_empty');
    } catch (e) {
      if (out) out.textContent = e.message;
      showToast(e.message, 'error', 8000);
    }
  }

  // ── AI addon ─────────────────────────────────────────
  // The prompt is generated from the capabilities the server reports, so it
  // cannot describe a hook this build would refuse. WRM does not launch other
  // apps: the text goes to the user's AI client, or to the built-in assistant.
  function aiPromptText() {
    const capLines = Object.entries(caps || {}).map(([k, v]) => '- ' + k + ': ' + v).join('\n');
    return [
      'You are writing an addon for WRM PRO, a self-hosted server manager.',
      'Deliver it as ONE zip with addon.json at the root plus the files it declares.',
      '',
      'addon.json:',
      '{',
      '  "id": "lowercase.dotted-id",',
      '  "name": "Short name",',
      '  "version": "1.0.0",',
      '  "description": "One sentence",',
      '  "author": "Your name",',
      '  "license": "MIT",',
      '  "homepage": "https://example.com/addon",',
      '  "capabilities": ["ui.panel"],',
      '  "ui": {"panel": "panel.html", "theme": "theme.css", "title": "Window title"},',
      '  "server": {"module": "plugin.wasm", "exports": ["status"]}',
      '}',
      '',
      'Capabilities this WRM build accepts:',
      capLines,
      '',
      'Rules:',
      '- ui.panel is HTML opened in an iframe with sandbox="allow-scripts" and no same-origin. It cannot read WRM cookies or the page. It receives {type:"wrm:context", user, lang, theme, connections} through postMessage, may send {type:"wrm:notify", text} back, and may ask for the context again with {type:"wrm:context"} (a page it navigates to must ask; only the first page gets it automatically).',
      '- ui.theme is a CSS file applied while the addon is enabled.',
      '- server.wasm is WebAssembly run by wazero with no WASI: no files, no network, no environment, no clock. The only import is wrm.log(ptr, len). The module must export its memory, wrm_alloc(size) -> ptr, and each declared export as (ptr, len) -> ptr returning a pointer to a NUL-terminated UTF-8 string. WRM writes the JSON request at the pointer wrm_alloc returned. 64 KiB input, 1 MB output, 5 seconds per call, 64 MiB memory.',
      '- Nothing runs at install time; the module is compiled and its imports and exports are checked.',
      '',
      'Write the smallest addon that does the job, state its license, and tell me how to zip it for install.',
      'If you are the WRM assistant and have the addon tools, call addon_create with addon.json and the files it declares instead of returning a ZIP; the draft waits for my review in Settings → Addons.'
    ].join('\n');
  }

  function aiPrompt() {
    let m = $('addon-ai-modal');
    if (!m) {
      m = document.createElement('div');
      m.id = 'addon-ai-modal';
      m.className = 'modal-overlay';
      m.style.zIndex = 9660;
      document.body.appendChild(m);
    }
    m.innerHTML = `<div class="modal" style="width:760px;max-width:calc(100vw - 32px);max-height:calc(100vh - 32px);overflow:auto;">
      <h3>✨ ${esc(t('add_ai_q'))}</h3>
      <div class="hint-box">${esc(t('add_ai_help'))}</div>
      <textarea id="addon-ai-text" spellcheck="false" style="width:100%;min-height:320px;background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;padding:8px;">${esc(aiPromptText())}</textarea>
      <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:10px;flex-wrap:wrap;">
        <button class="btn-sec" onclick="wrmAddons.copyPrompt()">${esc(t('add_ai_copy'))}</button>
        <button class="btn-sec" onclick="wrmAddons.openAssistant()">${esc(t('add_ai_open'))}</button>
        <button onclick="wrmAddons.closeAIPrompt()">${esc(t('add_ai_close'))}</button>
      </div>
    </div>`;
    m.classList.add('open');
  }

  function closeAIPrompt() {
    const m = $('addon-ai-modal');
    if (m) m.classList.remove('open');
  }

  async function copyPrompt() {
    const ta = $('addon-ai-text');
    if (!ta) return;
    try {
      await navigator.clipboard.writeText(ta.value);
      showToast(t('add_ai_copied'), 'success');
    } catch (_) {
      ta.focus();
      ta.select();
      showToast(t('add_ai_copy_fail'), 'warning', 7000);
    }
  }

  function openAssistant() {
    closeAIPrompt();
    if (typeof openAIHome === 'function') openAIHome();
    else showToast(t('add_ai_copy_fail'), 'warning', 7000);
  }

  // Themes are applied once the app is signed in: initApp is wrapped here
  // because this script loads after the main one.
  async function boot() {
    try { await refresh(); applyThemes(); } catch (_) {}
    if ($('addons-pane') && $('st-addons') && $('st-addons').classList.contains('active')) render();
  }

  window.wrmAddons = {boot, pane, install, toggle, remove, openPanel, run, aiPrompt, closeAIPrompt, copyPrompt, openAssistant, communityInstall, applyThemes, reload: after,
    draftInstall, draftDiff, draftDiscard, closeDraftDiff};

  if (typeof window.initApp === 'function' && !window.initApp.__addons) {
    const orig = window.initApp;
    window.initApp = function () {
      const r = orig.apply(this, arguments);
      boot();
      return r;
    };
    window.initApp.__addons = true;
  } else if (typeof currentUser !== 'undefined' && currentUser) {
    boot();
  }
})();
