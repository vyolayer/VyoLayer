package service

import (
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/types"
	"vyolayer/internal/repository"
	"vyolayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// ProjectService defines the interface for project-related business logic.
type ProjectService interface {
	Create(
		ctx *fiber.Ctx,
		userID types.UserID,
		orgID types.OrganizationID,
		name, description string,
	) (*domain.Project, *errors.AppError)

	GetByID(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
	) (*domain.Project, *errors.AppError)

	ListByOrganizationID(
		ctx *fiber.Ctx,
		userID types.UserID,
		orgID types.OrganizationID,
	) ([]domain.Project, *errors.AppError)

	Update(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
		name, description *string,
	) (*domain.Project, *errors.AppError)

	Archive(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
	) *errors.AppError

	Restore(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
	) *errors.AppError

	Delete(
		ctx *fiber.Ctx,
		userID types.UserID,
		projectID types.ProjectID,
		confirmName string,
	) *errors.AppError
}

type projectService struct {
	projectRepo repository.ProjectRepository
	memberRepo  repository.ProjectMemberRepository
	orgRepo     repository.OrganizationRepository
	auditRepo   repository.AuditLogRepository
}

func NewProjectService(
	projectRepo repository.ProjectRepository,
	memberRepo repository.ProjectMemberRepository,
	orgRepo repository.OrganizationRepository,
	auditRepo repository.AuditLogRepository,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
		orgRepo:     orgRepo,
		auditRepo:   auditRepo,
	}
}

func (s *projectService) Create(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
	name, description string,
) (*domain.Project, *errors.AppError) {
	// Verify org exists and user is a member
	org, err := s.orgRepo.FindByID(ctx.Context(), orgID)
	if err != nil {
		return nil, err
	}

	if !org.IsActive {
		return nil, domain.OrganizationNotActiveError()
	}

	// Check org project limit
	count, err := s.projectRepo.CountByOrganizationID(ctx.Context(), orgID)
	if err != nil {
		return nil, err
	}
	if int(count) >= org.MaxProjects {
		return nil, domain.ProjectLimitReachedError()
	}

	// Create domain project (creator becomes admin)
	creator := &domain.User{ID: userID}
	project := domain.NewProject(orgID, creator, name, description, nil, nil)

	if validErr := project.Validate(); validErr != nil {
		return nil, validErr
	}

	created, err := s.projectRepo.Create(ctx.Context(), project)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *projectService) GetByID(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) (*domain.Project, *errors.AppError) {
	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return nil, err
	}

	// Verify user is a member of this project
	_, memberErr := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if memberErr != nil {
		return nil, errors.Forbidden("You are not a member of this project")
	}

	return project, nil
}

func (s *projectService) ListByOrganizationID(
	ctx *fiber.Ctx,
	userID types.UserID,
	orgID types.OrganizationID,
) ([]domain.Project, *errors.AppError) {
	projects, err := s.projectRepo.FindByOrganizationID(ctx.Context(), orgID)
	if err != nil {
		return nil, err
	}

	// Filter to only projects the user is a member of
	var accessible []domain.Project
	for _, p := range projects {
		if p.IsMember(userID) {
			accessible = append(accessible, p)
		}
	}

	return accessible, nil
}

func (s *projectService) Update(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
	name, description *string,
) (*domain.Project, *errors.AppError) {
	// Verify membership and admin role
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return nil, err
	}
	if !member.IsAdmin() {
		return nil, errors.Forbidden("Only project admins can update project settings")
	}

	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		// Check slug uniqueness
		project.UpdateName(*name)
		exists, slugErr := s.projectRepo.SlugExists(ctx.Context(), project.OrganizationID, project.Slug, projectID)
		if slugErr != nil {
			return nil, slugErr
		}
		if exists {
			return nil, domain.ProjectSlugConflictError(project.Slug)
		}
	}

	if description != nil {
		project.UpdateDescription(*description)
	}

	if updateErr := s.projectRepo.Update(ctx.Context(), project); updateErr != nil {
		return nil, updateErr
	}

	return s.projectRepo.FindByID(ctx.Context(), projectID)
}

func (s *projectService) Archive(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) *errors.AppError {
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return err
	}
	if !member.IsAdmin() {
		return errors.Forbidden("Only project admins can archive a project")
	}

	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return err
	}

	if deactivateErr := project.Deactivate(); deactivateErr != nil {
		return deactivateErr
	}

	return s.projectRepo.Update(ctx.Context(), project)
}

func (s *projectService) Restore(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) *errors.AppError {
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return err
	}
	if !member.IsAdmin() {
		return errors.Forbidden("Only project admins can restore a project")
	}

	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return err
	}

	if restoreErr := project.Reactivate(); restoreErr != nil {
		return restoreErr
	}

	return s.projectRepo.Update(ctx.Context(), project)
}

func (s *projectService) Delete(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
	confirmName string,
) *errors.AppError {
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return err
	}
	if !member.IsAdmin() {
		return errors.Forbidden("Only project admins can delete a project")
	}

	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return err
	}

	if project.Name != confirmName {
		return domain.ProjectDeleteConfirmationError()
	}

	return s.projectRepo.Delete(ctx.Context(), projectID)
}
