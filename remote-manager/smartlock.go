package main

import (
	"net"
	"net/http"
)

// ─── SMART SCREEN LOCK (v14.13.0) ────────────────────
//
// The screen lock (lock.go, static/lock.js) adapts to what the user is doing:
//   - sensitive views (Admin, credentials, keys, the vault, password dialogs) lock sooner;
//     screen_lock_sensitive_minutes caps that time for everyone;
//   - live output on screen (a terminal or remote desktop that is receiving data, a following
//     log) may keep WRM unlocked for a while when the user allows it (screen_lock_watch) and
//     the administrator's longest idle time is not exceeded;
//   - a sign-in session whose address moves to another network is locked on the server
//     (screen_lock_network_change): a laptop taken from the office to a café, or a stolen cookie.

// smartLockPolicy is the administrator's part, sent to the page in mePayload.screen_lock.smart.
func smartLockPolicy() map[string]interface{} {
	return map[string]interface{}{
		"sensitive_minutes": settingInt("screen_lock_sensitive_minutes"),
		"watch":             settingBool("screen_lock_watch"),
		"network_change":    settingBool("screen_lock_network_change"),
	}
}

// sameNetwork: both addresses are in the same /24 (IPv4) or /64 (IPv6). Unparseable addresses
// count as the same network, so a proxy that hides addresses never locks anyone.
func sameNetwork(a, b string) bool {
	ia, ib := net.ParseIP(a), net.ParseIP(b)
	if ia == nil || ib == nil {
		return true
	}
	if a4, b4 := ia.To4(), ib.To4(); a4 != nil || b4 != nil {
		if a4 == nil || b4 == nil {
			return false
		}
		return a4[0] == b4[0] && a4[1] == b4[1] && a4[2] == b4[2]
	}
	m := net.CIDRMask(64, 128)
	return ia.Mask(m).Equal(ib.Mask(m))
}

// lockOnNetworkChange locks a sign-in session that is used from another network than before.
// It records the new address, so the unlock (with the password, PIN or passkey) holds there.
func lockOnNetworkChange(r *http.Request, s *authSession, storedIP string) {
	if s.LockedAt != "" || storedIP == "" || !settingBool("screen_lock_network_change") || screenLockMode() == "browser" {
		return
	}
	ip := clientIP(r)
	if sameNetwork(storedIP, ip) {
		return
	}
	now := nowRFC()
	res, err := db.Exec(`UPDATE auth_sessions SET locked_at=?, ip=? WHERE id=? AND locked_at=''`, now, truncateStr(ip, 60), s.ID)
	if err != nil {
		return
	}
	if n, _ := res.RowsAffected(); n == 1 {
		s.LockedAt = now
		auditLog(r, s.UserID, "auth.locked", "", map[string]string{"reason": "network_change", "from": storedIP, "to": ip})
	}
}
