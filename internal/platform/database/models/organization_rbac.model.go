package models

import (
	"time"
	"vyolayer/internal/platform/database/types"

	"github.com/google/uuid"
)

// OrganizationRBAC model
type OrganizationRole struct {
	BaseModel

	// Info
	Name        string `gorm:"size:100;not null;uniqueIndex:idx_org_role_name_unique,priority:2"`
	Description string `gorm:"type:text"`

	// Flags
	IsSystem  bool `gorm:"default:false"` // This role is a system role, cannot be deleted
	IsDefault bool `gorm:"default:false"` // This role is the default role, auto assigned to new members

	// relationships
	Permissions []OrganizationPermission `gorm:"many2many:organization_role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
}

func (OrganizationRole) TableName() string {
	return "organization_roles"
}

// PublicID returns the public ID of the organization role
func (or *OrganizationRole) PublicID() types.OrganizationRoleID {
	id, _ := types.ReconstructOrganizationRoleID(or.ID.String())
	return id
}

// Organization permissions
type OrganizationPermission struct {
	BaseModel

	// Info
	Resource    string `gorm:"size:100;not null;uniqueIndex:idx_org_perm_unique,priority:1;index:idx_org_permissions_resource"`
	Action      string `gorm:"size:50;not null;uniqueIndex:idx_org_perm_unique,priority:2"`
	Description string `gorm:"type:text"`
	Group       string `gorm:"size:20;not null;index:idx_org_permissions_group"`

	// Flags
	IsSystem bool `gorm:"default:false"` // This permission is a system permission, cannot be deleted
}

func (OrganizationPermission) TableName() string {
	return "organization_permissions"
}

// PublicID returns the public ID of the organization permission
func (op *OrganizationPermission) PublicID() types.OrganizationPermissionID {
	id, _ := types.ReconstructOrganizationPermissionID(op.ID.String())
	return id
}

// Organization role <-> permission mapping
type OrganizationRolePermission struct {
	BaseModel

	// foreign keys
	RoleID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:1"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:2"`

	// relationships
	Role       OrganizationRole       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission OrganizationPermission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

func (OrganizationRolePermission) TableName() string {
	return "organization_role_permissions"
}

// Organization role <-> user mapping
type MemberOrganizationRole struct {
	BaseModel

	// foreign keys
	MemberID       uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID         uuid.UUID `gorm:"type:uuid;not null;index"`

	// relationships
	// Member       OrganizationMember `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
	// Organization Organization       `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Role OrganizationRole `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`

	GrantedBy uuid.UUID `gorm:"type:uuid"` // Member ID of the user who granted the role
	GrantedAt time.Time `gorm:"autoCreateTime"`
	RevokedBy uuid.UUID `gorm:"type:uuid"` // Member ID of the user who revoked the role
	RevokedAt time.Time `gorm:"index"`
}

func (MemberOrganizationRole) TableName() string {
	return "member_organization_roles"
}

// PublicID returns the public ID of the member organization role
func (mor *MemberOrganizationRole) PublicID() types.MemberOrganizationRoleID {
	id, _ := types.ReconstructMemberOrganizationRoleID(mor.ID.String())
	return id
}
