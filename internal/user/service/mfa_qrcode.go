package service

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strconv"

	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/user"
	"github.com/pquerna/otp/totp"
)

// MFAQrcode 生成mfa的密钥和二维码
func (s *UserService) MFAQrcode(uid uint) (*user.MFAQrcodeInfo, error) {
	var err error

	// generate totp secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "StreamCore",
		AccountName: strconv.FormatUint(uint64(uid), 10),
	})
	if err != nil {
		return nil, err
	}
	secret := key.Secret()

	// generate totp qrcode
	w := constants.MFAQrcodeWidth
	h := constants.MFAQrcodeHeight
	img, err := key.Image(w, h)
	if err != nil {
		return nil, err
	}
	// encode img to png (binary)
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return nil, err
	}
	// base64 wrapper
	qrcode := base64.StdEncoding.EncodeToString(buf.Bytes())

	// cache secret
	err = s.cache.SetTOTPPending(s.ctx, uid, secret, constants.TOTPSecretExpiry)
	if err != nil {
		return nil, err
	}

	data := new(user.MFAQrcodeInfo)
	data.Secret = secret
	data.Qrcode = qrcode
	return data, nil
}
