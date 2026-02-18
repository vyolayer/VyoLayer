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

	GetOrgMemberByMemberID(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		memberID types.OrganizationMemberID,
	) (dto.OrganizationMemberDTO, error)

	GetCurrentMember(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) (dto.OrganizationMemberWithRBACDTO, error)
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

func (service *organizationMemberService) GetCurrentMember(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) (dto.OrganizationMemberWithRBACDTO, error) {
	member, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if err != nil {
		return dto.OrganizationMemberWithRBACDTO{}, WrapRepositoryError(err, "getting organization member")
	}

	return dto.FromDomainOrganizationMemberWithRBAC(member), nil
}

func (service *organizationMemberService) GetOrgMemberByMemberID(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	memberID types.OrganizationMemberID,
) (dto.OrganizationMemberDTO, error) {
	member, err := service.orgMemberRepo.GetByOrgIDAndMemberID(ctx.Context(), orgID, memberID)
	if err != nil {
		return dto.OrganizationMemberDTO{}, WrapRepositoryError(err, "getting organization member")
	}

	return dto.FromDomainOrganizationMember(member), nil
}
