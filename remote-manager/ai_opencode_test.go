package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeOpencodeCatalog serves a models.dev-style catalog whose OpenCode models point at base.
func fakeOpencodeCatalog(t *testing.T, base string) {
	doc := fmt.Sprintf(`{"other":{"id":"other","models":{"x":{"id":"x"}}},
	"opencode":{"id":"opencode","npm":"@ai-sdk/openai-compatible","api":"%[1]s/o","models":{
	  "free-chat":{"id":"free-chat","name":"Free Chat","tool_call":true,"cost":{"input":0,"output":0},"provider":{"api":"%[1]s/o"}},
	  "old-free":{"id":"old-free","tool_call":true,"status":"deprecated","cost":{"input":0,"output":0}},
	  "claude-x":{"id":"claude-x","tool_call":true,"cost":{"input":3,"output":15},"provider":{"npm":"@ai-sdk/anthropic","api":"%[1]s/a"}},
	  "gpt-x":{"id":"gpt-x","tool_call":true,"cost":{"input":2,"output":8},"provider":{"npm":"@ai-sdk/openai","api":"%[1]s/o"}},
	  "no-tools-free":{"id":"no-tools-free","tool_call":false,"cost":{"input":0,"output":0}}}}}`, base)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(doc))
	}))
	t.Cleanup(srv.Close)
	old := opencodeCatalogURL
	opencodeCatalogURL = srv.URL
	resetOpencodeCatalog()
	t.Cleanup(func() { opencodeCatalogURL = old; resetOpencodeCatalog() })
}

func resetOpencodeCatalog() {
	opencodeCatalog.Lock()
	opencodeCatalog.models, opencodeCatalog.err, opencodeCatalog.fetched = nil, nil, time.Time{}
	opencodeCatalog.Unlock()
}

func TestOpencodeCatalogAndModelList(t *testing.T) {
	fakeOpencodeCatalog(t, "http://127.0.0.1:1")
	ctx := context.Background()
	p := aiProvider{Kind: "opencode"}
	list, err := opencodeListModels(ctx, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "free-chat" || !list[0].Free {
		t.Fatalf("keyless list should hold only the current free model with tools: %+v", list)
	}
	p.secret.APIKey = "zen-key"
	list, _ = opencodeListModels(ctx, p)
	var ids []string
	for _, m := range list {
		ids = append(ids, m.ID)
	}
	// free first; Responses-API and tool-less models are left out
	if strings.Join(ids, ",") != "free-chat,claude-x" {
		t.Fatalf("keyed list: %v", ids)
	}
	if in, out, ok := opencodePrice("claude-x"); !ok || in != 3 || out != 15 {
		t.Fatalf("price: %v %v %v", in, out, ok)
	}
	if c := aiCost(aiProvider{Kind: "opencode"}, "claude-x", aiUsage{In: 1000000}); c != 3 {
		t.Fatalf("cost: %v", c)
	}
}

func TestOpencodeCompleteRouting(t *testing.T) {
	f := newFakeLLM(t, nil)
	fakeOpencodeCatalog(t, f.srv.URL)
	ctx := context.Background()
	req := aiLLMRequest{Model: "free-chat", MaxTokens: 50, Messages: []aiMsg{{Role: "user", Text: "hi"}}}

	// no key: a free model goes to chat completions with the public key
	if _, err := aiComplete(ctx, aiProvider{Kind: "opencode"}, req, nil); err != nil {
		t.Fatal(err)
	}
	if f.paths[0] != "/o/chat/completions" || f.headers[0].Get("Authorization") != "Bearer public" {
		t.Fatalf("free request: %s %q", f.paths[0], f.headers[0].Get("Authorization"))
	}
	// no key: a paid model is refused before any request
	req.Model = "claude-x"
	if _, err := aiComplete(ctx, aiProvider{Kind: "opencode"}, req, nil); err == nil || !strings.Contains(err.Error(), "not a free") {
		t.Fatalf("paid model without a key: %v", err)
	}
	// with a key: Claude goes to the Messages API with the key
	p := aiProvider{Kind: "opencode"}
	p.secret.APIKey = "zen-secret-key"
	if _, err := aiComplete(ctx, p, req, nil); err != nil {
		t.Fatal(err)
	}
	if n := len(f.paths); n != 2 || f.paths[1] != "/a/messages" || f.headers[1].Get("x-api-key") != "zen-secret-key" {
		t.Fatalf("claude request: %v %q", f.paths, f.headers[len(f.headers)-1].Get("x-api-key"))
	}
	// Responses-API models are reported, not sent
	req.Model = "gpt-x"
	if _, err := aiComplete(ctx, p, req, nil); err == nil || !strings.Contains(err.Error(), "Responses") {
		t.Fatalf("responses model: %v", err)
	}
	// a base URL override wins over the catalog
	p.BaseURL = f.srv.URL + "/gw"
	req.Model = "free-chat"
	if _, err := aiComplete(ctx, p, req, nil); err != nil {
		t.Fatal(err)
	}
	if f.paths[len(f.paths)-1] != "/gw/chat/completions" {
		t.Fatalf("override: %v", f.paths)
	}
}

func TestOpencodeFallbackWithoutCatalog(t *testing.T) {
	zen := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"id":"big-pickle"},{"id":"glm-9"},{"id":"claude-y"},{"id":"gpt-z"},{"id":"kimi-free"}]}`))
	}))
	defer zen.Close()
	old := opencodeCatalogURL
	opencodeCatalogURL = "http://127.0.0.1:1/api.json"
	resetOpencodeCatalog()
	defer func() { opencodeCatalogURL = old; resetOpencodeCatalog() }()
	p := aiProvider{Kind: "opencode", BaseURL: zen.URL}
	list, err := opencodeListModels(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, m := range list {
		ids = append(ids, m.ID)
	}
	if strings.Join(ids, ",") != "big-pickle,kimi-free" {
		t.Fatalf("fallback keyless list: %v", ids)
	}
	if opencodeGuessAPI("claude-y") != "messages" || opencodeGuessAPI("gpt-z") != "responses" || opencodeGuessAPI("glm-9") != "chat" {
		t.Fatal("name rule")
	}
}

func TestOpencodeProviderValidation(t *testing.T) {
	name, kind, models := "Zen", "opencode", "big-pickle"
	var p aiProvider
	if err := (aiProviderInput{Name: &name, Kind: &kind, Models: &models}).apply(&p); err != nil {
		t.Fatalf("a keyless OpenCode provider must be valid: %v", err)
	}
}
