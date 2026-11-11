package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ─── MONITORING CHECKS (v14.9.0) ─────────────────────
//
// Besides systemd units and Docker containers a server can check that a web
// address answers (http), a port accepts connections (tcp), a process runs
// (process) and a TLS certificate is not about to expire (cert). Every check
// runs ON the monitored server, over its SSH connection, in the same pass as
// the metrics: WRM itself never connects to the target, so a check cannot be
// used to reach hosts from WRM's network. The targets are validated before
// they are stored and reach the remote shell only as heredoc data.

var (
	monitorProcessPattern = regexp.MustCompile(`^[A-Za-z0-9_.@+\-]{1,64}$`)
	monitorHostPattern    = regexp.MustCompile(`^[A-Za-z0-9._:\-]{1,253}$`)
	monitorExpectPattern  = regexp.MustCompile(`^([1-5][0-9][0-9]|[1-5]xx)(-[1-5][0-9][0-9])?$`)
)

type monitorCheckResult struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Label  string `json:"label,omitempty"`
	State  string `json:"state"` // ok | warning | problem | unknown
	Raw    string `json:"raw,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func monitorHostPort(v, defPort string) (string, string, error) {
	h, p, err := net.SplitHostPort(v)
	if err != nil && defPort != "" && !strings.Contains(strings.Trim(v, "[]"), ":") {
		h, p, err = v, defPort, nil
	}
	if err != nil || h == "" || !monitorHostPattern.MatchString(h) || strings.HasPrefix(h, "-") {
		return "", "", fmt.Errorf("use host:port")
	}
	n, e := strconv.Atoi(p)
	if e != nil || n < 1 || n > 65535 {
		return "", "", fmt.Errorf("invalid port")
	}
	return h, strconv.Itoa(n), nil
}

func (x *monitorService) normalize() error {
	x.Name = strings.TrimSpace(x.Name)
	x.Label = strings.TrimSpace(x.Label)
	x.Expect = strings.TrimSpace(x.Expect)
	if x.Kind == "" {
		x.Kind = "systemd"
	}
	if len(x.Label) > 80 || strings.IndexFunc(x.Label, unicode.IsControl) >= 0 {
		return fmt.Errorf("check label is too long or has control characters")
	}
	if x.Kind != "http" {
		x.Expect = ""
	}
	switch x.Kind {
	case "systemd":
		if !monitorUnitPattern.MatchString(x.Name) {
			return fmt.Errorf("invalid systemd unit name %q", x.Name)
		}
	case "docker":
		if !monitorContainerPattern.MatchString(x.Name) {
			return fmt.Errorf("invalid Docker container name %q", x.Name)
		}
	case "process":
		if !monitorProcessPattern.MatchString(x.Name) {
			return fmt.Errorf("invalid process name %q (letters, digits, . _ - @ +)", x.Name)
		}
	case "tcp":
		h, p, err := monitorHostPort(x.Name, "")
		if err != nil {
			return fmt.Errorf("TCP check %q: %v", x.Name, err)
		}
		x.Name = net.JoinHostPort(h, p)
	case "cert":
		h, p, err := monitorHostPort(x.Name, "443")
		if err != nil || strings.Contains(h, ":") {
			return fmt.Errorf("certificate check %q: use a host name and optional port", x.Name)
		}
		x.Name = net.JoinHostPort(h, p)
	case "http":
		u, err := url.Parse(x.Name)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || len(x.Name) > 500 ||
			strings.ContainsAny(x.Name, " \t\"'`\\$") || strings.IndexFunc(x.Name, unicode.IsControl) >= 0 {
			return fmt.Errorf("HTTP check %q: use an http:// or https:// address without spaces, quotes or credentials", x.Name)
		}
		if x.Expect != "" && !monitorExpectPattern.MatchString(x.Expect) {
			return fmt.Errorf("expected HTTP status: 200, 2xx or 200-399")
		}
	default:
		return fmt.Errorf("a check is systemd, docker, http, tcp, process or cert")
	}
	return nil
}

func (x monitorService) key() string { return x.Kind + ":" + x.Name }

// monitorServiceScript prints one "kind:target=state" line per check. The
// targets come from a quoted heredoc, so none of them is parsed by the shell;
// they are used only as quoted variables. A target may contain "=", so the
// parser splits at the last one (states never contain "=").
func monitorServiceScript(list []monitorService) string {
	var b strings.Builder
	b.WriteString(`while IFS= read -r s; do `)
	b.WriteString(`case "$s" in `)
	b.WriteString(`systemd:*) n=${s#systemd:}; st=$(systemctl is-active -- "$n" 2>/dev/null) ;; `)
	b.WriteString(`docker:*) n=${s#docker:}; st=$(docker inspect -f '{{.State.Status}}' -- "$n" 2>/dev/null) || st=missing ;; `)
	b.WriteString(`http:*) n=${s#http:}; if command -v curl >/dev/null 2>&1; then st=$(curl -s -o /dev/null -m 8 -w '%{http_code} %{time_total}' -- "$n" 2>/dev/null); [ -n "$st" ] || st='000 0'; else st=nocurl; fi ;; `)
	b.WriteString(`tcp:*) n=${s#tcp:}; h=${n%:*}; p=${n##*:}; h=${h#[}; h=${h%]}; `)
	b.WriteString(`if command -v nc >/dev/null 2>&1; then if nc -z -w 5 "$h" "$p" >/dev/null 2>&1; then st=open; else st=closed; fi; `)
	b.WriteString(`elif command -v bash >/dev/null 2>&1 && command -v timeout >/dev/null 2>&1; then if timeout 5 bash -c 'exec 3<>"/dev/tcp/$1/$2"' _ "$h" "$p" >/dev/null 2>&1; then st=open; else st=closed; fi; `)
	b.WriteString(`else st=notool; fi ;; `)
	b.WriteString(`process:*) n=${s#process:}; n=$(printf '%.15s' "$n"); if pgrep -x -- "$n" >/dev/null 2>&1; then st=running; else st=stopped; fi ;; `)
	b.WriteString(`cert:*) n=${s#cert:}; h=${n%:*}; if command -v openssl >/dev/null 2>&1; then `)
	b.WriteString(`st=$(echo | timeout 10 openssl s_client -connect "$n" -servername "$h" 2>/dev/null | openssl x509 -noout -enddate 2>/dev/null); st=${st#notAfter=}; [ -n "$st" ] || st=noanswer; else st=notool; fi ;; `)
	b.WriteString(`*) st=unknown ;; esac; `)
	b.WriteString(`printf '%s=%s\n' "$s" "${st:-unknown}"; done <<'WRM_SERVICES'` + "\n")
	for _, x := range list {
		b.WriteString(x.key() + "\n")
	}
	b.WriteString("WRM_SERVICES\n")
	return b.String()
}

func parseMonitorServiceStates(out string) map[string]string {
	m := map[string]string{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		at := strings.LastIndex(line, "=")
		if at <= 0 {
			continue
		}
		k, v := line[:at], line[at+1:]
		if !strings.Contains(k, ":") {
			continue
		}
		m[truncateStr(k, 520)] = truncateStr(strings.TrimSpace(v), 80)
	}
	return m
}

// httpExpectOK reports whether code satisfies expect ("", "200", "2xx",
// "200-399").
func httpExpectOK(code int, expect string) bool {
	if expect == "" {
		return code >= 200 && code < 400
	}
	lo, hi, ranged := strings.Cut(expect, "-")
	if strings.HasSuffix(lo, "xx") {
		d, _ := strconv.Atoi(lo[:1])
		return code/100 == d
	}
	a, _ := strconv.Atoi(lo)
	if !ranged {
		return code == a
	}
	b, _ := strconv.Atoi(hi)
	return code >= a && code <= b
}

// judgeMonitorCheck turns the raw remote answer into ok/warning/problem.
func judgeMonitorCheck(x monitorService, raw string, certWarnDays float64, now time.Time) monitorCheckResult {
	r := monitorCheckResult{Kind: x.Kind, Name: x.Name, Label: x.Label, Raw: raw, State: "problem"}
	switch x.Kind {
	case "systemd":
		switch raw {
		case "active":
			r.State = "ok"
		case "activating", "reloading":
			r.State = "warning"
		case "unknown", "":
			r.State = "unknown"
		}
	case "docker":
		switch raw {
		case "running":
			r.State = "ok"
		case "restarting", "paused", "created":
			r.State = "warning"
		}
	case "http":
		if raw == "nocurl" {
			r.State, r.Detail = "unknown", "curl is not installed on the server"
			break
		}
		f := strings.Fields(raw)
		code := 0
		if len(f) > 0 {
			code, _ = strconv.Atoi(f[0])
		}
		secs := 0.0
		if len(f) > 1 {
			secs, _ = strconv.ParseFloat(f[1], 64)
		}
		switch {
		case code == 0:
			r.Detail = "no answer"
		case httpExpectOK(code, x.Expect):
			r.State = "ok"
			r.Detail = fmt.Sprintf("HTTP %d · %.0f ms", code, secs*1000)
			if secs > 5 {
				r.State = "warning"
			}
		default:
			r.Detail = fmt.Sprintf("HTTP %d · %.0f ms", code, secs*1000)
		}
	case "tcp":
		switch raw {
		case "open":
			r.State = "ok"
		case "notool":
			r.State, r.Detail = "unknown", "neither nc nor bash with timeout is installed"
		}
	case "process":
		if raw == "running" {
			r.State = "ok"
		}
	case "cert":
		if raw == "notool" {
			r.State, r.Detail = "unknown", "openssl is not installed on the server"
			break
		}
		at, err := time.Parse("Jan _2 15:04:05 2006 MST", strings.Join(strings.Fields(raw), " "))
		if err != nil {
			at, err = time.Parse("Jan 2 15:04:05 2006 MST", strings.Join(strings.Fields(raw), " "))
		}
		if err != nil {
			r.Detail = "no certificate received"
			break
		}
		days := at.Sub(now).Hours() / 24
		r.Detail = fmt.Sprintf("expires %s (%.0f days)", at.UTC().Format("2006-01-02"), days)
		switch {
		case days < 0:
			r.Detail = "expired " + at.UTC().Format("2006-01-02")
		case days < certWarnDays/2:
		case days < certWarnDays:
			r.State = "warning"
		default:
			r.State = "ok"
		}
	}
	return r
}

// monitorServiceProblem names the checks whose raw state is not active/running.
func monitorServiceProblem(states map[string]string) string {
	bad := []string{}
	for k, v := range states {
		if v != "active" && v != "running" {
			bad = append(bad, k+"="+v)
		}
	}
	if len(bad) == 0 {
		return ""
	}
	sort.Strings(bad)
	if len(bad) > 5 {
		bad = append(bad[:5], "…")
	}
	return "service problem: " + strings.Join(bad, ", ")
}

// ─── suggestions: what runs on the server ─────────────

type monitorSuggestion struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Label     string `json:"label,omitempty"`
	Suggested bool   `json:"suggested"`
}

const monitorSuggestScript = `LC_ALL=C; export LC_ALL
echo @@systemd; systemctl list-units --type=service --state=running --no-legend --plain 2>/dev/null | awk '{print $1}' | head -n 300
echo @@docker; docker ps --format '{{.Names}}' 2>/dev/null | head -n 200
echo @@ports; (ss -Htln 2>/dev/null || netstat -tln 2>/dev/null) | awk '{for(i=1;i<=NF;i++) if ($i ~ /:[0-9]+$/) {print $i; break}}' | head -n 300
echo @@host; hostname -f 2>/dev/null || hostname
echo @@end`

// monitorImportant are services worth watching on almost any server; they are
// pre-selected in the suggestion list.
var monitorImportant = regexp.MustCompile(`^(nginx|apache2|httpd|haproxy|caddy|traefik|mysql|mysqld|mariadb|postgresql.*|redis.*|mongod|docker|containerd|php.*-fpm|sshd|ssh|cron|crond|named|bind9|unbound|postfix|dovecot|keepalived|zabbix-agent2?|elasticsearch|opensearch|rabbitmq-server|memcached|tomcat.*|jenkins|gitlab-runsvdir|kubelet|wrm|wrm-pro|.*-api|.*-web|.*-worker)\.service$`)

func parseMonitorSuggestions(out, connHost string) []monitorSuggestion {
	section := ""
	seen := map[string]bool{}
	list := []monitorSuggestion{}
	add := func(x monitorSuggestion) {
		svc := monitorService{Kind: x.Kind, Name: x.Name}
		if svc.normalize() != nil || seen[svc.key()] {
			return
		}
		x.Name = svc.Name
		seen[svc.key()] = true
		list = append(list, x)
	}
	host := ""
	ports := []int{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@@") {
			section = line[2:]
			continue
		}
		if line == "" {
			continue
		}
		switch section {
		case "systemd":
			add(monitorSuggestion{Kind: "systemd", Name: line, Suggested: monitorImportant.MatchString(line)})
		case "docker":
			add(monitorSuggestion{Kind: "docker", Name: line, Suggested: true})
		case "ports":
			at := strings.LastIndex(line, ":")
			addr := strings.Trim(line[:at], "[]")
			p, err := strconv.Atoi(line[at+1:])
			if err != nil || p < 1 || p > 65535 {
				continue
			}
			// loopback-only listeners are internal; still worth a TCP check
			if addr == "" || addr == "*" || addr == "0.0.0.0" || addr == "::" || strings.HasPrefix(addr, "127.") || addr == "::1" {
				ports = append(ports, p)
			}
		case "host":
			if host == "" && monitorHostPattern.MatchString(line) {
				host = line
			}
		}
	}
	sort.Ints(ports)
	for i, p := range ports {
		if i > 0 && ports[i-1] == p {
			continue
		}
		add(monitorSuggestion{Kind: "tcp", Name: "127.0.0.1:" + strconv.Itoa(p), Suggested: p != 22 && p < 10000})
		switch p {
		case 80, 8080, 8000, 3000, 5000, 8081, 9000:
			add(monitorSuggestion{Kind: "http", Name: "http://127.0.0.1:" + strconv.Itoa(p) + "/", Suggested: false})
		case 443:
			name := host
			if strings.Contains(connHost, ".") && net.ParseIP(connHost) == nil {
				name = connHost
			}
			if strings.Contains(name, ".") {
				add(monitorSuggestion{Kind: "cert", Name: name + ":443", Suggested: true})
				add(monitorSuggestion{Kind: "http", Name: "https://" + name + "/", Suggested: false})
			}
		}
	}
	return list
}

// POST /api/monitoring/services/discover?id= — the owner's suggestion list.
func monitorSuggestHandler(w http.ResponseWriter, r *http.Request, c Connection, acc *connAccess) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if acc.Share != nil {
		jsonError(w, "Only the owner can configure monitoring", 403)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	cl, err := dialSSH(c, nil)
	if err != nil {
		jsonError(w, "SSH connection failed: "+truncateStr(redactText(err.Error()), 300), 502)
		return
	}
	defer cl.Close()
	out, err := monitorSSHCommand(ctx, cl, monitorSuggestScript)
	if err != nil && !strings.Contains(out, "@@end") {
		jsonError(w, "Could not list services on the server", 502)
		return
	}
	auditFileOp(r, acc, c, "monitor.services_discovered", c.Name, nil)
	jsonOK(w, map[string]interface{}{"suggestions": parseMonitorSuggestions(out, hostOnly(c.Host))})
}
