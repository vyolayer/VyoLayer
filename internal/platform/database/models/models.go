package models

import (
	"time"

	"gorm.io/gorm"
)

// IAMUser struct
type User struct {
	gorm.Model
	Email           string `gorm:"uniqueIndex;not null"`
	PasswordHash    string `gorm:"not null"`
	FullName        string `gorm:"not null,size:100"`
	IsEmailVerified bool
	IsActive        bool `gorm:"default:true"`
	LastLoginAt     *time.Time
}

type UserSession struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null"`
	TokenHash string `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time
	Revoked   bool   `gorm:"default:false"`
	IpAddress string `gorm:"size:50"`
	UserAgent string `gorm:"size:255"`
}
