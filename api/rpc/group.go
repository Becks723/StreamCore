package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/group"
	"StreamCore/kitex_gen/group/groupservice"
)

func initGroupRPC() {
	c, err := initRPCClient(constants.GroupServiceName, groupservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init group rpc client: %v", err)
	}
	groupClient = *c
}

func CreateGroupRPC(ctx context.Context, req *group.CreateGroupReq) (*group.CreateGroupResp, error) {
	if groupClient == nil {
		return nil, errors.New("group rpc client not initialized")
	}
	resp, err := groupClient.CreateGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create group rpc call failed: %w", err)
	}
	return resp, nil
}

func ApplyJoinGroupRPC(ctx context.Context, req *group.ApplyJoinGroupReq) (*group.ApplyJoinGroupResp, error) {
	if groupClient == nil {
		return nil, errors.New("group rpc client not initialized")
	}
	resp, err := groupClient.ApplyJoinGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("apply join group rpc call failed: %w", err)
	}
	return resp, nil
}

func RespondGroupApplyRPC(ctx context.Context, req *group.RespondGroupApplyReq) (*group.RespondGroupApplyResp, error) {
	if groupClient == nil {
		return nil, errors.New("group rpc client not initialized")
	}
	resp, err := groupClient.RespondGroupApply(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("respond group apply rpc call failed: %w", err)
	}
	return resp, nil
}
