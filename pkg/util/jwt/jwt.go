package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims[T any] struct {
	Payload T
	jwt.RegisteredClaims
}

func GenerateAccessToken[T any](payload T, secret string, expiresIn time.Duration) (result string, err error) {
	exp := time.Now().Add(expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims[T]{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "StreamCore",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken[T any](payload T, secret string, expiresIn time.Duration) (result string, err error) {
	exp := time.Now().Add(expiresIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims[T]{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "StreamCore",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	return token.SignedString([]byte(secret))
}

func ParseToken[T any](token string, secret string) (payload T, expiresAt time.Time, err error) {
	tk, err := jwt.ParseWithClaims(token, &customClaims[T]{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return
	}

	claims, ok := tk.Claims.(*customClaims[T])
	if !ok {
		err = errors.New("unknown claims type")
		return
	}
	return claims.Payload, claims.ExpiresAt.Time, nil
}
