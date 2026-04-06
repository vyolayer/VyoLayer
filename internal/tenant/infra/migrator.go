package infra

import (
	"context"
	"fmt"

	accountmodelv1 "github.com/vyolayer/vyolayer/internal/account/models/v1"
	"github.com/vyolayer/vyolayer/pkg/tenant"
	"gorm.io/gorm"
)

type Migrator interface {
	Run(ctx context.Context, schema string) error
}

type GormMigrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) Migrator {
	return &GormMigrator{db: db}
}

func (m *GormMigrator) Run(ctx context.Context, schema string) error {
	if err := tenant.ValidateSchemaName(schema); err != nil {
		return err
	}

	tx := m.db.WithContext(ctx).Exec(fmt.Sprintf(`SET search_path TO "%s"`, schema))
	if tx.Error != nil {
		return tx.Error
	}

	// add your tenant tables here
	return m.db.AutoMigrate(
		&accountmodelv1.ServiceUser{},
	// &Session{},
	// add more
	)
}
