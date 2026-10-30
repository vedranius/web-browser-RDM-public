package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConnectionTestDoesNotReuseStoredSecretAfterTargetChange(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "conntest-target", false)

	storedHost := startTestSSHServer(t)
	changedHost := startTestSSHServer(t)
	id := u.addConnection("saved", storedHost)

	// The saved password still works for the connection's original target.
	var result map[string]interface{}
	if code := u.jsonDo("POST", "/api/connections/test", map[string]interface{}{
		"id": id, "name": "saved", "protocol": "SSH", "host": storedHost,
		"username": testSSHUser, "auth_method": "PASSWORD",
	}, &result); code != 200 || result["ok"] != true {
		t.Fatalf("test using saved credentials for the original target: status=%d result=%v", code, result)
	}

	// Changing only the target must not silently carry the stored password to it.
	result = nil
	if code := u.jsonDo("POST", "/api/connections/test", map[string]interface{}{
		"id": id, "name": "saved", "protocol": "SSH", "host": changedHost,
		"username": testSSHUser, "auth_method": "PASSWORD",
	}, &result); code != 409 || !strings.Contains(result["error"].(string), "not reused for a changed target") {
		t.Fatalf("test with changed target reused saved credentials: status=%d result=%v", code, result)
	}
}
