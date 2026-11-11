package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─── CHAINS: OPERATIONS ───────────────────────────────────────────────────────────────────
//
// What an operator needs around a running chain, besides the checks themselves:
//
//   - problems: the failing servers and checks, worst first, with how long they have been
//     failing and, through the links, which upstream problem probably causes them;
//   - a status word per chain (ok, warning, critical, stale, maintenance, paused, unknown) and
//     a 24-hour strip, for the overview of all chains;
//   - maintenance: a time window in which changes are still recorded but nobody is notified
//     and automatic actions do not run;
//   - acknowledging a problem ("I am on it"), cleared when the problem goes away;
//   - discovery on a live chain: who talks to whom (established TCP connections between the
//     chain's servers, with ports), outside peers, and services that run but are not checked;
//   - actions: restart or start a systemd unit or a Docker container on a chain's server, from
//     a fixed catalog, by hand (with confirmation, step-up and audit) or automatically when a
//     chosen check fails, if the user switched that on for the action and the administrator
//     allows automatic actions.

const chainOpsSchema = `CREATE TABLE IF NOT EXISTS chain_ops (
		chain_id INTEGER PRIMARY KEY,
		maint_until TEXT NOT NULL DEFAULT '',
		maint_by TEXT NOT NULL DEFAULT '',
		maint_note TEXT NOT NULL DEFAULT '',
		acks TEXT NOT NULL DEFAULT '{}',
		actions TEXT NOT NULL DEFAULT '{}')`

func initChainOpsSchema() {
	ensureTable("chain_ops", chainOpsSchema, "chain_id", "maint_until", "maint_by", "maint_note", "acks", "actions")
}

type chainAck struct {
	By   string `json:"by"`
	At   string `json:"at"`
	Note string `json:"note,omitempty"`
}

type chainActionRun struct {
	At     string `json:"at"`
	By     string `json:"by"`
	OK     bool   `json:"ok"`
	Auto   bool   `json:"auto,omitempty"`
	Output string `json:"output,omitempty"`
}

type chainOpsRec struct {
	MaintUntil string                    `json:"maint_until,omitempty"`
	MaintBy    string                    `json:"maint_by,omitempty"`
	MaintNote  string                    `json:"maint_note,omitempty"`
	Acks       map[string]chainAck       `json:"acks"`
	Actions    map[string]chainActionRun `json:"actions"`
}

var chainOpsMu sync.Mutex

func loadChainOps(id int) chainOpsRec {
	o := chainOpsRec{Acks: map[string]chainAck{}, Actions: map[string]chainActionRun{}}
	var acks, actions string
	if db.QueryRow(`SELECT maint_until, maint_by, maint_note, acks, actions FROM chain_ops WHERE chain_id=?`, id).
		Scan(&o.MaintUntil, &o.MaintBy, &o.MaintNote, &acks, &actions) == nil {
		json.Unmarshal([]byte(acks), &o.Acks)
		json.Unmarshal([]byte(actions), &o.Actions)
	}
	if o.Acks == nil {
		o.Acks = map[string]chainAck{}
	}
	if o.Actions == nil {
		o.Actions = map[string]chainActionRun{}
	}
	if !o.inMaintenance() {
		o.MaintUntil, o.MaintBy, o.MaintNote = "", "", ""
	}
	return o
}

func saveChainOps(id int, o chainOpsRec) {
	db.Exec(`INSERT INTO chain_ops (chain_id, maint_until, maint_by, maint_note, acks, actions) VALUES (?,?,?,?,?,?)
		ON CONFLICT(chain_id) DO UPDATE SET maint_until=excluded.maint_until, maint_by=excluded.maint_by,
		maint_note=excluded.maint_note, acks=excluded.acks, actions=excluded.actions`,
		id, o.MaintUntil, o.MaintBy, o.MaintNote, string(jsonMarshal(o.Acks)), string(jsonMarshal(o.Actions)))
}

// updateChainOps changes the operations record of a chain under one lock.
func updateChainOps(id int, fn func(o *chainOpsRec)) chainOpsRec {
	chainOpsMu.Lock()
	defer chainOpsMu.Unlock()
	o := loadChainOps(id)
	fn(&o)
	saveChainOps(id, o)
	return o
}

func (o chainOpsRec) inMaintenance() bool {
	t, err := time.Parse(time.RFC3339, o.MaintUntil)
	return err == nil && time.Now().Before(t)
}

func chainInMaintenance(id int) bool { return loadChainOps(id).inMaintenance() }

// ─── problems and status ──────────────────────────────────────────────────────────────────

type chainProblem struct {
	Key      string    `json:"key"` // node:<id> or check:<id>
	Kind     string    `json:"kind"`
	Severity string    `json:"severity"` // critical, warning
	Node     string    `json:"node"`
	NodeName string    `json:"node_name"`
	Label    string    `json:"label"`
	Value    string    `json:"value,omitempty"`
	Detail   string    `json:"detail,omitempty"`
	Since    string    `json:"since,omitempty"`
	Cause    string    `json:"cause,omitempty"` // the upstream node whose problem probably explains this one
	Ack      *chainAck `json:"ack,omitempty"`
	Actions  []string  `json:"actions,omitempty"` // actions that are meant to fix it
}

func chainNodeName(d chainDef, id string) string {
	n := d.node(id)
	if n.Label != "" {
		return n.Label
	}
	return id
}

// chainSince returns, per node:/check: key, when its value last changed.
func chainSince(id int) map[string]string {
	out := map[string]string{}
	rows, err := db.Query(`SELECT kind, subject, MAX(at) FROM chain_samples WHERE chain_id=? AND kind IN ('node','check') GROUP BY kind, subject`, id)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var kind, subject, at string
		if rows.Scan(&kind, &subject, &at) == nil {
			out[kind+":"+subject] = at
		}
	}
	return out
}

// chainProblems lists what is failing now. Checks on a server that is down are folded into the
// server's problem (they cannot run), and a problem whose upstream node (through a link) has a
// problem itself names that node as the probable cause.
func chainProblems(c chainRow, o chainOpsRec, since map[string]string) []chainProblem {
	d, s := c.Def, c.State
	var out []chainProblem
	down := map[string]bool{}
	for _, n := range d.Nodes {
		if n.Kind != "server" {
			continue
		}
		ns, ok := s.Nodes[n.ID]
		if ok && ns.State == "down" {
			down[n.ID] = true
			out = append(out, chainProblem{Key: "node:" + n.ID, Kind: "node", Severity: "critical", Node: n.ID,
				NodeName: chainNodeName(d, n.ID), Label: chainNodeName(d, n.ID), Value: "down", Detail: ns.Error})
		}
	}
	for _, x := range d.Checks {
		r, ok := s.Checks[x.ID]
		if !ok || r.OK || down[x.Node] {
			continue
		}
		sev := "critical"
		if x.Severity == "warning" {
			sev = "warning"
		}
		label := x.Label
		if label == "" {
			label = x.ID
		}
		detail := r.Detail
		if r.Error != "" {
			detail = r.Error
		}
		out = append(out, chainProblem{Key: "check:" + x.ID, Kind: "check", Severity: sev, Node: x.Node,
			NodeName: chainNodeName(d, x.Node), Label: label, Value: r.Value, Detail: truncateStr(detail, 300)})
	}
	bad := map[string]bool{}
	for _, p := range out {
		if p.Severity == "critical" {
			bad[p.Node] = true
		}
	}
	for i := range out {
		p := &out[i]
		p.Since = since[p.Key]
		if a, ok := o.Acks[p.Key]; ok {
			a := a
			p.Ack = &a
		}
		p.Cause = chainUpstreamCause(d, p.Node, bad)
		for _, a := range d.Actions {
			if (p.Kind == "check" && a.On == strings.TrimPrefix(p.Key, "check:")) || (a.On == "" && a.Node == p.Node) {
				p.Actions = append(p.Actions, a.ID)
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		ri, rj := chainProblemRank(out[i]), chainProblemRank(out[j])
		if ri != rj {
			return ri < rj
		}
		return out[i].Since < out[j].Since
	})
	return out
}

func chainProblemRank(p chainProblem) int {
	r := 0
	if p.Severity != "critical" {
		r += 4
	}
	if p.Cause != "" {
		r += 2 // a consequence after its cause
	}
	if p.Ack != nil {
		r++
	}
	return r
}

// chainUpstreamCause walks the links backwards from node and returns the nearest upstream
// node that has a critical problem itself.
func chainUpstreamCause(d chainDef, node string, bad map[string]bool) string {
	seen := map[string]bool{node: true}
	queue := []string{node}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, l := range d.Links {
			if l.To != cur || seen[l.From] || l.Kind == "failover" {
				continue
			}
			if bad[l.From] {
				return chainNodeName(d, l.From)
			}
			seen[l.From] = true
			queue = append(queue, l.From)
		}
	}
	return ""
}

// chainStatus is the one word the overview shows for a chain.
func chainStatus(c chainRow, o chainOpsRec, problems []chainProblem) string {
	switch {
	case !c.Enabled:
		return "paused"
	case c.LastAt == "":
		return "unknown"
	case o.inMaintenance():
		return "maintenance"
	}
	if at, err := time.Parse(time.RFC3339, c.LastAt); err == nil && time.Since(at) > 3*c.Def.interval()+time.Minute {
		return "stale"
	}
	status := "ok"
	for _, p := range problems {
		if p.Severity == "critical" {
			return "critical"
		}
		status = "warning"
	}
	return status
}

// chainHistory24 is a 48-slot strip (30 minutes each) of the last 24 hours: "ok", "warning",
// "critical" or "" (no data yet), rebuilt from the recorded changes.
func chainHistory24(c chainRow) []string {
	const slots = 48
	out := make([]string, slots)
	now := time.Now().UTC()
	start := now.Add(-24 * time.Hour)
	sev := map[string]string{}
	for _, x := range c.Def.Checks {
		if x.Severity == "warning" {
			sev["check:"+x.ID] = "warning"
		}
	}
	bad := map[string]bool{}
	apply := func(kind, subject, newV string, ok bool) {
		key := kind + ":" + subject
		switch kind {
		case "node":
			bad[key] = newV == "down"
		case "check":
			bad[key] = !ok
		}
	}
	// the state at the start of the window: the last change of every subject before it
	if rows, err := db.Query(`SELECT kind, subject, new_value, ok, MAX(id) FROM chain_samples WHERE chain_id=? AND at < ? AND kind IN ('node','check') GROUP BY kind, subject`,
		c.ID, start.Format(time.RFC3339)); err == nil {
		for rows.Next() {
			var kind, subject, newV string
			var ok, id int
			if rows.Scan(&kind, &subject, &newV, &ok, &id) == nil {
				apply(kind, subject, newV, ok == 1)
			}
		}
		rows.Close()
	}
	type ev struct {
		at                  time.Time
		kind, subject, newV string
		ok                  bool
	}
	var evs []ev
	if rows, err := db.Query(`SELECT at, kind, subject, new_value, ok FROM chain_samples WHERE chain_id=? AND at >= ? AND kind IN ('node','check') ORDER BY id`,
		c.ID, start.Format(time.RFC3339)); err == nil {
		for rows.Next() {
			var at, kind, subject, newV string
			var ok int
			if rows.Scan(&at, &kind, &subject, &newV, &ok) == nil {
				if t, err := time.Parse(time.RFC3339, at); err == nil {
					evs = append(evs, ev{t, kind, subject, newV, ok == 1})
				}
			}
		}
		rows.Close()
	}
	created, _ := time.Parse(time.RFC3339, c.CreatedAt)
	worst := func() string {
		w := "ok"
		for k, b := range bad {
			if !b {
				continue
			}
			if sev[k] == "warning" {
				w = "warning"
				continue
			}
			return "critical"
		}
		return w
	}
	rank := map[string]int{"": 0, "ok": 1, "warning": 2, "critical": 3}
	i := 0
	for s := 0; s < slots; s++ {
		from := start.Add(time.Duration(s) * 30 * time.Minute)
		to := from.Add(30 * time.Minute)
		if to.Before(created) || c.LastAt == "" {
			continue
		}
		cur := worst()
		for i < len(evs) && evs[i].at.Before(to) {
			apply(evs[i].kind, evs[i].subject, evs[i].newV, evs[i].ok)
			if w := worst(); rank[w] > rank[cur] {
				cur = w
			}
			i++
		}
		out[s] = cur
	}
	return out
}

// ─── API ──────────────────────────────────────────────────────────────────────────────────

// chainOpsView is what the chain window and the overview show beyond the definition and state.
func chainOpsView(c chainRow) map[string]interface{} {
	o := loadChainOps(c.ID)
	problems := chainProblems(c, o, chainSince(c.ID))
	if problems == nil {
		problems = []chainProblem{}
	}
	runs := map[string]chainActionRun{}
	for _, a := range c.Def.Actions {
		if r, ok := o.Actions[a.ID]; ok {
			runs[a.ID] = r
		}
	}
	return map[string]interface{}{
		"status": chainStatus(c, o, problems), "problems": problems,
		"maintenance": map[string]interface{}{"active": o.inMaintenance(), "until": o.MaintUntil, "by": o.MaintBy, "note": o.MaintNote},
		"history24":   chainHistory24(c), "action_runs": runs,
		"actions_allowed": getSetting("chains_actions_who") != "off", "auto_actions_allowed": settingBool("chains_auto_actions"),
	}
}

func chainUserName(r *http.Request, userID int) string {
	var name string
	if db.QueryRow(`SELECT username FROM users WHERE id=?`, userID).Scan(&name) == nil && name != "" {
		return name
	}
	return "user " + strconv.Itoa(userID)
}

// apiChainOps answers /api/chains/{id}/{ops|maintenance|ack|discover-links|actions/...}.
func apiChainOps(w http.ResponseWriter, r *http.Request, c chainRow, userID int, parts []string) {
	switch parts[1] {
	case "ops":
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		jsonOK(w, chainOpsView(c))
	case "maintenance":
		apiChainMaintenance(w, r, c, userID)
	case "ack":
		apiChainAck(w, r, c, userID)
	case "discover-links":
		apiChainDiscoverLinks(w, r, c, userID)
	case "actions":
		if len(parts) != 4 || parts[3] != "run" {
			jsonError(w, "Not found", http.StatusNotFound)
			return
		}
		apiChainActionRun(w, r, c, userID, parts[2])
	default:
		jsonError(w, "Not found", http.StatusNotFound)
	}
}

func apiChainMaintenance(w http.ResponseWriter, r *http.Request, c chainRow, userID int) {
	switch r.Method {
	case http.MethodPost:
		var in struct {
			Minutes int    `json:"minutes"`
			Note    string `json:"note"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&in); err != nil {
			jsonError(w, "The request is not valid JSON", http.StatusBadRequest)
			return
		}
		if in.Minutes < 5 || in.Minutes > 7*24*60 {
			jsonError(w, "Maintenance lasts 5 minutes to 7 days", http.StatusBadRequest)
			return
		}
		in.Note = strings.TrimSpace(in.Note)
		if len(in.Note) > 300 || chainHasControl(in.Note) {
			jsonError(w, "The note is at most 300 characters", http.StatusBadRequest)
			return
		}
		until := time.Now().UTC().Add(time.Duration(in.Minutes) * time.Minute).Format(time.RFC3339)
		by := chainUserName(r, userID)
		updateChainOps(c.ID, func(o *chainOpsRec) { o.MaintUntil, o.MaintBy, o.MaintNote = until, by, in.Note })
		db.Exec(`INSERT INTO chain_samples (chain_id, at, kind, subject, old_value, new_value, ok) VALUES (?,?,?,?,?,?,?)`,
			c.ID, nowRFC(), "maintenance", "maintenance", "", "on until "+until, 1)
		auditLog(r, userID, "chain.maintenance_on", c.Def.Name, map[string]interface{}{"chain": c.ID, "minutes": in.Minutes})
	case http.MethodDelete:
		updateChainOps(c.ID, func(o *chainOpsRec) { o.MaintUntil, o.MaintBy, o.MaintNote = "", "", "" })
		db.Exec(`INSERT INTO chain_samples (chain_id, at, kind, subject, old_value, new_value, ok) VALUES (?,?,?,?,?,?,?)`,
			c.ID, nowRFC(), "maintenance", "maintenance", "", "off", 1)
		auditLog(r, userID, "chain.maintenance_off", c.Def.Name, map[string]interface{}{"chain": c.ID})
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hub.sendTo(c.UserID, []byte(fmt.Sprintf(`{"type":"chain_changed","id":%d}`, c.ID)))
	jsonOK(w, chainOpsView(c))
}

var chainAckKey = regexp.MustCompile(`^(node|check):[a-z0-9][a-z0-9_-]{0,39}$`)

func apiChainAck(w http.ResponseWriter, r *http.Request, c chainRow, userID int) {
	var in struct {
		Key  string `json:"key"`
		Note string `json:"note"`
	}
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&in); err != nil || !chainAckKey.MatchString(in.Key) {
		jsonError(w, "key must be node:<id> or check:<id>", http.StatusBadRequest)
		return
	}
	in.Note = strings.TrimSpace(in.Note)
	if len(in.Note) > 300 || chainHasControl(in.Note) {
		jsonError(w, "The note is at most 300 characters", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodPost {
		by := chainUserName(r, userID)
		updateChainOps(c.ID, func(o *chainOpsRec) { o.Acks[in.Key] = chainAck{By: by, At: nowRFC(), Note: in.Note} })
		auditLog(r, userID, "chain.acknowledged", c.Def.Name, map[string]interface{}{"chain": c.ID, "subject": in.Key})
	} else {
		updateChainOps(c.ID, func(o *chainOpsRec) { delete(o.Acks, in.Key) })
	}
	jsonOK(w, chainOpsView(c))
}

// chainOpsAfterPoll runs after the changes of a poll were recorded: an acknowledgement ends when
// its problem is gone, and an automatic action runs when its check turned bad.
func chainOpsAfterPoll(c chainRow, changes []chainChange) {
	o := loadChainOps(c.ID)
	if len(o.Acks) > 0 {
		var gone []string
		for key := range o.Acks {
			kind, id, _ := strings.Cut(key, ":")
			switch kind {
			case "node":
				if c.State.Nodes[id].State != "down" {
					gone = append(gone, key)
				}
			case "check":
				if r, ok := c.State.Checks[id]; !ok || r.OK {
					gone = append(gone, key)
				}
			}
		}
		if len(gone) > 0 {
			updateChainOps(c.ID, func(o *chainOpsRec) {
				for _, k := range gone {
					delete(o.Acks, k)
				}
			})
		}
	}
	if o.inMaintenance() || !settingBool("chains_auto_actions") || !policyWhoAllows(c.UserID, "chains_actions_who") {
		return
	}
	for _, ch := range changes {
		if ch.Kind != "check" || !ch.OKFlip || ch.OK {
			continue
		}
		for _, a := range c.Def.Actions {
			if !a.Auto || a.On != ch.Subject {
				continue
			}
			if last, ok := o.Actions[a.ID]; ok {
				if t, err := time.Parse(time.RFC3339, last.At); err == nil && time.Since(t) < chainAutoCooldown {
					continue
				}
			}
			go func(a chainAction) {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				run := runChainAction(ctx, c, a, "automatic", true)
				auditLog(nil, c.UserID, "chain.action_auto", c.Def.Name, map[string]interface{}{"chain": c.ID, "action": a.ID, "ok": run.OK})
				label := a.Label
				if label == "" {
					label = a.ID
				}
				notifyUser(c.UserID, "chain.action_auto", func(lang string) (string, string) {
					res := "succeeded"
					if lang == "hr" {
						res = "uspjela"
						if !run.OK {
							res = "nije uspjela"
						}
						return "Lanac " + c.Def.Name + ": automatska radnja " + label + " " + res, run.Output
					}
					if !run.OK {
						res = "failed"
					}
					return "Chain " + c.Def.Name + ": automatic action " + label + " " + res, run.Output
				})
			}(a)
		}
	}
}

const chainAutoCooldown = 15 * time.Minute

// ─── actions ──────────────────────────────────────────────────────────────────────────────

// chainAction is one entry of a fixed catalog: what to do, on which server node, to which unit
// or container. There is no free command.
type chainAction struct {
	ID     string `json:"id"`
	Node   string `json:"node"`
	Type   string `json:"type"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
	On     string `json:"on,omitempty"`   // the check this action is meant to fix
	Auto   bool   `json:"auto,omitempty"` // run when that check turns bad (needs on)
}

var chainActionTypes = map[string]string{
	"restart_unit":      "restart a systemd unit",
	"start_unit":        "start a systemd unit",
	"restart_container": "restart a Docker container",
	"start_container":   "start a Docker container",
}

const chainMaxActions = 20

func validateChainActions(d *chainDef, nodes map[string]chainNode, checks map[string]bool) []string {
	var errs []string
	if len(d.Actions) > chainMaxActions {
		return []string{fmt.Sprintf("actions: at most %d", chainMaxActions)}
	}
	seen := map[string]bool{}
	for i := range d.Actions {
		a := &d.Actions[i]
		p := fmt.Sprintf("actions[%d]", i)
		if !chainIDPattern.MatchString(a.ID) || seen[a.ID] {
			errs = append(errs, p+": id must be unique, 1-40 characters a-z, 0-9, '-' or '_'")
			continue
		}
		seen[a.ID] = true
		if _, ok := chainActionTypes[a.Type]; !ok {
			errs = append(errs, p+": type must be restart_unit, start_unit, restart_container or start_container")
			continue
		}
		if n, ok := nodes[a.Node]; !ok || n.Kind != "server" {
			errs = append(errs, p+": node must be a server node of the chain")
		}
		a.Target = strings.TrimSpace(a.Target)
		if strings.HasSuffix(a.Type, "_unit") && !monitorUnitPattern.MatchString(a.Target) {
			errs = append(errs, p+": target is not a valid unit name")
		}
		if strings.HasSuffix(a.Type, "_container") && !monitorContainerPattern.MatchString(a.Target) {
			errs = append(errs, p+": target is not a valid container name")
		}
		if len(a.Label) > 120 || chainHasControl(a.Label) {
			errs = append(errs, p+": label at most 120 characters")
		}
		if a.On != "" && !checks[a.On] {
			errs = append(errs, p+": unknown check "+strconv.Quote(a.On))
		}
		if a.Auto && a.On == "" {
			errs = append(errs, p+": an automatic action needs the check it fixes (on)")
		}
	}
	return errs
}

// chainActionScript runs one catalog action. The type and the target arrive on standard input
// as data; the target was validated against the unit or container name pattern. Without root
// it uses sudo -n, which fails at once instead of asking for a password.
const chainActionScript = `IFS=$(printf '\t') read -r ty a
S=""; [ "$(id -u)" = 0 ] || S="sudo -n"
case "$ty" in
restart_unit) out=$($S systemctl restart -- "$a" 2>&1); rc=$? ;;
start_unit) out=$($S systemctl start -- "$a" 2>&1); rc=$? ;;
restart_container) out=$(docker restart -- "$a" 2>&1) || out=$($S docker restart -- "$a" 2>&1); rc=$? ;;
start_container) out=$(docker start -- "$a" 2>&1) || out=$($S docker start -- "$a" 2>&1); rc=$? ;;
*) out="unknown action"; rc=2 ;;
esac
printf '%s\n' "$out" | tail -c 1500
printf '@rc %s\n' "$rc"`

func runChainAction(ctx context.Context, c chainRow, a chainAction, by string, auto bool) chainActionRun {
	run := chainActionRun{At: nowRFC(), By: by, Auto: auto}
	n := c.Def.node(a.Node)
	conn, err := loadConnection(n.ConnID)
	if err != nil || conn.UserID != c.UserID || !isSSHProtocol(conn) {
		run.Output = "the connection is not available to the chain's owner"
	} else if cl, err := dialSSH(conn, nil); err != nil {
		run.Output = chainShortErr(err)
	} else {
		out, _ := monitorSSHCommandInput(ctx, cl, "sh -c "+shellQuote(chainActionScript), a.Type+"\t"+a.Target+"\n")
		cl.Close()
		body, rc, _ := strings.Cut(out, "@rc ")
		run.OK = strings.TrimSpace(rc) == "0"
		run.Output = truncateStr(strings.TrimSpace(redactText(body)), 600)
		if run.Output == "" && run.OK {
			run.Output = "done"
		}
	}
	updateChainOps(c.ID, func(o *chainOpsRec) { o.Actions[a.ID] = run })
	db.Exec(`INSERT INTO chain_samples (chain_id, at, kind, subject, old_value, new_value, ok) VALUES (?,?,?,?,?,?,?)`,
		c.ID, run.At, "action", a.ID, by, truncateStr(run.Output, 300), boolInt(run.OK))
	// look again shortly, so the result of the action shows without waiting for the interval
	go func() {
		time.Sleep(8 * time.Second)
		pctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()
		pollChain(pctx, c.ID)
	}()
	return run
}

func apiChainActionRun(w http.ResponseWriter, r *http.Request, c chainRow, userID int, actionID string) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !policyWhoAllows(userID, "chains_actions_who") {
		jsonError(w, "Chain actions are not allowed for you by the administrator", http.StatusForbidden)
		return
	}
	var a *chainAction
	for i := range c.Def.Actions {
		if c.Def.Actions[i].ID == actionID {
			a = &c.Def.Actions[i]
		}
	}
	if a == nil {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	if !requireStepUp(w, r, "chain.action") {
		return
	}
	if !mcpRateOK(fmt.Sprintf("chain-action:%d", userID), 10, time.Minute) {
		jsonError(w, "Too many actions; wait a moment", http.StatusTooManyRequests)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	run := runChainAction(ctx, c, *a, chainUserName(r, userID), false)
	auditLog(r, userID, "chain.action_run", c.Def.Name, map[string]interface{}{"chain": c.ID, "action": a.ID, "type": a.Type, "target": a.Target, "ok": run.OK})
	jsonOK(w, run)
}

// ─── discovery on a live chain ────────────────────────────────────────────────────────────

type chainLinkSuggestion struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Port   int    `json:"port"`
	Conns  int    `json:"conns"`
	Exists bool   `json:"exists"`
}

type chainPeerSuggestion struct {
	Node  string `json:"node"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Conns int    `json:"conns"`
	Dir   string `json:"dir"` // out: the node connects to it; in: clients on a listening port
}

type chainCheckSuggestion struct {
	Node   string `json:"node"`
	Type   string `json:"type"`
	Target string `json:"target,omitempty"`
	Port   int    `json:"port,omitempty"`
	Label  string `json:"label"`
}

type chainStaleCheck struct {
	Check  string `json:"check"`
	Reason string `json:"reason"`
}

func chainCleanIP(s string) string {
	s = strings.Trim(strings.TrimSpace(s), "[]")
	if i := strings.Index(s, "%"); i >= 0 {
		s = s[:i]
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return ""
	}
	if v4 := ip.To4(); v4 != nil {
		return v4.String()
	}
	return ip.String()
}

// chainSplitAddr parses "1.2.3.4:80", "[::1]:80" or "::ffff:1.2.3.4:80" (ss and netstat style).
func chainSplitAddr(s string) (string, int, bool) {
	i := strings.LastIndexAny(s, ":.")
	if i <= 0 {
		return "", 0, false
	}
	p, err := strconv.Atoi(s[i+1:])
	ip := chainCleanIP(s[:i])
	if err != nil || ip == "" || p < 1 || p > 65535 {
		return "", 0, false
	}
	return ip, p, true
}

func chainLoopback(ip string) bool {
	p := net.ParseIP(ip)
	return p == nil || p.IsLoopback() || p.IsUnspecified() || p.IsLinkLocalUnicast()
}

// chainDiscoverLinks turns the facts of the chain's servers into suggestions. A connection
// from node A to node B on a port B listens on is the link A → B; the same connection seen
// from B's side is counted once.
func chainDiscoverLinks(d chainDef, facts map[string]chainServerFacts, hostIPs map[string][]string) ([]chainLinkSuggestion, []chainPeerSuggestion) {
	owner := map[string]string{}
	for id, f := range facts {
		for _, ip := range f.Addrs {
			if !chainLoopback(ip) {
				owner[ip] = id
			}
		}
	}
	for id, ips := range hostIPs {
		for _, ip := range ips {
			if _, ok := owner[ip]; !ok && !chainLoopback(ip) {
				owner[ip] = id
			}
		}
	}
	exists := map[string]bool{}
	for _, l := range d.Links {
		exists[l.From+">"+l.To] = true
	}
	links := map[string]*chainLinkSuggestion{}
	out := map[string]*chainPeerSuggestion{}
	in := map[string]map[string]bool{}
	ids := make([]string, 0, len(facts))
	for id := range facts {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		f := facts[id]
		for _, e := range f.Est {
			if chainLoopback(e.PeerIP) {
				continue
			}
			client := !f.Listen[e.LocalP] // our side is the ephemeral one
			other, isNode := owner[e.PeerIP]
			if isNode && other == id {
				continue
			}
			switch {
			case isNode && client:
				k := id + ">" + other + ">" + strconv.Itoa(e.PeerPort)
				if links[k] == nil {
					links[k] = &chainLinkSuggestion{From: id, To: other, Port: e.PeerPort, Exists: exists[id+">"+other]}
				}
				links[k].Conns++
			case isNode && !client:
				k := other + ">" + id + ">" + strconv.Itoa(e.LocalP)
				if links[k] == nil {
					links[k] = &chainLinkSuggestion{From: other, To: id, Port: e.LocalP, Exists: exists[other+">"+id]}
				}
				// counted from the client's side when that node is in the chain too
				if _, seen := facts[other]; !seen {
					links[k].Conns++
				}
			case client:
				k := id + ">" + e.PeerIP + ":" + strconv.Itoa(e.PeerPort)
				if out[k] == nil {
					out[k] = &chainPeerSuggestion{Node: id, Host: e.PeerIP, Port: e.PeerPort, Dir: "out"}
				}
				out[k].Conns++
			default:
				k := id + ">" + strconv.Itoa(e.LocalP)
				if in[k] == nil {
					in[k] = map[string]bool{}
				}
				in[k][e.PeerIP] = true
			}
		}
	}
	var ls []chainLinkSuggestion
	for _, l := range links {
		if l.Conns == 0 {
			l.Conns = 1
		}
		ls = append(ls, *l)
	}
	sort.Slice(ls, func(i, j int) bool {
		if ls[i].From != ls[j].From {
			return ls[i].From < ls[j].From
		}
		if ls[i].To != ls[j].To {
			return ls[i].To < ls[j].To
		}
		return ls[i].Port < ls[j].Port
	})
	var ps []chainPeerSuggestion
	for _, p := range out {
		ps = append(ps, *p)
	}
	for k, peers := range in {
		node, port, _ := strings.Cut(k, ">")
		n, _ := strconv.Atoi(port)
		if n == 22 {
			continue // WRM's own SSH sessions
		}
		ps = append(ps, chainPeerSuggestion{Node: node, Port: n, Conns: len(peers), Dir: "in"})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].Conns != ps[j].Conns {
			return ps[i].Conns > ps[j].Conns
		}
		return ps[i].Node+ps[i].Host < ps[j].Node+ps[j].Host
	})
	if len(ps) > 30 {
		ps = ps[:30]
	}
	return ls, ps
}

// chainDiscoverChecks lists what runs on a server but is not checked, and checks whose service
// or container is not there any more.
func chainDiscoverChecks(d chainDef, facts map[string]chainServerFacts, state chainState) ([]chainCheckSuggestion, []chainStaleCheck) {
	have := map[string]bool{}
	for _, x := range d.Checks {
		have[x.Node+"|"+x.Type+"|"+chainNodeTarget(x)] = true
	}
	var sugg []chainCheckSuggestion
	var stale []chainStaleCheck
	ids := make([]string, 0, len(facts))
	for id := range facts {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		f := facts[id]
		if f.Error != "" {
			continue
		}
		running := map[string]bool{}
		for _, s := range f.Services {
			running["systemd|"+s] = true
			if !have[id+"|systemd|"+s] {
				sugg = append(sugg, chainCheckSuggestion{Node: id, Type: "systemd", Target: s, Label: strings.TrimSuffix(s, ".service")})
			}
		}
		for _, c := range f.Containers {
			running["docker|"+c] = true
			if !have[id+"|docker|"+c] {
				sugg = append(sugg, chainCheckSuggestion{Node: id, Type: "docker", Target: c, Label: c})
			}
		}
		for _, p := range f.Ports {
			if !have[id+"|port|"+strconv.Itoa(p)] {
				sugg = append(sugg, chainCheckSuggestion{Node: id, Type: "port", Port: p, Label: chainPortLabel(p)})
			}
		}
		for _, k := range []string{"disk", "memory", "load"} {
			target := ""
			if k == "disk" {
				target = "/"
			}
			if !have[id+"|"+k+"|"+target] {
				sugg = append(sugg, chainCheckSuggestion{Node: id, Type: k, Target: target, Label: map[string]string{"disk": "Disk /", "memory": "Memory", "load": "Load"}[k]})
			}
		}
		for _, x := range d.Checks {
			if x.Node != id || (x.Type != "systemd" && x.Type != "docker") {
				continue
			}
			r, ok := state.Checks[x.ID]
			if ok && r.OK {
				continue
			}
			if !running[x.Type+"|"+x.Target] {
				what := "the service is not running"
				if x.Type == "docker" {
					what = "the container is not running"
				}
				stale = append(stale, chainStaleCheck{Check: x.ID, Reason: what})
			}
		}
	}
	if len(sugg) > 60 {
		sugg = sugg[:60]
	}
	return sugg, stale
}

func apiChainDiscoverLinks(w http.ResponseWriter, r *http.Request, c chainRow, userID int) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !mcpRateOK("chain-discover:"+strconv.Itoa(userID), 20, 10*time.Minute) {
		jsonError(w, "Too many discoveries in a short time; wait a moment", http.StatusTooManyRequests)
		return
	}
	var conns []Connection
	nodeOf := map[int]string{}
	for _, n := range c.Def.Nodes {
		if n.Kind != "server" {
			continue
		}
		conn, err := loadConnectionRaw(n.ConnID)
		if err != nil || conn.UserID != c.UserID || !isSSHProtocol(conn) {
			continue
		}
		if _, dup := nodeOf[conn.ID]; dup {
			continue
		}
		nodeOf[conn.ID] = n.ID
		conns = append(conns, conn)
	}
	if len(conns) == 0 {
		jsonError(w, "The chain has no server nodes to look at", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	facts := map[string]chainServerFacts{}
	hostIPs := map[string][]string{}
	var servers []map[string]interface{}
	for _, f := range discoverChainServers(ctx, conns) {
		id := nodeOf[f.ConnID]
		facts[id] = f
		servers = append(servers, map[string]interface{}{"node": id, "hostname": f.Hostname, "os": f.OS, "error": f.Error, "addrs": f.Addrs})
	}
	for _, conn := range conns {
		host := conn.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		if ip := chainCleanIP(host); ip != "" {
			hostIPs[nodeOf[conn.ID]] = []string{ip}
			continue
		}
		lctx, lcancel := context.WithTimeout(ctx, 3*time.Second)
		ips, _ := net.DefaultResolver.LookupHost(lctx, host)
		lcancel()
		for _, ip := range ips {
			if ip = chainCleanIP(ip); ip != "" {
				hostIPs[nodeOf[conn.ID]] = append(hostIPs[nodeOf[conn.ID]], ip)
			}
		}
	}
	links, peers := chainDiscoverLinks(c.Def, facts, hostIPs)
	checks, stale := chainDiscoverChecks(c.Def, facts, c.State)
	auditLog(r, userID, "chain.discovered", c.Def.Name, map[string]interface{}{"chain": c.ID, "servers": len(conns)})
	jsonOK(w, map[string]interface{}{"servers": servers, "links": chainNonNil(links), "peers": chainNonNil(peers),
		"checks": chainNonNil(checks), "stale": chainNonNil(stale)})
}

func chainNonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
