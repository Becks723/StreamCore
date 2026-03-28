package video

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *videodb) GetById(vid uint) (v *domain.Video, err error) {
	var po model.VideoModel
	err = repo.db.
		Model(&model.VideoModel{}).
		Where("id = ?", vid).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.Video(&po), nil
}
