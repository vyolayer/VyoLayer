package models

import (
	"time"

	"github.com/google/uuid"
)

// ApiKeyUsageLog records each API request made using an API key for analytics and rate-limiting.
type ApiKeyUsageLog struct {
	ID        uuid.UUID `gorm:"<-:create;type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`

	ApiKeyID       uuid.UUID `gorm:"type:uuid;not null;index"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`

	Endpoint   string `gorm:"size:255;not null"`
	Method     string `gorm:"size:10;not null"`
	StatusCode int    `gorm:"not null"`
	IPAddress  string `gorm:"size:45"`
	UserAgent  string `gorm:"size:255"`
}

func (ApiKeyUsageLog) TableName() string {
	return "api_key_usage_logs"
}
