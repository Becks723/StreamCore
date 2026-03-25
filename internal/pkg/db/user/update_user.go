package user

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *userdb) UpdateAvatar(id uint, url string) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.
		Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("avatar_url", url).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.User(&po), nil
}

func (repo *userdb) UpdateTOTPSecret(uid uint, secret string) error {
	return repo.db.
		Model(&model.UserModel{}).
		Where("id = ?", uid).
		Update("totp_secret", secret).
		Error
}

func (repo *userdb) UpdateTokenId(ctx context.Context, uid uint, id string) error {
	return repo.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", uid).
		Update("token_id", id).
		Error
}
