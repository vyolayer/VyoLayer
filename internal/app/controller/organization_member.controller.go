package controller

import (
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMemberController interface {
	GetAllMembersByOrgID(ctx *fiber.Ctx) error
}

type organizationMemberController struct {
	orgMemberService service.OrganizationMemberService
}

func NewOrganizationMemberController(orgMemberService service.OrganizationMemberService) OrganizationMemberController {
	return &organizationMemberController{orgMemberService: orgMemberService}
}

// GetAllMembersByOrgID godoc
// @Summary Get all members of an organization
// @Description Get all members of an organization
// @Tags organization_member
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationMemberDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members [get]
func (ctrl *organizationMemberController) GetAllMembersByOrgID(ctx *fiber.Ctx) error {
	// Get user id from context
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization id
	orgIDStr := ctx.Params("orgId")
	if orgIDStr == "" {
		return response.Error(ctx, errors.BadRequest("Organization ID is required"))
	}

	orgID, err := types.ReconstructOrganizationID(orgIDStr)
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	membersResp, err := ctrl.orgMemberService.ListByOrgAndUserId(ctx, orgID, localUserID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, membersResp)
}
