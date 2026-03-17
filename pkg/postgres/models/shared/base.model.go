package sharedmodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeStamps struct {
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type UUID struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
}

type TimeStampsWithSoftDelete struct {
	CreatedAt time.Time      `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;default:null"`
}
