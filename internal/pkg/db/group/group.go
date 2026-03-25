package group

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
)

type GroupDatabase interface {
	CreateGroup(ctx context.Context, ownerUid uint, name string) (groupId uint, err error)
	GroupExists(ctx context.Context, groupId uint) (bool, error)
	AddGroupMember(ctx context.Context, groupId, uid uint, role int) error
	IsGroupMember(ctx context.Context, groupId, uid uint) (bool, error)
	ListGroupMemberIds(ctx context.Context, groupId uint) ([]uint, error)
	GetGroupMemberRole(ctx context.Context, groupId, uid uint) (role int, isMember bool, err error)

	CreateGroupApply(ctx context.Context, groupId, applicantUid uint, reason string, ts int64) (applyId int64, err error)
	HasPendingGroupApply(ctx context.Context, groupId, applicantUid uint) (bool, error)
	GetGroupApplyByID(ctx context.Context, applyId int64) (*model.ChatGroupApplyModel, error)
	UpdateGroupApplyStatus(ctx context.Context, applyId int64, status int, ts int64) error
	RunInTx(ctx context.Context, fn func(tx GroupDatabase) error) error
}

func NewGroupDatabase(gdb *gorm.DB) GroupDatabase {
	return &groupdb{
		db: gdb,
	}
}

type groupdb struct {
	db *gorm.DB
}
