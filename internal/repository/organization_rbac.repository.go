package repository

import (
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"

	"gorm.io/gorm"
)

type OrganizationRBACRepository interface {
	GetAllPermissions(orgID types.OrganizationID) ([]models.OrganizationPermission, *errors.AppError)
	GetAllRoles(orgID types.OrganizationID) ([]models.OrganizationRole, *errors.AppError)
}

// Current all is common permission and role repository, it will be changed in future
// for now use simple cache
var (
	allPermissionsCache []models.OrganizationPermission
	allRolesCache       []models.OrganizationRole
)

type organizationRBACRepository struct {
	db *gorm.DB
}

func NewOrganizationRBACRepository(db *gorm.DB) OrganizationRBACRepository {
	return &organizationRBACRepository{db: db}
}

func (repo *organizationRBACRepository) GetAllPermissions(orgID types.OrganizationID) ([]models.OrganizationPermission, *errors.AppError) {
	if permissions := allPermissionsCache; len(permissions) > 0 {
		return permissions, nil
	}

	var permissions []models.OrganizationPermission
	result := repo.db.Find(&permissions)
	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "Failed to get all permissions")
	}

	allPermissionsCache = permissions
	return permissions, nil
}

func (repo *organizationRBACRepository) GetAllRoles(orgID types.OrganizationID) ([]models.OrganizationRole, *errors.AppError) {
	if roles := allRolesCache; len(roles) > 0 {
		return roles, nil
	}

	var roles []models.OrganizationRole
	result := repo.db.Find(&roles)
	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "Failed to get all roles")
	}

	allRolesCache = roles
	return roles, nil
}
