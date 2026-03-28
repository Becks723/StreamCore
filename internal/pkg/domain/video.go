package domain

import "time"

type Video struct {
	Id           uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Author       *User
	AuthorId     uint
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

// VideoQuery - for es querying
type VideoQuery struct {
	TitleMatches    string
	DescMatches     string
	FromDate        *string
	ToDate          *string
	AuthorIdIsExact *uint
	UsernameMatches *string
}
