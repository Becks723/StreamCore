package pack

import (
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
)

func VideoInfo(v *domain.Video) *common.VideoInfo {
	return &common.VideoInfo{
		CreatedAt:    v.CreatedAt.String(),
		UpdatedAt:    v.UpdatedAt.String(),
		DeletedAt:    util.TimePtr2String(v.DeletedAt),
		Id:           util.Uint2String(v.Id),
		UserId:       util.Uint2String(v.AuthorId),
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   int32(v.VisitCount),
		LikeCount:    int32(v.LikeCount),
		CommentCount: int32(v.CommentCount),
	}
}
