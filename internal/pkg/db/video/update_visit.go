package video

import (
	"context"

	"StreamCore/internal/pkg/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *videodb) BatchUpdateVisits(ctx context.Context, batch map[uint]int64) error {
	var po []*model.VisitCountModel
	for k, v := range batch {
		po = append(po, &model.VisitCountModel{
			Vid:        k,
			VisitCount: v,
		})
	}
	// upsert
	return repo.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "vid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"visit_count": gorm.Expr("visit_count + VALUES(visit_count)"),
		}),
	}).Create(&po).Error
}
