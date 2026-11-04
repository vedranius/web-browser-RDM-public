package main

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"testing"
)

// TestStorageShareEscaping: object names with markup or quotes are escaped on the share page
// and cannot add parameters to the Content-Disposition header.
func TestStorageShareEscaping(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-escape", nil)
	evil := `pub/<script>x</script>"; filename=evil.exe.txt`
	stub.put("team-a/"+evil, []byte("x"))

	var created storageShare
	base := "/api/storages/" + strconv.Itoa(id) + "/shares"
	if code := u.jsonDo("POST", base, map[string]interface{}{"prefix": "pub", "hours": 1}, &created); code != http.StatusOK {
		t.Fatalf("create share: %d", code)
	}
	page, err := http.Get(u.srv.URL + created.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(page.Body)
	page.Body.Close()
	if page.StatusCode != http.StatusOK || strings.Contains(string(body), "<script") {
		t.Fatalf("share page: %d %s", page.StatusCode, body)
	}
	if !strings.Contains(string(body), "&lt;script&gt;") {
		t.Fatalf("object name missing from the page: %s", body)
	}

	dl, err := http.Get(u.srv.URL + created.URL + "/object?key=" + url.QueryEscape(strings.TrimPrefix(evil, "pub/")))
	if err != nil {
		t.Fatal(err)
	}
	dl.Body.Close()
	cd := dl.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	if dl.StatusCode != http.StatusOK || err != nil || len(params) != 1 || params["filename"] != path.Base(evil) {
		t.Fatalf("download: %d %q %v", dl.StatusCode, cd, params)
	}
}

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
