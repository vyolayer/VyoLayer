package models

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationRBAC model
type OrganizationRole struct {
	BaseModel

	OrgID        uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_name_unique,priority:1"`
	Organization Organization `gorm:"foreignKey:OrgID;constraints:OnDelete:CASCADE;"`

	Name        string `gorm:"size:100;not null;uniqueIndex:idx_org_role_name_unique,priority:2"`
	Description string `gorm:"type:text"`
}

func (OrganizationRole) TableName() string {
	return "organization_roles"
}

// Organization permissions
type OrganizationPermission struct {
	BaseModel

	OrgID        uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_perm_unique,priority:1"`
	Organization Organization `gorm:"foreignKey:OrgID;constraint:OnDelete:CASCADE"`

	Resource    string `gorm:"size:100;not null;uniqueIndex:idx_org_perm_unique,priority:2;index:idx_org_permissions_resource"`
	Action      string `gorm:"size:50;not null;uniqueIndex:idx_org_perm_unique,priority:3"`
	Description string `gorm:"type:text"`
}

func (OrganizationPermission) TableName() string {
	return "organization_permissions"
}

// Organization role <-> permission mapping
type OrganizationRolePermission struct {
	BaseModel

	OrgRoleID              uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:1"`
	OrgPermissionID        uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_role_perm_unique,priority:2"`
	OrganizationRole       OrganizationRole       `gorm:"foreignKey:OrgRoleID;constraint:OnDelete:CASCADE"`
	OrganizationPermission OrganizationPermission `gorm:"foreignKey:OrgPermissionID;constraint:OnDelete:CASCADE"`
}

func (OrganizationRolePermission) TableName() string {
	return "organization_role_permissions"
}

// Organization role <-> user mapping
type UserOrganizationRole struct {
	BaseModel

	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_org_role_unique,priority:1"`
	OrgID     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_org_role_unique,priority:2"`
	OrgRoleID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_org_role_unique,priority:3"`

	User             User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Organization     Organization     `gorm:"foreignKey:OrgID;constraint:OnDelete:CASCADE"`
	OrganizationRole OrganizationRole `gorm:"foreignKey:OrgRoleID;constraint:OnDelete:CASCADE"`

	GrantedBy uuid.UUID `gorm:"type:uuid"`
	GrantedAt time.Time `gorm:"autoCreateTime"`
	RevokedBy uuid.UUID `gorm:"type:uuid"`
	RevokedAt time.Time `gorm:"index"`
}

func (UserOrganizationRole) TableName() string {
	return "user_organization_roles"
}
