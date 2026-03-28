package service

import (
	"errors"
	"fmt"

	"StreamCore/config"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
)

func (s *VideoService) Popular(query *video.PopularQuery) (*video.PopularRespData, error) {
	var limit, page int
	if query.PageSize == nil {
		limit = config.Instance().General.PageSize
	} else {
		limit = int(*query.PageSize)
	}

	if query.PageNum == nil {
		page = 0
	} else {
		page = int(*query.PageNum)
	}

	// Get video IDs from cache (descending order - most popular first)
	vids, err := s.cache.GetVisitRank(s.ctx, limit, page, true)
	if err != nil { // cache unavailable
		// TODO: log cache unavailable
		// rebuild rank cache
		if err = s.rebuildVisitRankCache(); err != nil {
			return nil, err
		}
		vids, err = s.cache.GetVisitRank(s.ctx, limit, page, true)
		if err != nil {
			return nil, fmt.Errorf("error cache.GetVisitRank: %w", err)
		}
	}

	var videos []*common.VideoInfo
	var failCount int
	for _, vid := range vids {
		v, err := s.db.GetById(vid)
		if err != nil { // if fail to fetch data, fill in with only vid
			failCount++
			videos = append(videos, &common.VideoInfo{
				Id: util.Uint2String(vid),
			})
		} else {
			videos = append(videos, pack.VideoInfo(v))
		}
	}

	if failCount == len(vids) {
		return nil, errors.New("failed to fetch all video data, something might go wrong")
	}
	data := new(video.PopularRespData)
	data.Items = videos
	return data, nil
}

// rebuildVisitRankCache rebuilds the visit ranking cache from database
func (s *VideoService) rebuildVisitRankCache() error {
	// Fetch top N videos from database (larger than typical page size to populate cache)
	const cacheSize = 1000
	m, err := s.db.FetchVideoIdsByVisit(s.ctx, cacheSize, 0)
	if err != nil {
		return fmt.Errorf("error db.FetchVideoIdsByVisit: %w", err)
	}

	// Rebuild cache
	if err := s.cache.RebuildVisitRank(s.ctx, m); err != nil {
		return fmt.Errorf("error cache.RebuildVisitRank: %w", err)
	}
	return nil
}
