package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantrepo "github.com/vyolayer/vyolayer/internal/tenant/repo"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberUseCase interface {
	GetByUserID(ctx context.Context, orgID, userID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error)
	GetAllMembersByOrg(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationMemberWithRoles, error)
	GetById(ctx context.Context, orgID, memberID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error)
	RemoveMember(ctx context.Context, orgID, memberID, currentUserID uuid.UUID) error
}

type OrganizationMemberUseCaseImpl struct {
	logger     *logger.AppLogger
	memberRepo tenantrepo.OrganizationMemberRepository
}

func NewOrganizationMemberUseCase(
	logger *logger.AppLogger,
	memberRepo tenantrepo.OrganizationMemberRepository,
) OrganizationMemberUseCase {
	return &OrganizationMemberUseCaseImpl{
		logger:     logger,
		memberRepo: memberRepo,
	}
}

func (uc *OrganizationMemberUseCaseImpl) GetByUserID(ctx context.Context, orgID, userID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error) {
	member, err := uc.memberRepo.GetByUserIdAndOrgId(ctx, userID, orgID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Organization member not found")
	}
	return member, nil
}

func (uc *OrganizationMemberUseCaseImpl) GetAllMembersByOrg(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationMemberWithRoles, error) {
	members, err := uc.memberRepo.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (uc *OrganizationMemberUseCaseImpl) GetById(ctx context.Context, orgID, memberID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error) {
	member, err := uc.memberRepo.GetById(ctx, memberID)
	if err != nil {
		return nil, err
	}
	if member.OrganizationID != orgID {
		return nil, status.Errorf(codes.NotFound, "Organization member not found")
	}

	return member, nil
}

func (uc *OrganizationMemberUseCaseImpl) RemoveMember(ctx context.Context, orgID, memberID, currentUserID uuid.UUID) error {
	member, err := uc.memberRepo.GetById(ctx, memberID)
	if err != nil {
		return err
	}

	if member.OrganizationID != orgID {
		return status.Errorf(codes.NotFound, "Organization member not found")
	}

	if member.UserID == currentUserID {
		return status.Errorf(codes.PermissionDenied, "You cannot remove yourself")
	}

	currentMember, err := uc.memberRepo.GetByUserIdAndOrgId(ctx, currentUserID, orgID)
	if err != nil {
		return err
	}

	err = checkRoleAndPermissionPolicy(currentMember, member)
	if err != nil {
		return err
	}

	now := time.Now()
	member.SetRemovedBy(&currentMember.ID)
	member.SetRemovedAt(&now)
	member.SetIsActive(false)
	member.SetUpdatedAt(now)

	err = uc.memberRepo.RemoveMember(ctx, &member.OrganizationMember)
	if err != nil {
		return err
	}

	return nil
}

func checkRoleAndPermissionPolicy(currentMember, member *domain.OrganizationMemberWithRolesAndPermissions) error {
	currentHighestLevel := getHighestRoleLevel(currentMember.Roles)
	memberHighestLevel := getHighestRoleLevel(member.Roles)

	if currentHighestLevel >= memberHighestLevel {
		return status.Errorf(codes.PermissionDenied, "You do not have permission to remove this member")
	}

	return nil
}

func getHighestRoleLevel(roles []domain.OrganizationRole) uint32 {
	var highestLevel uint32 = 0
	for _, role := range roles {
		if role.HierarchyLevel > highestLevel {
			highestLevel = role.HierarchyLevel
		}
	}
	return highestLevel
}
