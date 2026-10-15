package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed medium_wrapper.py
var mediumWrapper []byte

const mediumPrefix = "/medium"
const mediumArchiveLimit = 20 << 20

// Only an administrator installs executable code. Each operator's config, working
// tree, SSH profile and process live in a separate private runtime directory.
type mediumStatus struct {
	Allowed      bool            `json:"allowed"`
	Admin        bool            `json:"admin"`
	Tool         bool            `json:"tool"`
	Imported     bool            `json:"tool_imported"`
	Backend      string          `json:"backend"`
	Config       bool            `json:"config"`
	Running      bool            `json:"running"`
	Busy         bool            `json:"busy"`
	Stage        string          `json:"stage"`
	Error        string          `json:"error,omitempty"`
	Digest       string          `json:"tool_digest,omitempty"`
	Dependencies map[string]bool `json:"dependencies"`
}
type mediumProcess struct {
	cmd          *exec.Cmd
	cancel       context.CancelFunc
	base, secret string
	done         chan struct{}
}
type mediumManager struct {
	mu        sync.Mutex
	processes map[string]*mediumProcess
	statuses  map[string]mediumStatus
}

var mediumOperation sync.Mutex

// mediumStarts serializes starts of one user's instance without holding mediumMgr.mu
// (a start can wait minutes for Docker and readiness; status and proxy calls must not).
var mediumStarts = struct {
	sync.Mutex
	m map[string]*sync.Mutex
}{m: map[string]*sync.Mutex{}}

func mediumStartLock(dir string) *sync.Mutex {
	mediumStarts.Lock()
	defer mediumStarts.Unlock()
	if mediumStarts.m[dir] == nil {
		mediumStarts.m[dir] = &sync.Mutex{}
	}
	return mediumStarts.m[dir]
}

var mediumMgr = mediumManager{processes: map[string]*mediumProcess{}, statuses: map[string]mediumStatus{}}

func mediumRoot() string {
	p, _ := filepath.Abs(resolveDBPath())
	return p + ".medium"
}
func mediumUserDir(uid int) string { return filepath.Join(mediumRoot(), "users", strconv.Itoa(uid)) }
func mediumToolDir() string        { return filepath.Join(mediumRoot(), "tool") }
func mediumPython() string {
	if p := os.Getenv("WRM_MEDIUM_PYTHON"); p != "" {
		return p
	}
	for _, p := range []string{"python3", "python"} {
		if v, e := exec.LookPath(p); e == nil {
			return v
		}
	}
	return ""
}
func mediumVenvPython() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(mediumRoot(), "venv", "Scripts", "python.exe")
	}
	return filepath.Join(mediumRoot(), "venv", "bin", "python")
}
func mediumFileExists(p string) bool { s, e := os.Stat(p); return e == nil && !s.IsDir() }
func mediumAllowed(uid int) bool     { return gitActionAllowed(uid, "git_medium") }
func mediumState(uid int) mediumStatus {
	mediumMgr.mu.Lock()
	defer mediumMgr.mu.Unlock()
	s := mediumMgr.statuses[mediumUserDir(uid)]
	s.Allowed, s.Admin = mediumAllowed(uid), isAdminUser(uid)
	s.Backend = mediumBackend()
	s.Imported = mediumFileExists(filepath.Join(mediumToolDir(), "app", "main.py"))
	s.Tool = s.Imported && mediumRuntimeReady()
	s.Config = mediumFileExists(filepath.Join(mediumConfigDir(uid), "servers.yaml")) || mediumPreferences(uid)["DEPLOY_CONFIG_REPO"] != ""
	if p := mediumMgr.processes[mediumUserDir(uid)]; p != nil {
		select {
		case <-p.done:
		default:
			s.Running = true
		}
	}
	s.Dependencies = map[string]bool{"python": mediumPython() != "", "git": false, "ssh": false, "rsync": false}
	for _, n := range []string{"git", "ssh", "rsync"} {
		_, e := exec.LookPath(n)
		s.Dependencies[n] = e == nil
	}
	if s.Backend == "docker" {
		s.Dependencies = map[string]bool{"docker": mediumDockerCLI() != ""}
	}
	if b, e := os.ReadFile(filepath.Join(mediumRoot(), "tool.sha256")); e == nil {
		s.Digest = strings.TrimSpace(string(b))
	}
	return s
}
func mediumSetStage(uid int, stage string, busy bool, err error) {
	mediumMgr.mu.Lock()
	defer mediumMgr.mu.Unlock()
	s := mediumMgr.statuses[mediumUserDir(uid)]
	s.Stage, s.Busy = stage, busy
	s.Error = ""
	if err != nil {
		s.Error = truncateStr(redactText(err.Error()), 2000)
	}
	mediumMgr.statuses[mediumUserDir(uid)] = s
}
func mediumCommand(ctx context.Context, dir, binary string, args ...string) error {
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PIP_DISABLE_PIP_VERSION_CHECK=1")
	// Never return package output or arbitrary source code to the browser.
	if e := cmd.Run(); e != nil {
		return fmt.Errorf("runtime command failed (%s): %w", filepath.Base(binary), e)
	}
	return nil
}

// Extract a ZIP's selected files, with one optional top-level wrapper directory.
// Links, traversal, duplicates and decompression bombs are rejected before writes.
func mediumExtract(data []byte, dest, kind string) error {
	z, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if e != nil {
		return fmt.Errorf("invalid ZIP")
	}
	if len(z.File) > 5000 {
		return fmt.Errorf("ZIP has too many entries")
	}
	type entry struct {
		f   *zip.File
		rel string
	}
	entries := []entry{}
	seen := map[string]bool{}
	var total uint64
	for _, f := range z.File {
		n := f.Name
		if len(n) > 1024 || strings.Contains(n, "\\") || strings.HasPrefix(n, "/") || path.Clean(n) != strings.TrimSuffix(n, "/") || strings.Contains(n, ":") {
			return fmt.Errorf("unsafe ZIP path")
		}
		if f.Mode()&os.ModeSymlink != 0 || (!f.FileInfo().IsDir() && !f.Mode().IsRegular()) {
			return fmt.Errorf("ZIP links/devices are not allowed")
		}
		if f.FileInfo().IsDir() {
			continue
		}
		rel := n
		first := strings.SplitN(rel, "/", 2)
		if len(first) == 2 && first[0] != "app" && first[0] != "projects" && first[0] != "scripts" {
			rel = first[1]
		}
		allowed := false
		switch kind {
		case "tool":
			allowed = strings.HasPrefix(rel, "app/") || strings.HasPrefix(rel, "scripts/") || rel == "requirements.lock" || rel == "pyproject.toml"
		case "config":
			allowed = rel == "servers.yaml" || rel == "environments.yaml" || (strings.HasPrefix(rel, "projects/") && (strings.HasSuffix(rel, ".yaml") || strings.HasSuffix(rel, ".yml"))) || rel == ".gitattributes"
		case "ssh":
			allowed = rel == "config" || rel == "known_hosts" || regexp.MustCompile(`^[A-Za-z0-9_.-]+$`).MatchString(rel)
		}
		if !allowed {
			continue
		}
		if rel == "." || rel == ".." {
			continue
		}
		// case-insensitive: Windows and macOS file systems would overwrite one with the other
		if seen[strings.ToLower(rel)] {
			return fmt.Errorf("duplicate ZIP path")
		}
		seen[strings.ToLower(rel)] = true
		total += f.UncompressedSize64
		if total > mediumArchiveLimit || f.UncompressedSize64 > mediumArchiveLimit {
			return fmt.Errorf("expanded ZIP exceeds 20 MiB")
		}
		entries = append(entries, entry{f, rel})
	}
	required := []string{}
	if kind == "tool" {
		required = []string{"app/main.py", "app/config.py", "app/__init__.py", "requirements.lock"}
	}
	if kind == "config" {
		required = []string{"servers.yaml", "environments.yaml"}
	}
	for _, n := range required {
		if !seen[n] {
			return fmt.Errorf("ZIP is missing %s", n)
		}
	}
	if len(entries) == 0 {
		return fmt.Errorf("ZIP contains no supported files")
	}
	for _, v := range entries {
		dst := filepath.Join(dest, filepath.FromSlash(v.rel))
		if e = os.MkdirAll(filepath.Dir(dst), 0700); e != nil {
			return e
		}
		r, e := v.f.Open()
		if e != nil {
			return e
		}
		b, e := io.ReadAll(io.LimitReader(r, mediumArchiveLimit+1))
		r.Close()
		if e != nil || len(b) > mediumArchiveLimit {
			return fmt.Errorf("cannot read ZIP entry")
		}
		mode := os.FileMode(0600)
		if strings.HasPrefix(v.rel, "scripts/") {
			mode = 0700
		}
		if e = os.WriteFile(dst, b, mode); e != nil {
			return e
		}
	}
	return nil
}
func mediumStop(uid int) {
	mediumMgr.mu.Lock()
	p := mediumMgr.processes[mediumUserDir(uid)]
	delete(mediumMgr.processes, mediumUserDir(uid))
	mediumMgr.mu.Unlock()
	if p != nil {
		p.cancel()
		select {
		case <-p.done:
		case <-time.After(3 * time.Second):
		}
	}
}
func mediumStopAll() {
	mediumMgr.mu.Lock()
	ps := []*mediumProcess{}
	for k, p := range mediumMgr.processes {
		ps = append(ps, p)
		delete(mediumMgr.processes, k)
	}
	mediumMgr.mu.Unlock()
	for _, p := range ps {
		p.cancel()
	}
}
func mediumStart(uid int) error { return mediumStartForSetup(uid, false) }
func mediumStartForSetup(uid int, setup bool) error {
	if !mediumAllowed(uid) {
		return fmt.Errorf("Medium access is not allowed")
	}
	s := mediumState(uid)
	if s.Busy && !setup {
		return fmt.Errorf("runtime setup is still running")
	}
	if !s.Tool || !s.Config {
		return fmt.Errorf("import the deploy tool and deploy config first")
	}
	if s.Backend == "docker" {
		if e := mediumDockerCheck(context.Background()); e != nil {
			return e
		}
	} else {
		for _, n := range []string{"git", "ssh", "rsync"} {
			if !s.Dependencies[n] {
				return fmt.Errorf("Install %s on the WRM host or use the managed Docker runtime (WRM_MEDIUM_RUNTIME=docker)", n)
			}
		}
	}
	dir := mediumUserDir(uid)
	startLock := mediumStartLock(dir)
	startLock.Lock()
	defer startLock.Unlock()
	mediumMgr.mu.Lock()
	if p := mediumMgr.processes[dir]; p != nil {
		select {
		case <-p.done:
			delete(mediumMgr.processes, dir)
		default:
			mediumMgr.mu.Unlock()
			return nil
		}
	}
	mediumMgr.mu.Unlock()
	for _, d := range []string{dir, filepath.Join(dir, "data"), filepath.Join(dir, "home", ".ssh")} {
		if e := os.MkdirAll(d, 0700); e != nil {
			return e
		}
	}
	// earlier versions kept replaced profiles (and their keys) next to the current ones
	os.RemoveAll(filepath.Join(dir, "home", ".ssh.previous"))
	os.RemoveAll(filepath.Join(dir, "config.previous"))
	sshConfig := filepath.Join(dir, "home", ".ssh", "config")
	if !mediumFileExists(sshConfig) {
		if e := os.WriteFile(sshConfig, []byte("Host *\n  BatchMode yes\n  StrictHostKeyChecking yes\n"), 0600); e != nil {
			return e
		}
	}
	wrapper := filepath.Join(dir, "medium_wrapper.py")
	if e := os.WriteFile(wrapper, mediumWrapper, 0600); e != nil {
		return e
	}
	if s.Backend == "docker" {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if e := mediumDockerSync(ctx, uid); e != nil {
			return fmt.Errorf("Cannot synchronize private Linux runtime data; check Docker volume access and disk space")
		}
	}
	l, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		return e
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	ctx, cancel := context.WithCancel(context.Background())
	p := &mediumProcess{cancel: cancel, base: fmt.Sprintf("http://127.0.0.1:%d", port), secret: randomToken(32), done: make(chan struct{})}
	cmd := exec.CommandContext(ctx, mediumVenvPython(), "-m", "uvicorn", "medium_wrapper:app", "--host", "127.0.0.1", "--port", strconv.Itoa(port), "--no-proxy-headers", "--no-access-log")
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + filepath.Join(dir, "home"), "USERPROFILE=" + filepath.Join(dir, "home"), "LANG=C.UTF-8", "PYTHONUNBUFFERED=1", "WRM_MEDIUM_TOOL=" + mediumToolDir(), "WRM_MEDIUM_KEY=" + p.secret, "DEPLOY_DATA_DIR=" + filepath.Join(dir, "data"), "DEPLOY_CONFIG_DIR=" + mediumConfigDir(uid), "DEPLOY_BIND=127.0.0.1", "DEPLOY_PORT=" + strconv.Itoa(port), "DEPLOY_USER=" + usernameOf(uid), "DEPLOY_ALLOWED_HOSTS=127.0.0.1,localhost", "DEPLOY_SSH_OPTIONS=" + mediumSSHCommand(uid), "GIT_SSH_COMMAND=ssh " + mediumSSHCommand(uid)}
	for _, n := range []string{"SYSTEMROOT", "WINDIR", "TMP", "TEMP", "LD_LIBRARY_PATH", "SSL_CERT_FILE", "SSL_CERT_DIR"} {
		if v := os.Getenv(n); v != "" {
			cmd.Env = append(cmd.Env, n+"="+v)
		}
	}
	for k, v := range mediumPreferences(uid) {
		if v == "" {
			continue // an empty preference means the original default (and keeps WRM's SSH options)
		}
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	if sock := os.Getenv("WRM_MEDIUM_SSH_AUTH_SOCK"); sock != "" && serverKeysAllowed(uid) {
		cmd.Env = append(cmd.Env, "SSH_AUTH_SOCK="+sock)
	}
	if s.Backend == "docker" {
		env := []string{}
		// The pinned image's unprivileged wrm account also serves as the native
		// build user. Original builds drop privileges and cannot read root's SSH profile.
		env = append(env, "WRM_MEDIUM_ENV_DEPLOY_BUILD_USER=wrm")
		for _, value := range cmd.Env {
			key, val, _ := strings.Cut(value, "=")
			if key == "PATH" || key == "SYSTEMROOT" || key == "WINDIR" || key == "TMP" || key == "TEMP" || key == "LD_LIBRARY_PATH" || key == "SSL_CERT_FILE" || key == "SSL_CERT_DIR" || key == "SSH_AUTH_SOCK" {
				continue
			}
			switch key {
			case "DEPLOY_PORT":
				val = "8080"
			case "DEPLOY_BIND":
				val = "0.0.0.0"
			case "DEPLOY_ALLOWED_HOSTS":
				val = "127.0.0.1,localhost"
			case "DEPLOY_BUILD_USER":
				if val == "" {
					val = "wrm"
				}
			default:
				val = mediumContainerPath(val)
				if key == "DEPLOY_SSH_OPTIONS" || key == "GIT_SSH_COMMAND" {
					val = strings.ReplaceAll(val, mediumUserDir(uid), mediumContainerPath(mediumUserDir(uid)))
				}
			}
			env = append(env, "WRM_MEDIUM_ENV_"+key+"="+val)
		}
		keys := []string{}
		for _, value := range env {
			key, _, _ := strings.Cut(value, "=")
			keys = append(keys, key)
		}
		name := "wrm-medium-" + strconv.Itoa(uid) + "-" + randomToken(8)
		cmd = exec.CommandContext(ctx, mediumDockerCLI(), mediumDockerStartArgs(uid, port, name, keys)...)
		cmd.Env = append(os.Environ(), env...)
		var stop sync.Once
		p.cancel = func() { stop.Do(func() { mediumDockerStop(name); cancel() }) }
	}
	p.cmd = cmd
	if e = cmd.Start(); e != nil {
		cancel()
		return e
	}
	mediumMgr.mu.Lock()
	mediumMgr.processes[dir] = p
	mediumMgr.mu.Unlock()
	go func() { _ = cmd.Wait(); close(p.done) }()
	forget := func() {
		p.cancel()
		mediumMgr.mu.Lock()
		if mediumMgr.processes[dir] == p {
			delete(mediumMgr.processes, dir)
		}
		mediumMgr.mu.Unlock()
	}
	client := &http.Client{Timeout: time.Second}
	for i := 0; i < 300; i++ {
		select {
		case <-p.done:
			forget()
			return fmt.Errorf("deploy tool failed to start; check its Python dependencies and config")
		default:
		}
		req, _ := http.NewRequest("GET", p.base+"/__wrm_ready", nil)
		req.Header.Set("X-WRM-Medium-Key", p.secret)
		if r, e := client.Do(req); e == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	forget()
	return fmt.Errorf("deploy tool startup timed out")
}

func apiMediumHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !mediumAllowed(uid) {
		jsonError(w, "Medium access is disabled by policy", 403)
		return
	}
	action := strings.TrimPrefix(r.URL.Path, "/api/git/medium")
	if action == "" && r.Method == "GET" {
		jsonOK(w, mediumState(uid))
		return
	}
	if action == "/settings" && r.Method == "GET" {
		jsonOK(w, mediumPublicPreferences(uid))
		return
	}
	if !requireMethod(w, r, "POST") {
		return
	}
	if action == "/setup" {
		apiMediumSetup(w, r, uid)
		return
	}
	if action == "/stop" {
		// Stopping one's own instance never waits for another user's setup or start.
		// force=1 stops it even when its state cannot be read (a hung or crashed tool).
		if e := mediumCanStop(uid); e != nil && r.URL.Query().Get("force") != "1" {
			jsonError(w, e.Error(), 409)
			return
		}
		mediumStop(uid)
		auditLog(r, uid, "git.medium_stopped", "", map[string]interface{}{"forced": r.URL.Query().Get("force") == "1"})
		jsonOK(w, mediumState(uid))
		return
	}
	if !mediumOperation.TryLock() {
		jsonError(w, "Another Medium setup/start operation is running", 409)
		return
	}
	backgroundOperation := false
	defer func() {
		if !backgroundOperation {
			mediumOperation.Unlock()
		}
	}()
	if action == "/prepare" {
		if !isAdminUser(uid) || !mediumState(uid).Imported {
			jsonError(w, "An administrator must import the trusted tool ZIP first", 403)
			return
		}
		mediumMgr.mu.Lock()
		active := false
		for _, process := range mediumMgr.processes {
			select {
			case <-process.done:
			default:
				active = true
			}
		}
		mediumMgr.mu.Unlock()
		if active {
			jsonError(w, "Stop all Medium instances before preparing the shared runtime", 409)
			return
		}
		backgroundOperation = true
		mediumSetStage(uid, "dependencies", true, nil)
		go func() {
			defer mediumOperation.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			e := mediumPrepareRuntime(ctx)
			mediumSetStage(uid, "ready", false, e)
		}()
		auditLog(r, uid, "git.medium_runtime_prepared", "", nil)
		w.WriteHeader(202)
		jsonOK(w, mediumState(uid))
		return
	}
	if action == "/settings" {
		var values map[string]string
		if json.NewDecoder(r.Body).Decode(&values) != nil {
			jsonError(w, "Bad JSON", 400)
			return
		}
		if mediumState(uid).Running {
			jsonError(w, "Stop Medium before changing runtime settings", 409)
			return
		}
		saved := mediumPreferences(uid)
		for k, v := range values {
			if !mediumPreferenceKeys[k] || len(v) > 8192 || strings.ContainsAny(v, "\x00\r\n") {
				jsonError(w, "Unsupported runtime setting", 400)
				return
			}
			if v == "" {
				delete(saved, k) // empty: the original default
			} else if v != "[REDACTED]" {
				saved[k] = v
			}
		}
		if e := validateMediumPreferences(saved); e != nil {
			jsonError(w, e.Error(), 400)
			return
		}
		if e := os.MkdirAll(mediumUserDir(uid), 0700); e != nil {
			jsonError(w, "Cannot save runtime settings", 500)
			return
		}
		if e := os.WriteFile(filepath.Join(mediumUserDir(uid), "preferences.enc"), []byte(encryptValue(string(jsonMarshal(saved)))), 0600); e != nil {
			jsonError(w, "Cannot save runtime settings", 500)
			return
		}
		auditLog(r, uid, "git.medium_settings_saved", "", nil)
		jsonOK(w, map[string]bool{"ok": true})
		return
	}
	if action == "/start" {
		if e := mediumStart(uid); e != nil {
			jsonError(w, e.Error(), 400)
		} else {
			mediumSetStage(uid, "", false, nil)
			auditLog(r, uid, "git.medium_started", "", nil)
			jsonOK(w, mediumState(uid))
		}
		return
	}
	if action != "/tool" && action != "/config" && action != "/ssh" {
		jsonError(w, "Not found", 404)
		return
	}
	if action == "/tool" && !isAdminUser(uid) {
		jsonError(w, "Only administrators can install deploy-tool code", 403)
		return
	}
	// Trusted executable source and configuration may contain build/deploy commands.
	if r.URL.Query().Get("confirm") != "1" {
		jsonError(w, "Confirm that you trust this deploy tool/configuration", 400)
		return
	}
	s := mediumState(uid)
	if s.Busy {
		jsonError(w, "Runtime setup is already running", 409)
		return
	}
	if e := r.ParseMultipartForm(mediumArchiveLimit); e != nil {
		jsonError(w, "Cannot read ZIP upload (20 MiB maximum)", 400)
		return
	}
	defer r.MultipartForm.RemoveAll()
	f, _, e := r.FormFile("file")
	if e != nil {
		jsonError(w, "Choose a ZIP file", 400)
		return
	}
	defer f.Close()
	data, e := io.ReadAll(io.LimitReader(f, mediumArchiveLimit+1))
	if e != nil || len(data) > mediumArchiveLimit {
		jsonError(w, "ZIP exceeds 20 MiB", 413)
		return
	}
	if e = os.MkdirAll(mediumRoot(), 0700); e != nil {
		jsonError(w, "Cannot create runtime directory", 500)
		return
	}
	stage, e := os.MkdirTemp(mediumRoot(), "import-")
	if e != nil {
		jsonError(w, "Cannot prepare import", 500)
		return
	}
	kind := strings.TrimPrefix(action, "/")
	if e = mediumExtract(data, stage, kind); e != nil {
		os.RemoveAll(stage)
		jsonError(w, e.Error(), 400)
		return
	}
	if kind == "tool" {
		mediumMgr.mu.Lock()
		active := false
		for _, p := range mediumMgr.processes {
			select {
			case <-p.done:
			default:
				active = true
			}
		}
		mediumMgr.mu.Unlock()
		if active {
			os.RemoveAll(stage)
			jsonError(w, "Stop all Medium instances before replacing the shared tool runtime", 409)
			return
		}
		// Retain the approved source even if runtime setup is unavailable. A user
		// can start Docker and retry preparation without uploading the ZIP again.
		if e = mediumReplaceDir(stage, mediumToolDir()); e != nil {
			jsonError(w, "Cannot save imported tool", 500)
			return
		}
		if e = os.WriteFile(filepath.Join(mediumRoot(), "setup-pending"), []byte("pending"), 0600); e != nil {
			jsonError(w, "Cannot prepare runtime state", 500)
			return
		}
		backgroundOperation = true
		mediumSetStage(uid, "dependencies", true, nil)
		digest := sha256.Sum256(data)
		go func() {
			defer mediumOperation.Unlock()
			defer os.RemoveAll(stage)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			e := os.WriteFile(filepath.Join(mediumRoot(), "tool.sha256"), []byte(hex.EncodeToString(digest[:])), 0600)
			if e == nil {
				e = mediumPrepareRuntime(ctx)
			}
			mediumSetStage(uid, "ready", false, e)
		}()
		auditLog(r, uid, "git.medium_tool_imported", "", map[string]interface{}{"sha256": hex.EncodeToString(digest[:])})
		w.WriteHeader(202)
		jsonOK(w, mediumState(uid))
		return
	}
	defer os.RemoveAll(stage)
	if s.Running {
		jsonError(w, "Stop Medium before replacing its configuration or SSH profile", 409)
		return
	}
	mediumStop(uid)
	dest := filepath.Join(mediumUserDir(uid), "config")
	if kind == "config" {
		if !s.Tool {
			jsonError(w, "Prepare the deploy tool first", 400)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		if e = mediumValidateConfig(ctx, stage); e != nil {
			jsonError(w, "Deploy config validation failed; inspect YAML and project/server references", 400)
			return
		}
	} else {
		dest = filepath.Join(mediumUserDir(uid), "home", ".ssh")
	}
	if e = mediumReplaceDir(stage, dest); e != nil {
		jsonError(w, "Cannot save runtime files", 500)
		return
	}
	if e = os.WriteFile(filepath.Join(mediumUserDir(uid), kind+"-revision"), []byte(randomToken(16)), 0600); e != nil {
		jsonError(w, "Cannot save runtime import revision", 500)
		return
	}
	// This import replaces a pending one from an earlier one-click setup, and an imported
	// config replaces a config repository (as in one-click setup).
	os.RemoveAll(filepath.Join(mediumUserDir(uid), "pending-"+kind))
	if kind == "config" {
		if prefs := mediumPreferences(uid); prefs["DEPLOY_CONFIG_REPO"] != "" {
			delete(prefs, "DEPLOY_CONFIG_REPO")
			if e = os.WriteFile(filepath.Join(mediumUserDir(uid), "preferences.enc"), []byte(encryptValue(string(jsonMarshal(prefs)))), 0600); e != nil {
				jsonError(w, "Cannot save runtime settings", 500)
				return
			}
		}
	}
	mediumSetStage(uid, "", false, nil)
	auditLog(r, uid, "git.medium_"+kind+"_imported", "", nil)
	jsonOK(w, mediumState(uid))
}
func mediumReplaceDir(src, dst string) error {
	if e := os.MkdirAll(filepath.Dir(dst), 0700); e != nil {
		return e
	}
	backup := dst + ".previous"
	if e := os.RemoveAll(backup); e != nil {
		return e
	}
	existed := false
	if _, e := os.Stat(dst); e == nil {
		if e = os.Rename(dst, backup); e != nil {
			return e
		}
		existed = true
	}
	if e := os.Rename(src, dst); e != nil {
		if existed {
			_ = os.Rename(backup, dst)
		}
		return e
	}
	// A replaced SSH profile holds private keys: do not keep it on disk.
	return os.RemoveAll(backup)
}

var mediumURLAttr = regexp.MustCompile(`(?i)((?:href|src|action|hx-[a-z-]+|sse-connect)\s*=\s*["'])(/[^"']*)`)

func mediumRewriteHTML(b []byte) []byte {
	return mediumURLAttr.ReplaceAllFunc(b, func(m []byte) []byte {
		parts := mediumURLAttr.FindSubmatch(m)
		v := string(parts[2])
		if strings.HasPrefix(v, "//") || strings.HasPrefix(v, mediumPrefix+"/") {
			return m
		}
		return append(append([]byte{}, parts[1]...), []byte(mediumPrefix+v)...)
	})
}
func mediumProxyHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := requireAuth(w, r)
	if !ok {
		return
	}
	if !mediumAllowed(uid) {
		jsonError(w, "Medium access is disabled by policy", 403)
		return
	}
	site := r.Header.Get("Sec-Fetch-Site")
	// Original native forms use no-referrer, which can produce Origin:null.
	// Sec-Fetch-Site is a browser-controlled header; allow this only for same-origin.
	same := sameOriginRequest(r) || (r.Header.Get("Origin") == "null" && site == "same-origin")
	if !same || site == "cross-site" || site == "same-site" {
		jsonError(w, "Cross-origin Medium access is blocked", 403)
		return
	}
	route := strings.TrimPrefix(r.URL.Path, mediumPrefix)
	if !safeMethod(r.Method) {
		r.Body = http.MaxBytesReader(w, r.Body, mediumArchiveLimit)
		permission := "git_update"
		if strings.HasPrefix(route, "/povijest/rollback") {
			permission = "git_rollback"
		}
		if strings.Contains(route, "/deploy") && !gitActionAllowed(uid, "git_install") {
			jsonError(w, "Installation permission is required", 403)
			return
		}
		if !gitActionAllowed(uid, permission) {
			jsonError(w, "Deployment permission is required", 403)
			return
		}
		auditLog(r, uid, "git.medium_action", route, nil)
	}
	if !gitChecksAllowed(uid) {
		jsonError(w, "Git check permission is required", 403)
		return
	}
	mediumMgr.mu.Lock()
	p := mediumMgr.processes[mediumUserDir(uid)]
	mediumMgr.mu.Unlock()
	if p == nil {
		jsonError(w, "Start Medium from the Git workspace first", 409)
		return
	}
	select {
	case <-p.done:
		jsonError(w, "Medium process has stopped; start it again", 409)
		return
	default:
	}
	target, _ := url.Parse(p.base)
	proxy := httputil.NewSingleHostReverseProxy(target)
	original := proxy.Director
	proxy.Director = func(req *http.Request) {
		original(req)
		req.URL.Path = route
		req.URL.RawPath = ""
		req.Host = target.Host
		req.Header.Set("Origin", p.base)
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("X-WRM-Medium-Key", p.secret)
		req.Header.Del("Cookie")
		req.Header.Del("Authorization")
		req.Header.Del("Accept-Encoding")
	}
	proxy.Transport = mediumTransport
	proxy.ModifyResponse = func(resp *http.Response) error {
		for _, h := range []string{"Location", "HX-Redirect", "HX-Push-Url", "HX-Replace-Url"} {
			v := resp.Header.Get(h)
			v = strings.TrimPrefix(v, p.base)
			if strings.HasPrefix(v, "/") && !strings.HasPrefix(v, "//") {
				resp.Header.Set(h, mediumPrefix+v)
			}
		}
		resp.Header.Del("Set-Cookie")
		resp.Header.Set("X-Frame-Options", "SAMEORIGIN")
		resp.Header.Set("Cache-Control", "no-store")
		resp.Header.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; frame-ancestors 'self'; img-src 'self' data:; base-uri 'none'; form-action 'self'")
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			b, e := io.ReadAll(io.LimitReader(resp.Body, 32<<20+1))
			resp.Body.Close()
			if e != nil {
				return e
			}
			if len(b) > 32<<20 {
				return fmt.Errorf("Medium HTML response exceeds 32 MiB")
			}
			b = mediumRewriteHTML(b)
			resp.Body = io.NopCloser(bytes.NewReader(b))
			resp.ContentLength = int64(len(b))
			resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
			resp.Header.Del("Content-Encoding")
		}
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		jsonError(w, "Medium runtime did not answer; inspect its status and try again", 502)
	}
	w.Header().Del("X-Frame-Options")
	w.Header().Del("Content-Security-Policy")
	proxy.ServeHTTP(w, r)
}

var mediumPreferenceKeys = map[string]bool{
	"DEPLOY_CONFIG_REPO": true, "DEPLOY_CONFIG_PULL": true, "DEPLOY_CONFIG_CHECK_MINUTES": true,
	"DEPLOY_USER_EMAIL": true, "DEPLOY_HOSTNAME": true, "DEPLOY_SSH_OPTIONS": true, "DEPLOY_TIMEOUT": true,
	"DEPLOY_BUILD_USER": true, "GIT_SSH_COMMAND": true,
}

func mediumPreferences(uid int) map[string]string {
	values := map[string]string{}
	if b, e := os.ReadFile(filepath.Join(mediumUserDir(uid), "preferences.enc")); e == nil {
		json.Unmarshal([]byte(decryptValue(string(b))), &values)
	}
	for k := range values {
		if !mediumPreferenceKeys[k] {
			delete(values, k)
		}
	}
	return values
}
func mediumPublicPreferences(uid int) map[string]string {
	values := mediumPreferences(uid)
	for k, v := range values {
		if ((k == "GIT_SSH_COMMAND" || k == "DEPLOY_SSH_OPTIONS") && v != "") || redactText(v) != v {
			values[k] = "[REDACTED]"
			continue
		}
		if k == "DEPLOY_CONFIG_REPO" {
			if u, e := url.Parse(v); e == nil && u.User != nil {
				values[k] = "[REDACTED]"
			}
		}
	}
	return values
}

func mediumConfigDir(uid int) string {
	if mediumPreferences(uid)["DEPLOY_CONFIG_REPO"] != "" {
		return filepath.Join(mediumUserDir(uid), "git-config")
	}
	return filepath.Join(mediumUserDir(uid), "config")
}
func mediumSSHCommand(uid int) string {
	if mediumBackend() == "docker" {
		dir := "/wrm-medium/users/" + strconv.Itoa(uid) + "/home/.ssh"
		return "-o BatchMode=yes -o ConnectTimeout=15 -F " + shellQuote(dir+"/config") + " -o UserKnownHostsFile=" + shellQuote(dir+"/known_hosts")
	}
	dir := filepath.Join(mediumUserDir(uid), "home", ".ssh")
	return "-o BatchMode=yes -o ConnectTimeout=15 -F " + shellQuote(filepath.Join(dir, "config")) + " -o UserKnownHostsFile=" + shellQuote(filepath.Join(dir, "known_hosts"))
}

var mediumTransport = &http.Transport{Proxy: nil, DisableCompression: true, MaxIdleConns: 32, MaxIdleConnsPerHost: 8, IdleConnTimeout: 90 * time.Second, ResponseHeaderTimeout: 30 * time.Minute}

func mediumCanStop(uid int) error {
	mediumMgr.mu.Lock()
	p := mediumMgr.processes[mediumUserDir(uid)]
	mediumMgr.mu.Unlock()
	if p == nil {
		return nil
	}
	select {
	case <-p.done:
		return nil
	default:
	}
	req, _ := http.NewRequest("GET", p.base+"/__wrm_state", nil)
	req.Header.Set("X-WRM-Medium-Key", p.secret)
	client := &http.Client{Timeout: 3 * time.Second}
	r, e := client.Do(req)
	if e != nil {
		return fmt.Errorf("Medium status is unavailable; wait for outstanding operations to finish")
	}
	defer r.Body.Close()
	var status struct {
		Busy bool `json:"busy"`
	}
	if r.StatusCode != 200 || json.NewDecoder(io.LimitReader(r.Body, 4096)).Decode(&status) != nil {
		return fmt.Errorf("Cannot verify whether Medium has active operations")
	}
	if status.Busy {
		return fmt.Errorf("A Medium comparison/deployment is running; wait for it to finish before stopping")
	}
	return nil
}

func validateMediumPreferences(values map[string]string) error {
	for _, v := range []struct {
		k      string
		lo, hi int
	}{{"DEPLOY_TIMEOUT", 1, 86400}, {"DEPLOY_CONFIG_CHECK_MINUTES", 0, 1440}} {
		if values[v.k] != "" {
			n, e := strconv.Atoi(values[v.k])
			if e != nil || n < v.lo || n > v.hi {
				return fmt.Errorf("%s must be between %d and %d", v.k, v.lo, v.hi)
			}
		}
	}
	if v := strings.ToLower(values["DEPLOY_CONFIG_PULL"]); v != "" {
		valid := map[string]bool{"1": true, "0": true, "true": true, "false": true, "yes": true, "no": true, "da": true, "ne": true, "on": true, "off": true}
		if !valid[v] {
			return fmt.Errorf("DEPLOY_CONFIG_PULL must be a boolean")
		}
	}
	if strings.HasPrefix(values["DEPLOY_CONFIG_REPO"], "-") {
		return fmt.Errorf("Invalid config repository")
	}
	return nil
}
