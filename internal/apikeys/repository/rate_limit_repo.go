// internal/apikey/repository/rate_limit_repo.go
package repository

import (
	"context"

	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
	"gorm.io/gorm/clause"
)

func (r *repository) UpsertRateLimit(
	ctx context.Context,
	limit *apikeymodelv1.APIKeyRateLimit,
) error {

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "api_key_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"rate_limit",
				"request_limit",
				"burst_limit",
				"updated_at",
			}),
		}).
		Create(limit).Error
}
