package service

import (
	"context"
	"sync"

	"StreamCore/internal/pkg/base"
	cache "StreamCore/internal/pkg/cache/video"
	db "StreamCore/internal/pkg/db/video"
	es "StreamCore/internal/pkg/es/video"
	mq "StreamCore/internal/pkg/mq/video"
)

type VideoService struct {
	ctx   context.Context
	db    db.VideoDatabase
	cache cache.VideoCache
	es    es.VideoElastic
	mq    mq.VideoMQ
	infra *base.InfraSet
}

func NewVideoService(ctx context.Context, infra *base.InfraSet) *VideoService {
	svc := &VideoService{
		ctx:   ctx,
		db:    infra.DB.Video,
		cache: infra.Cache.Video,
		es:    es.NewVideoElastic(infra.ES),
		mq:    infra.MQ.Video,
		infra: infra,
	}
	svc.initConsumer()
	return svc
}

var consumeOnce sync.Once

func (s *VideoService) initConsumer() {
	consumeOnce.Do(func() {
		go s.consumeVisit(context.Background())
	})
}
