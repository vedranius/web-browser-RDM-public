package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ─── CHAINS: READ-ONLY SERVER DISCOVERY FOR THE WIZARD ───────────────────────────────────
//
// Before the assistant drafts a chain, WRM looks at the selected servers itself: the OS, the
// running non-system services, Docker containers, listening ports and the mounts (NFS, S3,
// rclone, CIFS, SSHFS). That lets the AI choose real checks instead of inventing unit names,
// and lets the deterministic fallback still produce a useful chain.
//
// The script is fixed text, read-only and bounded. It returns names, ports and mount points,
// never file contents, environment variables or credentials. A mount is a fact about the
// server; WRM does not attach to it or store any storage credentials for the chain.

const chainDiscoveryScript = `printf '@host %s\n' "$(hostname 2>/dev/null || uname -n 2>/dev/null)"
printf '@os %s\n' "$(uname -srm 2>/dev/null)"
{ systemctl list-units --type=service --state=running --no-legend --no-pager 2>/dev/null || true; } | awk '{print $1}' | head -n 40 | while IFS= read -r u; do [ -n "$u" ] && printf '@service %s\n' "$u"; done
{ docker ps --format '{{.Names}}' 2>/dev/null || true; } | head -n 20 | while IFS= read -r c; do [ -n "$c" ] && printf '@container %s\n' "$c"; done
{ ss -Hltn 2>/dev/null || netstat -ltn 2>/dev/null || true; } | awk '{print $4}' | sed 's/.*[:.]//' | grep -E '^[0-9]+$' | sort -n -u | head -n 30 | while IFS= read -r p; do [ -n "$p" ] && printf '@port %s\n' "$p"; done
{ findmnt -rn -o TARGET,SOURCE,FSTYPE,OPTIONS 2>/dev/null || true; } | head -n 40 | while IFS= read -r line; do [ -n "$line" ] && printf '@mount %s\n' "$(printf '%s' "$line" | base64 | tr -d '\n')"; done
{ hostname -I 2>/dev/null || true; } | tr ' ' '\n' | head -n 16 | while IFS= read -r a; do [ -n "$a" ] && printf '@addr %s\n' "$a"; done
{ ss -Hltn 2>/dev/null || true; } | awk '{print $4}' | sed 's/.*[:.]//' | grep -E '^[0-9]+$' | sort -n -u | head -n 200 | while IFS= read -r p; do printf '@listen %s\n' "$p"; done
{ ss -Htn state established 2>/dev/null || netstat -tn 2>/dev/null | awk '$6=="ESTABLISHED"{print 0, 0, $4, $5}' || true; } | awk '{print $3, $4}' | head -n 600 | while IFS= read -r l p; do printf '@est %s %s\n' "$l" "$p"; done`

type chainMountFact struct {
	Target  string `json:"target"`
	Source  string `json:"source"`
	FSType  string `json:"fstype,omitempty"`
	Options string `json:"options,omitempty"`
}

type chainServerFacts struct {
	ConnID     int              `json:"conn_id"`
	Name       string           `json:"name"`
	Hostname   string           `json:"hostname,omitempty"`
	OS         string           `json:"os,omitempty"`
	Services   []string         `json:"services,omitempty"`
	Containers []string         `json:"containers,omitempty"`
	Ports      []int            `json:"ports,omitempty"`
	Mounts     []chainMountFact `json:"mounts,omitempty"`
	Error      string           `json:"error,omitempty"`
	// Used by the link discovery of an existing chain only (chains_ops.go), not sent to the AI.
	Addrs  []string       `json:"-"`
	Listen map[int]bool   `json:"-"`
	Est    []chainEstConn `json:"-"`
}

// chainEstConn is one established TCP connection seen on a server: its own and the peer's
// address and port.
type chainEstConn struct {
	LocalIP  string
	LocalP   int
	PeerIP   string
	PeerPort int
}

var chainDiscoveryNoise = map[string]bool{
	"systemd-journald.service": true, "systemd-logind.service": true, "systemd-udevd.service": true,
	"dbus.service": true, "ssh.service": true, "sshd.service": true, "cron.service": true,
	"crond.service": true, "rsyslog.service": true, "chrony.service": true, "chronyd.service": true,
	"polkit.service": true, "accounts-daemon.service": true, "udisks2.service": true,
	"upower.service": true, "avahi-daemon.service": true, "NetworkManager.service": true,
	"wpa_supplicant.service": true, "ModemManager.service": true, "bluetooth.service": true,
	"cups.service": true, "snapd.service": true, "packagekit.service": true,
	"unattended-upgrades.service": true, "apparmor.service": true, "irqbalance.service": true,
	"multipathd.service": true, "getty@tty1.service": true, "user@1000.service": true,
}

// chainDiscoveryInterestingMount keeps network and object-storage mounts; a local ext4 root
// says nothing about a system's storage links.
func chainDiscoveryInterestingMount(m chainMountFact) bool {
	fst := strings.ToLower(m.FSType)
	if strings.HasPrefix(fst, "fuse.") || fst == "nfs" || fst == "nfs4" || fst == "cifs" || fst == "smbfs" || fst == "ceph" || fst == "glusterfs" {
		return true
	}
	if strings.Contains(m.Source, ":") || strings.HasPrefix(m.Source, "//") {
		return true
	}
	return strings.Contains(m.Options, "_netdev")
}

// parseChainDiscovery turns the fixed script's output into facts, with the same rules the
// fallback and the prompt rely on.
func parseChainDiscovery(out string, connID int, name string) chainServerFacts {
	f := chainServerFacts{ConnID: connID, Name: name}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	seenService := map[string]bool{}
	seenPort := map[int]bool{}
	seenMount := map[string]bool{}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "@") {
			continue
		}
		key, val, _ := strings.Cut(line[1:], " ")
		val = strings.TrimSpace(val)
		switch key {
		case "host":
			f.Hostname = truncateStr(val, 120)
		case "os":
			f.OS = truncateStr(val, 120)
		case "service":
			if val == "" || seenService[val] || chainDiscoveryNoise[val] || len(f.Services) >= 20 {
				continue
			}
			seenService[val] = true
			f.Services = append(f.Services, truncateStr(val, 120))
		case "container":
			if val == "" || len(f.Containers) >= 12 {
				continue
			}
			f.Containers = append(f.Containers, truncateStr(val, 120))
		case "port":
			n, err := strconv.Atoi(val)
			if err != nil || n < 1 || n > 65535 || n == 22 || seenPort[n] || len(f.Ports) >= 20 {
				continue
			}
			seenPort[n] = true
			f.Ports = append(f.Ports, n)
		case "addr":
			if ip := chainCleanIP(val); ip != "" && len(f.Addrs) < 16 {
				f.Addrs = append(f.Addrs, ip)
			}
		case "listen":
			if n, err := strconv.Atoi(val); err == nil && n > 0 && n < 65536 {
				if f.Listen == nil {
					f.Listen = map[int]bool{}
				}
				f.Listen[n] = true
			}
		case "est":
			a, b, _ := strings.Cut(val, " ")
			lip, lp, ok1 := chainSplitAddr(a)
			pip, pp, ok2 := chainSplitAddr(b)
			if ok1 && ok2 && len(f.Est) < 600 {
				f.Est = append(f.Est, chainEstConn{LocalIP: lip, LocalP: lp, PeerIP: pip, PeerPort: pp})
			}
		case "mount":
			raw, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				continue
			}
			parts := strings.Fields(string(raw))
			if len(parts) < 2 || seenMount[parts[0]] || len(f.Mounts) >= 12 {
				continue
			}
			seenMount[parts[0]] = true
			m := chainMountFact{Target: truncateStr(parts[0], 300), Source: truncateStr(parts[1], 300)}
			if len(parts) > 2 {
				m.FSType = truncateStr(parts[2], 40)
			}
			if len(parts) > 3 {
				m.Options = truncateStr(parts[3], 200)
			}
			if chainDiscoveryInterestingMount(m) {
				f.Mounts = append(f.Mounts, m)
			}
		}
	}
	return f
}

// discoverChainServers runs the fixed script on the selected connections, at most four at
// once, sharing the monitoring slots so a wizard cannot flood the servers.
func discoverChainServers(ctx context.Context, conns []Connection) []chainServerFacts {
	out := make([]chainServerFacts, len(conns))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, c := range conns {
		wg.Add(1)
		go func(i int, c Connection) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			select {
			case monitorSlots <- struct{}{}:
				defer func() { <-monitorSlots }()
			case <-ctx.Done():
				out[i] = chainServerFacts{ConnID: c.ID, Name: c.Name, Error: "timed out"}
				return
			}
			runCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
			defer cancel()
			cl, err := dialSSH(c, nil)
			if err != nil {
				out[i] = chainServerFacts{ConnID: c.ID, Name: c.Name, Error: chainShortErr(err)}
				return
			}
			defer cl.Close()
			raw, err := monitorSSHCommandInput(runCtx, cl, "sh -c "+shellQuote(chainDiscoveryScript), "")
			if err != nil && !strings.Contains(raw, "@os ") && !strings.Contains(raw, "@host ") {
				out[i] = chainServerFacts{ConnID: c.ID, Name: c.Name, Error: chainShortErr(err)}
				return
			}
			out[i] = parseChainDiscovery(raw, c.ID, c.Name)
		}(i, c)
	}
	wg.Wait()
	return out
}

// apiChainDiscover previews the facts the wizard and the assistant will use.
func apiChainDiscover(w http.ResponseWriter, r *http.Request, userID int) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !mcpRateOK("chain-discover:"+strconv.Itoa(userID), 20, 10*time.Minute) {
		jsonError(w, "Too many discoveries in a short time; wait a moment", http.StatusTooManyRequests)
		return
	}
	var in struct {
		ConnIDs []int `json:"conn_ids"`
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		jsonError(w, "The request is not valid JSON", http.StatusBadRequest)
		return
	}
	if len(in.ConnIDs) < 1 || len(in.ConnIDs) > chainAssistConns {
		jsonError(w, "Choose 1-10 servers", http.StatusBadRequest)
		return
	}
	var conns []Connection
	seen := map[int]bool{}
	for _, id := range in.ConnIDs {
		if seen[id] || id <= 0 {
			continue
		}
		seen[id] = true
		c, err := loadConnectionRaw(id)
		if err != nil || c.UserID != userID || !isSSHProtocol(c) {
			jsonError(w, "One of the chosen servers is not one of your SSH connections", http.StatusBadRequest)
			return
		}
		conns = append(conns, c)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	jsonOK(w, map[string]interface{}{"servers": discoverChainServers(ctx, conns)})
}
