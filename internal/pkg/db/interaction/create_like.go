package interaction

import (
	"context"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
)

func (repo *iactiondb) CreateLike(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error {
	// TODO: atomic
	repo.db.Model(&model.LikeRelationModel{}).Create(&model.LikeRelationModel{
		Uid:        uid,
		TargetType: tarType,
		TargetId:   tarId,
		Status:     constants.LikeAction_Like,
		Time:       time,
	})
	repo.db.Model(&model.LikeCountModel{}).Create(&model.LikeCountModel{
		TargetType:  tarType,
		TargetId:    tarId,
		LikeCount:   1,
		UnlikeCount: 0,
	})
	return nil
}
