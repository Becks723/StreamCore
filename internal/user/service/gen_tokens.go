package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/domain"
	"StreamCore/pkg/util/jwt"
)

// generateTokens 生成访问和刷新令牌
func (s *UserService) generateTokens(uid uint) (string, string, error) {
	// 生成访问令牌
	access, err := jwt.GenerateAccessToken(uid, constants.JWT_AccessSecret, constants.JWT_AccessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed gen accessToken: %w", err)
	}

	// 生成刷新令牌
	tokenId := s.generateTokenId() // 生成令牌id
	refresh, err := jwt.GenerateRefreshToken(&domain.RefreshToken{
		Uid: uid,
		Id:  tokenId,
	}, constants.JWT_RefreshSecret, constants.JWT_RefreshTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed gen refreshToken: %w", err)
	}
	// 令牌id存入db
	if err = s.db.UpdateTokenId(s.ctx, uid, tokenId); err != nil {
		return "", "", fmt.Errorf("update token id failed: %w", err)
	}

	return access, refresh, nil
}

func (s *UserService) generateTokenId() string {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return base64.RawStdEncoding.EncodeToString(buf)
}
