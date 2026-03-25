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
)

func (s *UserService) Login(req *user.LoginReq) (*common.UserInfo, *user.MFAInfo, *common.TokenInfo, error) {
	var err error

	// find user in db
	u, err := s.db.GetByUsername(req.Username)
	if err != nil {
		return nil, nil, nil, errors.New("用户不存在")
	}

	// password correct?
	if !util.CheckPassword(req.Password, u.Password) {
		return nil, nil, nil, errors.New("密码错误")
	}

	// deal with mfa token
	mfaToken := ""
	if u.TOTPBound {
		mfaToken = s.generateMFAToken()
		// cache mfa token
		if err = s.cache.SetMFATokenUser(s.ctx, mfaToken, u.Id, constants.MFATokenExpiry); err != nil {
			return nil, nil, nil, fmt.Errorf("failed cache.SetMFATokenUser: %w", err)
		}
	}
	auth := &user.MFAInfo{
		MfaRequired: u.TOTPBound,
		MfaToken:    mfaToken,
	}

	// jwt (only if MFA disabled)
	var token *common.TokenInfo
	if u.TOTPBound {
		token = nil
	} else {
		access, refresh, err := s.generateTokens(u.Id)
		if err != nil {
			return nil, nil, nil, err
		}
		token = &common.TokenInfo{
			AccessToken:  access,
			RefreshToken: refresh,
		}
	}
	return pack.UserInfo(u), auth, token, nil
}

func (s *UserService) generateMFAToken() string {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return base64.RawStdEncoding.EncodeToString(buf)
}
