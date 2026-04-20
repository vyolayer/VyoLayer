package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/console/model"
)

type ProjectServiceRepository interface {
	GetActiveByProjectAndKey(ctx context.Context, projectID uuid.UUID, serviceKey string) (*model.ProjectService, error)
	ListActiveByProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectService, error)
}

type projectServiceRepository struct {
	db *gorm.DB
}

func NewProjectServiceRepository(db *gorm.DB) ProjectServiceRepository {
	return &projectServiceRepository{db: db}
}

func (r *projectServiceRepository) GetActiveByProjectAndKey(ctx context.Context, projectID uuid.UUID, serviceKey string) (*model.ProjectService, error) {

	var projectService model.ProjectService

	err := r.db.WithContext(ctx).
		Preload("Service").
		Joins("LEFT JOIN services ON project_services.service_id = services.id").
		Where("project_services.project_id = ?", projectID).
		Where("services.key = ?", serviceKey).
		Where("project_services.status = ?", "active").
		First(&projectService).Error

	if err != nil {
		return nil, err
	}

	return &projectService, nil
}

func (r *projectServiceRepository) ListActiveByProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectService, error) {

	var projectServices []model.ProjectService

	err := r.db.WithContext(ctx).
		Preload("Service").
		Where("project_id = ?", projectID).
		Where("status = ?", "active").
		Find(&projectServices).Error

	if err != nil {
		return nil, err
	}

	return projectServices, nil
}
