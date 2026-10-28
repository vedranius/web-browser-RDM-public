package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
