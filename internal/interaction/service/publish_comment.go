package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
)

func (s *InteractionService) PublishComment(uid uint, req *interaction.PublishCommentReq) error {
	if req.VideoId == nil && req.CommentId == nil {
		return errors.New("either videoId or commentId should have value")
	}

	if req.VideoId != nil {
		vid, err := util.ParseUint(*req.VideoId)
		if err != nil {
			return errors.New("bad videoId format")
		}
		err = s.db.CreateComment(s.ctx, &domain.Comment{
			AuthorId: uid,
			VideoId:  vid,
			Content:  req.Content,
		})
		if err != nil {
			return fmt.Errorf("error db.CreateComment: %w", err)
		}
	} else {
		cid, err := util.ParseUint(*req.CommentId)
		if err != nil {
			return errors.New("bad commentId format")
		}
		err = s.db.CreateComment(s.ctx, &domain.Comment{
			AuthorId: uid,
			ParentId: &cid,
			Content:  req.Content,
		})
		if err != nil {
			return fmt.Errorf("error db.CreateComment: %w", err)
		}
	}
	return nil
}
