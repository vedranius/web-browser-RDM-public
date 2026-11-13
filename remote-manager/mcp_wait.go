package main

import (
	"context"
	"errors"
	"time"
)

// ─── MCP: WAITING FOR AN APPROVAL ────────────────────
//
// AI apps give up on a tool call after their own timeout (Codex after 60 s by default), often
// long before the user sees WRM's approval request. WRM then keeps the request open, and
// wait_for_approval lets the AI wait for the decision in short calls that stay under any
// app's timeout, instead of carrying on without the step or assuming an answer.

// mcpWaitApproval waits up to secs (5-55, default 45) for a decision on approval id and
// says what the AI does next.
func mcpWaitApproval(ctx context.Context, s *aiSession, id string, secs int) (string, bool) {
	if id == "" {
		return "approval_id is required: use the approval_id from the tool result that said the call waits for an approval.", true
	}
	if secs <= 0 {
		secs = 45
	}
	secs = max(5, min(secs, 55))
	deadline := time.Now().Add(time.Duration(secs) * time.Second)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		s.mu.Lock()
		a := s.pending[id]
		ended := s.ended
		var expires time.Time
		if a != nil {
			expires = a.Expires
		}
		s.mu.Unlock()
		if a == nil {
			return mcpApprovalOutcome(s, id)
		}
		if ended {
			return "The WRM session ended; the request is closed. Ask the user how to go on.", true
		}
		if time.Now().After(deadline) {
			return "Still waiting for the user's decision in WRM (the request stays open until " + expires.UTC().Format("15:04 UTC") +
				"). Call wait_for_approval again with the same approval_id. Do not continue with other steps meanwhile.", false
		}
		select {
		case <-ctx.Done():
			return "Stopped waiting. The request stays open in WRM: call wait_for_approval again with the same approval_id.", false
		case <-tick.C:
		}
	}
}

// mcpApprovalOutcome says what happened to an approval that no longer waits.
func mcpApprovalOutcome(s *aiSession, id string) (string, bool) {
	var np *aiNotPendingErr
	outcome := ""
	if errors.As(aiNotPending(s.ID, id), &np) {
		outcome, _ = np.Info["outcome"].(string)
	}
	switch outcome {
	case "approved":
		return "Approved by the user. Now call the original tool again with exactly the same arguments to get its result; it does not run twice.", false
	case "denied":
		return "Denied by the user. Do not run it or a variant of it: ask the user how to go on.", true
	}
	return "The approval request " + id + " is not open: it expired or does not belong to this session. Ask the user how to go on.", true
}
