package pack

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func Video(po *model.VideoModel) *domain.Video {
	return &domain.Video{
		Id:           po.ID,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
		DeletedAt:    packDeletedAt(po.DeletedAt),
		AuthorId:     po.AuthorId,
		VideoUrl:     po.VideoUrl,
		CoverUrl:     po.CoverUrl,
		Title:        po.Title,
		Description:  po.Description,
		VisitCount:   po.VisitCount,
		LikeCount:    po.LikeCount,
		CommentCount: po.CommentCount,
		PublishedAt:  po.PublishedAt,
		EditedAt:     po.EditedAt,
	}
}
