package repository

import (
	"context"

	"github.com/google/uuid"
	apikeymodelv1 "github.com/vyolayer/vyolayer/internal/apikeys/models/v1"
)

func (r *repository) ListScopes(
	ctx context.Context,
	apiKeyID uuid.UUID,
) ([]string, error) {

	var rows []apikeymodelv1.APIKeyScope

	err := r.db.WithContext(ctx).
		Where("api_key_id = ?", apiKeyID).
		Order("scope asc").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.Scope)
	}

	return out, nil
}
