package models

import (
	"time"
	"worklayer/internal/platform/database/types"

	"github.com/google/uuid"
)

type OrganizationMember struct {
	BaseModel

	// foreign keys
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:2"`
	User   User      `gorm:"foreignKey:UserID;constraints:OnDelete:CASCADE;"`

	InvitedAt *time.Time
	InvitedBy *uuid.UUID `gorm:"type:uuid"`
	JoinedAt  *time.Time `gorm:"autoCreateTime"`
	RemovedAt *time.Time `gorm:"index"`
	RemovedBy *uuid.UUID `gorm:"type:uuid"`
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

func (om *OrganizationMember) IsOwner() bool {
	return om.Organization.OwnerID == om.UserID
}

// PublicID returns the public ID of the organization member
func (om *OrganizationMember) PublicID() types.OrganizationMemberID {
	id, _ := types.ReconstructOrganizationMemberID(om.ID.String())
	return id
}
