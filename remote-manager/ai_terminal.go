package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/crypto/ssh"
)

// ─── AI IN THE LIVE TERMINAL (v12.2.0) ───────────────
//
// Every SSH shell terminal has a termShare: the terminal's tracker (termtrack.go, fed with
// the shell's output by ssh_ws.go), a way to type into it, and — only while the user turns
// it on — the AI sessions and AI connections (MCP grants) that may use it:
//
//   share     the user turns sharing on or off for this terminal window (never automatic);
//             the policy ai_terminal_share + ai_terminal_share_rules decides whether it is
//             allowed and at which level: read (terminal_read) or run (also terminal_run);
//   bind      a built-in panel session of the same user and connection uses the terminal;
//   attach    an AI connection (MCP) is bound to the terminal: its sessions on the
//             connection use it (attach turns sharing on);
//   items     selections the user sends to an attached AI app ("Ask AI"), fetched with
//             terminal_read.
//
// Sharing ends with the window (the WebSocket), with the terminal's SSH session, when the
// policy no longer allows it, and for an AI connection when it is revoked or expires.
// Everything the AI reads from the terminal is redacted and framed as untrusted data;
// terminal_run goes through the permission engine (CallTool) like run_command, and is typed
// only when the terminal is at a shell prompt (no full-screen program, no password prompt,
// the shell in the foreground).
//
// Browser API (/api/ai/terminals/{id}/…, the id is sent on the terminal WebSocket):
//   GET  /                    state, the AI connections that can be attached, the policy
//   POST share   {on}         turn sharing on / off
//   POST bind    {session_id} bind a panel AI session
//   POST attach  {grant_id}   attach an AI connection (MCP)
//   POST detach  {grant_id | session_id}
//   POST selection {grant_id, text, ask}   a pending context item for an attached AI app

const (
	aiTermItemMax      = 20
	aiTermItemMaxBytes = 16 << 10
	termExitUnknown    = -999
)

type termAttach struct {
	GrantID int64
	Name    string
	Client  string
	Since   time.Time
}

type termItem struct {
	ID         string
	GrantID    int64
	Text       string
	Ask        string
	Redactions int
	Created    time.Time
}

type termShare struct {
	mu        sync.Mutex
	ID        string // the browser's handle
	PubID     string // the id MCP clients see
	UserID    int
	Username  string
	ConnID    int
	ConnName  string
	Host      string
	ShareID   int
	TermID    int64 // terminal_sessions row
	tr        *termTracker
	typeFn    func([]byte) bool
	sendCtl   func(interface{})
	ssh       *ssh.Client
	on        bool
	level     string
	since     time.Time
	lastInput time.Time
	running   bool
	sessions  map[int64]bool
	grants    map[int64]*termAttach
	items     []termItem
	closed    bool
	cwdCache  string
	cwdAt     time.Time
}

var aiTerms = struct {
	sync.Mutex
	m map[string]*termShare
}{m: map[string]*termShare{}}

func aiTermGet(id string) *termShare {
	aiTerms.Lock()
	defer aiTerms.Unlock()
	return aiTerms.m[id]
}

func aiTermList(pred func(*termShare) bool) []*termShare {
	aiTerms.Lock()
	defer aiTerms.Unlock()
	var out []*termShare
	for _, ts := range aiTerms.m {
		if pred == nil || pred(ts) {
			out = append(out, ts)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TermID < out[j].TermID })
	return out
}

// registerAITerm creates the share of a new SSH shell terminal (sharing off).
func registerAITerm(userID, shareID int, c Connection, termID int64, cl *ssh.Client, typeFn func([]byte) bool, sendCtl func(interface{})) *termShare {
	ts := &termShare{ID: randomToken(16), PubID: "t" + randomToken(5), UserID: userID, Username: usernameOf(userID), ConnID: c.ID, ConnName: c.Name,
		Host: c.Host, ShareID: shareID, TermID: termID, tr: newTermTracker(), typeFn: typeFn, sendCtl: sendCtl, ssh: cl,
		sessions: map[int64]bool{}, grants: map[int64]*termAttach{}}
	aiTerms.Lock()
	aiTerms.m[ts.ID] = ts
	aiTerms.Unlock()
	return ts
}

func (ts *termShare) feed(p []byte) { ts.tr.Write(p) }

// noteInput records the user's own keystrokes (WRM does not type while the user types).
func (ts *termShare) noteInput() {
	ts.mu.Lock()
	ts.lastInput = time.Now()
	ts.mu.Unlock()
}

func (ts *termShare) audit(r *http.Request, action string, d map[string]interface{}) {
	if d == nil {
		d = map[string]interface{}{}
	}
	d["terminal"] = ts.PubID
	d["terminal_session"] = ts.TermID
	auditLogRef(r, ts.UserID, ts.Username, action, ts.ConnName, d, auditRef{ConnID: ts.ConnID, SessionID: int(ts.TermID)})
}

func (ts *termShare) isOn() (bool, string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.on && !ts.closed, ts.level
}

// ─── policy ──────────────────────────────────────────

// aiShareRule limits sharing of matching connections' terminals: off, read or run.
type aiShareRule struct {
	Match string `json:"match"`
	Value string `json:"value"`
	Share string `json:"share"`
}

var aiShareLevels = []string{"off", "read", "run"}

func aiShareRank(l string) int {
	for i, x := range aiShareLevels {
		if x == l {
			return i
		}
	}
	return 0
}

func parseAIShareRules(v string) ([]aiShareRule, error) {
	if strings.TrimSpace(v) == "" {
		return nil, nil
	}
	var rules []aiShareRule
	if err := json.Unmarshal([]byte(v), &rules); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}
	if len(rules) > 200 {
		return nil, fmt.Errorf("at most 200 rules")
	}
	for i := range rules {
		r := &rules[i]
		mr := aiModeRule{Match: r.Match, Value: r.Value}
		if err := normalizeRuleTarget(&mr, i); err != nil {
			return nil, err
		}
		r.Match, r.Value = mr.Match, mr.Value
		r.Share = strings.ToLower(strings.TrimSpace(r.Share))
		if !oneOf(r.Share, aiShareLevels...) {
			return nil, fmt.Errorf("rule %d: share must be off, read or run", i+1)
		}
	}
	return rules, nil
}

func validateAIShareRules(v string) (string, error) {
	rules, err := parseAIShareRules(v)
	if err != nil {
		return "", err
	}
	if rules == nil {
		rules = []aiShareRule{}
	}
	b, _ := json.Marshal(rules)
	return string(b), nil
}

// aiTerminalShareLevel is what the policy allows for the user's terminals of a connection:
// off, read or run. The most restrictive of ai_terminal_share and every matching rule wins.
func aiTerminalShareLevel(userID, connID int) (string, string) {
	switch getSetting("ai_terminal_share") {
	case "off":
		return "off", "sharing terminals with AI is turned off (policy ai_terminal_share)"
	case "admins":
		if !isAdminUser(userID) {
			return "off", "sharing terminals with AI is limited to administrators"
		}
	}
	if settingBool("ai_kill_switch") {
		return "off", "the AI assistant is stopped by the administrator (kill switch)"
	}
	if ok1, _ := aiUserAllowed(userID); !ok1 {
		if ok2, why := mcpUserAllowed(userID); !ok2 {
			return "off", why
		}
	}
	level, why := "run", ""
	rules, _ := parseAIShareRules(getSetting("ai_terminal_share_rules"))
	f := aiFactsOf(connID)
	for _, r := range rules {
		if (aiModeRule{Match: r.Match, Value: r.Value}).matches(f) && aiShareRank(r.Share) < aiShareRank(level) {
			level, why = r.Share, "rule "+r.Match+" "+r.Value+": "+r.Share
		}
	}
	return level, why
}

func aiTerminalShareFlag(userID int) bool {
	if getSetting("ai_terminal_share") == "off" || (getSetting("ai_terminal_share") == "admins" && !isAdminUser(userID)) {
		return false
	}
	return aiAllowedFlag(userID) || mcpAllowedFlag(userID)
}

func aiContextLines() int {
	return max(1, min(settingInt("ai_terminal_context_lines"), termTrackMaxLines))
}

// ─── share, bind, attach ─────────────────────────────

// setShare turns sharing on or off at the user's request.
func (ts *termShare) setShare(r *http.Request, on bool, why string) error {
	if on {
		level, reason := aiTerminalShareLevel(ts.UserID, ts.ConnID)
		if level == "off" {
			ts.audit(r, "ai.terminal_share_refused", map[string]interface{}{"reason": reason})
			return fmt.Errorf("%s", reason)
		}
		ts.mu.Lock()
		if ts.closed {
			ts.mu.Unlock()
			return fmt.Errorf("the terminal is closed")
		}
		was, old := ts.on, ts.level
		ts.on, ts.level = true, level
		if !was {
			ts.since = time.Now()
		}
		ts.mu.Unlock()
		if !was || old != level {
			ts.audit(r, "ai.terminal_share_on", map[string]interface{}{"level": level, "lines": aiContextLines(), "why": why, "rule": reason})
		}
		ts.push()
		return nil
	}
	ts.stopSharing(r, why)
	return nil
}

// stopSharing turns sharing off: AI sessions are unbound, AI connections detached, pending
// items dropped.
func (ts *termShare) stopSharing(r *http.Request, why string) {
	ts.mu.Lock()
	was := ts.on
	ts.on = false
	sids := make([]int64, 0, len(ts.sessions))
	for id := range ts.sessions {
		sids = append(sids, id)
	}
	grants := make([]int64, 0, len(ts.grants))
	for id := range ts.grants {
		grants = append(grants, id)
	}
	ts.sessions = map[int64]bool{}
	ts.grants = map[int64]*termAttach{}
	ts.items = nil
	ts.mu.Unlock()
	for _, id := range sids {
		if s := aiGet(id); s != nil {
			s.setTerm(nil, "the user stopped sharing the terminal")
		}
	}
	if was {
		ts.audit(r, "ai.terminal_share_off", map[string]interface{}{"why": why, "sessions": len(sids), "ai_connections": len(grants)})
	}
	for _, g := range grants {
		mcpNotifyUserOf(g)
	}
	ts.push()
}

// close ends the share with the terminal window or its SSH session.
func (ts *termShare) close(why string) {
	ts.stopSharing(nil, why)
	ts.mu.Lock()
	ts.closed = true
	ts.mu.Unlock()
	aiTerms.Lock()
	delete(aiTerms.m, ts.ID)
	aiTerms.Unlock()
}

// bindSession lets an AI session use the terminal (sharing must be on).
func (ts *termShare) bindSession(r *http.Request, s *aiSession, via string) error {
	on, _ := ts.isOn()
	if !on {
		return fmt.Errorf("the terminal is not shared")
	}
	if s.UserID != ts.UserID || s.ConnID != ts.ConnID {
		return fmt.Errorf("the AI session belongs to another connection")
	}
	if cur := s.boundTerm(); cur == ts {
		return nil
	} else if cur != nil {
		cur.unbindSession(s.ID)
	}
	ts.mu.Lock()
	ts.sessions[s.ID] = true
	ts.mu.Unlock()
	s.setTerm(ts, "")
	ts.audit(r, "ai.terminal_attached", map[string]interface{}{"ai_session": s.ID, "transport": s.Transport, "via": via, "mcp_grant": s.TokenID})
	ts.push()
	return nil
}

func (ts *termShare) unbindSession(id int64) {
	ts.mu.Lock()
	had := ts.sessions[id]
	delete(ts.sessions, id)
	ts.mu.Unlock()
	if had {
		ts.push()
	}
}

// attachGrant binds an AI connection (MCP) to the terminal; its sessions on the connection
// use it. Sharing is turned on for it.
func (ts *termShare) attachGrant(r *http.Request, g *mcpGrant) error {
	if g.UserID != ts.UserID || !g.active(time.Now()) {
		return fmt.Errorf("AI connection not found")
	}
	if ok, why := mcpUserAllowed(ts.UserID); !ok {
		return fmt.Errorf("%s", why)
	}
	if ts.ShareID != 0 || !g.hasConn(ts.ConnID) || !userOwnsConnection(ts.ConnID, ts.UserID) {
		return fmt.Errorf("this AI connection does not include %s", ts.ConnName)
	}
	if !oneOf("terminal_read", g.Scopes...) && !oneOf("terminal_with_approval", g.Scopes...) {
		return fmt.Errorf("this AI connection has neither the terminal_read nor the terminal_with_approval scope; create one with a terminal scope")
	}
	if err := ts.setShare(r, true, "attach "+g.Name); err != nil {
		return err
	}
	ts.mu.Lock()
	ts.grants[g.ID] = &termAttach{GrantID: g.ID, Name: g.Name, Client: g.ClientName, Since: time.Now()}
	ts.mu.Unlock()
	n := 0
	for _, s := range aiActiveSessions(func(s *aiSession) bool { return s.TokenID == g.ID && s.ConnID == ts.ConnID && s.boundTerm() == nil }) {
		if ts.bindSession(r, s, "attach") == nil {
			n++
		}
	}
	ts.audit(r, "ai.terminal_attached", map[string]interface{}{"mcp_grant": g.ID, "name": g.Name, "client": g.ClientName, "scopes": strings.Join(g.Scopes, ","), "sessions": n})
	mcpNotifyUser(ts.UserID)
	ts.push()
	return nil
}

// detachGrant ends an AI connection's binding (one click, revoke, expiry, policy).
func (ts *termShare) detachGrant(r *http.Request, gid int64, why string) bool {
	ts.mu.Lock()
	_, had := ts.grants[gid]
	delete(ts.grants, gid)
	var keep []termItem
	for _, it := range ts.items {
		if it.GrantID != gid {
			keep = append(keep, it)
		}
	}
	ts.items = keep
	var sids []int64
	for id := range ts.sessions {
		if s := aiGet(id); s != nil && s.TokenID == gid {
			sids = append(sids, id)
			delete(ts.sessions, id)
		}
	}
	ts.mu.Unlock()
	for _, id := range sids {
		if s := aiGet(id); s != nil {
			s.setTerm(nil, why)
		}
	}
	if had || len(sids) > 0 {
		ts.audit(r, "ai.terminal_detached", map[string]interface{}{"mcp_grant": gid, "why": why, "sessions": len(sids)})
		mcpNotifyUser(ts.UserID)
		ts.push()
	}
	return had
}

func (ts *termShare) detachSession(r *http.Request, s *aiSession, why string) {
	ts.unbindSession(s.ID)
	s.setTerm(nil, why)
	ts.audit(r, "ai.terminal_detached", map[string]interface{}{"ai_session": s.ID, "why": why})
}

func (ts *termShare) hasGrant(gid int64) bool {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.grants[gid] != nil
}

// aiTermDetachGrant detaches an AI connection from every terminal (revoke, expiry).
func aiTermDetachGrant(gid int64, why string) {
	for _, ts := range aiTermList(func(ts *termShare) bool { return ts.hasGrant(gid) }) {
		ts.detachGrant(nil, gid, why)
	}
}

// aiTermsOfGrant lists the terminals an AI connection is attached to (on a connection).
func aiTermsOfGrant(gid int64, connID int) []*termShare {
	return aiTermList(func(ts *termShare) bool { return ts.hasGrant(gid) && (connID == 0 || ts.ConnID == connID) })
}

// aiTermJanitor applies policy changes and expiry to shared terminals.
func aiTermJanitor() {
	for _, ts := range aiTermList(nil) {
		on, level := ts.isOn()
		if !on {
			continue
		}
		ts.mu.Lock()
		var gids []int64
		for id := range ts.grants {
			gids = append(gids, id)
		}
		ts.mu.Unlock()
		for _, gid := range gids {
			if !mcpGrantActive(gid) {
				ts.detachGrant(nil, gid, "the AI connection was revoked or has expired")
			}
		}
		now, why := aiTerminalShareLevel(ts.UserID, ts.ConnID)
		if now == "off" {
			ts.stopSharing(nil, "the policy no longer allows sharing: "+why)
		} else if now != level {
			ts.mu.Lock()
			ts.level = now
			ts.mu.Unlock()
			ts.audit(nil, "ai.terminal_share_on", map[string]interface{}{"level": now, "why": "policy changed", "rule": why})
			ts.push()
		}
	}
}

// view is the share as the browser sees it.
func (ts *termShare) view() map[string]interface{} {
	ts.mu.Lock()
	on, level, since := ts.on, ts.level, ts.since
	var sids []int64
	for id := range ts.sessions {
		sids = append(sids, id)
	}
	apps := []map[string]interface{}{}
	for _, a := range ts.grants {
		n := 0
		for _, it := range ts.items {
			if it.GrantID == a.GrantID {
				n++
			}
		}
		apps = append(apps, map[string]interface{}{"grant_id": a.GrantID, "name": a.Name, "client": a.Client, "since": a.Since.UTC().Format(time.RFC3339), "pending_items": n})
	}
	ts.mu.Unlock()
	sort.Slice(sids, func(i, j int) bool { return sids[i] < sids[j] })
	sort.Slice(apps, func(i, j int) bool { return apps[i]["grant_id"].(int64) < apps[j]["grant_id"].(int64) })
	sessions := []map[string]interface{}{}
	for _, id := range sids {
		if s := aiGet(id); s != nil {
			sessions = append(sessions, map[string]interface{}{"id": s.ID, "transport": s.Transport, "model": s.Model, "mode": s.Mode, "mcp_grant": s.TokenID})
		}
	}
	v := map[string]interface{}{"type": "ai_share", "id": ts.ID, "terminal": ts.PubID, "on": on, "level": level, "lines": aiContextLines(),
		"sessions": sessions, "apps": apps}
	if on {
		v["since"] = since.UTC().Format(time.RFC3339)
	}
	return v
}

// push sends the share's state to the terminal window.
func (ts *termShare) push() {
	if ts.sendCtl != nil {
		ts.sendCtl(ts.view())
	}
}

// pushApproval shows a pending terminal_run approval in the terminal window.
func (ts *termShare) pushApproval(s *aiSession, a *aiApproval) {
	if ts.sendCtl != nil {
		ts.sendCtl(map[string]interface{}{"type": "ai_approval", "session": s.ID, "transport": s.Transport, "client": s.Model, "approval": a})
	}
}

func (ts *termShare) pushApprovalDone(sid int64, aid, outcome string) {
	if ts.sendCtl != nil {
		ts.sendCtl(map[string]interface{}{"type": "ai_approval_done", "session": sid, "approval_id": aid, "outcome": outcome})
	}
}

// ─── what the AI reads ───────────────────────────────

func (ts *termShare) cwd() (string, string) {
	snap := ts.tr.snapshot(1)
	if snap.Cwd != "" {
		return snap.Cwd, "OSC 7"
	}
	ts.mu.Lock()
	c, at := ts.cwdCache, ts.cwdAt
	ts.mu.Unlock()
	if c != "" && time.Since(at) < 5*time.Second {
		return c, "the server's process table"
	}
	p, err := terminalCwd(ts.ssh)
	if err != nil {
		return "", ""
	}
	ts.mu.Lock()
	ts.cwdCache, ts.cwdAt = p, time.Now()
	ts.mu.Unlock()
	return p, "the server's process table"
}

func termStateText(st string) string {
	return map[string]string{"prompt": "at a shell prompt", "running": "a command is running", "full_screen": "a full-screen program is open (alternate screen)",
		"password_prompt": "a password prompt is shown", "unknown": "unknown (not at a recognised shell prompt)"}[st]
}

// context is the shared terminal as terminal_read returns it (before redaction). Pending
// items of the AI connection are included and marked delivered.
func (ts *termShare) context(gid int64) string {
	n := aiContextLines()
	snap := ts.tr.snapshot(n)
	cwd, from := ts.cwd()
	ts.mu.Lock()
	since := ts.since
	level := ts.level
	var items []termItem
	if gid > 0 {
		var keep []termItem
		for _, it := range ts.items {
			if it.GrantID == gid {
				items = append(items, it)
			} else {
				keep = append(keep, it)
			}
		}
		ts.items = keep
	}
	ts.mu.Unlock()
	var b strings.Builder
	fmt.Fprintf(&b, "Terminal %s on %s (%s), shared by %s since %s; this session may %s.\n", ts.PubID, ts.ConnName, hostOnly(ts.Host), ts.Username,
		since.UTC().Format(time.RFC3339), map[string]string{"read": "read it (terminal_read)", "run": "read it and type commands with terminal_run"}[level])
	if cwd != "" {
		fmt.Fprintf(&b, "Current directory: %s (from %s)\n", cwd, from)
	} else {
		b.WriteString("Current directory: unknown\n")
	}
	fmt.Fprintf(&b, "State: %s\n", termStateText(snap.State))
	if snap.OSC133 {
		b.WriteString("Shell integration: OSC 133 markers present (commands and their output are exact)\n")
	} else {
		b.WriteString("Shell integration: not detected (the last command is found heuristically)\n")
	}
	if snap.LastCmd != "" {
		exit := snap.LastExit
		if exit == "" {
			exit = "unknown"
		}
		fmt.Fprintf(&b, "\nLast command: %s (exit code %s%s)\n", snap.LastCmd, exit, map[bool]string{true: ", heuristic", false: ""}[snap.Heuristic])
		fmt.Fprintf(&b, "Its output (%d lines):\n%s\n", len(snap.LastOut), strings.Join(snap.LastOut, "\n"))
	}
	fmt.Fprintf(&b, "\nThe last %d lines of the screen and scrollback (at most %d):\n%s\n", len(snap.Lines), n, strings.Join(snap.Lines, "\n"))
	if len(items) > 0 {
		b.WriteString("\nSelections the user sent you from this terminal (\"Ask AI\"):\n")
		for _, it := range items {
			fmt.Fprintf(&b, "<selection id=%q at=%q>\n%s\n</selection>\n", it.ID, it.Created.UTC().Format(time.RFC3339), it.Text)
			if it.Ask != "" {
				fmt.Fprintf(&b, "The user asks: %s\n", it.Ask)
			}
		}
		ts.audit(nil, "ai.terminal_selection_fetched", map[string]interface{}{"mcp_grant": gid, "items": len(items)})
		ts.push()
	}
	return b.String()
}

// ─── typing into the terminal ────────────────────────

// validTerminalCommand allows one line of printable text: no control characters, so a
// command can never send Enter early, Ctrl+C, escape sequences or a second line.
func validTerminalCommand(cmd string) error {
	if strings.TrimSpace(cmd) == "" {
		return fmt.Errorf("the command is empty")
	}
	if len(cmd) > 4000 {
		return fmt.Errorf("the command is longer than 4000 characters")
	}
	for _, r := range cmd {
		if r == unicode.ReplacementChar || unicode.IsControl(r) || r == ' ' || r == ' ' {
			return fmt.Errorf("terminal_run takes one line without control characters (no newlines, tabs or escape sequences)")
		}
	}
	return nil
}

// ready checks that WRM may type now: the screen (tracker) and, unless quick, the
// foreground process over an exec channel.
func (ts *termShare) ready(quick bool) *termRefusal {
	ts.mu.Lock()
	closed, running, last := ts.closed, ts.running, ts.lastInput
	ts.mu.Unlock()
	if closed {
		return &termRefusal{Rule: "closed", Why: "the terminal window is closed", Hard: true}
	}
	if running {
		return &termRefusal{Rule: "busy", Why: "another terminal_run is still waiting for its command", Hard: true}
	}
	rf := ts.tr.readyToType(last)
	if rf != nil && !(rf.markers && !quick) {
		return rf
	}
	if quick {
		return nil
	}
	known, frf := termForeground(ts.ssh)
	if frf != nil {
		return frf
	}
	if rf != nil {
		// the markers say a command runs, but the shell is in the foreground: they are stale
		if !known {
			return rf
		}
		ts.tr.dropMarkers()
		return ts.tr.readyToType(last)
	}
	return nil
}

// waitReady waits up to d for the terminal to become ready (soft refusals only).
func (ts *termShare) waitReady(ctx context.Context, d time.Duration) *termRefusal {
	deadline := time.Now().Add(d)
	for {
		rf := ts.ready(true)
		if rf == nil || rf.markers {
			if rf = ts.ready(false); rf == nil {
				return nil
			}
		}
		if rf.Hard || time.Now().After(deadline) {
			return rf
		}
		select {
		case <-ctx.Done():
			return &termRefusal{Rule: "stopped", Why: "stopped by the user", Hard: true}
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// run types cmd into the terminal and captures its output until the prompt returns (OSC
// 133 or a quiet prompt), a full-screen program or a password prompt appears, or timeout.
func (ts *termShare) run(ctx context.Context, cmd string, timeout time.Duration, r *http.Request, sid int64) (string, int, error) {
	if err := validTerminalCommand(cmd); err != nil {
		return "", -1, err
	}
	if rf := ts.waitReady(ctx, 3*time.Second); rf != nil {
		ts.audit(r, "ai.terminal_refused", map[string]interface{}{"ai_session": sid, "rule": rf.Rule, "reason": rf.Why, "command": redactText(truncateStr(cmd, 500))})
		return "", -1, fmt.Errorf("Not typed by WRM: %s. Wait until the user is back at the shell prompt, or ask the user", rf.Why)
	}
	ts.mu.Lock()
	if ts.running {
		ts.mu.Unlock()
		return "", -1, &termRefusal{Rule: "busy", Why: "another terminal_run is still waiting for its command", Hard: true}
	}
	ts.running = true
	ts.mu.Unlock()
	defer func() {
		ts.mu.Lock()
		ts.running = false
		ts.mu.Unlock()
	}()
	m := ts.tr.mark()
	if ts.typeFn == nil || !ts.typeFn([]byte(cmd+"\r")) {
		return "", -1, fmt.Errorf("the terminal is closed")
	}
	ts.audit(r, "ai.terminal_typed", map[string]interface{}{"ai_session": sid, "command": redactText(truncateStr(cmd, 1000))})
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()
	finish := func(out []string, lost bool, exit, note string) (string, int, error) {
		text := strings.Join(out, "\n")
		if lost {
			text = "[earlier output scrolled out of WRM's terminal buffer]\n" + text
		}
		code := termExitUnknown
		if exit != "" {
			code, _ = strconv.Atoi(exit)
		}
		if note != "" {
			return text, code, fmt.Errorf("%s", note)
		}
		return text, code, nil
	}
	for {
		select {
		case <-tick.C:
			done, stop, out, lost, exit := ts.tr.commandResult(m, 500*time.Millisecond)
			if done {
				return finish(out, lost, exit, "")
			}
			if stop != "" {
				return finish(out, lost, "", stop)
			}
			ts.mu.Lock()
			closed := ts.closed
			ts.mu.Unlock()
			if closed {
				return finish(out, lost, "", "the terminal window was closed")
			}
		case <-deadline.C:
			_, _, out, lost, _ := ts.tr.commandResult(m, 0)
			return finish(out, lost, "", fmt.Sprintf("the command is still running in the terminal after %s; WRM did not interrupt it (the user sees it and can press Ctrl+C)", timeout))
		case <-ctx.Done():
			_, _, out, lost, _ := ts.tr.commandResult(m, 0)
			return finish(out, lost, "", "stopped waiting (the command keeps running in the terminal)")
		}
	}
}

// ─── the AI session's side ───────────────────────────

func (s *aiSession) boundTerm() *termShare {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.term
}

// setTerm binds or unbinds the session's terminal; the model is told on its next turn.
func (s *aiSession) setTerm(ts *termShare, why string) {
	s.mu.Lock()
	old := s.term
	s.term = ts
	switch {
	case ts != nil && old != ts:
		s.termNotice = "[WRM] The user shared terminal " + ts.PubID + " with you: terminal_read shows it" +
			map[bool]string{true: "; terminal_run types a command into it (the user sees it)", false: ""}[ts.level == "run"] + "."
	case ts == nil && old != nil:
		s.termNotice = "[WRM] The terminal is no longer shared with you (" + nonEmpty(why, "ended") + "); terminal_read and terminal_run do not work any more."
	}
	s.mu.Unlock()
	if old != ts {
		d := map[string]interface{}{"on": ts != nil}
		if ts != nil {
			d["terminal"] = ts.PubID
			d["level"] = ts.level
		} else {
			d["why"] = why
		}
		s.emit("terminal", d)
		s.notifyMCP()
	}
}

// terminalTool checks that a terminal tool may be used now; it returns the share.
func (s *aiSession) terminalTool(name string) (*termShare, string, string) {
	ts := s.boundTerm()
	if ts == nil {
		return nil, "terminal_not_shared", "no terminal is shared with this AI session. The user shares one in WRM (\"Share this terminal with AI\" or \"Connect to AI app…\" in the terminal window); only the user can do that"
	}
	on, level := ts.isOn()
	if !on {
		return nil, "terminal_not_shared", "the user stopped sharing the terminal"
	}
	if s.TokenID > 0 && !ts.hasGrant(s.TokenID) {
		return nil, "terminal_not_shared", "the AI connection is no longer attached to the terminal"
	}
	if name == "terminal_run" && level != "run" {
		return nil, "terminal_policy", "the policy allows this terminal to be read, not typed into"
	}
	return ts, "", ""
}

// toolDefsLocked lists the tools offered to the panel's model: the terminal tools only while
// a terminal is shared with the session (s.mu held).
func (s *aiSession) toolDefsLocked() []aiToolDef {
	defs := aiToolDefsFor(s.Mode)
	ts := s.term
	if ts == nil {
		return defs
	}
	on, level := ts.isOn()
	if !on {
		return defs
	}
	for _, t := range aiTools {
		if t.Terminal && (t.Def.Name == "terminal_read" || level == "run") {
			defs = append(defs, t.Def)
		}
	}
	return defs
}

// ─── HTTP API ────────────────────────────────────────

func aiTerminalsHandler(w http.ResponseWriter, r *http.Request, userID int, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		jsonError(w, "Not found", 404)
		return
	}
	ts := aiTermGet(parts[0])
	if ts == nil || ts.UserID != userID {
		jsonError(w, "Terminal not found or closed", 404)
		return
	}
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		v := ts.view()
		level, why := aiTerminalShareLevel(userID, ts.ConnID)
		v["policy"] = map[string]interface{}{"level": level, "why": why, "lines": aiContextLines()}
		apps := []map[string]interface{}{}
		if ok, _ := mcpUserAllowed(userID); ok && ts.ShareID == 0 {
			for _, g := range mcpListGrants(`user_id=? AND revoked_at='' AND expires_at > ?`, []interface{}{userID, nowRFC()}, 100) {
				gid, _ := g["id"].(int64)
				lg, err := loadGrant(gid)
				if err != nil || !lg.hasConn(ts.ConnID) {
					continue
				}
				g["terminal_scope"] = oneOf("terminal_read", lg.Scopes...) || oneOf("terminal_with_approval", lg.Scopes...)
				g["attached"] = ts.hasGrant(gid)
				apps = append(apps, g)
			}
		}
		v["available_apps"] = apps
		jsonOK(w, v)
		return
	}
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	var in struct {
		On        bool   `json:"on"`
		SessionID int64  `json:"session_id"`
		GrantID   int64  `json:"grant_id"`
		Text      string `json:"text"`
		Ask       string `json:"ask"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&in); err != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	switch parts[1] {
	case "share":
		if err := ts.setShare(r, in.On, "by the user"); err != nil {
			jsonError(w, err.Error(), 403)
			return
		}
	case "bind":
		s := aiGet(in.SessionID)
		if s == nil || s.UserID != userID || s.Transport == "mcp" {
			jsonError(w, "AI session not found", 404)
			return
		}
		if ok, why := aiUserAllowed(userID); !ok {
			jsonError(w, why, 403)
			return
		}
		if err := ts.bindSession(r, s, "panel"); err != nil {
			jsonError(w, err.Error(), 409)
			return
		}
	case "attach":
		g, err := loadGrant(in.GrantID)
		if err != nil {
			jsonError(w, "AI connection not found", 404)
			return
		}
		if err := ts.attachGrant(r, g); err != nil {
			jsonError(w, err.Error(), 403)
			return
		}
	case "detach":
		switch {
		case in.GrantID > 0:
			ts.detachGrant(r, in.GrantID, "detached by the user")
		case in.SessionID > 0:
			s := aiGet(in.SessionID)
			if s == nil || s.UserID != userID {
				jsonError(w, "AI session not found", 404)
				return
			}
			ts.detachSession(r, s, "detached by the user")
		default:
			jsonError(w, "grant_id or session_id is required", 400)
			return
		}
	case "selection":
		if !ts.hasGrant(in.GrantID) {
			jsonError(w, "The AI app is not attached to this terminal", 409)
			return
		}
		text, n, err := aiTerminalSelection(in.Text)
		if err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		ask := redactText(truncateStr(strings.TrimSpace(in.Ask), 1000))
		ts.mu.Lock()
		cnt := 0
		for _, it := range ts.items {
			if it.GrantID == in.GrantID {
				cnt++
			}
		}
		if cnt >= aiTermItemMax {
			ts.mu.Unlock()
			jsonError(w, fmt.Sprintf("%d selections are already waiting for the AI app", aiTermItemMax), 409)
			return
		}
		it := termItem{ID: randomToken(6), GrantID: in.GrantID, Text: text, Ask: ask, Redactions: n, Created: time.Now()}
		ts.items = append(ts.items, it)
		ts.mu.Unlock()
		ts.audit(r, "ai.terminal_selection", map[string]interface{}{"target": "mcp", "mcp_grant": in.GrantID, "chars": len(text), "redactions": n, "item": it.ID})
		ts.push()
		mcpNotifyUser(userID)
	default:
		jsonError(w, "Not found", 404)
		return
	}
	jsonOK(w, ts.view())
}

// aiTerminalSelection checks and redacts a selection from the terminal.
func aiTerminalSelection(text string) (string, int, error) {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	if text == "" {
		return "", 0, fmt.Errorf("the selection is empty")
	}
	if len(text) > aiTermItemMaxBytes {
		return "", 0, fmt.Errorf("the selection is larger than %d KB", aiTermItemMaxBytes>>10)
	}
	red, n := redactSecrets(text)
	return red, n, nil
}

// mcpNotifyUserOf tells the owner of an AI connection that something changed.
func mcpNotifyUserOf(gid int64) {
	var uid int
	if db.QueryRow(`SELECT user_id FROM mcp_tokens WHERE id=?`, gid).Scan(&uid) == nil {
		mcpNotifyUser(uid)
	}
}

// mcpTerminalsInfo lists the attached terminals of an AI connection on a connection, as
// list_connections and open_session show them.
func mcpTerminalsInfo(gid int64, connID int) []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, ts := range aiTermsOfGrant(gid, connID) {
		on, level := ts.isOn()
		if !on {
			continue
		}
		ts.mu.Lock()
		n := 0
		for _, it := range ts.items {
			if it.GrantID == gid {
				n++
			}
		}
		since := ts.since
		ts.mu.Unlock()
		out = append(out, map[string]interface{}{"terminal_id": ts.PubID, "connection_id": ts.ConnID, "shared_since": since.UTC().Format(time.RFC3339),
			"can_type": level == "run", "pending_selections": n})
	}
	return out
}
