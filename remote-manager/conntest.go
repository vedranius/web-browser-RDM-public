package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// POST /api/connections/test
// Body: a Connection (as in the edit form). When "id" is set and password/key fields are
// empty, the stored secrets of that connection are used.
func apiConnectionTestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	var c Connection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	if c.ID > 0 && userOwnsConnection(c.ID, userID) {
		if stored, err := loadConnection(c.ID); err == nil {
			if c.Password == "" {
				c.Password = stored.Password
			}
			if c.PrivateKey == "" {
				c.PrivateKey = stored.PrivateKey
			}
		}
	}
	c.Host = strings.TrimSpace(c.Host)
	if c.Host == "" {
		jsonError(w, "Host is required", 400)
		return
	}
	if c.Protocol == "" {
		c.Protocol = "SSH"
	}
	if c.AuthMethod == "" {
		c.AuthMethod = "PASSWORD"
	}
	c.Host = ensurePort(c.Host, c.Protocol)

	start := time.Now()
	var detail string
	var err error
	if isFTP(c) {
		fc, e := dialFTP(c)
		if e == nil {
			if wd, e2 := fc.CurrentDir(); e2 == nil {
				detail = "FTP login OK, working directory " + wd
			}
			fc.Quit()
		}
		err = e
	} else {
		client, e := getSSHClient(c)
		if e == nil {
			detail = "SSH login OK (" + string(client.ServerVersion()) + ")"
			client.Close()
		}
		err = e
	}
	ms := time.Since(start).Milliseconds()
	if err != nil {
		jsonOK(w, map[string]interface{}{"ok": false, "message": err.Error(), "latency_ms": ms})
		return
	}
	jsonOK(w, map[string]interface{}{"ok": true, "message": detail, "latency_ms": ms})
}
