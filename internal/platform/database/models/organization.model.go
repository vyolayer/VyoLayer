package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

type Organization struct {
	BaseModel

	// Info
	Name        string `gorm:"size:100;not null"`
	Slug        string `gorm:"size:100;not null;uniqueIndex"`
	Description string `gorm:"type:text"`

	// Ownership
	OwnerID uuid.UUID `gorm:"index;not null;index"` // User ID of the user who owns the organization
	Owner   User      `gorm:"foreignKey:OwnerID;constraints:OnDelete:RESTRICT;"`

	// Status
	IsActive      bool       `gorm:"default:true;index:idx_organizations_active"`
	DeactivatedBy *uuid.UUID `gorm:"type:uuid"` // User ID of the user who deactivated the organization
	DeactivatedAt *time.Time

	// Configuration and Limits
	MaxProjects int `gorm:"default:1;check:max_projects > 0 AND max_projects <= 100"`
	MaxMembers  int `gorm:"default:5;check:max_members > 0 AND max_members <= 100"`

	// Relationships
	Members []OrganizationMember `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (o *Organization) PublicID() types.OrganizationID {
	id, _ := types.ReconstructOrganizationID(o.ID.String())
	return id
}
