package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestAPIRoutesDenyByDefault covers ARCH-01: /api and /ws are protected unless the path is
// public, so a handler (or a route added later) that forgets requireAuth is refused.
func TestAPIRoutesDenyByDefault(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()

	for _, p := range []string{"/api/connections", "/api/admin/users", "/api/does-not-exist", "/ws/ssh", "/ws/events"} {
		resp, err := http.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s without a session: %d, want 401", p, resp.StatusCode)
		}
	}

	// Sign-in and token-based routes must stay reachable without a session cookie.
	for _, p := range []string{"/api/version", "/api/auth/config", "/api/share/missing"} {
		resp, err := http.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized {
			t.Errorf("%s is public but answered 401", p)
		}
	}

	u := newTestUser(t, srv, "deny-by-default", false)
	if code := u.jsonDo("GET", "/api/connections", nil, nil); code != http.StatusOK {
		t.Fatalf("signed-in /api/connections: %d, want 200", code)
	}
}

func TestAuthGateAllowsShareCookieOnlyOnShareAwareRoutes(t *testing.T) {
	srv := newTestServer(t)
	owner := newTestUser(t, srv, "share-gate-owner", true)
	conn := owner.addConnection("share-gate-conn", "host.example.com")
	share := mustShare(t, owner.userID, "link", RoleOperator, "")
	if _, err := db.Exec(`INSERT INTO share_items (share_id, connection_id) VALUES (?,?)`, share.ID, conn); err != nil {
		t.Fatal(err)
	}
	seedReq := httptest.NewRequest(http.MethodGet, "http://wrm.example/", nil)
	seedRes := httptest.NewRecorder()
	ensureGuestCookie(seedRes, seedReq)
	setShareAccessCookie(seedRes, seedReq, share)
	var accessCookie, guestCookie *http.Cookie
	for _, c := range seedRes.Result().Cookies() {
		if c.Name == shareAccessCookieName(share.ID) {
			accessCookie = c
		}
		if c.Name == guestCookieName {
			guestCookie = c
		}
	}

	gate := authGate(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	allowed := []string{"/api/remote/list", "/api/remote/list-recursive", "/api/remote/download", "/api/remote/preview",
		"/api/remote/download-dir", "/api/remote/upload", "/api/remote/mkdir", "/api/remote/chmod", "/api/remote/delete",
		"/api/remote/rename", "/api/remote/transfer", "/api/remote/search", "/ws/ssh", "/ws/desktop"}
	for _, p := range allowed {
		req := httptest.NewRequest(http.MethodGet, p+"?share_id="+strconv.Itoa(share.ID), nil)
		req.AddCookie(accessCookie)
		res := httptest.NewRecorder()
		gate.ServeHTTP(res, req)
		if res.Code != http.StatusNoContent {
			t.Errorf("%s with valid share cookie: status=%d, want 204", p, res.Code)
		}
	}
	connReq := httptest.NewRequest(http.MethodGet, "/api/terminal?id="+strconv.Itoa(conn)+"&share_id="+strconv.Itoa(share.ID), nil)
	connReq.AddCookie(accessCookie)
	connReq.AddCookie(guestCookie)
	if _, code, _ := authorizeConnection(connReq, conn, PermTerminal); code != 0 {
		t.Fatalf("operator share cookie did not authorize its connection: %d", code)
	}
	if _, err := db.Exec(`UPDATE share_links SET guest_role=? WHERE id=?`, RoleViewer, share.ID); err != nil {
		t.Fatal(err)
	}
	if _, code, _ := authorizeConnection(connReq, conn, PermTerminal); code != http.StatusForbidden {
		t.Fatalf("share cookie bypassed the current role: %d, want 403", code)
	}
	for _, target := range []string{"/api/remote/list", "/api/admin/users", "/api/does-not-exist"} {
		req := httptest.NewRequest(http.MethodGet, target+"?share_id="+strconv.Itoa(share.ID), nil)
		if target == "/api/remote/list" {
			req.AddCookie(&http.Cookie{Name: accessCookie.Name, Value: "invalid"})
		} else {
			req.AddCookie(accessCookie)
		}
		res := httptest.NewRecorder()
		gate.ServeHTTP(res, req)
		if res.Code != http.StatusUnauthorized {
			t.Errorf("%s with invalid or unrelated share cookie: status=%d, want 401", target, res.Code)
		}
	}
}

func TestSharePageServesCleanAndLegacyPaths(t *testing.T) {
	for _, path := range []string{"/share", "/share/example-token"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res := httptest.NewRecorder()
		newRouter().ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("GET %s: status %d, want %d", path, res.Code, http.StatusOK)
		}
		if !strings.Contains(res.Body.String(), "WRM") {
			t.Errorf("GET %s did not serve the app", path)
		}
	}
}
