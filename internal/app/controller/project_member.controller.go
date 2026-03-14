package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/app/dto"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/internal/service"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
)

type ProjectMemberController interface {
	ListMembers(ctx *fiber.Ctx) error
	GetCurrentMember(ctx *fiber.Ctx) error
	AddMember(ctx *fiber.Ctx) error
	ChangeRole(ctx *fiber.Ctx) error
	RemoveMember(ctx *fiber.Ctx) error
	LeaveProject(ctx *fiber.Ctx) error
}

type projectMemberController struct {
	memberService service.ProjectMemberService
}

func NewProjectMemberController(memberService service.ProjectMemberService) ProjectMemberController {
	return &projectMemberController{memberService: memberService}
}

// @Summary List all members of a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {array} dto.ProjectMemberDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members [get]
func (c *projectMemberController) ListMembers(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	members, svcErr := c.memberService.ListMembers(ctx, localUserID, projectID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	dtos := make([]dto.ProjectMemberDTO, 0, len(members))
	for _, m := range members {
		dtos = append(dtos, dto.FromDomainProjectMember(&m))
	}

	return response.Success(ctx, dtos)
}

// @Summary Get the current member of a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ProjectMemberDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members/current [get]
func (c *projectMemberController) GetCurrentMember(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	member, svcErr := c.memberService.GetCurrentMember(ctx, localUserID, projectID)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Success(ctx, dto.FromDomainProjectMember(member))
}

// @Summary Add a new member to a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param member body dto.AddProjectMemberRequestDTO true "Member to add"
// @Success 201 {object} dto.ProjectMemberDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members [post]
func (c *projectMemberController) AddMember(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	var req dto.AddProjectMemberRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	targetUserID, parseErr := types.ReconstructUserID(req.UserID)
	if parseErr != nil {
		return response.Error(ctx, errors.BadRequest("Invalid user ID format"))
	}

	member, svcErr := c.memberService.AddMember(ctx, localUserID, projectID, targetUserID, req.Role)
	if svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.Created(ctx, dto.FromDomainProjectMember(member))
}

// @Summary Change the role of a member in a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param memberId path string true "Member ID"
// @Param role body dto.ChangeProjectMemberRoleRequestDTO true "New role for the member"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members/{memberId}/role [post]
func (c *projectMemberController) ChangeRole(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	memberID, err := types.ReconstructProjectMemberID(ctx.Params("memberId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid member ID format"))
	}

	var req dto.ChangeProjectMemberRoleRequestDTO
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	if svcErr := c.memberService.UpdateRole(ctx, localUserID, projectID, memberID, req.Role); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Member role updated successfully")
}

// @Summary Remove a member from a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Param memberId path string true "Member ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members/{memberId} [delete]
func (c *projectMemberController) RemoveMember(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	memberID, err := types.ReconstructProjectMemberID(ctx.Params("memberId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid member ID format"))
	}

	if svcErr := c.memberService.RemoveMember(ctx, localUserID, projectID, memberID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Member removed successfully")
}

// @Summary Leave a project
// @Tags Project Members
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param projectId path string true "Project ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /organizations/{orgId}/projects/{projectId}/members/leave [post]
func (c *projectMemberController) LeaveProject(ctx *fiber.Ctx) error {
	localUserID, ok := ctx.Locals("user_id").(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
	if err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid project ID format"))
	}

	if svcErr := c.memberService.Leave(ctx, localUserID, projectID); svcErr != nil {
		return response.Error(ctx, svcErr)
	}

	return response.SuccessMessage(ctx, "Left project successfully")
}
