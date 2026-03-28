package pack

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func Like(po *model.LikeRelationModel) *domain.Like {
	return &domain.Like{
		Uid:        po.Uid,
		TargetType: po.TargetType,
		TargetId:   po.TargetId,
		Status:     po.Status,
		Time:       po.Time,
	}
}

func Comment(po *model.CommentModel) *domain.Comment {
	return &domain.Comment{
		Id:         po.ID,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
		DeletedAt:  packDeletedAt(po.DeletedAt),
		AuthorId:   po.AuthorId,
		VideoId:    po.VideoId,
		Content:    po.Content,
		ParentId:   po.ParentId,
		LikeCount:  po.LikeCount,
		ChildCount: po.ChildCount,
	}
}
