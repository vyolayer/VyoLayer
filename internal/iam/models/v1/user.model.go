package v1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserID = uuid.UUID
type User struct {
	ID              UserID `gorm:"<-:create;type:uuid;primaryKey"`
	Email           string `gorm:"type:varchar(255);not null;"`
	IsEmailVerified bool   `gorm:"default:false;"`
	PasswordHash    string `gorm:"type:text;"`
	FullName        string `gorm:"type:varchar(100);not null;"`
	Status          string `gorm:"type:varchar(10);not null"`
	LastLoginAt     *time.Time

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt

	// Avatar
	AvatarID AvatarID `gorm:"type:text;"`
	Avatar   Avatar
}

func (u User) TableName() string {
	return "users"
}

// Avatar is a struct that contains the avatar URL and fallback character and color
type AvatarID = int64
type Avatar struct {
	ID AvatarID `gorm:"<-:create;primaryKey;autoIncrement"`

	URL           string `gorm:"type:text;"`
	FallbackChar  string `gorm:"type:varchar(1);"`
	FallbackColor string `gorm:"type:varchar(7);"`

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
}

func (a Avatar) TableName() string {
	return "avatars"
}
