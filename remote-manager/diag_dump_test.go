package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDiagDumpWrittenAndPruned(t *testing.T) {
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "diag.db"))
	dir := diagDir()
	os.MkdirAll(dir, 0o700)
	for i := 0; i < diagMaxFiles+3; i++ {
		os.WriteFile(filepath.Join(dir, "stop-20200101-0000"+string(rune('a'+i))+".txt"), []byte("x"), 0o600)
	}
	name := writeDiagDump("hang")
	if name == "" {
		t.Fatal("no dump written")
	}
	data, err := os.ReadFile(name)
	if err != nil || !strings.Contains(string(data), "goroutines:") || !strings.Contains(string(data), "TestDiagDumpWrittenAndPruned") {
		t.Fatalf("dump content: %v\n%.300s", err, data)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != diagMaxFiles {
		t.Fatalf("%d files kept, want %d", len(entries), diagMaxFiles)
	}
	if _, err := os.Stat(name); err != nil {
		t.Fatalf("newest dump was pruned: %v", err)
	}
}

func TestLocalHealthURL(t *testing.T) {
	for addr, want := range map[string]string{
		":8080":             "http://127.0.0.1:8080/healthz",
		"0.0.0.0:9000":      "http://127.0.0.1:9000/healthz",
		"[::]:9000":         "http://127.0.0.1:9000/healthz",
		"192.168.5.10:8780": "http://192.168.5.10:8780/healthz",
		"[fd00::5]:8780":    "http://[fd00::5]:8780/healthz",
	} {
		t.Setenv("LISTEN_ADDR", addr)
		if got, _, ok := localHealthURL(time.Second); !ok || got != want {
			t.Errorf("%s: %q, want %q", addr, got, want)
		}
	}
}
