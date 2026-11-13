// Import from another tool: one picker for every format the server understands
// (GET /api/config/import/formats). The file is sent to POST /api/config/import/auto with
// "preview": true first; the dialog shows what would be created (connections, folders, jump
// hosts, skipped entries) and only "Import" writes. Passwords of encrypted files are asked
// inside the dialog and are never kept in the browser after it closes.
'use strict';

Object.assign(LANGS.en, {
  import_other: 'Import from another tool',
  import_other_desc: 'mRemoteNG (XML or CSV), Remote Desktop Manager CSV, RDCMan (.rdg), .rdp files, MobaXterm, WinSCP, Remmina, OpenSSH config and PuTTY (.reg). You see a preview before anything is saved.',
  imp_title: 'Import connections',
  imp_format: 'Format',
  imp_auto: 'Detect automatically',
  imp_file: 'File',
  imp_base: 'Put everything into folder (optional)',
  imp_base_ph: 'e.g. Imported',
  imp_password: 'Master password of the file',
  imp_password_hint: 'The file is encrypted or its passwords are. The password is used only for this import.',
  imp_preview: 'Preview',
  imp_do: 'Import',
  imp_detected: 'Detected format: {f}',
  imp_unknown: 'The format was not recognised. Choose it in the list above.',
  imp_summary: '{n} connections, {f} new folders, {j} jump hosts, {t} tunnels',
  imp_dupes_note: 'A connection whose name, host, protocol and folder match an existing one is skipped, so importing again is safe.',
  imp_col_name: 'Name', imp_col_proto: 'Protocol', imp_col_host: 'Host', imp_col_user: 'User', imp_col_folder: 'Folder', imp_col_jump: 'Jump host',
  imp_secret: 'has a saved password or key',
  imp_more: '… and {n} more',
  imp_nothing: 'Nothing to import in this file.',
  imp_choose_file: 'Choose a file first.',
  import_updated_n: '{n} existing connections updated',
});
Object.assign(LANGS.hr, {
  import_other: 'Uvoz iz drugog alata',
  import_other_desc: 'mRemoteNG (XML ili CSV), Remote Desktop Manager CSV, RDCMan (.rdg), .rdp datoteke, MobaXterm, WinSCP, Remmina, OpenSSH config i PuTTY (.reg). Prije spremanja vidite pregled.',
  imp_title: 'Uvoz veza',
  imp_format: 'Format',
  imp_auto: 'Prepoznaj automatski',
  imp_file: 'Datoteka',
  imp_base: 'Sve stavi u mapu (neobavezno)',
  imp_base_ph: 'npr. Uvezeno',
  imp_password: 'Glavna lozinka datoteke',
  imp_password_hint: 'Datoteka ili lozinke u njoj su šifrirane. Lozinka se koristi samo za ovaj uvoz.',
  imp_preview: 'Pregled',
  imp_do: 'Uvezi',
  imp_detected: 'Prepoznat format: {f}',
  imp_unknown: 'Format nije prepoznat. Odaberite ga na popisu iznad.',
  imp_summary: '{n} veza, {f} novih mapa, {j} jump hostova, {t} tunela',
  imp_dupes_note: 'Veza čije se ime, host, protokol i mapa poklapaju s postojećom se preskače, pa je ponovni uvoz siguran.',
  imp_col_name: 'Ime', imp_col_proto: 'Protokol', imp_col_host: 'Host', imp_col_user: 'Korisnik', imp_col_folder: 'Mapa', imp_col_jump: 'Jump host',
  imp_secret: 'ima spremljenu lozinku ili ključ',
  imp_more: '… i još {n}',
  imp_nothing: 'U ovoj datoteci nema ničega za uvoz.',
  imp_choose_file: 'Najprije odaberite datoteku.',
  import_updated_n: 'Ažurirano postojećih veza: {n}',
});

let impFormats = null;

function impB64(bytes) {
  let bin = '';
  for (let i = 0; i < bytes.length; i += 0x8000) bin += String.fromCharCode.apply(null, bytes.subarray(i, i + 0x8000));
  return btoa(bin);
}

// openImportPicker opens the dialog; file (optional) is preselected, e.g. from the JSON import row.
async function openImportPicker(file) {
  if (!impFormats) {
    try { impFormats = await (await fetch('/api/config/import/formats')).json(); } catch (e) { impFormats = []; }
    if (!Array.isArray(impFormats)) impFormats = [];
  }
  document.getElementById('imp-modal')?.remove();
  const m = document.createElement('div');
  m.id = 'imp-modal'; m.className = 'modal-overlay open'; m.style.zIndex = 9700;
  const opts = `<option value="">${esc(t('imp_auto'))}</option>` + impFormats.map(f => `<option value="${esc(f.ID)}">${esc(f.Name)}</option>`).join('');
  m.innerHTML = `<div class="modal" style="width:min(760px,calc(100vw - 24px))">
    <h3>${esc(t('imp_title'))}</h3>
    <div class="field"><label for="imp-file">${esc(t('imp_file'))}</label><input type="file" id="imp-file"></div>
    <div class="field"><label for="imp-format">${esc(t('imp_format'))}</label><select id="imp-format">${opts}</select></div>
    <div class="field"><label for="imp-base">${esc(t('imp_base'))}</label><input id="imp-base" maxlength="200" placeholder="${esc(t('imp_base_ph'))}"></div>
    <div class="field" id="imp-pw-f" style="display:none;"><label for="imp-pw">${esc(t('imp_password'))}</label><input id="imp-pw" type="password" autocomplete="off">
      <div style="font-size:11.5px;color:var(--text2);margin-top:4px;" id="imp-pw-hint">${esc(t('imp_password_hint'))}</div></div>
    <div id="imp-out" style="max-height:45vh;overflow:auto;font-size:12.5px;"></div>
    <div class="modal-footer"><button class="btn-sec" id="imp-cancel">${esc(t('cancel'))}</button>
      <button class="btn-sec" id="imp-prev">${esc(t('imp_preview'))}</button><button id="imp-go" disabled>${esc(t('imp_do'))}</button></div>
  </div>`;
  m.addEventListener('click', ev => { if (ev.target === m) m.remove(); });
  document.body.appendChild(m);
  const $ = s => m.querySelector(s);
  let data = null, name = '';
  const reset = () => { $('#imp-go').disabled = true; $('#imp-out').innerHTML = ''; };
  const useFile = async f => {
    reset(); data = null; name = '';
    if (!f) return;
    if (importFileTooBig(f)) { $('#imp-file').value = ''; return; }
    data = impB64(new Uint8Array(await f.arrayBuffer())); name = f.name;
    preview();
  };
  $('#imp-file').onchange = () => useFile($('#imp-file').files[0]);
  $('#imp-format').onchange = () => { if (data) preview(); };
  $('#imp-base').onchange = () => { if (data) preview(); };
  $('#imp-cancel').onclick = () => m.remove();
  $('#imp-prev').onclick = () => preview();
  $('#imp-go').onclick = () => run();

  async function call(dry) {
    const body = {data, name, format: $('#imp-format').value, folder: $('#imp-base').value.trim(), password: $('#imp-pw').value, preview: dry};
    const r = await fetch('/api/config/import/auto', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)});
    return {r, d: await r.json().catch(() => ({}))};
  }
  async function preview() {
    if (!data) { showToast(t('imp_choose_file'), 'error'); return; }
    reset();
    const {r, d} = await call(true);
    const out = $('#imp-out');
    if (!r.ok) {
      if (d.need_password || /master password/i.test(d.error || '')) {
        $('#imp-pw-f').style.display = '';
        $('#imp-pw-hint').textContent = $('#imp-pw').value ? t('mrng_wrong_pw') : t('imp_password_hint');
        $('#imp-pw').focus();
        $('#imp-pw').onkeydown = ev => { if (ev.key === 'Enter') preview(); };
      }
      out.innerHTML = `<p style="color:var(--red,#e5534b);">${esc(d.error || t('failed'))}</p>` + (/^unknown file format/i.test(d.error || '') ? `<p>${esc(t('imp_unknown'))}</p>` : '');
      return;
    }
    out.innerHTML = impRender(d);
    $('#imp-go').disabled = !(d.items || []).length;
  }
  async function run() {
    $('#imp-go').disabled = true;
    const {r, d} = await call(false);
    if (!r.ok) { showToast(d.error || t('failed'), 'error', 8000); $('#imp-go').disabled = false; return; }
    m.remove();
    const lines = [t('import_result').replace('{n}', d.imported).replace('{f}', d.folders).replace('{j}', d.jump_hosts).replace('{t}', d.tunnels)];
    if (d.updated) lines.push(t('import_updated_n').replace('{n}', d.updated));
    if (d.proxies || d.proxy_links) lines.push(t('import_result_proxies').replace('{p}', d.proxies || 0).replace('{l}', d.proxy_links || 0));
    if ((d.skipped || []).length) lines.push('', t('import_skipped_h') + ':', ...d.skipped.slice(0, 30).map(x => '• ' + x.name + ' — ' + x.reason));
    if ((d.notes || []).length) lines.push('', t('import_notes_h') + ':', ...d.notes.slice(0, 30).map(x => '• ' + x.name + ' — ' + x.reason));
    loadData();
    await uiDialog({title: t('import_done_title'), message: lines.join('\n'), okText: t('ok'), cancelText: t('close')});
  }
  if (file) useFile(file);
}

function impRender(d) {
  const fmt = (impFormats || []).find(f => f.ID === d.format);
  const items = d.items || [];
  let h = `<p><b>${esc(t('imp_detected').replace('{f}', fmt ? fmt.Name : (d.format || '?')))}</b><br>` +
    esc(t('imp_summary').replace('{n}', items.length).replace('{f}', d.folders || 0).replace('{j}', d.jump_hosts || 0).replace('{t}', d.tunnels || 0)) + '</p>';
  if (!items.length) h += `<p>${esc(t('imp_nothing'))}</p>`;
  else {
    const max = 300;
    h += `<table class="imp-table" style="width:100%;border-collapse:collapse;"><thead><tr>${['imp_col_name', 'imp_col_proto', 'imp_col_host', 'imp_col_user', 'imp_col_folder', 'imp_col_jump']
      .map(k => `<th style="text-align:left;padding:3px 6px;border-bottom:1px solid var(--border);">${esc(t(k))}</th>`).join('')}</tr></thead><tbody>` +
      items.slice(0, max).map(x => `<tr><td style="padding:2px 6px;">${esc(x.name)}${x.has_secret ? ` <span title="${esc(t('imp_secret'))}">🔑</span>` : ''}</td><td style="padding:2px 6px;">${esc(x.protocol)}</td>` +
        `<td style="padding:2px 6px;">${esc(x.host)}</td><td style="padding:2px 6px;">${esc(x.username || '')}</td><td style="padding:2px 6px;">${esc(x.folder || '')}</td><td style="padding:2px 6px;">${esc(x.jump || '')}</td></tr>`).join('') +
      '</tbody></table>';
    if (items.length > max) h += `<p>${esc(t('imp_more').replace('{n}', items.length - max))}</p>`;
    h += `<p style="color:var(--text2);font-size:11.5px;">${esc(t('imp_dupes_note'))}</p>`;
  }
  if ((d.skipped || []).length) h += `<p><b>${esc(t('import_skipped_h'))}</b></p><ul>` + d.skipped.slice(0, 50).map(x => `<li>${esc(x.name)} — ${esc(x.reason)}</li>`).join('') + '</ul>';
  if ((d.notes || []).length) h += `<p><b>${esc(t('import_notes_h'))}</b></p><ul>` + d.notes.slice(0, 50).map(x => `<li>${esc(x.name)} — ${esc(x.reason)}</li>`).join('') + '</ul>';
  return h;
}
