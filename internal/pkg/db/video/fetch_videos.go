package video

import (
	"context"
	"time"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/db/util"
	"StreamCore/internal/pkg/domain"
)

func (repo *videodb) Fetch(after *time.Time) (videos []*domain.Video, err error) {
	var records []*model.VideoModel
	if after == nil {
		err = repo.db.Find(&records).Error
	} else {
		err = repo.db.Where("published_at > ?", after).Find(&records).Error
	}
	if err != nil {
		return nil, err
	}

	for _, po := range records {
		videos = append(videos, pack.Video(po))
	}
	return videos, nil
}

func (repo *videodb) FetchByUid(uid uint, limit, page int) (videos []*domain.Video, total int, err error) {
	var records []*model.VideoModel
	var cnt int64
	tx := repo.db.
		Model(&model.VideoModel{}).
		Where("author_id = ?", uid).
		Count(&cnt)
	if err = tx.Error; err != nil {
		return nil, -1, err
	}

	if util.IsPageParamsValid(limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return nil, -1, err
	}

	for _, po := range records {
		videos = append(videos, pack.Video(po))
	}
	return videos, int(cnt), nil
}

func (repo *videodb) FetchVideoIdsByVisit(ctx context.Context, limit, page int) (map[uint]int64, error) {
	type kv struct {
		vid   uint
		visit int64
	}
	var pairs []kv
	err := repo.db.Model(&model.VisitCountModel{}).
		Select("vid, visit_count").
		Order("visit_count DESC").
		Limit(limit).
		Offset(limit * page).
		Scan(&pairs).
		Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint]int64, len(pairs))
	for _, p := range pairs {
		m[p.vid] = p.visit
	}
	return m, nil
}
