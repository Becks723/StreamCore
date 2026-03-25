package interaction

import (
	"context"
	"time"

	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

type InteractionDatabase interface {
	CreateLike(ctx context.Context, tarType int, uid, tarId uint, time time.Time) error
	GetLike(ctx context.Context, tarType int, uid, tarId uint) (*domain.Like, error)
	ToggleLikeStatus(ctx context.Context, tarType int, uid, tarId uint) error
	FetchUserLikedVideos(ctx context.Context, uid uint, limit, offset int) ([]*domain.Like, error)
	CreateComment(ctx context.Context, c *domain.Comment) error
	GetCommentById(cid uint) (*domain.Comment, error)
	ListRootComments(vid uint, limit, page int) ([]*domain.Comment, error)
	ListSubComments(cid uint, limit, page int) ([]*domain.Comment, error)
	DeleteCommentById(cid, authorId uint) error
}

func NewInteractionDatabase(gdb *gorm.DB) InteractionDatabase {
	return &iactiondb{
		db: gdb,
	}
}

type iactiondb struct {
	db *gorm.DB
}
