package service

import (
	"errors"

	"StreamCore/internal/pkg/domain"
	"StreamCore/pkg/util"
)

func (s *UserService) Register(username, password string) error {
	var err error

	// check potential duplicated username
	if _, err = s.db.GetByUsername(username); err == nil {
		return errors.New("用户名已存在")
	}

	// encrypt password
	var encrypted string
	if encrypted, err = util.EncryptPassword(password); err != nil {
		return errors.New("请换一个密码")
	}

	// repo stuff
	do := domain.User{
		Username:  username,
		Password:  encrypted,
		AvatarUrl: "", // TODO: offer default avatar
	}
	if err = s.db.Create(&do); err != nil {
		return err
	}
	return nil
}
