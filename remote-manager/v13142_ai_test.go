package main

import "testing"

// TestAIDestructiveBypasses covers AI-01: wrapper options with a value, relative parent
// traversal, and here-strings must not hide a destructive command.
func TestAIDestructiveBypasses(t *testing.T) {
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
		"Format-Volume -DriveLetter C":                         "mkfs",
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
	}
	for _, cmd := range allowed {
		if b := aiDestructive(cmd); b != nil {
			t.Errorf("%q was blocked as %q (%s), want allowed", cmd, b.Rule, b.Why)
		}
	}
}
