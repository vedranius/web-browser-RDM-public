package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ─── AUDIT LOG ───────────────────────────────────────
//
// Security-relevant events (sign-ins, account and policy changes, connections, shares,
// terminals, file changes, collaboration moderation) are written to the audit_log table
// and to the server log as a single "AUDIT" line, so they can also be shipped to a SIEM
// via journald/syslog. Old entries are removed after audit_retention_days.

func auditLog(r *http.Request, userID int, action, target string, details interface{}) {
	username := ""
	if userID > 0 {
		db.QueryRow(`SELECT username FROM users WHERE id=?`, userID).Scan(&username)
	}
	auditLogAs(r, userID, username, action, target, details)
}

func auditLogAs(r *http.Request, userID int, username, action, target string, details interface{}) {
	ip := ""
	if r != nil {
		ip = clientIP(r)
	}
	det := ""
	if details != nil {
		if s, ok := details.(string); ok {
			det = s
		} else if b, err := json.Marshal(details); err == nil && string(b) != "null" && string(b) != "{}" {
			det = string(b)
		}
	}
	if len(det) > 4000 {
		det = det[:4000]
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	var uid interface{}
	if userID > 0 {
		uid = userID
	}
	if _, err := db.Exec(`INSERT INTO audit_log (ts, user_id, username, ip, action, target, details) VALUES (?,?,?,?,?,?,?)`,
		ts, uid, username, ip, action, target, det); err != nil {
		log.Printf("audit: %v", err)
	}
	log.Printf("AUDIT action=%s user=%q ip=%s target=%q %s", action, username, ip, target, det)
}

func cleanupAuditLog() {
	days := settingInt("audit_retention_days")
	cutoff := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)
	db.Exec(`DELETE FROM audit_log WHERE ts < ?`, cutoff)
	db.Exec(`DELETE FROM collab_messages WHERE ts < ?`, cutoff)
}

type auditEntry struct {
	ID       int    `json:"id"`
	TS       string `json:"ts"`
	UserID   *int   `json:"user_id"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	Action   string `json:"action"`
	Target   string `json:"target"`
	Details  string `json:"details"`
}

func queryAudit(r *http.Request, limit int) ([]auditEntry, error) {
	q := r.URL.Query()
	where := []string{"1=1"}
	args := []interface{}{}
	if s := strings.TrimSpace(q.Get("q")); s != "" {
		where = append(where, "(action LIKE ? OR target LIKE ? OR details LIKE ? OR username LIKE ? OR ip LIKE ?)")
		like := "%" + s + "%"
		args = append(args, like, like, like, like, like)
	}
	if s := strings.TrimSpace(q.Get("action")); s != "" {
		where = append(where, "action LIKE ?")
		args = append(args, s+"%")
	}
	if s := strings.TrimSpace(q.Get("user")); s != "" {
		where = append(where, "username = ?")
		args = append(args, s)
	}
	if s := strings.TrimSpace(q.Get("from")); s != "" {
		where = append(where, "ts >= ?")
		args = append(args, s)
	}
	if s := strings.TrimSpace(q.Get("to")); s != "" {
		where = append(where, "ts <= ?")
		args = append(args, s)
	}
	if before, err := strconv.Atoi(q.Get("before_id")); err == nil && before > 0 {
		where = append(where, "id < ?")
		args = append(args, before)
	}
	args = append(args, limit)
	rows, err := db.Query(`SELECT id, ts, user_id, username, ip, action, target, details FROM audit_log WHERE `+
		strings.Join(where, " AND ")+` ORDER BY id DESC LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []auditEntry{}
	for rows.Next() {
		var e auditEntry
		var uid *int
		rows.Scan(&e.ID, &e.TS, &uid, &e.Username, &e.IP, &e.Action, &e.Target, &e.Details)
		e.UserID = uid
		out = append(out, e)
	}
	return out, nil
}

// GET /api/admin/audit?q&action&user&from&to&before_id&limit   (JSON)
// GET /api/admin/audit?format=csv                                (CSV download)
func apiAdminAuditHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	if r.Method != http.MethodGet {
		jsonError(w, "Method not allowed", 405)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		entries, err := queryAudit(r, 100000)
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		auditLog(r, adminID, "admin.audit_export", "", map[string]interface{}{"rows": len(entries)})
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", attachmentHeader("wrm-audit-"+time.Now().Format("20060102-150405")+".csv"))
		cw := csv.NewWriter(w)
		cw.Write([]string{"id", "time_utc", "user", "ip", "action", "target", "details"})
		for _, e := range entries {
			cw.Write([]string{strconv.Itoa(e.ID), e.TS, csvSafe(e.Username), e.IP, e.Action, csvSafe(e.Target), csvSafe(e.Details)})
		}
		cw.Flush()
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	entries, err := queryAudit(r, limit)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonOK(w, entries)
}

// csvSafe neutralises values that spreadsheet programs would interpret as formulas.
func csvSafe(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}
