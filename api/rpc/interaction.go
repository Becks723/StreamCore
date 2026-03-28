package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/kitex_gen/interaction/interactionservice"
)

func initInteractionRPC() {
	c, err := initRPCClient(constants.InteractionServiceName, interactionservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init interaction rpc client: %v", err)
	}
	iaClient = *c
}

func PublishLikeRPC(ctx context.Context, req *interaction.PublishLikeReq) (*interaction.PublishLikeResp, error) {
	if iaClient == nil {
		return nil, errors.New("interaction rpc client not initialized")
	}

	resp, err := iaClient.PublishLike(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("publish like rpc call failed: %w", err)
	}

	return resp, nil
}

func ListLikeRPC(ctx context.Context, req *interaction.ListLikeQuery) (*interaction.ListLikeResp, error) {
	if iaClient == nil {
		return nil, errors.New("interaction rpc client not initialized")
	}

	resp, err := iaClient.ListLike(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list like rpc call failed: %w", err)
	}

	return resp, nil
}

func PublishCommentRPC(ctx context.Context, req *interaction.PublishCommentReq) (*interaction.PublishCommentResp, error) {
	if iaClient == nil {
		return nil, errors.New("interaction rpc client not initialized")
	}

	resp, err := iaClient.PublishComment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("publish comment rpc call failed: %w", err)
	}

	return resp, nil
}

func ListCommentRPC(ctx context.Context, req *interaction.ListCommentQuery) (*interaction.ListCommentResp, error) {
	if iaClient == nil {
		return nil, errors.New("interaction rpc client not initialized")
	}

	resp, err := iaClient.ListComment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list comment rpc call failed: %w", err)
	}

	return resp, nil
}

func DeleteCommentRPC(ctx context.Context, req *interaction.DeleteCommentReq) (*interaction.DeleteCommentResp, error) {
	if iaClient == nil {
		return nil, errors.New("interaction rpc client not initialized")
	}

	resp, err := iaClient.DeleteComment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("delete comment rpc call failed: %w", err)
	}

	return resp, nil
}
