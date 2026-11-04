package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestLinkSharePolicyDoesNotWiden covers SEC-07: turning guest links off must deny the
// link, not silently turn it into "any signed-in account".
func TestLinkSharePolicyDoesNotWiden(t *testing.T) {
	old := getSetting("allow_link_shares")
	if _, err := setSetting("allow_link_shares", "0"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { setSetting("allow_link_shares", old) })

	share := ShareLink{ID: 1, OwnerID: 999999, AccessMode: "link", Active: true, Name: "test"}
	_, d := resolveShareAccessCtx(share, accessCtx{})
	if d == nil {
		t.Fatal("a disabled guest link was accepted")
	}
	if d.Reason != "link_shares_disabled" {
		t.Fatalf("denied with %q, want link_shares_disabled", d.Reason)
	}

	if _, err := setSetting("allow_link_shares", "1"); err != nil {
		t.Fatal(err)
	}
	// With the policy on, the same link is not refused by the policy branch (it may still
	// ask for the guest cookie or a password).
	if _, d := resolveShareAccessCtx(share, accessCtx{}); d != nil && d.Reason == "link_shares_disabled" {
		t.Fatal("a link is refused while guest links are enabled")
	}
}

func TestRestrictedSessionCannotFallBackToShareLink(t *testing.T) {
	srv := newTestServer(t)
	owner := newTestUser(t, srv, "restricted-share-owner", true)
	restricted := newTestUser(t, srv, "restricted-share-user", false)
	conn := owner.addConnection("restricted-share-conn", "host.example.com")
	share := mustShare(t, owner.userID, "link", RoleOperator, "")
	if _, err := db.Exec(`INSERT INTO share_items (share_id, connection_id) VALUES (?,?)`, share.ID, conn); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE auth_sessions SET restricted='password' WHERE user_id=?`, restricted.userID); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/api/terminal?id=1&share_token="+share.Token, nil)
	req.AddCookie(restricted.cookie)
	if _, code, reason := authorizeConnection(req, conn, "terminal"); code != 403 {
		t.Fatalf("restricted session fell back to share access: status=%d reason=%q", code, reason)
	}
}

func TestShareAccessExchangeIssuesPerShareCookieAndRotationInvalidatesIt(t *testing.T) {
	srv := newTestServer(t)
	owner := newTestUser(t, srv, "share-cookie-owner", true)
	conn := owner.addConnection("share-cookie-conn", "host.example.com")
	share := mustShare(t, owner.userID, "link", RoleOperator, "")
	if _, err := db.Exec(`INSERT INTO share_items (share_id, connection_id) VALUES (?,?)`, share.ID, conn); err != nil {
		t.Fatal(err)
	}

	client := srv.Client()
	blockedReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/share/access", strings.NewReader(`{"token":"`+share.Token+`"}`))
	blockedResp, err := client.Do(blockedReq)
	if err != nil {
		t.Fatal(err)
	}
	blockedResp.Body.Close()
	if blockedResp.StatusCode != http.StatusForbidden {
		t.Fatalf("exchange without CSRF proof returned %d", blockedResp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/share/access", strings.NewReader(`{"token":"`+share.Token+`"}`))
	req.Header.Set("X-WRM-Request", "1")
	req.Header.Set("Origin", srv.URL)
	rr, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rr.Body.Close()
	body, _ := io.ReadAll(rr.Body)
	if rr.StatusCode != http.StatusOK {
		t.Fatalf("exchange status=%d body=%s", rr.StatusCode, body)
	}
	var out struct {
		ShareID int `json:"share_id"`
	}
	if err := json.Unmarshal(body, &out); err != nil || out.ShareID != share.ID {
		t.Fatalf("exchange response=%q err=%v", body, err)
	}
	var accessCookie, guestCookie *http.Cookie
	for _, c := range rr.Cookies() {
		if c.Name == shareAccessCookieName(share.ID) {
			accessCookie = c
		}
		if c.Name == guestCookieName {
			guestCookie = c
		}
	}
	if accessCookie == nil || !accessCookie.HttpOnly || accessCookie.SameSite != http.SameSiteStrictMode {
		t.Fatalf("missing hardened share cookie: %#v", accessCookie)
	}
	if guestCookie == nil {
		t.Fatal("exchange did not initialize guest identity")
	}

	accessReq := httptest.NewRequest(http.MethodGet, "/api/terminal?id=1&share_id="+strconv.Itoa(share.ID), nil)
	accessReq.AddCookie(accessCookie)
	accessReq.AddCookie(guestCookie)
	if _, code, reason := authorizeConnection(accessReq, conn, PermTerminal); code != 0 {
		t.Fatalf("cookie share access denied: %d %s", code, reason)
	}

	if _, err := db.Exec(`UPDATE share_links SET token=? WHERE id=?`, "rotated-token", share.ID); err != nil {
		t.Fatal(err)
	}
	if _, code, _ := authorizeConnection(accessReq, conn, PermTerminal); code == 0 {
		t.Fatal("share cookie survived token rotation")
	}
}

func TestExactShareWSOriginChecksSchemeHostAndPort(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://wrm.example:8080/ws/share/12", nil)
	r.Host = "wrm.example:8080"
	r.Header.Set("Origin", "http://wrm.example:8080")
	if !exactShareWSOrigin(r) {
		t.Fatal("exact same-origin request was rejected")
	}
	r.Header.Set("Origin", "http://wrm.example:8081")
	if exactShareWSOrigin(r) {
		t.Fatal("same host on another port was accepted")
	}
	r.Header.Set("Origin", "https://wrm.example:8080")
	if exactShareWSOrigin(r) {
		t.Fatal("different origin scheme was accepted")
	}
	r.Header.Set("Origin", "")
	if exactShareWSOrigin(r) {
		t.Fatal("missing Origin was accepted")
	}
}
