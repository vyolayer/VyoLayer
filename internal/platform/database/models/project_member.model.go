package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

// ProjectMember represents a user's membership in a project with a minimal RBAC role.
type ProjectMember struct {
	BaseModel

	// Foreign keys
	ProjectID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_project_member_unique,priority:1"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraints:OnDelete:CASCADE;"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_project_member_unique,priority:2"`
	User   User      `gorm:"foreignKey:UserID;constraints:OnDelete:CASCADE;"`

	// Minimal RBAC role: "admin", "member", "viewer"
	Role string `gorm:"size:20;not null;default:'member';index"`

	// Tracking
	AddedBy uuid.UUID `gorm:"type:uuid;not null"`
	Adder   User      `gorm:"foreignKey:AddedBy;constraints:OnDelete:RESTRICT;"`

	JoinedAt  time.Time  `gorm:"autoCreateTime"`
	RemovedAt *time.Time `gorm:"index"`
	RemovedBy *uuid.UUID `gorm:"type:uuid"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}

func (pm *ProjectMember) IsRemoved() bool {
	return pm.RemovedAt != nil || pm.RemovedBy != nil
}

func (pm *ProjectMember) IsActive() bool {
	return !pm.IsRemoved() && !pm.DeletedAt.Valid
}

func (pm *ProjectMember) PublicID() types.ProjectMemberID {
	id, _ := types.ReconstructProjectMemberID(pm.ID.String())
	return id
}

func (pm *ProjectMember) ProjectPublicID() types.ProjectID {
	id, _ := types.ReconstructProjectID(pm.ProjectID.String())
	return id
}

func (pm *ProjectMember) UserPublicID() types.UserID {
	id, _ := types.ReconstructUserID(pm.UserID.String())
	return id
}
