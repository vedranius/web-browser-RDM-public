package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ─── STORAGES: S3 ACTIVITY (v14.4.0) ──────────────────────────────────────────────────────
//
// Every HTTP request WRM sends to a bucket is recorded here, one entry per attempt (a retried
// request is two billable requests at a provider that bills them): the S3 operation, the key,
// the answer, the time it took, the bytes and what caused it (the user in the file window, the
// background check, a chain, an AI connection). The file window shows the list live, with the
// count per operation and what the provider bills for it (storage_pricing.go). Memory only: a
// restart starts the counts again, and nothing here holds a secret or object content.

const storageActivityKept = 300

type storageOp struct {
	At     string `json:"at"`
	Op     string `json:"op"`
	Method string `json:"method"`
	Key    string `json:"key,omitempty"`
	Status int    `json:"status"`
	Ms     int64  `json:"ms"`
	Bytes  int64  `json:"bytes,omitempty"`
	Source string `json:"source"`
	Error  string `json:"error,omitempty"`
}

type storageOpCount struct {
	Total int `json:"total"`
	Today int `json:"today"`
}

type storageActivityLog struct {
	ops    []storageOp
	counts map[string]*storageOpCount
	day    string
	since  string
}

var storageActivity = struct {
	sync.Mutex
	m map[int]*storageActivityLog
}{m: map[int]*storageActivityLog{}}

type storageSourceKey struct{}

// withStorageSource names what caused the S3 requests made with this context.
func withStorageSource(ctx context.Context, source string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, storageSourceKey{}, source)
}

func storageSourceOf(ctx context.Context) string {
	if ctx != nil {
		if s, ok := ctx.Value(storageSourceKey{}).(string); ok && s != "" {
			return s
		}
	}
	return "user"
}

// s3OpName maps a request to the S3 API operation name a provider's price list uses.
func s3OpName(method, key string, q url.Values, hdr http.Header) string {
	switch method {
	case http.MethodHead:
		if key == "" {
			return "HeadBucket"
		}
		return "HeadObject"
	case http.MethodGet:
		if key == "" && q.Get("list-type") == "2" {
			return "ListObjectsV2"
		}
		if key == "" {
			return "ListObjects"
		}
		return "GetObject"
	case http.MethodPut:
		if hdr.Get("x-amz-copy-source") != "" {
			return "CopyObject"
		}
		return "PutObject"
	case http.MethodDelete:
		return "DeleteObject"
	}
	return method
}

// recordStorageOp stores one request and tells the owner's open windows.
func recordStorageOp(cfg s3Config, op storageOp) {
	if cfg.StorageID <= 0 {
		return // a draft being tested is not a saved storage yet
	}
	now := time.Now()
	day := now.Format("2006-01-02")
	storageActivity.Lock()
	log := storageActivity.m[cfg.StorageID]
	if log == nil {
		log = &storageActivityLog{counts: map[string]*storageOpCount{}, day: day, since: now.UTC().Format(time.RFC3339)}
		storageActivity.m[cfg.StorageID] = log
	}
	if log.day != day {
		for _, c := range log.counts {
			c.Today = 0
		}
		log.day = day
	}
	c := log.counts[op.Op]
	if c == nil {
		c = &storageOpCount{}
		log.counts[op.Op] = c
	}
	c.Total++
	c.Today++
	log.ops = append(log.ops, op)
	if len(log.ops) > storageActivityKept {
		log.ops = append([]storageOp(nil), log.ops[len(log.ops)-storageActivityKept:]...)
	}
	storageActivity.Unlock()
	if cfg.OwnerID > 0 && hub != nil {
		b, _ := json.Marshal(map[string]interface{}{"type": "storage_op", "storage_id": cfg.StorageID, "op": op})
		hub.sendTo(cfg.OwnerID, b)
	}
}

// storageActivityView is the answer of GET /api/storages/{id}/activity.
func storageActivityView(s Storage) map[string]interface{} {
	storageActivity.Lock()
	defer storageActivity.Unlock()
	ops := []storageOp{}
	counts := map[string]storageOpCount{}
	since := ""
	if log := storageActivity.m[s.ID]; log != nil {
		ops = append(ops, log.ops...)
		today := time.Now().Format("2006-01-02")
		for k, v := range log.counts {
			c := *v
			if log.day != today {
				c.Today = 0
			}
			counts[k] = c
		}
		since = log.since
	}
	p := storageProviderFor(s)
	return map[string]interface{}{"storage_id": s.ID, "since": since, "ops": ops, "counts": counts,
		"provider": p.ID, "pricing": p, "estimate": storageCostEstimate(p, counts)}
}

func forgetStorageActivity(id int) {
	storageActivity.Lock()
	delete(storageActivity.m, id)
	storageActivity.Unlock()
}

// storageOpKey strips the storage prefix, so the list shows the key as the file window does.
func storageOpKey(cfg s3Config, key string) string {
	base := storageKeyPrefix(Storage{Prefix: cfg.Prefix})
	return truncateStr(strings.TrimPrefix(key, base), 300)
}

func storageOpError(err error) string {
	if err == nil {
		return ""
	}
	return truncateStr(fmt.Sprint(err), 160)
}
