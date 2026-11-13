/* Browser-local Web Serial manager. No serial data is sent to WRM. */
(function (root) {
  'use strict';
  const ids = {
    '2341:0043': ['Arduino Uno', 'arduino'], '2341:0001': ['Arduino Uno', 'arduino'],
    '2341:0010': ['Arduino Mega 2560', 'arduino'], '2341:0042': ['Arduino Mega 2560', 'arduino'],
    '2341:0058': ['Arduino Nano', 'arduino'], '2a03:0043': ['Arduino Uno', 'arduino'],
    '2a03:0010': ['Arduino Mega 2560', 'arduino'], '303a:1001': ['ESP32 native USB', 'esp32'],
    '0483:5740': ['STM32 virtual COM port', 'stm32'],
    '2e8a:000a': ['Raspberry Pi Pico', 'pico'], '2e8a:0005': ['Raspberry Pi Pico', 'pico'],
    '10c4:ea60': ['CP210x USB-UART bridge (often ESP32/ESP8266)', 'bridge'],
    '1a86:7523': ['CH340 USB-UART bridge (often Arduino/ESP)', 'bridge'],
    '0403:6001': ['FTDI USB-UART bridge', 'bridge']
  };
  const hex = n => Number(n || 0).toString(16).padStart(4, '0');
  function profileKey(info) { return `${hex(info.usbVendorId)}:${hex(info.usbProductId)}`; }
  function recognizeDevice(info) {
    const key = profileKey(info);
    const found = ids[key] || (key.startsWith('303a:') ? ['Espressif USB device', 'esp32'] :
      key.startsWith('0483:') ? ['STMicroelectronics USB device', 'stm32'] :
      key.startsWith('2e8a:') ? ['Raspberry Pi USB device', 'pico'] :
      key.startsWith('2341:') || key.startsWith('2a03:') ? ['Arduino USB device', 'arduino'] : null);
    return {name: found ? found[0] : 'Unknown USB serial device', type: found ? found[1] : 'unknown', vendorId: hex(info.usbVendorId), productId: hex(info.usbProductId)};
  }
  function suggestSettings(type) {
    return {baudRate: ['esp32', 'esp8266', 'stm32', 'pico'].includes(type) ? 115200 : 9600,
      dataBits: 8, stopBits: 1, parity: 'none', flowControl: 'none'};
  }
  function validateSettings(s) {
    const baudRate = Number(s.baudRate), dataBits = Number(s.dataBits), stopBits = Number(s.stopBits);
    if (!Number.isSafeInteger(baudRate) || baudRate < 1 || baudRate > 4000000 || ![7, 8].includes(dataBits) ||
        ![1, 2].includes(stopBits) || !['none', 'even', 'odd'].includes(s.parity) ||
        !['none', 'hardware'].includes(s.flowControl)) throw new Error('Invalid serial settings');
    return {baudRate, dataBits, stopBits, parity: s.parity, flowControl: s.flowControl};
  }
  if (typeof module !== 'undefined' && module.exports) module.exports = {recognizeDevice, suggestSettings, validateSettings, profileKey};
  if (!root || !root.document) return;

  Object.assign(LANGS.en, {
    serial_lbl:'Serial', serial_title:'Serial / USB Device Manager', serial_close:'Close',
    serial_unsupported:'Web Serial requires a supported Chromium browser and HTTPS (or localhost).',
    serial_select:'Select USB device', serial_known:'Previously authorized ports', serial_none:'No authorized devices yet.',
    serial_connect:'Connect', serial_disconnect:'Disconnect', serial_remove:'Forget profile', serial_save:'Save profile',
    serial_suggest:'Suggest settings', serial_name:'Profile name', serial_baud:'Baud rate', serial_bits:'Data bits',
    serial_stop:'Stop bits', serial_parity:'Parity', serial_flow:'Flow control', serial_line:'Line ending',
    serial_send:'Send', serial_input:'Text to send', serial_status:'Status', serial_offline:'Disconnected',
    serial_online:'Connected', serial_rx:'Received', serial_tx:'Sent', serial_errors:'Errors',
    serial_ai_recognize:'AI recognition prompt', serial_ai_code:'AI coding prompt',
    serial_ai_note:'Copies a prompt for the existing AI assistant. Review generated code before flashing hardware.',
    serial_prompt_copy:'Prompt copied. Paste it into your AI assistant.', serial_profile_saved:'Profile saved in this browser.',
    serial_unknown:'Unknown device', serial_device_removed:'Device disconnected', serial_choose:'Choose a device first.',
    serial_none_opt:'None', serial_even:'Even', serial_odd:'Odd', serial_hw:'Hardware (RTS/CTS)'
  });
  Object.assign(LANGS.hr, {
    serial_lbl:'Serijski', serial_title:'Upravitelj serijskih / USB uređaja', serial_close:'Zatvori',
    serial_unsupported:'Web Serial zahtijeva podržani Chromium preglednik i HTTPS (ili localhost).',
    serial_select:'Odaberi USB uređaj', serial_known:'Prethodno dopušteni priključci', serial_none:'Nema dopuštenih uređaja.',
    serial_connect:'Spoji', serial_disconnect:'Odspoji', serial_remove:'Obriši profil', serial_save:'Spremi profil',
    serial_suggest:'Predloži postavke', serial_name:'Naziv profila', serial_baud:'Brzina prijenosa', serial_bits:'Podatkovni bitovi',
    serial_stop:'Stop bitovi', serial_parity:'Paritet', serial_flow:'Kontrola protoka', serial_line:'Kraj retka',
    serial_send:'Pošalji', serial_input:'Tekst za slanje', serial_status:'Stanje', serial_offline:'Odspojeno',
    serial_online:'Spojeno', serial_rx:'Primljeno', serial_tx:'Poslano', serial_errors:'Greške',
    serial_ai_recognize:'Upit za AI prepoznavanje', serial_ai_code:'Upit za AI kodiranje',
    serial_ai_note:'Kopira upit za postojećeg AI asistenta. Provjerite kod prije učitavanja na uređaj.',
    serial_prompt_copy:'Upit je kopiran. Zalijepite ga u AI asistenta.', serial_profile_saved:'Profil je spremljen u ovom pregledniku.',
    serial_unknown:'Nepoznat uređaj', serial_device_removed:'Uređaj je odspojen', serial_choose:'Najprije odaberite uređaj.',
    serial_none_opt:'Bez', serial_even:'Paran', serial_odd:'Neparan', serial_hw:'Hardverska (RTS/CTS)'
  });
  const $ = id => root.document.getElementById(id);
  const tr = k => root.t(k);
  const safe = s => root.esc(s);
  const storeKey = () => 'wrm_serial_profiles:' + ((typeof currentUser !== 'undefined' && currentUser && (currentUser.id || currentUser.username)) || 'local');
  function profiles() { try { return JSON.parse(root.localStorage.getItem(storeKey()) || '{}'); } catch (_) { return {}; } }
  const state = {ports: [], port: null, terminal: null, profile: null, selected: null, rx: 0, tx: 0, errors: 0};
  function profileFor(port) { return profiles()[profileKey(port.getInfo())] || null; }
  const css = `#serial-ws{position:fixed;inset:0;z-index:900;background:var(--bg);color:var(--text);display:flex;flex-direction:column;font-size:13px}
    #serial-ws .s-head{display:flex;align-items:center;gap:12px;padding:12px 16px;background:var(--bg2);border-bottom:1px solid var(--border)}
    #serial-ws .s-head h2{font-size:16px;margin:0;flex:1} #serial-ws .s-body{padding:18px;overflow:auto;display:grid;grid-template-columns:minmax(240px,320px) minmax(0,1fr);gap:18px;flex:1;min-height:0}
    #serial-ws .s-side,#serial-ws .s-main{min-width:0;display:flex;flex-direction:column;gap:12px} #serial-ws .s-card{border:1px solid var(--border);border-radius:8px;background:var(--bg2);padding:14px}
    #serial-ws .s-device{display:block;width:100%;text-align:left;margin:6px 0;padding:9px;background:var(--bg);color:var(--text);border:1px solid var(--border);border-radius:6px;cursor:pointer}
    #serial-ws .s-device.on{border-color:var(--accent)} #serial-ws .s-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(125px,1fr));gap:10px} #serial-ws label{display:flex;flex-direction:column;gap:4px;color:var(--text2)}
    #serial-ws input,#serial-ws select{width:100%;min-width:0;padding:7px;background:var(--bg);color:var(--text);border:1px solid var(--border);border-radius:5px}
    #serial-ws .s-actions{display:flex;gap:8px;flex-wrap:wrap;margin-top:12px} #serial-ws .s-term{min-height:200px;flex:1;background:#111;overflow:hidden;padding:8px} #serial-ws .s-stats{display:flex;gap:15px;flex-wrap:wrap}
    @media(max-width:760px){#serial-ws .s-body{grid-template-columns:1fr}#serial-ws .s-term{height:300px;flex:none}}`;
  function option(values, value) { return values.map(([v, label]) => `<option value="${safe(v)}" ${String(v) === String(value) ? 'selected' : ''}>${safe(label)}</option>`).join(''); }
  function settingsForm(s, name) {
    return `<div class="s-grid"><label>${safe(tr('serial_name'))}<input id="serial-name" maxlength="80" value="${safe(name)}"></label>
      <label>${safe(tr('serial_baud'))}<input id="serial-baud" type="number" min="1" max="4000000" value="${s.baudRate}"></label>
      <label>${safe(tr('serial_bits'))}<select id="serial-bits">${option([[8,8],[7,7]],s.dataBits)}</select></label>
      <label>${safe(tr('serial_stop'))}<select id="serial-stop">${option([[1,1],[2,2]],s.stopBits)}</select></label>
      <label>${safe(tr('serial_parity'))}<select id="serial-parity">${option([['none',tr('serial_none_opt')],['even',tr('serial_even')],['odd',tr('serial_odd')]],s.parity)}</select></label>
      <label>${safe(tr('serial_flow'))}<select id="serial-flow">${option([['none',tr('serial_none_opt')],['hardware',tr('serial_hw')]],s.flowControl)}</select></label>
      <label>${safe(tr('serial_line'))}<select id="serial-line">${option([['',tr('serial_none_opt')],['\n','LF'],['\r\n','CRLF'],['\r','CR']],'\n')}</select></label></div>`;
  }
  function readForm() { return validateSettings({baudRate:$('serial-baud').value,dataBits:$('serial-bits').value,stopBits:$('serial-stop').value,parity:$('serial-parity').value,flowControl:$('serial-flow').value}); }
  function stats() { if ($('serial-stats')) $('serial-stats').textContent = `${tr('serial_status')}: ${tr(state.port ? 'serial_online':'serial_offline')} · ${tr('serial_rx')}: ${state.rx} B · ${tr('serial_tx')}: ${state.tx} B · ${tr('serial_errors')}: ${state.errors}`; }
  function render() {
    const chosen = state.selected;
    const d = chosen ? recognizeDevice(chosen.getInfo()) : null;
    const p = chosen ? profileFor(chosen) : null;
    const s = p ? p.settings : suggestSettings(d ? d.type : 'unknown');
    $('serial-ports').innerHTML = state.ports.length ? state.ports.map((port,i) => {
      const info = recognizeDevice(port.getInfo()); const saved = profileFor(port);
      return `<button class="s-device ${port === chosen ? 'on':''}" data-index="${i}">${safe(saved?.name || info.name)}<br><small>VID:PID ${info.vendorId}:${info.productId}</small></button>`;
    }).join('') : `<p>${safe(tr('serial_none'))}</p>`;
    $('serial-detail').innerHTML = d ? `<h3>${safe(p?.name || d.name)}</h3><p>USB VID:PID ${d.vendorId}:${d.productId} · ${safe(d.type)}</p>
      ${settingsForm(s,p?.name || d.name)}<div class="s-actions"><button data-action="suggest">${safe(tr('serial_suggest'))}</button><button data-action="save">${safe(tr('serial_save'))}</button><button data-action="remove">${safe(tr('serial_remove'))}</button>
      <button data-action="${state.port ? 'disconnect':'connect'}">${safe(tr(state.port ? 'serial_disconnect':'serial_connect'))}</button></div>
      <div class="s-actions"><button data-action="recognize">${safe(tr('serial_ai_recognize'))}</button><button data-action="code">${safe(tr('serial_ai_code'))}</button></div><p>${safe(tr('serial_ai_note'))}</p>` : `<p>${safe(tr('serial_choose'))}</p>`;
    stats();
    $('serial-input').disabled = $('serial-send').disabled = !state.port;
    if (state.terminal) setTimeout(() => state.terminal.fit(), 0);
  }
  async function refresh() {
    if (!root.navigator.serial) return;
    try { state.ports = await root.navigator.serial.getPorts(); if (state.selected && !state.ports.includes(state.selected)) state.selected = null; if (!state.selected) state.selected = state.ports[0] || null; render(); }
    catch (e) { root.showToast(e.message, 'error'); }
  }
  async function choose() {
    try { const port = await root.navigator.serial.requestPort(); state.ports = await root.navigator.serial.getPorts(); if (!state.ports.includes(port)) state.ports.push(port); state.selected = port; render(); }
    catch (e) { if (e.name !== 'NotFoundError') root.showToast(e.message, 'error'); }
  }
  async function disconnect() {
    const term = state.terminal; state.port = null; state.terminal = null;
    if (term) try { await term.disconnect(); } catch (e) { state.errors++; root.showToast(e.message, 'error'); }
    render();
  }
  async function connect() {
    if (!state.selected || state.port) return;
    try {
      const port = state.selected, settings = readForm();
      await port.open(settings);
      state.port = port;
      state.rx = state.tx = state.errors = 0;
      state.terminal = root.createSerialTerminal($('serial-terminal'), port, {
        onReceive: n => { state.rx += n; stats(); }, onSend: n => { state.tx += n; stats(); },
        onError: e => { state.errors++; stats(); root.showToast(e.message, 'error'); },
        onDisconnect: () => { if (state.port === port) disconnect(); }
      });
      render();
    } catch (e) { if (state.port) { try { await state.port.close(); } catch (_) {} state.port = null; } state.errors++; stats(); root.showToast(e.message, 'error'); }
  }
  async function prompt(kind) {
    const d = recognizeDevice(state.selected.getInfo()); const settings = readForm();
    const text = kind === 'recognize' ? `Identify this USB serial device. VID 0x${d.vendorId}, PID 0x${d.productId}, detected family: ${d.name}. Distinguish USB-UART bridge from the attached board; say when the board cannot be known from descriptors. Suggest baud/data bits/parity/stop bits and safe ways to verify.` :
      `Generate a minimal Arduino or PlatformIO serial example for ${d.name} (USB VID:PID ${d.vendorId}:${d.productId}), baud ${settings.baudRate}. Explain board selection and serial monitor configuration; do not assume a USB-UART bridge identifies the attached MCU.`;
    try { await root.navigator.clipboard.writeText(text); root.showToast(tr('serial_prompt_copy'), 'success'); }
    catch (e) { root.showToast(e.message, 'error'); }
  }
  async function action(a) {
    if (a === 'suggest') { const s = suggestSettings(recognizeDevice(state.selected.getInfo()).type); for (const [id,v] of Object.entries({baud:s.baudRate,bits:s.dataBits,stop:s.stopBits,parity:s.parity,flow:s.flowControl})) $('serial-'+id).value = v; return; }
    if (a === 'save') { try { const all = profiles(); all[profileKey(state.selected.getInfo())] = {name:$('serial-name').value.trim().slice(0,80),settings:readForm()}; root.localStorage.setItem(storeKey(),JSON.stringify(all)); render(); root.showToast(tr('serial_profile_saved'),'success'); } catch(e) { root.showToast(e.message,'error'); } return; }
    if (a === 'remove') { const all = profiles(); delete all[profileKey(state.selected.getInfo())]; root.localStorage.setItem(storeKey(),JSON.stringify(all)); render(); return; }
    if (a === 'connect') return connect();
    if (a === 'disconnect') return disconnect();
    if (a === 'recognize' || a === 'code') return prompt(a);
  }
  function open() {
    if (!$('serial-ws')) {
      const style = root.document.createElement('style'); style.textContent = css; root.document.head.appendChild(style);
      const ws = root.document.createElement('div'); ws.id = 'serial-ws';
      ws.innerHTML = `<div class="s-head"><h2>${safe(tr('serial_title'))}</h2><button id="serial-close" aria-label="${safe(tr('serial_close'))}">✕</button></div>
        <div class="s-body"><div class="s-side"><div class="s-card"><button id="serial-choose">${safe(tr('serial_select'))}</button><h3>${safe(tr('serial_known'))}</h3><div id="serial-ports"></div></div><div class="s-card" id="serial-stats"></div></div>
        <div class="s-main"><div class="s-card" id="serial-detail"></div><div id="serial-terminal" class="s-term"></div><div class="s-actions"><input id="serial-input" aria-label="${safe(tr('serial_input'))}" placeholder="${safe(tr('serial_input'))}"><button id="serial-send">${safe(tr('serial_send'))}</button></div></div></div>`;
      root.document.body.appendChild(ws);
      $('serial-close').onclick = close;
      $('serial-choose').onclick = choose;
      $('serial-ports').onclick = e => { const btn = e.target.closest('[data-index]'); if (!btn) return; if (state.port) return; state.selected = state.ports[Number(btn.dataset.index)]; render(); };
      $('serial-detail').onclick = e => { const btn = e.target.closest('[data-action]'); if (btn) action(btn.dataset.action); };
      $('serial-send').onclick = send;
      $('serial-input').onkeydown = e => { if (e.key === 'Enter') send(); };
      root.navigator.serial?.addEventListener('disconnect', e => { if (e.target === state.port) disconnect(); refresh(); root.showToast(tr('serial_device_removed'),'warning'); });
      root.navigator.serial?.addEventListener('connect', refresh);
    }
    $('serial-ws').style.display = 'flex';
    $('serial-choose').disabled = !root.isSecureContext || !root.navigator.serial;
    if ($('serial-warning')) $('serial-warning').remove();
    if ($('serial-choose').disabled) { const warning = root.document.createElement('p'); warning.id = 'serial-warning'; warning.textContent = tr('serial_unsupported'); $('serial-ports').before(warning); }
    refresh();
  }
  function send() { const input = $('serial-input'); if (!state.port || !input.value) return; const bytes = new TextEncoder().encode(input.value + $('serial-line').value); state.terminal.send(bytes).then(() => { input.value = ''; }, e => { state.errors++; stats(); root.showToast(e.message,'error'); }); }
  function close() { if ($('serial-ws')) $('serial-ws').style.display = 'none'; }
  root.serialWorkspace = {open, close, refresh};
})(typeof window !== 'undefined' ? window : null);
