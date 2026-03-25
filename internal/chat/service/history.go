package service

import (
	"errors"
	"fmt"

	"StreamCore/config"
	"StreamCore/internal/pkg/domain"
)

func (s *ChatService) ListWhisperMessages(uid uint, req *domain.WhisperHistoryQuery) (*domain.WhisperHistory, error) {
	if req.PeerUid == 0 {
		return nil, errors.New("peer_uid is required")
	}

	isFriend, err := s.db.IsFriend(s.ctx, uid, req.PeerUid)
	if err != nil {
		return nil, fmt.Errorf("db.IsFriend failed: %w", err)
	}
	if !isFriend {
		return nil, errors.New("private history requires mutual follow (friend)")
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = config.Instance().General.PageSize
	}

	items, hasMore, nextCursor, err := s.db.ListWhisperMessages(s.ctx, uid, req.PeerUid, pageSize, req.CursorMsgId)
	if err != nil {
		return nil, fmt.Errorf("db.ListWhisperMessages failed: %w", err)
	}
	return &domain.WhisperHistory{
		Items:           items,
		HasMore:         hasMore,
		NextCursorMsgId: nextCursor,
	}, nil
}

func (s *ChatService) ListWhisperMessagesAll(uid uint, peerUid uint) ([]*domain.WhisperMessage, error) {
	if peerUid == 0 {
		return nil, errors.New("peer_uid is required")
	}

	isFriend, err := s.db.IsFriend(s.ctx, uid, peerUid)
	if err != nil {
		return nil, fmt.Errorf("db.IsFriend failed: %w", err)
	}
	if !isFriend {
		return nil, errors.New("private history requires mutual follow (friend)")
	}

	items, err := s.db.ListWhisperMessagesAll(s.ctx, uid, peerUid)
	if err != nil {
		return nil, fmt.Errorf("db.ListWhisperMessagesAll failed: %w", err)
	}
	return items, nil
}

func (s *ChatService) ListGroupMessages(uid uint, req *domain.GroupHistoryQuery) (*domain.GroupHistory, error) {
	if req.GroupId == 0 {
		return nil, errors.New("group_id is required")
	}

	isMember, err := s.isGroupMember(req.GroupId, uid)
	if err != nil {
		return nil, fmt.Errorf("s.isGroupMember failed: %w", err)
	}
	if !isMember {
		return nil, errors.New("group history requires membership")
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = config.Instance().General.PageSize
	}

	items, hasMore, nextCursor, err := s.db.ListGroupMessages(s.ctx, req.GroupId, pageSize, req.CursorMsgId)
	if err != nil {
		return nil, fmt.Errorf("db.ListGroupMessages failed: %w", err)
	}
	return &domain.GroupHistory{
		Items:           items,
		HasMore:         hasMore,
		NextCursorMsgId: nextCursor,
	}, nil
}

func (s *ChatService) ListGroupMessagesAll(uid uint, groupId uint) ([]*domain.GroupMessage, error) {
	if groupId == 0 {
		return nil, errors.New("group_id is required")
	}

	isMember, err := s.isGroupMember(groupId, uid)
	if err != nil {
		return nil, fmt.Errorf("s.isGroupMember failed: %w", err)
	}
	if !isMember {
		return nil, errors.New("group history requires membership")
	}

	items, err := s.db.ListGroupMessagesAll(s.ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("db.ListGroupMessagesAll failed: %w", err)
	}
	return items, nil
}
