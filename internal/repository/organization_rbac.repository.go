package repository

import (
	"sync"

	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

type OrganizationRBACRepository interface {
	GetAllPermissions(orgID types.OrganizationID) ([]models.OrganizationPermission, *errors.AppError)
	GetAllRoles(orgID types.OrganizationID) ([]models.OrganizationRole, *errors.AppError)
}

// Current all is common permission and role repository, it will be changed in future
// for now use simple cache
var (
	cacheMutex          sync.RWMutex
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
	cacheMutex.RLock()
	if permissions := allPermissionsCache; len(permissions) > 0 {
		cacheMutex.RUnlock()
		return permissions, nil
	}
	cacheMutex.RUnlock()

	var permissions []models.OrganizationPermission
	result := repo.db.Find(&permissions)
	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "Failed to get all permissions")
	}

	cacheMutex.Lock()
	allPermissionsCache = permissions
	cacheMutex.Unlock()
	return permissions, nil
}

func (repo *organizationRBACRepository) GetAllRoles(orgID types.OrganizationID) ([]models.OrganizationRole, *errors.AppError) {
	cacheMutex.RLock()
	if roles := allRolesCache; len(roles) > 0 {
		cacheMutex.RUnlock()
		return roles, nil
	}
	cacheMutex.RUnlock()

	var roles []models.OrganizationRole
	result := repo.db.Find(&roles)
	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "Failed to get all roles")
	}

	cacheMutex.Lock()
	allRolesCache = roles
	cacheMutex.Unlock()
	return roles, nil
}
