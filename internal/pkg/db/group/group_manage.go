package group

import (
	"context"
	"errors"

	"StreamCore/internal/pkg/constants"
	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
)

func (repo *groupdb) CreateGroup(ctx context.Context, ownerUid uint, name string) (uint, error) {
	g := &model.ChatGroupModel{
		Name:     name,
		OwnerUid: ownerUid,
	}
	if err := repo.db.WithContext(ctx).Create(g).Error; err != nil {
		return 0, err
	}
	return g.Id, nil
}

func (repo *groupdb) GroupExists(ctx context.Context, groupId uint) (bool, error) {
	var cnt int64
	err := repo.db.WithContext(ctx).Model(&model.ChatGroupModel{}).
		Where("id = ?", groupId).
		Count(&cnt).
		Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (repo *groupdb) AddGroupMember(ctx context.Context, groupId, uid uint, role int) error {
	var po model.ChatGroupMemberModel
	err := repo.db.WithContext(ctx).
		Model(&model.ChatGroupMemberModel{}).
		Where("group_id = ? AND user_id = ?", groupId, uid).
		First(&po).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return repo.db.WithContext(ctx).Create(&model.ChatGroupMemberModel{
			GroupId: groupId,
			UserId:  uid,
			Role:    role,
			Status:  constants.ChatGroupMemberStatusActive,
		}).Error
	}

	po.Role = role
	po.Status = constants.ChatGroupMemberStatusActive
	return repo.db.WithContext(ctx).Save(&po).Error
}

func (repo *groupdb) IsGroupMember(ctx context.Context, groupId, uid uint) (bool, error) {
	var cnt int64
	err := repo.db.WithContext(ctx).Model(&model.ChatGroupMemberModel{}).
		Where("group_id = ? AND user_id = ? AND status = ?", groupId, uid, constants.ChatGroupMemberStatusActive).
		Count(&cnt).
		Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (repo *groupdb) ListGroupMemberIds(ctx context.Context, groupId uint) ([]uint, error) {
	var ids []uint
	err := repo.db.WithContext(ctx).Model(&model.ChatGroupMemberModel{}).
		Select("user_id").
		Where("group_id = ? AND status = ?", groupId, constants.ChatGroupMemberStatusActive).
		Scan(&ids).
		Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (repo *groupdb) GetGroupMemberRole(ctx context.Context, groupId, uid uint) (int, bool, error) {
	var po model.ChatGroupMemberModel
	err := repo.db.WithContext(ctx).
		Model(&model.ChatGroupMemberModel{}).
		Where("group_id = ? AND user_id = ? AND status = ?", groupId, uid, constants.ChatGroupMemberStatusActive).
		First(&po).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return po.Role, true, nil
}

func (repo *groupdb) CreateGroupApply(ctx context.Context, groupId, applicantUid uint, reason string, ts int64) (int64, error) {
	po := &model.ChatGroupApplyModel{
		GroupId:      groupId,
		ApplicantUid: applicantUid,
		Reason:       reason,
		Status:       constants.ChatGroupApplyStatusPending,
		CreatedAtTs:  ts,
		UpdatedAtTs:  ts,
	}
	if err := repo.db.WithContext(ctx).Create(po).Error; err != nil {
		return 0, err
	}
	return po.Id, nil
}

func (repo *groupdb) HasPendingGroupApply(ctx context.Context, groupId, applicantUid uint) (bool, error) {
	var cnt int64
	err := repo.db.WithContext(ctx).Model(&model.ChatGroupApplyModel{}).
		Where("group_id = ? AND applicant_uid = ? AND status = ?", groupId, applicantUid, constants.ChatGroupApplyStatusPending).
		Count(&cnt).
		Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (repo *groupdb) GetGroupApplyByID(ctx context.Context, applyId int64) (*model.ChatGroupApplyModel, error) {
	var po model.ChatGroupApplyModel
	err := repo.db.WithContext(ctx).
		Model(&model.ChatGroupApplyModel{}).
		Where("id = ?", applyId).
		First(&po).
		Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (repo *groupdb) UpdateGroupApplyStatus(ctx context.Context, applyId int64, status int, ts int64) error {
	return repo.db.WithContext(ctx).
		Model(&model.ChatGroupApplyModel{}).
		Where("id = ?", applyId).
		Updates(map[string]interface{}{
			"status":        status,
			"updated_at_ts": ts,
		}).
		Error
}

func (repo *groupdb) RunInTx(ctx context.Context, fn func(tx GroupDatabase) error) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&groupdb{db: tx})
	})
}
