package service

import (
	"fmt"
	"time"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util/jwt"
)

// RefreshToken 刷新令牌
func (s *UserService) RefreshToken(token string) (*common.TokenInfo, error) {
	// 解析令牌
	tk, expiresAt, err := jwt.ParseToken[*domain.RefreshToken](token, constants.JWT_RefreshSecret)
	// 解析失败
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	// 令牌过期
	if expiresAt.UnixMilli() < time.Now().UnixMilli() {
		return nil, fmt.Errorf("token expired at %v", expiresAt)
	}
	// 校验令牌id
	tkId := s.db.GetTokenId(s.ctx, tk.Uid)
	if tk.Id != tkId {
		return nil, fmt.Errorf("invalid token id")
	}

	// 生成新令牌，刷新令牌也要生成新的，
	// 旧的刷新令牌可以丢弃，也可以继续保存在db里，用来检测攻击者
	access, refresh, err := s.generateTokens(tk.Uid)
	if err != nil {
		return nil, err
	}
	data := new(common.TokenInfo)
	data.AccessToken = access
	data.RefreshToken = refresh
	return data, nil
}
