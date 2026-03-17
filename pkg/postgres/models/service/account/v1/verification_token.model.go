package servicemodelv1

import (
	"time"

	"github.com/google/uuid"
)

type ServiceUserVerificationToken struct {
	UUID

	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`

	TokenHash string    `gorm:"type:text;not null;"`
	Type      string    `gorm:"size:20;not null;"`
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;"`
	UsedAt    *time.Time
}

func (ServiceUserVerificationToken) TableName() string {
	return "account_service.user_verification_tokens"
}
