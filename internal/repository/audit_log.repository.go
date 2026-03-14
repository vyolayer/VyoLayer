package repository

import (
	"context"
	"encoding/json"

	"vyolayer/internal/platform/database/models"
	"vyolayer/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(ctx context.Context, entry *AuditLogEntry) *errors.AppError
}

// AuditLogEntry is the convenience input struct for creating audit logs.
type AuditLogEntry struct {
	OrganizationID uuid.UUID

	// Actor
	ActorID   uuid.UUID
	ActorType string // "user" | "member"

	Action string // e.g. "org.updated", "member.removed"

	// Primary resource
	ResourceType string    // "organization" | "member" | "invitation"
	ResourceID   uuid.UUID // internal UUID of the resource

	// Optional secondary resource
	SecondaryResourceType string
	SecondaryResourceID   uuid.UUID

	// Metadata (arbitrary key-value pairs)
	Metadata map[string]interface{}

	// Request context (optional)
	IPAddress string
	UserAgent string
	RequestID uuid.UUID

	// Severity: "info" (default) | "warning" | "critical"
	Severity string
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, entry *AuditLogEntry) *errors.AppError {
	var metadataJSON datatypes.JSON
	if entry.Metadata != nil {
		data, err := json.Marshal(entry.Metadata)
		if err != nil {
			return errors.InternalWrap(err, "marshaling audit log metadata")
		}
		metadataJSON = datatypes.JSON(data)
	}

	severity := entry.Severity
	if severity == "" {
		severity = "info"
	}

	log := &models.AuditLog{
		OrganizationID:        entry.OrganizationID,
		ActorID:               entry.ActorID,
		ActorType:             entry.ActorType,
		Action:                entry.Action,
		ResourceType:          entry.ResourceType,
		ResourceID:            entry.ResourceID,
		SecondaryResourceType: entry.SecondaryResourceType,
		SecondaryResourceID:   entry.SecondaryResourceID,
		Metadata:              metadataJSON,
		IPAddress:             entry.IPAddress,
		UserAgent:             entry.UserAgent,
		RequestID:             entry.RequestID,
		Severity:              severity,
	}

	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return ConvertDBError(err, "creating audit log")
	}

	return nil
}
