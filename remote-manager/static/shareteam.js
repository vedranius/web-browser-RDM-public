// Team-managed shares (#192): members with "operator" or higher may add and edit the
// shared connections and folders, start/stop their tunnels and add/edit those tunnels
// when the owner allows it (three separate switches).
'use strict';

Object.assign(LANGS.en, {
  share_team_title: 'Team management',
  share_allow_edit: 'Operators may add and edit connections and folders',
  share_allow_tunnels: 'Operators may enable and disable (start and stop) tunnels',
  share_allow_tunnel_edit: 'Operators may add, edit and remove tunnels',
  st_edit_tunnels: 'Edit tunnels…',
  st_edit_tunnels_of: 'Tunnels of',
  st_tun_add: '+ Add tunnel',
  st_tun_kind: 'Type',
  st_tun_listen: 'Listen (address : port, 0 = automatic)',
  st_tun_target: 'Target (host : port)',
  st_tun_start: 'Start',
  st_tun_hint: 'Tunnels stay in the owner’s account. Listening on other addresses than 127.0.0.1 and remote forwarding follow the tunnel policies of both the owner and you.',
  st_tun_none: 'No tunnels yet.',
  share_team_hint: 'Applies to signed-in people with the operator role or higher, never to guests. New connections and folders stay in your account. Members never see or choose your saved secrets; when a member changes a host, its saved password is removed.',
  st_edit: 'Edit connection…',
  st_new_conn: 'New connection here…',
  st_new_folder: 'New folder…',
  st_rename_folder: 'Rename folder…',
  st_tunnels: 'Tunnels…',
  st_tunnels_of: 'Tunnels of the share',
  st_no_tunnels: 'The shared connections have no tunnels yet. The owner configures them.',
  st_tunnels_denied: 'Tunnels are not allowed for your account or the owner’s account.',
  st_host_hint: 'Changing the host or protocol removes the saved password (and saved key), so it is never sent to a different server.',
  st_password_keep: 'leave empty to keep it',
  st_saved: 'Saved',
  st_folder: 'Folder',
  st_name: 'Name',
  st_host: 'Host',
  st_notes: 'Notes',
  st_start: 'Start',
  st_stop: 'Stop',
  st_running: 'running',
  st_stopped: 'stopped',
});
Object.assign(LANGS.hr, {
  share_team_title: 'Timsko upravljanje',
  share_allow_edit: 'Operateri smiju dodavati i uređivati veze i mape',
  share_allow_tunnels: 'Operateri smiju uključivati i isključivati (pokretati i zaustavljati) tunele',
  share_allow_tunnel_edit: 'Operateri smiju dodavati, uređivati i brisati tunele',
  st_edit_tunnels: 'Uredi tunele…',
  st_edit_tunnels_of: 'Tuneli veze',
  st_tun_add: '+ Dodaj tunel',
  st_tun_kind: 'Vrsta',
  st_tun_listen: 'Slušanje (adresa : port, 0 = automatski)',
  st_tun_target: 'Cilj (host : port)',
  st_tun_start: 'Pokretanje',
  st_tun_hint: 'Tuneli ostaju na računu vlasnika. Slušanje na adresama osim 127.0.0.1 i udaljeno prosljeđivanje slijede pravila tunela i vlasnika i vas.',
  st_tun_none: 'Još nema tunela.',
  share_team_hint: 'Vrijedi za prijavljene osobe s ulogom operatera ili višom, nikad za goste. Nove veze i mape ostaju na vašem računu. Članovi nikad ne vide niti biraju vaše spremljene tajne; kad član promijeni host, spremljena lozinka se briše.',
  st_edit: 'Uredi vezu…',
  st_new_conn: 'Nova veza ovdje…',
  st_new_folder: 'Nova mapa…',
  st_rename_folder: 'Preimenuj mapu…',
  st_tunnels: 'Tuneli…',
  st_tunnels_of: 'Tuneli dijeljenja',
  st_no_tunnels: 'Dijeljene veze još nemaju tunele. Postavlja ih vlasnik.',
  st_tunnels_denied: 'Tuneli nisu dopušteni vašem računu ili računu vlasnika.',
  st_host_hint: 'Promjena hosta ili protokola briše spremljenu lozinku (i spremljeni ključ), da se nikad ne pošalje drugom poslužitelju.',
  st_password_keep: 'ostavite prazno da ostane',
  st_saved: 'Spremljeno',
  st_folder: 'Mapa',
  st_name: 'Naziv',
  st_host: 'Host',
  st_notes: 'Bilješke',
  st_start: 'Pokreni',
  st_stop: 'Zaustavi',
  st_running: 'radi',
  st_stopped: 'zaustavljen',
});

const ST_PROTOCOLS = ['SSH', 'SFTP', 'FTP', 'FTPS', 'RDP', 'VNC', 'TELNET', 'HTTP', 'HTTPS'];

function stGroup(token) {
  if (typeof isCollabMode !== 'undefined' && isCollabMode && token === collab.token) return {folders: collabFolders || [], share: collab.share || {}};
  return (sharedGroups || []).find(g => g.share && g.share.token === token) || {folders: [], share: {}};
}

// Reloads the tree after a change: the main tree, or the share room's tree.
async function stReload() {
  if (typeof isCollabMode !== 'undefined' && isCollabMode) {
    try {
      const r = await fetch('/api/share/' + encodeURIComponent(collab.id));
      if (r.ok) { const d = await r.json(); collabConns = d.connections || []; collabFolders = d.folders || []; renderCollabTree(); }
    } catch (e) { /* the next reload shows it */ }
    return;
  }
  loadData();
}

async function stFetch(token, path, method, body) {
  const r = await fetch('/api/share-team/' + encodeURIComponent(token) + path, {
    method: method || 'GET',
    headers: {'Content-Type': 'application/json', 'X-WRM-Request': '1'},
    body: body ? JSON.stringify(body) : undefined,
  });
  if (!r.ok) throw new Error(await apiError(r));
  return r.json();
}

// Buttons on the share's label in the tree.
function shareTeamGroupButtons(sh) {
  const p = sh.perms || {};
  const tok = esc(JSON.stringify(sh.token));
  let h = '';
  if (p.manage_tunnels) h += `<button class="icon-btn" style="padding:2px 5px;font-size:11px;" title="${esc(t('st_tunnels'))}" onclick="shareTeamTunnels(${tok})">⇄</button>`;
  return h;
}

// Items for the context menu of a shared connection.
function shareTeamCtxItems(conn, token, perms) {
  const tok = esc(JSON.stringify(token));
  let h = '';
  if (perms.manage_connections) h += `<div class="ctx-item" onclick="shareTeamEditConn(${tok}, ${Number(conn.id)})">✏️ ${esc(t('st_edit'))}</div>`;
  const ssh = conn.protocol === 'SSH' || conn.protocol === 'SFTP';
  if (perms.manage_tunnels && ssh) h += `<div class="ctx-item" onclick="shareTeamTunnels(${tok})">⇄ ${esc(t('st_tunnels'))}</div>`;
  if (perms.edit_tunnels && ssh) h += `<div class="ctx-item" onclick="shareTeamEditTunnels(${tok}, ${Number(conn.id)}, ${esc(JSON.stringify(conn.name || ''))})">🛠 ${esc(t('st_edit_tunnels'))}</div>`;
  return h ? '<div class="ctx-hr"></div>' + h : '';
}

// Context menu of a shared folder.
function shareTeamFolderCtx(e, f, sh) {
  const p = sh.perms || {};
  if (!p.manage_connections) return;
  e.preventDefault(); e.stopPropagation();
  const tok = esc(JSON.stringify(sh.token));
  const menu = document.getElementById('dyn-ctx-menu');
  menu.innerHTML = `<div class="ctx-section">${esc(f.name)}</div>
    <div class="ctx-item" onclick="shareTeamEditConn(${tok}, 0, ${Number(f.id)})">＋ ${esc(t('st_new_conn'))}</div>
    <div class="ctx-item" onclick="shareTeamNewFolder(${tok}, ${Number(f.id)})">📁 ${esc(t('st_new_folder'))}</div>
    <div class="ctx-item" onclick="shareTeamRenameFolder(${tok}, ${Number(f.id)})">✏️ ${esc(t('st_rename_folder'))}</div>`;
  showCtx('dyn-ctx-menu', e.clientX, e.clientY);
}

function stModal(id, inner) {
  document.getElementById(id)?.remove();
  const m = document.createElement('div');
  m.id = id; m.className = 'modal-overlay open'; m.style.zIndex = 9700;
  m.innerHTML = `<div class="modal" style="width:min(520px,calc(100vw - 24px))">${inner}</div>`;
  m.addEventListener('click', ev => { if (ev.target === m) m.remove(); });
  document.body.appendChild(m);
  return m;
}

async function shareTeamEditConn(token, id, folderId) {
  let c = {name: '', protocol: 'SSH', host: '', username: '', folder_id: folderId || null, tags: [], notes: '', has_password: false};
  if (id) {
    try { c = await stFetch(token, '/connections/' + id); } catch (e) { showToast(e.message, 'error'); return; }
  }
  const g = stGroup(token);
  const folders = (g.folders || []).map(f => `<option value="${f.id}"${f.id === c.folder_id ? ' selected' : ''}>${esc(f.name)}</option>`).join('');
  const protos = ST_PROTOCOLS.concat(ST_PROTOCOLS.includes(c.protocol) ? [] : [c.protocol])
    .map(p => `<option${p === c.protocol ? ' selected' : ''}>${esc(p)}</option>`).join('');
  const m = stModal('st-conn-modal', `
    <h3>${esc(id ? t('st_edit') : t('st_new_conn'))} <span style="color:var(--text3);font-weight:500;">· ${esc((g.share || {}).name || '')}</span></h3>
    <div class="field"><label for="st-name">${esc(t('st_name'))}</label><input id="st-name" value="${esc(c.name)}"></div>
    <div class="field"><label for="st-proto">${esc(t('protocol'))}</label><select id="st-proto">${protos}</select></div>
    <div class="field"><label for="st-host">${esc(t('st_host'))}</label><input id="st-host" value="${esc(c.host)}" autocomplete="off"></div>
    <div class="field"><label for="st-user">${esc(t('username'))}</label><input id="st-user" value="${esc(c.username)}" autocomplete="off"></div>
    <div class="field"><label for="st-pass">${esc(t('password'))}</label><input id="st-pass" type="password" autocomplete="new-password" placeholder="${c.has_password ? esc(t('st_password_keep')) : ''}"></div>
    <div class="field"><label for="st-folder">${esc(t('st_folder'))}</label><select id="st-folder">${folders}</select></div>
    <div class="field"><label for="st-notes">${esc(t('st_notes'))}</label><textarea id="st-notes" rows="3">${esc(c.notes || '')}</textarea></div>
    ${id ? `<p style="font-size:11.5px;color:var(--text2);">${esc(t('st_host_hint'))}</p>` : ''}
    <div class="modal-footer"><button class="btn-sec" id="st-cancel">${esc(t('cancel'))}</button><button id="st-save">${esc(t('save'))}</button></div>`);
  m.querySelector('#st-cancel').onclick = () => m.remove();
  m.querySelector('#st-save').onclick = async () => {
    const body = {
      name: m.querySelector('#st-name').value.trim(), protocol: m.querySelector('#st-proto').value,
      host: m.querySelector('#st-host').value.trim(), username: m.querySelector('#st-user').value.trim(),
      folder_id: Number(m.querySelector('#st-folder').value) || null, notes: m.querySelector('#st-notes').value,
    };
    const pw = m.querySelector('#st-pass').value;
    if (pw) body.password = pw;
    try {
      await stFetch(token, '/connections' + (id ? '/' + id : ''), id ? 'PUT' : 'POST', body);
      m.remove(); showToast(t('st_saved'), 'success'); stReload();
    } catch (e) { showToast(e.message, 'error'); }
  };
}

async function shareTeamNewFolder(token, parentId) {
  const name = await uiPrompt(t('st_new_folder'), '');
  if (!name) return;
  try { await stFetch(token, '/folders', 'POST', {name, parent_id: parentId}); stReload(); } catch (e) { showToast(e.message, 'error'); }
}

async function shareTeamRenameFolder(token, id) {
  const f = (stGroup(token).folders || []).find(x => x.id === id) || {};
  const name = await uiPrompt(t('st_rename_folder'), f.name || '');
  if (!name) return;
  try { await stFetch(token, '/folders/' + id, 'PUT', {name}); stReload(); } catch (e) { showToast(e.message, 'error'); }
}

async function shareTeamTunnels(token) {
  let data;
  try { data = await stFetch(token, '/tunnels'); } catch (e) { showToast(e.message, 'error'); return; }
  const tok = esc(JSON.stringify(token));
  const rows = (data.tunnels || []).map(x => {
    const on = x.state && x.state !== 'stopped';
    const bind = x.kind === 'remote' ? '' : esc((x.bind_host || '127.0.0.1') + ':' + (x.bind_port || ''));
    return `<tr><td>${esc(x.name || x.key)}<div class="sub">${esc(x.conn_name || '')}</div></td>
      <td>${bind} → ${esc((x.target_host || '') + ':' + (x.target_port || ''))}</td>
      <td>${esc(on ? t('st_running') : t('st_stopped'))}</td>
      <td>${data.allowed ? `<button class="btn-sec" onclick="shareTeamTunnelAction(${tok}, ${esc(JSON.stringify(x.key))}, '${on ? 'stop' : 'start'}')">${esc(on ? t('st_stop') : t('st_start'))}</button>` : ''}</td></tr>`;
  }).join('');
  const m = stModal('st-tun-modal', `<h3>${esc(t('st_tunnels_of'))} · ${esc((stGroup(token).share || {}).name || '')}</h3>
    ${data.allowed ? '' : `<p style="color:var(--yellow);">${esc(t('st_tunnels_denied'))}</p>`}
    ${rows ? `<div style="overflow:auto;"><table class="data-table" style="width:100%;">${rows}</table></div>` : `<p>${esc(t('st_no_tunnels'))}</p>`}
    <div class="modal-footer"><button class="btn-sec" id="st-tun-close">${esc(t('close'))}</button></div>`);
  m.querySelector('#st-tun-close').onclick = () => m.remove();
}

async function shareTeamTunnelAction(token, key, action) {
  try { await stFetch(token, '/tunnels/' + encodeURIComponent(key) + '/' + action, 'POST'); } catch (e) { showToast(e.message, 'error'); }
  shareTeamTunnels(token);
}

// Editor for the tunnels of one shared SSH connection ("add and edit tunnels" switch).
async function shareTeamEditTunnels(token, connId, connName) {
  let defs;
  try { defs = await stFetch(token, '/connections/' + connId + '/tunnels'); } catch (e) { showToast(e.message, 'error'); return; }
  const kinds = ['local', 'remote', 'dynamic'], modes = ['manual', 'connect', 'always'];
  const opt = (list, v) => list.map(k => `<option${k === v ? ' selected' : ''}>${k}</option>`).join('');
  const row = d => `<div class="st-tun-row" data-id="${Number(d.id) || 0}">
      <div class="st-t-line"><input class="st-t-name" placeholder="${esc(t('st_name'))}" value="${esc(d.name || '')}" aria-label="${esc(t('st_name'))}">
        <button class="btn-sec st-t-del" title="${esc(t('delete'))}" aria-label="${esc(t('delete'))}">✕</button></div>
      <div class="st-t-line"><select class="st-t-kind" aria-label="${esc(t('st_tun_kind'))}">${opt(kinds, d.kind || 'local')}</select>
        <select class="st-t-mode" aria-label="${esc(t('st_tun_start'))}">${opt(modes, d.start_mode || 'manual')}</select></div>
      <div class="st-t-line" title="${esc(t('st_tun_listen'))}"><input class="st-t-bh" value="${esc(d.bind_host || '127.0.0.1')}" aria-label="${esc(t('st_tun_listen'))}"><input class="st-t-bp st-t-port" type="number" min="0" max="65535" value="${Number(d.bind_port) || 0}" aria-label="port"></div>
      <div class="st-t-line" title="${esc(t('st_tun_target'))}"><input class="st-t-th" value="${esc(d.target_host || '')}" placeholder="host" aria-label="${esc(t('st_tun_target'))}"><input class="st-t-tp st-t-port" type="number" min="0" max="65535" value="${Number(d.target_port) || ''}" aria-label="port"></div>
    </div>`;
  const m = stModal('st-tunedit-modal', `<h3>${esc(t('st_edit_tunnels_of'))} · ${esc(connName)}</h3>
    <style>#st-tunedit-modal .st-tun-row{border:1px solid var(--border);border-radius:6px;padding:6px;margin-bottom:6px;display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:4px 8px}
    #st-tunedit-modal .st-t-line{display:flex;gap:4px;min-width:0}#st-tunedit-modal .st-t-line>*{min-width:0;flex:1 1 auto;width:auto}
    #st-tunedit-modal .st-t-port{flex:0 0 76px}#st-tunedit-modal .st-t-del{flex:0 0 auto;padding:2px 8px}</style>
    <div style="font-size:11px;color:var(--text3);margin-bottom:6px;">${esc(t('st_tun_listen'))} · ${esc(t('st_tun_target'))}</div>
    <div id="st-tun-rows" style="max-height:55vh;overflow:auto;">${defs.length ? defs.map(row).join('') : ''}</div>
    <p id="st-tun-none" style="${defs.length ? 'display:none;' : ''}">${esc(t('st_tun_none'))}</p>
    <button class="btn-sec" id="st-tun-add">${esc(t('st_tun_add'))}</button>
    <p style="font-size:11.5px;color:var(--text2);">${esc(t('st_tun_hint'))}</p>
    <div class="modal-footer"><button class="btn-sec" id="st-tun-cancel">${esc(t('cancel'))}</button><button id="st-tun-save">${esc(t('save'))}</button></div>`);
  const rows = m.querySelector('#st-tun-rows');
  const bindDel = () => rows.querySelectorAll('.st-t-del').forEach(b => { b.onclick = () => { b.closest('.st-tun-row').remove(); }; });
  bindDel();
  m.querySelector('#st-tun-add').onclick = () => {
    rows.insertAdjacentHTML('beforeend', row({kind: 'local', bind_host: '127.0.0.1', target_host: '127.0.0.1'}));
    m.querySelector('#st-tun-none').style.display = 'none'; bindDel();
  };
  m.querySelector('#st-tun-cancel').onclick = () => m.remove();
  m.querySelector('#st-tun-save').onclick = async () => {
    const body = [...rows.querySelectorAll('.st-tun-row')].map(r => ({
      id: Number(r.dataset.id) || 0, name: r.querySelector('.st-t-name').value.trim(),
      kind: r.querySelector('.st-t-kind').value, start_mode: r.querySelector('.st-t-mode').value,
      bind_host: r.querySelector('.st-t-bh').value.trim(), bind_port: Number(r.querySelector('.st-t-bp').value) || 0,
      target_host: r.querySelector('.st-t-th').value.trim(), target_port: Number(r.querySelector('.st-t-tp').value) || 0,
    }));
    try {
      await stFetch(token, '/connections/' + connId + '/tunnels', 'PUT', body);
      m.remove(); showToast(t('st_saved'), 'success'); stReload();
    } catch (e) { showToast(e.message, 'error'); }
  };
}
