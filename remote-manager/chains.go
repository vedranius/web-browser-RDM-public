package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─── CHAINS: LIVE SERVICE-CHAIN MONITORING ────────────────────────────────────────────────
//
// A chain is the set of servers (and storages, and outside endpoints) that form one working
// system, for example ingest → transcoding → content → CDN with a primary and a backup path.
// It has nodes, checks per node from a fixed catalog (chains_checks.go), links between nodes
// and scenarios: ordered rules over the check results, the first that matches being the
// active one ("primary", "ingest on backup", …). WRM checks a chain in the background with the
// SSH connections it already has and fixed read-only checks; nothing is installed on a server.
//
// A change is reported only after it was seen chain.confirm times in a row, so a flapping
// value does not notify on every poll. Changes are kept as samples (the timeline), the latest
// state for display, and the confirmed state for the change detection.

const (
	chainMaxNodes     = 30
	chainMaxChecks    = 200
	chainMaxLinks     = 100
	chainMaxScenarios = 20
	chainMaxCondDepth = 4
	chainMaxCondItems = 20
	chainMaxDefBytes  = 64 << 10
	chainSamplesKept  = 2000
)

var chainIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,39}$`)

var chainLinkKinds = map[string]bool{"data": true, "read": true, "write": true, "failover": true, "dependency": true}

type chainDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	IntervalS   int             `json:"interval_s,omitempty"`
	Confirm     int             `json:"confirm,omitempty"`
	Nodes       []chainNode     `json:"nodes"`
	Checks      []chainCheck    `json:"checks"`
	Links       []chainLink     `json:"links,omitempty"`
	Scenarios   []chainScenario `json:"scenarios,omitempty"`
}

// chainNode is a server (a saved SSH connection), a storage (from Storages) or an outside
// endpoint that only WRM-side checks (DNS, HTTP, TCP, TLS) look at.
type chainNode struct {
	ID        string  `json:"id"`
	Kind      string  `json:"kind"`
	Label     string  `json:"label,omitempty"`
	Role      string  `json:"role,omitempty"`
	ConnID    int     `json:"conn_id,omitempty"`
	StorageID int     `json:"storage_id,omitempty"`
	X         float64 `json:"x,omitempty"`
	Y         float64 `json:"y,omitempty"`
}

type chainCheck struct {
	ID       string     `json:"id"`
	Node     string     `json:"node"`
	Type     string     `json:"type"`
	Label    string     `json:"label,omitempty"`
	Target   string     `json:"target,omitempty"`
	Port     int        `json:"port,omitempty"`
	Format   string     `json:"format,omitempty"`
	Key      string     `json:"key,omitempty"`
	Record   string     `json:"record,omitempty"`
	Resolver string     `json:"resolver,omitempty"`
	Prefix   string     `json:"prefix,omitempty"`
	MaxAgeS  int        `json:"max_age_s,omitempty"`
	WarnDays int        `json:"warn_days,omitempty"`
	Expect   *chainCond `json:"expect,omitempty"`
}

type chainLink struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Kind  string `json:"kind,omitempty"`
	Label string `json:"label,omitempty"`
	Check string `json:"check,omitempty"` // optional: the check whose result colours the link
	// Scenarios, when set, are the scenarios in which the link carries traffic; in any other
	// scenario the diagram shows it idle (for example the backup path while on primary).
	Scenarios []string `json:"scenarios,omitempty"`
}

type chainScenario struct {
	ID    string    `json:"id"`
	Label string    `json:"label,omitempty"`
	When  chainCond `json:"when"`
}

// chainCond is a rule: all / any of nested rules, or one comparison on a check (its value or
// ok) or on a server node (its state). The operators are fixed; there is no expression
// language and no eval.
type chainCond struct {
	All    []chainCond `json:"all,omitempty"`
	Any    []chainCond `json:"any,omitempty"`
	Check  string      `json:"check,omitempty"`
	Node   string      `json:"node,omitempty"`
	Field  string      `json:"field,omitempty"`
	Op     string      `json:"op,omitempty"`
	Value  string      `json:"value,omitempty"`
	Values []string    `json:"values,omitempty"`
}

var chainOps = map[string]bool{"==": true, "!=": true, "in": true, "contains": true, "exists": true}

func (c chainCond) empty() bool {
	return len(c.All) == 0 && len(c.Any) == 0 && c.Check == "" && c.Node == ""
}

// matchValue applies the comparison to one value (an expect on a check).
func (c chainCond) matchValue(v string) bool {
	switch c.Op {
	case "", "==":
		return v == c.Value
	case "!=":
		return v != c.Value
	case "in":
		for _, x := range c.Values {
			if v == x {
				return true
			}
		}
		return false
	case "contains":
		return strings.Contains(v, c.Value)
	case "exists":
		return v != "" && v != "missing" && v != "unknown" && v != "none"
	}
	return false
}

func (c chainCond) eval(s chainState) bool {
	switch {
	case len(c.All) > 0:
		for _, x := range c.All {
			if !x.eval(s) {
				return false
			}
		}
		return true
	case len(c.Any) > 0:
		for _, x := range c.Any {
			if x.eval(s) {
				return true
			}
		}
		return false
	case c.Check != "":
		r, ok := s.Checks[c.Check]
		if !ok {
			return false
		}
		if c.Field == "ok" {
			return c.matchValue(strconv.FormatBool(r.OK))
		}
		return c.matchValue(r.Value)
	case c.Node != "":
		n, ok := s.Nodes[c.Node]
		if !ok {
			return false
		}
		return c.matchValue(n.State)
	}
	return true // an empty rule matches: the fallback scenario
}

func validateChainCond(c *chainCond, depth int, checks, nodes map[string]bool, path string, own bool) []string {
	var errs []string
	if depth > chainMaxCondDepth {
		return []string{path + ": rules are nested too deeply"}
	}
	if len(c.All) > chainMaxCondItems || len(c.Any) > chainMaxCondItems {
		return []string{path + ": too many rules in one list"}
	}
	parts := 0
	for _, x := range []bool{len(c.All) > 0, len(c.Any) > 0, c.Check != "", c.Node != ""} {
		if x {
			parts++
		}
	}
	if parts > 1 {
		errs = append(errs, path+": use one of all, any, check or node per rule")
	}
	for i := range c.All {
		errs = append(errs, validateChainCond(&c.All[i], depth+1, checks, nodes, fmt.Sprintf("%s.all[%d]", path, i), own)...)
	}
	for i := range c.Any {
		errs = append(errs, validateChainCond(&c.Any[i], depth+1, checks, nodes, fmt.Sprintf("%s.any[%d]", path, i), own)...)
	}
	if c.Check != "" && !checks[c.Check] {
		errs = append(errs, fmt.Sprintf("%s: unknown check %q", path, c.Check))
	}
	if c.Node != "" && !nodes[c.Node] {
		errs = append(errs, fmt.Sprintf("%s: unknown server node %q", path, c.Node))
	}
	if c.Check != "" || c.Node != "" || own {
		if !chainOps[c.Op] && c.Op != "" {
			errs = append(errs, fmt.Sprintf("%s: op must be ==, !=, in, contains or exists", path))
		}
		if c.Field != "" && c.Field != "value" && c.Field != "ok" {
			errs = append(errs, path+": field must be value or ok")
		}
		if len(c.Value) > 300 || len(c.Values) > 50 {
			errs = append(errs, path+": value is too long")
		}
	}
	return errs
}

// validateChain normalises a definition and lists every problem, with a path the editor can
// show. The owner must own every connection and storage the chain uses.
func validateChain(d *chainDef, userID int) []string {
	var errs []string
	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" || len(d.Name) > 120 || chainHasControl(d.Name) {
		errs = append(errs, "name: 1-120 characters")
	}
	if len(d.Description) > 2000 {
		errs = append(errs, "description: at most 2000 characters")
	}
	if d.IntervalS != 0 && (d.IntervalS < 15 || d.IntervalS > 3600) {
		errs = append(errs, "interval_s: 15-3600 seconds")
	}
	if d.Confirm != 0 && (d.Confirm < 1 || d.Confirm > 10) {
		errs = append(errs, "confirm: 1-10")
	}
	if len(d.Nodes) == 0 || len(d.Nodes) > chainMaxNodes {
		errs = append(errs, fmt.Sprintf("nodes: 1-%d nodes", chainMaxNodes))
	}
	if len(d.Checks) > chainMaxChecks || len(d.Links) > chainMaxLinks || len(d.Scenarios) > chainMaxScenarios {
		errs = append(errs, fmt.Sprintf("at most %d checks, %d links and %d scenarios", chainMaxChecks, chainMaxLinks, chainMaxScenarios))
		return errs
	}
	nodes := map[string]chainNode{}
	servers := map[string]bool{}
	for i := range d.Nodes {
		n := &d.Nodes[i]
		p := fmt.Sprintf("nodes[%d]", i)
		if !chainIDPattern.MatchString(n.ID) {
			errs = append(errs, p+": id must be 1-40 characters a-z, 0-9, '-' or '_'")
			continue
		}
		p += " (" + n.ID + ")"
		if _, dup := nodes[n.ID]; dup {
			errs = append(errs, p+": the id is used twice")
		}
		if len(n.Label) > 120 || len(n.Role) > 60 || chainHasControl(n.Label+n.Role) {
			errs = append(errs, p+": label at most 120 and role at most 60 characters")
		}
		if !validChainPos(n.X) || !validChainPos(n.Y) {
			errs = append(errs, p+": x and y must be between 0 and 20000")
		}
		switch n.Kind {
		case "server":
			c, err := loadConnectionRaw(n.ConnID)
			if err != nil || c.UserID != userID {
				errs = append(errs, p+": conn_id is not one of your connections")
			} else if !isSSHProtocol(c) {
				errs = append(errs, p+": the connection must use SSH")
			}
			n.StorageID = 0
			servers[n.ID] = true
		case "storage":
			s, err := loadStorage(n.StorageID, userID)
			if err != nil || s.UserID != userID {
				errs = append(errs, p+": storage_id is not one of your storages")
			}
			n.ConnID = 0
		case "external":
			n.ConnID, n.StorageID = 0, 0
		default:
			errs = append(errs, p+": kind must be server, storage or external")
		}
		nodes[n.ID] = *n
	}
	checks := map[string]bool{}
	for i := range d.Checks {
		c := &d.Checks[i]
		p := fmt.Sprintf("checks[%d]", i)
		if !chainIDPattern.MatchString(c.ID) {
			errs = append(errs, p+": id must be 1-40 characters a-z, 0-9, '-' or '_'")
			continue
		}
		p += " (" + c.ID + ")"
		if checks[c.ID] {
			errs = append(errs, p+": the id is used twice")
		}
		checks[c.ID] = true
		if len(c.Label) > 120 || chainHasControl(c.Label) {
			errs = append(errs, p+": label at most 120 characters")
		}
		n, ok := nodes[c.Node]
		if !ok {
			errs = append(errs, p+": unknown node "+strconv.Quote(c.Node))
			continue
		}
		if err := validateChainCheck(c, n); err != nil {
			errs = append(errs, p+": "+err.Error())
		}
		if c.Expect != nil {
			if !c.Expect.empty() {
				errs = append(errs, p+".expect: an expect compares the check's own value (op and value only)")
			}
			errs = append(errs, validateChainCond(c.Expect, 0, nil, nil, p+".expect", true)...)
		}
	}
	for i, l := range d.Links {
		p := fmt.Sprintf("links[%d]", i)
		if _, ok := nodes[l.From]; !ok {
			errs = append(errs, p+": unknown node "+strconv.Quote(l.From))
		}
		if _, ok := nodes[l.To]; !ok {
			errs = append(errs, p+": unknown node "+strconv.Quote(l.To))
		}
		if l.Kind != "" && !chainLinkKinds[l.Kind] {
			errs = append(errs, p+": kind must be data, read, write, failover or dependency")
		}
		if l.Check != "" && !checks[l.Check] {
			errs = append(errs, p+": unknown check "+strconv.Quote(l.Check))
		}
		if len(l.Label) > 120 || chainHasControl(l.Label) {
			errs = append(errs, p+": label at most 120 characters")
		}
	}
	seen := map[string]bool{}
	for i := range d.Scenarios {
		sc := &d.Scenarios[i]
		p := fmt.Sprintf("scenarios[%d]", i)
		if !chainIDPattern.MatchString(sc.ID) || seen[sc.ID] {
			errs = append(errs, p+": id must be unique, 1-40 characters a-z, 0-9, '-' or '_'")
		}
		seen[sc.ID] = true
		if len(sc.Label) > 120 || chainHasControl(sc.Label) {
			errs = append(errs, p+": label at most 120 characters")
		}
		errs = append(errs, validateChainCond(&sc.When, 0, checks, servers, p+".when", false)...)
	}
	for i, l := range d.Links {
		if len(l.Scenarios) > chainMaxScenarios {
			errs = append(errs, fmt.Sprintf("links[%d].scenarios: at most %d", i, chainMaxScenarios))
			continue
		}
		for _, id := range l.Scenarios {
			if !seen[id] {
				errs = append(errs, fmt.Sprintf("links[%d].scenarios: unknown scenario %q", i, id))
			}
		}
	}
	return errs
}

func (d chainDef) interval() time.Duration {
	s := d.IntervalS
	if s == 0 {
		s = settingInt("chain_interval_s")
	}
	if s < 15 {
		s = 15
	}
	return time.Duration(s) * time.Second
}

func (d chainDef) confirmN() int {
	if d.Confirm > 0 {
		return d.Confirm
	}
	return 2
}

func (d chainDef) node(id string) chainNode {
	for _, n := range d.Nodes {
		if n.ID == id {
			return n
		}
	}
	return chainNode{}
}

func (d chainDef) scenarioLabel(id string) string {
	for _, s := range d.Scenarios {
		if s.ID == id {
			if s.Label != "" {
				return s.Label
			}
			return s.ID
		}
	}
	return id
}

// ─── state and change detection ───────────────────────────────────────────────────────────

type chainNodeState struct {
	State string `json:"state"` // up, down (server nodes only)
	Error string `json:"error,omitempty"`
}

type chainState struct {
	At       string                    `json:"at"`
	Nodes    map[string]chainNodeState `json:"nodes"`
	Checks   map[string]chainResult    `json:"checks"`
	Scenario string                    `json:"scenario"`
}

type chainChange struct {
	Kind    string // node, check, scenario
	Subject string
	Old     string
	New     string
	OK      bool
	OKFlip  bool
}

type chainPending struct {
	sig   string
	count int
}

func chainCheckSig(r chainResult) string { return strconv.FormatBool(r.OK) + "|" + r.Value }

// chainApply moves the confirmed state towards what was observed. A key changes once its new
// value was seen confirm times in a row; a key the confirmed state does not know yet is
// adopted at once (a new chain or a new check), and only a server already down is reported.
func chainApply(confirmed *chainState, observed chainState, pending map[string]chainPending, confirm int) []chainChange {
	var out []chainChange
	if confirmed.Nodes == nil {
		confirmed.Nodes = map[string]chainNodeState{}
	}
	if confirmed.Checks == nil {
		confirmed.Checks = map[string]chainResult{}
	}
	step := func(key, sig string) bool {
		p := pending[key]
		if p.sig == sig {
			p.count++
		} else {
			p = chainPending{sig: sig, count: 1}
		}
		if p.count >= confirm {
			delete(pending, key)
			return true
		}
		pending[key] = p
		return false
	}
	for id, n := range observed.Nodes {
		old, known := confirmed.Nodes[id]
		switch {
		case !known:
			confirmed.Nodes[id] = n
			if n.State == "down" {
				out = append(out, chainChange{Kind: "node", Subject: id, New: "down"})
			}
		case old.State != n.State:
			if step("node:"+id, n.State) {
				confirmed.Nodes[id] = n
				out = append(out, chainChange{Kind: "node", Subject: id, Old: old.State, New: n.State})
			}
		default:
			delete(pending, "node:"+id)
			confirmed.Nodes[id] = n
		}
	}
	for id := range confirmed.Nodes {
		if _, ok := observed.Nodes[id]; !ok {
			delete(confirmed.Nodes, id)
		}
	}
	for id, r := range observed.Checks {
		old, known := confirmed.Checks[id]
		switch {
		case !known:
			confirmed.Checks[id] = r
		case chainCheckSig(old) != chainCheckSig(r):
			if step("check:"+id, chainCheckSig(r)) {
				confirmed.Checks[id] = r
				out = append(out, chainChange{Kind: "check", Subject: id, Old: old.Value, New: r.Value, OK: r.OK, OKFlip: old.OK != r.OK})
			}
		default:
			delete(pending, "check:"+id)
			confirmed.Checks[id] = r
		}
	}
	for id := range confirmed.Checks {
		if _, ok := observed.Checks[id]; !ok {
			delete(confirmed.Checks, id)
		}
	}
	switch {
	case confirmed.At == "":
		confirmed.Scenario = observed.Scenario
		out = append(out, chainChange{Kind: "scenario", Subject: "scenario", New: observed.Scenario})
	case confirmed.Scenario != observed.Scenario:
		if step("scenario", observed.Scenario) {
			out = append(out, chainChange{Kind: "scenario", Subject: "scenario", Old: confirmed.Scenario, New: observed.Scenario})
			confirmed.Scenario = observed.Scenario
		}
	default:
		delete(pending, "scenario")
	}
	confirmed.At = observed.At
	return out
}

func chainScenarioOf(d chainDef, s chainState) string {
	for _, sc := range d.Scenarios {
		if sc.When.eval(s) {
			return sc.ID
		}
	}
	return ""
}

// ─── storage ──────────────────────────────────────────────────────────────────────────────

const chainSchema = `CREATE TABLE IF NOT EXISTS chains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL DEFAULT '',
		definition TEXT NOT NULL DEFAULT '{}',
		enabled INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT '',
		last_at TEXT NOT NULL DEFAULT '',
		last_state TEXT NOT NULL DEFAULT '',
		confirmed TEXT NOT NULL DEFAULT '',
		last_error TEXT NOT NULL DEFAULT '')`

const chainSampleSchema = `CREATE TABLE IF NOT EXISTS chain_samples (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chain_id INTEGER NOT NULL,
		at TEXT NOT NULL DEFAULT '',
		kind TEXT NOT NULL DEFAULT '',
		subject TEXT NOT NULL DEFAULT '',
		old_value TEXT NOT NULL DEFAULT '',
		new_value TEXT NOT NULL DEFAULT '',
		ok INTEGER NOT NULL DEFAULT 0)`

type chainRow struct {
	ID        int        `json:"id"`
	UserID    int        `json:"-"`
	Enabled   bool       `json:"enabled"`
	Def       chainDef   `json:"definition"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
	LastAt    string     `json:"last_at"`
	LastError string     `json:"last_error,omitempty"`
	State     chainState `json:"state"`
	Confirmed chainState `json:"-"`
}

func loadChain(id int) (chainRow, error) {
	var c chainRow
	var def, last, conf string
	var en int
	err := db.QueryRow(`SELECT id, user_id, enabled, definition, created_at, updated_at, last_at, last_state, confirmed, last_error
		FROM chains WHERE id=?`, id).Scan(&c.ID, &c.UserID, &en, &def, &c.CreatedAt, &c.UpdatedAt, &c.LastAt, &last, &conf, &c.LastError)
	if err != nil {
		return c, err
	}
	c.Enabled = en == 1
	json.Unmarshal([]byte(def), &c.Def)
	json.Unmarshal([]byte(last), &c.State)
	json.Unmarshal([]byte(conf), &c.Confirmed)
	return c, nil
}

// chainRuns keeps the not-yet-confirmed changes per chain and prevents two polls of the same
// chain at once.
var chainRuns = struct {
	sync.Mutex
	pending map[int]map[string]chainPending
	busy    map[int]bool
}{pending: map[int]map[string]chainPending{}, busy: map[int]bool{}}

// ─── polling ──────────────────────────────────────────────────────────────────────────────

// observeChain runs every check of a chain once and returns what it saw.
func observeChain(ctx context.Context, c chainRow) chainState {
	s := chainState{At: nowRFC(), Nodes: map[string]chainNodeState{}, Checks: map[string]chainResult{}}
	var mu sync.Mutex
	put := func(id string, r chainResult) {
		mu.Lock()
		s.Checks[id] = r
		mu.Unlock()
	}
	byNode := map[string][]chainCheck{}
	var wrmChecks []chainCheck
	for _, ch := range c.Def.Checks {
		if chainCheckTypes[ch.Type].Where == "server" && ch.Type != "monitor" {
			byNode[ch.Node] = append(byNode[ch.Node], ch)
		} else {
			wrmChecks = append(wrmChecks, ch)
		}
	}
	var wg sync.WaitGroup
	for _, n := range c.Def.Nodes {
		if n.Kind != "server" {
			continue
		}
		wg.Add(1)
		go func(n chainNode) {
			defer wg.Done()
			st, results := observeChainServer(ctx, c.UserID, n, byNode[n.ID])
			mu.Lock()
			s.Nodes[n.ID] = st
			mu.Unlock()
			for id, r := range results {
				put(id, r)
			}
		}(n)
	}
	sem := make(chan struct{}, 4)
	for _, ch := range wrmChecks {
		wg.Add(1)
		go func(ch chainCheck) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			put(ch.ID, runChainWRMCheck(ctx, ch, c.UserID, c.Def.node(ch.Node)))
		}(ch)
	}
	wg.Wait()
	s.Scenario = chainScenarioOf(c.Def, s)
	return s
}

func observeChainServer(ctx context.Context, owner int, n chainNode, checks []chainCheck) (chainNodeState, map[string]chainResult) {
	fail := func(msg string) (chainNodeState, map[string]chainResult) {
		res := map[string]chainResult{}
		for _, ch := range checks {
			res[ch.ID] = chainResult{Value: "unknown", Error: "the server is not reachable"}
		}
		return chainNodeState{State: "down", Error: msg}, res
	}
	conn, err := loadConnection(n.ConnID)
	if err != nil || conn.UserID != owner || !isSSHProtocol(conn) {
		return fail("the connection is not available to the chain's owner")
	}
	select {
	case monitorSlots <- struct{}{}:
		defer func() { <-monitorSlots }()
	case <-ctx.Done():
		return fail("timed out waiting for a free check slot")
	}
	cl, err := dialSSH(conn, nil)
	if err != nil {
		return fail(chainShortErr(err))
	}
	defer cl.Close()
	if len(checks) == 0 {
		return chainNodeState{State: "up"}, nil
	}
	runCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	res, err := runChainNode(runCtx, cl, checks)
	if err != nil {
		st, r := fail(err.Error())
		st.State = "up" // SSH answered; the checks themselves failed
		return st, r
	}
	return chainNodeState{State: "up"}, res
}

// pollChain observes a chain, updates its confirmed state, records the changes and notifies.
func pollChain(ctx context.Context, id int) (chainRow, error) {
	chainRuns.Lock()
	if chainRuns.busy[id] {
		chainRuns.Unlock()
		return chainRow{}, fmt.Errorf("this chain is being checked right now")
	}
	chainRuns.busy[id] = true
	pending := chainRuns.pending[id]
	if pending == nil {
		pending = map[string]chainPending{}
		chainRuns.pending[id] = pending
	}
	chainRuns.Unlock()
	defer func() { chainRuns.Lock(); delete(chainRuns.busy, id); chainRuns.Unlock() }()

	c, err := loadChain(id)
	if err != nil {
		return c, err
	}
	observed := observeChain(ctx, c)
	chainRuns.Lock()
	changes := chainApply(&c.Confirmed, observed, pending, c.Def.confirmN())
	chainRuns.Unlock()
	c.State, c.LastAt = observed, observed.At
	db.Exec(`UPDATE chains SET last_at=?, last_state=?, confirmed=?, last_error='' WHERE id=?`,
		observed.At, string(jsonMarshal(observed)), string(jsonMarshal(c.Confirmed)), id)
	for _, ch := range changes {
		db.Exec(`INSERT INTO chain_samples (chain_id, at, kind, subject, old_value, new_value, ok) VALUES (?,?,?,?,?,?,?)`,
			id, observed.At, ch.Kind, ch.Subject, truncateStr(ch.Old, 300), truncateStr(ch.New, 300), boolInt(ch.OK))
		notifyChainChange(c, ch)
	}
	if len(changes) > 0 {
		db.Exec(`DELETE FROM chain_samples WHERE chain_id=? AND id NOT IN (SELECT id FROM chain_samples WHERE chain_id=? ORDER BY id DESC LIMIT ?)`,
			id, id, chainSamplesKept)
	}
	hub.sendTo(c.UserID, []byte(fmt.Sprintf(`{"type":"chain_changed","id":%d}`, id)))
	return c, nil
}

// notifyChainChange sends the events of one confirmed change. Every change is in the timeline;
// a check notifies only when it turns good or bad, not on every new value.
func notifyChainChange(c chainRow, ch chainChange) {
	name := c.Def.Name
	switch ch.Kind {
	case "scenario":
		if ch.Old == "" && ch.New != "" {
			return // the first state of a new chain
		}
		from, to := c.Def.scenarioLabel(ch.Old), c.Def.scenarioLabel(ch.New)
		if to == "" {
			to = "—"
		}
		notifyUser(c.UserID, "chain.scenario_changed", func(lang string) (string, string) {
			if lang == "hr" {
				return "Lanac " + name + ": " + to, "Aktivni scenarij lanca " + name + " promijenio se iz „" + from + "“ u „" + to + "“."
			}
			return "Chain " + name + ": " + to, "The active scenario of " + name + " changed from \"" + from + "\" to \"" + to + "\"."
		})
	case "node":
		label := c.Def.node(ch.Subject).Label
		if label == "" {
			label = ch.Subject
		}
		event := "chain.node_down"
		if ch.New == "up" {
			if ch.Old != "down" {
				return
			}
			event = "chain.node_up"
		}
		errText := c.State.Nodes[ch.Subject].Error
		notifyUser(c.UserID, event, func(lang string) (string, string) {
			if lang == "hr" {
				if event == "chain.node_up" {
					return "Lanac " + name + ": " + label + " je ponovno dostupan", label + " ponovno odgovara."
				}
				return "Lanac " + name + ": " + label + " nije dostupan", label + ": " + errText
			}
			if event == "chain.node_up" {
				return "Chain " + name + ": " + label + " is up again", label + " answers again."
			}
			return "Chain " + name + ": " + label + " is down", label + ": " + errText
		})
	case "check":
		if !ch.OKFlip {
			return
		}
		label := ch.Subject
		for _, x := range c.Def.Checks {
			if x.ID == ch.Subject && x.Label != "" {
				label = x.Label
			}
		}
		notifyUser(c.UserID, "chain.check_changed", func(lang string) (string, string) {
			if lang == "hr" {
				state := "nije u redu"
				if ch.OK {
					state = "je u redu"
				}
				return "Lanac " + name + ": " + label + " " + state, label + ": " + ch.Old + " → " + ch.New
			}
			state := "failed"
			if ch.OK {
				state = "is OK again"
			}
			return "Chain " + name + ": " + label + " " + state, label + ": " + ch.Old + " → " + ch.New
		})
	}
}

// chainsLoop checks every enabled chain whose interval has passed.
func chainsLoop() {
	time.Sleep(30 * time.Second)
	sem := make(chan struct{}, 2)
	for {
		if settingBool("chains_enabled") {
			for _, id := range dueChains() {
				sem <- struct{}{}
				go func(id int) {
					defer func() { <-sem }()
					ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
					defer cancel()
					pollChain(ctx, id)
				}(id)
			}
		}
		time.Sleep(15 * time.Second)
	}
}

func dueChains() []int {
	rows, err := db.Query(`SELECT c.id, c.definition, c.last_at FROM chains c JOIN users u ON u.id=c.user_id AND u.disabled=0 WHERE c.enabled=1`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		var def, last string
		if rows.Scan(&id, &def, &last) != nil {
			continue
		}
		var d chainDef
		json.Unmarshal([]byte(def), &d)
		at, _ := time.Parse(time.RFC3339, last)
		if time.Since(at) >= d.interval() {
			ids = append(ids, id)
		}
	}
	return ids
}

// ─── API ──────────────────────────────────────────────────────────────────────────────────

// chainSummary is one entry of the list: enough for a status line without the definition.
type chainSummary struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Enabled       bool   `json:"enabled"`
	Scenario      string `json:"scenario"`
	ScenarioLabel string `json:"scenario_label"`
	NodesUp       int    `json:"nodes_up"`
	NodesTotal    int    `json:"nodes_total"`
	ChecksOK      int    `json:"checks_ok"`
	ChecksTotal   int    `json:"checks_total"`
	LastAt        string `json:"last_at"`
}

func summarizeChain(c chainRow) chainSummary {
	s := chainSummary{ID: c.ID, Name: c.Def.Name, Enabled: c.Enabled, Scenario: c.Confirmed.Scenario,
		ScenarioLabel: c.Def.scenarioLabel(c.Confirmed.Scenario), ChecksTotal: len(c.Def.Checks), LastAt: c.LastAt}
	for _, n := range c.Def.Nodes {
		if n.Kind == "server" {
			s.NodesTotal++
			if c.State.Nodes[n.ID].State == "up" {
				s.NodesUp++
			}
		}
	}
	for _, r := range c.State.Checks {
		if r.OK {
			s.ChecksOK++
		}
	}
	return s
}

// chainErrors answers a definition that did not validate, with every problem listed.
func chainErrors(w http.ResponseWriter, errs []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": errs[0], "errors": errs})
}

// requireChain authenticates and loads a chain its owner may see.
func requireChain(w http.ResponseWriter, r *http.Request, id int) (chainRow, int, bool) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return chainRow{}, 0, false
	}
	c, err := loadChain(id)
	if err != nil || c.UserID != userID {
		jsonError(w, "Not found", http.StatusNotFound)
		return chainRow{}, 0, false
	}
	return c, userID, true
}

func decodeChainDef(w http.ResponseWriter, r *http.Request) (chainDef, bool, bool) {
	var in struct {
		Definition chainDef `json:"definition"`
		Enabled    *bool    `json:"enabled"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, chainMaxDefBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		jsonError(w, "The chain is not valid JSON: "+truncateStr(err.Error(), 200), http.StatusBadRequest)
		return chainDef{}, false, false
	}
	enabled := in.Enabled == nil || *in.Enabled
	return in.Definition, enabled, true
}

func apiChainsHandler(w http.ResponseWriter, r *http.Request) {
	if !settingBool("chains_enabled") {
		if _, ok := requireAuth(w, r); ok {
			jsonError(w, "Chains are turned off by the administrator", http.StatusForbidden)
		}
		return
	}
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/chains"), "/")
	switch rest {
	case "":
		apiChainsRoot(w, r)
		return
	case "catalog":
		if _, ok := requireAuth(w, r); ok {
			jsonOK(w, map[string]interface{}{"types": chainCheckTypes, "ops": []string{"==", "!=", "in", "contains", "exists"},
				"link_kinds": []string{"data", "read", "write", "failover", "dependency"}, "formats": []string{"ini", "json", "yaml", "env", "kv", "regex"}})
		}
		return
	case "templates":
		apiChainTemplates(w, r)
		return
	case "validate":
		userID, ok := requireAuth(w, r)
		if !ok || !requireMethod(w, r, http.MethodPost) {
			return
		}
		d, _, ok := decodeChainDef(w, r)
		if !ok {
			return
		}
		errs := validateChain(&d, userID)
		jsonOK(w, map[string]interface{}{"ok": len(errs) == 0, "errors": errs})
		return
	}
	parts := strings.Split(rest, "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	c, userID, ok := requireChain(w, r, id)
	if !ok {
		return
	}
	if len(parts) == 1 {
		apiChainByID(w, r, c, userID)
		return
	}
	switch parts[1] {
	case "layout":
		apiChainLayout(w, r, c)
	case "check":
		if !requireMethod(w, r, http.MethodPost) {
			return
		}
		if !mcpRateOK(fmt.Sprintf("chain-check:%d", userID), 20, time.Minute) {
			jsonError(w, "Too many checks; wait a moment", http.StatusTooManyRequests)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()
		got, err := pollChain(ctx, c.ID)
		if err != nil {
			jsonError(w, err.Error(), http.StatusConflict)
			return
		}
		auditLog(r, userID, "chain.checked", c.Def.Name, map[string]interface{}{"chain": c.ID})
		jsonOK(w, chainView(got))
	case "timeline":
		rows, err := db.Query(`SELECT at, kind, subject, old_value, new_value, ok FROM chain_samples WHERE chain_id=? ORDER BY id DESC LIMIT 300`, c.ID)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		out := []map[string]interface{}{}
		for rows.Next() {
			var at, kind, subject, oldV, newV string
			var okv int
			if rows.Scan(&at, &kind, &subject, &oldV, &newV, &okv) == nil {
				out = append(out, map[string]interface{}{"at": at, "kind": kind, "subject": subject, "old": oldV, "new": newV, "ok": okv == 1})
			}
		}
		jsonOK(w, out)
	case "export":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", contentDisposition("attachment", "chain-"+chainFileName(c.Def.Name)+".json"))
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(map[string]interface{}{"definition": c.Def, "enabled": c.Enabled})
	default:
		jsonError(w, "Not found", http.StatusNotFound)
	}
}

func chainFileName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteByte('-')
		}
	}
	if s := strings.Trim(b.String(), "-"); s != "" {
		return truncateStr(s, 60)
	}
	return "export"
}

func chainView(c chainRow) map[string]interface{} {
	return map[string]interface{}{"id": c.ID, "enabled": c.Enabled, "definition": c.Def, "state": c.State,
		"scenario": c.Confirmed.Scenario, "scenario_label": c.Def.scenarioLabel(c.Confirmed.Scenario),
		"last_at": c.LastAt, "created_at": c.CreatedAt, "updated_at": c.UpdatedAt}
}

func apiChainsRoot(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query(`SELECT id FROM chains WHERE user_id=? ORDER BY name COLLATE NOCASE`, userID)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		var ids []int
		for rows.Next() {
			var id int
			if rows.Scan(&id) == nil {
				ids = append(ids, id)
			}
		}
		rows.Close()
		out := []chainSummary{}
		for _, id := range ids {
			if c, err := loadChain(id); err == nil {
				out = append(out, summarizeChain(c))
			}
		}
		jsonOK(w, out)
	case http.MethodPost:
		d, enabled, ok := decodeChainDef(w, r)
		if !ok {
			return
		}
		if errs := validateChain(&d, userID); len(errs) > 0 {
			chainErrors(w, errs)
			return
		}
		var n int
		db.QueryRow(`SELECT COUNT(1) FROM chains WHERE user_id=?`, userID).Scan(&n)
		if n >= 50 {
			jsonError(w, "At most 50 chains per user", http.StatusBadRequest)
			return
		}
		now := nowRFC()
		res, err := db.Exec(`INSERT INTO chains (user_id, name, definition, enabled, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
			userID, d.Name, string(jsonMarshal(d)), boolInt(enabled), now, now)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		auditLog(r, userID, "chain.created", d.Name, map[string]interface{}{"chain": id, "nodes": len(d.Nodes), "checks": len(d.Checks)})
		c, _ := loadChain(int(id))
		jsonOK(w, chainView(c))
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func apiChainByID(w http.ResponseWriter, r *http.Request, c chainRow, userID int) {
	switch r.Method {
	case http.MethodGet:
		jsonOK(w, chainView(c))
	case http.MethodPut:
		d, enabled, ok := decodeChainDef(w, r)
		if !ok {
			return
		}
		if errs := validateChain(&d, userID); len(errs) > 0 {
			chainErrors(w, errs)
			return
		}
		if _, err := db.Exec(`UPDATE chains SET name=?, definition=?, enabled=?, updated_at=? WHERE id=?`,
			d.Name, string(jsonMarshal(d)), boolInt(enabled), nowRFC(), c.ID); err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		chainRuns.Lock()
		delete(chainRuns.pending, c.ID)
		chainRuns.Unlock()
		auditLog(r, userID, "chain.updated", d.Name, map[string]interface{}{"chain": c.ID, "nodes": len(d.Nodes), "checks": len(d.Checks), "enabled": enabled})
		got, _ := loadChain(c.ID)
		jsonOK(w, chainView(got))
	case http.MethodDelete:
		db.Exec(`DELETE FROM chain_samples WHERE chain_id=?`, c.ID)
		if _, err := db.Exec(`DELETE FROM chains WHERE id=?`, c.ID); err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		chainRuns.Lock()
		delete(chainRuns.pending, c.ID)
		chainRuns.Unlock()
		auditLog(r, userID, "chain.deleted", c.Def.Name, map[string]interface{}{"chain": c.ID})
		jsonOK(w, map[string]bool{"ok": true})
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
