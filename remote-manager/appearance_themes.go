package main

// Shared (organisation) themes and the appearance lock (v14.11.8).
//
// An administrator publishes a theme made in Settings → Appearance to everyone
// (appearance_shared_themes) and may lock the colours (appearance_lock): then every user
// sees the default theme, while the accessibility options (colour vision, contrast,
// fonts, size, motion) stay theirs. Values are checked here as strictly as addon themes:
// only token names and short CSS values without anything that loads or runs code.

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const maxSharedThemes = 12

var (
	sharedThemeIDRe    = regexp.MustCompile(`^org-[a-z0-9][a-z0-9-]{0,27}$`)
	sharedThemeTokenRe = regexp.MustCompile(`^[a-z0-9-]{1,40}$`)
	sharedThemeTermRe  = regexp.MustCompile(`^[a-zA-Z]{1,30}$`)
	themeUnsafeValueRe = regexp.MustCompile(`(?i)[;{}<>\\]|url\(|expression|@import`)
)

type sharedTheme struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Mode   string            `json:"mode"`
	Tokens map[string]string `json:"tokens"`
	Term   map[string]string `json:"term,omitempty"`
}

func safeThemeValue(v string) bool {
	return v != "" && len(v) <= 200 && !themeUnsafeValueRe.MatchString(v)
}

// validateSharedThemes checks and normalises the appearance_shared_themes setting.
func validateSharedThemes(v string) (string, error) {
	var in []sharedTheme
	if err := json.Unmarshal([]byte(v), &in); err != nil {
		return "", fmt.Errorf("invalid theme list: %v", err)
	}
	if len(in) > maxSharedThemes {
		return "", fmt.Errorf("at most %d shared themes", maxSharedThemes)
	}
	seen := map[string]bool{}
	out := make([]sharedTheme, 0, len(in))
	for _, th := range in {
		if !sharedThemeIDRe.MatchString(th.ID) || seen[th.ID] {
			return "", fmt.Errorf("invalid or repeated theme id %q (use org-<name>)", th.ID)
		}
		seen[th.ID] = true
		name := strings.TrimSpace(th.Name)
		if name == "" || len([]rune(name)) > 60 || strings.ContainsAny(name, "<>") {
			return "", fmt.Errorf("theme %s: name must be 1 to 60 characters without < or >", th.ID)
		}
		if th.Mode != "light" {
			th.Mode = "dark"
		}
		if len(th.Tokens) == 0 || len(th.Tokens) > 64 || len(th.Term) > 32 {
			return "", fmt.Errorf("theme %s: 1 to 64 colours and at most 32 terminal colours", th.ID)
		}
		for k, val := range th.Tokens {
			if !sharedThemeTokenRe.MatchString(k) || !safeThemeValue(val) {
				return "", fmt.Errorf("theme %s: invalid colour %q", th.ID, k)
			}
		}
		for k, val := range th.Term {
			if !sharedThemeTermRe.MatchString(k) || !safeThemeValue(val) {
				return "", fmt.Errorf("theme %s: invalid terminal colour %q", th.ID, k)
			}
		}
		out = append(out, sharedTheme{ID: th.ID, Name: name, Mode: th.Mode, Tokens: th.Tokens, Term: th.Term})
	}
	b, _ := json.Marshal(out)
	return string(b), nil
}

// sharedThemes returns the published themes for the page (an empty list on any error).
func sharedThemes() []sharedTheme {
	out := []sharedTheme{}
	if v := getSetting("appearance_shared_themes"); v != "" {
		_ = json.Unmarshal([]byte(v), &out)
	}
	return out
}

// appearancePolicy is the "appearance" object of /api/me.
func appearancePolicy() map[string]interface{} {
	return map[string]interface{}{
		"default_theme": getSetting("appearance_default_theme"),
		"custom_css":    settingBool("appearance_custom_css"),
		"lock":          settingBool("appearance_lock"),
		"shared":        sharedThemes(),
	}
}
