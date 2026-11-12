package main

import (
	"testing"
)

// TestShareTunnelEditSwitch covers the #192 follow-up: "add and edit tunnels" is a share
// switch of its own, separate from "start and stop tunnels".
func TestShareTunnelEditSwitch(t *testing.T) {
	srv := newTestServer(t)
	owner := newTestUser(t, srv, "tun-owner", false)
	op := newTestUser(t, srv, "tun-operator", false)
	viewer := newTestUser(t, srv, "tun-viewer", false)

	res, err := db.Exec(`INSERT INTO folders (name, user_id) VALUES ('demo-team', ?)`, owner.userID)
	if err != nil {
		t.Fatal(err)
	}
	fid64, _ := res.LastInsertId()
	conn := owner.addConnection("demo-api", "demo-api.example.com")
	db.Exec(`UPDATE connections SET folder_id=? WHERE id=?`, fid64, conn)
	private := owner.addConnection("private-db", "db.example.com")

	share := mustShare(t, owner.userID, "members", RoleViewer, "")
	db.Exec(`INSERT INTO share_items (share_id, folder_id) VALUES (?,?)`, share.ID, fid64)
	db.Exec(`INSERT INTO share_members (share_id, user_id, role) VALUES (?,?,?),(?,?,?)`, share.ID, op.userID, RoleOperator, share.ID, viewer.userID, RoleViewer)
	url := "/api/share-team/" + share.Token + "/connections/" + itoa(conn) + "/tunnels"
	one := []map[string]interface{}{{"name": "demo-web", "kind": "local", "bind_port": 0, "target_host": "127.0.0.1", "target_port": 8080}}

	// Start/stop alone does not allow adding or editing.
	db.Exec(`UPDATE share_links SET allow_tunnels=1 WHERE id=?`, share.ID)
	if code := op.jsonDo("PUT", url, one, nil); code != 403 {
		t.Fatalf("edit allowed with only start/stop: %d", code)
	}
	if code := op.jsonDo("GET", "/api/share-team/"+share.Token+"/tunnels", nil, nil); code != 200 {
		t.Fatalf("start/stop list: %d", code)
	}

	// Add/edit alone: tunnels can be configured, but not started.
	db.Exec(`UPDATE share_links SET allow_tunnels=0, allow_tunnel_edit=1 WHERE id=?`, share.ID)
	if code := op.jsonDo("GET", "/api/share-team/"+share.Token+"/tunnels", nil, nil); code != 403 {
		t.Fatalf("start/stop allowed with only add/edit: %d", code)
	}
	var defs []tunnelDef
	if code := op.jsonDo("PUT", url, one, &defs); code != 200 || len(defs) != 1 {
		t.Fatalf("add tunnel: %d %v", code, defs)
	}
	var uid int
	db.QueryRow(`SELECT user_id FROM connection_tunnels WHERE id=?`, defs[0].ID).Scan(&uid)
	if uid != owner.userID {
		t.Fatalf("tunnel not owned by the share owner: %d", uid)
	}
	// Edit keeps the id; listen addresses beyond loopback stay admin-only for the member.
	edit := []map[string]interface{}{{"id": defs[0].ID, "name": "demo-web-2", "kind": "local", "target_host": "127.0.0.1", "target_port": 8081}}
	if code := op.jsonDo("PUT", url, edit, &defs); code != 200 || len(defs) != 1 || defs[0].Name != "demo-web-2" {
		t.Fatalf("edit tunnel: %d %v", code, defs)
	}
	edit[0]["bind_host"] = "0.0.0.0"
	if code := op.jsonDo("PUT", url, edit, nil); code != 400 {
		t.Fatalf("member listened on 0.0.0.0: %d", code)
	}

	// Viewers, guests and connections outside the share stay out.
	if code := viewer.jsonDo("PUT", url, one, nil); code != 403 {
		t.Fatalf("viewer edited tunnels: %d", code)
	}
	if code := op.jsonDo("GET", "/api/share-team/"+share.Token+"/connections/"+itoa(private)+"/tunnels", nil, nil); code != 404 {
		t.Fatalf("private connection reachable: %d", code)
	}
	s2, _ := getShareByID(share.ID)
	if !s2.AllowTunnelEdit || sharePerms(s2, RoleOperator, true)[PermEditTunnels] || !sharePerms(s2, RoleOperator, false)[PermEditTunnels] {
		t.Fatal("sharePerms edit_tunnels wrong")
	}

	// Removing all tunnels works too, and the owner's switch is in the share API.
	if code := op.jsonDo("PUT", url, []interface{}{}, &defs); code != 200 || len(defs) != 0 {
		t.Fatalf("remove tunnels: %d %v", code, defs)
	}
	var sh map[string]interface{}
	if code := owner.jsonDo("PUT", "/api/shares/"+itoa(share.ID), map[string]interface{}{"allow_tunnel_edit": false}, &sh); code != 200 {
		t.Fatalf("owner update: %d", code)
	}
	if code := op.jsonDo("PUT", url, one, nil); code != 403 {
		t.Fatalf("edit allowed after the owner turned it off: %d", code)
	}
}
