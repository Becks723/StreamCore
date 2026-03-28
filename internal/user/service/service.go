package service

import (
	"context"

	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/user"
	db "StreamCore/internal/pkg/db/user"
)

type UserService struct {
	ctx   context.Context
	db    db.UserDatabase
	cache cache.UserCache
	infra *base.InfraSet
}

func NewUserService(ctx context.Context, infra *base.InfraSet) *UserService {
	return &UserService{
		ctx:   ctx,
		db:    infra.DB.User,
		cache: infra.Cache.User,
		infra: infra,
	}
}
