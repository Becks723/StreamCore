package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"StreamCore/internal/pkg/ai/llm"
)

// Client communicates with MCP servers over HTTP (streamable HTTP transport).
type Client interface {
	DiscoverTools(ctx context.Context, serverURL string) ([]llm.ToolDef, error)
	CallTool(ctx context.Context, serverURL, toolName string, args json.RawMessage) (string, error)
}

type httpClient struct {
	hc        *http.Client
	sessionID string // cached session ID per server instance
}

// NewClient creates an MCP Client with the given timeout.
func NewClient(timeout time.Duration) Client {
	return &httpClient{
		hc: &http.Client{Timeout: timeout},
	}
}

// ---- request/response types ----

type jsonRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
	Params  any    `json:"params,omitempty"`
}

type initializeParams struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    struct{}   `json:"capabilities"`
	ClientInfo      clientInfo `json:"clientInfo"`
}

type clientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type toolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type callToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ---- public methods ----

func (c *httpClient) DiscoverTools(ctx context.Context, serverURL string) ([]llm.ToolDef, error) {
	if err := c.ensureSession(ctx, serverURL); err != nil {
		return nil, fmt.Errorf("discover tools: %w", err)
	}

	resp, err := c.doRPC(ctx, serverURL, "tools/list", nil)
	if err != nil {
		return nil, fmt.Errorf("discover tools: %w", err)
	}

	var result struct {
		Tools []toolDef `json:"tools"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("discover tools: unmarshal: %w", err)
	}

	tools := make([]llm.ToolDef, len(result.Tools))
	for i, t := range result.Tools {
		tools[i] = llm.ToolDef{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		}
	}
	return tools, nil
}

func (c *httpClient) CallTool(ctx context.Context, serverURL, toolName string, args json.RawMessage) (string, error) {
	if err := c.ensureSession(ctx, serverURL); err != nil {
		return "", fmt.Errorf("call tool %s: %w", toolName, err)
	}

	resp, err := c.doRPC(ctx, serverURL, "tools/call", callToolParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return "", fmt.Errorf("call tool %s: %w", toolName, err)
	}

	var result struct {
		Content []toolContent `json:"content"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("call tool %s: unmarshal: %w", toolName, err)
	}

	var text string
	for _, c := range result.Content {
		if c.Type == "text" {
			text += c.Text
		}
	}
	return text, nil
}

// ---- session management ----

func (c *httpClient) ensureSession(ctx context.Context, serverURL string) error {
	if c.sessionID != "" {
		return nil
	}

	resp, err := c.doRPC(ctx, serverURL, "initialize", initializeParams{
		ProtocolVersion: "2024-11-05",
		ClientInfo:      clientInfo{Name: "streamcore", Version: "1.0"},
	})
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}
	_ = resp // initialize response is not needed

	return nil
}

// ---- low-level HTTP ----

func (c *httpClient) doRPC(ctx context.Context, serverURL, method string, params any) (json.RawMessage, error) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		ID:      1,
		Params:  params,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")

	if c.sessionID != "" {
		httpReq.Header.Set("Mcp-Session-Id", c.sessionID)
	}

	httpResp, err := c.hc.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer httpResp.Body.Close()

	// Capture session ID from response header
	if sid := httpResp.Header.Get("Mcp-Session-Id"); sid != "" {
		c.sessionID = sid
	}

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", httpResp.StatusCode, string(data))
	}

	body := data
	if strings.HasPrefix(string(data), "event:") || strings.HasPrefix(string(data), "data:") {
		if sseData := extractSSEData(data); sseData != nil {
			body = sseData
		}
	}

	var rpcResp struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return rpcResp.Result, nil
}

// extractSSEData extracts the JSON payload from an SSE event stream.
func extractSSEData(raw []byte) []byte {
	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	var lastData string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			lastData = strings.TrimPrefix(line, "data: ")
		}
	}
	if lastData != "" {
		return []byte(lastData)
	}
	return nil
}
