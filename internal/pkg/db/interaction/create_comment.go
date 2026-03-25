package interaction

import (
	"context"
	"fmt"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

func (repo *iactiondb) CreateComment(ctx context.Context, c *domain.Comment) (err error) {
	var deletedAt gorm.DeletedAt
	if c.DeletedAt != nil {
		deletedAt.Valid = true
		deletedAt.Time = *c.DeletedAt
	} else {
		deletedAt.Valid = false
	}

	po := &model.CommentModel{
		Model: gorm.Model{
			ID:        c.Id,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: deletedAt,
		},
		AuthorId:   c.AuthorId,
		VideoId:    c.VideoId,
		Content:    c.Content,
		ParentId:   c.ParentId,
		LikeCount:  c.LikeCount,
		ChildCount: c.ChildCount,
	}

	if po.ParentId != nil { // is sub, ensure videoId
		var parent model.CommentModel
		err = repo.db.First(&parent, *po.ParentId).Error // call First to throw an error if not found
		if err != nil {
			return fmt.Errorf("parent comment(id:%d) not found", *po.ParentId)
		}
		po.VideoId = parent.VideoId
	}
	return repo.db.
		Model(&model.CommentModel{}).
		Create(&po).
		Error
}
