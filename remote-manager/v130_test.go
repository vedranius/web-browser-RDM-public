package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestV130LogDirectoryNames(t *testing.T) {
	yes := []string{"log", "logs", "LOG", "log-2026", "logs_archive", "service_log", "app-logs"}
	no := []string{"catalog", "dialog", "blog", "node_modules", "bin", "logging.conf"}
	for _, n := range yes {
		if !monitorLogDirName(n) {
			t.Fatalf("%q should be treated as a log folder", n)
		}
	}
	for _, n := range no {
		if monitorLogDirName(n) {
			t.Fatalf("%q should not be treated as a log folder", n)
		}
	}
}

func TestV130JournalRootConfig(t *testing.T) {
	cfg := defaultMonitorConfig()
	cfg.Roots = []monitorLogRoot{{Kind: "journal", Unit: "media-api.service", Service: "media-api"}}
	if err := cfg.normalize(); err != nil {
		t.Fatal(err)
	}
	if cfg.Roots[0].Path != "journal:media-api.service" {
		t.Fatalf("journal root was not canonicalized: %+v", cfg.Roots[0])
	}
	if cfg.Roots[0].Unit != "media-api.service" || cfg.Roots[0].Kind != "journal" {
		t.Fatalf("journal root lost its unit: %+v", cfg.Roots[0])
	}
	// A unit name reaches a remote command line, so anything outside the strict
	// pattern must be refused instead of quoted and hoped for.
	for _, unit := range []string{"", "bad unit", "x; rm -rf /", "../../etc/passwd", "a|b", "$(id)"} {
		cfg := defaultMonitorConfig()
		cfg.Roots = []monitorLogRoot{{Kind: "journal", Unit: unit}}
		if err := cfg.normalize(); err == nil {
			t.Fatalf("accepted unsafe journal unit %q", unit)
		}
	}
	// A file root that only looks like a journal address keeps the file rules.
	cfg = defaultMonitorConfig()
	cfg.Roots = []monitorLogRoot{{Path: "journal:media-api.service"}}
	if err := cfg.normalize(); err == nil {
		t.Fatal("a relative file root must be refused")
	}
}

func TestV130LogsNeedExplicitDiscovery(t *testing.T) {
	if runtime.GOOS == "windows" {
		// The fake host exposes POSIX paths; like the other SSH log tests this
		// runs on the Linux CI runner.
		t.Skip("POSIX log roots")
	}
	h := startFakeHost(t, testSSHUser, testSSHPass)
	srv := newTestServer(t)
	u := newTestUser(t, srv, "log-explicit", false)
	id := u.addConnection("demo-stream", h.addr)
	root := filepath.Join(h.home, "demo-operator", "media-api", "logs")
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "app.log")
	writeFile(t, file, "2026-10-15T10:00:00Z INFO ready\n")
	cfg := defaultMonitorConfig()
	cfg.Roots = []monitorLogRoot{}
	if code := u.jsonDo("PUT", fmt.Sprintf("/api/monitoring/config?id=%d", id), cfg, nil); code != 200 {
		t.Fatalf("config: %d", code)
	}
	saveGitWorkspace(u.userID, &gitCatalog{ServerRoots: []string{h.home}}, &gitSettings{})
	var first struct {
		Files      []monitorLogFile `json:"files"`
		Discovered bool             `json:"discovered"`
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/logs?id=%d", id), nil, &first); code != 200 {
		t.Fatalf("list: %d", code)
	}
	if first.Discovered || len(first.Files) != 0 {
		t.Fatalf("opening the log list scanned the host: %+v", first)
	}
	var scan struct {
		Roots      []monitorLogRoot     `json:"roots"`
		Services   []monitorServiceInfo `json:"services"`
		Discovered bool                 `json:"discovered"`
	}
	if code := u.jsonDo("POST", fmt.Sprintf("/api/monitoring/logs/discover?id=%d", id), nil, &scan); code != 200 {
		t.Fatalf("discover: %d", code)
	}
	if !scan.Discovered || len(scan.Roots) == 0 {
		t.Fatalf("discovery found nothing: %+v", scan)
	}
	found := false
	for _, r := range scan.Roots {
		if r.Path == root && r.Service == "media-api" {
			found = true
		}
	}
	if !found {
		t.Fatalf("service log folder not discovered: %+v", scan.Roots)
	}
	var second struct {
		Files      []monitorLogFile `json:"files"`
		Discovered bool             `json:"discovered"`
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/logs?id=%d", id), nil, &second); code != 200 {
		t.Fatalf("list after discovery: %d", code)
	}
	if !second.Discovered || len(second.Files) != 1 || second.Files[0].Path != file {
		t.Fatalf("discovered logs were not listed: %+v", second)
	}
	if second.Files[0].Service != "media-api" {
		t.Fatalf("service label lost: %+v", second.Files[0])
	}
}

func TestV130JournalSourceIsListed(t *testing.T) {
	h := startFakeHost(t, testSSHUser, testSSHPass)
	srv := newTestServer(t)
	u := newTestUser(t, srv, "log-journal", false)
	id := u.addConnection("demo-stream", h.addr)
	cfg := defaultMonitorConfig()
	cfg.Roots = []monitorLogRoot{{Kind: "journal", Unit: "media-api.service", Service: "media-api"}}
	if code := u.jsonDo("PUT", fmt.Sprintf("/api/monitoring/config?id=%d", id), cfg, nil); code != 200 {
		t.Fatalf("config: %d", code)
	}
	var list struct {
		Files []monitorLogFile `json:"files"`
		Roots []monitorLogRoot `json:"roots"`
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/logs?id=%d", id), nil, &list); code != 200 {
		t.Fatalf("list: %d", code)
	}
	if len(list.Files) != 1 || list.Files[0].Kind != "journal" || list.Files[0].Unit != "media-api.service" {
		t.Fatalf("journal source missing from the catalog: %+v", list)
	}
	if !strings.HasPrefix(list.Files[0].Path, "journal:") {
		t.Fatalf("unexpected journal path: %+v", list.Files[0])
	}
}
