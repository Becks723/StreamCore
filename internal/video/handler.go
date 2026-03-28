package video

import (
	"context"
	"fmt"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/pack"
	"StreamCore/internal/video/service"
	"StreamCore/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	infra *base.InfraSet
}

func NewVideoHandler(infra *base.InfraSet) video.VideoService {
	return &VideoServiceImpl{
		infra: infra,
	}
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedQuery) (resp *video.FeedResp, err error) {
	resp = new(video.FeedResp)

	data, err := service.NewVideoService(ctx, s.infra).GetVideoFeed(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// Publish implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Publish(ctx context.Context, req *video.PublishReq) (resp *video.PublishResp, err error) {
	resp = new(video.PublishResp)

	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("VideoService.Feed: get login uid failed: %w", err)
	}

	err = service.NewVideoService(ctx, s.infra).Publish(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// List implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) List(ctx context.Context, req *video.ListQuery) (resp *video.ListResp, err error) {
	resp = new(video.ListResp)

	data, err := service.NewVideoService(ctx, s.infra).List(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// Popular implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Popular(ctx context.Context, req *video.PopularQuery) (resp *video.PopularResp, err error) {
	resp = new(video.PopularResp)

	data, err := service.NewVideoService(ctx, s.infra).Popular(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// Search implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Search(ctx context.Context, req *video.SearchReq) (resp *video.SearchResp, err error) {
	resp = new(video.SearchResp)

	data, err := service.NewVideoService(ctx, s.infra).SearchES(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// Visit implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Visit(ctx context.Context, req *video.VisitQuery) (resp *video.VisitResp, err error) {
	resp = new(video.VisitResp)

	var uidOptional *uint
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		uidOptional = nil
	} else {
		uidOptional = &uid
	}
	data, err := service.NewVideoService(ctx, s.infra).Visit(uidOptional, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}
