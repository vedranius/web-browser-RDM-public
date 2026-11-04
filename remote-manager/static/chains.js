/* Copyright 2026 vedranius and WRM contributors
   SPDX-License-Identifier: Apache-2.0 */

// Chains: live monitoring of a service chain (chains.go). Loaded on demand by openChains() in
// index.html; uses its helpers (t, esc, showToast). A live SVG diagram (automatic layout, dragged
// positions saved, animated flows, node/link details, a scenario band over time) and a list
// view, generic starter templates, and a JSON editor with validation, import and export.
'use strict';

Object.assign(LANGS.en, {
  chn_title: 'Chains', chn_new: 'New chain', chn_none: 'No chains yet. A chain links the servers of one system and checks them together.',
  chn_pick: 'Choose a chain on the left, or create one.', chn_check_now: 'Check now', chn_checking: 'Checking…', chn_edit: 'Edit JSON',
  chn_export: 'Export', chn_import: 'Import', chn_delete: 'Delete', chn_delete_q: 'Delete the chain {name}? Its timeline is deleted too; the servers are not touched.',
  chn_scenario: 'Scenario', chn_no_scenario: 'no scenario matches', chn_not_checked: 'not checked yet', chn_last: 'Last check',
  chn_nodes: 'Nodes', chn_checks: 'Checks', chn_links: 'Links', chn_scenarios: 'Scenarios', chn_timeline: 'Changes',
  chn_timeline_none: 'No changes recorded yet.', chn_up: 'up', chn_down: 'down', chn_disabled: 'paused',
  chn_kind_server: 'server', chn_kind_storage: 'storage', chn_kind_external: 'outside endpoint',
  chn_editor_t: 'Chain definition (JSON)', chn_validate: 'Validate', chn_valid: 'The chain is valid.', chn_save: 'Save', chn_cancel: 'Cancel',
  chn_saved: 'Saved. The first check runs within a minute.', chn_enabled: 'Checked in the background', chn_help: 'Reference',
  chn_help_conns: 'Your SSH connections (conn_id)', chn_help_storages: 'Your storages (storage_id)', chn_help_types: 'Check types',
  chn_help_rules: 'Rules: {"all":[…]}, {"any":[…]}, {"check":"id","field":"value|ok","op":"==|!=|in|contains|exists","value":"…"} or {"node":"id","op":"==","value":"up"}. The first scenario that matches is active; an empty rule {} matches always.',
  chn_import_bad: 'The file is not a chain export.', chn_failed: 'The request failed.', chn_ok_count: '{ok}/{n} checks OK',
  chn_nodes_count: '{up}/{n} servers up',
  chn_type_systemd: "systemd unit is active (target: unit name)",
  chn_type_timer: "systemd timer is active (target: timer name)",
  chn_type_docker: "Docker container is running (target: container name)",
  chn_type_process: "number of processes whose command line matches (target: pattern)",
  chn_type_port: "a TCP port is listening on the server (port)",
  chn_type_file: "file exists, optionally not older than max_age_s (target: absolute path)",
  chn_type_dir_count: "number of files in a directory (target: absolute path)",
  chn_type_config: "one value from a configuration file (format: ini, json, yaml, env, kv or regex; key: section.key, a.b.c or a regular expression)",
  chn_type_mount: "a mount point is mounted, read-write or read-only (target: mount point)",
  chn_type_cron: "a crontab line with the text is enabled, disabled (commented out) or missing (target: text)",
  chn_type_monitor: "the state of WRM's own monitoring for the server",
  chn_type_dns: "DNS answer for a name (target: name; record: A, AAAA, CNAME, TXT, MX, NS; resolver: optional ip[:port])",
  chn_type_tcp: "a TCP port answers, checked from WRM (target: host:port)",
  chn_type_http: "HTTP(S) status, or a JSON value with key (target: URL)",
  chn_type_tls: "days until the TLS certificate expires (target: host[:port]; warn_days, default 14)",
  chn_type_storage: "the storage's bucket answers; with prefix, the newest object under it (and max_age_s)",
});
Object.assign(LANGS.hr, {
  chn_title: 'Lanci', chn_new: 'Novi lanac', chn_none: 'Nema lanaca. Lanac povezuje servere jednog sustava i provjerava ih zajedno.',
  chn_pick: 'Odaberite lanac lijevo ili napravite novi.', chn_check_now: 'Provjeri sada', chn_checking: 'Provjeravam…', chn_edit: 'Uredi JSON',
  chn_export: 'Izvezi', chn_import: 'Uvezi', chn_delete: 'Obriši', chn_delete_q: 'Obrisati lanac {name}? Briše se i njegova povijest promjena; serveri se ne diraju.',
  chn_scenario: 'Scenarij', chn_no_scenario: 'nijedan scenarij ne odgovara', chn_not_checked: 'još nije provjereno', chn_last: 'Zadnja provjera',
  chn_nodes: 'Čvorovi', chn_checks: 'Provjere', chn_links: 'Veze', chn_scenarios: 'Scenariji', chn_timeline: 'Promjene',
  chn_timeline_none: 'Još nema zabilježenih promjena.', chn_up: 'radi', chn_down: 'ne radi', chn_disabled: 'pauzirano',
  chn_kind_server: 'server', chn_kind_storage: 'storage', chn_kind_external: 'vanjski endpoint',
  chn_editor_t: 'Definicija lanca (JSON)', chn_validate: 'Provjeri', chn_valid: 'Lanac je ispravan.', chn_save: 'Spremi', chn_cancel: 'Odustani',
  chn_saved: 'Spremljeno. Prva provjera kreće unutar minute.', chn_enabled: 'Provjerava se u pozadini', chn_help: 'Pomoć',
  chn_help_conns: 'Vaše SSH veze (conn_id)', chn_help_storages: 'Vaši storagei (storage_id)', chn_help_types: 'Vrste provjera',
  chn_help_rules: 'Pravila: {"all":[…]}, {"any":[…]}, {"check":"id","field":"value|ok","op":"==|!=|in|contains|exists","value":"…"} ili {"node":"id","op":"==","value":"up"}. Aktivan je prvi scenarij koji odgovara; prazno pravilo {} uvijek odgovara.',
  chn_import_bad: 'Datoteka nije izvoz lanca.', chn_failed: 'Zahtjev nije uspio.', chn_ok_count: '{ok}/{n} provjera u redu',
  chn_nodes_count: '{up}/{n} servera radi',
  chn_type_systemd: "systemd unit je aktivan (target: naziv unita)",
  chn_type_timer: "systemd timer je aktivan (target: naziv timera)",
  chn_type_docker: "Docker kontejner radi (target: naziv kontejnera)",
  chn_type_process: "broj procesa čija naredba odgovara uzorku (target: uzorak)",
  chn_type_port: "TCP port sluša na serveru (port)",
  chn_type_file: "datoteka postoji, po želji ne starija od max_age_s (target: apsolutna putanja)",
  chn_type_dir_count: "broj datoteka u direktoriju (target: apsolutna putanja)",
  chn_type_config: "jedna vrijednost iz konfiguracijske datoteke (format: ini, json, yaml, env, kv ili regex; key: section.key, a.b.c ili regularni izraz)",
  chn_type_mount: "mount point je montiran, za čitanje i pisanje ili samo za čitanje (target: mount point)",
  chn_type_cron: "crontab linija s tekstom je uključena, isključena (zakomentirana) ili ne postoji (target: tekst)",
  chn_type_monitor: "stanje WRM-ovog nadzora za server",
  chn_type_dns: "DNS odgovor za ime (target: ime; record: A, AAAA, CNAME, TXT, MX, NS; resolver: po želji ip[:port])",
  chn_type_tcp: "TCP port odgovara, provjereno iz WRM-a (target: host:port)",
  chn_type_http: "HTTP(S) status ili JSON vrijednost s key (target: URL)",
  chn_type_tls: "dana do isteka TLS certifikata (target: host[:port]; warn_days, zadano 14)",
  chn_type_storage: "bucket storagea odgovara; s prefix, najnoviji objekt ispod njega (i max_age_s)",
});

// Step 2: the live diagram, templates and the scenario timeline.
Object.assign(LANGS.en, {
  chn_view_diagram: 'Diagram', chn_view_list: 'List', chn_zoom_in: 'Zoom in', chn_zoom_out: 'Zoom out', chn_fit: 'Fit to window',
  chn_auto_layout: 'Arrange automatically', chn_auto_layout_q: 'Forget the saved positions and arrange the diagram automatically?',
  chn_drag_hint: 'Drag nodes to arrange them; the positions are saved. Click a node or a link for details.',
  chn_pos_failed: 'The positions were not saved.', chn_summary: 'Overview', chn_link: 'Link',
  chn_link_idle: 'idle in this scenario', chn_link_active_in: 'Carries traffic in', chn_link_check: 'Coloured by check',
  chn_history: 'History', chn_no_history: 'No changes recorded for this yet.', chn_definition: 'Definition',
  chn_conn: 'SSH connection', chn_storage: 'Storage', chn_band_t: 'Scenario over time', chn_band_none: 'No scenario recorded yet.',
  chn_now: 'now', chn_since: 'since {at}', chn_flow: 'Traffic flows', chn_legend_idle: 'idle', chn_legend_bad: 'failing',
  chn_lk_data: 'data', chn_lk_read: 'reads', chn_lk_write: 'writes', chn_lk_failover: 'failover', chn_lk_dependency: 'depends on',
  chn_tpl_pick: 'Start from a template', chn_tpl_empty: 'Empty chain', chn_tpl_empty_d: 'Write the definition in JSON yourself.',
  chn_tpl_single: 'One service', chn_tpl_single_d: 'One server with a service and its port, and a public HTTPS check.',
  chn_tpl_failover: 'Primary and backup path', chn_tpl_failover_d: 'Ingest on a primary and a backup server, a content server writing to a bucket, and a CDN in front.',
  chn_tpl_name: 'Name', chn_tpl_servers: 'Choose the server for each node', chn_tpl_storages: 'Choose the storage for each storage node',
  chn_tpl_choose: '— choose —', chn_tpl_no_ssh: 'You have no SSH connections yet. Add one first.',
  chn_tpl_no_storage: 'You have no storages yet. Add one under Storages, or remove the node in the JSON editor.',
  chn_tpl_missing: 'Choose a server or storage for every node.', chn_tpl_after: 'Service names, paths and addresses are examples: adjust them in the JSON editor.',
  chn_create: 'Create', chn_back: 'Back', chn_nodes_n: '{n} nodes',
  chn_help_links: 'Links: {"from":"a","to":"b","kind":"data|read|write|failover|dependency","check":"id","scenarios":["primary"]}. With scenarios, the link carries traffic only in those scenarios; otherwise the diagram shows it idle.',
});
Object.assign(LANGS.hr, {
  chn_view_diagram: 'Dijagram', chn_view_list: 'Popis', chn_zoom_in: 'Povećaj', chn_zoom_out: 'Smanji', chn_fit: 'Prilagodi prozoru',
  chn_auto_layout: 'Rasporedi automatski', chn_auto_layout_q: 'Zaboraviti spremljene položaje i automatski rasporediti dijagram?',
  chn_drag_hint: 'Povucite čvorove da ih rasporedite; položaji se spremaju. Kliknite čvor ili vezu za detalje.',
  chn_pos_failed: 'Položaji nisu spremljeni.', chn_summary: 'Pregled', chn_link: 'Veza',
  chn_link_idle: 'neaktivna u ovom scenariju', chn_link_active_in: 'Promet nosi u', chn_link_check: 'Boju određuje provjera',
  chn_history: 'Povijest', chn_no_history: 'Za ovo još nema zabilježenih promjena.', chn_definition: 'Definicija',
  chn_conn: 'SSH veza', chn_storage: 'Storage', chn_band_t: 'Scenarij kroz vrijeme', chn_band_none: 'Još nema zabilježenog scenarija.',
  chn_now: 'sada', chn_since: 'od {at}', chn_flow: 'Promet teče', chn_legend_idle: 'neaktivno', chn_legend_bad: 'ne radi',
  chn_lk_data: 'podaci', chn_lk_read: 'čita', chn_lk_write: 'piše', chn_lk_failover: 'rezervni put', chn_lk_dependency: 'ovisi o',
  chn_tpl_pick: 'Počnite od predloška', chn_tpl_empty: 'Prazan lanac', chn_tpl_empty_d: 'Definiciju napišite sami u JSON-u.',
  chn_tpl_single: 'Jedan servis', chn_tpl_single_d: 'Jedan server sa servisom i njegovim portom te javna HTTPS provjera.',
  chn_tpl_failover: 'Primarni i rezervni put', chn_tpl_failover_d: 'Ingest na primarnom i rezervnom serveru, content server koji piše u bucket i CDN ispred.',
  chn_tpl_name: 'Naziv', chn_tpl_servers: 'Odaberite server za svaki čvor', chn_tpl_storages: 'Odaberite storage za svaki storage čvor',
  chn_tpl_choose: '— odaberite —', chn_tpl_no_ssh: 'Još nemate SSH veza. Najprije dodajte jednu.',
  chn_tpl_no_storage: 'Još nemate storagea. Dodajte ga pod Storages ili uklonite čvor u JSON editoru.',
  chn_tpl_missing: 'Odaberite server ili storage za svaki čvor.', chn_tpl_after: 'Nazivi servisa, putanje i adrese su primjeri: prilagodite ih u JSON editoru.',
  chn_create: 'Napravi', chn_back: 'Natrag', chn_nodes_n: '{n} čvorova',
  chn_help_links: 'Veze: {"from":"a","to":"b","kind":"data|read|write|failover|dependency","check":"id","scenarios":["primary"]}. Sa scenarios veza nosi promet samo u tim scenarijima; inače je dijagram prikazuje neaktivnom.',
});

(function () {
  const state = { list: [], current: null, timeline: [], styleAdded: false, mode: 'diagram', sel: null, views: {}, pos: {}, openCheck: null };
  const $ = (sel) => document.querySelector(sel);
  const el = (tag, cls, text) => { const e = document.createElement(tag); if (cls) e.className = cls; if (text != null) e.textContent = text; return e; };
  const fmt = (s, v) => s.replace(/\{(\w+)\}/g, (_, k) => (v[k] != null ? String(v[k]) : ''));
  const when = (at) => {
    if (!at) return t('chn_not_checked');
    const d = new Date(at);
    return isNaN(d) ? at.slice(0, 19).replace('T', ' ') : d.toLocaleString();
  };
  const trunc = (s, n) => { s = String(s || ''); return s.length > n ? s.slice(0, n - 1) + '…' : s; };
  const NW = 184, NH = 78, GX = 250, GY = 112;
  const SCN_COLORS = ['#3b82f6', '#f59e0b', '#8b5cf6', '#10b981', '#ec4899', '#14b8a6', '#eab308', '#6366f1'];
  try { if (localStorage.getItem('wrm.chains.view') === 'list') state.mode = 'list'; } catch (_) { /* storage blocked */ }

  function css() {
    if (state.styleAdded) return;
    state.styleAdded = true;
    const s = document.createElement('style');
    s.textContent = `
.chn-wrap { display:flex; width:min(1320px,97vw); height:88vh; padding:0; overflow:hidden; }
.chn-side { width:250px; min-width:200px; border-right:1px solid var(--border2); display:flex; flex-direction:column; }
.chn-side-head { padding:14px 14px 8px; display:flex; align-items:center; justify-content:space-between; gap:8px; flex-wrap:wrap; }
.chn-side-head h3 { margin:0; font-size:14px; }
.chn-list { flex:1; overflow:auto; padding:6px; }
.chn-item { padding:9px 10px; border-radius:8px; cursor:pointer; border:1px solid transparent; margin-bottom:4px; }
.chn-item:hover, .chn-item.sel { background:var(--bg3); }
.chn-item.sel { border-color:var(--accent); }
.chn-item .n { font-size:12.5px; font-weight:600; }
.chn-item .m { font-size:10.5px; color:var(--text2); margin-top:3px; }
.chn-dot { display:inline-block; width:8px; height:8px; border-radius:50%; background:var(--text3); opacity:.6; margin-right:7px; vertical-align:middle; flex:none; }
.chn-dot.ok { background:#3fb950; opacity:1; } .chn-dot.bad { background:#f85149; opacity:1; } .chn-dot.warn { background:#d29922; opacity:1; }
.chn-main { flex:1; display:flex; flex-direction:column; min-width:0; position:relative; }
.chn-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; align-items:center; gap:8px; flex-wrap:wrap; }
.chn-head h3 { margin:0; font-size:14px; flex:1; min-width:140px; }
.chn-seg { display:inline-flex; border:1px solid var(--border2); border-radius:7px; overflow:hidden; }
.chn-wrap .chn-seg button { border:0; border-radius:0; }
.chn-wrap .chn-seg button + button { border-left:1px solid var(--border2); }
.chn-body { flex:1; overflow:auto; padding:12px 14px; min-height:0; }
.chn-body.dg { padding:0; overflow:hidden; display:flex; flex-direction:column; }
.chn-empty { padding:30px 14px; text-align:center; color:var(--text2); font-size:12px; }
.chn-badge { font-size:11px; padding:2px 8px; border-radius:10px; background:var(--bg3); color:var(--text2); }
.chn-badge.active { background:var(--accent); color:#fff; }
.chn-sec { margin:4px 0 16px; }
.chn-sec h4 { margin:0 0 8px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-nodes { display:grid; grid-template-columns:repeat(auto-fill,minmax(280px,1fr)); gap:10px; }
.chn-node { border:1px solid var(--border2); border-radius:10px; padding:10px 12px; background:var(--bg2); }
.chn-node .h { display:flex; align-items:center; gap:6px; font-weight:600; font-size:12.5px; }
.chn-node .r { font-size:10.5px; color:var(--text2); margin:2px 0 6px; }
.chn-node .e, .chn-det .e { font-size:10.5px; color:#f85149; margin-bottom:6px; overflow-wrap:anywhere; }
.chn-chk { display:flex; flex-wrap:wrap; align-items:baseline; gap:2px 6px; font-size:11.5px; padding:3px 0; border-top:1px dashed var(--border2); }
.chn-chk .l { flex:1; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .v { font-family:var(--font-mono); font-size:11px; max-width:45%; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .d { flex-basis:100%; font-size:10px; color:var(--text2); padding-left:14px; overflow-wrap:anywhere; }
.chn-chk.click { cursor:pointer; } .chn-chk.click:hover { background:var(--bg3); }
.chn-chk-more { font-size:10.5px; padding:4px 0 6px 14px; color:var(--text2); }
.chn-chk-more pre { font-family:var(--font-mono); font-size:10.5px; white-space:pre-wrap; overflow-wrap:anywhere; background:var(--bg); border-radius:6px; padding:6px 8px; margin:4px 0; color:var(--text); }
.chn-row { display:flex; gap:8px; align-items:center; font-size:11.5px; padding:4px 0; flex-wrap:wrap; }
.chn-row .k { color:var(--text2); font-size:10.5px; }
.chn-tl { font-size:11px; padding:3px 0; display:flex; gap:8px; align-items:baseline; }
.chn-tl .a { color:var(--text2); width:150px; flex:none; font-family:var(--font-mono); font-size:10.5px; }
.chn-dg-top { flex:1; display:flex; min-height:0; }
.chn-canvas { flex:1; position:relative; min-width:0; background:var(--bg); background-image:radial-gradient(var(--border2) 1px, transparent 1px); background-size:22px 22px; }
.chn-canvas svg { width:100%; height:100%; display:block; touch-action:none; cursor:grab; user-select:none; }
.chn-canvas svg.panning { cursor:grabbing; }
.chn-tools { position:absolute; top:8px; left:8px; display:flex; gap:4px; flex-wrap:wrap; }
.chn-tools button { min-width:30px; }
.chn-hint { position:absolute; bottom:6px; left:10px; right:10px; font-size:10.5px; color:var(--text2); pointer-events:none; }
.chn-legend { position:absolute; top:8px; right:10px; font-size:10.5px; color:var(--text2); display:flex; gap:10px; background:var(--bg2); border:1px solid var(--border2); border-radius:7px; padding:3px 8px; }
.chn-legend i { display:inline-block; width:18px; height:0; border-top:2px solid; vertical-align:middle; margin-right:4px; }
.chn-det { width:310px; flex:none; border-left:1px solid var(--border2); overflow:auto; padding:12px; font-size:11.5px; background:var(--bg2); }
.chn-det h4 { margin:0 0 4px; font-size:13px; display:flex; align-items:center; gap:6px; }
.chn-det h5 { margin:14px 0 6px; font-size:10.5px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-det .r { font-size:10.5px; color:var(--text2); margin-bottom:8px; }
.chn-det .act { display:flex; gap:6px; flex-wrap:wrap; margin:8px 0; }
.chn-band { flex:none; max-height:30%; overflow:auto; border-top:1px solid var(--border2); padding:8px 14px 10px; }
.chn-band h4 { margin:0 0 6px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-bar { display:flex; height:16px; border-radius:5px; overflow:hidden; background:var(--bg3); }
.chn-bar span { min-width:3px; height:100%; }
.chn-marks { position:relative; height:10px; margin-top:2px; }
.chn-marks span { position:absolute; top:1px; width:3px; height:8px; border-radius:1px; transform:translateX(-1px); }
.chn-axis { display:flex; justify-content:space-between; font-size:10px; color:var(--text2); margin:2px 0 4px; }
.chn-keys { display:flex; flex-wrap:wrap; gap:4px 12px; font-size:10.5px; color:var(--text2); margin-bottom:6px; }
.chn-keys b { display:inline-block; width:10px; height:10px; border-radius:2px; vertical-align:middle; margin-right:4px; }
.chn-node-g { cursor:pointer; }
.chn-node-g text { font-family:var(--font-ui); pointer-events:none; }
.chn-link-hit { stroke:transparent; stroke-width:14; fill:none; cursor:pointer; }
.chn-link { fill:none; stroke-width:2.2; }
.chn-link.ok { stroke:#3fb950; } .chn-link.bad { stroke:#f85149; } .chn-link.unknown { stroke:var(--text2); opacity:.7; } .chn-link.idle { stroke:var(--text2); opacity:.45; stroke-dasharray:4 5; }
.chn-link.k-failover:not(.idle), .chn-link.k-dependency:not(.idle) { stroke-dasharray:8 5; }
.chn-link.sel { stroke-width:3.4; }
.chn-flow { fill:none; stroke-width:2.4; stroke:#7ee787; stroke-dasharray:3 15; animation:chn-flow 1.1s linear infinite; pointer-events:none; }
.chn-flow.rev { animation-direction:reverse; }
@keyframes chn-flow { to { stroke-dashoffset:-18; } }
.chn-ll { font-size:10px; fill:var(--text2); paint-order:stroke; stroke:var(--bg); stroke-width:3px; pointer-events:none; }
.chn-editor { position:absolute; inset:0; background:var(--bg2); display:flex; flex-direction:column; z-index:5; }
.chn-editor-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; gap:8px; align-items:center; flex-wrap:wrap; }
.chn-editor-head h3 { margin:0; font-size:13px; flex:1; min-width:160px; }
.chn-editor-body { flex:1; display:flex; min-height:0; }
.chn-editor textarea { flex:1; border:0; border-radius:0; font-family:var(--font-mono); font-size:12px; resize:none; padding:12px 14px; background:var(--bg); color:var(--text); min-width:0; }
.chn-help { width:320px; overflow:auto; border-left:1px solid var(--border2); padding:10px 12px; font-size:11px; color:var(--text2); }
.chn-help h5 { margin:10px 0 4px; font-size:10.5px; text-transform:uppercase; color:var(--text2); }
.chn-help code { font-family:var(--font-mono); font-size:10.5px; }
.chn-errs { max-height:120px; overflow:auto; padding:6px 14px; font-size:11px; color:#f85149; border-top:1px solid var(--border2); }
.chn-errs:empty { display:none; }
.chn-tpls { flex:1; overflow:auto; padding:14px; }
.chn-tpls > h4 { margin:0 0 10px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); }
.chn-tpl-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(240px,1fr)); gap:10px; }
.chn-tpl { border:1px solid var(--border2); border-radius:10px; padding:12px; cursor:pointer; background:var(--bg); text-align:left; color:var(--text); font:inherit; }
.chn-tpl:hover, .chn-tpl:focus-visible { border-color:var(--accent); }
.chn-tpl b { display:block; font-size:12.5px; margin-bottom:4px; } .chn-tpl span { font-size:11px; color:var(--text2); }
.chn-form { max-width:560px; display:flex; flex-direction:column; gap:10px; font-size:12px; }
.chn-form label { display:flex; flex-direction:column; gap:4px; font-size:11.5px; color:var(--text2); }
.chn-form input, .chn-form select { font-size:12.5px; background:var(--bg); border:1px solid var(--border2); color:var(--text); border-radius:7px; padding:6px 8px; outline:none; font-family:inherit; }
.chn-form input:focus, .chn-form select:focus { border-color:var(--accent); }
.chn-form h4 { font-size:13px; margin:0; } .chn-form h5 { font-size:10.5px; text-transform:uppercase; letter-spacing:.05em; color:var(--text2); margin:6px 0 0; }
.chn-form .note { font-size:11px; color:var(--text2); }
@media (max-width: 900px) { .chn-det { width:250px; } }
@media (max-width: 760px) {
  .chn-wrap { flex-direction:column; height:92vh; }
  .chn-side { width:auto; max-height:24vh; border-right:0; border-bottom:1px solid var(--border2); }
  .chn-help { display:none; }
  .chn-body.dg { overflow:auto; }
  .chn-dg-top { flex-direction:column; flex:none; }
  .chn-canvas { height:46vh; flex:none; }
  .chn-det { width:auto; border-left:0; border-top:1px solid var(--border2); }
  .chn-band { max-height:none; }
  .chn-legend { display:none; }
}
`;
    document.head.appendChild(s);
  }

  function api(method, path, body) {
    const opts = { method, headers: { 'X-WRM-Request': '1' }, credentials: 'same-origin' };
    if (body !== undefined) { opts.body = JSON.stringify(body); opts.headers['Content-Type'] = 'application/json'; }
    return fetch(path, opts).then(async (res) => {
      const text = await res.text();
      let data = null;
      try { data = text ? JSON.parse(text) : null; } catch (_) { data = null; }
      if (!res.ok) {
        const err = new Error((data && data.error) || t('chn_failed'));
        err.status = res.status; err.errors = (data && data.errors) || [];
        throw err;
      }
      return data;
    });
  }

  function overlay() { return document.getElementById('chn-overlay'); }

  function shell() {
    let ov = overlay();
    if (ov) return ov;
    ov = el('div', 'modal-overlay');
    ov.id = 'chn-overlay';
    ov.innerHTML = `<div class="modal chn-wrap" role="dialog" aria-label="${esc(t('chn_title'))}">
      <div class="chn-side"><div class="chn-side-head"><h3>⛓ ${esc(t('chn_title'))}</h3>
        <span><button class="btn-sec btn-sm" id="chn-import">${esc(t('chn_import'))}</button>
        <button class="btn-sec btn-sm" id="chn-new">＋ ${esc(t('chn_new'))}</button></span></div>
        <div class="chn-list" id="chn-list"></div></div>
      <div class="chn-main"><div class="chn-head"><h3 id="chn-name"></h3>
          <span class="chn-badge" id="chn-scn" style="display:none;"></span>
          <span class="chn-seg" id="chn-mode" role="group" style="display:none;"><button class="btn-sec btn-sm" data-m="diagram">${esc(t('chn_view_diagram'))}</button><button class="btn-sec btn-sm" data-m="list">${esc(t('chn_view_list'))}</button></span>
          <button class="btn-sec btn-sm" id="chn-check" style="display:none;">↻ ${esc(t('chn_check_now'))}</button>
          <button class="btn-sec btn-sm" id="chn-edit" style="display:none;">{ } ${esc(t('chn_edit'))}</button>
          <button class="btn-sec btn-sm" id="chn-export" style="display:none;">⬇ ${esc(t('chn_export'))}</button>
          <button class="btn-sec btn-sm" id="chn-del" style="display:none;">${esc(t('chn_delete'))}</button>
          <button class="btn-sec btn-sm" id="chn-close" aria-label="${esc(t('chn_cancel'))}">✕</button>
        </div>
        <div class="chn-body" id="chn-body"></div></div></div>`;
    document.body.appendChild(ov);
    ov.addEventListener('mousedown', (e) => { if (e.target === ov) close(); });
    ov.querySelector('#chn-close').onclick = close;
    ov.querySelector('#chn-new').onclick = openNew;
    ov.querySelector('#chn-import').onclick = importFile;
    ov.querySelector('#chn-check').onclick = (e) => checkNow(e.currentTarget);
    ov.querySelector('#chn-edit').onclick = () => state.current && openEditor(state.current.id, { definition: state.current.definition, enabled: state.current.enabled });
    ov.querySelector('#chn-export').onclick = () => { if (state.current) location.href = `/api/chains/${state.current.id}/export`; };
    ov.querySelector('#chn-del').onclick = deleteChain;
    ov.querySelectorAll('#chn-mode button').forEach((b) => { b.onclick = () => setMode(b.dataset.m); });
    return ov;
  }

  function setMode(m) {
    state.mode = m === 'list' ? 'list' : 'diagram';
    try { localStorage.setItem('wrm.chains.view', state.mode); } catch (_) { /* storage blocked */ }
    renderDetail();
  }

  function open() {
    css();
    shell().classList.add('open');
    api('GET', '/api/connections').then((l) => { state.conns = l || []; renderSide(); }).catch(() => {});
    loadList();
  }

  function close() {
    const ov = overlay();
    if (ov) ov.classList.remove('open');
  }

  function isOpen() { const ov = overlay(); return !!(ov && ov.classList.contains('open')); }

  function loadList() {
    return api('GET', '/api/chains').then((list) => {
      state.list = list || [];
      renderList();
      if (state.current && state.list.some((c) => c.id === state.current.id)) return select(state.current.id, true);
      if (state.list.length) return select(state.list[0].id);
      state.current = null;
      renderDetail();
      return null;
    }).catch((e) => showToast(e.message, 'error'));
  }

  function statusClass(c) {
    if (!c.enabled || !c.last_at) return '';
    if (c.nodes_up < c.nodes_total || c.checks_ok < c.checks_total) return c.nodes_up < c.nodes_total ? 'bad' : 'warn';
    return 'ok';
  }

  function renderList() {
    const box = $('#chn-list');
    if (!box) return;
    box.textContent = '';
    if (!state.list.length) { box.appendChild(el('div', 'chn-empty', t('chn_none'))); return; }
    state.list.forEach((c) => {
      const item = el('div', 'chn-item' + (state.current && state.current.id === c.id ? ' sel' : ''));
      item.tabIndex = 0;
      const n = el('div', 'n');
      n.appendChild(el('span', 'chn-dot ' + statusClass(c)));
      n.appendChild(document.createTextNode(c.name));
      const m = el('div', 'm');
      m.textContent = !c.enabled ? t('chn_disabled')
        : [c.scenario_label || t('chn_no_scenario'), c.nodes_total ? fmt(t('chn_nodes_count'), { up: c.nodes_up, n: c.nodes_total }) : '',
          fmt(t('chn_ok_count'), { ok: c.checks_ok, n: c.checks_total })].filter(Boolean).join(' · ');
      item.appendChild(n); item.appendChild(m);
      item.onclick = () => select(c.id);
      item.onkeydown = (e) => { if (e.key === 'Enter') select(c.id); };
      box.appendChild(item);
    });
  }

  // keep: a live refresh of the same chain keeps the selection and the zoom.
  function select(id, keep) {
    return api('GET', `/api/chains/${id}`).then((c) => {
      if (!keep || !state.current || state.current.id !== c.id) { state.sel = null; state.openCheck = null; state.timeline = []; }
      state.current = c;
      state.pos = {};
      renderList();
      renderDetail();
      return api('GET', `/api/chains/${id}/timeline`).then((tl) => { state.timeline = tl || []; renderTimeline(state.timeline); renderSide(); }).catch(() => {});
    }).catch((e) => showToast(e.message, 'error'));
  }

  // ─── status of nodes and links ──────────────────────────────────────────────────────────

  function checkRes(id) { const st = (state.current && state.current.state) || {}; return (st.checks || {})[id]; }

  function nodeStatus(d, n) {
    const ns = ((state.current.state || {}).nodes || {})[n.id];
    if (ns && ns.state === 'down') return 'bad';
    const res = (d.checks || []).filter((x) => x.node === n.id).map((x) => checkRes(x.id)).filter(Boolean);
    if (res.length) {
      const ok = res.filter((r) => r.ok).length;
      if (ok === res.length) return 'ok';
      return n.kind !== 'server' && ok === 0 ? 'bad' : 'warn';
    }
    return ns ? 'ok' : '';
  }

  function linkStatus(d, l) {
    const c = state.current;
    if (!c.enabled || !c.last_at) return 'unknown';
    if ((l.scenarios || []).length && l.scenarios.indexOf(c.scenario) < 0) return 'idle';
    if (l.check) { const r = checkRes(l.check); return r ? (r.ok ? 'ok' : 'bad') : 'unknown'; }
    const a = nodeStatus(d, node(d, l.from)), b = nodeStatus(d, node(d, l.to));
    if (a === 'bad' || b === 'bad') return 'bad';
    return a && b ? 'ok' : 'unknown';
  }

  function node(d, id) { return (d.nodes || []).find((x) => x.id === id) || { id, kind: 'external' }; }

  // ─── layout ─────────────────────────────────────────────────────────────────────────────

  // autoLayout places nodes in columns by their longest path along the links (sources left),
  // and orders each column by the average row of its predecessors to keep crossings low.
  function autoLayout(d) {
    const ids = (d.nodes || []).map((n) => n.id);
    const layer = {};
    ids.forEach((id) => { layer[id] = 0; });
    const links = (d.links || []).filter((l) => l.from in layer && l.to in layer && l.from !== l.to);
    for (let i = 0; i < ids.length; i++) {
      let moved = false;
      links.forEach((l) => {
        if (layer[l.to] < layer[l.from] + 1 && layer[l.from] + 1 < ids.length) { layer[l.to] = layer[l.from] + 1; moved = true; }
      });
      if (!moved) break;
    }
    const cols = [];
    ids.forEach((id) => { (cols[layer[id]] = cols[layer[id]] || []).push(id); });
    const row = {};
    cols.forEach((col, li) => {
      if (li > 0) {
        const bary = (id) => {
          const ps = links.filter((l) => l.to === id && layer[l.from] < li).map((l) => row[l.from]);
          return ps.length ? ps.reduce((a, b) => a + b, 0) / ps.length : 1e6;
        };
        const order = col.map((id, i) => ({ id, i, b: bary(id) }));
        order.sort((a, b) => (a.b - b.b) || (a.i - b.i));
        cols[li] = order.map((o) => o.id);
      }
      cols[li].forEach((id, i) => { row[id] = i; });
    });
    const tall = Math.max(1, ...cols.map((c) => (c ? c.length : 0)));
    const out = {};
    cols.forEach((col, li) => {
      if (!col) return;
      const off = (tall - col.length) * GY / 2;
      col.forEach((id, i) => { out[id] = { x: 40 + li * GX, y: 40 + off + i * GY }; });
    });
    return out;
  }

  function positions(d) {
    const auto = autoLayout(d);
    const out = {};
    (d.nodes || []).forEach((n) => {
      out[n.id] = state.pos[n.id] || ((n.x || n.y) ? { x: n.x, y: n.y } : auto[n.id]);
    });
    return out;
  }

  // ─── the diagram ────────────────────────────────────────────────────────────────────────

  function anchor(p, q) {
    const cx = p.x + NW / 2, cy = p.y + NH / 2, dx = q.x - p.x, dy = q.y - p.y;
    if (Math.abs(dx) * NH >= Math.abs(dy) * NW) return { x: dx >= 0 ? p.x + NW : p.x, y: cy, nx: dx >= 0 ? 1 : -1, ny: 0 };
    return { x: cx, y: dy >= 0 ? p.y + NH : p.y, nx: 0, ny: dy >= 0 ? 1 : -1 };
  }

  function linkPath(p, q) {
    const a = anchor(p, q), b = anchor(q, p);
    const k = Math.max(40, Math.hypot(b.x - a.x, b.y - a.y) / 2.5);
    return {
      d: `M${a.x},${a.y} C${a.x + a.nx * k},${a.y + a.ny * k} ${b.x + b.nx * k},${b.y + b.ny * k} ${b.x},${b.y}`,
      mx: (a.x + b.x) / 2 + (a.nx + b.nx) * k * 0.375, my: (a.y + b.y) / 2 + (a.ny + b.ny) * k * 0.375,
    };
  }

  const ARROW = { ok: '#3fb950', bad: '#f85149', unknown: '#8892a0', idle: '#8892a0' };

  function nodeShape(n, cls, selected) {
    const stroke = { ok: '#3fb950', warn: '#d29922', bad: '#f85149' }[cls] || 'var(--border2)';
    const sw = selected ? 3 : 1.6;
    const sel = selected ? ' stroke-dasharray="none"' : '';
    if (n.kind === 'storage') {
      return `<path d="M0,10 a${NW / 2},10 0 0,0 ${NW},0 v${NH - 20} a${NW / 2},10 0 0,1 ${-NW},0 z" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"${sel}/>
        <path d="M0,10 a${NW / 2},10 0 0,0 ${NW},0 a${NW / 2},10 0 0,0 ${-NW},0 z" fill="var(--bg3)" stroke="${stroke}" stroke-width="${sw}"/>`;
    }
    if (n.kind === 'external') {
      return `<rect width="${NW}" height="${NH}" rx="${NH / 2}" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"${selected ? '' : ' stroke-dasharray="6 4"'}/>`;
    }
    return `<rect width="${NW}" height="${NH}" rx="10" fill="var(--bg2)" stroke="${stroke}" stroke-width="${sw}"/>`;
  }

  function renderDiagram(wrap) {
    const c = state.current, d = c.definition;
    const pos = positions(d);
    const sel = state.sel || {};
    let defs = '<defs>';
    Object.keys(ARROW).forEach((k) => {
      defs += `<marker id="chn-a-${k}" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0,0 L10,5 L0,10 z" fill="${ARROW[k]}"/></marker>`;
    });
    defs += '</defs>';
    let links = '', flows = '', labels = '';
    (d.links || []).forEach((l, i) => {
      const p = pos[l.from], q = pos[l.to];
      if (!p || !q || l.from === l.to) return;
      const st = linkStatus(d, l);
      const g = linkPath(p, q);
      const kind = l.kind || 'data';
      const isSel = sel.link === i;
      const rev = kind === 'read';
      const mark = kind === 'dependency' ? '' : (rev ? ` marker-start="url(#chn-a-${st})"` : ` marker-end="url(#chn-a-${st})"`);
      links += `<path class="chn-link ${st} k-${kind}${isSel ? ' sel' : ''}" d="${g.d}"${mark}/>`;
      links += `<path class="chn-link-hit" data-link="${i}" d="${g.d}"><title>${esc(label(d, l.from) + ' → ' + label(d, l.to))}</title></path>`;
      if (st === 'ok' && kind !== 'dependency') flows += `<path class="chn-flow${rev ? ' rev' : ''}" d="${g.d}"/>`;
      const txt = l.label || (kind !== 'data' ? t('chn_lk_' + kind) : '');
      if (txt) labels += `<text class="chn-ll" x="${g.mx}" y="${g.my - 4}" text-anchor="middle">${esc(trunc(txt, 28))}</text>`;
    });
    let nodes = '';
    (d.nodes || []).forEach((n) => {
      const p = pos[n.id];
      const cls = nodeStatus(d, n);
      const checks = (d.checks || []).filter((x) => x.node === n.id);
      const res = checks.map((x) => checkRes(x.id));
      const okN = res.filter((r) => r && r.ok).length;
      const icon = { server: '🖥', storage: '🪣', external: '🌐' }[n.kind] || '•';
      const ns = ((c.state || {}).nodes || {})[n.id];
      const sub = [n.role || t('chn_kind_' + n.kind), ns && ns.state === 'down' ? t('chn_down') : ''].filter(Boolean).join(' · ');
      const top = n.kind === 'storage' ? 30 : 24;
      let dots = '';
      res.slice(0, 12).forEach((r, i) => {
        const col = r ? (r.ok ? '#3fb950' : '#f85149') : '#8892a0';
        dots += `<circle cx="${18 + i * 11}" cy="${NH - 13}" r="3.6" fill="${col}"/>`;
      });
      const count = checks.length ? `<text x="${NW - 12}" y="${NH - 9}" text-anchor="end" font-size="10.5" fill="var(--text2)">${okN}/${checks.length}</text>` : '';
      nodes += `<g class="chn-node-g" data-node="${esc(n.id)}" transform="translate(${p.x},${p.y})" tabindex="0" role="button" aria-label="${esc(n.label || n.id)}">
        ${nodeShape(n, cls, sel.node === n.id)}
        <text x="12" y="${top}" font-size="13" fill="var(--text)" font-weight="600">${icon} ${esc(trunc(n.label || n.id, 20))}</text>
        <text x="12" y="${top + 16}" font-size="10.5" fill="var(--text2)">${esc(trunc(sub, 30))}</text>
        ${dots}${count}<title>${esc([n.label || n.id, n.role, n.id].filter(Boolean).join(' · '))}</title></g>`;
    });
    let svg = wrap.querySelector('svg');
    if (!svg) {
      svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
      svg.setAttribute('preserveAspectRatio', 'xMidYMid meet');
      svg.setAttribute('role', 'img');
      wrap.insertBefore(svg, wrap.firstChild);
      wireSvg(svg);
    }
    svg.setAttribute('aria-label', d.name);
    svg.innerHTML = defs + `<g class="chn-links">${links}</g><g>${flows}</g><g>${labels}</g><g class="chn-nodes-g">${nodes}</g>`;
    if (!state.views[c.id]) state.views[c.id] = fitView(pos);
    applyView(svg);
  }

  function fitView(pos) {
    const pts = Object.values(pos).filter(Boolean);
    if (!pts.length) return { x: 0, y: 0, w: 800, h: 400 };
    const minX = Math.min(...pts.map((p) => p.x)) - 40, minY = Math.min(...pts.map((p) => p.y)) - 50;
    const maxX = Math.max(...pts.map((p) => p.x + NW)) + 40, maxY = Math.max(...pts.map((p) => p.y + NH)) + 40;
    return { x: minX, y: minY, w: Math.max(420, maxX - minX), h: Math.max(240, maxY - minY) };
  }

  function applyView(svg) {
    const v = state.views[state.current.id];
    svg.setAttribute('viewBox', `${v.x} ${v.y} ${v.w} ${v.h}`);
  }

  function zoom(factor, cx, cy) {
    const c = state.current;
    const svg = $('.chn-canvas svg');
    if (!c || !svg) return;
    const v = state.views[c.id];
    const nw = Math.min(10000, Math.max(200, v.w * factor));
    const f = nw / v.w;
    const px = cx == null ? v.x + v.w / 2 : cx, py = cy == null ? v.y + v.h / 2 : cy;
    state.views[c.id] = { x: px - (px - v.x) * f, y: py - (py - v.y) * f, w: v.w * f, h: v.h * f };
    applyView(svg);
  }

  function svgPoint(svg, e) {
    const m = svg.getScreenCTM();
    if (!m) return { x: 0, y: 0 };
    const p = new DOMPoint(e.clientX, e.clientY).matrixTransform(m.inverse());
    return { x: p.x, y: p.y };
  }

  // Dragging a node moves it (and saves every position, so the picture stays as seen); a
  // press that does not move selects. Dragging the background pans; the wheel zooms.
  function wireSvg(svg) {
    let drag = null;
    svg.addEventListener('pointerdown', (e) => {
      if (e.button !== 0 || !state.current) return;
      const g = e.target.closest('[data-node]');
      const lk = e.target.closest('[data-link]');
      const p = svgPoint(svg, e);
      const d = state.current.definition;
      if (g) {
        const id = g.getAttribute('data-node');
        const start = positions(d)[id];
        drag = { kind: 'node', id, p, start, moved: false };
      } else if (lk) {
        drag = { kind: 'link', i: Number(lk.getAttribute('data-link')), p, moved: false };
      } else {
        drag = { kind: 'pan', cx: e.clientX, cy: e.clientY, v: Object.assign({}, state.views[state.current.id]), moved: false };
        svg.classList.add('panning');
      }
      svg.setPointerCapture(e.pointerId);
    });
    svg.addEventListener('pointermove', (e) => {
      if (!drag || !state.current) return;
      if (drag.kind === 'pan') {
        const r = svg.getBoundingClientRect();
        const v = drag.v;
        const s = Math.max(v.w / r.width, v.h / r.height);
        if (Math.abs(e.clientX - drag.cx) + Math.abs(e.clientY - drag.cy) > 3) drag.moved = true;
        state.views[state.current.id] = { x: v.x - (e.clientX - drag.cx) * s, y: v.y - (e.clientY - drag.cy) * s, w: v.w, h: v.h };
        applyView(svg);
        return;
      }
      if (drag.kind !== 'node') return;
      const p = svgPoint(svg, e);
      const dx = p.x - drag.p.x, dy = p.y - drag.p.y;
      if (!drag.moved && Math.hypot(dx, dy) < 4) return;
      if (!drag.moved) {
        drag.moved = true;
        // freeze the current picture: every node keeps where it is now
        state.pos = positions(state.current.definition);
      }
      state.pos[drag.id] = { x: Math.max(1, Math.round(drag.start.x + dx)), y: Math.max(1, Math.round(drag.start.y + dy)) };
      renderDiagram(svg.parentNode);
    });
    const end = (e) => {
      if (!drag) return;
      const was = drag;
      drag = null;
      svg.classList.remove('panning');
      try { svg.releasePointerCapture(e.pointerId); } catch (_) { /* already released */ }
      if (e.type === 'pointercancel') return;
      if (was.kind === 'node' && was.moved) { saveLayout(); return; }
      if (was.kind === 'node') state.sel = { node: was.id };
      else if (was.kind === 'link') state.sel = { link: was.i };
      else if (!was.moved) state.sel = null;
      else return;
      state.openCheck = null;
      renderDiagram(svg.parentNode);
      renderSide();
    };
    svg.addEventListener('pointerup', end);
    svg.addEventListener('pointercancel', end);
    svg.addEventListener('wheel', (e) => {
      e.preventDefault();
      const p = svgPoint(svg, e);
      zoom(e.deltaY > 0 ? 1.12 : 1 / 1.12, p.x, p.y);
    }, { passive: false });
    svg.addEventListener('keydown', (e) => {
      const g = e.target.closest && e.target.closest('[data-node]');
      if (g && (e.key === 'Enter' || e.key === ' ')) {
        e.preventDefault();
        state.sel = { node: g.getAttribute('data-node') };
        state.openCheck = null;
        renderDiagram(svg.parentNode);
        renderSide();
      }
    });
  }

  function saveLayout() {
    const c = state.current;
    if (!c) return;
    const positionsNow = Object.assign({}, state.pos);
    api('PUT', `/api/chains/${c.id}/layout`, { positions: positionsNow }).then((got) => {
      if (state.current && state.current.id === got.id) {
        state.current.definition = got.definition;
      }
    }).catch((e) => showToast(t('chn_pos_failed') + ' ' + e.message, 'error'));
  }

  function resetLayout() {
    const c = state.current;
    if (!c || !confirm(t('chn_auto_layout_q'))) return;
    api('PUT', `/api/chains/${c.id}/layout`, { reset: true }).then((got) => {
      state.current.definition = got.definition;
      state.pos = {};
      delete state.views[c.id];
      renderDetail();
    }).catch((e) => showToast(e.message, 'error'));
  }

  // ─── the detail panel next to the diagram ──────────────────────────────────────────────

  function renderSide() {
    const box = $('#chn-det');
    const c = state.current;
    if (!box || !c) return;
    const d = c.definition;
    box.textContent = '';
    const sel = state.sel || {};
    if (sel.node && (d.nodes || []).some((n) => n.id === sel.node)) { sideNode(box, d, node(d, sel.node)); return; }
    if (sel.link != null && (d.links || [])[sel.link]) { sideLink(box, d, d.links[sel.link]); return; }
    sideSummary(box, d);
  }

  function sideSummary(box, d) {
    const c = state.current;
    const h = el('h4');
    const sum = state.list.find((x) => x.id === c.id) || {};
    h.appendChild(el('span', 'chn-dot ' + statusClass(sum)));
    h.appendChild(document.createTextNode(t('chn_summary')));
    box.appendChild(h);
    box.appendChild(el('div', 'r', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) box.appendChild(el('div', 'r', t('chn_disabled')));
    if (d.description) box.appendChild(el('div', '', d.description));
    if (sum.checks_total != null) {
      const row = el('div', 'chn-row');
      if (sum.nodes_total) row.appendChild(el('span', 'chn-badge', fmt(t('chn_nodes_count'), { up: sum.nodes_up, n: sum.nodes_total })));
      row.appendChild(el('span', 'chn-badge', fmt(t('chn_ok_count'), { ok: sum.checks_ok, n: sum.checks_total })));
      box.appendChild(row);
    }
    if ((d.scenarios || []).length) {
      box.appendChild(el('h5', '', t('chn_scenarios')));
      d.scenarios.forEach((s, i) => {
        const row = el('div', 'chn-row');
        const b = el('span', 'chn-badge' + (s.id === c.scenario ? ' active' : ''), s.label || s.id);
        if (s.id !== c.scenario) b.style.borderLeft = '3px solid ' + SCN_COLORS[i % SCN_COLORS.length];
        row.appendChild(b);
        ruleText(row, s.when);
        box.appendChild(row);
      });
    }
    box.appendChild(el('h5', '', t('chn_nodes')));
    (d.nodes || []).forEach((n) => {
      const row = el('div', 'chn-chk click');
      row.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
      row.appendChild(el('span', 'l', n.label || n.id));
      row.appendChild(el('span', 'v', n.role || t('chn_kind_' + n.kind)));
      row.onclick = () => { state.sel = { node: n.id }; renderDiagramOnly(); renderSide(); };
      box.appendChild(row);
    });
  }

  function ruleText(parent, when) {
    const k = el('span', 'k', JSON.stringify(when || {}));
    k.style.fontFamily = 'var(--font-mono)';
    k.style.overflowWrap = 'anywhere';
    parent.appendChild(k);
  }

  function renderDiagramOnly() {
    const wrap = $('.chn-canvas');
    if (wrap && state.current) renderDiagram(wrap);
  }

  function sideNode(box, d, n) {
    const c = state.current;
    const ns = ((c.state || {}).nodes || {})[n.id];
    const h = el('h4');
    h.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
    h.appendChild(document.createTextNode(n.label || n.id));
    box.appendChild(h);
    box.appendChild(el('div', 'r', [t('chn_kind_' + n.kind), n.role, n.id].filter(Boolean).join(' · ')
      + (ns ? ' · ' + t(ns.state === 'up' ? 'chn_up' : 'chn_down') : '')));
    if (n.kind === 'server' && n.conn_id) {
      const conn = (state.conns || []).find((x) => x.id === n.conn_id);
      box.appendChild(el('div', 'r', t('chn_conn') + ': ' + (conn ? conn.name : '#' + n.conn_id)));
    }
    if (n.kind === 'storage' && n.storage_id) box.appendChild(el('div', 'r', t('chn_storage') + ': #' + n.storage_id));
    if (ns && ns.error) box.appendChild(el('div', 'e', ns.error));
    const act = el('div', 'act');
    const b = el('button', 'btn-sec btn-sm', '↻ ' + t('chn_check_now'));
    b.onclick = () => checkNow(b);
    act.appendChild(b);
    box.appendChild(act);
    const checks = (d.checks || []).filter((x) => x.node === n.id);
    box.appendChild(el('h5', '', t('chn_checks')));
    checks.forEach((x) => box.appendChild(checkRow(x, true)));
    const links = (d.links || []).map((l, i) => ({ l, i })).filter((o) => o.l.from === n.id || o.l.to === n.id);
    if (links.length) {
      box.appendChild(el('h5', '', t('chn_links')));
      links.forEach(({ l, i }) => {
        const row = el('div', 'chn-chk click');
        row.appendChild(el('span', 'chn-dot ' + dotOf(linkStatus(d, l))));
        row.appendChild(el('span', 'l', `${label(d, l.from)} → ${label(d, l.to)}`));
        row.appendChild(el('span', 'v', t('chn_lk_' + (l.kind || 'data'))));
        row.onclick = () => { state.sel = { link: i }; renderDiagramOnly(); renderSide(); };
        box.appendChild(row);
      });
    }
    const ids = new Set(checks.map((x) => x.id));
    history(box, state.timeline.filter((e) => (e.kind === 'node' && e.subject === n.id) || (e.kind === 'check' && ids.has(e.subject))));
  }

  function dotOf(st) { return { ok: 'ok', bad: 'bad', warn: 'warn' }[st] || ''; }

  function checkRow(x, expandable) {
    const wrap = el('div');
    const r = checkRes(x.id);
    const row = el('div', 'chn-chk' + (expandable ? ' click' : ''));
    row.appendChild(el('span', 'chn-dot ' + (r ? (r.ok ? 'ok' : 'bad') : '')));
    const l = el('span', 'l', (x.label || x.id) + ' · ' + x.type);
    l.title = x.target || '';
    row.appendChild(l);
    row.appendChild(el('span', 'v', r ? r.value : '—'));
    if (r && (r.detail || r.error)) row.appendChild(el('span', 'd', r.error || r.detail));
    wrap.appendChild(row);
    if (expandable) {
      row.tabIndex = 0;
      row.setAttribute('role', 'button');
      row.setAttribute('aria-expanded', state.openCheck === x.id ? 'true' : 'false');
      const toggle = () => { state.openCheck = state.openCheck === x.id ? null : x.id; renderSide(); };
      row.onclick = toggle;
      row.onkeydown = (e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggle(); } };
      if (state.openCheck === x.id) {
        const more = el('div', 'chn-chk-more');
        more.appendChild(el('div', '', t('chn_definition')));
        const def = Object.assign({}, x);
        delete def.node;
        more.appendChild(el('pre', '', JSON.stringify(def, null, 2)));
        if (r && r.at) more.appendChild(el('div', '', t('chn_last') + ': ' + when(r.at)));
        history(more, state.timeline.filter((e) => e.kind === 'check' && e.subject === x.id), true);
        wrap.appendChild(more);
      }
    }
    return wrap;
  }

  function history(box, list, small) {
    box.appendChild(el(small ? 'div' : 'h5', '', t('chn_history')));
    if (!list.length) { box.appendChild(el('div', 'k', t('chn_no_history'))); return; }
    list.slice(0, 30).forEach((e) => box.appendChild(timelineRow(e)));
  }

  function sideLink(box, d, l) {
    const st = linkStatus(d, l);
    const h = el('h4');
    h.appendChild(el('span', 'chn-dot ' + dotOf(st)));
    h.appendChild(document.createTextNode(`${label(d, l.from)} → ${label(d, l.to)}`));
    box.appendChild(h);
    box.appendChild(el('div', 'r', [t('chn_link'), t('chn_lk_' + (l.kind || 'data')), l.label].filter(Boolean).join(' · ')
      + (st === 'idle' ? ' · ' + t('chn_link_idle') : '')));
    if ((l.scenarios || []).length) {
      const row = el('div', 'chn-row');
      row.appendChild(el('span', 'k', t('chn_link_active_in') + ':'));
      l.scenarios.forEach((id) => {
        const s = (d.scenarios || []).find((x) => x.id === id);
        row.appendChild(el('span', 'chn-badge' + (id === state.current.scenario ? ' active' : ''), s ? (s.label || s.id) : id));
      });
      box.appendChild(row);
    }
    if (l.check) {
      const x = (d.checks || []).find((c) => c.id === l.check);
      if (x) {
        box.appendChild(el('h5', '', t('chn_link_check')));
        box.appendChild(checkRow(x, true));
      }
    }
    const act = el('div', 'act');
    [l.from, l.to].forEach((id) => {
      const b = el('button', 'btn-sec btn-sm', label(d, id));
      b.onclick = () => { state.sel = { node: id }; renderDiagramOnly(); renderSide(); };
      act.appendChild(b);
    });
    box.appendChild(act);
  }

  // ─── the main area ──────────────────────────────────────────────────────────────────────

  function renderDetail() {
    const c = state.current;
    const show = (id, on) => { const b = $(id); if (b) b.style.display = on ? '' : 'none'; };
    ['#chn-check', '#chn-edit', '#chn-export', '#chn-del', '#chn-scn', '#chn-mode'].forEach((id) => show(id, !!c));
    document.querySelectorAll('#chn-mode button').forEach((b) => { const on = b.dataset.m === state.mode; b.className = (on ? 'btn-pri' : 'btn-sec') + ' btn-sm'; b.setAttribute('aria-pressed', on ? 'true' : 'false'); });
    const body = $('#chn-body');
    if (!body) return;
    body.textContent = '';
    body.classList.toggle('dg', !!c && state.mode === 'diagram');
    $('#chn-name').textContent = c ? c.definition.name : t('chn_title');
    if (!c) { body.appendChild(el('div', 'chn-empty', state.list.length ? t('chn_pick') : t('chn_none'))); return; }
    const scn = $('#chn-scn');
    scn.textContent = t('chn_scenario') + ': ' + (c.scenario ? c.scenario_label : t('chn_no_scenario'));
    scn.className = 'chn-badge' + (c.scenario ? ' active' : '');
    if (state.mode === 'diagram') renderDiagramView(body); else renderListView(body);
    renderTimeline(state.timeline);
  }

  function renderDiagramView(body) {
    const top = el('div', 'chn-dg-top');
    const canvas = el('div', 'chn-canvas');
    const tools = el('div', 'chn-tools');
    const tool = (txt, title, fn) => { const b = el('button', 'btn-sec btn-sm', txt); b.title = title; b.setAttribute('aria-label', title); b.onclick = fn; tools.appendChild(b); };
    tool('−', t('chn_zoom_out'), () => zoom(1.25));
    tool('+', t('chn_zoom_in'), () => zoom(1 / 1.25));
    tool('⤢', t('chn_fit'), () => { delete state.views[state.current.id]; renderDiagramOnly(); });
    tool('⟲', t('chn_auto_layout'), resetLayout);
    canvas.appendChild(tools);
    const legend = el('div', 'chn-legend');
    [['#3fb950', t('chn_flow')], ['#f85149', t('chn_legend_bad')], ['#8892a0', t('chn_legend_idle')]].forEach(([col, txt]) => {
      const s = el('span'); const i = el('i'); i.style.borderColor = col; s.appendChild(i); s.appendChild(document.createTextNode(txt)); legend.appendChild(s);
    });
    canvas.appendChild(legend);
    canvas.appendChild(el('div', 'chn-hint', t('chn_drag_hint')));
    top.appendChild(canvas);
    const det = el('div', 'chn-det');
    det.id = 'chn-det';
    top.appendChild(det);
    body.appendChild(top);
    const band = el('div', 'chn-band');
    band.id = 'chn-band';
    body.appendChild(band);
    renderDiagram(canvas);
    renderSide();
    const tl = el('div', 'chn-sec');
    tl.id = 'chn-tl';
    band.appendChild(tl);
  }

  function renderListView(body) {
    const c = state.current, d = c.definition;
    const meta = el('div', 'chn-row');
    meta.appendChild(el('span', 'k', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) meta.appendChild(el('span', 'chn-badge', t('chn_disabled')));
    if (d.description) meta.appendChild(el('span', '', d.description));
    body.appendChild(meta);
    const nodesSec = section(body, t('chn_nodes'));
    const grid = el('div', 'chn-nodes');
    nodesSec.appendChild(grid);
    (d.nodes || []).forEach((n) => {
      const ns = ((c.state || {}).nodes || {})[n.id];
      const card = el('div', 'chn-node');
      const h = el('div', 'h');
      h.appendChild(el('span', 'chn-dot ' + nodeStatus(d, n)));
      h.appendChild(document.createTextNode(n.label || n.id));
      card.appendChild(h);
      card.appendChild(el('div', 'r', [t('chn_kind_' + n.kind), n.role, n.id].filter(Boolean).join(' · ')
        + (ns ? ' · ' + t(ns.state === 'up' ? 'chn_up' : 'chn_down') : '')));
      if (ns && ns.error) card.appendChild(el('div', 'e', ns.error));
      (d.checks || []).filter((x) => x.node === n.id).forEach((x) => card.appendChild(checkRow(x, false)));
      grid.appendChild(card);
    });
    if ((d.links || []).length) {
      const ls = section(body, t('chn_links'));
      d.links.forEach((l) => {
        const row = el('div', 'chn-row');
        row.appendChild(el('span', 'chn-dot ' + dotOf(linkStatus(d, l))));
        row.appendChild(el('span', '', `${label(d, l.from)} → ${label(d, l.to)}`));
        row.appendChild(el('span', 'k', [t('chn_lk_' + (l.kind || 'data')), l.label].filter(Boolean).join(' · ')));
        ls.appendChild(row);
      });
    }
    if ((d.scenarios || []).length) {
      const ss = section(body, t('chn_scenarios'));
      d.scenarios.forEach((s) => {
        const row = el('div', 'chn-row');
        row.appendChild(el('span', 'chn-badge' + (s.id === c.scenario ? ' active' : ''), s.label || s.id));
        ruleText(row, s.when);
        ss.appendChild(row);
      });
    }
    const tl = section(body, t('chn_timeline'));
    tl.id = 'chn-tl';
  }

  function label(d, id) {
    const n = (d.nodes || []).find((x) => x.id === id);
    return n ? (n.label || n.id) : id;
  }

  function section(parent, title) {
    const s = el('div', 'chn-sec');
    s.appendChild(el('h4', '', title));
    parent.appendChild(s);
    return s;
  }

  // ─── timeline: the scenario band and the list of changes ───────────────────────────────

  function scenarioLabel(d, id) { const s = (d.scenarios || []).find((x) => x.id === id); return s ? (s.label || s.id) : (id || t('chn_no_scenario')); }
  function scenarioColor(d, id) { const i = (d.scenarios || []).findIndex((x) => x.id === id); return i < 0 ? '#8892a0' : SCN_COLORS[i % SCN_COLORS.length]; }

  function renderBand(d, list) {
    const band = $('#chn-band');
    if (!band) return;
    band.querySelectorAll('.chn-band-head').forEach((x) => x.remove());
    const head = el('div', 'chn-band-head');
    head.appendChild(el('h4', '', t('chn_band_t')));
    band.insertBefore(head, band.firstChild);
    const scn = list.filter((e) => e.kind === 'scenario').slice().reverse();
    if (!scn.length) { head.appendChild(el('div', 'chn-empty', t('chn_band_none'))); return; }
    const now = Date.now();
    const start = new Date(scn[0].at).getTime();
    const span = Math.max(1, now - start);
    const bar = el('div', 'chn-bar');
    bar.setAttribute('role', 'img');
    const seen = [];
    scn.forEach((e, i) => {
      const from = new Date(e.at).getTime();
      const to = i + 1 < scn.length ? new Date(scn[i + 1].at).getTime() : now;
      const seg = el('span');
      seg.style.flexGrow = String(Math.max(0, to - from));
      seg.style.background = scenarioColor(d, e.new);
      seg.title = scenarioLabel(d, e.new) + ' · ' + fmt(t('chn_since'), { at: when(e.at) });
      bar.appendChild(seg);
      if (seen.indexOf(e.new) < 0) seen.push(e.new);
    });
    bar.setAttribute('aria-label', t('chn_band_t') + ': ' + scn.map((e) => scenarioLabel(d, e.new)).join(' → '));
    head.appendChild(bar);
    const marks = el('div', 'chn-marks');
    list.filter((e) => e.kind !== 'scenario').slice(0, 200).forEach((e) => {
      const at = new Date(e.at).getTime();
      if (!(at >= start)) return;
      const m = el('span');
      m.style.left = (100 * (at - start) / span).toFixed(2) + '%';
      m.style.background = e.kind === 'node' ? (e.new === 'up' ? '#3fb950' : '#f85149') : (e.ok ? '#3fb950' : '#f85149');
      m.title = when(e.at) + ' · ' + changeText(d, e);
      marks.appendChild(m);
    });
    head.appendChild(marks);
    const axis = el('div', 'chn-axis');
    axis.appendChild(el('span', '', when(scn[0].at)));
    axis.appendChild(el('span', '', t('chn_now')));
    head.appendChild(axis);
    const keys = el('div', 'chn-keys');
    seen.forEach((id) => {
      const k = el('span'); const b = el('b'); b.style.background = scenarioColor(d, id);
      k.appendChild(b); k.appendChild(document.createTextNode(scenarioLabel(d, id))); keys.appendChild(k);
    });
    head.appendChild(keys);
  }

  function changeText(d, e) {
    let what = e.subject;
    if (e.kind === 'node') what = label(d, e.subject);
    if (e.kind === 'scenario') what = t('chn_scenario');
    if (e.kind === 'check') { const x = (d.checks || []).find((c) => c.id === e.subject); if (x) what = (x.label || x.id) + ' (' + label(d, x.node) + ')'; }
    const from = e.kind === 'scenario' ? (e.old ? scenarioLabel(d, e.old) : '—') : (e.old || '—');
    const to = e.kind === 'scenario' ? scenarioLabel(d, e.new) : e.new;
    return `${what}: ${from} → ${to}`;
  }

  function timelineRow(e) {
    const d = state.current.definition;
    const row = el('div', 'chn-tl');
    row.appendChild(el('span', 'a', when(e.at)));
    if (e.kind === 'check') row.appendChild(el('span', 'chn-dot ' + (e.ok ? 'ok' : 'bad')));
    if (e.kind === 'node') row.appendChild(el('span', 'chn-dot ' + (e.new === 'up' ? 'ok' : 'bad')));
    row.appendChild(el('span', '', changeText(d, e)));
    return row;
  }

  function renderTimeline(list) {
    if (!state.current) return;
    list = list || [];
    renderBand(state.current.definition, list);
    const box = $('#chn-tl');
    if (!box) return;
    box.textContent = '';
    box.appendChild(el('h4', '', t('chn_timeline')));
    if (!list.length) { box.appendChild(el('div', 'chn-empty', t('chn_timeline_none'))); return; }
    list.forEach((e) => box.appendChild(timelineRow(e)));
  }

  // ─── actions ────────────────────────────────────────────────────────────────────────────

  function checkNow(btn) {
    const c = state.current;
    if (!c) return;
    const label0 = btn.textContent;
    btn.disabled = true; btn.textContent = t('chn_checking');
    api('POST', `/api/chains/${c.id}/check`).then(() => loadList())
      .catch((e) => showToast(e.message, 'error'))
      .finally(() => { btn.disabled = false; btn.textContent = label0; });
  }

  function deleteChain() {
    const c = state.current;
    if (!c || !confirm(fmt(t('chn_delete_q'), { name: c.definition.name }))) return;
    api('DELETE', `/api/chains/${c.id}`).then(() => { delete state.views[c.id]; state.current = null; return loadList(); })
      .catch((e) => showToast(e.message, 'error'));
  }

  // An empty starting point for the JSON editor: one server, one service, two scenarios.
  function emptyChain() {
    return {
      enabled: true,
      definition: {
        name: 'New chain', interval_s: 30, confirm: 2,
        nodes: [{ id: 'app', kind: 'server', label: 'App server', role: 'app', conn_id: 0 }],
        checks: [{ id: 'web', node: 'app', type: 'systemd', target: 'nginx.service', label: 'Web server' }],
        links: [],
        scenarios: [
          { id: 'normal', label: 'Normal', when: { check: 'web', field: 'ok', op: '==', value: 'true' } },
          { id: 'degraded', label: 'Degraded', when: {} },
        ],
      },
    };
  }

  function editorPane(title) {
    const main = $('.chn-main');
    if (!main) return null;
    const old = $('.chn-editor');
    if (old) old.remove();
    const ed = el('div', 'chn-editor');
    const head = el('div', 'chn-editor-head');
    head.appendChild(el('h3', '', title));
    ed.appendChild(head);
    main.appendChild(ed);
    return ed;
  }

  // New chain: pick a generic template, choose a server or storage for each of its nodes,
  // then create it at once or review the JSON first.
  function openNew() {
    const ed = editorPane(t('chn_new'));
    if (!ed) return;
    const head = ed.querySelector('.chn-editor-head');
    const cancel = el('button', 'btn-sec btn-sm', t('chn_cancel'));
    cancel.onclick = () => ed.remove();
    head.appendChild(cancel);
    const body = el('div', 'chn-tpls');
    ed.appendChild(body);
    Promise.all([
      api('GET', '/api/chains/templates').catch(() => []),
      api('GET', '/api/connections').catch(() => []),
      api('GET', '/api/storages').catch(() => []),
    ]).then(([tpls, conns, storages]) => {
      body.appendChild(el('h4', '', t('chn_tpl_pick')));
      const grid = el('div', 'chn-tpl-grid');
      body.appendChild(grid);
      const card = (name, desc, extra, fn) => {
        const b = el('button', 'chn-tpl');
        b.type = 'button';
        b.appendChild(el('b', '', name));
        b.appendChild(el('span', '', desc));
        if (extra) { const x = el('span', '', ' · ' + extra); b.appendChild(x); }
        b.onclick = fn;
        grid.appendChild(b);
      };
      (tpls || []).forEach((tp) => {
        const k = 'chn_tpl_' + tp.id;
        card(LANGS.en[k] ? t(k) : tp.name, LANGS.en[k + '_d'] ? t(k + '_d') : tp.description,
          fmt(t('chn_nodes_n'), { n: tp.definition.nodes.length }), () => templateForm(ed, body, tp, conns || [], storages || []));
      });
      card(t('chn_tpl_empty'), t('chn_tpl_empty_d'), '', () => openEditor(null, emptyChain()));
    });
  }

  function templateForm(ed, body, tp, conns, storages) {
    body.textContent = '';
    const ssh = conns.filter((c) => (c.protocol || 'SSH').toUpperCase() === 'SSH');
    const def = JSON.parse(JSON.stringify(tp.definition));
    const form = el('div', 'chn-form');
    body.appendChild(form);
    const k = 'chn_tpl_' + tp.id;
    form.appendChild(el('h4', '', LANGS.en[k] ? t(k) : tp.name));
    const nameL = el('label', '', t('chn_tpl_name'));
    const name = el('input');
    name.value = def.name;
    nameL.appendChild(name);
    form.appendChild(nameL);
    const pick = (kind, items, labelOf) => {
      const nodes = def.nodes.filter((n) => n.kind === kind);
      if (!nodes.length) return [];
      form.appendChild(el('h5', '', t(kind === 'server' ? 'chn_tpl_servers' : 'chn_tpl_storages')));
      if (!items.length) form.appendChild(el('div', 'note', t(kind === 'server' ? 'chn_tpl_no_ssh' : 'chn_tpl_no_storage')));
      return nodes.map((n) => {
        const l = el('label', '', [n.label || n.id, n.role].filter(Boolean).join(' · '));
        const s = el('select');
        s.appendChild(el('option', '', t('chn_tpl_choose'))).value = '';
        items.forEach((it) => { const o = el('option', '', labelOf(it)); o.value = String(it.id); s.appendChild(o); });
        l.appendChild(s);
        form.appendChild(l);
        return { n, s };
      });
    };
    const picks = pick('server', ssh, (c) => c.name + (c.host ? ' (' + c.host + ')' : ''))
      .concat(pick('storage', storages, (s) => s.name + (s.bucket ? ' · ' + s.bucket : '')));
    form.appendChild(el('div', 'note', t('chn_tpl_after')));
    const errs = el('div', 'chn-errs');
    errs.setAttribute('role', 'status');
    form.appendChild(errs);
    const act = el('div', 'chn-row');
    const back = el('button', 'btn-sec btn-sm', t('chn_back'));
    back.onclick = () => { ed.remove(); openNew(); };
    const json = el('button', 'btn-sec btn-sm', '{ } ' + t('chn_edit'));
    const create = el('button', 'btn-pri btn-sm', t('chn_create'));
    act.appendChild(back); act.appendChild(json); act.appendChild(create);
    form.appendChild(act);
    const build = (strict) => {
      let missing = false;
      picks.forEach(({ n, s }) => {
        const v = Number(s.value) || 0;
        if (!v) missing = true;
        if (n.kind === 'server') n.conn_id = v; else n.storage_id = v;
      });
      def.name = name.value.trim() || def.name;
      if (strict && missing) { errs.textContent = t('chn_tpl_missing'); return null; }
      return { definition: def, enabled: true };
    };
    json.onclick = () => { const v = build(false); if (v) openEditor(null, v); };
    create.onclick = () => {
      const v = build(true);
      if (!v) return;
      create.disabled = true;
      api('POST', '/api/chains', v).then((c) => {
        ed.remove();
        showToast(t('chn_saved'), 'success');
        state.current = c;
        state.sel = null;
        return loadList();
      }).catch((e) => {
        errs.textContent = '';
        (e.errors.length ? e.errors : [e.message]).forEach((x) => errs.appendChild(el('div', '', '• ' + x)));
      }).finally(() => { create.disabled = false; });
    };
  }

  function importFile() {
    const inp = document.createElement('input');
    inp.type = 'file'; inp.accept = '.json,application/json';
    inp.onchange = () => {
      const f = inp.files && inp.files[0];
      if (!f) return;
      f.text().then((txt) => {
        let data = null;
        try { data = JSON.parse(txt); } catch (_) { data = null; }
        if (!data || typeof data !== 'object') { showToast(t('chn_import_bad'), 'error'); return; }
        if (!data.definition && data.nodes) data = { definition: data, enabled: true };
        if (!data.definition) { showToast(t('chn_import_bad'), 'error'); return; }
        openEditor(null, data);
      });
    };
    inp.click();
  }

  function openEditor(id, value) {
    const ed = editorPane(t('chn_editor_t'));
    if (!ed) return;
    ed.querySelector('.chn-editor-head').insertAdjacentHTML('beforeend', `
        <label style="font-size:11.5px;display:flex;gap:6px;align-items:center;"><input type="checkbox" id="chn-ed-en"> ${esc(t('chn_enabled'))}</label>
        <button class="btn-sec btn-sm" id="chn-ed-val">${esc(t('chn_validate'))}</button>
        <button class="btn-sec btn-sm" id="chn-ed-cancel">${esc(t('chn_cancel'))}</button>
        <button class="btn-pri btn-sm" id="chn-ed-save">${esc(t('chn_save'))}</button>`);
    ed.insertAdjacentHTML('beforeend', `<div class="chn-editor-body"><textarea id="chn-ed-text" spellcheck="false" aria-label="${esc(t('chn_editor_t'))}"></textarea>
        <div class="chn-help" id="chn-ed-help"><b>${esc(t('chn_help'))}</b></div></div>
      <div class="chn-errs" id="chn-ed-errs" role="status" aria-live="polite"></div>`);
    ed.querySelector('#chn-ed-text').value = JSON.stringify(value.definition, null, 2);
    ed.querySelector('#chn-ed-en').checked = value.enabled !== false;
    ed.querySelector('#chn-ed-cancel').onclick = () => ed.remove();
    ed.querySelector('#chn-ed-val').onclick = () => submit(ed, id, true);
    ed.querySelector('#chn-ed-save').onclick = () => submit(ed, id, false);
    fillHelp(ed.querySelector('#chn-ed-help'));
  }

  function showErrors(ed, list, okText) {
    const box = ed.querySelector('#chn-ed-errs');
    box.textContent = '';
    if (okText) { box.style.color = '#3fb950'; box.appendChild(el('div', '', okText)); return; }
    box.style.color = '';
    list.forEach((e) => box.appendChild(el('div', '', '• ' + e)));
  }

  function submit(ed, id, validateOnly) {
    let def = null;
    try { def = JSON.parse(ed.querySelector('#chn-ed-text').value); } catch (e) { showErrors(ed, [String(e.message || e)]); return; }
    const body = { definition: def, enabled: ed.querySelector('#chn-ed-en').checked };
    if (validateOnly) {
      api('POST', '/api/chains/validate', body).then((r) => {
        if (r.ok) showErrors(ed, [], t('chn_valid')); else showErrors(ed, r.errors || []);
      }).catch((e) => showErrors(ed, e.errors.length ? e.errors : [e.message]));
      return;
    }
    const req = id ? api('PUT', `/api/chains/${id}`, body) : api('POST', '/api/chains', body);
    req.then((c) => {
      ed.remove();
      showToast(t('chn_saved'), 'success');
      state.current = c;
      state.sel = null;
      return loadList();
    }).catch((e) => showErrors(ed, e.errors.length ? e.errors : [e.message]));
  }

  function fillHelp(box) {
    const add = (title, lines) => {
      box.appendChild(el('h5', '', title));
      lines.forEach((l) => { const d = el('div'); const c = el('code', '', l[0]); d.appendChild(c); if (l[1]) d.appendChild(document.createTextNode(' ' + l[1])); box.appendChild(d); });
    };
    Promise.all([
      api('GET', '/api/connections').catch(() => []),
      api('GET', '/api/storages').catch(() => []),
      api('GET', '/api/chains/catalog').catch(() => ({ types: {} })),
    ]).then(([conns, storages, cat]) => {
      const ssh = (conns || []).filter((c) => (c.protocol || 'SSH').toUpperCase() === 'SSH');
      add(t('chn_help_conns'), ssh.slice(0, 200).map((c) => [String(c.id), c.name]));
      if ((storages || []).length) add(t('chn_help_storages'), storages.map((s) => [String(s.id), s.name]));
      add(t('chn_help_types'), Object.keys(cat.types || {}).sort().map((k) => [k, LANGS.en['chn_type_' + k] ? t('chn_type_' + k) : cat.types[k].help]));
      box.appendChild(el('h5', '', t('chn_scenarios')));
      box.appendChild(el('div', '', t('chn_help_rules')));
      box.appendChild(el('h5', '', t('chn_links')));
      box.appendChild(el('div', '', t('chn_help_links')));
    });
  }

  // A poll finished somewhere: refresh the list, and the details when that chain is open.
  function changed(id) {
    if (!isOpen()) return;
    if ($('.chn-editor')) { loadListOnly(); return; }
    if (state.current && state.current.id === id) loadList(); else loadListOnly();
  }

  function loadListOnly() {
    api('GET', '/api/chains').then((list) => { state.list = list || []; renderList(); }).catch(() => {});
  }

  window.wrmChains = { open, close, changed };
})();
