package social

import (
	"context"
	"fmt"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/pack"
	"StreamCore/internal/social/service"
	"StreamCore/kitex_gen/social"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct {
	infra *base.InfraSet
}

func NewSocialHandler(infra *base.InfraSet) social.SocialService {
	return &SocialServiceImpl{
		infra: infra,
	}
}

// Follow implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) Follow(ctx context.Context, req *social.FollowReq) (resp *social.FollowResp, err error) {
	resp = new(social.FollowResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("SocialService.Follow: get login uid failed: %w", err)
	}

	err = service.NewSocialService(ctx, s.infra).Follow(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// ListFollows implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFollows(ctx context.Context, req *social.ListFollowsQuery) (resp *social.ListFollowsResp, err error) {
	resp = new(social.ListFollowsResp)

	data, err := service.NewSocialService(ctx, s.infra).ListFollows(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// ListFollowers implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFollowers(ctx context.Context, req *social.ListFollowersQuery) (resp *social.ListFollowersResp, err error) {
	resp = new(social.ListFollowersResp)

	data, err := service.NewSocialService(ctx, s.infra).ListFollowers(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// ListFriends implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) ListFriends(ctx context.Context, req *social.ListFriendsQuery) (resp *social.ListFriendsResp, err error) {
	resp = new(social.ListFriendsResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("SocialService.ListFriends: get login uid failed: %w", err)
	}

	data, err := service.NewSocialService(ctx, s.infra).ListFriends(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}
