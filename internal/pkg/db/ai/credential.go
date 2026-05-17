package ai

import (
	"context"

	"StreamCore/internal/pkg/db/model"
)

func (a *aidb) SaveCredential(ctx context.Context, cred *model.CredentialModel) error {
	return a.db.WithContext(ctx).Save(cred).Error
}

func (a *aidb) DeleteCredential(ctx context.Context, credID uint) error {
	return a.db.WithContext(ctx).Delete(&model.CredentialModel{}, credID).Error
}

func (a *aidb) GetCredential(ctx context.Context, userID uint, serviceName string) (*model.CredentialModel, error) {
	var cred model.CredentialModel
	err := a.db.WithContext(ctx).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		First(&cred).Error
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (a *aidb) ListCredentials(ctx context.Context, userID uint) ([]*model.CredentialModel, error) {
	var creds []*model.CredentialModel
	err := a.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&creds).Error
	if err != nil {
		return nil, err
	}
	return creds, nil
}
