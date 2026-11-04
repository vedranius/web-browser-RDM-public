/* Copyright 2026 vedranius and WRM contributors
   SPDX-License-Identifier: Apache-2.0 */

// Chains: live monitoring of a service chain (chains.go). Loaded on demand by openChains() in
// index.html; uses its helpers (t, esc, showToast). Step 1: list, details per node and check,
// scenarios, timeline and a JSON editor with validation, import and export.
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

(function () {
  const state = { list: [], current: null, styleAdded: false };
  const $ = (sel) => document.querySelector(sel);
  const el = (tag, cls, text) => { const e = document.createElement(tag); if (cls) e.className = cls; if (text != null) e.textContent = text; return e; };
  const fmt = (s, v) => s.replace(/\{(\w+)\}/g, (_, k) => (v[k] != null ? String(v[k]) : ''));
  const when = (at) => (at ? at.slice(0, 19).replace('T', ' ') : t('chn_not_checked'));

  function css() {
    if (state.styleAdded) return;
    state.styleAdded = true;
    const s = document.createElement('style');
    s.textContent = `
.chn-wrap { display:flex; width:min(1240px,96vw); height:86vh; padding:0; overflow:hidden; }
.chn-side { width:260px; min-width:200px; border-right:1px solid var(--border2); display:flex; flex-direction:column; }
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
.chn-head h3 { margin:0; font-size:14px; flex:1; min-width:140px; }
.chn-body { flex:1; overflow:auto; padding:12px 14px; }
.chn-empty { padding:30px 14px; text-align:center; color:var(--text3); font-size:12px; }
.chn-badge { font-size:11px; padding:2px 8px; border-radius:10px; background:var(--bg3); color:var(--text2); }
.chn-badge.active { background:var(--accent); color:#fff; }
.chn-sec { margin:4px 0 16px; }
.chn-sec h4 { margin:0 0 8px; font-size:11px; text-transform:uppercase; letter-spacing:.05em; color:var(--text3); }
.chn-nodes { display:grid; grid-template-columns:repeat(auto-fill,minmax(280px,1fr)); gap:10px; }
.chn-node { border:1px solid var(--border2); border-radius:10px; padding:10px 12px; background:var(--bg2); }
.chn-node .h { display:flex; align-items:center; gap:6px; font-weight:600; font-size:12.5px; }
.chn-node .r { font-size:10.5px; color:var(--text3); margin:2px 0 6px; }
.chn-node .e { font-size:10.5px; color:#f85149; margin-bottom:6px; }
.chn-chk { display:flex; flex-wrap:wrap; align-items:baseline; gap:2px 6px; font-size:11.5px; padding:3px 0; border-top:1px dashed var(--border2); }
.chn-chk .l { flex:1; min-width:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .v { font-family:var(--font-mono); font-size:11px; max-width:45%; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.chn-chk .d { flex-basis:100%; font-size:10px; color:var(--text3); padding-left:14px; }
.chn-row { display:flex; gap:8px; align-items:center; font-size:11.5px; padding:4px 0; flex-wrap:wrap; }
.chn-row .k { color:var(--text3); font-size:10.5px; }
.chn-tl { font-size:11px; padding:3px 0; display:flex; gap:8px; }
.chn-tl .a { color:var(--text3); width:130px; flex:none; font-family:var(--font-mono); font-size:10.5px; }
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
@media (max-width: 760px) { .chn-wrap { flex-direction:column; } .chn-side { width:auto; max-height:34vh; border-right:0; border-bottom:1px solid var(--border2); } .chn-help { display:none; } }
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
          <button class="btn-sec btn-sm" id="chn-edit" style="display:none;">{ } ${esc(t('chn_edit'))}</button>
          <button class="btn-sec btn-sm" id="chn-export" style="display:none;">⬇ ${esc(t('chn_export'))}</button>
          <button class="btn-sec btn-sm" id="chn-del" style="display:none;">${esc(t('chn_delete'))}</button>
          <button class="btn-sec btn-sm" id="chn-close" aria-label="${esc(t('chn_cancel'))}">✕</button>
        </div>
        <div class="chn-body" id="chn-body"></div></div></div>`;
    document.body.appendChild(ov);
    ov.addEventListener('mousedown', (e) => { if (e.target === ov) close(); });
    ov.querySelector('#chn-close').onclick = close;
    ov.querySelector('#chn-new').onclick = () => openEditor(null, template());
    ov.querySelector('#chn-import').onclick = importFile;
    ov.querySelector('#chn-check').onclick = checkNow;
    ov.querySelector('#chn-edit').onclick = () => state.current && openEditor(state.current.id, { definition: state.current.definition, enabled: state.current.enabled });
    ov.querySelector('#chn-export').onclick = () => { if (state.current) location.href = `/api/chains/${state.current.id}/export`; };
    ov.querySelector('#chn-del').onclick = deleteChain;
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
    return api('GET', `/api/chains/${id}`).then((c) => {
      state.current = c;
      renderList();
      renderDetail();
      return api('GET', `/api/chains/${id}/timeline`).then(renderTimeline).catch(() => {});
    }).catch((e) => showToast(e.message, 'error'));
  }

  function renderDetail() {
    const c = state.current;
    const show = (id, on) => { const b = $(id); if (b) b.style.display = on ? '' : 'none'; };
    ['#chn-check', '#chn-edit', '#chn-export', '#chn-del', '#chn-scn'].forEach((id) => show(id, !!c));
    const body = $('#chn-body');
    if (!body) return;
    body.textContent = '';
    $('#chn-name').textContent = c ? c.definition.name : t('chn_title');
    if (!c) { body.appendChild(el('div', 'chn-empty', state.list.length ? t('chn_pick') : t('chn_none'))); return; }
    const d = c.definition;
    const st = c.state || {};
    const scn = $('#chn-scn');
    scn.textContent = t('chn_scenario') + ': ' + (c.scenario ? c.scenario_label : t('chn_no_scenario'));
    scn.className = 'chn-badge' + (c.scenario ? ' active' : '');
    const meta = el('div', 'chn-row');
    meta.appendChild(el('span', 'k', t('chn_last') + ': ' + when(c.last_at)));
    if (!c.enabled) meta.appendChild(el('span', 'chn-badge', t('chn_disabled')));
    if (d.description) meta.appendChild(el('span', '', d.description));
    body.appendChild(meta);

    const nodesSec = section(body, t('chn_nodes'));
    const grid = el('div', 'chn-nodes');
    nodesSec.appendChild(grid);
    (d.nodes || []).forEach((n) => {
      const ns = (st.nodes || {})[n.id];
      const card = el('div', 'chn-node');
      const h = el('div', 'h');
      const checks = (d.checks || []).filter((x) => x.node === n.id);
      const results = checks.map((x) => (st.checks || {})[x.id]).filter(Boolean);
      let cls = '';
      if (ns) cls = ns.state === 'up' ? 'ok' : 'bad';
      if (cls !== 'bad' && results.length) cls = results.every((r) => r.ok) ? 'ok' : 'warn';
      h.appendChild(el('span', 'chn-dot ' + cls));
      h.appendChild(document.createTextNode(n.label || n.id));
      card.appendChild(h);
      card.appendChild(el('div', 'r', [t('chn_kind_' + n.kind), n.role, n.id].filter(Boolean).join(' · ')
        + (ns ? ' · ' + t(ns.state === 'up' ? 'chn_up' : 'chn_down') : '')));
      if (ns && ns.error) card.appendChild(el('div', 'e', ns.error));
      checks.forEach((x) => {
        const r = (st.checks || {})[x.id];
        const row = el('div', 'chn-chk');
        row.appendChild(el('span', 'chn-dot ' + (r ? (r.ok ? 'ok' : 'bad') : '')));
        const l = el('span', 'l', (x.label || x.id) + ' · ' + x.type);
        l.title = x.target || '';
        row.appendChild(l);
        row.appendChild(el('span', 'v', r ? r.value : '—'));
        if (r && (r.detail || r.error)) row.appendChild(el('span', 'd', r.error || r.detail));
        card.appendChild(row);
      });
      grid.appendChild(card);
    });

    if ((d.links || []).length) {
      const ls = section(body, t('chn_links'));
      d.links.forEach((l) => {
        const row = el('div', 'chn-row');
        const r = l.check ? (st.checks || {})[l.check] : null;
        row.appendChild(el('span', 'chn-dot ' + (r ? (r.ok ? 'ok' : 'bad') : '')));
        row.appendChild(el('span', '', `${label(d, l.from)} → ${label(d, l.to)}`));
        row.appendChild(el('span', 'k', [l.kind, l.label].filter(Boolean).join(' · ')));
        ls.appendChild(row);
      });
    }
    if ((d.scenarios || []).length) {
      const ss = section(body, t('chn_scenarios'));
      d.scenarios.forEach((s) => {
        const row = el('div', 'chn-row');
        row.appendChild(el('span', 'chn-badge' + (s.id === c.scenario ? ' active' : ''), s.label || s.id));
        row.appendChild(el('span', 'k', JSON.stringify(s.when || {})));
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

  function renderTimeline(list) {
    const box = $('#chn-tl');
    if (!box) return;
    while (box.children.length > 1) box.removeChild(box.lastChild);
    if (!list || !list.length) { box.appendChild(el('div', 'chn-empty', t('chn_timeline_none'))); return; }
    const d = state.current.definition;
    list.forEach((e) => {
      const row = el('div', 'chn-tl');
      row.appendChild(el('span', 'a', when(e.at)));
      let what = e.subject;
      if (e.kind === 'node') what = label(d, e.subject);
      if (e.kind === 'scenario') what = t('chn_scenario');
      const sl = (id) => { const s = (d.scenarios || []).find((x) => x.id === id); return s ? (s.label || s.id) : (id || '—'); };
      const from = e.kind === 'scenario' ? sl(e.old) : (e.old || '—');
      const to = e.kind === 'scenario' ? sl(e.new) : e.new;
      if (e.kind === 'check') row.appendChild(el('span', 'chn-dot ' + (e.ok ? 'ok' : 'bad')));
      row.appendChild(el('span', '', `${what}: ${from} → ${to}`));
      box.appendChild(row);
    });
  }

  function checkNow() {
    const c = state.current;
    if (!c) return;
    const b = $('#chn-check');
    b.disabled = true; b.textContent = t('chn_checking');
    api('POST', `/api/chains/${c.id}/check`).then(() => loadList())
      .catch((e) => showToast(e.message, 'error'))
      .finally(() => { b.disabled = false; b.textContent = '↻ ' + t('chn_check_now'); });
  }

  function deleteChain() {
    const c = state.current;
    if (!c || !confirm(fmt(t('chn_delete_q'), { name: c.definition.name }))) return;
    api('DELETE', `/api/chains/${c.id}`).then(() => { state.current = null; return loadList(); })
      .catch((e) => showToast(e.message, 'error'));
  }

  // A generic starting point: one server, one service, two scenarios.
  function template() {
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
      state.current = c;
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
    if (state.current && state.current.id === id) loadList(); else loadListOnly();
  }

  function loadListOnly() {
    api('GET', '/api/chains').then((list) => { state.list = list || []; renderList(); }).catch(() => {});
  }

  window.wrmChains = { open, close, changed };
})();
