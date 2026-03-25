package service

import (
	"fmt"

	"StreamCore/internal/pkg/mq/model"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
)

func (s *VideoService) Visit(uid *uint, query *video.VisitQuery) (*common.VideoInfo, error) {
	vid, err := util.ParseUint(query.VideoId)
	if err != nil {
		return nil, fmt.Errorf("bad videoId format: %w", err)
	}

	// get video metadata from db
	v, err := s.db.GetById(vid)
	if err != nil {
		return nil, fmt.Errorf("error db GetById: %w", err)
	}

	// cache OnVisited
	if err = s.cache.OnVisited(s.ctx, vid); err != nil {
		return nil, fmt.Errorf("error cache.OnVisited: %w", err)
	}

	// mq
	err = s.mq.PublishVisitEvent(s.ctx, &model.VisitEvent{
		VisitorId: uid,
		Vid:       vid,
	})
	if err != nil {
		return nil, fmt.Errorf("error mq.PublishVisitEvent: %w", err)
	}
	// TODO: mq fail fallback

	return pack.VideoInfo(v), nil
}
