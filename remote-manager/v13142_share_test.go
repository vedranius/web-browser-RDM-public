package main

import (
	"net/http/httptest"
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
