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

func ListGroupBotsRPC(ctx context.Context, req *ai.ListGroupBotsReq) (*ai.ListGroupBotsResp, error) {
	if aiClient == nil {
		return nil, errors.New("ai rpc client not initialized")
	}
	resp, err := aiClient.ListGroupBots(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list group bots rpc call failed: %w", err)
	}
	return resp, nil
}
