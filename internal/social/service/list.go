package service

import (
	"errors"
	"fmt"

	"StreamCore/config"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/social"
	"StreamCore/pkg/util"
)

func (s *SocialService) ListFollowers(query *social.ListFollowersQuery) (*social.SocialList, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	var limit, page int
	if query.PageSize == nil {
		limit = config.Instance().General.PageSize
	} else {
		limit = int(*query.PageSize)
	}

	if query.PageNum == nil {
		page = 0
	} else {
		page = int(*query.PageNum)
	}

	// cache-aside
	// try cache
	followers, total, err := s.cache.GetFollowers(s.ctx, uid, limit, page)
	if err != nil { // cache failed
		// TODO: log cache unavailable
		// fetch from db
		followers, total, err = s.db.FetchFollowers(s.ctx, uid, limit, page)
		if err != nil {
			return nil, fmt.Errorf("error db.FetchFollowers: %w", err)
		}
		if err = s.cache.SetFollowers(s.ctx, uid, limit, page, followers, total); err != nil {
			return nil, fmt.Errorf("error cache.SetFollowers: %w", err)
		}
	}
	return s.saturateSocialList(followers, total)
}

func (s *SocialService) ListFollows(query *social.ListFollowsQuery) (*social.SocialList, error) {
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	var limit, page int
	if query.PageSize == nil {
		limit = config.Instance().General.PageSize
	} else {
		limit = int(*query.PageSize)
	}

	if query.PageNum == nil {
		page = 0
	} else {
		page = int(*query.PageNum)
	}

	// cache-aside
	// try cache
	follows, total, err := s.cache.GetFollows(s.ctx, uid, limit, page)
	if err != nil { // cache failed
		// TODO: log cache unavailable
		// fetch from db
		follows, total, err = s.db.FetchFollows(s.ctx, uid, limit, page)
		if err != nil {
			return nil, fmt.Errorf("error db.FetchFollows: %w", err)
		}
		if err = s.cache.SetFollows(s.ctx, uid, limit, page, follows, total); err != nil {
			return nil, fmt.Errorf("error cache.SetFollows: %w", err)
		}
	}
	return s.saturateSocialList(follows, total)
}

func (s *SocialService) ListFriends(uid uint, query *social.ListFriendsQuery) (*social.SocialList, error) {
	var limit, page int
	if query.PageSize == nil {
		limit = config.Instance().General.PageSize
	} else {
		limit = int(*query.PageSize)
	}

	if query.PageNum == nil {
		page = 0
	} else {
		page = int(*query.PageNum)
	}

	// cache-aside
	// try cache
	friends, total, err := s.cache.GetFriends(s.ctx, uid, limit, page)
	if err != nil { // cache failed
		// TODO: log cache unavailable
		// fetch from db
		friends, total, err = s.db.FetchFriends(s.ctx, uid, limit, page)
		if err != nil {
			return nil, fmt.Errorf("error db.FetchFriends: %w", err)
		}
		if err = s.cache.SetFriends(s.ctx, uid, limit, page, friends, total); err != nil {
			return nil, fmt.Errorf("error cache.SetFriends: %w", err)
		}
	}
	return s.saturateSocialList(friends, total)
}

func (s *SocialService) saturateSocialList(uids []uint, total int) (*social.SocialList, error) {
	data := new(social.SocialList)
	data.Total = int32(total)
	for _, uid := range uids {
		data.Items = append(data.Items, s.getSocialInfo(uid))
	}
	return data, nil
}

func (s *SocialService) getSocialInfo(uid uint) *common.SocialUserInfo {
	u, err := s.infra.DB.User.GetById(uid)
	if err != nil {
		return &common.SocialUserInfo{
			Id: util.Uint2String(uid),
		}
	} else {
		return pack.SocialUserInfo(u)
	}
}
