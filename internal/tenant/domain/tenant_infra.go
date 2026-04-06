package domain

import (
	"time"

	"github.com/google/uuid"
)

type TenantInfraStatus string

func (s TenantInfraStatus) String() string {
	return string(s)
}

const (
	TenantInfraStatusCreating TenantInfraStatus = "creating"
	TenantInfraStatusReady    TenantInfraStatus = "ready"
	TenantInfraStatusDeleting TenantInfraStatus = "deleting"
	TenantInfraStatusDeleted  TenantInfraStatus = "deleted"
)

type TenantInfra struct {
	ID             int64
	OrganizationID uuid.UUID
	Schema         string
	Status         TenantInfraStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewTenantInfra(organizationID uuid.UUID, schema string) *TenantInfra {
	now := time.Now()
	return &TenantInfra{
		OrganizationID: organizationID,
		Schema:         schema,
		Status:         TenantInfraStatusCreating,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Getter

func (t *TenantInfra) GetID() int64                 { return t.ID }
func (t *TenantInfra) GetOrganizationID() uuid.UUID { return t.OrganizationID }
func (t *TenantInfra) GetSchema() string            { return t.Schema }
func (t *TenantInfra) GetStatus() TenantInfraStatus { return t.Status }
func (t *TenantInfra) GetCreatedAt() time.Time      { return t.CreatedAt }
func (t *TenantInfra) GetUpdatedAt() time.Time      { return t.UpdatedAt }

// Setter

func (t *TenantInfra) SetID(id int64)                             { t.ID = id }
func (t *TenantInfra) SetOrganizationID(organizationID uuid.UUID) { t.OrganizationID = organizationID }
func (t *TenantInfra) SetSchema(schema string)                    { t.Schema = schema }
func (t *TenantInfra) SetStatus(status TenantInfraStatus)         { t.Status = status }
func (t *TenantInfra) SetCreatedAt(createdAt time.Time)           { t.CreatedAt = createdAt }
func (t *TenantInfra) SetUpdatedAt(updatedAt time.Time)           { t.UpdatedAt = updatedAt }

func (t *TenantInfra) CompareStatus(status TenantInfraStatus) bool {
	return t.Status.String() == status.String()
}
