package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ─── PASSWORD RESET ON THE SERVER (wrm -reset-password) ───
//
// Whoever can run WRM's binary against its database can reset a password; that cannot be
// prevented entirely (the same person can read the database and the key). What WRM decides is
// how loud and how slow it is, with the setting cli_password_reset (raised by the security level):
//
//   alert    (Home)       the reset happens at once; every administrator gets a security alert.
//   stopped  (Business)   as alert, but only while the WRM server is not running, so a silent
//                         takeover of a running system is impossible (the outage is visible).
//   confirm  (Enterprise) the command files a request. It takes effect when another administrator
//                         approves it in Admin → Security, when the user's own 2FA code is given
//                         (-reset-code), or when the command is run again after cliResetDelay.
//
// In every mode an AI working through WRM cannot run the command (ai_destructive.go, rule
// wrm_recovery) and the temporary password is redacted from AI output (ai_redact.go).
//
// The running server writes a heartbeat (server_heartbeat) every 20 s; the command treats
// WRM as running while the heartbeat is younger than a minute.

const cliResetDelay = time.Hour

func initCLIResetSchema() {
	ensureTable("server_heartbeat", `CREATE TABLE IF NOT EXISTS server_heartbeat (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		pid INTEGER NOT NULL DEFAULT 0,
		host TEXT NOT NULL DEFAULT '',
		at TEXT NOT NULL DEFAULT ''
	)`, "pid", "host", "at")
	ensureTable("cli_reset_requests", `CREATE TABLE IF NOT EXISTS cli_reset_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		reset_2fa INTEGER NOT NULL DEFAULT 0,
		host TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		decided_by TEXT NOT NULL DEFAULT '',
		decided_at TEXT NOT NULL DEFAULT ''
	)`, "username", "reset_2fa", "host", "created_at", "status", "decided_by", "decided_at")
}

// runServerHeartbeat marks this server as running until stop is closed.
func runServerHeartbeat(stop <-chan struct{}) {
	host, _ := os.Hostname()
	beat := func() {
		db.Exec(`INSERT INTO server_heartbeat (id, pid, host, at) VALUES (1,?,?,?)
			ON CONFLICT(id) DO UPDATE SET pid=excluded.pid, host=excluded.host, at=excluded.at`, os.Getpid(), host, time.Now().UTC().Format(time.RFC3339))
	}
	beat()
	t := time.NewTicker(20 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			beat()
		case <-stop:
			db.Exec(`DELETE FROM server_heartbeat WHERE id=1 AND pid=?`, os.Getpid())
			return
		}
	}
}

// serverRunning reports whether a WRM server uses this database (heartbeat younger than a minute).
func serverRunning() (bool, string) {
	var at, host string
	var pid int
	if db.QueryRow(`SELECT at, host, pid FROM server_heartbeat WHERE id=1`).Scan(&at, &host, &pid) != nil {
		return false, ""
	}
	ts, err := time.Parse(time.RFC3339, at)
	if err != nil || time.Since(ts) > time.Minute {
		return false, ""
	}
	return true, fmt.Sprintf("%s, process %d", host, pid)
}

// cliResetMode is the stricter of the stored setting and a WRM_CLI_PASSWORD_RESET forced in
// the environment. The command runs in the environment of whoever starts it, so a variable
// set there may make the reset stricter but never looser than what the server stores.
func cliResetMode() string {
	order := map[string]int{"alert": 0, "stopped": 1, "confirm": 2}
	mode := "alert"
	for _, v := range []string{cliStoredSetting("cli_password_reset"), getSetting("cli_password_reset")} {
		if r, ok := order[v]; ok && r > order[mode] {
			mode = v
		}
	}
	return mode
}

// cliResetAlertOn: the alert is on when the stored setting or the environment turns it on
// (again, the caller's environment cannot turn it off).
func cliResetAlertOn() bool {
	return cliStoredSetting("alert_cli_reset") == "1" || getSetting("alert_cli_reset") == "1"
}

// cliStoredSetting is the value in the database (or the default), ignoring the environment.
func cliStoredSetting(key string) string {
	settingsStore.RLock()
	v, ok := settingsStore.values[key]
	settingsStore.RUnlock()
	if ok {
		return v
	}
	s, _ := settingSpecFor(key)
	return s.Default
}

// resetPasswordCLI implements "wrm -reset-password <username>" for locked-out users.
func resetPasswordCLI(username string, reset2FA bool) error {
	return resetPasswordCLIWith(username, reset2FA, "")
}

func resetPasswordCLIWith(username string, reset2FA bool, code string) error {
	u, err := loadUser("username=?", username)
	if err == sql.ErrNoRows || u == nil {
		return fmt.Errorf("user %q not found", username)
	}
	if err != nil {
		return err
	}
	host, _ := os.Hostname()
	switch cliResetMode() {
	case "stopped":
		if on, where := serverRunning(); on {
			cliResetAlertText(u.Username, "refused",
				fmt.Sprintf("Someone tried to reset the password of %s on the server %s with wrm -reset-password while WRM was running; it was refused (it works only while WRM is stopped).", u.Username, host),
				fmt.Sprintf("Netko je pokušao resetirati lozinku korisnika %s na poslužitelju %s naredbom wrm -reset-password dok je WRM radio; odbijeno je (radi samo dok je WRM zaustavljen).", u.Username, host))
			return fmt.Errorf("WRM is running (%s). At this security level a password can be reset on the server only while WRM is stopped: stop the WRM service, run this command again, then start WRM. The administrators were told about this attempt", where)
		}
	case "confirm":
		ok, msg, err := cliResetConfirmed(u, reset2FA, code, host)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println(msg)
			return nil
		}
	}
	return cliDoReset(u, reset2FA)
}

// cliResetConfirmed decides a reset in the confirm mode. It returns true when the reset may
// happen now; otherwise the message says what is needed.
func cliResetConfirmed(u *userRow, reset2FA bool, code, host string) (bool, string, error) {
	now := time.Now().UTC()
	// the user's own 2FA code proves the user asked for it
	if code != "" {
		if !u.TOTPEnabled {
			return false, "", fmt.Errorf("%s has no two-factor authentication, so -reset-code cannot confirm the reset", u.Username)
		}
		if !cliResetCodeOK(u, code) {
			cliResetAlertText(u.Username, "bad_code",
				fmt.Sprintf("A password reset of %s on the server %s was refused: the 2FA code given with wrm -reset-password was wrong.", u.Username, host),
				fmt.Sprintf("Reset lozinke korisnika %s na poslužitelju %s je odbijen: 2FA kod uz wrm -reset-password nije točan.", u.Username, host))
			return false, "", fmt.Errorf("the 2FA code is not valid; the administrators were told")
		}
		db.Exec(`UPDATE cli_reset_requests SET status='done', decided_by=?, decided_at=? WHERE username=? AND status IN ('pending','approved')`, "2fa:"+u.Username, now.Format(time.RFC3339), u.Username)
		return true, "", nil
	}
	var id int64
	var created, status string
	var r2 int
	err := db.QueryRow(`SELECT id, created_at, status, reset_2fa FROM cli_reset_requests WHERE username=? AND status IN ('pending','approved') ORDER BY id DESC LIMIT 1`, u.Username).Scan(&id, &created, &status, &r2)
	if err == nil {
		ts, _ := time.Parse(time.RFC3339, created)
		if status == "approved" || now.Sub(ts) >= cliResetDelay {
			if reset2FA && r2 == 0 {
				return false, "", fmt.Errorf("request %d did not ask to turn off 2FA; run the command without -reset-2fa, or reject that request in Admin → Security and file a new one", id)
			}
			db.Exec(`UPDATE cli_reset_requests SET status='done', decided_at=CASE WHEN decided_at='' THEN ? ELSE decided_at END WHERE id=?`, now.Format(time.RFC3339), id)
			return true, "", nil
		}
		left := cliResetDelay - now.Sub(ts)
		return false, fmt.Sprintf("Reset request %d for %s is waiting. It takes effect when another administrator approves it in Admin → Security, when you run this command again with -reset-code and %s's current 2FA code, or when you run it again after %s (at %s).",
			id, u.Username, u.Username, left.Round(time.Minute), ts.Add(cliResetDelay).Local().Format("15:04")), nil
	}
	res, err := db.Exec(`INSERT INTO cli_reset_requests (username, reset_2fa, host, created_at) VALUES (?,?,?,?)`, u.Username, boolInt(reset2FA), host, now.Format(time.RFC3339))
	if err != nil {
		return false, "", err
	}
	id, _ = res.LastInsertId()
	auditLogAs(nil, 0, "cli", "admin.cli_reset_requested", u.Username, map[string]interface{}{"request": id, "reset_2fa": reset2FA, "host": host})
	cliResetAlertText(u.Username, "requested",
		fmt.Sprintf("A password reset of %s was requested on the server %s with wrm -reset-password (request %d). Approve or reject it in Admin → Security; unless it is rejected, it can take effect after one hour.", u.Username, host, id),
		fmt.Sprintf("Na poslužitelju %s naredbom wrm -reset-password zatražen je reset lozinke korisnika %s (zahtjev %d). Odobri ga ili odbij u Admin → Sigurnost; ako se ne odbije, može se provesti nakon sat vremena.", host, u.Username, id))
	return false, fmt.Sprintf("At this security level a reset on the server needs a confirmation. Request %d for %s was filed and every administrator was told.\nIt takes effect when another administrator approves it in Admin → Security, when you run this command again with -reset-code and %s's current 2FA code, or when you run it again after one hour (at %s) unless it was rejected.",
		id, u.Username, u.Username, now.Add(cliResetDelay).Local().Format("15:04")), nil
}

// cliDoReset sets a new temporary password and tells the administrators.
func cliDoReset(u *userRow, reset2FA bool) error {
	pw := randomPassword()
	hash, err := hashPassword(pw)
	if err != nil {
		return err
	}
	db.Exec(`UPDATE users SET password_hash=?, must_change_password=1, failed_logins=0, locked_until='', disabled=0 WHERE id=?`, hash, u.ID)
	if reset2FA {
		db.Exec(`UPDATE users SET totp_secret='', totp_enabled=0, totp_pending='', recovery_codes='' WHERE id=?`, u.ID)
	}
	revokeUserSessions(u.ID, 0)
	auditLogAs(nil, 0, "cli", "admin.user_password_reset", u.Username, map[string]interface{}{"cli": true, "reset_2fa": reset2FA, "mode": cliResetMode()})
	alerted := cliResetAlert(u.Username, reset2FA)
	// tell the user too; WRM may be stopped, so try to deliver at once
	if id := mailSecurityNotice(u.ID, "password_reset", "Your WRM password was reset on the server with wrm -reset-password", "Tvoja WRM lozinka resetirana je na poslužitelju naredbom wrm -reset-password", ""); id > 0 {
		mailDeliverDue(id)
	}
	fmt.Printf("New temporary password for %s: %s\n(The user must change it after signing in.)\n", u.Username, pw)
	if alerted {
		fmt.Println("Every administrator was told about this reset (security alert).")
	}
	return nil
}

// cliResetAlertText raises a cli_reset alert with the given texts (setting alert_cli_reset).
func cliResetAlertText(username, kind, en, hr string) bool {
	if !cliResetAlertOn() {
		return false
	}
	in := alertInput{Action: "admin.user_password_reset", Username: "cli", Target: username, Details: map[string]interface{}{"cli": true, "kind": kind}}
	raiseAlert("cli_reset", in, en, hr, true)
	return true
}

// cliResetAlert tells every administrator that a password was reset on the server. It runs in
// the command's own process, so it raises the alert directly instead of through the audit
// queue: the inbox entries and the pending notifications are rows in the database, which the
// running server delivers.
func cliResetAlert(username string, reset2FA bool) bool {
	host, _ := os.Hostname()
	what, whatHR := "", ""
	if reset2FA {
		what, whatHR = " and its two-factor authentication was turned off", " i isključena joj je dvofaktorska autentifikacija"
	}
	return cliResetAlertText(username, "done",
		fmt.Sprintf("The password of %s was reset on the server %s with wrm -reset-password%s. If nobody you know did this, someone has access to the WRM server.", username, host, what),
		fmt.Sprintf("Lozinka korisnika %s resetirana je na poslužitelju %s naredbom wrm -reset-password%s. Ako to nije napravio netko poznat, netko ima pristup WRM poslužitelju.", username, host, whatHR))
}

// ─── Admin → Security: reset requests ───
//
//   GET  /api/admin/cli-resets              pending and recent requests
//   POST /api/admin/cli-resets/{id}/approve  another administrator (step-up), never the user it resets
//   POST /api/admin/cli-resets/{id}/reject

func apiAdminCLIResetsHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/cli-resets"), "/")
	if rest == "" {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		out := []map[string]interface{}{}
		rows, err := db.Query(`SELECT id, username, reset_2fa, host, created_at, status, decided_by, decided_at FROM cli_reset_requests ORDER BY id DESC LIMIT 50`)
		if err == nil {
			for rows.Next() {
				var id int64
				var r2 int
				var user, host, created, status, by, at string
				if rows.Scan(&id, &user, &r2, &host, &created, &status, &by, &at) == nil {
					ts, _ := time.Parse(time.RFC3339, created)
					out = append(out, map[string]interface{}{"id": id, "username": user, "reset_2fa": r2 == 1, "host": host, "created_at": created,
						"status": status, "decided_by": by, "decided_at": at, "effective_at": ts.Add(cliResetDelay).UTC().Format(time.RFC3339)})
				}
			}
			rows.Close()
		}
		jsonOK(w, map[string]interface{}{"mode": cliResetMode(), "requests": out})
		return
	}
	parts := strings.Split(rest, "/")
	if len(parts) != 2 || r.Method != http.MethodPost || (parts[1] != "approve" && parts[1] != "reject") {
		jsonError(w, "Not found", 404)
		return
	}
	id, _ := strconv.ParseInt(parts[0], 10, 64)
	var user, status string
	if db.QueryRow(`SELECT username, status FROM cli_reset_requests WHERE id=?`, id).Scan(&user, &status) != nil {
		jsonError(w, "Request not found", 404)
		return
	}
	if status != "pending" {
		jsonError(w, "This request was already decided ("+status+")", 409)
		return
	}
	me := usernameOf(adminID)
	if strings.EqualFold(me, user) {
		jsonError(w, "Another administrator has to decide a reset of your own password", 403)
		return
	}
	// approving hands a password to whoever sits at the server: always a fresh confirmation
	if parts[1] == "approve" && !requireStepUpAlways(w, r, "cli_reset_decide") {
		return
	}
	if parts[1] == "reject" && !requireStepUp(w, r, "cli_reset_decide") {
		return
	}
	newStatus := map[string]string{"approve": "approved", "reject": "rejected"}[parts[1]]
	res, _ := db.Exec(`UPDATE cli_reset_requests SET status=?, decided_by=?, decided_at=? WHERE id=? AND status='pending'`, newStatus, me, time.Now().UTC().Format(time.RFC3339), id)
	if n, _ := res.RowsAffected(); n != 1 {
		jsonError(w, "This request was already decided", 409)
		return
	}
	auditLog(r, adminID, "admin.cli_reset_"+newStatus, user, map[string]interface{}{"request": id})
	jsonOK(w, map[string]interface{}{"ok": true, "status": newStatus})
}

// cliResetCodeOK checks the user's own second factor (a TOTP code or a recovery code) and
// consumes it, so the same code cannot confirm a second reset.
func cliResetCodeOK(u *userRow, code string) bool {
	method, _, err := consumeUserSecondFactor(u, code)
	return err == nil && method != ""
}
