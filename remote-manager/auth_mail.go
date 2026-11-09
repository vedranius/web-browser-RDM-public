package main

// Account links sent by e-mail (v14.6.2): "Forgot password", invitations of accounts an
// administrator creates, and confirmation of self-registered accounts.
//
// A link carries a random 192-bit token; the database keeps only its SHA-256. A token is
// single-use, expires, and a newer token of the same kind replaces the older one. Setting a
// password through a link never signs in by itself: the user then signs in normally, so
// two-factor authentication and every sign-in policy still apply. A user with 2FA must also
// give a code (or a recovery code) to set a password through a "Forgot password" link, so
// that access to the mailbox alone is not enough.

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	authLinkReset    = "reset"    // Forgot password
	authLinkActivate = "activate" // an administrator created the account and invited the user
	authLinkRegister = "register" // a self-registered account confirms its e-mail address

	authResetTTL    = time.Hour
	authActivateTTL = 72 * time.Hour
	authRegisterTTL = 24 * time.Hour
	// a user gets at most one new "Forgot password" mail in this interval
	authResetGap = 2 * time.Minute
	// self-registered accounts never confirmed are removed after this
	authPendingKeep = 7 * 24 * time.Hour
)

func initAuthMailSchema() {
	ensureTable("auth_links", `CREATE TABLE IF NOT EXISTS auth_links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		kind TEXT NOT NULL,
		token_hash TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL,
		used_at TEXT NOT NULL DEFAULT '',
		ip TEXT NOT NULL DEFAULT ''
	)`, "user_id", "kind", "token_hash", "email", "created_at", "expires_at", "used_at", "ip")
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_auth_links_user ON auth_links(user_id, kind)`)
	db.Exec(`ALTER TABLE users ADD COLUMN activation_pending INTEGER NOT NULL DEFAULT 0`) // fails harmlessly when it exists
}

type authLink struct {
	ID      int64
	UserID  int
	Kind    string
	Email   string
	Created time.Time
}

func authLinkHash(tok string) string {
	sum := sha256.Sum256([]byte(tok))
	return hex.EncodeToString(sum[:])
}

// authLinkCreate stores a new link for the user and returns its token. Older unused links of
// the same kind stop working.
func authLinkCreate(userID int, kind, email string, ttl time.Duration, ip string) string {
	tok, hash := mailToken()
	now := time.Now().UTC()
	db.Exec(`UPDATE auth_links SET used_at=? WHERE user_id=? AND kind=? AND used_at=''`, "replaced", userID, kind)
	db.Exec(`INSERT INTO auth_links (user_id, kind, token_hash, email, created_at, expires_at, ip) VALUES (?,?,?,?,?,?,?)`,
		userID, kind, hash, email, now.Format(time.RFC3339), now.Add(ttl).Format(time.RFC3339), ip)
	return tok
}

// authLinkLookup finds an unused link: status is "ok", "expired" or "invalid".
func authLinkLookup(tok string, kinds ...string) (*authLink, string) {
	if len(tok) != 48 {
		return nil, "invalid"
	}
	var l authLink
	var created, expires, used string
	if db.QueryRow(`SELECT id, user_id, kind, email, created_at, expires_at, used_at FROM auth_links WHERE token_hash=?`, authLinkHash(tok)).
		Scan(&l.ID, &l.UserID, &l.Kind, &l.Email, &created, &expires, &used) != nil {
		return nil, "invalid"
	}
	known := false
	for _, k := range kinds {
		known = known || k == l.Kind
	}
	if !known {
		return nil, "invalid"
	}
	if used != "" {
		return nil, "expired"
	}
	if t, ok := parseTS(expires); !ok || time.Now().After(t) {
		return nil, "expired"
	}
	l.Created, _ = parseTS(created)
	return &l, "ok"
}

// authLinkUse marks the link used; false when another request was first.
func authLinkUse(id int64) bool {
	res, err := db.Exec(`UPDATE auth_links SET used_at=? WHERE id=? AND used_at=''`, nowRFC(), id)
	if err != nil {
		return false
	}
	n, _ := res.RowsAffected()
	return n == 1
}

// authLinksCleanup removes old links and self-registered accounts never confirmed.
func authLinksCleanup() {
	old := time.Now().Add(-authPendingKeep).UTC().Format(time.RFC3339)
	db.Exec(`DELETE FROM auth_links WHERE expires_at < ?`, old)
	db.Exec(`DELETE FROM users WHERE activation_pending=1 AND disabled=1 AND COALESCE(last_login_at,'')='' AND created_at < ?`, old)
}

// mailLinkBase is WRM's public address for links in mail (Admin → E-mail).
func mailLinkBase() string {
	return strings.TrimRight(loadMailSettings().PublicURL, "/")
}

// mailLinksReady: WRM can send mail and knows the address to put in links.
func mailLinksReady() bool {
	return mailReady() && mailLinkBase() != ""
}

// passwordResetMode is the "Forgot password" policy: all, users (not administrators) or off.
func passwordResetMode() string {
	switch v := getSetting("password_reset_mail"); v {
	case "users", "off":
		return v
	}
	return "all"
}

// passwordResetOffered: the sign-in page shows "Forgot password?".
func passwordResetOffered() bool {
	return passwordResetMode() != "off" && !lockdownActive() && mailLinksReady()
}

// registrationNeedsEmail: self-registration asks for an e-mail address and creates the
// account only once the address is confirmed.
func registrationNeedsEmail() bool {
	return settingBool("registration_verify_email") && mailLinksReady()
}

func userActivationPending(userID int) bool {
	var v int
	db.QueryRow(`SELECT COALESCE(activation_pending,0) FROM users WHERE id=?`, userID).Scan(&v)
	return v == 1
}

// resetRefusal says why the account cannot use "Forgot password" ("" = it can).
func resetRefusal(u *userRow) string {
	switch {
	case passwordResetMode() == "off":
		return "turned off"
	case passwordResetMode() == "users" && u.IsAdmin:
		return "not for administrators"
	case !mailLinksReady():
		return "mail not set up"
	case u.Disabled:
		return "account disabled"
	case userLDAPIdentity(u.ID) != nil:
		return "directory account"
	case userSSOOnly(u.ID):
		return "single sign-on account"
	case localLoginLimited() && !u.BreakGlass:
		return "local sign-in limited"
	case lockdownRefuses(u.ID):
		return "lockdown"
	}
	if to, _ := userMailTo(u.ID); to == "" {
		return "no confirmed address"
	}
	return ""
}

// A small per-address limit for "Forgot password" requests, apart from the sign-in limit so
// that asking for a link does not lock anyone out of signing in.
var forgotLimit = struct {
	sync.Mutex
	m map[string][]time.Time
}{m: map[string][]time.Time{}}

const (
	forgotPerWindow = 5
	forgotWindow    = 15 * time.Minute
)

func forgotAllowed(ip string) bool {
	forgotLimit.Lock()
	defer forgotLimit.Unlock()
	now := time.Now()
	if len(forgotLimit.m) > 20000 {
		for k, v := range forgotLimit.m {
			if len(v) == 0 || now.Sub(v[len(v)-1]) > forgotWindow {
				delete(forgotLimit.m, k)
			}
		}
	}
	keep := forgotLimit.m[ip][:0]
	for _, t := range forgotLimit.m[ip] {
		if now.Sub(t) < forgotWindow {
			keep = append(keep, t)
		}
	}
	if len(keep) >= forgotPerWindow {
		forgotLimit.m[ip] = keep
		return false
	}
	forgotLimit.m[ip] = append(keep, now)
	return true
}

// POST /api/auth/forgot {login} — a username or a confirmed e-mail address. The answer is the
// same whether or not an account matched, so the form cannot be used to find accounts.
func apiAuthForgotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	var p struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&p); err != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	if msg := signInTransportError(r); msg != "" {
		jsonError(w, msg, 403)
		return
	}
	if !passwordResetOffered() {
		jsonError(w, "Password reset by e-mail is not available. Ask an administrator.", 403)
		return
	}
	ip := clientIP(r)
	if !forgotAllowed(ip) {
		jsonError(w, "Too many requests. Try again in a few minutes.", 429)
		return
	}
	login := strings.TrimSpace(p.Login)
	answer := map[string]interface{}{"ok": true}
	if login == "" || len(login) > 254 {
		jsonOK(w, answer)
		return
	}
	ids := []int{}
	if u, err := loadUser("username=?", login); err == nil {
		ids = append(ids, u.ID)
	} else if addr, ok := validMailAddr(login); ok {
		rows, err := db.Query(`SELECT id FROM users WHERE email_verified=1 AND lower(email)=lower(?) LIMIT 5`, addr)
		if err == nil {
			for rows.Next() {
				var id int
				rows.Scan(&id)
				ids = append(ids, id)
			}
			rows.Close()
		}
	}
	base := mailLinkBase()
	for _, id := range ids {
		u, err := loadUser("id=?", id)
		if err != nil {
			continue
		}
		if why := resetRefusal(u); why != "" {
			auditLogAs(r, u.ID, u.Username, "auth.reset_link_refused", "", map[string]string{"reason": why})
			continue
		}
		var last string
		db.QueryRow(`SELECT COALESCE(MAX(created_at),'') FROM auth_links WHERE user_id=? AND kind=?`, u.ID, authLinkReset).Scan(&last)
		if t, ok := parseTS(last); ok && time.Since(t) < authResetGap {
			continue
		}
		to, name := userMailTo(u.ID)
		tok := authLinkCreate(u.ID, authLinkReset, to, authResetTTL, ip)
		link := base + "/?reset=" + tok
		en := fmt.Sprintf("Hello %s,\n\nsomeone (%s) asked to set a new password for your WRM account on %s. To do it, open this link within one hour:\n\n%s\n\nThe link works once. If you did not ask for this, ignore this message: your password stays as it is.", name, ip, mailHostName(), link)
		hr := fmt.Sprintf("Pozdrav %s,\n\nnetko (%s) je zatražio novu lozinku za tvoj WRM račun na %s. Za to otvori ovu poveznicu u roku od sat vremena:\n\n%s\n\nPoveznica radi jednom. Ako to nisi tražio, zanemari ovu poruku: lozinka ostaje ista.", name, ip, mailHostName(), link)
		mailQueue(u.ID, to, "reset", "WRM: set a new password / nova lozinka", mailText(en, hr))
		auditLogAs(r, u.ID, u.Username, "auth.reset_link_requested", "", nil)
	}
	jsonOK(w, answer)
}

// /api/auth/reset — the page behind a reset or invitation link.
//
//	GET  ?t=TOKEN                  → {kind, username, totp, min_length}
//	POST {token, password, code}   → sets the password; the user then signs in
func apiAuthResetHandler(w http.ResponseWriter, r *http.Request) {
	if msg := signInTransportError(r); msg != "" {
		jsonError(w, msg, 403)
		return
	}
	ip := clientIP(r)
	if blocked, wait := loginBlocked(ip); blocked {
		jsonError(w, fmt.Sprintf("Too many attempts. Try again in %d s.", int(wait.Seconds())+1), 429)
		return
	}
	switch r.Method {
	case http.MethodGet:
		l, st := authLinkLookup(r.URL.Query().Get("t"), authLinkReset, authLinkActivate)
		if l == nil {
			if st == "invalid" {
				recordLoginResult(ip, false)
			}
			jsonOK(w, map[string]interface{}{"status": st})
			return
		}
		u, err := loadUser("id=?", l.UserID)
		if err != nil {
			jsonOK(w, map[string]interface{}{"status": "invalid"})
			return
		}
		jsonOK(w, map[string]interface{}{"status": "ok", "kind": l.Kind, "username": u.Username,
			"totp": l.Kind == authLinkReset && u.TOTPEnabled, "min_length": settingInt("password_min_length")})
	case http.MethodPost:
		var p struct {
			Token    string `json:"token"`
			Password string `json:"password"`
			Code     string `json:"code"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&p); err != nil {
			jsonError(w, "Bad JSON", 400)
			return
		}
		l, st := authLinkLookup(p.Token, authLinkReset, authLinkActivate)
		if l == nil {
			if st == "invalid" {
				recordLoginResult(ip, false)
			}
			jsonError(w, "This link is no longer valid. Ask for a new one.", 400)
			return
		}
		u, err := loadUser("id=?", l.UserID)
		if err != nil {
			jsonError(w, "This link is no longer valid. Ask for a new one.", 400)
			return
		}
		if l.Kind == authLinkReset {
			if why := resetRefusal(u); why != "" {
				jsonError(w, "A new password cannot be set this way for this account. Ask an administrator.", 403)
				return
			}
		} else if u.Disabled || userLDAPIdentity(u.ID) != nil || userSSOOnly(u.ID) || lockdownRefuses(u.ID) {
			jsonError(w, "This account cannot be activated this way. Ask an administrator.", 403)
			return
		}
		if err := validatePassword(p.Password); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		// Mail alone is not enough for an account protected by two-factor authentication.
		if l.Kind == authLinkReset && u.TOTPEnabled {
			code := strings.TrimSpace(p.Code)
			if code == "" {
				jsonError(w, "Enter a code from your authenticator app or a recovery code.", 400)
				return
			}
			if method, _, err := consumeUserSecondFactor(u, code); err != nil || method == "" {
				recordLoginResult(ip, false)
				auditLogAs(r, u.ID, u.Username, "auth.reset_link_refused", "", map[string]string{"reason": "wrong second factor"})
				jsonError(w, "The code is not valid.", 400)
				return
			}
		}
		hash, err := hashPassword(p.Password)
		if err != nil {
			jsonError(w, "Server error", 500)
			return
		}
		if !authLinkUse(l.ID) {
			jsonError(w, "This link is no longer valid. Ask for a new one.", 400)
			return
		}
		db.Exec(`UPDATE users SET password_hash=?, must_change_password=0, failed_logins=0, locked_until='', password_changed_at=? WHERE id=?`, hash, nowRFC(), u.ID)
		db.Exec(`UPDATE auth_links SET used_at=? WHERE user_id=? AND kind IN (?,?) AND used_at=''`, "replaced", u.ID, authLinkReset, authLinkActivate)
		if l.Kind == authLinkActivate && l.Email != "" {
			// opening the invitation proves the address
			db.Exec(`UPDATE users SET email=?, email_verified=1, email_pending='', email_token='', email_token_at='' WHERE id=? AND lower(email_pending)=lower(?)`, l.Email, u.ID, l.Email)
		}
		n := revokeUserSessions(u.ID, 0)
		recordLoginResult(ip, true)
		auditLogAs(r, u.ID, u.Username, "auth.reset_link_used", "", map[string]interface{}{"kind": l.Kind, "sessions_revoked": n})
		if l.Kind == authLinkReset {
			mailSecurityNotice(u.ID, "password_changed", "Your WRM password was changed with a link sent by e-mail", "Tvoja WRM lozinka promijenjena je poveznicom poslanom e-poštom", ip)
		}
		jsonOK(w, map[string]interface{}{"ok": true, "username": u.Username})
	default:
		jsonError(w, "Method not allowed", 405)
	}
}

// GET /api/auth/activate?t=TOKEN — the link that confirms a self-registered account.
func apiAuthActivateHandler(w http.ResponseWriter, r *http.Request) {
	result := "invalid"
	l, st := authLinkLookup(r.URL.Query().Get("t"), authLinkRegister)
	if l == nil {
		result = st
	} else if authLinkUse(l.ID) {
		res, err := db.Exec(`UPDATE users SET disabled=0, activation_pending=0, email=?, email_verified=1, email_pending='', email_token='', email_token_at='' WHERE id=? AND activation_pending=1`, l.Email, l.UserID)
		if err == nil {
			if n, _ := res.RowsAffected(); n == 1 {
				name := usernameOf(l.UserID)
				auditLogAs(r, l.UserID, name, "auth.account_activated", "", map[string]string{"email": l.Email})
				result = "ok"
			}
		}
	}
	http.Redirect(w, r, "/?activated="+result, http.StatusSeeOther)
}

// registerPending creates a self-registered account that waits for its e-mail address to be
// confirmed, and sends the link.
func registerPending(r *http.Request, username, hash, addr string) error {
	now := nowRFC()
	res, err := db.Exec(`INSERT INTO users (username, password_hash, is_admin, disabled, activation_pending, email_pending, created_at, password_changed_at) VALUES (?,?,0,1,1,?,?,?)`,
		username, hash, addr, now, now)
	if err != nil {
		return err
	}
	id64, _ := res.LastInsertId()
	id := int(id64)
	tok := authLinkCreate(id, authLinkRegister, addr, authRegisterTTL, clientIP(r))
	link := mailLinkBase() + "/api/auth/activate?t=" + tok
	en := fmt.Sprintf("Hello %s,\n\nto finish creating your WRM account on %s, confirm this e-mail address by opening this link within 24 hours:\n\n%s\n\nIf you did not register, ignore this message: the account is removed after a week.", username, mailHostName(), link)
	hr := fmt.Sprintf("Pozdrav %s,\n\nza dovršetak izrade tvog WRM računa na %s potvrdi ovu e-mail adresu otvaranjem ove poveznice u roku od 24 sata:\n\n%s\n\nAko se nisi registrirao, zanemari ovu poruku: račun se uklanja nakon tjedan dana.", username, mailHostName(), link)
	mailQueue(id, addr, "activate", "WRM: confirm your account / potvrdi račun", mailText(en, hr))
	auditLogAs(r, id, username, "auth.register", "", map[string]interface{}{"admin": false, "pending_confirmation": true})
	return nil
}

// inviteUser sends the invitation of an account an administrator created: the user sets the
// password through the link, and opening it confirms the address.
func inviteUser(r *http.Request, adminID, userID int, username, addr string) {
	db.Exec(`UPDATE users SET email_pending=?, email_token='', email_token_at='' WHERE id=?`, addr, userID)
	tok := authLinkCreate(userID, authLinkActivate, addr, authActivateTTL, clientIP(r))
	link := mailLinkBase() + "/?activate=" + tok
	admin := usernameOf(adminID)
	en := fmt.Sprintf("Hello,\n\n%s created the WRM account %s for you on %s. Choose your password by opening this link within 3 days:\n\n%s\n\nAfter that, sign in with the name %s.", admin, username, mailHostName(), link, username)
	hr := fmt.Sprintf("Pozdrav,\n\n%s ti je otvorio WRM račun %s na %s. Odaberi lozinku otvaranjem ove poveznice u roku od 3 dana:\n\n%s\n\nZatim se prijavi imenom %s.", admin, username, mailHostName(), link, username)
	mailQueue(userID, addr, "invite", "WRM: your new account / tvoj novi račun", mailText(en, hr))
	auditLog(r, adminID, "admin.user_invited", username, map[string]string{"email": addr})
}
