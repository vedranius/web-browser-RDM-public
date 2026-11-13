package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

// ─── CONNECTION ORDER (per user) ─────────────────────
//
// Connections in the tree are listed by name until the user drags one to another place.
// The order is kept per user in connection_order (user_id, conn_id → position inside its
// folder); connections without a position follow the placed ones, by name. "Sort A–Z" on a
// folder removes the positions of its connections.

func initConnOrderTable() {
	ensureTable("connection_order", `CREATE TABLE IF NOT EXISTS connection_order (
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		conn_id INTEGER NOT NULL REFERENCES connections(id) ON DELETE CASCADE,
		position INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (user_id, conn_id))`,
		"user_id", "conn_id", "position")
}

// connOrderOf returns the user's saved positions (conn id → position, 1 = first).
func connOrderOf(userID int) map[int]int {
	out := map[int]int{}
	rows, err := db.Query(`SELECT conn_id, position FROM connection_order WHERE user_id=?`, userID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, pos int
		if rows.Scan(&id, &pos) == nil && pos > 0 {
			out[id] = pos
		}
	}
	return out
}

// sortConnViews puts placed connections first (by position), the rest keep their order (by name).
func sortConnViews(views []connView) {
	sort.SliceStable(views, func(i, j int) bool {
		a, b := views[i].SortOrder, views[j].SortOrder
		if a == 0 || b == 0 {
			return a != 0 && b == 0
		}
		return a < b
	})
}

// POST /api/connections/order {"ids": [3, 1, 2]}               → positions 1..n in this order
// POST /api/connections/order {"ids": [3, 1, 2], "reset": true} → back to A–Z for these
//
// Only the caller's own connections are accepted; the request is refused when any id is not.
func apiConnectionOrderHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in struct {
		IDs   []int `json:"ids"`
		Reset bool  `json:"reset"`
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil || len(in.IDs) == 0 || len(in.IDs) > 5000 {
		jsonError(w, "Invalid connection order", 400)
		return
	}
	seen := map[int]bool{}
	for _, id := range in.IDs {
		if id <= 0 || seen[id] || !userOwnsConnection(id, uid) {
			jsonError(w, "Not found", 404)
			return
		}
		seen[id] = true
	}
	tx, err := db.Begin()
	if err != nil {
		jsonError(w, "Cannot change connection order", 500)
		return
	}
	defer tx.Rollback()
	for i, id := range in.IDs {
		if in.Reset {
			_, err = tx.Exec(`DELETE FROM connection_order WHERE user_id=? AND conn_id=?`, uid, id)
		} else {
			_, err = tx.Exec(`INSERT INTO connection_order (user_id, conn_id, position) VALUES (?,?,?)
				ON CONFLICT(user_id, conn_id) DO UPDATE SET position=excluded.position`, uid, id, i+1)
		}
		if err != nil {
			jsonError(w, "Cannot change connection order", 500)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		jsonError(w, "Cannot change connection order", 500)
		return
	}
	action := "connection.reordered"
	if in.Reset {
		action = "connection.order_reset"
	}
	auditLog(r, uid, action, "", map[string]interface{}{"ids": in.IDs})
	jsonOK(w, map[string]interface{}{"ok": true, "order": connOrderOf(uid)})
}
