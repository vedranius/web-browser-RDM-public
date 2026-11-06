package main

import (
	"strings"
	"testing"
)

// TestChainPortLabel: well-known ports get their service name in check labels.
func TestChainPortLabel(t *testing.T) {
	cases := map[int]string{5432: "PostgreSQL (5432)", 443: "HTTPS (443)", 8123: "Port 8123"}
	for p, want := range cases {
		if got := chainPortLabel(p); got != want {
			t.Errorf("chainPortLabel(%d) = %q, want %q", p, got, want)
		}
	}
}

// TestChainStarterNamesPorts: the deterministic starter labels a discovered port by its service.
func TestChainStarterNamesPorts(t *testing.T) {
	conns := []Connection{{ID: 7, Name: "demo-db"}}
	facts := []chainServerFacts{{ConnID: 7, Name: "demo-db", Ports: []int{5432}}}
	def := chainStarterFromSelection(0, conns, nil, facts, "", "")
	found := false
	for _, c := range def.Checks {
		if c.Type == "port" && c.Port == 5432 {
			found = true
			if !strings.HasPrefix(c.Label, "PostgreSQL") {
				t.Errorf("port check label %q, want the service name", c.Label)
			}
		}
	}
	if !found {
		t.Fatalf("no port check for 5432 in %+v", def.Checks)
	}
}
