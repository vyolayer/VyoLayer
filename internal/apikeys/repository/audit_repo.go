package repository

import (
	"context"

	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
)

func (r *repository) CreateAuditLog(
	ctx context.Context,
	log *apikeymodelv1.APIKeyAuditLog,
) error {
	return r.db.WithContext(ctx).
		Create(log).Error
}
