package rpccontext

import (
	"context"
	"errors"

	"StreamCore/pkg/util"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

// See https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/metainfo/
// and `context.WithValue()` is only for single process!

const (
	LoginUidKey = "uid"
)

func WithLoginUid(ctx context.Context, uid uint) context.Context {
	return metainfo.WithPersistentValue(ctx, LoginUidKey, util.Uint2String(uid))
}

func RetrieveLoginUid(ctx context.Context) (uint, error) {
	str, ok := metainfo.GetPersistentValue(ctx, LoginUidKey)
	if !ok {
		return 0, errors.New("error metainfo.GetPersistentValue")
	}
	return util.String2Uint(str), nil
}
