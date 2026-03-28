package social

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SocialCache interface {
	GetFollows(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
	SetFollows(ctx context.Context, uid uint, limit, page int, follows []uint, total int) error
	GetFollowers(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
	SetFollowers(ctx context.Context, uid uint, limit, page int, followers []uint, total int) error
	GetFriends(ctx context.Context, uid uint, limit, page int) ([]uint, int, error)
	SetFriends(ctx context.Context, uid uint, limit, page int, friends []uint, total int) error
	InvalidateUserCache(ctx context.Context, uid uint) error
}

func NewSocialCache(rdb *redis.Client) SocialCache {
	return &socialcache{
		rdb: rdb,
	}
}

type socialcache struct {
	rdb *redis.Client
}
