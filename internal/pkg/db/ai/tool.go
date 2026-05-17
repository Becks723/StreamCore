package ai

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm/clause"
)

func (a *aidb) CreateServer(ctx context.Context, server *model.MCPServerModel) error {
	return a.db.WithContext(ctx).Create(server).Error
}

func (a *aidb) UpdateServer(ctx context.Context, server *model.MCPServerModel) error {
	return a.db.WithContext(ctx).Save(server).Error
}

func (a *aidb) DeleteServer(ctx context.Context, serverID uint) error {
	return a.db.WithContext(ctx).Delete(&model.MCPServerModel{}, serverID).Error
}

func (a *aidb) GetServer(ctx context.Context, serverID uint) (*model.MCPServerModel, error) {
	var server model.MCPServerModel
	err := a.db.WithContext(ctx).First(&server, serverID).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (a *aidb) ListServers(ctx context.Context) ([]*model.MCPServerModel, error) {
	var servers []*model.MCPServerModel
	err := a.db.WithContext(ctx).Order("id ASC").Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (a *aidb) UpsertTool(ctx context.Context, tool *model.MCPToolModel) error {
	return a.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tool_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tool_name", "description", "input_schema", "server_id", "updated_at"}),
	}).Create(tool).Error
}

func (a *aidb) DeleteToolsByServer(ctx context.Context, serverID uint) error {
	return a.db.WithContext(ctx).Where("server_id = ?", serverID).Delete(&model.MCPToolModel{}).Error
}

func (a *aidb) ListTools(ctx context.Context) ([]*model.MCPToolModel, error) {
	var tools []*model.MCPToolModel
	err := a.db.WithContext(ctx).Order("id ASC").Find(&tools).Error
	if err != nil {
		return nil, err
	}
	return tools, nil
}

func (a *aidb) ListToolsByIDs(ctx context.Context, toolIDs []string) ([]*model.MCPToolModel, error) {
	var tools []*model.MCPToolModel
	err := a.db.WithContext(ctx).Where("tool_id IN ?", toolIDs).Find(&tools).Error
	if err != nil {
		return nil, err
	}
	return tools, nil
}
