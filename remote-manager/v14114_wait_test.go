package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

// An AI app that stopped waiting for an approval waits with wait_for_approval instead of
// carrying on: the tool result names the approval, the tool reports "still waiting", then
// the user's decision.
func TestV14114WaitForApproval(t *testing.T) {
	e := newMCPEnv(t)
	tok, _ := e.pat("ask", []string{"read_logs", "run_with_approval"}, nil)
	c := e.client(tok)
	c.initialize(true)
	c.onRequest = func(method string, params map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"action": "decline"}
	}
	sid := c.open("demo-api")
	s := mcpSessionOf(sid)
	tools := fmt.Sprint(c.toolNames())
	if !strings.Contains(tools, "wait_for_approval") {
		t.Fatal("wait_for_approval is not listed")
	}
	text, isErr := c.call("run_command", map[string]interface{}{"session_id": sid, "command": "touch ~/waited", "reason": "test"})
	m := regexp.MustCompile(`approval_id "([^"]+)"`).FindStringSubmatch(text)
	if !isErr || m == nil || !strings.Contains(text, "WAIT for the user's decision") {
		t.Fatalf("deferred result: %v %s", isErr, text)
	}
	id := m[1]
	if a := pendingOf(s); a == nil || a.ID != id {
		t.Fatalf("pending: %+v", a)
	}
	start := time.Now()
	if text, isErr := c.call("wait_for_approval", map[string]interface{}{"session_id": sid, "approval_id": id, "seconds": 5}); isErr || !strings.Contains(text, "Still waiting") || time.Since(start) < 4*time.Second {
		t.Fatalf("still waiting: %v %s after %s", isErr, text, time.Since(start))
	}
	go func() {
		time.Sleep(700 * time.Millisecond)
		e.user.aiDo("POST", fmt.Sprintf("/api/ai/sessions/%d/approvals/%s", s.ID, id), map[string]interface{}{"decision": "approve"})
	}()
	if text, isErr := c.call("wait_for_approval", map[string]interface{}{"session_id": sid, "approval_id": id, "seconds": 20}); isErr || !strings.Contains(text, "Approved by the user") {
		t.Fatalf("approved: %v %s", isErr, text)
	}
	// a denial and an unknown id
	text, _ = c.call("run_command", map[string]interface{}{"session_id": sid, "command": "touch ~/denied", "reason": "test"})
	id = regexp.MustCompile(`approval_id "([^"]+)"`).FindStringSubmatch(text)[1]
	e.user.aiDo("POST", fmt.Sprintf("/api/ai/sessions/%d/approvals/%s", s.ID, id), map[string]interface{}{"decision": "deny"})
	if text, isErr := c.call("wait_for_approval", map[string]interface{}{"session_id": sid, "approval_id": id}); !isErr || !strings.Contains(text, "Denied by the user") {
		t.Fatalf("denied: %v %s", isErr, text)
	}
	if text, isErr := c.call("wait_for_approval", map[string]interface{}{"session_id": sid, "approval_id": "nope"}); !isErr || !strings.Contains(text, "is not open") {
		t.Fatalf("unknown: %v %s", isErr, text)
	}
}
