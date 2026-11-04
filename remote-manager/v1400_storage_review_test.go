package main

import (
	"net/http"
	"strconv"
	"testing"
)

// Regression tests from the v14.0.0 release review of Storages.

// TestStorageTestAndUpdateWithScopedKey: "Test connection" on a saved storage uses the
// decrypted secret, and an update keeps working for a key that may only list its prefix.
func TestStorageTestAndUpdateWithScopedKey(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-review-test", nil)
	stub.listOnly = "team-a/"
	base := "/api/storages/" + strconv.Itoa(id)

	if code := u.jsonDo("POST", base+"/test", nil, nil); code != http.StatusOK {
		t.Fatalf("test connection on a saved storage: %d (%s)", code, stub.lastError)
	}
	var got Storage
	if code := u.jsonDo("GET", base, nil, &got); code != http.StatusOK {
		t.Fatalf("get: %d", code)
	}
	update := map[string]interface{}{
		"name": "renamed", "endpoint": got.Endpoint, "region": "us-east-1", "bucket": "wrm-test-bucket",
		"prefix": "team-a", "path_style": true, "access_key": "AKIASTORAGE", "allow_write": true,
	}
	if code := u.jsonDo("PUT", base, update, nil); code != http.StatusOK {
		t.Fatalf("update with a prefix-scoped key: %d", code)
	}

	// A new endpoint needs the secret again; the stored one is not used for another host.
	update["endpoint"] = "https://other.example.com"
	if code := u.jsonDo("PUT", base, update, nil); code != http.StatusBadRequest {
		t.Fatalf("endpoint change without the secret: %d, want 400", code)
	}
}

// TestStorageMoveSignsCopySource: a rename sends x-amz-copy-source, which must be signed.
func TestStorageMoveSignsCopySource(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-review-move", nil)
	stub.put("team-a/old.json", []byte(`{}`))
	if code := u.jsonDo("POST", "/api/storages/"+strconv.Itoa(id)+"/move",
		map[string]interface{}{"from": "old.json", "to": "new.json"}, nil); code != http.StatusOK {
		t.Fatalf("move: %d (%s)", code, stub.lastError)
	}
	stub.mu.Lock()
	_, oldThere := stub.objects["team-a/old.json"]
	_, newThere := stub.objects["team-a/new.json"]
	stub.mu.Unlock()
	if oldThere || !newThere {
		t.Fatalf("after the move: old=%v new=%v", oldThere, newThere)
	}
}

// TestStorageEditorCannotCreateWithoutUpload: saving a new key from the editor is an upload.
func TestStorageEditorCannotCreateWithoutUpload(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-review-upload", nil)
	base := "/api/storages/" + strconv.Itoa(id)
	if code := u.jsonDo("PUT", base+"/text?key=brand-new.json", map[string]interface{}{"content": "{}"}, nil); code != http.StatusForbidden {
		t.Fatalf("new object without allow_upload: %d, want 403", code)
	}
	stub.mu.Lock()
	_, created := stub.objects["team-a/brand-new.json"]
	stub.mu.Unlock()
	if created {
		t.Fatal("the editor created an object without the upload permission")
	}
}

// TestStorageKillSwitchClosesShares: storages_enabled=0 also closes the public share links.
func TestStorageKillSwitchClosesShares(t *testing.T) {
	u, _, id := storageTestEnv(t, "storage-review-kill", nil)
	var sh storageShare
	if code := u.jsonDo("POST", "/api/storages/"+strconv.Itoa(id)+"/shares", map[string]interface{}{"hours": 1}, &sh); code != http.StatusOK {
		t.Fatalf("create share: %d", code)
	}
	if _, err := setSetting("storages_enabled", "0"); err != nil {
		t.Fatal(err)
	}
	defer setSetting("storages_enabled", "1")
	resp, err := http.Get(u.srv.URL + sh.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("share link with storages off: %d, want 404", resp.StatusCode)
	}
}
