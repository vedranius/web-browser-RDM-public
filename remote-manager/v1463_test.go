package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestV1463LogSkipDirs(t *testing.T) {
	for _, n := range []string{"venv", "venv.bak.2026-07-16", ".venv311", "site-packages", "node_modules", "VirtualEnv-old", "__pycache__"} {
		if !logSkipDir(n, nil) {
			t.Errorf("%s not skipped", n)
		}
	}
	for _, n := range []string{"logs", "log", "FtpLocalDistribute", "novatv", "inventory"} {
		if logSkipDir(n, nil) {
			t.Errorf("%s skipped", n)
		}
	}
	if !logSkipDir("cache-x", []string{"cache-*"}) {
		t.Error("user exclude pattern not applied")
	}
	if s := logFindScript([]string{"/home"}, defaultLogSettings()); !strings.Contains(s, "-name 'venv*'") {
		t.Errorf("find does not prune renamed virtual environments: %s", s)
	}
}

// The viewer and search start a new discovery for the owner when the old one expired
// (#log viewer "Find logs first" after an hour), the discovery stays valid while used, and
// a renamed virtual environment is not walked.
func TestV1463LogRediscoverAndSlidingExpiry(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell")
	}
	h := startFakeHost(t, testSSHUser, testSSHPass)
	srv := newTestServer(t)
	u := newTestUser(t, srv, "log-1463", false)
	id := u.addConnection("demo-1463", h.addr)
	base := filepath.Join(h.home, "skripte")
	svc := filepath.Join(base, "demo", "FtpMirror", "logs")
	venv := filepath.Join(base, "demo", "venv.bak.2026-07-16", "lib", "python3.12", "site-packages", "pkg", "logs")
	for _, d := range []string{svc, venv} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(svc, "mirror.log"), "2026-11-09T10:00:00Z INFO started\n")
	writeFile(t, filepath.Join(venv, "pkg.log"), "not a service log\n")
	set := defaultLogSettings()
	set.ScanPaths = []string{base}
	if code := u.jsonDo("PUT", fmt.Sprintf("/api/monitoring/logs/settings?id=%d", id), set, nil); code != 200 {
		t.Fatalf("settings: %d", code)
	}
	if code := u.jsonDo("POST", fmt.Sprintf("/api/monitoring/logs/discover?id=%d", id), nil, &map[string]interface{}{}); code != 200 {
		t.Fatalf("discover: %d", code)
	}
	c, _ := loadConnection(id)
	view, ok := monitorDiscoveryView(c)
	if !ok {
		t.Fatal("no discovery")
	}
	path := ""
	for _, f := range view.Files {
		if strings.Contains(f.Path, "venv.bak") {
			t.Errorf("virtual environment listed: %s", f.Path)
		}
		if strings.HasSuffix(f.Path, "mirror.log") {
			path = f.Path
		}
	}
	if path == "" {
		t.Fatalf("service log not found: %+v", view.Files)
	}
	// in use: the expiry moves forward
	view.Until = time.Now().Add(5 * time.Minute)
	monitorDiscoveredRoots.Store(id, view)
	if v, ok := monitorDiscoveryView(c); !ok || time.Until(v.Until) < 55*time.Minute {
		t.Fatalf("discovery not kept alive: %v %v", ok, v.Until)
	}
	// expired: the viewer discovers again instead of answering 409
	view.Until = time.Now().Add(-time.Minute)
	monitorDiscoveredRoots.Store(id, view)
	var lines struct {
		Text string `json:"text"`
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/logs/lines?id=%d&path=%s&n=5", id, path), nil, &lines); code != 200 || !strings.Contains(lines.Text, "started") {
		t.Fatalf("viewer after expiry: %d %+v", code, lines)
	}
	monitorDiscoveredRoots.Delete(id)
	if code := u.jsonDo("POST", fmt.Sprintf("/api/monitoring/logs/grep?id=%d", id), map[string]interface{}{"files": []string{path}, "q": "started", "timezone": "UTC"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("search without a discovery: %d", code)
	}
	// a shared connection never starts a discovery
	if logRediscover(nil, c, &connAccess{Share: &ShareAccess{}}) {
		t.Fatal("share started a discovery")
	}
}

// #188: a dynamic tunnel (-D) of the user's connection is offered as a proxy until a proxy
// uses it; another user's tunnel is never offered.
func TestV1463ProxyTunnelOffers(t *testing.T) {
	srv := newTestServer(t)
	u := newTestUser(t, srv, "px-1463", false)
	other := newTestUser(t, srv, "px-1463-b", false)
	cid := u.addConnection("demo-socks", "127.0.0.1:22")
	oid := other.addConnection("demo-other", "127.0.0.1:22")
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := db.Exec(`INSERT INTO connection_tunnels (conn_id, user_id, name, kind, bind_host, bind_port, target_host, target_port, open_scheme, open_path, start_mode, sort, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		cid, u.userID, "socks", "dynamic", "127.0.0.1", 0, "", 0, "", "", "manual", 0, now)
	if err != nil {
		t.Fatal(err)
	}
	tid, _ := res.LastInsertId()
	db.Exec(`INSERT INTO connection_tunnels (conn_id, user_id, name, kind, bind_host, bind_port, target_host, target_port, open_scheme, open_path, start_mode, sort, created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		oid, other.userID, "theirs", "dynamic", "127.0.0.1", 0, "", 0, "", "", "manual", 0, now)
	var out struct {
		Offers []struct {
			TunnelID   int    `json:"tunnel_id"`
			Connection string `json:"connection"`
		} `json:"tunnel_offers"`
	}
	if code := u.jsonDo("GET", "/api/proxies", nil, &out); code != 200 || len(out.Offers) != 1 || out.Offers[0].TunnelID != int(tid) || out.Offers[0].Connection != "demo-socks" {
		t.Fatalf("offers: %d %+v", code, out.Offers)
	}
	var saved map[string]interface{}
	if code := u.jsonDo("POST", "/api/proxies", map[string]interface{}{"name": "demo-socks · socks", "kind": "wrm_tunnel", "tunnel_id": tid, "remote_dns": true}, &saved); code != 200 {
		t.Fatalf("create from offer: %d %v", code, saved)
	}
	out.Offers = nil
	u.jsonDo("GET", "/api/proxies", nil, &out)
	if len(out.Offers) != 0 {
		t.Fatalf("tunnel still offered after it became a proxy: %+v", out.Offers)
	}
}
