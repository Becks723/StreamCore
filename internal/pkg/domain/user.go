package domain

import (
	"time"
)

type User struct {
	Id         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Username   string
	Password   string
	AvatarUrl  string
	TOTPBound  bool
	TOTPSecret string
}

func (u *User) GetId() uint {
	return u.Id
}

func (u *User) GetUsername() string {
	return u.Username
}

type RefreshToken struct {
	Uid uint
	Id  string // token id
}
