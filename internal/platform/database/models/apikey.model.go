package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

// ApiKey represents a programmatic access key scoped to a project.
// Keys operate in either "dev" or "live" mode with different rate/request limits.
// Only the SHA-256 hash of the key is stored; the raw key is returned once at creation.
type ApiKey struct {
	BaseModel

	// Organization relationship
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraints:OnDelete:CASCADE;"`

	// Project relationship
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraints:OnDelete:CASCADE;"`

	// Key info
	Name      string `gorm:"size:100;not null"`             // user-friendly label
	KeyPrefix string `gorm:"size:16;not null;index"`        // first chars for identification (e.g. "wl_live_ab3f")
	KeyHash   string `gorm:"size:255;not null;uniqueIndex"` // SHA-256 hash of the full key

	// Mode: "live" or "dev"
	Mode string `gorm:"size:10;not null;default:'dev';index"`

	// Ownership
	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	Creator   User      `gorm:"foreignKey:CreatedBy;constraints:OnDelete:RESTRICT;"`

	// Lifecycle
	ExpiresAt  *time.Time `gorm:"index"`
	LastUsedAt *time.Time
	RevokedAt  *time.Time `gorm:"index"`
	RevokedBy  *uuid.UUID `gorm:"type:uuid"`

	// Limitations
	RequestLimit int `gorm:"default:1000"` // max requests per day (0 = unlimited)
	RateLimit    int `gorm:"default:60"`   // max requests per minute
}

func (ApiKey) TableName() string {
	return "api_keys"
}

func (ak *ApiKey) PublicID() types.ApiKeyID {
	id, _ := types.ReconstructApiKeyID(ak.ID.String())
	return id
}

func (ak *ApiKey) ProjectPublicID() types.ProjectID {
	id, _ := types.ReconstructProjectID(ak.ProjectID.String())
	return id
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
