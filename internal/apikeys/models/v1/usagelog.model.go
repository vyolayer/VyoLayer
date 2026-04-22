package apikeymodelv1

import (
	"time"

	"github.com/google/uuid"
)

/*
==================================================
Usage Logs (High Volume Optional)
==================================================
*/

type APIKeyUsageLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ApiKeyID       uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null;index"`

	Path       string `gorm:"size:255"`
	Method     string `gorm:"size:10"`
	StatusCode uint16

	IP        string `gorm:"size:64"`
	UserAgent string `gorm:"size:500"`

	LatencyMs uint32

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (APIKeyUsageLog) TableName() string {
	return "api_key_usage_logs"
}
