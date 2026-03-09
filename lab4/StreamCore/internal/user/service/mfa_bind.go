package service

import (
	"errors"
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/user"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// MFABind 绑定mfa
func (s *UserService) MFABind(uid uint, req *user.MFABindReq) error {
	var err error

	// 检查用户是否已绑定mfa，避免重复绑定
	u, _ := s.db.GetById(uid)
	if u.TOTPBound {
		return errors.New("重复绑定")
	}
	// 从缓存拿secret
	pending, err := s.cache.GetTOTPPending(s.ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get pending secret: %w", err)
	}

	// totp认证
	success, err := s.totpAuth(uid, pending, req.Code)
	// 认证失败
	if !success {
		return fmt.Errorf("failed to bind mfa, try to refresh the qrcode: %w", err)
	}
	// 认证成功
	// 密钥存db
	err = s.db.UpdateTOTPSecret(uid, pending)
	if err != nil {
		return err
	}
	return nil
}

// totpAuth totp认证逻辑
func (s *UserService) totpAuth(uid uint, secret string, code string) (bool, error) {
	failCount, err := s.cache.TOTPFailureCount(s.ctx, uid)
	if err != nil {
		return false, err
	}
	// 防爆破检验，每个用户设置失败次数限制（如10次/5分钟）
	if failCount > constants.TOTPFailureLimit {
		return false, errors.New("failure exceeds limit, please try again later")
	}
	// 检查用户端和服务器生成的验证码一致性
	success, _ := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    constants.TOTPInterval,
		Algorithm: otp.AlgorithmSHA1,
		Digits:    6,
		Skew:      1,
	})
	if !success { // 不一致
		// 失败次数+1
		err = s.cache.IncreaseTOTPFailure(s.ctx, uid, constants.TOTPFailureReset)
		if err != nil {
			return false, err
		}
		return false, errors.New("wrong code")
	}
	// 一致
	// 检测重放攻击
	replay, err := s.checkReplay(uid)
	if err != nil {
		return false, err
	}
	if replay {
		return false, errors.New("replay detected")
	}
	return true, nil
}

// checkReplay mfa防重放检测
func (s *UserService) checkReplay(uid uint) (bool, error) {
	marked, err := s.cache.IsTOTPTimestepMarked(s.ctx, uid)
	if err != nil {
		return false, err
	}
	// 重放攻击！
	if marked {
		return true, nil
	}

	// 标记该timestep，作为防重放的依据
	err = s.cache.MarkTOTPTimestep(s.ctx, uid)
	if err != nil {
		return false, err
	}
	return false, nil
}
