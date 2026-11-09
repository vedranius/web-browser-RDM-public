package main

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
)

// Folders that never hold logs worth listing and can be huge: Python virtual environments
// (also renamed copies such as venv.bak.2026-07-16), package caches and VCS data. Discovery
// skips them by name in the find command and in the SFTP walks, so a large tree cannot use
// up the time budget before the real log folders are reached.
var logAlwaysSkipDirs = []string{"venv*", ".venv*", "virtualenv*", "site-packages", "dist-packages", "node_modules", "__pycache__", ".git", ".tox", ".mypy_cache", ".pytest_cache"}

// logSkipDir: the folder name matches one of the built-in names or the user's exclude list.
func logSkipDir(name string, exclude []string) bool {
	n := strings.ToLower(name)
	for _, list := range [][]string{logAlwaysSkipDirs, exclude} {
		for _, pat := range list {
			if ok, _ := path.Match(strings.ToLower(pat), n); ok {
				return true
			}
		}
	}
	return false
}

// logPruneNames is the exclude list for the find command: the user's names plus the built-in ones.
func logPruneNames(exclude []string) []string {
	out := append([]string(nil), exclude...)
	seen := map[string]bool{}
	for _, e := range out {
		seen[strings.ToLower(e)] = true
	}
	for _, e := range logAlwaysSkipDirs {
		if !seen[e] {
			out = append(out, e)
		}
	}
	return out
}

const logReadDirTimeout = 15 * time.Second

// logReadDir lists a folder over SFTP but gives up on it after logReadDirTimeout (or when ctx
// ends), so one folder that does not answer (a hung network mount, a huge directory) is
// skipped with a warning instead of stopping the whole discovery.
func logReadDir(ctx context.Context, fs *sftp.Client, dir string) ([]os.FileInfo, error) {
	type res struct {
		e   []os.FileInfo
		err error
	}
	ch := make(chan res, 1)
	go func() {
		e, err := fs.ReadDir(dir)
		ch <- res{e, err}
	}()
	t := time.NewTimer(logReadDirTimeout)
	defer t.Stop()
	select {
	case r := <-ch:
		return r.e, r.err
	case <-t.C:
		return nil, errLogDirTimeout
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type logDirTimeoutError struct{}

func (logDirTimeoutError) Error() string { return "no answer, skipped" }

var errLogDirTimeout error = logDirTimeoutError{}

var logRediscoverLocks sync.Map // connection id → *sync.Mutex

// logRediscover runs "Find logs" again for the connection owner when the viewer or the search
// finds no current discovery (it expired, or the log settings changed), instead of failing
// with "Find logs first". One discovery runs per connection at a time; requests that arrive
// meanwhile wait for it. Shared connections never start one.
func logRediscover(r *http.Request, c Connection, acc *connAccess) bool {
	if acc == nil || acc.Share != nil {
		return false
	}
	v, _ := logRediscoverLocks.LoadOrStore(c.ID, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	if _, ok := monitorDiscoveryView(c); ok {
		return true // another request discovered meanwhile
	}
	start := time.Now()
	res, err := logDiscover(r.Context(), c)
	if err != nil {
		return false
	}
	auditFileOp(r, acc, c, "monitor.logs_discovered", c.Name, map[string]interface{}{"roots": len(res.Roots), "files": len(res.Files), "method": res.Method, "ms": time.Since(start).Milliseconds(), "automatic": true})
	_, ok := monitorDiscoveryView(c)
	return ok
}
