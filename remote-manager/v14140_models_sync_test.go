package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Adding an AI service connects to it and adds every model it offers, keeping the
// models already listed (and the default) first.
func TestAIProviderModelsSync(t *testing.T) {
	t.Setenv("WRM_AI_ASSISTANT", "all")
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "sync-admin", true)
	user := newTestUser(t, srv, "sync-user", false)
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`{"data":[{"id":"demo-b"},{"id":"demo-a"},{"id":"demo-c"},{"id":"bad name"}]}`))
	}))
	defer fake.Close()
	code, v := admin.aiDo("POST", "/api/admin/ai/providers", map[string]interface{}{"name": "Demo local", "kind": "openai_compatible",
		"base_url": fake.URL + "/v1", "models": "demo-c", "default_model": "demo-c"})
	if code != 200 {
		t.Fatalf("create: %d %v", code, v)
	}
	id := int(asMap(v)["id"].(float64))
	code, v = admin.aiDo("POST", fmt.Sprintf("/api/admin/ai/providers/%d/models/sync", id), map[string]interface{}{})
	m := asMap(v)
	if code != 200 || m["ok"] != true || len(m["added"].([]interface{})) != 2 || m["default_model"] != "demo-c" {
		t.Fatalf("sync: %d %v", code, v)
	}
	p, err := loadAIProvider(id)
	if err != nil || p.Models != "demo-c,demo-a,demo-b" || p.DefaultModel != "demo-c" {
		t.Fatalf("saved models: %q default %q err %v", p.Models, p.DefaultModel, err)
	}
	// a second sync adds nothing
	_, v = admin.aiDo("POST", fmt.Sprintf("/api/admin/ai/providers/%d/models/sync", id), map[string]interface{}{})
	if a := asMap(v)["added"].([]interface{}); len(a) != 0 {
		t.Fatalf("second sync: %v", v)
	}
	if ollamaAuditCount("ai.provider_models_synced") < 2 {
		t.Fatal("sync must be audited")
	}
	// a user cannot sync an organisation provider (not found in the user scope)
	code, _ = user.aiDo("POST", fmt.Sprintf("/api/ai/providers/%d/models/sync", id), map[string]interface{}{})
	if code != 404 && code != 403 {
		t.Fatalf("user sync of an org provider: %d", code)
	}
	// a provider without a model list answers ok=false and keeps its models
	_, v = admin.aiDo("POST", "/api/admin/ai/providers", map[string]interface{}{"name": "Demo Azure", "kind": "azure",
		"base_url": "https://demo.example.com", "api_key": "demo-key-0123456789", "models": "demo-deployment"})
	id = int(asMap(v)["id"].(float64))
	_, v = admin.aiDo("POST", fmt.Sprintf("/api/admin/ai/providers/%d/models/sync", id), map[string]interface{}{})
	if asMap(v)["ok"] != false {
		t.Fatalf("azure sync: %v", v)
	}
	if p, _ := loadAIProvider(id); p.Models != "demo-deployment" {
		t.Fatalf("azure models changed: %q", p.Models)
	}
}
