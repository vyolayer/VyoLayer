package domain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectRepository handles all database operations for Projects
type ProjectRepository interface {
	Create(ctx context.Context, tx *gorm.DB, project *Project) error

	GetByID(ctx context.Context, projectID uuid.UUID) (*Project, error)

	// Scoped by OrgID to ensure data isolation
	GetByOrgID(ctx context.Context, orgID, projectID uuid.UUID) (*Project, error)
	List(ctx context.Context, orgID uuid.UUID, limit, offset int32) ([]*Project, int32, error)

	Update(ctx context.Context, project *Project) error

	// Archive(ctx context.Context, orgID, projectID uuid.UUID) error
	// Restore(ctx context.Context, orgID, projectID uuid.UUID) error

	// Hard delete from the database
	Delete(ctx context.Context, orgID, projectID uuid.UUID) error
}

// ProjectMemberRepository handles database operations for Project Members
type ProjectMemberRepository interface {
	Add(ctx context.Context, tx *gorm.DB, member *ProjectMember) error

	// Fetching members
	GetByID(ctx context.Context, projectID, memberID uuid.UUID) (*ProjectMember, error)
	GetByUserID(ctx context.Context, projectID, userID uuid.UUID) (*ProjectMember, error)
	List(ctx context.Context, projectID uuid.UUID, limit, offset int32) ([]*ProjectMember, int32, error)

	// Role management
	UpdateRole(ctx context.Context, projectID, memberID uuid.UUID, newRole string) error

	// Soft delete (setting RemovedAt / RemovedBy)
	Remove(ctx context.Context, projectID, memberID, removedBy uuid.UUID) error
}

// ProjectUseCase orchestrates the business logic for Projects (e.g., generating Database URLs)
type ProjectUseCase interface {
	Create(ctx context.Context, orgID, createdBy uuid.UUID, name, description string) (*Project, error)

	GetByID(ctx context.Context, projectID uuid.UUID) (*Project, error)

	// Scoped by OrgID to ensure data isolation
	Get(ctx context.Context, orgID, projectID uuid.UUID) (*Project, error)
	List(ctx context.Context, orgID uuid.UUID, limit, offset int32) ([]*Project, int32, error)

	// We use pointers for name and description so we can handle partial updates (PATCH)
	Update(ctx context.Context, orgID, projectID uuid.UUID, name, description *string) (*Project, error)

	// Archive(ctx context.Context, orgID, projectID uuid.UUID) error
	// Restore(ctx context.Context, orgID, projectID uuid.UUID) error

	// Requires the project name to confirm the destructive action
	Delete(ctx context.Context, orgID, projectID uuid.UUID, confirmName string) error
}

// ProjectMemberUseCase orchestrates business logic for inviting/managing members
type ProjectMemberUseCase interface {
	// Adds a user to a project. Validates if the user exists and isn't already in the project.
	AddMember(ctx context.Context, orgID, projectID, userID, addedBy uuid.UUID, role string) (*ProjectMember, error)

	GetMember(ctx context.Context, projectID, memberID uuid.UUID) (*ProjectMember, error)
	GetCurrentMember(ctx context.Context, projectID, userID uuid.UUID) (*ProjectMember, error)
	ListMembers(ctx context.Context, projectID uuid.UUID, limit, offset int32) ([]*ProjectMember, int32, error)

	ChangeRole(ctx context.Context, projectID, memberID uuid.UUID, newRole string) error

	// Removes a member. Business logic must ensure you cannot remove the LAST Project Admin!
	RemoveMember(ctx context.Context, projectID, memberID, removedBy uuid.UUID) error
	LeaveProject(ctx context.Context, projectID, userID uuid.UUID) error
}
