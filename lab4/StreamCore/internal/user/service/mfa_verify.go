package service

import (
	"fmt"

	"StreamCore/kitex_gen/user"
)

// MFAVerify MFA校验
func (s *UserService) MFAVerify(uid uint, req *user.MFAVerifyReq) error {
	var err error
	// 从 cache 拿 mfa_token 对应的 uid，与登录的 uid 比较
	cacheUid, err := s.cache.GetMFATokenUser(s.ctx, req.MfaToken)
	if err != nil {
		return fmt.Errorf("error mfa token: %w", err)
	}
	if uid != cacheUid {
		return fmt.Errorf("login uid not match")
	}

	u, err := s.db.GetById(uid)
	if err != nil {
		return fmt.Errorf("err db.GetById(%d): %w", uid, err)
	}
	// 校验
	success, err := s.totpAuth(uid, u.TOTPSecret, req.Code)
	if !success {
		return fmt.Errorf("totp校验失败: %w", err)
	}
	return nil
}
