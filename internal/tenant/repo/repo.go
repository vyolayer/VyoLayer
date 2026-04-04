package tenantrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"gorm.io/gorm"
)

// --- Gorm Repository Interface ---
type GormRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB) error
}

// --- Organization Write Interface ---
type OrganizationWrite interface {
	Create(ctx context.Context, tx *gorm.DB, org *domain.Organization) error
	Update(ctx context.Context, org *domain.Organization) error

	Delete(ctx context.Context, orgID uuid.UUID, confirmName string) error
	Archive(ctx context.Context, orgID uuid.UUID, confirmName string) error
	Restore(ctx context.Context, orgID uuid.UUID) error

	UpdateOwner(ctx context.Context, tx *gorm.DB, orgID uuid.UUID, ownerID uuid.UUID) error
}

// --- Organization Read Interface ---
type OrganizationRead interface {
	GetByID(ctx context.Context, orgID uuid.UUID) (*domain.Organization, error)
	GetByIDWithMember(ctx context.Context, orgID uuid.UUID) (*domain.OrganizationWithMember, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	List(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Organization, error)
}

// --- Organization Repository Interface ---
type OrganizationRepository interface {
	GormRepository
	OrganizationWrite
	OrganizationRead
}

// --- Organization RBAC Interfaces ---

// Organization Role Interfaces
type OrganizationRoleRead interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationRole, error)
	GetByName(ctx context.Context, name string) (*domain.OrganizationRole, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationRole, error)
}

type OrganizationRoleRepository interface {
	OrganizationRoleRead
}

// Organization Permission Interfaces
type OrganizationPermissionRead interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationPermission, error)
	GetByCode(ctx context.Context, code string) (*domain.OrganizationPermission, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationPermission, error)
}

type OrganizationPermissionRepository interface {
	OrganizationPermissionRead
}

// Organization Role Permission Interfaces
type OrganizationRolePermissionRead interface {
	GetByRoleId(ctx context.Context, roleID uuid.UUID) ([]*domain.OrganizationRolePermission, error)
}

type OrganizationRolePermissionRepository interface {
	OrganizationRolePermissionRead
}

// Member Organization Role Interfaces
type MemberOrganizationRoleWrite interface {
	AddRole(ctx context.Context, tx *gorm.DB, role *domain.MemberOrganizationRole) error
	RemoveRole(ctx context.Context, role *domain.MemberOrganizationRole) error
	UpdateRole(ctx context.Context, tx *gorm.DB, memberID, roleID uuid.UUID) error
}

type MemberOrganizationRoleRead interface {
	GetByMemberId(ctx context.Context, memberID uuid.UUID) ([]*domain.MemberOrganizationRole, error)
}

type MemberOrganizationRoleRepository interface {
	GormRepository
	MemberOrganizationRoleWrite
	MemberOrganizationRoleRead
}

// Permission Checker Interface
type PermissionChecker interface {
	HasPermission(ctx context.Context, orgID, userID uuid.UUID, requiredPermissionCode string) (bool, error)
	IsMember(ctx context.Context, orgID, userID uuid.UUID) (bool, error)
}

// --- Organization Member Write Interface ---
type OrganizationMemberWrite interface {
	AddMember(ctx context.Context, tx *gorm.DB, member *domain.OrganizationMember) error
	RemoveMember(ctx context.Context, member *domain.OrganizationMember) error
	UpdateMember(ctx context.Context, member *domain.OrganizationMember) error
}

// --- Organization Member Read Interface ---
type OrganizationMemberRead interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error)
	GetByUserIdAndOrgId(ctx context.Context, userID, orgID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberWithRoles, error)
	GetByOrgIdAndEmail(ctx context.Context, orgID uuid.UUID, email string) (*domain.OrganizationMember, error)
}

// --- Organization Member Repository Interface ---
type OrganizationMemberRepository interface {
	GormRepository
	OrganizationMemberWrite
	OrganizationMemberRead
}

// --- Organization Member Invitation Interfaces ---
type OrganizationMemberInvitationWrite interface {
	Create(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error
	Delete(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error
	Accept(ctx context.Context, invitation *domain.OrganizationMemberInvitation) error
}

type OrganizationMemberInvitationRead interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationMemberInvitation, error)
	GetByToken(ctx context.Context, token string) (*domain.OrganizationMemberInvitation, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberInvitation, error)
	ListByUserEmail(ctx context.Context, email string) ([]*domain.OrganizationMemberInvitation, error)
	ListPendingByOrg(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberInvitationWithInviter, error)
	ListByInvitedBy(ctx context.Context, organizationID, invitedBy uuid.UUID) ([]*domain.OrganizationMemberInvitation, error)
}

type OrganizationMemberInvitationRepository interface {
	GormRepository
	OrganizationMemberInvitationWrite
	OrganizationMemberInvitationRead
}
