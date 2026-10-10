package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ─── AI: PERMISSION MODES AND POLICIES ───────────────
//
// Every AI session has one mode:
//   read_only – a fixed set of read tools and commands the read-only classifier accepts;
//   ask       – reads run, every change waits for an approval of that exact item;
//   auto      – after an explicit opt-in: changes that match the session's allow list (and
//               not its deny list) run without asking until the time or action limit is
//               reached; everything else still asks.
// Destructive patterns are refused in every mode (ai_destructive.go).
//
// Which modes are allowed is decided on the server: the global ai_modes policy, then
// every ai_mode_rules entry that matches the connection (tag, folder, connection, host).
// The most restrictive answer wins — the allowed set is the intersection.
//
// The evaluation (aiEvaluate) is independent of how the session is driven (the built-in
// panel today, MCP clients in v12.1.0): the transport only shows approvals and passes
// the user's decision back.

const (
	aiModeReadOnly = "read_only"
	aiModeAsk      = "ask"
	aiModeAuto     = "auto"

	aiAllProviderKinds = "anthropic,openai,azure,bedrock,vertex,openai_compatible"
)

var aiModeOrder = []string{aiModeReadOnly, aiModeAsk, aiModeAuto}

func aiModeRank(m string) int {
	for i, x := range aiModeOrder {
		if x == m {
			return i
		}
	}
	return -1
}

func normalizeModeList(v string) ([]string, error) {
	seen := map[string]bool{}
	for _, x := range splitPatterns(strings.ToLower(v)) {
		x = strings.ReplaceAll(strings.ReplaceAll(x, "-", "_"), " ", "_")
		if x == "readonly" {
			x = aiModeReadOnly
		}
		if aiModeRank(x) < 0 {
			return nil, fmt.Errorf("unknown mode %q (use read_only, ask, auto)", x)
		}
		seen[x] = true
	}
	var out []string
	for _, m := range aiModeOrder {
		if seen[m] {
			out = append(out, m)
		}
	}
	return out, nil
}

func normalizeAIModes(v string) (string, error) {
	l, err := normalizeModeList(v)
	if err != nil {
		return "", err
	}
	return strings.Join(l, ","), nil
}

func normalizeAIKinds(v string) (string, error) {
	seen := map[string]bool{}
	for _, x := range splitPatterns(strings.ToLower(v)) {
		if !oneOf(x, strings.Split(aiAllProviderKinds, ",")...) {
			return "", fmt.Errorf("unknown provider kind %q (use %s)", x, aiAllProviderKinds)
		}
		seen[x] = true
	}
	var out []string
	for _, k := range strings.Split(aiAllProviderKinds, ",") {
		if seen[k] {
			out = append(out, k)
		}
	}
	return strings.Join(out, ","), nil
}

// aiModeRule limits the modes of matching connections.
type aiModeRule struct {
	Match string   `json:"match"` // tag | folder | connection | host
	Value string   `json:"value"`
	Modes []string `json:"modes"`
}

func parseAIModeRules(v string) ([]aiModeRule, error) {
	var rules []aiModeRule
	if strings.TrimSpace(v) == "" {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(v), &rules); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}
	if len(rules) > 200 {
		return nil, fmt.Errorf("at most 200 rules")
	}
	for i := range rules {
		r := &rules[i]
		r.Match = strings.ToLower(strings.TrimSpace(r.Match))
		r.Value = strings.TrimSpace(r.Value)
		if !oneOf(r.Match, "tag", "folder", "connection", "host") {
			return nil, fmt.Errorf("rule %d: match must be tag, folder, connection or host", i+1)
		}
		if r.Value == "" || len(r.Value) > 200 {
			return nil, fmt.Errorf("rule %d: value is required", i+1)
		}
		if r.Match == "connection" {
			if n, err := strconv.Atoi(r.Value); err != nil || n <= 0 {
				return nil, fmt.Errorf("rule %d: connection must be a connection id", i+1)
			}
		}
		if r.Match == "tag" {
			r.Value = normalizeTag(r.Value)
		}
		modes, err := normalizeModeList(strings.Join(r.Modes, ","))
		if err != nil {
			return nil, fmt.Errorf("rule %d: %v", i+1, err)
		}
		r.Modes = modes
		if r.Modes == nil {
			r.Modes = []string{}
		}
	}
	return rules, nil
}

func validateAIModeRules(v string) (string, error) {
	rules, err := parseAIModeRules(v)
	if err != nil {
		return "", err
	}
	if rules == nil {
		rules = []aiModeRule{}
	}
	b, _ := json.Marshal(rules)
	return string(b), nil
}

// aiConnFacts are the connection properties the mode rules look at.
type aiConnFacts struct {
	ID     int
	Host   string
	Tags   []string
	Folder string
}

func aiFactsOf(connID int) aiConnFacts {
	f := aiConnFacts{ID: connID}
	var tags, folder string
	db.QueryRow(`SELECT c.host, c.tags, COALESCE(f.name,'') FROM connections c LEFT JOIN folders f ON f.id=c.folder_id WHERE c.id=?`, connID).
		Scan(&f.Host, &tags, &folder)
	f.Host = hostOnly(f.Host)
	f.Tags = parseTags(tags)
	f.Folder = folder
	return f
}

func (r aiModeRule) matches(f aiConnFacts) bool {
	switch r.Match {
	case "tag":
		for _, t := range f.Tags {
			if t == r.Value || globMatch(r.Value, t, false) {
				return true
			}
		}
	case "folder":
		return f.Folder != "" && (strings.EqualFold(f.Folder, r.Value) || globMatch(strings.ToLower(r.Value), strings.ToLower(f.Folder), false))
	case "connection":
		return strconv.Itoa(f.ID) == r.Value
	case "host":
		return globMatch(strings.ToLower(r.Value), strings.ToLower(f.Host), false)
	}
	return false
}

// aiAllowedModes returns the modes allowed for a connection (most restrictive wins) and
// the rules that limited them.
func aiAllowedModes(f aiConnFacts) ([]string, []string) {
	allowed, _ := normalizeModeList(getSetting("ai_modes"))
	set := map[string]bool{}
	for _, m := range allowed {
		set[m] = true
	}
	var why []string
	rules, _ := parseAIModeRules(getSetting("ai_mode_rules"))
	for _, r := range rules {
		if !r.matches(f) {
			continue
		}
		in := map[string]bool{}
		for _, m := range r.Modes {
			in[m] = true
		}
		for m := range set {
			if !in[m] {
				delete(set, m)
			}
		}
		why = append(why, r.Match+" "+r.Value+": "+strings.Join(r.Modes, ","))
	}
	var out []string
	for _, m := range aiModeOrder {
		if set[m] {
			out = append(out, m)
		}
	}
	return out, why
}

// aiDefaultMode picks the policy default when it is allowed, else the most restrictive.
func aiDefaultMode(allowed []string) string {
	if d := getSetting("ai_default_mode"); oneOf(d, allowed...) {
		return d
	}
	if len(allowed) > 0 {
		return allowed[0]
	}
	return ""
}

// aiUserAllowed reports whether the user may use the assistant at all.
func aiUserAllowed(userID int) (bool, string) {
	if userID <= 0 {
		return false, "sign in to use the AI assistant"
	}
	if settingBool("ai_kill_switch") {
		return false, "the AI assistant is stopped by the administrator (kill switch)"
	}
	switch getSetting("ai_assistant") {
	case "off":
		return false, "the AI assistant is turned off"
	case "admins":
		if !isAdminUser(userID) {
			return false, "the AI assistant is limited to administrators"
		}
	}
	var blocked int
	db.QueryRow(`SELECT COALESCE(ai_blocked,0) FROM users WHERE id=?`, userID).Scan(&blocked)
	if blocked == 1 {
		return false, "an administrator turned the AI assistant off for your account"
	}
	return true, ""
}

func aiPersonalKeysAllowed(userID int) bool {
	switch getSetting("ai_personal_keys") {
	case "all":
		return true
	case "admins":
		return isAdminUser(userID)
	}
	return false
}

func aiKindAllowed(kind string) bool {
	return oneOf(kind, splitPatterns(getSetting("ai_provider_kinds"))...)
}

// aiModelAllowed checks a model against the ai_models policy (glob patterns; empty = any).
func aiModelAllowed(model string) bool {
	pats := splitPatterns(getSetting("ai_models"))
	if len(pats) == 0 {
		return true
	}
	for _, p := range pats {
		if globMatch(p, model, false) {
			return true
		}
	}
	return false
}

// ─── automatic mode limits ───────────────────────────

type aiAutoLimits struct {
	Allow      []string  `json:"allow"`       // command patterns that may run without asking
	Deny       []string  `json:"deny"`        // command patterns that are refused
	PathAllow  []string  `json:"path_allow"`  // files that may be written without asking
	PathDeny   []string  `json:"path_deny"`   // files that are never written
	Minutes    int       `json:"minutes"`     // time limit
	MaxActions int       `json:"max_actions"` // 0 = up to the policy limit
	Until      time.Time `json:"until"`
	Used       int       `json:"used"`
}

func cleanPatternList(in []string, max int) ([]string, error) {
	var out []string
	for _, x := range in {
		x = strings.TrimSpace(x)
		if x == "" {
			continue
		}
		if len(x) > 300 {
			return nil, fmt.Errorf("pattern too long: %s…", x[:40])
		}
		out = append(out, x)
	}
	if len(out) > max {
		return nil, fmt.Errorf("at most %d patterns", max)
	}
	sort.Strings(out)
	return out, nil
}

// normalizeAutoLimits validates a user's opt-in against the policy limits.
func normalizeAutoLimits(a *aiAutoLimits) error {
	var err error
	if a.Allow, err = cleanPatternList(a.Allow, 100); err != nil {
		return err
	}
	if a.Deny, err = cleanPatternList(a.Deny, 100); err != nil {
		return err
	}
	if a.PathAllow, err = cleanPatternList(a.PathAllow, 100); err != nil {
		return err
	}
	if a.PathDeny, err = cleanPatternList(a.PathDeny, 100); err != nil {
		return err
	}
	if len(a.Allow) == 0 && len(a.PathAllow) == 0 {
		return fmt.Errorf("automatic mode needs at least one allowed command or path pattern")
	}
	for _, p := range a.PathAllow {
		if !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "~/") {
			return fmt.Errorf("path patterns must be absolute: %s", p)
		}
	}
	maxMin := settingInt("ai_auto_max_minutes")
	if a.Minutes <= 0 || a.Minutes > maxMin {
		return fmt.Errorf("the time limit must be between 1 and %d minutes", maxMin)
	}
	maxAct := settingInt("ai_auto_max_actions")
	if a.MaxActions <= 0 || a.MaxActions > maxAct {
		a.MaxActions = maxAct
	}
	a.Until = time.Now().Add(time.Duration(a.Minutes) * time.Minute)
	a.Used = 0
	return nil
}

// active reports whether the limits still allow automatic actions (and why not).
func (a *aiAutoLimits) active(now time.Time) (bool, string) {
	if a == nil {
		return false, "automatic mode is off"
	}
	if now.After(a.Until) {
		return false, "the automatic mode time limit is over"
	}
	if a.Used >= a.MaxActions {
		return false, "the automatic mode action limit is reached"
	}
	return true, ""
}

func matchAny(pats []string, s string, pathMode bool) string {
	for _, p := range pats {
		if globMatch(p, s, pathMode) {
			return p
		}
	}
	return ""
}

// ─── the decision ────────────────────────────────────

const (
	aiAllow   = "allow"
	aiApprove = "approve"
	aiDeny    = "deny"
)

// aiAction is one thing the assistant wants to do, as the engine sees it.
type aiAction struct {
	Tool    string // tool name
	Kind    string // read | command | write
	Command string // command (Kind command)
	Path    string // file (Kind read / write)
}

type aiDecision struct {
	Action   string // allow | approve | deny
	Reason   string
	Rule     string // destructive rule id, or "readonly", "auto", "mode" …
	ReadOnly bool
	Auto     bool // allowed by the automatic mode limits (counts as an action)
}

// aiEvaluate decides what happens to an action in a mode. home expands ~/ in paths.
func aiEvaluate(mode string, auto *aiAutoLimits, act aiAction, home string, now time.Time) aiDecision {
	switch act.Kind {
	case "read":
		if why := aiReadDenied(act.Path); why != "" {
			return aiDecision{Action: aiDeny, Rule: "sensitive_path", Reason: act.Path + " is not readable by the assistant (" + why + ")"}
		}
		return aiDecision{Action: aiAllow, Rule: "readonly", ReadOnly: true}
	case "write":
		p := aiExpandHome(act.Path, home)
		if b := aiDestructivePath(p); b != nil {
			return aiDecision{Action: aiDeny, Rule: b.Rule, Reason: "blocked: " + b.Why}
		}
		switch mode {
		case aiModeReadOnly:
			return aiDecision{Action: aiDeny, Rule: "mode", Reason: "the session is read-only"}
		case aiModeAuto:
			if auto == nil {
				return aiDecision{Action: aiApprove, Rule: "ask"}
			}
			if m := matchAny(expandAll(auto.PathDeny, home), p, true); m != "" {
				return aiDecision{Action: aiDeny, Rule: "auto_deny", Reason: "the path matches the session's deny list (" + m + ")"}
			}
			if ok, why := auto.active(now); !ok {
				return aiDecision{Action: aiApprove, Rule: "auto_limit", Reason: why}
			}
			if m := matchAny(expandAll(auto.PathAllow, home), p, true); m != "" {
				return aiDecision{Action: aiAllow, Rule: "auto", Auto: true, Reason: "allowed by " + m}
			}
			return aiDecision{Action: aiApprove, Rule: "auto_outside", Reason: "the path is not in the session's allow list"}
		}
		return aiDecision{Action: aiApprove, Rule: "ask"}
	}
	// a command
	cmd := act.Command
	if b := aiDestructive(cmd); b != nil {
		if b.Power && getSetting("ai_allow_power") == "ask" && mode != aiModeReadOnly {
			return aiDecision{Action: aiApprove, Rule: "power", Reason: "power command (" + b.Why + "): always needs an approval"}
		}
		return aiDecision{Action: aiDeny, Rule: b.Rule, Reason: "blocked: " + b.Why}
	}
	v := classifyReadOnly(cmd)
	if v.ReadOnly {
		return aiDecision{Action: aiAllow, Rule: "readonly", ReadOnly: true}
	}
	switch mode {
	case aiModeReadOnly:
		return aiDecision{Action: aiDeny, Rule: "mode", Reason: "the session is read-only and the command is " + v.String()}
	case aiModeAuto:
		return aiEvaluateAuto(auto, cmd, home, now, v)
	}
	return aiDecision{Action: aiApprove, Rule: "ask", Reason: v.Reason}
}

func aiEvaluateAuto(auto *aiAutoLimits, cmd, home string, now time.Time, v roVerdict) aiDecision {
	if auto == nil {
		return aiDecision{Action: aiApprove, Rule: "ask", Reason: v.Reason}
	}
	trimmed := strings.Join(strings.Fields(cmd), " ")
	if m := matchAny(auto.Deny, trimmed, false); m != "" {
		return aiDecision{Action: aiDeny, Rule: "auto_deny", Reason: "the command matches the session's deny list (" + m + ")"}
	}
	p := shParse(cmd)
	for _, s := range p.Segs {
		if m := matchAny(auto.Deny, s.text(), false); m != "" {
			return aiDecision{Action: aiDeny, Rule: "auto_deny", Reason: "the command matches the session's deny list (" + m + ")"}
		}
	}
	if ok, why := auto.active(now); !ok {
		return aiDecision{Action: aiApprove, Rule: "auto_limit", Reason: why}
	}
	if p.opaque() || p.Vars {
		return aiDecision{Action: aiApprove, Rule: "auto_outside", Reason: "the command cannot be checked against the allow list (" + v.Reason + ")"}
	}
	for _, s := range p.Segs {
		for _, r := range s.Redirs {
			if r.writes() && r.Target != "/dev/null" {
				t := aiExpandHome(r.Target, home)
				if m := matchAny(expandAll(auto.PathDeny, home), t, true); m != "" {
					return aiDecision{Action: aiDeny, Rule: "auto_deny", Reason: "writes to a path on the session's deny list (" + m + ")"}
				}
				if matchAny(expandAll(auto.PathAllow, home), t, true) == "" {
					return aiDecision{Action: aiApprove, Rule: "auto_outside", Reason: "writes to " + r.Target + ", which is not in the allow list"}
				}
			}
		}
		if classifySegment(shSegment{Words: s.Words}).ReadOnly {
			continue
		}
		for _, w := range s.Words[1:] {
			if why := aiReadDenied(w); why != "" {
				return aiDecision{Action: aiApprove, Rule: "auto_outside", Reason: w + " is a sensitive file (" + why + ")"}
			}
		}
		if matchAny(auto.Allow, strings.Join(s.Words, " "), false) == "" {
			return aiDecision{Action: aiApprove, Rule: "auto_outside", Reason: "\"" + truncateStr(strings.Join(s.Words, " "), 80) + "\" is not in the session's allow list"}
		}
	}
	return aiDecision{Action: aiAllow, Rule: "auto", Auto: true, Reason: "allowed by the session's allow list"}
}

func aiExpandHome(p, home string) string {
	if home != "" && (p == "~" || strings.HasPrefix(p, "~/")) {
		return strings.TrimSuffix(home, "/") + strings.TrimPrefix(p, "~")
	}
	if strings.HasPrefix(p, "/") && !hasGlob(p) {
		return cleanRemotePath(p)
	}
	return p
}

func expandAll(pats []string, home string) []string {
	out := make([]string, 0, len(pats))
	for _, p := range pats {
		out = append(out, aiExpandHome(p, home))
	}
	return out
}

func aiAllowedFlag(userID int) bool {
	ok, _ := aiUserAllowed(userID)
	return ok
}
