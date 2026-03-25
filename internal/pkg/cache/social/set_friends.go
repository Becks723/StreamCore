package social

import (
	"context"

	"StreamCore/internal/pkg/constants"
	"github.com/bytedance/sonic"
)

func (c *socialcache) SetFriends(ctx context.Context, uid uint, limit, page int, friends []uint, total int) error {
	key := c.friendsKey(uid, limit, page)
	data, err := sonic.Marshal(&f{
		uids:  friends,
		total: total,
	})
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, constants.SocialCacheExpiration).Err()
}
