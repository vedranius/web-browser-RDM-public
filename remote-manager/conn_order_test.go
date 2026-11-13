package main

import (
	"net/http/httptest"
	"testing"
)

func connIDsOf(t *testing.T, u *testClient, folder *int) []int {
	t.Helper()
	var list []connView
	if code := u.jsonDo("GET", "/api/connections", nil, &list); code != 200 {
		t.Fatal(code)
	}
	out := []int{}
	for _, c := range list {
		if intPtrEq(c.FolderID, folder) {
			out = append(out, c.ID)
		}
	}
	return out
}

func TestConnectionOrderPerUser(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "conn-order", false)
	other := newTestUser(t, srv, "conn-order-other", false)
	var fd Folder
	if code := u.jsonDo("POST", "/api/folders", map[string]string{"name": "Demo"}, &fd); code != 201 {
		t.Fatal(code)
	}
	ids := map[string]int{}
	for _, n := range []string{"demo-c", "demo-a", "demo-b"} {
		var c connView
		if code := u.jsonDo("POST", "/api/connections", map[string]interface{}{"name": n, "protocol": "SSH", "host": n + ".example.com", "username": "root", "auth_method": "PASSWORD"}, &c); code != 200 && code != 201 {
			t.Fatal(n, code)
		}
		ids[n] = c.ID
	}
	a, b, c := ids["demo-a"], ids["demo-b"], ids["demo-c"]
	if got := connIDsOf(t, u, nil); !sameInts(got, []int{a, b, c}) {
		t.Fatalf("default order by name: %v", got)
	}
	if code := u.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{c, a}}, nil); code != 200 {
		t.Fatal(code)
	}
	if got := connIDsOf(t, u, nil); !sameInts(got, []int{c, a, b}) {
		t.Fatalf("placed first, rest by name: %v", got)
	}
	// someone else's connection, duplicates and empty lists are refused without a change
	if code := other.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{b, a}}, nil); code != 404 {
		t.Fatal("ownership", code)
	}
	if code := u.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{a, a}}, nil); code != 404 {
		t.Fatal("duplicate", code)
	}
	if code := u.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{}}, nil); code != 400 {
		t.Fatal("empty", code)
	}
	if got := connIDsOf(t, u, nil); !sameInts(got, []int{c, a, b}) {
		t.Fatalf("refused request changed the order: %v", got)
	}
	// moving into a folder drops the old position; the order there is its own
	if code := u.jsonDo("POST", "/api/connections/bulk", map[string]interface{}{"action": "move", "ids": []int{c}, "folder_id": fd.ID}, nil); code != 200 {
		t.Fatal(code)
	}
	if n := connOrderOf(u.userID); n[c] != 0 || n[a] != 2 {
		t.Fatalf("positions after move: %v", n)
	}
	if got := connIDsOf(t, u, &fd.ID); !sameInts(got, []int{c}) {
		t.Fatalf("folder: %v", got)
	}
	// reset → by name again; deleting a connection removes its position
	if code := u.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{b, a}}, nil); code != 200 {
		t.Fatal(code)
	}
	if got := connIDsOf(t, u, nil); !sameInts(got, []int{b, a}) {
		t.Fatalf("custom: %v", got)
	}
	if code := u.jsonDo("POST", "/api/connections/order", map[string]interface{}{"ids": []int{b}, "reset": true}, nil); code != 200 {
		t.Fatal(code)
	}
	if got := connIDsOf(t, u, nil); !sameInts(got, []int{a, b}) {
		t.Fatalf("after reset of b: %v", got)
	}
	if code := u.jsonDo("POST", "/api/connections/bulk", map[string]interface{}{"action": "delete", "ids": []int{a}}, nil); code != 200 {
		t.Fatal(code)
	}
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM connection_order WHERE conn_id=?`, a).Scan(&n)
	if n != 0 {
		t.Fatal("position of a deleted connection kept")
	}
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
