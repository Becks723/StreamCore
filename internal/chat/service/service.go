package service

import (
	"context"

	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/db/chat"
	"StreamCore/kitex_gen/group/groupservice"
)

type ChatService struct {
	ctx         context.Context
	db          chat.ChatDatabase
	infra       *base.InfraSet
	groupClient groupservice.Client
}

func NewChatService(ctx context.Context, infra *base.InfraSet) *ChatService {
	return &ChatService{
		ctx:         ctx,
		db:          infra.DB.Chat,
		infra:       infra,
		groupClient: infra.GroupClient,
	}
}
