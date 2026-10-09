// Git workspace of WRM PRO: compare services on servers with GitLab / GitHub or offline bundles.
// Loaded on demand by openGit() in index.html; uses its helpers (t, esc, showToast, apiError,
// uiConfirm, fmtTime, conns, folders, currentUser).
'use strict';

Object.assign(LANGS.en, {
  git_title: 'Git workspace', git_close: 'Close', git_tab_overview: 'Overview', git_tab_services: 'Services', git_tab_installs: 'Installations', git_tab_settings: 'Settings',
  git_check_now: 'Check now', git_checking: 'Checking…', git_last_check: 'Last check', git_never: 'never', git_no_checks: 'Checks are not allowed for your account.',
  git_st_ok: 'Up to date', git_st_update: 'Needs update', git_st_review: 'Review', git_st_error: 'Error', git_st_gone: 'Gone', git_st_unknown: 'Not checked',
  git_fs_ok: 'ok', git_fs_old: 'old', git_fs_modified: 'modified', git_fs_missing: 'missing', git_fs_extra: 'extra', git_fs_protected: 'protected',
  git_behind: '{n} behind', git_vmd_missing: 'missing', git_version_md: 'VERSION.md', git_target: 'Target', git_latest_tag: 'Latest tag',
  git_all_states: 'All states', git_all_envs: 'All environments', git_all_folders: 'All folders', git_all_tags: 'All tags', git_search: 'Search…',
  git_ref_choice: 'Ref', git_ref_catalog: 'As in the catalog', git_ref_default: 'Default branch of each project', git_ref_branch: 'One branch for all',
  git_ref_override: 'Ref for this run', git_ref_override_ph: 'tag:latest, tag:v2 or a branch (empty = as chosen above)',
  git_server: 'Server', git_service: 'Service', git_path: 'Path', git_env: 'Environment', git_state: 'State', git_files: 'Files', git_checked: 'Checked', git_actions: 'Actions',
  git_empty_overview: 'Nothing to show yet. Add a source or import a bundle, fill in the catalog, choose the servers in Settings and run a check.',
  git_targets: 'Targets', git_refresh_targets: 'Refresh targets', git_source: 'Source', git_warnings: 'Warnings', git_no_target: 'no target yet',
  git_add_service: 'Add service', git_suggest: 'Suggest from repository', git_import_catalog: 'Import catalog', git_export_catalog: 'Export catalog',
  git_import_merge: 'Merge with the current catalog (OK) or replace it (Cancel)?', git_export_bundle: 'Export bundle', git_export_bundle_none: 'Select services with a target first.',
  git_from_bundle: 'From imported bundles (not in the catalog)', git_edit: 'Edit', git_delete: 'Delete', git_delete_q: 'Delete service {name} from the catalog?', git_save: 'Save', git_cancel: 'Cancel',
  git_f_name: 'Name', git_f_project: 'Project (group/name)', git_f_ref: 'Ref', git_f_ref_h: 'tag:latest, tag:v1 or a branch', git_f_branch: 'Branch (for tags)', git_f_tag_filter: 'Tag filter (regex)',
  git_f_subdir: 'Subdirectory = installation root', git_f_kind: 'Kind', git_kind_app: 'App', git_kind_library: 'Library', git_kind_tool: 'Tool',
  git_f_include: 'Include globs', git_f_exclude: 'Exclude globs', git_f_protected: 'Protected (never overwritten, per host)', git_f_fingerprint: 'Fingerprint files', git_f_hint: 'Install hints (directory names)',
  git_lines_h: 'One per line or comma-separated', git_suggest_project: 'Project', git_suggest_ref: 'Branch (empty = default)', git_suggest_go: 'Suggest', git_repo_files: 'Repository files',
  git_discover: 'Discover', git_add_install: 'Add by hand', git_details: 'Details', git_forget: 'Forget', git_forget_q: 'Forget this installation? Discovery may find it again.',
  git_units: 'Units', git_no_units: 'No systemd or supervisor unit mentions this directory.', git_diff: 'Diff', git_diff_title: 'Server → target', git_binary: 'Binary file',
  git_no_diff: 'No differences in the text (line endings or whitespace only).', git_capped: 'Only the first 5000 files were checked.', git_files_only: 'Show only changes',
  git_sources: 'Git sources', git_add_source: 'Add source', git_kind: 'Kind', git_url: 'Address', git_token: 'Read-only token', git_token_keep: 'leave empty to keep the stored token',
  git_test: 'Test', git_projects_found: '{n} projects', git_bundles: 'Offline bundles', git_import_bundle: 'Import bundle (.tar.gz)', git_bundle_age: 'created {d}',
  git_bundle_imported: 'Bundle imported: {n} files', git_source_mode: 'Targets from', git_mode_auto: 'Automatic (Git API, otherwise the newest bundle)', git_mode_live: 'Git API',
  git_mode_bundle: 'Offline bundle', git_newest_bundle: 'Newest bundle with the service', git_history_depth: 'History depth (earlier versions per file)',
  git_servers: 'Servers to check', git_servers_h: 'SSH connections. Checks only read files, nothing is changed on servers.', git_filter: 'Filter…', git_select_all: 'All', git_select_none: 'None',
  git_catalog_settings: 'Catalog settings', git_groups: 'Groups / organisations', git_fallback: 'Fallback branches', git_roots: 'Server roots (searched to depth 3)', git_ignore: 'Ignored directories (globs)',
  git_monitor: 'Periodic checks', git_monitor_on: 'Check my servers every {n} minutes while WRM runs and notify me', git_monitor_off: 'Periodic checks are turned off by the administrator.',
  git_monitor_h: 'Choose the channels for "Git" events in Settings → Notifications. Checks never update anything.', git_saved: 'Saved',
  git_n_installs: '{n} installations', git_none: '—', git_api_calls: '{n} API calls', git_from: 'from', git_bundle_only: 'bundle',
});
Object.assign(LANGS.hr, {
  git_title: 'Git radni prostor', git_close: 'Zatvori', git_tab_overview: 'Pregled', git_tab_services: 'Servisi', git_tab_installs: 'Instalacije', git_tab_settings: 'Postavke',
  git_check_now: 'Provjeri sada', git_checking: 'Provjera…', git_last_check: 'Zadnja provjera', git_never: 'nikad', git_no_checks: 'Provjere nisu dopuštene vašem računu.',
  git_st_ok: 'Ažurno', git_st_update: 'Treba ažurirati', git_st_review: 'Pregledati', git_st_error: 'Greška', git_st_gone: 'Nema je', git_st_unknown: 'Nije provjereno',
  git_fs_ok: 'ok', git_fs_old: 'staro', git_fs_modified: 'izmijenjeno', git_fs_missing: 'nedostaje', git_fs_extra: 'višak', git_fs_protected: 'zaštićeno',
  git_behind: '{n} iza', git_vmd_missing: 'nema', git_version_md: 'VERSION.md', git_target: 'Cilj', git_latest_tag: 'Zadnji tag',
  git_all_states: 'Sva stanja', git_all_envs: 'Sva okruženja', git_all_folders: 'Sve mape', git_all_tags: 'Sve oznake', git_search: 'Traži…',
  git_ref_choice: 'Ref', git_ref_catalog: 'Kao u katalogu', git_ref_default: 'Zadana grana svakog projekta', git_ref_branch: 'Jedna grana za sve',
  git_ref_override: 'Ref za ovu provjeru', git_ref_override_ph: 'tag:latest, tag:v2 ili grana (prazno = kao gore)',
  git_server: 'Server', git_service: 'Servis', git_path: 'Putanja', git_env: 'Okruženje', git_state: 'Stanje', git_files: 'Datoteke', git_checked: 'Provjereno', git_actions: 'Akcije',
  git_empty_overview: 'Još nema ničega. Dodajte izvor ili uvezite bundle, popunite katalog, odaberite servere u Postavkama i pokrenite provjeru.',
  git_targets: 'Ciljevi', git_refresh_targets: 'Osvježi ciljeve', git_source: 'Izvor', git_warnings: 'Upozorenja', git_no_target: 'još nema cilja',
  git_add_service: 'Dodaj servis', git_suggest: 'Prijedlog iz repozitorija', git_import_catalog: 'Uvezi katalog', git_export_catalog: 'Izvezi katalog',
  git_import_merge: 'Spojiti s trenutnim katalogom (U redu) ili ga zamijeniti (Odustani)?', git_export_bundle: 'Izvezi bundle', git_export_bundle_none: 'Najprije odaberite servise koji imaju cilj.',
  git_from_bundle: 'Iz uvezenih bundleova (nisu u katalogu)', git_edit: 'Uredi', git_delete: 'Obriši', git_delete_q: 'Obrisati servis {name} iz kataloga?', git_save: 'Spremi', git_cancel: 'Odustani',
  git_f_name: 'Naziv', git_f_project: 'Projekt (grupa/naziv)', git_f_ref: 'Ref', git_f_ref_h: 'tag:latest, tag:v1 ili grana', git_f_branch: 'Grana (za tagove)', git_f_tag_filter: 'Filtar tagova (regex)',
  git_f_subdir: 'Poddirektorij = korijen instalacije', git_f_kind: 'Vrsta', git_kind_app: 'Aplikacija', git_kind_library: 'Biblioteka', git_kind_tool: 'Alat',
  git_f_include: 'Uključi (globovi)', git_f_exclude: 'Isključi (globovi)', git_f_protected: 'Zaštićeno (nikad se ne prepisuje, po hostu)', git_f_fingerprint: 'Datoteke za prepoznavanje', git_f_hint: 'Nazivi direktorija instalacije',
  git_lines_h: 'Jedno po retku ili odvojeno zarezom', git_suggest_project: 'Projekt', git_suggest_ref: 'Grana (prazno = zadana)', git_suggest_go: 'Predloži', git_repo_files: 'Datoteke repozitorija',
  git_discover: 'Pronađi', git_add_install: 'Dodaj ručno', git_details: 'Detalji', git_forget: 'Zaboravi', git_forget_q: 'Zaboraviti ovu instalaciju? Pretraga je može ponovno pronaći.',
  git_units: 'Jedinice', git_no_units: 'Nijedna systemd ili supervisor jedinica ne spominje ovaj direktorij.', git_diff: 'Razlike', git_diff_title: 'Server → cilj', git_binary: 'Binarna datoteka',
  git_no_diff: 'Nema razlika u tekstu (samo završeci redaka ili razmaci).', git_capped: 'Provjereno je samo prvih 5000 datoteka.', git_files_only: 'Prikaži samo promjene',
  git_sources: 'Git izvori', git_add_source: 'Dodaj izvor', git_kind: 'Vrsta', git_url: 'Adresa', git_token: 'Token samo za čitanje', git_token_keep: 'ostavite prazno za zadržavanje spremljenog tokena',
  git_test: 'Testiraj', git_projects_found: '{n} projekata', git_bundles: 'Offline bundleovi', git_import_bundle: 'Uvezi bundle (.tar.gz)', git_bundle_age: 'izrađen {d}',
  git_bundle_imported: 'Bundle uvezen: {n} datoteka', git_source_mode: 'Ciljevi iz', git_mode_auto: 'Automatski (Git API, inače najnoviji bundle)', git_mode_live: 'Git API',
  git_mode_bundle: 'Offline bundle', git_newest_bundle: 'Najnoviji bundle sa servisom', git_history_depth: 'Dubina povijesti (ranije verzije po datoteci)',
  git_servers: 'Serveri za provjeru', git_servers_h: 'SSH veze. Provjere samo čitaju datoteke, na serverima se ništa ne mijenja.', git_filter: 'Filtar…', git_select_all: 'Sve', git_select_none: 'Ništa',
  git_catalog_settings: 'Postavke kataloga', git_groups: 'Grupe / organizacije', git_fallback: 'Zamjenske grane', git_roots: 'Korijeni na serverima (pretraga do dubine 3)', git_ignore: 'Zanemareni direktoriji (globovi)',
  git_monitor: 'Periodične provjere', git_monitor_on: 'Provjeravaj moje servere svakih {n} minuta dok WRM radi i obavijesti me', git_monitor_off: 'Administrator je isključio periodične provjere.',
  git_monitor_h: 'Kanale za "Git" događaje odaberite u Postavke → Obavijesti. Provjere nikad ništa ne ažuriraju.', git_saved: 'Spremljeno',
  git_n_installs: '{n} instalacija', git_none: '—', git_api_calls: '{n} API poziva', git_from: 'iz', git_bundle_only: 'bundle',
});

(function () {
  const G = {state: null, tab: 'overview', f: {state: '', env: '', folder: '', tag: '', q: ''}, onlyChanges: true, busy: false};
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(b); return s; };
  const lines = s => String(s || '').split(/[\n,]/).map(x => x.trim()).filter(Boolean);
  const json = (method, body) => ({method, headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)});

  const css = `
#git-ws { position: fixed; inset: 0; z-index: 900; background: var(--bg); color: var(--text); display: flex; flex-direction: column; font-size: 13px; }
#git-ws .g-head { display: flex; align-items: center; gap: 10px; padding: 10px 16px; border-bottom: 1px solid var(--border); background: var(--bg2); flex-wrap: wrap; }
#git-ws .g-head h2 { font-size: 15px; margin: 0; white-space: nowrap; }
#git-ws .g-tabs { display: flex; gap: 4px; flex-wrap: wrap; }
#git-ws .g-tabs button { background: transparent; border: 1px solid transparent; color: var(--text2); padding: 6px 12px; border-radius: 8px; cursor: pointer; }
#git-ws .g-tabs button.on { background: var(--accent-d); border-color: var(--accent); color: var(--text); }
#git-ws .g-sp { flex: 1; }
#git-ws .g-body { flex: 1; overflow: auto; padding: 14px 16px 40px; }
#git-ws .g-bar { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; margin: 0 0 12px; }
#git-ws .g-bar select, #git-ws .g-bar input { width: auto; min-width: 120px; flex: 0 1 auto; }
#git-ws .g-card { background: var(--bg2); border: 1px solid var(--border); border-radius: 10px; padding: 12px 14px; margin-bottom: 14px; }
#git-ws .g-card h3 { margin: 0 0 10px; font-size: 13px; }
#git-ws .g-scroll { overflow: auto; max-width: 100%; }
#git-ws table { border-collapse: collapse; width: 100%; font-size: 12px; }
#git-ws th { text-align: left; font-size: 10px; text-transform: uppercase; letter-spacing: .05em; color: var(--text3); padding: 6px 8px; border-bottom: 1px solid var(--border); white-space: nowrap; }
#git-ws td { padding: 6px 8px; border-bottom: 1px solid var(--border); vertical-align: top; }
#git-ws .g-matrix td.g-cell { min-width: 150px; }
#git-ws .g-inst { display: block; border-radius: 7px; padding: 4px 7px; margin: 2px 0; cursor: pointer; border: 1px solid var(--border2); background: var(--bg3); }
#git-ws .g-inst:hover { border-color: var(--accent); }
#git-ws .g-st { display: inline-block; font-size: 10.5px; padding: 1px 7px; border-radius: 10px; font-weight: 600; white-space: nowrap; }
#git-ws .s-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-ws .s-update, #git-ws .s-old, #git-ws .s-missing { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-ws .s-review, #git-ws .s-modified, #git-ws .s-error { background: rgba(239,68,68,.14); color: var(--red); }
#git-ws .s-gone, #git-ws .s-unknown, #git-ws .s-extra, #git-ws .s-protected { background: var(--bg4); color: var(--text2); }
#git-ws .g-env { font-size: 10px; color: var(--purple); margin-left: 4px; }
#git-ws .g-mut { color: var(--text3); font-size: 11px; }
#git-ws .g-mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; word-break: break-all; }
#git-ws .g-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 10px; }
#git-ws .g-grid label { display: block; font-size: 11px; color: var(--text2); margin-bottom: 3px; }
#git-ws textarea { width: 100%; min-height: 54px; font-family: inherit; }
#git-ws .g-conns { max-height: 260px; overflow: auto; border: 1px solid var(--border); border-radius: 8px; padding: 6px; }
#git-ws .g-conns label { display: flex; gap: 6px; align-items: center; padding: 2px 4px; }
#git-ws .g-conns input { width: auto; }
#git-modal { position: fixed; inset: 0; z-index: 950; background: rgba(0,0,0,.55); display: flex; align-items: flex-start; justify-content: center; padding: 30px 12px; overflow: auto; }
#git-modal .g-dlg { background: var(--bg2); border: 1px solid var(--border2); border-radius: 12px; width: min(980px, 100%); padding: 16px 18px; }
#git-modal .g-dlg h3 { margin: 0 0 10px; font-size: 14px; word-break: break-all; }
#git-modal .g-diff { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11.5px; background: var(--bg); border: 1px solid var(--border); border-radius: 8px; overflow: auto; max-height: 60vh; }
#git-modal .g-diff div { white-space: pre; padding: 0 8px; }
#git-modal .g-diff .d-add { background: rgba(34,197,94,.13); color: #bbf7d0; }
#git-modal .g-diff .d-del { background: rgba(239,68,68,.13); color: #fecaca; }
#git-modal .g-diff .d-hunk { background: var(--accent-d); color: var(--text2); }
#git-modal table { border-collapse: collapse; width: 100%; font-size: 12px; }
#git-modal td, #git-modal th { padding: 4px 6px; border-bottom: 1px solid var(--border); text-align: left; }
#git-modal .g-st { display: inline-block; font-size: 10.5px; padding: 1px 7px; border-radius: 10px; font-weight: 600; }
#git-modal .s-ok { background: rgba(34,197,94,.14); color: var(--green); }
#git-modal .s-old, #git-modal .s-missing, #git-modal .s-update { background: rgba(245,158,11,.14); color: var(--yellow); }
#git-modal .s-modified, #git-modal .s-review, #git-modal .s-error { background: rgba(239,68,68,.14); color: var(--red); }
#git-modal .s-extra, #git-modal .s-protected, #git-modal .s-gone, #git-modal .s-unknown { background: var(--bg4); color: var(--text2); }
#git-modal .g-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 10px; }
#git-modal .g-grid label { display: block; font-size: 11px; color: var(--text2); margin-bottom: 3px; }
#git-modal textarea { width: 100%; min-height: 48px; font-family: inherit; }
#git-modal .g-btns { display: flex; gap: 8px; justify-content: flex-end; margin-top: 12px; flex-wrap: wrap; }
#git-ws input:not([type=checkbox]), #git-ws select, #git-ws textarea, #git-modal input:not([type=checkbox]), #git-modal select, #git-modal textarea {
  background: var(--bg3); color: var(--text); border: 1px solid var(--border2); border-radius: 7px; padding: 6px 9px; font-size: 12.5px; box-sizing: border-box; outline: none; }
#git-ws .g-grid input, #git-ws .g-grid select, #git-modal .g-grid input, #git-modal .g-grid select { width: 100%; }
#git-ws input:focus, #git-ws select:focus, #git-ws textarea:focus, #git-modal input:focus, #git-modal select:focus, #git-modal textarea:focus { border-color: var(--accent); }
@media (max-width: 640px) { #git-ws .g-head { padding: 8px 10px; } #git-ws .g-body { padding: 10px; } #git-ws .g-bar select, #git-ws .g-bar input { min-width: 0; flex: 1 1 140px; } }
`;

  function stBadge(s) { return `<span class="g-st s-${esc(s)}">${esc(t('git_st_' + s))}</span>`; }
  function fsBadge(s) { return `<span class="g-st s-${esc(s)}">${esc(t('git_fs_' + s))}</span>`; }

  async function load() {
    const r = await fetch('/api/git');
    if (!r.ok) { showToast(await apiError(r), 'error'); return false; }
    G.state = await r.json();
    return true;
  }

  function open() {
    if (!$('git-style')) { const s = document.createElement('style'); s.id = 'git-style'; s.textContent = css; document.head.appendChild(s); }
    let ws = $('git-ws');
    if (!ws) {
      ws = document.createElement('div');
      ws.id = 'git-ws';
      ws.innerHTML = `<div class="g-head"><h2>⎇ ${esc(t('git_title'))}</h2><div class="g-tabs" id="git-tabs"></div><span class="g-sp"></span>
        <span class="g-mut" id="git-last"></span><button id="git-check" onclick="gitWorkspace.check()">${esc(t('git_check_now'))}</button>
        <button class="btn-sec" onclick="gitWorkspace.close()" title="${esc(t('git_close'))}">✕</button></div><div class="g-body" id="git-body"></div>`;
      document.body.appendChild(ws);
      document.addEventListener('keydown', e => { if (e.key === 'Escape' && $('git-ws') && !$('git-modal')) close(); });
    }
    ws.style.display = 'flex';
    load().then(ok => ok && render());
  }

  function close() { const ws = $('git-ws'); if (ws) ws.style.display = 'none'; closeModal(); }

  function render() {
    const s = G.state;
    $('git-tabs').innerHTML = ['overview', 'services', 'installs', 'settings'].map(k =>
      `<button class="${G.tab === k ? 'on' : ''}" onclick="gitWorkspace.tab('${k}')">${esc(t('git_tab_' + k))}</button>`).join('');
    $('git-last').textContent = `${t('git_last_check')}: ${s.last_check_at ? fmtTime(s.last_check_at) : t('git_never')}`;
    $('git-check').disabled = !s.can_check || G.busy;
    $('git-check').title = s.can_check ? '' : t('git_no_checks');
    $('git-check').textContent = G.busy ? t('git_checking') : t('git_check_now');
    const body = $('git-body');
    if (G.tab === 'overview') body.innerHTML = overviewHTML();
    else if (G.tab === 'services') body.innerHTML = servicesHTML();
    else if (G.tab === 'installs') body.innerHTML = installsHTML();
    else body.innerHTML = settingsHTML();
  }

  function setTab(k) { G.tab = k; render(); }

  // ── overview ──
  function filtered() {
    const f = G.f, q = f.q.toLowerCase();
    return G.state.installs.filter(i => (!f.state || i.state === f.state) && (!f.env || i.env === f.env) &&
      (!f.folder || String(i.folder_id || '') === f.folder) && (!f.tag || (i.tags || []).includes(f.tag)) &&
      (!q || (i.conn_name + ' ' + i.app + ' ' + i.path).toLowerCase().includes(q)));
  }

  function filterBar() {
    const s = G.state, f = G.f;
    const envs = [...new Set(s.installs.map(i => i.env).filter(Boolean))].sort();
    const fids = [...new Set(s.installs.map(i => i.folder_id).filter(Boolean))];
    const tags = [...new Set(s.installs.flatMap(i => i.tags || []))].sort();
    const opt = (v, label, cur) => `<option value="${esc(v)}" ${String(cur) === String(v) ? 'selected' : ''}>${esc(label)}</option>`;
    const fname = id => ((folders || []).find(x => x.id === id) || {}).name || ('#' + id);
    return `<select onchange="gitWorkspace.filter('state', this.value)">${opt('', t('git_all_states'), f.state)}${['update', 'review', 'ok', 'error', 'gone', 'unknown'].map(x => opt(x, t('git_st_' + x), f.state)).join('')}</select>
      <select onchange="gitWorkspace.filter('env', this.value)">${opt('', t('git_all_envs'), f.env)}${envs.map(x => opt(x, x, f.env)).join('')}</select>
      <select onchange="gitWorkspace.filter('folder', this.value)">${opt('', t('git_all_folders'), f.folder)}${fids.map(x => opt(x, fname(x), f.folder)).join('')}</select>
      <select onchange="gitWorkspace.filter('tag', this.value)">${opt('', t('git_all_tags'), f.tag)}${tags.map(x => opt(x, x, f.tag)).join('')}</select>
      <input type="search" placeholder="${esc(t('git_search'))}" value="${esc(f.q)}" oninput="gitWorkspace.filter('q', this.value, true)">`;
  }

  function refBar() {
    const st = G.state.settings;
    return `<label class="g-mut">${esc(t('git_ref_choice'))}</label>
      <select onchange="gitWorkspace.setRef(this.value)">${['catalog', 'default', 'branch'].map(m => `<option value="${m}" ${st.ref_mode === m ? 'selected' : ''}>${esc(t('git_ref_' + m))}</option>`).join('')}</select>
      ${st.ref_mode === 'branch' ? `<input id="git-ref-branch" value="${esc(st.ref_branch || '')}" placeholder="main" onchange="gitWorkspace.setRef('branch')" style="min-width:110px">` : ''}`;
  }

  function overviewHTML() {
    const s = G.state;
    const apps = [...s.catalog.apps.map(a => a.name), ...(s.bundle_apps || [])];
    const list = filtered();
    const targets = Object.fromEntries(s.targets.map(x => [x.app, x]));
    let html = `<div class="g-bar">${refBar()}<button class="btn-sec btn-sm" onclick="gitWorkspace.refresh()" ${s.can_check ? '' : 'disabled'}>↻ ${esc(t('git_refresh_targets'))}</button></div>`;
    html += `<div class="g-card"><h3>${esc(t('git_targets'))}</h3><div class="g-scroll"><table><tr><th>${esc(t('git_service'))}</th><th>${esc(t('git_target'))}</th><th>${esc(t('git_latest_tag'))}</th><th>${esc(t('git_source'))}</th><th>${esc(t('git_warnings'))}</th></tr>` +
      (apps.length ? apps.map(a => { const x = targets[a]; return `<tr><td><b>${esc(a)}</b></td><td>${x ? `<span class="g-mono">${esc(x.version || '—')}</span> <span class="g-mut">${esc(x.branch || '')} · ${esc(x.commit_date || '')}</span>` : `<span class="g-mut">${esc(t('git_no_target'))}</span>`}</td>
        <td class="g-mono">${esc((x && x.latest_tag) || '—')}</td><td class="g-mut">${esc(x ? x.source : '')}</td><td>${x && x.error ? `<span class="g-st s-error">${esc(x.error)}</span>` : esc(((x && x.warnings) || []).join('; '))}</td></tr>`; }).join('')
        : `<tr><td colspan="5" class="g-mut">${esc(t('git_empty_overview'))}</td></tr>`) + '</table></div></div>';
    html += `<div class="g-bar">${filterBar()}</div>`;
    const servers = [...new Map(list.map(i => [i.conn_id, i.conn_name])).entries()].sort((a, b) => a[1].localeCompare(b[1]));
    if (!servers.length) return html + `<div class="hint-box">${esc(t('git_empty_overview'))}</div>`;
    const cols = apps.filter(a => list.some(i => i.app === a));
    html += `<div class="g-card g-scroll"><table class="g-matrix"><tr><th>${esc(t('git_server'))}</th>${cols.map(a => `<th>${esc(a)}</th>`).join('')}</tr>` +
      servers.map(([cid, name]) => `<tr><td><b>${esc(name)}</b></td>${cols.map(a => `<td class="g-cell">${list.filter(i => i.conn_id === cid && i.app === a).map(instChip).join('')}</td>`).join('')}</tr>`).join('') + '</table></div>';
    return html;
  }

  function instChip(i) {
    const v = (i.version_md && i.version_md.version) || '';
    return `<span class="g-inst" onclick="gitWorkspace.details(${i.id})" title="${esc(i.path)}">${stBadge(i.state)}${i.env ? `<span class="g-env">${esc(i.env)}</span>` : ''}
      <div class="g-mut">${esc(t('git_version_md'))}: <span class="g-mono">${v ? esc(v) : esc(t('git_vmd_missing'))}</span>${i.behind ? ' · ' + esc(tf('git_behind', {n: i.behind})) : ''}</div></span>`;
  }

  // ── services ──
  function servicesHTML() {
    const s = G.state;
    const targets = Object.fromEntries(s.targets.map(x => [x.app, x]));
    const row = (a, fromBundle) => { const x = targets[a.name || a]; const name = a.name || a;
      return `<tr><td><input type="checkbox" class="git-sel" value="${esc(name)}" ${x && !x.error ? '' : 'disabled'} style="width:auto"></td><td><b>${esc(name)}</b>${fromBundle ? ` <span class="g-mut">(${esc(t('git_bundle_only'))})</span>` : ''}</td>
        <td class="g-mono">${esc(a.project || (x && x.project) || '')}</td><td class="g-mono">${esc(a.ref || (x && x.ref) || '')}${a.subdir ? ' · ' + esc(a.subdir) + '/' : ''}</td>
        <td>${x ? `<span class="g-mono">${esc(x.version)}</span> <span class="g-mut">${x.files} ${esc(t('git_files').toLowerCase())}</span>` : '<span class="g-mut">—</span>'}</td>
        <td style="white-space:nowrap">${fromBundle ? '' : `<button class="btn-sec btn-sm" onclick="gitWorkspace.editApp('${esc(name)}')">${esc(t('git_edit'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.delApp('${esc(name)}')">✕</button>`}</td></tr>`; };
    return `<div class="g-bar"><button onclick="gitWorkspace.editApp('')">＋ ${esc(t('git_add_service'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.suggest()">✨ ${esc(t('git_suggest'))}</button>
      <button class="btn-sec" onclick="document.getElementById('git-cat-file').click()">⬆ ${esc(t('git_import_catalog'))}</button><input type="file" id="git-cat-file" accept=".json,application/json" style="display:none" onchange="gitWorkspace.importCatalog(this)">
      <button class="btn-sec" onclick="gitWorkspace.exportCatalog()">⬇ ${esc(t('git_export_catalog'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.exportBundle()">📦 ${esc(t('git_export_bundle'))}</button>
      <button class="btn-sec" onclick="gitWorkspace.refresh()" ${s.can_check ? '' : 'disabled'}>↻ ${esc(t('git_refresh_targets'))}</button></div>
      <div class="g-card g-scroll"><table><tr><th></th><th>${esc(t('git_f_name'))}</th><th>${esc(t('git_f_project'))}</th><th>${esc(t('git_f_ref'))}</th><th>${esc(t('git_target'))}</th><th></th></tr>
      ${s.catalog.apps.map(a => row(a, false)).join('')}${(s.bundle_apps || []).map(a => row(a, true)).join('')}</table></div>`;
  }

  function appForm(a) {
    const f = (k, label, v, ph) => `<div><label>${esc(t(label))}</label><input id="ga-${k}" value="${esc(v || '')}" placeholder="${esc(ph || '')}"></div>`;
    const ta = (k, label, v) => `<div><label>${esc(t(label))}</label><textarea id="ga-${k}" placeholder="${esc(t('git_lines_h'))}">${esc((v || []).join('\n'))}</textarea></div>`;
    return `<div class="g-grid">${f('name', 'git_f_name', a.name)}${f('project', 'git_f_project', a.project, 'group/project')}${f('ref', 'git_f_ref', a.ref || 'tag:latest', t('git_f_ref_h'))}
      ${f('branch', 'git_f_branch', a.branch, 'main')}${f('tag_filter', 'git_f_tag_filter', a.tag_filter, '^v\\d')}${f('subdir', 'git_f_subdir', a.subdir, 'linux')}
      <div><label>${esc(t('git_f_kind'))}</label><select id="ga-kind">${['app', 'library', 'tool'].map(k => `<option value="${k}" ${a.kind === k ? 'selected' : ''}>${esc(t('git_kind_' + k))}</option>`).join('')}</select></div></div>
      <div class="g-grid" style="margin-top:10px">${ta('fingerprint', 'git_f_fingerprint', a.fingerprint)}${ta('protected', 'git_f_protected', a.protected)}${ta('install_hint', 'git_f_hint', a.install_hint)}
      ${ta('include', 'git_f_include', a.include)}${ta('exclude', 'git_f_exclude', a.exclude)}</div>`;
  }

  function readAppForm() {
    const v = k => $('ga-' + k).value.trim();
    return {name: v('name'), project: v('project'), ref: v('ref'), branch: v('branch'), tag_filter: v('tag_filter'), subdir: v('subdir'), kind: v('kind'),
      fingerprint: lines(v('fingerprint')), protected: lines(v('protected')), install_hint: lines(v('install_hint')), include: lines(v('include')), exclude: lines(v('exclude'))};
  }

  function editApp(name, preset, files) {
    const a = preset || G.state.catalog.apps.find(x => x.name === name) || {kind: 'app', ref: 'tag:latest', protected: ['config*', '*.ini', '.env']};
    modal(`<h3>${esc(name || t('git_add_service'))}</h3>${appForm(a)}${files ? `<details style="margin-top:10px"><summary class="g-mut">${esc(t('git_repo_files'))} (${files.length})</summary><div class="g-mono" style="max-height:200px;overflow:auto">${files.map(esc).join('<br>')}</div></details>` : ''}
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="ga-save">${esc(t('git_save'))}</button></div>`);
    $('ga-save').onclick = async () => {
      const app = readAppForm();
      const r = await fetch('/api/git/catalog/apps/' + encodeURIComponent(name || app.name), json('PUT', app));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); closeModal(); render(); showToast(t('git_saved'), 'success');
    };
  }

  async function delApp(name) {
    if (!(await uiConfirm(tf('git_delete_q', {name}), {danger: true, okText: t('git_delete')}))) return;
    const r = await fetch('/api/git/catalog/apps/' + encodeURIComponent(name), {method: 'DELETE'});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render();
  }

  async function suggest() {
    const src = G.state.sources.find(x => x.kind !== 'bundle');
    let projects = [];
    if (src) { const r = await fetch(`/api/git/sources/${src.id}/test`, {method: 'POST'}); if (r.ok) projects = (await r.json()).projects || []; }
    modal(`<h3>✨ ${esc(t('git_suggest'))}</h3><div class="g-grid"><div><label>${esc(t('git_suggest_project'))}</label><input id="gs-project" list="gs-projects" placeholder="group/project">
      <datalist id="gs-projects">${projects.map(p => `<option value="${esc(p)}">`).join('')}</datalist></div>
      <div><label>${esc(t('git_suggest_ref'))}</label><input id="gs-ref"></div></div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gs-go">${esc(t('git_suggest_go'))}</button></div>`);
    $('gs-go').onclick = async () => {
      $('gs-go').disabled = true;
      const r = await fetch('/api/git/suggest', json('POST', {project: $('gs-project').value.trim(), ref: $('gs-ref').value.trim()}));
      $('gs-go').disabled = false;
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      const j = await r.json();
      editApp('', j.app, j.files);
    };
  }

  function importCatalog(input) {
    const file = input.files[0]; input.value = '';
    if (!file) return;
    file.text().then(async text => {
      let body;
      try { body = JSON.parse(text); } catch (e) { showToast(e.message, 'error'); return; }
      const merge = await uiConfirm(t('git_import_merge'), {okText: 'OK'});
      const r = await fetch('/api/git/catalog/import?mode=' + (merge ? 'merge' : 'replace'), json('POST', body));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); render(); showToast(t('git_saved'), 'success');
    });
  }

  async function download(r, fallback) {
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const m = /filename="([^"]+)"/.exec(r.headers.get('Content-Disposition') || '');
    const url = URL.createObjectURL(await r.blob());
    const a = document.createElement('a'); a.href = url; a.download = m ? m[1] : fallback; document.body.appendChild(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 2000);
  }

  function exportCatalog() { fetch('/api/git/catalog/export').then(r => download(r, 'catalog.json')); }

  function exportBundle() {
    const apps = [...document.querySelectorAll('.git-sel:checked')].map(x => x.value);
    if (!apps.length) { showToast(t('git_export_bundle_none'), 'warning'); return; }
    fetch('/api/git/bundles/export', json('POST', {apps})).then(r => download(r, 'bundle.tar.gz'));
  }

  // ── installations ──
  function installsHTML() {
    const list = filtered();
    const sshConns = (conns || []).filter(c => c.protocol === 'SSH').sort((a, b) => a.name.localeCompare(b.name));
    const apps = [...G.state.catalog.apps.map(a => a.name), ...(G.state.bundle_apps || [])];
    return `<div class="g-bar">${filterBar()}<button class="btn-sec" onclick="gitWorkspace.check(true)" ${G.state.can_check ? '' : 'disabled'}>🔎 ${esc(t('git_discover'))}</button></div>
      <div class="g-card g-scroll"><table><tr><th>${esc(t('git_server'))}</th><th>${esc(t('git_service'))}</th><th>${esc(t('git_path'))}</th><th>${esc(t('git_env'))}</th><th>${esc(t('git_state'))}</th>
      <th>${esc(t('git_version_md'))}</th><th>${esc(t('git_target'))}</th><th>${esc(t('git_checked'))}</th><th></th></tr>
      ${list.map(i => `<tr><td>${esc(i.conn_name)}</td><td><b>${esc(i.app)}</b></td><td class="g-mono">${esc(i.path)}</td><td>${esc(i.env || '')}</td>
        <td>${stBadge(i.state)}${i.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: i.behind}))}</span>` : ''}${i.error ? `<div class="g-mut">${esc(i.error)}</div>` : ''}</td>
        <td class="g-mono">${esc((i.version_md && i.version_md.version) || t('git_vmd_missing'))}</td><td class="g-mono">${esc(i.target || '—')}</td><td class="g-mut">${esc(i.checked_at ? fmtTime(i.checked_at) : '—')}</td>
        <td style="white-space:nowrap"><button class="btn-sec btn-sm" onclick="gitWorkspace.details(${i.id})">${esc(t('git_details'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.forget(${i.id})">✕</button></td></tr>`).join('')}</table></div>
      <div class="g-card"><h3>${esc(t('git_add_install'))}</h3><div class="g-bar"><select id="gi-conn">${sshConns.map(c => `<option value="${c.id}">${esc(c.name)}</option>`).join('')}</select>
      <select id="gi-app">${apps.map(a => `<option>${esc(a)}</option>`).join('')}</select><input id="gi-path" placeholder="/opt/service" style="min-width:200px"><input id="gi-env" placeholder="${esc(t('git_env'))}">
      <button class="btn-sec" onclick="gitWorkspace.addInstall()">＋</button></div></div>`;
  }

  async function addInstall() {
    const r = await fetch('/api/git/installs', json('POST', {conn_id: Number($('gi-conn').value), app: $('gi-app').value, path: $('gi-path').value.trim(), env: $('gi-env').value.trim()}));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render();
  }

  async function forget(id) {
    if (!(await uiConfirm(t('git_forget_q'), {danger: true, okText: t('git_forget')}))) return;
    const r = await fetch('/api/git/installs/' + id, {method: 'DELETE'});
    if (r.ok) { G.state = await r.json(); render(); }
  }

  async function details(id) {
    const r = await fetch('/api/git/installs/' + id);
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const i = await r.json();
    G.detail = i;
    const v = i.version_md || {};
    const files = (i.files || []).filter(f => !G.onlyChanges || (f.state !== 'ok' && f.state !== 'protected'));
    modal(`<h3>${esc(i.conn_name)} · ${esc(i.app)} · <span class="g-mono">${esc(i.path)}</span></h3>
      <div class="g-bar">${stBadge(i.state)}${i.env ? `<span class="g-env">${esc(i.env)}</span>` : ''}${Object.entries(i.counts || {}).map(([k, n]) => `${fsBadge(k)} ${n}`).join(' ')}</div>
      ${i.error ? `<div class="hint-box bad">${esc(i.error)}</div>` : ''}${i.capped ? `<div class="hint-box warn">${esc(t('git_capped'))}</div>` : ''}
      <div class="g-grid"><div><label>${esc(t('git_version_md'))}</label><span class="g-mono">${v.version ? esc(v.version) : esc(t('git_vmd_missing'))}</span>
        ${v.ref ? `<div class="g-mut">${esc(v.ref)} · ${esc(v.commit || '')}</div>` : ''}${v.updated ? `<div class="g-mut">${esc(v.updated)} ${esc(v.user || '')}${v.bundle ? ' · ' + esc(v.bundle) : ''}</div>` : ''}</div>
        <div><label>${esc(t('git_target'))}</label><span class="g-mono">${esc(i.target || '—')}</span></div>
        <div><label>${esc(t('git_units'))}</label>${(i.units || []).length ? i.units.map(u => `<div class="g-mono">${esc(u.kind)}: ${esc(u.name)}</div>`).join('') : `<span class="g-mut">${esc(t('git_no_units'))}</span>`}</div></div>
      <label style="display:flex;gap:6px;align-items:center;margin:12px 0 6px"><input type="checkbox" style="width:auto" ${G.onlyChanges ? 'checked' : ''} onchange="gitWorkspace.toggleChanges(this.checked)">${esc(t('git_files_only'))}</label>
      <div style="max-height:45vh;overflow:auto"><table><tr><th>${esc(t('git_files'))}</th><th>${esc(t('git_state'))}</th><th></th></tr>
      ${files.map(f => `<tr><td class="g-mono">${esc(f.path)}</td><td>${fsBadge(f.state)}${f.behind ? ` <span class="g-mut">${esc(tf('git_behind', {n: f.behind}))} (${esc(f.tag || '')} ${esc(f.date || '')})</span>` : ''}${f.note ? ` <span class="g-mut">${esc(f.note)}</span>` : ''}</td>
        <td>${f.diffable && f.state !== 'ok' && G.state.can_check ? `<button class="btn-sec btn-sm" onclick="gitWorkspace.diff(${i.id}, '${esc(encodeURIComponent(f.path))}')">${esc(t('git_diff'))}</button>` : ''}</td></tr>`).join('')}</table></div>
      <div id="git-diff"></div><div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_close'))}</button></div>`);
  }

  async function diff(id, path) {
    const out = $('git-diff');
    out.innerHTML = `<div class="hint-box">⏳</div>`;
    const r = await fetch(`/api/git/installs/${id}/diff?path=${path}`);
    if (!r.ok) { out.innerHTML = `<div class="hint-box bad">${esc(await apiError(r))}</div>`; return; }
    const d = await r.json();
    if (d.binary) { out.innerHTML = `<div class="hint-box">${esc(t('git_binary'))}</div>`; return; }
    if (!d.lines.length) { out.innerHTML = `<div class="hint-box">${esc(t('git_no_diff'))}</div>`; return; }
    out.innerHTML = `<h3 style="margin-top:12px">${esc(t('git_diff_title'))}: <span class="g-mono">${esc(d.path)}</span> → ${esc(d.target || '')}</h3><div class="g-diff">` +
      d.lines.map(l => l.op === '@' ? `<div class="d-hunk">@@ -${l.a} +${l.b} @@</div>` : `<div class="${l.op === '+' ? 'd-add' : l.op === '-' ? 'd-del' : ''}">${esc(l.op + ' ' + l.s)}</div>`).join('') + '</div>';
    out.scrollIntoView({block: 'nearest'});
  }

  // ── settings ──
  function settingsHTML() {
    const s = G.state, st = s.settings, cat = s.catalog;
    const apiSources = s.sources.filter(x => x.kind !== 'bundle'), bundles = s.sources.filter(x => x.kind === 'bundle');
    const sshConns = (conns || []).filter(c => c.protocol === 'SSH').sort((a, b) => a.name.localeCompare(b.name));
    const sel = new Set(st.conn_ids || []);
    return `<div class="g-card"><h3>${esc(t('git_sources'))}</h3>
      ${apiSources.length ? `<table>${apiSources.map(x => `<tr><td><b>${esc(x.name)}</b> <span class="g-mut">${esc(x.kind)}</span><div class="g-mono g-mut">${esc(x.url)}</div>${x.last_error ? `<div class="g-st s-error">${esc(x.last_error)}</div>` : ''}</td>
        <td style="white-space:nowrap"><button class="btn-sec btn-sm" onclick="gitWorkspace.testSource(${x.id})">${esc(t('git_test'))}</button> <button class="btn-sec btn-sm" onclick="gitWorkspace.editSource(${x.id})">${esc(t('git_edit'))}</button> <button class="btn-danger btn-sm" onclick="gitWorkspace.delSource(${x.id})">✕</button></td></tr>`).join('')}</table>` : ''}
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec" onclick="gitWorkspace.editSource(0)">＋ ${esc(t('git_add_source'))}</button></div></div>
      <div class="g-card"><h3>${esc(t('git_bundles'))}</h3>
      ${bundles.length ? `<table>${bundles.map(x => `<tr><td><b>${esc(x.bundle_id || x.name)}</b> <span class="g-mut">${esc(tf('git_bundle_age', {d: x.created || '?'}))} · ${esc(x.created_by || '')}</span><div class="g-mono g-mut">${esc(x.origin || '')}</div><div class="g-mut">${esc((x.apps || []).join(', '))}</div></td>
        <td><button class="btn-danger btn-sm" onclick="gitWorkspace.delSource(${x.id})">✕</button></td></tr>`).join('')}</table>` : ''}
      <div class="g-bar" style="margin-top:8px"><button class="btn-sec" onclick="document.getElementById('git-bundle-file').click()">⬆ ${esc(t('git_import_bundle'))}</button><input type="file" id="git-bundle-file" accept=".gz,.tgz,application/gzip" style="display:none" onchange="gitWorkspace.importBundle(this)"></div></div>
      <div class="g-card"><h3>${esc(t('git_source_mode'))}</h3><div class="g-grid">
        <div><label>${esc(t('git_source_mode'))}</label><select id="gx-mode">${['auto', 'live', 'bundle'].map(m => `<option value="${m}" ${st.mode === m ? 'selected' : ''}>${esc(t('git_mode_' + m))}</option>`).join('')}</select></div>
        <div><label>${esc(t('git_mode_bundle'))}</label><select id="gx-bundle"><option value="0">${esc(t('git_newest_bundle'))}</option>${bundles.map(b => `<option value="${b.id}" ${st.bundle_id === b.id ? 'selected' : ''}>${esc(b.bundle_id || b.name)}</option>`).join('')}</select></div>
        ${apiSources.length ? `<div><label>${esc(t('git_source'))}</label><select id="gx-source">${apiSources.map(x => `<option value="${x.id}" ${st.source_id === x.id ? 'selected' : ''}>${esc(x.name)}</option>`).join('')}</select></div>` : ''}
        <div><label>${esc(t('git_history_depth'))}</label><input id="gx-depth" type="number" min="1" max="50" value="${st.history_depth}"></div></div></div>
      <div class="g-card"><h3>${esc(t('git_servers'))}</h3><div class="g-mut" style="margin-bottom:6px">${esc(t('git_servers_h'))}</div>
        <div class="g-bar"><input type="search" placeholder="${esc(t('git_filter'))}" oninput="gitWorkspace.filterConns(this.value)">
        <button class="btn-sec btn-sm" onclick="gitWorkspace.selConns(true)">${esc(t('git_select_all'))}</button><button class="btn-sec btn-sm" onclick="gitWorkspace.selConns(false)">${esc(t('git_select_none'))}</button></div>
        <div class="g-conns" id="gx-conns">${sshConns.map(c => `<label data-q="${esc((c.name + ' ' + c.host + ' ' + (c.tags || []).join(' ')).toLowerCase())}"><input type="checkbox" value="${c.id}" ${sel.has(c.id) ? 'checked' : ''}>${esc(c.name)} <span class="g-mut">${esc(c.host)}</span></label>`).join('')}</div></div>
      <div class="g-card"><h3>${esc(t('git_monitor'))}</h3>${s.interval_minutes > 0 ? `<label style="display:flex;gap:6px;align-items:center"><input type="checkbox" id="gx-monitor" style="width:auto" ${st.monitor ? 'checked' : ''}>${esc(tf('git_monitor_on', {n: s.interval_minutes}))}</label>` : `<div class="g-mut">${esc(t('git_monitor_off'))}</div>`}
        <div class="g-mut" style="margin-top:6px">${esc(t('git_monitor_h'))}</div></div>
      <div class="g-bar"><button onclick="gitWorkspace.saveSettings()">${esc(t('git_save'))}</button></div>
      <div class="g-card"><h3>${esc(t('git_catalog_settings'))}</h3><div class="g-grid">
        <div><label>${esc(t('git_url'))}</label><input id="gc-url" value="${esc(cat.gitlab.url || '')}" placeholder="https://git.example.com"></div>
        <div><label>${esc(t('git_groups'))}</label><textarea id="gc-groups">${esc((cat.gitlab.groups || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_fallback'))}</label><textarea id="gc-fallback">${esc((cat.fallback_branches || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_roots'))}</label><textarea id="gc-roots">${esc((cat.server_roots || []).join('\n'))}</textarea></div>
        <div><label>${esc(t('git_ignore'))}</label><textarea id="gc-ignore">${esc((cat.ignore_dirs || []).join('\n'))}</textarea></div></div>
        <div class="g-bar" style="margin-top:8px"><button onclick="gitWorkspace.saveCatalogSettings()">${esc(t('git_save'))}</button></div></div>`;
  }

  function filterConns(q) { q = q.toLowerCase(); document.querySelectorAll('#gx-conns label').forEach(l => { l.style.display = l.dataset.q.includes(q) ? '' : 'none'; }); }
  function selConns(on) { document.querySelectorAll('#gx-conns label').forEach(l => { if (l.style.display !== 'none') l.querySelector('input').checked = on; }); }

  async function putSettings(patch) {
    const body = Object.assign({}, G.state.settings, patch);
    const r = await fetch('/api/git/settings', json('PUT', body));
    if (!r.ok) { showToast(await apiError(r), 'error'); return false; }
    G.state = await r.json(); render(); return true;
  }

  function saveSettings() {
    const mon = $('gx-monitor');
    putSettings({mode: $('gx-mode').value, bundle_id: Number($('gx-bundle').value) || 0, source_id: Number(($('gx-source') || {}).value) || 0,
      history_depth: Number($('gx-depth').value) || 0, conn_ids: [...document.querySelectorAll('#gx-conns input:checked')].map(x => Number(x.value)),
      monitor: mon ? mon.checked : G.state.settings.monitor}).then(ok => ok && showToast(t('git_saved'), 'success'));
  }

  async function saveCatalogSettings() {
    const cat = Object.assign({}, G.state.catalog, {gitlab: {url: $('gc-url').value.trim(), groups: lines($('gc-groups').value)},
      fallback_branches: lines($('gc-fallback').value), server_roots: lines($('gc-roots').value), ignore_dirs: lines($('gc-ignore').value)});
    const r = await fetch('/api/git/catalog', json('PUT', cat));
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    G.state = await r.json(); render(); showToast(t('git_saved'), 'success');
  }

  function setRef(mode) {
    const b = $('git-ref-branch');
    if (mode === 'branch' && !(b && b.value.trim())) { G.state.settings.ref_mode = 'branch'; render(); setTimeout(() => $('git-ref-branch') && $('git-ref-branch').focus(), 30); return; }
    putSettings({ref_mode: mode, ref_branch: b ? b.value.trim() : G.state.settings.ref_branch});
  }

  function editSource(id) {
    const x = G.state.sources.find(s => s.id === id) || {kind: 'gitlab', name: '', url: G.state.catalog.gitlab.url || ''};
    modal(`<h3>${esc(id ? x.name : t('git_add_source'))}</h3><div class="g-grid">
      <div><label>${esc(t('git_kind'))}</label><select id="gsrc-kind"><option value="gitlab" ${x.kind === 'gitlab' ? 'selected' : ''}>GitLab</option><option value="github" ${x.kind === 'github' ? 'selected' : ''}>GitHub</option></select></div>
      <div><label>${esc(t('git_f_name'))}</label><input id="gsrc-name" value="${esc(x.name)}"></div>
      <div><label>${esc(t('git_url'))}</label><input id="gsrc-url" value="${esc(x.url)}" placeholder="https://git.example.com"></div>
      <div><label>${esc(t('git_token'))}</label><input id="gsrc-token" type="password" autocomplete="new-password" placeholder="${id && x.has_token ? esc(t('git_token_keep')) : ''}"></div></div>
      <div class="g-btns"><button class="btn-sec" onclick="gitWorkspace.closeModal()">${esc(t('git_cancel'))}</button><button id="gsrc-save">${esc(t('git_save'))}</button></div>`);
    $('gsrc-save').onclick = async () => {
      const body = {kind: $('gsrc-kind').value, name: $('gsrc-name').value.trim(), url: $('gsrc-url').value.trim()};
      const tok = $('gsrc-token').value;
      if (tok || !id) body.token = tok;
      const r = await fetch(id ? '/api/git/sources/' + id : '/api/git/sources', json(id ? 'PUT' : 'POST', body));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = await r.json(); closeModal(); render();
    };
  }

  async function testSource(id) {
    const r = await fetch(`/api/git/sources/${id}/test`, {method: 'POST'});
    if (!r.ok) showToast(await apiError(r), 'error');
    else showToast(tf('git_projects_found', {n: ((await r.json()).projects || []).length}), 'success');
    load().then(ok => ok && render());
  }

  async function delSource(id) {
    if (!(await uiConfirm(t('git_delete') + '?', {danger: true, okText: t('git_delete')}))) return;
    const r = await fetch('/api/git/sources/' + id, {method: 'DELETE'});
    if (r.ok) { G.state = await r.json(); render(); }
  }

  async function importBundle(input) {
    const file = input.files[0]; input.value = '';
    if (!file) return;
    const r = await fetch('/api/git/bundles', {method: 'POST', headers: {'Content-Type': 'application/gzip'}, body: file});
    if (!r.ok) { showToast(await apiError(r), 'error'); return; }
    const j = await r.json();
    G.state = j.state; render();
    showToast(tf('git_bundle_imported', {n: j.result.files}) + (j.result.warnings.length ? ' · ⚠ ' + j.result.warnings[0] : ''), j.result.warnings.length ? 'warning' : 'success', 5000);
  }

  // ── runs ──
  async function check(discoverOnly) {
    if (G.busy) return;
    G.busy = true; render();
    try {
      const r = await fetch('/api/git/check', json('POST', {discover: true, refresh: !discoverOnly}));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      const j = await r.json();
      G.state = j.state;
      const errs = j.result.servers.filter(s => s.error);
      if (errs.length) showToast(errs.map(s => s.name + ': ' + s.error).join('\n'), 'warning', 7000);
    } finally { G.busy = false; render(); }
  }

  async function refresh() {
    if (G.busy) return;
    G.busy = true; render();
    try {
      const r = await fetch('/api/git/targets/refresh', json('POST', {}));
      if (!r.ok) { showToast(await apiError(r), 'error'); return; }
      G.state = (await r.json()).state;
    } finally { G.busy = false; render(); }
  }

  function modal(html) {
    closeModal();
    const m = document.createElement('div');
    m.id = 'git-modal';
    m.innerHTML = `<div class="g-dlg">${html}</div>`;
    m.addEventListener('mousedown', e => { if (e.target === m) closeModal(); });
    m.addEventListener('keydown', e => { if (e.key === 'Escape') { e.stopPropagation(); closeModal(); } });
    document.body.appendChild(m);
    const first = m.querySelector('input, select, textarea');
    if (first) first.focus();
  }
  function closeModal() { const m = $('git-modal'); if (m) m.remove(); }

  let qTimer = null;
  window.gitWorkspace = {
    open, close, tab: setTab, check, refresh, details, diff, forget, addInstall, editApp, delApp, suggest, importCatalog, exportCatalog, exportBundle,
    editSource, testSource, delSource, importBundle, saveSettings, saveCatalogSettings, setRef, filterConns, selConns, closeModal,
    filter(k, v, debounce) { G.f[k] = v; clearTimeout(qTimer); if (debounce) qTimer = setTimeout(() => { render(); const i = document.querySelector('#git-body input[type=search]'); if (i) { i.focus(); i.setSelectionRange(i.value.length, i.value.length); } }, 250); else render(); },
    toggleChanges(on) { G.onlyChanges = on; if (G.detail) details(G.detail.id); },
  };
})();
