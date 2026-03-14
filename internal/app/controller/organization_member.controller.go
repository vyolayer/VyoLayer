package controller

import (
	"vyolayer/internal/app/dto"
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/service"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMemberController interface {
	GetAllMembersByOrgID(ctx *fiber.Ctx) error
	GetMemberByOrgIDAndMemberID(ctx *fiber.Ctx) error
	CurrentMember(ctx *fiber.Ctx) error
	RemoveMember(ctx *fiber.Ctx) error
	ChangeRole(ctx *fiber.Ctx) error
	LeaveOrganization(ctx *fiber.Ctx) error
	TransferOwnership(ctx *fiber.Ctx) error
}

type organizationMemberController struct {
	orgMemberService service.OrganizationMemberService
}

func NewOrganizationMemberController(orgMemberService service.OrganizationMemberService) OrganizationMemberController {
	return &organizationMemberController{orgMemberService: orgMemberService}
}

// GetAllMembersByOrgID godoc
// @Summary Get all members of an organization
// @Tags Organization Member
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationMemberDTO}
// @Failure 400,401,403,500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members [get]
func (ctrl *organizationMemberController) GetAllMembersByOrgID(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

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

// GetMemberByOrgIDAndMemberID godoc
// @Summary Get a member by ID
// @Tags Organization Member
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param memberId path string true "Member ID"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationMemberDTO}
// @Failure 400,401,403,500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/{memberId} [get]
func (ctrl *organizationMemberController) GetMemberByOrgIDAndMemberID(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	memberID, err := types.ReconstructOrganizationMemberID(ctx.Params("memberId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid member ID format"))
	}

	memberResp, err := ctrl.orgMemberService.GetOrgMemberByMemberID(ctx, orgID, memberID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, memberResp)
}

// CurrentMember godoc
// @Summary Get the current member
// @Tags Organization Member
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationMemberWithRBACDTO}
// @Failure 400,401,403,500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/me [get]
func (ctrl *organizationMemberController) CurrentMember(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	memberResp, err := ctrl.orgMemberService.GetCurrentMember(ctx, orgID, localUserID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, dto.FromDomainOrganizationMemberWithRBAC(&memberResp))
}

// RemoveMember godoc
// @Summary Remove a member from the organization
// @Tags Organization Member
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param memberId path string true "Member ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/{memberId} [delete]
func (ctrl *organizationMemberController) RemoveMember(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	memberID, err := types.ReconstructOrganizationMemberID(ctx.Params("memberId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid member ID format"))
	}

	if svcErr := ctrl.orgMemberService.RemoveMember(ctx, orgID, localUserID, memberID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Member removed successfully")
}

// ChangeRole godoc
// @Summary Change a member's role
// @Tags Organization Member
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param memberId path string true "Member ID"
// @Param request body dto.ChangeMemberRoleRequestDTO true "New role"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/{memberId}/role [patch]
func (ctrl *organizationMemberController) ChangeRole(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	memberID, err := types.ReconstructOrganizationMemberID(ctx.Params("memberId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid member ID format"))
	}

	var req dto.ChangeMemberRoleRequestDTO
	if parseErr := ctx.BodyParser(&req); parseErr != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	roleID, roleErr := types.ReconstructOrganizationRoleID(req.RoleID)
	if roleErr != nil {
		return response.Error(ctx, errors.BadRequest("Invalid role ID format"))
	}

	if svcErr := ctrl.orgMemberService.ChangeMemberRole(ctx, orgID, localUserID, memberID, roleID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Role changed successfully")
}

// LeaveOrganization godoc
// @Summary Leave an organization
// @Tags Organization Member
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/leave [post]
func (ctrl *organizationMemberController) LeaveOrganization(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	if svcErr := ctrl.orgMemberService.LeaveOrganization(ctx, orgID, localUserID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "You have left the organization")
}

// TransferOwnership godoc
// @Summary Transfer ownership to another member
// @Tags Organization Member
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param request body dto.TransferOwnershipRequestDTO true "New owner details"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId}/members/transfer-ownership [post]
func (ctrl *organizationMemberController) TransferOwnership(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	var req dto.TransferOwnershipRequestDTO
	if parseErr := ctx.BodyParser(&req); parseErr != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	newOwnerID, ownerErr := types.ReconstructOrganizationMemberID(req.NewOwnerMemberID)
	if ownerErr != nil {
		return response.Error(ctx, errors.BadRequest("Invalid new owner member ID format"))
	}

	if svcErr := ctrl.orgMemberService.TransferOwnership(ctx, orgID, localUserID, newOwnerID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Ownership transferred successfully")
}
