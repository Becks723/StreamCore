package service

import (
	"errors"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	groupdb "StreamCore/internal/pkg/db/group"
	"StreamCore/internal/pkg/domain"
)

func (s *GroupService) CreateGroup(uid uint, req *domain.GroupCreateReq) (*domain.GroupCreateResp, error) {
	if req.Name == "" {
		return nil, errors.New("group name is required")
	}

	groupId, err := s.db.CreateGroup(s.ctx, uid, req.Name)
	if err != nil {
		return nil, fmt.Errorf("db.CreateGroup failed: %w", err)
	}
	if err = s.db.AddGroupMember(s.ctx, groupId, uid, constants.ChatGroupRoleOwner); err != nil {
		return nil, fmt.Errorf("db.AddGroupMember(owner) failed: %w", err)
	}

	return &domain.GroupCreateResp{GroupId: groupId}, nil
}

func (s *GroupService) ApplyJoinGroup(uid uint, req *domain.GroupApplyReq) (*domain.GroupApplyResp, error) {
	if req.GroupId == 0 {
		return nil, errors.New("group_id is required")
	}

	exists, err := s.db.GroupExists(s.ctx, req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("db.GroupExists failed: %w", err)
	}
	if !exists {
		return nil, errors.New("group does not exist")
	}

	isMember, err := s.db.IsGroupMember(s.ctx, req.GroupId, uid)
	if err != nil {
		return nil, fmt.Errorf("db.IsGroupMember failed: %w", err)
	}
	if isMember {
		return nil, errors.New("user is already a group member")
	}

	hasPending, err := s.db.HasPendingGroupApply(s.ctx, req.GroupId, uid)
	if err != nil {
		return nil, fmt.Errorf("db.HasPendingGroupApply failed: %w", err)
	}
	if hasPending {
		return nil, errors.New("there is already a pending join request")
	}

	applyId, err := s.db.CreateGroupApply(s.ctx, req.GroupId, uid, req.Reason, time.Now().UnixMilli())
	if err != nil {
		return nil, fmt.Errorf("db.CreateGroupApply failed: %w", err)
	}
	return &domain.GroupApplyResp{ApplyId: applyId}, nil
}

func (s *GroupService) IsGroupMember(groupId, uid uint) (bool, error) {
	if groupId == 0 {
		return false, errors.New("group_id is required")
	}
	if uid == 0 {
		return false, errors.New("user_id is required")
	}

	isMember, err := s.db.IsGroupMember(s.ctx, groupId, uid)
	if err != nil {
		return false, fmt.Errorf("db.IsGroupMember failed: %w", err)
	}
	return isMember, nil
}

func (s *GroupService) ListGroupMemberIds(groupId uint) ([]uint, error) {
	if groupId == 0 {
		return nil, errors.New("group_id is required")
	}

	memberUids, err := s.db.ListGroupMemberIds(s.ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("db.ListGroupMemberIds failed: %w", err)
	}
	return memberUids, nil
}

func (s *GroupService) RespondGroupApply(uid uint, req *domain.GroupApplyRespondReq) (*domain.GroupApplyRespondResp, error) {
	if req.ApplyId <= 0 {
		return nil, errors.New("apply_id is required")
	}
	if req.Action != constants.ChatGroupApplyActionApprove && req.Action != constants.ChatGroupApplyActionReject {
		return nil, errors.New("action must be approve or reject")
	}

	apply, err := s.db.GetGroupApplyByID(s.ctx, req.ApplyId)
	if err != nil {
		return nil, fmt.Errorf("db.GetGroupApplyByID failed: %w", err)
	}
	if apply.Status != constants.ChatGroupApplyStatusPending {
		return nil, errors.New("apply is not pending")
	}

	role, isMember, err := s.db.GetGroupMemberRole(s.ctx, apply.GroupId, uid)
	if err != nil {
		return nil, fmt.Errorf("db.GetGroupMemberRole failed: %w", err)
	}
	if !isMember || role != constants.ChatGroupRoleOwner {
		return nil, errors.New("only group owner can respond to join requests")
	}

	newStatus := constants.ChatGroupApplyStatusRejected
	if req.Action == constants.ChatGroupApplyActionApprove {
		newStatus = constants.ChatGroupApplyStatusApproved
	}

	ts := time.Now().UnixMilli()
	if err = s.db.RunInTx(s.ctx, func(txdb groupdb.GroupDatabase) error {
		if newStatus == constants.ChatGroupApplyStatusApproved {
			if err := txdb.AddGroupMember(s.ctx, apply.GroupId, apply.ApplicantUid, constants.ChatGroupRoleMember); err != nil {
				return fmt.Errorf("db.AddGroupMember failed: %w", err)
			}
		}
		if err := txdb.UpdateGroupApplyStatus(s.ctx, apply.Id, newStatus, ts); err != nil {
			return fmt.Errorf("db.UpdateGroupApplyStatus failed: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &domain.GroupApplyRespondResp{
		ApplyId: apply.Id,
		Status:  newStatus,
	}, nil
}
