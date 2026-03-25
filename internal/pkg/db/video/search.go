package video

import (
	"time"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/db/util"
	"StreamCore/internal/pkg/domain"
)

func (repo *videodb) Search(keywords string, limit, page int, from, to *time.Time, username *string) (videos []*domain.Video, total int, err error) {
	var records []*model.VideoModel

	tx := repo.db.Table("video_models v")
	if keywords != "" {
		tx = tx.Where("title LIKE ? OR description LIKE ?",
			"%"+keywords+"%", "%"+keywords+"%")
	}
	if from != nil {
		tx = tx.Where("published_at > ?", from)
	}
	if to != nil {
		tx = tx.Where("published_at < ?", to)
	}
	if username != nil {
		tx = tx.Joins("JOIN user_models u ON u.id = v.author_id").
			Where("u.username LIKE ?", "%"+*username+"%")
	}
	var cnt int64
	tx = tx.Count(&cnt)
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
