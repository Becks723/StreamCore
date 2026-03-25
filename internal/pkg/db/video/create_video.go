package video

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

func (repo *videodb) Create(v *domain.Video) (err error) {
	var deletedAt gorm.DeletedAt
	if v.DeletedAt != nil {
		deletedAt.Valid = true
		deletedAt.Time = *v.DeletedAt
	} else {
		deletedAt.Valid = false
	}

	po := &model.VideoModel{
		Model: gorm.Model{
			ID:        v.Id,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			DeletedAt: deletedAt,
		},
		AuthorId:     v.AuthorId,
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   v.VisitCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		PublishedAt:  v.PublishedAt,
		EditedAt:     v.EditedAt,
	}
	return repo.db.
		Model(&model.VideoModel{}).
		Create(po).
		Error
}
