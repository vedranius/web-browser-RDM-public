package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func chainDraftTestDef(connID int) chainDef {
	def := chainDef{
		Name: "Demo AI chain", IntervalS: 30, Confirm: 2,
		Nodes: []chainNode{
			{ID: "app", Kind: "server", Label: "demo-app", Role: "app", ConnID: connID},
			{ID: "site", Kind: "external", Label: "Public site", Role: "endpoint"},
		},
		Checks: []chainCheck{
			{ID: "svc", Node: "app", Type: "systemd", Target: "demo.service", Label: "Demo service"},
			{ID: "cert", Node: "site", Type: "tls", Target: "example.com", Label: "Certificate"},
		},
		Links: []chainLink{{From: "app", To: "site", Kind: "data", Check: "svc"}},
		Scenarios: []chainScenario{
			{ID: "normal", Label: "Normal", When: chainCond{All: []chainCond{{Check: "svc", Field: "ok", Op: "==", Value: "true"}}}},
			{ID: "degraded", Label: "Degraded", When: chainCond{}},
		},
	}
	return def
}

// TestChainDraftAPI: a draft can be read and listed by its owner, is not a chain, and only the
// approve call validates it and starts monitoring; another user sees nothing.
func TestChainDraftAPI(t *testing.T) {
	u, _, _ := storageTestEnv(t, "chain-draft-api", nil)
	conn := chainTestConn(t, u.userID, "SSH")
	def := chainDraftTestDef(conn)
	id, err := saveChainDraft(u.userID, "sess-demo", def)
	if err != nil || id == 0 {
		t.Fatalf("save draft: %v", err)
	}

	var list []chainDraftSummary
	if code := u.jsonDo("GET", "/api/chains/drafts", nil, &list); code != http.StatusOK || len(list) != 1 || list[0].Name != def.Name || list[0].Nodes != 2 {
		t.Fatalf("draft list: %d %+v", code, list)
	}
	var got struct {
		ID         int      `json:"id"`
		Definition chainDef `json:"definition"`
	}
	if code := u.jsonDo("GET", "/api/chains/drafts/"+strconv.Itoa(id), nil, &got); code != http.StatusOK || got.Definition.Name != def.Name || len(got.Definition.Checks) != 2 {
		t.Fatalf("draft read: %d %+v", code, got)
	}
	// A draft is not a chain yet.
	var chains []chainSummary
	u.jsonDo("GET", "/api/chains", nil, &chains)
	if len(chains) != 0 {
		t.Fatalf("a draft appeared as a chain: %+v", chains)
	}

	other := newTestUser(t, u.srv, "chain-draft-api-other", false)
	if code := other.jsonDo("GET", "/api/chains/drafts/"+strconv.Itoa(id), nil, nil); code != http.StatusNotFound {
		t.Fatalf("another user's draft: %d", code)
	}
	if code := other.jsonDo("POST", "/api/chains/drafts/"+strconv.Itoa(id)+"/approve", nil, nil); code != http.StatusNotFound {
		t.Fatalf("another user approved a draft: %d", code)
	}

	var chain struct {
		ID         int      `json:"id"`
		Enabled    bool     `json:"enabled"`
		Definition chainDef `json:"definition"`
	}
	if code := u.jsonDo("POST", "/api/chains/drafts/"+strconv.Itoa(id)+"/approve", nil, &chain); code != http.StatusOK || chain.ID == 0 || !chain.Enabled {
		t.Fatalf("approve: %d %+v", code, chain)
	}
	if _, err := loadChainDraft(id); err == nil {
		t.Fatal("the draft was not removed after approval")
	}
	u.jsonDo("GET", "/api/chains", nil, &chains)
	if len(chains) != 1 || chains[0].Name != def.Name {
		t.Fatalf("the approved chain is missing: %+v", chains)
	}
	var events int
	db.QueryRow(`SELECT COUNT(1) FROM audit_log WHERE user_id=? AND action='chain.created' AND details LIKE ?`, u.userID, "%from_draft%").Scan(&events)
	if events == 0 {
		t.Fatal("approving a draft was not audited")
	}
}

// TestChainDraftApproveRefusesInvalid: an invalid draft cannot become a chain, and discarding
// it removes it.
func TestChainDraftApproveRefusesInvalid(t *testing.T) {
	u, _, _ := storageTestEnv(t, "chain-draft-invalid", nil)
	def := chainDraftTestDef(999999) // no such connection
	id, err := saveChainDraft(u.userID, "sess-bad", def)
	if err != nil {
		t.Fatal(err)
	}
	if code := u.jsonDo("POST", "/api/chains/drafts/"+strconv.Itoa(id)+"/approve", nil, nil); code != http.StatusBadRequest {
		t.Fatalf("invalid draft approved: %d", code)
	}
	if _, err := loadChainDraft(id); err != nil {
		t.Fatal("an invalid draft was removed instead of kept for fixing")
	}
	if code := u.jsonDo("DELETE", "/api/chains/drafts/"+strconv.Itoa(id), nil, nil); code != http.StatusOK {
		t.Fatalf("discard: %d", code)
	}
	if _, err := loadChainDraft(id); err == nil {
		t.Fatal("the discarded draft is still there")
	}
}

// TestChainDraftTools: the built-in assistant gets chain_status and chain_draft, they act on
// the signed-in user (not only administrators), a draft replaces the session's previous one,
// and a definition that would use a foreign connection is refused with the reason.
func TestChainDraftTools(t *testing.T) {
	u, _, storageID := storageTestEnv(t, "chain-draft-tools", nil)
	conn := chainTestConn(t, u.userID, "SSH")
	s := &aiSession{ID: 7, UID: "sess-tools", UserID: u.userID, Username: "tester"}

	// The tools are part of the built-in assistant for every user, never an MCP scope.
	for _, name := range []string{"chain_status", "chain_draft"} {
		if !aiToolInScopes(name, nil) || aiToolInScopes(name, []string{"read_logs", "addons.manage"}) {
			t.Errorf("%s scope rule", name)
		}
	}
	defs := aiToolDefsForUser(aiModeAuto, false)
	found := map[string]bool{}
	for _, d := range defs {
		found[d.Name] = true
	}
	if !found["chain_status"] || !found["chain_draft"] {
		t.Fatalf("a non-administrator does not get the chain tools: %+v", found)
	}
	if found["addon_create"] {
		t.Fatal("a non-administrator got the addon tools")
	}

	out, code, err := s.runChainTool(t.Context(), "chain_status", aiToolInput{})
	if err != nil || code != 0 {
		t.Fatalf("chain_status: %v %d", err, code)
	}
	for _, want := range []string{`"connections"`, `"storages"`, `"check_types"`, `"storage_id": ` + strconv.Itoa(storageID), `"conn_id": ` + strconv.Itoa(conn)} {
		if !strings.Contains(out, want) {
			t.Errorf("chain_status does not contain %s:\n%s", want, out)
		}
	}

	raw, _ := json.Marshal(chainDraftTestDef(conn))
	pl, err := s.planChainTool("chain_draft", aiToolInput{Definition: raw, Reason: "demo"})
	if err != nil || pl == nil || pl.Host == nil || !strings.Contains(pl.Display, "Demo AI chain") {
		t.Fatalf("plan chain_draft: %+v %v", pl, err)
	}
	first, code, err := s.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: raw, Reason: "demo"})
	if err != nil || code != 0 || !strings.Contains(first, `"draft_id"`) {
		t.Fatalf("run chain_draft: %v %d %s", err, code, first)
	}
	var draftID int
	db.QueryRow(`SELECT id FROM chain_drafts WHERE user_id=? AND session_id=?`, u.userID, s.UID).Scan(&draftID)
	if draftID == 0 {
		t.Fatal("the draft was not stored")
	}
	// A second call from the same session replaces it instead of piling drafts up.
	def2 := chainDraftTestDef(conn)
	def2.Name = "Demo AI chain v2"
	raw2, _ := json.Marshal(def2)
	second, code, err := s.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: raw2, Reason: "demo"})
	if err != nil || code != 0 || !strings.Contains(second, "v2") {
		t.Fatalf("replace: %v %d %s", err, code, second)
	}
	var n int
	db.QueryRow(`SELECT COUNT(1) FROM chain_drafts WHERE user_id=? AND session_id=?`, u.userID, s.UID).Scan(&n)
	if n != 1 {
		t.Fatalf("drafts of one session: %d", n)
	}

	// A foreign connection is refused with a message the model can act on.
	other := newTestUser(t, u.srv, "chain-draft-tools-other", false)
	foreign := chainDraftTestDef(chainTestConn(t, other.userID, "SSH"))
	raw3, _ := json.Marshal(foreign)
	if _, _, err := s.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: raw3}); err == nil || !strings.Contains(err.Error(), "conn_id") {
		t.Fatalf("a foreign connection was accepted: %v", err)
	}
	// A definition with an unknown field is refused, not silently cleaned.
	bad := json.RawMessage(`{"name":"x","nodes":[{"id":"a","kind":"external","conn_id":` + strconv.Itoa(conn) + `}],"commands":["rm -rf /"]}`)
	if _, _, err := s.runChainTool(t.Context(), "chain_draft", aiToolInput{Definition: bad}); err == nil {
		t.Fatal("an unknown definition field was accepted")
	}
}
