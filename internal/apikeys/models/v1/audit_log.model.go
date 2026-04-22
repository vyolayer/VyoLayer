package apikeymodelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

/*
==================================================
Audit Logs
==================================================
*/

type APIKeyAuditLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ApiKeyID       uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null;index"`

	// created / rotated / revoked / scope_updated / used
	Action string `gorm:"size:50;not null;index"`

	ActorID *uuid.UUID `gorm:"type:uuid"`

	IP        string `gorm:"size:64"`
	UserAgent string `gorm:"size:500"`

	Metadata datatypes.JSON `gorm:"type:jsonb;serializer:json;not null;default:'{}'"`

	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (APIKeyAuditLog) TableName() string {
	return "api_key_audit_logs"
}
