package service

import (
	"errors"
	"fmt"

	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/util"
)

func (s *UserService) GetInfo(query *user.InfoQuery) (*common.UserInfo, error) {
	// convert string id to uint
	uid, err := util.ParseUint(query.UserId)
	if err != nil {
		return nil, errors.New("bad uid format")
	}

	// find user in db
	u, err := s.db.GetById(uid)
	if err != nil {
		return nil, fmt.Errorf("cannot find user (uid=%d)", uid)
	}

	return pack.UserInfo(u), nil
}
