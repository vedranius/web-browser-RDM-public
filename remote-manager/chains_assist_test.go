package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// TestChainStarterFromSelection: the deterministic wizard result uses the selected servers and
// storages, gives the first server the app role, the second the backup role, and validates.
func TestChainStarterFromSelection(t *testing.T) {
	u, _, storageID := storageTestEnv(t, "chain-starter", nil)
	c1 := chainTestConn(t, u.userID, "SSH")
	c2 := chainTestConn(t, u.userID, "SSH")
	var conns []Connection
	for _, id := range []int{c1, c2} {
		c, err := loadConnectionRaw(id)
		if err != nil {
			t.Fatal(err)
		}
		conns = append(conns, c)
	}
	s, err := loadStorage(storageID, u.userID)
	if err != nil {
		t.Fatal(err)
	}
	def := chainStarterFromSelection(u.userID, conns, []Storage{s}, nil, "Demo system", "two servers and a bucket")
	if errs := validateChain(&def, u.userID); len(errs) > 0 {
		t.Fatalf("starter does not validate:\n%s", strings.Join(errs, "\n"))
	}
	if len(def.Nodes) != 3 || def.Nodes[0].Kind != "server" || def.Nodes[1].Role != "backup" || def.Nodes[2].Kind != "storage" {
		t.Fatalf("starter nodes: %+v", def.Nodes)
	}
	if def.Nodes[0].ID == def.Nodes[1].ID {
		t.Fatal("node ids are not unique")
	}
	if len(def.Links) != 2 || len(def.Scenarios) != 3 {
		t.Fatalf("starter links/scenarios: %+v %+v", def.Links, def.Scenarios)
	}
	if def.Scenarios[len(def.Scenarios)-1].When.empty() == false {
		t.Fatal("the starter needs a fallback scenario")
	}
}

// TestChainDiscoveryParse: the fixed discovery output keeps services, containers, ports and
// network/object-storage mounts, drops noise, SSH itself and local disks, and never leaks text.
func TestChainDiscoveryParse(t *testing.T) {
	out := strings.Join([]string{
		"@host demo-app",
		"@os Linux 6.1 x86_64",
		"@service nginx.service",
		"@service ssh.service",
		"@service systemd-journald.service",
		"@service demo-api.service",
		"@container demo-worker",
		"@port 22",
		"@port 443",
		"@port 8080",
		"@mount " + strings.TrimPrefix(b64("/mnt/records  mybucket:/records  fuse.rclone  rw,nosuid"), "b64:"),
		"@mount " + strings.TrimPrefix(b64("/mnt/data  /dev/sdb1  ext4  rw"), "b64:"),
		"@mount " + strings.TrimPrefix(b64("/srv/nfs  s3.example.com:/bucket  nfs4  rw"), "b64:"),
	}, "\n")
	f := parseChainDiscovery(out, 7, "demo-app")
	if f.ConnID != 7 || f.Hostname != "demo-app" || !strings.Contains(f.OS, "Linux") {
		t.Fatalf("facts: %+v", f)
	}
	if len(f.Services) != 2 || f.Services[0] != "nginx.service" || f.Services[1] != "demo-api.service" {
		t.Fatalf("services: %+v", f.Services)
	}
	if len(f.Containers) != 1 || f.Containers[0] != "demo-worker" {
		t.Fatalf("containers: %+v", f.Containers)
	}
	if len(f.Ports) != 2 || f.Ports[0] != 443 || f.Ports[1] != 8080 {
		t.Fatalf("ports: %+v", f.Ports)
	}
	if len(f.Mounts) != 2 || f.Mounts[0].Target != "/mnt/records" || f.Mounts[1].Source != "s3.example.com:/bucket" {
		t.Fatalf("mounts: %+v", f.Mounts)
	}
	if host, ok := chainMountHost("s3.example.com:/bucket"); !ok || host != "s3.example.com" {
		t.Fatalf("mount host: %q %v", host, ok)
	}
	if _, ok := chainMountHost("mybucket:/records"); ok {
		t.Fatal("a bare rclone remote name became a host")
	}
}

// TestChainStarterUsesDiscovery: without any storage, the starter checks the services, ports
// and mounts WRM found, and adds the mount's host as an outside endpoint (no credentials).
func TestChainStarterUsesDiscovery(t *testing.T) {
	u, _, _ := storageTestEnv(t, "chain-starter-facts", nil)
	conn := chainTestConn(t, u.userID, "SSH")
	c, err := loadConnectionRaw(conn)
	if err != nil {
		t.Fatal(err)
	}
	facts := []chainServerFacts{{
		ConnID: conn, Name: c.Name, OS: "Linux 6.1",
		Services: []string{"nginx.service", "demo-api.service"},
		Ports:    []int{22, 443, 8080},
		Mounts:   []chainMountFact{{Target: "/mnt/records", Source: "s3.example.com:/bucket", FSType: "nfs4"}},
	}}
	def := chainStarterFromSelection(u.userID, []Connection{c}, nil, facts, "Demo", "")
	if errs := validateChain(&def, u.userID); len(errs) > 0 {
		t.Fatalf("starter does not validate:\n%s", strings.Join(errs, "\n"))
	}
	types := map[string]int{}
	for _, x := range def.Checks {
		types[x.Type]++
	}
	if types["systemd"] != 2 || types["port"] < 1 || types["mount"] != 1 || types["dns"] != 1 {
		t.Fatalf("checks do not follow the discovered facts: %+v", types)
	}
	if len(def.Nodes) != 2 || def.Nodes[1].Kind != "external" || def.Nodes[1].Label != "s3.example.com" {
		t.Fatalf("external node for the mount host: %+v", def.Nodes)
	}
	if len(def.Scenarios) < 2 || def.Scenarios[len(def.Scenarios)-1].When.empty() == false {
		t.Fatalf("scenarios: %+v", def.Scenarios)
	}
}

// TestChainStarterManyServers: four servers with many discovered facts still validate; check
// ids stay inside 40 characters and the "all OK" rule is nested instead of exceeding 20 items.
func TestChainStarterManyServers(t *testing.T) {
	u, _, _ := storageTestEnv(t, "chain-many-servers", nil)
	var conns []Connection
	var facts []chainServerFacts
	for i := 0; i < 4; i++ {
		id := chainTestConn(t, u.userID, "SSH")
		c, err := loadConnectionRaw(id)
		if err != nil {
			t.Fatal(err)
		}
		conns = append(conns, c)
		f := chainServerFacts{ConnID: id, Name: c.Name, OS: "Linux", Ports: []int{22, 443, 8080, 9000, 5432}}
		for j := 0; j < 8; j++ {
			f.Services = append(f.Services, fmt.Sprintf("very-long-service-name-that-should-be-truncated-%d-%d.service", i, j))
		}
		for j := 0; j < 3; j++ {
			f.Containers = append(f.Containers, fmt.Sprintf("demo-container-%d-%d", i, j))
		}
		for j := 0; j < 3; j++ {
			f.Mounts = append(f.Mounts, chainMountFact{Target: fmt.Sprintf("/mnt/storage-%d-%d", i, j), Source: "s3.example.com:/bucket", FSType: "nfs4"})
		}
		facts = append(facts, f)
	}
	def := chainStarterFromSelection(u.userID, conns, nil, facts, "Four servers", "")
	if errs := validateChain(&def, u.userID); len(errs) > 0 {
		t.Fatalf("starter for four servers does not validate:\n%s", strings.Join(errs, "\n"))
	}
	if len(def.Checks) <= chainMaxCondItems {
		t.Fatalf("test did not reach the rule limit: %d checks", len(def.Checks))
	}
	for _, x := range def.Checks {
		if !chainIDPattern.MatchString(x.ID) || len(x.ID) > 40 {
			t.Errorf("invalid check id %q (%d chars)", x.ID, len(x.ID))
		}
	}
	normal := def.Scenarios[0]
	if len(normal.When.All) > chainMaxCondItems {
		t.Fatalf("normal rule has %d top-level items", len(normal.When.All))
	}
	// The nested groups still contain only checks, none lost.
	var flat []string
	for _, g := range normal.When.All {
		if g.Check != "" {
			flat = append(flat, g.Check)
		}
		for _, x := range g.All {
			flat = append(flat, x.Check)
		}
	}
	if len(flat) != len(def.Checks) {
		t.Fatalf("normal rule covers %d of %d checks", len(flat), len(def.Checks))
	}
}

// TestUniqueCheckIDBounded: repeated collisions never produce an id above 40 characters.
func TestUniqueCheckIDBounded(t *testing.T) {
	used := map[string]bool{}
	base := "a-very-long-node-and-service-name-that-will-collide-every-time"
	for i := 0; i < 2000; i++ {
		id := uniqueCheckID(used, base)
		if !chainIDPattern.MatchString(id) || len(id) > 40 {
			t.Fatalf("id %q is invalid (%d chars)", id, len(id))
		}
	}
}

// TestChainAssistFallback: without an AI provider the wizard still returns a valid draft built
// from the selection, and it refuses servers/storages the caller does not own.
func TestChainAssistFallback(t *testing.T) {
	u, _, storageID := storageTestEnv(t, "chain-assist-fallback", nil)
	c1 := chainTestConn(t, u.userID, "SSH")
	c2 := chainTestConn(t, u.userID, "SSH")
	var got struct {
		DraftID    int      `json:"draft_id"`
		AI         bool     `json:"ai"`
		Reason     string   `json:"reason"`
		Definition chainDef `json:"definition"`
	}
	body := map[string]interface{}{
		"name": "Demo ingest", "description": "ingest na dva servera i bucket",
		"conn_ids": []int{c1, c2}, "storage_ids": []int{storageID},
	}
	if code := u.jsonDo("POST", "/api/chains/ai-draft", body, &got); code != http.StatusOK || got.DraftID == 0 || got.AI {
		t.Fatalf("fallback: %d %+v", code, got)
	}
	if got.Reason == "" || len(got.Definition.Nodes) != 3 {
		t.Fatalf("fallback reason/definition: %q %+v", got.Reason, got.Definition)
	}
	d, err := loadChainDraft(got.DraftID)
	if err != nil || d.Definition.Name != "Demo ingest" || len(d.Definition.Checks) == 0 {
		t.Fatalf("stored draft: %+v %v", d, err)
	}
	if errs := validateChain(&d.Definition, u.userID); len(errs) > 0 {
		t.Fatalf("stored draft does not validate:\n%s", strings.Join(errs, "\n"))
	}

	other := newTestUser(t, u.srv, "chain-assist-other", false)
	foreign := chainTestConn(t, other.userID, "SSH")
	if code := u.jsonDo("POST", "/api/chains/ai-draft", map[string]interface{}{"conn_ids": []int{foreign}}, nil); code != http.StatusBadRequest {
		t.Fatalf("foreign connection: %d", code)
	}
	if code := u.jsonDo("POST", "/api/chains/ai-draft", map[string]interface{}{"conn_ids": []int{}}, nil); code != http.StatusBadRequest {
		t.Fatalf("empty selection: %d", code)
	}
	if code := u.jsonDo("GET", "/api/chains/ai-draft", nil, nil); code != http.StatusMethodNotAllowed {
		t.Fatalf("GET ai-draft: %d", code)
	}
}

// TestChainAssistAI: the wizard uses the configured provider, parses the JSON, repairs one
// validation error and stores the result as a draft.
func TestChainAssistAI(t *testing.T) {
	var connID int
	env := newAIEnv(t, func(n int, body map[string]interface{}) fakeTurn {
		def := fmt.Sprintf(`{"name":"AI demo","description":"two paths","nodes":[{"id":"app","kind":"server","label":"demo-api","role":"app","conn_id":%d}],"checks":[{"id":"svc","node":"app","type":"systemd","target":"demo.service","label":"Demo"}],"links":[],"scenarios":[{"id":"normal","label":"Normal","when":{"all":[{"check":"svc","field":"ok","op":"==","value":"true"}]}},{"id":"degraded","label":"Degraded","when":{}}]}`, connID)
		if n == 0 {
			// Valid JSON but a check that validateChain refuses (no target).
			return fakeTurn{Text: `{"name":"AI demo","nodes":[{"id":"app","kind":"server","label":"demo-api","role":"app","conn_id":` + strconv.Itoa(connID) + `}],"checks":[{"id":"svc","node":"app","type":"systemd"}],"scenarios":[{"id":"degraded","when":{}}]}`}
		}
		return fakeTurn{Text: "Here is the chain:\n```json\n" + def + "\n```"}
	})
	connID = env.connID
	// Other tests leave org providers pointing at closed fake servers; point this user at the
	// provider this env created instead of relying on list order.
	prefs, _ := json.Marshal(aiUserPrefs{ProviderID: env.provID, Model: "claude-opus-5-5"})
	db.Exec(`DELETE FROM user_ui WHERE user_id=? AND key='ai'`, env.user.userID)
	if _, err := db.Exec(`INSERT INTO user_ui (user_id, key, value, updated_at) VALUES (?, 'ai', ?, ?)`,
		env.user.userID, string(prefs), nowRFC()); err != nil {
		t.Fatal(err)
	}
	var got struct {
		DraftID    int      `json:"draft_id"`
		AI         bool     `json:"ai"`
		Reason     string   `json:"reason"`
		Definition chainDef `json:"definition"`
	}
	code, v := env.user.aiDo("POST", "/api/chains/ai-draft", map[string]interface{}{
		"name": "AI demo", "description": "jedan server", "conn_ids": []int{env.connID},
	})
	if code != http.StatusOK || asMap(v)["ai"] != true {
		t.Fatalf("ai draft: %d %v", code, v)
	}
	got.DraftID = int(asMap(v)["draft_id"].(float64))
	got.AI = asMap(v)["ai"].(bool)
	if !got.AI || got.DraftID == 0 || env.llm.requests() < 2 {
		t.Fatalf("expected a repaired AI draft and two model calls, got %d calls, %+v", env.llm.requests(), v)
	}
	d, err := loadChainDraft(got.DraftID)
	if err != nil || d.Definition.Name != "AI demo" || len(d.Definition.Checks) != 1 || d.Definition.Checks[0].Target != "demo.service" {
		t.Fatalf("stored AI draft: %+v %v", d, err)
	}
	if errs := validateChain(&d.Definition, env.user.userID); len(errs) > 0 {
		t.Fatalf("AI draft does not validate:\n%s", strings.Join(errs, "\n"))
	}
}

// TestExtractJSONObject: fences, prose and braces inside strings do not confuse the parser.
func TestExtractJSONObject(t *testing.T) {
	cases := []string{
		"```json\n{\"a\":1}\n```",
		"Here: {\"a\":{\"b\":\"} not the end\"}} trailing",
		`{"a":1}`,
	}
	for _, in := range cases {
		raw := extractJSONObject(in)
		if raw == nil || !strings.Contains(string(raw), `"a"`) {
			t.Errorf("extract %q → %s", in, raw)
		}
	}
	if extractJSONObject("no json here") != nil {
		t.Error("found JSON in plain prose")
	}
}
