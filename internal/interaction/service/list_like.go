package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
)

func (s *InteractionService) ListLikedVideos(query *interaction.ListLikeQuery) (*interaction.ListLikeRespData, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, fmt.Errorf("bad uid format")
	}
	pageSize := int(query.PageSize)
	pageNum := int(query.PageNum)
	if pageSize < 1 {
		pageSize = 1
	}
	if pageNum < 0 {
		pageNum = 0
	}

	offset := pageNum * pageSize
	limit := constants.UserLikesCacheLimit
	var vids []uint
	var cachedCount int

	if offset < limit {
		start := int64(offset)
		stop := int64(min(offset+pageSize-1, limit-1))
		cached, err := s.cache.GetUserLikedVideosRange(s.ctx, uid, start, stop)
		if err == nil {
			vids = append(vids, cached...)
			cachedCount = len(cached)
		}
	}

	if len(vids) < pageSize {
		needed := pageSize - len(vids)
		dbOffset := offset + len(vids)
		likes, err := s.db.FetchUserLikedVideos(s.ctx, uid, needed, dbOffset)
		if err != nil {
			return nil, fmt.Errorf("error db.FetchUserLikedVideos: %w", err)
		}
		for _, like := range likes {
			vids = append(vids, like.TargetId)
		}

		if cachedCount == 0 && offset < limit {
			if rebuildLikes, err := s.db.FetchUserLikedVideos(s.ctx, uid, limit, 0); err == nil && len(rebuildLikes) > 0 {
				_ = s.cache.SetUserLikedVideos(s.ctx, uid, rebuildLikes)
			}
		}
	}

	if len(vids) == 0 {
		return &interaction.ListLikeRespData{}, nil
	}
	if len(vids) > pageSize {
		vids = vids[:pageSize]
	}

	return s.saturateVideos(vids)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *InteractionService) saturateVideos(vids []uint) (*interaction.ListLikeRespData, error) {
	var videos []*common.VideoInfo
	var failCount int
	for _, vid := range vids {
		v, err := s.infra.DB.Video.GetById(vid)
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
	data := new(interaction.ListLikeRespData)
	data.Items = videos
	return data, nil
}
