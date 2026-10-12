#!/usr/bin/env bash
set -e

echo "=== [1/5] Kreiranje AI Sanitizer modula za Ollama tool-calling ==="
cat << 'GOCODE' > remote-manager/ai_sanitizer.go
package main

import (
	"encoding/json"
	"strings"
)

// SanitizeAIHistory sprječava da model upadne u sirovi JSON loop nakon učitavanja povijesti iz baze
func SanitizeAIHistory(messages []map[string]interface{}) []map[string]interface{} {
	var clean []map[string]interface{}
	for _, m := range messages {
		content, ok := m["content"].(string)
		role, _ := m["role"].(string)
		if ok && role == "assistant" {
			trimmed := strings.TrimSpace(content)
			if strings.HasPrefix(trimmed, `{"name":`) && strings.Contains(trimmed, `"arguments":`) {
				var parsed struct {
					Name      string          `json:"name"`
					Arguments json.RawMessage `json:"arguments"`
				}
				if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil && parsed.Name != "" {
					m["content"] = ""
					m["tool_calls"] = []map[string]interface{}{
						{
							"id":   "call_" + parsed.Name,
							"type": "function",
							"function": map[string]interface{}{
								"name":      parsed.Name,
								"arguments": string(parsed.Arguments),
							},
						},
					}
				}
			}
		}
		clean = append(clean, m)
	}
	return clean
}
GOCODE

echo "=== [2/5] Kreiranje Smart Safe Policy (zaštita i auto-naredbe) ==="
cat << 'GOCODE' > remote-manager/ai_policy.go
package main

import (
	"regexp"
	"strings"
)

var safeAutoRegex = []*regexp.Regexp{
	regexp.MustCompile(`^df(\s+-[a-zA-Z0-9]+)*(\s+[/\w.-]+)*$`),
	regexp.MustCompile(`^du(\s+-[a-zA-Z0-9]+)*(\s+[/\w.-]+)*$`),
	regexp.MustCompile(`^free(\s+-[a-zA-Z0-9]+)*$`),
	regexp.MustCompile(`^uptime$`),
	regexp.MustCompile(`^top\s+-b\s+-n\s+1.*$`),
	regexp.MustCompile(`^ps\s+(aux|ef|-ef).*$`),
	regexp.MustCompile(`^ls(\s+-[a-zA-Z0-9]+)*(\s+[/\w.-]+)*$`),
	regexp.MustCompile(`^docker\s+(ps|stats\s+--no-stream)$`),
}

var blockedDangerousRegex = []*regexp.Regexp{
	regexp.MustCompile(`\brm(\s+-[a-zA-Z0-9]*r|\s+--recursive)`),
	regexp.MustCompile(`\bmkfs`),
	regexp.MustCompile(`\bdd\s+`),
	regexp.MustCompile(`>\s*/dev/`),
	regexp.MustCompile(`\b(reboot|shutdown|poweroff)\b`),
	regexp.MustCompile(`\bpasswd\b`),
}

func EvaluateAutoCommand(cmd string) (allowed bool, reason string) {
	cmd = strings.TrimSpace(cmd)
	for _, blk := range blockedDangerousRegex {
		if blk.MatchString(cmd) {
			return false, "Destruktivna naredba je blokirana zaštitom sustava."
		}
	}
	for _, safe := range safeAutoRegex {
		if safe.MatchString(cmd) {
			return true, ""
		}
	}
	return false, "Naredba traži ručnu potvrdu korisnika."
}
GOCODE

echo "=== [3/5] Podizanje verzije na v12.3.0 kroz cijeli codebase ==="
# Zamjena verzije 12.2.9 u 12.3.0 u svim relevantnim datotekama
grep -rl "12.2.9" remote-manager/ README.md 2>/dev/null | xargs sed -i 's/12\.2\.9/12.3.0/g' || true
grep -rl "v12.2.9" remote-manager/ README.md 2>/dev/null | xargs sed -i 's/v12\.2\.9/v12.3.0/g' || true

# Dodavanje v12.3.0 na vrh CHANGELOG.md
TODAY=$(date +%Y-%m-%d)
CHANGELOG_ENTRY="## [12.3.0] - ${TODAY}

### Added
- **Smart Safe Auto-Mode**: Automatsko odobravanje sigurnih dijagnostičkih naredbi (\`df\`, \`free\`, \`top\`, \`ps\`, \`docker ps\`, \`ls\`) u auto modu.
- **AI UI Thinking Status**: Prikaz stanja obrade i izvođenja naredbi na serveru u stvarnom vremenu.
- **Ollama Context Sanitizer**: Sprječavanje slanja tekstualnog JSON-a u model nakon učitavanja povijesti iz baze.

### Fixed
- Ispravljen loop odbijanja alata kod Qwen modela na lokalnoj Ollami.
- Session 0 poboljšanja stabilnosti servisa na Windows poslužiteljima.

"

if [ -f CHANGELOG.md ]; then
  echo -e "${CHANGELOG_ENTRY}\n$(cat CHANGELOG.md)" > CHANGELOG.md
fi

echo "=== [4/5] Git staging, commit i tag v12.3.0 ==="
git add -A
git commit -m "feat(ai): release v12.3.0 - fix Ollama tool-calling history loop, safe auto-mode and bump version"
git tag -a v12.3.0 -m "Release v12.3.0"

echo "=== [5/5] Pushanje na GitHub i pokretanje GitHub Actions ==="
git push origin HEAD:main
git push origin v12.3.0

rm -f update_and_release.sh
echo -e "\n\033[1;32m>>> GOTOVO! Promjene su pushane i tag v12.3.0 je postavljen.\033[0m"
echo -e "\033[1;33mPrati izradu releasa uživo na:\033[0m https://github.com/vedranius/web-browser-RDM-public/actions"
