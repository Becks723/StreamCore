package model

import (
	"time"

	"gorm.io/gorm"
)

type LikeRelationModel struct {
	Id         uint `gorm:"primaryKey"`
	Uid        uint
	User       *UserModel `gorm:"foreignKey:Uid;reference:ID"`
	TargetType int        // 1-video 2-comment    多态关联
	TargetId   uint
	Status     int // 1-like 2-unlike
	Time       time.Time
}

type LikeCountModel struct {
	Id          uint `gorm:"primaryKey"`
	TargetType  int  // 1-video 2-comment
	TargetId    uint
	LikeCount   int64
	UnlikeCount int64
}

type CommentModel struct {
	gorm.Model

	AuthorId   uint
	Author     *UserModel `gorm:"foreignKey:AuthorId;references:ID"`
	VideoId    uint
	Video      *VideoModel `gorm:"foreignKey:VideoId;references:ID"`
	Content    string
	ParentId   *uint // parent comment id, null if root
	LikeCount  int
	ChildCount int
}
