package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func chainAITestDef(storageID int) chainDef {
	return chainDef{
		Name: "Demo bucket chain", IntervalS: 30, Confirm: 1,
		Nodes:  []chainNode{{ID: "bucket", Kind: "storage", Label: "records", Role: "storage", StorageID: storageID}},
		Checks: []chainCheck{{ID: "s3", Node: "bucket", Type: "storage", Label: "Bucket"}},
		Scenarios: []chainScenario{
			{ID: "normal", Label: "Normal", When: chainCond{All: []chainCond{{Check: "s3", Field: "ok", Op: "==", Value: "true"}}}},
			{ID: "degraded", Label: "Degraded", When: chainCond{}},
		},
	}
}

// TestChainAIChangeIsADraftUntilApplied: the assistant reads a chain, tests a change and drafts
// it with chain_id; the live chain stays as it is until the user applies the draft, which keeps
// the chain's id.
func TestChainAIChangeIsADraftUntilApplied(t *testing.T) {
	setSetting("chains_enabled", "1")
	setSetting("chain_storage_check_s", "0")
	t.Cleanup(func() { setSetting("chain_storage_check_s", "300") })
	u, _, storageID := storageTestEnv(t, "chain-ai-edit", nil)
	def := chainAITestDef(storageID)
	var created struct {
		ID int `json:"id"`
	}
	if code := u.jsonDo("POST", "/api/chains", map[string]interface{}{"definition": def, "enabled": true}, &created); code != http.StatusOK || created.ID == 0 {
		t.Fatalf("create chain: %d", code)
	}
	s := &aiSession{ID: 9, UID: "sess-edit", UserID: u.userID, Username: "tester"}

	got, code, err := s.runChainTool(t.Context(), "chain_get", aiToolInput{ChainID: created.ID})
	if err != nil || code != 0 || !strings.Contains(got, `"Demo bucket chain"`) || !strings.Contains(got, `"results"`) {
		t.Fatalf("chain_get: %v %s", err, got)
	}

	changed := def
	changed.Name = "Demo bucket chain v2"
	changed.Checks = append([]chainCheck{}, def.Checks...)
	changed.Checks[0].Label = "Bucket reachable"
	raw, _ := json.Marshal(changed)

	report, code, err := s.runChainTool(t.Context(), "chain_test", aiToolInput{Definition: raw})
	if err != nil || code != 0 || !strings.Contains(report, `"checks_ok": 1`) || !strings.Contains(report, "Normal") {
		t.Fatalf("chain_test: %v %s", err, report)
	}

	out, code, err := s.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: raw, ChainID: created.ID, Reason: "rename"})
	if err != nil || code != 0 || !strings.Contains(out, `"changes_chain_id": `+strconv.Itoa(created.ID)) {
		t.Fatalf("chain_draft with chain_id: %v %s", err, out)
	}
	var draftID int
	db.QueryRow(`SELECT id FROM chain_drafts WHERE user_id=? AND base_chain_id=?`, u.userID, created.ID).Scan(&draftID)
	if draftID == 0 {
		t.Fatal("the change draft was not stored with its base chain")
	}
	if c, _ := loadChain(created.ID); c.Def.Name != def.Name {
		t.Fatalf("the live chain changed before the user applied the draft: %q", c.Def.Name)
	}
	var view map[string]interface{}
	if code := u.jsonDo("GET", "/api/chains/drafts/"+strconv.Itoa(draftID), nil, &view); code != http.StatusOK || int(view["base_chain_id"].(float64)) != created.ID || view["base_chain_name"] != def.Name {
		t.Fatalf("draft view: %d %+v", code, view)
	}

	var applied struct {
		ID         int      `json:"id"`
		Definition chainDef `json:"definition"`
	}
	if code := u.jsonDo("POST", "/api/chains/drafts/"+strconv.Itoa(draftID)+"/approve", nil, &applied); code != http.StatusOK || applied.ID != created.ID || applied.Definition.Name != changed.Name {
		t.Fatalf("apply: %d %+v", code, applied)
	}
	var n int
	db.QueryRow(`SELECT COUNT(1) FROM chains WHERE user_id=?`, u.userID).Scan(&n)
	if n != 1 {
		t.Fatalf("applying a change created another chain: %d chains", n)
	}
	if _, err := loadChainDraft(draftID); err == nil {
		t.Fatal("the applied draft is still there")
	}

	// Another user can neither read nor change the chain through the tools.
	other := newTestUser(t, u.srv, "chain-ai-edit-other", false)
	os := &aiSession{ID: 10, UID: "sess-other", UserID: other.userID, Username: "other"}
	if _, _, err := os.runChainTool(t.Context(), "chain_get", aiToolInput{ChainID: created.ID}); err == nil {
		t.Fatal("another user read the chain")
	}
	if _, _, err := os.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: raw, ChainID: created.ID}); err == nil {
		t.Fatal("another user drafted a change to the chain")
	}

	// The builder's test endpoint runs a definition once without saving it.
	var res map[string]interface{}
	if code := u.jsonDo("POST", "/api/chains/test", map[string]interface{}{"definition": changed}, &res); code != http.StatusOK || res["checks_ok"].(float64) != 1 {
		t.Fatalf("test endpoint: %d %+v", code, res)
	}
}

// TestChainMCPTools: the chains scope gives an AI connection the chain tools without a server
// session, read-only connections cannot draft, and the scope is not in the default policy.
func TestChainMCPTools(t *testing.T) {
	setSetting("chains_enabled", "1")
	if strings.Contains(mcpDefaultScopes, "chains") {
		t.Fatal("the chains scope must be opt-in")
	}
	if !oneOf("chains", mcpScopeOrder...) {
		t.Fatal("chains is not a known scope")
	}
	names := func(g *mcpGrant) string {
		var out []string
		for _, tool := range mcpChainTools(g) {
			out = append(out, tool["name"].(string))
		}
		return strings.Join(out, ",")
	}
	if got := names(&mcpGrant{ID: 1, Mode: aiModeAsk, Scopes: []string{"read_logs"}}); got != "" {
		t.Fatalf("tools without the scope: %s", got)
	}
	if got := names(&mcpGrant{ID: 2, Mode: aiModeAsk, Scopes: []string{"chains"}}); got != "chain_status,chain_get,chain_test,chain_draft" {
		t.Fatalf("tools with the scope: %s", got)
	}
	if got := names(&mcpGrant{ID: 3, Mode: aiModeReadOnly, Scopes: []string{"chains"}}); got != "chain_status,chain_get,chain_test" {
		t.Fatalf("read-only tools: %s", got)
	}
}

// TestChainMCPCall runs chain_status and a refused chain_draft through the MCP dispatcher.
func TestChainMCPCall(t *testing.T) {
	setSetting("chains_enabled", "1")
	u, _, _ := storageTestEnv(t, "chain-mcp-call", nil)
	call := func(g *mcpGrant, name string) string {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/mcp", nil)
		mcpChainCall(rec, req, g, rpcMsg{ID: json.RawMessage(`1`)}, name, map[string]interface{}{})
		return rec.Body.String()
	}
	g := &mcpGrant{ID: 77, UserID: u.userID, Mode: aiModeAsk, Scopes: []string{"chains"}}
	if out := call(g, "chain_status"); !strings.Contains(out, "check_types") || strings.Contains(out, `"isError":true`) {
		t.Fatalf("chain_status over MCP: %s", out)
	}
	ro := &mcpGrant{ID: 78, UserID: u.userID, Mode: aiModeReadOnly, Scopes: []string{"chains"}}
	if out := call(ro, "chain_draft"); !strings.Contains(out, "read-only") {
		t.Fatalf("a read-only connection drafted: %s", out)
	}
	none := &mcpGrant{ID: 79, UserID: u.userID, Mode: aiModeAsk, Scopes: []string{"read_logs"}}
	if out := call(none, "chain_status"); !strings.Contains(out, "no chains scope") {
		t.Fatalf("a connection without the scope: %s", out)
	}
}
