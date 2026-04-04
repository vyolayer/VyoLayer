package tenantrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type organizationRepo struct {
	gormRepo
}

func NewOrganizationRepo(
	db *gorm.DB,
	logger *logger.AppLogger,
) OrganizationRepository {
	return &organizationRepo{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

func (r *organizationRepo) Create(ctx context.Context, tx *gorm.DB, org *domain.Organization) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := toOrgModel(org)
	err := db.WithContext(ctx).Create(model).Error
	if err != nil {
		r.logger.Error("Failed to create organization", map[string]any{"error": err})
		return ConvertDBError(err, "Failed to create organization")
	}

	// Update the domain model with generated data (like CreatedAt)
	org.CreatedAt = model.CreatedAt
	org.UpdatedAt = model.UpdatedAt

	r.logger.Debug("Organization created", map[string]any{"org": org})
	return nil
}

func (r *organizationRepo) Update(ctx context.Context, org *domain.Organization) error {

	// Omit CreatedAt to prevent accidental overwrites during updates
	err := r.db.WithContext(ctx).
		Model(&Organization{}).
		Where("id = ?", org.ID).
		// Update only name and description
		Updates(map[string]interface{}{
			"name":        org.Name,
			"description": org.Description,
		}).Error

	if err != nil {
		r.logger.Error("Failed to update organization", map[string]any{"error": err})
		return ConvertDBError(err, "Failed to update organization")
	}

	r.logger.Debug("Organization updated", map[string]any{"org": org})
	return nil
}

func (r *organizationRepo) Delete(ctx context.Context, orgID uuid.UUID, confirmName string) error {
	// Delete performs a soft delete
	err := r.db.WithContext(ctx).
		Where("id = ?", orgID).
		Where("name = ?", confirmName).
		Delete(&Organization{}).Error

	if err != nil {
		r.logger.Error("Failed to delete organization", map[string]any{"error": err})
		return ConvertDBError(err, "Failed to delete organization")
	}

	r.logger.Debug("Organization deleted", map[string]any{"orgID": orgID})
	return nil
}

// Archive soft deletes an organization by setting the archived_at field
func (r *organizationRepo) Archive(ctx context.Context, orgID uuid.UUID, confirmName string) error {
	result := r.db.WithContext(ctx).
		Model(&Organization{}).
		Where("id = ?", orgID).
		Where("name = ?", confirmName).
		Update("archived_at", time.Now())

	if result.Error != nil {
		r.logger.Error("Failed to archive organization", map[string]any{"error": result.Error})
		return ConvertDBError(result.Error, "failed to archive organization")
	}

	if result.RowsAffected == 0 {
		r.logger.Error("Failed to archive organization", map[string]any{"orgID": orgID, "confirmName": confirmName})
		// This means either the ID doesn't exist, OR the confirmName didn't match
		return NotFoundError("organization not found or confirmation name mismatch", "")
	}

	r.logger.Debug("Organization archived", map[string]any{"orgID": orgID})
	return nil
}

// Restore restores a soft-deleted organization by setting the archived_at field to nil
func (r *organizationRepo) Restore(ctx context.Context, orgID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&Organization{}).
		Where("id = ?", orgID).
		Update("archived_at", nil).Error

	if err != nil {
		r.logger.Error("Failed to restore organization", map[string]any{"error": err})
		return ConvertDBError(err, "Failed to restore organization")
	}

	r.logger.Debug("Organization restored", map[string]any{"orgID": orgID})
	return nil
}

// --- Read Implementation ---

func (r *organizationRepo) GetByID(ctx context.Context, orgID uuid.UUID) (*domain.Organization, error) {
	var model Organization
	err := r.db.WithContext(ctx).
		Where("id = ?", orgID).
		First(&model).Error

	if err != nil {
		r.logger.Error("Failed to get organization", map[string]any{"error": err})
		return nil, ConvertDBError(err, "Failed to get organization")
	}

	r.logger.Debug("Organization found", map[string]any{"org": model})
	return toOrgDomain(&model), nil
}

func (r *organizationRepo) GetByIDWithMember(ctx context.Context, orgID uuid.UUID) (*domain.OrganizationWithMember, error) {
	var model Organization

	err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("Members.User").
		Where("id = ? AND deleted_at IS NULL AND archived_at IS NULL", orgID).
		First(&model).Error

	if err != nil {
		r.logger.Error("Failed to get organization with member", map[string]any{"error": err})
		return nil, ConvertDBError(err, "Failed to get organization with member")
	}

	r.logger.Debug("Organization with member found", map[string]any{"org": model})
	return toOrgWithMemberDomain(&model), nil
}

func (r *organizationRepo) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	var model Organization
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL AND archived_at IS NULL").
		Where("slug = ?", slug).
		First(&model).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		// return nil, err
		r.logger.Error("Failed to get organization by slug", map[string]any{"error": err})
		return nil, ConvertDBError(err, "Failed to get organization by slug")
	}

	r.logger.Debug("Organization found", map[string]any{"org": model})
	return toOrgDomain(&model), nil
}

func (r *organizationRepo) List(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Organization, error) {
	var models []Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN tenant.organization_members om ON om.organization_id = organizations.id").
		Where("om.user_id = ?", userID).
		Where("archived_at IS NULL").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		r.logger.Error("Failed to list organizations", map[string]any{"error": err})
		return nil, ConvertDBError(err, "Failed to list organizations")
	}

	var orgs []*domain.Organization
	for _, model := range models {
		orgs = append(orgs, toOrgDomain(&model))
	}

	r.logger.Debug("Organizations found", map[string]any{"orgs": orgs, "limit": limit, "offset": offset})
	return orgs, nil
}

func (r *organizationRepo) UpdateOwner(ctx context.Context, tx *gorm.DB, orgID uuid.UUID, ownerID uuid.UUID) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Model(&Organization{}).
		Where("id = ?", orgID).
		Update("owner_id", ownerID).Error

	if err != nil {
		r.logger.Error("Failed to update organization owner", map[string]any{"error": err})
		return ConvertDBError(err, "Failed to update organization owner")
	}

	r.logger.Debug("Organization owner updated", map[string]any{"orgID": orgID, "ownerID": ownerID})
	return nil
}
