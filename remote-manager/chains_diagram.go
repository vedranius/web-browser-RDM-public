package main

import (
	"encoding/json"
	"math"
	"net/http"
)

// ─── CHAINS: LIVE DIAGRAM (step 2) ────────────────────────────────────────────────────────
//
// The diagram itself is drawn in the browser (static/chains.js). The server keeps the node
// positions the user dragged (x and y on each node of the definition; 0 means "place it
// automatically") and offers generic starter templates.

const chainPosMax = 20000

// validChainPos reports whether a saved coordinate is usable: finite and inside the canvas.
func validChainPos(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 && v <= chainPosMax
}

type chainLayoutReq struct {
	Positions map[string]struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"positions"`
	Reset bool `json:"reset"`
}

// apiChainLayout saves node positions only: the rest of the definition is untouched, so a
// drag never needs (or can bypass) the full validation of the chain.
func apiChainLayout(w http.ResponseWriter, r *http.Request, c chainRow) {
	if !requireMethod(w, r, http.MethodPut) {
		return
	}
	var in chainLayoutReq
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		jsonError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if len(in.Positions) > chainMaxNodes {
		jsonError(w, "Too many positions", http.StatusBadRequest)
		return
	}
	d := c.Def
	changed := 0
	for i := range d.Nodes {
		n := &d.Nodes[i]
		if in.Reset {
			n.X, n.Y = 0, 0
			changed++
			continue
		}
		p, ok := in.Positions[n.ID]
		if !ok {
			continue
		}
		if !validChainPos(p.X) || !validChainPos(p.Y) {
			jsonError(w, "x and y must be between 0 and 20000", http.StatusBadRequest)
			return
		}
		n.X, n.Y = math.Round(p.X), math.Round(p.Y)
		changed++
	}
	if changed > 0 {
		if _, err := db.Exec(`UPDATE chains SET definition=?, updated_at=? WHERE id=? AND user_id=?`,
			string(jsonMarshal(d)), nowRFC(), c.ID, c.UserID); err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
	}
	got, err := loadChain(c.ID)
	if err != nil {
		jsonError(w, "Not found", http.StatusNotFound)
		return
	}
	jsonOK(w, chainView(got))
}

// chainTemplate is a generic starting point. Server and storage nodes come without a
// connection or storage; the user picks them before the chain is saved.
type chainTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Definition  chainDef `json:"definition"`
}

func chainTemplates() []chainTemplate {
	return []chainTemplate{
		{
			ID: "single", Name: "One service",
			Description: "One server with a service and its port, and a public HTTPS check.",
			Definition: chainDef{
				Name: "My service", IntervalS: 30, Confirm: 2,
				Nodes: []chainNode{
					{ID: "app", Kind: "server", Label: "demo-app", Role: "app"},
					{ID: "site", Kind: "external", Label: "Public site", Role: "endpoint"},
				},
				Checks: []chainCheck{
					{ID: "web", Node: "app", Type: "systemd", Target: "nginx.service", Label: "Web server"},
					{ID: "https", Node: "app", Type: "port", Port: 443, Label: "HTTPS port"},
					{ID: "health", Node: "site", Type: "http", Target: "https://www.example.com/", Label: "Site answers"},
					{ID: "cert", Node: "site", Type: "tls", Target: "www.example.com", Label: "Certificate"},
				},
				Links: []chainLink{{From: "app", To: "site", Kind: "data", Check: "health"}},
				Scenarios: []chainScenario{
					{ID: "normal", Label: "Normal", When: chainCond{All: []chainCond{
						{Check: "web", Field: "ok", Op: "==", Value: "true"},
						{Check: "health", Field: "ok", Op: "==", Value: "true"},
					}}},
					{ID: "degraded", Label: "Degraded", When: chainCond{}},
				},
			},
		},
		{
			ID: "failover", Name: "Primary and backup path",
			Description: "Ingest on a primary and a backup server, a content server writing to a bucket, and a CDN in front.",
			Definition: chainDef{
				Name: "Media chain", IntervalS: 30, Confirm: 2,
				Nodes: []chainNode{
					{ID: "ingest-a", Kind: "server", Label: "demo-ingest-a", Role: "ingest (primary)"},
					{ID: "ingest-b", Kind: "server", Label: "demo-ingest-b", Role: "ingest (backup)"},
					{ID: "content", Kind: "server", Label: "demo-content", Role: "content"},
					{ID: "bucket", Kind: "storage", Label: "Recordings", Role: "storage"},
					{ID: "cdn", Kind: "external", Label: "CDN", Role: "delivery"},
				},
				Checks: []chainCheck{
					{ID: "ingest-a-svc", Node: "ingest-a", Type: "systemd", Target: "demo-ingest.service", Label: "Ingest service"},
					{ID: "ingest-b-svc", Node: "ingest-b", Type: "systemd", Target: "demo-ingest.service", Label: "Ingest service"},
					{ID: "content-svc", Node: "content", Type: "systemd", Target: "demo-content.service", Label: "Content service"},
					{ID: "content-mount", Node: "content", Type: "mount", Target: "/mnt/recordings", Label: "Bucket mount"},
					{ID: "bucket-new", Node: "bucket", Type: "storage", Prefix: "recordings/", MaxAgeS: 900, Label: "New recordings"},
					{ID: "cdn-health", Node: "cdn", Type: "http", Target: "https://cdn.example.com/health", Label: "CDN health"},
				},
				Links: []chainLink{
					{From: "ingest-a", To: "content", Kind: "data", Check: "ingest-a-svc", Label: "primary", Scenarios: []string{"primary"}},
					{From: "ingest-b", To: "content", Kind: "failover", Check: "ingest-b-svc", Label: "backup", Scenarios: []string{"backup"}},
					{From: "content", To: "bucket", Kind: "write", Check: "bucket-new"},
					{From: "content", To: "cdn", Kind: "data", Check: "cdn-health"},
				},
				Scenarios: []chainScenario{
					{ID: "primary", Label: "Primary", When: chainCond{All: []chainCond{
						{Node: "ingest-a", Op: "==", Value: "up"},
						{Check: "ingest-a-svc", Field: "ok", Op: "==", Value: "true"},
					}}},
					{ID: "backup", Label: "Ingest on backup", When: chainCond{Check: "ingest-b-svc", Field: "ok", Op: "==", Value: "true"}},
					{ID: "down", Label: "No ingest", When: chainCond{}},
				},
			},
		},
	}
}

func apiChainTemplates(w http.ResponseWriter, r *http.Request) {
	if _, ok := requireAuth(w, r); !ok || !requireMethod(w, r, http.MethodGet) {
		return
	}
	jsonOK(w, chainTemplates())
}
