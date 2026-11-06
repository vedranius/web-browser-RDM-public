package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

// A storage can be renamed, have its permissions changed and be removed while its endpoint is
// down; only a change to the connection fields is proven against the bucket again.
func TestStorageEditWithoutReachingTheBucket(t *testing.T) {
	if _, err := setSetting("storages_enabled", "1"); err != nil {
		t.Fatal(err)
	}
	stub := newS3Stub(t, "wrm-edit-bucket", "AKIAEDIT", "secret-edit-key")
	s3srv := httptest.NewServer(stub)
	app := httptest.NewServer(newRouter())
	t.Cleanup(app.Close)
	u := newTestUser(t, app, "storage-edit-offline", false)
	body := map[string]interface{}{
		"name": "do-spaces", "endpoint": s3srv.URL, "region": "us-east-1", "bucket": "wrm-edit-bucket",
		"prefix": "", "path_style": true, "access_key": "AKIAEDIT", "secret": "secret-edit-key",
		"allow_write": true,
	}
	var created Storage
	if code := u.jsonDo("POST", "/api/storages", body, &created); code != http.StatusOK {
		t.Fatalf("create: %d", code)
	}
	s3srv.Close() // the endpoint is gone from now on

	path := "/api/storages/" + strconv.Itoa(created.ID)
	body["secret"] = ""
	body["name"] = "do-spaces-renamed"
	body["allow_delete"] = true
	var updated Storage
	if code := u.jsonDo("PUT", path, body, &updated); code != http.StatusOK {
		t.Fatalf("rename while the endpoint is down: %d", code)
	}
	if updated.Name != "do-spaces-renamed" || !updated.AllowDelete || !updated.SecretSet {
		t.Fatalf("update result: %+v", updated)
	}

	body["region"] = "eu-central-1"
	if code := u.jsonDo("PUT", path, body, nil); code != http.StatusBadRequest {
		t.Fatalf("a region change must be checked against the bucket, got %d", code)
	}

	if code := u.jsonDo("DELETE", path, nil, nil); code != http.StatusOK {
		t.Fatalf("delete while the endpoint is down: %d", code)
	}
}

func TestStorageCheckInterval(t *testing.T) {
	t.Cleanup(func() { setSetting("storage_check_minutes", "15") })
	if _, err := setSetting("storage_check_minutes", "0"); err != nil {
		t.Fatal(err)
	}
	if d := storageCheckInterval(); d != 0 {
		t.Fatalf("0 must turn the check off, got %v", d)
	}
	if _, err := setSetting("storage_check_minutes", "30"); err != nil {
		t.Fatal(err)
	}
	if d := storageCheckInterval(); d != 30*time.Minute {
		t.Fatalf("interval: %v", d)
	}
}

// A chain's storage check reuses its answer for chain_storage_check_s, so a 30-second chain does
// not ask the bucket on every poll.
func TestChainStorageCheckIsCached(t *testing.T) {
	if _, err := setSetting("storages_enabled", "1"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { setSetting("chain_storage_check_s", "300") })
	setSetting("chain_storage_check_s", "300")
	stub := newS3Stub(t, "wrm-chain-cache", "AKIACHAIN", "secret-chain-key")
	var calls int
	var mu sync.Mutex
	s3srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls++
		mu.Unlock()
		stub.ServeHTTP(w, r)
	}))
	t.Cleanup(s3srv.Close)
	app := httptest.NewServer(newRouter())
	t.Cleanup(app.Close)
	u := newTestUser(t, app, "chain-storage-cache", false)
	var st Storage
	if code := u.jsonDo("POST", "/api/storages", map[string]interface{}{
		"name": "cache", "endpoint": s3srv.URL, "region": "us-east-1", "bucket": "wrm-chain-cache",
		"path_style": true, "access_key": "AKIACHAIN", "secret": "secret-chain-key",
	}, &st); code != http.StatusOK {
		t.Fatalf("create: %d", code)
	}
	mu.Lock()
	calls = 0
	mu.Unlock()
	node := chainNode{ID: "b", Kind: "storage", StorageID: st.ID}
	check := chainCheck{ID: "b_ok", Node: "b", Type: "storage"}
	for i := 0; i < 3; i++ {
		if r := chainStorage(context.Background(), check, u.userID, node); !r.OK {
			t.Fatalf("storage check: %+v", r)
		}
	}
	mu.Lock()
	got := calls
	mu.Unlock()
	if got != 1 {
		t.Fatalf("three polls sent %d requests, want 1", got)
	}
	setSetting("chain_storage_check_s", "0")
	chainStorage(context.Background(), check, u.userID, node)
	mu.Lock()
	got = calls
	mu.Unlock()
	if got != 2 {
		t.Fatalf("with caching off the poll must ask again, got %d requests", got)
	}
}
