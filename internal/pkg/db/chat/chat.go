package chat

import (
	"context"

	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

type ChatDatabase interface {
	CreateWhisperMessage(ctx context.Context, msg *domain.WhisperMessage) error
	CreateGroupMessage(ctx context.Context, msg *domain.GroupMessage) error

	ListWhisperMessages(ctx context.Context, uid, peer uint, pageSize int, cursorMsgId int64) (items []*domain.WhisperMessage, hasMore bool, nextCursorMsgId int64, err error)
	ListWhisperMessagesAll(ctx context.Context, uid, peer uint) ([]*domain.WhisperMessage, error)

	ListGroupMessages(ctx context.Context, groupId uint, pageSize int, cursorMsgId int64) (items []*domain.GroupMessage, hasMore bool, nextCursorMsgId int64, err error)
	ListGroupMessagesAll(ctx context.Context, groupId uint) ([]*domain.GroupMessage, error)

	IsFriend(ctx context.Context, uid, peer uint) (bool, error)
}

func NewChatDatabase(gdb *gorm.DB) ChatDatabase {
	return &chatdb{
		db: gdb,
	}
}

type chatdb struct {
	db *gorm.DB
}
