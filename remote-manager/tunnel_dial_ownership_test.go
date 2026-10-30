package main

import (
	"net"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func registerTestTunnelOwner(t *testing.T, addr string, owner int) {
	t.Helper()
	tr := &tunnelRun{key: t.Name() + ":" + addr, def: tunnelDef{Kind: "local"}, ownerID: owner, boundAddr: addr}
	tunnelMgr.mu.Lock()
	tunnelMgr.runs[tr.key] = tr
	tunnelMgr.mu.Unlock()
	t.Cleanup(func() {
		tunnelMgr.mu.Lock()
		delete(tunnelMgr.runs, tr.key)
		tunnelMgr.mu.Unlock()
	})
}

func localAlias(t *testing.T, addr string) string {
	t.Helper()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	return net.JoinHostPort("localhost", port)
}

func TestLocalDialChecksResolvedTunnelAddress(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	registerTestTunnelOwner(t, ln.Addr().String(), 101)
	// The dial resolves a hostname first; the ownership check sees only its
	// actual RemoteAddr, so every DNS alias receives the same protection.
	alias := localAlias(t, ln.Addr().String())
	if conn, err := dialLocalForUser(202, "tcp", alias, time.Second); err == nil || !strings.Contains(err.Error(), "another user") {
		if conn != nil {
			conn.Close()
		}
		t.Fatalf("another user's resolved local tunnel: conn=%v err=%v", conn, err)
	}
	conn, err := dialLocalForUser(101, "tcp", alias, time.Second)
	if err != nil {
		t.Fatalf("owner's resolved local tunnel: %v", err)
	}
	conn.Close()

	// A remote forward's target is also dialed on WRM, rather than on the SSH
	// server, so it must receive the same check.
	_, portText, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portText)
	remote := &tunnelRun{def: tunnelDef{Kind: "remote", TargetHost: "localhost", TargetPort: port}, ownerID: 202}
	in, peer := net.Pipe()
	defer peer.Close()
	done := make(chan struct{})
	go func() { remote.handle(in); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("remote forward did not refuse another user's local tunnel")
	}
	if !strings.Contains(remote.lastErr, "another user") {
		t.Fatalf("remote forward entered another user's local tunnel: %q", remote.lastErr)
	}
}

func TestDirectProtocolsCheckTunnelOwner(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	owner := newTestUser(t, srv, "local-tunnel-owner", false)
	other := newTestUser(t, srv, "local-tunnel-other", false)
	sshAddr := startTestSSHServer(t)
	registerTestTunnelOwner(t, sshAddr, owner.userID)
	otherID := other.addConnection("local alias", localAlias(t, sshAddr))
	c, err := loadConnection(otherID)
	if err != nil {
		t.Fatal(err)
	}
	if client, err := dialSSH(c, nil); err == nil || !strings.Contains(err.Error(), "another user") {
		if client != nil {
			client.Close()
		}
		t.Fatalf("SSH entered another user's tunnel: %v", err)
	}
	if pooled, err := acquireSFTP(c); err == nil || !strings.Contains(err.Error(), "another user") {
		if pooled != nil {
			releaseSFTP(pooled, nil)
		}
		t.Fatalf("SFTP entered another user's tunnel: %v", err)
	}
	if err := testWebReachable(c); err == nil || !strings.Contains(err.Error(), "another user") {
		t.Fatalf("web test entered another user's tunnel: %v", err)
	}

	ownerID := owner.addConnection("owned local alias", localAlias(t, sshAddr))
	owned, _ := loadConnection(ownerID)
	client, err := dialSSH(owned, nil)
	if err != nil {
		t.Fatalf("SSH owner could not use own tunnel: %v", err)
	}
	client.Close()
	if err := testWebReachable(owned); err != nil {
		t.Fatalf("web test owner could not use own tunnel: %v", err)
	}

	ftpAddr := startFakeFTP(t)
	registerTestTunnelOwner(t, ftpAddr, owner.userID)
	ftpConn := Connection{Host: localAlias(t, ftpAddr), Protocol: "FTP", Username: "anon", Password: "x", UserID: other.userID}
	if fc, done, err := dialFTP(ftpConn); err == nil || !strings.Contains(err.Error(), "another user") {
		if fc != nil {
			done()
		}
		t.Fatalf("FTP entered another user's tunnel: %v", err)
	}
	ftpConn.UserID = owner.userID
	_, done, err := dialFTP(ftpConn)
	if err != nil {
		t.Fatalf("FTP owner could not use own tunnel: %v", err)
	}
	done()
}

func TestJumpTargetLoopbackIsRemote(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()
	u := newTestUser(t, srv, "remote-loopback-user", false)
	jumpAddr := startTestSSHServer(t)
	targetAddr := startTestSSHServer(t)
	registerTestTunnelOwner(t, targetAddr, u.userID+1000)
	jump := u.addConnection("jump for remote loopback", jumpAddr)
	target := u.addConnectionVia("remote loopback", targetAddr, jump)
	c, err := loadConnection(target)
	if err != nil {
		t.Fatal(err)
	}
	client, err := dialSSH(c, nil)
	if err != nil {
		t.Fatalf("remote loopback target incorrectly treated as a local tunnel: %v", err)
	}
	client.Close()
}
