package model

import (
	"time"

	"gorm.io/gorm"
)

type VideoModel struct {
	gorm.Model

	AuthorId     uint
	Author       *UserModel `gorm:"foreignKey:AuthorId;references:ID"`
	VideoUrl     string
	CoverUrl     string
	Title        string
	Description  string
	VisitCount   int
	LikeCount    int
	CommentCount int
	PublishedAt  time.Time
	EditedAt     time.Time
}

type VisitCountModel struct {
	Vid        uint `gorm:"primaryKey"`
	VisitCount int64
}
