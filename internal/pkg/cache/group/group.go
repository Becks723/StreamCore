package group

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GroupCache interface {
	Ping(ctx context.Context) error
}

func NewGroupCache(rdb *redis.Client) GroupCache {
	return &groupcache{
		rdb: rdb,
	}
}

type groupcache struct {
	rdb *redis.Client
}

func (c *groupcache) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}
