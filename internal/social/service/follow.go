package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/util"
)

func (s *SocialService) Follow(uid uint, req *social.FollowReq) error {
	followee, err := util.ParseUint(req.ToUserId)
	if err != nil {
		return errors.New("bad uid format")
	}
	if req.ActionType != constants.FollowAction_Follow && req.ActionType != constants.FollowAction_Unfollow {
		return fmt.Errorf("unknown follow actionType (value=%d)", req.ActionType)
	}

	if req.ActionType == constants.FollowAction_Follow {
		if err = s.writeFollowToDB(s.ctx, uid, followee, time.Now()); err != nil {
			return err
		}
	} else {
		if err = s.writeUnfollowToDB(s.ctx, uid, followee, time.Now()); err != nil {
			return err
		}
	}
	// invalidate cache
	if err = s.cache.InvalidateUserCache(s.ctx, uid); err != nil {
		return fmt.Errorf("error cache.InvalidateUserCache: %w", err)
	}
	return nil
}

func (s *SocialService) writeFollowToDB(ctx context.Context, follower, followee uint, time time.Time) error {
	f, err := s.db.GetFollow(ctx, follower, followee)
	if err != nil { // first follow
		if err = s.db.CreateFollow(ctx, follower, followee, time); err != nil {
			return fmt.Errorf("error db.CreateFollow: %w", err)
		}
		return nil
	}

	if f.Time.Before(time) && // time correct
		f.Status == constants.FollowAction_Unfollow { // status correct
		if err = s.db.UpdateFollowStatus(ctx, follower, followee, constants.FollowAction_Follow, time); err != nil {
			return fmt.Errorf("error db.UpdateFollowStatus: %w", err)
		}
	}
	return nil
}

func (s *SocialService) writeUnfollowToDB(ctx context.Context, follower, followee uint, time time.Time) error {
	f, err := s.db.GetFollow(ctx, follower, followee)
	if err != nil { // first unfollow, ignore
		return nil
	}

	if f.Time.Before(time) && // time correct
		f.Status == constants.FollowAction_Follow { // status correct
		if err = s.db.UpdateFollowStatus(ctx, follower, followee, constants.FollowAction_Unfollow, time); err != nil {
			return fmt.Errorf("error db.UpdateFollowStatus: %w", err)
		}
	}
	return nil
}
