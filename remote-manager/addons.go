package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// ─── ADDONS ──────────────────────────────────────────
//
// An addon is a folder that is dropped in and discovered at runtime, in the
// spirit of Omarchy: <DB_PATH>.addons/<id>/addon.json plus the files it
// declares. Installing a ZIP is the only way in, so nothing executes at
// install time. The manifest declares capabilities; this build accepts
// interface hooks (a sandboxed panel and a theme) and WebAssembly server code
// (server.wasm, run through wazero in addons_wasm.go).

const (
	addonMaxArchive = 20 << 20
	addonMaxFiles   = 500
	addonMaxFile    = 8 << 20
)

var (
	addonIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{1,63}$`)
	// addonCapabilities is the whole trusted surface of an addon. A capability
	// that is not listed here refuses the install instead of being ignored.
	addonCapabilities = map[string]string{
		"ui.panel":    "show a panel in the WRM interface",
		"ui.theme":    "apply a colour theme",
		"server.wasm": "run a WebAssembly hook in the sandbox",
	}
	addonExportPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,63}$`)
)

type addonUI struct {
	Panel string `json:"panel,omitempty"` // panel.html
	Theme string `json:"theme,omitempty"` // theme.css
	Title string `json:"title,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

type addonServer struct {
	Module  string   `json:"module"`  // plugin.wasm
	Exports []string `json:"exports"` // functions the host may call
}

type addonManifest struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Description  string       `json:"description,omitempty"`
	Author       string       `json:"author,omitempty"`
	License      string       `json:"license,omitempty"`
	Homepage     string       `json:"homepage,omitempty"`
	Capabilities []string     `json:"capabilities"`
	UI           *addonUI     `json:"ui,omitempty"`
	Server       *addonServer `json:"server,omitempty"`
}

type addonRecord struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Description  string       `json:"description,omitempty"`
	Author       string       `json:"author,omitempty"`
	License      string       `json:"license,omitempty"`
	Homepage     string       `json:"homepage,omitempty"`
	Capabilities []string     `json:"capabilities"`
	UI           *addonUI     `json:"ui,omitempty"`
	Enabled      bool         `json:"enabled"`
	InstalledAt  string       `json:"installed_at"`
	InstalledBy  string       `json:"installed_by,omitempty"`
	Source       string       `json:"source,omitempty"`
	Checksum     string       `json:"checksum"`
	Server       *addonServer `json:"server,omitempty"`
}

func initAddonSchema() {
	ensureTable("addons", `CREATE TABLE IF NOT EXISTS addons (
		id TEXT PRIMARY KEY,
		version TEXT NOT NULL DEFAULT '',
		enabled INTEGER NOT NULL DEFAULT 0,
		manifest TEXT NOT NULL DEFAULT '',
		source TEXT NOT NULL DEFAULT '',
		checksum TEXT NOT NULL DEFAULT '',
		installed_at TEXT NOT NULL DEFAULT '',
		installed_by TEXT NOT NULL DEFAULT ''
	)`)
}

func addonsRoot() string { return resolveDBPath() + ".addons" }

func addonDir(id string) string { return filepath.Join(addonsRoot(), id) }

func addonRecords(enabledOnly bool) []addonRecord {
	q := `SELECT id, enabled, manifest, source, checksum, installed_at, installed_by FROM addons`
	if enabledOnly {
		q += ` WHERE enabled=1`
	}
	rows, err := db.Query(q + ` ORDER BY id`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := []addonRecord{}
	for rows.Next() {
		var id, manifest, source, checksum, installedAt, installedBy string
		var enabled int
		if rows.Scan(&id, &enabled, &manifest, &source, &checksum, &installedAt, &installedBy) != nil {
			continue
		}
		var m addonManifest
		if json.Unmarshal([]byte(manifest), &m) != nil {
			continue
		}
		out = append(out, addonRecord{
			ID: id, Name: m.Name, Version: m.Version, Description: m.Description, Author: m.Author,
			License: m.License, Homepage: m.Homepage, Capabilities: m.Capabilities, UI: m.UI, Server: m.Server,
			Enabled: enabled == 1, InstalledAt: installedAt, InstalledBy: installedBy, Source: source, Checksum: checksum,
		})
	}
	return out
}

func addonGet(id string) (addonRecord, bool) {
	for _, a := range addonRecords(false) {
		if a.ID == id {
			return a, true
		}
	}
	return addonRecord{}, false
}

// validate checks a manifest against the capabilities this build implements.
func (m *addonManifest) validate() error {
	if !addonIDPattern.MatchString(m.ID) {
		return fmt.Errorf("the addon id must be 2-64 characters: lowercase letters, digits, dot, dash or underscore")
	}
	m.Name = strings.TrimSpace(m.Name)
	m.Version = strings.TrimSpace(m.Version)
	if m.Name == "" || len(m.Name) > 80 {
		return fmt.Errorf("a name of 1-80 characters is required")
	}
	if m.Version == "" || len(m.Version) > 32 {
		return fmt.Errorf("a version of 1-32 characters is required")
	}
	if len(m.Capabilities) == 0 {
		return fmt.Errorf("the addon declares no capabilities")
	}
	if len(m.Capabilities) > 8 {
		return fmt.Errorf("too many capabilities")
	}
	seen := map[string]bool{}
	for _, c := range m.Capabilities {
		if _, ok := addonCapabilities[c]; !ok {
			return fmt.Errorf("capability %q is not supported by this version", c)
		}
		if seen[c] {
			return fmt.Errorf("capability %q is listed twice", c)
		}
		seen[c] = true
	}
	if m.Server != nil {
		if !seen["server.wasm"] {
			return fmt.Errorf("a server block requires the server.wasm capability")
		}
		if !addonSafeRel(m.Server.Module) || !strings.HasSuffix(strings.ToLower(m.Server.Module), ".wasm") {
			return fmt.Errorf("server.module must be a relative .wasm file")
		}
		if len(m.Server.Exports) == 0 {
			return fmt.Errorf("server.exports must name at least one exported function")
		}
		if len(m.Server.Exports) > 8 {
			return fmt.Errorf("at most 8 exported functions are allowed")
		}
		xseen := map[string]bool{}
		for _, e := range m.Server.Exports {
			if !addonExportPattern.MatchString(e) {
				return fmt.Errorf("export %q must be letters, digits and underscore, starting with a letter or underscore", e)
			}
			if xseen[e] {
				return fmt.Errorf("export %q is listed twice", e)
			}
			xseen[e] = true
		}
	} else if seen["server.wasm"] {
		return fmt.Errorf("the server.wasm capability needs a server block")
	}
	if m.UI != nil {
		if m.UI.Panel != "" {
			if !seen["ui.panel"] {
				return fmt.Errorf("ui.panel requires the ui.panel capability")
			}
			if !addonSafeRel(m.UI.Panel) || !strings.HasSuffix(strings.ToLower(m.UI.Panel), ".html") {
				return fmt.Errorf("ui.panel must be a relative .html file")
			}
		}
		if m.UI.Theme != "" {
			if !seen["ui.theme"] {
				return fmt.Errorf("ui.theme requires the ui.theme capability")
			}
			if !addonSafeRel(m.UI.Theme) || !strings.HasSuffix(strings.ToLower(m.UI.Theme), ".css") {
				return fmt.Errorf("ui.theme must be a relative .css file")
			}
		}
		if m.UI.Panel == "" && m.UI.Theme == "" {
			return fmt.Errorf("the ui block declares neither a panel nor a theme")
		}
	}
	if m.UI == nil && m.Server == nil {
		return fmt.Errorf("the addon declares neither a ui panel/theme nor server code")
	}
	return nil
}

// addonSafeRel accepts a relative path inside an addon folder and nothing else.
func addonSafeRel(p string) bool {
	if p == "" || len(p) > 300 || strings.ContainsAny(p, "\\\x00\r\n") || strings.HasPrefix(p, "/") {
		return false
	}
	clean := path.Clean(p)
	return clean == p && clean != "." && !strings.HasPrefix(clean, "../")
}

// servedFiles are the files the asset endpoint may return. Nothing else leaves
// an addon folder, so a stray note or key file stays private.
func addonAssetType(name string) (string, bool) {
	switch strings.ToLower(path.Ext(name)) {
	case ".html", ".htm":
		return "text/html; charset=utf-8", true
	case ".css":
		return "text/css; charset=utf-8", true
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8", true
	case ".json", ".map":
		return "application/json; charset=utf-8", true
	case ".svg":
		return "image/svg+xml", true
	case ".png":
		return "image/png", true
	case ".jpg", ".jpeg":
		return "image/jpeg", true
	case ".gif":
		return "image/gif", true
	case ".webp":
		return "image/webp", true
	case ".ico":
		return "image/x-icon", true
	case ".woff2":
		return "font/woff2", true
	case ".txt", ".md":
		return "text/plain; charset=utf-8", true
	}
	return "", false
}

// ─── install ─────────────────────────────────────────

// addonReadArchive extracts an uploaded ZIP into a temporary folder and returns
// its manifest. Entry names are checked before anything is written: absolute
// paths, "..", backslashes, duplicates and oversized entries refuse the whole
// archive.
func addonReadArchive(body []byte) (*addonManifest, string, error) {
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, "", fmt.Errorf("not a ZIP archive")
	}
	if len(zr.File) == 0 || len(zr.File) > addonMaxFiles {
		return nil, "", fmt.Errorf("an addon must contain 1-%d files", addonMaxFiles)
	}
	names := map[string]bool{}
	for _, f := range zr.File {
		n := strings.TrimSuffix(f.Name, "/")
		if !addonSafeRel(n) {
			return nil, "", fmt.Errorf("unsafe path in the archive: %q", f.Name)
		}
		key := strings.ToLower(n)
		if names[key] {
			return nil, "", fmt.Errorf("duplicate file in the archive: %q", f.Name)
		}
		names[key] = true
		if f.FileInfo().IsDir() {
			continue
		}
		if f.UncompressedSize64 > addonMaxFile {
			return nil, "", fmt.Errorf("%s is larger than 8 MB", n)
		}
	}
	// The manifest sits at the root, or one folder down (a ZIP made from a
	// folder). No other layout is guessed.
	prefix := ""
	if !names["addon.json"] {
		top := ""
		for n := range names {
			i := strings.Index(n, "/")
			if i < 0 {
				return nil, "", fmt.Errorf("addon.json is missing")
			}
			if top == "" {
				top = n[:i]
			} else if top != n[:i] {
				return nil, "", fmt.Errorf("addon.json is missing")
			}
		}
		if top == "" || !names[strings.ToLower(top+"/addon.json")] {
			return nil, "", fmt.Errorf("addon.json is missing")
		}
		prefix = top + "/"
	}
	var raw []byte
	for _, f := range zr.File {
		n := strings.TrimSuffix(f.Name, "/")
		if n != prefix+"addon.json" {
			continue
		}
		if f.UncompressedSize64 > 256<<10 {
			return nil, "", fmt.Errorf("addon.json is larger than 256 KiB")
		}
		rc, err := f.Open()
		if err != nil {
			return nil, "", fmt.Errorf("addon.json cannot be read")
		}
		raw, err = io.ReadAll(io.LimitReader(rc, 256<<10))
		rc.Close()
		if err != nil {
			return nil, "", fmt.Errorf("addon.json cannot be read")
		}
	}
	if raw == nil {
		return nil, "", fmt.Errorf("addon.json is missing")
	}
	// Editors on Windows like to add a byte order mark; JSON does not allow it.
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	var m addonManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, "", fmt.Errorf("addon.json is not valid JSON: %v", err)
	}
	if err := m.validate(); err != nil {
		return nil, "", err
	}
	// Declared entries must be in the archive, or the addon would install and
	// then show nothing.
	declared := []string{}
	if m.UI != nil {
		declared = append(declared, m.UI.Panel, m.UI.Theme)
	}
	if m.Server != nil {
		declared = append(declared, m.Server.Module)
	}
	for _, want := range declared {
		if want != "" && !names[strings.ToLower(prefix+want)] {
			return nil, "", fmt.Errorf("%s is declared in addon.json but missing from the archive", want)
		}
	}
	return &m, prefix, nil
}

func addonWriteFiles(body []byte, prefix, dst string) error {
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		n := strings.TrimSuffix(f.Name, "/")
		if !strings.HasPrefix(n, prefix) {
			continue
		}
		rel := strings.TrimPrefix(n, prefix)
		if rel == "" || f.FileInfo().IsDir() {
			continue
		}
		target := filepath.Join(dst, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, io.LimitReader(rc, addonMaxFile+1))
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func addonInstall(w http.ResponseWriter, r *http.Request, adminID int) {
	body, err := io.ReadAll(io.LimitReader(r.Body, addonMaxArchive+1))
	if err != nil || len(body) == 0 {
		jsonError(w, "Send the addon as a ZIP file", 400)
		return
	}
	if len(body) > addonMaxArchive {
		jsonError(w, "The archive is larger than 20 MB", 413)
		return
	}
	m, prefix, err := addonReadArchive(body)
	if err != nil {
		jsonError(w, err.Error(), 400)
		return
	}
	sum := sha256.Sum256(body)
	checksum := "sha256:" + hex.EncodeToString(sum[:])
	if err := os.MkdirAll(addonsRoot(), 0o755); err != nil {
		jsonError(w, "The addon folder cannot be created", 500)
		return
	}
	tmp, err := os.MkdirTemp(addonsRoot(), ".install-")
	if err != nil {
		jsonError(w, "The addon folder cannot be written", 500)
		return
	}
	defer os.RemoveAll(tmp)
	if err := addonWriteFiles(body, prefix, tmp); err != nil {
		jsonError(w, "The archive cannot be written: "+err.Error(), 500)
		return
	}
	// Compiling is not executing, and it turns a broken module into an install
	// error instead of a surprise on the first call.
	if m.Server != nil {
		vctx, vcancel := context.WithTimeout(r.Context(), 10*time.Second)
		err := addonCheckWASM(vctx, m, tmp)
		vcancel()
		if err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
	}
	final := addonDir(m.ID)
	old := final + ".old"
	os.RemoveAll(old)
	previous, hadPrevious := addonGet(m.ID)
	if _, err := os.Stat(final); err == nil {
		if err := os.Rename(final, old); err != nil {
			jsonError(w, "The previous version cannot be replaced", 500)
			return
		}
	}
	if err := os.Rename(tmp, final); err != nil {
		if _, statErr := os.Stat(old); statErr == nil {
			os.Rename(old, final)
		}
		jsonError(w, "The addon folder cannot be written", 500)
		return
	}
	os.RemoveAll(old)
	manifest, _ := json.Marshal(m)
	source := strings.TrimSpace(r.URL.Query().Get("source"))
	if len(source) > 300 {
		source = source[:300]
	}
	enabled := 0
	if hadPrevious && previous.Enabled {
		enabled = 1
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := db.Exec(`INSERT INTO addons (id, version, enabled, manifest, source, checksum, installed_at, installed_by)
		VALUES (?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET version=excluded.version, manifest=excluded.manifest,
			source=excluded.source, checksum=excluded.checksum, installed_at=excluded.installed_at, installed_by=excluded.installed_by`,
		m.ID, m.Version, enabled, string(manifest), source, checksum, now, usernameOf(adminID)); err != nil {
		jsonError(w, "The addon cannot be recorded", 500)
		return
	}
	auditLog(r, adminID, "addon.installed", m.ID, map[string]interface{}{
		"version": m.Version, "capabilities": m.Capabilities, "checksum": checksum, "enabled": enabled == 1,
	})
	jsonOK(w, map[string]interface{}{"id": m.ID, "name": m.Name, "version": m.Version, "enabled": enabled == 1, "checksum": checksum})
}

func addonSetEnabled(w http.ResponseWriter, r *http.Request, adminID int, id string, on bool) {
	rec, ok := addonGet(id)
	if !ok {
		jsonError(w, "Addon not found", 404)
		return
	}
	if on {
		if _, err := os.Stat(filepath.Join(addonDir(id), "addon.json")); err != nil {
			jsonError(w, "The addon files are missing; install the archive again", 409)
			return
		}
	}
	v := 0
	if on {
		v = 1
	}
	if _, err := db.Exec(`UPDATE addons SET enabled=? WHERE id=?`, v, id); err != nil {
		jsonError(w, "The addon cannot be updated", 500)
		return
	}
	action := "addon.disabled"
	if on {
		action = "addon.enabled"
	}
	auditLog(r, adminID, action, id, map[string]interface{}{"version": rec.Version})
	jsonOK(w, map[string]bool{"ok": true, "enabled": on})
}

func addonRemove(w http.ResponseWriter, r *http.Request, adminID int, id string) {
	rec, ok := addonGet(id)
	if !ok {
		jsonError(w, "Addon not found", 404)
		return
	}
	if _, err := db.Exec(`DELETE FROM addons WHERE id=?`, id); err != nil {
		jsonError(w, "The addon cannot be removed", 500)
		return
	}
	os.RemoveAll(addonDir(id))
	auditLog(r, adminID, "addon.removed", id, map[string]interface{}{"version": rec.Version})
	jsonOK(w, map[string]bool{"ok": true})
}

// addonAssetHandler serves a file from an addon folder. Only an enabled addon
// is served, only known asset types, and only from inside its own folder.
func addonAssetHandler(w http.ResponseWriter, r *http.Request, id, rel string) {
	rec, ok := addonGet(id)
	if !ok || !rec.Enabled {
		jsonError(w, "Addon not found", 404)
		return
	}
	rel = path.Clean("/" + rel)[1:]
	if !addonSafeRel(rel) {
		jsonError(w, "Invalid file", 400)
		return
	}
	typ, ok := addonAssetType(rel)
	if !ok {
		jsonError(w, "Invalid file", 400)
		return
	}
	root, err := filepath.Abs(addonDir(id))
	if err != nil {
		jsonError(w, "Addon not found", 404)
		return
	}
	full, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil || !strings.HasPrefix(full, root+string(os.PathSeparator)) {
		jsonError(w, "Invalid file", 400)
		return
	}
	st, err := os.Stat(full)
	if err != nil || !st.Mode().IsRegular() {
		jsonError(w, "Addon not found", 404)
		return
	}
	f, err := os.Open(full)
	if err != nil {
		jsonError(w, "Addon not found", 404)
		return
	}
	defer f.Close()
	h := w.Header()
	h.Set("Content-Type", typ)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Cache-Control", "no-cache")
	// The global middleware sends X-Frame-Options: DENY. A panel is the one
	// document WRM frames itself, and only from the same origin.
	h.Set("X-Frame-Options", "SAMEORIGIN")
	if strings.HasPrefix(typ, "text/html") {
		// Belt and braces: the panel already runs in a sandboxed iframe.
		h.Set("Content-Security-Policy", "sandbox allow-scripts")
	}
	http.ServeContent(w, r, rel, st.ModTime(), f)
}

func apiAddonsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuth(w, r)
	if !ok {
		return
	}
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/addons"), "/")
	parts := strings.Split(rest, "/")
	admin := isAdminUser(userID)
	if rest == "" {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		list := addonRecords(!admin)
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		jsonOK(w, map[string]interface{}{"addons": list, "capabilities": addonCapabilities})
		return
	}
	if parts[0] == "install" {
		if r.Method != http.MethodPost {
			jsonError(w, "Method not allowed", 405)
			return
		}
		if !admin {
			jsonError(w, "Admin access required", 403)
			return
		}
		addonInstall(w, r, userID)
		return
	}
	id := parts[0]
	if len(parts) >= 3 && parts[1] == "asset" {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		addonAssetHandler(w, r, id, strings.Join(parts[2:], "/"))
		return
	}
	if len(parts) == 3 && parts[1] == "call" {
		if r.Method != http.MethodPost {
			jsonError(w, "Method not allowed", 405)
			return
		}
		if !admin {
			jsonError(w, "Admin access required", 403)
			return
		}
		addonCallHandler(w, r, userID, id, parts[2])
		return
	}
	if !admin {
		jsonError(w, "Admin access required", 403)
		return
	}
	if len(parts) == 1 && r.Method == http.MethodDelete {
		addonRemove(w, r, userID, id)
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		jsonError(w, "Method not allowed", 405)
		return
	}
	switch parts[1] {
	case "enable", "disable":
		addonSetEnabled(w, r, userID, id, parts[1] == "enable")
	default:
		jsonError(w, "Not found", 404)
	}
}
