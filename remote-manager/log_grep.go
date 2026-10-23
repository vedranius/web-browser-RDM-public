package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// logGrepRequest searches many log files at once on the server. grep runs
// remotely (one SSH command), so a week of logs is searched without copying it
// to WRM; the time range is applied here on the matched lines.
type logGrepRequest struct {
	Files    []string `json:"files"`
	Q        string   `json:"q"`
	Regex    bool     `json:"regex"`
	Case     bool     `json:"case"`
	Word     bool     `json:"word"`
	Context  int      `json:"context"`
	Max      int      `json:"max"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Timezone string   `json:"timezone"`
}

type logGrepHit struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Text    string `json:"text"`
	Time    string `json:"time,omitempty"`
	Context bool   `json:"context,omitempty"`
}

func logGrepHandler(w http.ResponseWriter, r *http.Request, c Connection, acc *connAccess) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in logGrepRequest
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	in.Q = strings.TrimSpace(in.Q)
	if in.Q == "" || len(in.Q) > 500 || strings.ContainsAny(in.Q, "\x00\r\n") {
		jsonError(w, "Enter the text to search (up to 500 characters, one line)", 400)
		return
	}
	if in.Regex {
		if _, err := regexp.Compile(in.Q); err != nil {
			jsonError(w, "Invalid regular expression: "+err.Error(), 400)
			return
		}
	}
	if len(in.Files) == 0 || len(in.Files) > 2000 {
		jsonError(w, "Choose 1-2000 log files", 400)
		return
	}
	in.Context = max(0, min(in.Context, 5))
	if in.Max <= 0 {
		in.Max = 200
	}
	in.Max = min(in.Max, 2000)
	tz := nonEmpty(in.Timezone, "UTC")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		jsonError(w, "Unknown log timezone", 400)
		return
	}
	from, err1 := monitorLogTimeBound(in.From, loc)
	to, err2 := monitorLogTimeBound(in.To, loc)
	if err1 != nil || err2 != nil || (!from.IsZero() && !to.IsZero() && from.After(to)) {
		jsonError(w, "Invalid time range", 400)
		return
	}
	files, units, ok := logSearchable(c, in.Files)
	if !ok {
		jsonError(w, "Find logs first: the search uses the discovered file list", 409)
		return
	}
	if len(files)+len(units) == 0 {
		jsonError(w, "None of the chosen files is a discovered log", 400)
		return
	}
	flags := []string{"-H", "-n", "-I", "-s"}
	if !in.Regex {
		flags = append(flags, "-F")
	} else {
		flags = append(flags, "-E")
	}
	if !in.Case {
		flags = append(flags, "-i")
	}
	if in.Word {
		flags = append(flags, "-w")
	}
	if in.Context > 0 {
		flags = append(flags, "-C", strconv.Itoa(in.Context))
	}
	flags = append(flags, "-m", strconv.Itoa(in.Max))
	var sb strings.Builder
	sb.WriteString("LC_ALL=C; export LC_ALL\n")
	if len(files) > 0 {
		// find -P drops files that became symbolic links since discovery.
		sb.WriteString("find -P")
		for _, f := range files {
			sb.WriteString(" " + shellQuote(f))
		}
		sb.WriteString(" -maxdepth 0 -type f -print0 2>/dev/null | xargs -0 grep " + strings.Join(flags, " ") + " -e " + shellQuote(in.Q) + " -- 2>/dev/null | head -n 20000\n")
	}
	for _, u := range units {
		since := "-n 50000"
		if !from.IsZero() {
			since = "--since=@" + strconv.FormatInt(from.Unix(), 10)
		}
		if !to.IsZero() {
			since += " --until=@" + strconv.FormatInt(to.Unix(), 10)
		}
		jflags := []string{}
		for _, f := range flags {
			if f != "-H" {
				jflags = append(jflags, f)
			}
		}
		sb.WriteString("journalctl --no-pager --output=short-iso --unit=" + u + " " + since + " 2>/dev/null | grep " + strings.Join(jflags, " ") + " -e " + shellQuote(in.Q) + " 2>/dev/null | head -n " + strconv.Itoa(in.Max*(1+2*in.Context)) + " | sed 's|^|journal:" + u + ":|'\n")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	cl, err := dialSSH(c, nil)
	if err != nil {
		jsonError(w, "SSH log access failed: "+err.Error(), 502)
		return
	}
	defer cl.Close()
	start := time.Now()
	raw, err := monitorSSHCommand(ctx, cl, sb.String())
	if err != nil && strings.TrimSpace(raw) == "" && ctx.Err() == nil {
		// grep exits 1 when nothing matched; that is not an error.
		raw = ""
	}
	raw = strings.ToValidUTF8(raw, "�")
	known := append(append([]string{}, files...), func() []string {
		v := []string{}
		for _, u := range units {
			v = append(v, "journal:"+u)
		}
		return v
	}()...)
	sort.Slice(known, func(i, j int) bool { return len(known[i]) > len(known[j]) })
	hits := []logGrepHit{}
	perFile := map[string]int{}
	truncated := ctx.Err() != nil || len(raw) >= 1<<20-1024
	for _, line := range strings.Split(raw, "\n") {
		if line == "" || line == "--" {
			continue
		}
		var h logGrepHit
		okLine := false
		for _, k := range known {
			rest, ok := strings.CutPrefix(line, k)
			if !ok || rest == "" || (rest[0] != ':' && rest[0] != '-') {
				continue
			}
			sep := rest[0]
			num, text, ok := strings.Cut(rest[1:], string(sep))
			n, err := strconv.Atoi(num)
			if !ok || err != nil {
				continue
			}
			h = logGrepHit{Path: k, Line: n, Text: truncateStr(text, 4000), Context: sep == '-'}
			okLine = true
			break
		}
		if !okLine {
			continue
		}
		if at, ok := parseMonitorLogTime(h.Text, loc); ok {
			if (!from.IsZero() && at.Before(from)) || (!to.IsZero() && at.After(to)) {
				continue
			}
			h.Time = at.Format(time.RFC3339)
		}
		if !h.Context {
			perFile[h.Path]++
		}
		hits = append(hits, h)
		if len(hits) >= 20000 {
			truncated = true
			break
		}
	}
	auditFileOp(r, acc, c, "monitor.logs_searched", c.Name, map[string]interface{}{"files": len(files), "units": len(units), "matches": len(hits), "regex": in.Regex})
	jsonOK(w, map[string]interface{}{"hits": hits, "per_file": perFile, "truncated": truncated, "ms": time.Since(start).Milliseconds(), "timezone": tz})
}

// logSearchable keeps the requested paths that this server's discovery listed
// (or, after an SFTP discovery, paths inside a configured root) and the
// journal units it found. ok is false when nothing was discovered yet.
func logSearchable(c Connection, req []string) (files, units []string, ok bool) {
	view, ok := monitorDiscoveryView(c)
	if !ok {
		return nil, nil, false
	}
	allowed := map[string]bool{}
	journal := map[string]bool{}
	for _, f := range view.Files {
		if f.Kind == "journal" {
			journal[f.Unit] = true
		} else {
			allowed[f.Path] = true
		}
	}
	for _, rt := range view.Roots {
		if rt.Kind == "journal" {
			journal[rt.Unit] = true
		}
	}
	files, units = []string{}, []string{}
	for _, p := range req {
		if u, ok := strings.CutPrefix(p, "journal:"); ok {
			if journal[u] && monitorJournalUnit.MatchString(u) {
				units = append(units, u)
			}
			continue
		}
		p = path.Clean(p)
		if !strings.HasPrefix(p, "/") || strings.ContainsAny(p, "\x00\r\n") || aiReadDenied(p) != "" {
			continue
		}
		if view.Files == nil {
			inRoot := false
			for _, rt := range view.Roots {
				if rt.Kind == "" && strings.HasPrefix(p, strings.TrimSuffix(rt.Path, "/")+"/") {
					inRoot = true
					break
				}
			}
			if !inRoot {
				continue
			}
		} else if !allowed[p] {
			continue
		}
		files = append(files, p)
	}
	return files, units, true
}

// logLinesHandler returns numbered lines around one line of a discovered log
// file (GET ?id=&path=&line=&before=&after=), so a search hit opens at its
// place with real line numbers. sed runs on the server; nothing else is read.
func logLinesHandler(w http.ResponseWriter, r *http.Request, c Connection, acc *connAccess) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	q := r.URL.Query()
	files, _, ok := logSearchable(c, []string{q.Get("path")})
	if !ok {
		jsonError(w, "Find logs first: the viewer uses the discovered file list", 409)
		return
	}
	if len(files) != 1 {
		jsonError(w, "Not a discovered log file", 400)
		return
	}
	line, _ := strconv.Atoi(q.Get("line"))
	before, _ := strconv.Atoi(nonEmpty(q.Get("before"), "200"))
	after, _ := strconv.Atoi(nonEmpty(q.Get("after"), "300"))
	line = max(1, line)
	before = max(0, min(before, 5000))
	after = max(0, min(after, 5000))
	first := max(1, line-before)
	last := line + after
	cmd := "LC_ALL=C; export LC_ALL\nf=" + shellQuote(files[0]) + "\n" +
		"find -P \"$f\" -maxdepth 0 -type f 2>/dev/null | grep -q . || { echo 'WRM-NOFILE'; exit 0; }\n" +
		"printf 'WRM-TOTAL %s\\n' \"$(wc -l < \"$f\" 2>/dev/null)\"\n" +
		"sed -n '" + strconv.Itoa(first) + "," + strconv.Itoa(last) + "p;" + strconv.Itoa(last+1) + "q' \"$f\" 2>/dev/null | head -c 4000000\n"
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	cl, err := dialSSH(c, nil)
	if err != nil {
		jsonError(w, "SSH log access failed: "+err.Error(), 502)
		return
	}
	defer cl.Close()
	raw, err := monitorSSHCommand(ctx, cl, cmd)
	if err != nil && raw == "" {
		jsonError(w, "Reading the log failed: "+err.Error(), 502)
		return
	}
	if strings.HasPrefix(raw, "WRM-NOFILE") {
		jsonError(w, "The log file is gone or is no longer a regular file", 404)
		return
	}
	total := 0
	if head, rest, ok := strings.Cut(raw, "\n"); ok && strings.HasPrefix(head, "WRM-TOTAL ") {
		total, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(head, "WRM-TOTAL ")))
		raw = rest
	}
	text := strings.ToValidUTF8(strings.TrimSuffix(raw, "\n"), "\uFFFD")
	auditFileOp(r, acc, c, "monitor.log_lines", c.Name, map[string]interface{}{"path": files[0], "line": line})
	jsonOK(w, map[string]interface{}{"path": files[0], "first": first, "text": text, "total": total})
}
