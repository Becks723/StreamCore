package pack

import (
	"StreamCore/internal/pkg/db/model"
	"StreamCore/internal/pkg/domain"
)

func Follow(po *model.FollowModel) *domain.Follow {
	return &domain.Follow{
		FollowerId: po.FollowerId,
		FolloweeId: po.FolloweeId,
		Status:     po.Status,
		Time:       po.Time,
	}
}
