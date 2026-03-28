package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model

	Username   string `gorm:"column:username;unique"` // username should be unique
	Password   string `gorm:"column:password"`
	AvatarUrl  string `gorm:"column:avatar_url"` // gorm default maps struct names as snake_case, so this column tag is optional.
	TOTPSecret string `gorm:"column:totp_secret"`
	TokenId    string `gorm:"column:token_id"`
}
