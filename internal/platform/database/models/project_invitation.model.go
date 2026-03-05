package models

import (
	"time"
	"worklayer/internal/platform/database/types"

	"github.com/google/uuid"
)

// ProjectInvitation tracks pending invitations to join a project.
type ProjectInvitation struct {
	BaseModel

	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	InvitedBy uuid.UUID `gorm:"type:uuid;not null"` // project member who created the invitation

	Email string `gorm:"type:varchar(255);not null;uniqueIndex:idx_project_invitation_email_project"` // email of the person being invited
	Role  string `gorm:"size:20;not null;default:'member'"`                                           // role to assign on acceptance
	Token string `gorm:"type:varchar(64);not null;uniqueIndex"`                                       // unique invitation token

	InvitedAt  time.Time `gorm:"autoCreateTime"`
	IsAccepted bool      `gorm:"default:false"`
	AcceptedAt *time.Time
	ExpiredAt  time.Time  `gorm:"not null"` // expiration time for the invitation
	DeletedBy  *uuid.UUID `gorm:"type:uuid"`
}

func (ProjectInvitation) TableName() string {
	return "project_invitations"
}

func (pi *ProjectInvitation) PublicID() types.ProjectInvitationID {
	id, _ := types.ReconstructProjectInvitationID(pi.ID.String())
	return id
}

func (pi *ProjectInvitation) ProjectPublicID() types.ProjectID {
	id, _ := types.ReconstructProjectID(pi.ProjectID.String())
	return id
}

func (pi *ProjectInvitation) InvitedByPublicID() types.ProjectMemberID {
	id, _ := types.ReconstructProjectMemberID(pi.InvitedBy.String())
	return id
}

// IsExpired checks if the invitation has expired
func (pi *ProjectInvitation) IsExpired() bool {
	return time.Now().After(pi.ExpiredAt)
}

// IsPending checks if the invitation is still pending (not accepted and not expired)
func (pi *ProjectInvitation) IsPending() bool {
	return !pi.IsAccepted && !pi.IsExpired() && !pi.DeletedAt.Valid
}
