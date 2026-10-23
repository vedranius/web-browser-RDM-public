// Screen lock (#115). After a period without activity the interface is covered
// by a lock screen that asks for the signed-in account's password. Sessions,
// tunnels and AI runs keep going; the lock only detaches the view. Every unlock
// attempt is checked and audited by the server (/api/auth/lock/verify).
'use strict';

Object.assign(LANGS.en, {
  lock_title: 'WRM PRO is locked',
  lock_hint: 'Enter your password to continue. Your sessions keep running in the background.',
  lock_password: 'Password',
  lock_unlock: 'Unlock',
  lock_now: 'Lock now',
  lock_after: 'Lock after inactivity',
  lock_off: 'Off',
  lock_minutes: '{n} min',
  lock_wrong: 'The password is not correct',
  lock_throttled: 'Too many attempts; wait a minute',
  lock_note: 'This is a screen lock, not a new sign-in: the browser session stays valid.',
});

Object.assign(LANGS.hr, {
  lock_title: 'WRM PRO je zaključan',
  lock_hint: 'Upišite lozinku za nastavak. Sesije nastavljaju raditi u pozadini.',
  lock_password: 'Lozinka',
  lock_unlock: 'Otključaj',
  lock_now: 'Zaključaj sada',
  lock_after: 'Zaključaj nakon neaktivnosti',
  lock_off: 'Isključeno',
  lock_minutes: '{n} min',
  lock_wrong: 'Lozinka nije ispravna',
  lock_throttled: 'Previše pokušaja; pričekajte minutu',
  lock_note: 'Ovo je zaključavanje zaslona, ne nova prijava: sesija preglednika ostaje valjana.',
});

(function () {
  const $ = id => document.getElementById(id);
  const KEY = 'wrm_lock_minutes';
  const DEFAULT = 5;
  let last = Date.now(), timer = null, overlay = null, busy = false;

  function minutes() {
    const v = parseInt(localStorage.getItem(KEY), 10);
    return Number.isFinite(v) && v >= 0 && v <= 60 ? v : DEFAULT;
  }

  function touch() { last = Date.now(); }

  function lock() {
    if (overlay || !(typeof currentUser !== 'undefined' && currentUser)) return;
    overlay = document.createElement('div');
    overlay.id = 'lock-screen';
    overlay.style.cssText = 'position:fixed;inset:0;z-index:2147483000;background:rgba(6,8,12,.82);backdrop-filter:blur(6px);display:flex;align-items:center;justify-content:center;padding:24px;';
    overlay.innerHTML = `<div style="width:360px;max-width:100%;background:var(--bg2,#11151c);border:1px solid var(--border,#2a3240);border-radius:12px;padding:22px;color:var(--text,#e5e7eb);box-shadow:0 20px 60px #000a;font:13px/1.45 system-ui,sans-serif;">
      <div style="font-size:24px;margin-bottom:6px;">🔒</div>
      <h3 style="margin:0 0 6px;font-size:16px;">${esc(t('lock_title'))}</h3>
      <div style="color:var(--text3,#94a3b8);font-size:11.5px;margin-bottom:12px;">${esc(t('lock_hint'))}</div>
      <label for="lock-pass" style="font-size:11px;color:var(--text3,#94a3b8);">${esc(t('lock_password'))}</label>
      <input id="lock-pass" type="password" autocomplete="current-password" style="width:100%;margin:4px 0 10px;padding:8px;border-radius:6px;border:1px solid var(--border,#2a3240);background:var(--bg3,#0b0f15);color:var(--text,#e5e7eb);">
      <div id="lock-err" style="color:#fca5a5;font-size:11.5px;min-height:16px;margin-bottom:8px;"></div>
      <button id="lock-go" style="width:100%;padding:8px;border:0;border-radius:6px;background:var(--accent,#4f8cff);color:#fff;cursor:pointer;">${esc(t('lock_unlock'))}</button>
      <div style="display:flex;align-items:center;gap:6px;margin-top:14px;font-size:11px;color:var(--text3,#94a3b8);">
        <span>${esc(t('lock_after'))}</span>
        <select id="lock-min" style="background:var(--bg3,#0b0f15);color:var(--text,#e5e7eb);border:1px solid var(--border,#2a3240);border-radius:5px;font-size:11px;padding:2px 4px;">
          ${[0, 2, 5, 10, 30, 60].map(n => `<option value="${n}" ${minutes() === n ? 'selected' : ''}>${n ? esc(tf('lock_minutes', {n})) : esc(t('lock_off'))}</option>`).join('')}
        </select>
      </div>
      <div style="margin-top:8px;font-size:10.5px;color:var(--text3,#94a3b8);">${esc(t('lock_note'))}</div>
    </div>`;
    document.body.appendChild(overlay);
    const inp = $('lock-pass');
    inp.focus();
    inp.addEventListener('keydown', e => { if (e.key === 'Enter') unlock(); });
    $('lock-go').addEventListener('click', unlock);
    $('lock-min').addEventListener('change', e => {
      localStorage.setItem(KEY, e.target.value);
      touch(); arm();
    });
  }

  async function unlock() {
    if (busy) return;
    const pass = ($('lock-pass') || {}).value || '';
    const err = $('lock-err'), btn = $('lock-go');
    if (!pass) { if (err) err.textContent = t('lock_password'); return; }
    busy = true;
    if (btn) btn.disabled = true;
    try {
      const res = await fetch('/api/auth/lock/verify', {
        method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({password: pass})
      });
      if (!res.ok) throw new Error(res.status === 429 ? t('lock_throttled') : t('lock_wrong'));
      if (overlay) overlay.remove();
      overlay = null;
      touch();
      arm();
    } catch (e) {
      if (err) err.textContent = e.message;
      const inp = $('lock-pass');
      if (inp) { inp.value = ''; inp.focus(); }
    } finally {
      busy = false;
      if (btn) btn.disabled = false;
    }
  }

  function arm() {
    clearInterval(timer);
    timer = setInterval(() => {
      const m = minutes();
      if (!m || overlay) return;
      if (Date.now() - last >= m * 60000) lock();
    }, 15000);
  }

  function button() {
    let b = $('btn-lock');
    if (!b) {
      b = document.createElement('button');
      b.id = 'btn-lock';
      b.className = 'btn-top';
      b.innerHTML = '<span>🔒</span><span class="lbl"></span>';
      b.addEventListener('click', lock);
      if ($('btn-logout')) $('btn-logout').parentNode.insertBefore(b, $('btn-logout'));
    }
    b.title = t('lock_now');
    b.setAttribute('aria-label', t('lock_now'));
    const lbl = b.querySelector('.lbl');
    if (lbl) lbl.textContent = t('lock_now');
  }

  function boot() {
    if (typeof currentUser === 'undefined' || !currentUser) return;
    button();
    ['pointerdown', 'keydown', 'wheel', 'touchstart'].forEach(ev => window.addEventListener(ev, touch, {passive: true}));
    document.addEventListener('visibilitychange', () => { if (!document.hidden) touch(); });
    arm();
  }

  window.wrmLock = {lock, arm, minutes};

  if (typeof window.initApp === 'function' && !window.initApp.__lock) {
    const orig = window.initApp;
    window.initApp = function () {
      const r = orig.apply(this, arguments);
      boot();
      return r;
    };
    window.initApp.__lock = true;
  } else if (typeof currentUser !== 'undefined' && currentUser) {
    boot();
  }
})();
