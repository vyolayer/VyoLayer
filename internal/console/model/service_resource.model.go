package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
------------------------------------
Service resources
Ex: users, sessions, audit_logs
------------------------------------
*/
type ServiceResource struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ServiceID uint64  `gorm:"not null;uniqueIndex:idx_service_resource"`
	Service   Service `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:CASCADE"`

	Key         string `gorm:"size:63;not null;uniqueIndex:idx_service_resource"`
	Label       string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`

	Icon   string `gorm:"size:63"`
	Route  string `gorm:"size:255"`
	Method string `gorm:"size:10;default:'GET'"`

	SortOrder uint32 `gorm:"default:0"`
	IsVisible bool   `gorm:"default:true"`

	Supports datatypes.JSON `gorm:"type:jsonb;serializer:json;default:'{}'"`

	CreatedAt time.Time      `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ServiceResource) TableName() string {
	return "service_resources"
}
