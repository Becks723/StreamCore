package service

import (
	"context"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"testing"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/cache"
	userCache "StreamCore/internal/pkg/cache/user"
	"StreamCore/internal/pkg/db"
	"github.com/bytedance/mockey"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUserService_MFAQrcode(t *testing.T) {
	type testCase struct {
		name             string
		expectingError   bool
		expectedErrorMsg string
		expectedSecret   string
		expectedQrcode   string
		totpGenError     error
		imageGenError    error
		pngEncodeError   error
		cacheSetError    error
	}

	testCases := []testCase{
		{
			name:           "成功",
			expectingError: false,
			expectedSecret: "secret",
			expectedQrcode: "qrcode",
		},
		{
			name:           "totp密钥生成失败",
			expectingError: true,
			totpGenError:   errors.New("internal error"),
		},
		{
			name:           "二维码图片生成失败",
			expectingError: true,
			imageGenError:  errors.New("internal error"),
		},
		{
			name:           "PNG编码失败",
			expectingError: true,
			pngEncodeError: errors.New("internal error"),
		},
		{
			name:           "缓存写入secret失败",
			expectingError: true,
			cacheSetError:  errors.New("internal error"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// mock
			mockey.Mock(totp.Generate).
				To(func(opts totp.GenerateOpts) (*otp.Key, error) {
					if tc.totpGenError != nil {
						return nil, tc.totpGenError
					}
					return &otp.Key{}, nil
				}).Build()

			mockey.Mock((*otp.Key).Secret).
				Return(tc.expectedSecret).Build()

			mockey.Mock((*otp.Key).Image).
				To(func(width int, height int) (image.Image, error) {
					if tc.imageGenError != nil {
						return nil, tc.imageGenError
					}
					return image.NewRGBA(image.Rect(0, 0, 10, 10)), nil
				}).Build()

			mockey.Mock(png.Encode).
				Return(tc.pngEncodeError).Build()

			mockey.Mock((*base64.Encoding).EncodeToString).
				Return(tc.expectedQrcode).Build()

			_SetTOTPPending := mockey.GetMethod(userCache.NewUserCache(nil), "SetTOTPPending")
			mockey.Mock(_SetTOTPPending).
				Return(tc.cacheSetError).Build()

			// init
			mockInfraSet := &base.InfraSet{
				DB:    db.NewDatabaseSet(nil),
				Cache: cache.NewCacheSet(nil),
			}
			userService := NewUserService(context.Background(), mockInfraSet)

			// call
			qr, err := userService.MFAQrcode(234)

			// assert
			if tc.expectingError {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.expectedErrorMsg)
			} else {
				So(err, ShouldBeNil)
				So(qr.Secret, ShouldEqual, tc.expectedSecret)
				So(qr.Qrcode, ShouldEqual, tc.expectedQrcode)
			}
		})
	}
}
