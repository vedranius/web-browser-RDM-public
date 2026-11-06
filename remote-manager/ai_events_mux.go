package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GET /api/ai/events?s=<id>:<since>,<id>:<since>,… streams the events of several of the
// user's own AI sessions over ONE server-sent-events connection.
//
// A browser keeps at most six HTTP/1.1 connections per host, shared by all its tabs. With one
// EventSource per watched session (the MCP live view opens one for every session an AI app
// starts), six sessions took every connection and each further request of the page waited,
// so WRM looked frozen while an AI worked through MCP. The page now opens this one stream per
// tab and reopens it when the set of watched sessions changes.
//
// Each event is sent as `event: ai` with data {"sid": <session>, "ev": <the session event>}.
// A session that is not running (ended, or another user's) gets one synthetic "ended" event,
// as the per-session endpoint does, and is then dropped from the stream.

const aiEventsMuxMax = 64

type aiMuxSub struct {
	id    int64
	since int64
	s     *aiSession
}

func aiParseMuxSubs(q string) []aiMuxSub {
	var out []aiMuxSub
	seen := map[int64]bool{}
	for _, part := range strings.Split(q, ",") {
		idStr, sinceStr, _ := strings.Cut(strings.TrimSpace(part), ":")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 || seen[id] {
			continue
		}
		since, _ := strconv.ParseInt(sinceStr, 10, 64)
		if since < 0 {
			since = 0
		}
		seen[id] = true
		out = append(out, aiMuxSub{id: id, since: since})
		if len(out) >= aiEventsMuxMax {
			break
		}
	}
	return out
}

// aiEndedEvent is the event a stream sends for a session that is not running.
func aiEndedEvent(id int64, userID int) aiEvent {
	status, reason := "ended", "ended"
	var st, rs string
	if db.QueryRow(`SELECT status, end_reason FROM ai_sessions WHERE id=? AND user_id=?`, id, userID).Scan(&st, &rs) == nil {
		status, reason = st, nonEmpty(rs, "ended")
		if status == "active" {
			status, reason = "ended", "the session is no longer running (WRM was restarted)"
		}
	} else {
		reason = "AI session not found or ended"
	}
	return aiEvent{Seq: 1 << 40, Type: "ended", TS: nowRFC(), Data: map[string]interface{}{"reason": reason, "status": status, "resumable": true}}
}

func aiEventsMuxHandler(w http.ResponseWriter, r *http.Request, userID int) {
	if r.Method != http.MethodGet {
		jsonError(w, "Method not allowed", 405)
		return
	}
	fl, ok := w.(http.Flusher)
	if !ok {
		jsonError(w, "Streaming is not supported", 500)
		return
	}
	subs := aiParseMuxSubs(r.URL.Query().Get("s"))
	if len(subs) == 0 {
		jsonError(w, "No sessions to watch", 400)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(200)
	fmt.Fprintf(w, "retry: 3000\n\n")
	send := func(id int64, e aiEvent) {
		b, _ := json.Marshal(map[string]interface{}{"sid": id, "ev": e})
		fmt.Fprintf(w, "event: ai\ndata: %s\n\n", b)
	}
	live := subs[:0]
	for _, sub := range subs {
		if s := aiGet(sub.id); s != nil && s.UserID == userID {
			sub.s = s
			live = append(live, sub)
		} else {
			send(sub.id, aiEndedEvent(sub.id, userID))
		}
	}
	fl.Flush()
	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for len(live) > 0 {
		wake := make(chan struct{}, 1)
		stop := make(chan struct{})
		sent := false
		next := live[:0]
		for _, sub := range live {
			evs, wait, ended := sub.s.eventsSince(sub.since)
			for _, e := range evs {
				send(sub.id, e)
				sub.since = e.Seq
				sent = true
			}
			if ended {
				continue
			}
			next = append(next, sub)
			go func(wait <-chan struct{}) {
				select {
				case <-wait:
					select {
					case wake <- struct{}{}:
					default:
					}
				case <-stop:
				}
			}(wait)
		}
		live = next
		if sent {
			fl.Flush()
		}
		if len(live) == 0 {
			close(stop)
			return
		}
		select {
		case <-wake:
		case <-ping.C:
			fmt.Fprintf(w, ": ping\n\n")
			fl.Flush()
		case <-r.Context().Done():
			close(stop)
			return
		}
		close(stop)
	}
}
