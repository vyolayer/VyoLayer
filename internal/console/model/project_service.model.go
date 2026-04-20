package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProjectService struct {
	BaseModel

	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	ServiceID uint64    `gorm:"not null;index"`
	Service   Service   `gorm:"foreignKey:ServiceID;references:ID;constraint:OnDelete:RESTRICT"`

	Status string `gorm:"size:10;default:'pending'"`

	Plan string `gorm:"size:10;default:'free'"` // free, premium, enterprise

	Config   datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`
	Metadata datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`

	EnabledAt   *time.Time
	SuspendedAt *time.Time
}

func (ProjectService) TableName() string {
	return "project_services"
}
