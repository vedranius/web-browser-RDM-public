package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The origin check compares the port too: another service on the same host (for example a
// web interface published by a tunnel on 0.0.0.0) is a different origin.
func TestV1471OriginPortMatters(t *testing.T) {
	t.Setenv("WRM_ALLOW_ANY_ORIGIN", "")
	t.Setenv("WRM_ALLOWED_ORIGINS", "")
	cases := []struct {
		host, xfh, xfp, origin string
		want                   bool
	}{
		{"wrm.example.com:8780", "", "", "http://wrm.example.com:8780", true},
		{"wrm.example.com:8780", "", "", "http://wrm.example.com:9999", false},
		{"wrm.example.com:8780", "", "", "http://WRM.example.com:8780/path", true},
		{"wrm.example.com", "", "", "https://wrm.example.com", true},
		{"wrm.example.com", "", "", "http://wrm.example.com", true},
		{"wrm.example.com", "", "", "https://wrm.example.com:44321", false},
		{"wrm.example.com", "", "8443", "https://wrm.example.com:8443", true},
		{"127.0.0.1:8780", "wrm.example.com", "", "https://wrm.example.com", true},
		{"127.0.0.1:8780", "wrm.example.com", "", "https://wrm.example.com:9443", false},
		{"127.0.0.1:8780", "wrm.example.com:9443", "", "https://wrm.example.com:9443", true},
		{"[::1]:8780", "", "", "http://[::1]:8780", true},
		{"[::1]:8780", "", "", "http://[::1]:8781", false},
		{"wrm.example.com:8780", "", "", "http://other.example.com:8780", false},
		{"wrm.example.com:8780", "", "", "not a url", false},
	}
	for _, c := range cases {
		r := httptest.NewRequest("GET", "/ws/ssh", nil)
		r.Host = c.host
		if c.xfh != "" {
			r.Header.Set("X-Forwarded-Host", c.xfh)
		}
		if c.xfp != "" {
			r.Header.Set("X-Forwarded-Port", c.xfp)
		}
		if got := originAllowed(r, c.origin); got != c.want {
			t.Errorf("host %q xfh %q xfp %q origin %q: got %v, want %v", c.host, c.xfh, c.xfp, c.origin, got, c.want)
		}
	}
}

func TestV1471AllowedOriginsList(t *testing.T) {
	t.Setenv("WRM_ALLOW_ANY_ORIGIN", "")
	t.Setenv("WRM_ALLOWED_ORIGINS", "legacy.example.com, https://exact.example.com:8443")
	r := httptest.NewRequest("GET", "/ws/ssh", nil)
	r.Host = "127.0.0.1:8780"
	for origin, want := range map[string]bool{
		"https://legacy.example.com:1234": true, // bare host name: any port, as before
		"https://exact.example.com:8443":  true,
		"https://exact.example.com":       false,
		"http://exact.example.com:8443":   false,
	} {
		if got := originAllowed(r, origin); got != want {
			t.Errorf("%s: got %v, want %v", origin, got, want)
		}
	}
}

// #190: the administrator can turn HSTS off; WRM then sends max-age=0 so browsers forget it.
func TestV1471HSTSSetting(t *testing.T) {
	h := securityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	get := func() string {
		r := httptest.NewRequest("GET", "https://wrm.example.com/", nil)
		r.TLS = &tls.ConnectionState{}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w.Header().Get("Strict-Transport-Security")
	}
	setSetting("hsts", "1")
	if got := get(); got != "max-age=15552000" {
		t.Fatalf("hsts on: %q", got)
	}
	setSetting("hsts", "0")
	defer setSetting("hsts", "1")
	if got := get(); got != "max-age=0" {
		t.Fatalf("hsts off: %q", got)
	}
	if ips := wrmHostIPs(); ips == nil {
		t.Fatal("wrmHostIPs returned nil")
	}
}
