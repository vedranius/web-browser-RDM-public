package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"

	"golang.org/x/crypto/ssh"
)

func terminalDirectoryCommand(directory string, windows bool) (string, error) {
	if directory == "" {
		return "", nil
	}
	if len(directory) > 2048 || strings.ContainsAny(directory, "\x00\r\n") {
		return "", fmt.Errorf("unsupported terminal directory")
	}
	if windows {
		if len(directory) > 3 && directory[0] == '/' && directory[2] == ':' {
			directory = directory[1:]
		}
		script := "Set-Location -LiteralPath '" + strings.ReplaceAll(directory, "'", "''") + "' -ErrorAction Stop"
		units := utf16.Encode([]rune(script))
		data := make([]byte, len(units)*2)
		for i, v := range units {
			binary.LittleEndian.PutUint16(data[i*2:], v)
		}
		return "powershell.exe -NoLogo -NoExit -EncodedCommand " + base64.StdEncoding.EncodeToString(data), nil
	}
	script := "cd -- " + shellQuote(directory) + ` && exec "${SHELL:-/bin/sh}" -l`
	return "sh -c " + shellQuote(script), nil
}
func terminalDirectoryStart(cl *ssh.Client, c Connection, directory string) (string, error) {
	if directory == "" {
		return "", nil
	}
	windows := c.Options["terminal_shell"] == "powershell"
	if c.Options["terminal_shell"] == "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		out, err := monitorSSHCommand(ctx, cl, "cmd.exe /c ver")
		windows = err == nil && strings.Contains(strings.ToLower(out), "windows")
	}
	return terminalDirectoryCommand(directory, windows)
}
