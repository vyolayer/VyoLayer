package apikeymodelv1

import (
	"time"

	"github.com/google/uuid"
)

/*
==================================================
Scopes
==================================================
*/

type APIKeyScope struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ApiKeyID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_api_key_scope,priority:1"`
	Scope    string    `gorm:"size:100;not null;uniqueIndex:idx_api_key_scope,priority:2"`

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"<-:update;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (APIKeyScope) TableName() string {
	return "api_key_scopes"
}
