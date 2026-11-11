package main

import (
	"crypto/rand"
	"crypto/subtle"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// First-run setup code. On a fresh install the first account becomes the administrator.
// Without a check, whoever reaches the port first (anyone on the network) would own the
// server. serve() therefore creates a one-time code while no account exists: it is printed
// to the console/log and written next to the database. The first registration needs it,
// except from the WRM host itself (a direct loopback connection, not through a proxy).

const setupCodeFile = "first-admin-setup-code.txt"

var setupCode struct {
	sync.Mutex
	code string
	path string
}

// firstUserMu serialises the "is this the first account?" check with the insert, so two
// concurrent first registrations cannot both become administrators.
var firstUserMu sync.Mutex

// initSetupCode creates the code when the database has no account yet, and removes a
// stale code file once accounts exist.
func initSetupCode(dbPath string) {
	path := filepath.Join(filepath.Dir(dbPath), setupCodeFile)
	var n int
	db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n)
	setupCode.Lock()
	defer setupCode.Unlock()
	setupCode.path = path
	if n > 0 {
		setupCode.code = ""
		os.Remove(path)
		return
	}
	setupCode.code = newSetupCode()
	msg := "First-run setup code (needed to create the administrator account from another computer): " + setupCode.code
	if err := os.WriteFile(path, []byte(setupCode.code+"\n"), 0o600); err == nil {
		msg += " (also in " + path + ")"
	}
	log.Print(msg)
}

// newSetupCode returns 12 characters in three groups, without look-alike characters.
func newSetupCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	var sb strings.Builder
	for i, c := range b {
		if i > 0 && i%4 == 0 {
			sb.WriteByte('-')
		}
		sb.WriteByte(alphabet[int(c)%len(alphabet)])
	}
	return sb.String()
}

// setupCodeRequired reports whether a first registration from r must carry the code.
func setupCodeRequired(r *http.Request) bool {
	setupCode.Lock()
	active := setupCode.code != ""
	setupCode.Unlock()
	return active && !directLoopback(r)
}

// directLoopback is true for a connection from the WRM host itself that no proxy forwarded.
func directLoopback(r *http.Request) bool {
	if r.Header.Get("X-Forwarded-For") != "" || r.Header.Get("X-Real-IP") != "" || r.Header.Get("Forwarded") != "" {
		return false
	}
	ip := net.ParseIP(hostOnly(r.RemoteAddr))
	return ip != nil && ip.IsLoopback()
}

// setupCodeMatches compares a typed code (case, spaces and dashes ignored) in constant time.
func setupCodeMatches(typed string) bool {
	norm := func(s string) string {
		return strings.ToUpper(strings.NewReplacer("-", "", " ", "").Replace(strings.TrimSpace(s)))
	}
	setupCode.Lock()
	want := setupCode.code
	setupCode.Unlock()
	if want == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(norm(typed)), []byte(norm(want))) == 1
}

// clearSetupCode retires the code after the first account was created.
func clearSetupCode() {
	setupCode.Lock()
	defer setupCode.Unlock()
	if setupCode.code == "" {
		return
	}
	setupCode.code = ""
	if setupCode.path != "" {
		os.Remove(setupCode.path)
	}
}
