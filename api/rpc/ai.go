package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/ai"
	"StreamCore/kitex_gen/ai/aiservice"
)

func initAIRPC() {
	c, err := initRPCClient(constants.AIServiceName, aiservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init ai rpc client: %v", err)
	}
	aiClient = *c
}

// --- Bot CRUD ---

func CreateBotRPC(ctx context.Context, req *ai.CreateBotReq) (*ai.CreateBotResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.CreateBot(ctx, req)
}

func UpdateBotRPC(ctx context.Context, req *ai.UpdateBotReq) (*ai.UpdateBotResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.UpdateBot(ctx, req)
}

func ListBotsRPC(ctx context.Context, req *ai.ListBotsReq) (*ai.ListBotsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListBots(ctx, req)
}

func GetBotRPC(ctx context.Context, req *ai.GetBotReq) (*ai.GetBotResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.GetBot(ctx, req)
}

func DeleteBotRPC(ctx context.Context, req *ai.DeleteBotReq) (*ai.DeleteBotResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.DeleteBot(ctx, req)
}

// --- Bot-Group binding ---

func AddBotToGroupRPC(ctx context.Context, req *ai.AddBotToGroupReq) (*ai.AddBotToGroupResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.AddBotToGroup(ctx, req)
}

func RemoveBotFromGroupRPC(ctx context.Context, req *ai.RemoveBotFromGroupReq) (*ai.RemoveBotFromGroupResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.RemoveBotFromGroup(ctx, req)
}

func ListGroupBotsRPC(ctx context.Context, req *ai.ListGroupBotsReq) (*ai.ListGroupBotsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListGroupBots(ctx, req)
}

func ListBotGroupsRPC(ctx context.Context, req *ai.ListBotGroupsReq) (*ai.ListBotGroupsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListBotGroups(ctx, req)
}

// --- MCP Server / Tool management ---

func RegisterMCPServerRPC(ctx context.Context, req *ai.RegisterMCPServerReq) (*ai.RegisterMCPServerResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.RegisterMCPServer(ctx, req)
}

func RefreshMCPServerRPC(ctx context.Context, req *ai.RefreshMCPServerReq) (*ai.RefreshMCPServerResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.RefreshMCPServer(ctx, req)
}

func ListMCPServersRPC(ctx context.Context, req *ai.ListMCPServersReq) (*ai.ListMCPServersResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListMCPServers(ctx, req)
}

func DeleteMCPServerRPC(ctx context.Context, req *ai.DeleteMCPServerReq) (*ai.DeleteMCPServerResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.DeleteMCPServer(ctx, req)
}

func ListToolsRPC(ctx context.Context, req *ai.ListToolsReq) (*ai.ListToolsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListTools(ctx, req)
}

// --- Credentials ---

func SaveCredentialRPC(ctx context.Context, req *ai.SaveCredentialReq) (*ai.SaveCredentialResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.SaveCredential(ctx, req)
}

func DeleteCredentialRPC(ctx context.Context, req *ai.DeleteCredentialReq) (*ai.DeleteCredentialResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.DeleteCredential(ctx, req)
}

func ListCredentialsRPC(ctx context.Context, req *ai.ListCredentialsReq) (*ai.ListCredentialsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	return aiClient.ListCredentials(ctx, req)
}

// --- Message processing ---

func ProcessMessageRPC(ctx context.Context, req *ai.ProcessMessageReq) (*ai.ProcessMessageResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	resp, err := aiClient.ProcessMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("process message rpc call failed: %w", err)
	}
	return resp, nil
}
