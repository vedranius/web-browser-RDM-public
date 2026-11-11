package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var addonStudioSafeName = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// ─── ADDON STUDIO ────────────────────────────────────
//
// Writing an addon entirely inside WRM: start from a template or from an
// installed addon, edit its files, lint as you type, install it (disabled, as
// every draft) and send it to the community repository. The studio works on
// the same drafts as the assistant and MCP clients (addons_drafts.go). Unlike
// addon_create, a studio save keeps a draft that does not validate yet, so a
// half-written addon.json is never lost; the lint panel says what is wrong and
// Install stays disabled until it is right.

// addonStudioWrite writes (or with del removes) one file of a draft without
// validating the whole addon. The draft folder is opened as an os.Root, so the
// path cannot leave it.
func addonStudioWrite(id, rel string, content []byte, del bool) error {
	if !addonIDPattern.MatchString(id) {
		return fmt.Errorf("invalid addon id")
	}
	rel, err := cleanDraftRel(rel)
	if err != nil {
		return err
	}
	if strings.HasPrefix(path.Base(rel), ".") {
		return fmt.Errorf("hidden files are not part of an addon")
	}
	if len(content) > addonMaxFile {
		return fmt.Errorf("%s is larger than 8 MB", rel)
	}
	dir := addonDraftDir(id)
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("draft %s does not exist", id)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	if del {
		if rel == "addon.json" {
			return fmt.Errorf("addon.json cannot be removed")
		}
		return root.Remove(filepath.FromSlash(rel))
	}
	if files, total, err := addonDraftFiles(dir); err == nil {
		if len(files) >= addonMaxFiles {
			return fmt.Errorf("an addon has at most %d files", addonMaxFiles)
		}
		if total+int64(len(content)) > addonMaxArchive {
			return fmt.Errorf("the draft would be larger than 20 MB")
		}
	}
	if d := path.Dir(rel); d != "." {
		if err := root.MkdirAll(filepath.FromSlash(d), 0o755); err != nil {
			return err
		}
	}
	return root.WriteFile(filepath.FromSlash(rel), content, 0o644)
}

// addonStudioCreate starts a draft from a template, or from an installed
// addon when from names one ("id" or "~id").
func addonStudioCreate(userID int, id, template, from, name string) error {
	if !addonIDPattern.MatchString(id) || addonReservedIDs[id] {
		return fmt.Errorf("the addon id must be 2-64 characters: lowercase letters, digits, dot, dash or underscore (and not a reserved word)")
	}
	dir := addonDraftDir(id)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("a draft %s exists already: open it instead", id)
	}
	if err := os.MkdirAll(addonsDraftsRoot(), 0o755); err != nil {
		return fmt.Errorf("the drafts folder cannot be created")
	}
	if from != "" {
		ref, ok := addonParseKey(from, userID)
		if !ok {
			return fmt.Errorf("addon not found")
		}
		src := addonRefDir(ref)
		if _, err := os.Stat(filepath.Join(src, "addon.json")); err != nil {
			return fmt.Errorf("addon not found")
		}
		if ref.ID != id {
			return fmt.Errorf("a draft of an installed addon keeps its id (%s)", ref.ID)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		if err := addonDraftCopyTree(src, dir); err != nil {
			os.RemoveAll(dir)
			return err
		}
		return nil
	}
	if template == "" {
		template = "blank"
	}
	files, err := addonTemplateFiles(template, id)
	if err != nil {
		return err
	}
	if name = strings.TrimSpace(name); name != "" && len(name) <= 80 {
		var m map[string]interface{}
		if json.Unmarshal(files["addon.json"], &m) == nil {
			m["name"] = name
			if ui, ok := m["ui"].(map[string]interface{}); ok {
				ui["title"] = name
			}
			if b, err := json.MarshalIndent(m, "", "  "); err == nil {
				files["addon.json"] = append(b, '\n')
			}
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for rel, b := range files {
		if err := addonStudioWrite(id, rel, b, false); err != nil {
			os.RemoveAll(dir)
			return err
		}
	}
	return nil
}

// addonStudioRoute serves the studio's parts of /api/addons/drafts:
//
//	GET  /api/addons/drafts/templates
//	POST /api/addons/drafts                    {id, template | from, name}
//	GET  /api/addons/drafts/{id}/file?path=    one text file
//	PUT  /api/addons/drafts/{id}/file          {path, content}
//	DELETE /api/addons/drafts/{id}/file?path=
//	GET  /api/addons/drafts/{id}/lint
//	GET  /api/addons/drafts/{id}/zip
//
// It answers true when it handled the request.
func addonStudioRoute(w http.ResponseWriter, r *http.Request, adminID int, parts []string) bool {
	if len(parts) == 1 && parts[0] == "templates" && r.Method == http.MethodGet {
		jsonOK(w, map[string]interface{}{"templates": addonTemplateList()})
		return true
	}
	if len(parts) == 0 && r.Method == http.MethodPost {
		var req struct {
			ID       string `json:"id"`
			Template string `json:"template"`
			From     string `json:"from"`
			Name     string `json:"name"`
		}
		if json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&req) != nil {
			jsonError(w, "Invalid request", 400)
			return true
		}
		id := strings.TrimSpace(req.ID)
		if err := addonStudioCreate(adminID, id, strings.TrimSpace(req.Template), strings.TrimSpace(req.From), req.Name); err != nil {
			jsonError(w, err.Error(), 400)
			return true
		}
		auditLog(r, adminID, "addon.draft_created", id, map[string]interface{}{"template": truncateStr(req.Template, 40), "from": truncateStr(req.From, 70)})
		jsonOK(w, addonDraftViewOf(id))
		return true
	}
	if len(parts) != 2 || !addonIDPattern.MatchString(parts[0]) {
		return false
	}
	id := parts[0]
	dir := addonDraftDir(id)
	switch parts[1] {
	case "file":
		switch r.Method {
		case http.MethodGet:
			rel, err := cleanDraftRel(r.URL.Query().Get("path"))
			if err != nil {
				jsonError(w, err.Error(), 400)
				return true
			}
			text, ok := addonDraftRead(dir, rel)
			if !ok {
				jsonError(w, rel+" is binary, missing or larger than 2 MB", 404)
				return true
			}
			jsonOK(w, map[string]interface{}{"path": rel, "content": text})
		case http.MethodPut:
			var req struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if json.NewDecoder(io.LimitReader(r.Body, addonMaxFile+64<<10)).Decode(&req) != nil {
				jsonError(w, "Invalid request", 400)
				return true
			}
			if err := addonStudioWrite(id, req.Path, []byte(req.Content), false); err != nil {
				jsonError(w, err.Error(), 400)
				return true
			}
			jsonOK(w, map[string]interface{}{"draft": addonDraftViewOf(id), "lint": addonLintPath(r.Context(), dir)})
		case http.MethodDelete:
			if err := addonStudioWrite(id, r.URL.Query().Get("path"), nil, true); err != nil {
				jsonError(w, err.Error(), 400)
				return true
			}
			jsonOK(w, map[string]interface{}{"draft": addonDraftViewOf(id)})
		default:
			jsonError(w, "Method not allowed", 405)
		}
		return true
	case "lint":
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return true
		}
		if _, err := os.Stat(dir); err != nil {
			jsonError(w, "Draft not found", 404)
			return true
		}
		jsonOK(w, addonLintPath(r.Context(), dir))
		return true
	case "zip":
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return true
		}
		body, _, err := addonZipDir(dir)
		if err != nil {
			jsonError(w, "Draft not found", 404)
			return true
		}
		w.Header().Set("Content-Type", "application/zip")
		ver := addonStudioSafeName.ReplaceAllString(addonDraftVersion(id), "-")
		w.Header().Set("Content-Disposition", `attachment; filename="`+id+"-"+ver+`.zip"`)
		w.Write(body)
		return true
	}
	return false
}
