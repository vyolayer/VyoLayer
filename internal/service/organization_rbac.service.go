package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/internal/repository"
)

type OrganizationPermission struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Group       string `json:"group"`
	IsSystem    bool   `json:"is_system"`
}

func PermissionCode(resource string, action string) string {
	return fmt.Sprintf("%s.%s", resource, action)
}

type OrganizationRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
	IsDefault   bool   `json:"is_default"`
}

type OrganizationRBACService interface {
	GetAllPermissions(ctx *fiber.Ctx, orgID types.OrganizationID) ([]OrganizationPermission, error)
	GetAllRoles(ctx *fiber.Ctx, orgID types.OrganizationID) ([]OrganizationRole, error)
}

type organizationRBACService struct {
	rbacRepo repository.OrganizationRBACRepository
}

func NewOrganizationRBACService(rbacRepo repository.OrganizationRBACRepository) OrganizationRBACService {
	return &organizationRBACService{rbacRepo: rbacRepo}
}

func (service *organizationRBACService) GetAllPermissions(ctx *fiber.Ctx, orgID types.OrganizationID) ([]OrganizationPermission, error) {
	permissions, err := service.rbacRepo.GetAllPermissions(orgID)
	if err != nil {
		return nil, err
	}

	var orgPermissions []OrganizationPermission
	for _, permission := range permissions {
		orgPermissions = append(orgPermissions, OrganizationPermission{
			ID:          permission.PublicID().String(),
			Code:        PermissionCode(permission.Resource, permission.Action),
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
			Group:       permission.Group,
			IsSystem:    permission.IsSystem,
		})
	}

	return orgPermissions, nil
}

// GetAllRoles returns all roles for the given organization
func (service *organizationRBACService) GetAllRoles(ctx *fiber.Ctx, orgID types.OrganizationID) ([]OrganizationRole, error) {
	roles, err := service.rbacRepo.GetAllRoles(orgID)
	if err != nil {
		return nil, err
	}

	var orgRoles []OrganizationRole
	for _, role := range roles {
		orgRoles = append(orgRoles, OrganizationRole{
			ID:          role.PublicID().String(),
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			IsDefault:   role.IsDefault,
		})
	}

	return orgRoles, nil
}
