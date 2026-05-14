package llm

import (
	"context"
	"encoding/json"
)

// Message represents a single turn in the conversation.
type Message struct {
	Role      string // "user", "assistant", "tool_result"
	Content   string
	ToolUseID string // set when Role == "tool_result"
}

// ToolDef describes a tool the LLM can call.
type ToolDef struct {
	Name        string
	Description string
	InputSchema json.RawMessage
}

// ToolCall represents the LLM's decision to call a tool.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string // JSON string
}

// ChatParams bundles the parameters for a single LLM call.
type ChatParams struct {
	ModelName    string // e.g., "claude-sonnet-4-6"
	SystemPrompt string
	Messages     []Message
	Tools        []ToolDef
	MaxTokens    int64
	Temperature  float64
}

// ChatResult is the LLM's response.
// If StopReason is "tool_use", ToolCalls is non-empty and Text may be empty.
// If StopReason is "end_turn", Text contains the final response.
type ChatResult struct {
	Text      string
	ToolCalls []ToolCall
	StopUsage int64 // total input tokens used (for cost tracking)
}

// Client is the LLM provider interface.
type Client interface {
	Chat(ctx context.Context, params *ChatParams) (*ChatResult, error)
}
