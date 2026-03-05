package models

import (
	"worklayer/internal/platform/database/types"

	"github.com/google/uuid"
)

type Project struct {
	BaseModel

	// Organization relationship
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE;"`

	// Info
	Name        string `gorm:"size:100;not null;uniqueIndex:idx_project_org_slug,priority:2"`
	Slug        string `gorm:"size:100;not null;uniqueIndex:idx_project_org_slug,priority:1"`
	Description string `gorm:"type:text"`

	// Status
	IsActive bool `gorm:"default:true;index"`

	// Ownership
	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index"`
	Creator   User      `gorm:"foreignKey:CreatedBy;constraints:OnDelete:RESTRICT;"`

	// Configuration and Limits
	MaxApiKeys int `gorm:"default:5;check:max_api_keys > 0 AND max_api_keys <= 10"`
	MaxMembers int `gorm:"default:5;check:max_members > 0 AND max_members <= 10"`

	// Relationships
	Members []ProjectMember `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

func (Project) TableName() string {
	return "projects"
}

func (p *Project) PublicID() types.ProjectID {
	id, _ := types.ReconstructProjectID(p.ID.String())
	return id
}

func (p *Project) OrganizationPublicID() types.OrganizationID {
	id, _ := types.ReconstructOrganizationID(p.OrganizationID.String())
	return id
}
