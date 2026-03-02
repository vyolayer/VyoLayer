package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog represents an append-only audit trail for organization actions.
type AuditLog struct {
	ID uuid.UUID `gorm:"<-:create;type:uuid;primaryKey;default:gen_random_uuid()"`

	// OrganizationID is the ID of the organization associated with the audit log.
	OrganizationID uuid.UUID `gorm:"type:uuid;index"`

	// ActorId info
	ActorID   uuid.UUID `gorm:"type:uuid;index"`  // Public ID of the actor
	ActorType string    `gorm:"size:50;not null"` // e.g. "user", "member"

	Action string `gorm:"size:100;not null;index"` // e.g. "org.created", "member.invited"

	// Resource info
	ResourceType string    `gorm:"size:50;not null;index"` // e.g. "organization", "member", "invitation"
	ResourceID   uuid.UUID `gorm:"type:uuid;index"`        // Public ID of the resource

	// Secondary info
	SecondaryResourceType string    `gorm:"size:50"`
	SecondaryResourceID   uuid.UUID `gorm:"type:uuid"`

	// Metadata
	Metadata datatypes.JSON `gorm:"type:jsonb"`

	// Request info
	IPAddress string `gorm:"size:50"`
	UserAgent string `gorm:"size:255"`

	// RequestID
	RequestID uuid.UUID `gorm:"type:uuid;index"`

	// Severity classification
	Severity string `gorm:"size:20;default:'info';index"`

	// Timestamp
	CreatedAt time.Time `gorm:"<-:create;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
