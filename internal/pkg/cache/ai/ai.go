package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AICache interface {
	SetProactivePending(ctx context.Context, groupID uint, msgID int64, ttl time.Duration) error
	GetProactivePending(ctx context.Context, groupID uint) (int64, bool, error)
	IsProactiveCoolingDown(ctx context.Context, groupID, botID uint) (bool, error)
	ProactiveDailyCount(ctx context.Context, groupID, botID uint) (int64, error)
	MarkProactiveSent(ctx context.Context, groupID, botID uint, cooldownTTL, dailyTTL time.Duration) error
}

func NewAICache(rdb *redis.Client) AICache {
	return &aicache{
		rdb: rdb,
	}
}

func (c *aicache) SetProactivePending(ctx context.Context, groupID uint, msgID int64, ttl time.Duration) error {
	return c.rdb.Set(ctx, c.proactivePendingKey(groupID), msgID, ttl).Err()
}

func (c *aicache) GetProactivePending(ctx context.Context, groupID uint) (int64, bool, error) {
	msgID, err := c.rdb.Get(ctx, c.proactivePendingKey(groupID)).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return msgID, true, nil
}

func (c *aicache) IsProactiveCoolingDown(ctx context.Context, groupID, botID uint) (bool, error) {
	exists, err := c.rdb.Exists(ctx, c.proactiveCooldownKey(groupID, botID)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (c *aicache) ProactiveDailyCount(ctx context.Context, groupID, botID uint) (int64, error) {
	count, err := c.rdb.Get(ctx, c.proactiveDailyKey(groupID, botID)).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return count, err
}

func (c *aicache) MarkProactiveSent(ctx context.Context, groupID, botID uint, cooldownTTL, dailyTTL time.Duration) error {
	if err := c.rdb.Set(ctx, c.proactiveCooldownKey(groupID, botID), "1", cooldownTTL).Err(); err != nil {
		return err
	}
	dailyKey := c.proactiveDailyKey(groupID, botID)
	if err := c.rdb.Incr(ctx, dailyKey).Err(); err != nil {
		return err
	}
	return c.rdb.Expire(ctx, dailyKey, dailyTTL).Err()
}

func (c *aicache) proactivePendingKey(groupID uint) string {
	return fmt.Sprintf("ai:proactive:pending:%d", groupID)
}

func (c *aicache) proactiveCooldownKey(groupID, botID uint) string {
	return fmt.Sprintf("ai:proactive:last:%d:%d", groupID, botID)
}

func (c *aicache) proactiveDailyKey(groupID, botID uint) string {
	return fmt.Sprintf("ai:proactive:daily:%s:%d:%d", time.Now().Format("20060102"), groupID, botID)
}

type aicache struct {
	rdb *redis.Client
}
