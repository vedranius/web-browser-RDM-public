// Passkeys (v14.13.0): sign in with Windows Hello, Touch ID / Face ID, a phone's fingerprint or
// face (also a phone next to the computer through the browser's QR code) or a security key.
// A passkey signs in on its own, works as the second step after the password, unlocks the
// screen lock and confirms step-ups. Settings → Security lists, adds, renames and removes them.
'use strict';

Object.assign(LANGS.en, {
  pk_h: 'Passkeys',
  pk_intro: 'Sign in with your fingerprint, face or device PIN instead of a password: Windows Hello, Touch ID, Face ID, an Android or iPhone phone (also a phone next to this computer: the browser shows a QR code) or a security key. A passkey also unlocks the screen lock and confirms sensitive actions. It cannot be phished: it only works on this WRM address.',
  pk_add: 'Add a passkey',
  pk_none: 'No passkeys yet.',
  pk_name_q: 'Name for this passkey (for example "Work laptop" or "My phone"):',
  pk_added: 'Passkey added',
  pk_removed: 'Passkey removed',
  pk_remove_q: 'Remove the passkey "{n}"? You can no longer sign in with it.',
  pk_rename: 'Rename',
  pk_remove: 'Remove',
  pk_synced: 'synced',
  pk_synced_d: 'Saved in your password manager or cloud keychain, available on your other devices.',
  pk_device: 'this device only',
  pk_created: 'Added {d}',
  pk_used: 'last used {d}',
  pk_never: 'never used',
  pk_off_admin: 'Your administrator turned passkeys off.',
  pk_unavailable: 'Not available on this address: {r}',
  pk_browser_na: 'This browser does not support passkeys.',
  pk_cancelled: 'The passkey request was cancelled or timed out.',
  pk_sign_in: 'Sign in with a passkey',
  pk_sign_in_d: 'Fingerprint, face, device PIN or your phone',
  pk_use: 'Use a passkey',
  pk_or_code: 'or enter a code',
  pk_mfa_only: 'Confirm with your passkey to finish signing in.',
  pk_unlock: 'Unlock with a passkey',
  pk_login_mode: 'Sign-in: {m}',
  pk_mode_passwordless: 'a passkey signs in on its own',
  pk_mode_second_factor: 'a passkey is the second step after the password',
  pk_mode_off: 'only unlocking and confirmations (signing in uses the password)',
  pk_mfa_note: 'Your administrator requires 2FA: a passkey counts as your second factor.',
  pk_restricted_add: 'Add a passkey instead',
  // Admin → Policies
  pg_passkeys: 'Passkeys',
  p_passkeys: 'Allow passkeys',
  p_passkeys_d: 'Users may add passkeys (Windows Hello, Touch ID, phone, security key). Browsers allow them only over HTTPS on a host name.',
  p_passkey_login: 'What a passkey may do at sign-in',
  p_passkey_login_d: 'Passwordless: a passkey signs in on its own (it verifies the user with a fingerprint, face or PIN, so it also satisfies required 2FA). Second step: only after the password. Off: unlock and confirmations only.',
  po_passkey_login_passwordless: 'Sign in on its own',
  po_passkey_login_second_factor: 'Second step only',
  po_passkey_login_off: 'Not for sign-in',
  p_passkey_rp_id: 'Passkey domain',
  p_passkey_rp_id_d: 'Empty: the host name WRM is opened with. Set a parent domain (example.com) if WRM is reached under several subdomains. Changing it makes existing passkeys stop working.',
  p_passkey_origins: 'Extra passkey origins',
  p_passkey_origins_d: 'Comma-separated full origins (https://wrm.example.com:8443) accepted besides https pages inside the passkey domain.',
  pk_admin_reset: 'Remove passkeys',
  pk_admin_reset_q: 'Remove all passkeys of {u}? Use it for a lost phone or laptop.',
  pk_admin_reset_done: 'Removed {n} passkeys',
});

Object.assign(LANGS.hr, {
  pk_h: 'Passkeyi (pristupni ključevi)',
  pk_intro: 'Prijava otiskom prsta, licem ili PIN-om uređaja umjesto lozinke: Windows Hello, Touch ID, Face ID, Android ili iPhone mobitel (i mobitel pokraj ovog računala: preglednik pokaže QR kod) ili sigurnosni ključ. Passkey također otključava zaključan zaslon i potvrđuje osjetljive radnje. Ne može se ukrasti phishingom: radi samo na ovoj WRM adresi.',
  pk_add: 'Dodaj passkey',
  pk_none: 'Još nema passkeyeva.',
  pk_name_q: 'Naziv ovog passkeya (na primjer "Poslovni laptop" ili "Moj mobitel"):',
  pk_added: 'Passkey dodan',
  pk_removed: 'Passkey uklonjen',
  pk_remove_q: 'Ukloniti passkey "{n}"? S njim se više nećete moći prijaviti.',
  pk_rename: 'Preimenuj',
  pk_remove: 'Ukloni',
  pk_synced: 'sinkroniziran',
  pk_synced_d: 'Spremljen u upravitelju lozinki ili cloud privjesku, dostupan i na vašim drugim uređajima.',
  pk_device: 'samo na ovom uređaju',
  pk_created: 'Dodan {d}',
  pk_used: 'zadnje korišten {d}',
  pk_never: 'još nije korišten',
  pk_off_admin: 'Administrator je isključio passkeye.',
  pk_unavailable: 'Nije dostupno na ovoj adresi: {r}',
  pk_browser_na: 'Ovaj preglednik ne podržava passkeye.',
  pk_cancelled: 'Zahtjev za passkey je otkazan ili je isteklo vrijeme.',
  pk_sign_in: 'Prijava passkeyem',
  pk_sign_in_d: 'Otisak prsta, lice, PIN uređaja ili mobitel',
  pk_use: 'Koristi passkey',
  pk_or_code: 'ili upišite kod',
  pk_mfa_only: 'Potvrdite passkeyem za završetak prijave.',
  pk_unlock: 'Otključaj passkeyem',
  pk_login_mode: 'Prijava: {m}',
  pk_mode_passwordless: 'passkey prijavljuje sam',
  pk_mode_second_factor: 'passkey je drugi korak nakon lozinke',
  pk_mode_off: 'samo otključavanje i potvrde (prijava ide lozinkom)',
  pk_mfa_note: 'Administrator traži 2FA: passkey se računa kao drugi faktor.',
  pk_restricted_add: 'Umjesto toga dodaj passkey',
  pg_passkeys: 'Passkeyi',
  p_passkeys: 'Dopusti passkeye',
  p_passkeys_d: 'Korisnici smiju dodati passkeye (Windows Hello, Touch ID, mobitel, sigurnosni ključ). Preglednici ih dopuštaju samo preko HTTPS-a na nazivu hosta.',
  p_passkey_login: 'Što passkey smije pri prijavi',
  p_passkey_login_d: 'Bez lozinke: passkey prijavljuje sam (provjerava korisnika otiskom, licem ili PIN-om pa zadovoljava i obavezni 2FA). Drugi korak: samo nakon lozinke. Isključeno: samo otključavanje i potvrde.',
  po_passkey_login_passwordless: 'Prijava bez lozinke',
  po_passkey_login_second_factor: 'Samo drugi korak',
  po_passkey_login_off: 'Ne za prijavu',
  p_passkey_rp_id: 'Domena passkeyeva',
  p_passkey_rp_id_d: 'Prazno: naziv hosta kojim se otvara WRM. Upišite nadređenu domenu (example.com) ako se WRM otvara preko više poddomena. Promjena čini postojeće passkeye neupotrebljivima.',
  p_passkey_origins: 'Dodatna ishodišta passkeyeva',
  p_passkey_origins_d: 'Puna ishodišta odvojena zarezom (https://wrm.example.com:8443) koja se prihvaćaju uz https stranice unutar domene passkeyeva.',
  pk_admin_reset: 'Ukloni passkeye',
  pk_admin_reset_q: 'Ukloniti sve passkeye korisnika {u}? Za izgubljen mobitel ili laptop.',
  pk_admin_reset_done: 'Uklonjeno passkeyeva: {n}',
});

(function () {
  const $ = id => document.getElementById(id);
  const tf = (k, vars) => { let s = t(k); for (const [a, b] of Object.entries(vars || {})) s = s.split('{' + a + '}').join(String(b)); return s; };
  const supported = () => !!(window.PublicKeyCredential && navigator.credentials && window.isSecureContext);

  // ── base64url ↔ ArrayBuffer ──
  function fromB64(s) {
    s = String(s).replace(/-/g, '+').replace(/_/g, '/');
    while (s.length % 4) s += '=';
    const bin = atob(s), out = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
    return out.buffer;
  }
  function toB64(buf) {
    const b = new Uint8Array(buf);
    let s = '';
    for (let i = 0; i < b.length; i++) s += String.fromCharCode(b[i]);
    return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
  }

  async function post(url, body) {
    const res = await fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body || {})});
    const d = await res.json().catch(() => ({}));
    if (!res.ok) { const e = new Error(d.error || ('HTTP ' + res.status)); e.status = res.status; e.data = d; throw e; }
    return d;
  }

  function cancelled(e) { return e && (e.name === 'NotAllowedError' || e.name === 'AbortError'); }
  function message(e) { return cancelled(e) ? t('pk_cancelled') : (e && e.message) || t('failed'); }

  // get(): runs navigator.credentials.get with the server's options and returns the JSON answer.
  async function get(publicKey) {
    const o = Object.assign({}, publicKey, {challenge: fromB64(publicKey.challenge)});
    if (o.allowCredentials) o.allowCredentials = o.allowCredentials.map(c => Object.assign({}, c, {id: fromB64(c.id)}));
    const cred = await navigator.credentials.get({publicKey: o});
    const r = cred.response;
    return {id: cred.id, type: cred.type, response: {
      clientDataJSON: toB64(r.clientDataJSON), authenticatorData: toB64(r.authenticatorData),
      signature: toB64(r.signature), userHandle: r.userHandle ? toB64(r.userHandle) : ''}};
  }

  async function create(publicKey) {
    const o = Object.assign({}, publicKey, {challenge: fromB64(publicKey.challenge),
      user: Object.assign({}, publicKey.user, {id: new TextEncoder().encode(publicKey.user.id)})});
    o.excludeCredentials = (o.excludeCredentials || []).map(c => Object.assign({}, c, {id: fromB64(c.id)}));
    const cred = await navigator.credentials.create({publicKey: o});
    const r = cred.response;
    return {id: cred.id, response: {clientDataJSON: toB64(r.clientDataJSON), attestationObject: toB64(r.attestationObject),
      transports: typeof r.getTransports === 'function' ? r.getTransports() : []}};
  }

  // proof(): a passkey answer for the signed-in session (purpose unlock or step_up).
  async function proof(purpose) {
    const b = await post('/api/auth/webauthn/begin', {purpose});
    return {ceremony: b.ceremony, credential: await get(b.publicKey)};
  }

  // ── sign-in screen ──
  async function signIn(mfa) {
    const err = mfa ? 'auth-mfa-error' : 'auth-error';
    try {
      const body = mfa ? {mfa_token: mfa} : {};
      const b = await post('/api/auth/webauthn/login/begin', body);
      const credential = await get(b.publicKey);
      const me = await post('/api/auth/webauthn/login/finish', Object.assign({ceremony: b.ceremony, credential}, body));
      if (typeof onSignedIn === 'function') onSignedIn(me);
    } catch (e) {
      if (typeof authErr === 'function') authErr(err, message(e));
    }
  }

  function signInButton() {
    const host = $('auth-step-cred');
    if (!host) return;
    let box = $('auth-passkey');
    const cfg = typeof authConfig !== 'undefined' ? authConfig : {};
    const show = supported() && cfg.passkey_login === 'passwordless' && !cfg.first_run && (typeof authMode === 'undefined' || authMode === 'login');
    if (!show) { if (box) box.remove(); return; }
    if (!box) {
      box = document.createElement('div');
      box.id = 'auth-passkey';
      const sso = $('auth-sso');
      if (sso && sso.parentNode === host) host.insertBefore(box, sso); else host.insertBefore(box, host.firstChild);
    }
    box.innerHTML = `<button type="button" class="auth-btn auth-sso-btn" id="auth-passkey-btn" title="${esc(t('pk_sign_in_d'))}">🔑 ${esc(t('pk_sign_in'))}</button><div class="auth-or">${esc(t('sso_or'))}</div>`;
    $('auth-passkey-btn').onclick = () => signIn(null);
  }

  // The second step after the password: the passkey next to (or instead of) the code field.
  let mfaMethods = null;
  function mfaStep(token, methods) {
    mfaMethods = methods || ['totp'];
    const step = $('auth-step-mfa');
    if (!step) return;
    let box = $('auth-mfa-passkey');
    const hasPk = mfaMethods.includes('passkey') && supported(), hasCode = mfaMethods.includes('totp');
    if (!box) {
      box = document.createElement('div');
      box.id = 'auth-mfa-passkey';
      step.insertBefore(box, step.querySelector('.auth-field'));
    }
    box.innerHTML = hasPk ? `<button type="button" class="auth-btn" id="auth-mfa-pk">🔑 ${esc(t('pk_use'))}</button>${hasCode ? `<div class="auth-or">${esc(t('pk_or_code'))}</div>` : ''}` : '';
    if (hasPk) $('auth-mfa-pk').onclick = () => signIn(token);
    const note = step.querySelector('.auth-note'), field = step.querySelector('.auth-field'), btn = step.querySelector('.auth-field + .auth-btn');
    if (note) note.textContent = hasCode ? t('mfa_prompt') : t('pk_mfa_only');
    if (field) field.style.display = hasCode ? '' : 'none';
    if (btn) btn.style.display = hasCode ? '' : 'none';
    if (hasPk && !hasCode) signIn(token);
  }

  // ── step-up: "Use a passkey" in the confirmation dialog ──
  function stepUpButton(info) {
    const box = $('su-sso');
    if (!box || !info || !info.passkey || !supported()) return;
    const div = document.createElement('div');
    div.innerHTML = `<button class="btn-sec" style="width:100%;margin-bottom:8px;" id="su-passkey">🔑 ${esc(t('pk_use'))}</button>` +
      (info.password !== false ? `<div class="auth-or">${esc(t('sso_or'))}</div>` : '');
    box.insertBefore(div, box.firstChild);
    $('su-passkey').onclick = async () => {
      const e = $('su-error');
      try {
        const p = await proof('step_up');
        const res = await fetch('/api/auth/step-up', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({passkey: p, action: info.action || ''})});
        if (!res.ok) throw new Error(await apiError(res));
        if (typeof stepUpDone === 'function') stepUpDone(true);
      } catch (x) { if (e) { e.textContent = message(x); e.style.display = ''; } }
    };
  }

  // ── Settings → Security ──
  async function renderSettings(box) {
    if (!box) return;
    const res = await fetch('/api/me/passkeys');
    if (!res.ok) { box.innerHTML = ''; return; }
    const d = await res.json();
    const list = d.passkeys || [];
    let status = '';
    if (!d.enabled) status = t('pk_off_admin');
    else if (!supported()) status = t('pk_browser_na');
    else if (!d.available) status = tf('pk_unavailable', {r: d.reason || ''});
    const me = (typeof currentUser !== 'undefined' && currentUser) || {};
    box.innerHTML = `<div class="adm-h">${esc(t('pk_h'))}</div>
      <div class="setting-desc" style="margin-bottom:8px;">${esc(t('pk_intro'))}</div>
      ${d.enabled ? `<div class="setting-desc" style="margin-bottom:8px;">${esc(tf('pk_login_mode', {m: t('pk_mode_' + d.login)}))}${me.mfa_required ? ' ' + esc(t('pk_mfa_note')) : ''}</div>` : ''}
      ${status ? `<div class="hint-box" style="margin-bottom:8px;">${esc(status)}</div>` : ''}
      <div id="pk-list">${list.length ? list.map(p => `<div class="setting-row" data-pk="${Number(p.id)}">
        <div style="min-width:0;"><div class="setting-label">🔑 ${esc(p.name)} <span class="pill ${p.synced ? 'ok' : ''}" title="${esc(p.synced ? t('pk_synced_d') : '')}">${esc(p.synced ? t('pk_synced') : t('pk_device'))}</span></div>
          <div class="setting-desc">${esc(tf('pk_created', {d: fmt(p.created_at)}))} · ${esc(p.last_used_at ? tf('pk_used', {d: fmt(p.last_used_at)}) : t('pk_never'))}</div></div>
        <div class="setting-control" style="display:flex;gap:6px;flex-wrap:wrap;justify-content:flex-end;">
          <button class="btn-sec btn-sm" data-act="rename">${esc(t('pk_rename'))}</button>
          <button class="btn-danger btn-sm" data-act="remove">${esc(t('pk_remove'))}</button></div></div>`).join('') : `<div class="setting-desc">${esc(t('pk_none'))}</div>`}</div>
      ${d.enabled ? `<button class="btn-sm" id="pk-add" style="margin-top:8px;" ${d.available && supported() ? '' : 'disabled'}>＋ ${esc(t('pk_add'))}</button>` : ''}`;
    const add = $('pk-add');
    if (add) add.onclick = () => addPasskey().then(ok => { if (ok) renderSettings(box); });
    box.querySelectorAll('[data-pk]').forEach(row => {
      const id = row.dataset.pk, p = list.find(x => String(x.id) === id);
      row.querySelector('[data-act="rename"]').onclick = async () => {
        const name = await askText(t('pk_name_q'), p.name);
        if (!name) return;
        const r = await fetch('/api/me/passkeys/' + id, {method: 'PATCH', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({name})});
        if (!r.ok) { showToast(await apiError(r), 'error', 6000); return; }
        renderSettings(box);
      };
      row.querySelector('[data-act="remove"]').onclick = async () => {
        const ok = typeof uiConfirm === 'function' ? await uiConfirm(tf('pk_remove_q', {n: p.name}), {okText: t('pk_remove'), danger: true}) : confirm(tf('pk_remove_q', {n: p.name}));
        if (!ok) return;
        const r = await withStepUp(() => fetch('/api/me/passkeys/' + id, {method: 'DELETE'}));
        if (!r.ok) { showToast(await apiError(r), 'error', 6000); return; }
        showToast(t('pk_removed'), 'success', 2000);
        renderSettings(box);
      };
    });
  }

  function fmt(ts) { return typeof fmtTime === 'function' ? fmtTime(ts) : (ts || ''); }

  async function askText(q, value) {
    if (typeof uiPrompt === 'function') return uiPrompt(q, value || '');
    return prompt(q, value || '');
  }

  // withStepUp(): runs a request; when the server asks for a confirmation, shows it and retries once.
  async function withStepUp(run) {
    let res = await run();
    if (res.status === 403) {
      const d = await res.clone().json().catch(() => ({}));
      if (d.step_up && typeof window.wrmStepUp === 'function' && await window.wrmStepUp(d)) res = await run();
    }
    return res;
  }

  async function addPasskey() {
    if (!supported()) { showToast(t('pk_browser_na'), 'error', 6000); return false; }
    try {
      const begin = await withStepUp(() => fetch('/api/me/passkeys/begin', {method: 'POST'}));
      if (!begin.ok) throw new Error(await apiError(begin));
      const b = await begin.json();
      const credential = await create(b.publicKey);
      const name = (await askText(t('pk_name_q'), '')) || '';
      await post('/api/me/passkeys/finish', {ceremony: b.ceremony, name, credential});
      showToast(t('pk_added'), 'success', 2500);
      return true;
    } catch (e) {
      showToast(message(e), 'error', 7000);
      return false;
    }
  }

  // Admin → Users: remove every passkey of a user (lost device).
  async function adminReset(userId, username) {
    const q = tf('pk_admin_reset_q', {u: username});
    const ok = typeof uiConfirm === 'function' ? await uiConfirm(q, {okText: t('pk_admin_reset'), danger: true}) : confirm(q);
    if (!ok) return;
    const r = await withStepUp(() => fetch('/api/admin/passkeys/' + Number(userId), {method: 'DELETE'}));
    if (!r.ok) { showToast(await apiError(r), 'error', 6000); return; }
    const d = await r.json().catch(() => ({}));
    showToast(tf('pk_admin_reset_done', {n: d.removed || 0}), 'success', 3000);
  }

  // Hooks into the page's own functions (index.html): the sign-in buttons, the second step and
  // the step-up dialog.
  function hook() {
    if (typeof window.renderAuthSSO === 'function' && !window.renderAuthSSO._pk) {
      const orig = window.renderAuthSSO;
      window.renderAuthSSO = function () { const r = orig.apply(this, arguments); try { signInButton(); } catch (_) { /* sign-in still works */ } return r; };
      window.renderAuthSSO._pk = true;
    }
    if (typeof window.showMFAStep === 'function' && !window.showMFAStep._pk) {
      const orig = window.showMFAStep;
      window.showMFAStep = function (token, methods) { const r = orig.apply(this, arguments); try { mfaStep(token, methods); } catch (_) { /* code field stays */ } return r; };
      window.showMFAStep._pk = true;
    }
    if (typeof window.wrmStepUp === 'function' && !window.wrmStepUp._pk) {
      const orig = window.wrmStepUp;
      window.wrmStepUp = function (info) {
        const p = orig.apply(this, arguments);
        try {
          stepUpButton(info);
          if (info && info.passkey && info.password === false) { const s = $('su-submit'); if (s) s.style.display = 'none'; }
        } catch (_) { /* the password path stays */ }
        return p;
      };
      window.wrmStepUp._pk = true;
    }
    try { signInButton(); } catch (_) { /* not on the sign-in screen */ }
  }

  window.wrmPasskeys = {supported, get, create, proof, renderSettings, addPasskey, adminReset, signIn, hook};
  hook();
})();
