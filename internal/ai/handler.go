package ai

import (
	"context"
	"fmt"
	"time"

	"StreamCore/internal/ai/service"
	"StreamCore/internal/pkg/ai/agent"
	"StreamCore/internal/pkg/ai/mcp"
	"StreamCore/internal/pkg/ai/provider"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	kitexai "StreamCore/kitex_gen/ai"
)

type AIServiceImpl struct {
	infra   *base.InfraSet
	agent   *agent.Agent
	toolReg *mcp.ToolRegistry
}

func NewAIHandler(infra *base.InfraSet) kitexai.AIService {
	mcpClient := mcp.NewClient(30 * time.Second)
	toolReg := mcp.NewToolRegistry(infra.DB.AI, mcpClient)

	// Start periodic MCP tool sync
	go mcp.StartPeriodicSync(context.Background(), toolReg, infra.DB.AI)

	return &AIServiceImpl{
		infra:   infra,
		agent:   agent.New(infra.DB.AI, toolReg, provider.GlobalRegistry()),
		toolReg: toolReg,
	}
}

func (s *AIServiceImpl) ProcessMessage(ctx context.Context, req *kitexai.ProcessMessageReq) (*kitexai.ProcessMessageResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ProcessMessage: %w", err)
	}

	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ProcessMessage(svc, req)
}

func (s *AIServiceImpl) CreateBot(ctx context.Context, req *kitexai.CreateBotReq) (*kitexai.CreateBotResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.CreateBot: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.CreateBot(svc, req)
}

func (s *AIServiceImpl) UpdateBot(ctx context.Context, req *kitexai.UpdateBotReq) (*kitexai.UpdateBotResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.UpdateBot: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.UpdateBot(svc, req)
}

func (s *AIServiceImpl) ListBots(ctx context.Context, req *kitexai.ListBotsReq) (*kitexai.ListBotsResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListBots: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListBots(svc, req)
}

func (s *AIServiceImpl) GetBot(ctx context.Context, req *kitexai.GetBotReq) (*kitexai.GetBotResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.GetBot: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.GetBot(svc, req)
}

func (s *AIServiceImpl) DeleteBot(ctx context.Context, req *kitexai.DeleteBotReq) (*kitexai.DeleteBotResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.DeleteBot: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.DeleteBot(svc, req)
}

func (s *AIServiceImpl) AddBotToGroup(ctx context.Context, req *kitexai.AddBotToGroupReq) (*kitexai.AddBotToGroupResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.AddBotToGroup: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.AddBotToGroup(svc, req)
}

func (s *AIServiceImpl) RemoveBotFromGroup(ctx context.Context, req *kitexai.RemoveBotFromGroupReq) (*kitexai.RemoveBotFromGroupResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.RemoveBotFromGroup: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.RemoveBotFromGroup(svc, req)
}

func (s *AIServiceImpl) ListGroupBots(ctx context.Context, req *kitexai.ListGroupBotsReq) (*kitexai.ListGroupBotsResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListGroupBots: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListGroupBots(svc, req)
}

func (s *AIServiceImpl) ListBotGroups(ctx context.Context, req *kitexai.ListBotGroupsReq) (*kitexai.ListBotGroupsResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListBotGroups: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListBotGroups(svc, req)
}

func (s *AIServiceImpl) RegisterMCPServer(ctx context.Context, req *kitexai.RegisterMCPServerReq) (*kitexai.RegisterMCPServerResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.RegisterMCPServer: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.RegisterMCPServer(svc, req)
}

func (s *AIServiceImpl) RefreshMCPServer(ctx context.Context, req *kitexai.RefreshMCPServerReq) (*kitexai.RefreshMCPServerResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.RefreshMCPServer: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.RefreshMCPServer(svc, req)
}

func (s *AIServiceImpl) ListMCPServers(ctx context.Context, req *kitexai.ListMCPServersReq) (*kitexai.ListMCPServersResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListMCPServers: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListMCPServers(svc, req)
}

func (s *AIServiceImpl) DeleteMCPServer(ctx context.Context, req *kitexai.DeleteMCPServerReq) (*kitexai.DeleteMCPServerResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.DeleteMCPServer: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.DeleteMCPServer(svc, req)
}

func (s *AIServiceImpl) ListTools(ctx context.Context, req *kitexai.ListToolsReq) (*kitexai.ListToolsResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListTools: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListTools(svc, req)
}

func (s *AIServiceImpl) SaveCredential(ctx context.Context, req *kitexai.SaveCredentialReq) (*kitexai.SaveCredentialResp, error) {
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.SaveCredential: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.SaveCredential(svc, uid, req)
}

func (s *AIServiceImpl) DeleteCredential(ctx context.Context, req *kitexai.DeleteCredentialReq) (*kitexai.DeleteCredentialResp, error) {
	_, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.DeleteCredential: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.DeleteCredential(svc, req)
}

func (s *AIServiceImpl) ListCredentials(ctx context.Context, req *kitexai.ListCredentialsReq) (*kitexai.ListCredentialsResp, error) {
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("AIService.ListCredentials: %w", err)
	}
	svc := service.NewAIService(ctx, s.infra, s.agent, s.toolReg)
	return service.ListCredentials(svc, uid, req)
}
