package interaction

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) GetCommentById(cid uint) (c *domain.Comment, err error) {
	var po model.CommentModel
	err = repo.db.
		Model(&model.CommentModel{}).
		Where("id = ?", cid).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.Comment(&po), nil
}
