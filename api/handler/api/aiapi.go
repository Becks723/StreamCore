package api

import (
	"context"

	aimodel "StreamCore/api/model/ai"
	"StreamCore/api/pack"
	"StreamCore/api/rpc"
	kitexai "StreamCore/kitex_gen/ai"
	"github.com/cloudwego/hertz/pkg/app"
)

// --- Bot CRUD ---

// CreateBot .
// @router /ai/bot [POST]
func CreateBot(ctx context.Context, c *app.RequestContext) {
	var req aimodel.CreateBotReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.CreateBotRPC(ctx, &kitexai.CreateBotReq{
		BotName:      req.BotName,
		SystemPrompt: req.SystemPrompt,
		Provider:     req.Provider,
		ModelName:    req.ModelName,
		TriggerMode:  req.TriggerMode,
		Description:  req.Description,
		AvatarUrl:    req.AvatarURL,
		ToolIds:      req.ToolIds,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// UpdateBot .
// @router /ai/bot [PUT]
func UpdateBot(ctx context.Context, c *app.RequestContext) {
	var req aimodel.UpdateBotReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.UpdateBotRPC(ctx, &kitexai.UpdateBotReq{
		BotId:        req.BotID,
		BotName:      req.BotName,
		SystemPrompt: req.SystemPrompt,
		Provider:     req.Provider,
		ModelName:    req.ModelName,
		TriggerMode:  req.TriggerMode,
		Description:  req.Description,
		AvatarUrl:    req.AvatarURL,
		ToolIds:      req.ToolIds,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// ListBots .
// @router /ai/bots [GET]
func ListBots(ctx context.Context, c *app.RequestContext) {
	var req aimodel.ListBotsReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.ListBotsRPC(ctx, &kitexai.ListBotsReq{
		PageSize: req.PageSize,
		Page:     req.Page,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// GetBot .
// @router /ai/bot/:bot_id [GET]
func GetBot(ctx context.Context, c *app.RequestContext) {
	var req aimodel.GetBotReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.GetBotRPC(ctx, &kitexai.GetBotReq{
		BotId: req.BotID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// DeleteBot .
// @router /ai/bot/:bot_id [DELETE]
func DeleteBot(ctx context.Context, c *app.RequestContext) {
	var req aimodel.DeleteBotReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.DeleteBotRPC(ctx, &kitexai.DeleteBotReq{
		BotId: req.BotID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// --- Bot-Group binding ---

// AddBotToGroup .
// @router /ai/bot/group [POST]
func AddBotToGroup(ctx context.Context, c *app.RequestContext) {
	var req aimodel.AddBotToGroupReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.AddBotToGroupRPC(ctx, &kitexai.AddBotToGroupReq{
		BotId:   req.BotID,
		GroupId: req.GroupID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// RemoveBotFromGroup .
// @router /ai/bot/group [DELETE]
func RemoveBotFromGroup(ctx context.Context, c *app.RequestContext) {
	var req aimodel.RemoveBotFromGroupReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.RemoveBotFromGroupRPC(ctx, &kitexai.RemoveBotFromGroupReq{
		BotId:   req.BotID,
		GroupId: req.GroupID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// ListGroupBots .
// @router /ai/group/:group_id/bots [GET]
func ListGroupBots(ctx context.Context, c *app.RequestContext) {
	var req aimodel.ListGroupBotsReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.ListGroupBotsRPC(ctx, &kitexai.ListGroupBotsReq{
		GroupId: req.GroupID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// --- MCP Server / Tool management ---

// RegisterMCPServer .
// @router /ai/mcp/server [POST]
func RegisterMCPServer(ctx context.Context, c *app.RequestContext) {
	var req aimodel.RegisterMCPServerReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.RegisterMCPServerRPC(ctx, &kitexai.RegisterMCPServerReq{
		ServerName:      req.ServerName,
		ServerUrl:       req.ServerURL,
		AuthToken:       req.AuthToken,
		SyncIntervalSec: req.SyncIntervalSec,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// RefreshMCPServer .
// @router /ai/mcp/server/refresh [POST]
func RefreshMCPServer(ctx context.Context, c *app.RequestContext) {
	var req aimodel.RefreshMCPServerReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.RefreshMCPServerRPC(ctx, &kitexai.RefreshMCPServerReq{
		ServerId: req.ServerID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// ListMCPServers .
// @router /ai/mcp/servers [GET]
func ListMCPServers(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.ListMCPServersRPC(ctx, &kitexai.ListMCPServersReq{})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// DeleteMCPServer .
// @router /ai/mcp/server/:server_id [DELETE]
func DeleteMCPServer(ctx context.Context, c *app.RequestContext) {
	var req aimodel.DeleteMCPServerReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.DeleteMCPServerRPC(ctx, &kitexai.DeleteMCPServerReq{
		ServerId: req.ServerID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// ListTools .
// @router /ai/tools [GET]
func ListTools(ctx context.Context, c *app.RequestContext) {
	resp, err := rpc.ListToolsRPC(ctx, &kitexai.ListToolsReq{})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// --- Credentials ---

// SaveCredential .
// @router /ai/credential [POST]
func SaveCredential(ctx context.Context, c *app.RequestContext) {
	var req aimodel.SaveCredentialReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.SaveCredentialRPC(ctx, &kitexai.SaveCredentialReq{
		ServiceName: req.ServiceName,
		Username:    req.Username,
		Password:    req.Password,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// DeleteCredential .
// @router /ai/credential/:credential_id [DELETE]
func DeleteCredential(ctx context.Context, c *app.RequestContext) {
	var req aimodel.DeleteCredentialReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.DeleteCredentialRPC(ctx, &kitexai.DeleteCredentialReq{
		CredentialId: req.CredentialID,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}

// ListCredentials .
// @router /ai/credentials [GET]
func ListCredentials(ctx context.Context, c *app.RequestContext) {
	var req aimodel.ListCredentialsReq
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespParamError(c, err)
		return
	}
	resp, err := rpc.ListCredentialsRPC(ctx, &kitexai.ListCredentialsReq{
		ServiceName: req.ServiceName,
	})
	if err != nil {
		pack.RespRPCError(c, err)
		return
	}
	if !pack.RespBizError(c, resp.Base) {
		pack.RespWithData(c, resp)
	}
}
