// Unit tests of the browser-local serial manager: node --test tools/test-serial.js
const test = require('node:test');
const assert = require('node:assert/strict');
const {createSerialTerminal} = require('../remote-manager/static/serial_terminal.js');
const {recognizeDevice, suggestSettings, validateSettings, profileKey} = require('../remote-manager/static/serial.js');

test('recognizes common boards by USB identity', () => {
  assert.match(recognizeDevice({usbVendorId: 0x2341, usbProductId: 0x0043}).name, /Uno/);
  assert.match(recognizeDevice({usbVendorId: 0x303a, usbProductId: 0x1001}).name, /ESP32/);
  assert.match(recognizeDevice({usbVendorId: 0x0483, usbProductId: 0x5740}).name, /STM32/);
  assert.match(recognizeDevice({usbVendorId: 0x2e8a, usbProductId: 0x000a}).name, /Pico/);
  assert.match(recognizeDevice({usbVendorId: 0x10c4, usbProductId: 0xea60}).name, /USB-UART/);
});
test('unknown USB identity is not fabricated', () => {
  const d = recognizeDevice({usbVendorId: 0x9999, usbProductId: 0x1111});
  assert.equal(d.name, 'Unknown USB serial device');
  assert.equal(d.vendorId, '9999');
});
test('settings are suggested and validated before opening a port', () => {
  assert.equal(suggestSettings('esp32').baudRate, 115200);
  assert.equal(suggestSettings('arduino').baudRate, 9600);
  assert.throws(() => validateSettings({baudRate: 0, dataBits: 8, stopBits: 1, parity: 'none', flowControl: 'none'}));
  assert.throws(() => validateSettings({baudRate: 9600, dataBits: 9, stopBits: 1, parity: 'none', flowControl: 'none'}));
  assert.equal(validateSettings(suggestSettings('pico')).dataBits, 8);
});
test('profile keys preserve full VID/PID pair', () => {
  assert.equal(profileKey({usbVendorId: 0x2341, usbProductId: 0x0043}), '2341:0043');
});
test('terminal reads bytes and sends bytes with a writer lock', async () => {
  let received = 0, released = false;
  const sent = [];
  const reader = {read: async () => ({done: true}), releaseLock: () => { released = true; }};
  const port = {readable: {getReader: () => reader}, writable: {getWriter: () => ({write: async b => sent.push([...b]), releaseLock: () => {}})}, close: async () => {}};
  const term = createSerialTerminal(null, port, {onReceive: n => { received += n; }});
  await term.send(Uint8Array.from([65]));
  await term.disconnect();
  assert.deepEqual(sent, [[65]]);
  assert.equal(released, true);
  assert.equal(received, 0);
  await assert.rejects(term.send(Uint8Array.from([66])));
});
