package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/video"
	"StreamCore/kitex_gen/video/videoservice"
)

func initVideoRPC() {
	c, err := initRPCClient(constants.VideoServiceName, videoservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init video rpc client: %v", err)
	}
	videoClient = *c
}

func FeedRPC(ctx context.Context, req *video.FeedQuery) (*video.FeedResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.Feed(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("feed rpc call failed: %w", err)
	}

	return resp, nil
}

func PublishRPC(ctx context.Context, req *video.PublishReq) (*video.PublishResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.Publish(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("publish rpc call failed: %w", err)
	}

	return resp, nil
}

func ListRPC(ctx context.Context, req *video.ListQuery) (*video.ListResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list rpc call failed: %w", err)
	}

	return resp, nil
}

func PopularRPC(ctx context.Context, req *video.PopularQuery) (*video.PopularResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.Popular(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("popular rpc call failed: %w", err)
	}

	return resp, nil
}

func SearchRPC(ctx context.Context, req *video.SearchReq) (*video.SearchResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search rpc call failed: %w", err)
	}

	return resp, nil
}

func VisitRPC(ctx context.Context, req *video.VisitQuery) (*video.VisitResp, error) {
	if videoClient == nil {
		return nil, errors.New("video rpc client not initialized")
	}

	resp, err := videoClient.Visit(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("visit rpc call failed: %w", err)
	}

	return resp, nil
}
