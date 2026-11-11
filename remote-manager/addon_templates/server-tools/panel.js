(function () {
  'use strict';
  const $ = (id) => document.getElementById(id);
  let conns = [];

  async function showMetrics() {
    const id = Number($('conn').value);
    if (!id) return;
    $('metrics').textContent = wrm.t('Reading metrics…', 'Čitam metrike…');
    try {
      const m = await wrm.metrics(id);
      const parts = [];
      if (m.cpu_percent !== undefined) parts.push('CPU ' + Math.round(m.cpu_percent) + ' %');
      if (m.memory_percent !== undefined) parts.push('RAM ' + Math.round(m.memory_percent) + ' %');
      if (m.load !== undefined) parts.push('load ' + m.load);
      $('metrics').textContent = m.error || parts.join(' · ') || wrm.t('No metrics', 'Nema metrika');
    } catch (e) {
      $('metrics').textContent = e.code + ': ' + e.message;
    }
  }

  wrm.on('context', async () => {
    $('conn-label').textContent = wrm.t('Server', 'Poslužitelj');
    $('cmd-label').textContent = wrm.t('Command', 'Naredba');
    $('run').textContent = wrm.t('Run', 'Pokreni');
    try {
      conns = (await wrm.connections()).filter((c) => c.exec);
    } catch (e) {
      $('out').textContent = e.code + ': ' + e.message;
      return;
    }
    const last = await wrm.storeGet('connection', 0);
    $('conn').innerHTML = '';
    for (const c of conns) {
      const o = document.createElement('option');
      o.value = c.id;
      o.textContent = c.name + ' (' + c.host + ')';
      if (c.id === last) o.selected = true;
      $('conn').appendChild(o);
    }
    showMetrics();
  });

  $('conn').addEventListener('change', () => {
    wrm.storeSet('connection', Number($('conn').value));
    showMetrics();
  });

  $('run').addEventListener('click', async () => {
    $('out').textContent = wrm.t('Waiting for your approval in WRM…', 'Čekam vaše odobrenje u WRM-u…');
    try {
      const r = await wrm.exec(Number($('conn').value), $('cmd').value, {reason: 'Command typed in the Server tools window'});
      $('out').textContent = r.output + (r.exit ? '\n[exit ' + r.exit + ']' : '');
    } catch (e) {
      $('out').textContent = e.code + ': ' + e.message;
    }
  });

  wrm.ready();
})();
