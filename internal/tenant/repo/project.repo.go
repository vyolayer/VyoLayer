package tenantrepo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantmodelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type projectRepo struct {
	db     *gorm.DB
	logger *logger.AppLogger
}

// NewProjectRepository creates a new instance of the ProjectRepository
func NewProjectRepository(db *gorm.DB, logger *logger.AppLogger) domain.ProjectRepository {
	return &projectRepo{
		db:     db,
		logger: logger,
	}
}

// Create creates a new project in the database
func (r *projectRepo) Create(ctx context.Context, tx *gorm.DB, project *domain.Project) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := &Project{
		BaseModel: tenantmodelv1.BaseModel{ID: project.GetID(),
			TimeStamps: tenantmodelv1.TimeStamps{
				CreatedAt: project.GetCreatedAt(),
				UpdatedAt: project.GetUpdatedAt(),
			},
		},
		OrganizationID: project.GetOrganizationID(),
		Name:           project.GetName(),
		Slug:           project.GetSlug(),
		Description:    project.GetDescription(),
		IsActive:       project.GetIsActive(),
		CreatedBy:      project.GetCreatedBy(),
		MaxAPIKeys:     uint8(project.GetMaxAPIKeys()),
		MaxMembers:     uint8(project.GetMaxMembers()),
	}

	err := db.WithContext(ctx).Create(model).Error
	if err != nil {
		return ConvertDBError(err, "failed to create project")
	}

	return nil
}

func (r *projectRepo) GetByID(ctx context.Context, projectID uuid.UUID) (*domain.Project, error) {

	var model Project

	err := r.db.WithContext(ctx).
		Where("id = ?", projectID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific domain.ErrNotFound
		}
		r.logger.ErrorWithErr("failed to get project", err)
		return nil, err
	}

	return toProjectDomain(&model), nil
}

// GetByID fetches a project, strictly scoped to the Organization
func (r *projectRepo) GetByOrgID(ctx context.Context, orgID, projectID uuid.UUID) (*domain.Project, error) {
	var model Project

	err := r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", projectID, orgID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific domain.ErrNotFound
		}
		r.logger.ErrorWithErr("failed to get project", err)
		return nil, err
	}

	return toProjectDomain(&model), nil
}

// List fetches a paginated array of projects for an Organization
func (r *projectRepo) List(ctx context.Context, orgID uuid.UUID, limit, offset int32) ([]*domain.Project, int32, error) {
	var (
		models     []*Project
		totalCount int64
	)

	baseQuery := r.db.WithContext(ctx).
		Model(&Project{}).
		Where("organization_id = ?", orgID)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		r.logger.ErrorWithErr("failed to count projects", err)
		return nil, 0, err
	}

	err := baseQuery.
		Limit(int(limit)).
		Offset(int(offset)).
		Order("created_at DESC"). // Standard to show newest first
		Find(&models).Error

	if err != nil {
		r.logger.ErrorWithErr("failed to list projects", err)
		return nil, 0, err
	}

	result := make([]*domain.Project, 0, len(models))
	for _, m := range models {
		if mapped := toProjectDomain(m); mapped != nil {
			result = append(result, mapped)
		}
	}

	return result, int32(totalCount), nil
}

// Update modifies an existing project
func (r *projectRepo) Update(ctx context.Context, project *domain.Project) error {
	// We use Updates() with a map to only update specific fields, avoiding accidentally overwriting
	// sensitive fields like DatabaseURL if they aren't meant to change here.
	updates := map[string]interface{}{
		"name":        project.GetName(),
		"description": project.GetDescription(),
		"is_active":   project.GetIsActive(),
		"updated_at":  time.Now(),
	}

	result := r.db.WithContext(ctx).
		Model(&Project{}).
		Where("id = ? AND organization_id = ?", project.GetID(), project.GetOrganizationID()).
		Updates(updates)

	if result.Error != nil {
		r.logger.ErrorWithErr("failed to update project", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found or not owned by organization")
	}

	return nil
}

// Archive safely soft-deletes a project
func (r *projectRepo) Archive(ctx context.Context, orgID, projectID uuid.UUID) error {
	// If you use standard GORM Soft Deletes via `DeletedAt`, use .Delete()
	// If you use a custom `IsActive` or `ArchivedAt` flag, use .Update()

	result := r.db.WithContext(ctx).
		Where("id = ? AND organization_id = ?", projectID, orgID).
		Delete(&Project{}) // Triggers GORM Soft Delete

	if result.Error != nil {
		r.logger.Error("failed to archive project", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found or not owned by organization")
	}

	return nil
}

// Restore un-archives a project
func (r *projectRepo) Restore(ctx context.Context, orgID, projectID uuid.UUID) error {
	// Unscoped is required to find soft-deleted records in GORM
	result := r.db.WithContext(ctx).
		Unscoped().
		Model(&Project{}).
		Where("id = ? AND organization_id = ?", projectID, orgID).
		Update("deleted_at", nil) // Clear the soft-delete timestamp

	if result.Error != nil {
		r.logger.Error("failed to restore project", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found or not owned by organization")
	}

	return nil
}

// Delete performs a HARD delete (permanent destruction of data)
func (r *projectRepo) Delete(ctx context.Context, orgID, projectID uuid.UUID) error {
	// Unscoped() + Delete() bypasses Soft Deletes and permanently drops the row
	result := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ? AND organization_id = ?", projectID, orgID).
		Delete(&Project{})

	if result.Error != nil {
		r.logger.Error("failed to hard delete project", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("project not found or not owned by organization")
	}

	return nil
}

// --- Internal Mappers ---

// toProjectDomain converts a GORM model to our pure business logic Domain model safely
func toProjectDomain(m *Project) *domain.Project {
	if m == nil {
		return nil
	}

	domainProj := &domain.Project{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Name:           m.Name,
		Slug:           m.Slug,
		Description:    m.Description,
		IsActive:       m.IsActive,
		CreatedBy:      m.CreatedBy,
		MaxAPIKeys:     m.MaxAPIKeys,
		MaxMembers:     m.MaxMembers,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	// Safely map pointers (e.g., if you map the DeletedAt/ArchivedAt field)
	if m.DeletedAt.Valid {
		archivedTime := m.DeletedAt.Time
		domainProj.ArchivedAt = &archivedTime
	}

	return domainProj
}
