package video

import (
	"context"
	"fmt"
	"strconv"

	"StreamCore/pkg/util"
	"github.com/redis/go-redis/v9"
)

type VideoCache interface {
	OnVisited(ctx context.Context, vid uint) error
	GetVisitRank(ctx context.Context, limit, page int, desc bool) ([]uint, error)
	RebuildVisitRank(ctx context.Context, videos map[uint]int64) error
}

func NewVideoCache(rdb *redis.Client) VideoCache {
	return &videocache{
		rdb: rdb,
	}
}

func (c *videocache) OnVisited(ctx context.Context, vid uint) error {
	member := strconv.FormatUint(uint64(vid), 10)

	// incr zset score
	_, err := c.rdb.ZIncrBy(ctx, c.zVideoRankKey(), 1, member).Result()
	if err != nil {
		return err
	}

	// incr video visit
	_, err = c.rdb.Incr(ctx, c.videoVisitCountKey(vid)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *videocache) GetVisitRank(ctx context.Context, limit, page int, desc bool) ([]uint, error) {
	res, err := c.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   c.zVideoRankKey(),
		Start: limit * page,
		Stop:  limit*(page+1) - 1,
		Rev:   desc,
	}).Result()
	if err != nil {
		return nil, err
	}

	var videos []uint
	for _, s := range res {
		vid := util.String2Uint(s)
		videos = append(videos, vid)
	}
	return videos, nil
}

// RebuildVisitRank rebuilds the video ranking zset from database data
func (c *videocache) RebuildVisitRank(ctx context.Context, videos map[uint]int64) error {
	if len(videos) == 0 {
		return nil
	}

	// Use pipeline for batch operations
	pipe := c.rdb.Pipeline()

	// Clear existing zset
	pipe.Del(ctx, c.zVideoRankKey())

	// Add all videos to zset
	for vid, visitCount := range videos {
		member := strconv.FormatUint(uint64(vid), 10)
		pipe.ZAdd(ctx, c.zVideoRankKey(), redis.Z{
			Score:  float64(visitCount),
			Member: member,
		})
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (c *videocache) zVideoRankKey() string {
	return "zVideoRank"
}

func (c *videocache) videoVisitCountKey(vid uint) string {
	return fmt.Sprintf("video:visit_count:%d", vid)
}

type videocache struct {
	rdb *redis.Client
}
