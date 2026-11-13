/* xterm display and Web Serial stream lifecycle. Serial bytes stay in the browser. */
(function (root) {
  'use strict';
  function createSerialTerminal(host, port, hooks = {}) {
    const term = root && root.Terminal && host ? new root.Terminal({convertEol: true, scrollback: 5000, fontSize: 13}) : null;
    const fitAddon = term && root.FitAddon ? new root.FitAddon.FitAddon() : null;
    if (term) { if (fitAddon) term.loadAddon(fitAddon); term.open(host); fitAddon?.fit(); }
    let stopped = false, reader = null, pending = Promise.resolve();
    const decoder = new TextDecoder();
    const reading = (async () => {
      try {
        while (!stopped && port.readable) {
          const stream = port.readable;
          reader = stream.getReader();
          try {
            while (!stopped) {
              const {value, done} = await reader.read();
              if (done) break;
              if (value) { term?.write(decoder.decode(value, {stream:true})); hooks.onReceive?.(value.byteLength); }
            }
          } finally { reader.releaseLock(); reader = null; }
          if (!port.readable || stopped || port.readable === stream) break;
          // Framing errors can replace readable with a new stream.
          await new Promise(resolve => setTimeout(resolve, 50));
        }
      } catch (e) { if (!stopped) hooks.onError?.(e); }
      finally { if (!stopped) hooks.onDisconnect?.(); }
    })();
    function send(bytes) {
      if (stopped || !port.writable) return Promise.reject(new Error('Serial port is disconnected'));
      const operation = pending.then(async () => {
        if (stopped || !port.writable) throw new Error('Serial port is disconnected');
        const writer = port.writable.getWriter();
        try { await writer.write(bytes); hooks.onSend?.(bytes.byteLength); }
        finally { writer.releaseLock(); }
      });
      pending = operation.catch(() => {});
      return operation;
    }
    async function disconnect() {
      if (stopped) return;
      stopped = true;
      if (reader) { try { await reader.cancel(); } catch (_) {} }
      await reading;
      await pending;
      try { await port.close(); } finally { term?.dispose(); }
    }
    return {send, disconnect, fit: () => fitAddon?.fit()};
  }
  if (typeof module !== 'undefined' && module.exports) module.exports = {createSerialTerminal};
  if (root) root.createSerialTerminal = createSerialTerminal;
})(typeof window !== 'undefined' ? window : null);
