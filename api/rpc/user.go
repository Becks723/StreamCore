package rpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/user"
	"StreamCore/kitex_gen/user/userservice"
)

func initUserRPC() {
	c, err := initRPCClient(constants.UserServiceName, userservice.NewClient)
	if err != nil {
		log.Fatalf("failed to init user rpc client: %v", err)
	}
	userClient = *c
}

func RegisterRPC(ctx context.Context, req *user.RegisterReq) (*user.RegisterResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("register rpc call failed: %w", err)
	}

	return resp, nil
}

func LoginRPC(ctx context.Context, req *user.LoginReq) (*user.LoginResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("login rpc call failed: %w", err)
	}

	return resp, nil
}

func GetInfoRPC(ctx context.Context, req *user.InfoQuery) (*user.InfoResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.GetInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get info rpc call failed: %w", err)
	}

	return resp, nil
}

func UploadAvatarRPC(ctx context.Context, req *user.AvatarReq) (*user.AvatarResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.UploadAvatar(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("upload avatar rpc call failed: %w", err)
	}

	return resp, nil
}

func RefreshTokenRPC(ctx context.Context, req *user.RefreshTokenReq) (*user.RefreshTokenResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.RefreshToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("refresh token rpc call failed: %w", err)
	}

	return resp, nil
}

func MFAQrcodeRPC(ctx context.Context, req *user.MFAQrcodeReq) (*user.MFAQrcodeResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.MFAQrcode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("mfa qrcode rpc call failed: %w", err)
	}

	return resp, nil
}

func MFABindRPC(ctx context.Context, req *user.MFABindReq) (*user.MFABindResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.MFABind(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("mfa bind rpc call failed: %w", err)
	}

	return resp, nil
}

func MFAVerifyRPC(ctx context.Context, req *user.MFAVerifyReq) (*user.MFAVerifyResp, error) {
	if userClient == nil {
		return nil, errors.New("user rpc client not initialized")
	}

	resp, err := userClient.MFAVerify(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("mfa verify rpc call failed: %w", err)
	}

	return resp, nil
}
