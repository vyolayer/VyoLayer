package service

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// OrganizationService defines the interface for organization-related business logic.
type OrganizationService interface {
	ListByUserID(ctx *fiber.Ctx, userID types.UserID) ([]domain.Organization, *errors.AppError)
	Create(ctx *fiber.Ctx, userID types.UserID, name, description string) (*domain.Organization, *errors.AppError)
	GetByID(ctx *fiber.Ctx, userID types.UserID, orgID types.OrganizationID) (*domain.Organization, *errors.AppError)
	GetBySlug(ctx *fiber.Ctx, userID types.UserID, slug string) (*domain.Organization, *errors.AppError)
}

// organizationService defines the dependencies for the organization service.
type organizationService struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
}

// NewOrganizationService creates a new instance of organizationService.
func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	userRepo repository.UserRepository,
) OrganizationService {
	return &organizationService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// ListByUserID lists all organizations that the user is a member of.
func (os *organizationService) ListByUserID(
	ctx *fiber.Ctx,
	userID types.UserID,
) ([]domain.Organization, *errors.AppError) {
	// Fetch organizations from repository
	orgs, err := os.orgRepo.ListByUserID(ctx.Context(), userID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organizations by user ID")
	}

	return orgs, nil
}

// Create creates a new organization with the given name and description.
//
// It also creates an organization member with the owner role.
func (os *organizationService) Create(
	ctx *fiber.Ctx,
	userID types.UserID,
	name, description string,
) (*domain.Organization, *errors.AppError) {
	// Fetch the user
	user, err := os.userRepo.FindById(userID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch user for organization creation")
	}

	if user == nil {
		return nil, domain.UserNotFoundError(userID.String())
	}

	// Create organization domain entity
	org := domain.NewOrganization(user, name, description, nil, nil)

	// Validate the organization
	if err := org.Validate(); err != nil {
		return nil, err
	}

	// Persist to database (this also creates owner member and assigns Owner role)
	_, err = os.orgRepo.Create(ctx.Context(), org)
	if err != nil {
		return nil, WrapRepositoryError(err, "create organization")
	}

	// 5. Fetch the created organization with members
	createdOrg, err := os.GetByID(ctx, userID, org.ID)
	if err != nil {
		return nil, err
	}

	return createdOrg, nil
}

// GetByID fetches an organization by its ID.
//
// It also verifies that the user is a member of the organization.
func (os *organizationService) GetByID(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
) (*domain.Organization, *errors.AppError) {
	// Fetch organization from repository
	org, err := os.orgRepo.FindByID(ctx.Context(), orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organization by ID")
	}

	if org == nil {
		return nil, domain.OrganizationNotFoundError(orgID.String())
	}

	// Verify user is a member (authorization check)
	if !org.IsMember(userID) {
		return nil, errors.Forbidden("You are not a member of this organization")
	}

	return org, nil
}

// GetBySlug fetches an organization by its slug.
//
// It also verifies that the user is a member of the organization.
func (os *organizationService) GetBySlug(
	ctx *fiber.Ctx,
	userID types.UserID,
	slug string,
) (*domain.Organization, *errors.AppError) {
	// Fetch organization from repository
	org, err := os.orgRepo.FindBySlug(ctx.Context(), slug)
	if err != nil {
		return nil, WrapRepositoryError(err, "fetch organization by slug")
	}

	if org == nil {
		return nil, domain.OrganizationNotFoundError(slug)
	}

	// Use GetByID for the rest (includes member check)
	return os.GetByID(ctx, userID, org.ID)
}
