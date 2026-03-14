package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

type OrganizationMember struct {
	BaseModel

	// foreign keys
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_member_unique,priority:2"`
	User   User      `gorm:"foreignKey:UserID;constraints:OnDelete:CASCADE;"`

	InvitedAt *time.Time
	InvitedBy *uuid.UUID `gorm:"type:uuid"` // organization member id
	JoinedAt  *time.Time `gorm:"autoCreateTime"`
	RemovedAt *time.Time `gorm:"index"`
	RemovedBy *uuid.UUID `gorm:"type:uuid"` // organization member id

	Roles []MemberOrganizationRole `gorm:"foreignKey:MemberID;references:ID;constraint:OnDelete:CASCADE"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}

func (om *OrganizationMember) IsActive() bool {
	return om.JoinedAt != nil && // joined
		om.RemovedAt == nil && // not removed
		!om.DeletedAt.Valid
}

func (om *OrganizationMember) IsOwner() bool {
	return om.Organization.OwnerID == om.UserID
}

// PublicID returns the public ID of the organization member
func (om *OrganizationMember) PublicID() types.OrganizationMemberID {
	id, _ := types.ReconstructOrganizationMemberID(om.ID.String())
	return id
}

func (om *OrganizationMember) GetInvitedBy() *types.OrganizationMemberID {
	if om.InvitedBy == nil {
		return nil
	}
	id, _ := types.ReconstructOrganizationMemberID(om.InvitedBy.String())
	return &id
}

func (om *OrganizationMember) GetRemovedBy() *types.OrganizationMemberID {
	if om.RemovedBy == nil {
		return nil
	}
	id, _ := types.ReconstructOrganizationMemberID(om.RemovedBy.String())
	return &id
}

// MemberInvitation Model
type OrganizationMemberInvitation struct {
	BaseModel

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`                                            // organization id
	InvitedBy      uuid.UUID `gorm:"type:uuid;not null"`                                                  // organization member id who created the invitation
	Email          string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_org_invitation_email_org"` // email of the person being invited
	Token          string    `gorm:"type:varchar(64);not null;uniqueIndex"`                               // unique invitation token
	RoleIDs        string    `gorm:"type:text"`                                                           // JSON array of role IDs to assign on acceptance
	InvitedAt      time.Time `gorm:"autoCreateTime"`
	IsAccepted     bool      `gorm:"default:false"`
	AcceptedAt     *time.Time
	ExpiredAt      time.Time  `gorm:"not null"`  // expiration time for the invitation
	DeletedBy      *uuid.UUID `gorm:"type:uuid"` // organization member id who deleted the invitation
}

func (OrganizationMemberInvitation) TableName() string {
	return "organization_member_invitations"
}

// OrganizationMemberInvitationPublicID returns the public ID of the organization member invitation
func (omi *OrganizationMemberInvitation) OrganizationPublicID() types.OrganizationID {
	id, _ := types.ReconstructOrganizationID(omi.OrganizationID.String())
	return id
}

func (omi *OrganizationMemberInvitation) InvitedByPublicID() types.OrganizationMemberID {
	id, _ := types.ReconstructOrganizationMemberID(omi.InvitedBy.String())
	return id
}

func (omi *OrganizationMemberInvitation) DeletedByPublicID() *types.OrganizationMemberID {
	if omi.DeletedBy == nil {
		return nil
	}
	id, _ := types.ReconstructOrganizationMemberID(omi.DeletedBy.String())
	return &id
}

// PublicID returns the public ID of the organization member invitation
func (omi *OrganizationMemberInvitation) PublicID() types.OrganizationMemberInvitationID {
	id, _ := types.ReconstructOrganizationMemberInvitationID(omi.ID.String())
	return id
}

// IsExpired checks if the invitation has expired
func (omi *OrganizationMemberInvitation) IsExpired() bool {
	return time.Now().After(omi.ExpiredAt)
}

// IsPending checks if the invitation is still pending (not accepted and not expired)
func (omi *OrganizationMemberInvitation) IsPending() bool {
	return !omi.IsAccepted && !omi.IsExpired() && !omi.DeletedAt.Valid
}
