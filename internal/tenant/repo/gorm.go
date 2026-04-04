package tenantrepo

import (
	"context"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type gormRepo struct {
	db     *gorm.DB
	logger *logger.AppLogger
}

func NewGormRepo(db *gorm.DB, logger *logger.AppLogger) domain.GormRepository {
	return &gormRepo{
		db:     db,
		logger: logger,
	}
}

func (r *gormRepo) BeginTx(ctx context.Context) (*gorm.DB, error) {
	open := r.db.Begin().WithContext(ctx)

	return open, open.Error
}

func (r *gormRepo) CommitTx(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *gormRepo) RollbackTx(tx *gorm.DB) error {
	return tx.Rollback().Error
}
