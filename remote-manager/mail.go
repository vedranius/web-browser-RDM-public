package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─── WRM'S OWN MAIL (v14.6.1) ───
//
// One SMTP account (Admin → Mail) that WRM uses to write to its users: e-mail address
// verification and security notices (password changed or reset, two-factor authentication
// turned on or off, e-mail address changed). Presets fill in Gmail and Outlook /
// Microsoft 365; both need an app password (an account password with 2-step verification
// does not work over SMTP).
//
// Mail is written to the table mail_outbox and delivered by mailLoop in the running server
// with retries, so a slow or unreachable SMTP server never delays a request. The command
// line (wrm -reset-password) also tries to deliver what it queued at once, because WRM may
// be stopped while it runs.
//
// The SMTP password is encrypted at rest (encryptValue) and never sent to the browser.
// Links in mail use the configured WRM address (public_url), never the request's Host header.

type mailSettings struct {
	Enabled     bool   `json:"enabled"`
	Preset      string `json:"preset"` // gmail, outlook, office365, custom
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Security    string `json:"security"` // starttls, tls, none
	Username    string `json:"username"`
	PasswordSet bool   `json:"password_set"`
	From        string `json:"from"`
	FromName    string `json:"from_name"`
	PublicURL   string `json:"public_url"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
	password    string // decrypted
}

type mailPreset struct {
	ID       string `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Security string `json:"security"`
	Help     string `json:"help"` // documentation link: how to create an app password
}

var mailPresets = []mailPreset{
	{ID: "gmail", Host: "smtp.gmail.com", Port: 587, Security: "starttls", Help: "https://support.google.com/accounts/answer/185833"},
	{ID: "outlook", Host: "smtp-mail.outlook.com", Port: 587, Security: "starttls", Help: "https://support.microsoft.com/account-billing/how-to-get-and-use-app-passwords-5896ed9b-4263-e681-128a-a6f2979a7944"},
	{ID: "office365", Host: "smtp.office365.com", Port: 587, Security: "starttls", Help: "https://learn.microsoft.com/exchange/clients-and-mobile-in-exchange-online/authenticated-client-smtp-submission"},
	{ID: "custom", Port: 587, Security: "starttls"},
}

const (
	mailVerifyTTL   = 24 * time.Hour
	mailMaxAttempts = 6
)

func initMailSchema() {
	ensureTable("mail_settings", `CREATE TABLE IF NOT EXISTS mail_settings (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		enabled INTEGER NOT NULL DEFAULT 0,
		preset TEXT NOT NULL DEFAULT 'custom',
		host TEXT NOT NULL DEFAULT '',
		port INTEGER NOT NULL DEFAULT 587,
		security TEXT NOT NULL DEFAULT 'starttls',
		username TEXT NOT NULL DEFAULT '',
		password_enc TEXT NOT NULL DEFAULT '',
		from_addr TEXT NOT NULL DEFAULT '',
		from_name TEXT NOT NULL DEFAULT '',
		public_url TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT '',
		updated_by TEXT NOT NULL DEFAULT ''
	)`, "enabled", "preset", "host", "port", "security", "username", "password_enc", "from_addr", "from_name", "public_url", "updated_at", "updated_by")
	ensureTable("mail_outbox", `CREATE TABLE IF NOT EXISTS mail_outbox (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL DEFAULT 0,
		to_addr TEXT NOT NULL,
		kind TEXT NOT NULL DEFAULT '',
		subject TEXT NOT NULL,
		body TEXT NOT NULL,
		created_at TEXT NOT NULL,
		next_at TEXT NOT NULL DEFAULT '',
		attempts INTEGER NOT NULL DEFAULT 0,
		sent_at TEXT NOT NULL DEFAULT '',
		failed INTEGER NOT NULL DEFAULT 0,
		last_error TEXT NOT NULL DEFAULT ''
	)`, "user_id", "to_addr", "kind", "subject", "body", "created_at", "next_at", "attempts", "sent_at", "failed", "last_error")
	for _, q := range []string{
		`ALTER TABLE users ADD COLUMN email TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN email_verified INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN email_pending TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN email_token TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN email_token_at TEXT NOT NULL DEFAULT ''`,
	} {
		db.Exec(q) // fails harmlessly when the column exists
	}
}

func loadMailSettings() mailSettings {
	m := mailSettings{Preset: "custom", Port: 587, Security: "starttls"}
	var en int
	var pw string
	if db.QueryRow(`SELECT enabled, preset, host, port, security, username, password_enc, from_addr, from_name, public_url, updated_at, updated_by FROM mail_settings WHERE id=1`).
		Scan(&en, &m.Preset, &m.Host, &m.Port, &m.Security, &m.Username, &pw, &m.From, &m.FromName, &m.PublicURL, &m.UpdatedAt, &m.UpdatedBy) == nil {
		m.Enabled = en == 1
		m.PasswordSet = pw != ""
		if pw != "" {
			m.password = decryptValue(pw)
		}
	}
	return m
}

// mailReady: WRM can send mail.
func mailReady() bool {
	m := loadMailSettings()
	return m.Enabled && m.Host != "" && m.From != ""
}

// validMailAddr accepts one plain address (no display name, no line breaks).
func validMailAddr(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || len(s) > 254 || strings.ContainsAny(s, "\r\n<>\",; ") {
		return "", false
	}
	a, err := mail.ParseAddress(s)
	if err != nil || a.Address != s || !strings.Contains(s[strings.LastIndex(s, "@")+1:], ".") {
		return "", false
	}
	return s, true
}

func (m mailSettings) channel(to string) notifyChannel {
	from := m.From
	if m.FromName != "" {
		from = (&mail.Address{Name: m.FromName, Address: m.From}).String()
	}
	return notifyChannel{Kind: "email", cfg: map[string]string{
		"host": m.Host, "port": strconv.Itoa(m.Port), "security": m.Security,
		"username": m.Username, "password": m.password, "from": from, "to": to}}
}

// mailSendNow delivers one message with the current settings.
func mailSendNow(to, subject, body string) error {
	m := loadMailSettings()
	if !m.Enabled || m.Host == "" || m.From == "" {
		return fmt.Errorf("WRM mail is not set up (Admin → Mail)")
	}
	return sendSMTP(m.channel(to), subject, body)
}

// mailQueue stores a message for delivery and returns its id (0 when mail is off).
func mailQueue(userID int, to, kind, subject, body string) int64 {
	if !mailReady() {
		return 0
	}
	now := nowRFC()
	res, err := db.Exec(`INSERT INTO mail_outbox (user_id, to_addr, kind, subject, body, created_at, next_at) VALUES (?,?,?,?,?,?,?)`,
		userID, to, kind, subject, body, now, now)
	if err != nil {
		log.Printf("mail: queue: %v", err)
		return 0
	}
	id, _ := res.LastInsertId()
	mailWake()
	return id
}

var mailWakeCh = make(chan struct{}, 1)

func mailWake() {
	select {
	case mailWakeCh <- struct{}{}:
	default:
	}
}

var mailDeliverMu sync.Mutex

// mailDeliverDue sends the queued messages that are due; it returns how many were sent.
func mailDeliverDue(onlyID int64) int {
	mailDeliverMu.Lock()
	defer mailDeliverMu.Unlock()
	q := `SELECT id, to_addr, subject, body, attempts FROM mail_outbox WHERE sent_at='' AND failed=0 AND next_at<=? ORDER BY id LIMIT 20`
	args := []interface{}{nowRFC()}
	if onlyID > 0 {
		q = `SELECT id, to_addr, subject, body, attempts FROM mail_outbox WHERE id=? AND sent_at='' AND failed=0`
		args = []interface{}{onlyID}
	}
	type item struct {
		id                int64
		to, subject, body string
		attempts          int
	}
	var items []item
	rows, err := db.Query(q, args...)
	if err != nil {
		return 0
	}
	for rows.Next() {
		var it item
		if rows.Scan(&it.id, &it.to, &it.subject, &it.body, &it.attempts) == nil {
			items = append(items, it)
		}
	}
	rows.Close()
	sent := 0
	for _, it := range items {
		err := mailSendNow(it.to, it.subject, it.body)
		if err == nil {
			db.Exec(`UPDATE mail_outbox SET sent_at=?, attempts=attempts+1, last_error='' WHERE id=?`, nowRFC(), it.id)
			sent++
			continue
		}
		n := it.attempts + 1
		delay := []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 2 * time.Hour, 6 * time.Hour}[min(n-1, 4)]
		db.Exec(`UPDATE mail_outbox SET attempts=?, last_error=?, next_at=?, failed=? WHERE id=?`,
			n, truncateStr(err.Error(), 300), time.Now().UTC().Add(delay).Format(time.RFC3339), boolInt(n >= mailMaxAttempts), it.id)
		// the error stays in the outbox (Admin → E-mail); the log only says that it failed
		log.Printf("mail: message %d not delivered (attempt %d of %d), see Admin → E-mail", it.id, n, mailMaxAttempts)
	}
	return sent
}

// mailLoop delivers the outbox while the server runs.
func mailLoop() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		mailDeliverDue(0)
		db.Exec(`DELETE FROM mail_outbox WHERE created_at < ?`, time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339))
		select {
		case <-t.C:
		case <-mailWakeCh:
		}
	}
}

// ─── messages ───

// mailText builds a message in English and Croatian (WRM does not know a user's language
// on the server).
func mailText(en, hr string) string {
	return en + "\n\n— — —\n\n" + hr + "\n\n-- \nWeb Remote Manager PRO"
}

func mailHostName() string {
	m := loadMailSettings()
	if u, err := url.Parse(m.PublicURL); err == nil && u.Host != "" {
		return u.Host
	}
	h, _ := os.Hostname()
	return h
}

// userMailTo is the user's verified address, or "".
func userMailTo(userID int) (string, string) {
	var email, name string
	var verified int
	if db.QueryRow(`SELECT COALESCE(email,''), COALESCE(email_verified,0), username FROM users WHERE id=?`, userID).Scan(&email, &verified, &name) != nil || verified != 1 {
		return "", name
	}
	return email, name
}

// mailSecurityNotice tells a user about a change to their account, if they have a verified
// address and the policy mail_security_notices is on.
func mailSecurityNotice(userID int, kind, enWhat, hrWhat, ip string) int64 {
	if !settingBool("mail_security_notices") {
		return 0
	}
	to, name := userMailTo(userID)
	if to == "" {
		return 0
	}
	when := time.Now().Format("2006-01-02 15:04 MST")
	where := ""
	if ip != "" {
		where = " (" + ip + ")"
	}
	host := mailHostName()
	en := fmt.Sprintf("Hello %s,\n\n%s on %s at %s%s.\n\nIf this was you, there is nothing to do. If it was not, sign in, change your password and tell your administrator.", name, enWhat, host, when, where)
	hr := fmt.Sprintf("Pozdrav %s,\n\n%s na %s u %s%s.\n\nAko si to bio ti, ne trebaš ništa raditi. Ako nisi, prijavi se, promijeni lozinku i javi administratoru.", name, hrWhat, host, when, where)
	return mailQueue(userID, to, kind, "WRM: "+enWhat+" / "+hrWhat, mailText(en, hr))
}

// mailObserve turns audit entries into security notices. It runs on the alert worker.
func mailObserve(in alertInput) {
	byUser := func(name string) int {
		var id int
		db.QueryRow(`SELECT id FROM users WHERE username=?`, name).Scan(&id)
		return id
	}
	switch in.Action {
	case "auth.password_changed":
		mailSecurityNotice(byUser(in.Username), "password_changed", "Your WRM password was changed", "Tvoja WRM lozinka je promijenjena", in.IP)
	case "auth.2fa_enabled":
		mailSecurityNotice(byUser(in.Username), "2fa_enabled", "Two-factor authentication was turned on for your WRM account", "Za tvoj WRM račun uključena je dvofaktorska autentifikacija", in.IP)
	case "auth.2fa_disabled":
		mailSecurityNotice(byUser(in.Username), "2fa_disabled", "Two-factor authentication was turned off for your WRM account", "Za tvoj WRM račun isključena je dvofaktorska autentifikacija", in.IP)
	case "admin.user_password_reset":
		if in.Username != "cli" { // the command line sends its own notice
			mailSecurityNotice(byUser(in.Target), "password_reset", "An administrator ("+in.Username+") reset your WRM password", "Administrator ("+in.Username+") je resetirao tvoju WRM lozinku", in.IP)
		}
	case "admin.user_2fa_reset":
		mailSecurityNotice(byUser(in.Target), "2fa_reset", "An administrator ("+in.Username+") turned off two-factor authentication for your WRM account", "Administrator ("+in.Username+") je isključio dvofaktorsku autentifikaciju za tvoj WRM račun", in.IP)
	}
}

// ─── e-mail address of a user ───

func mailToken() (string, string) {
	b := make([]byte, 24)
	rand.Read(b)
	tok := hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(tok))
	return tok, hex.EncodeToString(sum[:])
}

// mailStartVerify stores addr as the user's pending address and sends the verification link.
func mailStartVerify(userID int, addr string) error {
	m := loadMailSettings()
	if !m.Enabled || m.Host == "" || m.From == "" {
		return fmt.Errorf("WRM cannot send mail yet: an administrator sets it up in Admin → Mail")
	}
	base := strings.TrimRight(m.PublicURL, "/")
	if base == "" {
		return fmt.Errorf("WRM's address for links is not set (Admin → Mail)")
	}
	tok, hash := mailToken()
	var name, old string
	var verified int
	db.QueryRow(`SELECT username, COALESCE(email,''), COALESCE(email_verified,0) FROM users WHERE id=?`, userID).Scan(&name, &old, &verified)
	db.Exec(`UPDATE users SET email_pending=?, email_token=?, email_token_at=? WHERE id=?`, addr, hash, nowRFC(), userID)
	link := base + "/api/auth/verify-email?t=" + tok
	en := fmt.Sprintf("Hello %s,\n\nconfirm that this address belongs to your WRM account on %s by opening this link within 24 hours:\n\n%s\n\nIf you did not ask for this, ignore this message.", name, mailHostName(), link)
	hr := fmt.Sprintf("Pozdrav %s,\n\npotvrdi da ova adresa pripada tvom WRM računu na %s otvaranjem ove poveznice u roku od 24 sata:\n\n%s\n\nAko to nisi tražio, zanemari ovu poruku.", name, mailHostName(), link)
	if mailQueue(userID, addr, "verify", "WRM: confirm your e-mail address / potvrdi e-mail adresu", mailText(en, hr)) == 0 {
		return fmt.Errorf("the message could not be queued")
	}
	if verified == 1 && old != "" && !strings.EqualFold(old, addr) {
		mailSecurityNotice(userID, "email_change", "A new e-mail address ("+addr+") was set for your WRM account", "Za tvoj WRM račun postavljena je nova e-mail adresa ("+addr+")", "")
	}
	return nil
}

// GET /api/auth/verify-email?t=TOKEN — the link in the verification mail.
func apiVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	tok := r.URL.Query().Get("t")
	result := "invalid"
	if len(tok) == 48 {
		sum := sha256.Sum256([]byte(tok))
		hash := hex.EncodeToString(sum[:])
		var id int
		var pending, at, name string
		err := db.QueryRow(`SELECT id, username, email_pending, email_token_at FROM users WHERE email_token=? AND email_token<>''`, hash).Scan(&id, &name, &pending, &at)
		if err == nil {
			ts, ok := parseTS(at)
			if ok && time.Since(ts) < mailVerifyTTL && pending != "" {
				db.Exec(`UPDATE users SET email=?, email_verified=1, email_pending='', email_token='', email_token_at='' WHERE id=?`, pending, id)
				auditLogAs(r, id, name, "user.email_verified", name, map[string]string{"email": pending})
				result = "ok"
			} else {
				db.Exec(`UPDATE users SET email_token='', email_token_at='' WHERE id=?`, id)
				result = "expired"
			}
		}
	}
	http.Redirect(w, r, "/?email_verified="+result, http.StatusSeeOther)
}

// /api/me/email: GET the state, PUT {email} (step-up) starts verification, DELETE removes it.
func apiMeEmailHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		var email, pending string
		var verified int
		db.QueryRow(`SELECT COALESCE(email,''), COALESCE(email_verified,0), COALESCE(email_pending,'') FROM users WHERE id=?`, userID).Scan(&email, &verified, &pending)
		jsonOK(w, map[string]interface{}{"email": email, "verified": verified == 1, "pending": pending, "mail_ready": mailReady(), "self_service": settingBool("user_email_self")})
	case http.MethodPut:
		if !settingBool("user_email_self") {
			jsonError(w, "Only an administrator can set e-mail addresses here", 403)
			return
		}
		var p struct {
			Email string `json:"email"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&p) != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		addr, ok := validMailAddr(p.Email)
		if !ok {
			jsonError(w, "That is not a valid e-mail address", 400)
			return
		}
		// the address will receive security notices and, later, password resets
		if !requireStepUpAlways(w, r, "email_change") {
			return
		}
		if err := mailStartVerify(userID, addr); err != nil {
			jsonError(w, err.Error(), 409)
			return
		}
		auditLog(r, userID, "user.email_change_requested", usernameOf(userID), map[string]string{"email": addr})
		jsonOK(w, map[string]interface{}{"ok": true, "pending": addr})
	case http.MethodDelete:
		if !requireStepUpAlways(w, r, "email_change") {
			return
		}
		db.Exec(`UPDATE users SET email='', email_verified=0, email_pending='', email_token='', email_token_at='' WHERE id=?`, userID)
		auditLog(r, userID, "user.email_removed", usernameOf(userID), nil)
		jsonOK(w, map[string]interface{}{"ok": true})
	default:
		jsonError(w, "Method not allowed", 405)
	}
}

// adminSetUserEmail sets or clears a user's address from Admin → Users. A new address is
// verified by the user, as when they set it themselves.
func adminSetUserEmail(r *http.Request, adminID, userID int, raw string) error {
	raw = strings.TrimSpace(raw)
	name := usernameOf(userID)
	if raw == "" {
		db.Exec(`UPDATE users SET email='', email_verified=0, email_pending='', email_token='', email_token_at='' WHERE id=?`, userID)
		auditLog(r, adminID, "admin.user_email_removed", name, nil)
		return nil
	}
	addr, ok := validMailAddr(raw)
	if !ok {
		return fmt.Errorf("that is not a valid e-mail address")
	}
	if err := mailStartVerify(userID, addr); err != nil {
		return err
	}
	auditLog(r, adminID, "admin.user_email_set", name, map[string]string{"email": addr})
	return nil
}

// ─── Admin → Mail ───
//
//   GET  /api/admin/mail          settings (never the password), presets, recent messages
//   PUT  /api/admin/mail          change (step-up); an empty password keeps the stored one
//   POST /api/admin/mail/test     send a test message now, to the administrator's own confirmed
//                                 address or else to the sender address (never to an address
//                                 typed in the request, so the endpoint cannot mail anyone)

func apiAdminMailHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/mail"), "/")
	switch {
	case rest == "" && r.Method == http.MethodGet:
		jsonOK(w, map[string]interface{}{"settings": loadMailSettings(), "presets": mailPresets, "recent": mailRecent(), "notices": settingBool("mail_security_notices")})
	case rest == "" && r.Method == http.MethodPut:
		var p struct {
			Enabled      bool   `json:"enabled"`
			Preset       string `json:"preset"`
			Host         string `json:"host"`
			Port         int    `json:"port"`
			Security     string `json:"security"`
			Username     string `json:"username"`
			Password     string `json:"password"`
			ForgetStored bool   `json:"clear_password"`
			From         string `json:"from"`
			FromName     string `json:"from_name"`
			PublicURL    string `json:"public_url"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&p) != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		preset := "custom"
		for _, x := range mailPresets {
			if x.ID == p.Preset {
				preset = x.ID
			}
		}
		security := map[string]string{"starttls": "starttls", "tls": "tls", "none": "none"}[p.Security]
		if security == "" {
			jsonError(w, "Security must be starttls, tls or none", 400)
			return
		}
		host := strings.TrimSpace(p.Host)
		if host != "" && (strings.ContainsAny(host, " /:\r\n") || net.ParseIP(host) == nil && !validHostName(host)) {
			jsonError(w, "Invalid SMTP server name", 400)
			return
		}
		if p.Port < 1 || p.Port > 65535 {
			jsonError(w, "Invalid port", 400)
			return
		}
		from := ""
		if strings.TrimSpace(p.From) != "" {
			var ok bool
			if from, ok = validMailAddr(p.From); !ok {
				jsonError(w, "The sender is not a valid e-mail address", 400)
				return
			}
		}
		fromName := strings.TrimSpace(p.FromName)
		if strings.ContainsAny(fromName, "\r\n") || len(fromName) > 100 {
			jsonError(w, "Invalid sender name", 400)
			return
		}
		username := strings.TrimSpace(p.Username)
		if strings.ContainsAny(username, "\r\n") || len(username) > 254 {
			jsonError(w, "Invalid user name", 400)
			return
		}
		pub := strings.TrimRight(strings.TrimSpace(p.PublicURL), "/")
		if pub != "" {
			u, err := url.Parse(pub)
			if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
				jsonError(w, "WRM's address must be an http(s) URL such as https://wrm.example.com", 400)
				return
			}
			pub = u.Scheme + "://" + u.Host + strings.TrimRight(u.EscapedPath(), "/")
		}
		if p.Enabled && (host == "" || from == "" || pub == "") {
			jsonError(w, "To turn mail on, fill in the SMTP server, the sender and WRM's address", 400)
			return
		}
		if security == "none" && username != "" && !isLoopbackHost(host) {
			jsonError(w, "A password is sent only over TLS: choose STARTTLS or TLS", 400)
			return
		}
		if !requireStepUp(w, r, "mail_settings") {
			return
		}
		old := loadMailSettings()
		var pwEnc string
		db.QueryRow(`SELECT password_enc FROM mail_settings WHERE id=1`).Scan(&pwEnc)
		// a stored password never goes to another server: changing the server asks for it again
		if p.Password == "" && pwEnc != "" && (host != old.Host || username != old.Username) {
			jsonError(w, "The SMTP server or user changed: enter the password again", 400)
			return
		}
		if p.Password != "" {
			pwEnc = encryptValue(p.Password)
		}
		if p.ForgetStored {
			pwEnc = ""
		}
		me := usernameOf(adminID)
		_, err := db.Exec(`INSERT INTO mail_settings (id, enabled, preset, host, port, security, username, password_enc, from_addr, from_name, public_url, updated_at, updated_by)
			VALUES (1,?,?,?,?,?,?,?,?,?,?,?,?)
			ON CONFLICT(id) DO UPDATE SET enabled=excluded.enabled, preset=excluded.preset, host=excluded.host, port=excluded.port, security=excluded.security,
			username=excluded.username, password_enc=excluded.password_enc, from_addr=excluded.from_addr, from_name=excluded.from_name, public_url=excluded.public_url,
			updated_at=excluded.updated_at, updated_by=excluded.updated_by`,
			boolInt(p.Enabled), preset, host, p.Port, security, username, pwEnc, from, fromName, pub, nowRFC(), me)
		if err != nil {
			jsonError(w, "Could not save", 500)
			return
		}
		auditLog(r, adminID, "admin.mail_settings_changed", host, map[string]interface{}{"enabled": p.Enabled, "preset": preset, "host": host, "port": p.Port, "security": security, "from": from, "cred_changed": p.Password != "" || p.ForgetStored})
		mailWake()
		jsonOK(w, map[string]interface{}{"ok": true, "settings": loadMailSettings()})
	case rest == "test" && r.Method == http.MethodPost:
		m := loadMailSettings()
		to, _ := userMailTo(adminID)
		if to == "" {
			to = m.From
		}
		if m.Host == "" || m.From == "" {
			jsonError(w, "Fill in and save the SMTP server and the sender first", 400)
			return
		}
		me := usernameOf(adminID)
		err := sendSMTP(m.channel(to), "WRM test message / probna poruka",
			mailText("This is a test message from WRM on "+mailHostName()+", sent by "+me+". Mail works.",
				"Ovo je probna poruka WRM-a na "+mailHostName()+", poslao ju je "+me+". Slanje e-pošte radi."))
		auditLog(r, adminID, "admin.mail_test", to, map[string]interface{}{"ok": err == nil})
		if err != nil {
			jsonError(w, mailErrorHint(m, err), 502)
			return
		}
		jsonOK(w, map[string]interface{}{"ok": true, "to": to})
	default:
		jsonError(w, "Not found", 404)
	}
}

// mailErrorHint adds what usually fixes a failed login with the presets.
func mailErrorHint(m mailSettings, err error) string {
	msg := err.Error()
	if strings.Contains(msg, "SMTP login") {
		switch m.Preset {
		case "gmail":
			msg += " — Gmail needs an app password (Google account → Security → 2-Step Verification → App passwords), not the account password."
		case "outlook", "office365":
			msg += " — Microsoft needs an app password, and in Microsoft 365 an administrator must allow SMTP AUTH for this mailbox."
		}
	}
	return msg
}

func mailRecent() []map[string]interface{} {
	out := []map[string]interface{}{}
	rows, err := db.Query(`SELECT id, to_addr, kind, subject, created_at, attempts, sent_at, failed, last_error FROM mail_outbox ORDER BY id DESC LIMIT 30`)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var to, kind, subject, created, sent, lastErr string
		var attempts, failed int
		if rows.Scan(&id, &to, &kind, &subject, &created, &attempts, &sent, &failed, &lastErr) == nil {
			state := "waiting"
			if sent != "" {
				state = "sent"
			} else if failed == 1 {
				state = "failed"
			}
			out = append(out, map[string]interface{}{"id": id, "to": to, "kind": kind, "subject": subject, "created_at": created,
				"attempts": attempts, "sent_at": sent, "state": state, "last_error": lastErr})
		}
	}
	return out
}
