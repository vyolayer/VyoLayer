package v1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionID = int64
type Session struct {
	ID SessionID `gorm:"<-:create;primaryKey;autoIncrement"`

	UserID    uuid.UUID `gorm:"type:uuid;index;not null;"`
	TokenHash string    `gorm:"type:text;uniqueIndex;not null;"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;"`
	RevokedAt *time.Time
	Reason    string `gorm:"size:255"`
	IpAddress string `gorm:"size:50"`
	UserAgent string `gorm:"size:255"`

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt
}

func (s Session) TableName() string {
	return "sessions"
}
