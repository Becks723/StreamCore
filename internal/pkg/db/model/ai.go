package model

import "time"

// BotGroupModel represents bot-group binding.
type BotGroupModel struct {
	ID        uint `gorm:"primaryKey"`
	BotID     uint `gorm:"uniqueIndex:idx_bot_group"`
	GroupID   uint `gorm:"uniqueIndex:idx_bot_group;index:idx_group_bot,priority:1"`
	Status    int  `gorm:"index:idx_group_bot,priority:2"` // 0: in, 1: removed
	CreatedAt time.Time
}

// MCPServerModel represents a registered MCP server.
type MCPServerModel struct {
	ID              uint   `gorm:"primaryKey"`
	ServerName      string `gorm:"size:128;uniqueIndex"`
	ServerURL       string
	AuthToken       string
	SyncIntervalSec int        `gorm:"default:300"`
	LastSyncedAt    *time.Time // NULL until first sync
	Status          int        `gorm:"default:1"` // 1: active, 0: disabled
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// MCPToolModel represents a tool discovered from an MCP server.
type MCPToolModel struct {
	ID          uint   `gorm:"primaryKey"`
	ToolID      string `gorm:"size:255;uniqueIndex"`
	ToolName    string
	Description string
	InputSchema string `gorm:"type:text"` // JSON Schema
	ServerID    uint   `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CredentialModel stores encrypted third-party service credentials.
type CredentialModel struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"uniqueIndex:idx_user_service"`
	ServiceName string `gorm:"size:64;uniqueIndex:idx_user_service"`
	Username    string
	Password    string // AES-256-GCM encrypted
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
