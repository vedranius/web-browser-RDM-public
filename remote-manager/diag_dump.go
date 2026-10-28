package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
)

// Freeze diagnostics. When WRM stops answering, the operator usually restarts the
// service and nothing in the log says what was stuck. WRM therefore writes the stack of
// every goroutine to <DB_PATH>.diag/ on every stop or restart and when its own watchdog
// sees /healthz or the database stall. The files hold code locations only (no request
// bodies, no secrets) and only the newest diagMaxFiles are kept.

const diagMaxFiles = 20

func diagDir() string { return resolveDBPath() + ".diag" }

// writeDiagDump writes all goroutine stacks to <DB_PATH>.diag/<kind>-<time>.txt and
// returns the file name ("" when it could not be written).
func writeDiagDump(kind string) string {
	dir := diagDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		log.Printf("Diagnostics: %v", err)
		return ""
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		log.Printf("Diagnostics: %v", err)
		return ""
	}
	defer root.Close()
	base := kind + "-" + time.Now().UTC().Format("20060102-150405") + ".txt"
	name := filepath.Join(dir, base)
	f, err := root.OpenFile(base, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		log.Printf("Diagnostics: %v", err)
		return ""
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Fprintf(f, "WRM %s %s/%s, %s (%s)\ngoroutines: %d, heap in use: %d MiB, system memory: %d MiB, open DB connections: %d (in use %d, waiting %d)\n\n",
		AppVersion, runtime.GOOS, runtime.GOARCH, time.Now().Format(time.RFC3339), kind,
		runtime.NumGoroutine(), ms.HeapInuse>>20, ms.Sys>>20, db.Stats().OpenConnections, db.Stats().InUse, db.Stats().WaitCount)
	if p := pprof.Lookup("goroutine"); p != nil {
		p.WriteTo(f, 2)
	}
	f.Close()
	pruneDiagDumps(root)
	log.Printf("Diagnostics: %d goroutines written to %s", runtime.NumGoroutine(), name)
	return name
}

func pruneDiagDumps(root *os.Root) {
	entries, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".txt") {
			names = append(names, e.Name())
		}
	}
	if len(names) <= diagMaxFiles {
		return
	}
	// Sort by time (the part after the kind), oldest first.
	sort.Slice(names, func(i, j int) bool { return diagStamp(names[i]) < diagStamp(names[j]) })
	for _, n := range names[:len(names)-diagMaxFiles] {
		root.Remove(n)
	}
}

func diagStamp(name string) string {
	if i := strings.Index(name, "-"); i >= 0 {
		return name[i+1:]
	}
	return name
}

// hangCheck reports why WRM looks stuck: its own /healthz does not answer in time, or
// the database cannot be reached. "" means healthy.
var hangCheck = func() string {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return "database: " + err.Error()
	}
	url, client, ok := localHealthURL(20 * time.Second)
	if !ok {
		return ""
	}
	resp, err := client.Get(url)
	if err != nil {
		return "web server: " + err.Error()
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "web server: " + resp.Status
	}
	return ""
}

// hangWatchdog checks every 30 s and writes a dump after two failed checks in a row, at
// most once every 15 minutes, so a freeze leaves evidence even when nobody restarts WRM.
func hangWatchdog() {
	time.Sleep(time.Minute) // let start-up finish
	var fails int
	var last time.Time
	for {
		if why := hangCheck(); why != "" {
			fails++
			log.Printf("Watchdog: WRM is not responding (%s)", why)
			if fails >= 2 && time.Since(last) > 15*time.Minute {
				writeDiagDump("hang")
				last = time.Now()
			}
		} else {
			fails = 0
		}
		time.Sleep(30 * time.Second)
	}
}
