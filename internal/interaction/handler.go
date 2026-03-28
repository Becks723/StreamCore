package interaction

import (
	"context"
	"fmt"

	"StreamCore/internal/interaction/service"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/pack"
	ia "StreamCore/kitex_gen/interaction"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct {
	infra *base.InfraSet
}

func NewInteractionHandler(infra *base.InfraSet) ia.InteractionService {
	return &InteractionServiceImpl{
		infra: infra,
	}
}

// PublishLike implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) PublishLike(ctx context.Context, req *ia.PublishLikeReq) (resp *ia.PublishLikeResp, err error) {
	resp = new(ia.PublishLikeResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("InteractionService.PublishLike: get login uid failed: %w", err)
	}

	err = service.NewInteractionService(ctx, s.infra).PublishLike(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// ListLike implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) ListLike(ctx context.Context, req *ia.ListLikeQuery) (resp *ia.ListLikeResp, err error) {
	resp = new(ia.ListLikeResp)

	data, err := service.NewInteractionService(ctx, s.infra).ListLikedVideos(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// PublishComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) PublishComment(ctx context.Context, req *ia.PublishCommentReq) (resp *ia.PublishCommentResp, err error) {
	resp = new(ia.PublishCommentResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("InteractionService.PublishComment: get login uid failed: %w", err)
	}

	err = service.NewInteractionService(ctx, s.infra).PublishComment(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}

// ListComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) ListComment(ctx context.Context, query *ia.ListCommentQuery) (resp *ia.ListCommentResp, err error) {
	resp = new(ia.ListCommentResp)

	data, err := service.NewInteractionService(ctx, s.infra).ListComment(query)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
		resp.Data = data
	}
	return resp, nil
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *ia.DeleteCommentReq) (resp *ia.DeleteCommentResp, err error) {
	resp = new(ia.DeleteCommentResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("InteractionService.DeleteComment: get login uid failed: %w", err)
	}

	err = service.NewInteractionService(ctx, s.infra).DeleteComment(uid, req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
	} else {
		resp.Base = pack.BuildSuccessResp()
	}
	return resp, nil
}
