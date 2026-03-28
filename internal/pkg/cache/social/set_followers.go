package social

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"github.com/bytedance/sonic"
)

func (c *socialcache) SetFollowers(ctx context.Context, uid uint, limit, page int, followers []uint, total int) error {
	key := c.followersKey(uid, limit, page)
	data, err := sonic.Marshal(&f{
		uids:  followers,
		total: total,
	})
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, constants.SocialCacheExpiration).Err()
}
