package controller

import (
	"vyolayer/internal/app/dto"
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/service"
	"vyolayer/internal/utils/validation"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProjectController interface {
	CreateProject(ctx *fiber.Ctx) error
	ListProjects(ctx *fiber.Ctx) error
	GetProjectByID(ctx *fiber.Ctx) error
	UpdateProject(ctx *fiber.Ctx) error
	ArchiveProject(ctx *fiber.Ctx) error
	RestoreProject(ctx *fiber.Ctx) error
	DeleteProject(ctx *fiber.Ctx) error
}

type projectController struct {
	projectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) ProjectController {
	return &projectController{projectService: projectService}
}

// @Summary Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param project body dto.CreateProjectRequestDTO true "Project to create"
// @Success 201 {object} dto.ProjectDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects [post]
func (pc *projectController) CreateProject(ctx *fiber.Ctx) error {
	// Get user from context
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	// Get organization ID from params
	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	// Bind request body
	var req dto.CreateProjectRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	// Validate request body
	if errs := validation.ValidateStruct(req); len(errs) > 0 {
		return validation.ValidationErrorsToResponse(ctx, errs)
	}

	project, svcErr := pc.projectService.Create(ctx, localUserID, orgID, req.Name, req.Description)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Created(ctx, dto.FromDomainProjectWithMembers(project))
}

// @Summary List all projects of an organization
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Success 200 {array} dto.ProjectDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects [get]
func (pc *projectController) ListProjects(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	orgID, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid organization ID format"))
	}

	projects, svcErr := pc.projectService.ListByOrganizationID(ctx, localUserID, orgID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	dtos := make([]dto.ProjectDTO, 0, len(projects))
	for _, p := range projects {
		if d := dto.FromDomainProject(&p); d != nil {
			dtos = append(dtos, *d)
		}
	}

	return response.Success(ctx, dtos)
}

// @Summary Get a specific project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ProjectDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId} [get]
func (pc *projectController) GetProjectByID(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	project, svcErr := pc.projectService.GetByID(ctx, localUserID, projectID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, dto.FromDomainProjectWithMembers(project))
}

// @Summary Update a specific project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param project body dto.UpdateProjectRequestDTO true "Project to update"
// @Success 200 {object} dto.ProjectDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId} [put]
func (pc *projectController) UpdateProject(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	var req dto.UpdateProjectRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	project, svcErr := pc.projectService.Update(ctx, localUserID, projectID, req.Name, req.Description)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, dto.FromDomainProjectWithMembers(project))
}

// @Summary Archive a specific project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/archive [post]
func (pc *projectController) ArchiveProject(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	if svcErr := pc.projectService.Archive(ctx, localUserID, projectID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Project archived successfully")
}

// @Summary Restore a specific project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/restore [post]
func (pc *projectController) RestoreProject(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	if svcErr := pc.projectService.Restore(ctx, localUserID, projectID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Project restored successfully")
}

// @Summary Delete a specific project by ID
// @Tags Projects
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param project body dto.DeleteProjectRequestDTO true "Project to delete"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId} [delete]
func (pc *projectController) DeleteProject(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	var req dto.DeleteProjectRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	if svcErr := pc.projectService.Delete(ctx, localUserID, projectID, req.ConfirmName); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Project deleted successfully")
}
