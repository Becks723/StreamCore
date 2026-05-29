package cache

import (
	"StreamCore/internal/pkg/cache/ai"
	"StreamCore/internal/pkg/cache/group"
	"StreamCore/internal/pkg/cache/interaction"
	"StreamCore/internal/pkg/cache/social"
	"StreamCore/internal/pkg/cache/user"
	"StreamCore/internal/pkg/cache/video"
	"github.com/redis/go-redis/v9"
)

type CacheSet struct {
	AI          ai.AICache
	User        user.UserCache
	Video       video.VideoCache
	Interaction interaction.InteractionCache
	Social      social.SocialCache
	Group       group.GroupCache
}

func NewCacheSet(rdb *redis.Client) *CacheSet {
	return &CacheSet{
		AI:          ai.NewAICache(rdb),
		User:        user.NewUserCache(rdb),
		Video:       video.NewVideoCache(rdb),
		Interaction: interaction.NewInteractionCache(rdb),
		Social:      social.NewSocialCache(rdb),
		Group:       group.NewGroupCache(rdb),
	}
}
