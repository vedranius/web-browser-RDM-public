package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ─── NETWORK SCANNER ─────────────────────────────────
//
// Finds the devices of a network (IPv4 range) and describes them: IP, MAC address and its
// vendor, names (reverse DNS, NetBIOS, mDNS, UPnP), open ports with their service, the
// SSH / HTTP / TLS banners and a guess of the device type. A scan runs from one of three
// places:
//
//   - the WRM server: ICMP echo (when the OS allows unprivileged ping), TCP probes, NetBIOS
//     node status, unicast mDNS, SSDP; MAC addresses come from the neighbour (ARP) table,
//     which the probes fill even for hosts that drop every packet;
//   - an SSH connection: a ping sweep (fping or ping) and the neighbour table run on the
//     server as plain shell commands, the TCP probes and banners go through SSH channels
//     (direct-tcpip) from WRM, so nothing is installed there;
//   - any computer, with "wrm -netscan RANGE": the same engine as on the WRM server, the
//     result is a JSON file that can be imported into WRM.
//
// Only TCP connects, single UDP queries and echo requests are sent: no SYN scans, no
// exploits, no brute force. Ranges are limited to netscanMaxHosts addresses, and public
// ranges need the administrator's setting network_scan_public.

const netscanMaxHosts = 1024

//go:embed oui.tsv.gz
var ouiGz []byte

var (
	ouiOnce sync.Once
	ouiMap  map[uint32]string
)

// ouiVendor returns the vendor registered for the first three bytes of mac (IEEE MA-L).
func ouiVendor(mac string) string {
	ouiOnce.Do(func() {
		ouiMap = map[uint32]string{}
		zr, err := gzip.NewReader(bytes.NewReader(ouiGz))
		if err != nil {
			return
		}
		sc := bufio.NewScanner(zr)
		for sc.Scan() {
			line := sc.Text()
			if len(line) < 8 || line[6] != '\t' {
				continue
			}
			v, err := strconv.ParseUint(line[:6], 16, 32)
			if err == nil {
				ouiMap[uint32(v)] = line[7:]
			}
		}
	})
	hw, err := net.ParseMAC(mac)
	if err != nil || len(hw) < 3 {
		return ""
	}
	return ouiMap[uint32(hw[0])<<16|uint32(hw[1])<<8|uint32(hw[2])]
}

// macIsRandom reports a locally administered address: phones and laptops use such a
// private (random) address per network.
func macIsRandom(mac string) bool {
	hw, err := net.ParseMAC(mac)
	return err == nil && len(hw) > 0 && hw[0]&0x02 != 0 && hw[0]&0x01 == 0
}

var macRe = regexp.MustCompile(`(?i)\b([0-9a-f]{1,2})([:-])([0-9a-f]{1,2})[:-]([0-9a-f]{1,2})[:-]([0-9a-f]{1,2})[:-]([0-9a-f]{1,2})[:-]([0-9a-f]{1,2})\b`)
var ipv4Re = regexp.MustCompile(`\b(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\b`)

// normMAC writes a MAC as AA:BB:CC:DD:EE:FF; "" for all zeros and broadcast.
func normMAC(m []string) string {
	parts := []string{m[1], m[3], m[4], m[5], m[6], m[7]}
	for i, p := range parts {
		if len(p) == 1 {
			p = "0" + p
		}
		parts[i] = strings.ToUpper(p)
	}
	s := strings.Join(parts, ":")
	if s == "00:00:00:00:00:00" || s == "FF:FF:FF:FF:FF:FF" {
		return ""
	}
	return s
}

// parseNeighbors reads a neighbour table in any of the usual formats (Linux /proc/net/arp
// and "ip neigh", Windows "arp -a", BSD/macOS "arp -an") and returns IP → MAC. Lines
// without a hardware address (incomplete, FAILED) are skipped.
func parseNeighbors(text string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "FAILED") || strings.Contains(line, "INCOMPLETE") || strings.Contains(line, "incomplete") {
			continue
		}
		ip := ipv4Re.FindString(line)
		m := macRe.FindStringSubmatch(line)
		if ip == "" || m == nil {
			continue
		}
		if mac := normMAC(m); mac != "" {
			if _, err := netip.ParseAddr(ip); err == nil {
				out[ip] = mac
			}
		}
	}
	return out
}

// ── targets ──

// parseScanTarget reads "192.168.1.0/24", "192.168.1.7", "192.168.1.10-60" or
// "192.168.1.10-192.168.1.60" and returns the addresses to scan. For a prefix shorter
// than /31 the network and broadcast addresses are left out.
func parseScanTarget(s string) ([]netip.Addr, string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, "", fmt.Errorf("enter a network such as 192.168.1.0/24")
	}
	tooMany := fmt.Errorf("at most %d addresses at once (a /22); narrow the range, e.g. to a /24", netscanMaxHosts)
	if strings.Contains(s, "/") {
		p, err := netip.ParsePrefix(s)
		if err != nil || !p.Addr().Is4() {
			return nil, "", fmt.Errorf("enter an IPv4 network such as 192.168.1.0/24 (IPv6 networks are too large to scan)")
		}
		p = p.Masked()
		bits := p.Bits()
		if 32-bits > 10 {
			return nil, "", tooMany
		}
		var out []netip.Addr
		a := p.Addr()
		n := 1 << (32 - bits)
		for i := 0; i < n; i++ {
			if bits < 31 && (i == 0 || i == n-1) {
				a = a.Next()
				continue
			}
			out = append(out, a)
			a = a.Next()
		}
		return out, p.String(), nil
	}
	if i := strings.Index(s, "-"); i > 0 {
		lo, err := netip.ParseAddr(strings.TrimSpace(s[:i]))
		if err != nil || !lo.Is4() {
			return nil, "", fmt.Errorf("enter a range such as 192.168.1.10-60")
		}
		rest := strings.TrimSpace(s[i+1:])
		hi, err := netip.ParseAddr(rest)
		if err != nil {
			n, err2 := strconv.Atoi(rest)
			if err2 != nil || n < 0 || n > 255 {
				return nil, "", fmt.Errorf("enter a range such as 192.168.1.10-60")
			}
			b := lo.As4()
			b[3] = byte(n)
			hi = netip.AddrFrom4(b)
		}
		if !hi.Is4() || hi.Less(lo) {
			return nil, "", fmt.Errorf("the end of the range is before its start")
		}
		var out []netip.Addr
		for a := lo; ; a = a.Next() {
			out = append(out, a)
			if len(out) > netscanMaxHosts {
				return nil, "", tooMany
			}
			if a == hi {
				break
			}
		}
		return out, lo.String() + "-" + hi.String(), nil
	}
	a, err := netip.ParseAddr(s)
	if err != nil || !a.Is4() {
		return nil, "", fmt.Errorf("enter an IPv4 address, range or network such as 192.168.1.0/24")
	}
	return []netip.Addr{a}, a.String(), nil
}

// scanAddrPrivate: loopback, RFC 1918, carrier-grade NAT, link-local.
func scanAddrPrivate(a netip.Addr) bool {
	if a.IsLoopback() || a.IsPrivate() || a.IsLinkLocalUnicast() {
		return true
	}
	return netip.MustParsePrefix("100.64.0.0/10").Contains(a)
}

// ── ports ──

var scanQuickPorts = []int{21, 22, 23, 53, 80, 135, 139, 443, 445, 554, 631, 1883, 3389, 5000, 5900, 8006, 8080, 8443, 9100, 62078}

var scanExtendedPorts = []int{21, 22, 23, 25, 53, 80, 81, 88, 110, 111, 135, 139, 143, 389, 443, 445, 465, 515, 548, 554, 587, 631, 636,
	873, 902, 993, 995, 1433, 1521, 1723, 1880, 1883, 2049, 2375, 2376, 3000, 3306, 3389, 4840, 5000, 5001, 5060, 5432, 5601, 5900, 5985, 5986,
	6379, 6443, 7547, 8000, 8006, 8008, 8009, 8080, 8081, 8123, 8291, 8443, 8728, 8883, 8888, 9000, 9090, 9100, 9200, 9443, 10000, 10250,
	27017, 32400, 49152, 62078}

// scanDiscoveryPorts are tried on every address to find hosts that drop ICMP.
var scanDiscoveryPorts = []int{80, 443, 22, 445, 139, 3389, 8080, 62078, 554, 53}

var portServices = map[int]string{21: "FTP", 22: "SSH", 23: "Telnet", 25: "SMTP", 53: "DNS", 80: "HTTP", 81: "HTTP", 88: "Kerberos", 110: "POP3",
	111: "RPC", 135: "MS RPC", 139: "NetBIOS", 143: "IMAP", 389: "LDAP", 443: "HTTPS", 445: "SMB", 465: "SMTPS", 515: "LPD", 548: "AFP", 554: "RTSP",
	587: "SMTP", 631: "IPP", 636: "LDAPS", 873: "rsync", 902: "VMware", 993: "IMAPS", 995: "POP3S", 1433: "MS SQL", 1521: "Oracle", 1723: "PPTP",
	1880: "Node-RED", 1883: "MQTT", 2049: "NFS", 2375: "Docker", 2376: "Docker TLS", 3000: "HTTP (dev)", 3306: "MySQL", 3389: "RDP", 4840: "OPC UA",
	5000: "HTTP / Synology", 5001: "HTTPS / Synology", 5060: "SIP", 5432: "PostgreSQL", 5601: "Kibana", 5900: "VNC", 5985: "WinRM", 5986: "WinRM TLS",
	6379: "Redis", 6443: "Kubernetes", 7547: "TR-069", 8000: "HTTP", 8006: "Proxmox", 8008: "HTTP / Chromecast", 8009: "Chromecast", 8080: "HTTP",
	8081: "HTTP", 8123: "Home Assistant", 8291: "MikroTik Winbox", 8443: "HTTPS", 8728: "MikroTik API", 8883: "MQTT TLS", 8888: "HTTP", 9000: "HTTP",
	9090: "HTTP / Cockpit", 9100: "Printer (JetDirect)", 9200: "Elasticsearch", 9443: "HTTPS", 10000: "Webmin", 10250: "Kubelet", 27017: "MongoDB",
	32400: "Plex", 49152: "UPnP", 62078: "Apple iDevice"}

// scanPortsFor returns the ports of a profile (quick, extended) or parses a custom list.
func scanPortsFor(profile, custom string) ([]int, error) {
	switch profile {
	case "", "quick":
		return scanQuickPorts, nil
	case "extended":
		return scanExtendedPorts, nil
	case "none":
		return nil, nil
	case "custom":
		return parsePorts(custom)
	}
	return nil, fmt.Errorf("unknown port profile %q", profile)
}

// ── results ──

type scanPort struct {
	Port    int    `json:"port"`
	Service string `json:"service,omitempty"`
	Banner  string `json:"banner,omitempty"`
}

type scanDevice struct {
	IP         string     `json:"ip"`
	MAC        string     `json:"mac,omitempty"`
	Vendor     string     `json:"vendor,omitempty"`
	RandomMAC  bool       `json:"random_mac,omitempty"`
	Name       string     `json:"name,omitempty"` // best name
	DNSName    string     `json:"dns_name,omitempty"`
	NetBIOS    string     `json:"netbios,omitempty"`
	Workgroup  string     `json:"workgroup,omitempty"`
	MDNSName   string     `json:"mdns_name,omitempty"`
	UPnP       string     `json:"upnp,omitempty"`  // friendly name
	Model      string     `json:"model,omitempty"` // UPnP manufacturer / model
	Ports      []scanPort `json:"ports,omitempty"`
	HTTPTitle  string     `json:"http_title,omitempty"`
	HTTPServer string     `json:"http_server,omitempty"`
	TLSName    string     `json:"tls_name,omitempty"`
	SSH        string     `json:"ssh,omitempty"`
	Kind       string     `json:"kind"`
	OS         string     `json:"os,omitempty"`
	RTT        float64    `json:"rtt_ms,omitempty"`
	Seen       []string   `json:"seen"` // what answered: icmp, tcp, arp, netbios, mdns, ssdp
	Gateway    bool       `json:"gateway,omitempty"`
	Self       bool       `json:"self,omitempty"`
	// Filled by the API from the previous scan and the user's labels.
	New   bool   `json:"new,omitempty"`
	Label string `json:"label,omitempty"`
	Note  string `json:"note,omitempty"`
}

func (d *scanDevice) seen(how string) {
	for _, s := range d.Seen {
		if s == how {
			return
		}
	}
	d.Seen = append(d.Seen, how)
}

func (d *scanDevice) hasPort(p int) bool {
	for _, x := range d.Ports {
		if x.Port == p {
			return true
		}
	}
	return false
}

// scanProgress is what a running scan reports.
type scanProgress struct {
	Phase string `json:"phase"` // discover | ports | names | done
	Done  int    `json:"done"`
	Total int    `json:"total"`
}

// scanner holds one scan. dial reaches the targets (directly or through SSH); remote is set
// for scans through an SSH connection.
type scanner struct {
	ctx    context.Context
	remote *ssh.Client
	dial   func(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error)
	addrs  []netip.Addr
	ports  []int

	mu      sync.Mutex
	devices map[string]*scanDevice
	logs    []string
	phase   string
	done    atomic.Int64
	total   atomic.Int64
	localIf map[string]string // own addresses → MAC
	gateway string
	inSet   map[string]bool
}

func newScanner(ctx context.Context, cl *ssh.Client, addrs []netip.Addr, ports []int) *scanner {
	s := &scanner{ctx: ctx, remote: cl, addrs: addrs, ports: ports, devices: map[string]*scanDevice{}, inSet: map[string]bool{}, localIf: map[string]string{}}
	for _, a := range addrs {
		s.inSet[a.String()] = true
	}
	if cl == nil {
		s.dial = func(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
			d := net.Dialer{Timeout: timeout}
			return d.DialContext(ctx, "tcp", addr)
		}
	} else {
		s.dial = func(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
			type res struct {
				c   net.Conn
				err error
			}
			ch := make(chan res, 1)
			go func() {
				c, err := cl.Dial("tcp", addr)
				ch <- res{c, err}
			}()
			select {
			case r := <-ch:
				return r.c, r.err
			case <-time.After(timeout):
			case <-ctx.Done():
			}
			go func() {
				if r := <-ch; r.c != nil {
					r.c.Close()
				}
			}()
			return nil, fmt.Errorf("i/o timeout")
		}
	}
	return s
}

func (s *scanner) logf(format string, args ...interface{}) {
	s.mu.Lock()
	if len(s.logs) < 200 {
		s.logs = append(s.logs, time.Now().Format("15:04:05")+" "+fmt.Sprintf(format, args...))
	}
	s.mu.Unlock()
}

func (s *scanner) setPhase(p string, total int) {
	s.mu.Lock()
	s.phase = p
	s.mu.Unlock()
	s.done.Store(0)
	s.total.Store(int64(total))
}

func (s *scanner) progress() scanProgress {
	s.mu.Lock()
	defer s.mu.Unlock()
	return scanProgress{Phase: s.phase, Done: int(s.done.Load()), Total: int(s.total.Load())}
}

// device returns the device of ip, creating it (only for addresses of the scan).
func (s *scanner) device(ip string) *scanDevice {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d := s.devices[ip]; d != nil {
		return d
	}
	if !s.inSet[ip] {
		return nil
	}
	d := &scanDevice{IP: ip, Seen: []string{}}
	s.devices[ip] = d
	return d
}

func (s *scanner) existing(ip string) *scanDevice {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.devices[ip]
}

// update runs f on the device of ip under the scanner's lock.
func (s *scanner) update(ip string, create bool, f func(d *scanDevice)) {
	var d *scanDevice
	if create {
		d = s.device(ip)
	} else {
		d = s.existing(ip)
	}
	if d == nil {
		return
	}
	s.mu.Lock()
	f(d)
	s.mu.Unlock()
}

// results returns copies of the devices sorted by address, with their vendor, kind and
// best name filled in (safe to use while the scan goes on).
func (s *scanner) results() []scanDevice {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]scanDevice, 0, len(s.devices))
	for _, src := range s.devices {
		d := *src
		d.Ports = append([]scanPort(nil), src.Ports...)
		d.Seen = append([]string{}, src.Seen...)
		if d.MAC == "" {
			d.MAC = s.localIf[d.IP]
		}
		if d.MAC != "" {
			d.Vendor = ouiVendor(d.MAC)
			d.RandomMAC = macIsRandom(d.MAC)
		}
		if _, ok := s.localIf[d.IP]; ok {
			d.Self = true
		}
		d.Gateway = d.IP == s.gateway
		sort.Slice(d.Ports, func(i, j int) bool { return d.Ports[i].Port < d.Ports[j].Port })
		d.Name = firstText(d.UPnP, d.NetBIOS, strings.TrimSuffix(d.MDNSName, ".local"), shortHost(d.DNSName), d.TLSName)
		d.Kind, d.OS = classifyDevice(&d)
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		a, _ := netip.ParseAddr(out[i].IP)
		b, _ := netip.ParseAddr(out[j].IP)
		return a.Less(b)
	})
	return out
}

func (s *scanner) logLines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string{}, s.logs...)
}

func firstText(v ...string) string {
	for _, s := range v {
		if s = strings.TrimSpace(s); s != "" {
			return s
		}
	}
	return ""
}

func shortHost(h string) string {
	return strings.TrimSuffix(h, ".")
}

// run performs the whole scan.
func (s *scanner) run() {
	if s.remote == nil {
		s.localInfo()
		s.discoverLocal()
	} else {
		s.discoverRemote()
	}
	if s.ctx.Err() != nil {
		return
	}
	s.probePorts()
	if s.ctx.Err() != nil {
		return
	}
	s.identify()
	s.setPhase("done", 0)
}

// ── own interfaces and gateway (WRM server / CLI) ──

func (s *scanner) localInfo() {
	ifs, _ := net.Interfaces()
	for _, ifc := range ifs {
		addrs, _ := ifc.Addrs()
		for _, a := range addrs {
			if ipn, ok := a.(*net.IPNet); ok && ipn.IP.To4() != nil {
				s.mu.Lock()
				s.localIf[ipn.IP.String()] = strings.ToUpper(ifc.HardwareAddr.String())
				s.mu.Unlock()
			}
		}
	}
	gw := localGateway()
	s.mu.Lock()
	s.gateway = gw
	s.mu.Unlock()
}

// localGateway returns the default gateway of the WRM server (Linux; elsewhere from
// "route print" / "netstat -rn" where available).
func localGateway() string {
	if b, err := os.ReadFile("/proc/net/route"); err == nil {
		for _, line := range strings.Split(string(b), "\n")[1:] {
			f := strings.Fields(line)
			if len(f) > 2 && f[1] == "00000000" {
				if v, err := strconv.ParseUint(f[2], 16, 32); err == nil && v != 0 {
					var ip [4]byte
					binary.LittleEndian.PutUint32(ip[:], uint32(v))
					return netip.AddrFrom4(ip).String()
				}
			}
		}
		return ""
	}
	var out []byte
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch runtime.GOOS {
	case "windows":
		out, _ = exec.CommandContext(ctx, "route", "print", "-4", "0.0.0.0").Output()
		for _, line := range strings.Split(string(out), "\n") {
			f := strings.Fields(line)
			if len(f) >= 3 && f[0] == "0.0.0.0" && f[1] == "0.0.0.0" && net.ParseIP(f[2]) != nil {
				return f[2]
			}
		}
	case "darwin", "freebsd", "openbsd", "netbsd":
		out, _ = exec.CommandContext(ctx, "netstat", "-rn", "-f", "inet").Output()
		for _, line := range strings.Split(string(out), "\n") {
			f := strings.Fields(line)
			if len(f) >= 2 && f[0] == "default" && net.ParseIP(f[1]) != nil {
				return f[1]
			}
		}
	}
	return ""
}

// localNeighbors reads the neighbour table of the WRM server.
func localNeighbors() map[string]string {
	if b, err := os.ReadFile("/proc/net/arp"); err == nil {
		return parseNeighbors(string(b))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	args := []string{"-an"}
	if runtime.GOOS == "windows" {
		args = []string{"-a"}
	}
	path, err := exec.LookPath("arp")
	if err != nil {
		return map[string]string{}
	}
	out, _ := exec.CommandContext(ctx, path, args...).Output()
	return parseNeighbors(string(out))
}

// ── discovery from the WRM server ──

func (s *scanner) discoverLocal() {
	s.setPhase("discover", len(s.addrs))
	s.logf("discovery of %d addresses: ICMP, TCP, NetBIOS, mDNS, SSDP", len(s.addrs))
	var wg sync.WaitGroup
	wg.Add(4)
	go func() { defer wg.Done(); s.icmpSweep() }()
	go func() { defer wg.Done(); s.udpQueries() }()
	go func() { defer wg.Done(); s.ssdpSearch() }()
	go func() { defer wg.Done(); s.tcpSweep() }()
	wg.Wait()
	n := 0
	for ip, mac := range localNeighbors() {
		if !s.inSet[ip] {
			continue
		}
		s.update(ip, true, func(d *scanDevice) {
			if d.MAC == "" {
				d.MAC = mac
			}
			d.seen("arp")
		})
		n++
	}
	s.mu.Lock()
	own := map[string]string{}
	for ip, mac := range s.localIf {
		own[ip] = mac
	}
	s.mu.Unlock()
	for ip, mac := range own {
		if s.inSet[ip] {
			s.update(ip, true, func(d *scanDevice) {
				if d.MAC == "" {
					d.MAC = mac
				}
				d.seen("self")
			})
		}
	}
	s.logf("neighbour table: %d addresses of the range", n)
}

// tcpSweep connects to a few common ports of every address; "refused" also proves a host.
func (s *scanner) tcpSweep() {
	sem := make(chan struct{}, 256)
	var wg sync.WaitGroup
	for _, a := range s.addrs {
		if s.ctx.Err() != nil {
			break
		}
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			defer s.done.Add(1)
			var hwg sync.WaitGroup
			for _, p := range scanDiscoveryPorts {
				hwg.Add(1)
				go func(p int) {
					defer hwg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()
					start := time.Now()
					c, err := s.dial(s.ctx, net.JoinHostPort(ip, strconv.Itoa(p)), 900*time.Millisecond)
					st, _ := classifyDialError(err)
					if c != nil {
						c.Close()
					}
					if st != "open" && st != "closed" {
						return
					}
					rtt := float64(time.Since(start).Microseconds()) / 1000
					s.update(ip, true, func(d *scanDevice) {
						d.seen("tcp")
						if st == "open" && !d.hasPort(p) {
							d.Ports = append(d.Ports, scanPort{Port: p, Service: portServices[p]})
						}
						if d.RTT == 0 || rtt < d.RTT {
							d.RTT = rtt
						}
					})
				}(p)
			}
			hwg.Wait()
		}(a.String())
	}
	wg.Wait()
}

// icmpSweep sends one echo request to every address over an unprivileged ICMP socket
// (Linux with net.ipv4.ping_group_range, macOS) or a raw one (root); without either it
// does nothing and the other probes find the hosts.
func (s *scanner) icmpSweep() {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		if conn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0"); err != nil {
			s.logf("ICMP echo not available on this system (%s); TCP and ARP find the hosts", shortNetError(err))
			return
		}
	}
	defer conn.Close()
	privileged := conn.LocalAddr().Network() == "ip"
	id := rand.IntN(0xffff)
	sent := map[string]time.Time{}
	var smu sync.Mutex
	stop := make(chan struct{})
	go func() {
		buf := make([]byte, 1500)
		for {
			conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
			n, peer, err := conn.ReadFrom(buf)
			select {
			case <-stop:
				return
			default:
			}
			if err != nil {
				continue
			}
			m, err := icmp.ParseMessage(1, buf[:n])
			if err != nil || m.Type != ipv4.ICMPTypeEchoReply {
				continue
			}
			ip := ""
			switch a := peer.(type) {
			case *net.UDPAddr:
				ip = a.IP.String()
			case *net.IPAddr:
				ip = a.IP.String()
			}
			smu.Lock()
			t0, ok := sent[ip]
			smu.Unlock()
			if !ok {
				continue
			}
			rtt := float64(time.Since(t0).Microseconds()) / 1000
			s.update(ip, true, func(d *scanDevice) {
				d.seen("icmp")
				if d.RTT == 0 || rtt < d.RTT {
					d.RTT = rtt
				}
			})
		}
	}()
	for i, a := range s.addrs {
		if s.ctx.Err() != nil {
			break
		}
		msg := icmp.Message{Type: ipv4.ICMPTypeEcho, Body: &icmp.Echo{ID: id, Seq: i & 0xffff, Data: []byte("WRM-netscan")}}
		b, _ := msg.Marshal(nil)
		var dst net.Addr = &net.UDPAddr{IP: a.AsSlice()}
		if privileged {
			dst = &net.IPAddr{IP: a.AsSlice()}
		}
		smu.Lock()
		sent[a.String()] = time.Now()
		smu.Unlock()
		conn.WriteTo(b, dst)
		if i%32 == 31 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	select {
	case <-time.After(1500 * time.Millisecond):
	case <-s.ctx.Done():
	}
	close(stop)
}

// udpQueries sends a NetBIOS node status query (port 137) and a unicast mDNS reverse
// lookup (port 5353) to every address and collects the answers.
func (s *scanner) udpQueries() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		s.logf("NetBIOS / mDNS: %s", shortNetError(err))
		return
	}
	defer conn.Close()
	stop := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		buf := make([]byte, 2048)
		for {
			conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
			n, peer, err := conn.ReadFromUDP(buf)
			select {
			case <-stop:
				return
			default:
			}
			if err != nil || peer == nil {
				continue
			}
			ip := peer.IP.To4().String()
			pkt := append([]byte(nil), buf[:n]...)
			switch peer.Port {
			case 137:
				if name, group, mac, ok := parseNBStat(pkt); ok {
					s.update(ip, true, func(d *scanDevice) {
						d.seen("netbios")
						d.NetBIOS, d.Workgroup = name, group
						if d.MAC == "" && mac != "" {
							d.MAC = mac
						}
					})
				}
			case 5353:
				if name := parseMDNSPTR(pkt); name != "" {
					s.update(ip, true, func(d *scanDevice) {
						d.seen("mdns")
						d.MDNSName = name
					})
				}
			}
		}
	}()
	nb := nbstatQuery(uint16(rand.IntN(0xffff)))
	for i, a := range s.addrs {
		if s.ctx.Err() != nil {
			break
		}
		ip := a.AsSlice()
		conn.WriteToUDP(nb, &net.UDPAddr{IP: ip, Port: 137})
		conn.WriteToUDP(mdnsPTRQuery(a), &net.UDPAddr{IP: ip, Port: 5353})
		if i%32 == 31 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	select {
	case <-time.After(2 * time.Second):
	case <-s.ctx.Done():
	}
	close(stop)
	<-finished
}

// nbstatQuery builds a NetBIOS node status request for the name "*".
func nbstatQuery(id uint16) []byte {
	b := make([]byte, 0, 50)
	b = binary.BigEndian.AppendUint16(b, id)
	b = append(b, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0)
	b = append(b, 0x20)
	name := append([]byte{'*'}, make([]byte, 15)...)
	for _, c := range name {
		b = append(b, 'A'+(c>>4), 'A'+(c&0x0f))
	}
	b = append(b, 0, 0, 0x21, 0, 1)
	return b
}

// parseNBStat reads a node status response: the workstation name, the workgroup and
// the MAC address.
func parseNBStat(p []byte) (name, group, mac string, ok bool) {
	if len(p) < 12 || binary.BigEndian.Uint16(p[6:8]) == 0 {
		return
	}
	i, err := skipDNSName(p, 12)
	if err != nil || i+10 > len(p) {
		return
	}
	if binary.BigEndian.Uint16(p[i:i+2]) != 0x21 {
		return
	}
	i += 10 // type, class, TTL, length
	if i >= len(p) {
		return
	}
	count := int(p[i])
	i++
	for k := 0; k < count; k++ {
		if i+18 > len(p) {
			return
		}
		raw := strings.TrimRight(string(p[i:i+15]), " \x00")
		suffix := p[i+15]
		flags := binary.BigEndian.Uint16(p[i+16 : i+18])
		i += 18
		if suffix != 0x00 || !printableASCII(raw) {
			continue
		}
		if flags&0x8000 != 0 {
			if group == "" {
				group = raw
			}
		} else if name == "" {
			name = raw
		}
	}
	if i+6 <= len(p) {
		mac = strings.ToUpper(net.HardwareAddr(p[i : i+6]).String())
		if mac == "00:00:00:00:00:00" {
			mac = ""
		}
	}
	return name, group, mac, name != "" || group != ""
}

func printableASCII(s string) bool {
	for _, r := range s {
		if r < 0x20 || r > 0x7e {
			return false
		}
	}
	return s != ""
}

// mdnsPTRQuery asks for the PTR record of a (d.c.b.a.in-addr.arpa).
func mdnsPTRQuery(a netip.Addr) []byte {
	v := a.As4()
	b := []byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	for _, l := range []string{strconv.Itoa(int(v[3])), strconv.Itoa(int(v[2])), strconv.Itoa(int(v[1])), strconv.Itoa(int(v[0])), "in-addr", "arpa"} {
		b = append(b, byte(len(l)))
		b = append(b, l...)
	}
	return append(b, 0, 0, 12, 0, 1)
}

// parseMDNSPTR returns the name of the first PTR answer of a DNS message.
func parseMDNSPTR(p []byte) string {
	if len(p) < 12 || p[2]&0x80 == 0 {
		return ""
	}
	qd, an := int(binary.BigEndian.Uint16(p[4:6])), int(binary.BigEndian.Uint16(p[6:8]))
	i := 12
	var err error
	for k := 0; k < qd; k++ {
		if i, err = skipDNSName(p, i); err != nil || i+4 > len(p) {
			return ""
		}
		i += 4
	}
	for k := 0; k < an; k++ {
		if i, err = skipDNSName(p, i); err != nil || i+10 > len(p) {
			return ""
		}
		typ := binary.BigEndian.Uint16(p[i : i+2])
		rdlen := int(binary.BigEndian.Uint16(p[i+8 : i+10]))
		i += 10
		if i+rdlen > len(p) {
			return ""
		}
		if typ == 12 {
			name, err := readDNSName(p, i)
			if err == nil && name != "" {
				return truncateStr(name, 120)
			}
		}
		i += rdlen
	}
	return ""
}

func skipDNSName(p []byte, i int) (int, error) {
	for {
		if i >= len(p) {
			return 0, io.ErrUnexpectedEOF
		}
		l := int(p[i])
		switch {
		case l == 0:
			return i + 1, nil
		case l&0xc0 == 0xc0:
			return i + 2, nil
		}
		i += l + 1
	}
}

func readDNSName(p []byte, i int) (string, error) {
	var labels []string
	for hops := 0; hops < 20; {
		if i >= len(p) {
			return "", io.ErrUnexpectedEOF
		}
		l := int(p[i])
		switch {
		case l == 0:
			return strings.Join(labels, "."), nil
		case l&0xc0 == 0xc0:
			if i+1 >= len(p) {
				return "", io.ErrUnexpectedEOF
			}
			i = int(binary.BigEndian.Uint16(p[i:i+2]) & 0x3fff)
			hops++
			continue
		}
		if i+1+l > len(p) {
			return "", io.ErrUnexpectedEOF
		}
		lab := string(p[i+1 : i+1+l])
		if !printableASCII(lab) {
			return "", errors.New("bad label")
		}
		labels = append(labels, lab)
		i += l + 1
	}
	return "", errors.New("too many pointers")
}

// ssdpSearch asks UPnP devices to announce themselves and reads their description.
func (s *scanner) ssdpSearch() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return
	}
	defer conn.Close()
	req := "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 1\r\nST: ssdp:all\r\n\r\n"
	dst := &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900}
	locations := map[string]string{}
	servers := map[string]string{}
	for k := 0; k < 2; k++ {
		conn.WriteToUDP([]byte(req), dst)
	}
	deadline := time.Now().Add(2500 * time.Millisecond)
	buf := make([]byte, 4096)
	for time.Now().Before(deadline) && s.ctx.Err() == nil {
		conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, peer, err := conn.ReadFromUDP(buf)
		if err != nil || peer == nil {
			continue
		}
		ip := peer.IP.To4().String()
		if !s.inSet[ip] {
			continue
		}
		for _, line := range strings.Split(string(buf[:n]), "\r\n") {
			k, v, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			switch strings.ToUpper(strings.TrimSpace(k)) {
			case "LOCATION":
				if _, ok := locations[ip]; !ok {
					locations[ip] = strings.TrimSpace(v)
				}
			case "SERVER":
				servers[ip] = truncateStr(strings.TrimSpace(v), 120)
			}
		}
		s.update(ip, true, func(d *scanDevice) { d.seen("ssdp") })
	}
	var wg sync.WaitGroup
	for ip, loc := range locations {
		wg.Add(1)
		go func(ip, loc string) {
			defer wg.Done()
			name, model := upnpDescription(s.ctx, ip, loc)
			s.update(ip, true, func(d *scanDevice) {
				d.UPnP, d.Model = name, model
				if d.Model == "" {
					d.Model = servers[ip]
				}
			})
		}(ip, loc)
	}
	wg.Wait()
}

// upnpDescription fetches the device description from the device that announced it (the
// URL must point at the same address, so a device cannot make WRM fetch other hosts).
func upnpDescription(ctx context.Context, ip, loc string) (string, string) {
	u, err := url.Parse(loc)
	if err != nil || u.Scheme != "http" || u.Hostname() != ip {
		return "", ""
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	hc := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		Transport: &http.Transport{Proxy: nil, DisableKeepAlives: true, DialContext: (&net.Dialer{Timeout: 2 * time.Second}).DialContext}}
	resp, err := hc.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()
	var desc struct {
		Device struct {
			FriendlyName string `xml:"friendlyName"`
			Manufacturer string `xml:"manufacturer"`
			ModelName    string `xml:"modelName"`
			ModelNumber  string `xml:"modelNumber"`
			DeviceType   string `xml:"deviceType"`
		} `xml:"device"`
	}
	if xml.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&desc) != nil {
		return "", ""
	}
	d := desc.Device
	model := strings.TrimSpace(strings.Join(strings.Fields(d.Manufacturer+" "+d.ModelName+" "+d.ModelNumber), " "))
	if t := upnpType(d.DeviceType); t != "" {
		model = strings.TrimSpace(model + " · " + t)
	}
	return truncateStr(strings.TrimSpace(d.FriendlyName), 80), truncateStr(model, 120)
}

// upnpType shortens urn:schemas-upnp-org:device:MediaRenderer:1 to MediaRenderer.
func upnpType(t string) string {
	parts := strings.Split(t, ":")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return ""
}

// ── discovery through an SSH connection ──

// remoteDiscoverScript reads the addresses from stdin, pings them (fping if installed,
// otherwise ping in batches of 64) and prints "P <ip> <ms>" for every answer, then the
// neighbour table, the default route and the server's own addresses.
const remoteDiscoverScript = `hosts=$(cat)
if command -v fping >/dev/null 2>&1; then
  printf '%s\n' "$hosts" | fping -a -e -r 0 -t 400 2>/dev/null | sed -n 's/^\([0-9.]*\) *(\([0-9.]*\) ms).*/P \1 \2/p'
elif command -v ping >/dev/null 2>&1; then
  n=0
  for h in $hosts; do
    ( t=$(ping -c 1 -W 1 "$h" 2>/dev/null | sed -n 's/.*time[=<]\([0-9.]*\).*/\1/p' | head -n 1); [ -n "$t" ] && echo "P $h $t" ) &
    n=$((n+1)); if [ "$n" -ge 64 ]; then wait; n=0; fi
  done
  wait
fi
echo "#neigh"
ip neigh show 2>/dev/null || arp -an 2>/dev/null || cat /proc/net/arp 2>/dev/null
echo "#route"
ip route show default 2>/dev/null || netstat -rn 2>/dev/null | grep -E '^(default|0\.0\.0\.0)'
echo "#self"
ip -o -4 addr show 2>/dev/null; ip -o link show 2>/dev/null
exit 0`

func (s *scanner) discoverRemote() {
	s.setPhase("discover", len(s.addrs))
	s.logf("ping sweep of %d addresses on the server, then its neighbour table", len(s.addrs))
	var list strings.Builder
	for _, a := range s.addrs {
		list.WriteString(a.String() + "\n")
	}
	timeout := time.Duration(30+len(s.addrs)/16) * time.Second
	out, err := runRemoteScriptLimit(s.remote, remoteDiscoverScript, list.String(), timeout, 4<<20)
	if err != nil && out == "" {
		s.logf("server: %s", err)
	}
	s.done.Store(int64(len(s.addrs)))
	s.applyRemoteDiscovery(out)
	// Hosts that answer neither ping nor ARP (another subnet behind a router) can still
	// have open ports: try the discovery ports through SSH for small ranges.
	if len(s.addrs) <= 256 && s.ctx.Err() == nil {
		s.setPhase("discover", len(s.addrs))
		s.tcpSweepRemote()
	}
}

// applyRemoteDiscovery reads the output of remoteDiscoverScript.
func (s *scanner) applyRemoteDiscovery(out string) {
	section := ""
	var neigh, self strings.Builder
	pinged := 0
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.HasPrefix(line, "#") {
			section = strings.TrimPrefix(line, "#")
			continue
		}
		switch section {
		case "":
			f := strings.Fields(line)
			if len(f) >= 2 && f[0] == "P" {
				rtt, _ := strconv.ParseFloat(firstOr(f[2:], ""), 64)
				s.update(f[1], true, func(d *scanDevice) {
					d.seen("icmp")
					if rtt > 0 {
						d.RTT = rtt
					}
				})
				pinged++
			}
		case "neigh":
			neigh.WriteString(line + "\n")
		case "route":
			if m := gatewayRe.FindStringSubmatch(line); m != nil {
				s.mu.Lock()
				if s.gateway == "" {
					s.gateway = m[1]
				}
				s.mu.Unlock()
			}
		case "self":
			self.WriteString(line + "\n")
		}
	}
	n := 0
	for ip, mac := range parseNeighbors(neigh.String()) {
		if !s.inSet[ip] {
			continue
		}
		n++
		s.update(ip, true, func(d *scanDevice) {
			d.MAC = mac
			d.seen("arp")
		})
	}
	for _, ifc := range parseIPAddrOutput(self.String()) {
		for _, c := range ifc.Addrs {
			ip := strings.Split(c, "/")[0]
			s.mu.Lock()
			s.localIf[ip] = ifc.MAC
			s.mu.Unlock()
			if s.inSet[ip] {
				s.update(ip, true, func(d *scanDevice) { d.seen("self") })
			}
		}
	}
	s.logf("server: %d answered ping, %d in the neighbour table", pinged, n)
}

var gatewayRe = regexp.MustCompile(`(?:via|default)\s+(\d+\.\d+\.\d+\.\d+)`)

func firstOr(v []string, def string) string {
	if len(v) > 0 {
		return v[0]
	}
	return def
}

// tcpSweepRemote tries a few ports of the addresses nothing answered for, through SSH.
func (s *scanner) tcpSweepRemote() {
	sem := make(chan struct{}, 24)
	var wg sync.WaitGroup
	for _, a := range s.addrs {
		ip := a.String()
		if s.existing(ip) != nil || s.ctx.Err() != nil {
			s.done.Add(1)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer s.done.Add(1)
			for _, p := range []int{22, 80, 443, 445, 3389} {
				sem <- struct{}{}
				c, err := s.dial(s.ctx, net.JoinHostPort(ip, strconv.Itoa(p)), 1500*time.Millisecond)
				<-sem
				st, _ := classifyDialError(err)
				if c != nil {
					c.Close()
				}
				if st == "open" || st == "closed" {
					s.update(ip, true, func(d *scanDevice) {
						d.seen("tcp")
						if st == "open" && !d.hasPort(p) {
							d.Ports = append(d.Ports, scanPort{Port: p, Service: portServices[p]})
						}
					})
					return
				}
				if st == "filtered" && strings.Contains(fmt.Sprint(err), "unreachable") {
					return
				}
			}
		}()
	}
	wg.Wait()
}

// ── port probes ──

func (s *scanner) aliveIPs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.devices))
	for ip := range s.devices {
		out = append(out, ip)
	}
	sort.Strings(out)
	return out
}

func (s *scanner) probePorts() {
	ips := s.aliveIPs()
	if len(s.ports) == 0 || len(ips) == 0 {
		return
	}
	s.setPhase("ports", len(ips))
	s.logf("port check of %d devices (%d ports each)", len(ips), len(s.ports))
	limit := 200
	if s.remote != nil {
		limit = 24
	}
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			defer s.done.Add(1)
			var hwg sync.WaitGroup
			for _, p := range s.ports {
				if d := s.existing(ip); d != nil {
					s.mu.Lock()
					has := d.hasPort(p)
					s.mu.Unlock()
					if has {
						continue
					}
				}
				hwg.Add(1)
				go func(p int) {
					defer hwg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()
					if s.ctx.Err() != nil {
						return
					}
					c, err := s.dial(s.ctx, net.JoinHostPort(ip, strconv.Itoa(p)), 1500*time.Millisecond)
					if err != nil {
						return
					}
					c.Close()
					s.update(ip, false, func(d *scanDevice) {
						if !d.hasPort(p) {
							d.Ports = append(d.Ports, scanPort{Port: p, Service: portServices[p]})
						}
					})
				}(p)
			}
			hwg.Wait()
		}(ip)
	}
	wg.Wait()
}

// ── names and banners ──

var htmlTitleRe = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

func (s *scanner) identify() {
	ips := s.aliveIPs()
	s.setPhase("names", len(ips))
	sem := make(chan struct{}, 32)
	var wg sync.WaitGroup
	var names map[string]string
	if s.remote != nil {
		names = s.remoteNames(ips)
	}
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			defer s.done.Add(1)
			sem <- struct{}{}
			defer func() { <-sem }()
			if s.ctx.Err() != nil {
				return
			}
			var name string
			if s.remote == nil {
				ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
				if n, err := net.DefaultResolver.LookupAddr(ctx, ip); err == nil && len(n) > 0 {
					name = n[0]
				}
				cancel()
			} else {
				name = names[ip]
			}
			d := s.existing(ip)
			if d == nil {
				return
			}
			s.mu.Lock()
			ports := append([]scanPort(nil), d.Ports...)
			s.mu.Unlock()
			b := s.banners(ip, ports)
			s.update(ip, false, func(d *scanDevice) {
				d.DNSName = truncateStr(name, 120)
				for i := range d.Ports {
					if v := b.ports[d.Ports[i].Port]; v != "" {
						d.Ports[i].Banner = v
					}
				}
				d.SSH, d.HTTPTitle, d.HTTPServer, d.TLSName = b.ssh, b.title, b.server, b.tlsName
			})
		}(ip)
	}
	wg.Wait()
}

// remoteNames resolves the addresses on the SSH server (its resolver and /etc/hosts).
func (s *scanner) remoteNames(ips []string) map[string]string {
	out := map[string]string{}
	script := `command -v getent >/dev/null 2>&1 || exit 0
while read -r ip; do n=$(getent hosts "$ip" 2>/dev/null | awk '{print $2; exit}'); [ -n "$n" ] && echo "$ip $n"; done; exit 0`
	res, _ := runRemoteScript(s.remote, script, strings.Join(ips, "\n")+"\n", time.Duration(20+len(ips)/4)*time.Second)
	for _, line := range strings.Split(res, "\n") {
		f := strings.Fields(line)
		if len(f) == 2 && s.inSet[f[0]] && f[1] != f[0] {
			out[f[0]] = f[1]
		}
	}
	return out
}

type bannerInfo struct {
	ssh, title, server, tlsName string
	ports                       map[int]string
}

// banners reads the SSH greeting and the answer of the first web port.
func (s *scanner) banners(ip string, ports []scanPort) bannerInfo {
	b := bannerInfo{ports: map[int]string{}}
	web := 0
	webTLS := false
	for _, p := range ports {
		switch p.Port {
		case 22:
			if v := s.readGreeting(ip, 22); strings.HasPrefix(v, "SSH-") {
				b.ssh = v
				b.ports[22] = v
			}
		case 21, 25, 110, 143:
			if v := s.readGreeting(ip, p.Port); v != "" {
				b.ports[p.Port] = v
			}
		case 80, 8080, 8000, 8081, 8888, 5000, 81:
			if web == 0 {
				web = p.Port
			}
		case 443, 8443, 9443, 5001, 8006:
			if web == 0 || !webTLS {
				web, webTLS = p.Port, true
			}
		}
	}
	if web > 0 {
		scheme := "http"
		if webTLS {
			scheme = "https"
		}
		b.title, b.server, b.tlsName = s.httpInfo(ip, web, scheme)
		if v := strings.TrimSpace(firstText(b.server, b.title)); v != "" {
			b.ports[web] = truncateStr(v, 100)
		}
	}
	return b
}

// readGreeting returns the first line a service sends on connect (SSH, FTP, SMTP …).
func (s *scanner) readGreeting(ip string, port int) string {
	c, err := s.dial(s.ctx, net.JoinHostPort(ip, strconv.Itoa(port)), 2*time.Second)
	if err != nil {
		return ""
	}
	defer c.Close()
	ch := make(chan string, 1)
	go func() {
		buf := make([]byte, 256)
		n, _ := io.ReadAtLeast(c, buf, 4)
		ch <- string(buf[:n])
	}()
	select {
	case v := <-ch:
		v, _, _ = strings.Cut(v, "\n")
		v = strings.TrimSpace(v)
		if !printableASCII(v) {
			return ""
		}
		return truncateStr(v, 100)
	case <-time.After(2500 * time.Millisecond):
		c.Close()
		return ""
	}
}

// httpInfo requests / once and returns the page title, the Server header and, for HTTPS,
// the certificate's name.
func (s *scanner) httpInfo(ip string, port int, scheme string) (title, server, tlsName string) {
	tr := &http.Transport{Proxy: nil, DisableKeepAlives: true, TLSHandshakeTimeout: 4 * time.Second, ResponseHeaderTimeout: 4 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true, //nolint:gosec // only the certificate's name is read, nothing is sent
			VerifyConnection: func(cs tls.ConnectionState) error {
				if len(cs.PeerCertificates) > 0 {
					c := cs.PeerCertificates[0]
					tlsName = c.Subject.CommonName
					if tlsName == "" && len(c.DNSNames) > 0 {
						tlsName = c.DNSNames[0]
					}
				}
				return nil
			}},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return s.dial(ctx, addr, 3*time.Second)
		}}
	hc := &http.Client{Transport: tr, Timeout: 6 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	req, _ := http.NewRequestWithContext(s.ctx, "GET", scheme+"://"+net.JoinHostPort(ip, strconv.Itoa(port))+"/", nil)
	req.Header.Set("User-Agent", "WRM-netscan/"+AppVersion)
	resp, err := hc.Do(req)
	if err != nil {
		return "", "", truncateStr(tlsName, 100)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if m := htmlTitleRe.FindSubmatch(body); m != nil {
		title = strings.Join(strings.Fields(htmlUnescape(string(m[1]))), " ")
	}
	return truncateStr(title, 100), truncateStr(resp.Header.Get("Server"), 80), truncateStr(tlsName, 100)
}

func htmlUnescape(s string) string {
	r := strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'", "&nbsp;", " ")
	return r.Replace(s)
}

// ── device type ──

var vendorKinds = []struct {
	re   *regexp.Regexp
	kind string
}{
	{regexp.MustCompile(`(?i)raspberry`), "sbc"},
	{regexp.MustCompile(`(?i)ubiquiti|mikrotik|routerboard|cisco|juniper|aruba|netgear|zyxel|d-link|tp-link|avm|fortinet|palo alto|meraki|ruckus|draytek|tenda|edgecore|extreme networks|allied telesis|sophos|watchguard|teltonika|cambium`), "network"},
	{regexp.MustCompile(`(?i)synology|qnap|western digital|buffalo|asustor|terramaster|netapp`), "nas"},
	{regexp.MustCompile(`(?i)hikvision|dahua|axis comm|reolink|hanwha|vivotek|uniview|amcrest|foscam`), "camera"},
	{regexp.MustCompile(`(?i)brother|canon|epson|kyocera|xerox|lexmark|ricoh|konica|sharp corp|oki data|zebra`), "printer"},
	{regexp.MustCompile(`(?i)espressif|tuya|shelly|sonoff|itead|xiaomi|philips lighting|signify|ikea|nest labs|ecobee|wiz|lifx|meross|tasmota`), "iot"},
	{regexp.MustCompile(`(?i)roku|sonos|chromecast|amazon|samsung electronics|lg electronics|sony|vizio|tcl|hisense|nvidia|bose|denon|yamaha|panasonic|philips`), "media"},
	{regexp.MustCompile(`(?i)vmware|proxmox|xensource|qemu|microsoft corp.*hyper|parallels|virtualbox|oracle virtualbox`), "vm"},
	{regexp.MustCompile(`(?i)apple`), "apple"},
	{regexp.MustCompile(`(?i)supermicro|dell|hewlett|hpe|lenovo|fujitsu|intel|asrock|gigabyte|micro-star|asustek|realtek|broadcom|super micro`), "computer"},
}

// classifyDevice guesses the device type (router, printer, windows, linux, …) and the OS.
func classifyDevice(d *scanDevice) (kind, os string) {
	has := d.hasPort
	text := strings.ToLower(d.SSH + " " + d.HTTPServer + " " + d.HTTPTitle + " " + d.Model + " " + d.UPnP)
	switch {
	case strings.Contains(text, "ubuntu"):
		os = "Ubuntu"
	case strings.Contains(text, "debian"):
		os = "Debian"
	case strings.Contains(text, "raspbian"):
		os = "Raspberry Pi OS"
	case strings.Contains(text, "freebsd"):
		os = "FreeBSD"
	case strings.Contains(text, "openwrt"):
		os = "OpenWrt"
	case strings.Contains(text, "routeros") || strings.Contains(text, "mikrotik"):
		os = "RouterOS"
	case strings.Contains(text, "dsm") || strings.Contains(text, "synology"):
		os = "Synology DSM"
	case strings.Contains(text, "proxmox"):
		os = "Proxmox VE"
	case strings.Contains(text, "esxi") || strings.Contains(text, "vmware"):
		os = "VMware ESXi"
	case strings.Contains(text, "microsoft") || strings.Contains(text, "windows") || (has(3389) && (has(445) || has(135))):
		os = "Windows"
	case strings.Contains(text, "openssh"):
		os = "Linux / Unix"
	}
	vendorKind := ""
	for _, vk := range vendorKinds {
		if d.Vendor != "" && vk.re.MatchString(d.Vendor) {
			vendorKind = vk.kind
			break
		}
	}
	model := strings.ToLower(d.Model)
	switch {
	case d.Self:
		return "self", os
	case d.Gateway || strings.Contains(model, "internetgatewaydevice") || strings.Contains(model, "wandevice"):
		return "router", os
	case has(9100) || has(631) || has(515) || vendorKind == "printer" || strings.Contains(model, "printer"):
		return "printer", os
	case has(8006) || os == "Proxmox VE" || os == "VMware ESXi" || has(902):
		return "hypervisor", os
	case vendorKind == "camera" || (has(554) && !has(445) && !has(22)):
		return "camera", os
	case vendorKind == "nas" || os == "Synology DSM" || (has(5000) && has(5001)) || has(548) || has(2049):
		return "nas", os
	case strings.Contains(model, "mediarenderer") || strings.Contains(model, "mediaserver") || has(8009) || has(32400) || vendorKind == "media":
		return "media", os
	case has(62078):
		return "phone", "iOS / iPadOS"
	case vendorKind == "network" || os == "RouterOS" || os == "OpenWrt" || has(8291) || has(7547):
		return "network", os
	case os == "Windows" || has(3389) || (has(445) && has(135)):
		return "windows", "Windows"
	case vendorKind == "iot" || has(1883) || has(8883):
		return "iot", os
	case vendorKind == "sbc":
		return "linux", firstText(os, "Linux")
	case has(22) || os != "":
		return "linux", firstText(os, "Linux / Unix")
	case vendorKind == "apple":
		return "apple", os
	case d.RandomMAC && len(d.Ports) == 0:
		return "phone", os
	case vendorKind == "vm":
		return "vm", os
	case vendorKind == "computer" || has(445):
		return "computer", os
	}
	return "unknown", os
}

// ── interfaces ──

// scanInterface is a network card of the WRM server or of an SSH server.
type scanInterface struct {
	Name    string   `json:"name"`
	MAC     string   `json:"mac,omitempty"`
	Addrs   []string `json:"addrs"`             // 192.168.1.10/24
	Suggest string   `json:"suggest,omitempty"` // network to scan (narrowed to /24 for larger ones)
	Virtual bool     `json:"virtual,omitempty"`
	Up      bool     `json:"up"`
}

var virtualIfRe = regexp.MustCompile(`(?i)^(docker|br-|veth|virbr|vmnet|vboxnet|cni|flannel|cali|tun|tap|wg|tailscale|zt|utun|vEthernet|lxc|lxd|podman|kube)`)

// scanSuggest returns the network of addr, narrowed to the /24 around it when larger
// than a scan may be.
func scanSuggest(cidr string) string {
	p, err := netip.ParsePrefix(cidr)
	if err != nil || !p.Addr().Is4() || p.Addr().IsLoopback() || p.Bits() >= 32 {
		return ""
	}
	if 32-p.Bits() > 10 {
		p = netip.PrefixFrom(p.Addr(), 24)
	}
	return p.Masked().String()
}

func localInterfaces() []scanInterface {
	out := []scanInterface{}
	ifs, _ := net.Interfaces()
	for _, ifc := range ifs {
		if ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		si := scanInterface{Name: ifc.Name, MAC: strings.ToUpper(ifc.HardwareAddr.String()), Up: ifc.Flags&net.FlagUp != 0, Virtual: virtualIfRe.MatchString(ifc.Name)}
		addrs, _ := ifc.Addrs()
		for _, a := range addrs {
			if ipn, ok := a.(*net.IPNet); ok && ipn.IP.To4() != nil {
				ones, _ := ipn.Mask.Size()
				c := ipn.IP.String() + "/" + strconv.Itoa(ones)
				si.Addrs = append(si.Addrs, c)
				if si.Suggest == "" {
					si.Suggest = scanSuggest(c)
				}
			}
		}
		if len(si.Addrs) > 0 {
			out = append(out, si)
		}
	}
	sortInterfaces(out)
	return out
}

func sortInterfaces(out []scanInterface) {
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Virtual != out[j].Virtual {
			return !out[i].Virtual
		}
		if out[i].Up != out[j].Up {
			return out[i].Up
		}
		return out[i].Name < out[j].Name
	})
}

// parseIPAddrOutput reads "ip -o -4 addr show" and "ip -o link show".
func parseIPAddrOutput(text string) []scanInterface {
	byName := map[string]*scanInterface{}
	var order []string
	get := func(name string) *scanInterface {
		name = strings.TrimSuffix(name, ":")
		if i := strings.Index(name, "@"); i > 0 {
			name = name[:i]
		}
		if byName[name] == nil {
			byName[name] = &scanInterface{Name: name, Addrs: []string{}, Virtual: virtualIfRe.MatchString(name)}
			order = append(order, name)
		}
		return byName[name]
	}
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(line)
		if len(f) < 4 {
			continue
		}
		switch {
		case f[2] == "inet":
			if f[1] == "lo" {
				continue
			}
			si := get(f[1])
			if _, err := netip.ParsePrefix(f[3]); err == nil {
				si.Addrs = append(si.Addrs, f[3])
				if si.Suggest == "" {
					si.Suggest = scanSuggest(f[3])
				}
			}
			si.Up = true
		case strings.HasSuffix(f[1], ":") && strings.HasPrefix(f[2], "<"):
			name := strings.TrimSuffix(f[1], ":")
			if name == "lo" {
				continue
			}
			si := get(name)
			si.Up = si.Up || strings.Contains(f[2], ",UP")
			for i := 3; i+1 < len(f); i++ {
				if f[i] == "link/ether" {
					si.MAC = strings.ToUpper(f[i+1])
				}
			}
		}
	}
	out := []scanInterface{}
	for _, n := range order {
		if si := byName[n]; len(si.Addrs) > 0 {
			out = append(out, *si)
		}
	}
	sortInterfaces(out)
	return out
}

func remoteInterfaces(cl *ssh.Client) ([]scanInterface, error) {
	out, err := runRemoteScript(cl, `if command -v ip >/dev/null 2>&1; then ip -o -4 addr show; ip -o link show; else echo "NOIP"; fi; exit 0`, "", 20*time.Second)
	if err != nil {
		return nil, err
	}
	if strings.Contains(out, "NOIP") {
		return nil, fmt.Errorf("the server has no \"ip\" command; enter the network by hand")
	}
	return parseIPAddrOutput(out), nil
}

// ── Wake-on-LAN ──

func parseWakeMAC(s string) (net.HardwareAddr, error) {
	m := macRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil || len(strings.TrimSpace(s)) > 17 {
		return nil, fmt.Errorf("enter a MAC address such as 00:11:22:33:44:55")
	}
	hw, _ := hex.DecodeString(strings.ReplaceAll(normMACAny(m), ":", ""))
	if len(hw) != 6 {
		return nil, fmt.Errorf("enter a MAC address such as 00:11:22:33:44:55")
	}
	return net.HardwareAddr(hw), nil
}

func normMACAny(m []string) string {
	parts := []string{m[1], m[3], m[4], m[5], m[6], m[7]}
	for i, p := range parts {
		if len(p) == 1 {
			parts[i] = "0" + p
		}
	}
	return strings.ToUpper(strings.Join(parts, ":"))
}

func magicPacket(hw net.HardwareAddr) []byte {
	p := bytes.Repeat([]byte{0xff}, 6)
	for i := 0; i < 16; i++ {
		p = append(p, hw...)
	}
	return p
}

// sendWake sends the magic packet from the WRM server to the broadcast address (UDP 9).
func sendWake(hw net.HardwareAddr, broadcast string) error {
	if broadcast == "" {
		broadcast = "255.255.255.255"
	}
	c, err := net.Dial("udp4", net.JoinHostPort(broadcast, "9"))
	if err != nil {
		return err
	}
	defer c.Close()
	p := magicPacket(hw)
	for i := 0; i < 3; i++ {
		if _, err := c.Write(p); err != nil {
			return err
		}
	}
	return nil
}

// sendWakeRemote sends it from an SSH server with wakeonlan, etherwake or python3.
func sendWakeRemote(cl *ssh.Client, hw net.HardwareAddr, broadcast string) (string, error) {
	if broadcast == "" {
		broadcast = "255.255.255.255"
	}
	mac, b := shellQuote(strings.ToLower(hw.String())), shellQuote(broadcast)
	script := `if command -v wakeonlan >/dev/null 2>&1; then wakeonlan -i ` + b + ` ` + mac + ` && echo "sent with wakeonlan"
elif command -v python3 >/dev/null 2>&1; then python3 -c 'import socket,sys
m=bytes.fromhex(sys.argv[1].replace(":",""));p=b"\xff"*6+m*16
s=socket.socket(socket.AF_INET,socket.SOCK_DGRAM);s.setsockopt(socket.SOL_SOCKET,socket.SO_BROADCAST,1)
[s.sendto(p,(sys.argv[2],9)) for _ in range(3)]' ` + mac + ` ` + b + ` && echo "sent with python3"
elif command -v etherwake >/dev/null 2>&1; then etherwake ` + mac + ` 2>&1 && echo "sent with etherwake"
else echo "install wakeonlan or python3 on this server"; exit 3; fi`
	return runRemoteScript(cl, script, "", 20*time.Second)
}
