package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestChainAIDraftAndApprove: the assistant lists what it may use, gets every validation
// error back, saves a draft that is checked once but never polled or notified, may change
// only drafts, and only the user's Approve starts monitoring.
func TestChainAIDraftAndApprove(t *testing.T) {
	e := newAIEnv(t, nil)
	if _, err := setSetting("chains_enabled", "1"); err != nil {
		t.Fatal(err)
	}
	s := e.start("read_only", nil) // drafts are inert, so a read-only session may write them
	call := func(name string, in interface{}) aiToolResult {
		raw, _ := json.Marshal(in)
		return s.CallTool(context.Background(), aiToolCall{ID: "c-" + randomToken(4), Name: name, Input: raw})
	}
	res := call("chain_status", map[string]interface{}{})
	if res.IsError || !strings.Contains(res.Content, `"conn_id": `+strconv.Itoa(e.connID)) || !strings.Contains(res.Content, "check_types") {
		t.Fatalf("chain_status: %s", res.Content)
	}

	other := e.admin.addConnection("demo-foreign", e.host.addr)
	bad := map[string]interface{}{"reason": "test", "definition": map[string]interface{}{
		"name":   "Demo chain",
		"nodes":  []map[string]interface{}{{"id": "api", "kind": "server", "conn_id": e.connID}, {"id": "x", "kind": "server", "conn_id": other}},
		"checks": []map[string]interface{}{{"id": "sh", "node": "api", "type": "shell", "target": "rm -rf /"}},
	}}
	res = call("chain_draft", bad)
	if !res.IsError || !strings.Contains(res.Content, "unknown type") || !strings.Contains(res.Content, "conn_id is not one of your connections") {
		t.Fatalf("an invalid draft: %s", res.Content)
	}
	var n int
	db.QueryRow(`SELECT COUNT(1) FROM chains WHERE user_id=?`, e.user.userID).Scan(&n)
	if n != 0 {
		t.Fatalf("an invalid draft was saved")
	}

	good := map[string]interface{}{"reason": "test", "note": "unsure about the port", "definition": map[string]interface{}{
		"name": "Demo chain", "confirm": 1,
		"nodes": []map[string]interface{}{{"id": "api", "kind": "server", "label": "demo-api", "conn_id": e.connID}, {"id": "edge", "kind": "external"},
			{"id": "web", "kind": "server", "conn_id": e.user.addConnection("demo-web", e.host.addr)}},
		"checks": []map[string]interface{}{
			{"id": "svc", "node": "api", "type": "systemd", "target": "demo.service"},
			{"id": "ssh", "node": "edge", "type": "tcp", "target": e.host.addr},
		},
		"links":     []map[string]interface{}{{"from": "edge", "to": "api", "kind": "dependency"}},
		"scenarios": []map[string]interface{}{{"id": "ok", "when": map[string]interface{}{"check": "ssh", "field": "ok", "op": "==", "value": "true"}}},
	}}
	res = call("chain_draft", good)
	var out struct {
		Saved   bool       `json:"saved"`
		ChainID int        `json:"chain_id"`
		State   chainState `json:"state"`
	}
	body := res.Content[strings.Index(res.Content, "{") : strings.LastIndex(res.Content, "}")+1]
	if res.IsError || json.Unmarshal([]byte(body), &out) != nil || !out.Saved || out.ChainID == 0 {
		t.Fatalf("a valid draft: %s", res.Content)
	}
	// The first look covers only the session's own server: other servers and checks that
	// leave WRM over the network wait for the user.
	if out.State.Nodes["api"].State != "up" || out.State.Nodes["web"].State != "unknown" ||
		out.State.Checks["ssh"].Value != "unknown" || !strings.Contains(out.State.Checks["ssh"].Detail, "Check now") {
		t.Fatalf("the draft's first check: %+v", out.State)
	}
	c, _ := loadChain(out.ChainID)
	if !c.AIDraft || c.Enabled || c.DraftNote != "unsure about the port" {
		t.Fatalf("draft row: %+v", c)
	}
	for _, id := range dueChains() {
		if id == c.ID {
			t.Fatal("a draft is due for background checks")
		}
	}
	var events int
	db.QueryRow(`SELECT COUNT(1) FROM user_events WHERE user_id=? AND event LIKE 'chain.%'`, e.user.userID).Scan(&events)
	if events != 0 {
		t.Fatalf("a draft sent %d notifications", events)
	}
	good["chain_id"] = c.ID
	if res = call("chain_draft", good); res.IsError {
		t.Fatalf("updating the draft: %s", res.Content)
	}

	// Only the owner approves; afterwards the assistant cannot change the chain.
	base := "/api/chains/" + strconv.Itoa(c.ID)
	if code := e.admin.jsonDo("POST", base+"/approve", nil, nil); code != http.StatusNotFound {
		t.Fatalf("another user approved: %d", code)
	}
	var v struct {
		Enabled bool `json:"enabled"`
		AIDraft bool `json:"ai_draft"`
	}
	if code := e.user.jsonDo("POST", base+"/approve", nil, &v); code != http.StatusOK || !v.Enabled || v.AIDraft {
		t.Fatalf("approve: %d %+v", code, v)
	}
	if code := e.user.jsonDo("POST", base+"/approve", nil, nil); code != http.StatusConflict {
		t.Fatalf("approving twice: %d", code)
	}
	if res = call("chain_draft", good); !res.IsError || !strings.Contains(res.Content, "approved") {
		t.Fatalf("the assistant changed an approved chain: %s", res.Content)
	}
	res = call("chain_status", map[string]interface{}{"chain_id": c.ID})
	if res.IsError || !strings.Contains(res.Content, `"ai_draft": false`) {
		t.Fatalf("chain_status of one chain: %s", res.Content)
	}
	found := false
	for _, id := range dueChains() {
		found = found || id == c.ID
	}
	if !found {
		t.Fatal("the approved chain is not checked in the background")
	}

	// Chains off: the tools refuse.
	setSetting("chains_enabled", "0")
	if res = call("chain_status", map[string]interface{}{}); !res.IsError {
		t.Fatalf("chain_status with chains off: %s", res.Content)
	}
	setSetting("chains_enabled", "1")
	e.user.jsonDo("DELETE", base, nil, nil)
}

// TestChainAIGrantConnections: over MCP a draft may use only the AI connection's servers.
func TestChainAIGrantConnections(t *testing.T) {
	e := newAIEnv(t, nil)
	setSetting("chains_enabled", "1")
	second := e.user.addConnection("demo-second", e.host.addr)
	res, err := db.Exec(`INSERT INTO mcp_tokens (user_id, name, token_hash, conn_ids, scopes, created_at, expires_at) VALUES (?,?,?,?,?,?,?)`,
		e.user.userID, "demo", "h-"+randomToken(8), strconv.Itoa(e.connID), "read_logs,edit_file_with_approval", nowRFC(), time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	tok, _ := res.LastInsertId()
	s := &aiSession{UserID: e.user.userID, TokenID: tok, Scopes: []string{"read_logs", "edit_file_with_approval"}}
	allowed, err := s.chainAIAllowedConns()
	if err != nil || !allowed[e.connID] || allowed[second] {
		t.Fatalf("grant connections: %v %v", allowed, err)
	}
	if s.chainAIStorages() {
		t.Fatal("storages allowed without a storage scope")
	}
	out, code, _ := s.chainAIDraft(context.Background(), aiToolInput{Definition: json.RawMessage(`{"name":"x","nodes":[{"id":"a","kind":"server","conn_id":` + strconv.Itoa(second) + `}],"checks":[]}`)})
	if code == 0 || !strings.Contains(out, "may use") {
		t.Fatalf("a server outside the grant was accepted: %d %s", code, out)
	}
	if !strings.Contains(strings.Join(aiScopesFor("chain_draft"), ","), "edit_file_with_approval") || !aiToolInScopes("chain_status", []string{"read_logs"}) {
		t.Fatal("chain tool scopes")
	}
}
