package main

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestS3OpNames(t *testing.T) {
	cp := http.Header{}
	cp.Set("x-amz-copy-source", "/b/a")
	cases := []struct {
		method, key string
		q           url.Values
		hdr         http.Header
		want        string
	}{
		{"HEAD", "", nil, nil, "HeadBucket"},
		{"HEAD", "a.txt", nil, nil, "HeadObject"},
		{"GET", "", url.Values{"list-type": {"2"}}, nil, "ListObjectsV2"},
		{"GET", "a.txt", nil, nil, "GetObject"},
		{"PUT", "a.txt", nil, http.Header{}, "PutObject"},
		{"PUT", "b.txt", nil, cp, "CopyObject"},
		{"DELETE", "a.txt", nil, nil, "DeleteObject"},
	}
	for _, c := range cases {
		hdr := c.hdr
		if hdr == nil {
			hdr = http.Header{}
		}
		q := c.q
		if q == nil {
			q = url.Values{}
		}
		if got := s3OpName(c.method, c.key, q, hdr); got != c.want {
			t.Errorf("%s %q: %s, want %s", c.method, c.key, got, c.want)
		}
	}
}

func TestStorageProviderPricing(t *testing.T) {
	for host, want := range map[string]string{
		"https://fsn1.your-objectstorage.com": "hetzner", "https://fra1.digitaloceanspaces.com": "digitalocean",
		"https://s3.eu-central-1.amazonaws.com": "aws", "https://s3.eu-central-003.backblazeb2.com": "backblaze",
		"https://abc.r2.cloudflarestorage.com": "r2", "https://minio.example.com": "other",
	} {
		if got := guessStorageProvider(host); got != want {
			t.Errorf("%s: %s, want %s", host, got, want)
		}
	}
	for _, p := range storageProviders {
		if p.Summary["en"] == "" || p.Summary["hr"] == "" {
			t.Errorf("%s has no summary in both languages", p.ID)
		}
		if p.Billed == "yes" {
			for _, op := range storageWRMOps {
				if _, ok := p.Classes[p.Ops[op]]; !ok {
					t.Errorf("%s: %s has no price class", p.ID, op)
				}
			}
		}
	}
	aws, _ := storageProviderByID("aws")
	est := storageCostEstimate(aws, map[string]storageOpCount{"ListObjectsV2": {Total: 1000}, "DeleteObject": {Total: 50}})
	if c := est["cost_total"].(float64); c < 0.00499 || c > 0.00501 {
		t.Fatalf("1,000 LIST at AWS: %v USD, want 0.005", c)
	}
	hz, _ := storageProviderByID("hetzner")
	if c := storageCostEstimate(hz, map[string]storageOpCount{"ListObjectsV2": {Total: 1000}})["cost_total"].(float64); c != 0 {
		t.Fatalf("Hetzner requests cost %v", c)
	}
}

// TestStorageActivityRecordsEveryRequest: each request to the bucket appears in the activity
// view with its operation and cause, and the counts add up.
func TestStorageActivityRecordsEveryRequest(t *testing.T) {
	u, stub, id := storageTestEnv(t, "storage-activity", nil)
	stub.put("team-a/a.txt", []byte("abc"))
	forgetStorageActivity(id)
	base := "/api/storages/" + strconv.Itoa(id)
	if code := u.jsonDo("GET", base+"/objects?prefix=", nil, nil); code != http.StatusOK {
		t.Fatalf("list: %d", code)
	}
	if code := u.jsonDo("GET", base+"/text?key=a.txt", nil, nil); code != http.StatusOK {
		t.Fatalf("read: %d", code)
	}
	var view struct {
		Provider string                    `json:"provider"`
		Ops      []storageOp               `json:"ops"`
		Counts   map[string]storageOpCount `json:"counts"`
	}
	if code := u.jsonDo("GET", base+"/activity", nil, &view); code != http.StatusOK {
		t.Fatalf("activity: %d", code)
	}
	if len(view.Ops) != 2 || view.Ops[0].Op != "ListObjectsV2" || view.Ops[1].Op != "GetObject" || view.Ops[1].Key != "a.txt" || view.Ops[0].Source != "user" {
		t.Fatalf("ops: %+v", view.Ops)
	}
	if view.Counts["ListObjectsV2"].Total != 1 || view.Counts["GetObject"].Today != 1 {
		t.Fatalf("counts: %+v", view.Counts)
	}
	if view.Provider != "other" {
		t.Fatalf("provider of a local endpoint: %q", view.Provider)
	}
	other := newTestUser(t, u.srv, "storage-activity-other", false)
	if code := other.jsonDo("GET", base+"/activity", nil, nil); code != http.StatusNotFound {
		t.Fatalf("another user's activity: %d", code)
	}
}
