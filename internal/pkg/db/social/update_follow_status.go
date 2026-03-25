package social

import (
	"context"
	"time"

	"StreamCore/internal/pkg/db/model"
)

func (repo *socialdb) UpdateFollowStatus(ctx context.Context, follower, followee uint, status int, time time.Time) error {
	return repo.db.Model(&model.FollowModel{}).
		Where("follower_id = ? AND followee_id = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"time":   time,
		}).
		Error
}
