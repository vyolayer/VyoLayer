package tenantrepo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type tenantInfraRepo struct {
	db     *gorm.DB
	logger *logger.AppLogger
}

func NewTenantInfraRepo(db *gorm.DB, logger *logger.AppLogger) domain.TenantInfraRepository {
	return &tenantInfraRepo{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new TenantInfra record.
func (r *tenantInfraRepo) Create(ctx context.Context, tx *gorm.DB, tenant *domain.TenantInfra) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := &TenantInfra{
		OrganizationID: tenant.OrganizationID,
		Schema:         tenant.Schema,
		Status:         tenant.Status.String(),
		CreatedAt:      tenant.CreatedAt,
		UpdatedAt:      tenant.UpdatedAt,
	}

	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.ErrorWithErr("tenantInfraRepo.Create: failed to insert tenant infra", err)
		return ConvertDBError(err, "tenant_infra")
	}

	// Reflect back the DB-generated ID to the caller.
	tenant.SetID(model.ID)
	tenant.SetCreatedAt(model.CreatedAt)
	tenant.SetUpdatedAt(model.UpdatedAt)
	return nil
}

// GetByOrgID fetches the TenantInfra record for the given organisation.
func (r *tenantInfraRepo) GetByOrgID(ctx context.Context, orgID uuid.UUID) (*domain.TenantInfra, error) {
	var model TenantInfra

	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.ErrorWithErr("tenantInfraRepo.GetByOrgID: failed to fetch tenant infra", err)
		return nil, ConvertDBError(err, "tenant_infra")
	}

	return toTenantInfraDomain(&model), nil
}

// GetByID fetches a TenantInfra record by its primary key.
func (r *tenantInfraRepo) GetByID(ctx context.Context, id int64) (*domain.TenantInfra, error) {
	var model TenantInfra

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.ErrorWithErr("tenantInfraRepo.GetByID: failed to fetch tenant infra", err)
		return nil, ConvertDBError(err, "tenant_infra")
	}

	return toTenantInfraDomain(&model), nil
}

// UpdateStatus updates only the status column (and updated_at) for a given ID.
func (r *tenantInfraRepo) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status domain.TenantInfraStatus) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).
		Model(&TenantInfra{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status.String(),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		r.logger.ErrorWithErr("tenantInfraRepo.UpdateStatus: failed to update status", result.Error)
		return ConvertDBError(result.Error, "tenant_infra")
	}

	if result.RowsAffected == 0 {
		return NotFoundError("TenantInfra", "id")
	}

	return nil
}

// Delete hard-deletes a TenantInfra record by its primary key.
func (r *tenantInfraRepo) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&TenantInfra{})

	if result.Error != nil {
		r.logger.ErrorWithErr("tenantInfraRepo.Delete: failed to delete tenant infra", result.Error)
		return ConvertDBError(result.Error, "tenant_infra")
	}

	if result.RowsAffected == 0 {
		return NotFoundError("TenantInfra", "id")
	}

	return nil
}

// --- Internal mapper ---

func toTenantInfraDomain(m *TenantInfra) *domain.TenantInfra {
	if m == nil {
		return nil
	}
	return &domain.TenantInfra{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Schema:         m.Schema,
		Status:         domain.TenantInfraStatus(m.Status),
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
