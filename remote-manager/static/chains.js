/* Copyright 2026 vedranius and WRM contributors
   SPDX-License-Identifier: Apache-2.0 */

// Chains: live monitoring of a service chain (chains.go). Loaded on demand by openChains() in
// index.html; uses its helpers (t, esc, showToast). Step 2: a live SVG diagram with automatic
// layout and saved positions, animated flows, storage nodes, details on click, the timeline and
// a generic starter template. The JSON editor from step 1 is still the way to edit a chain.
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
  chn_diagram: 'Live diagram', chn_auto_layout: 'Auto layout', chn_layout_saved: 'The layout was saved.', chn_layout_failed: 'The layout was not saved.',
  chn_select_hint: 'Choose a node or a check to see its details here.', chn_node_details: 'Node details', chn_check_details: 'Check details',
  chn_target: 'Target', chn_type: 'Type', chn_value: 'Value', chn_detail: 'Detail', chn_error: 'Error', chn_kind: 'Kind', chn_role: 'Role',
  chn_conn: 'Connection', chn_storage: 'Storage', chn_no_checks: 'No checks on this node.', chn_timeline_hint: 'Choose a change to see its node.',
  chn_zoom_in: 'Zoom in', chn_zoom_out: 'Zoom out', chn_zoom_fit: 'Fit', chn_checks_short: '{ok}/{n} OK', chn_links_count: '{n} links',
  chn_legend: 'Green: OK · yellow: warning · red: failing · grey: not checked yet',
  chn_check_ok: 'OK', chn_check_bad: 'failing', chn_check_unknown: 'not checked',
  chn_check_this: 'Check this chain now', chn_storage_role: 'bucket', chn_flow_hint: 'Arrows show the flow between nodes; drag a node to move it.',
  chn_node_up: 'up', chn_node_down: 'down', chn_node_unknown: 'not checked',
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
  chn_diagram: 'Dijagram uživo', chn_auto_layout: 'Automatski raspored', chn_layout_saved: 'Raspored je spremljen.', chn_layout_failed: 'Raspored nije spremljen.',
  chn_select_hint: 'Odaberite čvor ili provjeru da ovdje vidite detalje.', chn_node_details: 'Detalji čvora', chn_check_details: 'Detalji provjere',
  chn_target: 'Cilj', chn_type: 'Vrsta', chn_value: 'Vrijednost', chn_detail: 'Detalj', chn_error: 'Greška', chn_kind: 'Vrsta', chn_role: 'Uloga',
  chn_conn: 'Veza', chn_storage: 'Storage', chn_no_checks: 'Ovaj čvor nema provjera.', chn_timeline_hint: 'Odaberite promjenu da vidite njezin čvor.',
  chn_zoom_in: 'Povećaj', chn_zoom_out: 'Smanji', chn_zoom_fit: 'Prilagodi', chn_checks_short: '{ok}/{n} u redu', chn_links_count: '{n} veza',
  chn_legend: 'Zeleno: u redu · žuto: upozorenje · crveno: ne radi · sivo: još nije provjereno',
  chn_check_ok: 'u redu', chn_check_bad: 'nije u redu', chn_check_unknown: 'nije provjereno',
  chn_check_this: 'Provjeri ovaj lanac sada', chn_storage_role: 'bucket', chn_flow_hint: 'Strelice pokazuju tok između čvorova; povucite čvor da ga premjestite.',
  chn_node_up: 'radi', chn_node_down: 'ne radi', chn_node_unknown: 'nije provjereno',
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

(function () {
  const state = { list: [], current: null, timeline: [], styleAdded: false, pos: {}, selected: null, zoom: 1, drag: null, saveTimer: null, moved: false };
  const $ = (sel) => document.querySelector(sel);
  const el = (tag, cls, text) => { const e = document.createElement(tag); if (cls) e.className = cls; if (text != null) e.textContent = text; return e; };
  const fmt = (s, v) => s.replace(/\{(\w+)\}/g, (_, k) => (v[k] != null ? String(v[k]) : ''));
  const when = (at) => (at ? at.slice(0, 19).replace('T', ' ') : t('chn_not_checked'));
  const SVGNS = 'http://www.w3.org/2000/svg';
  const NODE_W = 182, NODE_H = 78, COL = 270, ROW = 126, PAD = 44;

  function css() {
    if (state.styleAdded) return;
    state.styleAdded = true;
    const s = document.createElement('style');
    s.textContent = `
.chn-wrap { display:flex; width:min(1340px,97vw); height:88vh; padding:0; overflow:hidden; }
.chn-side { width:240px; min-width:190px; border-right:1px solid var(--border2); display:flex; flex-direction:column; }
.chn-side-head { padding:14px 14px 8px; display:flex; align-items:center; justify-content:space-between; gap:8px; flex-wrap:wrap; }
.chn-side-head h3 { margin:0; font-size:14px; }
.chn-list { flex:1; overflow:auto; padding:6px; }
.chn-item { padding:9px 10px; border-radius:8px; cursor:pointer; border:1px solid transparent; margin-bottom:4px; }
.chn-item:hover, .chn-item.sel { background:var(--bg3); }
.chn-item.sel { border-color:var(--accent); }
.chn-item .n { font-size:12.5px; font-weight:600; }
.chn-item .m { font-size:10.5px; color:var(--text3); margin-top:3px; }
.chn-dot { display:inline-block; width:8px; height:8px; border-radius:50%; background:var(--text3); opacity:.5; margin-right:7px; vertical-align:middle; flex:none; }
.chn-dot.ok { background:#3fb950; opacity:1; } .chn-dot.bad { background:#f85149; opacity:1; } .chn-dot.warn { background:#d29922; opacity:1; }
.chn-main { flex:1; display:flex; flex-direction:column; min-width:0; position:relative; }
.chn-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; align-items:center; gap:8px; flex-wrap:wrap; }
.chn-head h3 { margin:0; font-size:14px; flex:1; min-width:120px; }
.chn-body { flex:1; overflow:auto; padding:12px 14px; display:flex; flex-direction:column; gap:10px; }
.chn-empty { padding:30px 14px; text-align:center; color:var(--text3); font-size:12px; }
.chn-badge { font-size:11px; padding:2px 8px; border-radius:10px; background:var(--bg3); color:var(--text2); }
.chn-badge.active { background:var(--accent); color:#fff; }
.chn-row { display:flex; gap:8px; align-items:center; font-size:11.5px; padding:2px 0; flex-wrap:wrap; }
.chn-row .k { color:var(--text3); font-size:10.5px; }
.chn-meta { font-size:11.5px; color:var(--text2); display:flex; gap:10px; align-items:center; flex-wrap:wrap; }
.chn-stage { display:flex; gap:10px; min-height:340px; flex:1; align-items:stretch; }
.chn-canvas { flex:1; min-width:0; position:relative; border:1px solid var(--border2); border-radius:10px; background:var(--bg1); overflow:auto; }
.chn-canvas svg { display:block; width:100%; height:100%; min-height:330px; touch-action:none; user-select:none; }
.chn-zoom { position:absolute; right:8px; top:8px; display:flex; gap:4px; background:var(--bg2); border:1px solid var(--border2); border-radius:8px; padding:3px; }
#chn-overlay .chn-zoom button { border:0; background:transparent; color:var(--text2); font-size:13px; width:24px; height:24px; padding:0; border-radius:5px; cursor:pointer; box-shadow:none; }
#chn-overlay .chn-zoom button:hover { background:var(--bg3); color:var(--text); }
.chn-inspector { width:290px; flex:none; border:1px solid var(--border2); border-radius:10px; background:var(--bg2); padding:10px 12px; overflow:auto; font-size:11.5px; }
.chn-inspector h4 { margin:0 0 4px; font-size:13px; }
.chn-inspector h5 { margin:12px 0 4px; font-size:10.5px; text-transform:uppercase; letter-spacing:.05em; color:var(--text3); }
.chn-inspector .sub { color:var(--text3); font-size:10.5px; margin-bottom:8px; }
.chn-inspector .hint { color:var(--text3); font-size:11.5px; }
.chn-chk { display:flex; flex-wrap:wrap; align-items:baseline; gap:2px 6px; font-size:11.5px; padding:5px 6px; border-top:1px dashed var(--border2); cursor:pointer; border-radius:6px; }
.chn-chk:hover { background:var(--bg3); }
.chn-chk .l { flex:1; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .v { font-family:var(--font-mono); font-size:11px; max-width:45%; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .d { flex-basis:100%; font-size:10px; color:var(--text3); padding-left:14px; }
.chn-kv { display:grid; grid-template-columns:auto 1fr; gap:3px 8px; font-size:11px; }
.chn-kv .k { color:var(--text3); }
.chn-kv .v { font-family:var(--font-mono); overflow-wrap:anywhere; }
.chn-tlsec { border:1px solid var(--border2); border-radius:10px; background:var(--bg2); }
.chn-tlhead { display:flex; gap:8px; align-items:center; padding:7px 12px; border-bottom:1px solid var(--border2); }
.chn-tlhead h4 { margin:0; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text3); flex:1; }
.chn-tlhead .hint { font-size:10px; color:var(--text3); }
.chn-tllist { max-height:170px; overflow:auto; padding:4px 6px; }
.chn-tl { font-size:11px; padding:3px 6px; display:flex; gap:8px; cursor:pointer; border-radius:6px; }
.chn-tl:hover { background:var(--bg3); }
.chn-tl .a { color:var(--text3); width:128px; flex:none; font-family:var(--font-mono); font-size:10.5px; }
.chn-tl .t { flex:1; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-legend { font-size:10px; color:var(--text3); }
.chn-editor { position:absolute; inset:0; background:var(--bg2); display:flex; flex-direction:column; z-index:5; }
.chn-editor-head { padding:10px 14px; border-bottom:1px solid var(--border2); display:flex; gap:8px; align-items:center; flex-wrap:wrap; }
.chn-editor-head h3 { margin:0; font-size:13px; flex:1; min-width:160px; }
.chn-editor-body { flex:1; display:flex; min-height:0; }
.chn-editor textarea { flex:1; border:0; border-radius:0; font-family:var(--font-mono); font-size:12px; resize:none; padding:12px 14px; background:var(--bg1); color:var(--text); min-width:0; }
.chn-help { width:320px; overflow:auto; border-left:1px solid var(--border2); padding:10px 12px; font-size:11px; color:var(--text2); }
.chn-help h5 { margin:10px 0 4px; font-size:10.5px; text-transform:uppercase; color:var(--text3); }
.chn-help code { font-family:var(--font-mono); font-size:10.5px; }
.chn-errs { max-height:120px; overflow:auto; padding:6px 14px; font-size:11px; color:#f85149; border-top:1px solid var(--border2); }
.chn-errs:empty { display:none; }
/* Diagram */
.chn-svg .box { fill:var(--bg2); stroke:var(--border2); stroke-width:1.2; }
.chn-svg .rim { fill:none; stroke:var(--border2); stroke-width:1; opacity:.85; }
.chn-svg .node-g { cursor:grab; outline:none; }
.chn-svg .node-g:active { cursor:grabbing; }
.chn-svg .node-g.sel .box { stroke:var(--accent); stroke-width:2.2; }
.chn-svg .node-g.ok .box { stroke:#3fb950; } .chn-svg .node-g.bad .box { stroke:#f85149; } .chn-svg .node-g.warn .box { stroke:#d29922; }
.chn-svg .bar { fill:var(--text3); opacity:.35; }
.chn-svg .node-g.ok .bar { fill:#3fb950; opacity:1; } .chn-svg .node-g.bad .bar { fill:#f85149; opacity:1; } .chn-svg .node-g.warn .bar { fill:#d29922; opacity:1; }
.chn-svg .icon { font-size:13px; fill:var(--text2); }
.chn-svg .label { font-size:12px; font-weight:600; fill:var(--text); }
.chn-svg .sub { font-size:9.5px; fill:var(--text3); }
.chn-svg .sum { font-size:9.5px; fill:var(--text3); font-family:var(--font-mono); }
.chn-svg .link { fill:none; stroke:var(--border2); stroke-width:1.7; }
.chn-svg .link.ok { stroke:#3fb950; } .chn-svg .link.bad { stroke:#f85149; } .chn-svg .link.warn { stroke:#d29922; }
.chn-svg .link.failover { stroke-dasharray:9 6; }
.chn-svg .link.dependency { stroke-dasharray:2 5; }
.chn-svg .link.flow { stroke-dasharray:7 5; animation:chn-flow 1.3s linear infinite; }
.chn-svg .link.failover.flow { animation-duration:2.1s; }
.chn-svg .link.bad.flow { animation:none; stroke-dasharray:5 4; }
@keyframes chn-flow { to { stroke-dashoffset: -24; } }
@media (prefers-reduced-motion: reduce) { .chn-svg .link.flow { animation:none; } }
@media (max-width: 820px) {
  .chn-wrap { flex-direction:column; }
  .chn-side { width:auto; max-height:30vh; border-right:0; border-bottom:1px solid var(--border2); }
  .chn-help { display:none; }
  .chn-stage { flex-direction:column; }
  .chn-inspector { width:auto; max-height:34vh; }
  .chn-canvas svg { min-height:280px; }
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
    ov.innerHTML = `<div class="modal chn-wrap">
      <div class="chn-side"><div class="chn-side-head"><h3>⛓ ${esc(t('chn_title'))}</h3>
        <span><button class="btn-sec btn-sm" id="chn-import">${esc(t('chn_import'))}</button>
        <button class="btn-sec btn-sm" id="chn-new">＋ ${esc(t('chn_new'))}</button></span></div>
        <div class="chn-list" id="chn-list"></div></div>
      <div class="chn-main"><div class="chn-head"><h3 id="chn-name"></h3>
          <span class="chn-badge" id="chn-scn" style="display:none;"></span>
          <button class="btn-sec btn-sm" id="chn-check" style="display:none;">↻ ${esc(t('chn_check_now'))}</button>
          <button class="btn-sec btn-sm" id="chn-layout" style="display:none;" title="${esc(t('chn_auto_layout'))}">⌗ ${esc(t('chn_auto_layout'))}</button>
          <button class="btn-sec btn-sm" id="chn-edit" style="display:none;">{ } ${esc(t('chn_edit'))}</button>
          <button class="btn-sec btn-sm" id="chn-export" style="display:none;">⬇ ${esc(t('chn_export'))}</button>
          <button class="btn-sec btn-sm" id="chn-del" style="display:none;">${esc(t('chn_delete'))}</button>
          <button class="btn-sec btn-sm" id="chn-close" aria-label="${esc(t('chn_cancel'))}">✕</button>
        </div>
        <div class="chn-body" id="chn-body"></div></div></div>`;
    document.body.appendChild(ov);
    ov.addEventListener('mousedown', (e) => { if (e.target === ov) close(); });
    ov.querySelector('#chn-close').onclick = close;
    ov.querySelector('#chn-new').onclick = newChain;
    ov.querySelector('#chn-import').onclick = importFile;
    ov.querySelector('#chn-check').onclick = checkNow;
    ov.querySelector('#chn-layout').onclick = () => autoLayout(true);
    ov.querySelector('#chn-edit').onclick = () => state.current && openEditor(state.current.id, { definition: state.current.definition, enabled: state.current.enabled });
    ov.querySelector('#chn-export').onclick = () => { if (state.current) location.href = `/api/chains/${state.current.id}/export`; };
    ov.querySelector('#chn-del').onclick = deleteChain;
    document.addEventListener('keydown', (e) => {
      if (e.key !== 'Escape' || !isOpen()) return;
      const ed = $('.chn-editor');
      if (ed) { ed.remove(); return; }
      close();
    });
    return ov;
  }

  function open() {
    css();
    shell().classList.add('open');
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
      if (state.current && state.list.some((c) => c.id === state.current.id)) return select(state.current.id);
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
      const n = el('div', 'n');
      n.appendChild(el('span', 'chn-dot ' + statusClass(c)));
      n.appendChild(document.createTextNode(c.name));
      const m = el('div', 'm');
      m.textContent = !c.enabled ? t('chn_disabled')
        : [c.scenario_label || t('chn_no_scenario'), c.nodes_total ? fmt(t('chn_nodes_count'), { up: c.nodes_up, n: c.nodes_total }) : '',
          fmt(t('chn_ok_count'), { ok: c.checks_ok, n: c.checks_total })].filter(Boolean).join(' · ');
      item.appendChild(n); item.appendChild(m);
      item.onclick = () => select(c.id);
      box.appendChild(item);
    });
  }

  function select(id) {
    return Promise.all([
      api('GET', `/api/chains/${id}`),
      api('GET', `/api/chains/${id}/timeline?limit=120`).catch(() => []),
    ]).then(([c, tl]) => {
      state.current = c;
      state.timeline = tl || [];
      if (state.selected) {
        const ok = state.selected.kind === 'node'
          ? (c.definition.nodes || []).some((n) => n.id === state.selected.id)
          : (c.definition.checks || []).some((x) => x.id === state.selected.id);
        if (!ok) state.selected = null;
      }
      state.pos = layoutPositions(c.definition);
      renderList();
      renderDetail();
      return c;
    }).catch((e) => showToast(e.message, 'error'));
  }

  function renderDetail() {
    const c = state.current;
    const show = (id, on) => { const b = $(id); if (b) b.style.display = on ? '' : 'none'; };
    ['#chn-check', '#chn-edit', '#chn-export', '#chn-del', '#chn-layout', '#chn-scn'].forEach((id) => show(id, !!c));
    const body = $('#chn-body');
    if (!body) return;
    body.textContent = '';
    $('#chn-name').textContent = c ? c.definition.name : t('chn_title');
    if (!c) { body.appendChild(el('div', 'chn-empty', state.list.length ? t('chn_pick') : t('chn_none'))); return; }
    const d = c.definition;
    const scn = $('#chn-scn');
    scn.textContent = t('chn_scenario') + ': ' + (c.scenario ? c.scenario_label : t('chn_no_scenario'));
    scn.className = 'chn-badge' + (c.scenario ? ' active' : '');

    const meta = el('div', 'chn-meta');
    meta.appendChild(el('span', '', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) meta.appendChild(el('span', 'chn-badge', t('chn_disabled')));
    if (d.description) meta.appendChild(el('span', '', d.description));
    meta.appendChild(el('span', 'chn-legend', t('chn_legend')));
    body.appendChild(meta);

    const stage = el('div', 'chn-stage');
    const canvas = el('div', 'chn-canvas');
    canvas.id = 'chn-canvas';
    stage.appendChild(canvas);
    const insp = el('aside', 'chn-inspector');
    insp.id = 'chn-inspector';
    stage.appendChild(insp);
    body.appendChild(stage);

    const tlsec = el('div', 'chn-tlsec');
    tlsec.innerHTML = `<div class="chn-tlhead"><h4>${esc(t('chn_timeline'))}</h4><span class="hint">${esc(t('chn_timeline_hint'))}</span></div>
      <div class="chn-tllist" id="chn-tl"></div>`;
    body.appendChild(tlsec);

    renderDiagram();
    renderInspector();
    renderTimeline();
  }

  // ─── layout ─────────────────────────────────────────────────────────────────────────────

  function autoPositions(def) {
    const nodes = def.nodes || [], links = def.links || [];
    const ids = nodes.map((n) => n.id);
    const rev = {};
    ids.forEach((id) => { rev[id] = []; });
    links.forEach((l) => { if (rev[l.to] && rev[l.from] !== undefined) rev[l.to].push(l.from); });
    const rank = {}, seen = {};
    const r = (id) => {
      if (rank[id] !== undefined) return rank[id];
      if (seen[id]) return 0;
      seen[id] = true;
      let best = 0;
      (rev[id] || []).forEach((p) => { best = Math.max(best, r(p) + 1); });
      delete seen[id];
      rank[id] = best;
      return best;
    };
    ids.forEach(r);
    const cols = {};
    ids.forEach((id) => { const k = rank[id] || 0; (cols[k] = cols[k] || []).push(id); });
    const keys = Object.keys(cols).map(Number).sort((a, b) => a - b);
    const order = {};
    ids.forEach((id, i) => { order[id] = i; });
    // Order each column by the average position of its predecessors, so links cross as little as
    // possible (a small barycentre pass; one rank is enough for a chain diagram).
    keys.forEach((k, ci) => {
      if (ci === 0) { cols[k].sort((a, b) => order[a] - order[b]); return; }
      const prev = cols[keys[ci - 1]];
      const bary = (id) => {
        const ps = (rev[id] || []).filter((p) => prev.indexOf(p) >= 0);
        if (!ps.length) return order[id];
        return ps.reduce((sum, p) => sum + prev.indexOf(p), 0) / ps.length;
      };
      cols[k].sort((a, b) => (bary(a) - bary(b)) || (order[a] - order[b]));
    });
    const rows = Math.max.apply(null, keys.map((k) => cols[k].length).concat([1]));
    const pos = {};
    keys.forEach((k, ci) => {
      const list = cols[k];
      list.forEach((id, ri) => {
        pos[id] = { x: PAD + ci * COL, y: PAD + (ri + (rows - list.length) / 2) * ROW };
      });
    });
    return pos;
  }

  // A saved position wins; a node without one (or still at 0,0) comes from the automatic layout.
  function layoutPositions(def) {
    const auto = autoPositions(def);
    const pos = {};
    (def.nodes || []).forEach((n) => {
      pos[n.id] = (n.x || n.y) ? { x: n.x, y: n.y } : (auto[n.id] || { x: PAD, y: PAD });
    });
    return pos;
  }

  function bounds() {
    const ps = Object.keys(state.pos).map((k) => state.pos[k]);
    let w = 420, h = 260;
    ps.forEach((p) => { w = Math.max(w, p.x + NODE_W + PAD); h = Math.max(h, p.y + NODE_H + PAD); });
    return { w, h };
  }

  // ─── diagram ────────────────────────────────────────────────────────────────────────────

  function nodeStatus(def, st, n) {
    const checks = (def.checks || []).filter((x) => x.node === n.id);
    const res = checks.map((x) => ((st.checks || {})[x.id])).filter(Boolean);
    const ns = (st.nodes || {})[n.id];
    if (ns) {
      if (ns.state !== 'up') return 'bad';
      return res.length && res.some((r) => !r.ok) ? 'warn' : 'ok';
    }
    if (res.length) {
      const ok = res.filter((r) => r.ok).length;
      return ok === res.length ? 'ok' : (ok === 0 ? 'bad' : 'warn');
    }
    return 'unknown';
  }

  function nodeChecks(def, n) { return (def.checks || []).filter((x) => x.node === n.id); }

  function nodeStateText(st, n) {
    const ns = (st.nodes || {})[n.id];
    if (ns) return t(ns.state === 'up' ? 'chn_node_up' : 'chn_node_down');
    return t('chn_node_unknown');
  }

  function linkState(c, l) {
    if (l.check) {
      const r = (c.state.checks || {})[l.check];
      return r ? (r.ok ? 'ok' : 'bad') : '';
    }
    const n = (c.state.nodes || {})[l.to];
    if (n) return n.state === 'up' ? 'ok' : 'bad';
    return '';
  }

  function svgEl(tag, attrs) {
    const e = document.createElementNS(SVGNS, tag);
    for (const k in attrs) e.setAttribute(k, attrs[k]);
    return e;
  }

  function renderDiagram() {
    const canvas = $('#chn-canvas');
    const c = state.current;
    if (!canvas || !c) return;
    canvas.textContent = '';
    const d = c.definition;
    const b = bounds();
    const svg = svgEl('svg', { class: 'chn-svg', viewBox: `0 0 ${b.w} ${b.h}`, role: 'img' });
    svg.setAttribute('aria-label', t('chn_diagram'));

    const defs = svgEl('defs');
    ['ok', 'bad', 'plain'].forEach((k) => {
      const m = svgEl('marker', { id: 'chn-arrow-' + k, viewBox: '0 0 10 10', refX: '9', refY: '5', markerWidth: '7', markerHeight: '7', orient: 'auto-start-reverse' });
      const p = svgEl('path', { d: 'M 0 0 L 10 5 L 0 10 z', fill: k === 'ok' ? '#3fb950' : k === 'bad' ? '#f85149' : 'currentColor', opacity: k === 'plain' ? '0.5' : '0.9' });
      m.appendChild(p); defs.appendChild(m);
    });
    svg.appendChild(defs);

    const g = svgEl('g', { id: 'chn-zoomg' });
    svg.appendChild(g);

    const linksG = svgEl('g', { class: 'chn-links' });
    const nodesG = svgEl('g', { class: 'chn-nodes' });
    g.appendChild(linksG); g.appendChild(nodesG);

    (d.links || []).forEach((l, i) => {
      if (!state.pos[l.from] || !state.pos[l.to]) return;
      const st = linkState(c, l);
      const p = svgEl('path', {
        class: 'chn-link link ' + (l.kind || '') + ' ' + st + ' flow',
        d: linkPath(state.pos[l.from], state.pos[l.to]),
        'marker-end': 'url(#chn-arrow-' + (st || 'plain') + ')',
        'data-from': l.from, 'data-to': l.to,
      });
      const title = svgEl('title');
      title.textContent = (l.label || ('link ' + (i + 1))) + (l.kind ? ' · ' + l.kind : '');
      p.appendChild(title);
      linksG.appendChild(p);
    });

    (d.nodes || []).forEach((n) => {
      const p = state.pos[n.id];
      if (!p) return;
      const st = nodeStatus(d, c.state || {}, n);
      const gnode = svgEl('g', { class: 'node-g ' + st + (state.selected && state.selected.kind === 'node' && state.selected.id === n.id ? ' sel' : ''), transform: `translate(${p.x},${p.y})`, 'data-id': n.id, tabindex: '0', role: 'button' });
      gnode.setAttribute('aria-label', (n.label || n.id) + ': ' + nodeStateText(c.state || {}, n));

      const title = svgEl('title');
      title.textContent = (n.label || n.id) + ' — ' + t('chn_kind_' + n.kind) + (n.role ? ' · ' + n.role : '');
      gnode.appendChild(title);

      gnode.appendChild(svgEl('rect', { class: 'bar', x: 2, y: 9, width: 4, height: NODE_H - 18, rx: 2 }));

      // A storage node is drawn as a bucket (a cylinder), a server as a box and an outside
      // endpoint as a diamond, so the diagram shows what a node is at a glance.
      const storage = n.kind === 'storage';
      if (storage) {
        gnode.appendChild(svgEl('path', { class: 'box', d: `M 0 10 L 0 ${NODE_H - 10} C 0 ${NODE_H} ${NODE_W} ${NODE_H} ${NODE_W} ${NODE_H - 10} L ${NODE_W} 10 C ${NODE_W} 0 0 0 0 10 Z` }));
        gnode.appendChild(svgEl('path', { class: 'rim', d: `M 0 10 C 0 20 ${NODE_W} 20 ${NODE_W} 10` }));
      } else {
        gnode.appendChild(svgEl('rect', { class: 'box', width: NODE_W, height: NODE_H, rx: 10 }));
      }

      const dy = storage ? 5 : 0;
      const icon = svgEl('text', { class: 'icon', x: 16, y: 25 + dy });
      icon.textContent = n.kind === 'storage' ? '▤' : n.kind === 'external' ? '◇' : '▣';
      gnode.appendChild(icon);

      const label = svgEl('text', { class: 'label', x: 34, y: 24 + dy });
      label.textContent = clip(n.label || n.id, 22);
      gnode.appendChild(label);

      const sub = svgEl('text', { class: 'sub', x: 34, y: 39 + dy });
      sub.textContent = [t('chn_kind_' + n.kind), n.role].filter(Boolean).join(' · ');
      gnode.appendChild(sub);

      const checks = nodeChecks(d, n);
      const ok = checks.filter((x) => ((c.state.checks || {})[x.id] || {}).ok).length;
      const sum = svgEl('text', { class: 'sum', x: 16, y: 61 + dy });
      const nodeState = (c.state.nodes || {})[n.id];
      sum.textContent = checks.length
        ? fmt(t('chn_checks_short'), { ok: ok, n: checks.length })
        : (nodeState ? nodeStateText(c.state || {}, n) : t('chn_no_checks'));
      gnode.appendChild(sum);

      gnode.addEventListener('pointerdown', (e) => startDrag(e, n.id, gnode, svg));
      gnode.addEventListener('click', () => { if (!state.moved) selectNode(n.id); });
      gnode.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); selectNode(n.id); }
      });
      nodesG.appendChild(gnode);
    });

    canvas.appendChild(svg);
    const zoom = el('div', 'chn-zoom');
    zoom.innerHTML = `<button type="button" aria-label="${esc(t('chn_zoom_out'))}" title="${esc(t('chn_zoom_out'))}">−</button>
      <button type="button" aria-label="${esc(t('chn_zoom_fit'))}" title="${esc(t('chn_zoom_fit'))}">⤢</button>
      <button type="button" aria-label="${esc(t('chn_zoom_in'))}" title="${esc(t('chn_zoom_in'))}">＋</button>`;
    const [out, fit, inn] = zoom.querySelectorAll('button');
    out.onclick = () => setZoom(state.zoom / 1.2);
    inn.onclick = () => setZoom(state.zoom * 1.2);
    fit.onclick = () => setZoom(1);
    canvas.appendChild(zoom);
    applyZoom(svg, b);
  }

  function clip(s, n) { s = String(s || ''); return s.length > n ? s.slice(0, n - 1) + '…' : s; }

  function linkPath(a, bNode) {
    const aC = { x: a.x + NODE_W / 2, y: a.y + NODE_H / 2 };
    const bC = { x: bNode.x + NODE_W / 2, y: bNode.y + NODE_H / 2 };
    const right = bC.x >= aC.x;
    const sx = right ? a.x + NODE_W : a.x;
    const ex = right ? bNode.x : bNode.x + NODE_W;
    const dx = Math.max(44, Math.abs(ex - sx) * 0.5);
    return `M ${sx} ${aC.y} C ${sx + (right ? dx : -dx)} ${aC.y}, ${ex + (right ? -dx : dx)} ${bC.y}, ${ex} ${bC.y}`;
  }

  function setZoom(z) {
    state.zoom = Math.max(0.4, Math.min(3, z));
    const svg = $('#chn-canvas svg');
    if (svg) applyZoom(svg, bounds());
  }

  function applyZoom(svg, b) {
    const z = state.zoom;
    const vw = b.w / z, vh = b.h / z;
    svg.setAttribute('viewBox', `${(b.w - vw) / 2} ${(b.h - vh) / 2} ${vw} ${vh}`);
  }

  function refreshLinks() {
    const svg = $('#chn-canvas svg');
    if (!svg) return;
    svg.querySelectorAll('.chn-link').forEach((p) => {
      const from = state.pos[p.getAttribute('data-from')], to = state.pos[p.getAttribute('data-to')];
      if (from && to) p.setAttribute('d', linkPath(from, to));
    });
  }

  function startDrag(e, id, gnode, svg) {
    if (e.button !== undefined && e.button !== 0) return;
    const zoomG = svg.querySelector('#chn-zoomg');
    if (!zoomG) return;
    const pt = (ev) => {
      const p = new DOMPoint(ev.clientX, ev.clientY);
      return p.matrixTransform(zoomG.getScreenCTM().inverse());
    };
    const start = pt(e);
    const origin = state.pos[id];
    state.drag = { id, gnode, offset: { x: start.x - origin.x, y: start.y - origin.y }, moved: false };
    state.moved = false;
    gnode.setPointerCapture(e.pointerId);
    const move = (ev) => {
      const p = pt(ev);
      const x = p.x - state.drag.offset.x, y = p.y - state.drag.offset.y;
      if (Math.abs(x - origin.x) + Math.abs(y - origin.y) > 3) { state.drag.moved = true; state.moved = true; }
      if (!state.drag.moved) return;
      state.pos[id] = { x, y };
      gnode.setAttribute('transform', `translate(${x},${y})`);
      refreshLinks();
    };
    const up = () => {
      gnode.removeEventListener('pointermove', move);
      gnode.removeEventListener('pointerup', up);
      gnode.removeEventListener('pointercancel', up);
      const moved = state.drag && state.drag.moved;
      state.drag = null;
      if (!moved) return;
      const p = state.pos[id];
      state.pos[id] = { x: Math.round(p.x / 10) * 10, y: Math.round(p.y / 10) * 10 };
      gnode.setAttribute('transform', `translate(${state.pos[id].x},${state.pos[id].y})`);
      refreshLinks();
      scheduleLayoutSave();
      setTimeout(() => { state.moved = false; }, 0);
    };
    gnode.addEventListener('pointermove', move);
    gnode.addEventListener('pointerup', up);
    gnode.addEventListener('pointercancel', up);
  }

  function scheduleLayoutSave() {
    const c = state.current;
    if (!c) return;
    clearTimeout(state.saveTimer);
    state.saveTimer = setTimeout(() => {
      const positions = {};
      (c.definition.nodes || []).forEach((n) => { if (state.pos[n.id]) positions[n.id] = state.pos[n.id]; });
      api('POST', `/api/chains/${c.id}/layout`, { positions }).then((got) => {
        if (got && got.definition) c.definition = got.definition;
        showToast(t('chn_layout_saved'), 'success');
      }).catch(() => showToast(t('chn_layout_failed'), 'error'));
    }, 350);
  }

  function autoLayout(save) {
    const c = state.current;
    if (!c) return;
    state.pos = autoPositions(c.definition);
    state.zoom = 1;
    renderDiagram();
    renderInspector();
    if (save) scheduleLayoutSave();
  }

  function selectNode(id) {
    state.selected = { kind: 'node', id };
    const svg = $('#chn-canvas svg');
    if (svg) svg.querySelectorAll('.node-g').forEach((g) => g.classList.toggle('sel', g.getAttribute('data-id') === id));
    renderInspector();
    const tl = $('#chn-tl');
    if (tl) Array.from(tl.children).forEach((r) => r.classList.toggle('hl', r.getAttribute('data-node') === id));
  }

  function selectCheck(id) {
    state.selected = { kind: 'check', id };
    const svg = $('#chn-canvas svg');
    if (svg) {
      const chk = (state.current.definition.checks || []).find((x) => x.id === id);
      svg.querySelectorAll('.node-g').forEach((g) => g.classList.toggle('sel', !!chk && g.getAttribute('data-id') === chk.node));
    }
    renderInspector();
  }

  // ─── inspector ──────────────────────────────────────────────────────────────────────────

  function renderInspector() {
    const box = $('#chn-inspector');
    const c = state.current;
    if (!box || !c) return;
    box.textContent = '';
    const sel = state.selected;
    if (!sel) {
      box.appendChild(el('div', 'hint', t('chn_select_hint')));
      box.appendChild(el('h5', '', t('chn_flow_hint')));
      return;
    }
    const d = c.definition;
    if (sel.kind === 'node') {
      const n = (d.nodes || []).find((x) => x.id === sel.id);
      if (!n) { state.selected = null; renderInspector(); return; }
      box.appendChild(el('h4', '', n.label || n.id));
      box.appendChild(el('div', 'sub', [t('chn_kind_' + n.kind), n.role, nodeStateText(c.state || {}, n)].filter(Boolean).join(' · ')));
      const ns = (c.state.nodes || {})[n.id];
      if (ns && ns.error) box.appendChild(el('div', 'chn-row', t('chn_error') + ': ' + ns.error));
      const kv = el('div', 'chn-kv');
      const add = (k, v) => { if (v === undefined || v === null || v === '') return; kv.appendChild(el('span', 'k', k)); kv.appendChild(el('span', 'v', String(v))); };
      add(t('chn_kind'), t('chn_kind_' + n.kind));
      if (n.role) add(t('chn_role'), n.role);
      if (n.kind === 'server') add(t('chn_conn'), '#' + n.conn_id);
      if (n.kind === 'storage') add(t('chn_storage'), '#' + n.storage_id);
      box.appendChild(kv);
      box.appendChild(el('h5', '', t('chn_checks')));
      const checks = nodeChecks(d, n);
      if (!checks.length) box.appendChild(el('div', 'hint', t('chn_no_checks')));
      checks.forEach((x) => box.appendChild(checkRow(c, x)));
    } else {
      const x = (d.checks || []).find((y) => y.id === sel.id);
      if (!x) { state.selected = null; renderInspector(); return; }
      const n = (d.nodes || []).find((y) => y.id === x.node);
      box.appendChild(el('h4', '', x.label || x.id));
      box.appendChild(el('div', 'sub', [n ? (n.label || n.id) : x.node, x.type].filter(Boolean).join(' · ')));
      const r = (c.state.checks || {})[x.id];
      const kv = el('div', 'chn-kv');
      const add = (k, v) => { if (v === undefined || v === null || v === '') return; kv.appendChild(el('span', 'k', k)); kv.appendChild(el('span', 'v', String(v))); };
      add(t('chn_type'), x.type);
      add(t('chn_target'), x.target);
      if (x.port) add('port', x.port);
      if (x.format) add('format', x.format);
      if (x.key) add('key', x.key);
      if (x.record) add('record', x.record);
      if (x.resolver) add('resolver', x.resolver);
      if (x.prefix) add('prefix', x.prefix);
      if (x.max_age_s) add('max_age_s', x.max_age_s);
      if (x.warn_days) add('warn_days', x.warn_days);
      if (x.expect) add('expect', JSON.stringify(x.expect));
      if (r) {
        add(t('chn_value'), r.value);
        add(t('chn_detail'), r.detail);
        if (r.error) add(t('chn_error'), r.error);
      }
      box.appendChild(kv);
      if (r) box.appendChild(el('div', 'sub', (r.ok ? '✓ ' + t('chn_check_ok') : '✕ ' + t('chn_check_bad'))));
    }
  }

  function checkRow(c, x) {
    const r = (c.state.checks || {})[x.id];
    const row = el('div', 'chn-chk');
    row.appendChild(el('span', 'chn-dot ' + (r ? (r.ok ? 'ok' : 'bad') : '')));
    const l = el('span', 'l', x.label || x.id);
    l.title = x.type + (x.target ? ' · ' + x.target : '');
    row.appendChild(l);
    row.appendChild(el('span', 'v', r ? r.value : '—'));
    if (r && (r.detail || r.error)) row.appendChild(el('span', 'd', r.error || r.detail));
    row.onclick = () => selectCheck(x.id);
    return row;
  }

  // ─── timeline ───────────────────────────────────────────────────────────────────────────

  function renderTimeline() {
    const box = $('#chn-tl');
    const c = state.current;
    if (!box || !c) return;
    box.textContent = '';
    if (!state.timeline.length) { box.appendChild(el('div', 'chn-empty', t('chn_timeline_none'))); return; }
    const d = c.definition;
    const label = (id) => { const n = (d.nodes || []).find((x) => x.id === id); return n ? (n.label || n.id) : id; };
    const sl = (id) => { const s = (d.scenarios || []).find((x) => x.id === id); return s ? (s.label || s.id) : (id || '—'); };
    state.timeline.forEach((e) => {
      const row = el('div', 'chn-tl');
      row.appendChild(el('span', 'a', when(e.at)));
      let what = e.subject;
      let target = null;
      if (e.kind === 'node') { what = label(e.subject); target = { kind: 'node', id: e.subject }; }
      if (e.kind === 'check') {
        const x = (d.checks || []).find((y) => y.id === e.subject);
        what = x ? (x.label || x.id) : e.subject;
        target = { kind: 'check', id: e.subject };
      }
      if (e.kind === 'scenario') what = t('chn_scenario');
      const from = e.kind === 'scenario' ? sl(e.old) : (e.old || '—');
      const to = e.kind === 'scenario' ? sl(e.new) : e.new;
      if (e.kind === 'check') row.appendChild(el('span', 'chn-dot ' + (e.ok ? 'ok' : 'bad')));
      const txt = el('span', 't', `${what}: ${from} → ${to}`);
      txt.title = txt.textContent;
      row.appendChild(txt);
      if (target) {
        row.setAttribute('data-node', target.kind === 'node' ? target.id : '');
        row.onclick = () => { (target.kind === 'node' ? selectNode : selectCheck)(target.id); };
        row.style.cursor = 'pointer';
      }
      box.appendChild(row);
    });
  }

  // ─── actions ────────────────────────────────────────────────────────────────────────────

  function checkNow() {
    const c = state.current;
    if (!c) return;
    const b = $('#chn-check');
    b.disabled = true; b.textContent = t('chn_checking');
    api('POST', `/api/chains/${c.id}/check`).then(() => loadList())
      .catch((e) => {
        // A background poll may be running; show its result instead of an error.
        if (e.status === 409) { setTimeout(() => loadList(), 1500); return; }
        showToast(e.message, 'error');
      })
      .finally(() => { b.disabled = false; b.textContent = '↻ ' + t('chn_check_now'); });
  }

  function deleteChain() {
    const c = state.current;
    if (!c || !confirm(fmt(t('chn_delete_q'), { name: c.definition.name }))) return;
    api('DELETE', `/api/chains/${c.id}`).then(() => { state.current = null; state.selected = null; return loadList(); })
      .catch((e) => showToast(e.message, 'error'));
  }

  function newChain() {
    api('GET', '/api/chains/template').then((tpl) => openEditor(null, tpl))
      .catch(() => openEditor(null, fallbackTemplate()));
  }

  // Only used when the template endpoint is not reachable: outside nodes need no saved connection.
  function fallbackTemplate() {
    return {
      enabled: true,
      definition: {
        name: 'New chain', interval_s: 30, confirm: 2,
        nodes: [
          { id: 'app', kind: 'external', label: 'App', role: 'app' },
          { id: 'edge', kind: 'external', label: 'Edge', role: 'cdn' },
        ],
        checks: [
          { id: 'app_dns', node: 'app', type: 'dns', target: 'app.example.com', label: 'DNS' },
          { id: 'edge_tls', node: 'edge', type: 'tls', target: 'example.com', label: 'Certificate' },
        ],
        links: [{ from: 'app', to: 'edge', kind: 'data', check: 'app_dns' }],
        scenarios: [
          { id: 'normal', label: 'Normal', when: { all: [{ check: 'app_dns', field: 'ok', op: '==', value: 'true' }, { check: 'edge_tls', field: 'ok', op: '==', value: 'true' }] } },
          { id: 'degraded', label: 'Degraded', when: {} },
        ],
      },
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
    const main = $('.chn-main');
    if (!main) return;
    const old = $('.chn-editor');
    if (old) old.remove();
    const ed = el('div', 'chn-editor');
    ed.innerHTML = `<div class="chn-editor-head"><h3>${esc(t('chn_editor_t'))}</h3>
        <label style="font-size:11.5px;display:flex;gap:6px;align-items:center;"><input type="checkbox" id="chn-ed-en"> ${esc(t('chn_enabled'))}</label>
        <button class="btn-sec btn-sm" id="chn-ed-val">${esc(t('chn_validate'))}</button>
        <button class="btn-sec btn-sm" id="chn-ed-cancel">${esc(t('chn_cancel'))}</button>
        <button class="btn-pri btn-sm" id="chn-ed-save">${esc(t('chn_save'))}</button></div>
      <div class="chn-editor-body"><textarea id="chn-ed-text" spellcheck="false" aria-label="${esc(t('chn_editor_t'))}"></textarea>
        <div class="chn-help" id="chn-ed-help"><b>${esc(t('chn_help'))}</b></div></div>
      <div class="chn-errs" id="chn-ed-errs" role="status" aria-live="polite"></div>`;
    main.appendChild(ed);
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
      state.selected = null;
      state.current = c;
      state.pos = layoutPositions(c.definition);
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
    });
  }

  // A poll finished somewhere: refresh the list, and the details when that chain is open.
  function changed(id) {
    if (!isOpen()) return;
    if ($('.chn-editor')) { loadListOnly(); return; }
    if (state.drag) return;
    if (state.current && state.current.id === id) select(id); else loadListOnly();
  }

  function loadListOnly() {
    api('GET', '/api/chains').then((list) => { state.list = list || []; renderList(); }).catch(() => {});
  }

  window.wrmChains = { open, close, changed };
})();
