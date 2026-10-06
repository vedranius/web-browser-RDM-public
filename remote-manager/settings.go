package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// ─── APPLICATION SETTINGS (admin policies) ───────────
//
// Settings live in the app_settings table and are edited in Admin panel → Policies.
// Every setting can also be forced with an environment variable WRM_<KEY IN UPPER CASE>
// (e.g. WRM_REGISTRATION=open). A setting forced by the environment is shown as locked
// in the admin panel and cannot be changed there — useful for configuration management.
// Some settings also accept a shorter alias (e.g. AUDIT_ENABLED).

type settingSpec struct {
	Key     string   `json:"key"`
	Default string   `json:"default"`
	Kind    string   `json:"kind"` // bool | int | enum | string | json
	Enum    []string `json:"enum,omitempty"`
	Min     int      `json:"min,omitempty"`
	Max     int      `json:"max,omitempty"`
	Restart bool     `json:"-"`
	Alias   string   `json:"-"` // additional environment variable name
}

var settingSpecs = []settingSpec{
	// Accounts & sign-in
	{Key: "registration", Default: "closed", Kind: "enum", Enum: []string{"closed", "open"}},
	{Key: "require_2fa", Default: "off", Kind: "enum", Enum: []string{"off", "admins", "all"}},
	{Key: "password_min_length", Default: "8", Kind: "int", Min: 6, Max: 128},
	{Key: "session_idle_hours", Default: "168", Kind: "int", Min: 1, Max: 8760},
	{Key: "session_max_days", Default: "30", Kind: "int", Min: 1, Max: 365},
	{Key: "login_max_failures", Default: "10", Kind: "int", Min: 3, Max: 100},
	// Connections
	{Key: "allow_server_keys", Default: "0", Kind: "bool"},
	{Key: "host_key_policy", Default: "tofu", Kind: "enum", Enum: []string{"tofu", "strict", "off"}},
	{Key: "max_upload_mb", Default: "0", Kind: "int", Min: 0, Max: 1048576},
	{Key: "allow_secret_export", Default: "1", Kind: "bool"},
	// Sharing & collaboration
	{Key: "allow_link_shares", Default: "1", Kind: "bool"},
	{Key: "chat_file_max_mb", Default: "5", Kind: "int", Min: 0, Max: 25},
	{Key: "voice_enabled", Default: "1", Kind: "bool"},
	{Key: "voice_max_participants", Default: "12", Kind: "int", Min: 2, Max: 50},
	{Key: "ice_servers", Default: `[{"urls":["stun:stun.l.google.com:19302","stun:stun.cloudflare.com:3478"]}]`, Kind: "json"},
	{Key: "turn_enabled", Default: "1", Kind: "bool", Restart: true},
	{Key: "turn_port", Default: "3478", Kind: "int", Min: 1, Max: 65535, Restart: true},
	{Key: "turn_public_ip", Default: "", Kind: "string", Restart: true},
	{Key: "turn_host", Default: "", Kind: "string"},
	{Key: "turn_relay_ports", Default: "49152-65535", Kind: "string", Restart: true},
	{Key: "turn_allow_private", Default: "0", Kind: "bool", Restart: true},
	// SSH tunnels (port forwarding)
	{Key: "tunnels_enabled", Default: "1", Kind: "bool", Alias: "TUNNELS_ENABLED"},
	{Key: "tunnel_users", Default: "all", Kind: "enum", Enum: []string{"all", "admins"}},
	{Key: "tunnel_bind_any", Default: "0", Kind: "bool"},
	{Key: "tunnel_remote_forward", Default: "admins", Kind: "enum", Enum: []string{"off", "admins", "all"}},
	{Key: "tunnel_idle_minutes", Default: "0", Kind: "int", Min: 0, Max: 10080},
	// Terminal productivity
	{Key: "broadcast_enabled", Default: "1", Kind: "bool", Alias: "BROADCAST_ENABLED"},
	// Remote desktop (RDP / VNC / Telnet through guacd)
	{Key: "desktop_enabled", Default: "1", Kind: "bool", Alias: "REMOTE_DESKTOP_ENABLED"},
	{Key: "guacd_address", Default: "127.0.0.1:4822", Kind: "string", Alias: "GUACD_ADDRESS"},
	{Key: "desktop_tunnel_bind", Default: "127.0.0.1", Kind: "string"},
	// Out-of-band management (BMC) and serial ports
	{Key: "bmc_enabled", Default: "1", Kind: "bool", Alias: "BMC_ENABLED"},
	{Key: "serial_ports", Default: "admins", Kind: "enum", Enum: []string{"off", "admins", "all"}},
	// Network tools (port check, ping, traceroute, DNS, HTTP/TLS)
	{Key: "network_tools", Default: "all", Kind: "enum", Enum: []string{"off", "admins", "all"}},
	// Live up/down status of connections
	{Key: "status_enabled", Default: "1", Kind: "bool", Alias: "STATUS_ENABLED"},
	{Key: "status_interval_seconds", Default: "60", Kind: "int", Min: 15, Max: 3600},
	{Key: "status_jump_checks", Default: "0", Kind: "bool"},
	// Notifications (e-mail, chat and push channels)
	{Key: "notifications_enabled", Default: "1", Kind: "bool", Alias: "NOTIFICATIONS_ENABLED"},
	{Key: "notify_min_interval_seconds", Default: "60", Kind: "int", Min: 0, Max: 86400},
	// Audit & session recording
	{Key: "audit_enabled", Default: "1", Kind: "bool", Alias: "AUDIT_ENABLED"},
	{Key: "audit_retention_days", Default: "365", Kind: "int", Min: 7, Max: 3650},
	{Key: "session_recording", Default: "1", Kind: "bool", Alias: "SESSION_RECORDING_ENABLED"},
	{Key: "session_recording_input", Default: "0", Kind: "bool"},
	{Key: "recording_retention_days", Default: "90", Kind: "int", Min: 7, Max: 3650},
	{Key: "recording_max_mb", Default: "100", Kind: "int", Min: 0, Max: 102400},
}

var settingsStore = struct {
	sync.RWMutex
	values map[string]string
}{values: map[string]string{}}

func settingSpecFor(key string) (settingSpec, bool) {
	for _, s := range settingSpecs {
		if s.Key == key {
			return s, true
		}
	}
	return settingSpec{}, false
}

func settingEnvName(key string) string { return "WRM_" + strings.ToUpper(key) }

// settingEnv returns the environment override of a setting and the variable it came from.
func settingEnv(key string) (string, string) {
	if v := os.Getenv(settingEnvName(key)); v != "" {
		return v, settingEnvName(key)
	}
	if s, ok := settingSpecFor(key); ok && s.Alias != "" {
		if v := os.Getenv(s.Alias); v != "" {
			return v, s.Alias
		}
	}
	return "", ""
}

// loadAppSettings reads every stored setting into memory (called once at startup).
func loadAppSettings() {
	rows, err := db.Query(`SELECT key, value FROM app_settings`)
	if err != nil {
		log.Printf("settings: %v", err)
		return
	}
	defer rows.Close()
	settingsStore.Lock()
	defer settingsStore.Unlock()
	for rows.Next() {
		var k, v string
		if rows.Scan(&k, &v) == nil {
			settingsStore.values[k] = v
		}
	}
	for _, s := range settingSpecs {
		if env, name := settingEnv(s.Key); env != "" {
			if _, err := validateSetting(s, env); err != nil {
				log.Printf("WARNING: ignoring %s=%q: %v", name, env, err)
			}
		}
	}
}

// getSetting returns the effective value: environment override, stored value or default.
func getSetting(key string) string {
	s, known := settingSpecFor(key)
	if known {
		if env, _ := settingEnv(key); env != "" {
			if v, err := validateSetting(s, env); err == nil {
				return v
			}
		}
	}
	settingsStore.RLock()
	v, ok := settingsStore.values[key]
	settingsStore.RUnlock()
	if ok {
		return v
	}
	return s.Default
}

func settingBool(key string) bool { return getSetting(key) == "1" }

func settingInt(key string) int {
	n, err := strconv.Atoi(getSetting(key))
	if err != nil {
		s, _ := settingSpecFor(key)
		n, _ = strconv.Atoi(s.Default)
	}
	return n
}

func settingLocked(key string) bool {
	env, _ := settingEnv(key)
	return env != ""
}

func validateSetting(s settingSpec, v string) (string, error) {
	v = strings.TrimSpace(v)
	switch s.Kind {
	case "bool":
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return "1", nil
		case "0", "false", "no", "off", "":
			return "0", nil
		}
		return "", fmt.Errorf("expected true/false")
	case "int":
		n, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("expected a number")
		}
		if n < s.Min || (s.Max > 0 && n > s.Max) {
			return "", fmt.Errorf("must be between %d and %d", s.Min, s.Max)
		}
		return strconv.Itoa(n), nil
	case "enum":
		for _, e := range s.Enum {
			if strings.EqualFold(v, e) {
				return e, nil
			}
		}
		return "", fmt.Errorf("must be one of %s", strings.Join(s.Enum, ", "))
	case "json":
		if v == "" {
			return "[]", nil
		}
		var tmp []map[string]interface{}
		if err := json.Unmarshal([]byte(v), &tmp); err != nil {
			return "", fmt.Errorf("invalid JSON: %v", err)
		}
		return v, nil
	}
	if len(v) > 4096 {
		return "", fmt.Errorf("value too long")
	}
	switch s.Key {
	case "turn_public_ip":
		if v != "" && net.ParseIP(v) == nil {
			return "", fmt.Errorf("not an IP address")
		}
	case "turn_relay_ports":
		if _, _, err := parsePortRange(v); err != nil {
			return "", err
		}
	}
	return v, nil
}

func parsePortRange(v string) (uint16, uint16, error) {
	parts := strings.SplitN(strings.TrimSpace(v), "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("use the form 49152-65535")
	}
	lo, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	hi, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || lo < 1024 || hi > 65535 || lo > hi {
		return 0, 0, fmt.Errorf("invalid port range (1024-65535, low-high)")
	}
	return uint16(lo), uint16(hi), nil
}

func setSetting(key, value string) (string, error) {
	s, ok := settingSpecFor(key)
	if !ok {
		return "", fmt.Errorf("unknown setting %q", key)
	}
	if _, name := settingEnv(key); name != "" {
		return "", fmt.Errorf("%s is set by the environment variable %s", key, name)
	}
	v, err := validateSetting(s, value)
	if err != nil {
		return "", fmt.Errorf("%s: %v", key, err)
	}
	if _, err := db.Exec(`INSERT INTO app_settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, v); err != nil {
		return "", err
	}
	settingsStore.Lock()
	settingsStore.values[key] = v
	settingsStore.Unlock()
	return v, nil
}

// Internal (non-admin-editable) values such as the server secret share the table.
func getInternal(key string) string {
	var v string
	db.QueryRow(`SELECT value FROM app_settings WHERE key=?`, "_"+key).Scan(&v)
	return v
}

func setInternal(key, value string) error {
	_, err := db.Exec(`INSERT INTO app_settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`, "_"+key, value)
	return err
}

// ─── ADMIN API ───────────────────────────────────────

// GET  /api/admin/settings → every setting with value, default and lock state
// PUT  /api/admin/settings → {"key": "value", ...}
func apiAdminSettingsHandler(w http.ResponseWriter, r *http.Request) {
	adminID, ok := requireAdmin(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		out := []map[string]interface{}{}
		for _, s := range settingSpecs {
			envName := settingEnvName(s.Key)
			if _, name := settingEnv(s.Key); name != "" {
				envName = name
			}
			out = append(out, map[string]interface{}{
				"key": s.Key, "value": getSetting(s.Key), "default": s.Default, "kind": s.Kind,
				"enum": s.Enum, "min": s.Min, "max": s.Max, "locked": settingLocked(s.Key),
				"env": envName, "restart": s.Restart,
			})
		}
		jsonOK(w, out)
	case http.MethodPut:
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			jsonError(w, "Bad JSON", 400)
			return
		}
		changed := map[string]string{}
		turnChanged := false
		for k, raw := range payload {
			var v string
			switch x := raw.(type) {
			case string:
				v = x
			case bool:
				v = map[bool]string{true: "1", false: "0"}[x]
			case float64:
				v = strconv.FormatFloat(x, 'f', -1, 64)
			default:
				b, _ := json.Marshal(x)
				v = string(b)
			}
			if getSetting(k) == strings.TrimSpace(v) {
				continue
			}
			nv, err := setSetting(k, v)
			if err != nil {
				jsonError(w, err.Error(), 400)
				return
			}
			changed[k] = nv
			if s, _ := settingSpecFor(k); strings.HasPrefix(s.Key, "turn_") {
				turnChanged = true
			}
		}
		if len(changed) > 0 {
			auditLog(r, adminID, "admin.settings", "", changed)
		}
		if turnChanged {
			restartTURN()
		}
		for k := range changed {
			if strings.HasPrefix(k, "tunnel") {
				tunnelMgr.applyPolicy()
				break
			}
		}
		jsonOK(w, map[string]interface{}{"ok": true, "changed": changed})
	default:
		jsonError(w, "Method not allowed", 405)
	}
}
