package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func egressTestPolicy(t *testing.T, policy, allow string) {
	t.Helper()
	oldPolicy, oldAllow := getSetting("egress_policy"), getSetting("egress_allow_hosts")
	t.Cleanup(func() {
		setSetting("egress_policy", oldPolicy)
		setSetting("egress_allow_hosts", oldAllow)
	})
	if _, err := setSetting("egress_policy", policy); err != nil {
		t.Fatal(err)
	}
	if _, err := setSetting("egress_allow_hosts", allow); err != nil {
		t.Fatal(err)
	}
}

func TestEgressDialBlocksPrivateDNSAnswers(t *testing.T) {
	egressTestPolicy(t, "public_only", "")
	for _, ip := range []string{"127.0.0.1", "10.1.2.3", "fd00::1", "169.254.169.254", "100.100.100.200", "fe80::1"} {
		t.Run(ip, func(t *testing.T) {
			lookups, dials := 0, 0
			_, err := dialEgressContext(context.Background(), "tcp", "service.example:443",
				func(context.Context, string) ([]net.IPAddr, error) {
					lookups++
					return []net.IPAddr{{IP: net.ParseIP(ip)}}, nil
				}, func(context.Context, string, string) (net.Conn, error) {
					dials++
					return nil, nil
				})
			if err == nil || lookups != 1 || dials != 0 {
				t.Fatalf("DNS answer %s: error=%v lookups=%d dials=%d", ip, err, lookups, dials)
			}
		})
	}
}

func TestEgressDialPinsCheckedDNSAnswer(t *testing.T) {
	egressTestPolicy(t, "public_only", "")
	lookups := 0
	var dialled string
	left, right := net.Pipe()
	defer right.Close()
	conn, err := dialEgressContext(context.Background(), "tcp", "changing.example:443",
		func(context.Context, string) ([]net.IPAddr, error) {
			lookups++
			if lookups > 1 {
				return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil
			}
			return []net.IPAddr{{IP: net.ParseIP("203.0.113.10")}}, nil
		}, func(_ context.Context, _ string, address string) (net.Conn, error) {
			dialled = address
			return left, nil
		})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if lookups != 1 || dialled != "203.0.113.10:443" {
		t.Fatalf("connection was not pinned to the checked IP: lookups=%d address=%q", lookups, dialled)
	}
}

func TestEgressDialAllowsExplicitPrivateDNSAnswers(t *testing.T) {
	egressTestPolicy(t, "public_only", "internal.example,10.0.0.0/8")
	for _, tc := range []struct{ host, ip string }{{"internal.example", "192.168.1.5"}, {"other.example", "10.1.2.3"}} {
		left, right := net.Pipe()
		conn, err := dialEgressContext(context.Background(), "tcp", tc.host+":443",
			func(context.Context, string) ([]net.IPAddr, error) {
				return []net.IPAddr{{IP: net.ParseIP(tc.ip)}}, nil
			},
			func(context.Context, string, string) (net.Conn, error) { return left, nil })
		if err != nil {
			left.Close()
			right.Close()
			t.Fatalf("allowlisted %s (%s): %v", tc.host, tc.ip, err)
		}
		conn.Close()
		right.Close()
	}
	// A hostname allowlist never permits a cloud metadata endpoint.
	if err := guardEgressIP("internal.example", net.ParseIP("169.254.169.254")); err == nil {
		t.Fatal("allowlisted hostname permitted metadata")
	}
}

func TestEgressDialAllowsPrivateDNSInLocalPolicy(t *testing.T) {
	egressTestPolicy(t, "allow_private", "")
	left, right := net.Pipe()
	defer right.Close()
	conn, err := dialEgressContext(context.Background(), "tcp", "local-service.example:443",
		func(context.Context, string) ([]net.IPAddr, error) {
			return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil
		},
		func(context.Context, string, string) (net.Conn, error) { return left, nil })
	if err != nil {
		left.Close()
		t.Fatal(err)
	}
	conn.Close()
	if err := guardEgressIP("local-service.example", net.ParseIP("169.254.169.254")); err == nil {
		t.Fatal("allow_private permitted metadata through DNS")
	}
}

func TestEgressHTTPRejectsCrossOriginRedirect(t *testing.T) {
	egressTestPolicy(t, "allow_private", "")
	var targetHits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetHits.Add(1) }))
	defer target.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cross" {
			http.Redirect(w, r, target.URL, http.StatusFound)
			return
		}
		if r.URL.Path == "/same" {
			http.Redirect(w, r, "/done", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer origin.Close()
	tr := egressHTTPTransport()
	defer tr.CloseIdleConnections()
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second, CheckRedirect: guardEgressRedirect}
	if _, err := client.Get(origin.URL + "/cross"); err == nil || !strings.Contains(err.Error(), "another origin") {
		t.Fatalf("cross-origin redirect was followed: %v", err)
	}
	if targetHits.Load() != 0 {
		t.Fatalf("redirect target received %d requests", targetHits.Load())
	}
	resp, err := client.Get(origin.URL + "/same")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("same-origin redirect failed: %d", resp.StatusCode)
	}
	from, _ := http.NewRequest(http.MethodGet, "https://id.example/discovery", nil)
	to, _ := http.NewRequest(http.MethodGet, "http://id.example/metadata", nil)
	if err := guardEgressRedirect(to, []*http.Request{from}); err == nil {
		t.Fatal("HTTPS-to-HTTP redirect was accepted")
	}
}

func TestOpencodeDownloadKeepsGuardedMirrorRedirects(t *testing.T) {
	egressTestPolicy(t, "allow_private", "")
	var tgz bytes.Buffer
	gz := gzip.NewWriter(&tgz)
	tw := tar.NewWriter(gz)
	binary := []byte("test opencode binary")
	if err := tw.WriteHeader(&tar.Header{Name: "package/bin/opencode", Mode: 0o755, Size: int64(len(binary)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	tw.Write(binary)
	tw.Close()
	gz.Close()
	sum := sha512.Sum512(tgz.Bytes())
	integrity := "sha512-" + base64.StdEncoding.EncodeToString(sum[:])
	var downloads atomic.Int32
	tarball := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloads.Add(1)
		w.Write(tgz.Bytes())
	}))
	defer tarball.Close()
	var registry *httptest.Server
	registry = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pkg.tgz" {
			http.Redirect(w, r, tarball.URL+"/pkg.tgz", http.StatusFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"version": "1", "dist": map[string]string{"tarball": registry.URL + "/pkg.tgz", "integrity": integrity}})
	}))
	defer registry.Close()
	oldRegistry := getSetting("ai_opencode_registry")
	t.Cleanup(func() { setSetting("ai_opencode_registry", oldRegistry) })
	if _, err := setSetting("ai_opencode_registry", registry.URL); err != nil {
		t.Fatal(err)
	}
	oldClient := aiHTTPClient
	tr := egressHTTPTransport()
	tr.TLSClientConfig = registry.Client().Transport.(*http.Transport).TLSClientConfig.Clone()
	aiHTTPClient = &http.Client{Transport: tr, CheckRedirect: guardEgressRedirect}
	t.Cleanup(func() { aiHTTPClient = oldClient; tr.CloseIdleConnections() })
	if err := os.MkdirAll(opencodeRoot(), 0o700); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "opencode")
	if _, err := opencodeFetch(context.Background(), "opencode-linux-amd64", "1", dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil || !bytes.Equal(got, binary) || downloads.Load() != 1 {
		t.Fatalf("mirror redirect download: data=%q error=%v requests=%d", got, err, downloads.Load())
	}
	from, _ := http.NewRequest(http.MethodGet, "https://registry.example/pkg.tgz", nil)
	to, _ := http.NewRequest(http.MethodGet, "http://registry.example/pkg.tgz", nil)
	if err := opencodeDownloadRedirect(to, []*http.Request{from}); err == nil {
		t.Fatal("installer accepted an HTTPS downgrade")
	}
}
