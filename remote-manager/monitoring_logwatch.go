package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/ssh"
)

// ─── LOG ALERTS (v14.9.0) ────────────────────────────
//
// A log alert watches one log file (or systemd journal unit) for a pattern.
// On every poll WRM reads only what was appended since the previous poll (at
// most 8 MiB, with rotation detected by a shrinking file), counts the
// matching lines on the server and keeps the last three. A hit turns the
// server state to warning or problem for the hold time, records an event and
// notifies through the normal state notification. The first poll starts at
// the end of the file, so old content never alerts.

type monitorLogWatch struct {
	Name     string `json:"name"`
	Path     string `json:"path"` // absolute file or journal:<unit>
	Pattern  string `json:"pattern"`
	Regex    bool   `json:"regex"`
	Case     bool   `json:"case"`
	Severity string `json:"severity"` // warning | problem
	Min      int    `json:"min"`      // matches per poll to alert
	Hold     int    `json:"hold"`     // minutes the alert stays after a hit
}

type monitorLogAlert struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Count    int      `json:"count"`
	Severity string   `json:"severity"`
	At       string   `json:"at"`
	Lines    []string `json:"lines,omitempty"`
	New      bool     `json:"new,omitempty"` // matched in this poll (not only held)
}

func (x *monitorLogWatch) normalize() error {
	x.Name = strings.TrimSpace(x.Name)
	x.Path = strings.TrimSpace(x.Path)
	x.Pattern = strings.TrimSpace(x.Pattern)
	if x.Severity == "" {
		x.Severity = "warning"
	}
	if x.Min <= 0 {
		x.Min = 1
	}
	if x.Hold <= 0 {
		x.Hold = 15
	}
	if !oneOf(x.Severity, "warning", "problem") || x.Min > 100000 || x.Hold > 1440 {
		return fmt.Errorf("log alert: severity is warning or problem, hold up to 1440 minutes")
	}
	if len(x.Name) > 80 || strings.IndexFunc(x.Name, unicode.IsControl) >= 0 {
		return fmt.Errorf("log alert name is too long")
	}
	if x.Pattern == "" || len(x.Pattern) > 300 || strings.IndexFunc(x.Pattern, unicode.IsControl) >= 0 {
		return fmt.Errorf("log alert pattern: one line, up to 300 characters")
	}
	if x.Regex {
		if _, err := regexp.Compile(x.Pattern); err != nil {
			return fmt.Errorf("log alert pattern: %v", err)
		}
	}
	if u, ok := strings.CutPrefix(x.Path, "journal:"); ok {
		if !monitorJournalUnit.MatchString(u) {
			return fmt.Errorf("log alert: invalid journal unit")
		}
	} else {
		x.Path = path.Clean(x.Path)
		if !strings.HasPrefix(x.Path, "/") || len(x.Path) > 1024 || strings.IndexFunc(x.Path, unicode.IsControl) >= 0 || aiReadDenied(x.Path) != "" {
			return fmt.Errorf("log alert: use the absolute path of a log file")
		}
	}
	if x.Name == "" {
		x.Name = truncateStr(path.Base(x.Path)+": "+x.Pattern, 80)
	}
	return nil
}

func (x monitorLogWatch) key() string {
	h := sha256.Sum256([]byte(x.Path + "\x00" + x.Pattern + "\x00" + strconv.FormatBool(x.Regex) + strconv.FormatBool(x.Case)))
	return hex.EncodeToString(h[:12])
}

func (x monitorLogWatch) grepMode() string {
	m := "F"
	if x.Regex {
		m = "E"
	}
	if !x.Case {
		m += "i"
	}
	return m
}

// The watches arrive on standard input, one per line: position, grep mode,
// path and pattern separated by tabs (none of them can contain a tab or a
// newline). Nothing user-supplied is part of the command text.
const monitorLogWatchScript = `LC_ALL=C; export LC_ALL
i=0
while IFS='	' read -r pos mode p pat; do
  i=$((i+1))
  case "$mode" in F|E|Fi|Ei) ;; *) mode=F ;; esac
  case "$pos" in ''|*[!0-9-]*) pos=-1 ;; esac
  case "$p" in
  journal:*)
    u=${p#journal:}; now=$(date +%s)
    echo "@@W $i $now"
    if [ "$pos" -gt 0 ]; then
      journalctl --no-pager -o short-iso -u "$u" --since "@$pos" --until "@$now" 2>/dev/null | head -n 200000 | grep -$mode -e "$pat" 2>/dev/null |
        awk '{n++; l[n%3]=$0} END {print n+0; for (k=n-2; k<=n; k++) if (k>0) print "> " substr(l[k%3],1,400)}'
    else echo 0; fi
    echo "@@E" ;;
  *)
    if [ ! -f "$p" ] || [ -h "$p" ] || [ ! -r "$p" ]; then echo "@@W $i -1"; echo "@@E"; continue; fi
    sz=$(wc -c < "$p" 2>/dev/null); sz=$((sz+0))
    echo "@@W $i $sz"
    if [ "$pos" -lt 0 ]; then echo 0; echo "@@E"; continue; fi
    if [ "$pos" -gt "$sz" ]; then pos=0; echo "@@R"; fi
    take=$((sz-pos))
    if [ "$take" -gt 8388608 ]; then pos=$((sz-8388608)); take=8388608; fi
    if [ "$take" -gt 0 ]; then
      tail -c +$((pos+1)) -- "$p" 2>/dev/null | head -c "$take" | grep -a -$mode -e "$pat" 2>/dev/null |
        awk '{n++; l[n%3]=$0} END {print n+0; for (k=n-2; k<=n; k++) if (k>0) print "> " substr(l[k%3],1,400)}'
    else echo 0; fi
    echo "@@E" ;;
  esac
done`

type monitorLogWatchResult struct {
	Pos     int64
	Missing bool
	Rotated bool
	Count   int
	Lines   []string
}

func parseMonitorLogWatch(out string, n int) []monitorLogWatchResult {
	res := make([]monitorLogWatchResult, n)
	seen := make([]bool, n)
	cur := -1
	gotCount := false
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case strings.HasPrefix(line, "@@W "):
			f := strings.Fields(line)
			cur = -1
			gotCount = false
			if len(f) == 3 {
				i, e1 := strconv.Atoi(f[1])
				pos, e2 := strconv.ParseInt(f[2], 10, 64)
				if e1 == nil && e2 == nil && i >= 1 && i <= n {
					cur = i - 1
					seen[cur] = true
					res[cur].Pos = pos
					res[cur].Missing = pos < 0
				}
			}
		case line == "@@R":
			if cur >= 0 {
				res[cur].Rotated = true
			}
		case line == "@@E":
			cur = -1
		case cur >= 0 && !gotCount:
			res[cur].Count, _ = strconv.Atoi(strings.TrimSpace(line))
			gotCount = true
		case cur >= 0 && strings.HasPrefix(line, "> "):
			if len(res[cur].Lines) < 3 {
				res[cur].Lines = append(res[cur].Lines, strings.ToValidUTF8(line[2:], "�"))
			}
		}
	}
	for i := range res {
		if !seen[i] {
			res[i].Pos = -2 // no answer: keep the stored position
		}
	}
	return res
}

func collectMonitorLogWatches(ctx context.Context, cl *ssh.Client, c Connection, watches []monitorLogWatch, s *monitorSample) {
	if len(watches) == 0 {
		return
	}
	type stored struct {
		pos            int64
		lastHit, lines string
		lastCount      int
		known          bool
	}
	prev := make([]stored, len(watches))
	var input strings.Builder
	for i, w := range watches {
		p := &prev[i]
		p.pos = -1
		if db.QueryRow(`SELECT pos,last_hit,last_count,last_lines FROM monitor_log_offsets WHERE conn_id=? AND watch=?`, c.ID, w.key()).Scan(&p.pos, &p.lastHit, &p.lastCount, &p.lines) == nil {
			p.known = true
		}
		fmt.Fprintf(&input, "%d\t%s\t%s\t%s\n", p.pos, w.grepMode(), w.Path, w.Pattern)
	}
	runCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	out, _ := monitorSSHCommandInput(runCtx, cl, monitorLogWatchScript, input.String())
	res := parseMonitorLogWatch(out, len(watches))
	now := time.Now().UTC()
	for i, w := range watches {
		r := res[i]
		p := prev[i]
		if r.Pos == -2 {
			continue
		}
		if r.Missing {
			s.LogAlerts = append(s.LogAlerts, monitorLogAlert{Name: w.Name, Path: w.Path, Severity: "warning", At: s.At, Lines: []string{"file missing or not readable"}})
			continue
		}
		lastHit, lastCount, lines := p.lastHit, p.lastCount, p.lines
		hit := r.Count >= w.Min && p.pos >= 0
		redacted := []string{}
		for _, l := range r.Lines {
			redacted = append(redacted, truncateStr(redactText(l), 400))
		}
		if hit {
			lastHit, lastCount, lines = s.At, r.Count, strings.Join(redacted, "\n")
			recordMonitorEvent(c, "log", w.Severity, "", w.Name+": "+strconv.Itoa(r.Count)+" new matching lines in "+w.Path, lines)
		}
		db.Exec(`INSERT INTO monitor_log_offsets(conn_id,watch,pos,last_hit,last_count,last_lines) VALUES(?,?,?,?,?,?) ON CONFLICT(conn_id,watch) DO UPDATE SET pos=excluded.pos,last_hit=excluded.last_hit,last_count=excluded.last_count,last_lines=excluded.last_lines`,
			c.ID, w.key(), r.Pos, lastHit, lastCount, truncateStr(lines, 2000))
		if lastHit == "" {
			continue
		}
		at, err := time.Parse(time.RFC3339, lastHit)
		if err != nil || now.Sub(at) > time.Duration(w.Hold)*time.Minute {
			continue
		}
		a := monitorLogAlert{Name: w.Name, Path: w.Path, Count: lastCount, Severity: w.Severity, At: lastHit, New: hit}
		if lines != "" {
			a.Lines = strings.Split(lines, "\n")
		}
		s.LogAlerts = append(s.LogAlerts, a)
		s.addReason(w.Severity, "log", w.Name, float64(lastCount), 0)
	}
}

// monitorSSHExtras runs the checks and log alerts over the poll's SSH client.
func monitorSSHExtras(ctx context.Context, cl *ssh.Client, c Connection, cfg monitorConfig, s *monitorSample) {
	if len(cfg.Services) > 0 {
		collectMonitorServices(ctx, cl, cfg, s)
	}
	if len(cfg.LogWatches) > 0 {
		collectMonitorLogWatches(ctx, cl, c, cfg.LogWatches, s)
	}
}
