package model

import (
	"gorm.io/datatypes"
)

// Filters
type ServiceResourceFilter struct {
	BaseModel

	ResourceID uint64          `gorm:"not null;uniqueIndex:idx_resource_filter"`
	Resource   ServiceResource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE"`

	Key   string `gorm:"size:63;not null;uniqueIndex:idx_resource_filter"`
	Label string `gorm:"size:100;not null"`

	// select, boolean, date, text
	Type string `gorm:"size:30;default:'text'"`

	Options datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'[]'"`

	SortOrder uint16 `gorm:"default:0"`
}

func (ServiceResourceFilter) TableName() string {
	return "service_resource_filters"
}
