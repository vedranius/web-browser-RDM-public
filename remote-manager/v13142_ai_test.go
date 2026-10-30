package main

import (
	"encoding/base64"
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

// TestAIDestructiveBypasses covers AI-01: wrapper options with a value, relative parent
// traversal, and here-strings must not hide a destructive command.
func TestAIDestructiveBypasses(t *testing.T) {
	encoded := func(script string) string {
		units := utf16.Encode([]rune(script))
		data := make([]byte, len(units)*2)
		for i, u := range units {
			binary.LittleEndian.PutUint16(data[i*2:], u)
		}
		return base64.StdEncoding.EncodeToString(data)
	}
	blocked := map[string]string{
		"sudo -u root rm -rf /":                                "rm_root",
		"sudo --user=root rm -rf /":                            "rm_root",
		"doas -u root rm -rf /":                                "rm_root",
		"sudo -u root mkfs.ext4 /dev/sda":                      "mkfs",
		"sudo -n rm -rf /var":                                  "rm_root",
		"rm -rf ../..":                                         "rm_root",
		"rm -rf ../../etc":                                     "rm_root",
		"chmod -R 777 ../..":                                   "chmod_root",
		"bash <<< 'rm -rf /'":                                  "rm_root",
		"sh <<< \"dd if=/dev/zero of=/dev/sda\"":               "dd_disk",
		"bash <<< 'sudo -u root rm -rf /'":                     "rm_root",
		"cd / && rm -rf ../..":                                 "rm_root",
		"Remove-Item -Recurse C:\\temp":                        "windows_recursive_delete",
		"rm -Recurse C:\\temp":                                 "windows_recursive_delete",
		"powershell -Command Remove-Item -Recurse C:\\Windows": "windows_recursive_delete",
		"rd /s /q C:\\":                                        "windows_recursive_delete",
		"Get-ChildItem -Recurse C:\\temp | Remove-Item -Force": "windows_recursive_delete",
		"powershell -Command \"Get-ChildItem -Recurse C:\\temp | Remove-Item -Force\"":              "windows_recursive_delete",
		"pwsh -Command \"Get-ChildItem -Recurse C:\\temp | Remove-Item -Force\"":                    "windows_recursive_delete",
		"powershell -NoProfile -Command \"Get-ChildItem -Recurse C:\\temp | Remove-Item -Force\"":   "windows_recursive_delete",
		"pwsh -NoLogo -NoProfile -Command \"Get-ChildItem -Recurse C:\\temp | Remove-Item -Force\"": "windows_recursive_delete",
		"gci -r C:\\temp | ri -Force":                                                    "windows_recursive_delete",
		"for /r C:\\temp %F in (*) do del %F":                                            "windows_recursive_delete",
		"for /r \"C:\\Program Files\" %F in (*) do del %F":                               "windows_recursive_delete",
		"cmd /c for /r C:\\temp %%F in (*) do del %%F":                                   "windows_recursive_delete",
		"cmd /c \"for /r C:\\temp %%F in (*) do del %%F\"":                               "windows_recursive_delete",
		"powershell -EncodedCommand " + encoded("Remove-Item -Recurse C:\\temp"):         "windows_recursive_delete",
		"powershell -en " + encoded("Remove-Item -Recurse C:\\temp"):                     "windows_recursive_delete",
		"pwsh -enco " + encoded("Remove-Item -Recurse C:\\temp"):                         "windows_recursive_delete",
		"powershell -e " + encoded("Remove-Item -Recurse C:\\temp"):                      "windows_recursive_delete",
		"sudo -u root pwsh -enc " + encoded("Remove-Item -Recurse C:\\temp"):             "windows_recursive_delete",
		"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe -enc not-base64": "windows_recursive_delete",
		"powershell -Command rm -r C:\\Windows":                                          "windows_recursive_delete",
		"pwsh -EncodedCommand not-base64":                                                "windows_recursive_delete",
		"Format-Volume -DriveLetter C":                                                   "mkfs",
	}
	for cmd, rule := range blocked {
		b := aiDestructive(cmd)
		if b == nil {
			t.Errorf("%q was not blocked", cmd)
			continue
		}
		if b.Rule != rule {
			t.Errorf("%q blocked as %q, want %q", cmd, b.Rule, rule)
		}
	}

	allowed := []string{
		"rm -rf ../cache",
		"rm -rf ./build",
		"sudo -n systemctl status nginx",
		"grep -q ok <<< \"ready\"",
		"Remove-Item C:\\temp",
		"rm -r /tmp/build",
		"echo Format-Volume",
		"Get-ChildItem -Recurse C:\\temp",
		"for /r C:\\temp %F in (*) do echo %F",
		"Get-ChildItem -Recurse C:\\temp | Measure-Object",
		"powershell -Command \"Get-ChildItem -Recurse C:\\temp | Measure-Object\"",
		"echo Get-ChildItem -Recurse C:\\temp | Remove-Item",
		"echo 'pwsh -enc not-base64'",
		"powershell -Command \"Write-Output 'pwsh -enc not-base64'\"",
		"pwsh -co \"Write-Output 'pwsh -enc not-base64'\"",
		"powershell -EncodedCommand " + encoded("Write-Output 'hello'"),
	}
	for _, cmd := range allowed {
		if b := aiDestructive(cmd); b != nil {
			t.Errorf("%q was blocked as %q (%s), want allowed", cmd, b.Rule, b.Why)
		}
	}
}
