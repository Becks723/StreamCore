package service

import (
	"fmt"
	"time"

	"StreamCore/internal/pkg/domain"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/video"
	"StreamCore/pkg/util"
)

func (s *VideoService) Search(query *video.SearchReq) (*video.SearchRespData, error) {
	var err error

	// resolve from/toDate
	var from, to *time.Time
	var tmp time.Time
	if query.FromDate != nil {
		tmp, err = util.FromTimestamp(*query.FromDate)
		if err != nil {
			return nil, fmt.Errorf("bad timestamp argument: %s", "FromDate")
		}
		from = &tmp
	}
	if query.ToDate != nil {
		tmp, err = util.FromTimestamp(*query.ToDate)
		if err != nil {
			return nil, fmt.Errorf("bad timestamp argument: %s", "ToDate")
		}
		to = &tmp
	}

	// core search
	videos, total, err := s.db.Search(query.Keywords, int(query.PageSize), int(query.PageNum), from, to, query.Username)
	if err != nil {
		return nil, fmt.Errorf("error db search: %w", err)
	}

	data := new(video.SearchRespData)
	data.Total = int32(total)
	for _, v := range videos {
		data.Items = append(data.Items, pack.VideoInfo(v))
	}
	return data, nil
}

// SearchES - search using es
func (s *VideoService) SearchES(query *video.SearchReq) (*video.SearchRespData, error) {
	var err error

	esquery := &domain.VideoQuery{
		TitleMatches:    query.Keywords,
		DescMatches:     query.Keywords,
		FromDate:        query.FromDate,
		ToDate:          query.ToDate,
		UsernameMatches: query.Username,
	}
	hits, total, err := s.es.SearchVideo(s.ctx, esquery)
	if err != nil {
		return nil, err
	}

	data := new(video.SearchRespData)
	data.Total = int32(total)
	failId := make([]uint, 0)
	for _, id := range hits {
		v, err := s.db.GetById(id)
		if err != nil {
			failId = append(failId, id)
		} else {
			data.Items = append(data.Items, pack.VideoInfo(v))
		}
	}

	if len(failId) == int(total) { // all failed -> throw error
		return nil, fmt.Errorf("StreamService.SearchEs failed: all %d hits fetch failed", total)
	} else { // partial fail is acceptable
		return data, nil
	}
}
