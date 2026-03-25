package social

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *socialdb) GetFollow(ctx context.Context, follower, followee uint) (*domain.Follow, error) {
	var po model.FollowModel
	err := repo.db.Model(&model.FollowModel{}).
		Where("follower_id = ? AND followee_id = ?", follower, followee).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.Follow(&po), nil
}
