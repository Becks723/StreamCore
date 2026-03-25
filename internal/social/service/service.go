package service

import (
	"context"

	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/social"
	db "StreamCore/internal/pkg/db/social"
)

type SocialService struct {
	ctx   context.Context
	db    db.SocialDatabase
	cache cache.SocialCache
	infra *base.InfraSet
}

func NewSocialService(ctx context.Context, infra *base.InfraSet) *SocialService {
	return &SocialService{
		ctx:   ctx,
		db:    infra.DB.Social,
		cache: infra.Cache.Social,
		infra: infra,
	}
}
