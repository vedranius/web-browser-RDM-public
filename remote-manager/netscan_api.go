package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// ─── NETWORK SCANNER: API, history, CLI ──────────────
//
// GET    /api/netscan                     scans of the user (newest first) and the policy
// POST   /api/netscan                     {source_id, target, profile, ports} → starts a scan
// GET    /api/netscan/interfaces?source_id=N   network cards of the WRM server or of an SSH server
// GET    /api/netscan/{id}                progress and devices (new/gone compared with the last scan of the same network)
// POST   /api/netscan/{id}/cancel
// DELETE /api/netscan/{id}
// POST   /api/netscan/import              a JSON file written by "wrm -netscan"
// PUT    /api/netscan/labels              {key, name, note} own name for a device (by MAC, else IP)
// GET    /api/netscan/vendor?mac=…        vendor of a MAC address
// POST   /api/netscan/wake                {mac, broadcast, source_id} Wake-on-LAN
// POST   /api/netscan/open                {source_id, host, port, scheme} open a found web interface (through a tunnel for SSH sources)

const netscanKeep = 30 // finished scans kept per user

func initNetScanSchema() {
	ensureTable("net_scans", `CREATE TABLE IF NOT EXISTS net_scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		source_id INTEGER NOT NULL DEFAULT 0,
		source_name TEXT NOT NULL DEFAULT '',
		target TEXT NOT NULL DEFAULT '',
		profile TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL DEFAULT 'running',
		error TEXT NOT NULL DEFAULT '',
		devices INTEGER NOT NULL DEFAULT 0,
		result TEXT NOT NULL DEFAULT '',
		started_at TEXT NOT NULL DEFAULT '',
		finished_at TEXT NOT NULL DEFAULT ''
	)`, "id", "user_id", "source_id", "source_name", "target", "profile", "status", "error", "devices", "result", "started_at", "finished_at")
	mustExec(`CREATE INDEX IF NOT EXISTS idx_net_scans_user ON net_scans(user_id, id)`)
	ensureTable("net_device_labels", `CREATE TABLE IF NOT EXISTS net_device_labels (
		user_id INTEGER NOT NULL,
		dkey TEXT NOT NULL,
		name TEXT NOT NULL DEFAULT '',
		note TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (user_id, dkey)
	)`, "user_id", "dkey", "name", "note", "updated_at")
	// Scans that were running when WRM stopped cannot continue.
	db.Exec(`UPDATE net_scans SET status='interrupted', error='WRM was restarted during the scan' WHERE status='running'`)
}

func netscanAllowed(userID int) bool {
	return policyWhoAllows(userID, "network_scan")
}

type netscanJob struct {
	id       int64
	userID   int
	sourceID int
	source   string
	target   string
	profile  string
	started  time.Time
	sc       *scanner
	cancel   context.CancelFunc
	finished chan struct{}
}

var netscanJobs = struct {
	sync.Mutex
	m map[int64]*netscanJob
}{m: map[int64]*netscanJob{}}

type netscanResult struct {
	Devices []scanDevice `json:"devices"`
	Gone    []scanDevice `json:"gone,omitempty"`
	Log     []string     `json:"log,omitempty"`
	Prev    int64        `json:"previous_id,omitempty"`
}

type netscanView struct {
	ID       int64        `json:"id"`
	SourceID int          `json:"source_id"`
	Source   string       `json:"source"`
	Target   string       `json:"target"`
	Profile  string       `json:"profile"`
	Status   string       `json:"status"`
	Error    string       `json:"error,omitempty"`
	Count    int          `json:"devices_count"`
	Started  string       `json:"started_at"`
	Finished string       `json:"finished_at,omitempty"`
	Progress scanProgress `json:"progress"`
	Devices  []scanDevice `json:"devices,omitempty"`
	Gone     []scanDevice `json:"gone,omitempty"`
	Log      []string     `json:"log,omitempty"`
	NewCount int          `json:"new_count"`
	Prev     int64        `json:"previous_id,omitempty"`
}

func deviceKey(d scanDevice) string {
	if d.MAC != "" && !d.RandomMAC {
		return d.MAC
	}
	if d.MAC != "" {
		return d.MAC + "@" + d.IP
	}
	return d.IP
}

// compareWithPrevious marks devices not in prev as new and returns those of prev that are gone.
func compareWithPrevious(cur []scanDevice, prev []scanDevice) []scanDevice {
	seen := map[string]bool{}
	for _, d := range prev {
		seen[deviceKey(d)] = true
		seen["ip:"+d.IP] = true
	}
	now := map[string]bool{}
	for i := range cur {
		now[deviceKey(cur[i])] = true
		now["ip:"+cur[i].IP] = true
		if !seen[deviceKey(cur[i])] {
			cur[i].New = true
		}
	}
	gone := []scanDevice{}
	for _, d := range prev {
		if !now[deviceKey(d)] {
			d.New = false
			gone = append(gone, d)
		}
	}
	return gone
}

func loadScanLabels(userID int) map[string][2]string {
	out := map[string][2]string{}
	rows, err := db.Query(`SELECT dkey, name, note FROM net_device_labels WHERE user_id=?`, userID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var k, n, note string
		if rows.Scan(&k, &n, &note) == nil {
			out[k] = [2]string{n, note}
		}
	}
	return out
}

func applyScanLabels(userID int, list []scanDevice) {
	labels := loadScanLabels(userID)
	for i := range list {
		l, ok := labels[deviceKey(list[i])]
		if !ok {
			l, ok = labels[list[i].IP]
		}
		if ok {
			list[i].Label, list[i].Note = l[0], l[1]
		}
	}
}

// previousScan returns the newest finished scan of the same network and source before id.
func previousScan(userID, sourceID int, target string, before int64) (int64, []scanDevice) {
	var id int64
	var raw string
	err := db.QueryRow(`SELECT id, result FROM net_scans WHERE user_id=? AND source_id=? AND target=? AND status='done' AND id<? ORDER BY id DESC LIMIT 1`,
		userID, sourceID, target, before).Scan(&id, &raw)
	if err != nil {
		return 0, nil
	}
	var res netscanResult
	json.Unmarshal([]byte(raw), &res)
	return id, res.Devices
}

func (j *netscanJob) view() netscanView {
	devs := j.sc.results()
	v := netscanView{ID: j.id, SourceID: j.sourceID, Source: j.source, Target: j.target, Profile: j.profile, Status: "running",
		Started: j.started.UTC().Format(time.RFC3339), Progress: j.sc.progress(), Devices: devs, Log: j.sc.logLines(), Count: len(devs)}
	return v
}

// runNetscanJob scans, stores the result and removes the job.
func runNetscanJob(j *netscanJob, cl *ssh.Client, r *http.Request) {
	defer close(j.finished)
	if cl != nil {
		defer cl.Close()
	}
	j.sc.run()
	status, errText := "done", ""
	if j.sc.ctx.Err() != nil {
		status, errText = "cancelled", "stopped by the user"
	}
	devs := j.sc.results()
	res := netscanResult{Devices: devs, Log: j.sc.logLines()}
	if status == "done" {
		var prev []scanDevice
		res.Prev, prev = previousScan(j.userID, j.sourceID, j.target, j.id)
		if res.Prev > 0 {
			res.Gone = compareWithPrevious(res.Devices, prev)
		}
	}
	raw, _ := json.Marshal(res)
	db.Exec(`UPDATE net_scans SET status=?, error=?, devices=?, result=?, finished_at=? WHERE id=?`,
		status, errText, len(devs), string(raw), time.Now().UTC().Format(time.RFC3339), j.id)
	db.Exec(`DELETE FROM net_scans WHERE user_id=? AND status!='running' AND id NOT IN (SELECT id FROM net_scans WHERE user_id=? ORDER BY id DESC LIMIT ?)`,
		j.userID, j.userID, netscanKeep)
	netscanJobs.Lock()
	delete(netscanJobs.m, j.id)
	netscanJobs.Unlock()
	auditLog(nil, j.userID, "netscan.finished", j.target, map[string]interface{}{"source": j.source, "devices": len(devs), "status": status})
	if newN := countNew(res.Devices); newN > 0 && res.Prev > 0 {
		recordUserEvent(j.userID, "netscan", fmt.Sprintf("Network scan %s: %d new device(s)", j.target, newN), "Source: "+j.source)
	}
}

func countNew(list []scanDevice) int {
	n := 0
	for _, d := range list {
		if d.New {
			n++
		}
	}
	return n
}

func loadScanView(id int64, userID int) (netscanView, error) {
	var v netscanView
	var raw string
	err := db.QueryRow(`SELECT id, source_id, source_name, target, profile, status, error, devices, result, started_at, finished_at FROM net_scans WHERE id=? AND user_id=?`,
		id, userID).Scan(&v.ID, &v.SourceID, &v.Source, &v.Target, &v.Profile, &v.Status, &v.Error, &v.Count, &raw, &v.Started, &v.Finished)
	if err != nil {
		return v, err
	}
	var res netscanResult
	json.Unmarshal([]byte(raw), &res)
	v.Devices, v.Gone, v.Log, v.Prev = res.Devices, res.Gone, res.Log, res.Prev
	if v.Devices == nil {
		v.Devices = []scanDevice{}
	}
	v.Progress = scanProgress{Phase: "done"}
	return v, nil
}

// netscanSource opens the SSH connection of source_id (0 = the WRM server).
func netscanSource(userID, sourceID int) (*ssh.Client, string, int, error) {
	if sourceID <= 0 {
		return nil, "WRM server", 404, nil
	}
	if !userOwnsConnection(sourceID, userID) {
		return nil, "", 404, fmt.Errorf("Connection not found")
	}
	c, err := loadConnection(sourceID)
	if err != nil || !isSSHProtocol(c) {
		return nil, "", 400, fmt.Errorf("The source must be an SSH connection")
	}
	cl, err := dialSSH(c, nil)
	if err != nil {
		return nil, "", 502, fmt.Errorf("Cannot log in to %s: %v", c.Name, err)
	}
	return cl, c.Name, 200, nil
}

func apiNetScanHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !netscanAllowed(userID) {
		jsonError(w, "The network scanner is turned off for this account (policy network_scan)", 403)
		return
	}
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/netscan"), "/")
	switch {
	case rest == "" && r.Method == http.MethodGet:
		netscanList(w, userID)
	case rest == "" && r.Method == http.MethodPost:
		netscanStart(w, r, userID)
	case rest == "interfaces" && r.Method == http.MethodGet:
		src, _ := strconv.Atoi(r.URL.Query().Get("source_id"))
		if src <= 0 {
			jsonOK(w, map[string]interface{}{"source": "WRM server", "interfaces": localInterfaces()})
			return
		}
		cl, name, code, err := netscanSource(userID, src)
		if err != nil {
			jsonError(w, err.Error(), code)
			return
		}
		defer cl.Close()
		ifs, err := remoteInterfaces(cl)
		if err != nil {
			jsonError(w, err.Error(), 502)
			return
		}
		jsonOK(w, map[string]interface{}{"source": name, "interfaces": ifs})
	case rest == "import" && r.Method == http.MethodPost:
		netscanImport(w, r, userID)
	case rest == "labels" && r.Method == http.MethodPut:
		var in struct{ Key, Name, Note string }
		if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil || strings.TrimSpace(in.Key) == "" || len(in.Key) > 80 {
			jsonError(w, "Bad request", 400)
			return
		}
		in.Name, in.Note = truncateStr(strings.TrimSpace(in.Name), 80), truncateStr(strings.TrimSpace(in.Note), 500)
		if in.Name == "" && in.Note == "" {
			db.Exec(`DELETE FROM net_device_labels WHERE user_id=? AND dkey=?`, userID, in.Key)
		} else {
			db.Exec(`INSERT INTO net_device_labels (user_id, dkey, name, note, updated_at) VALUES (?,?,?,?,?)
				ON CONFLICT(user_id, dkey) DO UPDATE SET name=excluded.name, note=excluded.note, updated_at=excluded.updated_at`,
				userID, in.Key, in.Name, in.Note, time.Now().UTC().Format(time.RFC3339))
		}
		jsonOK(w, map[string]bool{"ok": true})
	case rest == "vendor" && r.Method == http.MethodGet:
		mac := r.URL.Query().Get("mac")
		hw, err := parseWakeMAC(mac)
		if err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		jsonOK(w, map[string]interface{}{"mac": strings.ToUpper(hw.String()), "vendor": ouiVendor(hw.String()), "random": macIsRandom(hw.String())})
	case rest == "wake" && r.Method == http.MethodPost:
		netscanWake(w, r, userID)
	case rest == "open" && r.Method == http.MethodPost:
		netscanOpen(w, r, userID)
	default:
		parts := strings.Split(rest, "/")
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			jsonError(w, "Not found", 404)
			return
		}
		netscanByID(w, r, userID, id, parts[1:])
	}
}

func netscanList(w http.ResponseWriter, userID int) {
	out := []netscanView{}
	rows, err := db.Query(`SELECT id, source_id, source_name, target, profile, status, error, devices, started_at, finished_at, result FROM net_scans WHERE user_id=? ORDER BY id DESC LIMIT ?`, userID, netscanKeep)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var v netscanView
			var raw string
			if rows.Scan(&v.ID, &v.SourceID, &v.Source, &v.Target, &v.Profile, &v.Status, &v.Error, &v.Count, &v.Started, &v.Finished, &raw) == nil {
				var res struct {
					Devices []struct {
						New bool `json:"new"`
					} `json:"devices"`
					Gone []json.RawMessage `json:"gone"`
				}
				json.Unmarshal([]byte(raw), &res)
				for _, d := range res.Devices {
					if d.New {
						v.NewCount++
					}
				}
				v.Gone = make([]scanDevice, len(res.Gone)) // only the count is shown
				out = append(out, v)
			}
		}
	}
	jsonOK(w, map[string]interface{}{"scans": out, "max_hosts": netscanMaxHosts, "public": settingBool("network_scan_public"),
		"tunnels": tunnelsAllowedFor(userID) == nil})
}

func netscanStart(w http.ResponseWriter, r *http.Request, userID int) {
	var in struct {
		SourceID int    `json:"source_id"`
		Target   string `json:"target"`
		Profile  string `json:"profile"`
		Ports    string `json:"ports"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	addrs, target, err := parseScanTarget(in.Target)
	if err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	if !settingBool("network_scan_public") {
		for _, a := range addrs {
			if !scanAddrPrivate(a) {
				jsonError(w, "Only private networks (10.x, 172.16–31.x, 192.168.x, 100.64.x, link-local) can be scanned; an administrator can allow public ranges (setting network_scan_public)", 403)
				return
			}
		}
	}
	ports, err := scanPortsFor(in.Profile, in.Ports)
	if err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	if in.Profile == "" {
		in.Profile = "quick"
	}
	netscanJobs.Lock()
	running := 0
	for _, j := range netscanJobs.m {
		if j.userID == userID {
			netscanJobs.Unlock()
			jsonError(w, "Another scan of yours is still running", 429)
			return
		}
		running++
	}
	netscanJobs.Unlock()
	if running >= 4 {
		jsonError(w, "Four scans are running on this WRM already; try again in a minute", 429)
		return
	}
	cl, source, code, err := netscanSource(userID, in.SourceID)
	if err != nil {
		jsonError(w, err.Error(), code)
		return
	}
	now := time.Now()
	res, err := db.Exec(`INSERT INTO net_scans (user_id, source_id, source_name, target, profile, status, started_at) VALUES (?,?,?,?,?,'running',?)`,
		userID, max(in.SourceID, 0), source, target, in.Profile, now.UTC().Format(time.RFC3339))
	if err != nil {
		if cl != nil {
			cl.Close()
		}
		jsonError(w, err.Error(), 500)
		return
	}
	id, _ := res.LastInsertId()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	j := &netscanJob{id: id, userID: userID, sourceID: max(in.SourceID, 0), source: source, target: target, profile: in.Profile, started: now,
		sc: newScanner(ctx, cl, addrs, ports), cancel: cancel, finished: make(chan struct{})}
	netscanJobs.Lock()
	netscanJobs.m[id] = j
	netscanJobs.Unlock()
	auditLog(r, userID, "netscan.start", target, map[string]interface{}{"source": source, "addresses": len(addrs), "ports": len(ports)})
	go func() {
		defer cancel()
		runNetscanJob(j, cl, r)
	}()
	jsonOK(w, map[string]interface{}{"id": id, "target": target, "addresses": len(addrs)})
}

func netscanByID(w http.ResponseWriter, r *http.Request, userID int, id int64, action []string) {
	netscanJobs.Lock()
	j := netscanJobs.m[id]
	netscanJobs.Unlock()
	if j != nil && j.userID != userID {
		j = nil
	}
	switch {
	case len(action) == 0 && r.Method == http.MethodGet:
		if j != nil {
			v := j.view()
			applyScanLabels(userID, v.Devices)
			jsonOK(w, v)
			return
		}
		v, err := loadScanView(id, userID)
		if err != nil {
			jsonError(w, "Not found", 404)
			return
		}
		applyScanLabels(userID, v.Devices)
		applyScanLabels(userID, v.Gone)
		v.NewCount = countNew(v.Devices)
		jsonOK(w, v)
	case len(action) == 1 && action[0] == "cancel" && r.Method == http.MethodPost:
		if j != nil {
			j.cancel()
			select {
			case <-j.finished:
			case <-time.After(10 * time.Second):
			}
		}
		jsonOK(w, map[string]bool{"ok": true})
	case len(action) == 0 && r.Method == http.MethodDelete:
		if j != nil {
			jsonError(w, "Stop the scan first", 409)
			return
		}
		db.Exec(`DELETE FROM net_scans WHERE id=? AND user_id=?`, id, userID)
		jsonOK(w, map[string]bool{"ok": true})
	default:
		jsonError(w, "Not found", 404)
	}
}

// netscanFile is what "wrm -netscan" writes and /api/netscan/import reads.
type netscanFile struct {
	Format   string       `json:"format"` // wrm-netscan/1
	Version  string       `json:"wrm_version"`
	Host     string       `json:"host"`
	Target   string       `json:"target"`
	Profile  string       `json:"profile"`
	Started  string       `json:"started_at"`
	Finished string       `json:"finished_at"`
	Devices  []scanDevice `json:"devices"`
	Log      []string     `json:"log,omitempty"`
}

func netscanImport(w http.ResponseWriter, r *http.Request, userID int) {
	var f netscanFile
	if err := json.NewDecoder(io.LimitReader(r.Body, 8<<20)).Decode(&f); err != nil || f.Format != "wrm-netscan/1" {
		jsonError(w, "This is not a file written by wrm -netscan", 400)
		return
	}
	if _, target, err := parseScanTarget(f.Target); err != nil {
		jsonError(w, "Invalid network in the file: "+err.Error(), 400)
		return
	} else {
		f.Target = target
	}
	if len(f.Devices) > netscanMaxHosts {
		jsonError(w, "Too many devices in the file", 400)
		return
	}
	clean := make([]scanDevice, 0, len(f.Devices))
	for _, d := range f.Devices {
		if a, err := netip.ParseAddr(d.IP); err != nil || !a.Is4() {
			continue
		}
		d.New, d.Label, d.Note = false, "", ""
		if len(d.Ports) > 200 {
			d.Ports = d.Ports[:200]
		}
		clean = append(clean, sanitizeScanDevice(d))
	}
	host := truncateStr(strings.TrimSpace(f.Host), 60)
	if host == "" {
		host = "?"
	}
	source := "wrm -netscan @ " + host
	started := f.Started
	if _, err := time.Parse(time.RFC3339, started); err != nil {
		started = time.Now().UTC().Format(time.RFC3339)
	}
	res := netscanResult{Devices: clean}
	if len(f.Log) > 200 {
		f.Log = f.Log[:200]
	}
	for _, l := range f.Log {
		res.Log = append(res.Log, truncateStr(l, 300))
	}
	ins, err := db.Exec(`INSERT INTO net_scans (user_id, source_id, source_name, target, profile, status, started_at) VALUES (?,-1,?,?,?,'running',?)`,
		userID, source, f.Target, truncateStr(f.Profile, 20), started)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	id, _ := ins.LastInsertId()
	var prev []scanDevice
	if res.Prev, prev = previousScan(userID, -1, f.Target, id); res.Prev > 0 {
		res.Gone = compareWithPrevious(res.Devices, prev)
	}
	raw, _ := json.Marshal(res)
	db.Exec(`UPDATE net_scans SET status='done', devices=?, result=?, finished_at=? WHERE id=?`, len(clean), string(raw), time.Now().UTC().Format(time.RFC3339), id)
	auditLog(r, userID, "netscan.import", f.Target, map[string]interface{}{"source": source, "devices": len(clean)})
	jsonOK(w, map[string]interface{}{"id": id, "devices": len(clean)})
}

// sanitizeScanDevice shortens every text of an imported device.
func sanitizeScanDevice(d scanDevice) scanDevice {
	c := func(s string, n int) string { return truncateStr(strings.TrimSpace(s), n) }
	d.MAC, d.Name, d.DNSName, d.NetBIOS, d.Workgroup, d.MDNSName = c(d.MAC, 17), c(d.Name, 120), c(d.DNSName, 120), c(d.NetBIOS, 20), c(d.Workgroup, 20), c(d.MDNSName, 120)
	d.UPnP, d.Model, d.HTTPTitle, d.HTTPServer, d.TLSName, d.SSH = c(d.UPnP, 80), c(d.Model, 120), c(d.HTTPTitle, 100), c(d.HTTPServer, 80), c(d.TLSName, 100), c(d.SSH, 100)
	d.Kind, d.OS = c(d.Kind, 20), c(d.OS, 40)
	if d.MAC != "" {
		if hw, err := net.ParseMAC(d.MAC); err == nil && len(hw) == 6 {
			d.MAC = strings.ToUpper(hw.String())
			d.Vendor = ouiVendor(d.MAC)
			d.RandomMAC = macIsRandom(d.MAC)
		} else {
			d.MAC, d.Vendor = "", ""
		}
	}
	for i := range d.Ports {
		d.Ports[i].Service, d.Ports[i].Banner = c(d.Ports[i].Service, 40), c(d.Ports[i].Banner, 100)
	}
	if len(d.Seen) > 10 {
		d.Seen = d.Seen[:10]
	}
	for i := range d.Seen {
		d.Seen[i] = c(d.Seen[i], 10)
	}
	return d
}

func netscanWake(w http.ResponseWriter, r *http.Request, userID int) {
	var in struct {
		MAC       string `json:"mac"`
		Broadcast string `json:"broadcast"`
		SourceID  int    `json:"source_id"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	hw, err := parseWakeMAC(in.MAC)
	if err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	in.Broadcast = strings.TrimSpace(in.Broadcast)
	if in.Broadcast != "" {
		a, err := netip.ParseAddr(in.Broadcast)
		if err != nil || !a.Is4() || (!scanAddrPrivate(a) && in.Broadcast != "255.255.255.255" && !settingBool("network_scan_public")) {
			jsonError(w, "Enter a broadcast address such as 192.168.1.255", 400)
			return
		}
	}
	cl, source, code, err := netscanSource(userID, in.SourceID)
	if err != nil {
		jsonError(w, err.Error(), code)
		return
	}
	out := ""
	if cl == nil {
		err = sendWake(hw, in.Broadcast)
		out = "sent from the WRM server"
	} else {
		defer cl.Close()
		out, err = sendWakeRemote(cl, hw, in.Broadcast)
	}
	auditLog(r, userID, "netscan.wake", strings.ToUpper(hw.String()), map[string]interface{}{"source": source, "broadcast": in.Broadcast})
	if err != nil {
		jsonError(w, "Wake-on-LAN from "+source+": "+err.Error(), 502)
		return
	}
	jsonOK(w, map[string]interface{}{"ok": true, "source": source, "output": strings.TrimSpace(out)})
}

// netscanOpen opens a web interface found by a scan: directly for scans from the WRM
// server, through a temporary tunnel over the SSH connection otherwise.
func netscanOpen(w http.ResponseWriter, r *http.Request, userID int) {
	var in struct {
		SourceID int    `json:"source_id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Scheme   string `json:"scheme"`
	}
	if json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	a, err := netip.ParseAddr(in.Host)
	if err != nil || !a.Is4() || in.Port < 1 || in.Port > 65535 || (in.Scheme != "http" && in.Scheme != "https") {
		jsonError(w, "Bad request", 400)
		return
	}
	direct := in.Scheme + "://" + net.JoinHostPort(in.Host, strconv.Itoa(in.Port)) + "/"
	if in.SourceID <= 0 {
		jsonOK(w, map[string]interface{}{"url": direct, "direct": true})
		return
	}
	if err := tunnelsAllowedFor(userID); err != nil {
		jsonError(w, err.Error(), 403)
		return
	}
	if !userOwnsConnection(in.SourceID, userID) {
		jsonError(w, "Connection not found", 404)
		return
	}
	c, err := loadConnection(in.SourceID)
	if err != nil || !isSSHProtocol(c) {
		jsonError(w, "The source must be an SSH connection", 400)
		return
	}
	bind, allow := browserTunnelBind(r, userID)
	key := fmt.Sprintf("s%du%d-%s-%d-%s", c.ID, userID, in.Host, in.Port, allow)
	def := tunnelDef{ConnID: c.ID, UserID: userID, Name: in.Host + ":" + strconv.Itoa(in.Port), Kind: "local", BindHost: bind,
		TargetHost: in.Host, TargetPort: in.Port, OpenScheme: in.Scheme, OpenPath: "/", StartMode: "manual", AllowIP: allow}
	t, err := tunnelMgr.start(key, def, c, userID, "scan", r)
	if err != nil {
		jsonError(w, err.Error(), 502)
		return
	}
	_, lport, _ := net.SplitHostPort(t.view().Listen)
	auditLogRef(r, userID, usernameOf(userID), "netscan.open", c.Name, map[string]string{"target": net.JoinHostPort(in.Host, strconv.Itoa(in.Port)), "client_only": allow}, auditRef{ConnID: c.ID})
	jsonOK(w, map[string]interface{}{"url": in.Scheme + "://" + browserTunnelHost(r, allow) + ":" + lport + "/", "direct": false, "route": c.Name,
		"listen": t.view().Listen, "key": key, "local_only": allow == "", "client_only": allow})
}

// ── CLI: wrm -netscan RANGE ──

// runNetscanCLI scans from this computer and writes the JSON file for an import.
func runNetscanCLI(target, profile, outFile string, stdout, stderr io.Writer) int {
	addrs, norm, err := parseScanTarget(target)
	if err != nil {
		fmt.Fprintln(stderr, "wrm -netscan:", err)
		return 2
	}
	ports, err := scanPortsFor(cliPortProfile(profile))
	if err != nil {
		fmt.Fprintln(stderr, "wrm -netscan:", err)
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	sc := newScanner(ctx, nil, addrs, ports)
	started := time.Now()
	done := make(chan struct{})
	go func() { sc.run(); close(done) }()
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	last := ""
loop:
	for {
		select {
		case <-done:
			break loop
		case <-tick.C:
			p := sc.progress()
			if line := fmt.Sprintf("%s %d/%d", p.Phase, p.Done, p.Total); line != last {
				fmt.Fprintf(stderr, "\r%-40s", line)
				last = line
			}
		}
	}
	devs := sc.results()
	host, _ := os.Hostname()
	f := netscanFile{Format: "wrm-netscan/1", Version: AppVersion, Host: host, Target: norm, Profile: cliProfileName(profile),
		Started: started.UTC().Format(time.RFC3339), Finished: time.Now().UTC().Format(time.RFC3339), Devices: devs, Log: sc.logLines()}
	b, _ := json.MarshalIndent(f, "", "  ")
	fmt.Fprintf(stderr, "\r%-40s\n", fmt.Sprintf("%d devices found in %s", len(devs), time.Since(started).Round(time.Second)))
	for _, d := range devs {
		fmt.Fprintf(stderr, "  %-15s %-17s %-24s %s\n", d.IP, d.MAC, truncateStr(d.Vendor, 24), d.Name)
	}
	if outFile == "" || outFile == "-" {
		stdout.Write(append(b, '\n'))
		return 0
	}
	if err := os.WriteFile(outFile, b, 0o600); err != nil {
		fmt.Fprintln(stderr, "wrm -netscan:", err)
		return 1
	}
	fmt.Fprintf(stderr, "Saved to %s. Import it in WRM: Tools → Network scanner → Import.\n", outFile)
	return 0
}

func cliPortProfile(p string) (string, string) {
	switch p {
	case "", "quick", "extended", "none":
		return p, ""
	}
	return "custom", p
}

func cliProfileName(p string) string {
	prof, _ := cliPortProfile(p)
	if prof == "" {
		return "quick"
	}
	return prof
}
