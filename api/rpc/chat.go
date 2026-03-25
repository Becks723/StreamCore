package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/chat"
	"StreamCore/kitex_gen/chat/chatservice"
)

func initChatRPC() {
	c, err := initRPCClient(constants.ChatServiceName, chatservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init chat rpc client: %v", err)
	}
	chatClient = *c
}

func SendWhisperMessageRPC(ctx context.Context, req *chat.WhisperClientMsg) (*chat.WhisperServerMsg, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.SendWhisperMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("send whisper message rpc call failed: %w", err)
	}
	return resp, nil
}

func SendGroupMessageRPC(ctx context.Context, req *chat.GroupClientMsg) (*chat.GroupServerMsg, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.SendGroupMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("send group message rpc call failed: %w", err)
	}
	return resp, nil
}

func ListWhisperMessagesRPC(ctx context.Context, query *chat.ListWhisperMessagesQuery) (*chat.ListWhisperMessagesResp, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.ListWhisperMessages(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list whisper messages rpc call failed: %w", err)
	}
	return resp, nil
}

func ListWhisperMessagesAllRPC(ctx context.Context, query *chat.ListWhisperMessagesAllQuery) (*chat.ListWhisperMessagesAllResp, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.ListWhisperMessagesAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all whisper messages rpc call failed: %w", err)
	}
	return resp, nil
}

func ListGroupMessagesRPC(ctx context.Context, query *chat.ListGroupMessagesQuery) (*chat.ListGroupMessagesResp, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.ListGroupMessages(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list group messages rpc call failed: %w", err)
	}
	return resp, nil
}

func ListGroupMessagesAllRPC(ctx context.Context, query *chat.ListGroupMessagesAllQuery) (*chat.ListGroupMessagesAllResp, error) {
	if chatClient == nil {
		return nil, errors.New("chat rpc client not initialized")
	}
	resp, err := chatClient.ListGroupMessagesAll(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all group messages rpc call failed: %w", err)
	}
	return resp, nil
}
