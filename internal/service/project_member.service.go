package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/internal/repository"
	"github.com/vyolayer/vyolayer/pkg/errors"
)

// ProjectMemberService defines the interface for project member management.
type ProjectMemberService interface {
	AddMember(ctx *fiber.Ctx, actorUserID types.UserID, projectID types.ProjectID, targetUserID types.UserID, role string) (*domain.ProjectMember, *errors.AppError)
	ListMembers(ctx *fiber.Ctx, userID types.UserID, projectID types.ProjectID) ([]domain.ProjectMember, *errors.AppError)
	GetCurrentMember(ctx *fiber.Ctx, userID types.UserID, projectID types.ProjectID) (*domain.ProjectMember, *errors.AppError)
	UpdateRole(ctx *fiber.Ctx, actorUserID types.UserID, projectID types.ProjectID, memberID types.ProjectMemberID, newRole string) *errors.AppError
	RemoveMember(ctx *fiber.Ctx, actorUserID types.UserID, projectID types.ProjectID, memberID types.ProjectMemberID) *errors.AppError
	Leave(ctx *fiber.Ctx, userID types.UserID, projectID types.ProjectID) *errors.AppError
}

type projectMemberService struct {
	memberRepo  repository.ProjectMemberRepository
	projectRepo repository.ProjectRepository
	auditRepo   repository.AuditLogRepository
}

func NewProjectMemberService(
	memberRepo repository.ProjectMemberRepository,
	projectRepo repository.ProjectRepository,
	auditRepo repository.AuditLogRepository,
) ProjectMemberService {
	return &projectMemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
		auditRepo:   auditRepo,
	}
}

func (s *projectMemberService) AddMember(
	ctx *fiber.Ctx,
	actorUserID types.UserID,
	projectID types.ProjectID,
	targetUserID types.UserID,
	role string,
) (*domain.ProjectMember, *errors.AppError) {
	// Validate role
	if !domain.IsValidProjectRole(role) {
		return nil, domain.ValidationError("Invalid project role: must be 'admin', 'member', or 'viewer'")
	}

	// Verify actor is an admin
	actor, err := s.memberRepo.FindByUserAndProject(ctx.Context(), actorUserID, projectID)
	if err != nil {
		return nil, err
	}
	if !actor.IsAdmin() {
		return nil, errors.Forbidden("Only project admins can add members")
	}

	// Check if target is already a member
	_, existsErr := s.memberRepo.FindByUserAndProject(ctx.Context(), targetUserID, projectID)
	if existsErr == nil {
		return nil, domain.ProjectMemberAlreadyExistsError(targetUserID.String())
	}

	// Check project member limit
	project, err := s.projectRepo.FindByID(ctx.Context(), projectID)
	if err != nil {
		return nil, err
	}
	if !project.CanAddMember() {
		return nil, domain.ProjectFullError()
	}

	// Create member
	member := domain.NewProjectMember(projectID, targetUserID, actorUserID, role)
	created, createErr := s.memberRepo.AddMember(ctx.Context(), member)
	if createErr != nil {
		return nil, createErr
	}

	return created, nil
}

func (s *projectMemberService) ListMembers(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) ([]domain.ProjectMember, *errors.AppError) {
	// Verify user is a member
	_, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return nil, errors.Forbidden("You are not a member of this project")
	}

	return s.memberRepo.FindByProjectID(ctx.Context(), projectID)
}

func (s *projectMemberService) GetCurrentMember(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) (*domain.ProjectMember, *errors.AppError) {
	return s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
}

func (s *projectMemberService) UpdateRole(
	ctx *fiber.Ctx,
	actorUserID types.UserID,
	projectID types.ProjectID,
	memberID types.ProjectMemberID,
	newRole string,
) *errors.AppError {
	if !domain.IsValidProjectRole(newRole) {
		return domain.ValidationError("Invalid project role")
	}

	// Verify actor is admin
	actor, err := s.memberRepo.FindByUserAndProject(ctx.Context(), actorUserID, projectID)
	if err != nil {
		return err
	}
	if !actor.IsAdmin() {
		return errors.Forbidden("Only project admins can change roles")
	}

	// Verify target member exists
	target, err := s.memberRepo.FindByID(ctx.Context(), memberID)
	if err != nil {
		return err
	}

	// Cannot change own role
	if target.UserID.String() == actorUserID.String() {
		return errors.BadRequest("Cannot change your own role")
	}

	return s.memberRepo.UpdateRole(ctx.Context(), memberID, newRole)
}

func (s *projectMemberService) RemoveMember(
	ctx *fiber.Ctx,
	actorUserID types.UserID,
	projectID types.ProjectID,
	memberID types.ProjectMemberID,
) *errors.AppError {
	// Verify actor is admin
	actor, err := s.memberRepo.FindByUserAndProject(ctx.Context(), actorUserID, projectID)
	if err != nil {
		return err
	}
	if !actor.IsAdmin() {
		return errors.Forbidden("Only project admins can remove members")
	}

	// Verify target member exists
	target, err := s.memberRepo.FindByID(ctx.Context(), memberID)
	if err != nil {
		return err
	}

	// Cannot remove self
	if target.UserID.String() == actorUserID.String() {
		return errors.BadRequest("Cannot remove yourself, use leave instead")
	}

	return s.memberRepo.Remove(ctx.Context(), memberID, actorUserID)
}

func (s *projectMemberService) Leave(
	ctx *fiber.Ctx,
	userID types.UserID,
	projectID types.ProjectID,
) *errors.AppError {
	member, err := s.memberRepo.FindByUserAndProject(ctx.Context(), userID, projectID)
	if err != nil {
		return err
	}

	// Check if this is the last admin — cannot leave if so
	if member.IsAdmin() {
		members, listErr := s.memberRepo.FindByProjectID(ctx.Context(), projectID)
		if listErr != nil {
			return listErr
		}
		adminCount := 0
		for _, m := range members {
			if m.IsAdmin() && m.IsActive {
				adminCount++
			}
		}
		if adminCount <= 1 {
			return errors.BadRequest("Cannot leave: you are the last admin. Transfer admin role first.")
		}
	}

	return s.memberRepo.Remove(ctx.Context(), member.ID, userID)
}
