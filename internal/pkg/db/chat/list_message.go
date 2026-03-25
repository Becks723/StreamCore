package chat

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *chatdb) ListWhisperMessages(ctx context.Context, uid, peer uint, pageSize int, cursorMsgId int64) ([]*domain.WhisperMessage, bool, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}

	base := repo.db.WithContext(ctx).Model(&model.ChatMsgModel{}).
		Where("msg_type = ? AND ((sender_uid = ? AND receiver_uid = ?) OR (sender_uid = ? AND receiver_uid = ?))",
			constants.ChatMsgTypeWhisper,
			uid, peer,
			peer, uid,
		)
	if cursorMsgId > 0 {
		base = base.Where("id < ?", cursorMsgId)
	}

	var records []*model.ChatMsgModel
	if err := base.Order("id DESC").Limit(pageSize + 1).Find(&records).Error; err != nil {
		return nil, false, 0, err
	}

	hasMore := len(records) > pageSize
	if hasMore {
		records = records[:pageSize]
	}

	items := make([]*domain.WhisperMessage, 0, len(records))
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		items = append(items, &domain.WhisperMessage{
			MsgId:     r.Id,
			FromUid:   r.SenderUid,
			ToUid:     r.ReceiverUid,
			Payload:   r.Payload,
			Timestamp: r.Timestamp,
		})
	}

	var nextCursor int64
	if hasMore && len(records) > 0 {
		nextCursor = records[len(records)-1].Id
	}
	return items, hasMore, nextCursor, nil
}

func (repo *chatdb) ListWhisperMessagesAll(ctx context.Context, uid, peer uint) ([]*domain.WhisperMessage, error) {
	var records []*model.ChatMsgModel
	err := repo.db.WithContext(ctx).Model(&model.ChatMsgModel{}).
		Where("msg_type = ? AND ((sender_uid = ? AND receiver_uid = ?) OR (sender_uid = ? AND receiver_uid = ?))",
			constants.ChatMsgTypeWhisper,
			uid, peer,
			peer, uid,
		).
		Order("id ASC").
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WhisperMessage, 0, len(records))
	for _, r := range records {
		items = append(items, &domain.WhisperMessage{
			MsgId:     r.Id,
			FromUid:   r.SenderUid,
			ToUid:     r.ReceiverUid,
			Payload:   r.Payload,
			Timestamp: r.Timestamp,
		})
	}
	return items, nil
}

func (repo *chatdb) ListGroupMessages(ctx context.Context, groupId uint, pageSize int, cursorMsgId int64) ([]*domain.GroupMessage, bool, int64, error) {
	if pageSize <= 0 {
		pageSize = 20
	}

	base := repo.db.WithContext(ctx).Model(&model.ChatMsgModel{}).
		Where("msg_type = ? AND group_id = ?", constants.ChatMsgTypeGroup, groupId)
	if cursorMsgId > 0 {
		base = base.Where("id < ?", cursorMsgId)
	}

	var records []*model.ChatMsgModel
	if err := base.Order("id DESC").Limit(pageSize + 1).Find(&records).Error; err != nil {
		return nil, false, 0, err
	}

	hasMore := len(records) > pageSize
	if hasMore {
		records = records[:pageSize]
	}

	items := make([]*domain.GroupMessage, 0, len(records))
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		items = append(items, &domain.GroupMessage{
			MsgId:     r.Id,
			FromUid:   r.SenderUid,
			GroupId:   r.GroupId,
			Payload:   r.Payload,
			Timestamp: r.Timestamp,
		})
	}

	var nextCursor int64
	if hasMore && len(records) > 0 {
		nextCursor = records[len(records)-1].Id
	}
	return items, hasMore, nextCursor, nil
}

func (repo *chatdb) ListGroupMessagesAll(ctx context.Context, groupId uint) ([]*domain.GroupMessage, error) {
	var records []*model.ChatMsgModel
	err := repo.db.WithContext(ctx).Model(&model.ChatMsgModel{}).
		Where("msg_type = ? AND group_id = ?", constants.ChatMsgTypeGroup, groupId).
		Order("id ASC").
		Find(&records).
		Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.GroupMessage, 0, len(records))
	for _, r := range records {
		items = append(items, &domain.GroupMessage{
			MsgId:     r.Id,
			FromUid:   r.SenderUid,
			GroupId:   r.GroupId,
			Payload:   r.Payload,
			Timestamp: r.Timestamp,
		})
	}
	return items, nil
}
