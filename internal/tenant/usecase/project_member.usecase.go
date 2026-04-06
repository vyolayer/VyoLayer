package usecase

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

// Define allowed roles locally for validation
var validProjectRoles = map[string]bool{
	domain.ProjectRoleAdmin.String():  true,
	domain.ProjectRoleMember.String(): true,
	domain.ProjectRoleViewer.String(): true,
}

type projectMemberUC struct {
	logger     *logger.AppLogger
	memberRepo domain.ProjectMemberRepository
	projRepo   domain.ProjectRepository // Injected to check MaxMembers limits!
}

// NewProjectMemberUseCase creates the use case with its required dependencies
func NewProjectMemberUseCase(logger *logger.AppLogger, memberRepo domain.ProjectMemberRepository, projRepo domain.ProjectRepository) domain.ProjectMemberUseCase {
	return &projectMemberUC{
		logger:     logger,
		memberRepo: memberRepo,
		projRepo:   projRepo,
	}
}

func (u *projectMemberUC) AddMember(ctx context.Context, orgID, projectID, userID, addedBy uuid.UUID, role string) (*domain.ProjectMember, error) {
	// Validate the role
	if role == "" {
		role = domain.ProjectRoleViewer.String()
	}
	if !validProjectRoles[role] {
		return nil, status.Errorf(codes.InvalidArgument, "invalid project role specified")
	}

	// Get Project
	project, err := u.projRepo.GetByID(ctx, orgID, projectID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "project not found")
	}

	// Check if the user is already an active member
	existing, err := u.memberRepo.GetByUserID(ctx, projectID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing membership: %s", err.Error())
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user is already an active member of this project")
	}

	// Limit Check: Ensure we haven't hit the MaxMembers limit
	// Note: We need the Organization ID to fetch the project. Since it's not in the method signature,
	// we assume the API Gateway / Handler validated ownership, and we can fetch it.
	// For absolute safety, you might want to add orgID to this method signature later!
	_, totalCount, err := u.memberRepo.List(ctx, projectID, 1, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing membership: %s", err.Error())
	}

	// Assuming MaxMembers is roughly 100 based on our previous constraints
	if totalCount >= int32(project.GetMaxMembers()) {
		return nil, status.Errorf(codes.ResourceExhausted, "project has reached its maximum member capacity")
	}

	// Create Domain Object & Save
	member := domain.NewProjectMember(projectID, userID, addedBy, domain.ProjectRole(role))
	if err := u.memberRepo.Add(ctx, nil, member); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add member to project: %s", err.Error())
	}

	return member, nil
}

func (u *projectMemberUC) GetMember(ctx context.Context, projectID, memberID uuid.UUID) (*domain.ProjectMember, error) {
	member, err := u.memberRepo.GetByID(ctx, projectID, memberID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, status.Errorf(codes.NotFound, "member not found")
	}
	return member, nil
}

func (u *projectMemberUC) GetCurrentMember(ctx context.Context, projectID, userID uuid.UUID) (*domain.ProjectMember, error) {
	member, err := u.memberRepo.GetByUserID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, status.Errorf(codes.NotFound, "you are not an active member of this project")
	}
	return member, nil
}

func (u *projectMemberUC) ListMembers(ctx context.Context, projectID uuid.UUID, limit, offset int32) ([]*domain.ProjectMember, int32, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return u.memberRepo.List(ctx, projectID, limit, offset)
}

func (u *projectMemberUC) ChangeRole(ctx context.Context, projectID, memberID uuid.UUID, newRole string) error {
	// Validate role
	if !validProjectRoles[newRole] {
		return status.Errorf(codes.InvalidArgument, "invalid project role specified")
	}

	// Fetch the target member to see their current role
	member, err := u.GetMember(ctx, projectID, memberID)
	if err != nil {
		return err
	}

	// If they are already this role, do nothing (Idempotent)
	if member.GetRole() == newRole {
		return nil
	}

	// THE LAST ADMIN GUARD: If demoting an admin, ensure they aren't the last one
	if member.GetRole() == "project_admin" {
		if err := u.ensureMultipleAdminsExist(ctx, projectID); err != nil {
			return err // Blocks the demotion
		}
	}

	// Execute Update
	return u.memberRepo.UpdateRole(ctx, projectID, memberID, newRole)
}

func (u *projectMemberUC) RemoveMember(ctx context.Context, projectID, memberID, removedBy uuid.UUID) error {
	// Fetch the target member
	member, err := u.GetMember(ctx, projectID, memberID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get member: %s", err.Error())
	}

	// THE LAST ADMIN GUARD: If kicking an admin, ensure they aren't the last one
	if member.GetRole() == "project_admin" {
		if err := u.ensureMultipleAdminsExist(ctx, projectID); err != nil {
			return status.Errorf(codes.Internal, "failed to ensure multiple admins exist: %s", err.Error())
		}
	}

	// Execute Removal
	return u.memberRepo.Remove(ctx, projectID, memberID, removedBy)
}

func (u *projectMemberUC) LeaveProject(ctx context.Context, projectID, userID uuid.UUID) error {
	// Find the member record for the calling user
	member, err := u.GetCurrentMember(ctx, projectID, userID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get current member: %s", err.Error())
	}

	// THE LAST ADMIN GUARD: If the leaving user is an admin, ensure there is another one
	if member.GetRole() == "project_admin" {
		if err := u.ensureMultipleAdminsExist(ctx, projectID); err != nil {
			return status.Errorf(codes.Internal, "failed to ensure multiple admins exist: %s", err.Error())
		}
	}

	// Execute Removal (The user removes themselves)
	return u.memberRepo.Remove(ctx, projectID, member.GetID(), userID)
}

// --- Private Helper Methods ---

// ensureMultipleAdminsExist fetches the project members and counts the admins.
// It returns an error if there is 1 or fewer admins left.
func (u *projectMemberUC) ensureMultipleAdminsExist(ctx context.Context, projectID uuid.UUID) error {
	// Fetch up to 100 members (our max limit) to count the admins
	members, _, err := u.memberRepo.List(ctx, projectID, 100, 0)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to verify admin count: %s", err.Error())
	}

	adminCount := 0
	for _, m := range members {
		if m.GetRole() == "project_admin" {
			adminCount++
		}
	}

	if adminCount <= 1 {
		return status.Error(codes.Internal, "failed to ensure multiple admins exist")
	}

	return nil
}
