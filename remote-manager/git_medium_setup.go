package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func mediumAnyRunning() bool {
	mediumMgr.mu.Lock()
	defer mediumMgr.mu.Unlock()
	for _, p := range mediumMgr.processes {
		select {
		case <-p.done:
		default:
			return true
		}
	}
	return false
}

// Uploads are retained before dependency setup. A failed/closed browser request
// cannot lose approved archives; the same button resumes pending configuration.
func apiMediumSetup(w http.ResponseWriter, r *http.Request, uid int) {
	if !requireMethod(w, r, "POST") {
		return
	}
	if r.URL.Query().Get("confirm") != "1" {
		jsonError(w, "Confirm that you trust the deploy tool/configuration", 400)
		return
	}
	if !mediumOperation.TryLock() {
		jsonError(w, "Another Medium operation is running", 409)
		return
	}
	handed := false
	defer func() {
		if !handed {
			mediumOperation.Unlock()
		}
	}()
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		jsonError(w, "Cannot read ZIP uploads (20 MiB each)", 400)
		return
	}
	defer r.MultipartForm.RemoveAll()
	if err := os.MkdirAll(mediumRoot(), 0700); err != nil {
		jsonError(w, "Cannot create private runtime directory", 500)
		return
	}
	staged := map[string]string{}
	digests := map[string]string{}
	defer func() {
		for _, p := range staged {
			os.RemoveAll(p)
		}
	}()
	for _, kind := range []string{"tool", "config", "ssh"} {
		f, _, err := r.FormFile(kind)
		if err == http.ErrMissingFile {
			continue
		}
		if err != nil {
			jsonError(w, "Invalid ZIP upload", 400)
			return
		}
		data, err := io.ReadAll(io.LimitReader(f, mediumArchiveLimit+1))
		f.Close()
		if err != nil || len(data) > mediumArchiveLimit {
			jsonError(w, "ZIP exceeds 20 MiB", 413)
			return
		}
		if kind == "tool" && !isAdminUser(uid) {
			jsonError(w, "Only administrators can install deploy-tool code", 403)
			return
		}
		p, err := os.MkdirTemp(mediumRoot(), "one-click-")
		if err != nil {
			jsonError(w, "Cannot stage ZIP", 500)
			return
		}
		staged[kind] = p
		if err := mediumExtract(data, p, kind); err != nil {
			jsonError(w, err.Error(), 400)
			return
		}
		sum := sha256.Sum256(data)
		digests[kind] = hex.EncodeToString(sum[:])
	}
	if mediumAnyRunning() && len(staged) > 0 {
		jsonError(w, "Stop active Medium instances before replacing tool/configuration", 409)
		return
	}
	if staged["tool"] == "" && !mediumState(uid).Imported {
		jsonError(w, "Choose the trusted deploy-tool ZIP", 400)
		return
	}
	if !isAdminUser(uid) && !mediumRuntimeReady() {
		jsonError(w, "An administrator must prepare the shared runtime first", 403)
		return
	}
	if !mediumState(uid).Config && staged["config"] == "" && !mediumFileExists(filepath.Join(mediumUserDir(uid), "pending-config", "servers.yaml")) {
		jsonError(w, "Choose the deploy-config ZIP", 400)
		return
	}
	for _, kind := range []string{"tool", "config", "ssh"} {
		if p := staged[kind]; p != "" {
			dest := filepath.Join(mediumUserDir(uid), "pending-"+kind)
			if kind == "tool" {
				dest = mediumToolDir()
			}
			if err := mediumReplaceDir(p, dest); err != nil {
				jsonError(w, "Cannot retain approved ZIP files", 500)
				return
			}
			delete(staged, kind)
			if kind == "tool" {
				os.WriteFile(filepath.Join(mediumRoot(), "tool.sha256"), []byte(digests[kind]), 0600)
				os.WriteFile(filepath.Join(mediumRoot(), "setup-pending"), []byte("pending"), 0600)
			}
		}
	}
	mediumSetStage(uid, "docker", true, nil)
	handed = true
	go func() {
		defer mediumOperation.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		err := mediumSetup(ctx, uid)
		if err != nil {
			mediumSetStage(uid, "error", false, err)
		} else {
			mediumSetStage(uid, "ready", false, nil)
		}
	}()
	auditLog(r, uid, "git.medium_setup_started", "", map[string]interface{}{"archives": len(digests)})
	w.WriteHeader(http.StatusAccepted)
	jsonOK(w, mediumState(uid))
}
func mediumSetup(ctx context.Context, uid int) error {
	if mediumBackend() == "docker" {
		mediumSetStage(uid, "docker", true, nil)
		if err := mediumEnsureDocker(ctx); err != nil {
			return err
		}
	}
	if !mediumRuntimeReady() {
		mediumSetStage(uid, "dependencies", true, nil)
		if err := mediumPrepareRuntime(ctx); err != nil {
			return err
		}
	}
	mediumSetStage(uid, "configuration", true, nil)
	for _, kind := range []string{"config", "ssh"} {
		pending := filepath.Join(mediumUserDir(uid), "pending-"+kind)
		if st, e := os.Stat(pending); e != nil || !st.IsDir() {
			continue
		}
		dest := filepath.Join(mediumUserDir(uid), "config")
		if kind == "config" {
			if err := mediumValidateConfig(ctx, pending); err != nil {
				return fmt.Errorf("Config validation failed; approved ZIP is retained for retry. Check environments.yaml/servers.yaml and project and references")
			}
		} else {
			dest = filepath.Join(mediumUserDir(uid), "home", ".ssh")
		}
		if err := mediumReplaceDir(pending, dest); err != nil {
			return err
		}
		if kind == "config" {
			prefs := mediumPreferences(uid)
			delete(prefs, "DEPLOY_CONFIG_REPO")
			if err := os.WriteFile(filepath.Join(mediumUserDir(uid), "preferences.enc"), []byte(encryptValue(string(jsonMarshal(prefs)))), 0600); err != nil {
				return err
			}
		}
		if err := os.WriteFile(filepath.Join(mediumUserDir(uid), kind+"-revision"), []byte(randomToken(16)), 0600); err != nil {
			return err
		}
	}
	mediumSetStage(uid, "starting", true, nil)
	return mediumStartForSetup(uid, true)
}

func mediumEnsureDocker(ctx context.Context) error {
	if mediumDockerCheck(ctx) == nil {
		return nil
	}
	if mediumDockerCLI() == "" {
		if runtime.GOOS != "windows" {
			return fmt.Errorf("Docker Desktop/Engine is not installed. Install it once on the WRM host, then resume; your ZIPs are retained")
		}
		if serviceMode == "windows" {
			return fmt.Errorf("Install Docker Desktop for the WRM service account once in an interactive Windows session, then resume; ZIPs are retained")
		}
		if err := mediumInstallDesktop(ctx); err != nil {
			return err
		}
	}
	cli := mediumDockerCLI()
	if cli == "" {
		return fmt.Errorf("Docker Desktop installation did not finish. Complete its Windows installer, license/WSL prompts or restart, then resume; ZIPs are retained")
	}
	startCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	err := exec.CommandContext(startCtx, cli, "desktop", "start", "--timeout", "50").Run()
	cancel()
	if err != nil {
		switch runtime.GOOS {
		case "windows":
			for _, base := range []string{filepath.Join(os.Getenv("ProgramFiles"), "Docker", "Docker"), filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "DockerDesktop")} {
				p := filepath.Join(base, "Docker Desktop.exe")
				if mediumFileExists(p) {
					cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+strings.ReplaceAll(p, "'", "''")+"'")
					_ = cmd.Run()
					break
				}
			}
		case "darwin":
			_ = exec.CommandContext(ctx, "open", "-a", "Docker").Run()
		case "linux":
			for _, args := range [][]string{{"--user", "start", "docker"}, {"start", "docker"}} {
				cctx, stop := context.WithTimeout(ctx, 15*time.Second)
				_ = exec.CommandContext(cctx, "systemctl", args...).Run()
				stop()
				if mediumDockerCheck(ctx) == nil {
					return nil
				}
			}
		}
	}
	// Do not switch a host with running Windows containers or change the global
	// Docker context; those can belong to other applications.
	versionCtx, stop := context.WithTimeout(ctx, 5*time.Second)
	b, _ := exec.CommandContext(versionCtx, cli, "version", "--format", "{{.Server.Os}}").Output()
	stop()
	if strings.TrimSpace(string(b)) == "windows" {
		return fmt.Errorf("Docker is using Windows containers. Switch to Linux containers in Docker Desktop, then resume; no existing containers were stopped")
	}
	deadline := time.NewTimer(3 * time.Minute)
	defer deadline.Stop()
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		if mediumDockerCheck(ctx) == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("Docker Linux engine did not become ready. Complete Docker Desktop license/WSL/virtualization prompts or restart Windows, then resume; ZIPs are retained")
		case <-tick.C:
		}
	}
}

func mediumInstallDesktop(ctx context.Context) error {
	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}
	url := "https://desktop.docker.com/win/main/" + arch + "/Docker%20Desktop%20Installer.exe"
	dir, err := os.MkdirTemp(mediumRoot(), "docker-install-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	p := filepath.Join(dir, "Docker Desktop Installer.exe")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 15 * time.Minute, CheckRedirect: func(r *http.Request, via []*http.Request) error {
		if len(via) > 5 || r.URL.Scheme != "https" || !(r.URL.Hostname() == "docker.com" || strings.HasSuffix(r.URL.Hostname(), ".docker.com")) {
			return fmt.Errorf("Unexpected Docker download redirect")
		}
		return nil
	}}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Cannot download official Docker Desktop installer; check outbound HTTPS, then resume")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Docker installer download returned HTTP %d", resp.StatusCode)
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, io.LimitReader(resp.Body, 1<<30+1))
	closeErr := f.Close()
	if err != nil || closeErr != nil || n > 1<<30 {
		return fmt.Errorf("Docker installer download failed or exceeds 1 GiB")
	}
	quoted := "'" + strings.ReplaceAll(p, "'", "''") + "'"
	script := "$ErrorActionPreference='Stop'; $p=" + quoted + "; $s=Get-AuthenticodeSignature -LiteralPath $p; if($s.Status -ne 'Valid' -or $s.SignerCertificate.Subject -notmatch '(^|,\\s*)O=Docker Inc\\.?(,|$)'){throw 'Docker installer signature rejected'}; $x=Start-Process -FilePath $p -ArgumentList 'install','--user','--backend=wsl-2' -Wait -PassThru; if($x.ExitCode -ne 0){exit $x.ExitCode}"
	if err := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script).Run(); err != nil {
		return fmt.Errorf("Docker Desktop installer did not complete. Finish its Windows prompts or restart if requested, then resume; ZIPs are retained")
	}
	return nil
}
