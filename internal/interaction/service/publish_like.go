package service

import (
	"context"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/mq/model"
	"StreamCore/kitex_gen/interaction"
	"StreamCore/pkg/util"
	"github.com/bytedance/sonic"
)

func (s *InteractionService) PublishLike(uid uint, req *interaction.PublishLikeReq) error {
	var err error
	if req.VideoId == nil && req.CommentId == nil {
		return fmt.Errorf("either videoId or commentId should have value")
	}
	if req.ActionType != constants.LikeAction_Like && req.ActionType != constants.LikeAction_Unlike {
		return fmt.Errorf("unknown actionType (value=%d)", req.ActionType)
	}

	var tarType int
	var tarId uint
	if req.VideoId != nil {
		vid, err := util.ParseUint(*req.VideoId)
		if err != nil {
			return fmt.Errorf("bad videoId format: %w", err)
		}
		tarType = constants.LikeTarType_Video
		tarId = vid
	} else {
		cid, err := util.ParseUint(*req.CommentId)
		if err != nil {
			return fmt.Errorf("bad commentId format: %w", err)
		}
		tarType = constants.LikeTarType_Comment
		tarId = cid
	}

	// mq
	err = s.mq.PublishLikeEvent(s.ctx, &model.LikeEvent{
		TarType: tarType,
		TarId:   tarId,
		Uid:     uid,
		Action:  int(req.ActionType),
		Time:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error mq.PublishLikeEvent: %w", err)
	}
	return nil
}

func (s *InteractionService) consumeLikeEvent(ctx context.Context, ev *model.LikeEvent) error {
	if ev.Action == constants.LikeAction_Like {
		return s.consumeLike(ctx, ev)
	} else {
		return s.consumeUnlike(ctx, ev)
	}
}

func (s *InteractionService) consumeLike(ctx context.Context, ev *model.LikeEvent) error {
	tarType, tarId, uid, t := ev.TarType, ev.TarId, ev.Uid, ev.Time
	lastLike, err := s.db.GetLike(ctx, tarType, uid, tarId)
	if err != nil { // not recorded before
		if err = s.db.CreateLike(ctx, tarType, uid, tarId, t); err != nil {
			return fmt.Errorf("error db.CreateLike: %w", err)
		}
	} else { // recorded before
		if lastLike.Time.After(t) && // ignore former request
			lastLike.Status == constants.LikeAction_Like { // ignore repeated like
			return nil
		}
		err = s.db.ToggleLikeStatus(ctx, tarType, uid, tarId)
		if err != nil {
			return fmt.Errorf("error db.ToggleLikeStatus: %w", err)
		}
	}
	// write cache
	if err = s.cache.OnLiked(ctx, tarType, uid, tarId, t); err != nil {
		return fmt.Errorf("error cache.OnLiked: %w", err)
	}
	return nil
}

func (s *InteractionService) consumeUnlike(ctx context.Context, ev *model.LikeEvent) error {
	tarType, tarId, uid, t := ev.TarType, ev.TarId, ev.Uid, ev.Time
	lastLike, err := s.db.GetLike(ctx, tarType, uid, tarId)
	if err != nil { // not recorded before
		return nil
	} else {
		if lastLike.Time.After(t) && // ignore former request
			lastLike.Status == constants.LikeAction_Unlike { // ignore repeated unlike
			return nil
		}
		err = s.db.ToggleLikeStatus(ctx, tarType, uid, tarId)
		if err != nil {
			return fmt.Errorf("error db.ToggleLikeStatus: %w", err)
		}
	}
	// write cache
	if err = s.cache.OnUnliked(ctx, tarType, uid, tarId); err != nil {
		return fmt.Errorf("error db.OnUnliked: %w", err)
	}
	return nil
}

func (s *InteractionService) consumer(ctx context.Context) {
	c, err := s.mq.Consumer()
	if err != nil {
		// TODO: log consumer init error
		return
	}
	for {
		msg, err := c.Receive()
		if err != nil {
			break
		}
		var ev model.LikeEvent
		err = sonic.Unmarshal(msg.Body, &ev)
		if err != nil {
			continue
		}

		// >>>> consume like event <<<<
		if err = s.consumeLikeEvent(ctx, &ev); err != nil {

		} else {
			c.Ack(msg)
		}
		// =============================
	}
}
