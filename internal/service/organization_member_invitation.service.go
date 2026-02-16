package service

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type OrganizationMemberInvitationService interface {
	CreateInvitation(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		invitedByUserID types.UserID,
		req dto.CreateInvitationRequestDTO,
	) (*dto.OrganizationMemberInvitationDTO, error)

	ListInvitationsByOrgID(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) ([]dto.OrganizationMemberInvitationDTO, error)

	GetPendingInvitations(
		ctx *fiber.Ctx,
		userEmail string,
	) ([]dto.OrganizationMemberInvitationDTO, error)

	AcceptInvitation(
		ctx *fiber.Ctx,
		userID types.UserID,
		invitationToken string,
	) error

	CancelInvitation(
		ctx *fiber.Ctx,
		invitationID types.OrganizationMemberInvitationID,
		canceledByUserID types.UserID,
	) error
}

type organizationMemberInvitationService struct {
	invitationRepo repository.OrganizationMemberInvitationRepository
	memberRepo     repository.OrganizationMemberRepository
	userRepo       repository.UserRepository
}

func NewOrganizationMemberInvitationService(
	invitationRepo repository.OrganizationMemberInvitationRepository,
	memberRepo repository.OrganizationMemberRepository,
	userRepo repository.UserRepository,
) OrganizationMemberInvitationService {
	return &organizationMemberInvitationService{
		invitationRepo: invitationRepo,
		memberRepo:     memberRepo,
		userRepo:       userRepo,
	}
}

// CreateInvitation creates a new organization member invitation
func (s *organizationMemberInvitationService) CreateInvitation(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	invitedByUserID types.UserID,
	req dto.CreateInvitationRequestDTO,
) (*dto.OrganizationMemberInvitationDTO, error) {
	// Get the inviter's member record to get their member ID
	members, err := s.memberRepo.GetByOrgID(ctx.Context(), invitedByUserID, orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetching inviter membership")
	}

	if len(members) == 0 {
		return nil, errors.Forbidden("You are not a member of this organization")
	}

	inviterMember := &members[0]
	if !inviterMember.IsActive {
		return nil, errors.Forbidden("Your membership is not active")
	}

	// Check if invitation already exists for this email and organization
	exists, err := s.invitationRepo.ExistsByEmailAndOrg(ctx.Context(), req.Email, orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "checking existing invitations")
	}

	if exists {
		return nil, errors.InvitationAlreadyExists(req.Email, orgID.String())
	}

	// Check if user is already a member
	// We'll check by fetching a user with this email and checking their membership
	existingUser, userErr := s.userRepo.FindByEmail(req.Email)
	if userErr == nil && existingUser != nil {
		// User exists, check if already a member
		existingMembers, memberErr := s.memberRepo.GetByOrgID(ctx.Context(), existingUser.ID, orgID)
		if memberErr == nil && len(existingMembers) > 0 {
			return nil, errors.Conflict("User is already a member of this organization")
		}
	}

	// Create the invitation
	invitation, domainErr := domain.NewOrganizationMemberInvitation(
		orgID,
		inviterMember.ID,
		req.Email,
		req.RoleIDs,
		7, // 7 days expiration
	)
	if domainErr != nil {
		return nil, domainErr
	}

	// Save to database
	if saveErr := s.invitationRepo.Create(ctx.Context(), invitation); saveErr != nil {
		return nil, WrapRepositoryError(saveErr, "creating invitation")
	}

	invitationDTO := dto.FromDomainOrganizationMemberInvitation(invitation)
	return &invitationDTO, nil
}

// ListInvitationsByOrgID lists all invitations for an organization
func (s *organizationMemberInvitationService) ListInvitationsByOrgID(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) ([]dto.OrganizationMemberInvitationDTO, error) {
	// Verify user is a member of the organization
	members, err := s.memberRepo.GetByOrgID(ctx.Context(), userID, orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetching user membership")
	}

	if len(members) == 0 {
		return nil, errors.Forbidden("You are not a member of this organization")
	}

	// Get all invitations
	invitations, err := s.invitationRepo.GetByOrgID(ctx.Context(), orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "listing invitations")
	}

	invitationsDTOs := make([]dto.OrganizationMemberInvitationDTO, len(invitations))
	for i, inv := range invitations {
		invitationsDTOs[i] = dto.FromDomainOrganizationMemberInvitation(&inv)
	}

	return invitationsDTOs, nil
}

// GetPendingInvitations gets pending invitations for a user's email
func (s *organizationMemberInvitationService) GetPendingInvitations(
	ctx *fiber.Ctx,
	userEmail string,
) ([]dto.OrganizationMemberInvitationDTO, error) {
	invitations, err := s.invitationRepo.GetPendingByEmail(ctx.Context(), userEmail)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetching pending invitations")
	}

	invitationsDTOs := make([]dto.OrganizationMemberInvitationDTO, len(invitations))
	for i, inv := range invitations {
		invitationsDTOs[i] = dto.FromDomainOrganizationMemberInvitation(&inv)
	}

	return invitationsDTOs, nil
}

// AcceptInvitation accepts an invitation and creates a new organization member
func (s *organizationMemberInvitationService) AcceptInvitation(
	ctx *fiber.Ctx,
	userID types.UserID,
	invitationToken string,
) error {
	// Get the invitation
	invitation, err := s.invitationRepo.GetByToken(ctx.Context(), invitationToken)
	if err != nil {
		return WrapRepositoryError(err, "fetching invitation")
	}

	// Get the user to check email
	user, userErr := s.userRepo.FindById(userID)
	if userErr != nil {
		return WrapRepositoryError(userErr, "fetching user")
	}

	// Verify the email matches
	if user.Email != invitation.Email {
		return errors.Forbidden("This invitation is not for your email address")
	}

	// Check if invitation is expired
	if invitation.IsExpired() {
		return errors.Forbidden("This invitation has expired")
	}

	// Check is pending
	if !invitation.IsPending() {
		return errors.Forbidden("This invitation already accepted or cancelled")
	}

	// Check if user is already a member
	existingMembers, memberErr := s.memberRepo.GetByOrgID(ctx.Context(), userID, invitation.OrganizationID)
	if memberErr == nil && len(existingMembers) > 0 {
		return errors.Conflict("You are already a member of this organization")
	}

	// Accept the invitation (domain logic)
	if acceptErr := invitation.Accept(); acceptErr != nil {
		return acceptErr
	}

	// Update the invitation in the database
	if updateErr := s.invitationRepo.Update(ctx.Context(), invitation); updateErr != nil {
		return WrapRepositoryError(updateErr, "updating invitation")
	}

	// Create the organization member
	if _, saveErr := s.memberRepo.Create(
		ctx.Context(),
		user.ID,
		invitation.OrganizationID,
		invitation.ToRoleIDsString(),
	); saveErr != nil {
		return WrapRepositoryError(saveErr, "creating organization member")
	}

	return nil
}

// CancelInvitation cancels/deletes an invitation
func (s *organizationMemberInvitationService) CancelInvitation(
	ctx *fiber.Ctx,
	invitationID types.OrganizationMemberInvitationID,
	canceledByUserID types.UserID,
) error {
	// Get the invitation first to verify organization
	invitation, err := s.invitationRepo.GetByID(ctx.Context(), invitationID)
	if err != nil {
		return WrapRepositoryError(err, "fetching invitation")
	}

	// Check if invitation is already accepted
	if invitation.IsAccepted {
		return errors.Forbidden("Invitation is already accepted")
	}

	// Check if invitation is expired
	if invitation.IsExpired() {
		return errors.Forbidden("Invitation is expired")
	}

	// find inviter details by email
	theInvitedUser, err := s.userRepo.FindByEmail(invitation.Email)
	if err == nil && theInvitedUser != nil && theInvitedUser.ID.Compare(canceledByUserID) {
		deleteErr := s.invitationRepo.Delete(ctx.Context(), invitationID, canceledByUserID)
		if deleteErr != nil {
			return WrapRepositoryError(deleteErr, "canceling invitation")
		}

		return nil
	}

	// Verify user is a member of the organization
	members, memberErr := s.memberRepo.GetByOrgID(ctx.Context(), canceledByUserID, invitation.OrganizationID)
	if memberErr != nil {
		return WrapRepositoryError(memberErr, "fetching user membership")
	}

	if len(members) == 0 {
		return errors.Forbidden("You are not a member of this organization")
	}

	canceledByMember := &members[0]

	// Delete the invitation
	if deleteErr := s.invitationRepo.Delete(ctx.Context(), invitationID, canceledByMember.UserID); deleteErr != nil {
		return WrapRepositoryError(deleteErr, "canceling invitation")
	}

	return nil
}
