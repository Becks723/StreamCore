package video

import (
	"context"
	"time"

	"StreamCore/internal/pkg/domain"
	"gorm.io/gorm"
)

type VideoDatabase interface {
	Create(v *domain.Video) error
	GetById(vid uint) (*domain.Video, error)
	Fetch(after *time.Time) ([]*domain.Video, error)
	FetchByUid(uid uint, limit, page int) ([]*domain.Video, int, error)
	Search(keywords string, limit, page int, from, to *time.Time, username *string) ([]*domain.Video, int, error)
	BatchUpdateVisits(ctx context.Context, batch map[uint]int64) error
	FetchVideoIdsByVisit(ctx context.Context, limit, page int) (map[uint]int64, error)
}

func NewVideoDatabase(gdb *gorm.DB) VideoDatabase {
	return &videodb{
		db: gdb,
	}
}

type videodb struct {
	db *gorm.DB
}
