package group

import (
	"context"
	"fmt"
	"strconv"

	"StreamCore/internal/group/service"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/domain"
	pkgpack "StreamCore/internal/pkg/pack"
	kitexgroup "StreamCore/kitex_gen/group"
	"StreamCore/pkg/util"
)

// GroupServiceImpl implements the last service interface defined in the IDL.
type GroupServiceImpl struct {
	infra *base.InfraSet
}

func NewGroupHandler(infra *base.InfraSet) kitexgroup.GroupService {
	return &GroupServiceImpl{infra: infra}
}

func (s *GroupServiceImpl) CreateGroup(ctx context.Context, req *kitexgroup.CreateGroupReq) (*kitexgroup.CreateGroupResp, error) {
	resp := new(kitexgroup.CreateGroupResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("GroupService.CreateGroup: get login uid failed: %w", err)
	}

	data, err := service.NewGroupService(ctx, s.infra).CreateGroup(uid, &domain.GroupCreateReq{Name: req.Name})
	if err != nil {
		resp.Base = pkgpack.BuildBaseResp(err)
	} else {
		resp.Base = pkgpack.BuildSuccessResp()
		resp.Data = pkgpack.CreateGroupData(data)
	}
	return resp, nil
}

func (s *GroupServiceImpl) ApplyJoinGroup(ctx context.Context, req *kitexgroup.ApplyJoinGroupReq) (*kitexgroup.ApplyJoinGroupResp, error) {
	resp := new(kitexgroup.ApplyJoinGroupResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("GroupService.ApplyJoinGroup: get login uid failed: %w", err)
	}

	groupID, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("GroupService.ApplyJoinGroup: bad group_id format: %w", err)
	}

	reason := ""
	if req.IsSetReason() {
		reason = req.GetReason()
	}
	data, err := service.NewGroupService(ctx, s.infra).ApplyJoinGroup(uid, &domain.GroupApplyReq{
		GroupId: groupID,
		Reason:  reason,
	})
	if err != nil {
		resp.Base = pkgpack.BuildBaseResp(err)
	} else {
		resp.Base = pkgpack.BuildSuccessResp()
		resp.Data = pkgpack.ApplyJoinGroupData(data)
	}
	return resp, nil
}

func (s *GroupServiceImpl) IsGroupMember(ctx context.Context, req *kitexgroup.IsGroupMemberReq) (*kitexgroup.IsGroupMemberResp, error) {
	resp := new(kitexgroup.IsGroupMemberResp)

	groupId, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("GroupService.IsGroupMember: bad group_id format: %w", err)
	}
	userId, err := util.ParseUint(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("GroupService.IsGroupMember: bad user_id format: %w", err)
	}

	isMember, err := service.NewGroupService(ctx, s.infra).IsGroupMember(groupId, userId)
	if err != nil {
		resp.Base = pkgpack.BuildBaseResp(err)
	} else {
		resp.Base = pkgpack.BuildSuccessResp()
		resp.Data = pkgpack.IsGroupMemberData(isMember)
	}
	return resp, nil
}

func (s *GroupServiceImpl) ListGroupMemberIds(ctx context.Context, req *kitexgroup.ListGroupMemberIdsReq) (*kitexgroup.ListGroupMemberIdsResp, error) {
	resp := new(kitexgroup.ListGroupMemberIdsResp)

	groupId, err := util.ParseUint(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("GroupService.ListGroupMemberIds: bad group_id format: %w", err)
	}

	memberUids, err := service.NewGroupService(ctx, s.infra).ListGroupMemberIds(groupId)
	if err != nil {
		resp.Base = pkgpack.BuildBaseResp(err)
	} else {
		resp.Base = pkgpack.BuildSuccessResp()
		resp.Data = pkgpack.ListGroupMemberIdsData(memberUids)
	}
	return resp, nil
}

func (s *GroupServiceImpl) RespondGroupApply(ctx context.Context, req *kitexgroup.RespondGroupApplyReq) (*kitexgroup.RespondGroupApplyResp, error) {
	resp := new(kitexgroup.RespondGroupApplyResp)
	uid, err := rpccontext.RetrieveLoginUid(ctx)
	if err != nil {
		return nil, fmt.Errorf("GroupService.RespondGroupApply: get login uid failed: %w", err)
	}

	applyID, err := strconv.ParseInt(req.ApplyId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("GroupService.RespondGroupApply: bad apply_id format: %w", err)
	}

	data, err := service.NewGroupService(ctx, s.infra).RespondGroupApply(uid, &domain.GroupApplyRespondReq{
		ApplyId: applyID,
		Action:  int(req.Action),
	})
	if err != nil {
		resp.Base = pkgpack.BuildBaseResp(err)
	} else {
		resp.Base = pkgpack.BuildSuccessResp()
		resp.Data = pkgpack.RespondGroupApplyData(data)
	}
	return resp, nil
}
