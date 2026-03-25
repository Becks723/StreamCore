package db

import (
	"StreamCore/internal/pkg/db/chat"
	"StreamCore/internal/pkg/db/group"
	ia "StreamCore/internal/pkg/db/interaction"
	"StreamCore/internal/pkg/db/social"
	"StreamCore/internal/pkg/db/user"
	"StreamCore/internal/pkg/db/video"
	"gorm.io/gorm"
)

type DatabaseSet struct {
	User        user.UserDatabase
	Video       video.VideoDatabase
	Interaction ia.InteractionDatabase
	Social      social.SocialDatabase
	Chat        chat.ChatDatabase
	Group       group.GroupDatabase
}

func NewDatabaseSet(orm *gorm.DB) *DatabaseSet {
	return &DatabaseSet{
		User:        user.NewUserDataBase(orm),
		Video:       video.NewVideoDatabase(orm),
		Interaction: ia.NewInteractionDatabase(orm),
		Social:      social.NewSocialDatabase(orm),
		Chat:        chat.NewChatDatabase(orm),
		Group:       group.NewGroupDatabase(orm),
	}
}
