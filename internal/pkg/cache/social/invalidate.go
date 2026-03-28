package social

import (
	"context"
	"fmt"
)

func (c *socialcache) InvalidateUserCache(ctx context.Context, uid uint) error {
	keys := []string{
		fmt.Sprintf("social:follows:%d:*", uid),
		fmt.Sprintf("social:followers:%d:*", uid),
		fmt.Sprintf("social:friends:%d:*", uid),
	}
	return c.rdb.Del(ctx, keys...).Err()
}
