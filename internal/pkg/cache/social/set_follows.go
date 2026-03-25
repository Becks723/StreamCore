package social

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"github.com/bytedance/sonic"
)

func (c *socialcache) SetFollows(ctx context.Context, uid uint, limit, page int, follows []uint, total int) error {
	key := c.followsKey(uid, limit, page)
	data, err := sonic.Marshal(&f{
		uids:  follows,
		total: total,
	})
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, constants.SocialCacheExpiration).Err()
}
