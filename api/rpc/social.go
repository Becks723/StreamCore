package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/social"
	"StreamCore/kitex_gen/social/socialservice"
)

func initSocialRPC() {
	c, err := initRPCClient(constants.SocialServiceName, socialservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init social rpc client: %v", err)
	}
	socialClient = *c
}

func FollowRPC(ctx context.Context, req *social.FollowReq) (*social.FollowResp, error) {
	if socialClient == nil {
		return nil, errors.New("social rpc client not initialized")
	}

	resp, err := socialClient.Follow(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("follow rpc call failed: %w", err)
	}

	return resp, nil
}

func ListFollowsRPC(ctx context.Context, req *social.ListFollowsQuery) (*social.ListFollowsResp, error) {
	if socialClient == nil {
		return nil, errors.New("social rpc client not initialized")
	}

	resp, err := socialClient.ListFollows(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list follows rpc call failed: %w", err)
	}

	return resp, nil
}

func ListFollowersRPC(ctx context.Context, req *social.ListFollowersQuery) (*social.ListFollowersResp, error) {
	if socialClient == nil {
		return nil, errors.New("social rpc client not initialized")
	}

	resp, err := socialClient.ListFollowers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list followers rpc call failed: %w", err)
	}

	return resp, nil
}

func ListFriendsRPC(ctx context.Context, req *social.ListFriendsQuery) (*social.ListFriendsResp, error) {
	if socialClient == nil {
		return nil, errors.New("social rpc client not initialized")
	}

	resp, err := socialClient.ListFriends(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list friends rpc call failed: %w", err)
	}

	return resp, nil
}
