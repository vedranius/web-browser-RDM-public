package main

import (
	"testing"
)

// TestShareTeamManagement covers #192: the two share switches let signed-in operators
// add and edit the shared connections and folders and start/stop their tunnels.
func TestShareTeamManagement(t *testing.T) {
	srv := newTestServer(t)
	owner := newTestUser(t, srv, "team-owner", false)
	op := newTestUser(t, srv, "team-operator", false)
	viewer := newTestUser(t, srv, "team-viewer", false)
	outsider := newTestUser(t, srv, "team-outsider", false)

	res, err := db.Exec(`INSERT INTO folders (name, user_id) VALUES ('demo-team', ?)`, owner.userID)
	if err != nil {
		t.Fatal(err)
	}
	fid64, _ := res.LastInsertId()
	fid := int(fid64)
	res, _ = db.Exec(`INSERT INTO folders (name, user_id) VALUES ('private', ?)`, owner.userID)
	other64, _ := res.LastInsertId()
	conn := owner.addConnection("demo-api", "demo-api.example.com")
	db.Exec(`UPDATE connections SET folder_id=? WHERE id=?`, fid, conn)
	private := owner.addConnection("private-db", "db.example.com")
	db.Exec(`UPDATE connections SET folder_id=? WHERE id=?`, other64, private)

	share := mustShare(t, owner.userID, "members", RoleViewer, "")
	db.Exec(`INSERT INTO share_items (share_id, folder_id) VALUES (?,?)`, share.ID, fid)
	db.Exec(`INSERT INTO share_members (share_id, user_id, role) VALUES (?,?,?),(?,?,?)`, share.ID, op.userID, RoleOperator, share.ID, viewer.userID, RoleViewer)
	base := "/api/share-team/" + share.Token

	// Off by default: nobody may edit.
	if code := op.jsonDo("GET", base+"/connections/"+itoa(conn), nil, nil); code != 403 {
		t.Fatalf("edit allowed with the switch off: %d", code)
	}
	db.Exec(`UPDATE share_links SET allow_edit=1, allow_tunnels=1 WHERE id=?`, share.ID)

	// Viewers and outsiders stay out.
	if code := viewer.jsonDo("GET", base+"/connections/"+itoa(conn), nil, nil); code != 403 {
		t.Fatalf("viewer may edit: %d", code)
	}
	if code := outsider.jsonDo("GET", base+"/connections/"+itoa(conn), nil, nil); code < 400 {
		t.Fatalf("outsider may edit: %d", code)
	}
	// Connections outside the share stay hidden.
	if code := op.jsonDo("GET", base+"/connections/"+itoa(private), nil, nil); code != 404 {
		t.Fatalf("private connection reachable: %d", code)
	}

	var view map[string]interface{}
	if code := op.jsonDo("GET", base+"/connections/"+itoa(conn), nil, &view); code != 200 || view["has_password"] != true {
		t.Fatalf("GET: %d %v", code, view)
	}
	if _, ok := view["password"]; ok {
		t.Fatal("password sent to the member")
	}

	// Renaming keeps the password; changing the host clears it.
	if code := op.jsonDo("PUT", base+"/connections/"+itoa(conn), map[string]interface{}{"name": "demo-api-2"}, &view); code != 200 || view["has_password"] != true {
		t.Fatalf("rename: %d %v", code, view)
	}
	if code := op.jsonDo("PUT", base+"/connections/"+itoa(conn), map[string]interface{}{"host": "attacker.example.com"}, &view); code != 200 || view["has_password"] != false {
		t.Fatalf("host change kept the password: %d %v", code, view)
	}
	// Moving out of the share is refused.
	if code := op.jsonDo("PUT", base+"/connections/"+itoa(conn), map[string]interface{}{"folder_id": int(other64)}, nil); code != 400 {
		t.Fatalf("moved out of the share: %d", code)
	}
	// A credential reference of the owner stays untouchable.
	if _, err := db.Exec(`UPDATE connections SET auth_method='CREDENTIAL', credential_id=77, password='' WHERE id=?`, conn); err != nil {
		t.Fatal(err)
	}
	op.jsonDo("PUT", base+"/connections/"+itoa(conn), map[string]interface{}{"host": "evil.example.com"}, nil)
	var am string
	var cred *int
	db.QueryRow(`SELECT auth_method, credential_id FROM connections WHERE id=?`, conn).Scan(&am, &cred)
	if am != "PASSWORD" || cred != nil {
		t.Fatalf("host change kept the credential reference: %s %v", am, cred)
	}

	// Add a folder and a connection in it; both belong to the owner and are shared.
	var f map[string]interface{}
	if code := op.jsonDo("POST", base+"/folders", map[string]interface{}{"name": "demo-web", "parent_id": fid}, &f); code != 201 {
		t.Fatalf("folder: %d", code)
	}
	sub := int(f["id"].(float64))
	if code := op.jsonDo("POST", base+"/folders", map[string]interface{}{"name": "x", "parent_id": int(other64)}, nil); code != 400 {
		t.Fatalf("folder outside the share: %d", code)
	}
	if code := op.jsonDo("PUT", base+"/folders/"+itoa(sub), map[string]interface{}{"name": "demo-web-2"}, nil); code != 200 {
		t.Fatalf("folder rename: %d", code)
	}
	var nc map[string]interface{}
	if code := op.jsonDo("POST", base+"/connections", map[string]interface{}{"name": "demo-web", "protocol": "SSH", "host": "demo-web.example.com", "username": "ops", "password": "s3cret", "folder_id": sub}, &nc); code != 201 {
		t.Fatalf("create: %d %v", code, nc)
	}
	newID := int(nc["id"].(float64))
	var uid int
	db.QueryRow(`SELECT user_id FROM connections WHERE id=?`, newID).Scan(&uid)
	if uid != owner.userID || !shareHasConnection(share, newID) {
		t.Fatalf("new connection owner %d, shared %v", uid, shareHasConnection(share, newID))
	}
	if code := op.jsonDo("POST", base+"/connections", map[string]interface{}{"name": "x", "protocol": "SSH", "host": "x.example.com"}, nil); code != 400 {
		t.Fatalf("create without a shared folder: %d", code)
	}

	// Tunnels: listing works for the operator, not for the viewer; foreign keys are refused.
	var tl map[string]interface{}
	if code := op.jsonDo("GET", base+"/tunnels", nil, &tl); code != 200 {
		t.Fatalf("tunnels: %d", code)
	}
	if code := viewer.jsonDo("GET", base+"/tunnels", nil, nil); code != 403 {
		t.Fatalf("viewer tunnels: %d", code)
	}
	if code := op.jsonDo("POST", base+"/tunnels/t999999/start", nil, nil); code != 404 {
		t.Fatalf("unknown tunnel: %d", code)
	}

	// Perms reach the member's tree.
	p := sharePerms(share, RoleOperator, false)
	if p[PermManageConnections] {
		t.Fatal("stale share struct should have the switches off")
	}
	s2, _ := getShareByID(share.ID)
	if !sharePerms(s2, RoleOperator, false)[PermManageTunnels] || sharePerms(s2, RoleOperator, true)[PermManageTunnels] || sharePerms(s2, RoleViewer, false)[PermManageConnections] {
		t.Fatal("sharePerms wrong")
	}
}
