package ai

import (
	"context"

	"StreamCore/internal/pkg/db/model"
)

func (a *aidb) CreateBotUser(ctx context.Context, bot *model.UserModel) error {
	return a.db.WithContext(ctx).Create(bot).Error
}

func (a *aidb) UpdateBotUser(ctx context.Context, bot *model.UserModel) error {
	return a.db.WithContext(ctx).Save(bot).Error
}

func (a *aidb) DeleteBotUser(ctx context.Context, botID uint) error {
	return a.db.WithContext(ctx).Delete(&model.UserModel{}, botID).Error
}

func (a *aidb) GetBotUser(ctx context.Context, botID uint) (*model.UserModel, error) {
	var user model.UserModel
	err := a.db.WithContext(ctx).Where("id = ? AND is_bot = ?", botID, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *aidb) ListBotUsers(ctx context.Context, page, pageSize int) ([]*model.UserModel, int64, error) {
	var users []*model.UserModel
	var total int64

	db := a.db.WithContext(ctx).Where("is_bot = ?", true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (a *aidb) AddBotToGroup(ctx context.Context, botID, groupID uint) error {
	binding := &model.BotGroupModel{
		BotID:   botID,
		GroupID: groupID,
		Status:  0,
	}
	return a.db.WithContext(ctx).Create(binding).Error
}

func (a *aidb) RemoveBotFromGroup(ctx context.Context, botID, groupID uint) error {
	return a.db.WithContext(ctx).
		Model(&model.BotGroupModel{}).
		Where("bot_id = ? AND group_id = ?", botID, groupID).
		Update("status", 1).Error
}

func (a *aidb) ListGroupBots(ctx context.Context, groupID uint) ([]*model.UserModel, error) {
	var users []*model.UserModel
	err := a.db.WithContext(ctx).
		Table("users").
		Joins("JOIN bot_groups ON users.id = bot_groups.bot_id").
		Where("bot_groups.group_id = ? AND bot_groups.status = ? AND users.is_bot = ?", groupID, 0, true).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *aidb) ListBotGroups(ctx context.Context, botID uint) ([]uint, error) {
	var bindings []model.BotGroupModel
	err := a.db.WithContext(ctx).
		Where("bot_id = ? AND status = ?", botID, 0).
		Find(&bindings).Error
	if err != nil {
		return nil, err
	}
	groupIDs := make([]uint, len(bindings))
	for i, b := range bindings {
		groupIDs[i] = b.GroupID
	}
	return groupIDs, nil
}
