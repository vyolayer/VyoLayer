package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/console/model"
)

type OverrideRepository interface {
	ListOverrides(ctx context.Context, projectID uuid.UUID, resourceIDs []uint64) ([]model.ProjectResourceOverride, error)
}

type overrideRepository struct {
	db *gorm.DB
}

func NewOverrideRepository(db *gorm.DB) OverrideRepository {
	return &overrideRepository{db: db}
}

func (r *overrideRepository) ListOverrides(ctx context.Context, projectID uuid.UUID, resourceIDs []uint64) ([]model.ProjectResourceOverride, error) {
	var overrides []model.ProjectResourceOverride
	if len(resourceIDs) == 0 {
		return overrides, nil
	}

	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("resource_id IN ?", resourceIDs).
		Find(&overrides).Error
	return overrides, err
}
