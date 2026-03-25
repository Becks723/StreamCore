package service

import (
	"fmt"
	"time"

	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
)

func (s *VideoService) GetVideoFeed(query *video.FeedQuery) (*video.FeedRespData, error) {
	var err error
	var after *time.Time

	if query.LatestTime == nil {
		after = nil
	} else {
		var t time.Time
		t, err = util.FromTimestamp(*query.LatestTime)
		if err != nil {
			return nil, fmt.Errorf("bad timestamp argument: %s", "LatestTime")
		}
		after = &t
	}

	videos, err := s.db.Fetch(after)
	if err != nil {
		return nil, fmt.Errorf("error fetching db video feed: %w", err)
	}

	data := new(video.FeedRespData)
	for _, v := range videos {
		data.Items = append(data.Items, pack.VideoInfo(v))
	}
	return data, nil
}
