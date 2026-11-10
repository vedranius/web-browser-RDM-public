package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ─── ADDONS: READ ACCESS TO CHAINS (chains.read) ────────────────────────────────────────────
//
// An addon with the chains.read capability may read the chains of the person who looks at its
// window: names, nodes, roles, check definitions, the latest and the confirmed results, the
// active scenario and the timeline of changes. It can also ask for a check now, exactly like
// the Check now button (same rate limit, audited). It never gets connection ids, host names,
// commands, raw output or anything a chain does not already show its owner; Chains keep only
// extracted values and drop keys that look like secrets. A shared instance shows every viewer
// their own chains.

const addonChainTimeline = 200

type addonChainNode struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Label string `json:"label,omitempty"`
	Role  string `json:"role,omitempty"`
}

type addonChainCheck struct {
	ID     string `json:"id"`
	Node   string `json:"node"`
	Type   string `json:"type"`
	Label  string `json:"label,omitempty"`
	Target string `json:"target,omitempty"`
	Format string `json:"format,omitempty"`
	Key    string `json:"key,omitempty"`
}

type addonChainLink struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Kind  string `json:"kind,omitempty"`
	Label string `json:"label,omitempty"`
}

type addonChainScenario struct {
	ID    string `json:"id"`
	Label string `json:"label,omitempty"`
}

// addonChainView is what a panel sees of one chain.
func addonChainView(c chainRow) map[string]interface{} {
	nodes := make([]addonChainNode, 0, len(c.Def.Nodes))
	for _, n := range c.Def.Nodes {
		nodes = append(nodes, addonChainNode{ID: n.ID, Kind: n.Kind, Label: n.Label, Role: n.Role})
	}
	checks := make([]addonChainCheck, 0, len(c.Def.Checks))
	for _, k := range c.Def.Checks {
		checks = append(checks, addonChainCheck{ID: k.ID, Node: k.Node, Type: k.Type, Label: k.Label, Target: k.Target, Format: k.Format, Key: k.Key})
	}
	links := make([]addonChainLink, 0, len(c.Def.Links))
	for _, l := range c.Def.Links {
		links = append(links, addonChainLink{From: l.From, To: l.To, Kind: l.Kind, Label: l.Label})
	}
	scen := make([]addonChainScenario, 0, len(c.Def.Scenarios))
	for _, s := range c.Def.Scenarios {
		scen = append(scen, addonChainScenario{ID: s.ID, Label: s.Label})
	}
	return map[string]interface{}{"id": c.ID, "name": c.Def.Name, "description": c.Def.Description, "enabled": c.Enabled,
		"interval_s": int(c.Def.interval() / time.Second), "nodes": nodes, "checks": checks, "links": links, "scenarios": scen,
		"state": c.State, "confirmed": c.Confirmed, "scenario": c.Confirmed.Scenario,
		"scenario_label": c.Def.scenarioLabel(c.Confirmed.Scenario), "last_at": c.LastAt, "updated_at": c.UpdatedAt}
}

func addonChainTimelineOf(id int) []map[string]interface{} {
	out := []map[string]interface{}{}
	rows, err := db.Query(`SELECT at, kind, subject, old_value, new_value, ok FROM chain_samples WHERE chain_id=? ORDER BY id DESC LIMIT ?`, id, addonChainTimeline)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var at, kind, subject, oldV, newV string
		var okv int
		if rows.Scan(&at, &kind, &subject, &oldV, &newV, &okv) == nil {
			out = append(out, map[string]interface{}{"at": at, "kind": kind, "subject": subject, "old": oldV, "new": newV, "ok": okv == 1})
		}
	}
	return out
}

// addonChainsHandler serves /api/addons/{key}/instances/{iid}/chains:
//
//	GET  chains              the viewer's chains (summaries)
//	GET  chains/{id}         one chain with its timeline
//	POST chains/{id}/check   check it now (the Check now button's limit)
func addonChainsHandler(w http.ResponseWriter, r *http.Request, userID int, rec addonRecord, in addonInstance, rest []string) {
	if !addonHasCap(rec, "chains.read") {
		addonFail(w, addonErrf("E_CAPABILITY", "The addon does not declare the chains.read capability"), 403)
		return
	}
	if !settingBool("chains_enabled") {
		addonFail(w, addonErrf("E_CHAINS_OFF", "Chains are turned off by the administrator"), 403)
		return
	}
	if len(rest) == 0 {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		rows, err := db.Query(`SELECT id FROM chains WHERE user_id=? ORDER BY name COLLATE NOCASE`, userID)
		if err != nil {
			jsonError(w, "Server error", http.StatusInternalServerError)
			return
		}
		var ids []int
		for rows.Next() {
			var id int
			if rows.Scan(&id) == nil {
				ids = append(ids, id)
			}
		}
		rows.Close()
		out := []chainSummary{}
		for _, id := range ids {
			if c, err := loadChain(id); err == nil {
				out = append(out, summarizeChain(c))
			}
		}
		jsonOK(w, map[string]interface{}{"chains": out})
		return
	}
	id, err := strconv.Atoi(rest[0])
	if err != nil || len(rest) > 2 || (len(rest) == 2 && rest[1] != "check") {
		addonFail(w, addonErrf("E_NOT_FOUND", "Not found"), 404)
		return
	}
	c, err := loadChain(id)
	if err != nil || c.UserID != userID {
		addonFail(w, addonErrf("E_NOT_FOUND", "Chain not found"), 404)
		return
	}
	if len(rest) == 1 {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		jsonOK(w, map[string]interface{}{"chain": addonChainView(c), "timeline": addonChainTimelineOf(c.ID)})
		return
	}
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !mcpRateOK(fmt.Sprintf("chain-check:%d", userID), 20, time.Minute) {
		addonFail(w, addonErrf("E_LIMIT_RATE", "Too many checks; wait a moment"), http.StatusTooManyRequests)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	got, err := pollChain(ctx, c.ID)
	if err != nil {
		addonFail(w, addonErrf("E_CHAIN_BUSY", "%s", err.Error()), http.StatusConflict)
		return
	}
	auditLog(r, userID, "chain.checked", c.Def.Name, map[string]interface{}{"chain": c.ID, "addon": rec.Key, "instance": in.ID})
	jsonOK(w, map[string]interface{}{"chain": addonChainView(got), "timeline": addonChainTimelineOf(c.ID)})
}
