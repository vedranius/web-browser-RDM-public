package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// While several AI apps work through MCP at once, WRM's own pages must stay quick: MCP tool
// calls must not hold locks or the database long enough to slow the browser's requests.
func TestV1471MCPLoadKeepsUIResponsive(t *testing.T) {
	if testing.Short() {
		t.Skip("load test")
	}
	e := newMCPEnv(t)
	writeHostFile(t, e.host, "app.log", strings.Repeat("a log line\n", 2000))
	tok, _ := e.pat("read_only", []string{"read_logs", "run_readonly"}, nil)

	const clients, calls = 8, 15
	stop := make(chan struct{})
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var errs []string
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := e.client(tok)
			c.initialize(false)
			sid := c.open("demo-api")
			for j := 0; j < calls; j++ {
				var text string
				var isErr bool
				if j%2 == 0 {
					text, isErr = c.call("read_file", map[string]interface{}{"session_id": sid, "path": "~/app.log"})
				} else {
					text, isErr = c.call("run_command", map[string]interface{}{"session_id": sid, "command": "wc -l ~/app.log", "reason": "load"})
				}
				if isErr {
					errMu.Lock()
					errs = append(errs, text)
					errMu.Unlock()
				}
			}
			c.call("close_session", map[string]interface{}{"session_id": sid})
		}()
	}
	var lat []time.Duration
	done := make(chan struct{})
	go func() {
		defer close(done)
		paths := []string{"/api/connections", "/api/me", "/api/folders", "/api/tunnels"}
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
			}
			start := time.Now()
			resp := e.user.do("GET", paths[i%len(paths)], nil, "")
			resp.Body.Close()
			lat = append(lat, time.Since(start))
			time.Sleep(20 * time.Millisecond)
		}
	}()
	began := time.Now()
	wg.Wait()
	close(stop)
	<-done
	if len(errs) > 0 {
		for _, e := range errs {
			t.Log(e)
		}
		t.Fatalf("%d MCP calls failed", len(errs))
	}
	sort.Slice(lat, func(i, j int) bool { return lat[i] < lat[j] })
	p95 := lat[len(lat)*95/100]
	t.Logf("%d MCP calls in %v; %d UI requests meanwhile, p50 %v, p95 %v, max %v",
		clients*calls, time.Since(began).Round(time.Millisecond), len(lat), lat[len(lat)/2], p95, lat[len(lat)-1])
	if p95 > 1500*time.Millisecond {
		t.Fatalf("UI requests slowed down under MCP load: p95 %v", p95)
	}
	_ = fmt.Sprint
}
