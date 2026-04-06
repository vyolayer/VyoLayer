package tenantrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type optimizedPermissionChecker struct {
	gormRepo
}

func NewOptimizedPermissionChecker(db *gorm.DB, logger *logger.AppLogger) PermissionChecker {
	return &optimizedPermissionChecker{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

func (c *optimizedPermissionChecker) HasPermission(ctx context.Context, orgID, userID uuid.UUID, requiredPermissionCode string) (bool, error) {
	// 1. Safely split the permission code (e.g., "organization:update" -> "organization", "update")
	parts := strings.Split(requiredPermissionCode, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid permission format, expected resource:action, got %s", requiredPermissionCode)
	}
	resource := parts[0]
	action := parts[1]

	var count int64

	// 2. Execute the highly specific raw query
	err := c.db.WithContext(ctx).
		Table("organization_members AS om").
		Joins("JOIN member_organization_roles AS mor ON mor.member_id = om.id").
		Joins("JOIN organization_role_permissions AS orp ON orp.role_id = mor.role_id").
		Joins("JOIN organization_permissions AS op ON op.id = orp.permission_id").
		Where("om.organization_id = ?", orgID).
		Where("om.user_id = ?", userID).
		Where("om.removed_at IS NULL").  // Ensure member is not removed
		Where("mor.revoked_at IS NULL"). // Ensure role wasn't revoked
		Where("om.deleted_at IS NULL").  // Ensure member isn't soft-deleted
		Where("mor.deleted_at IS NULL"). // Ensure role assignment isn't soft-deleted
		Where("orp.deleted_at IS NULL"). // Ensure permission mapping isn't soft-deleted
		Where("op.deleted_at IS NULL").  // Ensure permission itself isn't soft-deleted
		Where("op.resource = ? AND op.action = ?", resource, action).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *optimizedPermissionChecker) IsMember(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
	var count int64

	err := c.db.WithContext(ctx).
		Table("organization_members").
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Where("removed_at IS NULL").
		Where("deleted_at IS NULL"). // Explicitly catch soft deletes
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
