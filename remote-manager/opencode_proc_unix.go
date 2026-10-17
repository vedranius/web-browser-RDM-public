//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// ocSysProcAttr puts OpenCode in its own process group, so a terminal signal to WRM does not
// reach it and WRM decides when it stops.
func ocSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
