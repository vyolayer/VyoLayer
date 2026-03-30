package tenantmodelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// --- Organization
type Organization struct {
	BaseModel

	Name        string `gorm:"size:50;not null"`
	Slug        string `gorm:"size:100;not null;uniqueIndex"`
	Description string `gorm:"type:text"`

	OwnerID uuid.UUID `gorm:"type:uuid;not null;index"` // User ID

	IsActive      bool       `gorm:"default:true;index:idx_organizations_active"`
	DeactivatedBy *uuid.UUID `gorm:"type:uuid"` // User ID
	DeactivatedAt *time.Time

	MaxProjects int `gorm:"default:1;check:max_projects > 0 AND max_projects <= 100"`
	MaxMembers  int `gorm:"default:5;check:max_members > 0 AND max_members <= 100"`

	Members     []OrganizationMember           `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Projects    []Project                      `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Invitations []OrganizationMemberInvitation `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

func (Organization) TableName() string {
	return "tenant.organizations"
}

// --- Organization Member
type OrganizationMember struct {
	BaseModel

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:2"`

	InvitedBy *uuid.UUID `gorm:"type:uuid"`
	InvitedAt *time.Time
	JoinedAt  *time.Time `gorm:"autoCreateTime"`
	RemovedBy *uuid.UUID `gorm:"type:uuid"`
	RemovedAt *time.Time `gorm:"index"`

	Roles []MemberOrganizationRole `gorm:"foreignKey:MemberID;references:ID;constraint:OnDelete:CASCADE"`
}

func (OrganizationMember) TableName() string {
	return "tenant.organization_members"
}

func (om *OrganizationMember) IsActive() bool {
	return om.JoinedAt != nil && om.RemovedAt == nil && !om.DeletedAt.Valid
}

// --- Organization Member Invitation
type OrganizationMemberInvitation struct {
	BaseModel

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_invitation_email_org,priority:2;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE"`

	InvitedBy uuid.UUID      `gorm:"type:uuid;not null;index"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_org_invitation_email_org,priority:1"`
	Token     string         `gorm:"type:varchar(128);not null;uniqueIndex"`
	RoleIDs   datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'"`

	InvitedAt  time.Time `gorm:"autoCreateTime"`
	IsAccepted bool      `gorm:"default:false"`
	AcceptedAt *time.Time
	ExpiredAt  time.Time  `gorm:"not null;index"`
	DeletedBy  *uuid.UUID `gorm:"type:uuid"`
}

func (OrganizationMemberInvitation) TableName() string {
	return "tenant.organization_member_invitations"
}

func (omi *OrganizationMemberInvitation) IsExpired() bool {
	return time.Now().After(omi.ExpiredAt)
}

func (omi *OrganizationMemberInvitation) IsPending() bool {
	return !omi.IsAccepted && !omi.IsExpired() && !omi.DeletedAt.Valid
}

// --- Organization Role
type OrganizationRole struct {
	BaseModel

	Name         string `gorm:"size:100;not null;uniqueIndex:idx_org_role_name_unique,priority:2"`
	Description  string `gorm:"type:text"`
	IsSystemRole bool   `gorm:"column:is_system;default:false"`
	IsDefault    bool   `gorm:"default:false"`

	Permissions []OrganizationPermission `gorm:"many2many:organization_role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
}

func (OrganizationRole) TableName() string {
	return "tenant.organization_roles"
}

// --- Organization Permission
type OrganizationPermission struct {
	BaseModel

	Resource    string `gorm:"size:100;not null;uniqueIndex:idx_org_perm_unique,priority:1;index:idx_org_permissions_resource"`
	Action      string `gorm:"size:50;not null;uniqueIndex:idx_org_perm_unique,priority:2"`
	Code        string `gorm:"-"`
	Group       string `gorm:"size:50;not null;index:idx_org_permissions_group"`
	Description string `gorm:"type:text"`
	IsSystem    bool   `gorm:"default:false"`
}

func (OrganizationPermission) TableName() string {
	return "tenant.organization_permissions"
}

type OrganizationRolePermission struct {
	BaseModel

	RoleID       uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:1"`
	PermissionID uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:2"`
	Role         OrganizationRole       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission   OrganizationPermission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}

func (OrganizationRolePermission) TableName() string {
	return "tenant.organization_role_permissions"
}

// --- Member Organization Role
type MemberOrganizationRole struct {
	BaseModel

	MemberID       uuid.UUID        `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;not null;index"`
	RoleID         uuid.UUID        `gorm:"type:uuid;not null;index"`
	Role           OrganizationRole `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	GrantedBy      *uuid.UUID       `gorm:"type:uuid"`
	GrantedAt      time.Time        `gorm:"autoCreateTime"`
	RevokedBy      *uuid.UUID       `gorm:"type:uuid"`
	RevokedAt      *time.Time       `gorm:"index"`
}

func (MemberOrganizationRole) TableName() string {
	return "tenant.member_organization_roles"
}
