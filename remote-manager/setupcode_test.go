package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// On a fresh install the first registration from the network needs the one-time setup code;
// a direct loopback request does not; concurrent first registrations yield one admin only.
func TestV1471FirstRegisterNeedsSetupCode(t *testing.T) {
	old := db
	dir := t.TempDir()
	t.Setenv("DB_PATH", filepath.Join(dir, "fresh.db"))
	initDB()
	fresh := db
	t.Cleanup(func() { fresh.Close(); db = old; clearSetupCode() })

	initSetupCode(filepath.Join(dir, "fresh.db"))
	file := filepath.Join(dir, setupCodeFile)
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("code file: %v", err)
	}
	code := strings.TrimSpace(string(raw))
	if len(code) != 14 || strings.Count(code, "-") != 2 {
		t.Fatalf("code format %q", code)
	}

	register := func(remote, xff, user, setup string) int {
		body, _ := json.Marshal(map[string]string{"username": user, "password": "First-Admin-Passw0rd-1471", "setup_code": setup})
		r := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
		r.RemoteAddr = remote
		if xff != "" {
			r.Header.Set("X-Forwarded-For", xff)
		}
		w := httptest.NewRecorder()
		apiAuthRegisterHandler(w, r)
		return w.Code
	}
	cfg := func(remote string) map[string]interface{} {
		r := httptest.NewRequest("GET", "/api/auth/config", nil)
		r.RemoteAddr = remote
		w := httptest.NewRecorder()
		apiAuthConfigHandler(w, r)
		var out map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &out)
		return out
	}
	if cfg("192.0.2.10:5000")["setup_code_required"] != true || cfg("127.0.0.1:5000")["setup_code_required"] != false {
		t.Fatal("config does not report the setup code correctly")
	}
	if c := register("192.0.2.10:5000", "", "intruder", ""); c != 403 {
		t.Fatalf("network register without code: %d", c)
	}
	if c := register("192.0.2.11:5000", "", "intruder", "AAAA-BBBB-CCCC"); c != 403 {
		t.Fatalf("network register with a wrong code: %d", c)
	}
	if c := register("127.0.0.1:5000", "203.0.113.5", "proxied", ""); c != 403 {
		t.Fatalf("proxied loopback register without code: %d", c)
	}

	// Two concurrent first registrations with the right code: exactly one admin.
	var wg sync.WaitGroup
	codes := make([]int, 2)
	for i := range codes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			codes[i] = register("192.0.2.12:5000", "", []string{"admin-a", "admin-b"}[i], strings.ToLower(code))
		}(i)
	}
	wg.Wait()
	var admins, users int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_admin=1`).Scan(&admins)
	db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&users)
	if admins != 1 || users != 1 {
		t.Fatalf("concurrent first registrations: codes %v, admins %d, users %d", codes, admins, users)
	}
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatal("code file left after the first account")
	}
	if setupCodeMatches(code) {
		t.Fatal("code still valid after the first account")
	}
}

func TestV1471DirectLoopback(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	for remote, want := range map[string]bool{"127.0.0.1:1": true, "[::1]:1": true, "192.0.2.1:1": false, "10.0.0.1:1": false} {
		r.RemoteAddr = remote
		if directLoopback(r) != want {
			t.Errorf("%s: want %v", remote, want)
		}
	}
	r.RemoteAddr = "127.0.0.1:1"
	r.Header.Set("Forwarded", "for=192.0.2.1")
	if directLoopback(r) {
		t.Error("forwarded request counted as loopback")
	}
}
