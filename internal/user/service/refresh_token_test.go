package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/cache"
	"StreamCore/internal/pkg/db"
	userDB "StreamCore/internal/pkg/db/user"
	"StreamCore/internal/pkg/domain"
	"StreamCore/pkg/util/jwt"
	"github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserService_RefreshToken(t *testing.T) {
	type testCase struct {
		name                 string
		expectingError       bool
		expectedErrorMsg     string
		expectedAccessToken  string
		expectedRefreshToken string
		parseTokenError      error
		tokenExpired         bool
		expectedToken        *domain.RefreshToken
		tokenId              string
		genTokensError       error
	}

	token := &domain.RefreshToken{
		Uid: 234,
		Id:  "fake_token_id",
	}
	testCases := []testCase{
		{
			name:                 "成功",
			expectingError:       false,
			expectedAccessToken:  "access_token",
			expectedRefreshToken: "refresh_token",
			expectedToken:        token,
			tokenId:              token.Id,
		},
		{
			name:             "token解析失败",
			expectingError:   true,
			expectedErrorMsg: "invalid token:",
			parseTokenError:  errors.New("internal error"),
		},
		{
			name:             "token过期",
			expectingError:   true,
			expectedErrorMsg: "token expired at ",
			expectedToken:    token,
			tokenExpired:     true,
		},
		{
			name:             "非法token id",
			expectingError:   true,
			expectedErrorMsg: "invalid token id",
			expectedToken:    token,
			tokenId:          "invalid_token_id",
		},
		{
			name:             "生成jwt令牌失败",
			expectingError:   true,
			expectedErrorMsg: "",
			expectedToken:    token,
			tokenId:          token.Id,
			genTokensError:   errors.New("internal error"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// mock
			mockey.MockGeneric(jwt.ParseToken[*domain.RefreshToken]).
				To(func(token string, secret string) (*domain.RefreshToken, time.Time, error) {
					if tc.parseTokenError != nil {
						return nil, time.Time{}, tc.parseTokenError
					}
					if tc.tokenExpired {
						return tc.expectedToken, time.Now().Add(-time.Hour), nil
					} else {
						return tc.expectedToken, time.Now().Add(time.Hour), nil
					}
				}).Build()

			_GetTokenId := mockey.GetMethod(userDB.NewUserDataBase(nil), "GetTokenId")
			mockey.Mock(_GetTokenId).
				Return(tc.tokenId).Build()

			mockey.Mock((*UserService).generateTokens).
				To(func(uid uint) (string, string, error) {
					if tc.genTokensError != nil {
						return "", "", tc.genTokensError
					}
					return tc.expectedAccessToken, tc.expectedRefreshToken, nil
				}).Build()

			// init
			mockInfraSet := &base.InfraSet{
				DB:    db.NewDatabaseSet(nil),
				Cache: cache.NewCacheSet(nil),
			}
			userService := NewUserService(context.Background(), mockInfraSet)

			// call
			info, err := userService.RefreshToken("old_refresh_token")

			// assert
			if tc.expectingError {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.expectedErrorMsg)
				So(info, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)
				So(info.AccessToken, ShouldEqual, tc.expectedAccessToken)
				So(info.RefreshToken, ShouldEqual, tc.expectedRefreshToken)
			}
		})
	}
}
