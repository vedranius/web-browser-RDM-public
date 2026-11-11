package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ─── LOG PATTERNS (v14.9.0) ──────────────────────────
//
// "What is in these logs?" The last lines of the chosen logs are read on the
// server (tail / journalctl, one SSH command) and grouped here into patterns:
// numbers, addresses, ids, hashes, quoted values and timestamps become
// placeholders, so ten thousand "connection from 10.0.0.7 timed out after
// 3012 ms" lines are one pattern with a count, a level, first and last time,
// the files it came from and a small timeline. Patterns seen only in the most
// recent fifth of the window are marked new, and a pattern whose rate in that
// last fifth is four times its earlier rate is marked rising.

type logPatternRequest struct {
	Files    []string `json:"files"`
	Lines    int      `json:"lines"`
	Timezone string   `json:"timezone"`
	Level    string   `json:"level"` // "", error, warn
}

type logPattern struct {
	Template string   `json:"template"`
	Example  string   `json:"example"`
	Count    int      `json:"count"`
	Level    string   `json:"level"`
	First    string   `json:"first,omitempty"`
	Last     string   `json:"last,omitempty"`
	Files    []string `json:"files"`
	Buckets  []int    `json:"buckets"`
	New      bool     `json:"new,omitempty"`
	Rising   bool     `json:"rising,omitempty"`
	firstAt  time.Time
	lastAt   time.Time
	times    []time.Time
	fileSet  map[string]bool
}

var (
	lpTimestamp = regexp.MustCompile(`^\s*(\[?\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?\]?|[A-Z][a-z]{2}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}|\d{2}[./]\d{2}[./]\d{4}\s+\d{2}:\d{2}:\d{2}|\[\d{2}/[A-Z][a-z]{2}/\d{4}:\d{2}:\d{2}:\d{2}\s[+-]\d{4}\])\s*`)
	lpRules     = []struct {
		re  *regexp.Regexp
		rep string
	}{
		{regexp.MustCompile(`\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b`), "<uuid>"},
		{regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?`), "<time>"},
		{regexp.MustCompile(`\b\d{2}:\d{2}:\d{2}(?:[.,]\d+)?\b`), "<time>"},
		{regexp.MustCompile(`\b[\w.+-]+@[\w-]+\.[\w.-]+\b`), "<email>"},
		{regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}(?::\d+)?\b`), "<ip>"},
		{regexp.MustCompile(`\b(?:[0-9a-fA-F]{1,4}:){3,7}[0-9a-fA-F]{1,4}\b`), "<ip>"},
		{regexp.MustCompile(`\b(?:0x)?[0-9a-fA-F]{12,}\b`), "<hex>"},
		{regexp.MustCompile(`"[^"]{0,200}"`), `"<str>"`},
		{regexp.MustCompile(`'[^']{0,200}'`), `'<str>'`},
		{regexp.MustCompile(`\b\d+(?:\.\d+)?(?:ms|s|m|h|kb|mb|gb|KB|MB|GB|B|%)?\b`), "<n>"},
		{regexp.MustCompile(`(?:<n>[\s,:/.-]*){2,}`), "<n> "},
		{regexp.MustCompile(`\s+`), " "},
	}
	lpLevel = []struct {
		re    *regexp.Regexp
		level string
	}{
		{regexp.MustCompile(`(?i)\b(fatal|panic|emerg|emergency|alert|crit|critical|severe)\b`), "error"},
		{regexp.MustCompile(`(?i)\b(error|err|exception|traceback|failed|failure)\b`), "error"},
		{regexp.MustCompile(`(?i)\b(warn|warning|deprecated)\b`), "warn"},
		{regexp.MustCompile(`(?i)\b(info|notice)\b`), "info"},
		{regexp.MustCompile(`(?i)\b(debug|trace|verbose)\b`), "debug"},
	}
)

// logTemplate turns one log line into its pattern.
func logTemplate(line string) string {
	t := lpTimestamp.ReplaceAllString(line, "")
	if len(t) > 600 {
		t = t[:600]
	}
	for _, r := range lpRules {
		t = r.re.ReplaceAllString(t, r.rep)
	}
	return strings.TrimSpace(t)
}

func logLineLevel(line string) string {
	head := line
	if len(head) > 300 {
		head = head[:300]
	}
	for _, l := range lpLevel {
		if l.re.MatchString(head) {
			return l.level
		}
	}
	return "other"
}

type logPatternResult struct {
	Patterns  []logPattern   `json:"patterns"`
	Levels    map[string]int `json:"levels"`
	Lines     int            `json:"lines"`
	Distinct  int            `json:"distinct"`
	From      string         `json:"from,omitempty"`
	To        string         `json:"to,omitempty"`
	Truncated bool           `json:"truncated"`
	Timeline  []int          `json:"timeline"`
	ErrorLine []int          `json:"error_timeline"`
}

const logPatternBuckets = 24

// clusterLogLines groups lines (path → lines) into patterns.
func clusterLogLines(byFile map[string][]string, loc *time.Location, level string) logPatternResult {
	res := logPatternResult{Levels: map[string]int{}, Patterns: []logPattern{}, Timeline: make([]int, logPatternBuckets), ErrorLine: make([]int, logPatternBuckets)}
	pats := map[string]*logPattern{}
	var minT, maxT time.Time
	type stamped struct {
		p  *logPattern
		at time.Time
		lv string
	}
	all := []stamped{}
	files := make([]string, 0, len(byFile))
	for f := range byFile {
		files = append(files, f)
	}
	sort.Strings(files)
	for _, f := range files {
		prev := false
		for _, line := range byFile[f] {
			if strings.TrimSpace(line) == "" {
				continue
			}
			// continuation lines (stack traces) belong to the previous entry
			if (line[0] == ' ' || line[0] == '\t') && prev {
				continue
			}
			prev = true
			lv := logLineLevel(line)
			if level != "" && lv != level && !(level == "warn" && lv == "error") {
				continue
			}
			res.Lines++
			res.Levels[lv]++
			tpl := logTemplate(line)
			if tpl == "" {
				continue
			}
			p := pats[tpl]
			if p == nil {
				if len(pats) >= 5000 {
					res.Truncated = true
					continue
				}
				p = &logPattern{Template: tpl, Example: truncateStr(redactText(line), 600), Level: lv, fileSet: map[string]bool{}}
				pats[tpl] = p
			}
			p.Count++
			p.fileSet[f] = true
			if levelRank(lv) > levelRank(p.Level) {
				p.Level = lv
			}
			if at, ok := parseMonitorLogTime(line, loc); ok {
				if p.firstAt.IsZero() || at.Before(p.firstAt) {
					p.firstAt = at
				}
				if at.After(p.lastAt) {
					p.lastAt = at
				}
				if minT.IsZero() || at.Before(minT) {
					minT = at
				}
				if at.After(maxT) {
					maxT = at
				}
				all = append(all, stamped{p, at, lv})
			}
		}
	}
	res.Distinct = len(pats)
	span := maxT.Sub(minT)
	recent := maxT.Add(-span / 5)
	for _, x := range all {
		b := 0
		if span > 0 {
			b = int(float64(x.at.Sub(minT)) / float64(span) * float64(logPatternBuckets-1))
		}
		b = max(0, min(b, logPatternBuckets-1))
		res.Timeline[b]++
		if x.lv == "error" {
			res.ErrorLine[b]++
		}
		if x.p.Buckets == nil {
			x.p.Buckets = make([]int, logPatternBuckets)
		}
		x.p.Buckets[b]++
	}
	list := make([]*logPattern, 0, len(pats))
	for _, p := range pats {
		list = append(list, p)
	}
	for _, p := range list {
		if !p.firstAt.IsZero() {
			p.First = p.firstAt.Format(time.RFC3339)
			p.Last = p.lastAt.Format(time.RFC3339)
			if span > time.Minute {
				if !p.firstAt.Before(recent) && p.Count > 0 {
					p.New = true
				} else if p.Buckets != nil {
					early, late := 0, 0
					cut := logPatternBuckets * 4 / 5
					for i, n := range p.Buckets {
						if i >= cut {
							late += n
						} else {
							early += n
						}
					}
					// rate per bucket in the last fifth versus before
					if late >= 5 && float64(late)/float64(logPatternBuckets-cut) >= 4*float64(early)/float64(cut) {
						p.Rising = true
					}
				}
			}
		}
		for f := range p.fileSet {
			p.Files = append(p.Files, f)
		}
		sort.Strings(p.Files)
		if len(p.Files) > 10 {
			p.Files = p.Files[:10]
		}
		if p.Buckets == nil {
			p.Buckets = []int{}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if ra, rb := levelRank(a.Level), levelRank(b.Level); ra != rb {
			return ra > rb
		}
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		return a.Template < b.Template
	})
	if len(list) > 300 {
		list = list[:300]
		res.Truncated = true
	}
	for _, p := range list {
		res.Patterns = append(res.Patterns, *p)
	}
	if !minT.IsZero() {
		res.From, res.To = minT.Format(time.RFC3339), maxT.Format(time.RFC3339)
	}
	return res
}

func levelRank(l string) int {
	switch l {
	case "error":
		return 4
	case "warn":
		return 3
	case "info":
		return 2
	case "debug":
		return 1
	}
	return 0
}

// POST /api/monitoring/logs/patterns?id= {files, lines, timezone, level}
func logPatternsHandler(w http.ResponseWriter, r *http.Request, c Connection, acc *connAccess) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var in logPatternRequest
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in) != nil {
		jsonError(w, "Bad JSON", 400)
		return
	}
	if len(in.Files) == 0 || len(in.Files) > 200 {
		jsonError(w, "Choose 1-200 log files", 400)
		return
	}
	if !oneOf(in.Level, "", "error", "warn") {
		jsonError(w, "Unknown level", 400)
		return
	}
	loc, err := time.LoadLocation(nonEmpty(in.Timezone, "UTC"))
	if err != nil {
		jsonError(w, "Unknown log timezone", 400)
		return
	}
	files, units, ok := logSearchable(c, in.Files)
	if !ok && logRediscover(r, c, acc) {
		files, units, ok = logSearchable(c, in.Files)
	}
	if !ok {
		jsonError(w, "Find logs first: patterns use the discovered file list", 409)
		return
	}
	if len(files)+len(units) == 0 {
		jsonError(w, "None of the chosen files is a discovered log", 400)
		return
	}
	if in.Lines <= 0 {
		in.Lines = 5000
	}
	in.Lines = min(in.Lines, 50000)
	// The output is capped at 1 MiB; share the budget between the logs.
	per := max(200, min(in.Lines, 8000/(len(files)+len(units))))
	tok := make([]byte, 8)
	rand.Read(tok)
	marker := "@@WRMLP" + hex.EncodeToString(tok)
	var sb strings.Builder
	sb.WriteString("LC_ALL=C; export LC_ALL\n")
	if len(files) > 0 {
		sb.WriteString("while IFS= read -r f; do [ -f \"$f\" ] && [ ! -h \"$f\" ] && { printf '%s %s\\n' '" + marker + "' \"$f\"; tail -n " + strconv.Itoa(per) + " -- \"$f\" 2>/dev/null | cut -c1-2000; }; done\n")
	}
	for _, u := range units {
		sb.WriteString("printf '%s %s\\n' '" + marker + "' 'journal:" + u + "'; journalctl --no-pager --output=short-iso --unit=" + u + " -n " + strconv.Itoa(per) + " 2>/dev/null | cut -c1-2000\n")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	cl, err := dialSSH(c, nil)
	if err != nil {
		jsonError(w, "SSH log access failed: "+truncateStr(redactText(err.Error()), 300), 502)
		return
	}
	defer cl.Close()
	start := time.Now()
	raw, _ := monitorSSHCommandInput(ctx, cl, sb.String(), strings.Join(files, "\n")+"\n")
	raw = strings.ToValidUTF8(raw, "�")
	byFile := map[string][]string{}
	cur := ""
	for _, line := range strings.Split(raw, "\n") {
		if rest, ok := strings.CutPrefix(line, marker+" "); ok {
			cur = rest
			continue
		}
		if cur != "" {
			byFile[cur] = append(byFile[cur], strings.TrimRight(line, "\r"))
		}
	}
	res := clusterLogLines(byFile, loc, in.Level)
	res.Truncated = res.Truncated || ctx.Err() != nil || len(raw) >= 1<<20-1024
	auditFileOp(r, acc, c, "monitor.logs_patterns", c.Name, map[string]interface{}{"files": len(files), "units": len(units), "patterns": res.Distinct})
	jsonOK(w, map[string]interface{}{"result": res, "ms": time.Since(start).Milliseconds(), "per_file": per, "timezone": loc.String()})
}
