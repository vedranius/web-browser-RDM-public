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
  stin_addons: 'Addons extend WRM from inside itself: an extra panel or a colour theme, installed as a ZIP. A panel runs sandboxed and cannot touch the rest of the interface. An administrator installs addons for everybody; each user may also install personal addons that only they see (policy addon_personal).',
  add_intro: 'An addon is a ZIP with addon.json. It can add a panel (a small sandboxed window) and a colour theme. Administrators install for everybody; personal addons are visible only to the user who installed them.',
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
  stin_addons: 'Dodaci proširuju WRM iz njega samoga: dodatni panel ili tema boja, instalirani kao ZIP. Panel radi u sandboxu i ne može dirati ostatak sučelja. Administrator instalira dodatke za sve; svaki korisnik može instalirati i osobne dodatke koje vidi samo on (pravilo addon_personal).',
  add_intro: 'Dodatak je ZIP s addon.json. Može dodati panel (mali sandbox prozor) i temu boja. Administratori instaliraju za sve; osobni dodatak vidi samo korisnik koji ga je instalirao.',
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

// v13.7: addon platform v2: personal addons, instances, the server-side store,
// remembered windows, the v2 bridge with a watchdog, and the permission sheet.
Object.assign(LANGS.en, {
  add_scope: 'Install for',
  add_scope_org: 'Everybody (organisation)',
  add_scope_me: 'Only me (personal addon)',
  add_install_me: 'Install a personal addon (.zip)',
  add_install_me_d: 'A personal addon is visible only to you. It runs in the same sandbox as any addon and cannot carry server code.',
  add_personal: 'Personal',
  add_personal_of: 'personal · {user}',
  add_personal_others: 'Personal addons of other users',
  add_personal_others_d: 'Each is visible only to its owner. You can turn one off or remove it.',
  add_personal_none: 'No user has a personal addon.',
  add_downgrade_q: '{error}\n\nInstall the older version anyway?',
  add_needs_approval: 'This update asks for more than was approved, so it stays off until you approve it again: {list}',
  add_installed_wait: 'Addon {name} is updated and waits for approval',
  add_approve_t: 'Approve addon {name}',
  add_approve_d: 'Enabling the addon allows exactly this. The list is generated by WRM from addon.json, not written by the author.',
  add_approve_new: 'new in this version',
  add_approve_go: 'Approve and enable',
  add_perm_panel: 'Show a window inside WRM. It runs sandboxed: no access to WRM, its cookies or other windows.',
  add_perm_panel_v2: 'Show a window inside WRM. It runs sandboxed with no network of its own: no access to WRM, its cookies, other windows or any other site.',
  add_perm_themes: 'Offer {n} colour theme(s) in Settings → Appearance; a theme is used only by people who pick it.',
  add_perm_store: 'Keep up to {kb} KB of settings per instance on the WRM server.',
  add_perm_store_shared: 'Keep up to {kb} KB of settings per instance on the WRM server, one set shared by everybody and written by administrators.',
  add_perm_notify: 'Show short notifications (at most 5 a minute, repeats are dropped).',
  add_perm_server: 'Run WebAssembly server code ({exports}) in the sandbox when an administrator calls it.',
  add_perm_other: 'Permission: {p}',
  add_inst: 'Instances',
  add_inst_d: 'Each instance is its own configuration and window (for example two dashboards). Shared instances are set up by an administrator for everybody; yours are visible only to you.',
  add_inst_new: 'New instance',
  add_inst_name_q: 'Name of the new instance',
  add_inst_shared: 'shared',
  add_inst_mine: 'mine',
  add_inst_shared_q: 'Share it with everybody who can use this addon?',
  add_inst_default: 'Default',
  add_inst_rename: 'Rename',
  add_inst_rename_q: 'New name',
  add_inst_delete_q: 'Remove the instance {name}? Its stored settings and window are deleted.',
  add_inst_export: 'Export',
  add_inst_import: 'Import',
  add_inst_imported: 'The configuration is imported ({n} values)',
  add_inst_none: 'No instances yet; opening the panel creates one.',
  add_wd_start: '{name} did not start: the panel sent no ready message within 5 seconds.',
  add_wd_hang: '{name} is not responding (no heartbeat for 45 seconds).',
  add_wd_reload: 'Reload',
  add_wd_diag: 'Diagnostics',
  add_wd_disable: 'Disable addon',
  add_wd_close: 'Close',
  add_diag_t: 'Diagnostics of {name}',
  add_diag_none: 'No errors reported.',
  add_reload_t: 'Reload the panel',
  add_diag_btn_t: 'Errors reported by the panel',
  pg_addons: 'Addons',
  po_addon_personal_off: 'Nobody', po_addon_personal_admins: 'Administrators', po_addon_personal_all: 'Everyone',
  p_addon_personal: 'Personal addons',
  p_addon_personal_d: 'Who may install personal addons, visible only to themselves (Settings → Addons). They run in the same sandbox as organisation addons and cannot carry server code. Off: only administrators install addons, for everybody.',
});

Object.assign(LANGS.hr, {
  add_scope: 'Instaliraj za',
  add_scope_org: 'Sve (organizacija)',
  add_scope_me: 'Samo mene (osobni dodatak)',
  add_install_me: 'Instaliraj osobni dodatak (.zip)',
  add_install_me_d: 'Osobni dodatak vidite samo vi. Radi u istom sandboxu kao svaki dodatak i ne može imati kôd na serveru.',
  add_personal: 'Osobni',
  add_personal_of: 'osobni · {user}',
  add_personal_others: 'Osobni dodaci drugih korisnika',
  add_personal_others_d: 'Svaki vidi samo njegov vlasnik. Možete ga isključiti ili ukloniti.',
  add_personal_none: 'Nijedan korisnik nema osobni dodatak.',
  add_downgrade_q: '{error}\n\nIpak instalirati stariju verziju?',
  add_needs_approval: 'Ovo ažuriranje traži više nego što je odobreno pa ostaje isključeno dok ga ponovno ne odobrite: {list}',
  add_installed_wait: 'Dodatak {name} je ažuriran i čeka odobrenje',
  add_approve_t: 'Odobri dodatak {name}',
  add_approve_d: 'Uključivanjem dodatka dopuštate točno ovo. Popis sastavlja WRM iz addon.json, a ne autor dodatka.',
  add_approve_new: 'novo u ovoj verziji',
  add_approve_go: 'Odobri i uključi',
  add_perm_panel: 'Prikazati prozor unutar WRM-a. Radi u sandboxu: nema pristup WRM-u, njegovim kolačićima ni drugim prozorima.',
  add_perm_panel_v2: 'Prikazati prozor unutar WRM-a. Radi u sandboxu bez vlastite mreže: nema pristup WRM-u, njegovim kolačićima, drugim prozorima ni ijednoj drugoj stranici.',
  add_perm_themes: 'Ponuditi {n} tema boja u Postavke → Izgled; temu koristi samo onaj tko je odabere.',
  add_perm_store: 'Čuvati do {kb} KB postavki po instanci na WRM serveru.',
  add_perm_store_shared: 'Čuvati do {kb} KB postavki po instanci na WRM serveru, jedan zajednički skup za sve koji pišu administratori.',
  add_perm_notify: 'Prikazivati kratke obavijesti (najviše 5 u minuti, ponavljanja se odbacuju).',
  add_perm_server: 'Izvršavati WebAssembly kôd na serveru ({exports}) u sandboxu kad ga administrator pozove.',
  add_perm_other: 'Dozvola: {p}',
  add_inst: 'Instance',
  add_inst_d: 'Svaka instanca ima svoju konfiguraciju i prozor (npr. dvije nadzorne ploče). Zajedničke instance administrator postavlja za sve; vaše vidite samo vi.',
  add_inst_new: 'Nova instanca',
  add_inst_name_q: 'Naziv nove instance',
  add_inst_shared: 'zajednička',
  add_inst_mine: 'moja',
  add_inst_shared_q: 'Podijeliti je sa svima koji mogu koristiti ovaj dodatak?',
  add_inst_default: 'Zadana',
  add_inst_rename: 'Preimenuj',
  add_inst_rename_q: 'Novi naziv',
  add_inst_delete_q: 'Ukloniti instancu {name}? Njezine spremljene postavke i prozor se brišu.',
  add_inst_export: 'Izvezi',
  add_inst_import: 'Uvezi',
  add_inst_imported: 'Konfiguracija je uvezena ({n} vrijednosti)',
  add_inst_none: 'Još nema instanci; otvaranje panela stvara jednu.',
  add_wd_start: '{name} se nije pokrenuo: panel nije poslao poruku spremnosti unutar 5 sekundi.',
  add_wd_hang: '{name} ne odgovara (nema otkucaja 45 sekundi).',
  add_wd_reload: 'Ponovno učitaj',
  add_wd_diag: 'Dijagnostika',
  add_wd_disable: 'Isključi dodatak',
  add_wd_close: 'Zatvori',
  add_diag_t: 'Dijagnostika: {name}',
  add_diag_none: 'Nema prijavljenih grešaka.',
  add_reload_t: 'Ponovno učitaj panel',
  add_diag_btn_t: 'Greške koje je panel prijavio',
  pg_addons: 'Dodaci',
  po_addon_personal_off: 'Nitko', po_addon_personal_admins: 'Administratori', po_addon_personal_all: 'Svi',
  p_addon_personal: 'Osobni dodaci',
  p_addon_personal_d: 'Tko smije instalirati osobne dodatke koje vidi samo on (Postavke → Dodaci). Rade u istom sandboxu kao dodaci organizacije i ne mogu imati kôd na serveru. Isključeno: dodatke instaliraju samo administratori, za sve.',
});

(function () {
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(String(b)); return s; };
  let list = [], caps = {}, loaded = false, loading = null;
  let community = [], trust = {}, communityError = '';
  let drafts = [], draftsError = '';
  let personalAllowed = false, others = [];
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
    loading = api('GET', '/api/addons' + (isAdmin() ? '?personal=all' : '')).then(d => {
      list = d.addons || [];
      caps = d.capabilities || {};
      personalAllowed = !!d.personal_allowed;
      others = d.personal_others || [];
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

  // a.asset_base is /api/addons/<key>/asset/: the key is the id of an
  // organisation addon or "~id" for one of your personal addons.
  function assetURL(a, file) {
    const base = a.asset_base || ('/api/addons/' + encodeURIComponent(a.id) + '/asset/');
    return base + String(file).split('/').map(encodeURIComponent).join('/');
  }
  const keyOf = a => a.key || a.id;
  const byKey = k => list.find(x => keyOf(x) === k);
  const isAdmin = () => !!(typeof currentUser !== 'undefined' && currentUser && currentUser.is_admin);
  const canManage = a => a && (a.personal ? true : isAdmin());
  const apiBase = a => '/api/addons/' + encodeURIComponent(keyOf(a));

  // ── the panel window ─────────────────────────────────
  function themeTokens() {
    const cs = getComputedStyle(document.documentElement), out = {};
    ['--bg', '--bg2', '--bg3', '--bg4', '--bg5', '--border', '--border2', '--text', '--text2', '--text3',
      '--accent', '--accent-h', '--green', '--red', '--yellow', '--purple', '--font-ui', '--font-mono', '--radius'].forEach(k => { const v = cs.getPropertyValue(k).trim(); if (v) out[k] = v; });
    return out;
  }

  function sendContext(frame, a) {
    let conn = [];
    try {
      conn = (typeof conns !== 'undefined' ? conns : []).map(c => ({id: c.id, name: c.name, host: c.host, folder: c.folder_id || null}));
    } catch (_) {}
    // The panel gets the display name only: the sign-in name stays inside WRM.
    const me = (typeof currentUser !== 'undefined' && currentUser) ? currentUser : {};
    const ctx = {type: 'wrm:context', addon: a.id, user: me.display_name || '', lang: typeof currentLang !== 'undefined' ? currentLang : 'en', theme: themeTokens(), connections: conn};
    try { frame.contentWindow.postMessage(ctx, '*'); } catch (_) {}
  }

  function newFrame(src) {
    const frame = document.createElement('iframe');
    frame.className = 'addon-frame';
    frame.setAttribute('sandbox', 'allow-scripts');
    frame.setAttribute('referrerpolicy', 'no-referrer');
    frame.style.cssText = 'width:100%;height:100%;border:0;display:block;background:#fff';
    frame.src = src;
    return frame;
  }

  function openPanel(key, instanceId) {
    const a = byKey(key);
    if (!a || !a.enabled || !a.ui || !a.ui.panel) return;
    if ((a.apiVersion || 1) >= 2) return openV2(a, instanceId);
    const have = openWins[key];
    if (have && typeof winObjects !== 'undefined' && winObjects.some(w => w.id === have.id && w.el)) {
      if (have.el.style.display === 'none') have.el.style.display = 'flex';
      focusWindow(have.id);
      return;
    }
    if (typeof addWindow !== 'function') return;
    const wo = addWindow({mode: 'addon-panel'});
    wo.addonPanel = key;
    wo.addonTitle = a.ui.title || a.name;
    openWins[key] = wo;
    const tt = $('wt-' + wo.id);
    if (tt) tt.textContent = '▣ ' + wo.addonTitle;
    const wc = $('wc-' + wo.id);
    wc.innerHTML = '';
    wc.style.padding = '0';
    wc.style.overflow = 'hidden';
    const frame = newFrame(assetURL(a, a.ui.panel));
    // The context goes out once, to the panel page WRM loaded. If the frame later
    // navigates (a link in the panel), the new page may be another site: it gets
    // nothing unless it is an addon page that asks (see the message listener).
    frame.addEventListener('load', () => { if (frame.dataset.ctxSent) return; frame.dataset.ctxSent = '1'; sendContext(frame, a); });
    wc.appendChild(frame);
    if (typeof renderTabs === 'function') renderTabs();
  }

  // ── apiVersion 2 panels ──────────────────────────────
  // One window per instance. The panel talks only to this page, only with the
  // message types below; the page does the work (store, notifications) with the
  // user's session, which the panel never sees. A watchdog shows a banner when
  // the panel does not start or stops answering, and the window can always be
  // reloaded or closed.
  const v2 = new Map(); // window id -> panel state
  const instCache = {}; // addon key -> instances
  const notifyLog = {}; // instance id -> [{text, at}]

  async function loadInstances(a) {
    const d = await api('GET', apiBase(a) + '/instances');
    instCache[keyOf(a)] = d.instances || [];
    return instCache[keyOf(a)];
  }

  async function openV2(a, instanceId, opts) {
    opts = opts || {};
    let insts;
    try { insts = await loadInstances(a); } catch (e) { showToast(e.message, 'error', 8000); return; }
    let inst = instanceId ? insts.find(i => i.id === instanceId) : (insts.find(i => i.mine) || insts.find(i => i.shared));
    if (!inst) {
      if (instanceId) return;
      try {
        const d = await api('POST', apiBase(a) + '/instances', {name: t('add_inst_default')});
        insts = await loadInstances(a);
        inst = insts.find(i => i.id === d.id);
      } catch (e) { showToast(e.message, 'error', 8000); return; }
      if (!inst) return;
    }
    for (const st of v2.values()) {
      if (st.inst.id === inst.id && st.wo.el) {
        if (st.wo.el.style.display === 'none') st.wo.el.style.display = 'flex';
        focusWindow(st.wo.id);
        return;
      }
    }
    if (typeof addWindow !== 'function') return;
    const wo = addWindow({mode: 'addon-panel'});
    wo.addonPanel = keyOf(a);
    wo.addonTitle = (a.ui.title || a.name) + (insts.length > 1 ? ' · ' + inst.name : '');
    const st = {a, inst, wo, frame: null, ready: false, started: 0, beat: 0, diag: [], dropped: 0, visible: true, saveT: 0, restoring: !!opts.restore};
    v2.set(wo.id, st);
    const tt = $('wt-' + wo.id);
    if (tt) tt.textContent = '▣ ' + wo.addonTitle;
    addTitleButtons(st);
    placeWindow(st, inst.window);
    const wc = $('wc-' + wo.id);
    wc.innerHTML = '';
    wc.style.padding = '0';
    wc.style.overflow = 'hidden';
    wc.style.position = 'relative';
    loadFrame(st);
    st.restoring = false;
    saveWindow(st, true);
    if (typeof renderTabs === 'function') renderTabs();
  }

  function addTitleButtons(st) {
    const hdr = $('wh-' + st.wo.id);
    const btns = hdr && hdr.querySelector('.win-btns');
    if (!btns) return;
    const min = btns.querySelector('.win-btn.minimize');
    const diag = document.createElement('button');
    diag.className = 'win-btn';
    diag.id = 'wb-adiag-' + st.wo.id;
    diag.title = t('add_diag_btn_t');
    diag.style.display = 'none';
    diag.onclick = e => { e.stopPropagation(); showDiag(st); };
    const rl = document.createElement('button');
    rl.className = 'win-btn';
    rl.title = t('add_reload_t');
    rl.textContent = '↻';
    rl.onclick = e => { e.stopPropagation(); loadFrame(st); };
    btns.insertBefore(diag, min);
    btns.insertBefore(rl, min);
  }

  // Saved geometry is a fraction of the canvas, so it fits another screen.
  function placeWindow(st, saved) {
    const el = st.wo.el, cont = $('win-container');
    if (!el || !cont || (typeof isCompact === 'function' && isCompact())) return;
    const cw = cont.offsetWidth, ch = cont.offsetHeight;
    const ui = st.a.ui || {};
    if (saved && typeof saved.w === 'number' && typeof saved.h === 'number') {
      el.style.left = Math.round(saved.l * cw) + 'px';
      el.style.top = Math.round(saved.t * ch) + 'px';
      el.style.width = Math.round(saved.w * cw) + 'px';
      el.style.height = Math.round(saved.h * ch) + 'px';
    } else if (ui.defaultSize) {
      el.style.width = Math.min(cw, ui.defaultSize.w) + 'px';
      el.style.height = Math.min(ch, ui.defaultSize.h + 34) + 'px';
    }
    if (ui.minSize) {
      el.style.minWidth = Math.min(cw, ui.minSize.w) + 'px';
      el.style.minHeight = Math.min(ch, ui.minSize.h + 34) + 'px';
    }
    if (typeof clampWinToBounds === 'function') clampWinToBounds(el);
    if (saved && saved.snap && typeof snapWindow === 'function') snapWindow(st.wo.id, saved.snap);
    if (saved && saved.min) el.style.display = 'none';
  }

  function saveWindow(st, open) {
    if (st.restoring) return;
    clearTimeout(st.saveT);
    st.saveT = setTimeout(() => {
      const el = st.wo.el, cont = $('win-container');
      const state = {open: !!open};
      if (el && cont && cont.offsetWidth && cont.offsetHeight) {
        const cw = cont.offsetWidth, ch = cont.offsetHeight;
        const r = (v, d) => Math.round(v / d * 10000) / 10000;
        Object.assign(state, {l: r(el.offsetLeft, cw), t: r(el.offsetTop, ch), w: r(el.offsetWidth, cw), h: r(el.offsetHeight, ch),
          min: el.style.display === 'none', snap: st.wo.snap || null});
        if (st.wo.snap && st.wo.savedPos) {
          const p = st.wo.savedPos, n = v => parseFloat(v) || 0;
          Object.assign(state, {l: r(n(p.left), cw), t: r(n(p.top), ch), w: r(n(p.width), cw), h: r(n(p.height), ch)});
        }
      }
      st.inst.window = state;
      fetch(apiBase(st.a) + '/instances/' + st.inst.id + '/window', {method: 'PUT', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(state)}).catch(() => {});
    }, open ? 500 : 0);
  }

  function loadFrame(st) {
    const wc = $('wc-' + st.wo.id);
    if (!wc) return;
    if (st.frame) st.frame.remove();
    banner(st, null);
    st.ready = false;
    st.beat = 0;
    st.started = Date.now();
    st.frame = newFrame(assetURL(st.a, st.a.ui.panel) + '?instance=' + st.inst.id);
    wc.appendChild(st.frame);
  }

  function post(st, msg) {
    if (!st.frame || !st.frame.contentWindow) return;
    try { st.frame.contentWindow.postMessage(Object.assign({v: 2}, msg), '*'); } catch (_) {}
  }

  const hasCap = (a, c) => (a.capabilities || []).includes(c);

  async function contextV2(st) {
    const a = st.a, me = (typeof currentUser !== 'undefined' && currentUser) ? currentUser : {};
    let config = {};
    if (hasCap(a, 'store')) {
      try { config = (await api('GET', apiBase(a) + '/instances/' + st.inst.id + '/store')).values || {}; } catch (_) {}
    }
    const wc = $('wc-' + st.wo.id);
    post(st, {type: 'wrm:context', user: {id: me.id || 0, name: me.display_name || ''},
      lang: typeof currentLang !== 'undefined' ? currentLang : 'en', theme: themeTokens(),
      addon: {id: a.id, key: keyOf(a), version: a.version, instance: st.inst.id, instanceName: st.inst.name},
      config, features: {store: hasCap(a, 'store'), notify: hasCap(a, 'notify'), http: false, jobs: false, history: false, visibility: true},
      size: {w: wc ? wc.clientWidth : 0, h: wc ? wc.clientHeight : 0}, visible: st.visible});
  }

  function banner(st, kind) {
    const wc = $('wc-' + st.wo.id);
    if (!wc) return;
    let b = wc.querySelector('.addon-banner');
    if (!kind) { if (b) b.remove(); st.banner = null; return; }
    if (st.banner === kind && b) return;
    st.banner = kind;
    if (!b) {
      b = document.createElement('div');
      b.className = 'addon-banner';
      b.style.cssText = 'position:absolute;left:0;right:0;top:0;z-index:2;background:var(--bg2);border-bottom:1px solid var(--red);color:var(--text);padding:8px 10px;font-size:12px;display:flex;gap:8px;align-items:center;flex-wrap:wrap;';
      wc.appendChild(b);
    }
    const name = st.a.ui.title || st.a.name;
    b.innerHTML = `<span style="flex:1;min-width:160px;">⚠ ${esc(tf(kind === 'start' ? 'add_wd_start' : 'add_wd_hang', {name}))}</span>
      <button class="btn-sec btn-sm" data-act="reload">${esc(t('add_wd_reload'))}</button>
      <button class="btn-sec btn-sm" data-act="diag">${esc(t('add_wd_diag'))}</button>
      ${canManage(st.a) ? `<button class="btn-danger btn-sm" data-act="disable">${esc(t('add_wd_disable'))}</button>` : ''}
      <button class="btn-sec btn-sm" data-act="close">${esc(t('add_wd_close'))}</button>`;
    b.querySelectorAll('button').forEach(btn => btn.onclick = () => {
      const act = btn.dataset.act;
      if (act === 'reload') loadFrame(st);
      else if (act === 'diag') showDiag(st);
      else if (act === 'disable') toggle(keyOf(st.a), false);
      else if (act === 'close' && typeof closeWindow === 'function') closeWindow(st.wo.id);
    });
  }

  function addDiag(st, entry) {
    st.diag.push(Object.assign({at: new Date().toISOString()}, entry));
    if (st.diag.length > 20) st.diag.shift();
    const b = $('wb-adiag-' + st.wo.id);
    if (b) { b.style.display = ''; b.textContent = '⚠ ' + st.diag.length; }
  }

  function showDiag(st) {
    let m = $('addon-diag-modal');
    if (!m) {
      m = document.createElement('div');
      m.id = 'addon-diag-modal';
      m.className = 'modal-overlay';
      m.style.zIndex = 9660;
      document.body.appendChild(m);
    }
    const rows = st.diag.map(d => `${d.at.slice(11, 19)}  ${d.kind || 'error'}  ${d.message || ''}${d.source ? '  (' + d.source + (d.line ? ':' + d.line : '') + ')' : ''}`).join('\n');
    m.innerHTML = `<div class="modal" style="width:760px;max-width:calc(100vw - 32px);max-height:calc(100vh - 32px);overflow:auto;">
      <h3>${esc(tf('add_diag_t', {name: st.a.ui.title || st.a.name}))}</h3>
      <pre style="margin:0;padding:8px;background:var(--bg2);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;white-space:pre-wrap;max-height:60vh;overflow:auto;color:var(--text2);">${esc(rows || t('add_diag_none'))}</pre>
      <div style="display:flex;justify-content:flex-end;margin-top:10px;"><button onclick="document.getElementById('addon-diag-modal').classList.remove('open')">${esc(t('add_ai_close'))}</button></div>
    </div>`;
    m.classList.add('open');
  }

  // The watchdog: no wrm:ready within 5 s, or three missed heartbeats (15 s each).
  setInterval(() => {
    const now = Date.now();
    for (const st of v2.values()) {
      if (!st.wo.el) continue;
      if (!st.ready) { if (now - st.started > 5000) banner(st, 'start'); }
      else if (now - st.beat > 45000) banner(st, 'hang');
      else if (st.banner === 'hang') banner(st, null);
    }
  }, 2500);

  const shortStr = (v, n) => typeof v === 'string' && v.length <= n;
  const reqID = v => (typeof v === 'string' && v.length <= 64) || (typeof v === 'number' && isFinite(v));

  async function handleV2(st, d) {
    const a = st.a;
    const fail = (type, id, e, extra) => post(st, Object.assign({type, id, error: e.message || String(e), code: e.code || ''}, extra || {}));
    switch (d.type) {
      case 'wrm:ready':
        st.ready = true; st.beat = Date.now(); banner(st, null);
        return contextV2(st);
      case 'wrm:heartbeat':
        st.beat = Date.now();
        return;
      case 'wrm:context':
        return contextV2(st);
      case 'wrm:notify': {
        if (!hasCap(a, 'notify')) return addDiag(st, {kind: 'E_CAPABILITY', message: 'wrm:notify needs the notify capability'});
        if (typeof d.text !== 'string' || !d.text.trim()) return addDiag(st, {kind: 'E_BRIDGE', message: 'wrm:notify needs text'});
        const text = d.text.trim().slice(0, 200), now = Date.now();
        const log = (notifyLog[st.inst.id] = (notifyLog[st.inst.id] || []).filter(x => now - x.at < 300000));
        if (log.some(x => x.text === text) || log.filter(x => now - x.at < 60000).length >= 5) return addDiag(st, {kind: 'E_LIMIT_RATE', message: 'notification dropped: ' + text});
        log.push({text, at: now});
        const level = {error: 'error', warn: 'warning', warning: 'warning', success: 'success'}[d.level || d.kind] || 'info';
        showToast((a.ui.title || a.name) + ': ' + text, level, 6000);
        return;
      }
      case 'wrm:store.get':
      case 'wrm:store.set': {
        if (!reqID(d.id) || !shortStr(d.key, 64)) return addDiag(st, {kind: 'E_BRIDGE', message: d.type + ' needs id and key'});
        const url = apiBase(a) + '/instances/' + st.inst.id + '/store/' + encodeURIComponent(d.key);
        try {
          if (d.type === 'wrm:store.get') {
            const r = await api('GET', url);
            post(st, {type: 'wrm:store', id: d.id, key: d.key, value: r.value === undefined ? null : r.value});
          } else {
            const body = JSON.stringify(d.value === undefined ? null : d.value);
            if (body.length > 32768) throw Object.assign(new Error('A value is larger than 32 KB'), {code: 'E_LIMIT_SIZE'});
            const res = await fetch(url, {method: 'PUT', headers: {'Content-Type': 'application/json'}, body});
            const j = await res.json().catch(() => ({}));
            if (!res.ok) throw Object.assign(new Error(j.error || ('HTTP ' + res.status)), {code: j.code || ''});
            post(st, {type: 'wrm:store', id: d.id, key: d.key, value: d.value, used: j.used});
          }
        } catch (e) {
          addDiag(st, {kind: e.code || 'E_STORE', message: e.message});
          fail('wrm:store', d.id, e, {key: d.key});
        }
        return;
      }
      case 'wrm:resize': {
        const w = Number(d.w), h = Number(d.h), el = st.wo.el, wc = $('wc-' + st.wo.id), cont = $('win-container');
        if (!isFinite(w) || !isFinite(h) || !el || !wc || !cont || st.wo.snap) return;
        if (typeof winMoveBlocked === 'function' && winMoveBlocked(st.wo)) return;
        const min = a.ui.minSize || {w: 160, h: 120};
        const chrome = el.offsetHeight - wc.clientHeight;
        el.style.width = Math.max(min.w, Math.min(cont.offsetWidth, Math.round(w))) + 'px';
        el.style.height = Math.max(min.h + chrome, Math.min(cont.offsetHeight, Math.round(h) + chrome)) + 'px';
        if (typeof clampWinToBounds === 'function') clampWinToBounds(el);
        saveWindow(st, true);
        return;
      }
      case 'wrm:error':
        return addDiag(st, {kind: 'panel', message: String(d.message || '').slice(0, 500), source: String(d.source || '').slice(0, 200), line: Number(d.line) || 0});
      default:
        st.dropped++;
        addDiag(st, {kind: 'E_BRIDGE', message: 'unknown message type dropped: ' + String(d.type).slice(0, 40)});
    }
  }

  // Minimize, restore, move and resize come here from the window code
  // (winLayoutChanged); a hidden panel is told so it can pause its polling.
  function windowChanged(wo) {
    const st = v2.get(wo.id);
    if (!st) return;
    syncVisibility(st);
    saveWindow(st, true);
  }

  function syncVisibility(st) {
    const el = st.wo.el;
    const vis = !!el && el.style.display !== 'none' && !el.classList.contains('folder-hidden') && document.visibilityState !== 'hidden';
    if (vis !== st.visible) {
      st.visible = vis;
      if (st.ready) post(st, {type: 'wrm:visibility', visible: vis});
    }
  }

  function windowGone(wo) {
    const st = v2.get(wo.id);
    if (!st) return;
    post(st, {type: 'wrm:visibility', visible: false});
    saveWindow(st, false);
    v2.delete(wo.id);
  }

  function closeAddonWindows(key, code) {
    for (const st of Array.from(v2.values())) {
      if (keyOf(st.a) !== key) continue;
      post(st, {type: 'wrm:fatal', code, message: code === 'E_DISABLED' ? 'The addon was disabled' : 'The addon was removed'});
      v2.delete(st.wo.id);
      try { if (st.wo.el) closeWindow(st.wo.id); } catch (_) {}
    }
    const wo = openWins[key];
    if (wo) { delete openWins[key]; try { if (wo.el) closeWindow(wo.id); } catch (_) {} }
  }

  document.addEventListener('visibilitychange', () => { for (const st of v2.values()) syncVisibility(st); });

  // A theme switch (Settings → Appearance, an addon theme) reaches open panels.
  let themeT = 0, lastTheme = '';
  function themeMaybeChanged() {
    clearTimeout(themeT);
    themeT = setTimeout(() => {
      const tok = themeTokens(), sig = JSON.stringify(tok);
      if (sig === lastTheme) return;
      lastTheme = sig;
      for (const st of v2.values()) if (st.ready) post(st, {type: 'wrm:theme', theme: tok});
    }, 150);
  }
  try {
    new MutationObserver(themeMaybeChanged).observe(document.documentElement, {attributes: true, attributeFilter: ['style', 'class', 'data-theme']});
    new MutationObserver(themeMaybeChanged).observe(document.head, {childList: true});
  } catch (_) {}
  window.addEventListener('wrm-theme-changed', themeMaybeChanged);

  // Windows that were open come back after a reload or on another device.
  async function restoreWindows() {
    for (const a of list) {
      if (!a.enabled || (a.apiVersion || 1) < 2 || !a.ui || !a.ui.panel) continue;
      let insts;
      try { insts = await loadInstances(a); } catch (_) { continue; }
      for (const inst of insts) if (inst.window && inst.window.open) await openV2(a, inst.id, {restore: true});
    }
  }

  // The only words a panel gets back into WRM. v1 panels: a toast, or a
  // request for the context again. v2 panels: the protocol above. Anything
  // else is ignored (v2: dropped and counted in the diagnostics).
  window.addEventListener('message', ev => {
    const d = ev.data;
    if (!d || typeof d !== 'object' || typeof d.type !== 'string') return;
    // Addon pages are served with the sandbox CSP, so their origin is opaque
    // ("null"); a page from any real origin in the frame is not an addon page.
    if (ev.origin !== 'null') return;
    const frame = Array.from(document.querySelectorAll('iframe.addon-frame')).find(f => f.contentWindow === ev.source);
    if (!frame) return;
    const st = Array.from(v2.values()).find(x => x.frame === frame);
    if (st) {
      if (d.v !== 2) { st.dropped++; addDiag(st, {kind: 'E_BRIDGE', message: 'message without v: 2 dropped: ' + d.type.slice(0, 40)}); return; }
      handleV2(st, d);
      return;
    }
    const a = list.find(x => { const wo = openWins[keyOf(x)]; return wo && wo.el && wo.el.contains(frame); });
    if (d.type === 'wrm:notify' && typeof d.text === 'string') {
      showToast(String(d.text).slice(0, 500), d.kind === 'error' ? 'error' : 'success', 6000);
    } else if (d.type === 'wrm:context' && a) {
      sendContext(frame, a);
    }
  });

  // ── permission sheet ─────────────────────────────────
  function permText(p, a) {
    const [kind, val, extra] = String(p).split(':');
    if (kind === 'panel') return t((a.apiVersion || 1) >= 2 ? 'add_perm_panel_v2' : 'add_perm_panel');
    if (kind === 'themes') return tf('add_perm_themes', {n: val});
    if (kind === 'store') return tf(extra === 'shared' ? 'add_perm_store_shared' : 'add_perm_store', {kb: val});
    if (kind === 'notify') return t('add_perm_notify');
    if (kind === 'server') return tf('add_perm_server', {exports: (val || '').split(',').join(', ')});
    return tf('add_perm_other', {p});
  }

  function approveSheet(a) {
    return new Promise(resolve => {
      let m = $('addon-approve-modal');
      if (!m) {
        m = document.createElement('div');
        m.id = 'addon-approve-modal';
        m.className = 'modal-overlay';
        m.style.zIndex = 9665;
        document.body.appendChild(m);
      }
      const added = new Set(a.permissions_added || []);
      const items = (a.permissions || []).map(p => `<li style="margin:4px 0;">${esc(permText(p, a))}${added.has(p) ? ` <span class="pill warn">${esc(t('add_approve_new'))}</span>` : ''}</li>`).join('');
      m.innerHTML = `<div class="modal" style="width:560px;max-width:calc(100vw - 32px);">
        <h3>${esc(tf('add_approve_t', {name: a.name}))} ${a.personal ? `<span class="pill">${esc(t('add_personal'))}</span>` : ''}</h3>
        <div class="hint-box">${esc(t('add_approve_d'))}</div>
        <ul style="margin:8px 0 0 18px;padding:0;font-size:12.5px;">${items}</ul>
        <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:12px;">
          <button class="btn-sec" data-act="no">${esc(t('cancel'))}</button>
          <button data-act="yes">${esc(t('add_approve_go'))}</button>
        </div></div>`;
      const done = v => { m.classList.remove('open'); resolve(v); };
      m.querySelector('[data-act="no"]').onclick = () => done(false);
      m.querySelector('[data-act="yes"]').onclick = () => done(true);
      m.classList.add('open');
    });
  }

  // ── settings pane ────────────────────────────────────
  function pill(text) { return `<span class="pill">${esc(text)}</span>`; }

  function instancesHTML(a) {
    if (!a.enabled || (a.apiVersion || 1) < 2 || !a.ui || !a.ui.panel) return '';
    const k = esc(keyOf(a)), insts = instCache[keyOf(a)] || [];
    const rows = insts.map(i => `<div style="display:flex;gap:6px;align-items:center;flex-wrap:wrap;margin:4px 0;">
        <b style="font-size:12px;">${esc(i.name)}</b> ${pill(i.shared ? t('add_inst_shared') : i.mine ? t('add_inst_mine') : (i.owner || ''))}
        <span style="flex:1;"></span>
        <button class="btn-sec btn-sm" onclick="wrmAddons.openPanel('${k}', ${i.id})">${esc(t('add_open'))}</button>
        ${i.mine || (i.shared && isAdmin()) ? `<button class="btn-sec btn-sm" onclick="wrmAddons.instRename('${k}', ${i.id})">${esc(t('add_inst_rename'))}</button>` : ''}
        ${hasCap(a, 'store') ? `<button class="btn-sec btn-sm" onclick="wrmAddons.instExport('${k}', ${i.id})">${esc(t('add_inst_export'))}</button>
        <button class="btn-sec btn-sm" onclick="wrmAddons.instImport('${k}', ${i.id})">${esc(t('add_inst_import'))}</button>` : ''}
        ${i.mine || isAdmin() ? `<button class="btn-danger btn-sm" onclick="wrmAddons.instDelete('${k}', ${i.id})">${esc(t('add_remove'))}</button>` : ''}
      </div>`).join('');
    return `<details style="margin-top:6px;" ${insts.length > 1 ? 'open' : ''}>
      <summary style="cursor:pointer;font-size:11.5px;color:var(--text2);">${esc(t('add_inst'))} (${insts.length})</summary>
      <div class="setting-desc" style="margin:4px 0;">${esc(t('add_inst_d'))}</div>
      ${rows || `<div class="sub" style="color:var(--text3);font-size:11.5px;">${esc(t('add_inst_none'))}</div>`}
      <button class="btn-sec btn-sm" style="margin-top:4px;" onclick="wrmAddons.instNew('${k}')">＋ ${esc(t('add_inst_new'))}</button>
    </details>`;
  }

  function rowHTML(a) {
    const admin = isAdmin(), k = esc(keyOf(a));
    const panel = a.ui && a.ui.panel;
    const btns = [];
    if (panel && a.enabled) btns.push(`<button class="btn-sec btn-sm" onclick="wrmAddons.openPanel('${k}')">${esc(t('add_open'))}</button>`);
    if (canManage(a)) {
      btns.push(`<label class="mcp-check" style="margin:0;display:flex;align-items:center;gap:6px;"><input type="checkbox" ${a.enabled ? 'checked' : ''} onchange="wrmAddons.toggle('${k}', this.checked)"><span>${esc(a.enabled ? t('add_disable') : t('add_enable'))}</span></label>`);
      btns.push(`<button class="btn-danger btn-sm" onclick="wrmAddons.remove('${k}')">${esc(t('add_remove'))}</button>`);
    }
    const meta = [a.author ? tf('add_by', {author: a.author}) : '', a.license || '', a.homepage || ''].filter(Boolean).map(esc).join(' · ');
    const exports = (a.server && a.server.exports) || [];
    const runner = admin && !a.personal && a.enabled && exports.length ? `<details style="margin-top:6px;">
        <summary style="cursor:pointer;font-size:11.5px;color:var(--text2);">${esc(t('add_run'))} · ${esc(a.server.module)}</summary>
        <div class="setting-desc" style="margin:4px 0;">${esc(t('add_run_d'))}</div>
        <div style="display:flex;gap:6px;align-items:flex-start;flex-wrap:wrap;">
          <select id="addon-export-${esc(a.id)}" style="background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font-size:11.5px;padding:4px 6px;">${exports.map(x => `<option value="${esc(x)}">${esc(x)}</option>`).join('')}</select>
          <textarea id="addon-in-${esc(a.id)}" placeholder="${esc(t('add_run_in'))}" style="flex:1;min-width:180px;min-height:44px;background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;padding:5px 7px;"></textarea>
          <button class="btn-sec btn-sm" onclick="wrmAddons.run('${esc(a.id)}')">${esc(t('add_run'))}</button>
        </div>
        <pre id="addon-out-${esc(a.id)}" style="margin:6px 0 0;padding:6px 8px;background:var(--bg2);border:1px solid var(--border);border-radius:6px;font:11.5px 'JetBrains Mono',monospace;white-space:pre-wrap;max-height:180px;overflow:auto;color:var(--text2);">${esc(t('add_run_out'))}</pre>
      </details>` : '';
    const wait = a.needs_approval ? `<div class="hint-box warn" style="margin-top:6px;">${esc(tf('add_needs_approval', {list: (a.permissions_added || []).map(p => permText(p, a)).join(' · ')}))}</div>` : '';
    return `<div style="border:1px solid var(--border);border-radius:8px;padding:10px;margin:8px 0;${a.enabled ? '' : 'opacity:.75;'}">
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        <b>${esc(a.name)}</b> ${pill(a.version)} ${a.personal ? `<span class="pill ok">${esc(t('add_personal'))}</span>` : ''} ${(a.apiVersion || 1) >= 2 ? pill('v2') : ''}
        ${(a.capabilities || []).map(c => pill(c)).join(' ')}
        <span style="flex:1;"></span>${btns.join(' ')}
      </div>
      ${a.description ? `<div class="sub" style="font-size:11.5px;color:var(--text2);margin-top:4px;">${esc(a.description)}</div>` : ''}
      ${meta ? `<div class="sub" style="font-size:11px;color:var(--text3);margin-top:2px;">${meta}</div>` : ''}
      <div class="sub" style="font-size:10.5px;color:var(--text3);margin-top:2px;">${esc(a.id)} · ${esc(a.checksum || '')}${a.installed_at ? ' · ' + esc(fmtTime(a.installed_at)) : ''}</div>
      ${wait}
      ${instancesHTML(a)}
      ${runner}
    </div>`;
  }

  function othersHTML() {
    if (!isAdmin()) return '';
    const rows = others.map(a => `<div style="border:1px solid var(--border);border-radius:8px;padding:8px 10px;margin:6px 0;${a.enabled ? '' : 'opacity:.75;'}">
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        <b>${esc(a.name)}</b> ${pill(a.version)} <span class="pill">${esc(tf('add_personal_of', {user: a.owner || ('#' + a.owner_id)}))}</span>
        ${(a.capabilities || []).map(c => pill(c)).join(' ')}
        <span style="flex:1;"></span>
        ${a.enabled ? `<button class="btn-sec btn-sm" onclick="wrmAddons.otherDisable(${a.owner_id}, '${esc(a.id)}')">${esc(t('add_disable'))}</button>` : ''}
        <button class="btn-danger btn-sm" onclick="wrmAddons.otherRemove(${a.owner_id}, '${esc(a.id)}')">${esc(t('add_remove'))}</button>
      </div></div>`).join('');
    return `<div class="adm-h">${esc(t('add_personal_others'))}</div>
      <div class="setting-desc">${esc(t('add_personal_others_d'))}</div>
      ${rows || `<div class="sub" style="color:var(--text3);font-size:11.5px;">${esc(t('add_personal_none'))}</div>`}`;
  }

  function render() {
    const el = $('addons-pane');
    if (!el) return;
    const admin = isAdmin();
    const scope = admin && personalAllowed ? `<select id="addon-scope" style="max-width:180px;background:var(--bg3);color:var(--text);border:1px solid var(--border);border-radius:6px;font-size:11.5px;padding:4px 6px;" title="${esc(t('add_scope'))}">
        <option value="org">${esc(t('add_scope_org'))}</option><option value="personal">${esc(t('add_scope_me'))}</option></select>` : '';
    const install = admin || personalAllowed ? `<div class="setting-row">
      <div><div class="setting-label">${esc(t(admin ? 'add_install' : 'add_install_me'))}</div><div class="setting-desc">${esc(t('add_install_d'))}${admin ? '' : ' ' + esc(t('add_install_me_d'))}</div></div>
      <div class="setting-control" style="display:flex;gap:8px;align-items:center;flex-wrap:wrap;">
        ${scope}
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
      ${othersHTML()}
      ${communityHTML()}`;
  }

  async function pane() {
    const el = $('addons-pane');
    if (!el) return;
    el.innerHTML = `<div class="sub" style="color:var(--text3);">…</div>`;
    try { await refresh(); } catch (e) { el.innerHTML = `<div class="fm-error">${esc(e.message)}</div>`; return; }
    await Promise.all(list.filter(a => a.enabled && (a.apiVersion || 1) >= 2 && a.ui && a.ui.panel).map(a => loadInstances(a).catch(() => {})));
    await loadCommunity();
    await loadDrafts();
    render();
  }

  async function after(msg) {
    if (msg) showToast(msg, 'success');
    try { await refresh(); } catch (_) {}
    await Promise.all(list.filter(a => a.enabled && (a.apiVersion || 1) >= 2 && a.ui && a.ui.panel).map(a => loadInstances(a).catch(() => {})));
    await loadDrafts();
    applyThemes();
    if ($('addons-pane') && $('st-addons') && $('st-addons').classList.contains('active')) render();
  }

  async function install(allowDowngrade) {
    const inp = $('addon-file'), btn = $('addon-install-btn'), sc = $('addon-scope');
    const file = inp && inp.files && inp.files[0];
    if (!file) { showToast(t('add_need_file'), 'error'); return; }
    if (file.size > 20 * 1024 * 1024) { showToast(t('add_too_big'), 'error'); return; }
    if (btn) btn.disabled = true;
    try {
      const q = [];
      if (sc && sc.value === 'personal') q.push('scope=personal');
      if (allowDowngrade === true) q.push('allow_downgrade=1');
      const res = await fetch('/api/addons/install' + (q.length ? '?' + q.join('&') : ''), {method: 'POST', headers: {'Content-Type': 'application/zip'}, body: file});
      const d = await res.json().catch(() => ({}));
      if (res.status === 409 && d.code === 'E_DOWNGRADE' && allowDowngrade !== true) {
        if (btn) btn.disabled = false;
        if (await uiConfirm(tf('add_downgrade_q', {error: d.error}), {okText: t('add_go')})) return install(true);
        return;
      }
      if (!res.ok) throw new Error((d.error || ('HTTP ' + res.status)) + (d.hint ? ' — ' + d.hint : ''));
      await after(tf(d.needs_approval ? 'add_installed_wait' : 'add_installed', {name: d.name || d.id || file.name}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    } finally {
      if (btn) btn.disabled = false;
    }
  }

  async function toggle(key, on) {
    const a = byKey(key);
    try {
      if (on && a && !(await approveSheet(a))) { render(); return; }
      await api('POST', '/api/addons/' + encodeURIComponent(key) + '/' + (on ? 'enable' : 'disable'));
      if (!on) closeAddonWindows(key, 'E_DISABLED');
      await after(tf(on ? 'add_enabled_toast' : 'add_disabled_toast', {name: a ? a.name : key}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
      render();
    }
  }

  async function remove(key) {
    const a = byKey(key);
    if (!(await uiConfirm(tf('add_remove_q', {name: a ? a.name : key}), {okText: t('add_remove'), danger: true}))) return;
    try {
      await api('DELETE', '/api/addons/' + encodeURIComponent(key));
      closeAddonWindows(key, 'E_REMOVED');
      await after(tf('add_removed', {name: a ? a.name : key}));
    } catch (e) {
      showToast(e.message, 'error', 8000);
    }
  }

  async function otherDisable(uid, id) {
    try { await api('POST', '/api/addons/users/' + uid + '/' + encodeURIComponent(id) + '/disable'); await after(tf('add_disabled_toast', {name: id})); }
    catch (e) { showToast(e.message, 'error', 8000); }
  }

  async function otherRemove(uid, id) {
    const a = others.find(x => x.owner_id === uid && x.id === id);
    if (!(await uiConfirm(tf('add_remove_q', {name: a ? a.name : id}), {okText: t('add_remove'), danger: true}))) return;
    try { await api('DELETE', '/api/addons/users/' + uid + '/' + encodeURIComponent(id)); await after(tf('add_removed', {name: a ? a.name : id})); }
    catch (e) { showToast(e.message, 'error', 8000); }
  }

  // ── instances ────────────────────────────────────────
  async function instNew(key) {
    const a = byKey(key);
    if (!a) return;
    const name = await uiPrompt(t('add_inst_name_q'), '');
    if (!name || !name.trim()) return;
    let shared = false;
    if (isAdmin() && !a.personal) shared = !!(await uiConfirm(t('add_inst_shared_q'), {okText: t('add_inst_shared'), cancelText: t('add_inst_mine')}));
    try {
      const d = await api('POST', apiBase(a) + '/instances', {name: name.trim(), shared});
      await after('');
      openPanel(key, d.id);
    } catch (e) { showToast(e.message, 'error', 8000); }
  }

  async function instRename(key, iid) {
    const a = byKey(key), i = (instCache[key] || []).find(x => x.id === iid);
    if (!a || !i) return;
    const name = await uiPrompt(t('add_inst_rename_q'), i.name);
    if (!name || !name.trim()) return;
    try { await api('PATCH', apiBase(a) + '/instances/' + iid, {name: name.trim()}); await after(''); }
    catch (e) { showToast(e.message, 'error', 8000); }
  }

  async function instDelete(key, iid) {
    const a = byKey(key), i = (instCache[key] || []).find(x => x.id === iid);
    if (!a || !i) return;
    if (!(await uiConfirm(tf('add_inst_delete_q', {name: i.name}), {okText: t('add_remove'), danger: true}))) return;
    for (const st of Array.from(v2.values())) if (st.inst.id === iid) { v2.delete(st.wo.id); try { closeWindow(st.wo.id); } catch (_) {} }
    try { await api('DELETE', apiBase(a) + '/instances/' + iid); await after(''); }
    catch (e) { showToast(e.message, 'error', 8000); }
  }

  async function instExport(key, iid) {
    const a = byKey(key);
    if (!a) return;
    try {
      const d = await api('GET', apiBase(a) + '/instances/' + iid + '/export');
      const blob = new Blob([JSON.stringify(d, null, 2)], {type: 'application/json'});
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = a.id + '-' + iid + '.json';
      document.body.appendChild(link);
      link.click();
      setTimeout(() => { URL.revokeObjectURL(link.href); link.remove(); }, 1000);
    } catch (e) { showToast(e.message, 'error', 8000); }
  }

  function instImport(key, iid) {
    const a = byKey(key);
    if (!a) return;
    const inp = document.createElement('input');
    inp.type = 'file';
    inp.accept = '.json,application/json';
    inp.onchange = async () => {
      const f = inp.files && inp.files[0];
      if (!f) return;
      try {
        const body = JSON.parse(await f.text());
        const d = await api('POST', apiBase(a) + '/instances/' + iid + '/import', body);
        for (const st of v2.values()) if (st.inst.id === iid && st.ready) contextV2(st);
        showToast(tf('add_inst_imported', {n: d.keys || 0}), 'success');
      } catch (e) { showToast(e.message, 'error', 8000); }
    };
    inp.click();
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
      'Deliver it as ONE zip with addon.json at the root plus the files it declares. Prefer apiVersion 2.',
      '',
      'addon.json (apiVersion 2; unknown fields are rejected, so spell every field exactly):',
      '{',
      '  "apiVersion": 2,',
      '  "id": "lowercase.dotted-id",',
      '  "name": "Short name",',
      '  "version": "1.0.0",',
      '  "description": "One sentence",',
      '  "author": "Your name",',
      '  "license": "MIT",',
      '  "homepage": "https://example.com/addon",',
      '  "capabilities": ["ui.panel", "store", "notify"],',
      '  "ui": {"panel": "panel.html", "title": "Window title", "defaultSize": {"w": 760, "h": 500}, "minSize": {"w": 320, "h": 220},',
      '         "themes": [{"id": "ocean", "name": "Ocean", "name_hr": "Ocean", "mode": "dark", "tokens": {"--bg": "#0b1220", "--accent": "#38bdf8"}, "css": "ocean.css"}]},',
      '  "store": {"quotaKb": 64, "shared": false}',
      '}',
      '',
      'Capabilities this WRM build accepts:',
      capLines,
      '',
      'Panel (ui.panel) rules:',
      '- The panel is HTML in an iframe with sandbox="allow-scripts", an opaque origin and this CSP: default-src \'none\'; script-src \'self\' plus the hash of each inline <script>; style-src \'self\' \'unsafe-inline\'; img-src \'self\' data: blob:; connect-src \'none\'. So: no fetch/XHR/WebSocket, no inline event handlers (onclick="…"), no eval, no remote scripts or fonts; put JS in <script> blocks or files of the zip and use addEventListener.',
      '- Talk to WRM only with window.parent.postMessage(msg, "*"); every message carries v: 2. Listen with window.addEventListener("message", e => …) and check e.data.type.',
      '- On load send {v:2, type:"wrm:ready"} (within 5 s) and then {v:2, type:"wrm:heartbeat"} every 15 s, or WRM shows a "not responding" banner.',
      '- WRM answers wrm:ready with {type:"wrm:context", v:2, user:{id,name}, lang, theme:{"--bg":…}, addon:{id, key, version, instance, instanceName}, config:{…stored values…}, features:{store, notify, http:false, visibility:true}, size:{w,h}, visible}. Ask again with {v:2, type:"wrm:context"}.',
      '- WRM also sends {type:"wrm:theme", theme} when the user switches theme, {type:"wrm:visibility", visible} when the window is minimized or restored (stop polling and timers while hidden), and {type:"wrm:fatal", code, message} before it closes the panel.',
      '- store capability: {v:2, type:"wrm:store.get", id, key} and {v:2, type:"wrm:store.set", id, key, value} (any JSON, 32 KB per value, quotaKb per instance). The answer is {type:"wrm:store", id, key, value} or with error and code (E_LIMIT_QUOTA, E_LIMIT_SIZE, E_STORE_KEY, E_STORE_READONLY). Keys: 1-64 letters, digits, dot, dash, underscore, not starting with wrm.',
      '- notify capability: {v:2, type:"wrm:notify", text, level:"info"|"success"|"warn"|"error"}; plain text, 200 characters, at most 5 a minute, repeats within 5 minutes are dropped.',
      '- {v:2, type:"wrm:resize", w, h} asks for a content size (WRM decides); {v:2, type:"wrm:error", message, source, line} reports an error into the window diagnostics. Use a ResizeObserver to redraw.',
      '- Each instance (the user may open several, e.g. two dashboards) has its own stored values and window. Data is fetched only while the window is open; WRM runs no background jobs for addons.',
      '- There is no network access from the panel in this version. Use data the user enters or pastes, or wait for the HTTP proxy capability.',
      '',
      'Themes (ui.theme capability): ui.themes lists themes users pick in Settings → Appearance (never applied to everybody). tokens are CSS variables of WRM (--bg --bg2 --bg3 --bg4 --bg5 --border --border2 --text --text2 --text3 --accent --green --red --yellow --purple --font-ui --font-mono --radius), values without ; { } < > url( or @import; css is an optional stylesheet in the zip. A v1 "theme": "theme.css" still works.',
      '',
      'server.wasm (organisation addons only, administrators): WebAssembly run by wazero with no WASI: no files, no network, no environment, no clock. The only import is wrm.log(ptr, len). Export memory, wrm_alloc(size) -> ptr, and each declared export as (ptr, len) -> ptr returning a NUL-terminated UTF-8 string. 64 KiB input, 1 MB output, 5 seconds per call, 64 MiB memory. "server": {"module": "plugin.wasm", "exports": ["status"]}.',
      '',
      'Nothing runs at install time. Every file is pinned by SHA-256. Enabling the addon shows the user a permission sheet generated from addon.json; an update that asks for more stays off until it is approved again. Errors come back with a code (E_MANIFEST, E_MANIFEST_FIELD, E_MISSING_FILE, E_LIMIT_SIZE, …) and a hint: fix exactly what they name.',
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
    try { await restoreWindows(); } catch (_) {}
    if ($('addons-pane') && $('st-addons') && $('st-addons').classList.contains('active')) render();
  }

  window.wrmAddons = {boot, pane, install, toggle, remove, openPanel, run, aiPrompt, closeAIPrompt, copyPrompt, openAssistant, communityInstall, applyThemes, reload: after,
    draftInstall, draftDiff, draftDiscard, closeDraftDiff, windowChanged, windowGone, list: () => list.slice(),
    instNew, instRename, instDelete, instExport, instImport, otherDisable, otherRemove};

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
