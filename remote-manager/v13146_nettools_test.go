package main

import "testing"

func TestHTTPNetToolRejectsURLCredentialsAndRedactsAuditTarget(t *testing.T) {
	if _, err := runHTTPCheck(nil, "https://alice:secret@example.test/health?token=private"); err == nil {
		t.Fatal("HTTP check accepted credentials embedded in the URL")
	}
	got := netToolAuditTarget("http", "https://alice:secret@example.test:8443/health?token=private#section")
	if want := "https://example.test:8443"; got != want {
		t.Fatalf("audit target = %q, want %q", got, want)
	}
	if got := netToolAuditTarget("http", "example.test/health?token=private"); got != "https://example.test" {
		t.Fatalf("default-scheme audit target = %q, want https://example.test", got)
	}
	if got := netToolAuditTarget("ports", "db.internal"); got != "db.internal" {
		t.Fatalf("non-HTTP audit target changed: %q", got)
	}
}
