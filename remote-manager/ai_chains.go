package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// ─── CHAINS: AI BUILDER (step 3) ──────────────────────────────────────────────────────────
//
// The AI assistant can describe a chain but never start one. chain_status lists what it may
// use (the user's chains, SSH connections, storages and the check catalog) or the live state
// of one chain; chain_draft validates a definition exactly like the Chains editor, saves it as
// a draft (ai_draft=1, never checked in the background, never notifies) and runs its
// read-only catalog checks once so the user sees real values on the diagram. Only the user
// starts monitoring, with Approve in the Chains window (POST /api/chains/{id}/approve) or by
// saving the draft with background checks on. The assistant can change only drafts.

const chainDraftMaxPerUser = 10

// chainAIAllowedConns is the set of SSH connections a session may put in a chain: the
// connections of its AI connection for an MCP session, every own SSH connection otherwise.
func (s *aiSession) chainAIAllowedConns() (map[int]bool, error) {
	out := map[int]bool{}
	if s.TokenID != 0 {
		g, err := loadGrant(s.TokenID)
		if err != nil {
			return nil, fmt.Errorf("the AI connection is not available")
		}
		for _, id := range g.ConnIDs {
			out[id] = true
		}
		return out, nil
	}
	rows, err := db.Query(`SELECT id FROM connections WHERE user_id=?`, s.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if rows.Scan(&id) == nil {
			out[id] = true
		}
	}
	return out, nil
}

// chainAIStorages reports whether the session may put storages in a chain: always in the
// panel, with a storage scope over MCP.
func (s *aiSession) chainAIStorages() bool {
	return s.Scopes == nil || oneOf("storage_read", s.Scopes...) || oneOf("storage_write", s.Scopes...)
}

func (s *aiSession) planChainTool(name string, in aiToolInput) (*aiPlanned, error) {
	if !settingBool("chains_enabled") {
		return nil, fmt.Errorf("Chains are turned off by the administrator")
	}
	cin := in
	pl := &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Chain: &cin, NoEdit: true}
	switch name {
	case "chain_status":
		if in.ChainID != 0 {
			c, err := loadChain(in.ChainID)
			if err != nil || c.UserID != s.UserID {
				return nil, fmt.Errorf("chain %d is not one of your chains", in.ChainID)
			}
			pl.Display = "status of chain " + c.Def.Name
		} else {
			pl.Display = "list chains, connections, storages and check types"
		}
	case "chain_draft":
		if len(in.Definition) == 0 {
			return nil, fmt.Errorf("give the chain in definition (the same JSON as the Chains editor)")
		}
		if in.ChainID != 0 {
			c, err := loadChain(in.ChainID)
			if err != nil || c.UserID != s.UserID {
				return nil, fmt.Errorf("chain %d is not one of your chains", in.ChainID)
			}
			if !c.AIDraft {
				return nil, fmt.Errorf("chain %d was approved by the user; the assistant changes only drafts. Write a new draft (without chain_id) or tell the user what to change", in.ChainID)
			}
			pl.Display = "update the chain draft " + c.Def.Name
		} else {
			pl.Display = "write a chain draft"
		}
	default:
		return nil, fmt.Errorf("unknown chain tool")
	}
	return pl, nil
}

func (s *aiSession) runChainTool(ctx context.Context, name string, in aiToolInput) (string, int, error) {
	switch name {
	case "chain_status":
		if in.ChainID != 0 {
			return s.chainAIStatusOne(in.ChainID)
		}
		return s.chainAIOverview()
	case "chain_draft":
		return s.chainAIDraft(ctx, in)
	}
	return "", 2, fmt.Errorf("unknown chain tool")
}

func (s *aiSession) chainAIOverview() (string, int, error) {
	allowed, err := s.chainAIAllowedConns()
	if err != nil {
		return "", 2, err
	}
	type conn struct {
		ID   int    `json:"conn_id"`
		Name string `json:"name"`
		Host string `json:"host"`
	}
	conns := []conn{}
	if rows, err := db.Query(`SELECT id, name, host, protocol FROM connections WHERE user_id=? ORDER BY name COLLATE NOCASE`, s.UserID); err == nil {
		for rows.Next() {
			var c conn
			var proto string
			if rows.Scan(&c.ID, &c.Name, &c.Host, &proto) == nil && allowed[c.ID] && isSSHProtocol(Connection{Protocol: proto}) {
				conns = append(conns, c)
			}
		}
		rows.Close()
	}
	storages := []map[string]interface{}{}
	if s.chainAIStorages() {
		for _, st := range storagesForUser(s.UserID) {
			storages = append(storages, map[string]interface{}{"storage_id": st.ID, "name": st.Name, "bucket": st.Bucket, "prefix": st.Prefix})
		}
	}
	chains := []chainSummary{}
	if rows, err := db.Query(`SELECT id FROM chains WHERE user_id=? ORDER BY name COLLATE NOCASE`, s.UserID); err == nil {
		var ids []int
		for rows.Next() {
			var id int
			if rows.Scan(&id) == nil {
				ids = append(ids, id)
			}
		}
		rows.Close()
		for _, id := range ids {
			if c, err := loadChain(id); err == nil {
				chains = append(chains, summarizeChain(c))
			}
		}
	}
	types := make([]string, 0, len(chainCheckTypes))
	for k := range chainCheckTypes {
		types = append(types, k)
	}
	sort.Strings(types)
	catalog := []map[string]string{}
	for _, k := range types {
		t := chainCheckTypes[k]
		catalog = append(catalog, map[string]string{"type": k, "runs": t.Where, "needs": t.Needs, "help": t.Help})
	}
	b, _ := json.MarshalIndent(map[string]interface{}{
		"chains": chains, "ssh_connections": conns, "storages": storages, "check_types": catalog,
		"link_kinds": []string{"data", "read", "write", "failover", "dependency"}, "ops": []string{"==", "!=", "in", "contains", "exists"},
		"formats":          []string{"ini", "json", "yaml", "env", "kv", "regex"},
		"definition_shape": chainAIShape,
	}, "", "  ")
	return string(b), 0, nil
}

// chainAIShape tells the model how a definition looks, with invented names.
var chainAIShape = map[string]interface{}{
	"name": "Media chain", "interval_s": 30, "confirm": 2,
	"nodes": []map[string]interface{}{
		{"id": "ingest-a", "kind": "server", "label": "demo-ingest-a", "role": "ingest (primary)", "conn_id": 1},
		{"id": "bucket", "kind": "storage", "label": "Recordings", "storage_id": 1},
		{"id": "cdn", "kind": "external", "label": "CDN"},
	},
	"checks": []map[string]interface{}{
		{"id": "ingest-svc", "node": "ingest-a", "type": "systemd", "target": "demo-ingest.service", "label": "Ingest service"},
		{"id": "mode", "node": "ingest-a", "type": "config", "target": "/etc/demo/app.ini", "format": "ini", "key": "storage.mode", "expect": map[string]string{"op": "==", "value": "s3"}},
		{"id": "cdn-health", "node": "cdn", "type": "http", "target": "https://cdn.example.com/health"},
	},
	"links": []map[string]interface{}{
		{"from": "ingest-a", "to": "bucket", "kind": "write", "check": "ingest-svc", "scenarios": []string{"primary"}},
	},
	"scenarios": []map[string]interface{}{
		{"id": "primary", "label": "Primary", "when": map[string]interface{}{"all": []map[string]string{{"node": "ingest-a", "op": "==", "value": "up"}, {"check": "ingest-svc", "field": "ok", "op": "==", "value": "true"}}}},
		{"id": "degraded", "label": "Degraded", "when": map[string]interface{}{}},
	},
}

func (s *aiSession) chainAIStatusOne(id int) (string, int, error) {
	c, err := loadChain(id)
	if err != nil || c.UserID != s.UserID {
		return "", 2, fmt.Errorf("chain %d is not one of your chains", id)
	}
	out := map[string]interface{}{"id": c.ID, "ai_draft": c.AIDraft, "enabled": c.Enabled, "definition": c.Def, "state": c.State,
		"scenario": c.Confirmed.Scenario, "last_at": c.LastAt}
	var changes []map[string]interface{}
	if rows, err := db.Query(`SELECT at, kind, subject, old_value, new_value FROM chain_samples WHERE chain_id=? ORDER BY id DESC LIMIT 30`, c.ID); err == nil {
		for rows.Next() {
			var at, kind, subject, o, n string
			if rows.Scan(&at, &kind, &subject, &o, &n) == nil {
				changes = append(changes, map[string]interface{}{"at": at, "kind": kind, "subject": subject, "old": o, "new": n})
			}
		}
		rows.Close()
	}
	out["recent_changes"] = changes
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b), 0, nil
}

// chainAIDraft validates and saves a draft, then checks it once.
func (s *aiSession) chainAIDraft(ctx context.Context, in aiToolInput) (string, int, error) {
	var d chainDef
	dec := json.NewDecoder(strings.NewReader(string(in.Definition)))
	dec.DisallowUnknownFields()
	if len(in.Definition) > chainMaxDefBytes {
		return "", 1, fmt.Errorf("the definition is larger than %d bytes", chainMaxDefBytes)
	}
	if err := dec.Decode(&d); err != nil {
		return "", 1, fmt.Errorf("the definition is not a valid chain: %v", err)
	}
	errs := validateChain(&d, s.UserID)
	allowed, err := s.chainAIAllowedConns()
	if err != nil {
		return "", 2, err
	}
	for i, n := range d.Nodes {
		if n.Kind == "server" && !allowed[n.ConnID] {
			errs = append(errs, fmt.Sprintf("nodes[%d] (%s): conn_id %d is not one of the connections this AI connection may use", i, n.ID, n.ConnID))
		}
		if n.Kind == "storage" && !s.chainAIStorages() {
			errs = append(errs, fmt.Sprintf("nodes[%d] (%s): this AI connection has no storage scope", i, n.ID))
		}
	}
	if len(errs) > 0 {
		b, _ := json.MarshalIndent(map[string]interface{}{"saved": false, "errors": errs,
			"hint": "Fix every error and call chain_draft again. Use chain_status for the connection and storage ids and the check types."}, "", "  ")
		return string(b), 1, nil
	}
	note := strings.TrimSpace(in.Note)
	if chainHasControl(strings.ReplaceAll(note, "\n", " ")) {
		note = ""
	}
	note = truncateStr(note, 1000)
	now := nowRFC()
	id := in.ChainID
	if id != 0 {
		c, err := loadChain(id)
		if err != nil || c.UserID != s.UserID || !c.AIDraft {
			return "", 2, fmt.Errorf("chain %d is not a draft of yours", id)
		}
		if _, err := db.Exec(`UPDATE chains SET name=?, definition=?, draft_note=?, updated_at=? WHERE id=? AND user_id=? AND ai_draft=1`,
			d.Name, string(jsonMarshal(d)), note, now, id, s.UserID); err != nil {
			return "", 2, fmt.Errorf("the draft was not saved")
		}
		chainRuns.Lock()
		delete(chainRuns.pending, id)
		chainRuns.Unlock()
		s.audit(nil, "chain.draft_updated", map[string]interface{}{"chain": id, "nodes": len(d.Nodes), "checks": len(d.Checks)})
	} else {
		var n, drafts int
		db.QueryRow(`SELECT COUNT(1), COALESCE(SUM(ai_draft),0) FROM chains WHERE user_id=?`, s.UserID).Scan(&n, &drafts)
		if n >= 50 {
			return "", 2, fmt.Errorf("the user already has 50 chains")
		}
		if drafts >= chainDraftMaxPerUser {
			return "", 2, fmt.Errorf("there are already %d unapproved drafts; update one with chain_id or ask the user to approve or delete them", drafts)
		}
		res, err := db.Exec(`INSERT INTO chains (user_id, name, definition, enabled, created_at, updated_at, ai_draft, draft_note) VALUES (?,?,?,0,?,?,1,?)`,
			s.UserID, d.Name, string(jsonMarshal(d)), now, now, note)
		if err != nil {
			return "", 2, fmt.Errorf("the draft was not saved")
		}
		id64, _ := res.LastInsertId()
		id = int(id64)
		s.audit(nil, "chain.draft_created", map[string]interface{}{"chain": id, "nodes": len(d.Nodes), "checks": len(d.Checks)})
	}
	// The first look shows only what the assistant could read anyway: the servers of this
	// session (the panel's own server, or the AI connection's servers) and storages; checks
	// that leave WRM over the network wait for the user's Check now or approval.
	readable := map[int]bool{}
	if s.TokenID != 0 {
		readable = allowed
	} else {
		readable[s.ConnID] = true
	}
	ctx = context.WithValue(ctx, chainLimitKey{}, &chainLimit{Conns: readable})
	pctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	got, err := pollChain(pctx, id)
	out := map[string]interface{}{"saved": true, "chain_id": id, "ai_draft": true,
		"next": "The draft is not monitored yet. Tell the user to open Chains (top bar), review it on the diagram with the values below and press Approve. Adjust it with chain_draft and chain_id if they ask."}
	if err != nil {
		out["check_error"] = err.Error()
	} else {
		out["state"] = got.State
		out["scenario"] = got.State.Scenario
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b), 0, nil
}

// apiChainApprove starts monitoring a draft: it becomes an ordinary enabled chain, and its
// first background check reports what is already down.
func apiChainApprove(w http.ResponseWriter, r *http.Request, c chainRow, userID int) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !c.AIDraft {
		jsonError(w, "This chain is not a draft", http.StatusConflict)
		return
	}
	if errs := validateChain(&c.Def, userID); len(errs) > 0 {
		chainErrors(w, errs)
		return
	}
	if _, err := db.Exec(`UPDATE chains SET ai_draft=0, enabled=1, draft_note='', confirmed='', last_at='', updated_at=? WHERE id=? AND user_id=?`,
		nowRFC(), c.ID, userID); err != nil {
		jsonError(w, "Server error", http.StatusInternalServerError)
		return
	}
	chainRuns.Lock()
	delete(chainRuns.pending, c.ID)
	chainRuns.Unlock()
	auditLog(r, userID, "chain.approved", c.Def.Name, map[string]interface{}{"chain": c.ID, "nodes": len(c.Def.Nodes), "checks": len(c.Def.Checks)})
	got, _ := loadChain(c.ID)
	hub.sendTo(userID, []byte(fmt.Sprintf(`{"type":"chain_changed","id":%d}`, c.ID)))
	jsonOK(w, chainView(got))
}
