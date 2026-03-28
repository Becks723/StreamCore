package pack

import (
	"StreamCore/internal/pkg/domain"
	"StreamCore/kitex_gen/common"
	"StreamCore/pkg/util"
)

func SocialUserInfo(u *domain.User) *common.SocialUserInfo {
	return &common.SocialUserInfo{
		Id:        util.Uint2String(u.Id),
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
}
