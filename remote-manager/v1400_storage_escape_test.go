package main

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"
)

func TestContentDisposition(t *testing.T) {
	cases := map[string]string{
		"a/report.json":     `attachment; filename=report.json`,
		`a/x"; y=1.txt`:     `attachment; filename="x\"; y=1.txt"`,
		"a/čćž titlovi.srt": `attachment; filename*=utf-8''%C4%8D%C4%87%C5%BE%20titlovi.srt`,
	}
	for key, want := range cases {
		if got := contentDisposition("attachment", key); got != want {
			t.Errorf("%q: got %q, want %q", key, got, want)
		}
	}
}

// TestStorageEditorLimit: the editor opens a typical subtitle or JSON file by default and
// refuses one above storage_edit_kb with 413.
func TestStorageEditorLimit(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-limit", nil)
	stub.put("team-a/movie.srt", bytes.Repeat([]byte("1\n00:00:01,000 --> 00:00:02,000\nLine\n\n"), 3000)) // ~120 KB
	base := "/api/storages/" + strconv.Itoa(id)
	if code := u.jsonDo("GET", base+"/text?key=movie.srt", nil, nil); code != http.StatusOK {
		t.Fatalf("a 120 KB subtitle: %d, want 200", code)
	}
	if _, err := setSetting("storage_edit_kb", "64"); err != nil {
		t.Fatal(err)
	}
	defer setSetting("storage_edit_kb", "2048")
	if code := u.jsonDo("GET", base+"/text?key=movie.srt", nil, nil); code != http.StatusRequestEntityTooLarge {
		t.Fatalf("above the limit: %d, want 413", code)
	}
}
