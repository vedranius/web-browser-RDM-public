const test = require('node:test');
const assert = require('node:assert/strict');
const {createSerialTerminal} = require('./serial_terminal.js');
test('reads bytes and sends bytes with a writer lock', async () => {
  let received = 0, sent = [], released = false;
  const reader = {read: async () => ({done:true}), releaseLock: () => {released=true;}};
  const port = {readable:{getReader:()=>reader}, writable:{getWriter:()=>({write:async b=>sent.push([...b]),releaseLock:()=>{}})},close:async()=>{}};
  const term = createSerialTerminal(null, port, {onReceive:n=>received+=n});
  await term.send(Uint8Array.from([65]));
  await term.disconnect();
  assert.deepEqual(sent, [[65]]);
  assert.equal(released,true);
  assert.equal(received,0);
});
