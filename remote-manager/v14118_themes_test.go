package main

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

// v14.11.8: themes an administrator shares with everyone and the appearance lock.
func TestV14118SharedThemesAndLock(t *testing.T) {
	srv := newTestServer(t)
	defer setSetting("appearance_shared_themes", "[]")
	defer setSetting("appearance_lock", "0")
	defer setSetting("appearance_default_theme", "graphite")

	good := `[{"id":"org-ops-blue","name":"Ops blue","mode":"dark","tokens":{"bg":"#101820","accent":"#3a7bd5"},"term":{"background":"#000000"}}]`
	if _, err := setSetting("appearance_shared_themes", good); err != nil {
		t.Fatalf("valid shared theme refused: %v", err)
	}
	bad := []string{
		`{"id":"x"}`, // not a list
		`[{"id":"ops","name":"x","tokens":{"bg":"#000"}}]`,                                                    // id without org-
		`[{"id":"org-a","name":"x","tokens":{"bg":"url(https://example.com/x)"}}]`,                            // loads something
		`[{"id":"org-a","name":"x","tokens":{"bg":"#000;}body{display:none"}}]`,                               // breaks out of the value
		`[{"id":"org-a","name":"<b>x</b>","tokens":{"bg":"#000"}}]`,                                           // markup in the name
		`[{"id":"org-a","name":"x","tokens":{}}]`,                                                             // no colours
		`[{"id":"org-a","name":"x","tokens":{"bg":"#000"}},{"id":"org-a","name":"y","tokens":{"bg":"#111"}}]`, // repeated id
		`[{"id":"org-a","name":"x","tokens":{"bg":"#000"},"term":{"bad key":"#fff"}}]`,
	}
	many := make([]string, maxSharedThemes+1)
	for i := range many {
		many[i] = `{"id":"org-t` + strings.Repeat("a", i+1) + `","name":"t","tokens":{"bg":"#000"}}`
	}
	bad = append(bad, "["+strings.Join(many, ",")+"]")
	for _, v := range bad {
		if _, err := setSetting("appearance_shared_themes", v); err == nil {
			t.Errorf("accepted %s", v)
		}
	}
	if _, err := setSetting("appearance_default_theme", "org-ops-blue"); err != nil {
		t.Errorf("shared theme as default refused: %v", err)
	}
	if _, err := setSetting("appearance_lock", "1"); err != nil {
		t.Fatal(err)
	}

	c := newTestUser(t, srv, "themes14117", false)
	resp := c.do("GET", "/api/auth/me", nil, "")
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var me struct {
		Appearance struct {
			DefaultTheme string        `json:"default_theme"`
			Lock         bool          `json:"lock"`
			Shared       []sharedTheme `json:"shared"`
		} `json:"appearance"`
	}
	if err := json.Unmarshal(b, &me); err != nil {
		t.Fatalf("me: %v %s", err, b)
	}
	a := me.Appearance
	if !a.Lock || a.DefaultTheme != "org-ops-blue" || len(a.Shared) != 1 || a.Shared[0].Tokens["accent"] != "#3a7bd5" {
		t.Errorf("appearance policy not in /api/auth/me: %+v", a)
	}
}
