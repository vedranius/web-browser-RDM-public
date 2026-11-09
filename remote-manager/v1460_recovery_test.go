package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// An AI with a session on the WRM host must not reset a WRM password or switch the
// lockdown: the temporary password would be in its output and the account would be its.
func TestV1460AIRefusesWRMRecovery(t *testing.T) {
	blocked := []string{
		"wrm -reset-password admin",
		"./wrm-pro-v14.5.1-linux-amd64 --reset-password=admin -reset-2fa",
		"sudo -u wrm /opt/wrm/wrm -reset-password admin",
		`sh -c "wrm -reset-password admin"`,
		"cd /opt/wrm && ./wrm -reset-2fa -reset-password admin",
		`& "C:\wrm\wrm.exe" -reset-password admin`,
		"wrm -lockdown off",
		"echo x; wrm -lockdown=on",
	}
	for _, c := range blocked {
		if b := aiDestructive(c); b == nil || b.Rule != "wrm_recovery" {
			t.Errorf("not refused as wrm_recovery: %q → %+v", c, b)
		}
	}
	for _, c := range []string{"grep reset-password /var/log/wrm.log", "wrm -version", "cat lockdown.txt", "man wrm | grep -i lockdown-reason"} {
		if b := aiDestructive(c); b != nil {
			t.Errorf("refused although harmless: %q → %+v", c, b)
		}
	}
	out, n := redactSecrets("New temporary password for admin: Xy7-secret-PW\n(The user must change it after signing in.)")
	if n == 0 || strings.Contains(out, "Xy7-secret-PW") {
		t.Fatalf("temporary password not redacted: %q", out)
	}
}

// wrm -reset-password tells every administrator (inbox and pending notifications) and
// leaves a security.alert entry, from its own process.
func TestV1460CLIResetAlertsAdmins(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1460-admin", true)
	newTestUser(t, srv, "v1460-victim", false)
	start := time.Now().UTC().Add(-time.Second).Format(time.RFC3339)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := resetPasswordCLI("v1460-victim", true)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Every administrator was told") {
		t.Fatalf("output does not say the administrators were told: %q", buf.String())
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE action='security.alert' AND target='cli_reset' AND details LIKE '%v1460-victim%'`).Scan(&n)
	if n != 1 {
		t.Fatalf("security.alert entries: %d", n)
	}
	db.QueryRow(`SELECT COUNT(*) FROM user_events WHERE user_id=? AND event='security.alert' AND body LIKE '%v1460-victim%' AND created_at>=?`, admin.userID, start).Scan(&n)
	if n != 1 {
		t.Fatalf("administrator inbox entries: %d", n)
	}
	var mc int
	db.QueryRow(`SELECT must_change_password FROM users WHERE username='v1460-victim'`).Scan(&mc)
	if mc != 1 {
		t.Fatal("the reset password must be changed at the next sign-in")
	}
	// turned off: no alert
	setSetting("alert_cli_reset", "0")
	t.Cleanup(func() { setSetting("alert_cli_reset", "1") })
	os.Stdout = nil
	if err := resetPasswordCLI("v1460-victim", false); err != nil {
		t.Fatal(err)
	}
	os.Stdout = old
	db.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE action='security.alert' AND target='cli_reset' AND details LIKE '%v1460-victim%'`).Scan(&n)
	if n != 1 {
		t.Fatalf("alert raised although alert_cli_reset is off: %d", n)
	}
}

// quietCLI runs a CLI function with its standard output captured.
func quietCLI(t *testing.T, f func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func cliHash(name string) string {
	var h string
	db.QueryRow(`SELECT password_hash FROM users WHERE username=?`, name).Scan(&h)
	return h
}

// Business level: a reset works only while no WRM server uses the database.
func TestV1460CLIResetOnlyWhileStopped(t *testing.T) {
	srv := newTestServer(t)
	newTestUser(t, srv, "v1460-stop", false)
	setSetting("cli_password_reset", "stopped")
	t.Cleanup(func() { setSetting("cli_password_reset", "alert"); db.Exec(`DELETE FROM server_heartbeat`) })
	db.Exec(`INSERT OR REPLACE INTO server_heartbeat (id, pid, host, at) VALUES (1, 4242, 'demo-host', ?)`, time.Now().UTC().Format(time.RFC3339))
	before := cliHash("v1460-stop")
	_, err := quietCLI(t, func() error { return resetPasswordCLI("v1460-stop", false) })
	if err == nil || !strings.Contains(err.Error(), "WRM is running") {
		t.Fatalf("reset allowed while running: %v", err)
	}
	if cliHash("v1460-stop") != before {
		t.Fatal("the password changed although the reset was refused")
	}
	// a stale heartbeat (crashed server) does not block
	db.Exec(`UPDATE server_heartbeat SET at=? WHERE id=1`, time.Now().UTC().Add(-5*time.Minute).Format(time.RFC3339))
	out, err := quietCLI(t, func() error { return resetPasswordCLI("v1460-stop", false) })
	if err != nil || !strings.Contains(out, "New temporary password") {
		t.Fatalf("reset while stopped: %v %q", err, out)
	}
	// the heartbeat goroutine writes and removes its row
	db.Exec(`DELETE FROM server_heartbeat`)
	stop, done := make(chan struct{}), make(chan struct{})
	go func() { runServerHeartbeat(stop); close(done) }()
	deadline := time.Now().Add(3 * time.Second)
	for on, _ := serverRunning(); !on; on, _ = serverRunning() {
		if time.Now().After(deadline) {
			t.Fatal("heartbeat not written")
		}
		time.Sleep(20 * time.Millisecond)
	}
	close(stop)
	<-done
	if on, _ := serverRunning(); on {
		t.Fatal("heartbeat still there after the server stopped")
	}
}

// Enterprise level: the command files a request, which another administrator approves,
// the user's own 2FA code confirms, or which takes effect after an hour.
func TestV1460CLIResetConfirm(t *testing.T) {
	srv := newTestServer(t)
	adminA := newTestUser(t, srv, "v1460-adm-a", true)
	adminB := newTestUser(t, srv, "v1460-adm-b", true)
	newTestUser(t, srv, "v1460-conf", false)
	setSetting("cli_password_reset", "confirm")
	t.Cleanup(func() { setSetting("cli_password_reset", "alert") })

	before := cliHash("v1460-conf")
	out, err := quietCLI(t, func() error { return resetPasswordCLI("v1460-conf", false) })
	if err != nil || !strings.Contains(out, "was filed") || cliHash("v1460-conf") != before {
		t.Fatalf("first run must only file a request: %v %q", err, out)
	}
	var id int64
	db.QueryRow(`SELECT id FROM cli_reset_requests WHERE username='v1460-conf' AND status='pending'`).Scan(&id)
	if id == 0 {
		t.Fatal("no pending request")
	}
	// running again too early: still waiting, no second request
	out, _ = quietCLI(t, func() error { return resetPasswordCLI("v1460-conf", false) })
	if !strings.Contains(out, "is waiting") || cliHash("v1460-conf") != before {
		t.Fatalf("second run before approval: %q", out)
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM cli_reset_requests WHERE username='v1460-conf'`).Scan(&n)
	if n != 1 {
		t.Fatalf("requests: %d", n)
	}
	// the list shows it; a non-admin may not see it
	var list struct {
		Mode     string                   `json:"mode"`
		Requests []map[string]interface{} `json:"requests"`
	}
	if code := adminA.jsonDo("GET", "/api/admin/cli-resets", nil, &list); code != 200 || list.Mode != "confirm" || len(list.Requests) == 0 {
		t.Fatalf("list: %d %+v", code, list)
	}
	user := newTestUser(t, srv, "v1460-plain", false)
	if code := user.jsonDo("POST", "/api/admin/cli-resets/"+itoa64(id)+"/approve", nil, nil); code != 403 {
		t.Fatalf("non-admin approve: %d", code)
	}
	// approve by another administrator
	if code := adminA.jsonDo("POST", "/api/admin/cli-resets/"+itoa64(id)+"/approve", nil, nil); code != 200 {
		t.Fatalf("approve: %d", code)
	}
	if code := adminB.jsonDo("POST", "/api/admin/cli-resets/"+itoa64(id)+"/reject", nil, nil); code != 409 {
		t.Fatalf("decided twice: %d", code)
	}
	out, err = quietCLI(t, func() error { return resetPasswordCLI("v1460-conf", false) })
	if err != nil || !strings.Contains(out, "New temporary password") || cliHash("v1460-conf") == before {
		t.Fatalf("approved reset: %v %q", err, out)
	}

	// an administrator cannot approve a reset of their own account
	quietCLI(t, func() error { return resetPasswordCLI("v1460-adm-a", false) })
	db.QueryRow(`SELECT id FROM cli_reset_requests WHERE username='v1460-adm-a' AND status='pending'`).Scan(&id)
	if code := adminA.jsonDo("POST", "/api/admin/cli-resets/"+itoa64(id)+"/approve", nil, nil); code != 403 {
		t.Fatalf("self approve: %d", code)
	}
	// rejected: refused also after an hour (a new request is filed instead)
	if code := adminB.jsonDo("POST", "/api/admin/cli-resets/"+itoa64(id)+"/reject", nil, nil); code != 200 {
		t.Fatalf("reject: %d", code)
	}
	db.Exec(`UPDATE cli_reset_requests SET created_at=? WHERE id=?`, time.Now().UTC().Add(-2*time.Hour).Format(time.RFC3339), id)
	beforeA := cliHash("v1460-adm-a")
	out, _ = quietCLI(t, func() error { return resetPasswordCLI("v1460-adm-a", false) })
	if cliHash("v1460-adm-a") != beforeA || !strings.Contains(out, "was filed") {
		t.Fatalf("rejected request took effect: %q", out)
	}
	// after an hour a pending request takes effect
	db.Exec(`UPDATE cli_reset_requests SET created_at=? WHERE username='v1460-adm-a' AND status='pending'`, time.Now().UTC().Add(-61*time.Minute).Format(time.RFC3339))
	out, err = quietCLI(t, func() error { return resetPasswordCLI("v1460-adm-a", false) })
	if err != nil || cliHash("v1460-adm-a") == beforeA {
		t.Fatalf("delayed reset: %v %q", err, out)
	}
	// a request without -reset-2fa cannot be used to turn off 2FA
	quietCLI(t, func() error { return resetPasswordCLI("v1460-adm-b", false) })
	db.Exec(`UPDATE cli_reset_requests SET created_at=? WHERE username='v1460-adm-b' AND status='pending'`, time.Now().UTC().Add(-61*time.Minute).Format(time.RFC3339))
	if _, err := quietCLI(t, func() error { return resetPasswordCLI("v1460-adm-b", true) }); err == nil {
		t.Fatal("2FA turned off with a request that did not ask for it")
	}

	// the user's own 2FA code confirms at once; a wrong one is refused
	secret := "JBSWY3DPEHPK3PXP"
	db.Exec(`UPDATE users SET totp_secret=?, totp_enabled=1, totp_last_step=0 WHERE username='v1460-conf'`, encryptValue(secret))
	before = cliHash("v1460-conf")
	if _, err := quietCLI(t, func() error { return resetPasswordCLIWith("v1460-conf", false, "000000") }); err == nil || cliHash("v1460-conf") != before {
		t.Fatal("wrong 2FA code accepted")
	}
	raw, _ := decodeTOTPSecret(secret)
	code := totpCodeAt(raw, time.Now().Unix()/totpPeriod)
	out, err = quietCLI(t, func() error { return resetPasswordCLIWith("v1460-conf", false, code) })
	if err != nil || cliHash("v1460-conf") == before {
		t.Fatalf("2FA-confirmed reset: %v %q", err, out)
	}
	// the same code cannot be used twice
	before = cliHash("v1460-conf")
	if _, err := quietCLI(t, func() error { return resetPasswordCLIWith("v1460-conf", false, code) }); err == nil || cliHash("v1460-conf") != before {
		t.Fatal("2FA code replayed")
	}
}

func itoa64(n int64) string { return strconv.FormatInt(n, 10) }

// The command runs in its caller's environment: WRM_CLI_PASSWORD_RESET and
// WRM_ALERT_CLI_RESET there may make it stricter, never looser.
func TestV1460CLIResetEnvOnlyStricter(t *testing.T) {
	setSetting("cli_password_reset", "confirm")
	t.Cleanup(func() { setSetting("cli_password_reset", "alert") })
	t.Setenv("WRM_CLI_PASSWORD_RESET", "alert")
	t.Setenv("WRM_ALERT_CLI_RESET", "0")
	if m := cliResetMode(); m != "confirm" {
		t.Fatalf("environment loosened the mode: %s", m)
	}
	if !cliResetAlertOn() {
		t.Fatal("environment turned the alert off")
	}
	os.Unsetenv("WRM_CLI_PASSWORD_RESET")
	setSetting("cli_password_reset", "alert")
	t.Setenv("WRM_CLI_PASSWORD_RESET", "stopped")
	if m := cliResetMode(); m != "stopped" {
		t.Fatalf("environment could not tighten the mode: %s", m)
	}
}
