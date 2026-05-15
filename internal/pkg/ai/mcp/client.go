package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"StreamCore/internal/pkg/ai/llm"
)

// Client communicates with MCP servers over HTTP.
type Client interface {
	// DiscoverTools fetches the tool list from an MCP server.
	DiscoverTools(ctx context.Context, serverURL string) ([]llm.ToolDef, error)
	// CallTool invokes a specific tool on an MCP server.
	CallTool(ctx context.Context, serverURL, toolName string, args json.RawMessage) (string, error)
}

type httpClient struct {
	hc *http.Client
}

// NewClient creates an MCP Client with the given timeout.
func NewClient(timeout time.Duration) Client {
	return &httpClient{
		hc: &http.Client{Timeout: timeout},
	}
}

type listToolsRequest struct {
	Method string `json:"method"` // "tools/list"
}

type listToolsResponse struct {
	Tools []toolDef `json:"tools"`
}

type toolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

func (c *httpClient) DiscoverTools(ctx context.Context, serverURL string) ([]llm.ToolDef, error) {
	body := listToolsRequest{Method: "tools/list"}
	resp, err := c.doRequest(ctx, serverURL, body)
	if err != nil {
		return nil, fmt.Errorf("discover tools: %w", err)
	}

	var result listToolsResponse
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

type callToolRequest struct {
	Method    string          `json:"method"` // "tools/call"
	ToolName  string          `json:"toolName"`
	Arguments json.RawMessage `json:"arguments"`
}

type callToolResponse struct {
	Content []toolContent `json:"content"`
}

type toolContent struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

func (c *httpClient) CallTool(ctx context.Context, serverURL, toolName string, args json.RawMessage) (string, error) {
	body := callToolRequest{
		Method:    "tools/call",
		ToolName:  toolName,
		Arguments: args,
	}
	resp, err := c.doRequest(ctx, serverURL, body)
	if err != nil {
		return "", fmt.Errorf("call tool %s: %w", toolName, err)
	}

	var result callToolResponse
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

func (c *httpClient) doRequest(ctx context.Context, serverURL string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}
