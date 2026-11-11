package main

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/twofish"
)

func autoImport(t *testing.T, u *testClient, name, content string, extra map[string]interface{}) (int, importResult, map[string]interface{}) {
	t.Helper()
	body := map[string]interface{}{"name": name, "data": base64.StdEncoding.EncodeToString([]byte(content))}
	for k, v := range extra {
		body[k] = v
	}
	var raw map[string]interface{}
	code := u.jsonDo("POST", "/api/config/import/auto", body, &raw)
	var res importResult
	if code == 200 {
		b := jsonMarshal(raw)
		if err := json.Unmarshal(b, &res); err != nil {
			t.Fatal(err)
		}
	}
	return code, res, raw
}

func TestImportDetectFormats(t *testing.T) {
	cases := map[string][2]string{
		"mremoteng":     {"confCons.xml", `<?xml version="1.0" encoding="utf-16"?><mrng:Connections xmlns:mrng="http://mremoteng.org" Name="Connections"></mrng:Connections>`},
		"mremoteng-csv": {"export.csv", "Name;Id;Parent;NodeType;Description\n"},
		"rdm-csv":       {"rdm.csv", "Name,Host,ConnectionType,Group\n"},
		"rdcman":        {"servers.rdg", `<?xml version="1.0" encoding="utf-8"?><RDCMan programVersion="2.7" schemaVersion="3"></RDCMan>`},
		"rdp":           {"x.rdp", "screen mode id:i:2\r\nfull address:s:demo-rdp.example.com\r\n"},
		"mobaxterm":     {"s.mxtsessions", "[Bookmarks]\r\nSubRep=\r\nImgNum=42\r\n"},
		"winscp":        {"WinSCP.ini", "[Sessions\\demo]\r\nHostName=demo.example.com\r\n"},
		"remmina":       {"a.remmina", "[remmina]\nname=demo\n"},
		"sshconfig":     {"config", "# mine\nHost demo\n  HostName demo.example.com\n  User ops\n"},
		"putty":         {"putty.reg", "Windows Registry Editor Version 5.00\r\n"},
	}
	for want, c := range cases {
		if got := detectImportFormat(c[0], c[1]); got != want {
			t.Errorf("%s: detected %q, want %q", c[0], got, want)
		}
	}
	if got := detectImportFormat("notes.txt", "hello world"); got != "" {
		t.Errorf("plain text detected as %q", got)
	}
	// UTF-16 LE without BOM (regedit, some Windows tools) is decoded.
	u16 := []byte{}
	for _, r := range "[remmina]\nname=x\n" {
		u16 = append(u16, byte(r), 0)
	}
	if got := importText(u16); !strings.HasPrefix(got, "[remmina]") {
		t.Errorf("UTF-16 without BOM: %q", got)
	}
}

// mrEncryptTwofish is mRemoteNG's AEAD format with the Twofish engine.
func mrEncryptTwofish(t *testing.T, plain, password string) string {
	salt, nonce := make([]byte, 16), make([]byte, 16)
	rand.Read(salt)
	rand.Read(nonce)
	key := pbkdf2.Key([]byte(password), salt, 1000, 32, sha1.New)
	block, err := twofish.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	aead, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		t.Fatal(err)
	}
	out := append(append(append([]byte{}, salt...), nonce...), aead.Seal(nil, nonce, []byte(plain), salt)...)
	return base64.StdEncoding.EncodeToString(out)
}

// winscpEncrypt is WinSCP's stored password obfuscation (the inverse of winscpDecrypt).
func winscpEncrypt(host, user, pw string) string {
	const magic = 0xA3
	var sb strings.Builder
	put := func(b int) { sb.WriteString(fmt.Sprintf("%02X", (^b^magic)&0xff)) }
	s := user + host + pw
	put(0xff)
	put(0)
	put(len(s))
	put(0) // nothing to skip
	for i := 0; i < len(s); i++ {
		put(int(s[i]))
	}
	return sb.String()
}

func TestImportMoreFormats(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()

	t.Run("preview writes nothing", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-preview", false)
		code, res, raw := autoImport(t, u, "confCons.xml", mrFile(t, "mR3m"), map[string]interface{}{"preview": true, "folder": "From mRemoteNG"})
		if code != 200 || !res.Dry || res.Format != "mremoteng" || res.Imported != 5 || len(res.Preview) != 5 || res.Jumps != 2 {
			t.Fatalf("preview: %d %+v", code, raw)
		}
		for _, p := range res.Preview {
			if !strings.HasPrefix(p.Folder, "From mRemoteNG") {
				t.Fatalf("preview folder: %+v", p)
			}
		}
		if strings.Contains(string(jsonMarshal(raw)), "secret1") {
			t.Fatal("preview leaks a password")
		}
		if n := len(loadUserConnections(u.userID)); n != 0 || len(loadFolders(u.userID)) != 0 {
			t.Fatalf("preview wrote %d connections", n)
		}
		code, res, _ = autoImport(t, u, "confCons.xml", mrFile(t, "mR3m"), map[string]interface{}{"folder": "From mRemoteNG"})
		if code != 200 || res.Imported != 5 {
			t.Fatalf("import: %d %+v", code, res)
		}
		if p := folderPathOf(u.userID, connByName(t, u.userID, "web01").FolderID); p != "From mRemoteNG / DC1 / Rack A" {
			t.Fatalf("base folder: %q", p)
		}
	})

	t.Run("mRemoteNG old full-file encryption and Twofish", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-mr-old", false)
		inner := `<?xml version="1.0" encoding="utf-8"?><Connections Name="Connections" Export="False" Protected="` + mrEncryptCBC("ThisIsNotProtected", "mR3m") + `" ConfVersion="2.4">` +
			`<Node Name="old01" Type="Connection" Protocol="SSH2" Hostname="old01.example.com" Port="22" Username="ops" Password="` + mrEncryptCBC("old-pw", "mR3m") + `" /></Connections>`
		file := mrEncryptCBC(inner, "mR3m")
		code, res, raw := autoImport(t, u, "confCons.xml", file, nil)
		if code != 200 || res.Imported != 1 {
			t.Fatalf("old encrypted file: %d %v", code, raw)
		}
		if c := connByName(t, u.userID, "old01"); c.Password != "old-pw" || c.Username != "ops" {
			t.Fatalf("old01: %+v", c)
		}
		// Protected with a master password: asked for, then decrypted.
		file = mrEncryptCBC(inner, "Old-Master")
		if code, _, raw = autoImport(t, u, "confCons.xml", file, nil); code != 400 || raw["need_password"] != true {
			t.Fatalf("old file without password: %d %v", code, raw)
		}

		tf := `<mrng:Connections xmlns:mrng="http://mremoteng.org" Name="Connections" EncryptionEngine="Twofish" BlockCipherMode="GCM" KdfIterations="1000" FullFileEncryption="false" Protected="` +
			mrEncryptTwofish(t, "ThisIsNotProtected", "mR3m") + `" ConfVersion="2.8"><Node Name="tf01" Type="Connection" Protocol="RDP" Hostname="tf01.example.com" Port="3389" Username="adm" Domain="CORP" Password="` +
			mrEncryptTwofish(t, "tf-pw", "mR3m") + `" /></mrng:Connections>`
		if code, res, raw = autoImport(t, u, "confCons.xml", tf, nil); code != 200 || res.Imported != 1 {
			t.Fatalf("twofish: %d %v", code, raw)
		}
		if c := connByName(t, u.userID, "tf01"); c.Password != "tf-pw" || loadDesktopOptions(c.ID)["domain"] != "CORP" {
			t.Fatalf("tf01: %+v", c)
		}
		bad := strings.Replace(tf, `"Twofish"`, `"Serpent"`, 1)
		if code, _, raw = autoImport(t, u, "confCons.xml", bad, nil); code != 400 || !strings.Contains(fmt.Sprint(raw["error"]), "Serpent") {
			t.Fatalf("serpent: %d %v", code, raw)
		}
	})

	t.Run("same server in two folders", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-dup", false)
		x := `<mrng:Connections xmlns:mrng="http://mremoteng.org" Name="Connections" EncryptionEngine="AES" BlockCipherMode="GCM" KdfIterations="1000" Protected="` + mrEncrypt(t, "ThisIsNotProtected", "mR3m") + `" ConfVersion="2.8">
<Node Name="Prod" Type="Container"><Node Name="db" Type="Connection" Protocol="SSH2" Hostname="db.example.com" Port="22" Username="u" /></Node>
<Node Name="Favourites" Type="Container"><Node Name="db" Type="Connection" Protocol="SSH2" Hostname="db.example.com" Port="22" Username="u" /></Node>
</mrng:Connections>`
		if code, res, raw := autoImport(t, u, "c.xml", x, nil); code != 200 || res.Imported != 2 {
			t.Fatalf("two folders: %d %v", code, raw)
		}
		if _, res, _ := autoImport(t, u, "c.xml", x, nil); res.Imported != 0 || len(res.Skipped) != 2 {
			t.Fatalf("re-import: %+v", res)
		}
	})

	t.Run("mRemoteNG CSV", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-mrcsv", false)
		csv := "Name;Id;Parent;NodeType;Description;Icon;Panel;Username;Password;Domain;Hostname;Port;Protocol;SSHTunnelConnectionName\r\n" +
			"web;c1;f2;Connection;;;General;ops;\"p;w\";;web.example.com;2222;SSH2;jump\r\n" +
			"jump;c2;f1;Connection;;;General;ops;;;jump.example.com;22;SSH2;\r\n" +
			"win;c3;f2;Connection;;;General;adm;pw3;CORP;win.example.com;3389;RDP;\r\n" +
			"ps;c4;f1;Connection;;;General;;;;ps.example.com;5985;PowerShell;\r\n" +
			"Rack;f2;f1;Container;;;General;;;;;;;\r\n" +
			"DC;f1;;Container;;;General;;;;;;;\r\n"
		code, res, raw := autoImport(t, u, "export.csv", csv, nil)
		if code != 200 || res.Format != "mremoteng-csv" || res.Imported != 3 || res.Jumps != 1 || len(res.Skipped) != 1 {
			t.Fatalf("mRemoteNG CSV: %d %v", code, raw)
		}
		w := connByName(t, u.userID, "web")
		if w.Host != "web.example.com:2222" || w.Password != "p;w" || folderPathOf(u.userID, w.FolderID) != "DC / Rack" || w.JumpID == nil {
			t.Fatalf("web: %+v %q", w, folderPathOf(u.userID, w.FolderID))
		}
	})

	t.Run("RDM CSV, RDCMan, rdp", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-win", false)
		rdm := "Name,Host,Port,ConnectionType,Group,CredentialUserName,CredentialDomain,CredentialPassword,Description\n" +
			"app01,app01.example.com,3390,RDP (Microsoft Remote Desktop),Site A\\Apps,adm,CORP,pw1,\n" +
			"lnx01,lnx01.example.com,22,SSH Shell,Site A,root,,pw2,\n"
		code, res, raw := autoImport(t, u, "rdm.csv", rdm, nil)
		if code != 200 || res.Format != "rdm-csv" || res.Imported != 2 {
			t.Fatalf("RDM: %d %v", code, raw)
		}
		if a := connByName(t, u.userID, "app01"); a.Host != "app01.example.com:3390" || a.Password != "pw1" || folderPathOf(u.userID, a.FolderID) != "Site A / Apps" {
			t.Fatalf("app01: %+v", a)
		}
		rdg := `<?xml version="1.0" encoding="utf-8"?>
<RDCMan programVersion="2.7" schemaVersion="3"><file><properties><name>Lab</name></properties>
<logonCredentials inherit="None"><profileName scope="Local">Custom</profileName><userName>CORP\admin</userName><password>AQAAANCMnd8B</password><domain></domain></logonCredentials>
<group><properties><name>Hyper-V</name></properties>
<server><properties><displayName>HV 1</displayName><name>hv1.example.com</name></properties><connectionSettings inherit="None"><port>3391</port></connectionSettings></server>
<group><properties><name>Inner</name></properties><server><properties><name>hv2.example.com</name></properties></server></group>
</group></file></RDCMan>`
		code, res, raw = autoImport(t, u, "lab.rdg", rdg, nil)
		if code != 200 || res.Imported != 2 || len(res.Notes) != 1 {
			t.Fatalf("RDCMan: %d %v", code, raw)
		}
		hv := connByName(t, u.userID, "HV 1")
		if hv.Host != "hv1.example.com:3391" || hv.Username != "admin" || loadDesktopOptions(hv.ID)["domain"] != "CORP" || folderPathOf(u.userID, hv.FolderID) != "Hyper-V" {
			t.Fatalf("HV 1: %+v", hv)
		}
		if p := folderPathOf(u.userID, connByName(t, u.userID, "hv2.example.com").FolderID); p != "Hyper-V / Inner" {
			t.Fatalf("hv2 folder %q", p)
		}
		rdp := "full address:s:jump.example.com:3392\r\nusername:s:CORP\\ops\r\ngatewayhostname:s:gw.example.com\r\ngatewayusagemethod:i:1\r\n"
		code, res, raw = autoImport(t, u, "C:\\Users\\me\\Desktop\\Jump box.rdp", rdp, nil)
		if code != 200 || res.Imported != 1 {
			t.Fatalf("rdp: %d %v", code, raw)
		}
		j := connByName(t, u.userID, "Jump box")
		if o := loadDesktopOptions(j.ID); j.Host != "jump.example.com:3392" || j.Username != "ops" || o["domain"] != "CORP" || o["gateway_host"] != "gw.example.com" {
			t.Fatalf("Jump box: %+v %v", j, o)
		}
	})

	t.Run("MobaXterm, WinSCP, Remmina", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-tools", false)
		moba := "[Bookmarks]\r\nSubRep=\r\nImgNum=42\r\n\r\n[Bookmarks_1]\r\nSubRep=Lab\\Linux\r\nImgNum=41\r\n" +
			"lnx=#109#0%lnx.example.com%2200%ops%%-1%-1%%%%%0%0%0%%%-1%0%0%0%%1080%%0%0%1#MobaFont%10%0%0%-1%15#0# #-1\r\n" +
			"win=#91#4%win.example.com%3389%adm%0%0%-1%0%-1%0#MobaFont%14#0# #-1\r\n" +
			"vnc=#128#5%vnc.example.com%5901%-1%0%%%%%-1%0#MobaFont%10#0# #-1\r\n" +
			"sh=#131#10%%%#MobaFont%10#0# #-1\r\n"
		code, res, raw := autoImport(t, u, "MobaXterm Sessions.mxtsessions", moba, nil)
		if code != 200 || res.Format != "mobaxterm" || res.Imported != 3 || len(res.Skipped) != 1 {
			t.Fatalf("MobaXterm: %d %v", code, raw)
		}
		if l := connByName(t, u.userID, "lnx"); l.Host != "lnx.example.com:2200" || l.Username != "ops" || folderPathOf(u.userID, l.FolderID) != "Lab / Linux" {
			t.Fatalf("lnx: %+v", l)
		}
		if v := connByName(t, u.userID, "vnc"); v.Protocol != "VNC" || v.Host != "vnc.example.com:5901" {
			t.Fatalf("vnc: %+v", v)
		}
		ini := "[Sessions\\Lab/files%20box]\r\nHostName=files.example.com\r\nUserName=ftpuser\r\nPortNumber=2121\r\nFSProtocol=5\r\nFtps=0\r\nPassword=" +
			winscpEncrypt("files.example.com", "ftpuser", "Ftp-Pw!") + "\r\n\r\n[Sessions\\sftp1]\r\nHostName=sftp1.example.com\r\nUserName=root\r\nPassword=" +
			winscpEncrypt("sftp1.example.com", "root", "x") + "\r\n\r\n[Sessions\\Default%20Settings]\r\nHostName=\r\n"
		code, res, raw = autoImport(t, u, "WinSCP.ini", ini, nil)
		if code != 200 || res.Format != "winscp" || res.Imported != 2 {
			t.Fatalf("WinSCP: %d %v", code, raw)
		}
		f := connByName(t, u.userID, "files box")
		if f.Protocol != "FTP" || f.Host != "files.example.com:2121" || f.Password != "Ftp-Pw!" || folderPathOf(u.userID, f.FolderID) != "Lab" {
			t.Fatalf("files box: %+v", f)
		}
		if s := connByName(t, u.userID, "sftp1"); s.Protocol != "SFTP" || s.Password != "x" {
			t.Fatalf("sftp1: %+v", s)
		}
		rem := "[remmina]\nname=Desk\nprotocol=RDP\nserver=desk.example.com:3390\nusername=me\ndomain=HOME\ngroup=Home/PCs\npassword=.\n"
		code, res, raw = autoImport(t, u, "desk.remmina", rem, nil)
		if code != 200 || res.Imported != 1 || len(res.Notes) != 1 {
			t.Fatalf("Remmina: %d %v", code, raw)
		}
		if d := connByName(t, u.userID, "Desk"); d.Host != "desk.example.com:3390" || folderPathOf(u.userID, d.FolderID) != "Home / PCs" {
			t.Fatalf("Desk: %+v", d)
		}
	})

	t.Run("unknown file", func(t *testing.T) {
		u := newTestUser(t, srv, "imp-unknown", false)
		if code, _, raw := autoImport(t, u, "notes.txt", "just some text", nil); code != 400 || !strings.Contains(fmt.Sprint(raw["error"]), "MobaXterm") {
			t.Fatalf("unknown: %d %v", code, raw)
		}
	})
}

func TestWinSCPDecrypt(t *testing.T) {
	if pw, ok := winscpDecrypt("h.example.com", "u", winscpEncrypt("h.example.com", "u", "pässwörd")); !ok || pw != "pässwörd" {
		t.Fatalf("roundtrip: %q %v", pw, ok)
	}
	if _, ok := winscpDecrypt("h", "u", "ZZ"); ok {
		t.Fatal("garbage accepted")
	}
}
