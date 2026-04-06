package tenantrepo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantmodelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

type projectMemberRepo struct {
	db     *gorm.DB
	logger *logger.AppLogger
}

// NewProjectMemberRepository creates a new instance of the ProjectMemberRepository
func NewProjectMemberRepository(db *gorm.DB, logger *logger.AppLogger) domain.ProjectMemberRepository {
	return &projectMemberRepo{
		db:     db,
		logger: logger,
	}
}

// Add inserts a new member into a project
func (r *projectMemberRepo) Add(ctx context.Context, tx *gorm.DB, member *domain.ProjectMember) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := &tenantmodelv1.ProjectMember{
		BaseModel: tenantmodelv1.BaseModel{
			ID: member.GetID(),
		},
		ProjectID: member.GetProjectID(),
		UserID:    member.GetUserID(),
		Role:      member.GetRole(),
		AddedBy:   member.AddedBy,
	}

	err := db.WithContext(ctx).Create(model).Error
	if err != nil {
		r.logger.Error("failed to add project member", map[string]any{"error": err, "projectID": member.GetProjectID()})
		return err
	}
	return nil
}

// GetByID fetches a specific member record, ensuring it belongs to the project
func (r *projectMemberRepo) GetByID(ctx context.Context, projectID, memberID uuid.UUID) (*domain.ProjectMember, error) {
	var model tenantmodelv1.ProjectMember

	err := r.db.WithContext(ctx).
		Preload("User"). // Assuming you have a User relationship to fetch Email/FullName
		Where("id = ? AND project_id = ?", memberID, projectID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil so the UseCase can handle "Not Found" cleanly
		}
		r.logger.Error("failed to get project member by id", err)
		return nil, err
	}

	return toProjectMemberDomain(&model), nil
}

// GetByUserID fetches an ACTIVE member record by their User ID
func (r *projectMemberRepo) GetByUserID(ctx context.Context, projectID, userID uuid.UUID) (*domain.ProjectMember, error) {
	var model tenantmodelv1.ProjectMember

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND project_id = ? AND removed_at IS NULL", userID, projectID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("failed to get project member by user id", err)
		return nil, err
	}

	return toProjectMemberDomain(&model), nil
}

// List fetches a paginated array of ACTIVE members for a specific project
func (r *projectMemberRepo) List(ctx context.Context, projectID uuid.UUID, limit, offset int32) ([]*domain.ProjectMember, int32, error) {
	var models []*tenantmodelv1.ProjectMember
	var totalCount int64

	// Base query explicitly filters out removed members
	baseQuery := r.db.WithContext(ctx).
		Model(&tenantmodelv1.ProjectMember{}).
		Where("project_id = ? AND removed_at IS NULL", projectID)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		r.logger.Error("failed to count project members", err)
		return nil, 0, err
	}

	err := baseQuery.
		Preload("User").
		Limit(int(limit)).
		Offset(int(offset)).
		Order("joined_at DESC").
		Find(&models).Error

	if err != nil {
		r.logger.Error("failed to list project members", err)
		return nil, 0, err
	}

	result := make([]*domain.ProjectMember, 0, len(models))
	for _, m := range models {
		if mapped := toProjectMemberDomain(m); mapped != nil {
			result = append(result, mapped)
		}
	}

	return result, int32(totalCount), nil
}

// UpdateRole changes a member's permission level within the project
func (r *projectMemberRepo) UpdateRole(ctx context.Context, projectID, memberID uuid.UUID, newRole string) error {
	result := r.db.WithContext(ctx).
		Model(&tenantmodelv1.ProjectMember{}).
		Where("id = ? AND project_id = ? AND removed_at IS NULL", memberID, projectID).
		Update("role", newRole)

	if result.Error != nil {
		r.logger.Error("failed to update project member role", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("active project member not found")
	}

	return nil
}

// Remove performs a business-logic soft delete by stamping RemovedAt and RemovedBy
func (r *projectMemberRepo) Remove(ctx context.Context, projectID, memberID, removedBy uuid.UUID) error {
	now := time.Now()

	updates := map[string]interface{}{
		"removed_at": &now,
		"removed_by": removedBy,
	}

	result := r.db.WithContext(ctx).
		Model(&tenantmodelv1.ProjectMember{}).
		Where("id = ? AND project_id = ? AND removed_at IS NULL", memberID, projectID).
		Updates(updates)

	if result.Error != nil {
		r.logger.Error("failed to remove project member", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("active project member not found")
	}

	return nil
}

// --- Internal Mappers ---

// toProjectMemberDomain safely converts the GORM model to the pure Domain model
func toProjectMemberDomain(m *tenantmodelv1.ProjectMember) *domain.ProjectMember {
	if m == nil {
		return nil
	}

	domainMem := &domain.ProjectMember{
		ID:        m.ID,
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		Role:      domain.ProjectRole(m.Role),
		IsActive:  m.IsActive(),
		AddedBy:   m.AddedBy,
		JoinedAt:  m.JoinedAt,
		RemovedAt: m.RemovedAt,
		RemovedBy: m.RemovedBy,
	}

	// Safe Preload Check: Assuming the User relationship is bound to the struct
	// This prevents panics and nil-string overwrites if `.Preload("User")` was missed
	// Adjust `m.User.ID` depending on how your IAM User struct is attached
	if m.User.ID != uuid.Nil {
		domainMem.Email = m.User.Email
		domainMem.FullName = m.User.FullName
	}

	return domainMem
}
