package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ─── AI: MODEL CALLS ─────────────────────────────────
//
// One provider-neutral conversation format (aiMsg) is translated per API:
//   Anthropic Messages (also Claude on Vertex and the Bedrock Messages endpoint),
//   OpenAI Chat Completions (also Azure OpenAI, Vertex's OpenAI-compatible endpoint and
//   any OpenAI-compatible server), and the Bedrock Converse API.
// Every call streams (server-sent events, or AWS event streams for Converse): text is
// passed to onText as it arrives, tool calls are assembled and returned with the usage.
// The context cancels the HTTP request (stop / kill), and aiRequestTimeout bounds a turn.

type aiToolDef struct {
	Name        string
	Description string
	Schema      map[string]interface{}
}

type aiToolCall struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Input    json.RawMessage `json:"input"`
	BadInput string          `json:"-"` // the arguments were not valid JSON
}

type aiToolResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}

type aiMsg struct {
	Role    string          `json:"role"` // user | assistant
	Text    string          `json:"text,omitempty"`
	Calls   []aiToolCall    `json:"calls,omitempty"`
	Results []aiToolResult  `json:"results,omitempty"`
	Raw     json.RawMessage `json:"raw,omitempty"` // Anthropic content blocks (thinking blocks are passed back unchanged)
}

type aiUsage struct {
	In, Out, CacheRead, CacheWrite int
}

type aiLLMRequest struct {
	System    string
	Messages  []aiMsg
	Tools     []aiToolDef
	Model     string
	MaxTokens int
}

type aiLLMResponse struct {
	Text    string
	Calls   []aiToolCall
	Raw     json.RawMessage
	Stop    string
	Refusal bool
	Usage   aiUsage
}

const aiRequestTimeout = 10 * time.Minute

var aiHTTPClient = &http.Client{Transport: &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	ResponseHeaderTimeout: 3 * time.Minute,
	TLSHandshakeTimeout:   20 * time.Second,
	IdleConnTimeout:       90 * time.Second,
	MaxIdleConnsPerHost:   4,
}}

// aiComplete runs one model request with streaming. Error messages never carry the
// provider's credentials (some APIs echo the key they refused).
func aiComplete(ctx context.Context, p aiProvider, req aiLLMRequest, onText func(string)) (aiLLMResponse, error) {
	res, err := aiCompleteRaw(ctx, p, req, onText)
	if err != nil && !errors.Is(err, context.Canceled) {
		msg := err.Error()
		for _, sec := range []string{p.secret.APIKey, p.secret.SecretKey, p.secret.SessionToken, p.secret.AccessKeyID} {
			if len(sec) >= 4 {
				msg = strings.ReplaceAll(msg, sec, aiRedacted)
			}
		}
		if msg != err.Error() {
			err = errors.New(msg)
		}
	}
	return res, err
}

func aiCompleteRaw(ctx context.Context, p aiProvider, req aiLLMRequest, onText func(string)) (aiLLMResponse, error) {
	if req.MaxTokens <= 0 {
		req.MaxTokens = p.MaxTokens
		if req.MaxTokens <= 0 {
			req.MaxTokens = 16000
		}
	}
	ctx, cancel := context.WithTimeout(ctx, aiRequestTimeout)
	defer cancel()
	switch p.Kind {
	case "anthropic":
		base := p.BaseURL
		if base == "" {
			base = "https://api.anthropic.com"
		}
		return anthropicStream(ctx, p, base+"/v1/messages", req, onText, func(h http.Header) {
			if p.secret.APIKey != "" {
				h.Set("x-api-key", p.secret.APIKey)
			}
		}, true)
	case "bedrock":
		if p.APIMode == "messages" {
			base := p.BaseURL
			if base == "" {
				base = "https://bedrock-mantle." + p.Region + ".api.aws"
			}
			return anthropicStream(ctx, p, base+"/anthropic/v1/messages", req, onText, func(h http.Header) {
				h.Set("Authorization", "Bearer "+p.secret.APIKey)
			}, false)
		}
		return bedrockConverse(ctx, p, req, onText)
	case "vertex":
		tok, err := vertexToken(ctx, p)
		if err != nil {
			return aiLLMResponse{}, err
		}
		base := p.BaseURL
		if base == "" {
			base = "https://" + p.Region + "-aiplatform.googleapis.com"
			if p.Region == "global" {
				base = "https://aiplatform.googleapis.com"
			}
		}
		loc := "/v1/projects/" + url.PathEscape(p.Project) + "/locations/" + url.PathEscape(p.Region)
		auth := func(h http.Header) { h.Set("Authorization", "Bearer "+tok) }
		if p.APIMode == "openai" {
			return openaiStream(ctx, p, base+loc+"/endpoints/openapi/chat/completions", req, onText, auth)
		}
		return anthropicStream(ctx, p, base+loc+"/publishers/anthropic/models/"+url.PathEscape(req.Model)+":streamRawPredict", req, onText, auth, false)
	case "azure":
		var u string
		model := req.Model
		if p.APIVersion != "" {
			dep := p.Deployment
			if dep == "" {
				dep = req.Model
			}
			u = p.BaseURL + "/openai/deployments/" + url.PathEscape(dep) + "/chat/completions?api-version=" + url.QueryEscape(p.APIVersion)
		} else {
			u = p.BaseURL + "/openai/v1/chat/completions"
		}
		req.Model = model
		return openaiStream(ctx, p, u, req, onText, func(h http.Header) { h.Set("api-key", p.secret.APIKey) })
	case "openai", "openai_compatible":
		base := p.BaseURL
		if base == "" {
			base = "https://api.openai.com/v1"
		}
		return openaiStream(ctx, p, base+"/chat/completions", req, onText, func(h http.Header) {
			if p.secret.APIKey != "" {
				h.Set("Authorization", "Bearer "+p.secret.APIKey)
			}
		})
	}
	return aiLLMResponse{}, fmt.Errorf("unknown provider kind %q", p.Kind)
}

// aiTestProvider sends one tiny request (no tools) to check the endpoint and the key.
func aiTestProvider(ctx context.Context, p aiProvider, model string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	res, err := aiComplete(ctx, p, aiLLMRequest{
		System: "You are a connectivity check.", Model: model, MaxTokens: 256,
		Messages: []aiMsg{{Role: "user", Text: "Reply with the single word OK."}},
	}, nil)
	if err != nil {
		return "", err
	}
	return res.Text, nil
}

// ─── HTTP helpers ────────────────────────────────────

type aiAPIError struct {
	Status int
	Msg    string
}

func (e *aiAPIError) Error() string {
	return fmt.Sprintf("the provider answered HTTP %d: %s", e.Status, e.Msg)
}

func aiPost(ctx context.Context, u string, body []byte, setHeaders func(http.Header), sign func(*http.Request) error) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "WRM-PRO/"+AppVersion)
	if setHeaders != nil {
		setHeaders(req.Header)
	}
	if sign != nil {
		if err := sign(req); err != nil {
			return nil, err
		}
	}
	resp, err := aiHTTPClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("cannot reach the provider: %v", err)
	}
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		resp.Body.Close()
		return nil, &aiAPIError{Status: resp.StatusCode, Msg: aiErrorMessage(b)}
	}
	return resp, nil
}

// aiErrorMessage extracts the message of a provider error body (never the request).
func aiErrorMessage(b []byte) string {
	var e struct {
		Error   json.RawMessage `json:"error"`
		Message string          `json:"message"`
	}
	if json.Unmarshal(b, &e) == nil {
		var inner struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		}
		if json.Unmarshal(e.Error, &inner) == nil && inner.Message != "" {
			return truncateStr(redactText(inner.Message), 300)
		}
		var s string
		if json.Unmarshal(e.Error, &s) == nil && s != "" {
			return truncateStr(redactText(s), 300)
		}
		if e.Message != "" {
			return truncateStr(redactText(e.Message), 300)
		}
	}
	return truncateStr(redactText(strings.TrimSpace(string(b))), 300)
}

// sseEvents reads server-sent events and calls fn(event, data) for each one.
func sseEvents(r io.Reader, fn func(event, data string) error) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64<<10), 8<<20)
	event := ""
	var data []string
	flush := func() error {
		if len(data) == 0 {
			event = ""
			return nil
		}
		err := fn(event, strings.Join(data, "\n"))
		event, data = "", nil
		return err
	}
	for sc.Scan() {
		line := sc.Text()
		switch {
		case line == "":
			if err := flush(); err != nil {
				return err
			}
		case strings.HasPrefix(line, ":"):
		case strings.HasPrefix(line, "event:"):
			event = strings.TrimSpace(line[6:])
		case strings.HasPrefix(line, "data:"):
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return flush()
}

// parseToolInput validates streamed tool arguments: the result must be a JSON object.
func parseToolInput(s string) (json.RawMessage, string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return json.RawMessage(`{}`), ""
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return json.RawMessage(`{}`), "the tool arguments are not valid JSON: " + err.Error()
	}
	return json.RawMessage(s), ""
}

// ─── Anthropic Messages API ──────────────────────────

func anthropicModelSupportsFallback(model string) bool {
	m := strings.TrimPrefix(model, "anthropic.")
	for _, x := range []string{"claude-fable-5-1", "claude-opus-5-5", "claude-opus-5", "claude-sonnet-5-5"} {
		if m == x {
			return true
		}
	}
	return false
}

func anthropicMessages(msgs []aiMsg) []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, m := range msgs {
		var blocks []interface{}
		if m.Role == "assistant" && len(m.Raw) > 0 {
			var raw []interface{}
			if json.Unmarshal(m.Raw, &raw) == nil && len(raw) > 0 {
				out = append(out, map[string]interface{}{"role": "assistant", "content": raw})
				continue
			}
		}
		for _, r := range m.Results {
			b := map[string]interface{}{"type": "tool_result", "tool_use_id": r.ID, "content": r.Content}
			if r.IsError {
				b["is_error"] = true
			}
			blocks = append(blocks, b)
		}
		if m.Text != "" {
			blocks = append(blocks, map[string]interface{}{"type": "text", "text": m.Text})
		}
		for _, c := range m.Calls {
			blocks = append(blocks, map[string]interface{}{"type": "tool_use", "id": c.ID, "name": c.Name, "input": c.Input})
		}
		if len(blocks) == 0 {
			blocks = append(blocks, map[string]interface{}{"type": "text", "text": "(empty)"})
		}
		out = append(out, map[string]interface{}{"role": m.Role, "content": blocks})
	}
	return out
}

// anthropicStream calls a Messages API endpoint. first marks Anthropic's own API
// (prompt caching, eager tool input streaming and refusal fallbacks are sent there only).
func anthropicStream(ctx context.Context, p aiProvider, endpoint string, req aiLLMRequest, onText func(string), auth func(http.Header), first bool) (aiLLMResponse, error) {
	body := map[string]interface{}{
		"max_tokens": req.MaxTokens,
		"messages":   anthropicMessages(req.Messages),
		"stream":     true,
	}
	if req.System != "" {
		body["system"] = []map[string]interface{}{{"type": "text", "text": req.System}}
	}
	if p.Kind == "vertex" {
		body["anthropic_version"] = "vertex-2023-10-16"
	} else {
		body["model"] = req.Model
	}
	if len(req.Tools) > 0 {
		tools := []map[string]interface{}{}
		for _, t := range req.Tools {
			td := map[string]interface{}{"name": t.Name, "description": t.Description, "input_schema": t.Schema}
			if first && p.BaseURL == "" {
				td["eager_input_streaming"] = true
			}
			tools = append(tools, td)
		}
		body["tools"] = tools
	}
	beta := ""
	if first && p.BaseURL == "" {
		body["cache_control"] = map[string]string{"type": "ephemeral"}
		if p.Fallback && anthropicModelSupportsFallback(req.Model) {
			body["fallbacks"] = "default"
			beta = "server-side-fallback-2026-07-01"
		}
	}
	b, _ := json.Marshal(body)
	resp, err := aiPost(ctx, endpoint, b, func(h http.Header) {
		h.Set("anthropic-version", "2023-06-01")
		h.Set("Accept", "text/event-stream")
		if beta != "" {
			h.Set("anthropic-beta", beta)
		}
		auth(h)
	}, nil)
	if err != nil {
		return aiLLMResponse{}, err
	}
	defer resp.Body.Close()
	var out aiLLMResponse
	type block struct {
		data    map[string]interface{}
		partial strings.Builder
	}
	blocks := map[int]*block{}
	maxIdx := -1
	var text strings.Builder
	err = sseEvents(resp.Body, func(event, data string) error {
		var ev map[string]json.RawMessage
		if json.Unmarshal([]byte(data), &ev) != nil {
			return nil
		}
		var typ string
		json.Unmarshal(ev["type"], &typ)
		if typ == "" {
			typ = event
		}
		switch typ {
		case "message_start":
			var m struct {
				Message struct {
					Usage struct {
						In         int `json:"input_tokens"`
						Out        int `json:"output_tokens"`
						CacheRead  int `json:"cache_read_input_tokens"`
						CacheWrite int `json:"cache_creation_input_tokens"`
					} `json:"usage"`
				} `json:"message"`
			}
			json.Unmarshal([]byte(data), &m)
			out.Usage.In += m.Message.Usage.In
			out.Usage.Out += m.Message.Usage.Out
			out.Usage.CacheRead += m.Message.Usage.CacheRead
			out.Usage.CacheWrite += m.Message.Usage.CacheWrite
		case "content_block_start":
			var m struct {
				Index int                    `json:"index"`
				Block map[string]interface{} `json:"content_block"`
			}
			json.Unmarshal([]byte(data), &m)
			blocks[m.Index] = &block{data: m.Block}
			if m.Index > maxIdx {
				maxIdx = m.Index
			}
		case "content_block_delta":
			var m struct {
				Index int `json:"index"`
				Delta struct {
					Type      string `json:"type"`
					Text      string `json:"text"`
					Partial   string `json:"partial_json"`
					Thinking  string `json:"thinking"`
					Signature string `json:"signature"`
				} `json:"delta"`
			}
			json.Unmarshal([]byte(data), &m)
			bl := blocks[m.Index]
			if bl == nil {
				bl = &block{data: map[string]interface{}{"type": "text", "text": ""}}
				blocks[m.Index] = bl
				if m.Index > maxIdx {
					maxIdx = m.Index
				}
			}
			switch m.Delta.Type {
			case "text_delta":
				s, _ := bl.data["text"].(string)
				bl.data["text"] = s + m.Delta.Text
				text.WriteString(m.Delta.Text)
				if onText != nil {
					onText(m.Delta.Text)
				}
			case "input_json_delta":
				bl.partial.WriteString(m.Delta.Partial)
			case "thinking_delta":
				s, _ := bl.data["thinking"].(string)
				bl.data["thinking"] = s + m.Delta.Thinking
			case "signature_delta":
				s, _ := bl.data["signature"].(string)
				bl.data["signature"] = s + m.Delta.Signature
			}
		case "message_delta":
			var m struct {
				Delta struct {
					Stop string `json:"stop_reason"`
				} `json:"delta"`
				Usage struct {
					Out int `json:"output_tokens"`
				} `json:"usage"`
			}
			json.Unmarshal([]byte(data), &m)
			if m.Delta.Stop != "" {
				out.Stop = m.Delta.Stop
			}
			if m.Usage.Out > 0 {
				out.Usage.Out = m.Usage.Out
			}
		case "error":
			var m struct {
				Error struct {
					Type    string `json:"type"`
					Message string `json:"message"`
				} `json:"error"`
			}
			json.Unmarshal([]byte(data), &m)
			return fmt.Errorf("the provider reported an error: %s", truncateStr(redactText(m.Error.Message+" ("+m.Error.Type+")"), 300))
		}
		return nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		return out, err
	}
	var raw []map[string]interface{}
	for i := 0; i <= maxIdx; i++ {
		bl := blocks[i]
		if bl == nil {
			continue
		}
		if bl.data["type"] == "tool_use" {
			in, bad := parseToolInput(bl.partial.String())
			if bl.partial.Len() == 0 {
				if v, ok := bl.data["input"]; ok && v != nil {
					b, _ := json.Marshal(v)
					in, bad = parseToolInput(string(b))
				}
			}
			var obj interface{}
			json.Unmarshal(in, &obj)
			bl.data["input"] = obj
			id, _ := bl.data["id"].(string)
			name, _ := bl.data["name"].(string)
			out.Calls = append(out.Calls, aiToolCall{ID: id, Name: name, Input: in, BadInput: bad})
		}
		raw = append(raw, bl.data)
	}
	out.Raw, _ = json.Marshal(raw)
	out.Text = text.String()
	if out.Stop == "refusal" {
		out.Refusal = true
		out.Calls = nil
	}
	if out.Stop == "max_tokens" {
		// a tool call cut off by the token limit is not run
		for i := range out.Calls {
			if out.Calls[i].BadInput == "" {
				continue
			}
			out.Calls[i].BadInput = "the response hit the token limit before the tool arguments were complete"
		}
	}
	return out, nil
}

// ─── OpenAI Chat Completions (and compatible APIs) ───

func openaiMessages(system string, msgs []aiMsg) []map[string]interface{} {
	out := []map[string]interface{}{}
	if system != "" {
		out = append(out, map[string]interface{}{"role": "system", "content": system})
	}
	for _, m := range msgs {
		for _, r := range m.Results {
			out = append(out, map[string]interface{}{"role": "tool", "tool_call_id": r.ID, "content": r.Content})
		}
		if m.Role == "user" && m.Text == "" {
			continue
		}
		if m.Role == "user" {
			out = append(out, map[string]interface{}{"role": "user", "content": m.Text})
			continue
		}
		msg := map[string]interface{}{"role": "assistant"}
		if m.Text != "" {
			msg["content"] = m.Text
		} else {
			msg["content"] = nil
		}
		if len(m.Calls) > 0 {
			calls := []map[string]interface{}{}
			for _, c := range m.Calls {
				calls = append(calls, map[string]interface{}{"id": c.ID, "type": "function",
					"function": map[string]interface{}{"name": c.Name, "arguments": string(c.Input)}})
			}
			msg["tool_calls"] = calls
		}
		out = append(out, msg)
	}
	return out
}

func openaiStream(ctx context.Context, p aiProvider, endpoint string, req aiLLMRequest, onText func(string), auth func(http.Header)) (aiLLMResponse, error) {
	body := map[string]interface{}{
		"model":    req.Model,
		"messages": openaiMessages(req.System, req.Messages),
		"stream":   true,
	}
	if p.Kind == "openai" || p.Kind == "azure" {
		body["max_completion_tokens"] = req.MaxTokens
		body["stream_options"] = map[string]bool{"include_usage": true}
	} else {
		body["max_tokens"] = req.MaxTokens
		body["stream_options"] = map[string]bool{"include_usage": true}
	}
	if len(req.Tools) > 0 {
		tools := []map[string]interface{}{}
		for _, t := range req.Tools {
			tools = append(tools, map[string]interface{}{"type": "function", "function": map[string]interface{}{
				"name": t.Name, "description": t.Description, "parameters": t.Schema}})
		}
		body["tools"] = tools
	}
	b, _ := json.Marshal(body)
	resp, err := aiPost(ctx, endpoint, b, func(h http.Header) {
		h.Set("Accept", "text/event-stream")
		auth(h)
	}, nil)
	if err != nil {
		return aiLLMResponse{}, err
	}
	defer resp.Body.Close()
	var out aiLLMResponse
	type call struct {
		id, name string
		args     strings.Builder
	}
	calls := map[int]*call{}
	maxIdx := -1
	var text strings.Builder
	err = sseEvents(resp.Body, func(event, data string) error {
		if data == "[DONE]" {
			return nil
		}
		var m struct {
			Choices []struct {
				Delta struct {
					Content   *string `json:"content"`
					ToolCalls []struct {
						Index    int    `json:"index"`
						ID       string `json:"id"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"delta"`
				Finish string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				Prompt     int `json:"prompt_tokens"`
				Completion int `json:"completion_tokens"`
				Details    struct {
					Cached int `json:"cached_tokens"`
				} `json:"prompt_tokens_details"`
			} `json:"usage"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal([]byte(data), &m) != nil {
			return nil
		}
		if m.Error != nil {
			return fmt.Errorf("the provider reported an error: %s", truncateStr(redactText(m.Error.Message), 300))
		}
		if m.Usage != nil {
			out.Usage.In = m.Usage.Prompt - m.Usage.Details.Cached
			out.Usage.CacheRead = m.Usage.Details.Cached
			out.Usage.Out = m.Usage.Completion
		}
		for _, ch := range m.Choices {
			if ch.Delta.Content != nil && *ch.Delta.Content != "" {
				text.WriteString(*ch.Delta.Content)
				if onText != nil {
					onText(*ch.Delta.Content)
				}
			}
			for _, tc := range ch.Delta.ToolCalls {
				c := calls[tc.Index]
				if c == nil {
					c = &call{}
					calls[tc.Index] = c
					if tc.Index > maxIdx {
						maxIdx = tc.Index
					}
				}
				if tc.ID != "" {
					c.id = tc.ID
				}
				if tc.Function.Name != "" {
					c.name += tc.Function.Name
				}
				c.args.WriteString(tc.Function.Arguments)
			}
			if ch.Finish != "" {
				out.Stop = ch.Finish
			}
		}
		return nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		return out, err
	}
	for i := 0; i <= maxIdx; i++ {
		c := calls[i]
		if c == nil {
			continue
		}
		if c.id == "" {
			c.id = "call_" + randomToken(8)
		}
		in, bad := parseToolInput(c.args.String())
		out.Calls = append(out.Calls, aiToolCall{ID: c.id, Name: c.name, Input: in, BadInput: bad})
	}
	out.Text = text.String()
	if out.Stop == "content_filter" {
		out.Refusal = true
	}
	return out, nil
}

// ─── AWS Bedrock Converse ────────────────────────────

func bedrockMessages(msgs []aiMsg) []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, m := range msgs {
		var content []interface{}
		for _, r := range m.Results {
			status := "success"
			if r.IsError {
				status = "error"
			}
			content = append(content, map[string]interface{}{"toolResult": map[string]interface{}{
				"toolUseId": r.ID, "status": status, "content": []map[string]string{{"text": nonEmpty(r.Content, "(no output)")}}}})
		}
		if m.Text != "" {
			content = append(content, map[string]string{"text": m.Text})
		}
		for _, c := range m.Calls {
			var in interface{}
			json.Unmarshal(c.Input, &in)
			content = append(content, map[string]interface{}{"toolUse": map[string]interface{}{"toolUseId": c.ID, "name": c.Name, "input": in}})
		}
		if len(content) == 0 {
			content = append(content, map[string]string{"text": "(empty)"})
		}
		out = append(out, map[string]interface{}{"role": m.Role, "content": content})
	}
	return out
}

func nonEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

func bedrockConverse(ctx context.Context, p aiProvider, req aiLLMRequest, onText func(string)) (aiLLMResponse, error) {
	body := map[string]interface{}{
		"messages":        bedrockMessages(req.Messages),
		"inferenceConfig": map[string]interface{}{"maxTokens": req.MaxTokens},
	}
	if req.System != "" {
		body["system"] = []map[string]string{{"text": req.System}}
	}
	if len(req.Tools) > 0 {
		tools := []map[string]interface{}{}
		for _, t := range req.Tools {
			tools = append(tools, map[string]interface{}{"toolSpec": map[string]interface{}{
				"name": t.Name, "description": t.Description, "inputSchema": map[string]interface{}{"json": t.Schema}}})
		}
		body["toolConfig"] = map[string]interface{}{"tools": tools}
	}
	b, _ := json.Marshal(body)
	base := p.BaseURL
	if base == "" {
		base = "https://bedrock-runtime." + p.Region + ".amazonaws.com"
	}
	u, err := url.Parse(base)
	if err != nil {
		return aiLLMResponse{}, err
	}
	rawPath := strings.TrimRight(u.EscapedPath(), "/") + "/model/" + awsURIEncode(req.Model) + "/converse-stream"
	u.RawPath = rawPath
	u.Path, _ = url.PathUnescape(rawPath)
	var sign func(*http.Request) error
	setH := func(h http.Header) { h.Set("Accept", "application/vnd.amazon.eventstream") }
	if p.secret.APIKey != "" {
		setH = func(h http.Header) {
			h.Set("Accept", "application/vnd.amazon.eventstream")
			h.Set("Authorization", "Bearer "+p.secret.APIKey)
		}
	} else {
		sign = func(r *http.Request) error {
			return sigV4Sign(r, b, p.secret.AccessKeyID, p.secret.SecretKey, p.secret.SessionToken, p.Region, "bedrock", time.Now().UTC())
		}
	}
	resp, err := aiPost(ctx, u.String(), b, setH, sign)
	if err != nil {
		return aiLLMResponse{}, err
	}
	defer resp.Body.Close()
	var out aiLLMResponse
	type tu struct {
		id, name string
		input    strings.Builder
	}
	tools := map[int]*tu{}
	maxIdx := -1
	var text strings.Builder
	err = readEventStream(resp.Body, func(msgType, eventType string, payload []byte) error {
		if msgType == "exception" || msgType == "error" {
			var e struct {
				Message string `json:"message"`
			}
			json.Unmarshal(payload, &e)
			return fmt.Errorf("the provider reported an error: %s %s", eventType, truncateStr(redactText(e.Message), 300))
		}
		switch eventType {
		case "contentBlockStart":
			var m struct {
				Index int `json:"contentBlockIndex"`
				Start struct {
					ToolUse *struct {
						ID   string `json:"toolUseId"`
						Name string `json:"name"`
					} `json:"toolUse"`
				} `json:"start"`
			}
			json.Unmarshal(payload, &m)
			if m.Start.ToolUse != nil {
				tools[m.Index] = &tu{id: m.Start.ToolUse.ID, name: m.Start.ToolUse.Name}
				if m.Index > maxIdx {
					maxIdx = m.Index
				}
			}
		case "contentBlockDelta":
			var m struct {
				Index int `json:"contentBlockIndex"`
				Delta struct {
					Text    *string `json:"text"`
					ToolUse *struct {
						Input string `json:"input"`
					} `json:"toolUse"`
				} `json:"delta"`
			}
			json.Unmarshal(payload, &m)
			if m.Delta.Text != nil {
				text.WriteString(*m.Delta.Text)
				if onText != nil {
					onText(*m.Delta.Text)
				}
			}
			if m.Delta.ToolUse != nil && tools[m.Index] != nil {
				tools[m.Index].input.WriteString(m.Delta.ToolUse.Input)
			}
		case "messageStop":
			var m struct {
				Stop string `json:"stopReason"`
			}
			json.Unmarshal(payload, &m)
			out.Stop = m.Stop
		case "metadata":
			var m struct {
				Usage struct {
					In         int `json:"inputTokens"`
					Out        int `json:"outputTokens"`
					CacheRead  int `json:"cacheReadInputTokens"`
					CacheWrite int `json:"cacheWriteInputTokens"`
				} `json:"usage"`
			}
			json.Unmarshal(payload, &m)
			out.Usage = aiUsage{In: m.Usage.In, Out: m.Usage.Out, CacheRead: m.Usage.CacheRead, CacheWrite: m.Usage.CacheWrite}
		}
		return nil
	})
	if err != nil {
		if ctx.Err() != nil {
			return out, ctx.Err()
		}
		return out, err
	}
	for i := 0; i <= maxIdx; i++ {
		t := tools[i]
		if t == nil {
			continue
		}
		in, bad := parseToolInput(t.input.String())
		out.Calls = append(out.Calls, aiToolCall{ID: t.id, Name: t.name, Input: in, BadInput: bad})
	}
	out.Text = text.String()
	if out.Stop == "guardrail_intervened" || out.Stop == "content_filtered" {
		out.Refusal = true
	}
	return out, nil
}

// ─── cost ────────────────────────────────────────────

// aiBuiltinPrices are Anthropic's first-party list prices in USD per million tokens
// (input, output). Other providers use the prices configured on the provider.
var aiBuiltinPrices = map[string][2]float64{
	"claude-fable-5-1":  {10, 50},
	"claude-mythos-5-1": {10, 50},
	"claude-fable-5":    {10, 50},
	"claude-opus-5-5":   {4, 20},
	"claude-opus-5":     {5, 25},
	"claude-opus-4-8":   {5, 25},
	"claude-opus-4-7":   {5, 25},
	"claude-opus-4-6":   {5, 25},
	"claude-sonnet-5-5": {2, 10},
	"claude-sonnet-5":   {2, 10},
	"claude-sonnet-4-6": {3, 15},
	"claude-haiku-5-5":  {0.10, 0.50},
	"claude-haiku-4-5":  {1, 5},
}

// aiCost returns the estimated cost of a request (0 when no price is known).
func aiCost(p aiProvider, model string, u aiUsage) float64 {
	in, out := p.PriceIn, p.PriceOut
	if in == 0 && out == 0 && p.Kind == "anthropic" {
		if pr, ok := aiBuiltinPrices[model]; ok {
			in, out = pr[0], pr[1]
		}
	}
	return (float64(u.In)*in + float64(u.CacheRead)*in*0.1 + float64(u.CacheWrite)*in*1.25 + float64(u.Out)*out) / 1e6
}

var errAIStopped = errors.New("stopped")
