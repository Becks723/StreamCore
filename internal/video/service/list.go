package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
)

func (s *VideoService) List(query *video.ListQuery) (*video.ListRespData, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	limit := int(query.PageSize)
	page := int(query.PageNum)
	videos, total, err := s.db.FetchByUid(uid, limit, page)
	if err != nil {
		return nil, fmt.Errorf("error db.FetchByUid: %w", err)
	}

	data := new(video.ListRespData)
	data.Total = int32(total)
	for _, v := range videos {
		data.Items = append(data.Items, pack.VideoInfo(v))
	}
	return data, nil
}
