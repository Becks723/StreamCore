package service

import (
	"context"

	"StreamCore/internal/pkg/ai/agent"
	"StreamCore/internal/pkg/ai/mcp"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/db/ai"
	"StreamCore/kitex_gen/chat/chatservice"
	"github.com/redis/go-redis/v9"
)

type AIService struct {
	ctx        context.Context
	db         ai.AIDatabase
	infra      *base.InfraSet
	chatClient chatservice.Client
	agent      *agent.Agent
	toolReg    *mcp.ToolRegistry
	rdb        *redis.Client
}

func NewAIService(ctx context.Context, infra *base.InfraSet, ag *agent.Agent, tr *mcp.ToolRegistry) *AIService {
	return &AIService{
		ctx:        ctx,
		db:         infra.DB.AI,
		infra:      infra,
		chatClient: infra.ChatClient,
		agent:      ag,
		toolReg:    tr,
		rdb:        infra.RDB,
	}
}
