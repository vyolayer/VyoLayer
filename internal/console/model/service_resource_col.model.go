package model

import (
	"gorm.io/datatypes"
)

type ServiceResourceColumn struct {
	BaseModel

	ResourceID uint64          `gorm:"not null;uniqueIndex:idx_resource_column"`
	Resource   ServiceResource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE"`

	Key   string `gorm:"size:63;not null;uniqueIndex:idx_resource_column"`
	Label string `gorm:"size:100;not null"`

	// text, badge, image, boolean, datetime, number
	Type string `gorm:"size:30;default:'text'"`

	Sortable bool `gorm:"default:false"`
	Visible  bool `gorm:"default:true"`

	Width     uint32 `gorm:"default:0"`
	SortOrder uint32 `gorm:"default:0"`

	Config datatypes.JSON `gorm:"type:jsonb;serializer:json;default:'{}'"`
}

func (ServiceResourceColumn) TableName() string {
	return "service_resource_columns"
}
