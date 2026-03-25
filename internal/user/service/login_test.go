package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/cache"
	userCache "StreamCore/internal/pkg/cache/user"
	"StreamCore/internal/pkg/db"
	userDB "StreamCore/internal/pkg/db/user"
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/user"
	"StreamCore/pkg/util"
	"github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserService_Login(t *testing.T) {
	type testCase struct {
		Name                 string
		ExpectingError       bool
		ExpectedErrorMsg     string
		Username             string
		Password             string
		ExpectedUser         *domain.User
		MFAEnabled           bool
		ExpectedAccessToken  string
		ExpectedRefreshToken string
		DBGetUserError       error
		MFATokenSetError     error
		GenTokensError       error
	}

	// user data
	pass, _ := util.EncryptPassword("admin")
	mfaEnabledUser := domain.User{
		Id:        234,
		Username:  "admin",
		Password:  pass,
		TOTPBound: true,
	}
	mfaDisabledUser := mfaEnabledUser
	mfaDisabledUser.TOTPBound = false

	testCases := []testCase{
		{
			Name:                 "登录成功-未开启MFA",
			ExpectingError:       false,
			Username:             "admin",
			Password:             "admin",
			MFAEnabled:           false,
			ExpectedUser:         &mfaDisabledUser,
			ExpectedAccessToken:  "access_token",
			ExpectedRefreshToken: "refresh_token",
		},
		{
			Name:                 "登录成功-已开启MFA",
			ExpectingError:       false,
			Username:             "admin",
			Password:             "admin",
			MFAEnabled:           true,
			ExpectedUser:         &mfaEnabledUser,
			ExpectedAccessToken:  "access_token",
			ExpectedRefreshToken: "refresh_token",
		},
		{
			Name:             "用户不存在",
			ExpectingError:   true,
			Username:         "admin111",
			Password:         "admin",
			DBGetUserError:   errors.New("internal error"),
			ExpectedErrorMsg: "用户不存在",
		},
		{
			Name:             "密码错误",
			ExpectingError:   true,
			Username:         "admin",
			Password:         "123456",
			ExpectedErrorMsg: "密码错误",
			ExpectedUser:     &mfaDisabledUser,
		},
		{
			Name:             "MFAToken写入失败",
			ExpectingError:   true,
			Username:         "admin",
			Password:         "admin",
			MFATokenSetError: errors.New("internal error"),
			ExpectedErrorMsg: "failed cache.SetMFATokenUser:",
			ExpectedUser:     &mfaEnabledUser,
		},
		{
			Name:             "生成jwt令牌失败",
			ExpectingError:   true,
			Username:         "admin",
			Password:         "admin",
			GenTokensError:   errors.New("internal error"),
			ExpectedErrorMsg: "",
			ExpectedUser:     &mfaDisabledUser,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// mock
			_SetMFATokenUser := mockey.GetMethod(userCache.NewUserCache(nil), "SetMFATokenUser")
			mockey.Mock(_SetMFATokenUser).
				To(func(ctx context.Context, token string, uid uint, ttl time.Duration) error {
					if tc.MFATokenSetError != nil {
						return tc.MFATokenSetError
					}
					return nil
				}).Build()

			_GetByUsername := mockey.GetMethod(userDB.NewUserDataBase(nil), "GetByUsername")
			mockey.Mock(_GetByUsername).
				To(func(username string) (*domain.User, error) {
					if tc.DBGetUserError != nil {
						return nil, tc.DBGetUserError
					}
					return tc.ExpectedUser, nil
				}).Build()

			mockey.Mock((*UserService).generateTokens).
				To(func(uid uint) (string, string, error) {
					if tc.GenTokensError != nil {
						return "", "", tc.GenTokensError
					}
					return tc.ExpectedAccessToken, tc.ExpectedRefreshToken, nil
				}).Build()

			// init
			mockInfraSet := &base.InfraSet{
				DB:    db.NewDatabaseSet(nil),
				Cache: cache.NewCacheSet(nil),
			}
			userService := NewUserService(context.Background(), mockInfraSet)

			// call
			info, auth, token, err := userService.Login(&user.LoginReq{
				Username: tc.Username,
				Password: tc.Password,
			})

			// assert
			if tc.ExpectingError {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.ExpectedErrorMsg)
			} else {
				So(err, ShouldBeNil)
				So(info.Id, ShouldEqual, "234")
				So(info.Username, ShouldEqual, "admin")
				if tc.MFAEnabled {
					So(token, ShouldBeNil)
					So(auth.MfaRequired, ShouldBeTrue)
					So(auth.MfaToken, ShouldNotBeEmpty)
				} else {
					So(token, ShouldNotBeNil)
					So(auth.MfaRequired, ShouldBeFalse)
					So(token.AccessToken, ShouldEqual, tc.ExpectedAccessToken)
					So(token.RefreshToken, ShouldEqual, tc.ExpectedRefreshToken)
				}
			}
		})
	}
}
