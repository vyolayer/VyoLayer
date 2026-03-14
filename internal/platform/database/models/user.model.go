package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"gorm.io/gorm"
)

type TimeStamps struct {
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
}

type BaseModel struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;primaryKey"`
	TimeStamps
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID.String() == "00000000-0000-0000-0000-000000000000" || b.ID.String() == "" {
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

func (u User) TableName() string {
	return "users"
}

func (u *User) PublicID() types.UserID {
	uid, _ := types.ReconstructUserID(u.ID.String())
	return uid
}

type UserSession struct {
	BaseModel
	UserID    uuid.UUID `gorm:"index;not null"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	Revoked   bool   `gorm:"default:false"`
	Reason    string `gorm:"size:255"`
	IpAddress string `gorm:"size:50"`
	UserAgent string `gorm:"size:255"`
}
