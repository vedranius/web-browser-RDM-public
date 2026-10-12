package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Gemini's OpenAI-compatible endpoint returns a thought_signature in extra_content
// of each tool call and rejects the next request when it is missing.
func TestOpenAIToolCallExtraContentRoundTrip(t *testing.T) {
	var sent []byte
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		if n == 2 {
			sent, _ = readAllBody(r)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, `data: {"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_1","function":{"name":"run_command","arguments":"{\"command\":\"uptime\"}"},"extra_content":{"google":{"thought_signature":"SIG123"}}}]},"finish_reason":"tool_calls"}]}`+"\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()
	p := aiProvider{Kind: "openai_compatible", BaseURL: srv.URL}
	req := aiLLMRequest{Model: "gemini", Messages: []aiMsg{{Role: "user", Text: "hi"}}}
	res, err := aiComplete(context.Background(), p, req, nil)
	if err != nil || len(res.Calls) != 1 || !strings.Contains(string(res.Calls[0].Extra), "SIG123") {
		t.Fatalf("%v %+v", err, res)
	}
	req.Messages = append(req.Messages, aiMsg{Role: "assistant", Calls: res.Calls},
		aiMsg{Role: "user", Results: []aiToolResult{{ID: "call_1", Name: "run_command", Content: "ok"}}})
	if _, err := aiComplete(context.Background(), p, req, nil); err != nil {
		t.Fatal(err)
	}
	var body struct {
		Messages []struct {
			ToolCalls []struct {
				Extra json.RawMessage `json:"extra_content"`
			} `json:"tool_calls"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(sent, &body); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range body.Messages {
		for _, c := range m.ToolCalls {
			found = found || strings.Contains(string(c.Extra), "SIG123")
		}
	}
	if !found {
		t.Fatalf("thought_signature was not sent back: %s", sent)
	}
}

func readAllBody(r *http.Request) ([]byte, error) {
	var b []byte
	buf := make([]byte, 4096)
	for {
		k, err := r.Body.Read(buf)
		b = append(b, buf[:k]...)
		if err != nil {
			return b, nil
		}
	}
}

func TestAIErrorClassification(t *testing.T) {
	if !aiRetryable(&aiAPIError{Status: 503}) || aiRetryable(&aiAPIError{Status: 401}) {
		t.Fatal("retryable statuses")
	}
	if !aiTooLarge(&aiAPIError{Status: 413}) || !aiTooLarge(&aiAPIError{Status: 400, Msg: "Prompt is too long"}) || aiTooLarge(&aiAPIError{Status: 401, Msg: "request too large"}) {
		t.Fatal("too large detection")
	}
}

func TestAIShrinkHistory(t *testing.T) {
	big := strings.Repeat("é", 5000)
	msgs := []aiMsg{{Role: "user", Results: []aiToolResult{{ID: "a", Content: big}, {ID: "b", Content: "short"}}}}
	if !aiShrinkHistory(msgs, 1000) {
		t.Fatal("expected a change")
	}
	got := msgs[0].Results[0].Content
	if len(got) > 1200 || !strings.Contains(got, "shortened") || msgs[0].Results[1].Content != "short" {
		t.Fatalf("%d %q", len(got), got[:40])
	}
	if aiShrinkHistory(msgs, 1000) {
		t.Fatal("already short")
	}
}
