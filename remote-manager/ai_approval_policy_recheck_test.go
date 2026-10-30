package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAIApprovalRechecksCurrentModeBeforeRunning(t *testing.T) {
	var input string
	e := newAIEnv(t, func(n int, body map[string]interface{}) fakeTurn {
		if n == 0 {
			return fakeTurn{Calls: []fakeCall{{"run_command", input}}}
		}
		return fakeTurn{Text: "done"}
	})
	target := filepath.Join(e.host.home, "must-not-run-after-policy-change")
	b, _ := json.Marshal(map[string]string{"command": "touch " + filepath.ToSlash(target), "reason": "test"})
	input = string(b)
	s := e.start("ask", nil)
	e.prompt(s, "run the command")
	a := e.waitApproval(s)

	if code, v := e.user.aiDo("POST", fmt.Sprintf("/api/ai/sessions/%d/mode", s.ID), map[string]interface{}{"mode": "read_only"}); code != 200 {
		t.Fatalf("change session to read-only: %d %v", code, v)
	}
	e.decide(s, a, map[string]interface{}{"decision": "approve"})
	e.waitIdle(s)
	if _, err := os.Stat(target); err == nil {
		t.Fatal("an unchanged pending approval ran after the session became read-only")
	}
	results := lastToolResults(e.llm.body(1))
	if len(results) != 1 || !strings.Contains(results[0], "Denied by WRM after approval") {
		t.Fatalf("approved command result does not report the current policy denial: %v", results)
	}
}
