package service

import (
	"context"
	"sync"

	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/interaction"
	db "StreamCore/internal/pkg/db/interaction"
	mq "StreamCore/internal/pkg/mq/interaction"
)

type InteractionService struct {
	ctx   context.Context
	db    db.InteractionDatabase
	cache cache.InteractionCache
	mq    mq.InteractionMQ
	infra *base.InfraSet
}

func NewInteractionService(ctx context.Context, infra *base.InfraSet) *InteractionService {
	svc := &InteractionService{
		ctx:   ctx,
		db:    infra.DB.Interaction,
		cache: infra.Cache.Interaction,
		mq:    infra.MQ.Interaction,
		infra: infra,
	}
	svc.initConsumer()
	return svc
}

var consumeOnce sync.Once

func (s *InteractionService) initConsumer() {
	consumeOnce.Do(func() {
		go s.consumer(context.Background())
	})
}
