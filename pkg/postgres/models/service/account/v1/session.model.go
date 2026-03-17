package servicemodelv1

import (
	"time"

	"github.com/google/uuid"
)

type ServiceUserSession struct {
	UUID
	TimeStamps

	ProjectID uuid.UUID `gorm:"type:uuid;not null;"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;"`
	TokenHash string    `gorm:"type:text;not null;"`

	ExpiresAt time.Time  `gorm:"type:timestamp;not null;"`
	RevokedAt *time.Time `gorm:"type:timestamp"`
	Reason    string     `gorm:"size:255"`
	IpAddress string     `gorm:"size:50"`
	UserAgent string     `gorm:"size:255"`
}

func (ServiceUserSession) TableName() string {
	return "account_service.session"
}

func (s ServiceUserSession) IsValid() bool {
	return s.RevokedAt.IsZero() && time.Now().Before(s.ExpiresAt)
}
