package social

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
)

func (repo *socialdb) FetchFollows(ctx context.Context, uid uint, limit, page int) ([]uint, int, error) {
	var follows []uint
	var cnt int64
	var err error
	baseQuery := repo.db.Model(&model.FollowModel{}).
		Select("followee_id").
		Where("follower_id = ? AND status = ?", uid, constants.FollowAction_Follow)
	if err = baseQuery.Count(&cnt).Error; err != nil {
		return nil, -1, err
	}
	err = baseQuery.
		Limit(limit).
		Offset(limit * page).
		Scan(&follows).
		Error
	if err != nil {
		return nil, -1, err
	}
	return follows, int(cnt), nil
}

func (repo *socialdb) FetchFollowers(ctx context.Context, uid uint, limit, page int) ([]uint, int, error) {
	var followers []uint
	var cnt int64
	var err error
	baseQuery := repo.db.Model(&model.FollowModel{}).
		Select("follower_id").
		Where("followee_id = ? AND status = ?", uid, constants.FollowAction_Follow)
	if err = baseQuery.Count(&cnt).Error; err != nil {
		return nil, -1, err
	}
	err = baseQuery.
		Limit(limit).
		Offset(limit * page).
		Scan(&followers).
		Error
	if err != nil {
		return nil, -1, err
	}
	return followers, int(cnt), nil
}

func (repo *socialdb) FetchFriends(ctx context.Context, uid uint, limit, page int) ([]uint, int, error) {
	var friends []uint
	var cnt int64
	var err error
	baseQuery := repo.db.Table("follows f1").
		Select("f2.follower_id").
		Joins("JOIN follows f2 ON f1.follower_id = f2.followee_id AND f2.follower_id = f1.followee_id").
		Where("f1.follower_id = ? AND f1.status = ? AND f2.status = ?",
			uid, constants.FollowAction_Follow, constants.FollowAction_Follow)
	if err = baseQuery.Count(&cnt).Error; err != nil {
		return nil, -1, err
	}
	err = baseQuery.
		Limit(limit).
		Offset(limit * page).
		Scan(&friends).
		Error
	if err != nil {
		return nil, -1, err
	}
	return friends, int(cnt), nil
}
