package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode/utf8"

	"github.com/pkg/sftp"
)

type monitorLogFile struct {
	Root     int    `json:"root"`
	Path     string `json:"path"`
	Service  string `json:"service"`
	Operator string `json:"operator"`
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
}
type monitorLogChunk struct {
	Text         string `json:"text"`
	Modified     int64  `json:"modified"`
	Rows         int    `json:"rows"`
	FirstTime    string `json:"first_time,omitempty"`
	LastTime     string `json:"last_time,omitempty"`
	Start        int64  `json:"start"`
	EventTime    string `json:"event_time,omitempty"`
	Offset       int64  `json:"offset"`
	Signature    string `json:"signature"`
	Size         int64  `json:"size"`
	Reset        bool   `json:"reset"`
	Truncated    bool   `json:"truncated"`
	UnknownTimes int    `json:"unknown_times"`
	Matches      int    `json:"matches"`
	Timezone     string `json:"timezone"`
	Note         string `json:"note"`
}

func openMonitorSFTP(ctx context.Context, c Connection) (*sftp.Client, func(), error) {
	cl, err := dialSSH(c, nil)
	if err != nil {
		return nil, nil, err
	}
	fs, err := sftp.NewClient(cl)
	if err != nil {
		cl.Close()
		return nil, nil, err
	}
	stop := context.AfterFunc(ctx, func() { cl.Close() })
	return fs, func() { stop(); fs.Close(); cl.Close() }, nil
}

// SFTP REALPATH implementations differ in whether they resolve symlinks. Reject
// symlink components explicitly before opening any root or log file.
func sftpNoSymlinks(fs *sftp.Client, p string) error {
	if !strings.HasPrefix(p, "/") {
		return fmt.Errorf("absolute path required")
	}
	current := "/"
	for _, part := range strings.Split(strings.TrimPrefix(path.Clean(p), "/"), "/") {
		if part == "" {
			continue
		}
		current = path.Join(current, part)
		info, err := fs.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink path components are not allowed")
		}
	}
	return nil
}
func monitorSafePath(fs *sftp.Client, root, p string) (string, error) {
	if err := sftpNoSymlinks(fs, root); err != nil {
		return "", err
	}
	if err := sftpNoSymlinks(fs, p); err != nil {
		return "", err
	}
	realRoot, err := fs.RealPath(root)
	if err != nil {
		return "", fmt.Errorf("log root unavailable")
	}
	real, err := fs.RealPath(p)
	if err != nil {
		return "", fmt.Errorf("log path unavailable")
	}
	if realRoot == "/" || !(real == realRoot || strings.HasPrefix(real, strings.TrimSuffix(realRoot, "/")+"/")) || aiReadDenied(real) != "" {
		return "", fmt.Errorf("log path escapes configured root or is sensitive")
	}
	return real, nil
}
func monitorLogsHandler(w http.ResponseWriter, r *http.Request, c Connection, acc *connAccess, rest string) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	cfg := loadMonitorConfig(c.ID)
	cfg.Roots = monitorAvailableRoots(c)
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	fs, closeFS, err := openMonitorSFTP(ctx, c)
	if err != nil {
		jsonError(w, "SSH/SFTP log access failed: "+err.Error(), 502)
		return
	}
	defer closeFS()
	if rest == "logs" {
		rootCtx, stopRoots := context.WithTimeout(ctx, 10*time.Second)
		roots, rootWarnings := monitorDiscoverRoots(rootCtx, fs, c)
		stopRoots()
		cfg.Roots = roots
		files := []monitorLogFile{}
		warnings := rootWarnings
		visited := 0
		seenFiles := map[string]bool{}
		capped := false
		type dir struct {
			name        string
			depth, root int
		}
		for root, x := range cfg.Roots {
			real, err := monitorSafePath(fs, x.Path, x.Path)
			if err != nil {
				warnings = append(warnings, x.Path+": "+err.Error())
				continue
			}
			queue := []dir{{real, 0, root}}
			for len(queue) > 0 && len(files) < 2000 && visited < 10000 {
				d := queue[0]
				queue = queue[1:]
				if ctx.Err() != nil {
					warnings = append(warnings, "Discovery time limit reached")
					capped = true
					break
				}
				entries, err := fs.ReadDir(d.name)
				if err != nil {
					warnings = append(warnings, d.name+": not readable")
					continue
				}
				for _, f := range entries {
					visited++
					if visited > 10000 || len(files) >= 2000 {
						capped = true
						break
					}
					p := path.Join(d.name, f.Name())
					if f.IsDir() {
						if d.depth < 8 {
							queue = append(queue, dir{p, d.depth + 1, root})
						} else {
							capped = true
						}
						continue
					}
					if !f.Mode().IsRegular() || !monitorLogName(f.Name()) || aiReadDenied(p) != "" {
						continue
					}
					if seenFiles[p] {
						continue
					}
					seenFiles[p] = true
					files = append(files, monitorLogFile{Root: root, Path: p, Service: x.Service, Operator: x.Operator, Size: f.Size(), Modified: f.ModTime().Unix()})
				}
			}
			if len(queue) > 0 {
				capped = true
			}
		}
		jsonOK(w, map[string]interface{}{"files": files, "roots": cfg.Roots, "warnings": warnings, "incomplete": capped, "limits": map[string]int{"files": 2000, "entries": 10000, "depth": 8}})
		return
	}
	root, err := strconv.Atoi(r.URL.Query().Get("root"))
	if rp := r.URL.Query().Get("root_path"); rp != "" {
		root = -1
		for i, x := range cfg.Roots {
			if x.Path == rp {
				root = i
				err = nil
				break
			}
		}
	}
	if err != nil || root < 0 || root >= len(cfg.Roots) {
		jsonError(w, "Choose a configured log root", 400)
		return
	}
	p := r.URL.Query().Get("path")
	if strings.ContainsAny(p, "\x00\r\n") || len(p) > 2048 {
		jsonError(w, "Invalid log path", 400)
		return
	}
	real, err := monitorSafePath(fs, cfg.Roots[root].Path, p)
	if err != nil {
		jsonError(w, err.Error(), 403)
		return
	}
	if oneOf(strings.ToLower(path.Ext(real)), ".gz", ".bz2", ".xz", ".zip", ".zst") {
		jsonError(w, "Compressed logs need to be extracted before text/live viewing; use file download", 415)
		return
	}
	f, err := fs.Open(real)
	if err != nil {
		jsonError(w, "Log is not readable by this SSH account", 403)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() {
		jsonError(w, "Log must be a regular file", 400)
		return
	}
	prefix := make([]byte, 256)
	n, _ := f.ReadAt(prefix, 0)
	sum := sha256.Sum256(prefix[:n])
	signature := hex.EncodeToString(sum[:])
	offset := -int64(1)
	if x := r.URL.Query().Get("offset"); x != "" {
		offset, err = strconv.ParseInt(x, 10, 64)
		if err != nil || offset < -1 {
			jsonError(w, "Invalid offset", 400)
			return
		}
	}
	oldSignature := r.URL.Query().Get("signature")
	previousSize, _ := strconv.ParseInt(r.URL.Query().Get("previous_size"), 10, 64)
	reset := offset > st.Size() || previousSize > st.Size()
	if oldSignature != "" && oldSignature != signature {
		changed := true
		parts := strings.SplitN(oldSignature, ":", 2)
		if len(parts) == 2 {
			oldN, _ := strconv.Atoi(parts[0])
			if oldN >= 0 && oldN <= n {
				oldSum := sha256.Sum256(prefix[:oldN])
				changed = hex.EncodeToString(oldSum[:]) != parts[1]
			}
		}
		reset = reset || changed
	}
	signature = fmt.Sprintf("%d:%s", n, signature)
	if reset {
		offset = -1
	}
	scan := r.URL.Query().Get("scan") == "1"
	capBytes := int64(64 << 10)
	if scan {
		capBytes = 8 << 20
		if offset < 0 {
			offset = 0
		}
	}
	if offset < 0 {
		capBytes = 256 << 10
		offset = max(int64(0), st.Size()-capBytes)
	}
	if offset > st.Size() {
		offset = st.Size()
	}
	b := make([]byte, int(min(capBytes, st.Size()-offset)))
	n, err = f.ReadAt(b, offset)
	if err != nil && err != io.EOF {
		jsonError(w, "Cannot read log chunk", 502)
		return
	}
	b = b[:n]
	if strings.ContainsRune(string(b), '\x00') {
		jsonError(w, "Binary logs cannot be displayed as text", 415)
		return
	}
	text := string(b)
	if scan {
		// Offsets count raw bytes, so a scan keeps the bytes as they are (invalid UTF-8,
		// such as a CP1250 letter, is replaced only in the returned lines). A segment
		// that does not reach the end of the file stops after its last complete line, so
		// the next segment starts at a line; without any newline, only a rune split at
		// the end is left for the next segment.
		if offset+int64(n) < st.Size() {
			if i := strings.LastIndexByte(text, '\n'); i >= 0 {
				text = text[:i+1]
			} else {
				text = trimSplitRune(text)
			}
		}
		n = len(text)
	} else {
		text = strings.ToValidUTF8(text, "�")
	}
	tz := r.URL.Query().Get("timezone")
	if tz == "" {
		tz = "UTC"
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		jsonError(w, "Unknown log timezone", 400)
		return
	}
	q := r.URL.Query().Get("q")
	if len(q) > 200 {
		jsonError(w, "Search text too long", 400)
		return
	}
	from, err := monitorLogTimeBound(r.URL.Query().Get("from"), location)
	if err != nil {
		jsonError(w, "Invalid from date/time", 400)
		return
	}
	to, err := monitorLogTimeBound(r.URL.Query().Get("to"), location)
	if err != nil {
		jsonError(w, "Invalid to date/time", 400)
		return
	}
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		jsonError(w, "From must not be after until", 400)
		return
	}
	chunk := monitorLogChunk{Modified: st.ModTime().Unix(), Start: offset, Offset: offset + int64(n), Signature: signature, Size: st.Size(), Reset: reset, Truncated: offset > 0 || offset+int64(n) < st.Size(), Timezone: tz, Note: "Bounded text chunk; timestamp-less lines are retained. Syslog dates use the current year in the selected timezone. Live offsets refer to raw bytes."}
	if q != "" || !from.IsZero() || !to.IsZero() {
		if scan {
			initialTime, _ := time.Parse(time.RFC3339, r.URL.Query().Get("event_time"))
			filtered, unknown, matches, consumed, lastTime := filterMonitorLogSegment(text, q, from, to, location, initialTime)
			chunk.Text = filtered
			chunk.UnknownTimes = unknown
			chunk.Matches = matches
			chunk.Offset = offset + int64(consumed)
			if !lastTime.IsZero() {
				chunk.EventTime = lastTime.Format(time.RFC3339)
			}
			if consumed < len(text) {
				chunk.Truncated = true
			}
		} else {
			chunk.Text, chunk.UnknownTimes, chunk.Matches = filterMonitorLog(text, q, from, to, location)
		}
	} else {
		if scan {
			filtered, unknown, matches, consumed, lastTime := filterMonitorLogSegment(text, q, from, to, location, time.Time{})
			chunk.Text = filtered
			chunk.UnknownTimes = unknown
			chunk.Matches = matches
			chunk.Offset = offset + int64(consumed)
			if !lastTime.IsZero() {
				chunk.EventTime = lastTime.Format(time.RFC3339)
			}
			if consumed < len(text) {
				chunk.Truncated = true
			}
		} else {
			chunk.Text = text
		}
	}
	lines, _ := strconv.Atoi(r.URL.Query().Get("lines"))
	if lines > 0 {
		parts := strings.Split(chunk.Text, "\n")
		keep := min(lines, 2000)
		if len(parts) > keep {
			chunk.Text = strings.Join(parts[len(parts)-keep:], "\n")
			chunk.Truncated = true
		}
	}
	for _, line := range strings.Split(chunk.Text, "\n") {
		if strings.TrimSpace(line) != "" {
			chunk.Rows++
		}
	}
	if q == "" {
		chunk.Matches = chunk.Rows
	}
	initial, _ := time.Parse(time.RFC3339, r.URL.Query().Get("event_time"))
	if !scan {
		initial = time.Time{}
	}
	chunk.FirstTime, chunk.LastTime = monitorEventRange(text[:min(len(text), int(chunk.Offset-offset))], q, from, to, location, initial)

	if r.URL.Query().Get("download") == "1" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Disposition", `attachment; filename="log-excerpt.txt"`)
		io.WriteString(w, chunk.Text)
		return
	}
	jsonOK(w, chunk)
}

var monitorISOTime = regexp.MustCompile(`^\s*(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?)`)

func parseMonitorLogTime(line string, loc *time.Location) (time.Time, bool) {
	if m := monitorISOTime.FindStringSubmatch(line); len(m) > 1 {
		v := strings.Replace(m[1], " ", "T", 1)
		for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999999-0700", "2006-01-02T15:04:05.999999999"} {
			if at, err := time.ParseInLocation(layout, v, loc); err == nil {
				return at, true
			}
		}
	}
	if len(line) >= 15 {
		if at, err := time.ParseInLocation("2006 Jan _2 15:04:05", fmt.Sprintf("%d %s", time.Now().In(loc).Year(), line[:15]), loc); err == nil {
			return at, true
		}
	}
	return time.Time{}, false
}
func monitorLogTimeBound(v string, loc *time.Location) (time.Time, error) {
	if v == "" {
		return time.Time{}, nil
	}
	for _, format := range []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02T15:04:05"} {
		if at, err := time.ParseInLocation(format, v, loc); err == nil {
			return at, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid timestamp")
}
func filterMonitorLog(text, q string, from, to time.Time, loc *time.Location) (string, int, int) {
	lines := strings.Split(text, "\n")
	out := []string{}
	unknown, matches := 0, 0
	var eventTime time.Time
	for _, line := range lines {
		at, ok := parseMonitorLogTime(line, loc)
		if ok {
			eventTime = at
		} else {
			unknown++
		}
		if !eventTime.IsZero() && ((!from.IsZero() && eventTime.Before(from)) || (!to.IsZero() && eventTime.After(to))) {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(q)) {
			continue
		}
		if q != "" {
			matches += strings.Count(strings.ToLower(line), strings.ToLower(q))
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n"), unknown, matches
}

// trimSplitRune drops an incomplete UTF-8 sequence at the end of text (at most 3 bytes).
func trimSplitRune(text string) string {
	for i := len(text) - 1; i >= 0 && i >= len(text)-utf8.UTFMax; i-- {
		if utf8.RuneStart(text[i]) {
			if !utf8.FullRuneInString(text[i:]) {
				return text[:i]
			}
			break
		}
	}
	return text
}

// Search advances a bounded byte segment. Stop before an undisplayed matching
// line when the output cap is reached, so Next segment can resume at that line.
func filterMonitorLogSegment(text, q string, from, to time.Time, loc *time.Location, eventTime time.Time) (string, int, int, int, time.Time) {
	var out strings.Builder
	unknown, matches, consumed, rows := 0, 0, 0, 0
	for _, raw := range strings.SplitAfter(text, "\n") {
		line := strings.TrimSuffix(raw, "\n")
		at, known := parseMonitorLogTime(line, loc)
		nextTime := eventTime
		if known {
			nextTime = at
		}
		keep := nextTime.IsZero() || ((from.IsZero() || !nextTime.Before(from)) && (to.IsZero() || !nextTime.After(to)))
		if q != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(q)) {
			keep = false
		}
		if keep && (out.Len()+len(raw) > 256<<10 || rows >= 2000) {
			break
		}
		eventTime = nextTime
		consumed += len(raw)
		if !known {
			unknown++
		}
		if keep {
			out.WriteString(strings.ToValidUTF8(raw, "�"))
			rows++
			if q != "" {
				matches += strings.Count(strings.ToLower(line), strings.ToLower(q))
			}
		}
	}
	return out.String(), unknown, matches, consumed, eventTime
}

// Matching continuation lines inherit their event timestamp even when the
// keyword appears only on the continuation, not on the timestamped header.
func monitorEventRange(text, q string, from, to time.Time, loc *time.Location, at time.Time) (string, string) {
	first, last := "", ""
	for _, line := range strings.Split(text, "\n") {
		if v, ok := parseMonitorLogTime(line, loc); ok {
			at = v
		}
		if at.IsZero() || strings.TrimSpace(line) == "" || (!from.IsZero() && at.Before(from)) || (!to.IsZero() && at.After(to)) || (q != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(q))) {
			continue
		}
		if first == "" {
			first = at.Format(time.RFC3339Nano)
		}
		last = at.Format(time.RFC3339Nano)
	}
	return first, last
}
