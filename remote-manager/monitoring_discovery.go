package main

import (
	"context"
	"encoding/json"
	"io"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
)

// A discovery grants no remote privileges: every read still checks the current
// connection permission, sensitive paths and all symlink components.
type monitorRootCache struct {
	Owner  int
	Config string
	Until  time.Time
	Roots  []monitorLogRoot
}

var monitorDiscoveredRoots sync.Map

func monitorRootConfig(id int) string {
	b, _ := json.Marshal(loadMonitorConfig(id).Roots)
	return string(b)
}
func monitorAvailableRoots(c Connection) []monitorLogRoot {
	roots := append([]monitorLogRoot(nil), loadMonitorConfig(c.ID).Roots...)
	if v, ok := monitorDiscoveredRoots.Load(c.ID); ok {
		x := v.(monitorRootCache)
		if x.Owner == c.UserID && x.Config == monitorRootConfig(c.ID) && time.Now().Before(x.Until) {
			roots = x.Roots
		}
	}
	return roots
}
func monitorLogName(name string) bool {
	n := strings.ToLower(name)
	return strings.Contains(n, ".log") || strings.HasSuffix(n, ".out") || strings.HasSuffix(n, ".err") || oneOf(n, "syslog", "messages", "auth.log", "kern.log", "daemon.log", "debug", "access", "error", "stdout", "stderr")
}
func monitorDiscoverRoots(ctx context.Context, fs *sftp.Client, c Connection) ([]monitorLogRoot, []string) {
	cfg := loadMonitorConfig(c.ID)
	encoded, _ := json.Marshal(cfg.Roots)
	configKey := string(encoded)
	roots := []monitorLogRoot{}
	seen := map[string]bool{}
	warnings := []string{}
	add := func(p, service, operator string) {
		p = path.Clean(p)
		if len(roots) >= 100 || seen[p] || !strings.HasPrefix(p, "/") || p == "/" || aiReadDenied(p) != "" || strings.ContainsAny(p, "\x00\r\n") {
			return
		}
		// Candidate metadata is cheap. Full ancestor/REALPATH validation runs
		// before directory listing and again before every content read.
		st, err := fs.Lstat(p)
		if err != nil || !st.IsDir() {
			return
		}
		seen[p] = true
		roots = append(roots, monitorLogRoot{Path: p, Service: service, Operator: operator})
	}
	for _, r := range cfg.Roots {
		add(r.Path, r.Service, r.Operator)
		if !seen[path.Clean(r.Path)] {
			warnings = append(warnings, r.Path+": directory unavailable or not readable")
		}
	}
	// Known installations provide labels without requiring another expensive Git
	// comparison. Their small log directories are preferred over broad roots.
	profiles := loadGitConfigProfiles(c.UserID)
	for _, in := range loadGitInstalls(c.UserID, false, "i.conn_id=?", c.ID) {
		if ctx.Err() != nil {
			break
		}
		operator := in.Env
		for _, p := range profiles {
			if p.App == in.App && p.Operator != "" {
				operator = p.Operator
				break
			}
		}
		for _, suffix := range []string{"log", "logs"} {
			add(path.Join(in.Path, suffix), in.App, operator)
		}
	}
	// Read only bounded Supervisor metadata; no commands or configurations run.
	for _, dir := range []string{"/etc/supervisor/conf.d", "/etc/supervisord.d"} {
		if ctx.Err() != nil {
			break
		}
		if _, err := monitorSafePath(fs, dir, dir); err != nil {
			continue
		}
		entries, err := fs.ReadDir(dir)
		if err != nil {
			continue
		}
		for i, e := range entries {
			if i >= 40 || ctx.Err() != nil {
				break
			}
			if !e.Mode().IsRegular() || !oneOf(strings.ToLower(path.Ext(e.Name())), ".conf", ".ini") {
				continue
			}
			p := path.Join(dir, e.Name())
			if _, err := monitorSafePath(fs, dir, p); err != nil {
				continue
			}
			f, err := fs.Open(p)
			if err != nil {
				continue
			}
			b, _ := io.ReadAll(io.LimitReader(f, 64<<10))
			f.Close()
			service := strings.TrimSuffix(e.Name(), path.Ext(e.Name()))
			work := ""
			logs := []string{}
			flush := func() {
				for _, v := range logs {
					v = strings.ReplaceAll(v, "%(here)s", dir)
					if strings.Contains(v, "%(") {
						continue
					}
					if !path.IsAbs(v) && work != "" {
						v = path.Join(work, v)
					}
					if path.IsAbs(v) {
						add(path.Dir(v), service, "")
					}
				}
				logs = nil
			}
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "[") {
					flush()
					work = ""
					if strings.HasPrefix(line, "[program:") {
						service = strings.TrimSuffix(strings.TrimPrefix(line, "[program:"), "]")
					}
					continue
				}
				k, v, ok := strings.Cut(line, "=")
				if !ok {
					continue
				}
				v = strings.Trim(strings.TrimSpace(v), "\"'")
				switch strings.TrimSpace(k) {
				case "directory":
					work = v
				case "stdout_logfile", "stderr_logfile":
					if v != "AUTO" && v != "NONE" {
						logs = append(logs, v)
					}
				}
			}
			flush()
		}
	}
	cat, _ := loadGitWorkspace(c.UserID)
	bases := append([]string{"/home/operators", "/srv", "/opt"}, cat.ServerRoots...)
	// Bounded operator/service discovery, never a recursive /home scan. On the
	// deployment layout /home/operators/<operator>/<service>/log(s) the first level
	// is the operator and the second level is the service.
	type candidate struct {
		p, service, operator string
		depth                int
	}
	queue := []candidate{}
	baseSeen := map[string]bool{}
	for _, p := range bases {
		p = path.Clean(p)
		if p == "/" || baseSeen[p] || !strings.HasPrefix(p, "/") {
			continue
		}
		baseSeen[p] = true
		queue = append(queue, candidate{p: p})
	}
	visited := 0
	for len(queue) > 0 && visited < 300 && len(roots) < 100 && ctx.Err() == nil {
		d := queue[0]
		queue = queue[1:]
		visited++
		if d.depth == 0 {
			if _, err := monitorSafePath(fs, d.p, d.p); err != nil {
				continue
			}
		}
		entries, err := fs.ReadDir(d.p)
		if err != nil {
			continue
		}
		for i, e := range entries {
			if i >= 300 {
				break
			}
			if !e.IsDir() {
				continue
			}
			n := e.Name()
			lower := strings.ToLower(n)
			if strings.HasPrefix(n, ".") || strings.HasSuffix(lower, ".egg-info") || oneOf(lower, "node_modules", "venv", "__pycache__", "backups", "build", "dist", "archive", "workingdir", "temp", "privkey") {
				continue
			}
			p := path.Join(d.p, n)
			if oneOf(lower, "log", "logs") {
				add(p, d.service, d.operator)
			} else if d.depth < 2 {
				service, operator := d.service, d.operator
				if d.depth == 0 {
					operator = n
				} else {
					service = n
				}
				queue = append(queue, candidate{p, service, operator, d.depth + 1})
			}
		}
	}
	if ctx.Err() != nil {
		warnings = append(warnings, "Automatic folder discovery time budget reached; readable known roots are listed below")
	}
	if len(queue) > 0 || len(roots) >= 100 {
		warnings = append(warnings, "Automatic root discovery reached its directory/root limit")
	}
	sort.SliceStable(roots, func(i, j int) bool { return len(roots[i].Path) > len(roots[j].Path) })
	if monitorRootConfig(c.ID) == configKey {
		monitorDiscoveredRoots.Store(c.ID, monitorRootCache{c.UserID, configKey, time.Now().Add(time.Hour), roots})
	} else {
		warnings = append(warnings, "Monitoring roots changed during discovery; find logs again")
	}
	return roots, warnings
}
