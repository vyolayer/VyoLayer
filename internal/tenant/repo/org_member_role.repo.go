package tenantrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

type memberOrganizationRoleRepository struct {
	gormRepo // Embed common transaction methods
}

func NewMemberOrganizationRoleRepository(db *gorm.DB, logger *logger.AppLogger) MemberOrganizationRoleRepository {
	return &memberOrganizationRoleRepository{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

func (r *memberOrganizationRoleRepository) AddRole(ctx context.Context, tx *gorm.DB, role *domain.MemberOrganizationRole) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := toMemberRoleModel(role)
	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		return ConvertDBError(err, "Failed to add role")
	}

	role.CreatedAt = model.CreatedAt
	role.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *memberOrganizationRoleRepository) RemoveRole(ctx context.Context, role *domain.MemberOrganizationRole) error {
	now := time.Now()
	role.RevokedAt = &now

	// Soft revoke: we keep the record for audit, but set revoked_at and revoked_by
	return r.db.WithContext(ctx).
		Model(&MemberOrganizationRole{}).
		Where("id = ?", role.ID).
		Updates(map[string]interface{}{
			"revoked_at": role.RevokedAt,
			"revoked_by": role.RevokedBy,
		}).Error
}

func (r *memberOrganizationRoleRepository) GetByMemberId(ctx context.Context, memberID uuid.UUID) ([]*domain.MemberOrganizationRole, error) {
	var models []MemberOrganizationRole

	// Only fetch active roles
	err := r.db.WithContext(ctx).
		Where("member_id = ? AND revoked_at IS NULL", memberID).
		Find(&models).Error

	if err != nil {
		return nil, ConvertDBError(err, "Failed to get roles")
	}

	var result []*domain.MemberOrganizationRole
	for _, m := range models {
		result = append(result, toMemberRoleDomain(&m))
	}
	return result, nil
}

func (r *memberOrganizationRoleRepository) UpdateRole(ctx context.Context, tx *gorm.DB, memberID uuid.UUID, roleID uuid.UUID) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Model(&MemberOrganizationRole{}).
		Where("member_id = ?", memberID).
		Update("role_id", roleID).Error

	if err != nil {
		return ConvertDBError(err, "Failed to update role")
	}

	return nil
}
