package main

import (
	"encoding/base64"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func b64(s string) string { return "b64:" + base64.StdEncoding.EncodeToString([]byte(s)) }

// TestChainNodeOutput: the script's answers become values, and only extracted values remain.
func TestChainNodeOutput(t *testing.T) {
	now := int64(1_800_000_000)
	checks := []chainCheck{
		{ID: "svc", Type: "systemd", Target: "nginx.service"},
		{ID: "ctr", Type: "docker", Target: "web"},
		{ID: "proc", Type: "process", Target: "worker"},
		{ID: "p", Type: "port", Port: 443},
		{ID: "f", Type: "file", Target: "/data/x", MaxAgeS: 60},
		{ID: "f2", Type: "file", Target: "/data/y", MaxAgeS: 60},
		{ID: "dc", Type: "dir_count", Target: "/data/in"},
		{ID: "ini", Type: "config", Target: "/etc/app.ini", Format: "ini", Key: "storage.recordings"},
		{ID: "js", Type: "config", Target: "/etc/app.json", Format: "json", Key: "backend.kind", Expect: &chainCond{Op: "==", Value: "s3"}},
		{ID: "mnt", Type: "mount", Target: "/mnt/a"},
		{ID: "cron", Type: "cron", Target: "rclone copy"},
		{ID: "lost", Type: "systemd", Target: "gone.service"},
	}
	out := strings.Join([]string{
		"@@now " + strconv.FormatInt(now, 10),
		"@svc active",
		"@ctr exited",
		"@proc 3",
		"@p 1",
		"@f 2048 " + strconv.FormatInt(now-10, 10),
		"@f2 2048 " + strconv.FormatInt(now-600, 10),
		"@dc 17",
		"@ini " + b64("[storage]\nRecordings = s3\npassword = hunter2\n"),
		"@js " + b64(`{"backend":{"kind":"fs"},"token":"x"}`),
		"@mnt remote:bucket fuse.rclone ro,nosuid",
		"@cron " + b64("# 0 */2 * * * rclone copy a b\n"),
		"noise that is not an answer",
	}, "\n")
	r := parseChainNodeOutput(out, checks, time.Unix(now, 0))
	want := map[string]struct {
		v  string
		ok bool
	}{
		"svc": {"active", true}, "ctr": {"exited", false}, "proc": {"3", true}, "p": {"listening", true},
		"f": {"present", true}, "f2": {"stale", false}, "dc": {"17", true}, "ini": {"s3", true},
		"js": {"fs", false}, "mnt": {"ro", true}, "cron": {"disabled", false}, "lost": {"unknown", false},
	}
	for id, w := range want {
		if r[id].Value != w.v || r[id].OK != w.ok {
			t.Errorf("%s: got %+v, want value %q ok %v", id, r[id], w.v, w.ok)
		}
	}
	if r["lost"].Error == "" {
		t.Error("a check without an answer must carry an error")
	}
	for id, res := range r {
		if strings.Contains(res.Value+res.Detail, "hunter2") {
			t.Fatalf("%s: a value outside the requested key leaked: %+v", id, res)
		}
	}
	if !strings.Contains(r["mnt"].Detail, "remote:bucket") {
		t.Errorf("mount detail: %q", r["mnt"].Detail)
	}
}

func TestChainExtractFormats(t *testing.T) {
	cases := []struct{ format, key, data, want string }{
		{"json", "a.list.1.name", `{"a":{"list":[{"name":"x"},{"name":"y"}]}}`, "y"},
		{"yaml", "server.port", "server:\n  port: 8080\n", "8080"},
		{"ini", "Main.Mode", "[main]\nmode = fast\n", "fast"},
		{"ini", "top", "top=1\n[s]\ntop=2\n", "1"},
		{"env", "BACKEND", "# c\nexport BACKEND=\"s3\"\n", "s3"},
		{"kv", "listen", "listen: 0.0.0.0:80\n", "0.0.0.0:80"},
		{"regex", `version (\d+\.\d+)`, "app version 2.7 build 1", "2.7"},
	}
	for _, c := range cases {
		got, found, err := chainExtract(c.format, c.key, []byte(c.data))
		if err != nil || !found || got != c.want {
			t.Errorf("%s %s: got %q found=%v err=%v, want %q", c.format, c.key, got, found, err, c.want)
		}
	}
	if _, found, _ := chainExtract("json", "missing", []byte(`{}`)); found {
		t.Error("a missing key was found")
	}
	if _, _, err := chainExtract("json", "db", []byte(`{"db":{"host":"x","password":"y"}}`)); err == nil {
		t.Error("a key pointing to an object returned the whole object")
	}
	if _, _, err := chainExtract("json", "a", []byte(`{bad`)); err == nil {
		t.Error("broken JSON did not report an error")
	}
	if chainCronState("0 1 * * * /usr/bin/job\n", "/usr/bin/job") != "enabled" || chainCronState("", "x") != "missing" {
		t.Error("cron state")
	}
}

func TestChainRulesAndScenarios(t *testing.T) {
	s := chainState{
		Nodes:  map[string]chainNodeState{"main": {State: "down"}, "bkp": {State: "up"}},
		Checks: map[string]chainResult{"dns": {Value: "10.0.0.2", OK: true}, "vxt": {Value: "inactive"}},
	}
	d := chainDef{Scenarios: []chainScenario{
		{ID: "primary", When: chainCond{All: []chainCond{{Node: "main", Op: "==", Value: "up"}, {Check: "dns", Value: "10.0.0.1"}}}},
		{ID: "backup", When: chainCond{Any: []chainCond{{Check: "dns", Op: "in", Values: []string{"10.0.0.2", "10.0.0.3"}}, {Check: "vxt", Field: "ok", Value: "true"}}}},
		{ID: "unknown", When: chainCond{}},
	}}
	if got := chainScenarioOf(d, s); got != "backup" {
		t.Fatalf("scenario: %q", got)
	}
	s.Checks["dns"] = chainResult{Value: "none"}
	if got := chainScenarioOf(d, s); got != "unknown" {
		t.Fatalf("fallback scenario: %q", got)
	}
	ops := []struct {
		c    chainCond
		v    string
		want bool
	}{
		{chainCond{Op: "!=", Value: "a"}, "b", true},
		{chainCond{Op: "contains", Value: "nova"}, "s3-novabackup", true},
		{chainCond{Op: "exists"}, "missing", false},
		{chainCond{Op: "exists"}, "x", true},
		{chainCond{Op: "bogus"}, "x", false},
	}
	for _, o := range ops {
		if o.c.matchValue(o.v) != o.want {
			t.Errorf("%+v on %q", o.c, o.v)
		}
	}
}

// TestChainChangeConfirmation: a value must hold for confirm polls before it is a change; a
// flapping value never is; new checks are adopted silently; a new server that is down is told.
func TestChainChangeConfirmation(t *testing.T) {
	obs := func(dns, node string) chainState {
		scn := "primary"
		if dns != "a" {
			scn = "backup"
		}
		return chainState{At: nowRFC(), Nodes: map[string]chainNodeState{"n": {State: node}},
			Checks: map[string]chainResult{"dns": {Value: dns, OK: true}}, Scenario: scn}
	}
	var conf chainState
	pend := map[string]chainPending{}
	ch := chainApply(&conf, obs("a", "down"), pend, 2)
	if len(ch) != 2 || ch[0].Kind != "node" || ch[0].New != "down" || ch[1].Kind != "scenario" || ch[1].Old != "" {
		t.Fatalf("first poll: %+v", ch)
	}
	// Flapping: b, a, b, a never holds twice in a row.
	for _, v := range []string{"b", "a", "b", "a"} {
		if ch := chainApply(&conf, obs(v, "down"), pend, 2); len(ch) != 0 {
			t.Fatalf("flapping %s reported: %+v", v, ch)
		}
	}
	if ch := chainApply(&conf, obs("b", "up"), pend, 2); len(ch) != 0 {
		t.Fatalf("first sight of a change reported: %+v", ch)
	}
	ch = chainApply(&conf, obs("b", "up"), pend, 2)
	kinds := map[string]chainChange{}
	for _, c := range ch {
		kinds[c.Kind] = c
	}
	if kinds["node"].New != "up" || kinds["check"].New != "b" || kinds["check"].OKFlip || kinds["scenario"].New != "backup" || kinds["scenario"].Old != "primary" {
		t.Fatalf("confirmed changes: %+v", ch)
	}
	// A new check is adopted without a change.
	o := obs("b", "up")
	o.Checks["new"] = chainResult{Value: "x"}
	if ch := chainApply(&conf, o, pend, 2); len(ch) != 0 {
		t.Fatalf("a new check reported a change: %+v", ch)
	}
}

func TestChainValidation(t *testing.T) {
	u, _, storageID := storageTestEnv(t, "chain-validate", nil)
	other := newTestUser(t, u.srv, "chain-validate-other", false)
	mine := chainTestConn(t, u.userID, "SSH")
	theirs := chainTestConn(t, other.userID, "SSH")
	rdp := chainTestConn(t, u.userID, "RDP")
	d := chainDef{
		Name: "demo",
		Nodes: []chainNode{
			{ID: "app", Kind: "server", ConnID: mine},
			{ID: "foreign", Kind: "server", ConnID: theirs},
			{ID: "desk", Kind: "server", ConnID: rdp},
			{ID: "web", Kind: "external"},
			{ID: "bucket", Kind: "storage", StorageID: storageID},
		},
		Checks: []chainCheck{
			{ID: "ok1", Node: "app", Type: "systemd", Target: "nginx.service"},
			{ID: "secret", Node: "app", Type: "config", Target: "/etc/app.ini", Format: "ini", Key: "db.password"},
			{ID: "wrongnode", Node: "web", Type: "systemd", Target: "nginx.service"},
			{ID: "meta", Node: "web", Type: "http", Target: "http://169.254.169.254/latest/"},
			{ID: "rel", Node: "app", Type: "file", Target: "relative/path"},
			{ID: "unit", Node: "app", Type: "systemd", Target: "a;b"},
			{ID: "s3", Node: "bucket", Type: "storage"},
		},
		Scenarios: []chainScenario{{ID: "x", When: chainCond{Check: "nope", Op: "=~"}}},
	}
	errs := strings.Join(validateChain(&d, u.userID), "\n")
	for _, want := range []string{"foreign", "desk", "looks like a secret", "wrongnode", "metadata", "absolute path", "unit name", "unknown check", "op must be"} {
		if !strings.Contains(errs, want) {
			t.Errorf("validation does not mention %q:\n%s", want, errs)
		}
	}
	if strings.Contains(errs, "(ok1)") || strings.Contains(errs, "(s3)") {
		t.Errorf("valid checks were refused:\n%s", errs)
	}
}

func chainTestConn(t *testing.T, userID int, protocol string) int {
	t.Helper()
	res, err := db.Exec(`INSERT INTO connections (name, protocol, host, username, user_id) VALUES (?,?,?,?,?)`,
		"conn-"+protocol+"-"+strconv.Itoa(userID), protocol, "127.0.0.1:1", "demo", userID)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

// TestChainAPI runs a chain end to end with checks WRM can reach in a test: HTTP with a JSON
// value, TCP, a storage, and a server whose SSH port is closed (node down).
func TestChainAPI(t *testing.T) {
	u, _, storageID := storageTestEnv(t, "chain-api", nil)
	if _, err := setSetting("chains_enabled", "1"); err != nil {
		t.Fatal(err)
	}
	mode := "primary"
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ingest":{"active":"` + mode + `"}}`))
	}))
	t.Cleanup(api.Close)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()
	conn := chainTestConn(t, u.userID, "SSH")
	def := chainDef{
		Name: "Demo chain", Confirm: 1,
		Nodes: []chainNode{
			{ID: "ingest", Kind: "server", Label: "demo-ingest", ConnID: conn},
			{ID: "api", Kind: "external", Label: "Status API"},
			{ID: "bucket", Kind: "storage", StorageID: storageID},
		},
		Checks: []chainCheck{
			{ID: "svc", Node: "ingest", Type: "systemd", Target: "demo.service"},
			{ID: "active", Node: "api", Type: "http", Target: api.URL, Key: "ingest.active"},
			{ID: "port", Node: "api", Type: "tcp", Target: ln.Addr().String()},
			{ID: "s3", Node: "bucket", Type: "storage"},
		},
		Links: []chainLink{{From: "ingest", To: "bucket", Kind: "write", Check: "s3"}},
		Scenarios: []chainScenario{
			{ID: "primary", Label: "Primary", When: chainCond{Check: "active", Value: "primary"}},
			{ID: "backup", Label: "Backup", When: chainCond{Check: "active", Value: "backup"}},
		},
	}
	var v struct {
		ID int `json:"id"`
	}
	if code := u.jsonDo("POST", "/api/chains/validate", map[string]interface{}{"definition": def}, nil); code != http.StatusOK {
		t.Fatalf("validate: %d", code)
	}
	if code := u.jsonDo("POST", "/api/chains", map[string]interface{}{"definition": def}, &v); code != http.StatusOK || v.ID == 0 {
		t.Fatalf("create: %d", code)
	}
	base := "/api/chains/" + strconv.Itoa(v.ID)
	var got struct {
		State    chainState `json:"state"`
		Scenario string     `json:"scenario"`
	}
	if code := u.jsonDo("POST", base+"/check", nil, &got); code != http.StatusOK {
		t.Fatalf("check: %d", code)
	}
	if got.Scenario != "primary" || !got.State.Checks["active"].OK || got.State.Checks["port"].Value != "open" ||
		got.State.Checks["s3"].Value != "reachable" || got.State.Nodes["ingest"].State != "down" || got.State.Checks["svc"].OK {
		t.Fatalf("first state: %+v scenario %q", got.State, got.Scenario)
	}

	mode = "backup"
	if code := u.jsonDo("POST", base+"/check", nil, &got); code != http.StatusOK || got.Scenario != "backup" {
		t.Fatalf("after the switch: %d %q", code, got.Scenario)
	}
	var tl []map[string]interface{}
	u.jsonDo("GET", base+"/timeline", nil, &tl)
	found := false
	for _, e := range tl {
		if e["kind"] == "scenario" && e["old"] == "primary" && e["new"] == "backup" {
			found = true
		}
	}
	if !found {
		t.Fatalf("the switch is not in the timeline: %+v", tl)
	}
	var events int
	db.QueryRow(`SELECT COUNT(1) FROM user_events WHERE user_id=? AND event IN ('chain.scenario_changed','chain.node_down')`, u.userID).Scan(&events)
	if events < 2 {
		t.Fatalf("expected scenario and node-down events in the Workboard, got %d", events)
	}

	var list []chainSummary
	u.jsonDo("GET", "/api/chains", nil, &list)
	if len(list) != 1 || list[0].Scenario != "backup" || list[0].NodesTotal != 1 || list[0].NodesUp != 0 {
		t.Fatalf("list: %+v", list)
	}
	resp := u.do("GET", base+"/export", nil, "")
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(resp.Header.Get("Content-Disposition"), "chain-demo-chain.json") {
		t.Fatalf("export: %d %q", resp.StatusCode, resp.Header.Get("Content-Disposition"))
	}

	// Another user cannot see it; the kill switch closes the API.
	other := newTestUser(t, u.srv, "chain-api-other", false)
	if code := other.jsonDo("GET", base, nil, nil); code != http.StatusNotFound {
		t.Fatalf("another user's chain: %d", code)
	}
	setSetting("chains_enabled", "0")
	if code := u.jsonDo("GET", "/api/chains", nil, nil); code != http.StatusForbidden {
		t.Fatalf("chains off: %d", code)
	}
	setSetting("chains_enabled", "1")
	if code := u.jsonDo("DELETE", base, nil, nil); code != http.StatusOK {
		t.Fatalf("delete: %d", code)
	}
}

// TestChainNodeScriptRuns runs the fixed script in a real POSIX shell against real files, so
// its syntax and the answers' format are checked, not only the parser.
func TestChainNodeScriptRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("needs a POSIX shell")
	}
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	dir := t.TempDir()
	cfg := filepath.Join(dir, "app.json")
	os.WriteFile(cfg, []byte(`{"backend":"s3"}`), 0o600)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("1"), 0o600)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("2"), 0o600)
	checks := []chainCheck{
		{ID: "cfg", Type: "config", Target: cfg, Format: "json", Key: "backend"},
		{ID: "cnt", Type: "dir_count", Target: dir},
		{ID: "file", Type: "file", Target: filepath.Join(dir, "a.txt")},
		{ID: "none", Type: "file", Target: filepath.Join(dir, "missing; echo injected")},
		{ID: "proc", Type: "process", Target: "this-process-does-not-exist-" + strconv.Itoa(os.Getpid())},
	}
	cmd := exec.Command(sh, "-c", chainNodeScript)
	cmd.Stdin = strings.NewReader(chainNodeInput(checks))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("script: %v\n%s", err, out)
	}
	if strings.Contains(string(out), "injected") {
		t.Fatalf("a target was executed: %s", out)
	}
	r := parseChainNodeOutput(string(out), checks, time.Now())
	if r["cfg"].Value != "s3" || r["cnt"].Value != "3" || r["file"].Value != "present" || r["none"].Value != "missing" || r["proc"].Value != "0" {
		t.Fatalf("results: %+v\n%s", r, out)
	}
}
