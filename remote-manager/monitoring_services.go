package main

import (
	"context"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// ─── MONITORED SERVICES ──────────────────────────────
//
// A monitored service is a systemd unit or a Docker container (v14.9.0 adds
// HTTP, TCP, process and certificate checks, monitoring_checks.go) whose state
// is read on every SSH poll, in the same round trip as the metrics. The command
// is one POSIX loop over a quoted heredoc, so a name cannot become shell
// syntax: the names are also validated against a strict pattern before they
// are stored.

var (
	monitorUnitPattern      = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.@:\-]{0,199}$`)
	monitorContainerPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.\-]{0,63}$`)
)

// collectMonitorServices fills s.Services and s.Checks and adds a reason for
// every check that is not ok, so the state and the notification name it.
func collectMonitorServices(ctx context.Context, cl *ssh.Client, cfg monitorConfig, s *monitorSample) {
	list := cfg.Services
	runCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	out, err := monitorSSHCommand(runCtx, cl, monitorServiceScript(list))
	s.Services = map[string]string{}
	s.Checks = map[string]monitorCheckResult{}
	if err != nil && strings.TrimSpace(out) == "" {
		for _, x := range list {
			s.Services[x.key()] = "unknown"
			s.Checks[x.key()] = monitorCheckResult{Kind: x.Kind, Name: x.Name, Label: x.Label, State: "unknown", Detail: "check failed"}
		}
		s.addReason("warning", "checks_failed", "", 0, 0)
		return
	}
	raw := parseMonitorServiceStates(out)
	now := time.Now()
	bad := map[string]string{}
	for _, x := range list {
		k := x.key()
		v, ok := raw[k]
		if !ok {
			v = "unknown"
		}
		s.Services[k] = v
		r := judgeMonitorCheck(x, v, cfg.Thresholds.CertDays, now)
		s.Checks[k] = r
		switch r.State {
		case "problem":
			bad[k] = v
			s.addReason("problem", "check", nonEmpty(x.Label, k), 0, 0)
		case "warning":
			s.addReason("warning", "check", nonEmpty(x.Label, k), 0, 0)
		}
	}
	if msg := monitorServiceProblem(bad); msg != "" {
		if s.Error != "" {
			s.Error += "; "
		}
		s.Error += msg
	}
}
