package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// ─── TEAM-MANAGED SHARES (#192) ──────────────────────
//
// A share owner can let the team manage the shared connections with two switches:
//   allow_edit    – operators and above may add and edit connections and folders
//   allow_tunnels – operators and above may start and stop the connections' tunnels
// Both are off by default and never apply to guests (people without an account).
//
// The connections and folders stay the owner's: new ones are created in the owner's
// account inside a shared folder, so they appear in the share automatically. A member
// can never pick or change the owner's stored secrets (key store keys, vault
// credentials, Bitwarden items, key files on the server); when a member changes the
// host or protocol of a connection, its saved password and key are cleared, so the
// owner's secret is never sent to a server the member chose.
//
// /api/share-team/{share}/connections          POST  add a connection to a shared folder
// /api/share-team/{share}/connections/{id}     GET · PUT
// /api/share-team/{share}/folders              POST  add a folder under a shared folder
// /api/share-team/{share}/folders/{id}         PUT   rename
// /api/share-team/{share}/tunnels              GET   tunnels of the shared connections
// /api/share-team/{share}/tunnels/{key}/start|stop
// {share} is the share's token or id (as /api/connections/… ?share_id= uses it).

type shareTeamConn struct {
	Name     *string   `json:"name"`
	Protocol *string   `json:"protocol"`
	Host     *string   `json:"host"`
	Username *string   `json:"username"`
	Password *string   `json:"password"`
	FolderID *int      `json:"folder_id"`
	WebPath  *string   `json:"web_path"`
	Tags     *[]string `json:"tags"`
	Notes    *string   `json:"notes"`
}

func shareTeamAccess(w http.ResponseWriter, r *http.Request, ident, perm string) (*ShareAccess, bool) {
	share, ok := shareFromRequest(r, ident)
	if !ok {
		jsonError(w, "Share not found", 404)
		return nil, false
	}
	sa, d := resolveShareAccess(r, share)
	if d != nil {
		jsonError(w, d.Message, d.Code)
		return nil, false
	}
	if sa.Guest || sa.UserID <= 0 || isAuditorUser(sa.UserID) || !sharePerms(sa.Share, sa.Role, sa.Guest)[perm] {
		jsonError(w, "This share does not allow you to do this", 403)
		return nil, false
	}
	return sa, true
}

// shareHasFolder reports whether the folder is a shared folder or below one.
func shareHasFolder(share ShareLink, folderID int) bool {
	var n int
	db.QueryRow(`SELECT COUNT(1) FROM folders WHERE id=? AND user_id=? AND id IN (`+shareFolderSubtreeSQL+`)`,
		folderID, share.OwnerID, share.ID, share.OwnerID, share.OwnerID).Scan(&n)
	return n > 0
}

func shareTeamConnView(c Connection) map[string]interface{} {
	v := c.view()
	return map[string]interface{}{"id": c.ID, "name": c.Name, "protocol": c.Protocol, "host": c.Host, "username": c.Username,
		"auth_method": c.AuthMethod, "folder_id": c.FolderID, "web_path": c.WebPath, "has_password": c.Password != "",
		"has_private_key": c.PrivateKey != "", "tags": v.Tags, "notes": loadNotes(c.ID), "tunnels": v.Tunnels}
}

func apiShareTeamHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAuth(w, r); !ok {
		return
	}
	parts := strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/share-team/"), "/"), "/")
	if len(parts) < 2 || parts[0] == "" {
		jsonError(w, "Not found", 404)
		return
	}
	switch parts[1] {
	case "connections":
		sa, ok := shareTeamAccess(w, r, parts[0], PermManageConnections)
		if !ok {
			return
		}
		if len(parts) == 2 && r.Method == http.MethodPost {
			shareTeamCreateConn(w, r, sa)
			return
		}
		if len(parts) != 3 {
			jsonError(w, "Not found", 404)
			return
		}
		id, _ := strconv.Atoi(parts[2])
		if id <= 0 || !shareHasConnection(sa.Share, id) {
			jsonError(w, "Connection not found", 404)
			return
		}
		cur, err := loadConnectionRaw(id)
		if err != nil {
			jsonError(w, "Connection not found", 404)
			return
		}
		switch r.Method {
		case http.MethodGet:
			jsonOK(w, shareTeamConnView(cur))
		case http.MethodPut:
			shareTeamUpdateConn(w, r, sa, cur)
		default:
			jsonError(w, "Method not allowed", 405)
		}
	case "folders":
		sa, ok := shareTeamAccess(w, r, parts[0], PermManageConnections)
		if !ok {
			return
		}
		shareTeamFolders(w, r, sa, parts[2:])
	case "tunnels":
		sa, ok := shareTeamAccess(w, r, parts[0], PermManageTunnels)
		if !ok {
			return
		}
		shareTeamTunnels(w, r, sa, parts[2:])
	default:
		jsonError(w, "Not found", 404)
	}
}

func decodeShareTeam(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(v); err != nil {
		jsonError(w, "Bad JSON", 400)
		return false
	}
	return true
}

// applyShareTeamInput copies the fields a member may set; everything else stays as it is.
func applyShareTeamInput(c *Connection, in shareTeamConn) {
	if in.Name != nil {
		c.Name = *in.Name
	}
	if in.Protocol != nil {
		c.Protocol = *in.Protocol
	}
	if in.Host != nil {
		c.Host = *in.Host
	}
	if in.Username != nil {
		c.Username = *in.Username
	}
	if in.WebPath != nil {
		c.WebPath = *in.WebPath
	}
	if in.Tags != nil {
		c.Tags = *in.Tags
	}
}

func shareTeamCreateConn(w http.ResponseWriter, r *http.Request, sa *ShareAccess) {
	var in shareTeamConn
	if !decodeShareTeam(w, r, &in) {
		return
	}
	if in.FolderID == nil || !shareHasFolder(sa.Share, *in.FolderID) {
		jsonError(w, "Choose a shared folder for the new connection", 400)
		return
	}
	c := Connection{AuthMethod: "PASSWORD", FolderID: in.FolderID}
	applyShareTeamInput(&c, in)
	if in.Password != nil {
		c.Password = *in.Password
	}
	if err := normalizeConnection(&c); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	if c.Protocol == "SERIAL" {
		jsonError(w, "Serial ports cannot be added through a share", 403)
		return
	}
	c.AuthMethod, c.KeyID, c.CredentialID, c.KeyPath, c.PrivateKey, c.JumpID, c.ProxyID = "PASSWORD", nil, nil, "", "", nil, nil
	if in.Notes != nil {
		if err := validNotes(*in.Notes); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
	}
	owner := sa.Share.OwnerID
	plain := c
	plain.UserID = owner
	encryptConnectionSecrets(&c)
	res, err := db.Exec(`INSERT INTO connections (name,protocol,host,username,auth_method,password,private_key,key_path,folder_id,user_id,web_path,tags) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		c.Name, c.Protocol, c.Host, c.Username, c.AuthMethod, c.Password, "", "", c.FolderID, owner, c.WebPath, tagsString(c.Tags))
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	id64, _ := res.LastInsertId()
	id := int(id64)
	if in.Notes != nil {
		db.Exec(`UPDATE connections SET notes=? WHERE id=?`, *in.Notes, id)
	}
	statusMon.poke()
	auditLogRef(r, sa.UserID, sa.Username, "connection.created", c.Name, map[string]interface{}{"host": c.Host, "protocol": c.Protocol,
		"share_id": sa.Share.ID, "owner": usernameOf(owner)}, auditRef{ConnID: id})
	saved, err := loadConnectionRaw(id)
	if err != nil {
		saved = plain
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, shareTeamConnView(saved))
}

func shareTeamUpdateConn(w http.ResponseWriter, r *http.Request, sa *ShareAccess, cur Connection) {
	var in shareTeamConn
	if !decodeShareTeam(w, r, &in) {
		return
	}
	c := cur
	applyShareTeamInput(&c, in)
	if in.FolderID != nil && *in.FolderID != intOr(cur.FolderID) {
		if !shareHasFolder(sa.Share, *in.FolderID) {
			jsonError(w, "The connection can only be moved to a shared folder", 400)
			return
		}
		c.FolderID = in.FolderID
	}
	if err := normalizeConnection(&c); err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	if c.Protocol == "SERIAL" && cur.Protocol != "SERIAL" {
		jsonError(w, "Serial ports cannot be added through a share", 403)
		return
	}
	// The member cannot choose the owner's stored secrets: these stay as they are.
	c.AuthMethod, c.KeyID, c.CredentialID, c.KeyPath = cur.AuthMethod, cur.KeyID, cur.CredentialID, cur.KeyPath
	changed := []string{}
	if !strings.EqualFold(c.Host, cur.Host) || c.Protocol != cur.Protocol {
		// A new target never receives the owner's saved password or key.
		c.Password, c.PrivateKey = "", ""
		if c.AuthMethod != "PASSWORD" && c.AuthMethod != "KEY" {
			c.AuthMethod, c.KeyID, c.CredentialID, c.KeyPath = "PASSWORD", nil, nil, ""
		}
		if cur.Password != "" || cur.PrivateKey != "" || cur.AuthMethod != c.AuthMethod {
			changed = append(changed, "secrets_cleared")
		}
	}
	if in.Password != nil && *in.Password != "" {
		c.Password = *in.Password
		if c.AuthMethod != "PASSWORD" && c.AuthMethod != "KEY" {
			c.AuthMethod, c.KeyID, c.CredentialID, c.KeyPath = "PASSWORD", nil, nil, ""
		}
		changed = append(changed, "password")
	}
	for f, same := range map[string]bool{"name": c.Name == cur.Name, "host": c.Host == cur.Host, "username": c.Username == cur.Username,
		"protocol": c.Protocol == cur.Protocol, "web_path": c.WebPath == cur.WebPath, "folder": intPtrEq(c.FolderID, cur.FolderID)} {
		if !same {
			changed = append(changed, f)
		}
	}
	if in.Notes != nil {
		if err := validNotes(*in.Notes); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
	}
	encryptConnectionSecrets(&c)
	if _, err := db.Exec(`UPDATE connections SET name=?,protocol=?,host=?,username=?,auth_method=?,password=?,private_key=?,key_path=?,key_id=?,credential_id=?,folder_id=?,web_path=? WHERE id=? AND user_id=?`,
		c.Name, c.Protocol, c.Host, c.Username, c.AuthMethod, c.Password, c.PrivateKey, c.KeyPath, c.KeyID, c.CredentialID, c.FolderID, c.WebPath, cur.ID, sa.Share.OwnerID); err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	if in.Tags != nil {
		db.Exec(`UPDATE connections SET tags=? WHERE id=?`, tagsString(c.Tags), cur.ID)
	}
	if in.Notes != nil && *in.Notes != loadNotes(cur.ID) {
		db.Exec(`UPDATE connections SET notes=? WHERE id=?`, *in.Notes, cur.ID)
		changed = append(changed, "notes")
	}
	statusMon.poke()
	auditLogRef(r, sa.UserID, sa.Username, "connection.updated", c.Name, map[string]interface{}{"id": cur.ID, "changed": changed,
		"share_id": sa.Share.ID, "owner": usernameOf(sa.Share.OwnerID)}, auditRef{ConnID: cur.ID})
	saved, _ := loadConnectionRaw(cur.ID)
	jsonOK(w, shareTeamConnView(saved))
}

func intOr(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func shareTeamFolders(w http.ResponseWriter, r *http.Request, sa *ShareAccess, rest []string) {
	var in struct {
		Name     string `json:"name"`
		ParentID int    `json:"parent_id"`
	}
	if !decodeShareTeam(w, r, &in) {
		return
	}
	in.Name = truncateStr(strings.TrimSpace(in.Name), 120)
	if in.Name == "" {
		jsonError(w, "Name is required", 400)
		return
	}
	owner := sa.Share.OwnerID
	switch {
	case len(rest) == 0 && r.Method == http.MethodPost:
		if !shareHasFolder(sa.Share, in.ParentID) {
			jsonError(w, "Choose a shared folder as the parent", 400)
			return
		}
		var sort int
		db.QueryRow(`SELECT COALESCE(MAX(sort_order),0)+1 FROM folders WHERE user_id=? AND parent_id=?`, owner, in.ParentID).Scan(&sort)
		res, err := db.Exec(`INSERT INTO folders (name,user_id,parent_id,sort_order) VALUES (?,?,?,?)`, in.Name, owner, in.ParentID, sort)
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		id, _ := res.LastInsertId()
		auditLog(r, sa.UserID, "folder.created", in.Name, map[string]interface{}{"id": id, "share_id": sa.Share.ID, "owner": usernameOf(owner)})
		w.WriteHeader(http.StatusCreated)
		jsonOK(w, map[string]interface{}{"id": id, "name": in.Name, "parent_id": in.ParentID})
	case len(rest) == 1 && r.Method == http.MethodPut:
		id, _ := strconv.Atoi(rest[0])
		if id <= 0 || !shareHasFolder(sa.Share, id) {
			jsonError(w, "Folder not found", 404)
			return
		}
		var old string
		db.QueryRow(`SELECT name FROM folders WHERE id=?`, id).Scan(&old)
		db.Exec(`UPDATE folders SET name=? WHERE id=? AND user_id=?`, in.Name, id, owner)
		auditLog(r, sa.UserID, "folder.updated", in.Name, map[string]interface{}{"id": id, "old_name": old, "share_id": sa.Share.ID, "owner": usernameOf(owner)})
		jsonOK(w, map[string]interface{}{"id": id, "name": in.Name})
	default:
		jsonError(w, "Not found", 404)
	}
}

// shareTeamTunnels lists, starts and stops the tunnels the owner configured on the
// shared connections. The owner's tunnel policy applies (the tunnel runs as the owner's),
// and the member's own account must be allowed to use tunnels too.
func shareTeamTunnels(w http.ResponseWriter, r *http.Request, sa *ShareAccess, rest []string) {
	owner := sa.Share.OwnerID
	shared := map[int]SharedConn{}
	for _, c := range sharedConnectionsFor(sa.Share) {
		shared[c.ID] = c
	}
	if len(rest) == 0 {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		out := []tunnelView{}
		for _, d := range loadTunnelDefs("user_id=?", owner) {
			c, ok := shared[d.ConnID]
			if !ok {
				continue
			}
			key := "t" + strconv.Itoa(d.ID)
			if t := tunnelMgr.get(key); t != nil {
				out = append(out, t.view())
				continue
			}
			out = append(out, tunnelView{Key: key, tunnelDef: d, ConnName: c.Name, ConnHost: c.Host, State: "stopped"})
		}
		allowed := tunnelsAllowedFor(sa.UserID) == nil && tunnelsAllowedFor(owner) == nil
		jsonOK(w, map[string]interface{}{"tunnels": out, "allowed": allowed})
		return
	}
	if len(rest) != 2 || r.Method != http.MethodPost || !strings.HasPrefix(rest[0], "t") {
		jsonError(w, "Not found", 404)
		return
	}
	id, err := strconv.Atoi(strings.TrimPrefix(rest[0], "t"))
	if err != nil {
		jsonError(w, "Not found", 404)
		return
	}
	def, err := loadTunnelDef(id)
	if err != nil || def.UserID != owner {
		jsonError(w, "Not found", 404)
		return
	}
	if _, ok := shared[def.ConnID]; !ok {
		jsonError(w, "Not found", 404)
		return
	}
	key := rest[0]
	switch rest[1] {
	case "stop":
		if tunnelMgr.get(key) != nil {
			tunnelMgr.stop(key, sa.UserID, "stopped by "+sa.Username+" (share)")
		}
		jsonOK(w, map[string]bool{"ok": true})
	case "start":
		for _, uid := range []int{sa.UserID, owner} {
			if err := tunnelsAllowedFor(uid); err != nil {
				jsonError(w, err.Error(), 403)
				return
			}
		}
		if err := validateTunnelDef(&def, owner); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		c, err := loadConnection(def.ConnID)
		if err != nil || c.UserID != owner {
			jsonError(w, "Not found", 404)
			return
		}
		t, err := tunnelMgr.start(key, def, c, sa.UserID, "manual", r)
		if err != nil {
			jsonError(w, err.Error(), 502)
			return
		}
		tunnelMgr.notify(sa.UserID)
		jsonOK(w, t.view())
	default:
		jsonError(w, "Not found", 404)
	}
}
