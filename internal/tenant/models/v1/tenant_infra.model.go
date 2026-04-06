package tenantmodelv1

import (
	"time"

	"github.com/google/uuid"
)

type TenantInfra struct {
	ID             int64     `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Schema         string    `gorm:"size:63;not null;uniqueIndex"`

	Status string `gorm:"size:10;default:'ready'"` // creating, ready, failed

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TenantInfra) TableName() string {
	return "tenant_infra"
}
