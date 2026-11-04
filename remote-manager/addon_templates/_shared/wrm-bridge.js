// WRM addon bridge helper for apiVersion 2 panels. Copy it into your addon and
// change it as you like. The panel has no network and no browser storage of
// its own: everything goes through these messages to the WRM page around it.
(function () {
  'use strict';
  const pending = new Map();
  const listeners = {context: [], theme: [], visibility: []};
  let seq = 0;
  let context = null;
  let bridgeToken = '';
  let bridgePort = null;
  let bridgeAccepted = false;
  let bridgeInit = null;
  let bridgeTimer = null;
  try { bridgeToken = new URLSearchParams(location.hash.slice(1)).get('wrm-bridge') || ''; } catch (_) {}

  function send(msg) {
    const payload = Object.assign({v: 2}, msg);
    if (bridgePort) {
      bridgePort.postMessage(payload);
      return;
    }
    if (msg.type === 'wrm:ready' && bridgeToken && typeof MessageChannel === 'function') {
      const channel = new MessageChannel();
      bridgePort = channel.port2;
      bridgeAccepted = false;
      bridgeInit = payload;
      bridgePort.onmessage = (e) => {
        const d = e.data || {};
        if (d.type === 'wrm:bridge:ready') {
          bridgeAccepted = true;
          bridgeInit = null;
          clearTimeout(bridgeTimer);
          bridgeTimer = null;
          return;
        }
        receive(d);
      };
      bridgePort.start();
      window.parent.postMessage(Object.assign({}, payload, {bridge_token: bridgeToken}), '*', [channel.port1]);
      bridgeToken = '';
      bridgeTimer = setTimeout(fallbackToLegacyBridge, 1000);
      return;
    }
    window.parent.postMessage(payload, '*');
  }

  function fallbackToLegacyBridge() {
    clearTimeout(bridgeTimer);
    bridgeTimer = null;
    if (bridgePort) { try { bridgePort.close(); } catch (_) {} }
    bridgePort = null;
    bridgeAccepted = false;
    if (bridgeInit) window.parent.postMessage(bridgeInit, '*');
    bridgeInit = null;
  }

  function request(type, fields) {
    const id = type + '-' + (++seq);
    return new Promise((resolve, reject) => {
      pending.set(id, {resolve, reject});
      send(Object.assign({type, id}, fields));
      setTimeout(() => {
        if (pending.delete(id)) reject(Object.assign(new Error('WRM did not answer'), {code: 'E_TIMEOUT'}));
      }, 15000);
    });
  }

  function applyTheme(theme) {
    if (!theme) return;
    for (const k of Object.keys(theme)) {
      if (/^--[a-z0-9-]+$/.test(k)) document.documentElement.style.setProperty(k, String(theme[k]));
    }
  }

  function receive(d) {
    if (d.id !== undefined && pending.has(d.id)) {
      const p = pending.get(d.id);
      pending.delete(d.id);
      if (d.error || d.ok === false) p.reject(Object.assign(new Error(d.error || 'Refused'), {code: d.code || ''}));
      else p.resolve(d);
      return;
    }
    if (d.type === 'wrm:context') {
      context = d;
      applyTheme(d.theme);
      listeners.context.forEach((f) => f(d));
    } else if (d.type === 'wrm:theme') {
      applyTheme(d.theme);
      listeners.theme.forEach((f) => f(d.theme));
    } else if (d.type === 'wrm:visibility') {
      listeners.visibility.forEach((f) => f(!!d.visible));
    }
  }

  window.addEventListener('message', (e) => {
    if (e.source !== window.parent) return;
    // A legacy WRM responds over window.postMessage. A current WRM acknowledges
    // the transferred port and sends all later data over that private channel.
    if (bridgePort && !bridgeAccepted) fallbackToLegacyBridge();
    receive(e.data || {});
  });

  window.addEventListener('error', (e) => {
    send({type: 'wrm:error', message: String(e.message || 'error'), source: String(e.filename || ''), line: e.lineno || 0});
  });
  window.addEventListener('unhandledrejection', (e) => {
    const r = e.reason || {};
    send({type: 'wrm:error', message: 'Unhandled: ' + String(r.message || r)});
  });

  window.wrm = {
    // ready tells WRM the panel started and keeps the watchdog quiet.
    ready() {
      send({type: 'wrm:ready'});
      setInterval(() => send({type: 'wrm:heartbeat'}), 15000);
    },
    on(event, fn) { (listeners[event] || []).push(fn); },
    get context() { return context; },
    get visible() { return !context || context.visible !== false; },
    t(en, hr) { return context && context.lang === 'hr' && hr ? hr : en; },
    // Settings of this instance, kept on the WRM server (store capability).
    storeGet(key, fallback) {
      return request('wrm:store.get', {key}).then((d) => (d.value === null || d.value === undefined ? fallback : d.value), () => fallback);
    },
    storeSet(key, value) { return request('wrm:store.set', {key, value}); },
    // A short notification (notify capability).
    notify(text, level) { send({type: 'wrm:notify', text: String(text).slice(0, 200), level: level || 'info'}); },
    // A request WRM sends to a site listed in http.origins (http capability).
    // secret names an entry of http.secrets; the panel never sees its value.
    http(req) { return request('wrm:http', req); },
    error(message) { send({type: 'wrm:error', message: String(message)}); },
    resize(w, h) { send({type: 'wrm:resize', w, h}); }
  };
})();
