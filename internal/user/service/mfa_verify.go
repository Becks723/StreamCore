package service

import (
	"fmt"

	"StreamCore/kitex_gen/common"
	"StreamCore/kitex_gen/user"
)

// MFAVerify MFA校验
func (s *UserService) MFAVerify(req *user.MFAVerifyReq) (*common.TokenInfo, error) {
	var err error
	// 从 cache 拿 mfa_token 对应的 uid
	uid, err := s.cache.GetMFATokenUser(s.ctx, req.MfaToken)
	if err != nil {
		return nil, fmt.Errorf("error get mfa token: %w", err)
	}
	u, err := s.db.GetById(uid)
	if err != nil {
		return nil, fmt.Errorf("err db.GetById(%d): %w", uid, err)
	}
	// 校验
	success, err := s.totpAuth(uid, u.TOTPSecret, req.Code)
	if !success {
		return nil, fmt.Errorf("totp校验失败: %w", err)
	}

	// 生成令牌
	access, refresh, err := s.generateTokens(u.Id)
	if err != nil {
		return nil, err
	}
	data := new(common.TokenInfo)
	data.AccessToken = access
	data.RefreshToken = refresh
	return data, nil
}
