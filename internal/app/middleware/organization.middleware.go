package middleware

import (
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/service"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMiddleware struct {
	memberService service.OrganizationMemberService
}

func NewOrganizationMiddleware(memberService service.OrganizationMemberService) OrganizationMiddleware {
	return OrganizationMiddleware{memberService: memberService}
}

func (om *OrganizationMiddleware) CheckOrganizationMembership() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		localUserId, ok := ctx.Locals("user_id").(types.UserID)
		if !ok || localUserId.IsNil() {
			return response.Error(ctx, errors.Unauthorized("Invalid user ID in token"))
		}

		orgId, err := types.ReconstructOrganizationID(ctx.Params("orgId"))
		if err != nil {
			return response.Error(ctx, errors.BadRequest("Invalid organization ID"))
		}

		member, err := om.memberService.GetCurrentMember(ctx, orgId, localUserId)
		if err != nil {
			return response.Error(ctx, err)
		}

		ctx.Locals("member", member)

		return ctx.Next()
	}
}

func (om *OrganizationMiddleware) HasRole(role string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		localMember := ctx.Locals("member")
		if localMember == nil {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		member, ok := localMember.(domain.OrganizationMemberWithRBAC)
		if !ok {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		if !member.HasRole(role) {
			return response.Error(ctx, errors.Forbidden("You don't have permission to perform this action"))
		}

		return ctx.Next()
	}
}

// Is admin
func (om *OrganizationMiddleware) IsAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		localMember := ctx.Locals("member")
		if localMember == nil {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		member, ok := localMember.(domain.OrganizationMemberWithRBAC)
		if !ok {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		if !member.IsAdmin() {
			return response.Error(ctx, errors.Forbidden("You don't have permission to perform this action"))
		}

		return ctx.Next()
	}
}

// Is owner
func (om *OrganizationMiddleware) IsOwner() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		localMember := ctx.Locals("member")
		if localMember == nil {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		member, ok := localMember.(domain.OrganizationMemberWithRBAC)
		if !ok {
			return response.Error(ctx, errors.Unauthorized("Invalid member"))
		}

		if !member.IsOwner() {
			return response.Error(ctx, errors.Forbidden("You don't have permission to perform this action"))
		}

		return ctx.Next()
	}
}
