package models

import (
	"time"
	"worklayer/internal/platform/database/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// IAMUser struct
type User struct {
	BaseModel
	Email           string `gorm:"uniqueIndex;not null"`
	PasswordHash    string `gorm:"not null"`
	FullName        string `gorm:"not null,size:100"`
	IsEmailVerified bool
	IsActive        bool `gorm:"default:true"`
	LastLoginAt     *time.Time
}

type UserSession struct {
	BaseModel
	UserID    types.InternalID `gorm:"index;not null"`
	TokenHash string           `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	Revoked   bool   `gorm:"default:false"`
	IpAddress string `gorm:"size:50"`
	UserAgent string `gorm:"size:255"`
}
