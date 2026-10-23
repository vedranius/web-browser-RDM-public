package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

// ─── SCREEN LOCK (#115) ──────────────────────────────
//
// The lock detaches the interface from the screen after a period without
// activity. It is a presence check, not a second session: the cookie and the
// running terminals, tunnels and AI sessions stay alive, and every request
// still passes the normal authentication and CSRF checks. Unlocking asks for
// the signed-in account's password; the attempt is audited either way.

var lockAttempts = struct {
	sync.Mutex
	m map[int][]time.Time
}{m: map[int][]time.Time{}}

// lockThrottled allows five attempts per minute per account.
func lockThrottled(userID int) bool {
	lockAttempts.Lock()
	defer lockAttempts.Unlock()
	now := time.Now()
	keep := []time.Time{}
	for _, t := range lockAttempts.m[userID] {
		if now.Sub(t) < time.Minute {
			keep = append(keep, t)
		}
	}
	if len(keep) >= 5 {
		lockAttempts.m[userID] = keep
		return true
	}
	lockAttempts.m[userID] = append(keep, now)
	return false
}

// POST /api/auth/lock/verify {"password": "..."}
func apiAuthLockVerifyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var p struct {
		Password string `json:"password"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&p) != nil || p.Password == "" {
		jsonError(w, "Enter your password", 400)
		return
	}
	if lockThrottled(userID) {
		jsonError(w, "Too many attempts; wait a minute", 429)
		return
	}
	var hash string
	if db.QueryRow(`SELECT password_hash FROM users WHERE id=? AND disabled=0`, userID).Scan(&hash) != nil {
		jsonError(w, "Account not found", 403)
		return
	}
	if !checkPassword(hash, p.Password) {
		auditLog(r, userID, "auth.unlock_failed", "", nil)
		jsonError(w, "The password is not correct", 403)
		return
	}
	auditLog(r, userID, "auth.unlocked", "", nil)
	jsonOK(w, map[string]bool{"ok": true})
}
