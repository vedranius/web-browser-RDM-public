//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

// ocSysProcAttr starts OpenCode without a console window.
func ocSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
}
