package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/pack"
	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/util"
	"StreamCore/pkg/util/jwt"
)

func (s *UserService) Login(req *user.LoginReq) (*common.UserInfo, *common.AuthenticationInfo, error) {
	var err error

	// find user in db
	u, err := s.db.GetByUsername(req.Username)
	if err != nil {
		return nil, nil, errors.New("用户不存在")
	}

	// password correct?
	if !util.CheckPassword(req.Password, u.Password) {
		return nil, nil, errors.New("密码错误")
	}

	// generate access, refresh tokens
	access, err := jwt.GenerateAccessToken(u.Id, constants.JWT_AccessSecret, constants.JWT_AccessTokenExpiration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen accessToken: %w", err)
	}
	refresh, err := jwt.GenerateRefreshToken(u.Id, constants.JWT_RefreshSecret, constants.JWT_RefreshTokenExpiration)
	if err != nil {
		return nil, nil, fmt.Errorf("failed gen refreshToken: %w", err)
	}

	// deal with mfa token
	mfaToken := ""
	if u.TOTPBound {
		mfaToken = s.generateMFAToken()
		// cache mfa token
		if err = s.cache.SetMFATokenUser(s.ctx, mfaToken, u.Id, constants.MFATokenExpiry); err != nil {
			return nil, nil, fmt.Errorf("failed cache.SetMFATokenUser: %w", err)
		}
	}

	auth := &common.AuthenticationInfo{
		AccessToken:  access,
		RefreshToken: refresh,
		MfaRequired:  u.TOTPBound,
		MfaToken:     mfaToken,
	}
	return pack.UserInfo(u), auth, nil
}

func (s *UserService) generateMFAToken() string {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return base64.RawStdEncoding.EncodeToString(buf)
}
