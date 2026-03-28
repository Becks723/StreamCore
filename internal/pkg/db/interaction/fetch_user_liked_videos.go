package interaction

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) FetchUserLikedVideos(ctx context.Context, uid uint, limit, offset int) ([]*domain.Like, error) {
	var records []model.LikeRelationModel
	err := repo.db.Model(&model.LikeRelationModel{}).
		Select("target_id", "time").
		Where("uid = ? AND target_type = ? AND status = ?", uid, constants.LikeTarType_Video, constants.LikeAction_Like).
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Scan(&records).
		Error
	if err != nil {
		return nil, err
	}

	likes := make([]*domain.Like, len(records))
	for i, rec := range records {
		likes[i] = &domain.Like{
			Uid:        uid,
			TargetType: constants.LikeTarType_Video,
			TargetId:   rec.TargetId,
			Status:     constants.LikeAction_Like,
			Time:       rec.Time,
		}
	}
	return likes, nil
}
