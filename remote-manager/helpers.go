package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
)

// loadConnectionRaw reads a connection by ID and decrypts its secrets.
func loadConnectionRaw(connID int) (Connection, error) {
	var c Connection
	var uid *int
	err := db.QueryRow(`SELECT id,name,protocol,host,username,auth_method,password,private_key,key_path,folder_id,user_id FROM connections WHERE id=?`, connID).
		Scan(&c.ID, &c.Name, &c.Protocol, &c.Host, &c.Username, &c.AuthMethod, &c.Password, &c.PrivateKey, &c.KeyPath, &c.FolderID, &uid)
	if err != nil {
		return c, fmt.Errorf("connection not found")
	}
	if uid != nil {
		c.UserID = *uid
	}
	decryptConnectionSecrets(&c)
	return c, nil
}

// loadConnection is loadConnectionRaw with the host normalised to host:port for dialing.
func loadConnection(connID int) (Connection, error) {
	c, err := loadConnectionRaw(connID)
	if err != nil {
		return c, err
	}
	c.Host = ensurePort(c.Host, c.Protocol)
	return c, nil
}

func isFTP(c Connection) bool {
	p := strings.ToUpper(c.Protocol)
	return p == "FTP" || p == "FTPS"
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// attachmentHeader builds a Content-Disposition value that is safe for any file name
// (quotes, non-ASCII characters, etc.).
func attachmentHeader(name string) string {
	v := mime.FormatMediaType("attachment", map[string]string{"filename": name})
	if v == "" {
		return `attachment; filename="download"`
	}
	return v
}

// joinRemote joins a remote directory and a name using forward slashes.
func joinRemote(dir, name string) string {
	if dir == "" || dir == "." {
		return name
	}
	return strings.TrimSuffix(dir, "/") + "/" + name
}

// ─── WEBSOCKET ORIGIN CHECK ──────────────────────────

// checkWSOrigin blocks cross-site WebSocket hijacking: the Origin host must match the
// Host (or X-Forwarded-Host) of the request. Set WRM_ALLOW_ANY_ORIGIN=1 to disable, or
// WRM_ALLOWED_ORIGINS=a.example.com,b.example.com to allow extra hosts.
func checkWSOrigin(r *http.Request) bool {
	return originAllowed(r, r.Header.Get("Origin"))
}

// originAllowed reports whether an Origin/Referer URL belongs to this server.
func originAllowed(r *http.Request, origin string) bool {
	if origin == "" || os.Getenv("WRM_ALLOW_ANY_ORIGIN") == "1" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	oh := strings.ToLower(hostOnly(u.Host))
	for _, h := range []string{r.Host, r.Header.Get("X-Forwarded-Host")} {
		for _, part := range strings.Split(h, ",") {
			if part = strings.TrimSpace(part); part != "" && strings.ToLower(hostOnly(part)) == oh {
				return true
			}
		}
	}
	for _, allowed := range strings.Split(os.Getenv("WRM_ALLOWED_ORIGINS"), ",") {
		if a := strings.TrimSpace(allowed); a != "" && strings.ToLower(hostOnly(a)) == oh {
			return true
		}
	}
	log.Printf("Cross-origin request rejected: %s (host %s)", origin, r.Host)
	return false
}

func hostOnly(hp string) string {
	if h, _, err := net.SplitHostPort(hp); err == nil {
		return strings.Trim(h, "[]")
	}
	return strings.Trim(hp, "[]")
}

// ─── LOGIN RATE LIMITING ─────────────────────────────

type loginAttempt struct {
	fails       int
	first       time.Time
	lockedUntil time.Time
}

var (
	loginAttempts   = map[string]*loginAttempt{}
	loginAttemptsMu sync.Mutex
)

// Per IP address; accounts additionally lock after login_max_failures (policy).
const (
	loginMaxFails   = 20
	loginFailWindow = 10 * time.Minute
	loginLockout    = 5 * time.Minute
)

// loginBlocked reports whether the IP is temporarily locked out and for how long.
func loginBlocked(ip string) (bool, time.Duration) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	a := loginAttempts[ip]
	if a == nil {
		return false, 0
	}
	if d := time.Until(a.lockedUntil); d > 0 {
		return true, d
	}
	return false, 0
}

func recordLoginResult(ip string, success bool) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	if success {
		delete(loginAttempts, ip)
		return
	}
	now := time.Now()
	if len(loginAttempts) > 20000 { // bound memory under a distributed attack
		for k, v := range loginAttempts {
			if now.Sub(v.first) > loginFailWindow && now.After(v.lockedUntil) {
				delete(loginAttempts, k)
			}
		}
	}
	a := loginAttempts[ip]
	if a == nil || now.Sub(a.first) > loginFailWindow {
		a = &loginAttempt{first: now}
		loginAttempts[ip] = a
	}
	a.fails++
	if a.fails >= loginMaxFails {
		a.lockedUntil = now.Add(loginLockout)
		a.fails = 0
		a.first = now
		log.Printf("Login locked for %s for %s after repeated failures", ip, loginLockout)
	}
}

// ─── ERROR CLASSIFICATION ────────────────────────────

func isConnLostErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sftp.ErrSSHFxConnectionLost) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	s := err.Error()
	for _, frag := range []string{"connection lost", "broken pipe", "connection reset", "use of closed network connection", "EOF"} {
		if strings.Contains(s, frag) {
			return true
		}
	}
	return false
}
