package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ─── MONITORING RULES, EVENTS AND HISTORY (v14.9.0) ──
//
// Every poll ends in one state: healthy, warning, problem (a critical limit is
// exceeded or a check fails) or failed (the server could not be read). The
// reasons say why, in a structured form the UI translates. CPU and load must
// stay above their limit for a number of polls in a row before they count,
// so one busy second does not page anybody. A change of state is an event
// (kept per server) and a notification; hourly roll-ups keep 90 days of
// history next to the last 2880 raw samples.

type monitorThresholds struct {
	CPUWarn   float64 `json:"cpu_warn"`
	CPUCrit   float64 `json:"cpu_crit"`
	MemWarn   float64 `json:"mem_warn"`
	MemCrit   float64 `json:"mem_crit"`
	SwapWarn  float64 `json:"swap_warn"`
	SwapCrit  float64 `json:"swap_crit"`
	DiskWarn  float64 `json:"disk_warn"`
	DiskCrit  float64 `json:"disk_crit"`
	InodeWarn float64 `json:"inode_warn"`
	InodeCrit float64 `json:"inode_crit"`
	LoadWarn  float64 `json:"load_warn"` // per CPU core
	LoadCrit  float64 `json:"load_crit"`
	Sustain   int     `json:"sustain"`   // polls in a row for CPU and load
	CertDays  float64 `json:"cert_days"` // certificate warning window
	Reboot    bool    `json:"reboot"`    // a pending reboot is a warning
}

func defaultMonitorThresholds() monitorThresholds {
	return monitorThresholds{CPUWarn: 85, CPUCrit: 95, MemWarn: 85, MemCrit: 95, SwapWarn: 50, SwapCrit: 80, DiskWarn: 85, DiskCrit: 95,
		InodeWarn: 85, InodeCrit: 95, LoadWarn: 1.5, LoadCrit: 3, Sustain: 2, CertDays: 14, Reboot: true}
}

// normalize fills missing limits with the defaults (a configuration saved
// before v14.9.0 has none) and refuses nonsense. A limit above 100 % turns
// that rule off.
func (t *monitorThresholds) normalize() error {
	d := defaultMonitorThresholds()
	if *t == (monitorThresholds{}) {
		*t = d
		return nil
	}
	fill := func(v *float64, def, hi float64) error {
		if *v == 0 {
			*v = def
		}
		if *v < 0 || *v > hi || math.IsNaN(*v) {
			return fmt.Errorf("alert limits must be between 0 and %g", hi)
		}
		return nil
	}
	for _, x := range []struct {
		v      *float64
		def, h float64
	}{{&t.CPUWarn, d.CPUWarn, 1000}, {&t.CPUCrit, d.CPUCrit, 1000}, {&t.MemWarn, d.MemWarn, 1000}, {&t.MemCrit, d.MemCrit, 1000},
		{&t.SwapWarn, d.SwapWarn, 1000}, {&t.SwapCrit, d.SwapCrit, 1000}, {&t.DiskWarn, d.DiskWarn, 1000}, {&t.DiskCrit, d.DiskCrit, 1000},
		{&t.InodeWarn, d.InodeWarn, 1000}, {&t.InodeCrit, d.InodeCrit, 1000}, {&t.LoadWarn, d.LoadWarn, 1000}, {&t.LoadCrit, d.LoadCrit, 1000},
		{&t.CertDays, d.CertDays, 365}} {
		if err := fill(x.v, x.def, x.h); err != nil {
			return err
		}
	}
	if t.Sustain == 0 {
		t.Sustain = d.Sustain
	}
	if t.Sustain < 1 || t.Sustain > 30 {
		return fmt.Errorf("polls in a row must be 1-30")
	}
	return nil
}

type monitorReason struct {
	Severity string  `json:"severity"` // warning | problem
	Code     string  `json:"code"`     // cpu, memory, swap, disk, inode, load, check, log, failed_units, reboot, time_sync, stale, missing, unreachable, checks_failed
	Subject  string  `json:"subject,omitempty"`
	Value    float64 `json:"value,omitempty"`
	Limit    float64 `json:"limit,omitempty"`
}

func (s *monitorSample) addReason(severity, code, subject string, value, limit float64) {
	if len(s.Reasons) >= 40 {
		return
	}
	s.Reasons = append(s.Reasons, monitorReason{Severity: severity, Code: code, Subject: truncateStr(subject, 200), Value: math.Round(value*10) / 10, Limit: limit})
}

// reasonText is the English (or Croatian) line used in notifications and
// events; the UI renders reasons itself.
func (x monitorReason) text(lang string) string {
	hr := lang == "hr"
	pick := func(en, h string) string {
		if hr {
			return h
		}
		return en
	}
	v := strconv.FormatFloat(x.Value, 'f', -1, 64)
	l := strconv.FormatFloat(x.Limit, 'f', -1, 64)
	switch x.Code {
	case "cpu":
		return pick("CPU "+v+" % (limit "+l+" %)", "CPU "+v+" % (granica "+l+" %)")
	case "memory":
		return pick("memory "+v+" % (limit "+l+" %)", "memorija "+v+" % (granica "+l+" %)")
	case "swap":
		return pick("swap "+v+" % (limit "+l+" %)", "swap "+v+" % (granica "+l+" %)")
	case "disk":
		return pick("disk "+x.Subject+" "+v+" % full (limit "+l+" %)", "disk "+x.Subject+" "+v+" % pun (granica "+l+" %)")
	case "inode":
		return pick("inodes on "+x.Subject+" "+v+" % used", "inodovi na "+x.Subject+" "+v+" % iskorišteno")
	case "load":
		return pick("load "+v+" per core (limit "+l+")", "opterećenje "+v+" po jezgri (granica "+l+")")
	case "check":
		return pick("check failed: "+x.Subject, "provjera ne prolazi: "+x.Subject)
	case "log":
		return pick("log alert: "+x.Subject, "log upozorenje: "+x.Subject)
	case "failed_units":
		return pick("failed systemd units: "+x.Subject, "palih systemd unitova: "+x.Subject)
	case "reboot":
		return pick("a reboot is pending", "čeka se ponovno pokretanje")
	case "time_sync":
		return pick("the clock is not synchronised", "sat nije sinkroniziran")
	case "stale":
		return pick("values are older than 5 minutes", "vrijednosti su starije od 5 minuta")
	case "missing":
		return pick("some values could not be read: "+x.Subject, "neke vrijednosti nisu pročitane: "+x.Subject)
	case "checks_failed":
		return pick("checks could not run", "provjere se nisu mogle pokrenuti")
	case "unreachable":
		return pick("the server could not be read: "+x.Subject, "poslužitelj nije dostupan: "+x.Subject)
	}
	return x.Code + " " + x.Subject
}

func monitorStateRank(s string) int {
	switch s {
	case "healthy":
		return 0
	case "warning":
		return 1
	case "problem":
		return 2
	case "failed":
		return 3
	}
	return -1
}

// applyMonitorRules compares the sample with the limits and sets its state.
// It runs after finishMonitor for every transport.
func applyMonitorRules(s *monitorSample, old monitorSample, t monitorThresholds) {
	if s.State == "failed" {
		if len(s.Reasons) == 0 {
			s.addReason("problem", "unreachable", truncateStr(s.Error, 200), 0, 0)
		}
		return
	}
	if s.Breach == nil {
		s.Breach = map[string]int{}
	}
	num := func(k string) (float64, bool) {
		m, ok := s.Metrics[k]
		if !ok || m.Value == nil || m.Error != "" {
			return 0, false
		}
		return *m.Value, true
	}
	// sustained: CPU and load count only after Sustain polls in a row
	sustained := func(key string, v, warn, crit float64, code string) {
		if v < warn {
			return
		}
		n := old.Breach[key] + 1
		s.Breach[key] = n
		if n < t.Sustain {
			return
		}
		if v >= crit {
			s.addReason("problem", code, "", v, crit)
		} else {
			s.addReason("warning", code, "", v, warn)
		}
	}
	simple := func(code, subject string, v, warn, crit float64) {
		if v >= crit {
			s.addReason("problem", code, subject, v, crit)
		} else if v >= warn {
			s.addReason("warning", code, subject, v, warn)
		}
	}
	if v, ok := num("cpu_percent"); ok {
		sustained("cpu", v, t.CPUWarn, t.CPUCrit, "cpu")
	}
	if v, ok := num("load"); ok {
		cores, _ := num("cpu_count")
		if cores < 1 {
			cores = 1
		}
		sustained("load", v/cores, t.LoadWarn, t.LoadCrit, "load")
	}
	if v, ok := num("memory_percent"); ok {
		simple("memory", "", v, t.MemWarn, t.MemCrit)
	}
	if v, ok := num("swap_percent"); ok {
		if total, ok := num("swap_total"); !ok || total > 0 {
			simple("swap", "", v, t.SwapWarn, t.SwapCrit)
		}
	}
	keys := make([]string, 0, len(s.Metrics))
	for k := range s.Metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v, ok := num(k)
		if !ok {
			continue
		}
		if m, found := strings.CutPrefix(k, "disk:"); found {
			simple("disk", m, v, t.DiskWarn, t.DiskCrit)
		} else if m, found := strings.CutPrefix(k, "inode:"); found {
			simple("inode", m, v, t.InodeWarn, t.InodeCrit)
		}
	}
	if v, ok := num("failed_units"); ok && v > 0 {
		s.addReason("warning", "failed_units", truncateStr(s.Metrics["failed_units"].Text, 200), v, 0)
	}
	if v, ok := num("reboot_required"); ok && v > 0 && t.Reboot {
		s.addReason("warning", "reboot", "", 0, 0)
	}
	if v, ok := num("time_sync"); ok && v == 0 {
		s.addReason("warning", "time_sync", "", 0, 0)
	}
	state := "healthy"
	if s.State == "warning" {
		state = "warning"
	}
	for _, r := range s.Reasons {
		if r.Severity == "problem" {
			state = "problem"
		} else if state == "healthy" {
			state = "warning"
		}
	}
	if monitorStateRank(s.State) > monitorStateRank(state) {
		state = s.State
	}
	s.State = state
}

// ─── events ──────────────────────────────────────────

type monitorEvent struct {
	ID       int64  `json:"id"`
	ConnID   int    `json:"conn_id"`
	At       string `json:"at"`
	Kind     string `json:"kind"` // state | log | check
	Severity string `json:"severity"`
	State    string `json:"state,omitempty"`
	Text     string `json:"text"`
	Detail   string `json:"detail,omitempty"`
}

func initMonitoringV149() {
	ensureTable("monitor_events", `CREATE TABLE IF NOT EXISTS monitor_events(id INTEGER PRIMARY KEY AUTOINCREMENT,conn_id INTEGER NOT NULL,user_id INTEGER NOT NULL,at TEXT NOT NULL,kind TEXT NOT NULL DEFAULT '',severity TEXT NOT NULL DEFAULT '',state TEXT NOT NULL DEFAULT '',text TEXT NOT NULL DEFAULT '',detail TEXT NOT NULL DEFAULT '')`, "id", "conn_id", "user_id", "at", "kind", "severity", "state", "text", "detail")
	db.Exec(`CREATE INDEX IF NOT EXISTS ix_monitor_events ON monitor_events(user_id,id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS ix_monitor_events_conn ON monitor_events(conn_id,id)`)
	ensureTable("monitor_rollups", `CREATE TABLE IF NOT EXISTS monitor_rollups(conn_id INTEGER NOT NULL,hour TEXT NOT NULL,data TEXT NOT NULL,PRIMARY KEY(conn_id,hour))`, "conn_id", "hour", "data")
	ensureTable("monitor_log_offsets", `CREATE TABLE IF NOT EXISTS monitor_log_offsets(conn_id INTEGER NOT NULL,watch TEXT NOT NULL,pos INTEGER NOT NULL DEFAULT -1,last_hit TEXT NOT NULL DEFAULT '',last_count INTEGER NOT NULL DEFAULT 0,last_lines TEXT NOT NULL DEFAULT '',PRIMARY KEY(conn_id,watch))`, "conn_id", "watch", "pos", "last_hit", "last_count", "last_lines")
}

func recordMonitorEvent(c Connection, kind, severity, state, text, detail string) {
	db.Exec(`INSERT INTO monitor_events(conn_id,user_id,at,kind,severity,state,text,detail) VALUES(?,?,?,?,?,?,?,?)`,
		c.ID, c.UserID, nowRFC(), kind, severity, state, truncateStr(text, 2000), truncateStr(redactText(detail), 4000))
	db.Exec(`DELETE FROM monitor_events WHERE conn_id=? AND id NOT IN (SELECT id FROM monitor_events WHERE conn_id=? ORDER BY id DESC LIMIT 1000)`, c.ID, c.ID)
}

func monitorReasonsText(list []monitorReason, lang string) string {
	parts := []string{}
	for _, r := range list {
		parts = append(parts, r.text(lang))
		if len(parts) == 8 {
			parts = append(parts, "…")
			break
		}
	}
	return strings.Join(parts, "; ")
}

// monitorStateChanged records the event and notifies the owner when the state
// of a server changed.
func monitorStateChanged(c Connection, old, s monitorSample) {
	if old.State == s.State || (old.At == "" && s.State == "healthy") {
		return
	}
	severity := "info"
	event := "monitor.recovered"
	if s.State != "healthy" {
		severity, event = s.State, "monitor.problem"
	}
	recordMonitorEvent(c, "state", severity, s.State, monitorReasonsText(s.Reasons, "en"), s.Error)
	notifyUser(c.UserID, event, func(lang string) (string, string) {
		why := monitorReasonsText(s.Reasons, lang)
		if lang == "hr" {
			names := map[string]string{"healthy": "u redu", "warning": "upozorenje", "problem": "problem", "failed": "nedostupan"}
			return "WRM nadzor: " + c.Name + " – " + names[s.State], "Poslužitelj " + c.Name + " (" + c.Host + "): " + names[s.State] + "\n" + why
		}
		return "WRM Monitoring: " + c.Name + " – " + s.State, "Server " + c.Name + " (" + c.Host + "): " + s.State + "\n" + why
	})
}

// ─── roll-ups ────────────────────────────────────────

// monitorRollupKey keeps the metrics worth a long history.
func monitorRollupKey(k string) bool {
	switch k {
	case "cpu_percent", "memory_percent", "swap_percent", "load", "iowait", "processes", "tcp_established":
		return true
	}
	return strings.HasPrefix(k, "disk:") || strings.HasPrefix(k, "rate:")
}

type monitorRollup struct {
	N  int                   `json:"n"`
	Up int                   `json:"up"` // samples that were not failed
	OK int                   `json:"ok"` // healthy samples
	M  map[string][4]float64 `json:"m"`  // min, max, sum, count
}

func addMonitorRollup(c Connection, s monitorSample) {
	at, err := time.Parse(time.RFC3339, s.At)
	if err != nil {
		return
	}
	hour := at.UTC().Truncate(time.Hour).Format(time.RFC3339)
	var r monitorRollup
	var raw string
	if db.QueryRow(`SELECT data FROM monitor_rollups WHERE conn_id=? AND hour=?`, c.ID, hour).Scan(&raw) == nil {
		json.Unmarshal([]byte(raw), &r)
	}
	if r.M == nil {
		r.M = map[string][4]float64{}
	}
	r.N++
	if s.State != "failed" {
		r.Up++
	}
	if s.State == "healthy" {
		r.OK++
	}
	for k, m := range s.Metrics {
		if m.Value == nil || m.Error != "" || !monitorRollupKey(k) || math.IsNaN(*m.Value) || math.IsInf(*m.Value, 0) {
			continue
		}
		if len(r.M) >= 64 {
			if _, ok := r.M[k]; !ok {
				continue
			}
		}
		v := *m.Value
		x, ok := r.M[k]
		if !ok {
			x = [4]float64{v, v, 0, 0}
		}
		x[0], x[1] = math.Min(x[0], v), math.Max(x[1], v)
		x[2] += v
		x[3]++
		r.M[k] = x
	}
	db.Exec(`INSERT INTO monitor_rollups(conn_id,hour,data) VALUES(?,?,?) ON CONFLICT(conn_id,hour) DO UPDATE SET data=excluded.data`, c.ID, hour, string(jsonMarshal(r)))
	db.Exec(`DELETE FROM monitor_rollups WHERE conn_id=? AND hour<?`, c.ID, time.Now().UTC().Add(-90*24*time.Hour).Format(time.RFC3339))
}

// ─── history API ─────────────────────────────────────

type monitorSeriesPoint [4]float64 // unix seconds, avg, min, max

func monitorRangeDuration(v string) (time.Duration, bool) {
	switch v {
	case "", "24h":
		return 24 * time.Hour, true
	case "1h":
		return time.Hour, true
	case "6h":
		return 6 * time.Hour, true
	case "7d":
		return 7 * 24 * time.Hour, true
	case "30d":
		return 30 * 24 * time.Hour, true
	case "90d":
		return 90 * 24 * time.Hour, true
	}
	return 0, false
}

// monitorHistory builds series for the range: raw samples (bucketed to at most
// 360 points) up to one day, hourly roll-ups beyond.
func monitorHistory(connID int, rng string) (map[string]interface{}, error) {
	d, ok := monitorRangeDuration(rng)
	if !ok {
		return nil, fmt.Errorf("unknown range")
	}
	now := time.Now().UTC()
	since := now.Add(-d)
	series := map[string][]monitorSeriesPoint{}
	total, up, healthy := 0, 0, 0
	if d <= 24*time.Hour {
		step := d.Seconds() / 360
		type acc struct{ min, max, sum, n float64 }
		buckets := map[string]map[int64]*acc{}
		rows, err := db.Query(`SELECT data FROM monitor_samples WHERE conn_id=? AND at>=? ORDER BY id`, connID, since.Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var raw string
			rows.Scan(&raw)
			var s monitorSample
			if json.Unmarshal([]byte(raw), &s) != nil {
				continue
			}
			at, err := time.Parse(time.RFC3339, s.At)
			if err != nil {
				continue
			}
			total++
			if s.State != "failed" {
				up++
			}
			if s.State == "healthy" {
				healthy++
			}
			b := int64(math.Floor(at.Sub(since).Seconds()/step)) * int64(step)
			for k, m := range s.Metrics {
				if m.Value == nil || m.Error != "" || !monitorRollupKey(k) {
					continue
				}
				if buckets[k] == nil {
					if len(buckets) >= 64 {
						continue
					}
					buckets[k] = map[int64]*acc{}
				}
				a := buckets[k][b]
				if a == nil {
					a = &acc{min: *m.Value, max: *m.Value}
					buckets[k][b] = a
				}
				a.min, a.max = math.Min(a.min, *m.Value), math.Max(a.max, *m.Value)
				a.sum += *m.Value
				a.n++
			}
		}
		rows.Close()
		for k, bs := range buckets {
			list := make([]monitorSeriesPoint, 0, len(bs))
			for b, a := range bs {
				list = append(list, monitorSeriesPoint{float64(since.Unix() + b), round2(a.sum / a.n), round2(a.min), round2(a.max)})
			}
			sort.Slice(list, func(i, j int) bool { return list[i][0] < list[j][0] })
			series[k] = list
		}
	} else {
		rows, err := db.Query(`SELECT hour,data FROM monitor_rollups WHERE conn_id=? AND hour>=? ORDER BY hour`, connID, since.Truncate(time.Hour).Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var hour, raw string
			rows.Scan(&hour, &raw)
			var r monitorRollup
			at, err := time.Parse(time.RFC3339, hour)
			if err != nil || json.Unmarshal([]byte(raw), &r) != nil {
				continue
			}
			total += r.N
			up += r.Up
			healthy += r.OK
			for k, x := range r.M {
				if x[3] == 0 {
					continue
				}
				series[k] = append(series[k], monitorSeriesPoint{float64(at.Unix()), round2(x[2] / x[3]), round2(x[0]), round2(x[1])})
			}
		}
		rows.Close()
	}
	out := map[string]interface{}{"range": nonEmpty(rng, "24h"), "from": since.Format(time.RFC3339), "to": now.Format(time.RFC3339), "series": series, "samples": total}
	if total > 0 {
		out["availability"] = round2(100 * float64(up) / float64(total))
		out["healthy"] = round2(100 * float64(healthy) / float64(total))
	}
	return out, nil
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }

// monitorAvailability is the share of the last day's polls that reached the
// server (roll-ups, so it is cheap for the overview).
func monitorAvailability(connID int) (float64, bool) {
	var n, up int
	rows, err := db.Query(`SELECT data FROM monitor_rollups WHERE conn_id=? AND hour>=?`, connID, time.Now().UTC().Add(-24*time.Hour).Truncate(time.Hour).Format(time.RFC3339))
	if err != nil {
		return 0, false
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		rows.Scan(&raw)
		var r monitorRollup
		if json.Unmarshal([]byte(raw), &r) == nil {
			n += r.N
			up += r.Up
		}
	}
	if n == 0 {
		return 0, false
	}
	return round2(100 * float64(up) / float64(n)), true
}

// monitorSpark is the last hour of CPU and memory for the overview cards.
func monitorSpark(connID int) map[string][]float64 {
	out := map[string][]float64{"cpu_percent": {}, "memory_percent": {}}
	rows, err := db.Query(`SELECT data FROM monitor_samples WHERE conn_id=? ORDER BY id DESC LIMIT 30`, connID)
	if err != nil {
		return out
	}
	list := []monitorSample{}
	for rows.Next() {
		var raw string
		rows.Scan(&raw)
		var s monitorSample
		if json.Unmarshal([]byte(raw), &s) == nil {
			list = append(list, s)
		}
	}
	rows.Close()
	for i := len(list) - 1; i >= 0; i-- {
		for k := range out {
			if m := list[i].Metrics[k]; m.Value != nil {
				out[k] = append(out[k], round2(*m.Value))
			} else {
				out[k] = append(out[k], -1)
			}
		}
	}
	return out
}

// ─── user-level endpoints (no connection id) ─────────

// GET /api/monitoring/events?limit=&conn_id= — events of the user's servers.
func monitorEventsHandler(w http.ResponseWriter, r *http.Request, uid int) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	q := `SELECT e.id,e.conn_id,e.at,e.kind,e.severity,e.state,e.text,e.detail FROM monitor_events e JOIN connections c ON c.id=e.conn_id AND c.user_id=e.user_id WHERE e.user_id=?`
	args := []interface{}{uid}
	if id, err := strconv.Atoi(r.URL.Query().Get("conn_id")); err == nil && id > 0 {
		q += ` AND e.conn_id=?`
		args = append(args, id)
	}
	q += ` ORDER BY e.id DESC LIMIT ?`
	args = append(args, limit)
	rows, err := db.Query(q, args...)
	out := []monitorEvent{}
	if err == nil {
		for rows.Next() {
			var e monitorEvent
			if rows.Scan(&e.ID, &e.ConnID, &e.At, &e.Kind, &e.Severity, &e.State, &e.Text, &e.Detail) == nil {
				out = append(out, e)
			}
		}
		rows.Close()
	}
	jsonOK(w, out)
}

// POST /api/monitoring/bulk {ids:[], enabled:true, interval:60} — turn
// agentless monitoring on (or off) for many of the user's own SSH servers at
// once. Servers with a saved configuration keep their transport and limits.
func monitorBulkHandler(w http.ResponseWriter, r *http.Request, uid int) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in struct {
		IDs      []int `json:"ids"`
		Enabled  bool  `json:"enabled"`
		Interval int   `json:"interval"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&in) != nil || len(in.IDs) == 0 || len(in.IDs) > 500 {
		jsonError(w, "Choose 1-500 servers", 400)
		return
	}
	done := []int{}
	skipped := []int{}
	for _, id := range in.IDs {
		c, err := loadConnection(id)
		if err != nil || c.UserID != uid || (c.Protocol != "SSH" && c.Protocol != "SFTP") {
			skipped = append(skipped, id)
			continue
		}
		cfg := loadMonitorConfig(c.ID)
		cfg.Enabled = in.Enabled
		if in.Interval != 0 {
			cfg.Interval = in.Interval
		}
		if err := cfg.normalize(); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		if _, err := db.Exec(`INSERT INTO monitor_profiles(conn_id,user_id,config,updated_at) VALUES(?,?,?,?) ON CONFLICT(conn_id) DO UPDATE SET config=excluded.config,updated_at=excluded.updated_at`, c.ID, uid, encryptValue(string(jsonMarshal(cfg))), nowRFC()); err != nil {
			jsonError(w, "Cannot store monitoring settings", 500)
			return
		}
		done = append(done, c.ID)
	}
	auditLog(r, uid, "monitor.bulk_configured", "", map[string]interface{}{"servers": len(done), "enabled": in.Enabled})
	jsonOK(w, map[string]interface{}{"updated": done, "skipped": skipped})
}
