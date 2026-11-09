package main

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

func v1462Anon(t *testing.T, c *testClient) *testClient {
	return &testClient{t: t, srv: c.srv, cookie: &http.Cookie{Name: "v1462", Value: "anon"}}
}

func v1462ResetLimits() {
	forgotLimit.Lock()
	forgotLimit.m = map[string][]time.Time{}
	forgotLimit.Unlock()
	loginAttemptsMu.Lock()
	loginAttempts = map[string]*loginAttempt{}
	loginAttemptsMu.Unlock()
}

// v1462LastLink delivers the queue and returns the token of the last link matching kind
// ("reset", "activate" or "api/auth/activate").
func v1462LastLink(t *testing.T, fake *fakeSMTP, kind string) (string, string) {
	t.Helper()
	mailDeliverDue(0)
	msgs := fake.messages()
	if len(msgs) == 0 {
		return "", ""
	}
	body := v1461Body(msgs[len(msgs)-1])
	re := regexp.MustCompile(`https://wrm\.example\.com/\?` + kind + `=([0-9a-f]{48})`)
	if kind == "api/auth/activate" {
		re = regexp.MustCompile(`https://wrm\.example\.com/api/auth/activate\?t=([0-9a-f]{48})`)
	}
	m := re.FindStringSubmatch(body)
	if m == nil {
		return "", msgs[len(msgs)-1]
	}
	return m[1], msgs[len(msgs)-1]
}

func v1462Login(c *testClient, user, pw string) (int, map[string]interface{}) {
	out := map[string]interface{}{}
	code := c.jsonDo("POST", "/api/auth/login", map[string]string{"username": user, "password": pw}, &out)
	return code, out
}

// Forgot password: one-time link to the confirmed address, same answer for unknown accounts,
// old sessions end, a notice follows, the link cannot be reused.
func TestV1462ForgotPassword(t *testing.T) {
	v1462ResetLimits()
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1462-admin", true)
	user := newTestUser(t, srv, "v1462-user", false)
	anon := v1462Anon(t, admin)
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")

	var cfg map[string]interface{}
	anon.jsonDo("GET", "/api/auth/config", nil, &cfg)
	if cfg["password_reset"] != false {
		t.Fatalf("reset offered without mail: %v", cfg["password_reset"])
	}
	if code := anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-user"}, &map[string]interface{}{}); code != 403 {
		t.Fatalf("forgot without mail: %d", code)
	}
	v1461MailSetup(t, admin, fake.addr)
	anon.jsonDo("GET", "/api/auth/config", nil, &cfg)
	if cfg["password_reset"] != true {
		t.Fatalf("reset not offered: %v", cfg)
	}

	// no confirmed address: same answer, no mail
	before := len(fake.messages())
	if code := anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-user"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("forgot: %d", code)
	}
	if code := anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "nobody-v1462"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("forgot unknown: %d", code)
	}
	mailDeliverDue(0)
	if len(fake.messages()) != before {
		t.Fatal("mail sent to an account without a confirmed address")
	}

	db.Exec(`UPDATE users SET email='v1462u@example.com', email_verified=1 WHERE id=?`, user.userID)
	// by e-mail address, case-insensitive
	if code := anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "V1462U@Example.com"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("forgot by mail: %d", code)
	}
	tok, raw := v1462LastLink(t, fake, "reset")
	if tok == "" || !strings.Contains(raw, "To: v1462u@example.com") {
		t.Fatalf("no reset link: %q", raw)
	}
	// a second request at once sends nothing new
	n := len(fake.messages())
	anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-user"}, &map[string]interface{}{})
	mailDeliverDue(0)
	if len(fake.messages()) != n {
		t.Fatal("a second link within the gap")
	}

	var info map[string]interface{}
	anon.jsonDo("GET", "/api/auth/reset?t="+tok, nil, &info)
	if info["status"] != "ok" || info["kind"] != "reset" || info["username"] != "v1462-user" || info["totp"] != false {
		t.Fatalf("token info: %v", info)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "short"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("weak password accepted: %d", code)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "A-new-Passw0rd-1462"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("reset: %d", code)
	}
	if code := user.jsonDo("GET", "/api/auth/me", nil, &map[string]interface{}{}); code != 401 {
		t.Fatalf("old session survived: %d", code)
	}
	if code, out := v1462Login(anon, "v1462-user", "A-new-Passw0rd-1462"); code != 200 {
		t.Fatalf("sign-in with the new password: %d %v", code, out)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "Another-Passw0rd-1462"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("link reused: %d", code)
	}
	anon.jsonDo("GET", "/api/auth/reset?t="+tok, nil, &info)
	if info["status"] != "expired" {
		t.Fatalf("used link: %v", info)
	}
	mailDeliverDue(0)
	msgs := fake.messages()
	if !strings.Contains(v1461Body(msgs[len(msgs)-1]), "changed with a link") {
		t.Fatalf("no notice after the reset: %q", msgs[len(msgs)-1])
	}
	// an expired link
	db.Exec(`DELETE FROM auth_links WHERE user_id=?`, user.userID)
	old := authLinkCreate(user.userID, authLinkReset, "v1462u@example.com", -time.Minute, "")
	anon.jsonDo("GET", "/api/auth/reset?t="+old, nil, &info)
	if info["status"] != "expired" {
		t.Fatalf("expired link: %v", info)
	}
}

// Who can use it: accounts with 2FA give a code, administrators are excluded by "users", the
// policy can turn it off, disabled and directory accounts never get a link.
func TestV1462ForgotRules(t *testing.T) {
	v1462ResetLimits()
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1462-adm2", true)
	user := newTestUser(t, srv, "v1462-u2", false)
	anon := v1462Anon(t, admin)
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")
	v1461MailSetup(t, admin, fake.addr)
	db.Exec(`UPDATE users SET email='v1462a2@example.com', email_verified=1 WHERE id=?`, admin.userID)
	db.Exec(`UPDATE users SET email='v1462u2@example.com', email_verified=1 WHERE id=?`, user.userID)
	t.Cleanup(func() { setSetting("password_reset_mail", "all") })

	// 2FA: the link alone is not enough
	secret := "JBSWY3DPEHPK3PXP"
	db.Exec(`UPDATE users SET totp_secret=?, totp_enabled=1, totp_last_step=0 WHERE id=?`, encryptValue(secret), user.userID)
	anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-u2"}, &map[string]interface{}{})
	tok, raw := v1462LastLink(t, fake, "reset")
	if tok == "" {
		t.Fatalf("no link: %q", raw)
	}
	var info map[string]interface{}
	anon.jsonDo("GET", "/api/auth/reset?t="+tok, nil, &info)
	if info["totp"] != true {
		t.Fatalf("2FA not asked: %v", info)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "A-new-Passw0rd-1462"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("reset without a code: %d", code)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "A-new-Passw0rd-1462", "code": "000000"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("reset with a wrong code: %d", code)
	}
	key, _ := decodeTOTPSecret(secret)
	code := totpCodeAt(key, time.Now().Unix()/totpPeriod)
	if c := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "A-new-Passw0rd-1462", "code": code}, &map[string]interface{}{}); c != 200 {
		t.Fatalf("reset with the code: %d", c)
	}

	// "users": administrators are excluded
	if _, err := setSetting("password_reset_mail", "users"); err != nil {
		t.Fatal(err)
	}
	mailDeliverDue(0) // the notice about the reset above
	n := len(fake.messages())
	anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-adm2"}, &map[string]interface{}{})
	mailDeliverDue(0)
	if len(fake.messages()) != n {
		t.Fatal("administrator got a link under 'users'")
	}
	// disabled accounts
	db.Exec(`UPDATE users SET totp_enabled=0, disabled=1 WHERE id=?`, user.userID)
	db.Exec(`DELETE FROM auth_links WHERE user_id=?`, user.userID)
	anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-u2"}, &map[string]interface{}{})
	mailDeliverDue(0)
	if len(fake.messages()) != n {
		t.Fatal("disabled account got a link")
	}
	// off
	setSetting("password_reset_mail", "off")
	var cfg map[string]interface{}
	anon.jsonDo("GET", "/api/auth/config", nil, &cfg)
	if cfg["password_reset"] != false {
		t.Fatal("offered while off")
	}
	if c := anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "v1462-u2"}, &map[string]interface{}{}); c != 403 {
		t.Fatalf("forgot while off: %d", c)
	}
	// the request limit per address
	setSetting("password_reset_mail", "all")
	v1462ResetLimits()
	limited := false
	for i := 0; i < forgotPerWindow+1; i++ {
		if anon.jsonDo("POST", "/api/auth/forgot", map[string]string{"login": "x"}, &map[string]interface{}{}) == 429 {
			limited = true
		}
	}
	if !limited {
		t.Fatal("no request limit")
	}
	v1462ResetLimits()
}

// An administrator invites a user: no password is shown, the link sets one and confirms the
// address.
func TestV1462Invite(t *testing.T) {
	v1462ResetLimits()
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1462-adm3", true)
	anon := v1462Anon(t, admin)
	var out map[string]interface{}
	if code := admin.jsonDo("POST", "/api/admin/users", map[string]interface{}{"username": "v1462-inv0", "email": "inv0@example.com", "invite": true}, &out); code != 400 {
		t.Fatalf("invitation without mail: %d %v", code, out)
	}
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")
	v1461MailSetup(t, admin, fake.addr)
	if code := admin.jsonDo("POST", "/api/admin/users", map[string]interface{}{"username": "v1462-inv1", "email": "bad\r\naddress", "invite": true}, &out); code != 400 {
		t.Fatalf("bad address: %d", code)
	}
	out = nil
	if code := admin.jsonDo("POST", "/api/admin/users", map[string]interface{}{"username": "v1462-inv", "email": "inv@example.com", "invite": true}, &out); code != 200 {
		t.Fatalf("invite: %d %v", code, out)
	}
	if out["generated_password"] != "" || out["invited"] != "inv@example.com" {
		t.Fatalf("invite answer: %v", out)
	}
	tok, raw := v1462LastLink(t, fake, "activate")
	if tok == "" || !strings.Contains(raw, "To: inv@example.com") {
		t.Fatalf("no invitation: %q", raw)
	}
	var info map[string]interface{}
	anon.jsonDo("GET", "/api/auth/reset?t="+tok, nil, &info)
	if info["kind"] != "activate" || info["username"] != "v1462-inv" {
		t.Fatalf("invitation info: %v", info)
	}
	if code := anon.jsonDo("POST", "/api/auth/reset", map[string]string{"token": tok, "password": "Invited-Passw0rd-1462"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("set password: %d", code)
	}
	var email string
	var verified, must int
	db.QueryRow(`SELECT email, email_verified, must_change_password FROM users WHERE username='v1462-inv'`).Scan(&email, &verified, &must)
	if email != "inv@example.com" || verified != 1 || must != 0 {
		t.Fatalf("after the invitation: %q %d %d", email, verified, must)
	}
	if code, o := v1462Login(anon, "v1462-inv", "Invited-Passw0rd-1462"); code != 200 {
		t.Fatalf("sign-in: %d %v", code, o)
	}
	// an e-mail address without an invitation starts the usual confirmation
	out = nil
	if code := admin.jsonDo("POST", "/api/admin/users", map[string]interface{}{"username": "v1462-mail", "email": "m@example.com"}, &out); code != 200 || out["generated_password"] == "" {
		t.Fatalf("create with address: %d %v", code, out)
	}
	var pending string
	db.QueryRow(`SELECT email_pending FROM users WHERE username='v1462-mail'`).Scan(&pending)
	if pending != "m@example.com" {
		t.Fatalf("pending address: %q", pending)
	}
}

// Self-registration with mail: the account starts once the address is confirmed; accounts
// never confirmed are removed.
func TestV1462RegisterConfirm(t *testing.T) {
	v1462ResetLimits()
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1462-adm4", true)
	anon := v1462Anon(t, admin)
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")
	v1461MailSetup(t, admin, fake.addr)
	prev := getSetting("registration")
	setSetting("registration", "open")
	t.Cleanup(func() { setSetting("registration", prev) })

	var cfg map[string]interface{}
	anon.jsonDo("GET", "/api/auth/config", nil, &cfg)
	if cfg["registration_email"] != true {
		t.Fatalf("registration does not ask for mail: %v", cfg)
	}
	if code := anon.jsonDo("POST", "/api/auth/register", map[string]string{"username": "v1462-reg", "password": "Reg-Passw0rd-1462"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("registration without an address: %d", code)
	}
	var out map[string]interface{}
	if code := anon.jsonDo("POST", "/api/auth/register", map[string]string{"username": "v1462-reg", "password": "Reg-Passw0rd-1462", "email": "reg@example.com"}, &out); code != 200 || out["pending_confirmation"] != true {
		t.Fatalf("register: %d %v", code, out)
	}
	if code, o := v1462Login(anon, "v1462-reg", "Reg-Passw0rd-1462"); code != 403 || !strings.Contains(o["error"].(string), "Confirm") {
		t.Fatalf("sign-in before confirmation: %d %v", code, o)
	}
	tok, raw := v1462LastLink(t, fake, "api/auth/activate")
	if tok == "" {
		t.Fatalf("no confirmation link: %q", raw)
	}
	noFollow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := noFollow.Get(srv.URL + "/api/auth/activate?t=" + tok)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 303 || !strings.Contains(resp.Header.Get("Location"), "activated=ok") {
		t.Fatalf("activate: %d %s", resp.StatusCode, resp.Header.Get("Location"))
	}
	resp, _ = noFollow.Get(srv.URL + "/api/auth/activate?t=" + tok)
	resp.Body.Close()
	if !strings.Contains(resp.Header.Get("Location"), "activated=expired") {
		t.Fatalf("link reused: %s", resp.Header.Get("Location"))
	}
	if code, o := v1462Login(anon, "v1462-reg", "Reg-Passw0rd-1462"); code != 200 {
		t.Fatalf("sign-in after confirmation: %d %v", code, o)
	}
	var email string
	var verified int
	db.QueryRow(`SELECT email, email_verified FROM users WHERE username='v1462-reg'`).Scan(&email, &verified)
	if email != "reg@example.com" || verified != 1 {
		t.Fatalf("address after confirmation: %q %d", email, verified)
	}

	// an account never confirmed is removed after a week
	old := time.Now().Add(-8 * 24 * time.Hour).UTC().Format(time.RFC3339)
	db.Exec(`INSERT INTO users (username, password_hash, is_admin, disabled, activation_pending, created_at) VALUES ('v1462-stale','x',0,1,1,?)`, old)
	authLinksCleanup()
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE username='v1462-stale'`).Scan(&n)
	if n != 0 {
		t.Fatal("stale registration kept")
	}
	// without mail, registration works as before
	db.Exec(`DELETE FROM mail_settings`)
	out = nil
	if code := anon.jsonDo("POST", "/api/auth/register", map[string]string{"username": "v1462-reg2", "password": "Reg-Passw0rd-1462"}, &out); code != 200 || out["pending_confirmation"] == true {
		t.Fatalf("register without mail: %d %v", code, out)
	}
}
