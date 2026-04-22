package apikeymodelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/*
==================================================
Constants
==================================================
*/

const (
	APIKeyModeDev  = "dev"
	APIKeyModeLive = "live"
)

const (
	APIKeyStatusActive   = "active"
	APIKeyStatusRevoked  = "revoked"
	APIKeyStatusExpired  = "expired"
	APIKeyStatusDisabled = "disabled"
)

/*
==================================================
API Keys
==================================================
*/

type APIKey struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null;index"`

	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`

	// Visible lookup prefix
	Prefix string `gorm:"size:32;not null;uniqueIndex"`

	// Hashed full secret
	SecretHash string `gorm:"size:255;not null"`

	// dev / live
	Environment string `gorm:"size:10;not null;default:'dev';index"`

	// active / revoked / expired / disabled
	Status string `gorm:"size:20;not null;default:'active';index"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index"`

	LastUsedAt *time.Time `gorm:"index"`
	LastUsedIP string     `gorm:"size:64"`
	LastUsedUA string     `gorm:"size:500"`

	ExpiresAt *time.Time `gorm:"index"`

	RevokedBy *uuid.UUID `gorm:"type:uuid"`
	RevokedAt *time.Time

	Metadata datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`

	CreatedAt time.Time      `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (APIKey) TableName() string {
	return "api_keys"
}

func (ak *APIKey) IsRevoked() bool {
	return ak.RevokedAt != nil || ak.Status == APIKeyStatusRevoked
}

func (ak *APIKey) IsExpired() bool {
	if ak.Status == APIKeyStatusExpired {
		return true
	}
	if ak.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ak.ExpiresAt)
}

func (ak *APIKey) IsUsable() bool {
	if ak.Status != APIKeyStatusActive {
		return false
	}
	if ak.IsRevoked() {
		return false
	}
	if ak.IsExpired() {
		return false
	}
	return true
}

func (ak *APIKey) Revoke(by uuid.UUID, at time.Time) {
	ak.Status = APIKeyStatusRevoked
	ak.RevokedBy = &by
	ak.RevokedAt = &at
}
