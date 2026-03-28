package user

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

func (repo *userdb) Create(u *domain.User) error {
	var deletedAt gorm.DeletedAt
	if u.DeletedAt != nil {
		deletedAt.Valid = true
		deletedAt.Time = *u.DeletedAt
	} else {
		deletedAt.Valid = false
	}

	po := &model.UserModel{
		Model: gorm.Model{
			ID:        u.Id,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: deletedAt,
		},
		Username:  u.Username,
		Password:  u.Password,
		AvatarUrl: u.AvatarUrl,
	}
	return repo.db.
		Model(&model.UserModel{}).
		Create(&po).
		Error
}
