package domain

import "time"

// BotConfig represents the JSON stored in users.bot_config.
type BotConfig struct {
	SystemPrompt string   `json:"system_prompt"`
	ModelName    string   `json:"model_name"`
	TriggerMode  int32    `json:"trigger_mode"`
	ToolIDs      []string `json:"tool_ids"`
}

// Bot represents a bot user in the system.
type Bot struct {
	ID           uint
	Username     string
	AvatarURL    string
	SystemPrompt string
	ModelName    string
	TriggerMode  int32
	ToolIDs      []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MCPServer represents a registered MCP server.
type MCPServer struct {
	ID              uint
	ServerName      string
	ServerURL       string
	AuthToken       string
	SyncIntervalSec int32
	LastSyncedAt    time.Time
	Status          int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// MCPTool represents a tool discovered from an MCP server.
type MCPTool struct {
	ID          uint
	ToolID      string
	ToolName    string
	Description string
	InputSchema string
	ServerID    uint
	ServerName  string
}

// Credential represents a user's third-party service credential.
type Credential struct {
	ID          uint
	UserID      uint
	ServiceName string
	Username    string
	Password    string // encrypted
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
