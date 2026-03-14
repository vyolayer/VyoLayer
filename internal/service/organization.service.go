package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/internal/repository"
	"github.com/vyolayer/vyolayer/pkg/errors"
)

// OrganizationService defines the interface for organization-related business logic.
type OrganizationService interface {
	ListByUserID(ctx *fiber.Ctx, userID types.UserID) ([]domain.Organization, *errors.AppError)
	Create(ctx *fiber.Ctx, userID types.UserID, name, description string) (*domain.Organization, *errors.AppError)
	GetByID(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID) (*domain.Organization, *errors.AppError)
	GetBySlug(ctx *fiber.Ctx, userID types.UserID, slug string) (*domain.Organization, *errors.AppError)
	Update(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID, name, description, slug *string) (*domain.Organization, *errors.AppError)
	Archive(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID) *errors.AppError
	Restore(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID) *errors.AppError
	Delete(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID, confirmName string) *errors.AppError
}

// organizationService defines the dependencies for the organization service.
type organizationService struct {
	orgRepo      repository.OrganizationRepository
	userRepo     repository.UserRepository
	memberRepo   repository.OrganizationMemberRepository
	auditLogRepo repository.AuditLogRepository
}

// NewOrganizationService creates a new instance of organizationService.
func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	userRepo repository.UserRepository,
	memberRepo repository.OrganizationMemberRepository,
	auditLogRepo repository.AuditLogRepository,
) OrganizationService {
	return &organizationService{
		orgRepo:      orgRepo,
		userRepo:     userRepo,
		memberRepo:   memberRepo,
		auditLogRepo: auditLogRepo,
	}
}

// ListByUserID lists all organizations that the user is a member of.
func (os *organizationService) ListByUserID(
	ctx *fiber.Ctx,
	userID types.UserID,
) ([]domain.Organization, *errors.AppError) {
	orgs, err := os.orgRepo.ListByUserID(ctx.Context(), userID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organizations by user ID")
	}

	return orgs, nil
}

// Create creates a new organization with the given name and description.
func (os *organizationService) Create(
	ctx *fiber.Ctx,
	userID types.UserID,
	name, description string,
) (*domain.Organization, *errors.AppError) {
	user, err := os.userRepo.FindById(userID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch user for organization creation")
	}

	if user == nil {
		return nil, domain.UserNotFoundError(userID.String())
	}

	org := domain.NewOrganization(user, name, description, nil, nil)

	if err := org.Validate(); err != nil {
		return nil, err
	}

	_, err = os.orgRepo.Create(ctx.Context(), org)
	if err != nil {
		return nil, WrapRepositoryError(err, "create organization")
	}

	createdOrg, err := os.GetByID(ctx, userID, org.ID)
	if err != nil {
		return nil, err
	}

	return createdOrg, nil
}

// GetByID fetches an organization by its ID and verifies membership.
func (os *organizationService) GetByID(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
) (*domain.Organization, *errors.AppError) {
	org, err := os.orgRepo.FindByID(ctx.Context(), orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organization by ID")
	}

	if org == nil {
		return nil, domain.OrganizationNotFoundError(orgID.String())
	}

	if !org.IsMember(userID) {
		return nil, errors.Forbidden("You are not a member of this organization")
	}

	return org, nil
}

// GetBySlug fetches an organization by its slug and verifies membership.
func (os *organizationService) GetBySlug(
	ctx *fiber.Ctx,
	userID types.UserID,
	slug string,
) (*domain.Organization, *errors.AppError) {
	org, err := os.orgRepo.FindBySlug(ctx.Context(), slug)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organization by slug")
	}

	if org == nil {
		return nil, domain.OrganizationNotFoundError(slug)
	}

	return os.GetByID(ctx, userID, org.ID)
}

// Update updates an organization's name, description, and/or slug.
func (os *organizationService) Update(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
	name, description, slug *string,
) (*domain.Organization, *errors.AppError) {
	org, err := os.GetByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}

	member, memberErr := os.memberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if memberErr != nil {
		return nil, WrapRepositoryError(memberErr, "get current member for update")
	}
	if !member.IsAdmin() {
		return nil, errors.Forbidden("You don't have permission to edit this organization")
	}

	if name != nil && *name != "" {
		org.UpdateName(*name)
	}
	if description != nil {
		org.UpdateDescription(*description)
	}
	if slug != nil && *slug != "" {
		exists, slugErr := os.orgRepo.SlugExists(ctx.Context(), *slug, orgID)
		if slugErr != nil {
			return nil, WrapRepositoryError(slugErr, "check slug uniqueness")
		}
		if exists {
			return nil, domain.OrganizationSlugConflictError(*slug)
		}
		org.Slug = *slug
	}

	if valErr := org.Validate(); valErr != nil {
		return nil, valErr
	}

	if updateErr := os.orgRepo.Update(ctx.Context(), org); updateErr != nil {
		return nil, WrapRepositoryError(updateErr, "update organization")
	}

	// Audit log
	_ = os.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        member.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "org.updated",
		ResourceType:   "organization",
		ResourceID:     orgID.InternalID().ID(),
	})

	return os.GetByID(ctx, userID, orgID)
}

// Archive deactivates (archives) an organization.
func (os *organizationService) Archive(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
) *errors.AppError {
	org, err := os.GetByID(ctx, userID, orgID)
	if err != nil {
		return err
	}

	member, memberErr := os.memberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if memberErr != nil {
		return WrapRepositoryError(memberErr, "get current member for archive")
	}
	if !member.IsAdmin() {
		return errors.Forbidden("You don't have permission to archive this organization")
	}

	if archiveErr := org.Deactivate(userID); archiveErr != nil {
		return archiveErr
	}

	if updateErr := os.orgRepo.Update(ctx.Context(), org); updateErr != nil {
		return WrapRepositoryError(updateErr, "archive organization")
	}

	_ = os.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        member.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "org.archived",
		ResourceType:   "organization",
		ResourceID:     orgID.InternalID().ID(),
	})

	return nil
}

// Restore reactivates an archived organization.
func (os *organizationService) Restore(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
) *errors.AppError {
	org, err := os.GetByID(ctx, userID, orgID)
	if err != nil {
		return err
	}

	member, memberErr := os.memberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if memberErr != nil {
		return WrapRepositoryError(memberErr, "get current member for restore")
	}
	if !member.IsAdmin() {
		return errors.Forbidden("You don't have permission to restore this organization")
	}

	if restoreErr := org.Reactivate(); restoreErr != nil {
		return restoreErr
	}

	if updateErr := os.orgRepo.Update(ctx.Context(), org); updateErr != nil {
		return WrapRepositoryError(updateErr, "restore organization")
	}

	_ = os.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        member.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "org.restored",
		ResourceType:   "organization",
		ResourceID:     orgID.InternalID().ID(),
	})

	return nil
}

// Delete permanently removes an organization. Requires Owner + name confirmation.
func (os *organizationService) Delete(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
	confirmName string,
) *errors.AppError {
	org, err := os.GetByID(ctx, userID, orgID)
	if err != nil {
		return err
	}

	member, memberErr := os.memberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if memberErr != nil {
		return WrapRepositoryError(memberErr, "get current member for delete")
	}
	if !member.IsOwner() {
		return errors.Forbidden("Only the organization owner can delete the organization")
	}

	if confirmName != org.Name {
		return domain.OrganizationDeleteConfirmationError()
	}

	// Audit log before deletion
	_ = os.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        member.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "org.deleted",
		ResourceType:   "organization",
		ResourceID:     orgID.InternalID().ID(),
		Severity:       "critical",
		Metadata: map[string]interface{}{
			"name": org.Name,
			"slug": org.Slug,
		},
	})

	if deleteErr := os.orgRepo.Delete(ctx.Context(), orgID); deleteErr != nil {
		return WrapRepositoryError(deleteErr, "delete organization")
	}

	return nil
}
