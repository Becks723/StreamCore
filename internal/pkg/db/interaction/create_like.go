package interaction

import (
	"context"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
)

func (repo *iactiondb) CreateLike(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.LikeRelationModel{}).Create(&model.LikeRelationModel{
			Uid:        uid,
			TargetType: tarType,
			TargetId:   tarId,
			Status:     constants.LikeAction_Like,
			Time:       time,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.LikeCountModel{}).Create(&model.LikeCountModel{
			TargetType:  tarType,
			TargetId:    tarId,
			LikeCount:   1,
			UnlikeCount: 0,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}
