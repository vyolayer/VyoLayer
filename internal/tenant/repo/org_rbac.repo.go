package tenantrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

// --- Role Repository ---

type organizationRoleRepository struct {
	gormRepo
}

func NewOrganizationRoleRepository(db *gorm.DB, logger *logger.AppLogger) OrganizationRoleRepository {
	return &organizationRoleRepository{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

func (r *organizationRoleRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationRole, error) {
	var model OrganizationRole
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toRoleDomain(&model), nil
}

func (r *organizationRoleRepository) GetByName(ctx context.Context, name string) (*domain.OrganizationRole, error) {
	var model OrganizationRole
	if err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toRoleDomain(&model), nil
}

func (r *organizationRoleRepository) List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationRole, error) {
	var models []OrganizationRole
	// If roles are global (System Roles), you might just return all of them.
	// If you eventually support custom roles per org, add a "WHERE organization_id = ?" here.
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	var result []*domain.OrganizationRole
	for _, m := range models {
		result = append(result, toRoleDomain(&m))
	}
	return result, nil
}

// --- Permission Repository ---

type organizationPermissionRepository struct {
	gormRepo
}

func NewOrganizationPermissionRepository(db *gorm.DB, logger *logger.AppLogger) OrganizationPermissionRepository {
	return &organizationPermissionRepository{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

func (r *organizationPermissionRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationPermission, error) {
	var model OrganizationPermission
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, ConvertDBError(err, "Failed to get permission")
	}
	return toPermissionDomain(&model), nil
}

func (r *organizationPermissionRepository) GetByCode(ctx context.Context, code string) (*domain.OrganizationPermission, error) {
	var model OrganizationPermission
	// Assuming you added a "code" column to the DB. If it's dynamically generated
	// from resource:action, use a WHERE CONCAT(resource, ':', action) = ? query instead.
	if err := r.db.WithContext(ctx).First(&model, "code = ?", code).Error; err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, ConvertDBError(err, "Failed to get permission")
	}
	return toPermissionDomain(&model), nil
}

func (r *organizationPermissionRepository) List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationPermission, error) {
	var models []OrganizationPermission
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, ConvertDBError(err, "Failed to list permissions")
	}

	var result []*domain.OrganizationPermission
	for _, m := range models {
		result = append(result, toPermissionDomain(&m))
	}
	return result, nil
}
