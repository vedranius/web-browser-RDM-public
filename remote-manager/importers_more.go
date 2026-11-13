package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ─── MORE IMPORT SOURCES AND AUTO-DETECTION (v14.13.0) ─────────────────────────────
//
// One file picker for every tool WRM imports from. The server looks at the bytes, picks the
// format, and can show a preview (dry run) before anything is written:
//
//   mremoteng      mRemoteNG confCons.xml (also old fully encrypted files)
//   mremoteng-csv  mRemoteNG "Export → mRemoteNG CSV" (semicolon separated, plain passwords)
//   rdm-csv        Devolutions Remote Desktop Manager CSV export (Host, ConnectionType, Group …)
//   rdcman         Microsoft Remote Desktop Connection Manager .rdg
//   rdp            a Windows .rdp file
//   mobaxterm      MobaXterm .mxtsessions export (or the [Bookmarks] of MobaXterm.ini)
//   winscp         WinSCP.ini (sessions with their stored passwords)
//   remmina        a Remmina .remmina profile
//   sshconfig      OpenSSH ~/.ssh/config
//   putty          PuTTY sessions exported from the registry (.reg)
//
// WRM's own JSON export keeps its own endpoint (/api/config/import).

// importFormats lists the auto import formats with a short description for the UI.
var importFormats = []struct{ ID, Name string }{
	{"mremoteng", "mRemoteNG (confCons.xml)"},
	{"mremoteng-csv", "mRemoteNG CSV export"},
	{"rdm-csv", "Remote Desktop Manager CSV"},
	{"rdcman", "Remote Desktop Connection Manager (.rdg)"},
	{"rdp", "Remote Desktop file (.rdp)"},
	{"mobaxterm", "MobaXterm (.mxtsessions)"},
	{"winscp", "WinSCP (WinSCP.ini)"},
	{"remmina", "Remmina (.remmina)"},
	{"sshconfig", "OpenSSH config"},
	{"putty", "PuTTY (.reg)"},
}

// importText decodes an uploaded file: UTF-16 (with or without BOM), UTF-8 (BOM removed) or
// Windows-1250 when the bytes are not valid UTF-8.
func importText(b []byte) string {
	switch {
	case bytes.HasPrefix(b, []byte{0xff, 0xfe}):
		return decodeUTF16(b[2:], false)
	case bytes.HasPrefix(b, []byte{0xfe, 0xff}):
		return decodeUTF16(b[2:], true)
	case len(b) > 4 && b[1] == 0 && b[3] == 0 && b[0] != 0:
		return decodeUTF16(b, false)
	}
	b = bytes.TrimPrefix(b, []byte{0xef, 0xbb, 0xbf})
	if utf8.Valid(b) {
		return string(b)
	}
	return decodeCP1250(b)
}

// detectImportFormat names the format of a file ("" = unknown).
func detectImportFormat(name, text string) string {
	ext := strings.ToLower(path.Ext(strings.ReplaceAll(name, "\\", "/")))
	head := strings.TrimSpace(text)
	if len(head) > 4096 {
		head = head[:4096]
	}
	low := strings.ToLower(head)
	firstLine := strings.TrimSpace(strings.SplitN(head, "\n", 2)[0])
	switch {
	case strings.HasPrefix(low, "windows registry editor") || strings.HasPrefix(low, "regedit4"):
		return "putty"
	case strings.Contains(low, "<rdcman"):
		return "rdcman"
	case regexpConnectionsRoot.MatchString(head):
		return "mremoteng"
	case strings.HasPrefix(firstLine, "Name;Id;Parent;NodeType;"):
		return "mremoteng-csv"
	case csvHasColumns(firstLine, "Host", "ConnectionType"):
		return "rdm-csv"
	case strings.Contains(low, "[bookmarks") && strings.Contains(head, "SubRep="):
		return "mobaxterm"
	case strings.Contains(low, "[sessions\\"):
		return "winscp"
	case strings.Contains(low, "[remmina]"):
		return "remmina"
	case strings.Contains(low, "full address:s:"):
		return "rdp"
	case ext == ".xml" && mrBase64Only.MatchString(head):
		return "mremoteng" // old fully encrypted confCons.xml
	case sshConfigLike(head):
		return "sshconfig"
	}
	switch ext {
	case ".rdp":
		return "rdp"
	case ".rdg":
		return "rdcman"
	case ".mxtsessions":
		return "mobaxterm"
	case ".remmina":
		return "remmina"
	case ".reg":
		return "putty"
	}
	return ""
}

var regexpConnectionsRoot = regexp.MustCompile(`^\s*(?:<\?xml[^>]*\?>\s*)?(?:<!--[\s\S]*?-->\s*)*<(?:[A-Za-z_][\w.-]*:)?Connections\b`)

func csvHasColumns(line string, cols ...string) bool {
	r := csv.NewReader(strings.NewReader(line))
	r.LazyQuotes = true
	rec, err := r.Read()
	if err != nil {
		return false
	}
	have := map[string]bool{}
	for _, c := range rec {
		have[strings.TrimSpace(c)] = true
	}
	for _, c := range cols {
		if !have[c] {
			return false
		}
	}
	return true
}

func sshConfigLike(text string) bool {
	sc := bufio.NewScanner(strings.NewReader(text))
	hosts := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, _ := splitSSHOption(line)
		switch strings.ToLower(k) {
		case "host":
			hosts++
		case "hostname", "user", "port", "identityfile", "proxyjump", "proxycommand", "localforward", "remoteforward",
			"dynamicforward", "match", "include", "serveraliveinterval", "forwardagent", "stricthostkeychecking", "userknownhostsfile":
		default:
			if hosts == 0 {
				return false
			}
		}
	}
	return hosts > 0
}

// importAuto runs the importer of a detected (or chosen) format.
func importAuto(im *importer, format string, raw []byte, password, fileName string) (*importResult, error) {
	text := importText(raw)
	var res *importResult
	var err error
	switch format {
	case "mremoteng":
		res, err = importMRemoteNGWith(im, []byte(text), password)
	case "mremoteng-csv":
		res, err = importMRemoteNGCSV(im, text)
	case "rdm-csv":
		res, err = importRDMCSV(im, text)
	case "rdcman":
		res, err = importRDCMan(im, text)
	case "rdp":
		res, err = importRDPFile(im, text, fileName)
	case "mobaxterm":
		res, err = importMobaXterm(im, text)
	case "winscp":
		res, err = importWinSCP(im, text)
	case "remmina":
		res, err = importRemmina(im, text)
	case "sshconfig":
		res, err = importSSHConfigWith(im, text, "")
	case "putty":
		res, err = importPuTTYWith(im, raw, "")
	default:
		return nil, fmt.Errorf("unknown file format: supported are mRemoteNG (XML or CSV), Remote Desktop Manager CSV, RDCMan (.rdg), .rdp, MobaXterm, WinSCP, Remmina, OpenSSH config, PuTTY (.reg) and WRM's own JSON export")
	}
	if res != nil {
		res.Format = format
		res.Dry = im.dry
	}
	return res, err
}

// hostWithPort adds a port to host unless it is the protocol's default or already present.
func hostWithPort(host, proto string, port int) string {
	host = strings.TrimSpace(host)
	if port <= 0 || port > 65535 || strings.Contains(strings.Trim(host, "[]"), ":") {
		return host
	}
	def := map[string]int{"SSH": 22, "SFTP": 22, "FTP": 21, "FTPS": 21, "HTTP": 80, "HTTPS": 443, "RDP": 3389, "VNC": 5900, "TELNET": 23}[proto]
	if port == def {
		return host
	}
	return net.JoinHostPort(strings.Trim(host, "[]"), strconv.Itoa(port))
}

// add appends a connection unless it has no host.
func (im *importer) add(c Connection, folder, jump string) *importItem {
	c.UserID = im.userID
	if c.AuthMethod == "" {
		c.AuthMethod = "PASSWORD"
	}
	if strings.TrimSpace(c.Name) == "" {
		c.Name = c.Host
	}
	if strings.TrimSpace(c.Host) == "" {
		im.skip(c.Name, "no host name")
		return nil
	}
	if c.Protocol == "SSH" && c.Username == "" {
		im.note(c.Name, "no user name — set it before connecting")
	}
	it := &importItem{c: c, folder: folder, jumpName: jump}
	im.items = append(im.items, it)
	return it
}

func (im *importer) finish() (*importResult, error) {
	if len(im.items) == 0 && len(im.res.Skipped) == 0 {
		return nil, fmt.Errorf("no connections found in the file")
	}
	if err := im.commit(); err != nil {
		return nil, err
	}
	return &im.res, nil
}

// splitUserDomain splits DOMAIN\user and user@domain (RDP) into user and domain.
func splitUserDomain(u string) (user, domain string) {
	u = strings.TrimSpace(u)
	if i := strings.Index(u, "\\"); i > 0 {
		return u[i+1:], u[:i]
	}
	return u, ""
}

// ── mRemoteNG CSV ──

func readCSVRecords(text string, sep rune) ([]map[string]string, error) {
	r := csv.NewReader(strings.NewReader(text))
	r.Comma, r.FieldsPerRecord, r.LazyQuotes = sep, -1, true
	recs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read the CSV file: %v", err)
	}
	if len(recs) < 2 {
		return nil, fmt.Errorf("the CSV file has no rows")
	}
	head := recs[0]
	var out []map[string]string
	for _, rec := range recs[1:] {
		if len(rec) == 1 && strings.TrimSpace(rec[0]) == "" {
			continue
		}
		m := map[string]string{}
		for i, h := range head {
			if i < len(rec) {
				m[strings.TrimSpace(h)] = rec[i]
			}
		}
		out = append(out, m)
	}
	return out, nil
}

// mrProtocol maps an mRemoteNG protocol name to WRM's; "" = not supported.
func mrProtocol(p string) string {
	switch strings.ToUpper(strings.TrimSpace(p)) {
	case "SSH1", "SSH2", "SSH":
		return "SSH"
	case "RDP":
		return "RDP"
	case "VNC":
		return "VNC"
	case "TELNET":
		return "TELNET"
	case "HTTP":
		return "HTTP"
	case "HTTPS":
		return "HTTPS"
	}
	return ""
}

func importMRemoteNGCSV(im *importer, text string) (*importResult, error) {
	rows, err := readCSVRecords(text, ';')
	if err != nil {
		return nil, err
	}
	type node struct{ name, parent string }
	containers := map[string]node{}
	for _, r := range rows {
		if strings.EqualFold(r["NodeType"], "Container") {
			containers[r["Id"]] = node{strings.TrimSpace(r["Name"]), r["Parent"]}
		}
	}
	pathOf := func(id string) string {
		var parts []string
		for seen := 0; id != "" && seen < 64; seen++ {
			n, ok := containers[id]
			if !ok {
				break
			}
			parts = append([]string{n.name}, parts...)
			id = n.parent
		}
		return strings.Join(parts, " / ")
	}
	for _, r := range rows {
		if !strings.EqualFold(r["NodeType"], "Connection") {
			continue
		}
		name := strings.TrimSpace(r["Name"])
		proto := mrProtocol(r["Protocol"])
		if proto == "" {
			im.skip(name, fmt.Sprintf("protocol %s is not supported by WRM (SSH, SFTP, RDP, VNC, Telnet, HTTP/HTTPS)", r["Protocol"]))
			continue
		}
		host := strings.TrimSpace(r["Hostname"])
		c := Connection{Name: name, Protocol: proto, Username: strings.TrimSpace(r["Username"]), Password: r["Password"]}
		if proto == "HTTP" || proto == "HTTPS" {
			host, c.Protocol, c.WebPath = splitWebURL(host, proto)
		}
		port, _ := strconv.Atoi(r["Port"])
		c.Host = hostWithPort(host, c.Protocol, port)
		if proto == "RDP" {
			c.Options = map[string]string{"ignore_cert": "true"}
			if d := strings.TrimSpace(r["Domain"]); d != "" {
				c.Options["domain"] = d
			}
		}
		im.add(c, pathOf(r["Parent"]), strings.TrimSpace(r["SSHTunnelConnectionName"]))
	}
	return im.finish()
}

// splitWebURL takes "https://host:port/path" apart.
func splitWebURL(raw, proto string) (host, protocol, webPath string) {
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "://") {
		if i := strings.Index(raw, "/"); i > 0 {
			return raw[:i], proto, raw[i:]
		}
		return raw, proto, ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return raw, proto, ""
	}
	webPath = u.EscapedPath()
	if u.RawQuery != "" {
		webPath += "?" + u.RawQuery
	}
	p := strings.ToUpper(u.Scheme)
	if p != "HTTP" && p != "HTTPS" {
		p = proto
	}
	return u.Host, p, webPath
}

// ── Remote Desktop Manager CSV ──

func importRDMCSV(im *importer, text string) (*importResult, error) {
	rows, err := readCSVRecords(text, ',')
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		name := strings.TrimSpace(r["Name"])
		typ := strings.ToLower(r["ConnectionType"])
		var proto string
		switch {
		case strings.Contains(typ, "rdp") || strings.Contains(typ, "remote desktop"):
			proto = "RDP"
		case strings.Contains(typ, "sftp"):
			proto = "SFTP"
		case strings.Contains(typ, "ssh"):
			proto = "SSH"
		case strings.Contains(typ, "vnc"):
			proto = "VNC"
		case strings.Contains(typ, "telnet"):
			proto = "TELNET"
		case strings.Contains(typ, "ftp"):
			proto = "FTP"
		case strings.Contains(typ, "web"):
			proto = "HTTPS"
		case typ == "group" || typ == "folder":
			continue
		default:
			im.skip(name, "connection type "+r["ConnectionType"]+" is not supported by WRM")
			continue
		}
		user, domain := splitUserDomain(r["CredentialUserName"])
		if d := strings.TrimSpace(r["CredentialDomain"]); d != "" {
			domain = d
		}
		c := Connection{Name: name, Protocol: proto, Username: user, Password: r["CredentialPassword"]}
		host := strings.TrimSpace(r["Host"])
		if proto == "HTTPS" {
			host, c.Protocol, c.WebPath = splitWebURL(host, proto)
		}
		port, _ := strconv.Atoi(r["Port"])
		c.Host = hostWithPort(host, c.Protocol, port)
		if proto == "RDP" {
			c.Options = map[string]string{"ignore_cert": "true"}
			if domain != "" {
				c.Options["domain"] = domain
			}
		}
		folder := strings.Join(strings.FieldsFunc(r["Group"], func(r rune) bool { return r == '\\' || r == '/' }), " / ")
		im.add(c, folder, "")
	}
	return im.finish()
}

// ── RDCMan (.rdg) ──

type rdgProps struct {
	Name        string `xml:"name"`
	DisplayName string `xml:"displayName"`
}

type rdgCreds struct {
	Inherit  string `xml:"inherit,attr"`
	User     string `xml:"userName"`
	Domain   string `xml:"domain"`
	Password string `xml:"password"`
}

type rdgConnSettings struct {
	Inherit string `xml:"inherit,attr"`
	Port    string `xml:"port"`
}

type rdgGateway struct {
	Inherit string `xml:"inherit,attr"`
	Enabled string `xml:"enabled"`
	Host    string `xml:"hostName"`
}

type rdgServer struct {
	Props rdgProps `xml:"properties"`
	// schema 1 (RDCMan 2.2) keeps the fields directly in <server>
	Name        string          `xml:"name"`
	DisplayName string          `xml:"displayName"`
	Creds       rdgCreds        `xml:"logonCredentials"`
	Conn        rdgConnSettings `xml:"connectionSettings"`
	Gateway     rdgGateway      `xml:"gatewaySettings"`
}

type rdgGroup struct {
	Props   rdgProps        `xml:"properties"`
	Name    string          `xml:"name"`
	Creds   rdgCreds        `xml:"logonCredentials"`
	Conn    rdgConnSettings `xml:"connectionSettings"`
	Servers []rdgServer     `xml:"server"`
	Groups  []rdgGroup      `xml:"group"`
}

type rdgFile struct {
	XMLName xml.Name
	File    rdgGroup `xml:"file"`
}

func importRDCMan(im *importer, text string) (*importResult, error) {
	var f rdgFile
	if err := mrUnmarshal([]byte(text), &f); err != nil {
		return nil, fmt.Errorf("not an RDCMan file (.rdg): %v", err)
	}
	if !strings.EqualFold(f.XMLName.Local, "RDCMan") {
		return nil, fmt.Errorf("not an RDCMan file (root element %q)", f.XMLName.Local)
	}
	encrypted := false
	var walk func(g rdgGroup, path []string, creds rdgCreds, port string)
	walk = func(g rdgGroup, path []string, creds rdgCreds, port string) {
		if g.Creds.Inherit == "None" {
			creds = g.Creds
		}
		if g.Conn.Inherit == "None" && g.Conn.Port != "" {
			port = g.Conn.Port
		}
		for _, s := range g.Servers {
			host := firstNonEmpty(s.Props.Name, s.Name)
			name := firstNonEmpty(s.Props.DisplayName, s.DisplayName, host)
			cr := creds
			if s.Creds.Inherit == "None" {
				cr = s.Creds
			}
			p := port
			if s.Conn.Inherit == "None" && s.Conn.Port != "" {
				p = s.Conn.Port
			}
			pn, _ := strconv.Atoi(p)
			user, domain := splitUserDomain(cr.User)
			if cr.Domain != "" {
				domain = cr.Domain
			}
			c := Connection{Name: name, Protocol: "RDP", Host: hostWithPort(host, "RDP", pn), Username: user,
				Options: map[string]string{"ignore_cert": "true"}}
			if domain != "" {
				c.Options["domain"] = domain
			}
			if s.Gateway.Inherit == "None" && strings.EqualFold(s.Gateway.Enabled, "true") && s.Gateway.Host != "" {
				c.Options["gateway_host"] = s.Gateway.Host
			}
			if cr.Password != "" {
				encrypted = true
			}
			im.add(c, strings.Join(path, " / "), "")
		}
		for _, sub := range g.Groups {
			walk(sub, append(append([]string{}, path...), firstNonEmpty(sub.Props.Name, sub.Name)), creds, port)
		}
	}
	walk(f.File, nil, rdgCreds{}, "")
	if encrypted {
		im.note("RDCMan", "passwords are encrypted with Windows DPAPI for one Windows user and cannot be read elsewhere — set them in WRM or use a vault credential")
	}
	return im.finish()
}

// ── .rdp ──

func importRDPFile(im *importer, text, name string) (*importResult, error) {
	vals := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), ":", 3)
		if len(parts) == 3 {
			vals[strings.ToLower(parts[0])] = strings.TrimSpace(parts[2])
		}
	}
	host := vals["full address"]
	if host == "" {
		return nil, fmt.Errorf("not a Remote Desktop file: no \"full address\"")
	}
	user, domain := splitUserDomain(vals["username"])
	if d := vals["domain"]; d != "" {
		domain = d
	}
	if port, _ := strconv.Atoi(vals["server port"]); port > 0 {
		host = hostWithPort(host, "RDP", port)
	}
	title := strings.TrimSuffix(path.Base(strings.ReplaceAll(name, "\\", "/")), path.Ext(name))
	if title == "" || title == "." || title == "/" {
		title = host
	}
	c := Connection{Name: title, Protocol: "RDP", Host: host, Username: user, Options: map[string]string{"ignore_cert": "true"}}
	if domain != "" {
		c.Options["domain"] = domain
	}
	if gw := vals["gatewayhostname"]; gw != "" && vals["gatewayusagemethod"] != "0" {
		c.Options["gateway_host"] = gw
	}
	if vals["administrative session"] == "1" || vals["connect to console"] == "1" {
		c.Options["console"] = "true"
	}
	if vals["password 51"] != "" {
		im.note(title, "the saved password is encrypted with Windows DPAPI and was not imported")
	}
	im.add(c, "", "")
	return im.finish()
}

// ── MobaXterm ──

// mobaTypes maps MobaXterm's session type number (second field) to WRM protocols.
var mobaTypes = map[string]string{"0": "SSH", "1": "TELNET", "4": "RDP", "5": "VNC", "6": "FTP", "7": "SFTP", "11": "HTTPS"}
var mobaTypeNames = map[string]string{"2": "Rsh", "3": "XDMCP", "8": "Serial", "9": "File", "10": "Shell", "12": "Mosh", "13": "S3", "14": "WSL"}

func iniSections(text string) []struct {
	name string
	kv   [][2]string
} {
	var out []struct {
		name string
		kv   [][2]string
	}
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimRight(strings.TrimLeft(raw, " \t"), "\r")
		if strings.HasPrefix(line, "[") && strings.HasSuffix(strings.TrimSpace(line), "]") {
			l := strings.TrimSpace(line)
			out = append(out, struct {
				name string
				kv   [][2]string
			}{name: l[1 : len(l)-1]})
			continue
		}
		if len(out) == 0 || line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(line, "="); i > 0 {
			cur := &out[len(out)-1]
			cur.kv = append(cur.kv, [2]string{strings.TrimSpace(line[:i]), line[i+1:]})
		}
	}
	return out
}

func importMobaXterm(im *importer, text string) (*importResult, error) {
	found := false
	for _, sec := range iniSections(text) {
		if !strings.HasPrefix(strings.ToLower(sec.name), "bookmarks") {
			continue
		}
		found = true
		folder := ""
		for _, kv := range sec.kv {
			if strings.EqualFold(kv[0], "SubRep") {
				folder = strings.Join(strings.FieldsFunc(kv[1], func(r rune) bool { return r == '\\' }), " / ")
			}
		}
		for _, kv := range sec.kv {
			name, v := kv[0], strings.TrimSpace(kv[1])
			if strings.EqualFold(name, "SubRep") || strings.EqualFold(name, "ImgNum") || !strings.HasPrefix(v, "#") {
				continue
			}
			// "#<icon>#<type>%host%port%user%…#<terminal settings>#…"
			parts := strings.SplitN(v, "#", 4)
			if len(parts) < 3 {
				im.skip(name, "session line not understood")
				continue
			}
			f := strings.Split(parts[2], "%")
			typ := f[0]
			proto, ok := mobaTypes[typ]
			if !ok {
				im.skip(name, fmt.Sprintf("MobaXterm %s session is not supported by WRM", firstNonEmpty(mobaTypeNames[typ], "type "+typ)))
				continue
			}
			field := func(i int) string {
				if i < len(f) {
					return strings.TrimSpace(f[i])
				}
				return ""
			}
			c := Connection{Name: name, Protocol: proto}
			host := field(1)
			if proto == "HTTPS" {
				host, c.Protocol, c.WebPath = splitWebURL(host, proto)
			} else {
				port, _ := strconv.Atoi(field(2))
				host = hostWithPort(host, proto, port)
				if u := field(3); u != "-1" && proto != "VNC" {
					c.Username = u
				}
			}
			c.Host = host
			if proto == "RDP" {
				c.Options = map[string]string{"ignore_cert": "true"}
			}
			im.add(c, folder, "")
		}
	}
	if !found {
		return nil, fmt.Errorf("not a MobaXterm sessions file: no [Bookmarks] section")
	}
	im.note("MobaXterm", "passwords are kept by MobaXterm in its own encrypted store and are not part of the export — set them in WRM or use a vault credential")
	return im.finish()
}

// ── WinSCP ──

// winscpDecrypt reverses WinSCP's stored password obfuscation (no master password).
func winscpDecrypt(host, user, enc string) (string, bool) {
	const magic = 0xA3
	enc = strings.ToUpper(strings.TrimSpace(enc))
	next := func() (int, bool) {
		if len(enc) < 2 {
			return 0, false
		}
		a := strings.IndexByte("0123456789ABCDEF", enc[0])
		b := strings.IndexByte("0123456789ABCDEF", enc[1])
		if a < 0 || b < 0 {
			return 0, false
		}
		enc = enc[2:]
		return ^((a<<4 + b) ^ magic) & 0xff, true
	}
	flag, ok := next()
	if !ok {
		return "", false
	}
	length := flag
	if flag == 0xff {
		if _, ok = next(); !ok {
			return "", false
		}
		if length, ok = next(); !ok {
			return "", false
		}
	}
	skip, ok := next()
	if !ok || skip*2 > len(enc) {
		return "", false
	}
	enc = enc[skip*2:]
	out := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		c, ok := next()
		if !ok {
			return "", false
		}
		out = append(out, byte(c))
	}
	s := string(out)
	if flag == 0xff {
		key := user + host
		if !strings.HasPrefix(s, key) {
			return "", false
		}
		s = s[len(key):]
	}
	return s, true
}

func importWinSCP(im *importer, text string) (*importResult, error) {
	found, master := false, false
	for _, sec := range iniSections(text) {
		if strings.EqualFold(sec.name, "Configuration\\Security") {
			for _, kv := range sec.kv {
				if strings.EqualFold(kv[0], "UseMasterPassword") && kv[1] == "1" {
					master = true
				}
			}
		}
	}
	for _, sec := range iniSections(text) {
		if !strings.HasPrefix(strings.ToLower(sec.name), "sessions\\") {
			continue
		}
		raw, _ := url.PathUnescape(sec.name[len("Sessions\\"):])
		if strings.EqualFold(raw, "Default Settings") {
			continue
		}
		found = true
		v := map[string]string{}
		for _, kv := range sec.kv {
			v[strings.ToLower(kv[0])] = kv[1]
		}
		parts := strings.Split(raw, "/")
		name := parts[len(parts)-1]
		folder := strings.Join(parts[:len(parts)-1], " / ")
		host, user := strings.TrimSpace(v["hostname"]), strings.TrimSpace(v["username"])
		host, _ = url.PathUnescape(host)
		user, _ = url.PathUnescape(user)
		proto := "SFTP"
		switch v["fsprotocol"] {
		case "5":
			proto = "FTP"
			if v["ftps"] != "" && v["ftps"] != "0" {
				proto = "FTPS"
			}
		case "6":
			im.skip(name, "WebDAV sessions are not supported by WRM")
			continue
		case "7":
			im.skip(name, "S3 sessions: add the bucket under Storages instead")
			continue
		}
		port, _ := strconv.Atoi(v["portnumber"])
		c := Connection{Name: name, Protocol: proto, Host: hostWithPort(host, proto, port), Username: user}
		if p := v["password"]; p != "" {
			if master {
				im.note(name, "password is protected by the WinSCP master password and was not imported")
			} else if pw, ok := winscpDecrypt(host, user, p); ok {
				c.Password = pw
			} else {
				im.note(name, "stored password could not be read")
			}
		}
		if k := strings.TrimSpace(v["publickeyfile"]); k != "" {
			im.note(name, "uses the key file "+k+" — add the key to the WRM vault")
		}
		im.add(c, folder, "")
	}
	if !found {
		return nil, fmt.Errorf("no WinSCP sessions found (export with Tools → Export/Backup Configuration)")
	}
	return im.finish()
}

// ── Remmina ──

func importRemmina(im *importer, text string) (*importResult, error) {
	for _, sec := range iniSections(text) {
		if !strings.EqualFold(sec.name, "remmina") {
			continue
		}
		v := map[string]string{}
		for _, kv := range sec.kv {
			v[strings.ToLower(kv[0])] = strings.TrimSpace(kv[1])
		}
		name := v["name"]
		proto := map[string]string{"RDP": "RDP", "VNC": "VNC", "SSH": "SSH", "SFTP": "SFTP", "TELNET": "TELNET"}[strings.ToUpper(v["protocol"])]
		if proto == "" {
			im.skip(name, "Remmina protocol "+v["protocol"]+" is not supported by WRM")
			continue
		}
		user := v["username"]
		if proto == "SSH" || proto == "SFTP" {
			user = firstNonEmpty(v["ssh_username"], user)
		}
		c := Connection{Name: name, Protocol: proto, Host: firstNonEmpty(v["server"], v["ssh_server"]), Username: user}
		if proto == "RDP" {
			c.Options = map[string]string{"ignore_cert": "true"}
			if d := v["domain"]; d != "" {
				c.Options["domain"] = d
			}
		}
		if v["password"] != "" {
			im.note(name, "the password is encrypted with this Remmina installation's secret and was not imported")
		}
		im.add(c, strings.Join(strings.FieldsFunc(v["group"], func(r rune) bool { return r == '/' }), " / "), "")
	}
	return im.finish()
}
