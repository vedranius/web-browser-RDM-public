package main

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// ─── TERMINAL TRACKER ────────────────────────────────
//
// A small model of what an SSH terminal shows, fed with the shell's output on the server
// (v12.2.0). It is what an AI may read when the user shares the terminal, and what WRM looks
// at before it types a command into it:
//
//   - logical lines (scrollback and screen together; carriage returns, backspaces, cursor
//     moves on the current line and erase-in-line are applied, colours and other escape
//     sequences are dropped), at most termTrackMaxLines;
//   - the alternate screen buffer (?1049h / ?1047h / ?47h: editors, pagers, top …): while it
//     is active its content is not recorded and nothing is typed;
//   - the working directory from OSC 7 (file://host/path);
//   - shell integration markers OSC 133 (and VS Code's OSC 633): A prompt start, B input
//     start, C output start, D;<exit> command finished;
//   - bracketed paste mode and the time of the last output.
//
// The tracker never interprets output as instructions; its text is untrusted data.

const (
	termTrackMaxLines = 2000 // the hard cap of ai_terminal_context_lines
	termTrackMaxLine  = 2000 // runes per line
)

const (
	tsGround = iota
	tsEsc
	tsEscInter
	tsCSI
	tsOSC
	tsOSCEsc
	tsStr
	tsStrEsc
)

type termTracker struct {
	mu        sync.Mutex
	lines     []string // committed lines; lines[0] has the absolute number first
	first     int64
	cur       []rune // the line the cursor is on
	col       int
	alt       bool
	bracketed bool
	cwd       string
	lastOut   time.Time

	// shell integration (OSC 133)
	osc133   bool
	state    byte  // 0, 'A', 'B', 'C' or 'D'
	marks    int64 // number of A and D markers seen
	bAbs     int64 // line of the last B marker
	bCol     int
	cAbs     int64 // line of the last C marker (-1: none since the last B)
	lastCmd  string
	lastExit string
	outFrom  int64 // output of the last finished command: lines [outFrom, outTo)
	outTo    int64
	cmdDone  bool

	// parser
	st  int
	seq []byte
	utf []byte
}

func newTermTracker() *termTracker { return &termTracker{cAbs: -1} }

// abs is the absolute number of the current (not yet committed) line.
func (t *termTracker) abs() int64 { return t.first + int64(len(t.lines)) }

func (t *termTracker) Write(p []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(p) > 0 {
		t.lastOut = time.Now()
	}
	for _, b := range p {
		t.byte(b)
	}
}

func (t *termTracker) byte(b byte) {
	switch t.st {
	case tsGround:
		if len(t.utf) > 0 || b >= 0x80 {
			t.utf = append(t.utf, b)
			if utf8.FullRune(t.utf) {
				r, _ := utf8.DecodeRune(t.utf)
				t.utf = t.utf[:0]
				if r >= 0xa0 || r == utf8.RuneError {
					t.put(r)
				}
			} else if len(t.utf) >= 4 {
				t.utf = t.utf[:0]
				t.put(utf8.RuneError)
			}
			return
		}
		switch {
		case b == 0x1b:
			t.st = tsEsc
		case b == '\r':
			t.col = 0
		case b == '\n':
			t.newline()
		case b == '\b':
			if t.col > 0 {
				t.col--
			}
		case b == '\t':
			t.col = (t.col/8 + 1) * 8
		case b < 0x20 || b == 0x7f:
			// BEL and other controls
		default:
			t.put(rune(b))
		}
	case tsEsc:
		switch b {
		case '[':
			t.st, t.seq = tsCSI, t.seq[:0]
		case ']':
			t.st, t.seq = tsOSC, t.seq[:0]
		case 'P', 'X', '^', '_':
			t.st = tsStr
		case '(', ')', '*', '+', '-', '.', '/', '#', '%', ' ':
			t.st = tsEscInter
		case 'c': // full reset
			t.alt, t.bracketed = false, false
			t.st = tsGround
		default:
			t.st = tsGround
		}
	case tsEscInter:
		t.st = tsGround
	case tsCSI:
		switch {
		case b >= 0x40 && b <= 0x7e:
			t.csi(string(t.seq), b)
			t.st = tsGround
		case b == 0x1b:
			t.st = tsEsc
		case b < 0x20:
			// controls inside a CSI sequence are executed (rare); ignored here
		default:
			if len(t.seq) < 64 {
				t.seq = append(t.seq, b)
			}
		}
	case tsOSC:
		switch b {
		case 0x07:
			t.osc(string(t.seq))
			t.st = tsGround
		case 0x1b:
			t.st = tsOSCEsc
		default:
			if len(t.seq) < 4096 {
				t.seq = append(t.seq, b)
			}
		}
	case tsOSCEsc:
		t.osc(string(t.seq))
		t.st = tsGround
		if b != '\\' {
			t.byte(b)
		}
	case tsStr:
		if b == 0x1b {
			t.st = tsStrEsc
		} else if b == 0x07 {
			t.st = tsGround
		}
	case tsStrEsc:
		t.st = tsGround
	}
}

func (t *termTracker) put(r rune) {
	if t.alt || t.col >= termTrackMaxLine {
		return
	}
	for len(t.cur) < t.col {
		t.cur = append(t.cur, ' ')
	}
	if t.col < len(t.cur) {
		t.cur[t.col] = r
	} else {
		t.cur = append(t.cur, r)
	}
	t.col++
}

func (t *termTracker) newline() {
	if t.alt {
		return
	}
	t.lines = append(t.lines, strings.TrimRight(string(t.cur), " "))
	if n := len(t.lines) - termTrackMaxLines; n > 0 {
		t.lines = append([]string(nil), t.lines[n:]...)
		t.first += int64(n)
	}
	t.cur = t.cur[:0]
	t.col = 0
}

func csiNums(params string) []int {
	params = strings.TrimLeft(params, "?>=<")
	if params == "" {
		return nil
	}
	var out []int
	for _, p := range strings.Split(params, ";") {
		n, _ := strconv.Atoi(strings.TrimSpace(p))
		out = append(out, n)
	}
	return out
}

func (t *termTracker) csi(params string, final byte) {
	nums := csiNums(params)
	arg := func(i, def int) int {
		if i < len(nums) && nums[i] > 0 {
			return nums[i]
		}
		return def
	}
	if strings.HasPrefix(params, "?") && (final == 'h' || final == 'l') {
		for _, n := range nums {
			switch n {
			case 1049, 1047, 47:
				t.alt = final == 'h'
			case 2004:
				t.bracketed = final == 'h'
			}
		}
		return
	}
	if t.alt {
		return
	}
	switch final {
	case 'K':
		mode := 0
		if len(nums) > 0 {
			mode = nums[0]
		}
		switch mode {
		case 0:
			if t.col < len(t.cur) {
				t.cur = t.cur[:t.col]
			}
		case 1:
			for i := 0; i <= t.col && i < len(t.cur); i++ {
				t.cur[i] = ' '
			}
		case 2:
			t.cur = t.cur[:0]
		}
	case 'J':
		if len(nums) > 0 && nums[0] == 3 { // clear the scrollback ("clear")
			t.first += int64(len(t.lines))
			t.lines = nil
		}
	case 'C':
		t.col = min(t.col+arg(0, 1), termTrackMaxLine)
	case 'D':
		t.col = max(0, t.col-arg(0, 1))
	case 'G', '`':
		t.col = min(arg(0, 1)-1, termTrackMaxLine)
	case 'H', 'f':
		t.col = min(arg(1, 1)-1, termTrackMaxLine)
	case 'P':
		if n := arg(0, 1); t.col < len(t.cur) {
			end := min(len(t.cur), t.col+n)
			t.cur = append(t.cur[:t.col], t.cur[end:]...)
		}
	case '@':
		if n := min(arg(0, 1), termTrackMaxLine); t.col < len(t.cur) {
			ins := make([]rune, n)
			for i := range ins {
				ins[i] = ' '
			}
			t.cur = append(t.cur[:t.col], append(ins, t.cur[t.col:]...)...)
			if len(t.cur) > termTrackMaxLine {
				t.cur = t.cur[:termTrackMaxLine]
			}
		}
	case 'X':
		for i := t.col; i < t.col+arg(0, 1) && i < len(t.cur); i++ {
			t.cur[i] = ' '
		}
	}
}

func (t *termTracker) osc(data string) {
	code, rest, _ := strings.Cut(data, ";")
	switch code {
	case "7":
		if strings.HasPrefix(rest, "file://") {
			p := rest[len("file://"):]
			if i := strings.Index(p, "/"); i >= 0 {
				p = p[i:]
				if d, err := url.PathUnescape(p); err == nil {
					p = d
				}
				if len(p) < 4096 && !strings.ContainsAny(p, "\x00\n\r") {
					t.cwd = p
				}
			}
		}
	case "133", "633":
		if rest == "" {
			return
		}
		args := strings.Split(rest, ";")
		switch args[0] {
		case "A":
			t.osc133 = true
			t.state = 'A'
			t.marks++
		case "B":
			t.osc133 = true
			t.state = 'B'
			t.bAbs, t.bCol, t.cAbs = t.abs(), t.col, -1
		case "C":
			t.osc133 = true
			t.state = 'C'
			t.cAbs = t.abs()
			t.lastCmd = t.inputText()
		case "D":
			t.osc133 = true
			if t.state == 'B' || t.state == 'C' {
				if t.cAbs < 0 {
					t.lastCmd = t.inputText()
				}
				t.outFrom = t.cAbs
				if t.outFrom < 0 {
					t.outFrom = t.bAbs + 1
				}
				t.outTo = t.abs()
				t.cmdDone = true
			}
			t.lastExit = ""
			if len(args) > 1 {
				if _, err := strconv.Atoi(args[1]); err == nil {
					t.lastExit = args[1]
				}
			}
			t.state = 'D'
			t.marks++
		}
	}
}

// inputText is what the user typed after the last B marker (the command line).
func (t *termTracker) inputText() string {
	var line []rune
	if t.bAbs == t.abs() {
		line = t.cur
	} else if i := t.bAbs - t.first; i >= 0 && i < int64(len(t.lines)) {
		line = []rune(t.lines[i])
	}
	if t.bCol < len(line) {
		return strings.TrimSpace(string(line[t.bCol:]))
	}
	return ""
}

// linesRange returns the committed lines [from, to) that are still in the buffer.
func (t *termTracker) linesRange(from, to int64) ([]string, bool) {
	lost := from < t.first
	if from < t.first {
		from = t.first
	}
	if to > t.abs() {
		to = t.abs()
	}
	if from >= to {
		return nil, lost
	}
	return append([]string(nil), t.lines[from-t.first:to-t.first]...), lost
}

// ─── what the AI may see ─────────────────────────────

type termSnapshot struct {
	Lines     []string // the last n lines, the current line last when it is not empty
	Current   string
	Cwd       string
	Alt       bool
	OSC133    bool
	State     string // prompt | running | full_screen | password_prompt | unknown
	LastCmd   string
	LastExit  string
	LastOut   []string
	Heuristic bool // the last command was found without shell integration markers
	Quiet     time.Duration
}

var (
	// a shell prompt at the end of the current line
	rePromptEnd = regexp.MustCompile(`[$#%>❯➜»λ]$`)
	// prompts of programs that are not a shell (REPLs, database clients, continuation lines)
	reNotShell = regexp.MustCompile(`^\s*(>|>>>|\.\.\.|\(gdb\)|\(Pdb\)|\(lldb\)|mysql>|MariaDB \[[^\]]*\]>|sqlite>|ftp>|sftp>|telnet>|irb\(.*|pry\(.*|In \[\d+\]:|[A-Za-z0-9_]*\*?[=-][#>]|redis[^ ]*>|mongo[^ ]*>|>>|\?)\s*$`)
	// a prompt for a password, passphrase or one-time code
	rePassword = regexp.MustCompile(`(?i)(password|passwort|passphrase|pass phrase|passcode|lozinka|\bpin\b|one-time|\botp\b|verification code|authenticator code|token code|\[sudo\])[^\n]{0,80}[:?]\s*$|^\s*password\s*$`)
	// "prompt command" in a committed line (heuristic for the last command)
	rePromptCmd = regexp.MustCompile(`^\S.{0,160}?[$#%❯➜»] (\S.*)$`)
)

func promptLike(line string) bool {
	l := strings.TrimRight(line, " ")
	if l == "" || reNotShell.MatchString(l) {
		return false
	}
	return rePromptEnd.MatchString(l)
}

func passwordPrompt(line string) bool {
	return rePassword.MatchString(strings.TrimRight(line, " "))
}

// snapshot returns the last n lines and what WRM knows about the terminal's state.
func (t *termTracker) snapshot(n int) termSnapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	n = max(1, min(n, termTrackMaxLines))
	cur := strings.TrimRight(string(t.cur), " ")
	s := termSnapshot{Current: cur, Cwd: t.cwd, Alt: t.alt, OSC133: t.osc133, LastExit: t.lastExit}
	if !t.lastOut.IsZero() {
		s.Quiet = time.Since(t.lastOut)
	}
	all := t.lines
	if cur != "" {
		all = append(append([]string(nil), t.lines...), cur)
	}
	if len(all) > n {
		all = all[len(all)-n:]
	}
	s.Lines = append([]string(nil), all...)
	s.State = t.stateLocked()
	if t.osc133 && t.cmdDone {
		s.LastCmd = t.lastCmd
		s.LastOut, _ = t.linesRange(t.outFrom, t.outTo)
	} else if !t.osc133 {
		// heuristic: the last committed line that looks like "prompt command"
		for i := len(t.lines) - 1; i >= 0 && i >= len(t.lines)-n; i-- {
			if m := rePromptCmd.FindStringSubmatch(t.lines[i]); m != nil && !reNotShell.MatchString(t.lines[i]) {
				s.LastCmd, s.Heuristic = strings.TrimSpace(m[1]), true
				s.LastOut = append([]string(nil), t.lines[i+1:]...)
				break
			}
		}
	}
	if len(s.LastOut) > n {
		s.LastOut = s.LastOut[len(s.LastOut)-n:]
	}
	return s
}

// stateLocked describes the terminal (t.mu held).
func (t *termTracker) stateLocked() string {
	cur := string(t.cur)
	switch {
	case t.alt:
		return "full_screen"
	case passwordPrompt(cur):
		return "password_prompt"
	case t.osc133:
		switch {
		case t.state == 'C', t.state == 'B' && t.abs() > t.bAbs:
			return "running"
		case t.state == 'A', t.state == 'B':
			return "prompt"
		}
		return "unknown"
	case promptLike(cur):
		return "prompt"
	}
	return "unknown"
}

// termRefusal is why WRM does not type into the terminal now. Hard refusals do not go
// away by waiting a moment (a full-screen program, a password prompt, another program).
type termRefusal struct {
	Rule string
	Why  string
	Hard bool
	// markers: the refusal rests only on OSC 133 markers, which any program's output can
	// imitate; the foreground check can prove them stale (see termShare.ready)
	markers bool
}

func (e *termRefusal) Error() string { return e.Why }

// readyToType checks the terminal's own state (the foreground process is checked
// separately over an exec channel, see termForeground).
func (t *termTracker) readyToType(lastInput time.Time) *termRefusal {
	t.mu.Lock()
	defer t.mu.Unlock()
	cur := strings.TrimRight(string(t.cur), " ")
	switch {
	case t.alt:
		return &termRefusal{Rule: "alt_screen", Why: "a full-screen program is running in the terminal (alternate screen buffer)", Hard: true}
	case passwordPrompt(cur):
		return &termRefusal{Rule: "password_prompt", Why: "the terminal is asking for a password; WRM never types into a password prompt", Hard: true}
	}
	if !lastInput.IsZero() && time.Since(lastInput) < 1500*time.Millisecond {
		return &termRefusal{Rule: "user_typing", Why: "the user is typing in the terminal", Hard: false}
	}
	if !t.lastOut.IsZero() && time.Since(t.lastOut) < 300*time.Millisecond {
		return &termRefusal{Rule: "busy", Why: "the terminal is still printing output", Hard: false}
	}
	if t.osc133 {
		switch {
		case t.state == 'C' || (t.state == 'B' && t.abs() > t.bAbs):
			return &termRefusal{Rule: "busy", Why: "a command is still running in the terminal", Hard: true, markers: true}
		case t.state == 'B':
			if in := t.inputText(); in != "" {
				return &termRefusal{Rule: "partial_input", Why: "the prompt has text the user has not sent", Hard: true}
			}
			return nil
		case t.state == 'A':
			// a shell that marks only prompt starts (A) and command ends (D)
			return nil
		}
		// D without a new prompt yet
		return &termRefusal{Rule: "busy", Why: "the shell has not shown its prompt yet", Hard: false}
	}
	if reNotShell.MatchString(cur) && cur != "" {
		return &termRefusal{Rule: "not_shell", Why: "the terminal shows the prompt of a program that is not a shell (" + truncateStr(cur, 40) + ")", Hard: true}
	}
	if !promptLike(cur) {
		return &termRefusal{Rule: "not_shell", Why: "the terminal is not at a shell prompt (or the prompt has unsent text)", Hard: true}
	}
	return nil
}

// dropMarkers forgets the shell integration state when it is proven stale (the shell is
// in the foreground although the markers say a command runs: markers printed by a
// program, e.g. a log of another terminal); the next real marker turns it on again.
func (t *termTracker) dropMarkers() {
	t.mu.Lock()
	t.osc133, t.state, t.cmdDone, t.cAbs = false, 0, false, -1
	t.mu.Unlock()
}

// mark is the position before WRM types a command.
type termMarkPos struct {
	abs   int64
	marks int64
	at    time.Time
}

func (t *termTracker) mark() termMarkPos {
	t.mu.Lock()
	defer t.mu.Unlock()
	return termMarkPos{abs: t.abs(), marks: t.marks, at: time.Now()}
}

// commandResult reports whether the command typed at m has finished, and its output.
// done: the prompt is back (OSC 133 D/A, or a quiet prompt-like line); stop: a reason to
// stop waiting (full-screen program, password prompt); exit is "" when unknown.
func (t *termTracker) commandResult(m termMarkPos, quietFor time.Duration) (done bool, stop string, out []string, lost bool, exit string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	cur := strings.TrimRight(string(t.cur), " ")
	from := m.abs + 1 // the line the command was echoed on is skipped
	switch {
	case t.alt:
		stop = "a full-screen program started in the terminal; its output is not captured"
	case t.abs() > m.abs && passwordPrompt(cur):
		stop = "the command is asking for a password; WRM never types passwords (the user can answer in the terminal)"
	}
	if t.osc133 {
		if t.marks > m.marks && t.cmdDone && t.outTo > m.abs && (t.state == 'D' || t.state == 'A' || t.state == 'B') {
			done = true
			exit = t.lastExit
			if t.outFrom > m.abs {
				from = t.outFrom
			}
			out, lost = t.linesRange(from, t.outTo)
			return
		}
	} else if t.abs() > m.abs && promptLike(cur) && time.Since(t.lastOut) >= quietFor {
		done = true
	}
	out, lost = t.linesRange(from, t.abs())
	if stop != "" && cur != "" {
		out = append(out, cur)
	}
	return
}
