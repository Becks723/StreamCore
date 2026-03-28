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

func TestUserService_generateTokens(t *testing.T) {
	type testCase struct {
		name                 string
		expectingError       bool
		expectedErrorMsg     string
		expectedAccessToken  string
		expectedRefreshToken string
		genAccessError       error
		genRefreshError      error
		updateTokenIdError   error
	}

	testCases := []testCase{
		{
			name:                 "成功",
			expectingError:       false,
			expectedAccessToken:  "access_token",
			expectedRefreshToken: "refresh_token",
		},
		{
			name:             "生成访问令牌失败",
			expectingError:   true,
			expectedErrorMsg: "failed gen accessToken:",
			genAccessError:   errors.New("internal error"),
		},
		{
			name:             "生成刷新令牌失败",
			expectingError:   true,
			expectedErrorMsg: "failed gen refreshToken:",
			genRefreshError:  errors.New("internal error"),
		},
		{
			name:               "db写入token_id失败",
			expectingError:     true,
			expectedErrorMsg:   "update token id failed:",
			updateTokenIdError: errors.New("internal error"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// mock
			mockey.MockGeneric(jwt.GenerateAccessToken[uint]).
				To(func(payload uint, secret string, expiresIn time.Duration) (string, error) {
					if tc.genAccessError != nil {
						return "", tc.genAccessError
					}
					return tc.expectedAccessToken, nil
				}).Build()

			mockey.MockGeneric(jwt.GenerateRefreshToken[*domain.RefreshToken]).
				To(func(payload *domain.RefreshToken, secret string, expiresIn time.Duration) (string, error) {
					if tc.genRefreshError != nil {
						return "", tc.genRefreshError
					}
					return tc.expectedRefreshToken, nil
				}).Build()

			_UpdateTokenId := mockey.GetMethod(userDB.NewUserDataBase(nil), "UpdateTokenId")
			mockey.Mock(_UpdateTokenId).
				To(func(ctx context.Context, uid uint, id string) error {
					if tc.updateTokenIdError != nil {
						return tc.updateTokenIdError
					}
					return nil
				}).Build()

			// init
			mockInfraSet := &base.InfraSet{
				DB:    db.NewDatabaseSet(nil),
				Cache: cache.NewCacheSet(nil),
			}
			userService := NewUserService(context.Background(), mockInfraSet)

			// call
			access, refresh, err := userService.generateTokens(234)

			// assert
			if tc.expectingError {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.expectedErrorMsg)
			} else {
				So(err, ShouldBeNil)
				So(access, ShouldEqual, tc.expectedAccessToken)
				So(refresh, ShouldEqual, tc.expectedRefreshToken)
			}
		})
	}
}
