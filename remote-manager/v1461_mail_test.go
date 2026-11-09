package main

import (
	"io"
	"mime/quotedprintable"
	"net"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

func v1461MailSetup(t *testing.T, admin *testClient, addr string) {
	t.Helper()
	host, port, _ := net.SplitHostPort(addr)
	var out map[string]interface{}
	code := admin.jsonDo("PUT", "/api/admin/mail", map[string]interface{}{"enabled": true, "preset": "custom", "host": host, "port": v1461Port(port),
		"security": "none", "username": "wrm", "password": "smtp-secret-61", "from": "wrm@example.com", "from_name": "WRM", "public_url": "https://wrm.example.com/"}, &out)
	if code != 200 {
		t.Fatalf("mail settings: %d %v", code, out)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM mail_settings`); db.Exec(`DELETE FROM mail_outbox`) })
}

// Admin → Mail: only administrators, the password never leaves the server, a changed
// server needs the password again, the test message arrives, header injection is refused.
func TestV1461AdminMail(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1461-admin", true)
	user := newTestUser(t, srv, "v1461-user", false)
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")
	if code := user.jsonDo("GET", "/api/admin/mail", nil, nil); code != 403 {
		t.Fatalf("user reads mail settings: %d", code)
	}
	v1461MailSetup(t, admin, fake.addr)
	resp := admin.do("GET", "/api/admin/mail", nil, "")
	var buf strings.Builder
	b := make([]byte, 65536)
	n, _ := resp.Body.Read(b)
	buf.Write(b[:n])
	resp.Body.Close()
	if strings.Contains(buf.String(), "smtp-secret-61") || !strings.Contains(buf.String(), `"password_set":true`) {
		t.Fatalf("GET leaks or misses the password state: %s", buf.String())
	}
	host, port, _ := net.SplitHostPort(fake.addr)
	if code := admin.jsonDo("PUT", "/api/admin/mail", map[string]interface{}{"enabled": true, "host": "localhost", "port": v1461Port(port), "security": "none", "username": "wrm",
		"from": "wrm@example.com", "public_url": "https://wrm.example.com"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("stored password reused for another server: %d", code)
	}
	if code := admin.jsonDo("PUT", "/api/admin/mail", map[string]interface{}{"enabled": true, "host": host, "port": v1461Port(port), "security": "none", "username": "wrm",
		"from": "wrm@example.com", "from_name": "x\r\nBcc: evil@example.com", "public_url": "https://wrm.example.com"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("header injection in the sender name accepted: %d", code)
	}
	if code := admin.jsonDo("PUT", "/api/admin/mail", map[string]interface{}{"enabled": true, "host": "smtp.example.com", "port": 587, "security": "none", "username": "wrm", "password": "x",
		"from": "wrm@example.com", "public_url": "https://wrm.example.com"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("password over plain SMTP to a remote server accepted: %d", code)
	}
	if code := admin.jsonDo("POST", "/api/admin/mail/test", nil, &map[string]interface{}{}); code != 200 {
		t.Fatalf("test mail: %d", code)
	}
	msgs := fake.messages()
	if len(msgs) == 0 || !strings.Contains(v1461Body(msgs[len(msgs)-1]), "Mail works") || !strings.Contains(msgs[len(msgs)-1], "To: wrm@example.com") {
		t.Fatalf("test message not delivered: %v", msgs)
	}
}

// A user's address is verified by a link; security notices go only to a verified address.
func TestV1461UserEmailVerifyAndNotices(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1461-adm2", true)
	user := newTestUser(t, srv, "v1461-u2", false)
	fake := newFakeSMTP(t, "wrm", "smtp-secret-61")
	v1461MailSetup(t, admin, fake.addr)

	if code := user.jsonDo("PUT", "/api/me/email", map[string]string{"email": "not an address"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("invalid address accepted: %d", code)
	}
	if code := user.jsonDo("PUT", "/api/me/email", map[string]string{"email": "u2@example.com\r\nBcc: x@example.com"}, &map[string]interface{}{}); code != 400 {
		t.Fatalf("address with a line break accepted: %d", code)
	}
	// nothing goes to an unverified address
	if id := mailSecurityNotice(user.userID, "t", "x", "y", ""); id != 0 {
		t.Fatal("notice queued without an address")
	}
	if code := user.jsonDo("PUT", "/api/me/email", map[string]string{"email": "u2@example.com"}, &map[string]interface{}{}); code != 200 {
		t.Fatalf("set address: %d", code)
	}
	if mailDeliverDue(0) < 1 {
		t.Fatal("verification mail not delivered")
	}
	msgs := fake.messages()
	m := regexp.MustCompile(`https://wrm\.example\.com/api/auth/verify-email\?t=([0-9a-f]{48})`).FindStringSubmatch(v1461Body(msgs[len(msgs)-1]))
	if m == nil {
		t.Fatalf("no link in the mail: %q", msgs[len(msgs)-1])
	}
	var st struct {
		Email    string `json:"email"`
		Verified bool   `json:"verified"`
		Pending  string `json:"pending"`
	}
	user.jsonDo("GET", "/api/me/email", nil, &st)
	if st.Verified || st.Pending != "u2@example.com" {
		t.Fatalf("before the link: %+v", st)
	}
	noFollow := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := noFollow.Get(srv.URL + "/api/auth/verify-email?t=" + m[1])
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 303 || !strings.Contains(resp.Header.Get("Location"), "email_verified=ok") {
		t.Fatalf("verify: %d %s", resp.StatusCode, resp.Header.Get("Location"))
	}
	// the link works once
	resp, _ = noFollow.Get(srv.URL + "/api/auth/verify-email?t=" + m[1])
	resp.Body.Close()
	if !strings.Contains(resp.Header.Get("Location"), "email_verified=invalid") {
		t.Fatalf("link reused: %s", resp.Header.Get("Location"))
	}
	user.jsonDo("GET", "/api/me/email", nil, &st)
	if !st.Verified || st.Email != "u2@example.com" {
		t.Fatalf("after the link: %+v", st)
	}

	// security notices
	before := len(fake.messages())
	mailObserve(alertInput{Action: "auth.2fa_disabled", Username: "v1461-u2", IP: "198.51.100.7"})
	mailObserve(alertInput{Action: "admin.user_password_reset", Username: "v1461-adm2", Target: "v1461-u2"})
	mailDeliverDue(0)
	msgs = fake.messages()
	if len(msgs) != before+2 || !strings.Contains(v1461Body(msgs[before]), "198.51.100.7") || !strings.Contains(v1461Body(msgs[before+1]), "v1461-adm2") {
		t.Fatalf("notices: %d new", len(msgs)-before)
	}
	setSetting("mail_security_notices", "0")
	t.Cleanup(func() { setSetting("mail_security_notices", "1") })
	if mailSecurityNotice(user.userID, "t", "x", "y", "") != 0 {
		t.Fatal("notice queued with notices off")
	}
	setSetting("mail_security_notices", "1")

	// an expired link does not verify
	user.jsonDo("PUT", "/api/me/email", map[string]string{"email": "u2-new@example.com"}, &map[string]interface{}{})
	mailDeliverDue(0)
	msgs = fake.messages()
	var link []string
	for i := len(msgs) - 1; i >= 0 && link == nil; i-- {
		link = regexp.MustCompile(`verify-email\?t=([0-9a-f]{48})`).FindStringSubmatch(v1461Body(msgs[i]))
	}
	db.Exec(`UPDATE users SET email_token_at=? WHERE id=?`, time.Now().UTC().Add(-25*time.Hour).Format(time.RFC3339), user.userID)
	resp, _ = noFollow.Get(srv.URL + "/api/auth/verify-email?t=" + link[1])
	resp.Body.Close()
	if !strings.Contains(resp.Header.Get("Location"), "email_verified=expired") {
		t.Fatalf("expired link: %s", resp.Header.Get("Location"))
	}
	user.jsonDo("GET", "/api/me/email", nil, &st)
	if st.Email != "u2@example.com" {
		t.Fatalf("address changed by an expired link: %+v", st)
	}
	// the previous verified address was told about the change
	found := false
	for _, x := range msgs {
		if strings.Contains(v1461Body(x), "u2-new@example.com") && strings.Contains(x, "To: u2@example.com") {
			found = true
		}
	}
	if !found {
		t.Fatal("old address not told about the new one")
	}
	// admins see the address
	var list []map[string]interface{}
	admin.jsonDo("GET", "/api/admin/users", nil, &list)
	for _, u := range list {
		if u["username"] == "v1461-u2" && u["email"] != "u2@example.com" {
			t.Fatalf("admin list: %v", u)
		}
	}
}

// Undeliverable mail is retried and finally given up, without blocking anything.
func TestV1461MailRetry(t *testing.T) {
	srv := newTestServer(t)
	admin := newTestUser(t, srv, "v1461-adm3", true)
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	ln.Close() // nothing listens there
	v1461MailSetup(t, admin, addr)
	id := mailQueue(admin.userID, "ops@example.com", "t", "s", "b")
	if id == 0 {
		t.Fatal("not queued")
	}
	for i := 0; i < mailMaxAttempts; i++ {
		db.Exec(`UPDATE mail_outbox SET next_at=? WHERE id=?`, nowRFC(), id)
		mailDeliverDue(0)
	}
	var attempts, failed int
	var lastErr string
	db.QueryRow(`SELECT attempts, failed, last_error FROM mail_outbox WHERE id=?`, id).Scan(&attempts, &failed, &lastErr)
	if attempts != mailMaxAttempts || failed != 1 || lastErr == "" {
		t.Fatalf("attempts %d failed %d err %q", attempts, failed, lastErr)
	}
}

func v1461Port(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

// v1461Body is the decoded text of a received message.
func v1461Body(msg string) string {
	_, body, _ := strings.Cut(msg, "\r\n\r\n")
	b, _ := io.ReadAll(quotedprintable.NewReader(strings.NewReader(body)))
	return string(b)
}
