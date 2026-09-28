package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// ─── SSH WEBSOCKET ────────────────────────────────────
//
// Protocol (v9.10):
//   server → client  binary frame : raw terminal output (never split-UTF-8 sensitive)
//                    text frame   : JSON control {"type":"status"|"error"|"exit", ...}
//   client → server  binary frame : keyboard input
//                    text frame   : JSON control {"type":"resize"|"pause"|"resume"|"ping"}
//                                   (non-JSON text is treated as keyboard input for
//                                    backwards compatibility)
//
// Close codes: 1000 = shell exited normally, 4001 = could not connect,
//              4002 = SSH connection to the server was lost.

const (
	wsCloseConnectFailed = 4001
	wsCloseSSHLost       = 4002

	wsReadTimeout    = 90 * time.Second
	wsPingInterval   = 25 * time.Second
	sshKeepaliveTick = 30 * time.Second
	flowAutoResume   = 5 * time.Second
)

type termControl struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

func sshHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	ws.SetReadLimit(4 << 20)

	var wsMu sync.Mutex
	writeWS := func(mt int, data []byte) error {
		wsMu.Lock()
		defer wsMu.Unlock()
		ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
		return ws.WriteMessage(mt, data)
	}
	printTerm := func(s string) { writeWS(websocket.BinaryMessage, []byte(s)) }
	sendCtl := func(v interface{}) { writeWS(websocket.TextMessage, jsonMarshal(v)) }
	closeWS := func(code int, reason string) {
		wsMu.Lock()
		defer wsMu.Unlock()
		ws.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, truncateStr(reason, 120)), time.Now().Add(2*time.Second))
	}
	fail := func(msg string) {
		printTerm("\r\n\x1b[31m" + msg + "\x1b[0m\r\n")
		sendCtl(map[string]interface{}{"type": "error", "message": msg})
		closeWS(wsCloseConnectFailed, msg)
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	acc, _, denyMsg := authorizeConnection(r, id, PermTerminal)
	if acc == nil {
		fail(denyMsg)
		return
	}
	c, err := loadConnection(id)
	if err != nil {
		fail("Connection not found.")
		return
	}
	if isFTP(c) {
		fail("FTP connections do not support a terminal. Use the File Manager instead.")
		return
	}
	am, err := buildAuthMethods(c)
	if err != nil {
		fail("Auth error: " + err.Error())
		return
	}

	sendCtl(map[string]interface{}{"type": "status", "state": "connecting", "host": c.Host})
	printTerm("\r\n\x1b[36mConnecting to " + c.Host + "...\x1b[0m\r\n")
	cfg := sshClientConfig(c, am, func(fp string) {
		printTerm("\x1b[33mNew host " + c.Host + " — key " + fp + " saved (trust on first use).\x1b[0m\r\n")
	})
	sshClient, err := ssh.Dial("tcp", c.Host, cfg)
	if err != nil {
		if hk := asHostKeyError(err); hk != nil {
			sendCtl(map[string]interface{}{"type": "hostkey", "kind": hk.Kind, "host": hk.Host, "key_type": hk.KeyType,
				"expected": hk.Expected, "fingerprint": hk.Got, "conn_id": c.ID, "can_accept": acc.Share == nil})
			printTerm("\r\n\x1b[41;97m " + map[string]string{"mismatch": "WARNING: HOST KEY CHANGED", "unknown": "UNKNOWN HOST"}[hk.Kind] + " \x1b[0m\r\n")
			fail(hk.Error())
			return
		}
		fail("Connection failed: " + err.Error())
		return
	}
	defer sshClient.Close()

	// Register the terminal so administrators can see/end it and access can be revoked.
	actorID, actorName := acc.actor()
	ts := &termSession{ID: randomToken(9), UserID: actorID, User: actorName, ConnID: c.ID, ConnName: c.Name, Host: c.Host, IP: clientIP(r), Started: time.Now()}
	if acc.Share != nil {
		ts.ShareID, ts.Share, ts.PKey, ts.Ctx = acc.Share.Share.ID, acc.Share.Share.Name, acc.Share.PKey, acc.Share.Ctx
	}
	var killOnce sync.Once
	ts.kill = func(reason string) {
		killOnce.Do(func() {
			printTerm("\r\n\x1b[31m" + reason + "\x1b[0m\r\n")
			sendCtl(map[string]interface{}{"type": "error", "message": reason})
			wsMu.Lock()
			closeWSRevoked(ws, reason)
			wsMu.Unlock()
			sshClient.Close()
			ws.Close()
		})
	}
	registerTerminal(ts)
	defer unregisterTerminal(ts.ID)
	auditLogAs(r, actorID, actorName, "terminal.open", c.Name, map[string]interface{}{"host": c.Host, "user": c.Username, "share": ts.Share})
	defer func() {
		auditLogAs(r, actorID, actorName, "terminal.close", c.Name, map[string]interface{}{"host": c.Host, "minutes": int(time.Since(ts.Started).Minutes())})
	}()

	session, err := sshClient.NewSession()
	if err != nil {
		fail("Session error: " + err.Error())
		return
	}
	defer session.Close()

	// The browser tells us its real size in the URL, so the PTY is correct from the very
	// first byte (nano/vim/less/htop draw correctly without waiting for a resize).
	cols, _ := strconv.Atoi(r.URL.Query().Get("cols"))
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	haveSize := cols > 0 && rows > 0
	if !haveSize {
		cols, rows = 80, 24
	}
	modes := ssh.TerminalModes{
		ssh.ECHO: 1, ssh.ICANON: 1, ssh.ISIG: 1, ssh.IEXTEN: 1,
		ssh.ICRNL: 1, ssh.IMAXBEL: 1, ssh.IXON: 1, ssh.IXANY: 1, ssh.IUTF8: 1,
		ssh.OPOST: 1, ssh.ONLCR: 1, ssh.OCRNL: 0, ssh.ONLRET: 0,
		ssh.CS8: 1, ssh.PARENB: 0,
		ssh.VINTR: 3, ssh.VQUIT: 28, ssh.VERASE: 127, ssh.VKILL: 21, ssh.VEOF: 4,
		ssh.VSTART: 17, ssh.VSTOP: 19, ssh.VSUSP: 26, ssh.VREPRINT: 18,
		ssh.VWERASE: 23, ssh.VLNEXT: 22, ssh.VDISCARD: 15,
		ssh.TTY_OP_ISPEED: 115200, ssh.TTY_OP_OSPEED: 115200,
	}
	term := r.URL.Query().Get("term")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, rows, cols, modes); err != nil {
		fail("PTY error: " + err.Error())
		return
	}
	session.Setenv("LANG", "C.UTF-8") // best effort, usually refused by AcceptEnv

	sshIn, _ := session.StdinPipe()
	sshOut, _ := session.StdoutPipe()
	sshErr, _ := session.StderrPipe()

	done := make(chan struct{})
	var doneOnce sync.Once
	finish := func() { doneOnce.Do(func() { close(done) }) }

	// ── Flow control: the browser asks us to pause when xterm.js falls behind ──
	var paused atomic.Bool
	resumeCh := make(chan struct{}, 1)
	resume := func() {
		paused.Store(false)
		select {
		case resumeCh <- struct{}{}:
		default:
		}
	}

	// ── WebSocket reader ──
	stdinCh := make(chan []byte, 1024)
	gotSize := make(chan struct{})
	var sizeOnce sync.Once
	readDone := make(chan struct{})
	ws.SetReadDeadline(time.Now().Add(wsReadTimeout))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(wsReadTimeout))
		return nil
	})
	go func() {
		defer close(readDone)
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			ws.SetReadDeadline(time.Now().Add(wsReadTimeout))
			if mt == websocket.TextMessage && len(msg) > 1 && msg[0] == '{' {
				var ctl termControl
				if json.Unmarshal(msg, &ctl) == nil && ctl.Type != "" {
					switch ctl.Type {
					case "resize":
						if ctl.Cols > 0 && ctl.Rows > 0 && ctl.Cols < 2000 && ctl.Rows < 1000 {
							session.WindowChange(ctl.Rows, ctl.Cols)
							sizeOnce.Do(func() { close(gotSize) })
						}
					case "pause":
						paused.Store(true)
					case "resume":
						resume()
					case "ping":
						sendCtl(map[string]string{"type": "pong"})
					}
					continue
				}
			}
			select {
			case stdinCh <- msg:
			case <-done:
				return
			}
		}
	}()

	// ── stdin writer (decoupled so a slow remote never blocks control messages) ──
	go func() {
		for {
			select {
			case msg := <-stdinCh:
				if _, err := sshIn.Write(msg); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Old clients do not send the size in the URL: wait briefly for the first resize.
	if !haveSize {
		select {
		case <-gotSize:
		case <-time.After(2 * time.Second):
		case <-readDone:
			return
		}
	}

	if err := session.Shell(); err != nil {
		fail("Shell error: " + err.Error())
		return
	}
	log.Printf("SSH connected: %s@%s", c.Username, c.Host)
	sendCtl(map[string]interface{}{"type": "status", "state": "connected", "host": c.Host})

	// ── Output pump ──
	outDone := make(chan struct{})
	pump := func(src interface{ Read([]byte) (int, error) }, signal chan struct{}) {
		if signal != nil {
			defer close(signal)
		}
		buf := make([]byte, 32*1024)
		for {
			if signal != nil && paused.Load() { // only the main stdout pump is flow-controlled
				select {
				case <-resumeCh:
				case <-time.After(flowAutoResume):
					paused.Store(false)
				case <-done:
					return
				}
				continue
			}
			n, err := src.Read(buf)
			if n > 0 {
				if writeWS(websocket.BinaryMessage, buf[:n]) != nil {
					finish()
					return
				}
			}
			if err != nil {
				return
			}
		}
	}
	go pump(sshOut, outDone)
	go pump(sshErr, nil)

	// ── WebSocket ping + SSH keepalive ──
	var sshLost atomic.Bool
	go func() {
		pingT := time.NewTicker(wsPingInterval)
		kaT := time.NewTicker(sshKeepaliveTick)
		defer pingT.Stop()
		defer kaT.Stop()
		misses := 0
		for {
			select {
			case <-done:
				return
			case <-pingT.C:
				wsMu.Lock()
				err := ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second))
				wsMu.Unlock()
				if err != nil {
					finish()
					return
				}
			case <-kaT.C:
				if sshAlive(sshClient, 15*time.Second) {
					misses = 0
					continue
				}
				misses++
				if misses >= 3 {
					log.Printf("SSH keepalive failed: %s@%s", c.Username, c.Host)
					sshLost.Store(true)
					sshClient.Close()
					return
				}
			}
		}
	}()

	select {
	case <-outDone:
		// Remote side closed: either the shell exited or the connection dropped.
		exitCode := -1
		waitErr := make(chan error, 1)
		go func() { waitErr <- session.Wait() }()
		select {
		case err := <-waitErr:
			var ee *ssh.ExitError
			if err == nil {
				exitCode = 0
			} else if errors.As(err, &ee) {
				exitCode = ee.ExitStatus()
			} else {
				sshLost.Store(true)
			}
		case <-time.After(2 * time.Second):
			sshLost.Store(true)
		}
		if sshLost.Load() {
			printTerm("\r\n\x1b[31mSSH connection to the server was lost.\x1b[0m\r\n")
			sendCtl(map[string]interface{}{"type": "exit", "lost": true})
			closeWS(wsCloseSSHLost, "SSH connection lost")
		} else {
			sendCtl(map[string]interface{}{"type": "exit", "code": exitCode})
			closeWS(websocket.CloseNormalClosure, "session ended")
		}
	case <-readDone:
	case <-done:
	}
	finish()
	log.Printf("SSH disconnected: %s@%s", c.Username, c.Host)
}
