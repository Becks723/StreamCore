package social

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
)

func (c *socialcache) GetFollows(ctx context.Context, uid uint, limit, page int) ([]uint, int, error) {
	key := c.followsKey(uid, limit, page)
	buffer, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, -1, err
	}
	var data f
	err = sonic.Unmarshal(buffer, &data)
	if err != nil {
		return nil, -1, err
	}
	return data.uids, data.total, nil
}

func (c *socialcache) followsKey(uid uint, limit, page int) string {
	return fmt.Sprintf("social:follows:%d:%d:%d", uid, limit, page)
}

type f struct {
	uids  []uint
	total int
}
