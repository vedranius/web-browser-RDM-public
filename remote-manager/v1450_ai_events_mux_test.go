package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

type muxEvent struct {
	SID int64   `json:"sid"`
	Ev  aiEvent `json:"ev"`
}

// readMux reads the multiplexed AI event stream until stop returns true or the time is up.
func readMux(t *testing.T, c *testClient, query string, d time.Duration, stop func([]muxEvent) bool) ([]muxEvent, http.Header, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", c.srv.URL+"/api/ai/events?s="+query, nil)
	req.AddCookie(c.cookie)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var out []muxEvent
	if res.StatusCode != 200 {
		return out, res.Header, res.StatusCode
	}
	sc := bufio.NewScanner(res.Body)
	sc.Buffer(make([]byte, 1<<20), 8<<20)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var m muxEvent
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &m) == nil {
			out = append(out, m)
		}
		if stop(out) {
			break
		}
	}
	return out, res.Header, res.StatusCode
}

// Several AI sessions are streamed over ONE connection: a browser keeps six HTTP/1.1
// connections per host, and one stream per MCP session froze the page at six sessions.
func TestV1450AIEventsOneStreamForManySessions(t *testing.T) {
	e := newMCPEnv(t)
	ids := []int{e.connID}
	for i := 1; i < 3; i++ {
		ids = append(ids, e.addConn(fmt.Sprintf("demo-mux-%d", i)))
	}
	tok, _ := e.pat("auto", []string{"read_logs", "run_with_approval", "edit_file_with_approval"}, map[string]interface{}{"conn_ids": ids, "confirm_auto": true,
		"auto": map[string]interface{}{"allow": []string{"*"}, "path_allow": []string{"/**"}, "until_grant": true}})
	c := e.client(tok)
	c.initialize(false)
	var sids []int64
	for _, id := range ids {
		var n int64
		fmt.Sscan(c.open(id), &n)
		sids = append(sids, n)
	}
	for _, sid := range sids {
		if txt, bad := c.call("run_command", map[string]interface{}{"session_id": fmt.Sprint(sid), "command": "echo mux-ok", "reason": "test"}); bad {
			t.Fatalf("run_command: %s", txt)
		}
	}
	q := fmt.Sprintf("%d:0,%d:0,%d:0,999999:0,%d:0", sids[0], sids[1], sids[2], sids[0])
	got := map[int64][]string{}
	evs, hdr, code := readMux(t, e.user, q, 10*time.Second, func(evs []muxEvent) bool {
		results := 0
		for _, m := range evs {
			if m.Ev.Type == "tool_result" {
				results++
			}
		}
		return results >= 3
	})
	if code != 200 || !strings.HasPrefix(hdr.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("stream: %d %v", code, hdr)
	}
	for _, m := range evs {
		got[m.SID] = append(got[m.SID], m.Ev.Type)
	}
	for _, sid := range sids {
		if !strings.Contains(strings.Join(got[sid], ","), "tool_result") {
			t.Fatalf("session %d: no tool_result in %v", sid, got)
		}
	}
	if fmt.Sprint(got[999999]) != "[ended]" {
		t.Fatalf("unknown session must end at once: %v", got[999999])
	}
	// the same session listed twice is streamed once
	n := 0
	for _, m := range evs {
		if m.SID == sids[0] && m.Ev.Type == "tool_result" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("duplicate session in the query streamed %d times", n)
	}
	// a position skips the events already seen
	last := int64(0)
	for _, m := range evs {
		if m.SID == sids[1] && m.Ev.Seq > last {
			last = m.Ev.Seq
		}
	}
	if txt, bad := c.call("run_command", map[string]interface{}{"session_id": fmt.Sprint(sids[1]), "command": "echo again", "reason": "test"}); bad {
		t.Fatalf("run_command: %s", txt)
	}
	evs, _, _ = readMux(t, e.user, fmt.Sprintf("%d:%d", sids[1], last), 10*time.Second, func(evs []muxEvent) bool {
		for _, m := range evs {
			if m.Ev.Type == "tool_result" {
				return true
			}
		}
		return false
	})
	for _, m := range evs {
		if m.Ev.Seq <= last {
			t.Fatalf("event %d replayed after position %d", m.Ev.Seq, last)
		}
	}
	if len(evs) == 0 {
		t.Fatal("no new events after the position")
	}
}

// Another user's session is never streamed: it looks ended, as on the per-session endpoint.
func TestV1450AIEventsOnlyOwnSessions(t *testing.T) {
	e := newMCPEnv(t)
	tok, _ := e.pat("ask", []string{"read_logs"}, nil)
	c := e.client(tok)
	c.initialize(false)
	var sid int64
	fmt.Sscan(c.open("demo-api"), &sid)
	evs, _, code := readMux(t, e.admin, fmt.Sprintf("%d:0", sid), 5*time.Second, func([]muxEvent) bool { return false })
	if code != 200 || len(evs) != 1 || evs[0].Ev.Type != "ended" {
		t.Fatalf("another user's session: %d %v", code, evs)
	}
	if _, _, code := readMux(t, e.user, "x,,-1:5", 5*time.Second, func([]muxEvent) bool { return true }); code != 400 {
		t.Fatalf("no valid session: %d", code)
	}
	if len(aiParseMuxSubs(strings.Repeat("1:0,", 10)+func() string {
		var b []string
		for i := 2; i < 200; i++ {
			b = append(b, fmt.Sprint(i))
		}
		return strings.Join(b, ",")
	}())) != aiEventsMuxMax {
		t.Fatal("the stream must cap the number of sessions")
	}
}
