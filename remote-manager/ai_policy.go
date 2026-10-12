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
