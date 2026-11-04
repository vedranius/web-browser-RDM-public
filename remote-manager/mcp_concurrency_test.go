package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMCPConcurrentToolCallsBoundedAndReleased(t *testing.T) {
	e := newMCPEnv(t)
	tok, _ := e.pat("ask", []string{"run_with_approval"}, nil)
	clients := []*mcpTestClient{e.client(tok), e.client(tok)}
	for _, c := range clients {
		c.initialize(false)
	}
	sid := clients[0].open("demo-api")
	s := mcpSessionOf(sid)
	finished := make(chan struct{}, mcpMaxConcurrentCalls)
	cancels := make([]context.CancelFunc, 0, mcpMaxConcurrentCalls)
	t.Cleanup(func() {
		for _, cancel := range cancels {
			cancel()
		}
	})
	for i := 0; i < mcpMaxConcurrentCalls; i++ {
		c := clients[i%len(clients)]
		body, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0", "id": 100 + i, "method": "tools/call",
			"params": map[string]interface{}{"name": "run_command", "arguments": map[string]interface{}{
				"session_id": sid, "command": fmt.Sprintf("touch ~/mcp-limit-%d", i), "reason": "concurrency test",
			}},
		})
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.srv+"/mcp", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		req.Header.Set("Mcp-Session-Id", c.session)
		req.Header.Set("MCP-Protocol-Version", c.version)
		if i%2 == 0 {
			req.Header.Set("Accept", "text/event-stream")
		}
		go func() {
			defer func() { finished <- struct{}{} }()
			res, err := http.DefaultClient.Do(req)
			if err == nil {
				io.Copy(io.Discard, res.Body)
				res.Body.Close()
			}
		}()
	}
	waitMCPCondition(t, func() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		return len(s.pending) == mcpMaxConcurrentCalls
	}, "all accepted calls to reach approval")

	started := time.Now()
	text, isErr := clients[0].call("run_command", map[string]interface{}{"session_id": sid, "command": "pwd"})
	if !isErr || !strings.Contains(text, "busy") || time.Since(started) > time.Second {
		t.Fatalf("overflow must fail promptly: error=%v text=%q elapsed=%s", isErr, text, time.Since(started))
	}
	text, isErr = clients[1].call("open_session", map[string]interface{}{"connection": "demo-api"})
	if !isErr || !strings.Contains(text, "busy") {
		t.Fatalf("open_session must share the work limit: error=%v text=%q", isErr, text)
	}
	// Lightweight discovery and the regular web API remain reachable while all work slots wait.
	if _, isErr := clients[1].call("list_connections", nil); isErr {
		t.Fatal("list_connections must remain available at the work limit")
	}
	res, err := http.Get(e.srv + "/api/version")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("regular API at MCP limit: %d", res.StatusCode)
	}
	for _, cancel := range cancels {
		cancel()
	}
	for range cancels {
		select {
		case <-finished:
		case <-time.After(5 * time.Second):
			t.Fatal("cancelled MCP request did not finish")
		}
	}
	waitMCPCondition(t, func() bool { return len(mcpToolSlots) == 0 }, "cancelled calls to release all slots")
	text, isErr = clients[0].call("open_session", map[string]interface{}{"connection": "demo-api"})
	if isErr {
		t.Fatalf("a call after cancellation must run: %s", text)
	}
	if len(mcpToolSlots) != 0 {
		t.Fatal("a completed open_session did not release its slot")
	}
}

func waitMCPCondition(t *testing.T, check func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for " + what)
}
