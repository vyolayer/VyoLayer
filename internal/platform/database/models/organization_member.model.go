package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationMember struct {
	BaseModel

	// foreign keys
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:2"`
	User   User      `gorm:"foreignKey:UserID;constraints:OnDelete:CASCADE;"`

	InvitedBy     *uuid.UUID `gorm:"type:uuid"`
	InvitedByUser *User      `gorm:"foreignKey:InvitedBy"`
	InvitedAt     *time.Time
	JoinedAt      *time.Time `gorm:"autoCreateTime"`

	RemovedBy *uuid.UUID `gorm:"type:uuid"`
	RemovedAt *time.Time `gorm:"index"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}

func (om *OrganizationMember) IsActive() bool {
	return om.JoinedAt != nil && // joined
		om.RemovedAt == nil && // not removed
		om.DeletedAt == nil && // not deleted
		om.Organization.IsActive // organization is active
}
