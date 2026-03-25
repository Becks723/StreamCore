package service

import (
	"errors"
	"fmt"

	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
)

func (s *InteractionService) DeleteComment(uid uint, req *interaction.DeleteCommentReq) error {
	cid, err := util.ParseUint(req.CommentId)
	if err != nil {
		return errors.New("bad commentId format")
	}

	err = s.db.DeleteCommentById(cid, uid)
	if err != nil {
		return fmt.Errorf("error db.DeleteCommentById: %w", err)
	}
	return nil
}
