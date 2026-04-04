package tenantmodelv1

import (
	"time"

	"github.com/google/uuid"
)

const (
	APIKeyModeDev  = "dev"
	APIKeyModeLive = "live"
)

type ApiKey struct {
	BaseModel

	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE"`

	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraints:OnDelete:CASCADE"`

	Name      string `gorm:"size:100;not null"`
	KeyPrefix string `gorm:"size:32;not null;index"`
	KeyHash   string `gorm:"size:255;not null;uniqueIndex"`
	Mode      string `gorm:"size:10;not null;default:'dev';index"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index"`

	ExpiresAt    *time.Time `gorm:"index"`
	LastUsedAt   *time.Time
	RevokedAt    *time.Time `gorm:"index"`
	RevokedBy    *uuid.UUID `gorm:"type:uuid"`    // this is id of the user
	RequestLimit uint32     `gorm:"default:1000"` // per minute
	RateLimit    uint32     `gorm:"default:60"`   // per minute
}

func (ApiKey) TableName() string {
	return "tenant.api_keys"
}

func (ak *ApiKey) IsRevoked() bool {
	return ak.RevokedAt != nil
}

func (ak *ApiKey) IsExpired() bool {
	if ak.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ak.ExpiresAt)
}
