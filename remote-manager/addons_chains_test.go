package main

import (
	"fmt"
	"strings"
	"testing"
)

// The chains.read capability: an addon panel reads the chains of whoever opens it.

func addonChainsInsert(t *testing.T, userID int, name, value string) int {
	t.Helper()
	def := chainDef{Name: name,
		Nodes:  []chainNode{{ID: "demo-content", Kind: "server", Label: "Content", Role: "content", ConnID: 4242}},
		Checks: []chainCheck{{ID: "keys-type", Node: "demo-content", Type: "config", Label: "keys", Target: "/etc/demo/storage.ini", Format: "ini", Key: "backend:keys.type"}}}
	st := chainState{At: nowRFC(), Nodes: map[string]chainNodeState{"demo-content": {State: "up"}},
		Checks: map[string]chainResult{"keys-type": {Value: value, OK: true}}}
	res, err := db.Exec(`INSERT INTO chains (user_id, name, definition, enabled, created_at, updated_at, last_at, last_state, confirmed) VALUES (?,?,?,1,?,?,?,?,?)`,
		userID, name, string(jsonMarshal(def)), st.At, st.At, st.At, string(jsonMarshal(st)), string(jsonMarshal(st)))
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	db.Exec(`INSERT INTO chain_samples (chain_id, at, kind, subject, old_value, new_value, ok) VALUES (?,?,?,?,?,?,1)`, id, st.At, "check", "keys-type", "fs", value)
	t.Cleanup(func() {
		db.Exec(`DELETE FROM chain_samples WHERE chain_id=?`, id)
		db.Exec(`DELETE FROM chains WHERE id=?`, id)
	})
	return int(id)
}

func TestAddonChainsRead(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "acr-admin", true)
	alice := newTestUser(t, srv, "acr-alice", false)
	if _, err := setSetting("chains_enabled", "1"); err != nil {
		t.Fatal(err)
	}
	mine := addonChainsInsert(t, alice.userID, "Demo storage", "s3")
	other := addonChainsInsert(t, admin.userID, "Admin chain", "fs")

	panel := map[string]string{"p.html": "<p>x</p>"}
	man := `{"apiVersion":2,"id":%q,"name":"Chains demo","version":"1.0.0","capabilities":["ui.panel"%s],"ui":{"panel":"p.html"}}`
	install := func(id, caps string) int {
		files := map[string]string{"addon.json": fmt.Sprintf(man, id, caps)}
		for k, v := range panel {
			files[k] = v
		}
		if r := addonPost(admin, "/api/addons/install", addonZip(t, files)); r.Code != 200 {
			t.Fatalf("install %s: %d %s", id, r.Code, r.Raw)
		}
		t.Cleanup(func() { addonRemoveCore(admin.userID, id) })
		if r := addonReq(admin, "POST", "/api/addons/"+id+"/enable", ""); r.Code != 200 {
			t.Fatalf("enable %s: %d %s", id, r.Code, r.Raw)
		}
		in := addonReq(admin, "POST", "/api/addons/"+id+"/instances", `{"name":"Wall","shared":true}`)
		return int(in.Body["id"].(float64))
	}
	iid := install("demo.chains", `,"chains.read"`)
	base := fmt.Sprintf("/api/addons/demo.chains/instances/%d/chains", iid)

	// The list holds only the viewer's chains.
	r := addonReq(alice, "GET", base, "")
	if r.Code != 200 || !strings.Contains(r.Raw, "Demo storage") || strings.Contains(r.Raw, "Admin chain") {
		t.Fatalf("alice's list: %d %s", r.Code, r.Raw)
	}
	if r := addonReq(admin, "GET", base, ""); r.Code != 200 || !strings.Contains(r.Raw, "Admin chain") || strings.Contains(r.Raw, "Demo storage") {
		t.Fatalf("admin's list: %d %s", r.Code, r.Raw)
	}

	// One chain: values, the timeline, no connection ids.
	r = addonReq(alice, "GET", fmt.Sprintf("%s/%d", base, mine), "")
	if r.Code != 200 || !strings.Contains(r.Raw, `"value":"s3"`) || !strings.Contains(r.Raw, "backend:keys.type") || !strings.Contains(r.Raw, `"new":"s3"`) {
		t.Fatalf("alice's chain: %d %s", r.Code, r.Raw)
	}
	if strings.Contains(r.Raw, "conn_id") || strings.Contains(r.Raw, "4242") {
		t.Fatalf("the chain view leaks the connection: %s", r.Raw)
	}
	if r := addonReq(alice, "GET", fmt.Sprintf("%s/%d", base, other), ""); r.Code != 404 {
		t.Fatalf("another user's chain: %d %s", r.Code, r.Raw)
	}
	if r := addonReq(alice, "GET", fmt.Sprintf("%s/%d/check", base, mine), ""); r.Code != 405 {
		t.Fatalf("check now with GET: %d", r.Code)
	}
	if r := addonReq(alice, "GET", fmt.Sprintf("%s/%d/other", base, mine), ""); r.Code != 404 {
		t.Fatalf("unknown sub-path: %d", r.Code)
	}

	// Chains turned off by the administrator.
	setSetting("chains_enabled", "0")
	if r := addonReq(alice, "GET", base, ""); r.Code != 403 || r.Body["code"] != "E_CHAINS_OFF" {
		t.Fatalf("chains off: %d %s", r.Code, r.Raw)
	}
	setSetting("chains_enabled", "1")

	// An addon without the capability.
	plain := install("demo.nochains", "")
	if r := addonReq(alice, "GET", fmt.Sprintf("/api/addons/demo.nochains/instances/%d/chains", plain), ""); r.Code != 403 || r.Body["code"] != "E_CAPABILITY" {
		t.Fatalf("without chains.read: %d %s", r.Code, r.Raw)
	}

	// A disabled addon.
	addonReq(admin, "POST", "/api/addons/demo.chains/disable", "")
	if r := addonReq(alice, "GET", base, ""); r.Code != 409 {
		t.Fatalf("disabled addon: %d %s", r.Code, r.Raw)
	}
}

func TestAddonChainsManifest(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "acm-admin", true)
	for name, m := range map[string]string{
		"v1":       `{"id":"acm.v1","name":"x","version":"1","capabilities":["ui.panel","chains.read"],"ui":{"panel":"p.html"}}`,
		"no panel": `{"apiVersion":2,"id":"acm.nopanel","name":"x","version":"1","capabilities":["ui.theme","chains.read"],"ui":{"themes":[{"id":"t","name":"T","mode":"dark","tokens":{"--bg":"#000"}}]}}`,
	} {
		if r := addonPost(admin, "/api/addons/install", addonZip(t, map[string]string{"addon.json": m, "p.html": "<p>x</p>"})); r.Code != 400 {
			t.Errorf("%s: want 400, got %d %s", name, r.Code, r.Raw)
		}
	}
	m := &addonManifest{APIVersion: 2, Capabilities: []string{"ui.panel", "chains.read"}, UI: &addonUI{Panel: "p.html"}}
	perms := addonPermissions(m)
	if strings.Join(perms, " ") != "panel chains" {
		t.Fatalf("permission atoms: %v", perms)
	}
	if a := addonPermsAdded(perms, []string{"panel"}); len(a) != 1 || a[0] != "chains" {
		t.Fatalf("adding chains.read did not ask again: %v", a)
	}
	if !strings.Contains(addonPermissionText("chains"), "chains") {
		t.Fatal("no permission sentence for chains")
	}
	if !strings.Contains(addonLLMsText(), "wrm:chains.get") {
		t.Fatal("llms.txt does not document chains.read")
	}
}
