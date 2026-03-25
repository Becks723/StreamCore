package interaction

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
)

func (repo *iactiondb) ToggleLikeStatus(ctx context.Context, tarType int, uid, tarId uint) error {
	// 1. status like > unlike or (unlike > like)
	// 2. incr like/unlike count

	lastLike, err := repo.GetLike(ctx, tarType, uid, tarId)
	if err != nil { // ignore unrecorded
		return nil
	}
	if lastLike.Status == constants.LikeAction_Like {
		return repo.switchToUnlike(ctx, tarType, uid, tarId)
	} else {
		return repo.switchToLike(ctx, tarType, uid, tarId)
	}
}

func (repo *iactiondb) switchToLike(ctx context.Context, tarType int, uid, tarId uint) error {
	// TODO: atomic
	repo.db.Model(&model.LikeRelationModel{}).
		Where("target_type = ? AND uid = ? AND target_id = ?", tarType, uid, tarId).
		Update("status", constants.LikeAction_Like)
	repo.db.Model(&model.LikeCountModel{}).
		Where("target_type = ? AND target_id = ?", tarType, tarId).
		Update("like_count", gorm.Expr("like_count + ?", 1))
	return nil
}

func (repo *iactiondb) switchToUnlike(ctx context.Context, tarType int, uid, tarId uint) error {
	// TODO: atomic
	// TODO: atomic
	repo.db.Model(&model.LikeRelationModel{}).
		Where("target_type = ? AND uid = ? AND target_id = ?", tarType, uid, tarId).
		Update("status", constants.LikeAction_Unlike)
	repo.db.Model(&model.LikeCountModel{}).
		Where("target_type = ? AND target_id = ?", tarType, tarId).
		Update("unlike_count", gorm.Expr("unlike_count + ?", 1))
	return nil
}
