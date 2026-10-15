package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// "*" in the automatic allow list allows every command, also chains, variables and
// redirections that a narrower list cannot check; blocked commands, deny lists, sensitive
// files and unparseable commands still apply.
func TestV1252AllowEverything(t *testing.T) {
	home := "/home/demo"
	now := time.Now()
	all := &aiAutoLimits{Allow: []string{"*"}, PathAllow: []string{"/**"}, PathDeny: []string{"/etc/shadow"}, Deny: []string{"shutdown*"}, UntilGrant: true}
	narrow := &aiAutoLimits{Allow: []string{"cp *", "ls *"}, Until: now.Add(time.Hour), MaxActions: 10}
	backup := `ts=$(date +%F); cp -a /opt/app/storage_backend.ini /opt/app/storage_backend.ini.bak-$ts && cp -a /srv/app/storage_backend.ini /srv/app/storage_backend.ini.bak-$ts && ls -la /opt/app/ 2>&1`
	cases := []struct {
		auto *aiAutoLimits
		cmd  string
		want string
	}{
		{all, backup, aiAllow},
		{all, "venv=/opt/venv; $venv/bin/pip install -r req.txt && systemctl restart app", aiAllow},
		{all, "echo done > /opt/app/state.txt", aiAllow},
		{all, "(cd /srv && make install)", aiAllow},
		{all, "echo x > /etc/shadow", aiDeny},          // path deny list
		{all, "shutdown -h now", aiDeny},               // deny list
		{all, "rm -rf --no-preserve-root /", aiDeny},   // blocked command
		{all, "cp /etc/shadow /tmp/x", aiApprove},      // sensitive file
		{all, `echo "unterminated`, aiApprove},         // cannot be parsed
		{narrow, backup, aiApprove},                    // a narrow list cannot check variables
		{narrow, "cp -a /a /b && ls -la /b", aiAllow},  // but allows a plain chain
		{narrow, "cp -a /a /b && touch /c", aiApprove}, // outside the list
	}
	for _, c := range cases {
		d := aiEvaluate(aiModeAuto, c.auto, aiAction{Tool: "run_command", Kind: "command", Command: c.cmd}, home, now)
		if d.Action != c.want {
			t.Errorf("%q: %s (%s: %s), want %s", c.cmd, d.Action, d.Rule, d.Reason, c.want)
		}
	}
	// without limits the action count and time never stop it, until the connection expires
	all.Used, all.Until = 100000, now.Add(time.Hour)
	if ok, _ := all.active(now); !ok {
		t.Fatal("until_grant stopped by the action count")
	}
	if ok, why := all.active(now.Add(2 * time.Hour)); ok || !strings.Contains(why, "expired") {
		t.Fatalf("after the connection's expiry: %v %s", ok, why)
	}
	// the refusal text names the rule and what the user can change
	why := aiApprovalWhy(aiDecision{Rule: "auto_limit", Reason: "the automatic mode time limit is over"}, true)
	if !strings.Contains(why, "time limit is over") || !strings.Contains(why, "Allow everything") || !strings.Contains(why, "rule auto_limit") {
		t.Fatalf("auto_limit hint: %s", why)
	}
}

// An AI connection with "Allow everything" and no limits: commands run without asking for
// as long as the connection is valid, the janitor never switches it back to ask, and only
// the administrator's policy and AI connections can use it.
func TestV1252AutoUntilGrantMCP(t *testing.T) {
	e := newMCPEnv(t)
	auto := map[string]interface{}{"allow": []string{"*"}, "path_allow": []string{"/**"}, "until_grant": true}
	tok, gid := e.pat("auto", []string{"read_logs", "run_with_approval", "edit_file_with_approval"}, map[string]interface{}{"auto": auto, "confirm_auto": true, "hours": 48})
	g, err := loadGrant(gid)
	if err != nil || g.Auto == nil || !g.Auto.UntilGrant || g.Auto.MaxActions != 0 {
		t.Fatalf("grant: %v %+v", err, g)
	}
	c := e.client(tok)
	c.initialize(true)
	c.onRequest = func(method string, params map[string]interface{}) map[string]interface{} {
		t.Errorf("asked although everything is allowed: %s %v", method, params["message"])
		return map[string]interface{}{"action": "decline"}
	}
	sid := c.open("demo-api")
	s := mcpSessionOf(sid)
	s.mu.Lock()
	until := s.Auto.Until
	s.mu.Unlock()
	if !until.Equal(g.ExpiresAt) {
		t.Fatalf("automatic mode ends at %v, the connection at %v", until, g.ExpiresAt)
	}
	cmd := "d=~/bk; mkdir -p $d && touch ~/conf.ini && cp -a ~/conf.ini $d/conf.ini.bak && ls -la $d"
	if text, isErr := c.call("run_command", map[string]interface{}{"session_id": sid, "command": cmd, "reason": "backup"}); isErr || !aiFileExists(filepath.Join(e.host.home, "bk", "conf.ini.bak")) {
		t.Fatalf("allowed command: %v %s", isErr, text)
	}
	// many actions later, and past the old time limit, it is still automatic
	s.mu.Lock()
	s.Auto.Used = 5000
	s.mu.Unlock()
	aiJanitorTick(time.Now().Add(50 * time.Minute)) // past any minute limit the old mode had, before the idle limit
	if s.Mode != aiModeAuto {
		t.Fatalf("mode after hours: %s", s.Mode)
	}
	if text, isErr := c.call("run_command", map[string]interface{}{"session_id": sid, "command": "touch ~/later", "reason": "x"}); isErr || !aiFileExists(filepath.Join(e.host.home, "later")) {
		t.Fatalf("after many actions: %v %s", isErr, text)
	}
	// the policy can switch it off; the browser panel cannot use it
	t.Setenv("WRM_AI_AUTO_UNTIL_GRANT", "0")
	body := map[string]interface{}{"name": "x", "conn_ids": []int{e.connID}, "mode": "auto", "scopes": []string{"run_with_approval"}, "auto": auto, "confirm_auto": true}
	if code, v := e.user.aiDo("POST", "/api/mcp/grants", body); code == 200 || !strings.Contains(fmt.Sprint(v), "without a time limit") {
		t.Fatalf("policy off: %d %v", code, v)
	}
}

// The browser panel's sessions keep their time and action limits.
func TestV1252PanelKeepsLimits(t *testing.T) {
	e := newAIEnv(t, func(n int, body map[string]interface{}) fakeTurn { return fakeTurn{Text: "ok"} })
	code, v := e.user.aiDo("POST", "/api/ai/sessions", map[string]interface{}{"conn_id": e.connID, "provider_id": e.provID, "mode": "auto", "confirm_auto": true,
		"auto": map[string]interface{}{"allow": []string{"*"}, "until_grant": true}})
	if code == 200 || !strings.Contains(fmt.Sprint(v), "only for AI connections") {
		t.Fatalf("panel session without limits: %d %v", code, v)
	}
}
