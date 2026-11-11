package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// v14.9.0: agentless monitoring over SSH, alert limits with reasons, checks
// that run on the monitored server, log alerts, events, history and log
// patterns.

const v149AgentlessOutput = `@@os
Linux
Linux 6.8.0-45-generic x86_64
Ubuntu 24.04.1 LTS
@@host
demo-node
@@load
3.10 2.00 1.50 2/300 4000
@@stat
cpu  1000 0 1000 8000 0 0 0 0 0 0
cpu  1050 0 1050 8000 100 0 0 0 0 0
@@mem
MemTotal:        8000000 kB
MemAvailable:    1000000 kB
SwapTotal:       1000000 kB
SwapFree:         900000 kB
@@df
Filesystem 1024-blocks Used Available Capacity Mounted on
/dev/sda1 100000 96000 4000 96% /
/dev/sdb1 200000 20000 180000 10% /srv data
tmpfs 1000 0 1000 0% /run
@@inode
Filesystem Inodes IUsed IFree IUse% Mounted on
/dev/sda1 1000 100 900 10% /
@@net
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 100 1 0 0 0 0 0 0 100 1 0 0 0 0 0 0
  eth0: 5000 10 0 0 0 0 0 0 7000 10 0 0 0 0 0 0
@@uptime
12345.67 1000.00
@@ncpu
2
@@procs
250 1
@@tcp
17
@@top
 1234 55.0  2.0 java
 99 1.0 20.5 postgres
 1234 55.0  2.0 java
@@failed
demo-worker.service
@@reboot
yes
@@ntp
no
@@end
`

func TestV149AgentlessParseAndRules(t *testing.T) {
	s := monitorSample{At: nowRFC(), Metrics: map[string]monitorMetric{}}
	ok, why := parseAgentless(v149AgentlessOutput, &s)
	if !ok {
		t.Fatal(why)
	}
	finishMonitor(&s)
	num := func(k string) float64 {
		m := s.Metrics[k]
		if m.Value == nil {
			t.Fatalf("metric %s missing: %+v", k, s.Metrics)
		}
		return *m.Value
	}
	if s.Hostname != "demo-node" || num("cpu_percent") != 50 || num("iowait") != 50 || num("memory_percent") != 87.5 || num("swap_percent") != 10 ||
		num("disk:/") != 96 || num("disk:/srv data") != 10 || num("inode:/") != 10 || num("network:eth0:in") != 5000 || num("load") != 3.1 ||
		num("cpu_count") != 2 || num("processes") != 250 || num("zombies") != 1 || num("tcp_established") != 17 || num("failed_units") != 1 ||
		num("reboot_required") != 1 || num("time_sync") != 0 || num("uptime") != 12345.67 {
		t.Fatalf("parsed: %+v", s.Metrics)
	}
	if _, ok := s.Metrics["disk:/run"]; ok {
		t.Fatal("tmpfs /run counted as a disk")
	}
	if _, ok := s.Metrics["network:lo:in"]; ok {
		t.Fatal("loopback counted")
	}
	if len(s.Top) != 2 || s.Top[0].Name != "java" {
		t.Fatalf("top: %+v", s.Top)
	}
	th := defaultMonitorThresholds()
	// CPU 50 % is fine; load 1.55 per core needs two polls in a row
	applyMonitorRules(&s, monitorSample{}, th)
	codes := map[string]string{}
	for _, r := range s.Reasons {
		codes[r.Code] = r.Severity
	}
	if s.State != "problem" || codes["disk"] != "problem" || codes["memory"] != "warning" || codes["failed_units"] != "warning" ||
		codes["reboot"] != "warning" || codes["time_sync"] != "warning" || codes["load"] != "" || codes["cpu"] != "" {
		t.Fatalf("state %s reasons %+v", s.State, s.Reasons)
	}
	next := monitorSample{At: nowRFC(), Metrics: map[string]monitorMetric{}}
	parseAgentless(v149AgentlessOutput, &next)
	finishMonitor(&next)
	applyMonitorRules(&next, s, th)
	found := false
	for _, r := range next.Reasons {
		if r.Code == "load" && r.Severity == "warning" {
			found = true
		}
		if r.text("hr") == "" || r.text("en") == "" {
			t.Fatal("empty reason text")
		}
	}
	if !found {
		t.Fatalf("sustained load not reported: %+v", next.Reasons)
	}
	if ok, why := parseAgentless("@@os\nFreeBSD\n@@end\n", &monitorSample{Metrics: map[string]monitorMetric{}}); ok || !strings.Contains(why, "Linux") {
		t.Fatal("non-Linux accepted")
	}
	// limits: zero means default, nonsense refused
	var empty monitorThresholds
	if empty.normalize() != nil || empty.CPUWarn != 85 {
		t.Fatal(empty)
	}
	bad := defaultMonitorThresholds()
	bad.Sustain = 99
	if bad.normalize() == nil {
		t.Fatal("sustain 99 accepted")
	}
}

func TestV149CheckValidationAndJudgement(t *testing.T) {
	good := []monitorService{{Kind: "http", Name: "https://demo.example.com/health", Expect: "2xx"}, {Kind: "tcp", Name: "127.0.0.1:5432"},
		{Kind: "tcp", Name: "[::1]:80"}, {Kind: "process", Name: "nginx"}, {Kind: "cert", Name: "demo.example.com"}, {Kind: "systemd", Name: "ssh.service"}}
	for _, x := range good {
		if err := x.normalize(); err != nil {
			t.Fatalf("%+v: %v", x, err)
		}
	}
	c := monitorService{Kind: "cert", Name: "demo.example.com"}
	c.normalize()
	if c.Name != "demo.example.com:443" {
		t.Fatal(c.Name)
	}
	for _, x := range []monitorService{{Kind: "http", Name: "ftp://x"}, {Kind: "http", Name: "http://a b"}, {Kind: "http", Name: "http://u:p@x/"},
		{Kind: "http", Name: "http://x/$(id)"}, {Kind: "tcp", Name: "host"}, {Kind: "tcp", Name: "-oProxy:22"}, {Kind: "tcp", Name: "h;id:22"},
		{Kind: "process", Name: "a b"}, {Kind: "cert", Name: "::1"}, {Kind: "http", Name: "http://x/", Expect: "abc"}} {
		if x.normalize() == nil {
			t.Fatalf("accepted %+v", x)
		}
	}
	// a URL with "=" survives the state parser
	st := parseMonitorServiceStates("http:https://x/?a=b=200 0.120\ntcp:127.0.0.1:22=open\n")
	if st["http:https://x/?a=b"] != "200 0.120" || st["tcp:127.0.0.1:22"] != "open" {
		t.Fatal(st)
	}
	now := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		x     monitorService
		raw   string
		state string
	}{
		{monitorService{Kind: "http"}, "200 0.1", "ok"},
		{monitorService{Kind: "http", Expect: "404"}, "404 0.1", "ok"},
		{monitorService{Kind: "http"}, "503 0.1", "problem"},
		{monitorService{Kind: "http"}, "000 0", "problem"},
		{monitorService{Kind: "tcp"}, "closed", "problem"},
		{monitorService{Kind: "process"}, "running", "ok"},
		{monitorService{Kind: "cert"}, "Dec 20 00:00:00 2026 GMT", "ok"},
		{monitorService{Kind: "cert"}, "Nov 10 00:00:00 2026 GMT", "warning"},
		{monitorService{Kind: "cert"}, "Nov  3 00:00:00 2026 GMT", "problem"},
		{monitorService{Kind: "docker"}, "restarting", "warning"},
	}
	for _, k := range cases {
		if r := judgeMonitorCheck(k.x, k.raw, 14, now); r.State != k.state {
			t.Fatalf("%+v %q: %s (%s)", k.x, k.raw, r.State, r.Detail)
		}
	}
	sug := parseMonitorSuggestions("@@systemd\nnginx.service\ncron.service\nfoo;bar.service\n@@docker\nweb-1\n@@ports\n0.0.0.0:80\n127.0.0.1:5432\n10.0.0.5:9100\n[::]:443\n@@host\ndemo-node\n@@end\n", "demo.example.com")
	kinds := map[string]bool{}
	for _, x := range sug {
		kinds[x.Kind+":"+x.Name] = x.Suggested
	}
	if !kinds["systemd:nginx.service"] || kinds["systemd:cron.service"] != true || !kinds["docker:web-1"] || !kinds["tcp:127.0.0.1:5432"] || !kinds["cert:demo.example.com:443"] {
		t.Fatalf("suggestions %+v", sug)
	}
	if _, ok := kinds["systemd:foo;bar.service"]; ok {
		t.Fatal("invalid unit suggested")
	}
	if _, ok := kinds["tcp:127.0.0.1:9100"]; ok {
		t.Fatal("a listener on one external address was suggested on loopback")
	}
}

func TestV149LogWatchAndPatterns(t *testing.T) {
	w := monitorLogWatch{Path: "/var/log/demo.log", Pattern: "ERROR"}
	if err := w.normalize(); err != nil || w.Severity != "warning" || w.Min != 1 || w.Hold != 15 || w.Name == "" {
		t.Fatal(w, err)
	}
	for _, b := range []monitorLogWatch{{Path: "relative.log", Pattern: "x"}, {Path: "/var/log/a", Pattern: "a\tb"}, {Path: "/var/log/a", Pattern: "("},
		{Path: "/var/log/a", Pattern: "x", Regex: true, Severity: "fatal"}, {Path: "/etc/shadow", Pattern: "x"}, {Path: "journal:a b", Pattern: "x"}} {
		if b.Pattern == "(" {
			b.Regex = true
		}
		if b.normalize() == nil {
			t.Fatalf("accepted %+v", b)
		}
	}
	res := parseMonitorLogWatch("@@W 1 500\n3\n> a\n> b\n> c\n@@E\n@@W 2 -1\n@@E\n", 3)
	if res[0].Pos != 500 || res[0].Count != 3 || len(res[0].Lines) != 3 || !res[1].Missing || res[2].Pos != -2 {
		t.Fatalf("%+v", res)
	}
	lines := []string{}
	for i := 0; i < 50; i++ {
		lines = append(lines, fmt.Sprintf("2026-11-10T10:%02d:00Z INFO request %d from 10.0.0.%d took %d ms", i, i, i, i*3))
	}
	lines = append(lines, "2026-11-10T10:59:00Z ERROR database connection lost to db-1 (id 7f3e2c1a-0000-4000-8000-123456789abc)",
		"   at demo.Main(Main.java:42)",
		"2026-11-10T10:59:30Z ERROR database connection lost to db-1 (id 7f3e2c1a-0000-4000-8000-123456789abd)")
	r := clusterLogLines(map[string][]string{"/var/log/demo.log": lines}, time.UTC, "")
	if r.Lines != 52 || r.Distinct != 2 || r.Patterns[0].Level != "error" || r.Patterns[0].Count != 2 || !r.Patterns[0].New || r.Patterns[1].Count != 50 {
		t.Fatalf("%+v", r)
	}
	if !strings.Contains(r.Patterns[0].Template, "<uuid>") || !strings.Contains(r.Patterns[1].Template, "<ip>") {
		t.Fatalf("templates %q %q", r.Patterns[0].Template, r.Patterns[1].Template)
	}
	only := clusterLogLines(map[string][]string{"/var/log/demo.log": lines}, time.UTC, "error")
	if len(only.Patterns) != 1 {
		t.Fatal(only.Patterns)
	}
}

func TestV149AgentlessEndToEnd(t *testing.T) {
	if _, err := os.Stat("/proc/stat"); err != nil {
		t.Skip("needs Linux /proc")
	}
	h := startFakeHost(t, testSSHUser, testSSHPass)
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "monitor-agentless", false)
	other := newTestUser(t, srv, "monitor-agentless-other", false)
	id := u.addConnection("demo-agentless", h.addr)
	logFile := filepath.Join(h.home, "app.log")
	writeFile(t, logFile, "2026-11-10T10:00:00Z ERROR old problem\n")
	cfg := defaultMonitorConfig()
	if cfg.Mode != "agentless" {
		t.Fatal(cfg.Mode)
	}
	cfg.Enabled = true
	cfg.Services = []monitorService{{Kind: "process", Name: "no-such-proc-x"}}
	cfg.LogWatches = []monitorLogWatch{{Path: logFile, Pattern: "ERROR", Severity: "problem"}}
	if code := u.jsonDo("PUT", fmt.Sprintf("/api/monitoring/config?id=%d", id), cfg, nil); code != 200 {
		t.Fatal(code)
	}
	var s monitorSample
	if code := u.jsonDo("POST", fmt.Sprintf("/api/monitoring/check?id=%d", id), map[string]string{}, &s); code != 200 {
		t.Fatal(code)
	}
	if s.Source != "agentless" || s.Metrics["cpu_percent"].Value == nil || s.Metrics["memory_percent"].Value == nil || s.Checks["process:no-such-proc-x"].State != "problem" {
		t.Fatalf("sample %+v", s)
	}
	if len(s.LogAlerts) != 0 {
		t.Fatalf("old log content alerted: %+v", s.LogAlerts)
	}
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0o600)
	f.WriteString("2026-11-10T10:05:00Z ERROR new problem token=abc\n2026-11-10T10:05:01Z INFO fine\n")
	f.Close()
	time.Sleep(1100 * time.Millisecond)
	if code := u.jsonDo("POST", fmt.Sprintf("/api/monitoring/check?id=%d", id), map[string]string{}, &s); code != 200 {
		t.Fatal(code)
	}
	if len(s.LogAlerts) != 1 || s.LogAlerts[0].Count != 1 || !s.LogAlerts[0].New || s.State != "problem" {
		t.Fatalf("log alert %+v state %s", s.LogAlerts, s.State)
	}
	var events []monitorEvent
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/events?conn_id=%d", id), nil, &events); code != 200 || len(events) == 0 {
		t.Fatal(code, events)
	}
	var foreign []monitorEvent
	other.jsonDo("GET", fmt.Sprintf("/api/monitoring/events?conn_id=%d", id), nil, &foreign)
	if len(foreign) != 0 {
		t.Fatal("events of another user's server were listed")
	}
	var hist map[string]interface{}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/history?id=%d&range=1h", id), nil, &hist); code != 200 || hist["series"] == nil || hist["availability"] == nil {
		t.Fatal(code, hist)
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/history?id=%d&range=7d", id), nil, &hist); code != 200 || hist["samples"].(float64) < 2 {
		t.Fatal(code, hist)
	}
	if code := u.jsonDo("GET", fmt.Sprintf("/api/monitoring/history?id=%d&range=nope", id), nil, nil); code != 400 {
		t.Fatal(code)
	}
	var bulk map[string][]int
	if code := other.jsonDo("POST", "/api/monitoring/bulk", map[string]interface{}{"ids": []int{id}, "enabled": false}, &bulk); code != 200 || len(bulk["updated"]) != 0 {
		t.Fatal("another user changed monitoring", code, bulk)
	}
	if code := u.jsonDo("POST", "/api/monitoring/bulk", map[string]interface{}{"ids": []int{id}, "enabled": false}, &bulk); code != 200 || len(bulk["updated"]) != 1 {
		t.Fatal(code, bulk)
	}
	if loadMonitorConfig(id).Enabled || len(loadMonitorConfig(id).LogWatches) != 1 {
		t.Fatal("bulk lost the configuration")
	}
}
