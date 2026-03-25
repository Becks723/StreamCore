package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/domain"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
)

func (s *InteractionService) ListComment(query *interaction.ListCommentQuery) (*interaction.ListCommentRespData, error) {
	var comments []*domain.Comment
	limit, page := int(query.PageSize), int(query.PageNum)

	if query.VideoId == nil && query.CommentId == nil {
		return nil, errors.New("either videoId or commentId should have value")
	}

	if query.VideoId != nil {
		vid, err := util.ParseUint(*query.VideoId)
		if err != nil {
			return nil, errors.New("bad videoId format")
		}
		comments, err = s.db.ListRootComments(vid, limit, page)
		if err != nil {
			return nil, fmt.Errorf("error db.ListRootComments: %w", err)
		}
	} else if query.CommentId != nil {
		cid, err := util.ParseUint(*query.CommentId)
		if err != nil {
			return nil, errors.New("bad commentId format")
		}
		comments, err = s.db.ListSubComments(cid, limit, page)
		if err != nil {
			return nil, fmt.Errorf("error db.ListSubComments: %w", err)
		}
	}

	data := new(interaction.ListCommentRespData)
	for _, c := range comments {
		data.Items = append(data.Items, pack.CommentInfo(c))
	}
	return data, nil
}
