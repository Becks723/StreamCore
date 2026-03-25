package infra

import (
	"context"
	"fmt"

	"StreamCore/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	conf := config.Instance().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return rdb, nil
}
