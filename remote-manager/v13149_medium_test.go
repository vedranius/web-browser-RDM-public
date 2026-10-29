package main

import (
	"os"
	"slices"
	"testing"
)

// TestMediumHostnameAndGitMode covers v13.14.9: DEPLOY_HOSTNAME defaults to the WRM
// computer (not the container ID) and the status says whether deploy-config is a Git clone.
func TestMediumHostnameAndGitMode(t *testing.T) {
	srv := newTestServer(t)
	u := newTestUser(t, srv, "md-host-git", true)
	os.RemoveAll(mediumUserDir(u.userID))
	t.Cleanup(func() { os.RemoveAll(mediumUserDir(u.userID)) })
	host := mediumHostName()
	if host == "" {
		t.Skip("no host name")
	}
	mdSetPreferences(t, u.userID, map[string]string{})
	if env := mediumPreferenceEnv(u.userID); !slices.Contains(env, "DEPLOY_HOSTNAME="+host) {
		t.Fatalf("default DEPLOY_HOSTNAME missing: %v", env)
	}
	if mediumState(u.userID).ConfigGit {
		t.Fatal("config_git without DEPLOY_CONFIG_REPO")
	}

	mdSetPreferences(t, u.userID, map[string]string{"DEPLOY_HOSTNAME": "demo-laptop", "DEPLOY_CONFIG_REPO": "git@gitlab.example.com:team/deploy-config.git", "DEPLOY_TIMEOUT": ""})
	env := mediumPreferenceEnv(u.userID)
	want := []string{"DEPLOY_CONFIG_REPO=git@gitlab.example.com:team/deploy-config.git", "DEPLOY_HOSTNAME=demo-laptop"}
	if !slices.Equal(env, want) {
		t.Fatalf("env = %v, want %v", env, want)
	}
	var st mediumStatus
	if code := u.jsonDo("GET", "/api/git/medium", nil, &st); code != 200 {
		t.Fatalf("status %d", code)
	}
	if !st.ConfigGit || !st.Config || st.HostName != host {
		t.Fatalf("status config_git=%v config=%v host=%q", st.ConfigGit, st.Config, st.HostName)
	}
}
