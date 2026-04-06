package domain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TenantInfraRepository interface {
	Create(ctx context.Context, tx *gorm.DB, tenant *TenantInfra) error
	GetByOrgID(ctx context.Context, orgID uuid.UUID) (*TenantInfra, error)
	GetByID(ctx context.Context, id int64) (*TenantInfra, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status TenantInfraStatus) error
	Delete(ctx context.Context, id int64) error
}
