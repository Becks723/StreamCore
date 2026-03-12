package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"StreamCore/internal/pkg/constants"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	SetTOTPPending(ctx context.Context, uid uint, secret string, ttl time.Duration) error
	GetTOTPPending(ctx context.Context, uid uint) (string, error)
	MarkTOTPTimestep(ctx context.Context, uid uint, code string, ttl time.Duration) error
	IsTOTPTimestepMarked(ctx context.Context, uid uint, code string) (bool, error)
	IncreaseTOTPFailure(ctx context.Context, uid uint, ttl time.Duration) error
	TOTPFailureCount(ctx context.Context, uid uint) (int, error)
	SetMFATokenUser(ctx context.Context, token string, uid uint, ttl time.Duration) error
	GetMFATokenUser(ctx context.Context, token string) (uint, error)
}

func NewUserCache(rdb *redis.Client) UserCache {
	return &usercache{
		rdb: rdb,
	}
}

func (c *usercache) SetTOTPPending(ctx context.Context, uid uint, secret string, ttl time.Duration) error {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Set(ctx, cacheKey, secret, ttl).Err()
}

func (c *usercache) GetTOTPPending(ctx context.Context, uid uint) (string, error) {
	cacheKey := c.totpPendingKey(uid)
	return c.rdb.Get(ctx, cacheKey).Result()
}

func (c *usercache) MarkTOTPTimestep(ctx context.Context, uid uint, code string, ttl time.Duration) error {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpTimestepKey(uid, timestep)
	return c.rdb.Set(ctx, key, code, ttl).Err()
}

func (c *usercache) IsTOTPTimestepMarked(ctx context.Context, uid uint, code string) (bool, error) {
	timestep := time.Now().UnixMilli() / constants.TOTPInterval
	key := c.totpTimestepKey(uid, timestep)
	cache, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return cache == code, nil
}

func (c *usercache) IncreaseTOTPFailure(ctx context.Context, uid uint, ttl time.Duration) error {
	key := c.totpFailureKey(uid)
	ok, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if ok == 0 { // key not exists, create key with ttl
		err = c.rdb.Set(ctx, key, 1, ttl).Err()
	} else { // key exists, INCR
		err = c.rdb.Incr(ctx, key).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *usercache) TOTPFailureCount(ctx context.Context, uid uint) (int, error) {
	key := c.totpFailureKey(uid)
	ok, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ok == 0 { // key not exists, return 0
		return 0, nil
	}
	// key exists, return value of key
	str, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	count, _ := strconv.ParseInt(str, 10, 64)
	return int(count), nil
}

func (c *usercache) GetMFATokenUser(ctx context.Context, token string) (uint, error) {
	key := c.mfaTokenKey(token)
	uid, err := c.rdb.Get(ctx, key).Uint64()
	return uint(uid), err
}

func (c *usercache) SetMFATokenUser(ctx context.Context, token string, uid uint, ttl time.Duration) error {
	key := c.mfaTokenKey(token)
	return c.rdb.Set(ctx, key, uid, ttl).Err()
}

func (c *usercache) totpPendingKey(uid uint) string {
	return fmt.Sprintf("totp:pending:%d", uid)
}

func (c *usercache) totpTimestepKey(uid uint, period int64) string {
	return fmt.Sprintf("totp:replay:%d:%d", uid, period)
}

func (c *usercache) totpFailureKey(uid uint) string {
	return fmt.Sprintf("totp:failure:%d", uid)
}

func (c *usercache) mfaTokenKey(token string) string {
	return fmt.Sprintf("mfa:token:%s", token)
}

type usercache struct {
	rdb *redis.Client
}
