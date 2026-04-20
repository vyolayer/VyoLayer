package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/console/model"
)

type ResourceRepository interface {
	ListByServiceID(ctx context.Context, serviceID uint64) ([]model.ServiceResource, error)
	ListColumns(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceColumn, error)
	ListActions(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceAction, error)
	ListFilters(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceFilter, error)
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) ListByServiceID(ctx context.Context, serviceID uint64) ([]model.ServiceResource, error) {
	var resources []model.ServiceResource
	err := r.db.WithContext(ctx).
		Where("service_id = ?", serviceID).
		Order("sort_order asc").
		Find(&resources).Error
	return resources, err
}

func (r *resourceRepository) ListColumns(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceColumn, error) {
	var columns []model.ServiceResourceColumn
	if len(resourceIDs) == 0 {
		return columns, nil
	}
	err := r.db.WithContext(ctx).
		Where("resource_id IN ?", resourceIDs).
		Order("sort_order asc").
		Find(&columns).Error
	return columns, err
}

func (r *resourceRepository) ListActions(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceAction, error) {
	var actions []model.ServiceResourceAction
	if len(resourceIDs) == 0 {
		return actions, nil
	}
	err := r.db.WithContext(ctx).
		Where("resource_id IN ?", resourceIDs).
		Order("sort_order asc").
		Find(&actions).Error
	return actions, err
}

func (r *resourceRepository) ListFilters(ctx context.Context, resourceIDs []uint64) ([]model.ServiceResourceFilter, error) {
	var filters []model.ServiceResourceFilter
	if len(resourceIDs) == 0 {
		return filters, nil
	}
	err := r.db.WithContext(ctx).
		Where("resource_id IN ?", resourceIDs).
		Order("sort_order asc").
		Find(&filters).Error
	return filters, err
}
