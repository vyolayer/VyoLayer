package controller

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationController interface {
	CreateOrganization(ctx *fiber.Ctx) error
	OnboardOrganization(ctx *fiber.Ctx) error
	GetOrganizationByID(ctx *fiber.Ctx) error
	GetOrganizationBySlug(ctx *fiber.Ctx) error
	ListOrganizations(ctx *fiber.Ctx) error
	UpdateOrganization(ctx *fiber.Ctx) error
	ArchiveOrganization(ctx *fiber.Ctx) error
	RestoreOrganization(ctx *fiber.Ctx) error
	DeleteOrganization(ctx *fiber.Ctx) error
}

type organizationController struct {
	orgService service.OrganizationService
}

func NewOrganizationController(orgService service.OrganizationService) OrganizationController {
	return &organizationController{
		orgService: orgService,
	}
}

// CreateOrganization godoc
// @Summary Create organization
// @Description Create a new organization. The authenticated user becomes the owner.
// @Tags Organizations
// @Accept json
// @Produce json
// @Param organization body dto.CreateOrganizationRequestDTO true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations [post]
func (oc *organizationController) CreateOrganization(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	var req dto.CreateOrganizationRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	org, err := oc.orgService.Create(ctx, localUserID, req.Name, req.Description)
	if err != nil {
		return response.Error(ctx, err)
	}

	responseDTO := dto.FromDomainOrganizationWithMembers(org)
	return response.Created(ctx, responseDTO)
}

// OnboardOrganization godoc
// @Summary Onboard organization
// @Description Guided onboarding flow for creating a new organization.
// @Tags Organizations
// @Accept json
// @Produce json
// @Param organization body dto.CreateOrganizationRequestDTO true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations/onboarding [post]
func (oc *organizationController) OnboardOrganization(ctx *fiber.Ctx) error {
	return oc.CreateOrganization(ctx)
}

// GetOrganizationByID godoc
// @Summary Get organization by ID
// @Description Retrieve organization details by ID. User must be a member.
// @Tags Organizations
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{id} [get]
func (oc *organizationController) GetOrganizationByID(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
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

	org, svcErr := oc.orgService.GetByID(ctx, localUserID, orgID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	responseDTO := dto.FromDomainOrganizationWithMembers(org)
	return response.Success(ctx, responseDTO)
}

// GetOrganizationBySlug godoc
// @Summary Get organization by slug
// @Description Retrieve organization details by slug. User must be a member.
// @Tags Organizations
// @Produce json
// @Param slug path string true "Organization slug"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/slug/{slug} [get]
func (oc *organizationController) GetOrganizationBySlug(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return response.Error(ctx, errors.BadRequest("Organization slug is required"))
	}

	org, err := oc.orgService.GetBySlug(ctx, localUserID, slug)
	if err != nil {
		return response.Error(ctx, err)
	}

	responseDTO := dto.FromDomainOrganizationWithMembers(org)
	return response.Success(ctx, responseDTO)
}

// ListOrganizations godoc
// @Summary List organizations
// @Description List organizations the user is a member of.
// @Tags Organizations
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationDTO}
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations [get]
func (oc *organizationController) ListOrganizations(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgs, err := oc.orgService.ListByUserID(ctx, localUserID)
	if err != nil {
		return response.Error(ctx, err)
	}

	responseDTOs := make([]dto.OrganizationDTO, 0, len(orgs))
	for _, org := range orgs {
		orgDTO := dto.FromDomainOrganization(&org)
		responseDTOs = append(responseDTOs, *orgDTO)
	}

	return response.Success(ctx, responseDTOs)
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization details (name, description, slug). Requires Admin+ role.
// @Tags Organizations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param organization body dto.UpdateOrganizationRequestDTO true "Updated fields"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400,401,403,404,409 {object} response.ErrorResponse
// @Router /organizations/{orgId} [patch]
func (oc *organizationController) UpdateOrganization(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	var req dto.UpdateOrganizationRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	org, svcErr := oc.orgService.Update(ctx, localUserID, orgID, req.Name, req.Description, req.Slug)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, dto.FromDomainOrganizationWithMembers(org))
}

// ArchiveOrganization godoc
// @Summary Archive organization
// @Description Deactivate (archive) an organization. Requires Admin+ role.
// @Tags Organizations
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId}/archive [post]
func (oc *organizationController) ArchiveOrganization(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	if svcErr := oc.orgService.Archive(ctx, localUserID, orgID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Organization archived successfully")
}

// RestoreOrganization godoc
// @Summary Restore organization
// @Description Reactivate a previously archived organization. Requires Admin+ role.
// @Tags Organizations
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId}/restore [post]
func (oc *organizationController) RestoreOrganization(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	if svcErr := oc.orgService.Restore(ctx, localUserID, orgID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Organization restored successfully")
}

// DeleteOrganization godoc
// @Summary Delete organization
// @Description Permanently delete an organization. Requires Owner role and name confirmation.
// @Tags Organizations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param request body dto.DeleteOrganizationRequestDTO true "Deletion confirmation"
// @Success 200 {object} response.SuccessResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse
// @Router /organizations/{orgId} [delete]
func (oc *organizationController) DeleteOrganization(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	var req dto.DeleteOrganizationRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	if svcErr := oc.orgService.Delete(ctx, localUserID, orgID, req.ConfirmName); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Organization deleted successfully")
}
