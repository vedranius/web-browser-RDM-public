package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// ─── AGENTLESS MONITORING (v14.9.0) ──────────────────
//
// The default transport: WRM reads a Linux server over its existing SSH
// connection with one read-only command (no agent, no installation, no root).
// It covers CPU (busy, iowait, steal over one second), memory and swap, load,
// filesystems and inodes, network counters (rates come from the previous
// sample), TCP connections, processes and zombies, failed systemd units, a
// pending reboot, clock synchronisation and the busiest processes.

const monitorAgentlessScript = `LC_ALL=C; export LC_ALL
echo @@os; uname -s; uname -srm; ( . /etc/os-release 2>/dev/null; echo "${PRETTY_NAME:-}" )
echo @@host; hostname 2>/dev/null || cat /proc/sys/kernel/hostname
echo @@load; cat /proc/loadavg
echo @@stat; head -n1 /proc/stat; sleep 1; head -n1 /proc/stat
echo @@mem; grep -E '^(MemTotal|MemAvailable|SwapTotal|SwapFree):' /proc/meminfo
echo @@df; df -P -k -x tmpfs -x devtmpfs -x squashfs -x overlay 2>/dev/null || df -P -k
echo @@inode; df -P -i -x tmpfs -x devtmpfs -x squashfs -x overlay 2>/dev/null
echo @@net; cat /proc/net/dev
echo @@uptime; cat /proc/uptime
echo @@ncpu; nproc 2>/dev/null || getconf _NPROCESSORS_ONLN
echo @@procs; ps -e -o stat= 2>/dev/null | awk '{n++; if ($1 ~ /^Z/) z++} END {print n+0, z+0}'
echo @@tcp; cat /proc/net/tcp /proc/net/tcp6 2>/dev/null | awk '$4=="01"{n++} END {print n+0}'
echo @@top; ps -e -o pid=,pcpu=,pmem=,comm= --sort=-pcpu 2>/dev/null | head -n 5; ps -e -o pid=,pcpu=,pmem=,comm= --sort=-pmem 2>/dev/null | head -n 5
echo @@failed; systemctl list-units --state=failed --no-legend --plain 2>/dev/null | awk '{print $1}' | head -n 20
echo @@reboot; [ -f /var/run/reboot-required ] && echo yes
echo @@ntp; timedatectl show -p NTPSynchronized --value 2>/dev/null
echo @@end`

type monitorProcess struct {
	PID  int     `json:"pid"`
	CPU  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
	Name string  `json:"name"`
}

func monitorSections(out string) map[string][]string {
	m := map[string][]string{}
	cur := ""
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, "@@") {
			cur = strings.TrimSpace(line[2:])
			m[cur] = []string{}
			continue
		}
		if cur != "" && strings.TrimSpace(line) != "" {
			m[cur] = append(m[cur], line)
		}
	}
	return m
}

// parseAgentless fills s from the script output; it returns false when the
// output is not from a Linux server.
func parseAgentless(out string, s *monitorSample) (bool, string) {
	sec := monitorSections(out)
	osl := sec["os"]
	if len(osl) == 0 {
		return false, "no answer from the server"
	}
	if strings.TrimSpace(osl[0]) != "Linux" {
		return false, "agentless monitoring reads Linux servers; this one reports " + truncateStr(strings.TrimSpace(osl[0]), 40) + " (use the Zabbix agent transport for it)"
	}
	at := s.At
	set := func(k string, v float64) {
		x := v
		s.Metrics[k] = monitorMetric{Value: &x, At: at}
	}
	text := func(k, v string) { s.Metrics[k] = monitorMetric{Text: truncateStr(v, 500), At: at} }
	if len(osl) > 1 {
		os := strings.TrimSpace(osl[1])
		if len(osl) > 2 && strings.TrimSpace(osl[2]) != "" {
			os = strings.TrimSpace(osl[2]) + " · " + os
		}
		text("os", os)
	}
	if h := sec["host"]; len(h) > 0 {
		s.Hostname = truncateStr(strings.TrimSpace(h[0]), 200)
		text("hostname", s.Hostname)
	}
	if l := sec["load"]; len(l) > 0 {
		f := strings.Fields(l[0])
		for i, k := range []string{"load", "load5", "load15"} {
			if i < len(f) {
				if v, err := strconv.ParseFloat(f[i], 64); err == nil {
					set(k, v)
				}
			}
		}
	}
	if st := sec["stat"]; len(st) == 2 {
		a, b := strings.Fields(st[0]), strings.Fields(st[1])
		if len(a) >= 5 && len(a) == len(b) && a[0] == "cpu" {
			var t0, t1 float64
			d := make([]float64, len(a))
			for i := 1; i < len(a); i++ {
				x, _ := strconv.ParseFloat(a[i], 64)
				y, _ := strconv.ParseFloat(b[i], 64)
				d[i] = y - x
				t0 += x
				t1 += y
			}
			if dt := t1 - t0; dt > 0 {
				idle := d[4]
				if len(d) > 5 {
					idle += d[5]
				}
				set("cpu_percent", round2(max(0, 100*(1-idle/dt))))
				if len(d) > 5 {
					set("iowait", round2(100*d[5]/dt))
				}
				if len(d) > 8 {
					set("steal", round2(100*d[8]/dt))
				}
			}
		}
	}
	mem := map[string]float64{}
	for _, ln := range sec["mem"] {
		f := strings.Fields(ln)
		if len(f) >= 2 {
			v, _ := strconv.ParseFloat(f[1], 64)
			mem[strings.TrimSuffix(f[0], ":")] = v * 1024
		}
	}
	if mem["MemTotal"] > 0 {
		set("memory_total", mem["MemTotal"])
		set("memory_available", mem["MemAvailable"])
	}
	if _, ok := mem["SwapTotal"]; ok {
		set("swap_total", mem["SwapTotal"])
		if mem["SwapTotal"] > 0 {
			set("swap_percent", round2(100*(1-mem["SwapFree"]/mem["SwapTotal"])))
		}
	}
	readDF := func(lines []string, prefix string) {
		seen := map[string]bool{}
		n := 0
		for i, ln := range lines {
			f := strings.Fields(ln)
			if i == 0 || len(f) < 6 {
				continue
			}
			size, _ := strconv.ParseInt(f[1], 10, 64)
			pct, err := strconv.ParseFloat(strings.TrimSuffix(f[4], "%"), 64)
			mount := strings.Join(f[5:], " ")
			if err != nil || size <= 0 || seen[mount] || strings.HasPrefix(mount, "/snap/") || strings.HasPrefix(mount, "/run") || strings.HasPrefix(mount, "/dev") || strings.HasPrefix(mount, "/sys") || strings.HasPrefix(mount, "/proc") || strings.HasPrefix(mount, "/var/lib/docker/") {
				continue
			}
			seen[mount] = true
			set(prefix+truncateStr(mount, 100), pct)
			if prefix == "disk:" {
				set("disksize:"+truncateStr(mount, 100), float64(size)*1024)
			}
			if n++; n >= 24 {
				break
			}
		}
	}
	readDF(sec["df"], "disk:")
	readDF(sec["inode"], "inode:")
	nif := 0
	for _, ln := range sec["net"] {
		name, rest, ok := strings.Cut(ln, ":")
		name = strings.TrimSpace(name)
		if !ok || name == "lo" || !monitorInterfaceName(name) || strings.HasPrefix(name, "veth") || strings.HasPrefix(name, "br-") || strings.HasPrefix(name, "docker") {
			continue
		}
		f := strings.Fields(rest)
		if len(f) < 9 {
			continue
		}
		in, _ := strconv.ParseFloat(f[0], 64)
		out, _ := strconv.ParseFloat(f[8], 64)
		if in == 0 && out == 0 {
			continue
		}
		set("network:"+name+":in", in)
		set("network:"+name+":out", out)
		if nif++; nif >= 8 {
			break
		}
	}
	if u := sec["uptime"]; len(u) > 0 {
		if f := strings.Fields(u[0]); len(f) > 0 {
			if v, err := strconv.ParseFloat(f[0], 64); err == nil {
				set("uptime", v)
			}
		}
	}
	if n := sec["ncpu"]; len(n) > 0 {
		if v, err := strconv.Atoi(strings.TrimSpace(n[0])); err == nil && v > 0 {
			set("cpu_count", float64(v))
		}
	}
	if p := sec["procs"]; len(p) > 0 {
		if f := strings.Fields(p[0]); len(f) == 2 {
			a, _ := strconv.ParseFloat(f[0], 64)
			b, _ := strconv.ParseFloat(f[1], 64)
			if a > 0 {
				set("processes", a)
				set("zombies", b)
			}
		}
	}
	if t := sec["tcp"]; len(t) > 0 {
		if v, err := strconv.ParseFloat(strings.TrimSpace(t[0]), 64); err == nil {
			set("tcp_established", v)
		}
	}
	seenPID := map[int]bool{}
	for _, ln := range sec["top"] {
		f := strings.Fields(ln)
		if len(f) < 4 {
			continue
		}
		pid, _ := strconv.Atoi(f[0])
		if pid == 0 || seenPID[pid] {
			continue
		}
		seenPID[pid] = true
		c, _ := strconv.ParseFloat(f[1], 64)
		m, _ := strconv.ParseFloat(f[2], 64)
		s.Top = append(s.Top, monitorProcess{PID: pid, CPU: c, Mem: m, Name: truncateStr(strings.Join(f[3:], " "), 64)})
	}
	if _, ok := sec["failed"]; ok {
		names := []string{}
		for _, ln := range sec["failed"] {
			if n := strings.TrimSpace(ln); monitorUnitPattern.MatchString(n) {
				names = append(names, n)
			}
		}
		set("failed_units", float64(len(names)))
		if len(names) > 0 {
			m := s.Metrics["failed_units"]
			m.Text = truncateStr(strings.Join(names, ", "), 500)
			s.Metrics["failed_units"] = m
		}
	}
	reboot := 0.0
	if r := sec["reboot"]; len(r) > 0 && strings.TrimSpace(r[0]) == "yes" {
		reboot = 1
	}
	set("reboot_required", reboot)
	if n := sec["ntp"]; len(n) > 0 {
		switch strings.TrimSpace(n[0]) {
		case "yes":
			set("time_sync", 1)
		case "no":
			set("time_sync", 0)
		}
	}
	return true, ""
}

// collectAgentless reads one sample over an open SSH client.
func collectAgentless(ctx context.Context, cl *ssh.Client, s *monitorSample) {
	runCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	out, err := monitorSSHCommand(runCtx, cl, monitorAgentlessScript)
	ok, why := parseAgentless(out, s)
	if !ok {
		s.State = "failed"
		if err != nil && why == "no answer from the server" {
			why = "the read-only command failed: " + truncateStr(redactText(err.Error()), 200)
		}
		s.Error = why
		return
	}
	s.State = "healthy"
	finishMonitor(s)
}
