package domain

import "time"

type Like struct {
	Uid        uint
	TargetType int
	TargetId   uint
	Status     int // 1-like 2-unlike
	Time       time.Time
}

type Comment struct {
	Id         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	AuthorId   uint
	VideoId    uint
	ParentId   *uint
	Content    string
	LikeCount  int
	ChildCount int
}
