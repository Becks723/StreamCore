package ai

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
)

type AIDatabase interface {
	// Bot user CRUD (operates on users table)
	CreateBotUser(ctx context.Context, bot *model.UserModel) error
	UpdateBotUser(ctx context.Context, bot *model.UserModel) error
	GetBotUser(ctx context.Context, botID uint) (*model.UserModel, error)
	ListBotUsers(ctx context.Context, page, pageSize int) ([]*model.UserModel, int64, error)
	DeleteBotUser(ctx context.Context, botID uint) error

	// Bot-Group binding
	AddBotToGroup(ctx context.Context, botID, groupID uint) error
	RemoveBotFromGroup(ctx context.Context, botID, groupID uint) error
	ListGroupBots(ctx context.Context, groupID uint) ([]*model.UserModel, error)
	ListBotGroups(ctx context.Context, botID uint) ([]uint, error)

	// MCP Server
	CreateServer(ctx context.Context, server *model.MCPServerModel) error
	UpdateServer(ctx context.Context, server *model.MCPServerModel) error
	DeleteServer(ctx context.Context, serverID uint) error
	GetServer(ctx context.Context, serverID uint) (*model.MCPServerModel, error)
	ListServers(ctx context.Context) ([]*model.MCPServerModel, error)

	// MCP Tools
	UpsertTool(ctx context.Context, tool *model.MCPToolModel) error
	DeleteToolsByServer(ctx context.Context, serverID uint) error
	ListTools(ctx context.Context) ([]*model.MCPToolModel, error)
	ListToolsByIDs(ctx context.Context, toolIDs []string) ([]*model.MCPToolModel, error)

	// Credentials
	SaveCredential(ctx context.Context, cred *model.CredentialModel) error
	DeleteCredential(ctx context.Context, credID uint) error
	GetCredential(ctx context.Context, userID uint, serviceName string) (*model.CredentialModel, error)
	ListCredentials(ctx context.Context, userID uint) ([]*model.CredentialModel, error)
}

func NewAIDatabase(gdb *gorm.DB) AIDatabase {
	return &aidb{db: gdb}
}

type aidb struct {
	db *gorm.DB
}
