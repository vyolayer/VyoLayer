package infra

import (
	"context"
	"fmt"

	"github.com/vyolayer/vyolayer/pkg/tenant"
	"gorm.io/gorm"
)

type Provisioner interface {
	EnsureSchema(ctx context.Context, schema string) error
}

type PostgresProvisioner struct {
	db *gorm.DB
}

func NewPostgresProvisioner(db *gorm.DB) Provisioner {
	return &PostgresProvisioner{db: db}
}

func (p *PostgresProvisioner) EnsureSchema(ctx context.Context, schema string) error {
	// Validate (prevents injection)
	if err := tenant.ValidateSchemaName(schema); err != nil {
		return err
	}

	// Create schema safely
	query := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, schema)

	return p.db.WithContext(ctx).Exec(query).Error
}
