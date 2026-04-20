package model

import (
	"gorm.io/datatypes"
)

type ServiceResourceAction struct {
	BaseModel

	ResourceID uint64          `gorm:"not null;uniqueIndex:idx_resource_action"`
	Resource   ServiceResource `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE"`

	Key   string `gorm:"size:63;not null;uniqueIndex:idx_resource_action"`
	Label string `gorm:"size:100;not null"`

	Scope   string `gorm:"size:63;default:'page'"`
	Variant string `gorm:"size:20;default:'secondary'"`

	Route  string `gorm:"size:255"`
	Method string `gorm:"size:10;default:'POST'"`

	IsDanger  bool   `gorm:"default:false"`
	SortOrder uint32 `gorm:"default:0"`

	Config datatypes.JSON `gorm:"type:jsonb;serializer:json;default:'{}'"`
}

func (ServiceResourceAction) TableName() string {
	return "service_resource_actions"
}
