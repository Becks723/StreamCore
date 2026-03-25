package social

import (
	"context"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
)

func (repo *socialdb) CreateFollow(ctx context.Context, follower, followee uint, time time.Time) error {
	return repo.db.Model(&model.FollowModel{}).Create(&model.FollowModel{
		FollowerId: follower,
		FolloweeId: followee,
		Status:     constants.FollowAction_Follow,
		Time:       time,
	}).Error
}
