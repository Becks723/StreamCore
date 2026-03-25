package pack

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func User(po *model.UserModel) *domain.User {
	return &domain.User{
		Id:         po.ID,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
		DeletedAt:  packDeletedAt(po.DeletedAt),
		Username:   po.Username,
		Password:   po.Password,
		AvatarUrl:  po.AvatarUrl,
		TOTPBound:  po.TOTPSecret != "",
		TOTPSecret: po.TOTPSecret,
	}
}
