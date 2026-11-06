package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

// ─── CHAINS: CHECK CATALOG ────────────────────────────────────────────────────────────────
//
// Every check is one entry of a fixed catalog: a chain can name a type and its parameters,
// never a command. Checks on a server run in one SSH round trip per node: the script below is
// fixed text, and the parameters reach it on standard input as data (one tab-separated line
// per check), so no parameter is ever part of a command line. Network checks (DNS, TCP, HTTP,
// TLS) and storage checks run from WRM itself, through the egress policy.
//
// Only extracted values are kept: a configuration file is parsed here and only the one key a
// check asks for leaves this file; raw command output is never stored.

// chainCheckTypes lists the catalog: where a type runs and what it needs.
var chainCheckTypes = map[string]chainCheckType{
	"systemd":   {Where: "server", Needs: "target", Help: "systemd unit is active (target: unit name)"},
	"timer":     {Where: "server", Needs: "target", Help: "systemd timer is active (target: timer name)"},
	"docker":    {Where: "server", Needs: "target", Help: "Docker container is running (target: container name)"},
	"process":   {Where: "server", Needs: "target", Help: "number of processes whose command line matches (target: pattern)"},
	"port":      {Where: "server", Needs: "port", Help: "a TCP port is listening on the server (port)"},
	"file":      {Where: "server", Needs: "target", Help: "file exists, optionally not older than max_age_s (target: absolute path)"},
	"dir_count": {Where: "server", Needs: "target", Help: "number of files in a directory (target: absolute path)"},
	"config":    {Where: "server", Needs: "target,format,key", Help: "one value from a configuration file (format: ini, json, yaml, env, kv or regex; key: section.key, a.b.c or a regular expression)"},
	"mount":     {Where: "server", Needs: "target", Help: "a mount point is mounted, read-write or read-only (target: mount point)"},
	"cron":      {Where: "server", Needs: "target", Help: "a crontab line containing the text is enabled, disabled (commented out) or missing (target: text)"},
	"monitor":   {Where: "server", Needs: "", Help: "the state of WRM's own monitoring for the server (healthy, warning, …)"},
	"dns":       {Where: "wrm", Needs: "target", Help: "DNS answer for a name (target: name; record: A, AAAA, CNAME, TXT, MX, NS; resolver: optional ip[:port])"},
	"tcp":       {Where: "wrm", Needs: "target", Help: "a TCP port answers, checked from WRM (target: host:port)"},
	"http":      {Where: "wrm", Needs: "target", Help: "HTTP(S) status, or a JSON value with key (target: URL)"},
	"tls":       {Where: "wrm", Needs: "target", Help: "days until the TLS certificate expires (target: host[:port]; warn_days, default 14)"},
	"storage":   {Where: "storage", Needs: "", Help: "the storage's bucket answers; with prefix, the newest object under it (and max_age_s)"},
}

type chainCheckType struct {
	Where string `json:"where"` // server: over SSH on a server node; wrm: from WRM; storage: a storage node
	Needs string `json:"needs"`
	Help  string `json:"help"`
}

// chainResult is what one check found. Value and OK are compared for changes; Detail is
// shown but does not count as a change (an age grows on every check).
type chainResult struct {
	Value  string `json:"value"`
	OK     bool   `json:"ok"`
	Detail string `json:"detail,omitempty"`
	Error  string `json:"error,omitempty"`
}

var (
	chainSecretKey  = regexp.MustCompile(`(?i)(pass|secret|token|private|credential|api[_-]?key|access[_-]?key|auth)`)
	chainDNSRecords = map[string]bool{"A": true, "AAAA": true, "CNAME": true, "TXT": true, "MX": true, "NS": true}
	chainFormats    = map[string]bool{"ini": true, "json": true, "yaml": true, "env": true, "kv": true, "regex": true}
)

const chainMaxConfigBytes = 256 << 10

// chainNodeScript is the fixed script run on a server node. It reads "id<TAB>type<TAB>target"
// lines from standard input and prints "@id value" per check; multi-line or binary results are
// base64-encoded.
const chainNodeScript = `now=$(date +%s); printf '@@now %s\n' "$now"
tab=$(printf '\t')
while IFS="$tab" read -r id ty a; do
  case "$ty" in
  systemd|timer) v=$(systemctl is-active -- "$a" 2>/dev/null); [ -n "$v" ] || v=unknown ;;
  docker) v=$(docker inspect -f '{{.State.Status}}' -- "$a" 2>/dev/null) || v=missing ;;
  process) v=$(pgrep -f -- "$a" 2>/dev/null | wc -l | tr -d ' ') ;;
  port) v=$( { ss -Hltn 2>/dev/null || netstat -ltn 2>/dev/null; } | awk '{print $4}' | grep -cE "[:.]$a\$") ;;
  file) v=$(stat -c '%s %Y' -- "$a" 2>/dev/null) || v=missing ;;
  dir_count) if [ -d "$a" ]; then v=$(find "$a" -mindepth 1 -maxdepth 1 -type f 2>/dev/null | wc -l | tr -d ' '); else v=missing; fi ;;
  config) if [ -r "$a" ] && [ -f "$a" ]; then v=b64:$(head -c 262144 -- "$a" | base64 | tr -d '\n'); else v=missing; fi ;;
  mount) v=$(findmnt -n -o SOURCE,FSTYPE,OPTIONS -M "$a" 2>/dev/null | head -n 1); [ -n "$v" ] || v=missing ;;
  cron) v=b64:$( { crontab -l 2>/dev/null; cat /etc/crontab /etc/cron.d/* 2>/dev/null; } | head -c 262144 | base64 | tr -d '\n') ;;
  *) v=unknown ;;
  esac
  printf '@%s %s\n' "$id" "$(printf '%s' "$v" | tr '\n\r' '  ')"
done`

// chainNodeInput is the data the script reads: one line per server check of the node.
func chainNodeInput(checks []chainCheck) string {
	var b strings.Builder
	for _, c := range checks {
		if c.Type == "monitor" {
			continue
		}
		target := c.Target
		if c.Type == "port" {
			target = strconv.Itoa(c.Port)
		}
		if c.Type == "timer" && !strings.HasSuffix(target, ".timer") {
			target += ".timer"
		}
		b.WriteString(c.ID + "\t" + c.Type + "\t" + target + "\n")
	}
	return b.String()
}

// runChainNode runs the server checks of one node over SSH and interprets the answers.
func runChainNode(ctx context.Context, cl *ssh.Client, checks []chainCheck) (map[string]chainResult, error) {
	input := chainNodeInput(checks)
	raw := ""
	if input != "" {
		out, err := monitorSSHCommandInput(ctx, cl, "sh -c "+shellQuote(chainNodeScript), input)
		if err != nil && !strings.Contains(out, "@@now ") {
			return nil, fmt.Errorf("the checks did not run: %v", err)
		}
		raw = out
	}
	return parseChainNodeOutput(raw, checks, time.Now()), nil
}

// parseChainNodeOutput turns the script's answer into results; a check without an answer is
// reported as an error rather than silently passing.
func parseChainNodeOutput(out string, checks []chainCheck, fallbackNow time.Time) map[string]chainResult {
	answers := map[string]string{}
	now := fallbackNow.Unix()
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Buffer(make([]byte, 64<<10), 2<<20)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "@") {
			continue
		}
		id, v, _ := strings.Cut(line[1:], " ")
		if id == "@now" {
			if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
				now = n
			}
			continue
		}
		answers[id] = strings.TrimSpace(v)
	}
	res := map[string]chainResult{}
	for _, c := range checks {
		if c.Type == "monitor" {
			continue
		}
		v, ok := answers[c.ID]
		if !ok {
			res[c.ID] = chainResult{Value: "unknown", Error: "no answer from the server"}
			continue
		}
		res[c.ID] = interpretChainAnswer(c, v, now)
	}
	return res
}

func interpretChainAnswer(c chainCheck, v string, now int64) chainResult {
	r := chainResult{Value: v}
	switch c.Type {
	case "systemd", "timer":
		r.OK = v == "active"
	case "docker":
		r.OK = v == "running"
	case "process":
		n, _ := strconv.Atoi(v)
		r.Value, r.OK = strconv.Itoa(n), n > 0
	case "port":
		n, _ := strconv.Atoi(v)
		r.Value, r.OK = "closed", false
		if n > 0 {
			r.Value, r.OK = "listening", true
		}
	case "file":
		if v == "missing" {
			r.OK = false
			break
		}
		f := strings.Fields(v)
		if len(f) != 2 {
			r.Value, r.Error = "unknown", "unexpected answer"
			break
		}
		size, _ := strconv.ParseInt(f[0], 10, 64)
		mtime, _ := strconv.ParseInt(f[1], 10, 64)
		age := now - mtime
		r.Value, r.OK = "present", true
		r.Detail = fmt.Sprintf("%s, %s old", humanSize(size), chainAge(age))
		if c.MaxAgeS > 0 && age > int64(c.MaxAgeS) {
			r.Value, r.OK = "stale", false
		}
	case "dir_count":
		r.OK = v != "missing"
	case "config":
		if v == "missing" {
			r.OK = false
			break
		}
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(v, "b64:"))
		if err != nil {
			r.Value, r.Error = "unknown", "unreadable answer"
			break
		}
		val, found, err := chainExtract(c.Format, c.Key, data)
		switch {
		case err != nil:
			r.Value, r.Error = "unknown", err.Error()
		case !found:
			r.Value = "missing"
		default:
			r.Value, r.OK = truncateStr(val, 200), true
		}
	case "mount":
		if v == "missing" {
			r.OK = false
			break
		}
		f := strings.Fields(v)
		r.Value, r.OK = "mounted", true
		if len(f) >= 3 {
			r.Value = "rw"
			for _, o := range strings.Split(f[2], ",") {
				if o == "ro" {
					r.Value = "ro"
				}
			}
			r.Detail = f[0] + " (" + f[1] + ")"
		}
	case "cron":
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(v, "b64:"))
		if err != nil {
			r.Value, r.Error = "unknown", "unreadable answer"
			break
		}
		r.Value = chainCronState(string(data), c.Target)
		r.OK = r.Value == "enabled"
	}
	if c.Expect != nil && r.Error == "" {
		r.OK = c.Expect.matchValue(r.Value)
	}
	return r
}

// chainCronState: "enabled" when an active line contains the text, "disabled" when only a
// commented-out line does, otherwise "missing".
func chainCronState(tab, text string) string {
	state := "missing"
	for _, line := range strings.Split(tab, "\n") {
		l := strings.TrimSpace(line)
		if !strings.Contains(l, text) {
			continue
		}
		if strings.HasPrefix(l, "#") {
			state = "disabled"
			continue
		}
		return "enabled"
	}
	return state
}

func chainAge(sec int64) string {
	switch {
	case sec < 0:
		return "0 s"
	case sec < 120:
		return strconv.FormatInt(sec, 10) + " s"
	case sec < 7200:
		return strconv.FormatInt(sec/60, 10) + " min"
	case sec < 172800:
		return strconv.FormatInt(sec/3600, 10) + " h"
	}
	return strconv.FormatInt(sec/86400, 10) + " d"
}

// ─── value extraction from configuration files ──────────────────────────────────────────

// chainExtract reads one key from a configuration file. It returns found=false for a key that
// is not there, and an error only for a file that cannot be parsed.
func chainExtract(format, key string, data []byte) (string, bool, error) {
	switch format {
	case "json":
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			return "", false, fmt.Errorf("the file is not valid JSON")
		}
		return chainPath(v, key)
	case "yaml":
		var v interface{}
		if err := yaml.Unmarshal(data, &v); err != nil {
			return "", false, fmt.Errorf("the file is not valid YAML")
		}
		return chainPath(v, key)
	case "ini":
		v := chainINI(string(data), key)
		return v, v != "\x00", nil
	case "env", "kv":
		v, ok := chainKV(string(data), key, format == "env")
		return v, ok, nil
	case "regex":
		re, err := regexp.Compile(key)
		if err != nil {
			return "", false, fmt.Errorf("the regular expression is not valid")
		}
		m := re.FindStringSubmatch(string(data))
		if m == nil {
			return "", false, nil
		}
		if len(m) > 1 {
			return m[1], true, nil
		}
		return m[0], true, nil
	}
	return "", false, fmt.Errorf("unknown format %q", format)
}

// chainPath walks a.b.0.c through maps and lists.
func chainPath(v interface{}, key string) (string, bool, error) {
	cur := v
	for _, part := range strings.Split(key, ".") {
		switch t := cur.(type) {
		case map[string]interface{}:
			next, ok := t[part]
			if !ok {
				return "", false, nil
			}
			cur = next
		case []interface{}:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(t) {
				return "", false, nil
			}
			cur = t[i]
		default:
			return "", false, nil
		}
	}
	switch t := cur.(type) {
	case nil:
		return "null", true, nil
	case string:
		return t, true, nil
	case map[string]interface{}, []interface{}:
		// A whole object could carry nested secrets; a check reads one value.
		return "", false, fmt.Errorf("the key points to an object or a list; choose one value inside it")
	}
	return fmt.Sprint(cur), true, nil
}

// chainINI finds "section.key" (or a key before the first section); "\x00" means not found.
// Keys are matched case-insensitively, the way Python's configparser reads them.
func chainINI(text, key string) string {
	section, name := "", key
	if i := strings.LastIndex(key, "."); i >= 0 {
		section, name = key[:i], key[i+1:]
	}
	cur := ""
	for _, line := range strings.Split(text, "\n") {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, ";") {
			continue
		}
		if strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			cur = strings.TrimSpace(l[1 : len(l)-1])
			continue
		}
		i := strings.IndexAny(l, "=:")
		if i < 0 {
			continue
		}
		k, v := strings.TrimSpace(l[:i]), strings.TrimSpace(l[i+1:])
		if strings.EqualFold(cur, section) && strings.EqualFold(k, name) {
			return v
		}
	}
	return "\x00"
}

// chainKV reads KEY=VALUE (env, with an optional "export" and quotes) or "key = value" /
// "key: value" lines.
func chainKV(text, key string, env bool) (string, bool) {
	for _, line := range strings.Split(text, "\n") {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if env {
			l = strings.TrimPrefix(l, "export ")
		}
		i := strings.IndexAny(l, "=:")
		if env {
			i = strings.Index(l, "=")
		}
		if i < 0 {
			continue
		}
		if strings.TrimSpace(l[:i]) != key {
			continue
		}
		v := strings.TrimSpace(l[i+1:])
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		return v, true
	}
	return "", false
}

// ─── checks that run from WRM ─────────────────────────────────────────────────────────────

func runChainWRMCheck(ctx context.Context, c chainCheck, owner int, node chainNode) chainResult {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var r chainResult
	switch c.Type {
	case "dns":
		r = chainDNS(ctx, c)
	case "tcp":
		conn, err := egressDialContext(ctx, "tcp", c.Target)
		if err != nil {
			r = chainResult{Value: "closed", Detail: chainShortErr(err)}
		} else {
			conn.Close()
			r = chainResult{Value: "open", OK: true}
		}
	case "http":
		r = chainHTTP(ctx, c)
	case "tls":
		r = chainTLS(ctx, c)
	case "storage":
		r = chainStorage(ctx, c, owner, node)
	case "monitor":
		if conn, err := loadConnectionRaw(node.ConnID); err != nil || conn.UserID != owner {
			r = chainResult{Value: "unknown", Error: "the connection is not available to the chain's owner"}
			break
		}
		s := latestMonitorSample(node.ConnID)
		r = chainResult{Value: s.State, OK: s.State == "healthy", Detail: truncateStr(s.Error, 200)}
		if s.At == "" {
			r = chainResult{Value: "unknown", Detail: "monitoring is not set up for this server"}
		}
	default:
		r = chainResult{Value: "unknown", Error: "unknown check type"}
	}
	if c.Expect != nil && r.Error == "" {
		r.OK = c.Expect.matchValue(r.Value)
	}
	return r
}

func chainShortErr(err error) string {
	s := err.Error()
	if i := strings.LastIndex(s, ": "); i >= 0 && len(s)-i < 80 {
		s = s[i+2:]
	}
	return truncateStr(s, 160)
}

func chainDNS(ctx context.Context, c chainCheck) chainResult {
	res := &net.Resolver{PreferGo: true}
	if c.Resolver != "" {
		addr := c.Resolver
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr = net.JoinHostPort(addr, "53")
		}
		res.Dial = func(ctx context.Context, network, _ string) (net.Conn, error) {
			return egressDialContext(ctx, network, addr)
		}
	}
	var out []string
	var err error
	switch c.Record {
	case "", "A", "AAAA":
		var ips []net.IPAddr
		ips, err = res.LookupIPAddr(ctx, c.Target)
		for _, ip := range ips {
			v4 := ip.IP.To4() != nil
			if c.Record == "" || (c.Record == "A") == v4 {
				out = append(out, ip.IP.String())
			}
		}
	case "CNAME":
		var name string
		name, err = res.LookupCNAME(ctx, c.Target)
		if name != "" {
			out = append(out, strings.TrimSuffix(name, "."))
		}
	case "TXT":
		out, err = res.LookupTXT(ctx, c.Target)
	case "MX":
		var mx []*net.MX
		mx, err = res.LookupMX(ctx, c.Target)
		for _, m := range mx {
			out = append(out, strings.TrimSuffix(m.Host, "."))
		}
	case "NS":
		var ns []*net.NS
		ns, err = res.LookupNS(ctx, c.Target)
		for _, n := range ns {
			out = append(out, strings.TrimSuffix(n.Host, "."))
		}
	}
	if err != nil && len(out) == 0 {
		return chainResult{Value: "none", Detail: chainShortErr(err)}
	}
	sort.Strings(out)
	return chainResult{Value: truncateStr(strings.Join(out, ","), 300), OK: len(out) > 0}
}

func chainHTTP(ctx context.Context, c chainCheck) chainResult {
	client := &http.Client{
		Transport: &http.Transport{DialContext: egressDialContext, TLSHandshakeTimeout: 8 * time.Second,
			ResponseHeaderTimeout: 8 * time.Second, DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Target, nil)
	if err != nil {
		return chainResult{Value: "unknown", Error: "the address is not valid"}
	}
	req.Header.Set("User-Agent", "WRM-Chains/1")
	resp, err := client.Do(req) // #nosec G704 -- the URL passed guardEgressURL and every dial is egress-checked
	if err != nil {
		return chainResult{Value: "unreachable", Detail: chainShortErr(err)}
	}
	defer resp.Body.Close()
	status := strconv.Itoa(resp.StatusCode)
	if c.Key == "" {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
		return chainResult{Value: status, OK: resp.StatusCode >= 200 && resp.StatusCode < 400}
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, chainMaxConfigBytes))
	v, found, err := chainExtract("json", c.Key, body)
	switch {
	case err != nil:
		return chainResult{Value: "unknown", Detail: "HTTP " + status + ": " + err.Error()}
	case !found:
		return chainResult{Value: "missing", Detail: "HTTP " + status}
	}
	return chainResult{Value: truncateStr(v, 200), OK: true, Detail: "HTTP " + status}
}

func chainTLS(ctx context.Context, c chainCheck) chainResult {
	addr := c.Target
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "443")
	}
	host, _, _ := net.SplitHostPort(addr)
	raw, err := egressDialContext(ctx, "tcp", addr)
	if err != nil {
		return chainResult{Value: "unreachable", Detail: chainShortErr(err)}
	}
	defer raw.Close()
	conn := tls.Client(raw, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	if err := conn.HandshakeContext(ctx); err != nil {
		return chainResult{Value: "invalid", Detail: chainShortErr(err)}
	}
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return chainResult{Value: "invalid", Detail: "no certificate"}
	}
	days := int(time.Until(certs[0].NotAfter).Hours() / 24)
	warn := c.WarnDays
	if warn <= 0 {
		warn = 14
	}
	return chainResult{Value: strconv.Itoa(days), OK: days > warn,
		Detail: "expires " + certs[0].NotAfter.UTC().Format("2006-01-02")}
}

// chainStorageCache keeps the last answer of each storage check, so a chain that is polled every
// 30 seconds does not send a HEAD or a LIST to the bucket every 30 seconds: the bucket is asked
// at most every chain_storage_check_s (default 300 s).
var chainStorageCache = struct {
	sync.Mutex
	m map[string]chainStorageCached
}{m: map[string]chainStorageCached{}}

type chainStorageCached struct {
	at  time.Time
	res chainResult
}

func chainStorageCacheKey(owner int, node chainNode, c chainCheck) string {
	return fmt.Sprintf("%d|%d|%s|%d", owner, node.StorageID, c.Prefix, c.MaxAgeS)
}

func chainStorage(ctx context.Context, c chainCheck, owner int, node chainNode) chainResult {
	key := chainStorageCacheKey(owner, node, c)
	every := time.Duration(settingInt("chain_storage_check_s")) * time.Second
	chainStorageCache.Lock()
	hit, ok := chainStorageCache.m[key]
	chainStorageCache.Unlock()
	if ok && every > 0 && time.Since(hit.at) < every {
		return hit.res
	}
	res := chainStorageAsk(withStorageSource(ctx, "chain"), c, owner, node)
	chainStorageCache.Lock()
	if len(chainStorageCache.m) > 2000 {
		chainStorageCache.m = map[string]chainStorageCached{}
	}
	chainStorageCache.m[key] = chainStorageCached{at: time.Now(), res: res}
	chainStorageCache.Unlock()
	return res
}

func chainStorageAsk(ctx context.Context, c chainCheck, owner int, node chainNode) chainResult {
	s, err := loadStorage(node.StorageID, owner)
	if err != nil || s.UserID != owner {
		return chainResult{Value: "unknown", Error: "the storage is not available to the chain's owner"}
	}
	cl, err := s.clientFor()
	if err != nil {
		return chainResult{Value: "unknown", Error: err.Error()}
	}
	if c.Prefix == "" {
		if err := cl.headBucket(ctx); err != nil {
			return chainResult{Value: "unreachable", Detail: chainShortErr(err)}
		}
		return chainResult{Value: "reachable", OK: true}
	}
	list, err := cl.list(ctx, strings.TrimSuffix(c.Prefix, "/")+"/", "", 1000)
	if err != nil {
		return chainResult{Value: "unreachable", Detail: chainShortErr(err)}
	}
	newest, at := "", time.Time{}
	for _, o := range list.Objects {
		t, err := time.Parse(time.RFC3339, o.LastModified)
		if err == nil && t.After(at) {
			newest, at = o.Key, t
		}
	}
	if newest == "" {
		return chainResult{Value: "empty"}
	}
	age := int64(time.Since(at).Seconds())
	r := chainResult{Value: "fresh", OK: true, Detail: newest + ", " + chainAge(age) + " old"}
	if c.MaxAgeS > 0 && age > int64(c.MaxAgeS) {
		r.Value, r.OK = "stale", false
	}
	return r
}

// ─── validation of one check's parameters ─────────────────────────────────────────────────

func chainHasControl(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}

func validateChainCheck(c *chainCheck, node chainNode) error {
	ct, ok := chainCheckTypes[c.Type]
	if !ok {
		return fmt.Errorf("unknown type %q", c.Type)
	}
	switch ct.Where {
	case "server":
		if node.Kind != "server" {
			return fmt.Errorf("a %s check needs a server node", c.Type)
		}
	case "storage":
		if node.Kind != "storage" {
			return fmt.Errorf("a storage check needs a storage node")
		}
	}
	c.Target = strings.TrimSpace(c.Target)
	if len(c.Target) > 512 || chainHasControl(c.Target) {
		return fmt.Errorf("target: at most 512 characters, no control characters")
	}
	if strings.Contains(ct.Needs, "target") && c.Target == "" {
		return fmt.Errorf("target is required")
	}
	if c.MaxAgeS < 0 || c.WarnDays < 0 {
		return fmt.Errorf("max_age_s and warn_days must not be negative")
	}
	abs := func() error {
		if !strings.HasPrefix(c.Target, "/") {
			return fmt.Errorf("target must be an absolute path")
		}
		return nil
	}
	switch c.Type {
	case "systemd", "timer":
		if !monitorUnitPattern.MatchString(c.Target) {
			return fmt.Errorf("target is not a valid unit name")
		}
	case "docker":
		if !monitorContainerPattern.MatchString(c.Target) {
			return fmt.Errorf("target is not a valid container name")
		}
	case "process":
		if len(c.Target) > 200 {
			return fmt.Errorf("the pattern is too long")
		}
	case "port":
		if c.Port < 1 || c.Port > 65535 {
			return fmt.Errorf("port must be 1-65535")
		}
	case "file", "dir_count", "mount":
		return abs()
	case "config":
		if err := abs(); err != nil {
			return err
		}
		if !chainFormats[c.Format] {
			return fmt.Errorf("format must be ini, json, yaml, env, kv or regex")
		}
		return validateChainKey(c.Format, c.Key)
	case "dns":
		c.Record = strings.ToUpper(strings.TrimSpace(c.Record))
		if c.Record != "" && !chainDNSRecords[c.Record] {
			return fmt.Errorf("record must be A, AAAA, CNAME, TXT, MX or NS")
		}
		if strings.ContainsAny(c.Target, " /:") {
			return fmt.Errorf("target must be a DNS name")
		}
		if c.Resolver != "" {
			h := c.Resolver
			if hh, _, err := net.SplitHostPort(h); err == nil {
				h = hh
			}
			if net.ParseIP(h) == nil {
				return fmt.Errorf("resolver must be an IP address with an optional port")
			}
			if err := guardEgressHost(h); err != nil {
				return err
			}
		}
	case "tcp":
		h, p, err := net.SplitHostPort(c.Target)
		if err != nil || h == "" || p == "" {
			return fmt.Errorf("target must be host:port")
		}
		return guardEgressHost(strings.ToLower(h))
	case "tls":
		h := c.Target
		if hh, _, err := net.SplitHostPort(h); err == nil {
			h = hh
		}
		return guardEgressHost(strings.ToLower(h))
	case "http":
		if err := guardEgressURL(c.Target); err != nil {
			return err
		}
		if c.Key != "" {
			return validateChainKey("json", c.Key)
		}
	case "storage":
		if c.Prefix != "" {
			if _, err := normalizePrefix(c.Prefix); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateChainKey refuses keys that would read a secret: a chain shows values to whoever
// opens it and in notifications, so passwords and tokens never become check values.
func validateChainKey(format, key string) error {
	key = strings.TrimSpace(key)
	if key == "" || len(key) > 200 || chainHasControl(key) {
		return fmt.Errorf("key is required (at most 200 characters)")
	}
	if chainSecretKey.MatchString(key) {
		return fmt.Errorf("key %q looks like a secret; WRM does not read secrets into a chain", key)
	}
	if format == "regex" {
		if _, err := regexp.Compile(key); err != nil {
			return fmt.Errorf("key is not a valid regular expression")
		}
	}
	return nil
}
