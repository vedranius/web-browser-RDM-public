package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// ─── MCP: CLAUDE DESKTOP EXTENSION (.mcpb) ───────────
//
//   POST /api/mcp/claude-desktop {token}  →  wrm-pro.mcpb
//
// Claude Desktop starts local (stdio) MCP servers only, and its custom connectors connect
// from Anthropic's cloud, which cannot reach a WRM on a private network or without a public
// HTTPS certificate. The one-click way in is an MCP bundle (a zip with manifest.json) that
// Claude Desktop installs with a double click: its server is the small Node.js bridge in
// mcpb/server/index.js (Claude Desktop ships the Node.js runtime), which talks to this WRM's
// /mcp over HTTP with the token. The address of this WRM is filled in; the token the wizard
// just created can be filled in too (it is then stored by Claude Desktop in the operating
// system's secure storage). The bundle carries nothing else of the user.

//go:embed mcpb/server/index.js
var mcpbBridgeJS string

const mcpbTokenMax = 512

func mcpbVersion() string {
	v := strings.TrimPrefix(AppVersion, "v")
	if v == "" {
		return "0.0.0"
	}
	return v
}

// mcpbManifest is the bundle's manifest.json (MCP bundle manifest version 0.2).
func mcpbManifest(mcpURL, token string) map[string]interface{} {
	urlCfg := map[string]interface{}{
		"type": "string", "title": "WRM MCP address", "required": true, "default": mcpURL,
		"description": "The MCP address of your WRM PRO (Settings → AI connections), for example https://wrm.example.com/mcp.",
	}
	tokCfg := map[string]interface{}{
		"type": "string", "title": "WRM token", "required": true, "sensitive": true,
		"description": "The token of an AI connection (wrm_pat_…), created in WRM → Settings → AI connections.",
	}
	if token != "" {
		tokCfg["default"] = token
	}
	return map[string]interface{}{
		"manifest_version": "0.2",
		"name":             "wrm-pro",
		"display_name":     "WRM PRO",
		"version":          mcpbVersion(),
		"description":      "Work on your servers through WRM PRO: logs, files, services and commands, under WRM's permissions and approvals.",
		"long_description": mcpServerDescription,
		"author":           map[string]interface{}{"name": "WRM PRO", "url": "https://github.com/vedranius/web-browser-RDM-public"},
		"homepage":         "https://github.com/vedranius/web-browser-RDM-public",
		"license":          "LicenseRef-WRM-PRO",
		"icon":             "icon.png",
		"server": map[string]interface{}{
			"type":        "node",
			"entry_point": "server/index.js",
			"mcp_config": map[string]interface{}{
				"command": "node",
				"args":    []string{"${__dirname}/server/index.js"},
				"env": map[string]string{
					"WRM_MCP_URL":     "${user_config.wrm_url}",
					"WRM_MCP_TOKEN":   "${user_config.wrm_token}",
					"WRM_MCP_CA_FILE": "${user_config.wrm_ca_file}",
				},
			},
		},
		"user_config": map[string]interface{}{
			"wrm_url":   urlCfg,
			"wrm_token": tokCfg,
			"wrm_ca_file": map[string]interface{}{
				"type": "file", "title": "CA certificate (optional)", "required": false,
				"description": "Only for a WRM with a self-signed or private-CA HTTPS certificate: the PEM file to trust.",
			},
		},
		"compatibility": map[string]interface{}{
			"platforms": []string{"darwin", "win32", "linux"},
			"runtimes":  map[string]string{"node": ">=18.0.0"},
		},
		"keywords": []string{"wrm", "ssh", "servers", "logs", "devops"},
	}
}

// buildMCPB returns the bundle: manifest.json, server/index.js and icon.png.
func buildMCPB(mcpURL, token string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	add := func(name string, data []byte) error {
		f, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, Modified: time.Now()})
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		return err
	}
	man, _ := json.MarshalIndent(mcpbManifest(mcpURL, token), "", "  ")
	js := strings.Replace(mcpbBridgeJS, "const VERSION = '0.0.0';", "const VERSION = '"+mcpbVersion()+"';", 1)
	icon, _ := staticFiles.ReadFile("static/brand/icon-512.png")
	for _, f := range []struct {
		name string
		data []byte
	}{{"manifest.json", man}, {"server/index.js", []byte(js)}, {"icon.png", icon}} {
		if err := add(f.name, f.data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mcpClaudeDesktopHandler(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	if ok, why := mcpUserAllowed(userID); !ok {
		jsonError(w, why, 403)
		return
	}
	var q struct {
		Token string `json:"token"`
	}
	if r.ContentLength != 0 {
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&q); err != nil {
			jsonError(w, "Bad JSON", 400)
			return
		}
	}
	q.Token = strings.TrimSpace(q.Token)
	// only a token of this WRM may be filled in; the bundle never carries anything else
	if q.Token != "" && (len(q.Token) > mcpbTokenMax || !strings.HasPrefix(q.Token, mcpPATPrefix) || strings.ContainsAny(q.Token, " \t\r\n\"\\")) {
		jsonError(w, "Not a WRM token", 400)
		return
	}
	data, err := buildMCPB(mcpBaseURL(r)+"/mcp", q.Token)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	auditLog(r, userID, "mcp.desktop_extension", "", map[string]interface{}{"with_token": q.Token != ""})
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="wrm-pro.mcpb"`)
	w.Header().Set("Cache-Control", "no-store")
	w.Write(data)
}
