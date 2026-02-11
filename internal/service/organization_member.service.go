package service

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMemberService interface {
	ListByOrgAndUserId(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) ([]dto.OrganizationMemberDTO, error)
}

type organizationMemberService struct {
	orgMemberRepo repository.OrganizationMemberRepository
}

func NewOrganizationMemberService(orgMemberRepo repository.OrganizationMemberRepository) OrganizationMemberService {
	return &organizationMemberService{orgMemberRepo: orgMemberRepo}
}

func (service *organizationMemberService) ListByOrgAndUserId(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) ([]dto.OrganizationMemberDTO, error) {
	members, err := service.orgMemberRepo.GetByOrgID(ctx.Context(), userID, orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "listing organization members")
	}

	membersDto := make([]dto.OrganizationMemberDTO, 0, len(members))
	for _, member := range members {
		membersDto = append(membersDto, dto.FromDomainOrganizationMember(&member))
	}

	return membersDto, nil
}
