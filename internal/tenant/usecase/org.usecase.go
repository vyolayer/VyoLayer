package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantrepo "github.com/vyolayer/vyolayer/internal/tenant/repo"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationUseCase interface {
	Create(ctx context.Context, name, description string) (*domain.Organization, *domain.OrganizationMember, error)
	Update(ctx context.Context, orgID uuid.UUID, name, description string) (*domain.Organization, error)
	Archive(ctx context.Context, orgID uuid.UUID, confirmName string) error
	Restore(ctx context.Context, orgID uuid.UUID) error
	Delete(ctx context.Context, orgID uuid.UUID, confirmName string) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationWithMember, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	List(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Organization, int, error)
	TransferOwnership(ctx context.Context, orgID uuid.UUID, currentUserID uuid.UUID, newOwnerMemberID uuid.UUID) error

	// RABC
	GetAllPermissions(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationPermission, error)
	GetAllRoles(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationRole, error)
}

type OrganizationUseCaseImpl struct {
	logger         *logger.AppLogger
	orgRepo        tenantrepo.OrganizationRepository
	memberRepo     tenantrepo.OrganizationMemberRepository
	memberRoleRepo tenantrepo.MemberOrganizationRoleRepository
	roleRepo       tenantrepo.OrganizationRoleRepository
	permRepo       tenantrepo.OrganizationPermissionRepository
}

func NewOrganizationUseCase(
	logger *logger.AppLogger,
	orgRepo tenantrepo.OrganizationRepository,
	memberRepo tenantrepo.OrganizationMemberRepository,
	memberRoleRepo tenantrepo.MemberOrganizationRoleRepository,
	roleRepo tenantrepo.OrganizationRoleRepository,
	permRepo tenantrepo.OrganizationPermissionRepository,
) OrganizationUseCase {
	return &OrganizationUseCaseImpl{
		logger:         logger,
		orgRepo:        orgRepo,
		memberRepo:     memberRepo,
		memberRoleRepo: memberRoleRepo,
		roleRepo:       roleRepo,
		permRepo:       permRepo,
	}
}

func (uc *OrganizationUseCaseImpl) Create(
	ctx context.Context,
	name, description string,
) (*domain.Organization, *domain.OrganizationMember, error) {
	userId, err := ctxutil.ExtractIAMUserUUID(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get owner role id
	ownerRole, err := uc.roleRepo.GetByName(ctx, "owner")
	if err != nil {
		return nil, nil, err
	}

	org := domain.NewOrganization(userId, name, description)
	member := domain.NewOrganizationMember(org.ID, userId)
	memberRole := domain.NewMemberOrganizationRole(org.ID, member.ID, ownerRole.ID, member.ID)

	uc.logger.Debug("Organization creating", map[string]any{
		"org":        org,
		"member":     member,
		"memberRole": memberRole,
	})

	// Create organization and member in a transaction
	tx, err := uc.orgRepo.BeginTx(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer uc.orgRepo.RollbackTx(tx)

	err = uc.orgRepo.Create(ctx, tx, org)
	if err != nil {
		return nil, nil, err
	}
	uc.logger.Debug("Organization created", "")

	err = uc.memberRepo.AddMember(ctx, tx, member)
	if err != nil {
		return nil, nil, err
	}
	uc.logger.Debug("Organization member added", "")

	err = uc.memberRoleRepo.AddRole(ctx, tx, memberRole)
	if err != nil {
		return nil, nil, err
	}
	uc.logger.Debug("Organization member role added", "")

	err = uc.orgRepo.CommitTx(tx)
	if err != nil {
		return nil, nil, err
	}
	uc.logger.Debug("Organization committed", "")

	return org, member, nil
}

func (uc *OrganizationUseCaseImpl) Update(
	ctx context.Context,
	orgID uuid.UUID,
	name, description string,
) (*domain.Organization, error) {
	_, err := ctxutil.ExtractIAMUserUUID(ctx)
	if err != nil {
		return nil, err
	}

	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	flag := false
	if name != "" && name != org.Name {
		org.Name = name
		flag = true
	}

	if description != "" && description != org.Description {
		org.Description = description
		flag = true
	}

	if !flag {
		return nil, status.Error(codes.InvalidArgument, "No changes found")
	}

	org.SetUpdatedAt(time.Now())

	err = uc.orgRepo.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (uc *OrganizationUseCaseImpl) Archive(ctx context.Context, orgID uuid.UUID, confirmName string) error {
	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	if org.Name != confirmName {
		return status.Error(codes.InvalidArgument, "Organization name does not match")
	}

	if org.GetArchivedAt() != nil {
		return status.Error(codes.AlreadyExists, "Organization is already archived")
	}

	return uc.orgRepo.Archive(ctx, orgID, confirmName)
}

func (uc *OrganizationUseCaseImpl) Restore(ctx context.Context, orgID uuid.UUID) error {
	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	if org.GetArchivedAt() == nil {
		return status.Error(codes.AlreadyExists, "Organization is not archived")
	}

	return uc.orgRepo.Restore(ctx, orgID)
}

func (uc *OrganizationUseCaseImpl) Delete(ctx context.Context, orgID uuid.UUID, confirmName string) error {
	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	if org.Name != confirmName {
		return status.Error(codes.InvalidArgument, "Organization name does not match")
	}

	return uc.orgRepo.Delete(ctx, orgID, confirmName)
}

func (uc *OrganizationUseCaseImpl) GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationWithMember, error) {
	return uc.orgRepo.GetByIDWithMember(ctx, id)
}

func (uc *OrganizationUseCaseImpl) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	return uc.orgRepo.GetBySlug(ctx, slug)
}

func (uc *OrganizationUseCaseImpl) List(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Organization, int, error) {
	orgs, err := uc.orgRepo.List(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	nextOffset := offset + limit
	if len(orgs) < limit {
		nextOffset = 0
	}
	return orgs, nextOffset, nil
}

func (uc *OrganizationUseCaseImpl) TransferOwnership(ctx context.Context, orgID uuid.UUID, currentUserID uuid.UUID, newOwnerMemberID uuid.UUID) error {
	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	if org.GetOwnerID().String() != currentUserID.String() {
		return status.Error(codes.PermissionDenied, "You are not the owner of this organization")
	}

	member, err := uc.memberRepo.GetById(ctx, newOwnerMemberID)
	if err != nil {
		return err
	}

	if member.GetOrganizationID().String() != orgID.String() {
		return status.Error(codes.PermissionDenied, "Member is not a member of this organization")
	}

	// Get owner role id
	ownerRole, err := uc.roleRepo.GetByName(ctx, "owner")
	if err != nil {
		return err
	}

	// Get member role id
	memberRole, err := uc.roleRepo.GetByName(ctx, "member")
	if err != nil {
		return err
	}

	tx, err := uc.orgRepo.BeginTx(ctx)
	if err != nil {
		return status.Error(codes.Internal, "Failed to transfer ownership")
	}
	defer uc.orgRepo.RollbackTx(tx)

	// Update member role
	err = uc.memberRoleRepo.UpdateRole(ctx, tx, newOwnerMemberID, ownerRole.ID)
	if err != nil {
		return err
	}

	// Update member role
	err = uc.memberRoleRepo.UpdateRole(ctx, tx, member.ID, memberRole.ID)
	if err != nil {
		return err
	}

	// Update organization owner
	err = uc.orgRepo.UpdateOwner(ctx, tx, orgID, newOwnerMemberID)
	if err != nil {
		return err
	}

	err = uc.orgRepo.CommitTx(tx)
	if err != nil {
		return status.Error(codes.Internal, "Failed to transfer ownership")
	}

	return nil
}

func (uc *OrganizationUseCaseImpl) GetAllPermissions(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationPermission, error) {
	return uc.permRepo.List(ctx, orgID)
}

func (uc *OrganizationUseCaseImpl) GetAllRoles(ctx context.Context, orgID uuid.UUID) ([]*domain.OrganizationRole, error) {
	return uc.roleRepo.List(ctx, orgID)
}
