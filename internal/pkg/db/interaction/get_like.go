package interaction

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) GetLike(ctx context.Context, tarType int, uid, tarId uint) (*domain.Like, error) {
	var po model.LikeRelationModel
	err := repo.db.Model(&model.LikeRelationModel{}).
		Where("target_type = ? AND uid = ? AND target_id = ?", tarType, uid, tarId).
		Find(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.Like(&po), nil
}
