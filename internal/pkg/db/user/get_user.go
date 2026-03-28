package user

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/db/pack"
	"StreamCore/internal/pkg/domain"
)

func (repo *userdb) GetByUsername(username string) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.Model(&model.UserModel{}).
		Where("username = ?", username).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.User(&po), nil
}

func (repo *userdb) GetById(id uint) (u *domain.User, err error) {
	po := model.UserModel{}
	err = repo.db.
		Where("id = ?", id).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return pack.User(&po), nil
}

func (repo *userdb) GetTokenId(ctx context.Context, uid uint) string {
	tokenId := ""
	err := repo.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Select("token_id").
		Where("id = ?", uid).
		Scan(&tokenId).
		Error
	if err != nil {
		// TODO: log
		return ""
	}
	return tokenId
}
