package pack

import (
	"time"

	"StreamCore/internal/pkg/db/model"
	kitexai "StreamCore/kitex_gen/ai"
	kitexcommon "StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
	"github.com/bytedance/sonic"
)

// BotConfigData is used for JSON serialization of bot_config.
type BotConfigData struct {
	SystemPrompt string   `json:"system_prompt"`
	ModelName    string   `json:"model_name"`
	TriggerMode  int32    `json:"trigger_mode"`
	ToolIDs      []string `json:"tool_ids"`
}

// BotConfigToJSON marshals BotConfigData to JSON string.
func BotConfigToJSON(cfg *BotConfigData) string {
	if cfg == nil {
		return "{}"
	}
	if cfg.ToolIDs == nil {
		cfg.ToolIDs = []string{}
	}
	data, _ := sonic.MarshalString(cfg)
	return data
}

// ParseBotConfig parses JSON string to BotConfigData.
func ParseBotConfig(raw string) *BotConfigData {
	cfg := &BotConfigData{}
	if raw == "" {
		return cfg
	}
	_ = sonic.UnmarshalString(raw, cfg)
	if cfg.ToolIDs == nil {
		cfg.ToolIDs = []string{}
	}
	return cfg
}

// BotInfo converts a UserModel (is_bot=1) to thrift BotInfo.
func BotInfo(user *model.UserModel) *kitexcommon.BotInfo {
	cfg := ParseBotConfig(user.BotConfig)
	return &kitexcommon.BotInfo{
		BotId:        util.Uint2String(user.ID),
		BotName:      user.Username,
		AvatarUrl:    user.AvatarUrl,
		Description:  "", // UserModel doesn't have a description field
		SystemPrompt: cfg.SystemPrompt,
		ModelName:    cfg.ModelName,
		TriggerMode:  cfg.TriggerMode,
		ToolIds:      cfg.ToolIDs,
		CreatedAt:    user.CreatedAt.Format(time.DateTime),
		UpdatedAt:    user.UpdatedAt.Format(time.DateTime),
	}
}

// MCPToolInfo converts MCPToolModel to thrift MCPToolInfo.
func MCPToolInfo(tool *model.MCPToolModel, serverName string) *kitexcommon.MCPToolInfo {
	return &kitexcommon.MCPToolInfo{
		ToolId:      tool.ToolID,
		ToolName:    tool.ToolName,
		Description: tool.Description,
		InputSchema: tool.InputSchema,
		ServerName:  serverName,
		ServerId:    int64(tool.ServerID),
	}
}

// MCPServerInfo converts MCPServerModel to thrift MCPServerInfo.
func MCPServerInfo(s *model.MCPServerModel) *kitexcommon.MCPServerInfo {
	return &kitexcommon.MCPServerInfo{
		ServerId:        int64(s.ID),
		ServerName:      s.ServerName,
		ServerUrl:       s.ServerURL,
		SyncIntervalSec: int32(s.SyncIntervalSec),
		LastSyncedAt:    s.LastSyncedAt.Format(time.DateTime),
		Status:          int32(s.Status),
		CreatedAt:       s.CreatedAt.Format(time.DateTime),
		UpdatedAt:       s.UpdatedAt.Format(time.DateTime),
	}
}

// CredentialInfo converts CredentialModel to thrift CredentialInfo.
func CredentialInfo(cred *model.CredentialModel) *kitexai.CredentialInfo {
	return &kitexai.CredentialInfo{
		CredentialId: util.Uint2String(cred.ID),
		ServiceName:  cred.ServiceName,
		Username:     cred.Username,
		SavedAt:      cred.CreatedAt.Format(time.DateTime),
	}
}
