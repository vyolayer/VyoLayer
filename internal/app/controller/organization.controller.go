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
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body dto.CreateOrganizationRequestDTO true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations [post]
func (oc *organizationController) CreateOrganization(ctx *fiber.Ctx) error {
	// 1. Extract user ID from context
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// 2. Parse request body
	var req dto.CreateOrganizationRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	// 3. Create organization via service
	org, err := oc.orgService.Create(
		ctx,
		localUserID,
		req.Name,
		req.Description,
	)
	if err != nil {
		return response.Error(ctx, err)
	}

	// 4. Convert to response DTO
	responseDTO := dto.FromDomainOrganizationWithMembers(org)

	// 5. Return 201 Created
	return response.Created(ctx, responseDTO)
}

// OnboardOrganization godoc
// @Summary Onboard organization
// @Description Guided onboarding flow for creating a new organization. Similar to create but may include additional steps.
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body dto.CreateOrganizationRequestDTO true "Organization details"
// @Success 201 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations/onboarding [post]
func (oc *organizationController) OnboardOrganization(ctx *fiber.Ctx) error {
	// For now, onboarding is the same as create
	// Future: Add onboarding-specific steps (tutorials, setup wizards, etc.)
	return oc.CreateOrganization(ctx)
}

// GetOrganizationByID godoc
// @Summary Get organization by ID
// @Description Retrieve organization details by ID. User must be a member.
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID" example:"org_550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /organizations/{id} [get]
func (oc *organizationController) GetOrganizationByID(ctx *fiber.Ctx) error {
	// 1. Extract user ID from context
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// 2. Extract organization ID from params
	orgIDStr := ctx.Params("orgId")
	if orgIDStr == "" {
		return response.Error(ctx, errors.BadRequest("Organization ID is required"))
	}

	orgID, err := types.ReconstructOrganizationID(orgIDStr)
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	// 3. Fetch organization via service
	org, svcErr := oc.orgService.GetByID(ctx, localUserID, orgID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	// 4. Convert to response DTO
	responseDTO := dto.FromDomainOrganizationWithMembers(org)

	// 5. Return 200 OK
	return response.Success(ctx, responseDTO)
}

// GetOrganizationBySlug godoc
// @Summary Get organization by slug
// @Description Retrieve organization details by slug. User must be a member.
// @Tags organizations
// @Accept json
// @Produce json
// @Param slug path string true "Organization slug" example:"acme-corp"
// @Success 200 {object} response.SuccessResponse{data=dto.OrganizationResponseDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /organizations/slug/{slug} [get]
func (oc *organizationController) GetOrganizationBySlug(ctx *fiber.Ctx) error {
	// 1. Extract user ID from context
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// 2. Extract slug from params
	slug := ctx.Params("slug")
	if slug == "" {
		return response.Error(ctx, errors.BadRequest("Organization slug is required"))
	}

	// 3. Fetch organization via service
	org, err := oc.orgService.GetBySlug(ctx, localUserID, slug)
	if err != nil {
		return response.Error(ctx, err)
	}

	// 4. Convert to response DTO
	responseDTO := dto.FromDomainOrganizationWithMembers(org)

	// 5. Return 200 OK
	return response.Success(ctx, responseDTO)
}

// ListOrganizations godoc
// @Summary List organizations
// @Description List organizations the user is a member of.
// @Tags organizations
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]dto.OrganizationDTO}
// @Failure 401 {object} response.ErrorResponse
// @Router /organizations [get]
func (oc *organizationController) ListOrganizations(ctx *fiber.Ctx) error {
	// 1. Extract user ID from context
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// 2. Fetch organizations via service
	orgs, err := oc.orgService.ListByUserID(ctx, localUserID)
	if err != nil {
		return response.Error(ctx, err)
	}

	// 3. Convert to response DTOs
	responseDTOs := make([]dto.OrganizationDTO, 0, len(orgs))
	for _, org := range orgs {
		orgDTO := dto.FromDomainOrganization(&org)
		responseDTOs = append(responseDTOs, *orgDTO)
	}

	// 4. Return 200 OK
	return response.Success(ctx, responseDTOs)
}
