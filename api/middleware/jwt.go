package middleware

import (
	"context"
	"errors"
	"time"

	"StreamCore/api/pack"
	"StreamCore/internal/pkg/base/rpccontext"
	"StreamCore/internal/pkg/constants"
	"StreamCore/pkg/util/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func JWTAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		const (
			AccessTokenKey  = "Access-Token"
			RefreshTokenKey = "Refresh-Token"
		)
		access := string(c.GetHeader(AccessTokenKey))
		refresh := string(c.GetHeader(RefreshTokenKey))

		if access == "" && refresh == "" {
			pack.RespUnauthorizedError(c, errors.New("token not provided"))
			c.Abort()
			return
		}

		uid, expiresAt, err := jwt.ParseToken[uint](access, constants.JWT_AccessSecret)
		// access failed
		if err != nil {
			pack.RespUnauthorizedError(c, errors.New("failed to resolve token"))
			c.Abort()
			return
		}
		// access expired
		if expiresAt.UnixMilli() < time.Now().UnixMilli() {
			pack.RespUnauthorizedError(c, errors.New("token expired"))
			c.Abort()
			return
		}

		ctx = rpccontext.WithLoginUid(ctx, uid)
		c.Next(ctx)
	}
}
