package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─── PASSKEYS (v14.13.0) ─────────────────────────────
//
// Passkeys sign in with the device's own unlock: Windows Hello (face, fingerprint, PIN), Touch ID
// and Face ID, Android and iOS biometrics, a FIDO2 security key, or a phone near the computer
// (the browser shows a QR code and checks over Bluetooth that the phone is close).
//
// A passkey can:
//   - sign in on its own (passkey_login=passwordless): it proves possession and, through the
//     device's biometric or PIN, the user, so it also satisfies require_2fa;
//   - be the second step after the password (instead of, or next to, an authenticator code);
//   - unlock the screen lock and confirm a step-up.
// Browsers allow WebAuthn only on https:// pages (and http://localhost), on a host name, never
// on an IP address. GET /api/me/passkeys says when this WRM address cannot use them, and why.

func initPasskeySchema() {
	ensureTable("passkeys", `CREATE TABLE IF NOT EXISTS passkeys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		credential_id TEXT NOT NULL UNIQUE,
		public_key BLOB NOT NULL,
		alg INTEGER NOT NULL DEFAULT 0,
		sign_count INTEGER NOT NULL DEFAULT 0,
		rp_id TEXT NOT NULL DEFAULT '',
		name TEXT NOT NULL DEFAULT '',
		aaguid TEXT NOT NULL DEFAULT '',
		transports TEXT NOT NULL DEFAULT '',
		backup INTEGER NOT NULL DEFAULT 0,
		uv INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL DEFAULT '',
		last_used_at TEXT NOT NULL DEFAULT '',
		last_used_ip TEXT NOT NULL DEFAULT '')`,
		"id", "user_id", "credential_id", "public_key", "alg", "sign_count", "rp_id", "name", "aaguid", "transports",
		"backup", "uv", "created_at", "last_used_at", "last_used_ip")
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_passkeys_user ON passkeys(user_id)`)
	// The WebAuthn user handle: random, so a passkey never carries the WRM user ID.
	db.Exec(`ALTER TABLE users ADD COLUMN webauthn_handle TEXT NOT NULL DEFAULT ''`)
}

func passkeysEnabled() bool { return settingBool("passkeys") }

// passkeyLoginMode: passwordless | second_factor | off (off: unlock and step-up only).
func passkeyLoginMode() string {
	if !passkeysEnabled() {
		return "off"
	}
	return getSetting("passkey_login")
}

func userPasskeyCount(userID int) int {
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM passkeys WHERE user_id=?`, userID).Scan(&n)
	return n
}

// userPasskeySecondFactor: a password sign-in of this account asks for a passkey (or a code).
func userPasskeySecondFactor(u *userRow) bool {
	return passkeyLoginMode() != "off" && userPasskeyCount(u.ID) > 0
}

// needsSecondStep decides whether a correct password still needs a second step, and which.
// A registered passkey is a second factor once 2FA is required for the account, or when the
// account has no authenticator code but chose a passkey for sign-in (passkey_login allows it).
func needsSecondStep(u *userRow) (bool, []string) {
	methods := []string{}
	if u.TOTPEnabled {
		methods = append(methods, "totp")
	}
	if userPasskeySecondFactor(u) {
		methods = append(methods, "passkey")
	}
	if u.TOTPEnabled {
		return true, methods
	}
	if len(methods) > 0 && (mfaRequiredByPolicy(u.IsAdmin) || u.BreakGlass) {
		return true, methods
	}
	return false, nil
}

// secondStepAnswer is the reply to a correct password when a second step follows.
func secondStepAnswer(u *userRow, token string, methods []string) map[string]interface{} {
	return map[string]interface{}{"mfa_required": true, "mfa_token": token, "mfa_methods": methods}
}

// ─── relying party ───────────────────────────────────

// passkeyRP is the relying party for this request: the configured passkey_rp_id, otherwise the
// host name the browser used. ok is false (with a reason for people) when browsers would refuse
// WebAuthn here.
func passkeyRP(r *http.Request) (rp waRelyingParty, ok bool, why string) {
	host := r.Host
	if trustProxy {
		if h := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Host"), ",")[0]); h != "" {
			host = h
		}
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.ToLower(strings.Trim(host, "[]"))
	id := strings.ToLower(strings.TrimSpace(getSetting("passkey_rp_id")))
	if id == "" {
		id = host
	}
	for _, o := range strings.Split(getSetting("passkey_origins"), ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			rp.Origins = append(rp.Origins, o)
		}
	}
	rp.ID = id
	switch {
	case net.ParseIP(host) != nil:
		return rp, false, "WRM is opened by its IP address. Browsers allow passkeys only on a host name (for example wrm.example.com) with HTTPS."
	case !validRPID(id):
		return rp, false, "The passkey domain (passkey_rp_id) is not a valid host name."
	case !(host == id || strings.HasSuffix(host, "."+id)):
		return rp, false, "This address (" + host + ") is not inside the passkey domain " + id + "."
	case !isHTTPS(r) && host != "localhost" && !strings.HasSuffix(host, ".localhost"):
		return rp, false, "Passkeys need HTTPS. Open WRM over https:// (or on localhost)."
	}
	return rp, true, ""
}

// ─── ceremonies in flight ────────────────────────────

type waCeremony struct {
	Challenge []byte
	Purpose   string // register | login | mfa | unlock | step_up
	UserID    int    // 0 for a passwordless sign-in (the passkey names the account)
	SessionID int
	MFANonce  string
	RP        waRelyingParty
	Expires   time.Time
}

var waCeremonies = struct {
	sync.Mutex
	m map[string]*waCeremony
}{m: map[string]*waCeremony{}}

const waCeremonyTTL = 5 * time.Minute

func waStart(c *waCeremony) (string, error) {
	c.Challenge = make([]byte, 32)
	if _, err := rand.Read(c.Challenge); err != nil {
		return "", err
	}
	c.Expires = time.Now().Add(waCeremonyTTL)
	id := randomToken(18)
	waCeremonies.Lock()
	defer waCeremonies.Unlock()
	now := time.Now()
	if len(waCeremonies.m) > 2000 {
		for k, v := range waCeremonies.m {
			if now.After(v.Expires) {
				delete(waCeremonies.m, k)
			}
		}
		if len(waCeremonies.m) > 5000 {
			return "", errors.New("too many sign-ins in progress; try again in a minute")
		}
	}
	waCeremonies.m[id] = c
	return id, nil
}

// waTake removes and returns a ceremony (each challenge is answered once).
func waTake(id, purpose string) *waCeremony {
	waCeremonies.Lock()
	defer waCeremonies.Unlock()
	c := waCeremonies.m[id]
	if c == nil {
		return nil
	}
	delete(waCeremonies.m, id)
	if c.Purpose != purpose || time.Now().After(c.Expires) {
		return nil
	}
	return c
}

type passkeyRow struct {
	ID        int
	UserID    int
	CredID    string
	PublicKey []byte
	SignCount uint32
	RPID      string
	Name      string
}

func loadPasskey(where string, args ...interface{}) (*passkeyRow, error) {
	var p passkeyRow
	var sc int64
	err := db.QueryRow(`SELECT id, user_id, credential_id, public_key, sign_count, rp_id, name FROM passkeys WHERE `+where, args...).
		Scan(&p.ID, &p.UserID, &p.CredID, &p.PublicKey, &sc, &p.RPID, &p.Name)
	if err != nil {
		return nil, err
	}
	p.SignCount = uint32(sc)
	return &p, nil
}

func userWebAuthnHandle(userID int) string {
	var h string
	db.QueryRow(`SELECT webauthn_handle FROM users WHERE id=?`, userID).Scan(&h)
	if h == "" {
		h = randomToken(32)
		db.Exec(`UPDATE users SET webauthn_handle=? WHERE id=? AND webauthn_handle=''`, h, userID)
		db.QueryRow(`SELECT webauthn_handle FROM users WHERE id=?`, userID).Scan(&h)
	}
	return h
}

func userPasskeyDescriptors(userID int, rpID string) []map[string]interface{} {
	out := []map[string]interface{}{}
	rows, err := db.Query(`SELECT credential_id, transports FROM passkeys WHERE user_id=? AND rp_id=?`, userID, rpID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, tr string
		rows.Scan(&id, &tr)
		d := map[string]interface{}{"type": "public-key", "id": id}
		if tr != "" {
			d["transports"] = strings.Split(tr, ",")
		}
		out = append(out, d)
	}
	return out
}

// assertionOptions is the publicKey argument of navigator.credentials.get().
func assertionOptions(c *waCeremony, allow []map[string]interface{}, uv string) map[string]interface{} {
	o := map[string]interface{}{
		"challenge": waB64.EncodeToString(c.Challenge), "rpId": c.RP.ID, "timeout": 120000, "userVerification": uv,
	}
	if allow != nil {
		o["allowCredentials"] = allow
	}
	return o
}

type waCredentialIn struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Response waAssertionResponse `json:"response"`
}

// checkAssertion verifies an answer to an assertion ceremony and returns the passkey it used.
// userID > 0 requires the passkey to belong to that account.
func checkAssertion(r *http.Request, c *waCeremony, cred waCredentialIn, userID int, requireUV bool) (*passkeyRow, *waAssertion, error) {
	raw, err := b64urlDecode(cred.ID)
	if err != nil || len(raw) == 0 || len(raw) > 1023 {
		return nil, nil, errors.New("unknown passkey")
	}
	pk, err := loadPasskey(`credential_id=?`, waB64.EncodeToString(raw))
	if err != nil {
		return nil, nil, errors.New("this passkey is not registered in WRM (it may have been removed)")
	}
	if userID > 0 && pk.UserID != userID {
		return nil, nil, errors.New("this passkey belongs to another account")
	}
	if cred.Response.UserHandle != "" {
		h, err := b64urlDecode(cred.Response.UserHandle)
		if err != nil || subtle.ConstantTimeCompare(h, []byte(userWebAuthnHandle(pk.UserID))) != 1 {
			return nil, nil, errors.New("the passkey names another account")
		}
	}
	rp := c.RP
	if pk.RPID != rp.ID {
		return nil, nil, errors.New("this passkey was made for " + pk.RPID + ", not " + rp.ID)
	}
	a, err := verifyAssertion(cred.Response, c.Challenge, rp, pk.PublicKey, requireUV)
	if err != nil {
		return nil, nil, err
	}
	// A signature counter that goes back points to a cloned authenticator (synced passkeys keep 0).
	if a.SignCount != 0 || pk.SignCount != 0 {
		if a.SignCount <= pk.SignCount {
			auditLogAs(r, pk.UserID, usernameOf(pk.UserID), "auth.passkey_counter", pk.Name, map[string]interface{}{"stored": pk.SignCount, "got": a.SignCount})
			return nil, nil, errors.New("this passkey's counter went back: it may have been copied. Remove it and add it again")
		}
	}
	db.Exec(`UPDATE passkeys SET sign_count=?, last_used_at=?, last_used_ip=? WHERE id=?`, int64(a.SignCount), nowRFC(), truncateStr(clientIP(r), 60), pk.ID)
	return pk, a, nil
}

// ─── signed-in user: list, add, rename, remove ───────

// GET    /api/me/passkeys                       → state and the account's passkeys
// POST   /api/me/passkeys/begin                 → creation options (needs a fresh confirmation)
// POST   /api/me/passkeys/finish {ceremony, name, credential}
// PATCH  /api/me/passkeys/{id} {name}
// DELETE /api/me/passkeys/{id}
func apiMePasskeysHandler(w http.ResponseWriter, r *http.Request) {
	s, err := lookupSession(r)
	if err != nil {
		jsonError(w, "Unauthorized", 401)
		return
	}
	// A session that must set up 2FA first may add a passkey as its second factor.
	if s.Restricted != "" && !(s.Restricted == "mfa" && strings.HasPrefix(r.URL.Path, "/api/me/passkeys")) {
		jsonError(w, restrictedMessage(s.Restricted), 403)
		return
	}
	if screenLockRefuses(s, r) {
		jsonError(w, "WRM is locked; unlock it first", http.StatusLocked)
		return
	}
	userID := s.UserID
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/me/passkeys"), "/")
	switch {
	case rest == "" && r.Method == http.MethodGet:
		rp, ok, why := passkeyRP(r)
		jsonOK(w, map[string]interface{}{
			"enabled": passkeysEnabled(), "available": ok && passkeysEnabled(), "reason": why, "rp_id": rp.ID,
			"login": passkeyLoginMode(), "passkeys": passkeyList(userID),
		})
	case rest == "begin" && r.Method == http.MethodPost:
		if !passkeysEnabled() {
			jsonError(w, "The administrator turned passkeys off", 403)
			return
		}
		rp, ok, why := passkeyRP(r)
		if !ok {
			jsonError(w, why, 400)
			return
		}
		if s.Restricted == "" && !requireStepUpAlways(w, r, "passkey_add") {
			return
		}
		if userPasskeyCount(userID) >= 20 {
			jsonError(w, "An account can have at most 20 passkeys; remove one first", 400)
			return
		}
		u, err := loadUser("id=?", userID)
		if err != nil {
			jsonError(w, "Unauthorized", 401)
			return
		}
		c := &waCeremony{Purpose: "register", UserID: userID, SessionID: s.ID, RP: rp}
		id, err := waStart(c)
		if err != nil {
			jsonError(w, err.Error(), 503)
			return
		}
		display := u.DisplayName
		if display == "" {
			display = u.Username
		}
		jsonOK(w, map[string]interface{}{"ceremony": id, "publicKey": map[string]interface{}{
			"challenge": waB64.EncodeToString(c.Challenge),
			"rp":        map[string]string{"id": rp.ID, "name": "WRM PRO"},
			"user":      map[string]string{"id": userWebAuthnHandle(userID), "name": u.Username, "displayName": display},
			"pubKeyCredParams": []map[string]interface{}{
				{"type": "public-key", "alg": coseAlgES256}, {"type": "public-key", "alg": coseAlgEdDSA}, {"type": "public-key", "alg": coseAlgRS256}},
			"timeout": 180000, "attestation": "none",
			"authenticatorSelection": map[string]interface{}{"residentKey": "preferred", "requireResidentKey": false, "userVerification": "preferred"},
			"excludeCredentials":     userPasskeyDescriptors(userID, rp.ID),
		}})
	case rest == "finish" && r.Method == http.MethodPost:
		var in struct {
			Ceremony   string `json:"ceremony"`
			Name       string `json:"name"`
			Credential struct {
				ID       string                `json:"id"`
				Response waAttestationResponse `json:"response"`
			} `json:"credential"`
		}
		if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil {
			jsonError(w, "Bad JSON", 400)
			return
		}
		c := waTake(in.Ceremony, "register")
		if c == nil || c.UserID != userID || c.SessionID != s.ID {
			jsonError(w, "The request expired. Start adding the passkey again.", 400)
			return
		}
		nc, err := verifyRegistration(in.Credential.Response, c.Challenge, c.RP, false)
		if err != nil {
			auditLog(r, userID, "auth.passkey_add_failed", "", map[string]string{"reason": truncateStr(err.Error(), 200)})
			jsonError(w, "The passkey was not accepted: "+err.Error(), 400)
			return
		}
		name := truncateStr(strings.TrimSpace(in.Name), 60)
		if name == "" {
			name = passkeyDefaultName(r, nc)
		}
		credID := waB64.EncodeToString(nc.ID)
		_, err = db.Exec(`INSERT INTO passkeys (user_id, credential_id, public_key, alg, sign_count, rp_id, name, aaguid, transports, backup, uv, created_at)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, userID, credID, nc.PublicKey, nc.Alg, int64(nc.SignCount), c.RP.ID, name,
			hex.EncodeToString(nc.AAGUID), strings.Join(nc.Transports, ","), boolToInt(nc.Backup), boolToInt(nc.UV), nowRFC())
		if err != nil {
			jsonError(w, "This passkey is already registered", 400)
			return
		}
		auditLog(r, userID, "auth.passkey_added", name, map[string]interface{}{"rp_id": c.RP.ID, "synced": nc.Backup, "verified": nc.UV})
		if s.Restricted == "mfa" {
			db.Exec(`UPDATE auth_sessions SET restricted='' WHERE id=?`, s.ID)
		}
		jsonOK(w, map[string]interface{}{"ok": true, "name": name, "passkeys": passkeyList(userID)})
	default:
		id, err := strconv.Atoi(rest)
		if err != nil || id <= 0 || s.Restricted != "" {
			jsonError(w, "Not found", 404)
			return
		}
		pk, err := loadPasskey(`id=? AND user_id=?`, id, userID)
		if err != nil {
			jsonError(w, "Passkey not found", 404)
			return
		}
		switch r.Method {
		case http.MethodPatch:
			var in struct {
				Name string `json:"name"`
			}
			json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&in)
			name := truncateStr(strings.TrimSpace(in.Name), 60)
			if name == "" {
				jsonError(w, "Enter a name", 400)
				return
			}
			db.Exec(`UPDATE passkeys SET name=? WHERE id=?`, name, pk.ID)
			jsonOK(w, map[string]bool{"ok": true})
		case http.MethodDelete:
			if !requireStepUp(w, r, "passkey_remove") {
				return
			}
			db.Exec(`DELETE FROM passkeys WHERE id=?`, pk.ID)
			auditLog(r, userID, "auth.passkey_removed", pk.Name, nil)
			jsonOK(w, map[string]interface{}{"ok": true, "passkeys": passkeyList(userID)})
		default:
			jsonError(w, "Method not allowed", 405)
		}
	}
}

func passkeyDefaultName(r *http.Request, nc *waNewCredential) string {
	ua := r.UserAgent()
	dev := "Passkey"
	for _, p := range []struct{ needle, name string }{
		{"Windows", "Windows Hello"}, {"iPhone", "iPhone"}, {"iPad", "iPad"}, {"Macintosh", "Mac"},
		{"Android", "Android"}, {"CrOS", "ChromeOS"}, {"Linux", "Linux"},
	} {
		if strings.Contains(ua, p.needle) {
			dev = p.name
			break
		}
	}
	for _, t := range nc.Transports {
		switch t {
		case "hybrid":
			dev = "Phone"
		case "usb", "nfc":
			if dev == "Passkey" || !nc.Backup {
				dev = "Security key"
			}
		}
	}
	return dev + " " + time.Now().Format("2006-01-02")
}

func passkeyList(userID int) []map[string]interface{} {
	out := []map[string]interface{}{}
	rows, err := db.Query(`SELECT id, name, rp_id, transports, backup, uv, created_at, last_used_at FROM passkeys WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, backup, uv int
		var name, rpID, tr, created, used string
		rows.Scan(&id, &name, &rpID, &tr, &backup, &uv, &created, &used)
		transports := []string{}
		if tr != "" {
			transports = strings.Split(tr, ",")
		}
		out = append(out, map[string]interface{}{"id": id, "name": name, "rp_id": rpID, "transports": transports,
			"synced": backup == 1, "verified": uv == 1, "created_at": created, "last_used_at": used})
	}
	return out
}

// ─── sign-in ─────────────────────────────────────────

// POST /api/auth/webauthn/login/begin  {mfa_token?}
// POST /api/auth/webauthn/login/finish {ceremony, credential, mfa_token?}
// Without mfa_token the passkey signs in on its own (the browser offers the passkeys saved for
// this site). With mfa_token it is the second step after the password.
func apiWebAuthnLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	if msg := signInTransportError(r); msg != "" {
		jsonError(w, msg, 403)
		return
	}
	ip := clientIP(r)
	if blocked, wait := loginBlocked(ip); blocked {
		jsonError(w, fmt.Sprintf("Too many failed attempts. Try again in %d s.", int(wait.Seconds())+1), 429)
		return
	}
	var in struct {
		Ceremony   string         `json:"ceremony"`
		MFAToken   string         `json:"mfa_token"`
		Credential waCredentialIn `json:"credential"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	mfaUser, mfaNonce := 0, ""
	if in.MFAToken != "" {
		uid, nonce, ok := parseMFAToken(in.MFAToken)
		if !ok {
			jsonError(w, "The sign-in step expired. Enter your password again.", 401)
			return
		}
		mfaUser, mfaNonce = uid, nonce
	}
	mode := passkeyLoginMode()
	switch strings.TrimPrefix(r.URL.Path, "/api/auth/webauthn/login/") {
	case "begin":
		if mode == "off" || (mfaUser == 0 && mode != "passwordless") {
			jsonError(w, "Signing in with a passkey is turned off", 403)
			return
		}
		rp, ok, why := passkeyRP(r)
		if !ok {
			jsonError(w, why, 400)
			return
		}
		c := &waCeremony{Purpose: "login", UserID: mfaUser, MFANonce: mfaNonce, RP: rp}
		var allow []map[string]interface{}
		uv := "required"
		if mfaUser > 0 {
			c.Purpose, uv = "mfa", "preferred"
			allow = userPasskeyDescriptors(mfaUser, rp.ID)
			if len(allow) == 0 {
				jsonError(w, "This account has no passkey for this address", 400)
				return
			}
		}
		id, err := waStart(c)
		if err != nil {
			jsonError(w, err.Error(), 503)
			return
		}
		jsonOK(w, map[string]interface{}{"ceremony": id, "publicKey": assertionOptions(c, allow, uv)})
	case "finish":
		purpose := "login"
		if mfaUser > 0 {
			purpose = "mfa"
		}
		c := waTake(in.Ceremony, purpose)
		if c == nil || c.UserID != mfaUser || c.MFANonce != mfaNonce {
			jsonError(w, "The sign-in request expired. Try again.", 400)
			return
		}
		pk, _, err := checkAssertion(r, c, in.Credential, mfaUser, mfaUser == 0)
		if err != nil {
			recordLoginResult(ip, false)
			uid := mfaUser
			if pk != nil {
				uid = pk.UserID
			}
			auditLogAs(r, uid, usernameOf(uid), "auth.login_failed", "", map[string]string{"method": "passkey", "reason": truncateStr(err.Error(), 200)})
			jsonError(w, "The passkey was not accepted: "+err.Error(), 401)
			return
		}
		u, err := loadUser("id=?", pk.UserID)
		if err != nil {
			jsonError(w, "Invalid sign-in", 401)
			return
		}
		if msg := passkeySignInRefusal(r, u, mfaUser == 0); msg != "" {
			auditLogAs(r, u.ID, u.Username, "auth.login_failed", "", map[string]string{"method": "passkey", "reason": msg})
			jsonError(w, msg, 403)
			return
		}
		method := "passkey"
		if mfaUser > 0 {
			method = "password+passkey"
		}
		if identityID := oidcMFAIdentity(mfaNonce); mfaUser > 0 && identityID > 0 {
			finishLogin(w, r, u, "oidc+passkey")
			db.Exec(`UPDATE auth_sessions SET identity_id=? WHERE id=(SELECT MAX(id) FROM auth_sessions WHERE user_id=?) AND identity_id=0`, identityID, u.ID)
			return
		}
		finishPasskeyLogin(w, r, u, method, pk)
	default:
		jsonError(w, "Not found", 404)
	}
}

// passkeySignInRefusal applies the checks of a password sign-in to a passkey sign-in.
func passkeySignInRefusal(r *http.Request, u *userRow, passwordless bool) string {
	if u.Disabled {
		return "This account is disabled. Contact an administrator."
	}
	if t, ok := parseTS(u.LockedUntil); ok && time.Now().Before(t) {
		return fmt.Sprintf("Account temporarily locked after repeated failed sign-ins. Try again in %d min or ask an administrator.", int(time.Until(t).Minutes())+1)
	}
	if u.BreakGlass && !breakGlassIPAllowed(clientIP(r)) {
		return "This emergency (break-glass) account may only sign in from the networks the administrator allowed."
	}
	if !passwordless {
		return ""
	}
	if localLoginLimited() && !u.BreakGlass {
		return "Sign in with your organisation's account. Sign-in with a local WRM credential is kept for emergency administrators only."
	}
	// A directory account must still be allowed by its directory.
	if id := userLDAPIdentity(u.ID); id != nil {
		if why := ldapStillAllowed(id); why != "" {
			return why
		}
	} else if userSSOOnly(u.ID) {
		return "Your account signs in through your organisation's identity provider: use its button. The passkey still unlocks the screen and confirms actions."
	}
	return ""
}

// ldapStillAllowed looks the linked directory account up with the service account ("" = allowed).
func ldapStillAllowed(id *ldapIdentity) string {
	if id.Revoked {
		return "Your directory account may no longer use WRM."
	}
	d, err := loadLDAPDirectory(id.DirectoryID)
	if err != nil {
		return "Your directory is not configured any more. Ask your administrator."
	}
	c, err := d.connect()
	if err != nil {
		log.Printf("LDAP %s: passkey sign-in check failed: %s", d.Name, ldapErrorSummary(err))
		return "WRM could not check your account with the directory " + d.Name + ". Try again later."
	}
	defer c.Close()
	e, err := d.findByID(c, id.Key)
	if err != nil {
		return "Your directory account was not found."
	}
	if why := ldapAccountBlocked(e); why != "" {
		return "Your directory account may not sign in (" + why + ")."
	}
	return ""
}

func finishPasskeyLogin(w http.ResponseWriter, r *http.Request, u *userRow, method string, pk *passkeyRow) {
	// A verified passkey is a second factor in itself: 2FA policies do not restrict this session.
	if method == "passkey" || method == "password+passkey" {
		restricted := ""
		if u.MustChange {
			restricted = "password"
		}
		if lockdownRefuses(u.ID) {
			auditLogAs(r, u.ID, u.Username, "auth.login_failed", "", map[string]string{"method": method, "reason": "lockdown"})
			jsonError(w, lockdownMessage, 403)
			return
		}
		if err := createAuthSession(w, r, u.ID, restricted); err != nil {
			jsonError(w, "Server error", 500)
			return
		}
		recordLoginResult(clientIP(r), true)
		db.Exec(`UPDATE users SET last_login_at=?, last_login_ip=?, failed_logins=0, locked_until='' WHERE id=?`, nowRFC(), clientIP(r), u.ID)
		auditLogAs(r, u.ID, u.Username, "auth.login", "", map[string]string{"method": method, "passkey": pk.Name, "user_agent": truncateStr(r.UserAgent(), 160)})
		if u.BreakGlass {
			breakGlassSignedIn(r, u)
		}
		jsonOK(w, mePayload(u.ID, restricted))
		return
	}
	finishLogin(w, r, u, method)
}

// ─── unlock and step-up ──────────────────────────────

// POST /api/auth/webauthn/begin {purpose: "unlock" | "step_up"} → assertion options for the
// signed-in account's passkeys. The answer goes to /api/auth/lock/verify or /api/auth/step-up
// as {"passkey": {ceremony, credential}}.
func apiWebAuthnBeginHandler(w http.ResponseWriter, r *http.Request) {
	s, err := lookupSession(r)
	if err != nil || s.Restricted != "" {
		jsonError(w, "Unauthorized", 401)
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in struct {
		Purpose string `json:"purpose"`
	}
	json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&in)
	if in.Purpose != "unlock" && in.Purpose != "step_up" {
		jsonError(w, "Unknown purpose", 400)
		return
	}
	if !passkeysEnabled() {
		jsonError(w, "The administrator turned passkeys off", 403)
		return
	}
	rp, ok, why := passkeyRP(r)
	if !ok {
		jsonError(w, why, 400)
		return
	}
	allow := userPasskeyDescriptors(s.UserID, rp.ID)
	if len(allow) == 0 {
		jsonError(w, "You have no passkey for this address. Add one in Settings → Security.", 400)
		return
	}
	c := &waCeremony{Purpose: in.Purpose, UserID: s.UserID, SessionID: s.ID, RP: rp}
	id, err := waStart(c)
	if err != nil {
		jsonError(w, err.Error(), 503)
		return
	}
	jsonOK(w, map[string]interface{}{"ceremony": id, "publicKey": assertionOptions(c, allow, "required")})
}

type passkeyProof struct {
	Ceremony   string         `json:"ceremony"`
	Credential waCredentialIn `json:"credential"`
}

// verifyPasskeyProof checks a passkey answer for the signed-in session (unlock, step-up).
func verifyPasskeyProof(r *http.Request, s *authSession, purpose string, p *passkeyProof) (*passkeyRow, error) {
	if p == nil || !passkeysEnabled() {
		return nil, errors.New("passkeys are turned off")
	}
	c := waTake(p.Ceremony, purpose)
	if c == nil || c.UserID != s.UserID || c.SessionID != s.ID {
		return nil, errors.New("the request expired; try again")
	}
	pk, _, err := checkAssertion(r, c, p.Credential, s.UserID, true)
	return pk, err
}

// userHasPasskeyForUnlock: the lock screen may offer the passkey button.
func userHasPasskeyForUnlock(userID int) bool {
	return passkeysEnabled() && userPasskeyCount(userID) > 0
}

// ─── administrators ──────────────────────────────────

// GET    /api/admin/passkeys/{userID} → the user's passkeys
// DELETE /api/admin/passkeys/{userID} → removes them all (lost phone, departed employee)
func apiAdminPasskeysHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	uid, err := strconv.Atoi(strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/passkeys"), "/"))
	if err != nil || uid <= 0 {
		jsonError(w, "Bad ID", 400)
		return
	}
	name := usernameOf(uid)
	if name == "" {
		jsonError(w, "User not found", 404)
		return
	}
	switch r.Method {
	case http.MethodGet:
		jsonOK(w, map[string]interface{}{"passkeys": passkeyList(uid)})
	case http.MethodDelete:
		if !requireStepUp(w, r, "passkey_reset") {
			return
		}
		res, _ := db.Exec(`DELETE FROM passkeys WHERE user_id=?`, uid)
		n, _ := res.RowsAffected()
		auditLog(r, adminID, "admin.user_passkeys_removed", name, map[string]int64{"removed": n})
		jsonOK(w, map[string]interface{}{"ok": true, "removed": n})
	default:
		jsonError(w, "Method not allowed", 405)
	}
}
