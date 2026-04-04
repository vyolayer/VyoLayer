package tenantmodelv1

import (
	"time"

	"github.com/google/uuid"
)

// --- Project
type Project struct {
	BaseModel

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE"`

	Name        string `gorm:"size:50;not null;uniqueIndex:idx_project_org_slug,priority:2"`
	Slug        string `gorm:"size:100;not null;uniqueIndex:idx_project_org_slug,priority:1"`
	Description string `gorm:"type:text"`

	IsActive bool `gorm:"default:true;index"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index"`

	MaxAPIKeys uint8 `gorm:"column:max_api_keys;default:5;check:max_api_keys > 0 AND max_api_keys <= 10"`
	MaxMembers uint8 `gorm:"default:5;check:max_members > 0 AND max_members <= 10"`

	Members []ProjectMember `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	APIKeys []ApiKey        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

func (Project) TableName() string {
	return "tenant.projects"
}

// --- Project Member
type ProjectMember struct {
	BaseModel

	ProjectID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_project_member_unique,priority:1"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraints:OnDelete:CASCADE"`

	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_project_member_unique,priority:2"`
	Role   string    `gorm:"size:20;not null;default:'member';index"`

	AddedBy uuid.UUID `gorm:"type:uuid;not null"`

	JoinedAt  time.Time  `gorm:"autoCreateTime"`
	RemovedAt *time.Time `gorm:"index"`
	RemovedBy *uuid.UUID `gorm:"type:uuid"`
}

func (ProjectMember) TableName() string {
	return "tenant.project_members"
}

func (pm *ProjectMember) IsRemoved() bool {
	return pm.RemovedAt != nil || pm.RemovedBy != nil
}

func (pm *ProjectMember) IsActive() bool {
	return !pm.IsRemoved() && !pm.DeletedAt.Valid
}
