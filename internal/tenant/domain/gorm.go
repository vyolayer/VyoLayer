package domain

import (
	"context"

	"gorm.io/gorm"
)

type GormRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB) error
}
