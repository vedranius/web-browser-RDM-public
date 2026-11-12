// Team-managed shares (#192): members with "operator" or higher may add and edit the
// shared connections and folders and start/stop their tunnels when the owner allows it.
'use strict';

Object.assign(LANGS.en, {
  share_team_title: 'Team management',
  share_allow_edit: 'Operators may add and edit connections and folders',
  share_allow_tunnels: 'Operators may start and stop tunnels',
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
  share_allow_tunnels: 'Operateri smiju pokretati i zaustavljati tunele',
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
  if (perms.manage_tunnels && (conn.protocol === 'SSH' || conn.protocol === 'SFTP')) h += `<div class="ctx-item" onclick="shareTeamTunnels(${tok})">⇄ ${esc(t('st_tunnels'))}</div>`;
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
