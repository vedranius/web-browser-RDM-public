package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/pkg/sftp"
)

// ─── ADDONS: WRM CAPABILITY LAYERS ───────────────────
//
// An addon may use the same building blocks as WRM itself: the viewer's
// connections, commands and files on their servers, and server metrics. It
// always acts as the person who looks at its window, with that person's own
// rights, never more:
//
//   - connections.read  lists the viewer's connections (name, protocol, host,
//     folder, tags). Logins, passwords, keys and notes never leave WRM.
//   - servers.exec      runs a command on one of the viewer's SSH/SFTP
//     connections. Every command waits for the viewer's approval in WRM
//     (outside the panel, so the addon cannot click it) unless an administrator
//     allowed this addon to run without asking. Output is capped and
//     secret-looking values are masked.
//   - files.read        reads a file over SFTP; sensitive files (keys, shadow,
//     .env …) need approval like a command, and secrets are masked.
//   - files.write       writes a file over SFTP; every write waits for approval.
//   - monitoring.read   reads the current metrics of a server.
//
// The administrator can turn these layers off or limit them to roles with the
// addon_server_access policy. Every call is audited.

const (
	addonWRMOutputMax  = 256 << 10
	addonWRMFileMax    = 1 << 20
	addonWRMApprovalTT = 2 * time.Minute
	addonWRMPerMinute  = 30
)

var addonWRMCaps = map[string]bool{"connections.read": true, "servers.exec": true, "files.read": true, "files.write": true, "monitoring.read": true}

func initAddonWRMSchema() {
	ensureTable("addon_wrm_grants", `CREATE TABLE IF NOT EXISTS addon_wrm_grants (
		owner_id INTEGER NOT NULL DEFAULT 0,
		addon_id TEXT NOT NULL,
		auto INTEGER NOT NULL DEFAULT 0,
		granted_by TEXT NOT NULL DEFAULT '',
		granted_at TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (owner_id, addon_id)
	)`)
}

func addonServerAccessAllowed(userID int) bool {
	v := getSetting("addon_server_access")
	return v != "off" && policyAllows(userID, v)
}

// addonWRMAuto tells whether an administrator allowed the addon to run
// commands and writes without asking each time.
func addonWRMAuto(ref addonRef) bool {
	var n int
	db.QueryRow(`SELECT auto FROM addon_wrm_grants WHERE owner_id=? AND addon_id=?`, ref.Owner, ref.ID).Scan(&n)
	return n == 1
}

func addonWRMSetAuto(ref addonRef, auto bool, by string) {
	v := 0
	if auto {
		v = 1
	}
	db.Exec(`INSERT INTO addon_wrm_grants (owner_id, addon_id, auto, granted_by, granted_at) VALUES (?,?,?,?,?)
		ON CONFLICT(owner_id, addon_id) DO UPDATE SET auto=excluded.auto, granted_by=excluded.granted_by, granted_at=excluded.granted_at`,
		ref.Owner, ref.ID, v, by, nowRFC())
}

func addonWRMForget(ref addonRef) {
	db.Exec(`DELETE FROM addon_wrm_grants WHERE owner_id=? AND addon_id=?`, ref.Owner, ref.ID)
}

// ── approvals ────────────────────────────────────────

type addonWRMPending struct {
	ID       string
	UserID   int
	Instance int
	Ref      addonRef
	Op       string
	Req      addonWRMRequest
	Summary  string
	Expires  time.Time
}

var addonWRMApprovals = struct {
	sync.Mutex
	m map[string]*addonWRMPending
}{m: map[string]*addonWRMPending{}}

var addonWRMRate = struct {
	sync.Mutex
	m map[int][]time.Time
}{m: map[int][]time.Time{}}

func addonWRMRateOK(instance int) bool {
	addonWRMRate.Lock()
	defer addonWRMRate.Unlock()
	now := time.Now()
	keep := addonWRMRate.m[instance][:0]
	for _, t := range addonWRMRate.m[instance] {
		if now.Sub(t) < time.Minute {
			keep = append(keep, t)
		}
	}
	if len(keep) >= addonWRMPerMinute {
		addonWRMRate.m[instance] = keep
		return false
	}
	addonWRMRate.m[instance] = append(keep, now)
	return true
}

func addonWRMPark(p *addonWRMPending) {
	addonWRMApprovals.Lock()
	defer addonWRMApprovals.Unlock()
	now := time.Now()
	for k, v := range addonWRMApprovals.m {
		if now.After(v.Expires) {
			delete(addonWRMApprovals.m, k)
		}
	}
	p.ID = randomToken(12)
	p.Expires = now.Add(addonWRMApprovalTT)
	addonWRMApprovals.m[p.ID] = p
}

// addonWRMTake hands out a parked request exactly once, to the person it was
// parked for, from the same instance.
func addonWRMTake(id string, userID, instance int) *addonWRMPending {
	addonWRMApprovals.Lock()
	defer addonWRMApprovals.Unlock()
	p := addonWRMApprovals.m[id]
	if p == nil || p.UserID != userID || p.Instance != instance || time.Now().After(p.Expires) {
		return nil
	}
	delete(addonWRMApprovals.m, id)
	return p
}

// ── requests ─────────────────────────────────────────

type addonWRMRequest struct {
	Connection int    `json:"connection"`
	Command    string `json:"command,omitempty"`
	Path       string `json:"path,omitempty"`
	Content    string `json:"content,omitempty"`
	Encoding   string `json:"encoding,omitempty"` // "base64" for bytes
	TimeoutS   int    `json:"timeout_s,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

type addonWRMConn struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Protocol string   `json:"protocol"`
	Host     string   `json:"host"`
	Folder   string   `json:"folder,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Exec     bool     `json:"exec"`
}

// addonWRMConnections lists what connections.read shows: nothing that logs in.
func addonWRMConnections(userID int) []addonWRMConn {
	out := []addonWRMConn{}
	rows, err := db.Query(`SELECT c.id, c.name, c.protocol, c.host, COALESCE(f.name,''), COALESCE(c.tags,'') FROM connections c
		LEFT JOIN folders f ON f.id=c.folder_id WHERE c.user_id=? ORDER BY c.name COLLATE NOCASE`, userID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var c addonWRMConn
		var tags string
		if rows.Scan(&c.ID, &c.Name, &c.Protocol, &c.Host, &c.Folder, &tags) != nil {
			continue
		}
		c.Tags = parseTags(tags)
		c.Host = hostOnly(c.Host)
		c.Exec = c.Protocol == "SSH" || c.Protocol == "SFTP"
		out = append(out, c)
	}
	return out
}

// addonWRMConnection loads one of the viewer's SSH/SFTP connections.
func addonWRMConnection(userID, id int) (Connection, error) {
	c, err := loadConnection(id)
	if err != nil || c.UserID != userID {
		return Connection{}, addonErrf("E_WRM_CONNECTION", "Connection %d is not one of yours", id)
	}
	if c.Protocol != "SSH" && c.Protocol != "SFTP" {
		return Connection{}, addonErrf("E_WRM_CONNECTION", "Connection %s is %s; commands and files need SSH or SFTP", c.Name, c.Protocol)
	}
	return c, nil
}

func addonWRMCleanPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" || len(p) > 1024 || strings.ContainsAny(p, "\x00\r\n") {
		return "", addonErrf("E_WRM_REQUEST", "Give an absolute path or one starting with ~/")
	}
	if strings.HasPrefix(p, "~/") {
		return "~/" + strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(p, "~/")), "/"), nil
	}
	if !strings.HasPrefix(p, "/") {
		return "", addonErrf("E_WRM_REQUEST", "Give an absolute path or one starting with ~/")
	}
	return path.Clean(p), nil
}

// addonWRMPrepare checks a request and says whether it needs the viewer's
// approval; summary is what the approval dialog shows.
func addonWRMPrepare(rec addonRecord, userID int, op string, req *addonWRMRequest) (summary string, ask bool, err error) {
	need := map[string]string{"exec": "servers.exec", "read": "files.read", "write": "files.write", "metrics": "monitoring.read"}[op]
	if need == "" {
		return "", false, addonErrf("E_WRM_REQUEST", "Unknown operation %q", op)
	}
	if !addonHasCap(rec, need) {
		return "", false, addonErrf("E_CAPABILITY", "The addon does not declare the %s capability", need)
	}
	c, err := addonWRMConnection(userID, req.Connection)
	if err != nil {
		return "", false, err
	}
	auto := addonWRMAuto(addonRef{Owner: rec.OwnerID, ID: rec.ID})
	switch op {
	case "exec":
		req.Command = strings.TrimSpace(req.Command)
		if req.Command == "" || len(req.Command) > 8192 || strings.ContainsRune(req.Command, 0) {
			return "", false, addonErrf("E_WRM_REQUEST", "Give a command of 1-8192 characters")
		}
		if req.TimeoutS <= 0 {
			req.TimeoutS = 30
		}
		if req.TimeoutS > 300 {
			req.TimeoutS = 300
		}
		return fmt.Sprintf("Run on %s:\n%s", c.Name, req.Command), !auto, nil
	case "read":
		p, err := addonWRMCleanPath(req.Path)
		if err != nil {
			return "", false, err
		}
		req.Path = p
		if why := aiReadDenied(p); why != "" {
			// a sensitive file is read only after the viewer says so, even for a trusted addon
			return fmt.Sprintf("Read the sensitive file %s on %s (%s); secret-looking values are masked", p, c.Name, why), true, nil
		}
		return "", false, nil
	case "write":
		p, err := addonWRMCleanPath(req.Path)
		if err != nil {
			return "", false, err
		}
		req.Path = p
		data, err := addonWRMBytes(req)
		if err != nil {
			return "", false, err
		}
		sum := sha256.Sum256(data)
		return fmt.Sprintf("Write %d bytes to %s on %s (SHA-256 %s…)", len(data), p, c.Name, hex.EncodeToString(sum[:])[:12]), !auto || aiReadDenied(p) != "", nil
	}
	return "", false, nil
}

func addonWRMBytes(req *addonWRMRequest) ([]byte, error) {
	var data []byte
	if req.Encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(req.Content)
		if err != nil {
			return nil, addonErrf("E_WRM_REQUEST", "content is not valid base64")
		}
		data = b
	} else {
		data = []byte(req.Content)
	}
	if len(data) > addonWRMFileMax {
		return nil, addonErrf("E_LIMIT_SIZE", "A file written by an addon is at most 1 MB")
	}
	return data, nil
}

// addonWRMRun does the work of an approved (or approval-free) request.
func addonWRMRun(ctx context.Context, userID int, op string, req addonWRMRequest) (map[string]interface{}, error) {
	c, err := addonWRMConnection(userID, req.Connection)
	if err != nil {
		return nil, err
	}
	if op == "metrics" {
		return addonWRMMetrics(ctx, c)
	}
	cl, err := dialSSH(c, nil)
	if err != nil {
		return nil, addonErrf("E_WRM_SSH", "SSH to %s failed: %s", c.Name, truncateStr(redactText(err.Error()), 300))
	}
	defer cl.Close()
	switch op {
	case "exec":
		cctx, cancel := context.WithTimeout(ctx, time.Duration(req.TimeoutS)*time.Second)
		defer cancel()
		sess, err := cl.NewSession()
		if err != nil {
			return nil, addonErrf("E_WRM_SSH", "SSH session failed: %s", truncateStr(err.Error(), 200))
		}
		defer sess.Close()
		buf := &lockedBuf{buf: capBuffer{max: addonWRMOutputMax}}
		sess.Stdout, sess.Stderr = buf, buf
		done := make(chan error, 1)
		go func() { done <- sess.Run(req.Command) }()
		exit := 0
		timedOut := false
		select {
		case err = <-done:
		case <-cctx.Done():
			sess.Close()
			err = cctx.Err()
			timedOut = true
		}
		if err != nil {
			exit = -1
			if ee, ok := err.(interface{ ExitStatus() int }); ok {
				exit = ee.ExitStatus()
			}
		}
		buf.mu.Lock()
		out := buf.buf.String()
		truncated := buf.buf.Len() >= addonWRMOutputMax
		buf.mu.Unlock()
		out, masked := redactSecrets(out)
		return map[string]interface{}{"output": out, "exit": exit, "timed_out": timedOut, "truncated": truncated, "masked": masked}, nil
	case "read", "write":
		fs, err := sftp.NewClient(cl)
		if err != nil {
			return nil, addonErrf("E_WRM_SSH", "SFTP failed: %s", truncateStr(err.Error(), 200))
		}
		defer fs.Close()
		p := req.Path
		if strings.HasPrefix(p, "~/") {
			if home, err := fs.Getwd(); err == nil {
				p = path.Join(home, strings.TrimPrefix(p, "~/"))
			} else {
				p = strings.TrimPrefix(p, "~/")
			}
		}
		if op == "read" {
			f, err := fs.Open(p)
			if err != nil {
				return nil, addonErrf("E_WRM_FILE", "Cannot open %s: %s", req.Path, truncateStr(err.Error(), 200))
			}
			defer f.Close()
			data, err := io.ReadAll(io.LimitReader(f, addonWRMFileMax+1))
			if err != nil {
				return nil, addonErrf("E_WRM_FILE", "Cannot read %s: %s", req.Path, truncateStr(err.Error(), 200))
			}
			truncated := len(data) > addonWRMFileMax
			if truncated {
				data = data[:addonWRMFileMax]
			}
			if !utf8.Valid(data) {
				return map[string]interface{}{"path": req.Path, "encoding": "base64", "content": base64.StdEncoding.EncodeToString(data), "size": len(data), "truncated": truncated}, nil
			}
			text, masked := redactSecrets(string(data))
			return map[string]interface{}{"path": req.Path, "content": text, "size": len(data), "truncated": truncated, "masked": masked}, nil
		}
		data, err := addonWRMBytes(&req)
		if err != nil {
			return nil, err
		}
		f, err := fs.OpenFile(p, 0x1|0x40|0x200) // O_WRONLY|O_CREATE|O_TRUNC
		if err != nil {
			return nil, addonErrf("E_WRM_FILE", "Cannot write %s: %s", req.Path, truncateStr(err.Error(), 200))
		}
		_, werr := f.Write(data)
		cerr := f.Close()
		if werr == nil {
			werr = cerr
		}
		if werr != nil {
			return nil, addonErrf("E_WRM_FILE", "Cannot write %s: %s", req.Path, truncateStr(werr.Error(), 200))
		}
		return map[string]interface{}{"path": req.Path, "size": len(data)}, nil
	}
	return nil, addonErrf("E_WRM_REQUEST", "Unknown operation")
}

func addonWRMMetrics(ctx context.Context, c Connection) (map[string]interface{}, error) {
	var m wbMetrics
	if loadMonitorConfig(c.ID).Enabled {
		s := latestMonitorSample(c.ID)
		if t, err := time.Parse(time.RFC3339, s.At); err == nil && time.Since(t) < 5*time.Minute && len(s.Metrics) > 0 {
			m = wbFromSample(s)
		}
	}
	if m.At == "" {
		if !settingBool("workboard_ssh_metrics") {
			return nil, addonErrf("E_WRM_OFF", "Reading server metrics over SSH is turned off by the administrator (workboard_ssh_metrics)")
		}
		cctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		m = wbCollect(cctx, c)
	}
	b, _ := json.Marshal(m)
	out := map[string]interface{}{}
	json.Unmarshal(b, &out)
	return out, nil
}

// addonWRMHandler serves /api/addons/{key}/instances/{iid}/wrm/{connections|exec|read|write|metrics|approve/{id}|deny/{id}}.
func addonWRMHandler(w http.ResponseWriter, r *http.Request, userID int, rec addonRecord, in addonInstance, rest []string) {
	if !rec.Enabled {
		addonFail(w, addonErrf("E_DISABLED", "The addon is disabled"), 409)
		return
	}
	if len(rest) == 0 {
		jsonError(w, "Not found", 404)
		return
	}
	if !addonServerAccessAllowed(userID) {
		addonFail(w, addonErrf("E_WRM_OFF", "Addon access to servers is turned off for you (addon_server_access)"), 403)
		return
	}
	ref := addonRef{Owner: rec.OwnerID, ID: rec.ID}
	op := rest[0]
	if op == "connections" {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		if !addonHasCap(rec, "connections.read") {
			addonFail(w, addonErrf("E_CAPABILITY", "The addon does not declare the connections.read capability"), 403)
			return
		}
		jsonOK(w, map[string]interface{}{"connections": addonWRMConnections(userID)})
		return
	}
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	if (op == "approve" || op == "deny") && len(rest) == 2 {
		p := addonWRMTake(rest[1], userID, in.ID)
		if p == nil {
			addonFail(w, addonErrf("E_WRM_APPROVAL", "The request expired or was already decided"), 410)
			return
		}
		if op == "deny" {
			auditLog(r, userID, "addon.wrm_denied", ref.Key(), map[string]interface{}{"instance": in.ID, "op": p.Op, "conn_id": p.Req.Connection})
			addonFail(w, addonErrf("E_WRM_DENIED", "The request was denied"), 403)
			return
		}
		addonWRMExecute(w, r, userID, ref, in, p.Op, p.Req, true)
		return
	}
	var req addonWRMRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, addonWRMFileMax*2+64<<10)).Decode(&req); err != nil {
		addonFail(w, addonErrf("E_WRM_REQUEST", "The request is not valid JSON"), 400)
		return
	}
	if !addonWRMRateOK(in.ID) {
		addonFail(w, addonErrf("E_LIMIT_RATE", "Too many server requests from this window"), 429)
		return
	}
	summary, ask, err := addonWRMPrepare(rec, userID, op, &req)
	if err != nil {
		status := 400
		if ae, ok := err.(*addonError); ok && ae.Code == "E_CAPABILITY" {
			status = 403
		}
		addonFail(w, err, status)
		return
	}
	if ask {
		p := &addonWRMPending{UserID: userID, Instance: in.ID, Ref: ref, Op: op, Req: req, Summary: summary}
		addonWRMPark(p)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(202)
		json.NewEncoder(w).Encode(map[string]interface{}{"approval": map[string]interface{}{
			"id": p.ID, "op": op, "summary": summary, "reason": truncateStr(req.Reason, 300), "expires_at": p.Expires.UTC().Format(time.RFC3339)}})
		return
	}
	addonWRMExecute(w, r, userID, ref, in, op, req, false)
}

func addonWRMExecute(w http.ResponseWriter, r *http.Request, userID int, ref addonRef, in addonInstance, op string, req addonWRMRequest, approved bool) {
	start := time.Now()
	res, err := addonWRMRun(r.Context(), userID, op, req)
	detail := map[string]interface{}{"instance": in.ID, "op": op, "conn_id": req.Connection, "approved": approved, "ms": time.Since(start).Milliseconds()}
	if op == "exec" {
		detail["command"] = truncateStr(redactText(req.Command), 300)
	}
	if req.Path != "" {
		detail["path"] = truncateStr(req.Path, 300)
	}
	if err != nil {
		detail["error"] = truncateStr(err.Error(), 200)
		auditLog(r, userID, "addon.wrm_"+op, ref.Key(), detail)
		addonFail(w, err, 502)
		return
	}
	if op != "read" && op != "metrics" {
		if x, ok := res["exit"]; ok {
			detail["exit"] = x
		}
		auditLog(r, userID, "addon.wrm_"+op, ref.Key(), detail)
	}
	res["ok"] = true
	jsonOK(w, res)
}

// addonWRMGrantHandler: GET/PUT /api/addons/{key}/wrm (administrators):
// whether the addon runs commands and writes without asking each time.
func addonWRMGrantHandler(w http.ResponseWriter, r *http.Request, adminID int, ref addonRef) {
	rec, ok := addonGetRef(ref)
	if !ok {
		jsonError(w, "Addon not found", 404)
		return
	}
	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
		var req struct {
			Auto bool `json:"auto"`
		}
		if json.NewDecoder(io.LimitReader(r.Body, 1024)).Decode(&req) != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		if req.Auto && !requireStepUp(w, r, "addon.wrm_auto") {
			return
		}
		addonWRMSetAuto(ref, req.Auto, usernameOf(adminID))
		auditLog(r, adminID, "addon.wrm_auto", ref.Key(), map[string]interface{}{"auto": req.Auto, "owner": ref.Owner})
	default:
		jsonError(w, "Method not allowed", 405)
		return
	}
	caps := []string{}
	for _, c := range rec.Capabilities {
		if addonWRMCaps[c] {
			caps = append(caps, c)
		}
	}
	sort.Strings(caps)
	jsonOK(w, map[string]interface{}{"auto": addonWRMAuto(ref), "capabilities": caps})
}
