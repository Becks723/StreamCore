package service

import (
	"context"

	"StreamCore/internal/pkg/base"
	groupcache "StreamCore/internal/pkg/cache/group"
	groupdb "StreamCore/internal/pkg/db/group"
)

type GroupService struct {
	ctx   context.Context
	infra *base.InfraSet
	db    groupdb.GroupDatabase
	cache groupcache.GroupCache
}

func NewGroupService(ctx context.Context, infra *base.InfraSet) *GroupService {
	return &GroupService{
		ctx:   ctx,
		infra: infra,
		db:    infra.DB.Group,
		cache: infra.Cache.Group,
	}
}
