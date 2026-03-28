package interaction

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/db/util"
	"StreamCore/internal/pkg/domain"
)

func (repo *iactiondb) ListRootComments(vid uint, limit, page int) (comments []*domain.Comment, err error) {
	var records []*model.CommentModel
	var cnt int64
	tx := repo.db.
		Model(&model.CommentModel{}).
		Where("video_id = ? AND parent_id IS NULL", vid).
		Count(&cnt)

	if util.IsPageParamsValid(limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return nil, err
	}
	for _, po := range records {
		comments = append(comments, pack.Comment(po))
	}
	return comments, nil
}

func (repo *iactiondb) ListSubComments(cid uint, limit, page int) (comments []*domain.Comment, err error) {
	var records []*model.CommentModel
	var cnt int64
	tx := repo.db.
		Model(&model.CommentModel{}).
		Where("parent_id = ?", cid).
		Count(&cnt)

	if util.IsPageParamsValid(limit, page) {
		tx = tx.Limit(limit).
			Offset(limit * page)
	}
	if err = tx.Find(&records).Error; err != nil {
		return nil, err
	}
	for _, po := range records {
		comments = append(comments, pack.Comment(po))
	}
	return comments, nil
}
