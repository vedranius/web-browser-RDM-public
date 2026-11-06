package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ─── MCP: CHAIN TOOLS (v14.4.0) ───────────────────────────────────────────────────────────
//
// An AI connection with the `chains` scope reads, tests and drafts the user's chains without a
// server session: the checks are WRM's own read-only catalog over the connections the chain
// already uses. The scope is not in the default policy, so an administrator turns it on
// (Settings → AI → allowed scopes) before a user can give it to an AI connection. A draft never
// changes a monitored chain: the user applies it in the Chains window. Read-only connections
// may read and test, not draft.

var mcpChainToolNames = []string{"chain_status", "chain_get", "chain_test", "chain_draft"}

func mcpChainTools(g *mcpGrant) []map[string]interface{} {
	if !mcpHasScope(g, "chains") || !settingBool("chains_enabled") {
		return nil
	}
	schemas := map[string]map[string]interface{}{
		"chain_status": obj(map[string]interface{}{}), "chain_get": chainGetSchema(),
		"chain_test": chainTestSchema(), "chain_draft": chainDraftSchemaJSON(),
	}
	titles := map[string]string{"chain_status": "List chains", "chain_get": "Read a chain", "chain_test": "Test a chain once", "chain_draft": "Draft a chain or a change"}
	out := []map[string]interface{}{}
	for _, name := range mcpChainToolNames {
		if name == "chain_draft" && g.Mode == aiModeReadOnly {
			continue
		}
		ann := map[string]interface{}{"readOnlyHint": name != "chain_draft", "openWorldHint": false}
		if name == "chain_draft" {
			ann["destructiveHint"] = false
		}
		out = append(out, map[string]interface{}{"name": name, "title": titles[name], "description": chainToolDesc[name], "inputSchema": schemas[name], "annotations": ann})
	}
	return out
}

func isMCPChainTool(name string) bool { return oneOf(name, mcpChainToolNames...) }

// mcpChainCall runs one chain tool for an AI connection and answers the RPC call.
func mcpChainCall(w http.ResponseWriter, r *http.Request, g *mcpGrant, m rpcMsg, name string, args map[string]interface{}) {
	fail := func(msg string) {
		auditLog(r, g.UserID, "mcp.chain_denied", name, map[string]interface{}{"mcp_grant": g.ID, "tool": name, "error": msg})
		writeRPC(w, 200, rpcResult(m.ID, toolText(msg, true)))
	}
	if !mcpHasScope(g, "chains") {
		fail("This AI connection has no chains scope. An administrator allows the scope (Settings → AI → allowed scopes), then the user enables it on the AI connection in WRM.")
		return
	}
	if name == "chain_draft" && g.Mode == aiModeReadOnly {
		fail("This AI connection is read-only: it may read and test chains, not draft them.")
		return
	}
	if !mcpRateOK(fmt.Sprintf("mcp-chain:%d", g.ID), 60, time.Minute) {
		fail("Too many chain calls; wait a moment.")
		return
	}
	raw, _ := json.Marshal(args)
	var in aiToolInput
	if err := json.Unmarshal(raw, &in); err != nil {
		fail("invalid arguments: " + truncateStr(err.Error(), 200))
		return
	}
	if !mcpStartToolCall(w, m.ID) {
		return
	}
	defer func() { <-mcpToolSlots }()
	ctx, cancel := context.WithTimeout(r.Context(), 100*time.Second)
	defer cancel()
	out, err := runChainAITool(ctx, g.UserID, "mcp:"+strconv.FormatInt(g.ID, 10), name, in, "mcp")
	if err != nil {
		fail(err.Error())
		return
	}
	auditLog(r, g.UserID, "mcp."+name, "", map[string]interface{}{"mcp_grant": g.ID, "chain": in.ChainID, "draft": in.DraftID})
	b, _ := json.MarshalIndent(out, "", "  ")
	writeRPC(w, 200, rpcResult(m.ID, toolText(string(b), false)))
}
