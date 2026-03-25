package pack

import (
	"time"

	"gorm.io/gorm"
)

func packDeletedAt(d gorm.DeletedAt) *time.Time {
	if d.Valid {
		return &d.Time
	}
	return nil
}
