package user

import (
	"context"
	"fmt"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/pack"
	"StreamCore/internal/user/service"
	"StreamCore/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	infra *base.InfraSet
}

func NewUserHandler(infra *base.InfraSet) user.UserService {
	return &UserServiceImpl{
		infra: infra,
	}
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp = new(user.RegisterResp)

	err = service.NewUserService(ctx, s.infra).Register(req.Username, req.Password)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp = new(user.LoginResp)

	data, auth, token, err := service.NewUserService(ctx, s.infra).Login(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
		resp.Auth = auth
		resp.Token = token
	}
	return resp, nil
}

// GetInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetInfo(ctx context.Context, req *user.InfoQuery) (resp *user.InfoResp, err error) {
	resp = new(user.InfoResp)

	data, err := service.NewUserService(ctx, s.infra).GetInfo(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// UploadAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, req *user.AvatarReq) (resp *user.AvatarResp, err error) {
	resp = new(user.AvatarResp)

	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserService.UploadAvatar: get login uid failed: %w", err)
	}

	data, err := service.NewUserService(ctx, s.infra).UploadAvatar(uid, req.Data)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// RefreshToken implements the UserServiceImpl interface.
func (s *UserServiceImpl) RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (resp *user.RefreshTokenResp, err error) {
	resp = new(user.RefreshTokenResp)

	data, err := service.NewUserService(ctx, s.infra).RefreshToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// MFAQrcode implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFAQrcode(ctx context.Context, req *user.MFAQrcodeReq) (resp *user.MFAQrcodeResp, err error) {
	resp = new(user.MFAQrcodeResp)

	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserService.MFAQrcode: get login uid failed: %w", err)
	}

	data, err := service.NewUserService(ctx, s.infra).MFAQrcode(uid)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// MFABind implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFABind(ctx context.Context, req *user.MFABindReq) (resp *user.MFABindResp, err error) {
	resp = new(user.MFABindResp)

	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserService.MFABind: get login uid failed: %w", err)
	}

	err = service.NewUserService(ctx, s.infra).MFABind(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// MFAVerify implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFAVerify(ctx context.Context, req *user.MFAVerifyReq) (resp *user.MFAVerifyResp, err error) {
	resp = new(user.MFAVerifyResp)

	data, err := service.NewUserService(ctx, s.infra).MFAVerify(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildSuccessResp()
	resp.Data = data
	return resp, nil
}
