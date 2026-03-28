package pack

import (
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
)

func CommentInfo(c *domain.Comment) *common.CommentInfo {
	return &common.CommentInfo{
		CreatedAt:  c.CreatedAt.String(),
		UpdatedAt:  c.UpdatedAt.String(),
		DeletedAt:  util.TimePtr2String(c.DeletedAt),
		Id:         util.Uint2String(c.Id),
		UserId:     util.Uint2String(c.AuthorId),
		VideoId:    util.Uint2String(c.VideoId),
		ParentId:   util.Uint2StringOrEmpty(c.ParentId),
		Content:    c.Content,
		LikeCount:  int32(c.LikeCount),
		ChildCount: int32(c.ChildCount),
	}
}
