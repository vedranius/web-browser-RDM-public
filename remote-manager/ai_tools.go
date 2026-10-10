package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ─── AI: TOOLS ───────────────────────────────────────
//
// The model never gets a shell. It can only call these tools; each call becomes an
// aiAction that the permission engine evaluates before WRM runs anything over its own SSH
// connection. Read tools run fixed scripts with the arguments passed as shell-quoted
// positional parameters; run_command runs exactly the command that was evaluated (or
// approved, or edited by the user); file edits are shown as a unified diff and written
// only if the file did not change since the diff was made.

type aiToolSpec struct {
	Def    aiToolDef
	Kind   string // read | command | write
	Modes  []string
	Expose bool // offered to the model
}

func obj(props map[string]interface{}, required ...string) map[string]interface{} {
	if required == nil {
		required = []string{}
	}
	return map[string]interface{}{"type": "object", "properties": props, "required": required, "additionalProperties": false}
}

func strProp(desc string) map[string]interface{} {
	return map[string]interface{}{"type": "string", "description": desc}
}

func intProp(desc string) map[string]interface{} {
	return map[string]interface{}{"type": "integer", "description": desc}
}

var aiTools = []aiToolSpec{
	{Kind: "read", Def: aiToolDef{Name: "system_info", Description: "Show the operating system, kernel, uptime, memory and CPU count of the server.", Schema: obj(map[string]interface{}{})}},
	{Kind: "read", Def: aiToolDef{Name: "read_file", Description: "Read a text file on the server (size-limited). Sensitive files (password hashes, private keys, credentials) are refused.",
		Schema: obj(map[string]interface{}{"path": strProp("Absolute path, or ~/ for the home directory"), "offset": intProp("Byte offset to start at (default 0)"), "max_bytes": intProp("Bytes to read (default and maximum: the policy limit)")}, "path")}},
	{Kind: "read", Def: aiToolDef{Name: "list_directory", Description: "List a directory with sizes, owners, permissions and times.",
		Schema: obj(map[string]interface{}{"path": strProp("Absolute path, or ~/")}, "path")}},
	{Kind: "read", Def: aiToolDef{Name: "tail_log", Description: "Show the last lines of a log file.",
		Schema: obj(map[string]interface{}{"path": strProp("Absolute path of the log file"), "lines": intProp("Number of lines (1-2000, default 200)")}, "path")}},
	{Kind: "read", Def: aiToolDef{Name: "journal", Description: "Read the systemd journal (journalctl), optionally for one unit, since a time, at a priority, matching a pattern.",
		Schema: obj(map[string]interface{}{"unit": strProp("systemd unit, e.g. nginx.service"), "lines": intProp("Number of lines (1-2000, default 200)"),
			"since": strProp("e.g. '1 hour ago', 'today', '2025-01-31 10:00'"), "priority": strProp("emerg, alert, crit, err, warning, notice, info or debug"), "grep": strProp("Only lines matching this pattern")})}},
	{Kind: "read", Def: aiToolDef{Name: "service_status", Description: "Show the status of a service, or list all services when no name is given.",
		Schema: obj(map[string]interface{}{"name": strProp("Service name, e.g. nginx")})}},
	{Kind: "read", Def: aiToolDef{Name: "list_packages", Description: "List installed packages (dpkg, rpm, apk or pacman), optionally filtered.",
		Schema: obj(map[string]interface{}{"filter": strProp("Only packages whose line contains this text")})}},
	{Kind: "read", Def: aiToolDef{Name: "processes", Description: "List the processes using the most CPU or memory.",
		Schema: obj(map[string]interface{}{"sort": strProp("cpu (default) or mem"), "limit": intProp("Number of processes (default 25)")})}},
	{Kind: "read", Def: aiToolDef{Name: "disk_usage", Description: "Show file system usage, and the largest entries of a directory when a path is given.",
		Schema: obj(map[string]interface{}{"path": strProp("Directory to summarise (optional)")})}},
	{Kind: "read", Def: aiToolDef{Name: "network_status", Description: "Show addresses, routes and listening sockets.", Schema: obj(map[string]interface{}{})}},
	{Kind: "command", Def: aiToolDef{Name: "run_command", Description: "Run one shell command on the server. Read-only commands run directly; commands that change the system need the user's approval (or the session's automatic-mode allow list) and are refused in read-only mode. Destructive commands are always refused. Non-interactive only: no pagers, editors or password prompts (use sudo -n).",
		Schema: obj(map[string]interface{}{"command": strProp("The exact command line"), "reason": strProp("One sentence for the user: why this command")}, "command", "reason")}},
	{Kind: "write", Def: aiToolDef{Name: "write_file", Description: "Create or replace a text file. The user sees a unified diff and must approve it (unless the session's automatic mode allows the path).",
		Schema: obj(map[string]interface{}{"path": strProp("Absolute path, or ~/"), "content": strProp("The complete new content"), "reason": strProp("One sentence for the user: why this change")}, "path", "content", "reason")}},
	{Kind: "write", Def: aiToolDef{Name: "edit_file", Description: "Replace one exact occurrence of a text in a file. The user sees a unified diff and must approve it (unless the session's automatic mode allows the path).",
		Schema: obj(map[string]interface{}{"path": strProp("Absolute path, or ~/"), "old_text": strProp("The exact text to replace; it must occur exactly once"),
			"new_text": strProp("The replacement"), "reason": strProp("One sentence for the user: why this change")}, "path", "old_text", "new_text", "reason")}},
}

func aiToolSpecFor(name string) (aiToolSpec, bool) {
	for _, t := range aiTools {
		if t.Def.Name == name {
			return t, true
		}
	}
	return aiToolSpec{}, false
}

// aiToolDefsFor lists the tools offered to the model in a mode: read-only sessions are
// not offered the write tools at all (run_command stays: read-only commands are allowed).
func aiToolDefsFor(mode string) []aiToolDef {
	var out []aiToolDef
	for _, t := range aiTools {
		if t.Kind == "write" && mode == aiModeReadOnly {
			continue
		}
		out = append(out, t.Def)
	}
	return out
}

// aiToolInput is the union of every tool's arguments.
type aiToolInput struct {
	Path     string `json:"path"`
	Offset   int64  `json:"offset"`
	MaxBytes int64  `json:"max_bytes"`
	Lines    int    `json:"lines"`
	Unit     string `json:"unit"`
	Since    string `json:"since"`
	Priority string `json:"priority"`
	Grep     string `json:"grep"`
	Name     string `json:"name"`
	Filter   string `json:"filter"`
	Sort     string `json:"sort"`
	Limit    int    `json:"limit"`
	Command  string `json:"command"`
	Reason   string `json:"reason"`
	Content  string `json:"content"`
	OldText  string `json:"old_text"`
	NewText  string `json:"new_text"`
}

var (
	reUnit     = regexp.MustCompile(`^[A-Za-z0-9@._:\\-]{1,200}$`)
	reSince    = regexp.MustCompile(`^[A-Za-z0-9 :.+-]{1,60}$`)
	rePriority = regexp.MustCompile(`^(emerg|alert|crit|err|warning|notice|info|debug|[0-7])(\.\.(emerg|alert|crit|err|warning|notice|info|debug|[0-7]))?$`)
)

func aiClamp(v, def, lo, hi int) int {
	if v == 0 {
		v = def
	}
	if v < lo {
		v = lo
	}
	if v > hi {
		v = hi
	}
	return v
}

// aiPlanned is a tool call turned into something to evaluate and run.
type aiPlanned struct {
	Act     aiAction
	Script  string // command line to run (read tools and run_command)
	Display string // what the user sees (command or tool summary)
	// writes
	Write     bool
	Path      string
	Old       string
	OldHash   string // "new" when the file does not exist
	NewText   string
	Diff      string
	maxOutput int
}

func aiRequirePath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", fmt.Errorf("a path is required")
	}
	if !strings.HasPrefix(p, "/") && p != "~" && !strings.HasPrefix(p, "~/") {
		return "", fmt.Errorf("use an absolute path (or ~/…)")
	}
	if strings.ContainsAny(p, "\x00\n\r") || len(p) > 1024 {
		return "", fmt.Errorf("invalid path")
	}
	return p, nil
}

// aiReadMax is the read limit for files (bytes).
func aiReadMax() int64 { return int64(settingInt("ai_output_max_kb")) << 10 }

// plan turns a tool call into an action. File edits read the current file first (through
// a read-only script) to build the diff.
func (s *aiSession) plan(ctx context.Context, name string, in aiToolInput) (*aiPlanned, error) {
	q := shellQuote
	switch name {
	case "system_info":
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "system information",
			Script: `uname -a; echo; cat /etc/os-release 2>/dev/null; echo; uptime; echo; free -m 2>/dev/null; echo; echo "CPUs: $(nproc 2>/dev/null)"`}, nil
	case "read_file", "tail_log":
		p, err := aiRequirePath(in.Path)
		if err != nil {
			return nil, err
		}
		p = aiExpandHome(p, s.home)
		var body string
		if name == "read_file" {
			max := aiReadMax()
			n := in.MaxBytes
			if n <= 0 || n > max {
				n = max
			}
			off := in.Offset
			if off < 0 {
				off = 0
			}
			body = `tail -c +` + strconv.FormatInt(off+1, 10) + ` -- "$r" | head -c ` + strconv.FormatInt(n, 10)
		} else {
			body = `tail -n ` + strconv.Itoa(aiClamp(in.Lines, 200, 1, 2000)) + ` -- "$r"`
		}
		// the real path is printed first so a symlink to a sensitive file is caught
		script := `f=` + q(p) + `; r=$(readlink -f -- "$f" 2>/dev/null || printf %s "$f"); printf 'WRM-REAL %s\n' "$r"; ` +
			`[ -e "$r" ] || { echo "No such file: $f"; exit 2; }; [ -d "$r" ] && { echo "Is a directory: $f"; exit 2; }; ` +
			`printf 'WRM-SIZE %s\n' "$(wc -c < "$r" 2>/dev/null)"; ` + body
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read", Path: p}, Display: name + " " + p, Script: script}, nil
	case "list_directory":
		p, err := aiRequirePath(in.Path)
		if err != nil {
			return nil, err
		}
		p = aiExpandHome(p, s.home)
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "ls -la " + p,
			Script: `ls -la --time-style=long-iso -- ` + q(p) + ` 2>/dev/null || ls -la -- ` + q(p)}, nil
	case "journal":
		args := []string{"journalctl", "--no-pager", "-o", "short-iso", "-n", strconv.Itoa(aiClamp(in.Lines, 200, 1, 2000))}
		if in.Unit != "" {
			if !reUnit.MatchString(in.Unit) {
				return nil, fmt.Errorf("invalid unit name")
			}
			args = append(args, "-u", q(in.Unit))
		}
		if in.Since != "" {
			if !reSince.MatchString(in.Since) {
				return nil, fmt.Errorf("invalid since value")
			}
			args = append(args, "--since", q(in.Since))
		}
		if in.Priority != "" {
			if !rePriority.MatchString(in.Priority) {
				return nil, fmt.Errorf("invalid priority")
			}
			args = append(args, "-p", in.Priority)
		}
		if in.Grep != "" {
			if len(in.Grep) > 200 || strings.ContainsAny(in.Grep, "\x00\n") {
				return nil, fmt.Errorf("invalid pattern")
			}
			args = append(args, "-g", q(in.Grep))
		}
		cmd := strings.Join(args, " ")
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: cmd, Script: cmd + " 2>&1"}, nil
	case "service_status":
		if in.Name == "" {
			return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "list services",
				Script: `systemctl list-units --type=service --all --no-pager --no-legend 2>/dev/null || service --status-all 2>&1 || rc-status -a 2>&1`}, nil
		}
		if !reUnit.MatchString(in.Name) {
			return nil, fmt.Errorf("invalid service name")
		}
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "systemctl status " + in.Name,
			Script: `systemctl status --no-pager --lines=30 -- ` + q(in.Name) + ` 2>&1 || service ` + q(in.Name) + ` status 2>&1`}, nil
	case "list_packages":
		filter := `cat`
		if in.Filter != "" {
			if len(in.Filter) > 100 || strings.ContainsAny(in.Filter, "\x00\n") {
				return nil, fmt.Errorf("invalid filter")
			}
			filter = `grep -i -F -- ` + q(in.Filter)
		}
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "list packages " + in.Filter,
			Script: `{ if command -v dpkg-query >/dev/null 2>&1; then dpkg-query -W -f='${Package} ${Version} ${Status}\n'; ` +
				`elif command -v rpm >/dev/null 2>&1; then rpm -qa --qf '%{NAME} %{VERSION}-%{RELEASE}\n'; ` +
				`elif command -v apk >/dev/null 2>&1; then apk info -v; elif command -v pacman >/dev/null 2>&1; then pacman -Q; ` +
				`else echo "no known package manager"; fi; } 2>/dev/null | ` + filter + ` | head -n 1000`}, nil
	case "processes":
		key := "-%cpu"
		if in.Sort == "mem" {
			key = "-%mem"
		}
		n := aiClamp(in.Limit, 25, 1, 200)
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "ps aux --sort=" + key,
			Script: `ps aux --sort=` + key + ` 2>/dev/null | head -n ` + strconv.Itoa(n+1) + ` || ps aux | head -n ` + strconv.Itoa(n+1)}, nil
	case "disk_usage":
		script := `df -hP 2>/dev/null || df -h`
		disp := "df -h"
		if in.Path != "" {
			p, err := aiRequirePath(in.Path)
			if err != nil {
				return nil, err
			}
			p = aiExpandHome(p, s.home)
			script += `; echo; du -xsh -- ` + q(p) + `/* ` + q(p) + `/.[!.]* 2>/dev/null | sort -h | tail -n 25`
			disp += "; du -xsh " + p + "/*"
		}
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: disp, Script: script}, nil
	case "network_status":
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "read"}, Display: "addresses, routes, listening sockets",
			Script: `(ip -brief addr 2>/dev/null || ifconfig -a 2>/dev/null); echo; (ip route 2>/dev/null || netstat -rn 2>/dev/null); echo; (ss -tulpn 2>/dev/null || netstat -tuln 2>/dev/null)`}, nil
	case "run_command":
		cmd := strings.TrimSpace(in.Command)
		if cmd == "" {
			return nil, fmt.Errorf("the command is empty")
		}
		if len(cmd) > 8000 || strings.ContainsRune(cmd, 0) {
			return nil, fmt.Errorf("the command is too long")
		}
		return &aiPlanned{Act: aiAction{Tool: name, Kind: "command", Command: cmd}, Display: cmd, Script: cmd}, nil
	case "write_file", "edit_file":
		p, err := aiRequirePath(in.Path)
		if err != nil {
			return nil, err
		}
		p = aiExpandHome(p, s.home)
		old, hash, real, err := s.readForEdit(ctx, p)
		if err != nil {
			return nil, err
		}
		newText := in.Content
		if name == "edit_file" {
			if hash == "new" {
				return nil, fmt.Errorf("the file does not exist (use write_file to create it)")
			}
			if in.OldText == "" {
				return nil, fmt.Errorf("old_text is empty")
			}
			switch strings.Count(old, in.OldText) {
			case 0:
				return nil, fmt.Errorf("old_text was not found in the file")
			case 1:
			default:
				return nil, fmt.Errorf("old_text occurs more than once; include more context")
			}
			newText = strings.Replace(old, in.OldText, in.NewText, 1)
		}
		if len(newText) > 2<<20 {
			return nil, fmt.Errorf("the new content is larger than 2 MB")
		}
		pl := &aiPlanned{Act: aiAction{Tool: name, Kind: "write", Path: real}, Write: true, Path: p, Old: old, OldHash: hash, NewText: newText}
		pl.Diff = aiUnifiedDiff(p, old, newText, hash == "new")
		pl.Display = "edit " + p
		if real != p {
			pl.Display += " (→ " + real + ")"
		}
		return pl, nil
	}
	return nil, fmt.Errorf("unknown tool %q", name)
}

// readForEdit reads a file before an edit: content, SHA-256 ("new" when missing) and the
// real path behind symlinks.
func (s *aiSession) readForEdit(ctx context.Context, p string) (string, string, string, error) {
	script := `f=` + shellQuote(p) + `; r=$(readlink -f -- "$f" 2>/dev/null || printf %s "$f"); printf 'WRM-REAL %s\n' "$r"; ` +
		`if [ ! -e "$r" ]; then echo WRM-NEW; exit 0; fi; [ -f "$r" ] || { echo "WRM-ERR not a regular file"; exit 0; }; ` +
		`n=$(wc -c < "$r"); [ "$n" -le 2097152 ] || { echo "WRM-ERR larger than 2 MB"; exit 0; }; echo WRM-BEGIN; cat -- "$r"`
	out, code, err := s.exec(ctx, script, "", 0)
	if err != nil {
		return "", "", "", err
	}
	if code != 0 {
		return "", "", "", fmt.Errorf("cannot read the file (exit %d)", code)
	}
	real := p
	if strings.HasPrefix(out, "WRM-REAL ") {
		i := strings.Index(out, "\n")
		if i < 0 {
			return "", "", "", fmt.Errorf("unexpected answer from the server")
		}
		real = out[len("WRM-REAL "):i]
		out = out[i+1:]
	}
	switch {
	case strings.HasPrefix(out, "WRM-NEW"):
		return "", "new", real, nil
	case strings.HasPrefix(out, "WRM-ERR "):
		return "", "", "", fmt.Errorf("%s", strings.TrimSpace(out[8:]))
	case strings.HasPrefix(out, "WRM-BEGIN\n"):
		content := out[len("WRM-BEGIN\n"):]
		if why := aiReadDenied(real); why != "" {
			return "", "", "", fmt.Errorf("the file is not readable by the assistant (%s)", why)
		}
		sum := sha256.Sum256([]byte(content))
		return content, hex.EncodeToString(sum[:]), real, nil
	}
	return "", "", "", fmt.Errorf("unexpected answer from the server")
}

// aiWriteScript writes stdin to the file if it still has the hash the diff was made from.
func aiWriteScript(p, oldHash string) string {
	return `f=` + shellQuote(p) + `; want=` + shellQuote(oldHash) + `; ` +
		`if [ "$want" = new ]; then [ -e "$f" ] && { echo "The file was created meanwhile; not written."; exit 3; }; ` +
		`d=$(dirname -- "$f"); [ -d "$d" ] || { echo "No such directory: $d"; exit 2; }; ` +
		`else cur=$( (sha256sum -- "$f" 2>/dev/null || shasum -a 256 -- "$f" 2>/dev/null || openssl dgst -sha256 -r -- "$f") | cut -d' ' -f1); ` +
		`[ "$cur" = "$want" ] || { echo "The file changed since the diff was made; not written."; exit 3; }; fi; ` +
		`cat > "$f" && echo "Written: $f"`
}

// aiUnifiedDiff renders a unified diff of a file change.
func aiUnifiedDiff(p, a, b string, created bool) string {
	var sb strings.Builder
	from := "a" + p
	if created {
		from = "/dev/null"
	}
	sb.WriteString("--- " + from + "\n+++ b" + p + "\n")
	lines, _ := unifiedDiff(splitLines(a), splitLines(b), 3)
	if len(lines) == 0 {
		sb.WriteString("(no changes)\n")
		return sb.String()
	}
	for _, l := range lines {
		switch l.Op {
		case "@":
			fmt.Fprintf(&sb, "@@ -%d +%d @@\n", l.A, l.B)
		default:
			sb.WriteString(l.Op + l.S + "\n")
		}
	}
	return sb.String()
}

// decodeToolInput parses a model's arguments strictly (unknown fields are refused).
func decodeToolInput(raw json.RawMessage) (aiToolInput, error) {
	var in aiToolInput
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return in, fmt.Errorf("invalid arguments: %v", err)
	}
	return in, nil
}
