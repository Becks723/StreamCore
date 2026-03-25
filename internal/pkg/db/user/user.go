package user

import (
	"context"

	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

type UserDatabase interface {
	Create(u *domain.User) error
	GetByUsername(username string) (u *domain.User, err error)
	GetById(id uint) (u *domain.User, err error)
	GetTokenId(ctx context.Context, uid uint) string
	UpdateAvatar(id uint, url string) (u *domain.User, err error)
	UpdateTOTPSecret(uid uint, secret string) error
	UpdateTokenId(ctx context.Context, uid uint, id string) error
}

func NewUserDataBase(gdb *gorm.DB) UserDatabase {
	return &userdb{
		db: gdb,
	}
}

type userdb struct {
	db *gorm.DB
}
