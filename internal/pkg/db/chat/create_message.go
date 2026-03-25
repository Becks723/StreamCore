package chat

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *chatdb) CreateWhisperMessage(ctx context.Context, msg *domain.WhisperMessage) error {
	po := &model.ChatMsgModel{
		SenderUid:     msg.FromUid,
		ReceiverUid:   msg.ToUid,
		Payload:       msg.Payload,
		MsgType:       constants.ChatMsgTypeWhisper,
		Timestamp:     msg.Timestamp,
		DeliverStatus: constants.ChatDeliverStatusDone,
	}
	if err := repo.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	msg.MsgId = po.Id
	return nil
}

func (repo *chatdb) CreateGroupMessage(ctx context.Context, msg *domain.GroupMessage) error {
	po := &model.ChatMsgModel{
		SenderUid:     msg.FromUid,
		GroupId:       msg.GroupId,
		Payload:       msg.Payload,
		MsgType:       constants.ChatMsgTypeGroup,
		Timestamp:     msg.Timestamp,
		DeliverStatus: constants.ChatDeliverStatusDone,
	}
	if err := repo.db.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	msg.MsgId = po.Id
	return nil
}

func (repo *chatdb) IsFriend(ctx context.Context, uid, peer uint) (bool, error) {
	var cnt int64
	err := repo.db.WithContext(ctx).Model(&model.FollowModel{}).
		Where("follower_id = ? AND followee_id = ? AND status = ?", uid, peer, constants.FollowAction_Follow).
		Count(&cnt).
		Error
	if err != nil {
		return false, err
	}
	if cnt == 0 {
		return false, nil
	}

	err = repo.db.WithContext(ctx).Model(&model.FollowModel{}).
		Where("follower_id = ? AND followee_id = ? AND status = ?", peer, uid, constants.FollowAction_Follow).
		Count(&cnt).
		Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
