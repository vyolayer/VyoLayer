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

type OrganizationMemberInvitationUseCase interface {
	Create(ctx context.Context, orgID uuid.UUID, email string, roleIDs []string, invitedBy uuid.UUID) error
	Accept(ctx context.Context, userID uuid.UUID, token string) error
	CancelByOrgMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	ListByUserEmail(ctx context.Context, email string) ([]domain.OrganizationMemberInvitation, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationMemberInvitation, error)
	ListPendingByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationMemberInvitationWithInviter, error)
	ListByOrgMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) ([]domain.OrganizationMemberInvitation, error)
}

type OrganizationMemberInvitationUseCaseImpl struct {
	logger            *logger.AppLogger
	orgRepo           tenantrepo.OrganizationRepository
	orgMemberRepo     tenantrepo.OrganizationMemberRepository
	orgRoleRepo       tenantrepo.OrganizationRoleRepository
	orgMemberRoleRepo tenantrepo.MemberOrganizationRoleRepository
	invitationRepo    tenantrepo.OrganizationMemberInvitationRepository
}

func NewOrganizationMemberInvitationUseCase(
	logger *logger.AppLogger,
	orgRepo tenantrepo.OrganizationRepository,
	orgMemberRepo tenantrepo.OrganizationMemberRepository,
	orgRoleRepo tenantrepo.OrganizationRoleRepository,
	orgMemberRoleRepo tenantrepo.MemberOrganizationRoleRepository,
	invitationRepo tenantrepo.OrganizationMemberInvitationRepository,
) OrganizationMemberInvitationUseCase {
	return &OrganizationMemberInvitationUseCaseImpl{
		logger:            logger,
		orgRepo:           orgRepo,
		orgMemberRepo:     orgMemberRepo,
		invitationRepo:    invitationRepo,
		orgRoleRepo:       orgRoleRepo,
		orgMemberRoleRepo: orgMemberRoleRepo,
	}
}

// Create sends an invitation email to the target email address.
// It validates that the org exists, the inviting member is active, and the email
// is not already a member of the organization.
func (uc *OrganizationMemberInvitationUseCaseImpl) Create(ctx context.Context, orgID uuid.UUID, email string, roleIDs []string, invitedBy uuid.UUID) error {
	org, err := uc.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}
	if org == nil {
		return status.Errorf(codes.NotFound, "organization not found")
	}

	// Ensure the requester is still an active member
	currentMember, err := uc.orgMemberRepo.GetByUserIdAndOrgId(ctx, invitedBy, orgID)
	if err != nil {
		return err
	}

	// Reject if the email is already a member
	existingMember, _ := uc.orgMemberRepo.GetByOrgIdAndEmail(ctx, orgID, email)
	if existingMember != nil {
		return status.Errorf(codes.AlreadyExists, "member with that email already exists in the organization")
	}

	roleUUIDs := make([]uuid.UUID, len(roleIDs))
	for i, roleID := range roleIDs {
		parsed, err := uuid.Parse(roleID)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid role id: %s", roleID)
		}
		roleUUIDs[i] = parsed
	}

	invitationToken := uuid.New().String()
	invitation := domain.NewOrganizationMemberInvitation(
		orgID,
		currentMember.ID,
		email,
		invitationToken,
		roleUUIDs,
		7*24*time.Hour, // 7 days
	)

	// TODO: publish InvitationCreatedEvent or call email service here

	return uc.invitationRepo.Create(ctx, invitation)
}

// Accept validates the invitation token and marks it as accepted.
// The caller is expected to subsequently create the OrganizationMember record
// (handled by a dedicated AcceptInvitation handler/use case in the future).
func (uc *OrganizationMemberInvitationUseCaseImpl) Accept(ctx context.Context, userID uuid.UUID, token string) error {
	invitation, err := uc.invitationRepo.GetByToken(ctx, token)
	if err != nil {
		return err
	}
	if invitation == nil {
		return status.Errorf(codes.NotFound, "invitation not found")
	}

	if invitation.GetIsAccepted() {
		return status.Errorf(codes.FailedPrecondition, "invitation has already been accepted")
	}

	now := time.Now()
	if invitation.GetExpiredAt().Before(now) {
		return status.Errorf(codes.FailedPrecondition, "invitation has expired")
	}

	invitation.SetIsAccepted(true)
	invitation.SetAcceptedAt(&now)

	if err := uc.invitationRepo.Accept(ctx, invitation); err != nil {
		return err
	}

	role, err := uc.orgRoleRepo.GetByName(ctx, "viewer")
	if err != nil {
		return err
	}

	member := domain.NewOrganizationMember(
		invitation.OrganizationID,
		userID,
	)
	member.SetInvitedBy(member.GetInvitedBy())

	memberRole := domain.NewMemberOrganizationRole(
		invitation.OrganizationID,
		member.ID,
		role.ID,
		invitation.GetInvitedBy(),
	)

	tx, err := uc.orgMemberRepo.BeginTx(ctx)
	if err != nil {
		return status.Error(codes.Internal, "something went wrong")
	}

	if err := uc.orgMemberRepo.AddMember(ctx, tx, member); err != nil {
		tx.Rollback()
		return err
	}

	if err := uc.orgMemberRoleRepo.AddRole(ctx, tx, memberRole); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return status.Error(codes.Internal, "failed to create organization member")
	}

	return nil
}

// CancelByOrgMember soft-deletes an invitation that was created by the given member.
// Only the member who created the invitation (or an admin – enforced at the handler
// level via PBAC interceptor) may cancel it.
func (uc *OrganizationMemberInvitationUseCaseImpl) CancelByOrgMember(ctx context.Context, orgID, userID uuid.UUID) error {
	invitations, err := uc.invitationRepo.ListByInvitedBy(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if len(invitations) == 0 {
		return status.Errorf(codes.NotFound, "no invitations found for that member")
	}

	for _, inv := range invitations {
		if inv.IsAccepted {
			continue // Already accepted – skip silently
		}
		if err := uc.invitationRepo.Delete(ctx, inv); err != nil {
			return err
		}
	}
	return nil
}

// ListByUserId returns all invitations (across all orgs) for the calling user.
func (uc *OrganizationMemberInvitationUseCaseImpl) ListByUserEmail(ctx context.Context, email string) ([]domain.OrganizationMemberInvitation, error) {
	invitations, err := uc.invitationRepo.ListByUserEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return derefInvitationSlice(invitations), nil
}

// ListByOrg returns all invitations (all states) for an organization.
func (uc *OrganizationMemberInvitationUseCaseImpl) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationMemberInvitation, error) {
	invitations, err := uc.invitationRepo.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return derefInvitationSlice(invitations), nil
}

// ListPendingByOrg returns only active (not accepted, not expired, not deleted) invitations,
// each wrapped with the inviter's member ID. Full name enrichment can be added once a
// member-lookup call is available in the use case.
func (uc *OrganizationMemberInvitationUseCaseImpl) ListPendingByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationMemberInvitationWithInviter, error) {
	invitations, err := uc.invitationRepo.ListPendingByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.OrganizationMemberInvitationWithInviter, 0, len(invitations))
	for _, inv := range invitations {
		if inv == nil {
			continue
		}
		result = append(result, *inv)
	}
	return result, nil
}

// ListByOrgMember returns all invitations created by a particular member within an org.
func (uc *OrganizationMemberInvitationUseCaseImpl) ListByOrgMember(ctx context.Context, orgID uuid.UUID, memberID uuid.UUID) ([]domain.OrganizationMemberInvitation, error) {
	invitations, err := uc.invitationRepo.ListByInvitedBy(ctx, orgID, memberID)
	if err != nil {
		return nil, err
	}
	return derefInvitationSlice(invitations), nil
}

// --- helpers ---

func derefInvitationSlice(ptrs []*domain.OrganizationMemberInvitation) []domain.OrganizationMemberInvitation {
	result := make([]domain.OrganizationMemberInvitation, 0, len(ptrs))
	for _, p := range ptrs {
		if p != nil {
			result = append(result, *p)
		}
	}
	return result
}
