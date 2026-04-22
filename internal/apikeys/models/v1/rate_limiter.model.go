package apikeymodelv1

import (
	"time"

	"github.com/google/uuid"
)

/*
==================================================
Rate Limits
==================================================
*/

type APIKeyRateLimit struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ApiKeyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	// requests per minute
	RateLimit uint32 `gorm:"type:integer;not null;default:60"`

	// monthly / total quota depending logic
	RequestLimit uint32 `gorm:"type:integer;not null;default:10000"`

	BurstLimit uint32 `gorm:"type:integer;not null;default:10"`

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (APIKeyRateLimit) TableName() string {
	return "api_key_rate_limits"
}
