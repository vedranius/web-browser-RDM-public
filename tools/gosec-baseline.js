#!/usr/bin/env node
// Compares a gosec JSON report with the accepted baseline and fails on new findings.
//
//   gosec -quiet -fmt json -out gosec.json -exclude=G104 ./...      (in remote-manager/)
//   node tools/gosec-baseline.js remote-manager/gosec.json           check (CI)
//   node tools/gosec-baseline.js remote-manager/gosec.json --update  accept the current findings
//
// A finding is identified by rule, file and the text of the flagged line (not its line
// number), so unrelated edits that move code do not count as new findings. The baseline
// lists findings that were reviewed and judged not exploitable; a new one must be fixed,
// or reviewed and then added with --update in the same pull request, with a reason in the
// pull request description.
'use strict';
const fs = require('fs');
const path = require('path');

const baselinePath = path.join(__dirname, '..', '.github', 'security', 'gosec-baseline.json');

function fingerprint(issue) {
  const file = issue.file.replace(/\\/g, '/').replace(/^.*?\/remote-manager\//, '');
  const first = String(issue.line).split('-')[0];
  let text = '';
  for (const l of String(issue.code || '').split('\n')) {
    const m = l.match(/^(\d+):\s?(.*)$/);
    if (m && m[1] === first) { text = m[2]; break; }
  }
  return `${issue.rule_id} ${file} ${text.replace(/\s+/g, ' ').trim()}`;
}

function counts(list) {
  const out = {};
  for (const f of list) out[f] = (out[f] || 0) + 1;
  return out;
}

const [report, flag] = process.argv.slice(2);
if (!report) {
  console.error('usage: gosec-baseline.js <gosec.json> [--update]');
  process.exit(2);
}
const data = JSON.parse(fs.readFileSync(report, 'utf8'));
// gosec keeps going when it cannot type-check a file (for example a Go release newer than
// the one it was built for) and then silently misses most taint findings: refuse that.
const loadErrors = Object.keys(data['Golang errors'] || {});
if (loadErrors.length) {
  console.log(`::error::gosec could not load ${loadErrors.length} files, so its results are incomplete (first: ${loadErrors[0]}). Run it with a Go version it supports.`);
  process.exit(1);
}
const issues = data.Issues || [];
const current = issues.map(fingerprint);

if (flag === '--update') {
  const sorted = current.slice().sort();
  fs.writeFileSync(baselinePath, JSON.stringify({ note: 'Reviewed gosec findings (rule, file, flagged line). See tools/gosec-baseline.js and SECURITY.md.', findings: sorted }, null, 1) + '\n');
  console.log(`Baseline updated: ${sorted.length} findings.`);
  process.exit(0);
}

const base = counts(JSON.parse(fs.readFileSync(baselinePath, 'utf8')).findings);
const seen = {};
const fresh = [];
issues.forEach((issue, i) => {
  const f = current[i];
  seen[f] = (seen[f] || 0) + 1;
  if (seen[f] > (base[f] || 0)) fresh.push(issue);
});
const gone = Object.keys(base).filter(f => !seen[f]).length;
console.log(`gosec: ${issues.length} findings, ${issues.length - fresh.length} in the baseline, ${fresh.length} new; ${gone} baseline entries no longer found.`);
for (const i of fresh) {
  console.log(`::error file=remote-manager/${i.file.replace(/\\/g, '/').replace(/^.*?\/remote-manager\//, '')},line=${String(i.line).split('-')[0]}::gosec ${i.rule_id} (${i.severity}/${i.confidence}): ${i.details}`);
}
if (fresh.length) {
  console.log('Fix the new findings, or review them and run: node tools/gosec-baseline.js remote-manager/gosec.json --update');
  process.exit(1);
}
