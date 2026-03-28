package interaction

import (
	"context"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/domain"
	"StreamCore/pkg/util"
	"github.com/redis/go-redis/v9"
)

type InteractionCache interface {
	OnLiked(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error
	OnUnliked(ctx context.Context, tarType int, uid, tarId uint) error
	GetUserLikedVideosRange(ctx context.Context, uid uint, start, stop int64) ([]uint, error)
	SetUserLikedVideos(ctx context.Context, uid uint, likes []*domain.Like) error
}

func NewInteractionCache(rdb *redis.Client) InteractionCache {
	return &iacache{
		rdb: rdb,
	}
}

func (c *iacache) OnLiked(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error {
	err := c.rdb.Incr(ctx, c.likeCountKey(tarType, tarId)).Err()
	if err != nil {
		return err
	}

	zaddIfExists := redis.NewScript(`
		if redis.call("EXISTS", KEYS[1]) == 1 then
			redis.call("ZADD", KEYS[1], ARGV[1], ARGV[2])
			local size = redis.call("ZCARD", KEYS[1])
			local limit = tonumber(ARGV[3])
			if size > limit then
				redis.call("ZREMRANGEBYRANK", KEYS[1], 0, size - limit - 1)
			end
			return 1
		else
			return 0
		end
	`)
	err = zaddIfExists.Run(ctx, c.rdb, []string{c.userLikesKey(tarType, uid)}, time.UnixMilli(), tarId, constants.UserLikesCacheLimit).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *iacache) OnUnliked(ctx context.Context, tarType int, uid, tarId uint) error {
	err := c.rdb.Incr(ctx, c.unlikeCountKey(tarType, tarId)).Err()
	if err != nil {
		return err
	}

	err = c.rdb.ZRem(ctx, c.userLikesKey(tarType, uid), tarId).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *iacache) GetUserLikedVideosRange(ctx context.Context, uid uint, start, stop int64) ([]uint, error) {
	members, err := c.rdb.ZRange(ctx, c.userLikesKey(constants.LikeTarType_Video, uid), start, stop).Result()
	if err != nil {
		return nil, err
	}

	vids := make([]uint, len(members))
	for i, m := range members {
		vids[i] = util.String2Uint(m)
	}
	return vids, nil
}

func (c *iacache) SetUserLikedVideos(ctx context.Context, uid uint, likes []*domain.Like) error {
	key := c.userLikesKey(constants.LikeTarType_Video, uid)
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}

	if len(likes) == 0 {
		return nil
	}

	zMembers := make([]redis.Z, len(likes))
	for i, like := range likes {
		zMembers[i] = redis.Z{
			Score:  float64(like.Time.UnixMilli()),
			Member: like.TargetId,
		}
	}
	if err := c.rdb.ZAdd(ctx, key, zMembers...).Err(); err != nil {
		return err
	}
	return nil
}

func (c *iacache) likeCountKey(tarType int, tarId uint) string {
	return fmt.Sprintf("like_count:%d:%d", tarId, tarType)
}

func (c *iacache) unlikeCountKey(tarType int, tarId uint) string {
	return fmt.Sprintf("unlike_count:%d:%d", tarId, tarType)
}

func (c *iacache) userLikesKey(tarType int, uid uint) string {
	return fmt.Sprintf("user_likes:%d:%d", uid, tarType)
}

type iacache struct {
	rdb *redis.Client
}
