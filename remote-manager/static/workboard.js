// WRM PRO Workboard (v12.6.9).
// A page of tiles over the window area, opened with the ▤ Workboard button. Every tile is an
// addon the user adds from the catalog (nothing is on by default): server status, server load,
// notifications, clock, recent Git commits, running tunnels, quick connect and notes. A tile can
// be pulled out of the Workboard into a panel of the page layout (static/layout.js), where it
// docks, stacks and floats like any other panel; hiding that panel puts it back on the Workboard.
// The tiles are stored per user on the server (/api/me/ui/workboard) with a copy in this browser.
(function () {
  'use strict';

  Object.assign(LANGS.en, {
    wb_title: 'Workboard', wb_title_t: 'Workboard: your tiles with server status, load, notifications, commits and more',
    wb_add: 'Add addon', wb_catalog: 'Addons', wb_catalog_hint: 'Add only what helps you. Every tile can be configured, resized, moved, pulled out into a panel and removed again.',
    wb_empty_title: 'The Workboard is empty', wb_empty: 'Add the addons you want to see here: server status, CPU, RAM and disks of chosen servers, notifications, clocks, recent commits and more. Nothing is turned on until you add it.',
    wb_close: 'Close the Workboard', wb_added: 'Added: {name}', wb_max: 'The Workboard holds at most {n} tiles.',
    wb_t_settings: 'Settings', wb_t_remove: 'Remove', wb_t_out: 'Pull out into a panel', wb_t_refresh: 'Refresh now', wb_t_size: 'Wide / narrow',
    wb_t_drag: 'Drag to reorder', wb_remove_q: 'Remove the tile {name}?', wb_out_here: 'This tile is shown as a panel ({where}).', wb_bring_back: 'Bring it back',
    wb_locked: 'The layout is locked: unlock it with the 🔒 button to move or pull out tiles.',
    wb_set_title: 'Title', wb_set_title_ph: 'Default: the addon name', wb_set_size: 'Width', wb_size_s: 'Narrow', wb_size_m: 'Wide',
    wb_set_conns: 'Servers', wb_set_filter: 'Filter…', wb_set_conns_all: 'None chosen: all servers with live status are shown.', wb_set_conns_need: 'Choose up to {n} servers.',
    wb_set_more: 'At most {n} servers.', wb_saved: 'Saved',
    wb_loading: 'Loading…', wb_never: 'not loaded yet', wb_updated: 'updated {when}', wb_none: 'Nothing to show.',
    // addons
    wb_a_status: 'Server status', wb_a_status_d: 'Up, down and slow servers from the live status checks, with latency and since when.',
    wb_status_off: 'Live status checks are turned off. An administrator turns them on in Settings (status checks).',
    wb_status_problems: 'Only show servers with problems', wb_status_all_ok: 'All {n} servers are up.', wb_status_none: 'No server with live status yet. Live status is shown for connections with checks turned on.',
    wb_a_metrics: 'Server load', wb_a_metrics_d: 'CPU, RAM, load and disk usage of the servers you choose, refreshed every minute while visible.',
    wb_metrics_choose: 'Choose servers in ⚙ settings (SSH or SFTP connections).', wb_m_cpu: 'CPU', wb_m_ram: 'RAM', wb_m_load: 'Load', wb_m_up: 'up {d}',
    wb_m_src_zabbix: 'Monitoring (Zabbix agent)', wb_m_src_ssh: 'read over SSH', wb_m_more_disks: '+{n} more',
    wb_metrics_help: 'Values come from Monitoring when the server has it set up, otherwise WRM reads them over SSH (Linux: /proc and df), as far as the administrator allows it.',
    wb_a_events: 'Notifications', wb_a_events_d: 'What WRM told you: servers going down, expiring keys, new versions, Git jobs. Also without e-mail or chat channels.',
    wb_ev_read_all: 'Mark all read', wb_ev_clear: 'Clear', wb_ev_clear_q: 'Clear all notifications?', wb_ev_none: 'No notifications.', wb_ev_unread: '{n} unread',
    wb_a_clock: 'Clock', wb_a_clock_d: 'Your local time and the time zones you work with.',
    wb_clock_zones: 'Time zones', wb_clock_zones_ph: 'e.g. UTC, Europe/London, America/New_York', wb_clock_zones_bad: 'Unknown time zone: {z}', wb_clock_seconds: 'Show seconds', wb_clock_local: 'Local',
    wb_a_commits: 'Recent commits', wb_a_commits_d: 'The last commits and pushes of a service from the Git workspace (GitLab, GitHub or Gitea).',
    wb_commits_service: 'Service (from the Git catalog)', wb_commits_branch: 'Branch', wb_commits_branch_ph: 'Empty: the branches the service tracks', wb_commits_count: 'How many',
    wb_commits_choose: 'Choose a service in ⚙ settings.', wb_commits_nogit: 'The Git workspace is not available to you.', wb_commits_noapps: 'The Git catalog has no services yet.', wb_commits_pushed: 'pushed', wb_commits_none: 'No commits in the last days.',
    wb_a_tunnels: 'Running tunnels', wb_a_tunnels_d: 'Tunnels that are running now, with their local address.', wb_tun_none: 'No tunnel is running.', wb_tun_open: 'Tunnels…',
    wb_a_quick: 'Quick connect', wb_a_quick_d: 'Buttons for the servers you open most often.', wb_quick_choose: 'Choose servers in ⚙ settings.',
    wb_a_notes: 'Notes', wb_a_notes_d: 'A short note to yourself, kept on the server.', wb_notes_ph: 'Write a note…',
    wb_unavailable: 'Not available', wb_add_btn: 'Add',
  });
  Object.assign(LANGS.hr, {
    wb_title: 'Workboard', wb_title_t: 'Workboard: tvoje pločice sa stanjem poslužitelja, opterećenjem, obavijestima, commitovima i drugim',
    wb_add: 'Dodaj addon', wb_catalog: 'Addoni', wb_catalog_hint: 'Dodaj samo ono što ti pomaže. Svaku pločicu možeš podesiti, proširiti, premjestiti, izvući u panel i opet maknuti.',
    wb_empty_title: 'Workboard je prazan', wb_empty: 'Dodaj addone koje želiš vidjeti ovdje: stanje poslužitelja, CPU, RAM i diskove odabranih poslužitelja, obavijesti, satove, zadnje commitove i drugo. Ništa nije uključeno dok ga ne dodaš.',
    wb_close: 'Zatvori Workboard', wb_added: 'Dodano: {name}', wb_max: 'Workboard ima najviše {n} pločica.',
    wb_t_settings: 'Postavke', wb_t_remove: 'Ukloni', wb_t_out: 'Izvuci u panel', wb_t_refresh: 'Osvježi sada', wb_t_size: 'Široko / usko',
    wb_t_drag: 'Povuci da promijeniš redoslijed', wb_remove_q: 'Ukloniti pločicu {name}?', wb_out_here: 'Ova pločica je prikazana kao panel ({where}).', wb_bring_back: 'Vrati je ovamo',
    wb_locked: 'Raspored je zaključan: otključaj ga gumbom 🔒 da bi micao ili izvlačio pločice.',
    wb_set_title: 'Naslov', wb_set_title_ph: 'Zadano: naziv addona', wb_set_size: 'Širina', wb_size_s: 'Usko', wb_size_m: 'Široko',
    wb_set_conns: 'Poslužitelji', wb_set_filter: 'Filtriraj…', wb_set_conns_all: 'Ništa odabrano: prikazuju se svi poslužitelji sa stanjem uživo.', wb_set_conns_need: 'Odaberi do {n} poslužitelja.',
    wb_set_more: 'Najviše {n} poslužitelja.', wb_saved: 'Spremljeno',
    wb_loading: 'Učitavam…', wb_never: 'još nije učitano', wb_updated: 'osvježeno {when}', wb_none: 'Nema ničega za prikaz.',
    wb_a_status: 'Stanje poslužitelja', wb_a_status_d: 'Poslužitelji koji rade, ne rade ili su spori, iz provjera stanja uživo, s latencijom i od kada.',
    wb_status_off: 'Provjere stanja uživo su isključene. Administrator ih uključuje u Postavkama (provjere stanja).',
    wb_status_problems: 'Prikaži samo poslužitelje s problemima', wb_status_all_ok: 'Svih {n} poslužitelja radi.', wb_status_none: 'Još nema poslužitelja sa stanjem uživo. Stanje se prikazuje za veze s uključenim provjerama.',
    wb_a_metrics: 'Opterećenje poslužitelja', wb_a_metrics_d: 'CPU, RAM, opterećenje i zauzeće diskova odabranih poslužitelja, osvježava se svake minute dok je vidljivo.',
    wb_metrics_choose: 'Odaberi poslužitelje u ⚙ postavkama (SSH ili SFTP veze).', wb_m_cpu: 'CPU', wb_m_ram: 'RAM', wb_m_load: 'Opterećenje', wb_m_up: 'radi {d}',
    wb_m_src_zabbix: 'Nadzor (Zabbix agent)', wb_m_src_ssh: 'pročitano preko SSH-a', wb_m_more_disks: '+{n} još',
    wb_metrics_help: 'Vrijednosti dolaze iz Nadzora kad ga poslužitelj ima postavljenog, inače ih WRM čita preko SSH-a (Linux: /proc i df), koliko administrator dopušta.',
    wb_a_events: 'Obavijesti', wb_a_events_d: 'Što ti je WRM javio: poslužitelji koji su pali, ključevi koji istječu, nove verzije, Git poslovi. I bez e-mail ili chat kanala.',
    wb_ev_read_all: 'Označi sve pročitanim', wb_ev_clear: 'Očisti', wb_ev_clear_q: 'Obrisati sve obavijesti?', wb_ev_none: 'Nema obavijesti.', wb_ev_unread: '{n} nepročitanih',
    wb_a_clock: 'Sat', wb_a_clock_d: 'Tvoje lokalno vrijeme i vremenske zone s kojima radiš.',
    wb_clock_zones: 'Vremenske zone', wb_clock_zones_ph: 'npr. UTC, Europe/London, America/New_York', wb_clock_zones_bad: 'Nepoznata vremenska zona: {z}', wb_clock_seconds: 'Prikaži sekunde', wb_clock_local: 'Lokalno',
    wb_a_commits: 'Zadnji commitovi', wb_a_commits_d: 'Zadnji commitovi i pushevi servisa iz Git radnog prostora (GitLab, GitHub ili Gitea).',
    wb_commits_service: 'Servis (iz Git kataloga)', wb_commits_branch: 'Grana', wb_commits_branch_ph: 'Prazno: grane koje servis prati', wb_commits_count: 'Koliko',
    wb_commits_choose: 'Odaberi servis u ⚙ postavkama.', wb_commits_nogit: 'Git radni prostor ti nije dostupan.', wb_commits_noapps: 'Git katalog još nema servisa.', wb_commits_pushed: 'push', wb_commits_none: 'Nema commitova u zadnjih nekoliko dana.',
    wb_a_tunnels: 'Aktivni tuneli', wb_a_tunnels_d: 'Tuneli koji sada rade, s lokalnom adresom.', wb_tun_none: 'Nijedan tunel ne radi.', wb_tun_open: 'Tuneli…',
    wb_a_quick: 'Brzo spajanje', wb_a_quick_d: 'Gumbi za poslužitelje koje najčešće otvaraš.', wb_quick_choose: 'Odaberi poslužitelje u ⚙ postavkama.',
    wb_a_notes: 'Bilješke', wb_a_notes_d: 'Kratka bilješka za tebe, sprema se na poslužitelj.', wb_notes_ph: 'Napiši bilješku…',
    wb_unavailable: 'Nije dostupno', wb_add_btn: 'Dodaj',
  });

  const $ = id => document.getElementById(id);
  const tn = (k, v) => { let s = t(k); for (const n in (v || {})) s = s.split('{' + n + '}').join(v[n]); return s; };
  const MAX_TILES = 24, NOTE_MAX = 2000, METRIC_MAX = 8, QUICK_MAX = 12, STATUS_MAX = 60;
  const user = () => (typeof currentUser !== 'undefined' && currentUser) || null;
  const pol = () => (user() && user().policy) || {};
  const allConns = () => (typeof conns !== 'undefined' && Array.isArray(conns) ? conns : []);
  const connById = id => allConns().find(c => c.id === id);
  const sshLike = c => c && (c.protocol === 'SSH' || c.protocol === 'SFTP');
  const compact = () => typeof isCompact === 'function' && isCompact();
  const visible = el => !!el && el.isConnected && el.getClientRects().length > 0 && document.visibilityState === 'visible';
  const pct = v => (v == null || !isFinite(v) ? null : Math.max(0, Math.min(100, Number(v))));
  const level = p => (p == null ? '' : p >= 90 ? 'crit' : p >= 75 ? 'warn' : 'ok');

  function ago(ts) {
    const d = new Date(ts).getTime();
    if (!d) return '';
    const s = Math.round((d - Date.now()) / 1000), a = Math.abs(s);
    let rtf;
    try { rtf = new Intl.RelativeTimeFormat(typeof currentLang !== 'undefined' && currentLang === 'hr' ? 'hr' : 'en', {numeric: 'auto'}); } catch (e) { return new Date(ts).toLocaleString(); }
    if (a < 60) return rtf.format(s, 'second');
    if (a < 3600) return rtf.format(Math.round(s / 60), 'minute');
    if (a < 86400) return rtf.format(Math.round(s / 3600), 'hour');
    return rtf.format(Math.round(s / 86400), 'day');
  }
  function dur(sec) {
    sec = Number(sec) || 0;
    const d = Math.floor(sec / 86400), h = Math.floor(sec % 86400 / 3600), m = Math.floor(sec % 3600 / 60);
    return d ? `${d}d ${h}h` : h ? `${h}h ${m}m` : `${m}m`;
  }
  function openConn(c) {
    if (!c) return;
    close();
    if (typeof closeDrawers === 'function') closeDrawers();
    if (typeof isWebConn === 'function' && isWebConn(c)) { openWebConn(c); return; }
    const wo = addWindow();
    wo.connId = c.id;
    setWinMode(wo.id, defaultWinMode(c));
  }
  function bar(label, p, side) {
    const v = pct(p);
    return `<div class="wb-bar ${level(v)}"><span class="l">${esc(label)}</span><span class="tr"><span class="fi" style="width:${v == null ? 0 : v.toFixed(1)}%"></span></span><span class="v">${v == null ? '—' : Math.round(v) + '%'}${side ? ' <span class="sd">' + esc(side) + '</span>' : ''}</span></div>`;
  }

  // ── addons ──
  // Each addon: ico, name/desc keys, avail() (with a reason when not), render(tile, el, ctx),
  // every (seconds between loads while visible), load(tile, ctx) and settings fields.
  const A = {
    status: {
      ico: '🟢', key: 'wb_a_status', desc: 'wb_a_status_d', fields: ['conns', 'problems'], connMax: STATUS_MAX, connFilter: () => true,
      sig: tile => JSON.stringify([typeof statusInfo !== 'undefined' && statusInfo.enabled, typeof connStatus !== 'undefined' ? connStatus : null, tile.cfg]),
      render(tile, el) {
        if (typeof statusInfo === 'undefined' || !statusInfo.enabled) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_status_off'))}</div>`; return; }
        const ids = (tile.cfg.ids || []).filter(connById);
        let list = (ids.length ? ids.map(connById) : allConns().filter(c => connStatus[c.id])).map(c => ({c, st: connStatus[c.id] || null}));
        const rank = x => (!x.st ? 3 : x.st.state === 'down' ? 0 : x.st.state === 'up' && x.st.latency_ms > 300 ? 1 : x.st.state === 'up' ? 4 : 2);
        list.sort((a, b) => rank(a) - rank(b) || a.c.name.localeCompare(b.c.name));
        const total = list.length;
        if (tile.cfg.problems) list = list.filter(x => rank(x) < 4);
        if (!total) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_status_none'))}</div>`; return; }
        if (!list.length) { el.innerHTML = `<div class="wb-ok">✓ ${esc(tn('wb_status_all_ok', {n: total}))}</div>`; return; }
        const up = list.filter(x => rank(x) === 4).length, down = list.filter(x => rank(x) === 0).length;
        el.innerHTML = `<div class="wb-sum">${down ? `<span class="b down">${down} ${esc(t('st_down'))}</span>` : ''}<span class="b up">${up} ${esc(t('st_up'))}</span></div>` +
          `<div class="wb-list">${list.slice(0, STATUS_MAX).map(x => {
            const cls = x.st ? statusClass(x.st) : 'unknown';
            const side = x.st && x.st.state === 'up' ? Number(x.st.latency_ms) + ' ms' : x.st && x.st.since ? ago(x.st.since) : '';
            return `<button type="button" class="wb-row" data-c="${x.c.id}" title="${esc(statusTitle(x.st))}"><span class="conn-st ${cls}"></span><span class="nm">${esc(x.c.name)}</span><span class="sd">${esc(side)}</span></button>`;
          }).join('')}</div>`;
        el.querySelectorAll('.wb-row').forEach(b => { b.onclick = () => openConn(connById(Number(b.dataset.c))); });
      },
    },
    metrics: {
      ico: '📈', key: 'wb_a_metrics', desc: 'wb_a_metrics_d', fields: ['conns'], connMax: METRIC_MAX, connFilter: sshLike, every: 60, refresh: true,
      async load(tile, ctx) {
        const ids = tile.cfg.ids || [];
        ctx.data = ctx.data || {};
        await Promise.all(ids.map(async id => {
          try {
            const r = await fetch('/api/workboard/metrics?id=' + id + (ctx.force ? '&refresh=1' : ''));
            const j = await r.json().catch(() => ({}));
            ctx.data[id] = r.ok ? j : {error: j.error || ('HTTP ' + r.status)};
          } catch (e) { ctx.data[id] = {error: String(e.message || e)}; }
        }));
      },
      render(tile, el, ctx) {
        const ids = (tile.cfg.ids || []).filter(id => sshLike(connById(id)));
        if (!ids.length) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_metrics_choose'))}</div><div class="wb-hint sm">${esc(t('wb_metrics_help'))}</div>`; return; }
        const data = ctx.data || {};
        el.innerHTML = `<div class="wb-srv-grid">${ids.map(id => {
          const c = connById(id), m = data[id];
          let body;
          if (!m) body = `<div class="wb-hint sm">${esc(t('wb_loading'))}</div>`;
          else if (m.error) body = `<div class="wb-err">⚠ ${esc(m.error)}</div>`;
          else {
            const disks = (m.disks || []).slice().sort((a, b) => b.percent - a.percent);
            const shown = disks.slice(0, tile.size === 'm' ? 6 : 3);
            const ld = m.load != null ? Number(m.load).toFixed(2) + (m.cpu_count ? ' / ' + m.cpu_count : '') : '';
            const loadPct = m.load != null && m.cpu_count ? (m.load / m.cpu_count) * 100 : null;
            body = bar(t('wb_m_cpu'), m.cpu_percent) +
              bar(t('wb_m_ram'), m.memory_percent, m.memory_total ? fmtBytes(m.memory_total) : '') +
              (m.load != null ? (loadPct != null ? bar(t('wb_m_load'), loadPct, ld) : `<div class="wb-kv"><span>${esc(t('wb_m_load'))}</span><b>${esc(ld)}</b></div>`) : '') +
              shown.map(d => bar(d.mount, d.percent, d.size_kb ? fmtBytes(d.size_kb * 1024) : '')).join('') +
              (disks.length > shown.length ? `<div class="wb-hint sm" title="${esc(disks.slice(shown.length).map(d => d.mount + ' ' + Math.round(d.percent) + '%').join('\n'))}">${esc(tn('wb_m_more_disks', {n: disks.length - shown.length}))}</div>` : '') +
              `<div class="wb-src">${esc(t(m.source === 'zabbix' ? 'wb_m_src_zabbix' : 'wb_m_src_ssh'))}${m.uptime ? ' · ' + esc(tn('wb_m_up', {d: dur(m.uptime)})) : ''}${m.at ? ' · ' + esc(ago(m.at)) : ''}</div>`;
          }
          const st = typeof connStatusOf === 'function' ? connStatusOf(c) : null;
          return `<div class="wb-srv"><button type="button" class="wb-srv-h" data-c="${id}"><span class="conn-st ${st ? statusClass(st) : 'unknown'}"></span><span class="nm">${esc(c.name)}</span></button>${body}</div>`;
        }).join('')}</div>`;
        el.querySelectorAll('.wb-srv-h').forEach(b => { b.onclick = () => openConn(connById(Number(b.dataset.c))); });
      },
    },
    events: {
      ico: '🔔', key: 'wb_a_events', desc: 'wb_a_events_d', fields: [], every: 300, refresh: true,
      async load(tile, ctx) {
        try {
          const r = await fetch('/api/me/events?limit=50');
          if (!r.ok) { ctx.error = 'HTTP ' + r.status; return; }
          const j = await r.json();
          ctx.events = j.events || []; ctx.unread = j.unread || 0; ctx.error = '';
          setUnread(ctx.unread);
        } catch (e) { ctx.error = String(e.message || e); }
      },
      render(tile, el, ctx) {
        if (ctx.error) { el.innerHTML = `<div class="wb-err">⚠ ${esc(ctx.error)}</div>`; return; }
        if (!ctx.events) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_loading'))}</div>`; return; }
        const ev = ctx.events;
        el.innerHTML = `<div class="wb-acts"><span class="wb-cnt">${ctx.unread ? esc(tn('wb_ev_unread', {n: ctx.unread})) : ''}</span>` +
          `<button type="button" class="btn-sec btn-sm" data-a="read" ${ctx.unread ? '' : 'disabled'}>${esc(t('wb_ev_read_all'))}</button><button type="button" class="btn-sec btn-sm" data-a="clear" ${ev.length ? '' : 'disabled'}>${esc(t('wb_ev_clear'))}</button></div>` +
          (ev.length ? `<div class="wb-list">${ev.map(e => `<div class="wb-ev${e.read ? '' : ' new'}" data-id="${e.id}"><div class="h"><span class="tt">${esc(e.title)}</span><span class="sd" title="${esc(fmtTime(e.at))}">${esc(ago(e.at))}</span></div>${e.body ? `<div class="bd">${esc(e.body)}</div>` : ''}</div>`).join('')}</div>`
            : `<div class="wb-hint">${esc(t('wb_ev_none'))}</div>`);
        el.querySelector('[data-a="read"]').onclick = () => markRead({all: true});
        el.querySelector('[data-a="clear"]').onclick = async () => {
          if (!await uiConfirm(t('wb_ev_clear_q'), {danger: true, okText: t('wb_ev_clear')})) return;
          await fetch('/api/me/events', {method: 'DELETE'}).catch(() => {});
          reloadType('events');
        };
        el.querySelectorAll('.wb-ev.new').forEach(d => { d.onclick = () => markRead({ids: [Number(d.dataset.id)]}); });
      },
    },
    clock: {
      ico: '🕒', key: 'wb_a_clock', desc: 'wb_a_clock_d', fields: ['zones', 'seconds'], tick: true,
      render(tile, el) {
        const sec = tile.cfg.seconds !== false, now = new Date(), loc = typeof currentLang !== 'undefined' && currentLang === 'hr' ? 'hr-HR' : undefined;
        const fmt = (tz, o) => { try { return new Intl.DateTimeFormat(loc, Object.assign({timeZone: tz || undefined}, o)).format(now); } catch (e) { return '—'; } };
        const tOpt = {hour: '2-digit', minute: '2-digit', second: sec ? '2-digit' : undefined, hour12: false};
        const zones = (tile.cfg.zones || []).filter(Boolean);
        const html = `<div class="wb-clock"><div class="big">${esc(fmt('', tOpt))}</div><div class="dt">${esc(fmt('', {weekday: 'long', day: 'numeric', month: 'long', year: 'numeric'}))}</div></div>` +
          (zones.length ? `<div class="wb-zones">${zones.map(z => `<div class="wb-kv"><span>${esc(z)}</span><b>${esc(fmt(z, {hour: '2-digit', minute: '2-digit', hour12: false, weekday: 'short'}))}</b></div>`).join('')}</div>` : '');
        if (el._last !== html) { el.innerHTML = html; el._last = html; }
      },
    },
    commits: {
      ico: '⎇', key: 'wb_a_commits', desc: 'wb_a_commits_d', fields: ['app', 'branch', 'count'], every: 300, refresh: true,
      avail: () => (pol().git === false ? t('wb_commits_nogit') : ''),
      async load(tile, ctx) {
        if (!tile.cfg.app) return;
        const q = new URLSearchParams({app: tile.cfg.app});
        if (tile.cfg.branch) q.set('branch', tile.cfg.branch);
        if (ctx.force) q.set('refresh', '1');
        try {
          const r = await fetch('/api/git/activity?' + q);
          const j = await r.json().catch(() => ({}));
          if (!r.ok) { ctx.error = j.error || ('HTTP ' + r.status); return; }
          ctx.items = (j.items || []).slice().sort((a, b) => String(b.at).localeCompare(String(a.at)));
          ctx.error = '';
        } catch (e) { ctx.error = String(e.message || e); }
      },
      render(tile, el, ctx) {
        if (!tile.cfg.app) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_commits_choose'))}</div>`; return; }
        const head = `<div class="wb-sub">⎇ ${esc(tile.cfg.app)}${tile.cfg.branch ? ' · ' + esc(tile.cfg.branch) : ''}</div>`;
        if (ctx.error) { el.innerHTML = head + `<div class="wb-err">⚠ ${esc(ctx.error)}</div>`; return; }
        if (!ctx.items) { el.innerHTML = head + `<div class="wb-hint">${esc(t('wb_loading'))}</div>`; return; }
        const items = ctx.items.slice(0, Number(tile.cfg.count) || 10);
        el.innerHTML = head + (items.length ? `<div class="wb-list">${items.map(i => {
          const title = i.kind === 'push' ? `${t('wb_commits_pushed')} ${i.ref || i.branch || ''}${i.count ? ' (' + i.count + ')' : ''}` : (i.title || '');
          const br = i.branch || (i.branches || []).join(', ') || (i.kind === 'push' ? '' : '');
          const sha = (i.sha || i.after || '').slice(0, 8);
          const inner = `<div class="h"><span class="tt">${i.kind === 'push' ? '⇪ ' : ''}${esc(title)}</span><span class="sd" title="${esc(fmtTime(i.at))}">${esc(ago(i.at))}</span></div><div class="mt">${esc(i.author || '')}${br ? ' · ' + esc(br) : ''}${sha ? ' · <code>' + esc(sha) + '</code>' : ''}</div>`;
          return /^https?:\/\//i.test(i.url || '') ? `<a class="wb-ev cm" href="${esc(i.url)}" target="_blank" rel="noopener noreferrer">${inner}</a>` : `<div class="wb-ev cm">${inner}</div>`;
        }).join('')}</div>` : `<div class="wb-hint">${esc(t('wb_commits_none'))}</div>`);
      },
    },
    tunnels: {
      ico: '🔀', key: 'wb_a_tunnels', desc: 'wb_a_tunnels_d', fields: [],
      sig: () => JSON.stringify(typeof tunnelState !== 'undefined' ? tunnelState.tunnels.filter(x => x.running).map(x => [x.id, x.listen, x.state, x.active]) : null),
      render(tile, el) {
        const run = (typeof tunnelState !== 'undefined' ? tunnelState.tunnels : []).filter(x => x.running);
        el.innerHTML = (run.length ? `<div class="wb-list">${run.map(v => {
          const c = connById(v.conn_id);
          return `<div class="wb-row ro"><span class="conn-st ${v.state === 'error' ? 'down' : 'up'}"></span><span class="nm">${esc(v.name || (c && c.name) || '')}</span><span class="sd"><code>${esc(v.listen || '')}</code></span></div>`;
        }).join('')}</div>` : `<div class="wb-hint">${esc(t('wb_tun_none'))}</div>`) +
          `<div class="wb-acts"><span class="wb-cnt"></span><button type="button" class="btn-sec btn-sm" data-a="tun">${esc(t('wb_tun_open'))}</button></div>`;
        el.querySelector('[data-a="tun"]').onclick = () => { close(); openTunnelsPanel(); };
      },
    },
    quick: {
      ico: '⚡', key: 'wb_a_quick', desc: 'wb_a_quick_d', fields: ['conns'], connMax: QUICK_MAX, connFilter: () => true,
      sig: tile => JSON.stringify([tile.cfg.ids, (tile.cfg.ids || []).map(id => { const c = connById(id); return c ? [c.name, connStatus[id] && connStatus[id].state] : null; })]),
      render(tile, el) {
        const list = (tile.cfg.ids || []).map(connById).filter(Boolean);
        if (!list.length) { el.innerHTML = `<div class="wb-hint">${esc(t('wb_quick_choose'))}</div>`; return; }
        el.innerHTML = `<div class="wb-quick">${list.map(c => { const st = connStatusOf(c); return `<button type="button" class="btn-sec" data-c="${c.id}" title="${esc(c.host || '')}"><span class="conn-st ${st ? statusClass(st) : 'unknown'}"></span>${esc(c.name)}</button>`; }).join('')}</div>`;
        el.querySelectorAll('[data-c]').forEach(b => { b.onclick = () => openConn(connById(Number(b.dataset.c))); });
      },
    },
    notes: {
      ico: '📝', key: 'wb_a_notes', desc: 'wb_a_notes_d', fields: [],
      render(tile, el) {
        if (el.querySelector('textarea')) return;
        el.innerHTML = `<textarea class="wb-note" maxlength="${NOTE_MAX}" placeholder="${esc(t('wb_notes_ph'))}"></textarea>`;
        const ta = el.querySelector('textarea');
        ta.value = tile.cfg.text || '';
        ta.oninput = () => { tile.cfg.text = ta.value.slice(0, NOTE_MAX); saveSoon(); };
      },
    },
  };
  const ORDER = ['status', 'metrics', 'events', 'clock', 'commits', 'tunnels', 'quick', 'notes'];

  // ── state and persistence ──
  let S = null, saveT = 0, open = false, ticker = 0, unread = 0, appsCache = null;
  const ctxs = {}, bodies = {}, panels = {}, registered = new Set();
  const lsKey = () => 'wrm_workboard_' + ((user() && user().id) || 0);
  const canSave = () => !!user() && !(typeof isCollabMode !== 'undefined' && isCollabMode);
  const newId = () => { let s = ''; const c = 'abcdefghijklmnopqrstuvwxyz0123456789'; for (let i = 0; i < 8; i++) s += c[Math.floor(Math.random() * c.length)]; return S && S.tiles.some(x => x.id === s) ? newId() : s; };
  const nameOf = tile => (tile.cfg.title || '').trim() || t(A[tile.type].key);

  function sanitize(x) {
    const out = {v: 1, t: x && Number(x.t) || 0, tiles: []};
    const seen = new Set();
    for (const tl of (x && Array.isArray(x.tiles) ? x.tiles : [])) {
      if (!tl || !A[tl.type] || typeof tl.id !== 'string' || !/^[a-z0-9]{1,16}$/.test(tl.id) || seen.has(tl.id)) continue;
      seen.add(tl.id);
      const c = tl.cfg && typeof tl.cfg === 'object' ? tl.cfg : {}, cfg = {};
      if (typeof c.title === 'string') cfg.title = c.title.slice(0, 60);
      if (Array.isArray(c.ids)) cfg.ids = c.ids.map(Number).filter(n => n > 0).slice(0, STATUS_MAX);
      if (c.problems) cfg.problems = true;
      if (Array.isArray(c.zones)) cfg.zones = c.zones.filter(z => typeof z === 'string').map(z => z.slice(0, 60)).slice(0, 8);
      if (c.seconds === false) cfg.seconds = false;
      if (typeof c.app === 'string') cfg.app = c.app.slice(0, 200);
      if (typeof c.branch === 'string') cfg.branch = c.branch.slice(0, 200);
      if ([5, 10, 20].includes(Number(c.count))) cfg.count = Number(c.count);
      if (typeof c.text === 'string') cfg.text = c.text.slice(0, NOTE_MAX);
      out.tiles.push({id: tl.id, type: tl.type, size: tl.size === 'm' ? 'm' : 's', cfg});
      if (out.tiles.length >= MAX_TILES) break;
    }
    return out;
  }
  function saveSoon() {
    if (!S || !canSave()) return;
    S.t = Date.now();
    try { localStorage.setItem(lsKey(), JSON.stringify(S)); } catch (e) { /* only on the server then */ }
    clearTimeout(saveT);
    saveT = setTimeout(push, 800);
  }
  function push(keepalive) {
    saveT = 0;
    if (!S || !canSave()) return;
    fetch('/api/me/ui/workboard', {method: 'PUT', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({value: S}), keepalive: !!keepalive}).catch(() => {});
  }
  window.addEventListener('pagehide', () => { if (saveT) { clearTimeout(saveT); push(true); } });

  async function init() {
    if (!canSave()) { const b = $('btn-workboard'); if (b) b.style.display = 'none'; return; }
    let local = null;
    try { local = JSON.parse(localStorage.getItem(lsKey()) || 'null'); } catch (e) { /* none */ }
    S = sanitize(local || {});
    reconcile();
    if (!ticker) ticker = setInterval(tick, 1000);
    try {
      const r = await fetch('/api/me/ui/workboard');
      if (r.ok) {
        const j = await r.json(), v = j && j.value, vt = (v && Number(v.t)) || 0;
        if (v && vt > (S.t || 0)) {
          S = sanitize(v);
          try { localStorage.setItem(lsKey(), JSON.stringify(S)); } catch (e) { /* not cached */ }
          reconcile();
        } else if (local && (S.t || 0) > vt) push();
      }
    } catch (e) { /* offline: the copy in this browser */ }
    if (S.tiles.some(x => x.type === 'events')) loadCtx(S.tiles.find(x => x.type === 'events'));
  }
  document.addEventListener('wrm-layout-loaded', () => { if (S) reconcile(); });

  // ── tile bodies, panels and the grid ──
  function bodyOf(tile) {
    let b = bodies[tile.id];
    if (!b) { b = document.createElement('div'); b.className = 'wb-body'; bodies[tile.id] = b; }
    return b;
  }
  const ctxOf = tile => ctxs[tile.id] || (ctxs[tile.id] = {});
  const pid = tile => 'wb:' + tile.id;
  const isOut = tile => !compact() && !!window.wrmLayout && registered.has(tile.id) && wrmLayout.where(pid(tile)) !== 'hidden';

  function panelOf(tile) {
    let p = panels[tile.id];
    if (!p) {
      p = document.createElement('div');
      p.id = 'pnl-wb-' + tile.id; p.className = 'wb-pnl';
      p.innerHTML = `<div class="wb-pbar"><span class="wb-upd"></span><button type="button" class="wb-ib" data-a="refresh" title="${esc(t('wb_t_refresh'))}">↻</button><button type="button" class="wb-ib" data-a="set" title="${esc(t('wb_t_settings'))}">⚙</button></div>`;
      p.querySelector('[data-a="set"]').onclick = () => { const tl = tileById(tile.id); if (tl) settings(tl); };
      p.querySelector('[data-a="refresh"]').onclick = () => { const tl = tileById(tile.id); if (tl) loadCtx(tl, true); };
      panels[tile.id] = p;
    }
    const def = A[tile.type];
    p.querySelector('[data-a="refresh"]').style.display = def.refresh ? '' : 'none';
    return p;
  }
  function registerPanel(tile, place) {
    const p = panelOf(tile), b = bodyOf(tile);
    if (b.parentElement !== p) p.appendChild(b);
    const ok = wrmLayout.register(pid(tile), {
      ico: A[tile.type].ico, el: p.id, node: p,
      label: () => { const tl = tileById(tile.id); return tl ? nameOf(tl) : ''; },
      shown: () => { const tl = tileById(tile.id); if (tl) paint(tl, true); },
      closed: () => { registered.delete(tile.id); renderGrid(); },
    }, place);
    if (ok) registered.add(tile.id);
  }
  // reconcile matches the tiles with the panels the layout places: a placed tile is a panel,
  // a placed id without a tile is dropped, a tile the layout no longer places is on the grid.
  function reconcile() {
    if (window.wrmLayout && wrmLayout.ready()) {
      const placed = new Set(wrmLayout.placed());
      for (const id of placed) if (!S.tiles.some(x => pid(x) === id)) wrmLayout.remove(id);
      for (const tile of S.tiles) {
        if (placed.has(pid(tile)) && !compact()) { if (!registered.has(tile.id)) registerPanel(tile); }
        else if (registered.has(tile.id)) { registered.delete(tile.id); wrmLayout.remove(pid(tile)); }
      }
    }
    for (const id of Object.keys(bodies)) if (!tileById(id)) { bodies[id].remove(); delete bodies[id]; delete ctxs[id]; if (panels[id]) { panels[id].remove(); delete panels[id]; } }
    renderGrid();
    updateBadge();
  }
  const tileById = id => S && S.tiles.find(x => x.id === id);

  function overlay() {
    let o = $('workboard');
    if (o) return o;
    o = document.createElement('div');
    o.id = 'workboard'; o.setAttribute('role', 'region'); o.setAttribute('aria-label', t('wb_title'));
    o.innerHTML = `<div class="wb-head"><span class="wb-ttl">▤ <span data-i18n="wb_title">${esc(t('wb_title'))}</span></span><span class="wb-hsp"></span>` +
      `<button type="button" class="wb-addbtn" data-a="add">＋ <span>${esc(t('wb_add'))}</span></button><button type="button" class="wb-ib wb-x" data-a="close" title="${esc(t('wb_close'))}" aria-label="${esc(t('wb_close'))}">✕</button></div><div class="wb-grid" id="wb-grid"></div>`;
    o.querySelector('[data-a="add"]').onclick = catalog;
    o.querySelector('[data-a="close"]').onclick = close;
    o.addEventListener('keydown', e => { if (e.key === 'Escape' && !document.querySelector('.modal-overlay.open')) { e.stopPropagation(); close(); } });
    const wc = $('win-container');
    (wc || document.body).appendChild(o);
    return o;
  }
  function renderGrid() {
    if (!S) return;
    const o = overlay(), g = $('wb-grid');
    o.querySelector('.wb-addbtn span').textContent = t('wb_add');
    o.querySelector('.wb-ttl [data-i18n]').textContent = t('wb_title');
    o.querySelector('.wb-x').title = t('wb_close');
    if (!S.tiles.length) {
      g.innerHTML = `<div class="wb-empty"><div class="ic">▤</div><div class="tt">${esc(t('wb_empty_title'))}</div><div class="tx">${esc(t('wb_empty'))}</div><button type="button" data-a="add">＋ ${esc(t('wb_add'))}</button></div>`;
      g.querySelector('[data-a="add"]').onclick = catalog;
      return;
    }
    const locked = window.wrmLayout && wrmLayout.locked();
    // Keep the bodies (their content and timers) while the cards are rebuilt.
    S.tiles.forEach(tl => { const b = bodies[tl.id]; if (b && b.parentElement && b.parentElement.closest('#wb-grid')) b.remove(); });
    g.innerHTML = '';
    S.tiles.forEach(tile => {
      const def = A[tile.type], out = isOut(tile), card = document.createElement('div');
      card.className = 'wb-tile' + (tile.size === 'm' ? ' m' : '') + (out ? ' out' : '');
      card.dataset.id = tile.id;
      card.innerHTML = `<div class="wb-th" ${locked || compact() ? '' : 'draggable="true"'} title="${locked ? '' : esc(t('wb_t_drag'))}"><span class="ico">${esc(def.ico)}</span><span class="nm">${esc(nameOf(tile))}</span><span class="wb-upd"></span>` +
        (def.refresh && !out ? `<button type="button" class="wb-ib" data-a="refresh" title="${esc(t('wb_t_refresh'))}">↻</button>` : '') +
        (compact() ? '' : `<button type="button" class="wb-ib" data-a="size" title="${esc(t('wb_t_size'))}">${tile.size === 'm' ? '⇤' : '⇥'}</button><button type="button" class="wb-ib" data-a="out" title="${esc(t('wb_t_out'))}" ${out ? 'disabled' : ''}>⤢</button>`) +
        `<button type="button" class="wb-ib" data-a="set" title="${esc(t('wb_t_settings'))}">⚙</button><button type="button" class="wb-ib" data-a="rm" title="${esc(t('wb_t_remove'))}">✕</button></div>`;
      if (out) {
        const w = wrmLayout.where(pid(tile));
        const ph = document.createElement('div');
        ph.className = 'wb-body wb-outph';
        ph.innerHTML = `<div class="wb-hint">⤢ ${esc(tn('wb_out_here', {where: t('dk_where_' + w)}))}</div><button type="button" class="btn-sec btn-sm">${esc(t('wb_bring_back'))}</button>`;
        ph.querySelector('button').onclick = () => wrmLayout.hide(pid(tile));
        card.appendChild(ph);
      } else card.appendChild(bodyOf(tile));
      card.querySelectorAll('.wb-ib').forEach(b => { b.onclick = e => { e.stopPropagation(); tileAction(tile, b.dataset.a); }; });
      const th = card.querySelector('.wb-th');
      th.addEventListener('dragstart', e => { e.dataTransfer.setData('text/x-wrm-tile', tile.id); e.dataTransfer.effectAllowed = 'move'; card.classList.add('drag'); });
      th.addEventListener('dragend', () => { card.classList.remove('drag'); g.querySelectorAll('.wb-tile.over').forEach(x => x.classList.remove('over')); });
      card.addEventListener('dragover', e => { if (![...e.dataTransfer.types].includes('text/x-wrm-tile')) return; e.preventDefault(); card.classList.add('over'); });
      card.addEventListener('dragleave', () => card.classList.remove('over'));
      card.addEventListener('drop', e => {
        const from = e.dataTransfer.getData('text/x-wrm-tile');
        card.classList.remove('over');
        if (!from || from === tile.id) return;
        e.preventDefault();
        const a = S.tiles.findIndex(x => x.id === from);
        if (a < 0) return;
        const [mv] = S.tiles.splice(a, 1);
        S.tiles.splice(S.tiles.findIndex(x => x.id === tile.id), 0, mv);
        saveSoon(); renderGrid();
      });
      g.appendChild(card);
      if (!out) paint(tile);
    });
  }
  async function tileAction(tile, a) {
    const locked = window.wrmLayout && wrmLayout.locked();
    if (a === 'set') settings(tile);
    else if (a === 'refresh') loadCtx(tile, true);
    else if (a === 'size') { if (locked) { showToast(t('wb_locked'), 'warning', 3500); return; } tile.size = tile.size === 'm' ? 's' : 'm'; saveSoon(); renderGrid(); }
    else if (a === 'out') {
      if (locked) { showToast(t('wb_locked'), 'warning', 3500); return; }
      if (!window.wrmLayout || compact()) return;
      registerPanel(tile, 'float');
      renderGrid();
    } else if (a === 'rm') {
      if (!await uiConfirm(tn('wb_remove_q', {name: nameOf(tile)}), {danger: true, okText: t('wb_t_remove')})) return;
      if (registered.has(tile.id)) { registered.delete(tile.id); wrmLayout.remove(pid(tile)); }
      S.tiles = S.tiles.filter(x => x !== tile);
      saveSoon(); reconcile();
    }
  }

  // ── painting and loading ──
  function paint(tile, force) {
    const def = A[tile.type], el = bodyOf(tile), ctx = ctxOf(tile);
    const why = def.avail ? def.avail() : '';
    if (why) { el.innerHTML = `<div class="wb-hint">${esc(why)}</div>`; el._last = ''; ctx.sig = ''; return; }
    if (def.sig) { const s = def.sig(tile); if (!force && s === ctx.sig && el.childElementCount) return; ctx.sig = s; }
    try { def.render(tile, el, ctx); } catch (e) { el.innerHTML = `<div class="wb-err">⚠ ${esc(String(e.message || e))}</div>`; }
    const upd = ctx.at ? tn('wb_updated', {when: ago(ctx.at)}) : '';
    [el.closest('.wb-tile'), el.closest('.wb-pnl')].forEach(h => { const u = h && h.querySelector('.wb-upd'); if (u) u.textContent = def.load ? upd : ''; });
    if (def.load && !ctx.at && !ctx.busy) loadCtx(tile);
  }
  async function loadCtx(tile, force) {
    const def = A[tile.type], ctx = ctxOf(tile);
    if (!def.load || ctx.busy || (def.avail && def.avail())) return;
    // tiles that name connections wait until the connection list has loaded
    if ((def.fields || []).includes('conns') && !allConns().length) return;
    ctx.busy = true; ctx.force = !!force;
    try { await def.load(tile, ctx); } finally { ctx.busy = false; ctx.force = false; ctx.at = new Date().toISOString(); }
    if (tileById(tile.id)) paint(tile, true);
  }
  function reloadType(type) { if (S) S.tiles.filter(x => x.type === type).forEach(x => loadCtx(x, true)); }
  async function markRead(body) {
    await fetch('/api/me/events/read', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)}).catch(() => {});
    reloadType('events');
  }
  function setUnread(n) { unread = n; updateBadge(); }
  function updateBadge() {
    const b = $('wb-badge');
    if (!b) return;
    const on = !!(S && S.tiles.some(x => x.type === 'events')) && unread > 0;
    b.style.display = on ? '' : 'none';
    b.textContent = on ? (unread > 99 ? '99+' : String(unread)) : '';
  }
  let secs = 0;
  function tick() {
    if (!S || !S.tiles.length) return;
    secs++;
    for (const tile of S.tiles) {
      const def = A[tile.type], el = bodies[tile.id];
      if (!el || !visible(el)) continue;
      const ctx = ctxOf(tile);
      if (def.tick) paint(tile);
      else if (def.sig && secs % 2 === 0) paint(tile);
      if (def.load && def.every && !ctx.busy && (!ctx.at || Date.now() - new Date(ctx.at).getTime() > def.every * 1000)) loadCtx(tile);
      else if (def.load && secs % 30 === 0) paint(tile, true);
    }
  }

  // ── open / close ──
  function show() {
    if (!S) return;
    open = true;
    const o = overlay();
    o.classList.add('open');
    renderGrid();
    const b = $('btn-workboard'); if (b) b.classList.add('wb-on');
    if (typeof closeDrawers === 'function') closeDrawers();
    setTimeout(() => { const f = o.querySelector('.wb-addbtn'); if (f) f.focus(); }, 30);
  }
  function close() {
    if (!open) return;
    open = false;
    const o = $('workboard'); if (o) o.classList.remove('open');
    const b = $('btn-workboard'); if (b) b.classList.remove('wb-on');
  }
  function toggle() { if (open) close(); else show(); }

  // ── catalog and settings dialogs ──
  function modal(id, html) {
    let m = $(id);
    if (!m) { m = document.createElement('div'); m.className = 'modal-overlay'; m.id = id; document.body.appendChild(m); m.addEventListener('mousedown', e => { if (e.target === m) closeModal(id); }); }
    m.innerHTML = `<div class="modal wb-modal">${html}</div>`;
    m.classList.add('open');
    return m;
  }
  function catalog() {
    const m = modal('modal-wb-catalog', `<h3>▤ ${esc(t('wb_catalog'))}</h3><div class="muted-text" style="line-height:1.5;margin-bottom:10px;">${esc(t('wb_catalog_hint'))}</div><div class="wb-cat">${ORDER.map(k => {
      const d = A[k], why = d.avail ? d.avail() : '', n = S.tiles.filter(x => x.type === k).length;
      return `<div class="wb-ci"><span class="ico">${esc(d.ico)}</span><div class="tx"><div class="nm">${esc(t(d.key))}${n ? ` <span class="wb-n">×${n}</span>` : ''}</div><div class="ds">${esc(why || t(d.desc))}</div></div><button type="button" class="btn-sm" data-k="${k}" ${why ? 'disabled' : ''}>${esc(why ? t('wb_unavailable') : t('wb_add_btn'))}</button></div>`;
    }).join('')}</div><div class="modal-footer"><button class="btn-sec" data-a="close">${esc(t('close'))}</button></div>`);
    m.querySelector('[data-a="close"]').onclick = () => closeModal('modal-wb-catalog');
    m.querySelectorAll('[data-k]').forEach(b => { b.onclick = () => add(b.dataset.k); });
  }
  function add(type) {
    if (S.tiles.length >= MAX_TILES) { showToast(tn('wb_max', {n: MAX_TILES}), 'warning', 3500); return; }
    const tile = {id: newId(), type, size: type === 'metrics' || type === 'commits' ? 'm' : 's', cfg: {}};
    if (type === 'clock') tile.cfg.zones = ['UTC'];
    if (type === 'commits') tile.cfg.count = 10;
    S.tiles.push(tile);
    saveSoon();
    closeModal('modal-wb-catalog');
    if (!open) show(); else renderGrid();
    showToast(tn('wb_added', {name: t(A[type].key)}), 'success', 2200);
    if ((A[type].fields || []).some(f => f === 'conns' || f === 'app') && type !== 'status') settings(tile);
    const card = document.querySelector(`.wb-tile[data-id="${tile.id}"]`);
    if (card) card.scrollIntoView({block: 'nearest'});
  }
  async function settings(tile) {
    const def = A[tile.type], f = def.fields || [], cfg = tile.cfg;
    let html = `<h3>${esc(def.ico)} ${esc(t(def.key))}</h3>` +
      `<div class="field"><label>${esc(t('wb_set_title'))}</label><input type="text" id="wbs-title" maxlength="60" placeholder="${esc(t('wb_set_title_ph'))}" value="${esc(cfg.title || '')}"></div>`;
    if (!compact()) html += `<div class="field"><label>${esc(t('wb_set_size'))}</label><select id="wbs-size"><option value="s">${esc(t('wb_size_s'))}</option><option value="m">${esc(t('wb_size_m'))}</option></select></div>`;
    if (f.includes('conns')) {
      const list = allConns().filter(def.connFilter).sort((a, b) => a.name.localeCompare(b.name));
      const sel = new Set(cfg.ids || []);
      html += `<div class="field"><label>${esc(t('wb_set_conns'))}</label><div class="muted-text" style="margin-bottom:6px;">${esc(tile.type === 'status' ? t('wb_set_conns_all') : tn('wb_set_conns_need', {n: def.connMax}))}</div>` +
        `<input type="text" id="wbs-filter" placeholder="${esc(t('wb_set_filter'))}"><div class="wb-pick" id="wbs-conns">${list.map(c => `<label class="wb-pi" data-q="${esc((c.name + ' ' + (c.host || '') + ' ' + (c.folder_path || '')).toLowerCase())}"><input type="checkbox" value="${c.id}" ${sel.has(c.id) ? 'checked' : ''}><span class="nm">${esc(c.name)}</span><span class="sd">${esc(c.protocol)} · ${esc(c.host || '')}</span></label>`).join('') || `<div class="wb-hint sm">${esc(t('wb_none'))}</div>`}</div></div>`;
      if (tile.type === 'metrics') html += `<div class="muted-text" style="line-height:1.5;">${esc(t('wb_metrics_help'))}</div>`;
    }
    if (f.includes('problems')) html += `<label class="wb-chk"><input type="checkbox" id="wbs-problems" ${cfg.problems ? 'checked' : ''}> ${esc(t('wb_status_problems'))}</label>`;
    if (f.includes('zones')) html += `<div class="field"><label>${esc(t('wb_clock_zones'))}</label><input type="text" id="wbs-zones" placeholder="${esc(t('wb_clock_zones_ph'))}" value="${esc((cfg.zones || []).join(', '))}"></div>`;
    if (f.includes('seconds')) html += `<label class="wb-chk"><input type="checkbox" id="wbs-seconds" ${cfg.seconds !== false ? 'checked' : ''}> ${esc(t('wb_clock_seconds'))}</label>`;
    if (f.includes('app')) html += `<div class="field"><label>${esc(t('wb_commits_service'))}</label><select id="wbs-app"><option value="">${esc(t('wb_loading'))}</option></select></div>` +
      `<div class="field"><label>${esc(t('wb_commits_branch'))}</label><input type="text" id="wbs-branch" maxlength="200" placeholder="${esc(t('wb_commits_branch_ph'))}" value="${esc(cfg.branch || '')}"></div>` +
      `<div class="field"><label>${esc(t('wb_commits_count'))}</label><select id="wbs-count">${[5, 10, 20].map(n => `<option value="${n}" ${Number(cfg.count || 10) === n ? 'selected' : ''}>${n}</option>`).join('')}</select></div>`;
    html += `<div class="modal-footer"><button class="btn-sec" data-a="cancel">${esc(t('cancel'))}</button><button data-a="save">${esc(t('save'))}</button></div>`;
    const m = modal('modal-wb-set', html);
    const q = s => m.querySelector(s);
    if (q('#wbs-size')) q('#wbs-size').value = tile.size;
    if (q('#wbs-filter')) q('#wbs-filter').oninput = e => { const v = e.target.value.trim().toLowerCase(); m.querySelectorAll('.wb-pi').forEach(l => { l.style.display = !v || l.dataset.q.includes(v) ? '' : 'none'; }); };
    q('[data-a="cancel"]').onclick = () => closeModal('modal-wb-set');
    setTimeout(() => q('#wbs-title').focus(), 30);
    if (f.includes('app')) {
      const sel = q('#wbs-app');
      try {
        if (!appsCache) { const r = await fetch('/api/git'); const j = r.ok ? await r.json() : {}; appsCache = ((j.catalog && j.catalog.apps) || []).map(a => a.name).filter(Boolean).sort(); setTimeout(() => { appsCache = null; }, 60000); }
      } catch (e) { appsCache = []; }
      const apps = appsCache.slice();
      if (cfg.app && !apps.includes(cfg.app)) apps.unshift(cfg.app);
      sel.innerHTML = apps.length ? `<option value="">—</option>` + apps.map(a => `<option value="${esc(a)}">${esc(a)}</option>`).join('') : `<option value="">${esc(t('wb_commits_noapps'))}</option>`;
      sel.value = cfg.app || '';
    }
    q('[data-a="save"]').onclick = () => {
      const n = {};
      const title = q('#wbs-title').value.trim();
      if (title) n.title = title.slice(0, 60);
      if (f.includes('conns')) {
        const ids = [...m.querySelectorAll('#wbs-conns input:checked')].map(i => Number(i.value));
        if (ids.length > def.connMax) { showToast(tn('wb_set_more', {n: def.connMax}), 'warning', 3500); return; }
        if (ids.length) n.ids = ids;
      }
      if (f.includes('problems') && q('#wbs-problems').checked) n.problems = true;
      if (f.includes('zones')) {
        const zs = q('#wbs-zones').value.split(/[,;\s]+/).map(z => z.trim()).filter(Boolean).slice(0, 8);
        for (const z of zs) { try { new Intl.DateTimeFormat('en', {timeZone: z}); } catch (e) { showToast(tn('wb_clock_zones_bad', {z}), 'error', 4000); return; } }
        n.zones = zs;
      }
      if (f.includes('seconds') && !q('#wbs-seconds').checked) n.seconds = false;
      if (f.includes('app')) { if (q('#wbs-app').value) n.app = q('#wbs-app').value; const br = q('#wbs-branch').value.trim(); if (br) n.branch = br; n.count = Number(q('#wbs-count').value) || 10; }
      if (cfg.text) n.text = cfg.text;
      const changedSrc = JSON.stringify([cfg.ids, cfg.app, cfg.branch]) !== JSON.stringify([n.ids, n.app, n.branch]);
      tile.cfg = n;
      if (q('#wbs-size')) tile.size = q('#wbs-size').value === 'm' ? 'm' : 's';
      if (changedSrc) delete ctxs[tile.id];
      const b = bodies[tile.id]; if (b) { b.innerHTML = ''; b._last = ''; }
      saveSoon();
      closeModal('modal-wb-set');
      renderGrid();
      if (registered.has(tile.id)) { wrmLayout.refresh(); paint(tile, true); }
      showToast(t('wb_saved'), 'success', 1600);
    };
  }

  // ── styles ──
  const css = document.createElement('style');
  css.textContent = `
#workboard { position: absolute; inset: 0; z-index: 7900; display: none; flex-direction: column; background: var(--bg); overflow: hidden; }
#workboard.open { display: flex; }
.wb-head { display: flex; align-items: center; gap: 8px; padding: 10px 14px; border-bottom: 1px solid var(--border); flex: none; }
.wb-ttl { font-weight: 700; font-size: 14px; letter-spacing: .02em; }
.wb-hsp { flex: 1; }
.wb-addbtn { padding: 6px 12px; font-size: 12px; }
.wb-grid { flex: 1; overflow: auto; padding: 14px; display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); grid-auto-rows: min-content; gap: 12px; align-content: start; align-items: start; }
.wb-tile { display: flex; flex-direction: column; background: var(--bg2); border: 1px solid var(--border2); border-radius: 12px; min-height: 120px; max-height: 420px; overflow: hidden; box-shadow: var(--shadow-md); }
.wb-tile.m { grid-column: span 2; }
.wb-tile.drag { opacity: .45; }
.wb-tile.over { outline: 2px solid var(--accent); outline-offset: 2px; }
.wb-th { display: flex; align-items: center; gap: 6px; padding: 7px 8px 7px 12px; border-bottom: 1px solid var(--border); flex: none; }
.wb-th[draggable="true"] { cursor: grab; }
.wb-th .ico { font-size: 13px; }
.wb-th .nm { font-weight: 700; font-size: 12px; flex: 1 1 auto; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.wb-upd { flex: 0 1000 auto; color: var(--text3); font-size: 10px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; min-width: 0; }
.wb-ib { background: transparent; border: 1px solid transparent; color: var(--text3); width: 24px; height: 24px; padding: 0; border-radius: 6px; font-size: 12px; line-height: 1; flex: none; cursor: pointer; }
.wb-ib:hover:not(:disabled) { background: var(--bg3); color: var(--text); border-color: var(--border2); }
.wb-ib:disabled { opacity: .35; cursor: default; }
.wb-body { flex: 1; min-height: 0; overflow: auto; padding: 10px 12px; font-size: 12px; }
.wb-outph { display: flex; flex-direction: column; gap: 8px; align-items: flex-start; }
.wb-pnl { display: flex; flex-direction: column; min-height: 0; }
.wb-pnl .wb-body { padding: 8px 10px; }
.wb-pbar { display: flex; align-items: center; gap: 4px; justify-content: flex-end; padding: 3px 6px 0; flex: none; }
.wb-pbar .wb-upd { flex: 1; text-align: left; padding-left: 4px; }
.wb-empty { grid-column: 1 / -1; max-width: 560px; margin: 40px auto; text-align: center; display: flex; flex-direction: column; gap: 10px; align-items: center; color: var(--text2); line-height: 1.55; }
.wb-empty .ic { font-size: 34px; opacity: .6; }
.wb-empty .tt { font-weight: 700; font-size: 15px; color: var(--text); }
.wb-hint { color: var(--text3); line-height: 1.5; }
.wb-hint.sm { font-size: 11px; margin-top: 4px; }
.wb-err { color: var(--red, #ef4444); line-height: 1.45; word-break: break-word; }
.wb-ok { color: var(--green, #22c55e); font-weight: 600; }
.wb-sum { display: flex; gap: 6px; margin-bottom: 8px; flex-wrap: wrap; }
.wb-sum .b { font-size: 11px; font-weight: 700; padding: 2px 8px; border-radius: 10px; background: var(--bg3); }
.wb-sum .b.down { color: var(--red, #ef4444); }
.wb-sum .b.up { color: var(--green, #22c55e); }
.wb-list { display: flex; flex-direction: column; gap: 2px; }
.wb-row { display: flex; align-items: center; gap: 8px; width: 100%; background: transparent; border: none; color: var(--text); padding: 5px 6px; border-radius: 6px; text-align: left; font-size: 12px; cursor: pointer; }
.wb-row:hover { background: var(--bg3); }
.wb-row.ro { cursor: default; }
.wb-row.ro:hover { background: transparent; }
.wb-row .nm, .wb-srv-h .nm { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.wb-row .sd, .wb-ev .sd { color: var(--text3); font-size: 11px; white-space: nowrap; }
.wb-row .conn-st, .wb-srv-h .conn-st, .wb-quick .conn-st { flex: none; position: static; }
.wb-srv-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 10px; }
.wb-srv { background: var(--bg3); border: 1px solid var(--border); border-radius: 9px; padding: 8px 10px; display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.wb-srv-h { display: flex; align-items: center; gap: 7px; background: transparent; border: none; color: var(--text); padding: 0; font-weight: 700; font-size: 12px; cursor: pointer; text-align: left; }
.wb-srv-h:hover .nm { text-decoration: underline; }
.wb-bar { display: grid; grid-template-columns: minmax(42px, 30%) 1fr auto; align-items: center; gap: 7px; font-size: 11px; }
.wb-bar .l { color: var(--text2); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-family: 'JetBrains Mono', monospace; }
.wb-bar .tr { height: 7px; background: var(--bg); border-radius: 4px; overflow: hidden; }
.wb-bar .fi { display: block; height: 100%; background: var(--green, #22c55e); border-radius: 4px; transition: width .4s; }
.wb-bar.warn .fi { background: var(--yellow, #f59e0b); }
.wb-bar.crit .fi { background: var(--red, #ef4444); }
.wb-bar .v { font-weight: 700; white-space: nowrap; text-align: right; }
.wb-bar .sd { color: var(--text3); font-weight: 500; }
.wb-src { color: var(--text3); font-size: 10px; margin-top: 2px; }
.wb-kv { display: flex; justify-content: space-between; gap: 8px; font-size: 12px; padding: 2px 0; }
.wb-kv span { color: var(--text2); }
.wb-clock { text-align: center; padding: 6px 0 8px; }
.wb-clock .big { font: 700 34px 'JetBrains Mono', monospace; letter-spacing: .02em; }
.wb-clock .dt { color: var(--text2); font-size: 12px; margin-top: 2px; }
.wb-zones { border-top: 1px solid var(--border); padding-top: 6px; }
.wb-acts { display: flex; align-items: center; gap: 6px; margin-bottom: 8px; }
.wb-acts .wb-cnt { flex: 1; color: var(--text2); font-size: 11px; font-weight: 600; }
.wb-ev { display: block; padding: 6px 8px; border-radius: 7px; border-left: 3px solid transparent; color: var(--text); text-decoration: none; }
.wb-ev + .wb-ev { margin-top: 2px; }
.wb-ev.new { background: var(--bg3); border-left-color: var(--accent); cursor: pointer; }
.wb-ev.cm:hover { background: var(--bg3); }
.wb-ev .h { display: flex; gap: 8px; align-items: baseline; }
.wb-ev .tt { flex: 1; min-width: 0; font-weight: 600; overflow-wrap: anywhere; }
.wb-ev .bd { color: var(--text2); font-size: 11px; margin-top: 2px; white-space: pre-wrap; overflow-wrap: anywhere; display: -webkit-box; -webkit-line-clamp: 4; -webkit-box-orient: vertical; overflow: hidden; }
.wb-ev .mt { color: var(--text3); font-size: 11px; margin-top: 2px; }
.wb-sub { color: var(--text2); font-size: 11px; font-weight: 600; margin-bottom: 6px; }
.wb-quick { display: flex; flex-wrap: wrap; gap: 6px; }
.wb-quick button { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; padding: 6px 10px; }
.wb-note { width: 100%; height: 100%; min-height: 110px; resize: none; background: var(--bg); color: var(--text); border: 1px solid var(--border); border-radius: 7px; padding: 8px; font-family: inherit; font-size: 12px; line-height: 1.5; box-sizing: border-box; }
.wb-modal { width: min(560px, calc(100vw - 24px)); max-height: calc(100vh - 32px); overflow: auto; }
.wb-cat { display: flex; flex-direction: column; gap: 6px; }
.wb-ci { display: flex; align-items: center; gap: 12px; padding: 9px 10px; background: var(--bg3); border: 1px solid var(--border); border-radius: 9px; }
.wb-ci .ico { font-size: 20px; width: 28px; text-align: center; flex: none; }
.wb-ci .tx { flex: 1; min-width: 0; }
.wb-ci .nm { font-weight: 700; font-size: 13px; }
.wb-ci .ds { color: var(--text2); font-size: 11.5px; line-height: 1.45; margin-top: 2px; }
.wb-n { color: var(--text3); font-weight: 600; font-size: 11px; }
.wb-pick { max-height: 220px; overflow: auto; border: 1px solid var(--border); border-radius: 7px; margin-top: 6px; padding: 4px; }
.field .wb-pi, .wb-pi { display: flex; align-items: center; gap: 8px; padding: 4px 6px; border-radius: 5px; cursor: pointer; font-size: 12px; text-transform: none; letter-spacing: 0; font-weight: 500; margin: 0; }
.wb-pi:hover { background: var(--bg3); }
.wb-pi input { width: auto; margin: 0; }
.wb-pi .nm { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: var(--text); }
.wb-pi .sd { color: var(--text3); font-size: 11px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 45%; }
.wb-chk { display: flex; align-items: center; gap: 8px; font-size: 12px; margin: 6px 0 10px; cursor: pointer; }
.wb-chk input { width: auto; margin: 0; }
#btn-workboard.wb-on { color: var(--accent); }
@media (max-width: 759px) {
  .wb-grid { grid-template-columns: 1fr; padding: 10px; }
  .wb-tile.m { grid-column: auto; }
  .wb-head { padding: 8px 10px; }
  .wb-addbtn span { display: none; }
}
`;
  document.head.appendChild(css);

  window.addEventListener('resize', (() => { let rt = 0, was = compact(); return () => { clearTimeout(rt); rt = setTimeout(() => { if (S && compact() !== was) { was = compact(); reconcile(); } }, 120); }; })());

  window.wrmWorkboard = {
    init, show, close, toggle, isOpen: () => open,
    // bus is called by the page's event bus for every pushed message type.
    bus(type) {
      if (!S) return;
      if (type === 'events_changed') reloadType('events');
      if (type === 'sessions_changed' || type === 'status_changed' || type === 'tunnels_changed') setTimeout(() => S.tiles.forEach(x => { if (A[x.type].sig && bodies[x.id] && visible(bodies[x.id])) paint(x); }), 600);
    },
    langChanged() { if (!S) return; Object.values(bodies).forEach(b => { b._last = ''; }); Object.values(ctxs).forEach(c => { c.sig = ''; }); renderGrid(); if (window.wrmLayout) wrmLayout.refresh(); },
    _state: () => S, _sanitize: sanitize,
  };
})();
