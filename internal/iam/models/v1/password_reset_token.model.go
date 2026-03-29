package v1

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken stores a one-time password-reset code tied to a user.
type PasswordResetToken struct {
	ID        int64     `gorm:"<-:create;primaryKey;autoIncrement"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null"`
	UsedAt    *time.Time
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (PasswordResetToken) TableName() string { return "iam.password_reset_tokens" }
