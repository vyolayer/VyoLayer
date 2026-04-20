package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

/*
------------------------------------
Project Resource Overrides
Per project customizations:
- hide columns
- rename labels
- disable actions
- reorder
------------------------------------
*/

type ProjectResourceOverride struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ProjectID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_project_resource_override"`
	ResourceID uint64    `gorm:"not null;uniqueIndex:idx_project_resource_override"`

	Resource ServiceResource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE"`

	IsVisible *bool

	CustomLabel string `gorm:"size:100"`

	ColumnOverrides datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`
	ActionOverrides datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`
	Settings        datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`
}

func (ProjectResourceOverride) TableName() string {
	return "project_resource_overrides"
}
