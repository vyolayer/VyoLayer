package controller

import (
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/service"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationRBACController interface {
	GetAllPermissions(ctx *fiber.Ctx) error
	GetAllRoles(ctx *fiber.Ctx) error
}

type organizationRBACController struct {
	rbacService service.OrganizationRBACService
}

func NewOrganizationRBACController(rbacService service.OrganizationRBACService) OrganizationRBACController {
	return &organizationRBACController{rbacService: rbacService}
}

func (ctrl *organizationRBACController) GetAllPermissions(ctx *fiber.Ctx) error {
	// Local user id
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization id
	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	permissions, err := ctrl.rbacService.GetAllPermissions(ctx, orgID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, permissions)
}

func (ctrl *organizationRBACController) GetAllRoles(ctx *fiber.Ctx) error {
	// Local user id
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization id
	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	roles, err := ctrl.rbacService.GetAllRoles(ctx, orgID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, roles)
}
