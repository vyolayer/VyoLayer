package middleware

import (
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProjectMiddleware struct {
	memberService service.ProjectMemberService
}

func NewProjectMiddleware(memberService service.ProjectMemberService) ProjectMiddleware {
	return ProjectMiddleware{memberService: memberService}
}

// CheckProjectMembership validates that the authenticated user is an active member
// of the project identified by the :projectId route param.
// On success it stores the project member in ctx.Locals("project_member").
func (pm *ProjectMiddleware) CheckProjectMembership() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		localUserID, ok := ctx.Locals("user_id").(types.UserID)
		if !ok || localUserID.IsNil() {
			return response.Error(ctx, errors.Unauthorized("Invalid user ID in token"))
		}

		projectID, err := types.ReconstructProjectID(ctx.Params("projectId"))
		if err != nil {
			return response.Error(ctx, errors.BadRequest("Invalid project ID"))
		}

		member, svcErr := pm.memberService.GetCurrentMember(ctx, localUserID, projectID)
		if svcErr != nil {
			return response.Error(ctx, svcErr)
		}

		ctx.Locals("project_member", member)
		return ctx.Next()
	}
}

// IsProjectAdmin checks that the user from ctx.Locals("project_member") has admin role.
func (pm *ProjectMiddleware) IsProjectAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		member, ok := ctx.Locals("project_member").(*service.ProjectMemberService)
		_ = member
		_ = ok

		// We rely on the domain struct directly
		localMember := ctx.Locals("project_member")
		if localMember == nil {
			return response.Error(ctx, errors.Unauthorized("Invalid project member"))
		}

		// The member is stored as *domain.ProjectMember
		type adminChecker interface {
			IsAdmin() bool
		}

		if m, ok := localMember.(adminChecker); ok {
			if !m.IsAdmin() {
				return response.Error(ctx, errors.Forbidden("Only project admins can perform this action"))
			}
		}

		return ctx.Next()
	}
}
