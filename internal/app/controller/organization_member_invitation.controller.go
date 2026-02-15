package controller

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMemberInvitationController interface {
	CreateInvitation(ctx *fiber.Ctx) error
	ListInvitations(ctx *fiber.Ctx) error
	GetPendingInvitations(ctx *fiber.Ctx) error
	AcceptInvitation(ctx *fiber.Ctx) error
	CancelInvitation(ctx *fiber.Ctx) error
}

type organizationMemberInvitationController struct {
	invitationService service.OrganizationMemberInvitationService
}

func NewOrganizationMemberInvitationController(
	invitationService service.OrganizationMemberInvitationService,
) OrganizationMemberInvitationController {
	return &organizationMemberInvitationController{
		invitationService: invitationService,
	}
}

// CreateInvitation godoc
// @Summary Create an organization member invitation
// @Description Create a new invitation to join an organization
// @Tags organization_invitations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param request body dto.CreateInvitationRequestDTO true "Invitation details"
// @Success 201 {object} response.SuccessResponse{data=dto.OrganizationMemberInvitationDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/invitations [post]
func (ctrl *organizationMemberInvitationController) CreateInvitation(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || userID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization ID from path
	orgIDStr := ctx.Params("orgId")
	if orgIDStr == "" {
		return response.Error(ctx, errors.BadRequest("Organization ID is required"))
	}

	orgID, err := types.ReconstructOrganizationID(orgIDStr)
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	// Parse request body
	var req dto.CreateInvitationRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	// Create invitation
	invitation, svcErr := ctrl.invitationService.CreateInvitation(ctx, orgID, userID, req)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	ctx.Status(fiber.StatusCreated)
	return response.Success(ctx, invitation)
}

// ListInvitations godoc
// @Summary List organization invitations
// @Description Get all invitations for an organization
// @Tags organization_invitations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationMemberInvitationDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/invitations [get]
func (ctrl *organizationMemberInvitationController) ListInvitations(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || userID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization ID from path
	orgIDStr := ctx.Params("orgId")
	if orgIDStr == "" {
		return response.Error(ctx, errors.BadRequest("Organization ID is required"))
	}

	orgID, err := types.ReconstructOrganizationID(orgIDStr)
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	// Get invitations
	invitations, svcErr := ctrl.invitationService.ListInvitationsByOrgID(ctx, orgID, userID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, invitations)
}

// GetPendingInvitations godoc
// @Summary Get pending invitations for current user
// @Description Get all pending invitations for the authenticated user's email
// @Tags organization_invitations
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationMemberInvitationDTO}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/invitations/pending [get]
func (ctrl *organizationMemberInvitationController) GetPendingInvitations(ctx *fiber.Ctx) error {
	// Get user email from context (assuming it's set by auth middleware)
	userEmail, ok := ctx.Locals("user_email").(string)
	if !ok || userEmail == "" {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get pending invitations
	invitations, svcErr := ctrl.invitationService.GetPendingInvitations(ctx, userEmail)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, invitations)
}

// AcceptInvitation godoc
// @Summary Accept an invitation
// @Description Accept an invitation to join an organization
// @Tags organization_invitations
// @Accept json
// @Produce json
// @QueryParam org-invite-token query string true "Invitation token"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/invitations/accept?org-invite-token={invitationToken} [post]
func (ctrl *organizationMemberInvitationController) AcceptInvitation(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || userID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// get invitation token from query
	invitationToken := ctx.Query("org-invite-token")
	if invitationToken == "" {
		return response.Error(ctx, errors.BadRequest("Invitation token is required"))
	}

	// Accept invitation
	if svcErr := ctrl.invitationService.AcceptInvitation(ctx, userID, invitationToken); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Invitation accepted successfully")
}

// CancelInvitation godoc
// @Summary Cancel an invitation
// @Description Cancel/delete an invitation
// @Tags organization_invitations
// @Accept json
// @Produce json
// @Param invitationId path string true "Invitation ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/invitations/{invitationId} [delete]
func (ctrl *organizationMemberInvitationController) CancelInvitation(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || userID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get invitation ID from path
	invitationIDStr := ctx.Params("invitationId")
	if invitationIDStr == "" {
		return response.Error(ctx, errors.BadRequest("Invitation ID is required"))
	}

	invitationID, err := types.ReconstructOrganizationMemberInvitationID(invitationIDStr)
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid invitation ID format"))
	}

	// Cancel invitation
	if svcErr := ctrl.invitationService.CancelInvitation(ctx, invitationID, userID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Invitation canceled successfully")
}
