package service

import (
	"errors"
	"fmt"
	"time"

	"StreamCore/internal/pkg/domain"
)

func (s *ChatService) SendWhisperMessage(uid uint, req *domain.WhisperMessageReq) (*domain.WhisperMessage, error) {
	if req.ToUid == 0 {
		return nil, errors.New("to_uid is required")
	}
	if req.ToUid == uid {
		return nil, errors.New("cannot send private message to yourself")
	}
	if req.Payload == "" {
		return nil, errors.New("payload is required")
	}

	isFriend, err := s.db.IsFriend(s.ctx, uid, req.ToUid)
	if err != nil {
		return nil, fmt.Errorf("db.IsFriend failed: %w", err)
	}
	if !isFriend {
		return nil, errors.New("private messaging requires mutual follow (friend)")
	}

	ts := req.Timestamp
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}

	msg := &domain.WhisperMessage{
		FromUid:   uid,
		ToUid:     req.ToUid,
		Payload:   req.Payload,
		Timestamp: ts,
	}
	if err = s.db.CreateWhisperMessage(s.ctx, msg); err != nil {
		return nil, fmt.Errorf("db.CreateWhisperMessage failed: %w", err)
	}
	return msg, nil
}

func (s *ChatService) SendGroupMessage(uid uint, req *domain.GroupMessageReq) (*domain.GroupMessage, error) {
	if req.GroupId == 0 {
		return nil, errors.New("group_id is required")
	}
	if req.Payload == "" {
		return nil, errors.New("payload is required")
	}

	isMember, err := s.isGroupMember(req.GroupId, uid)
	if err != nil {
		return nil, fmt.Errorf("s.isGroupMember failed: %w", err)
	}
	if !isMember {
		return nil, errors.New("sender is not a group member")
	}

	ts := req.Timestamp
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}

	memberUids, err := s.listGroupMemberIDs(req.GroupId)
	if err != nil {
		return nil, fmt.Errorf("s.listGroupMemberIDs failed: %w", err)
	}

	receiverUids := make([]uint, 0, len(memberUids))
	for _, memberUid := range memberUids {
		if memberUid != uid {
			receiverUids = append(receiverUids, memberUid)
		}
	}

	msg := &domain.GroupMessage{
		FromUid:      uid,
		GroupId:      req.GroupId,
		Payload:      req.Payload,
		Timestamp:    ts,
		ReceiverUids: receiverUids,
	}
	if err = s.db.CreateGroupMessage(s.ctx, msg); err != nil {
		return nil, fmt.Errorf("db.CreateGroupMessage failed: %w", err)
	}
	return msg, nil
}
