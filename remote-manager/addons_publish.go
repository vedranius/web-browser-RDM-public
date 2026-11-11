package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ─── ADDONS: PUBLISHING TO THE COMMUNITY REPOSITORY ──
//
// An addon written in WRM (Addon Studio, the assistant or an MCP client) is
// sent to the community repository in one step: WRM lints it, refuses
// anything that looks like a secret, forks the repository under the author's
// GitHub account when needed, commits the sources to addons/<id>/ on a branch
// and opens a pull request. The repository's CI then runs WRM's own lint and
// its security scan; nothing becomes public until a maintainer approves and
// merges the pull request, after which the repository's CI builds the archive
// and the index WRM reads (addons_community.go).
//
// The author's GitHub token is stored encrypted per WRM user, never sent to
// the browser and never logged. A fine-grained token with "Contents" and
// "Pull requests" read/write on the author's fork (or a classic token with
// public_repo) is enough.

var (
	// tests point this at a fake GitHub
	addonGitHubAPI       = "https://api.github.com"
	addonPublishRepoExpr = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,38})/[A-Za-z0-9._-]{1,100}$`)
	addonPublishForkWait = 30 * time.Second
)

func initAddonPublishSchema() {
	ensureTable("addon_publish_auth", `CREATE TABLE IF NOT EXISTS addon_publish_auth (
		user_id INTEGER PRIMARY KEY,
		token TEXT NOT NULL DEFAULT '',
		login TEXT NOT NULL DEFAULT '',
		saved_at TEXT NOT NULL DEFAULT ''
	)`)
	ensureTable("addon_submissions", `CREATE TABLE IF NOT EXISTS addon_submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL DEFAULT 0,
		addon_id TEXT NOT NULL DEFAULT '',
		version TEXT NOT NULL DEFAULT '',
		repo TEXT NOT NULL DEFAULT '',
		branch TEXT NOT NULL DEFAULT '',
		pr_number INTEGER NOT NULL DEFAULT 0,
		pr_url TEXT NOT NULL DEFAULT '',
		state TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL DEFAULT '',
		updated_at TEXT NOT NULL DEFAULT ''
	)`)
}

func addonPublishRepo() string {
	r := strings.TrimSpace(getSetting("addon_publish_repo"))
	if !addonPublishRepoExpr.MatchString(r) {
		return "vedranius/wrm-pro-community-addons"
	}
	return r
}

func addonPublishAllowed(userID int) bool {
	v := getSetting("addon_publish")
	return v != "off" && policyAllows(userID, v)
}

// ── GitHub client ────────────────────────────────────

type ghError struct {
	Status int
	Msg    string
}

func (e *ghError) Error() string { return fmt.Sprintf("GitHub answered %d: %s", e.Status, e.Msg) }

func ghClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}

// ghCall sends one GitHub REST request. The token goes only to addonGitHubAPI.
func ghCall(ctx context.Context, token, method, p string, body, out interface{}) error {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(addonGitHubAPI, "/")+p, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "WRM-Pro/"+AppVersion)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := ghClient().Do(req)
	if err != nil {
		return fmt.Errorf("GitHub cannot be reached: %s", truncateStr(redactText(err.Error()), 200))
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, 8<<20))
	if res.StatusCode >= 300 {
		var e struct {
			Message string `json:"message"`
			Errors  []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		json.Unmarshal(data, &e)
		msg := e.Message
		for _, x := range e.Errors {
			if x.Message != "" {
				msg += "; " + x.Message
			}
		}
		if msg == "" {
			msg = http.StatusText(res.StatusCode)
		}
		return &ghError{Status: res.StatusCode, Msg: truncateStr(redactText(msg), 300)}
	}
	if out != nil && len(data) > 0 {
		return json.Unmarshal(data, out)
	}
	return nil
}

func ghStatus(err error) int {
	if e, ok := err.(*ghError); ok {
		return e.Status
	}
	return 0
}

func ghPath(parts ...string) string {
	out := ""
	for _, p := range parts {
		for _, seg := range strings.Split(p, "/") {
			out += "/" + url.PathEscape(seg)
		}
	}
	return out
}

// ── token ────────────────────────────────────────────

func addonPublishToken(userID int) (token, login string) {
	var enc string
	if db.QueryRow(`SELECT token, login FROM addon_publish_auth WHERE user_id=?`, userID).Scan(&enc, &login) != nil || enc == "" {
		return "", ""
	}
	return decryptValue(enc), login
}

func addonPublishSaveToken(ctx context.Context, userID int, token string) (string, error) {
	token = strings.TrimSpace(token)
	if token == "" || len(token) > 400 || strings.ContainsAny(token, " \r\n\t") {
		return "", fmt.Errorf("paste a GitHub token (a fine-grained token, or a classic token with public_repo)")
	}
	var u struct {
		Login string `json:"login"`
	}
	if err := ghCall(ctx, token, http.MethodGet, "/user", nil, &u); err != nil {
		return "", fmt.Errorf("the token does not work: %v", err)
	}
	if u.Login == "" {
		return "", fmt.Errorf("GitHub did not say whose token this is")
	}
	_, err := db.Exec(`INSERT INTO addon_publish_auth (user_id, token, login, saved_at) VALUES (?,?,?,?)
		ON CONFLICT(user_id) DO UPDATE SET token=excluded.token, login=excluded.login, saved_at=excluded.saved_at`,
		userID, encryptValue(token), u.Login, nowRFC())
	return u.Login, err
}

// ── what is sent ─────────────────────────────────────

// addonPublishSource finds the folder of a draft ("draft:<id>") or of an
// installed addon the user may manage (its key).
func addonPublishSource(userID int, admin bool, src string) (string, error) {
	if id, ok := strings.CutPrefix(src, "draft:"); ok {
		if !admin {
			return "", fmt.Errorf("drafts belong to administrators")
		}
		if !addonIDPattern.MatchString(id) {
			return "", fmt.Errorf("invalid addon id")
		}
		dir := addonDraftDir(id)
		if _, err := os.Stat(dir); err != nil {
			return "", fmt.Errorf("draft %s does not exist", id)
		}
		return dir, nil
	}
	ref, ok := addonParseKey(src, userID)
	if !ok {
		return "", fmt.Errorf("addon not found")
	}
	if ref.Owner == 0 && !admin {
		return "", fmt.Errorf("only an administrator publishes an organisation addon")
	}
	dir := addonRefDir(ref)
	if _, err := os.Stat(filepath.Join(dir, "addon.json")); err != nil {
		return "", fmt.Errorf("addon not found")
	}
	return dir, nil
}

type addonPublishFile struct {
	Path string
	Data []byte
}

type addonPublishCheck struct {
	Manifest *addonManifest
	Lint     addonLintReport
	Files    []addonPublishFile
	Problems []string
}

// addonPublishPrepare lints the folder and collects its files; Problems lists
// what blocks a submission (lint errors, secret-looking values, size).
func addonPublishPrepare(ctx context.Context, dir string) addonPublishCheck {
	out := addonPublishCheck{Lint: addonLintPath(ctx, dir)}
	for _, d := range out.Lint.Diagnostics {
		if d.Level == "error" {
			where := d.File
			if d.Line > 0 {
				where += ":" + strconv.Itoa(d.Line)
			}
			out.Problems = append(out.Problems, strings.TrimSpace(d.Code+" "+where+" "+d.Message))
		}
	}
	m, err := addonDraftManifest(dir)
	if err != nil {
		out.Problems = append(out.Problems, "addon.json: "+err.Error())
		return out
	}
	out.Manifest = m
	files, total, err := addonDraftFiles(dir)
	if err != nil {
		out.Problems = append(out.Problems, err.Error())
		return out
	}
	if total > addonMaxArchive {
		out.Problems = append(out.Problems, "the addon is larger than 20 MB")
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		out.Problems = append(out.Problems, err.Error())
		return out
	}
	defer root.Close()
	for _, f := range files {
		if strings.HasPrefix(f.Path, ".") || strings.Contains(f.Path, "/.") {
			continue // hidden files (.git, editor state) are not published
		}
		fh, err := root.Open(filepath.FromSlash(f.Path))
		if err != nil {
			out.Problems = append(out.Problems, f.Path+": "+err.Error())
			continue
		}
		data, err := io.ReadAll(io.LimitReader(fh, addonMaxFile+1))
		fh.Close()
		if err != nil {
			out.Problems = append(out.Problems, f.Path+": "+err.Error())
			continue
		}
		if addonDraftText(data) {
			if _, n := redactSecrets(string(data)); n > 0 {
				out.Problems = append(out.Problems, f.Path+": contains a value that looks like a secret (token, password or private key); declare it in http.secrets instead")
			}
		}
		out.Files = append(out.Files, addonPublishFile{Path: f.Path, Data: data})
	}
	sort.Slice(out.Files, func(i, j int) bool { return out.Files[i].Path < out.Files[j].Path })
	return out
}

func addonPublishBody(m *addonManifest, lint addonLintReport, note, login string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s %s (`%s`)\n\n", m.Name, m.Version, m.ID)
	if m.Description != "" {
		b.WriteString(m.Description + "\n\n")
	}
	if note = strings.TrimSpace(note); note != "" {
		b.WriteString("### Note from the author\n\n" + truncateStr(note, 2000) + "\n\n")
	}
	b.WriteString("### What it asks for\n\n")
	for _, s := range lint.Summary {
		b.WriteString("- " + s + "\n")
	}
	if len(lint.Summary) == 0 {
		b.WriteString("- nothing beyond a window\n")
	}
	fmt.Fprintf(&b, "\n### Checks in WRM before sending\n\n- WRM lint: %d errors, %d warnings\n- secret scan: clean\n", lint.Errors, lint.Warnings)
	if m.License != "" {
		fmt.Fprintf(&b, "- license: %s\n", m.License)
	}
	fmt.Fprintf(&b, "\nSubmitted from WRM Pro %s by @%s. The repository's checks run on this pull request; the addon is published in the community library only after a maintainer approves and merges it.\n", AppVersion, login)
	return b.String()
}

type addonSubmission struct {
	ID        int    `json:"id"`
	AddonID   string `json:"addon_id"`
	Version   string `json:"version"`
	Repo      string `json:"repo"`
	Branch    string `json:"branch"`
	PRNumber  int    `json:"pr_number"`
	PRURL     string `json:"pr_url"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// addonPublishSubmit sends the checked files to GitHub and opens (or updates)
// the pull request. It returns the recorded submission.
func addonPublishSubmit(ctx context.Context, userID int, chk addonPublishCheck, note string) (addonSubmission, error) {
	token, login := addonPublishToken(userID)
	if token == "" {
		return addonSubmission{}, fmt.Errorf("connect your GitHub account first (Settings → Addons → Publish)")
	}
	if len(chk.Problems) > 0 {
		return addonSubmission{}, fmt.Errorf("fix these first: %s", strings.Join(chk.Problems, "; "))
	}
	m := chk.Manifest
	up := addonPublishRepo()
	owner, _, _ := strings.Cut(up, "/")
	var repo struct {
		DefaultBranch string `json:"default_branch"`
		Permissions   struct {
			Push bool `json:"push"`
		} `json:"permissions"`
	}
	if err := ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(up), nil, &repo); err != nil {
		return addonSubmission{}, err
	}
	base := repo.DefaultBranch
	if base == "" {
		base = "main"
	}
	var ref struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(up)+"/git/ref/heads"+ghPath(base), nil, &ref); err != nil {
		return addonSubmission{}, err
	}
	var commit struct {
		Tree struct {
			SHA string `json:"sha"`
		} `json:"tree"`
	}
	if err := ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(up)+"/git/commits/"+url.PathEscape(ref.Object.SHA), nil, &commit); err != nil {
		return addonSubmission{}, err
	}
	// where the branch goes: the repository itself for its maintainers, a fork for everyone else
	target := up
	if !repo.Permissions.Push && !strings.EqualFold(login, owner) {
		var fork struct {
			FullName string `json:"full_name"`
		}
		if err := ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(up)+"/forks", map[string]interface{}{"default_branch_only": true}, &fork); err != nil {
			return addonSubmission{}, fmt.Errorf("cannot fork %s: %v", up, err)
		}
		if !addonPublishRepoExpr.MatchString(fork.FullName) {
			return addonSubmission{}, fmt.Errorf("GitHub did not name the fork")
		}
		target = fork.FullName
		deadline := time.Now().Add(addonPublishForkWait)
		for {
			if ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(target), nil, nil) == nil {
				break
			}
			if time.Now().After(deadline) {
				return addonSubmission{}, fmt.Errorf("the fork %s is not ready yet; send again in a minute", target)
			}
			select {
			case <-ctx.Done():
				return addonSubmission{}, ctx.Err()
			case <-time.After(2 * time.Second):
			}
		}
		// bring the fork up to date so the base commit is there (best effort)
		ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(target)+"/merge-upstream", map[string]string{"branch": base}, nil)
	}
	// files: new blobs, and deletions of files the addon no longer has
	prefix := "addons/" + m.ID + "/"
	type entry struct {
		Path string  `json:"path"`
		Mode string  `json:"mode"`
		Type string  `json:"type"`
		SHA  *string `json:"sha"`
	}
	tree := []entry{}
	have := map[string]bool{}
	for _, f := range chk.Files {
		var blob struct {
			SHA string `json:"sha"`
		}
		if err := ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(target)+"/git/blobs",
			map[string]string{"content": base64.StdEncoding.EncodeToString(f.Data), "encoding": "base64"}, &blob); err != nil {
			return addonSubmission{}, err
		}
		sha := blob.SHA
		tree = append(tree, entry{Path: prefix + f.Path, Mode: "100644", Type: "blob", SHA: &sha})
		have[prefix+f.Path] = true
	}
	var old struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}
	if ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(up)+"/git/trees/"+url.PathEscape(commit.Tree.SHA)+"?recursive=1", nil, &old) == nil {
		for _, t := range old.Tree {
			if t.Type == "blob" && strings.HasPrefix(t.Path, prefix) && !have[t.Path] {
				tree = append(tree, entry{Path: t.Path, Mode: "100644", Type: "blob", SHA: nil})
			}
		}
	}
	var newTree struct {
		SHA string `json:"sha"`
	}
	if err := ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(target)+"/git/trees", map[string]interface{}{"base_tree": commit.Tree.SHA, "tree": tree}, &newTree); err != nil {
		return addonSubmission{}, err
	}
	var newCommit struct {
		SHA string `json:"sha"`
	}
	msg := fmt.Sprintf("Addon %s %s (%s)\n\nSubmitted from WRM Pro %s.", m.Name, m.Version, m.ID, AppVersion)
	if err := ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(target)+"/git/commits",
		map[string]interface{}{"message": msg, "tree": newTree.SHA, "parents": []string{ref.Object.SHA}}, &newCommit); err != nil {
		return addonSubmission{}, err
	}
	branch := "addon/" + m.ID + "-" + regexp.MustCompile(`[^A-Za-z0-9._-]`).ReplaceAllString(m.Version, "-")
	err := ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(target)+"/git/refs", map[string]string{"ref": "refs/heads/" + branch, "sha": newCommit.SHA}, nil)
	if ghStatus(err) == 422 {
		err = ghCall(ctx, token, http.MethodPatch, "/repos"+ghPath(target)+"/git/refs/heads"+ghPath(branch), map[string]interface{}{"sha": newCommit.SHA, "force": true}, nil)
	}
	if err != nil {
		return addonSubmission{}, err
	}
	head := branch
	if target != up {
		head = strings.SplitN(target, "/", 2)[0] + ":" + branch
	}
	var pr struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
	}
	err = ghCall(ctx, token, http.MethodPost, "/repos"+ghPath(up)+"/pulls", map[string]interface{}{
		"title": fmt.Sprintf("Addon: %s %s (%s)", m.Name, m.Version, m.ID), "head": head, "base": base,
		"body": addonPublishBody(m, chk.Lint, note, login), "maintainer_can_modify": true}, &pr)
	if ghStatus(err) == 422 {
		// a pull request for this branch is open already: the push above updated it
		var open []struct {
			Number  int    `json:"number"`
			HTMLURL string `json:"html_url"`
		}
		if ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(up)+"/pulls?state=open&head="+url.QueryEscape(head), nil, &open) == nil && len(open) > 0 {
			pr.Number, pr.HTMLURL, err = open[0].Number, open[0].HTMLURL, nil
		}
	}
	if err != nil {
		return addonSubmission{}, err
	}
	now := nowRFC()
	res, _ := db.Exec(`INSERT INTO addon_submissions (user_id, addon_id, version, repo, branch, pr_number, pr_url, state, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		userID, m.ID, m.Version, up, branch, pr.Number, pr.HTMLURL, "waiting", now, now)
	id, _ := res.LastInsertId()
	return addonSubmission{ID: int(id), AddonID: m.ID, Version: m.Version, Repo: up, Branch: branch, PRNumber: pr.Number, PRURL: pr.HTMLURL, State: "waiting", CreatedAt: now, UpdatedAt: now}, nil
}

func addonSubmissionsOf(userID int) []addonSubmission {
	out := []addonSubmission{}
	rows, err := db.Query(`SELECT id, addon_id, version, repo, branch, pr_number, pr_url, state, created_at, updated_at FROM addon_submissions WHERE user_id=? ORDER BY id DESC LIMIT 50`, userID)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var s addonSubmission
		if rows.Scan(&s.ID, &s.AddonID, &s.Version, &s.Repo, &s.Branch, &s.PRNumber, &s.PRURL, &s.State, &s.CreatedAt, &s.UpdatedAt) == nil {
			out = append(out, s)
		}
	}
	return out
}

// addonSubmissionRefresh asks GitHub where a pull request stands:
// waiting (open), published (merged) or declined (closed without a merge).
func addonSubmissionRefresh(ctx context.Context, userID, id int) (addonSubmission, error) {
	var s addonSubmission
	if db.QueryRow(`SELECT id, addon_id, version, repo, branch, pr_number, pr_url, state, created_at, updated_at FROM addon_submissions WHERE id=? AND user_id=?`, id, userID).
		Scan(&s.ID, &s.AddonID, &s.Version, &s.Repo, &s.Branch, &s.PRNumber, &s.PRURL, &s.State, &s.CreatedAt, &s.UpdatedAt) != nil {
		return s, fmt.Errorf("submission not found")
	}
	token, _ := addonPublishToken(userID)
	var pr struct {
		State  string `json:"state"`
		Merged bool   `json:"merged"`
	}
	if err := ghCall(ctx, token, http.MethodGet, "/repos"+ghPath(s.Repo)+"/pulls/"+strconv.Itoa(s.PRNumber), nil, &pr); err != nil {
		return s, err
	}
	switch {
	case pr.Merged:
		s.State = "published"
	case pr.State == "closed":
		s.State = "declined"
	default:
		s.State = "waiting"
	}
	s.UpdatedAt = nowRFC()
	db.Exec(`UPDATE addon_submissions SET state=?, updated_at=? WHERE id=?`, s.State, s.UpdatedAt, s.ID)
	return s, nil
}

// ── API ──────────────────────────────────────────────

// /api/addons/publish                       GET status, submissions
// /api/addons/publish/github                PUT {token} / DELETE
// /api/addons/publish/check                 POST {source}: lint + secret scan, nothing sent
// /api/addons/publish/submit                POST {source, note}
// /api/addons/publish/submissions/{id}      POST: refresh from GitHub
func addonPublishHandler(w http.ResponseWriter, r *http.Request, userID int, admin bool, parts []string) {
	allowed := addonPublishAllowed(userID)
	if len(parts) == 0 {
		if r.Method != http.MethodGet {
			jsonError(w, "Method not allowed", 405)
			return
		}
		_, login := addonPublishToken(userID)
		jsonOK(w, map[string]interface{}{"repo": addonPublishRepo(), "allowed": allowed, "github_login": login,
			"repo_url": "https://github.com/" + addonPublishRepo(), "submissions": addonSubmissionsOf(userID)})
		return
	}
	if !allowed {
		jsonError(w, "Publishing addons is turned off for you (addon_publish)", 403)
		return
	}
	switch {
	case parts[0] == "github" && len(parts) == 1:
		switch r.Method {
		case http.MethodPut:
			var req struct {
				Token string `json:"token"`
			}
			if json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&req) != nil {
				jsonError(w, "Invalid request", 400)
				return
			}
			login, err := addonPublishSaveToken(r.Context(), userID, req.Token)
			if err != nil {
				jsonError(w, err.Error(), 400)
				return
			}
			auditLog(r, userID, "addon.publish_account", "", map[string]interface{}{"github_login": login})
			jsonOK(w, map[string]interface{}{"github_login": login})
		case http.MethodDelete:
			db.Exec(`DELETE FROM addon_publish_auth WHERE user_id=?`, userID)
			auditLog(r, userID, "addon.publish_account_removed", "", nil)
			jsonOK(w, map[string]bool{"ok": true})
		default:
			jsonError(w, "Method not allowed", 405)
		}
	case (parts[0] == "check" || parts[0] == "submit") && len(parts) == 1:
		if r.Method != http.MethodPost {
			jsonError(w, "Method not allowed", 405)
			return
		}
		var req struct {
			Source string `json:"source"`
			Note   string `json:"note"`
		}
		if json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(&req) != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		dir, err := addonPublishSource(userID, admin, strings.TrimSpace(req.Source))
		if err != nil {
			jsonError(w, err.Error(), 404)
			return
		}
		chk := addonPublishPrepare(r.Context(), dir)
		if parts[0] == "check" {
			files := []string{}
			for _, f := range chk.Files {
				files = append(files, f.Path)
			}
			jsonOK(w, map[string]interface{}{"ok": len(chk.Problems) == 0, "problems": chk.Problems, "lint": chk.Lint, "files": files, "repo": addonPublishRepo()})
			return
		}
		sub, err := addonPublishSubmit(r.Context(), userID, chk, req.Note)
		if err != nil {
			auditLog(r, userID, "addon.publish_failed", strings.TrimSpace(req.Source), map[string]interface{}{"error": truncateStr(err.Error(), 300)})
			jsonError(w, err.Error(), 400)
			return
		}
		auditLog(r, userID, "addon.published", sub.AddonID, map[string]interface{}{"version": sub.Version, "repo": sub.Repo, "pr": sub.PRNumber})
		jsonOK(w, sub)
	case parts[0] == "submissions" && len(parts) == 2:
		if r.Method != http.MethodPost {
			jsonError(w, "Method not allowed", 405)
			return
		}
		id, _ := strconv.Atoi(parts[1])
		s, err := addonSubmissionRefresh(r.Context(), userID, id)
		if err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		jsonOK(w, s)
	default:
		jsonError(w, "Not found", 404)
	}
}
