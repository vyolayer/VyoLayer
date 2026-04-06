package v1

import (
	"time"

	"github.com/google/uuid"
)

// VerificationToken stores a one-time email-verification code tied to a user.
type VerificationToken struct {
	ID        int64     `gorm:"<-:create;primaryKey;autoIncrement"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null"`
	UsedAt    *time.Time
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (VerificationToken) TableName() string { return "verification_tokens" }
