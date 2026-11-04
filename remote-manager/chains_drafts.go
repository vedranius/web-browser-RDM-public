package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ─── CHAINS: AI DRAFTS (step 3) ───────────────────────────────────────────────────────────
//
// The AI assistant drafts a chain from the user's description, using the read-only tools it
// already has on the server it is attached to and the connections and storages this owner may
// use. A draft is a definition that is NOT monitored: it waits in the Chains window, where the
// user reviews it on the same live diagram. Only the user's approval creates the chain and
// starts its background checks.
//
// A draft carries only catalog check types (validateChain refuses anything else), can use only
// the owner's connections and storages, and never runs a command.

const chainDraftSchema = `CREATE TABLE IF NOT EXISTS chain_drafts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		session_id TEXT NOT NULL DEFAULT '',
		name TEXT NOT NULL DEFAULT '',
		definition TEXT NOT NULL DEFAULT '{}',
		created_at TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT '')`

type chainDraftRow struct {
	ID         int      `json:"id"`
	UserID     int      `json:"-"`
	SessionID  string   `json:"-"`
	Definition chainDef `json:"definition"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// chainDraftSummary is one entry of the draft list.
type chainDraftSummary struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Nodes     int    `json:"nodes"`
	Checks    int    `json:"checks"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func loadChainDraft(id int) (chainDraftRow, error) {
	var d chainDraftRow
	var def string
	err := db.QueryRow(`SELECT id, user_id, session_id, definition, created_at, updated_at FROM chain_drafts WHERE id=?`, id).
		Scan(&d.ID, &d.UserID, &d.SessionID, &def, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return d, err
	}
	if err := json.Unmarshal([]byte(def), &d.Definition); err != nil {
		return d, err
	}
	return d, nil
}

// decodeChainDefinition parses a definition strictly, so a typo in a field name is reported
// instead of silently dropped.
func decodeChainDefinition(raw json.RawMessage) (chainDef, error) {
	var d chainDef
	if len(bytes.TrimSpace(raw)) == 0 {
		return d, fmt.Errorf("definition is required")
	}
	if len(raw) > chainMaxDefBytes {
		return d, fmt.Errorf("definition is larger than %d KB", chainMaxDefBytes>>10)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&d); err != nil {
		return d, fmt.Errorf("definition: %v", truncateStr(err.Error(), 300))
	}
	return d, nil
}

// saveChainDraft stores the definition as the session's single draft: a new draft replaces the
// previous one from the same assistant session, so iterating never leaves a pile of drafts.
func saveChainDraft(userID int, sessionID string, def chainDef) (int, error) {
	now := nowRFC()
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	if sessionID != "" && userID > 0 {
		if _, err := tx.Exec(`DELETE FROM chain_drafts WHERE user_id=? AND session_id=?`, userID, sessionID); err != nil {
			return 0, err
		}
	}
	res, err := tx.Exec(`INSERT INTO chain_drafts (user_id, session_id, name, definition, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
		userID, sessionID, def.Name, string(jsonMarshal(def)), now, now)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func apiChainDrafts(w http.ResponseWriter, r *http.Request, userID int, rest string) {
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if parts[0] != "drafts" {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	if len(parts) == 1 {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		rows, err := db.Query(`SELECT id, definition, created_at, updated_at FROM chain_drafts WHERE user_id=? ORDER BY updated_at DESC`, userID)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		out := []chainDraftSummary{}
		for rows.Next() {
			var id int
			var def, created, updated string
			if rows.Scan(&id, &def, &created, &updated) != nil {
				continue
			}
			var d chainDef
			json.Unmarshal([]byte(def), &d)
			out = append(out, chainDraftSummary{ID: id, Name: d.Name, Nodes: len(d.Nodes), Checks: len(d.Checks), CreatedAt: created, UpdatedAt: updated})
		}
		jsonOK(w, out)
		return
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	d, err := loadChainDraft(id)
	if err != nil || d.UserID != userID {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			jsonOK(w, map[string]interface{}{"id": d.ID, "definition": d.Definition, "created_at": d.CreatedAt, "updated_at": d.UpdatedAt})
		case http.MethodDelete:
			if _, err := db.Exec(`DELETE FROM chain_drafts WHERE id=?`, d.ID); err != nil {
				jsonError(w, "Server error", http.StatusInternalServerError)
				return
			}
			auditLog(r, userID, "chain.draft_discarded", d.Definition.Name, map[string]interface{}{"draft": d.ID})
			jsonOK(w, map[string]bool{"ok": true})
		default:
			jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	if len(parts) != 3 || parts[2] != "approve" {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	apiChainApproveDraft(w, r, userID, d)
}

// apiChainApproveDraft is the only path from a draft to a monitored chain: the definition is
// validated again, the chain is created enabled, and the draft is removed in one transaction so
// two clicks cannot create the same chain twice.
func apiChainApproveDraft(w http.ResponseWriter, r *http.Request, userID int, d chainDraftRow) {
	def := d.Definition
	if errs := validateChain(&def, userID); len(errs) > 0 {
		chainErrors(w, errs)
		return
	}
	var n int
	db.QueryRow(`SELECT COUNT(1) FROM chains WHERE user_id=?`, userID).Scan(&n)
	if n >= 50 {
		jsonError(w, "At most 50 chains per user", http.StatusBadRequest)
		return
	}
	now := nowRFC()
	tx, err := db.Begin()
	if err != nil {
		jsonError(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	res, err := tx.Exec(`INSERT INTO chains (user_id, name, definition, enabled, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
		userID, def.Name, string(jsonMarshal(def)), 1, now, now)
	if err != nil {
		jsonError(w, "Server error", http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	del, err := tx.Exec(`DELETE FROM chain_drafts WHERE id=? AND user_id=?`, d.ID, userID)
	if err != nil {
		jsonError(w, "Server error", http.StatusInternalServerError)
		return
	}
	if gone, _ := del.RowsAffected(); gone != 1 {
		jsonError(w, "This draft was already approved or discarded", http.StatusConflict)
		return
	}
	if err := tx.Commit(); err != nil {
		jsonError(w, "Server error", http.StatusInternalServerError)
		return
	}
	auditLog(r, userID, "chain.created", def.Name, map[string]interface{}{"chain": id, "nodes": len(def.Nodes), "checks": len(def.Checks), "from_draft": d.ID})
	hub.sendTo(userID, []byte(fmt.Sprintf(`{"type":"chain_changed","id":%d}`, id)))
	hub.sendTo(userID, []byte(fmt.Sprintf(`{"type":"chain_draft_changed","id":%d,"removed":true}`, d.ID)))
	c, _ := loadChain(int(id))
	jsonOK(w, chainView(c))
}

// chainDraftViewForAI is the compact answer a drafting tool gets back.
func chainDraftViewForAI(id int) (map[string]interface{}, error) {
	d, err := loadChainDraft(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"draft_id": d.ID, "name": d.Definition.Name, "nodes": len(d.Definition.Nodes), "checks": len(d.Definition.Checks),
		"message": "Draft saved, not monitored. Open Chains to review it on the diagram; the chain starts only after the user approves it.",
	}, nil
}

// chainDraftStatus lists what a new chain may use, for the assistant that writes the draft.
func chainDraftStatus(userID int) map[string]interface{} {
	out := map[string]interface{}{}
	rows, err := db.Query(`SELECT id FROM chains WHERE user_id=? ORDER BY name COLLATE NOCASE`, userID)
	if err == nil {
		chains := []chainSummary{}
		for rows.Next() {
			var id int
			if rows.Scan(&id) == nil {
				if c, err := loadChain(id); err == nil {
					chains = append(chains, summarizeChain(c))
				}
			}
		}
		rows.Close()
		out["chains"] = chains
	}
	conns := []map[string]interface{}{}
	if rows, err := db.Query(`SELECT id, name, protocol FROM connections WHERE user_id=? ORDER BY name COLLATE NOCASE`, userID); err == nil {
		for rows.Next() {
			var id int
			var name, protocol string
			if rows.Scan(&id, &name, &protocol) == nil {
				conns = append(conns, map[string]interface{}{"conn_id": id, "name": name, "protocol": protocol})
			}
		}
		rows.Close()
	}
	out["connections"] = conns
	storages := []map[string]interface{}{}
	if rows, err := db.Query(`SELECT id, name FROM storages WHERE user_id=? ORDER BY name COLLATE NOCASE`, userID); err == nil {
		for rows.Next() {
			var id int
			var name string
			if rows.Scan(&id, &name) == nil {
				storages = append(storages, map[string]interface{}{"storage_id": id, "name": name})
			}
		}
		rows.Close()
	}
	out["storages"] = storages
	types := map[string]string{}
	for k, v := range chainCheckTypes {
		types[k] = v.Help
	}
	out["check_types"] = types
	out["link_kinds"] = []string{"data", "read", "write", "failover", "dependency"}
	out["rules"] = `{"all":[...]} / {"any":[...]} / {"check":"id","field":"value|ok","op":"==|!=|in|contains|exists","value":"..."} / {"node":"id","op":"==","value":"up"}`
	out["definition_example"] = map[string]interface{}{
		"name": "Example chain", "interval_s": 30, "confirm": 2,
		"nodes": []map[string]interface{}{
			{"id": "app", "kind": "server", "label": "app-01", "role": "app", "conn_id": 0},
			{"id": "site", "kind": "external", "label": "Public site", "role": "endpoint"},
		},
		"checks": []map[string]interface{}{
			{"id": "svc", "node": "app", "type": "systemd", "target": "nginx.service", "label": "Web server"},
			{"id": "cert", "node": "site", "type": "tls", "target": "example.com", "label": "Certificate"},
		},
		"links": []map[string]interface{}{{"from": "app", "to": "site", "kind": "data", "check": "svc"}},
		"scenarios": []map[string]interface{}{
			{"id": "normal", "label": "Normal", "when": map[string]interface{}{"all": []map[string]interface{}{{"check": "svc", "field": "ok", "op": "==", "value": "true"}}}},
			{"id": "degraded", "label": "Degraded", "when": map[string]interface{}{}},
		},
	}
	out["note"] = "A chain can only use the listed conn_id and storage_id values. The draft must use catalog check types and cannot run commands."
	return out
}

// ─── AI tools ────────────────────────────────────────────────────────────────────────────

func (s *aiSession) planChainTool(name string, in aiToolInput) (*aiPlanned, error) {
	if !settingBool("chains_enabled") {
		return nil, fmt.Errorf("Chains are turned off by the administrator")
	}
	host := in
	read := &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Host: &host}
	write := &aiPlanned{Act: aiAction{Tool: name, Kind: "write"}, Host: &host, NoEdit: true}
	switch name {
	case "chain_status":
		read.Display = "list chains, connections and check types"
		return read, nil
	case "chain_draft":
		def, err := decodeChainDefinition(in.Definition)
		if err != nil {
			return nil, err
		}
		if errs := validateChain(&def, s.UserID); len(errs) > 0 {
			return nil, fmt.Errorf("the chain is not valid yet: %s", strings.Join(errs, "; "))
		}
		pretty, _ := json.MarshalIndent(def, "", "  ")
		write.Path = "chain-draft:" + def.Name
		write.Diff = truncateStr(string(pretty), 12000)
		write.Display = fmt.Sprintf("draft chain %q (%d nodes, %d checks) — review in Chains", def.Name, len(def.Nodes), len(def.Checks))
		return write, nil
	}
	return nil, fmt.Errorf("unknown chain tool %q", name)
}

func (s *aiSession) runChainTool(ctx context.Context, name string, in aiToolInput) (string, int, error) {
	if !settingBool("chains_enabled") {
		return "", 2, fmt.Errorf("Chains are turned off by the administrator")
	}
	switch name {
	case "chain_status":
		b, _ := json.MarshalIndent(chainDraftStatus(s.UserID), "", "  ")
		return string(b), 0, nil
	case "chain_draft":
		def, err := decodeChainDefinition(in.Definition)
		if err != nil {
			return "", 2, err
		}
		if errs := validateChain(&def, s.UserID); len(errs) > 0 {
			return "", 2, fmt.Errorf("%s", strings.Join(errs, "; "))
		}
		id, err := saveChainDraft(s.UserID, s.UID, def)
		if err != nil {
			return "", 2, err
		}
		auditLog(nil, s.UserID, "chain.drafted", def.Name, map[string]interface{}{"draft": id, "nodes": len(def.Nodes), "checks": len(def.Checks), "session": s.ID})
		hub.sendTo(s.UserID, []byte(fmt.Sprintf(`{"type":"chain_draft_changed","id":%d}`, id)))
		view, err := chainDraftViewForAI(id)
		if err != nil {
			return "", 2, err
		}
		b, _ := json.MarshalIndent(view, "", "  ")
		return string(b), 0, nil
	}
	return "", 2, fmt.Errorf("unknown chain tool %q", name)
}
