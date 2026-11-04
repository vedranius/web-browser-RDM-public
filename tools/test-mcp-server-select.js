#!/usr/bin/env node
'use strict';

const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const {test} = require('node:test');

const source = fs.readFileSync(path.join(__dirname, '..', 'remote-manager', 'static', 'mcp.js'), 'utf8');
const start = source.indexOf('  function selectAll(prefix) {');
const end = source.indexOf('\n  function modeChanged(', start);
if (start < 0 || end < 0) throw new Error('MCP server selection helper not found');
const helper = source.slice(start, end);

test('select all checks every server in the current grant form and individual items stay editable', () => {
  const servers = [{value: '4', checked: false}, {value: '9', checked: false}, {value: '12', checked: false}];
  const selectors = [];
  const context = vm.createContext({document: {querySelectorAll(selector) { selectors.push(selector); return servers; }}});
  vm.runInContext(helper, context);
  context.selectAll('mcpc');
  assert.deepEqual(selectors, ['.mcpc-conn']);
  assert.deepEqual(servers.map(x => x.checked), [true, true, true]);
  servers[1].checked = false;
  assert.deepEqual(servers.map(x => x.checked), [true, false, true]);
});

test('server grant and setup-wizard forms both expose select all', () => {
  assert.match(source, /onclick="wrmMCP\.selectAll\('\$\{prefix\}'\)"/);
  assert.match(source, /onclick="wrmMCP\.selectAll\('mcpw'\)"/);
  assert.match(source, /mcp_select_all: 'Select all'/);
  assert.match(source, /mcp_select_all: 'Odaberi sve'/);
});
