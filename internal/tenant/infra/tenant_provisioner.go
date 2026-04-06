package infra

import (
	"context"

	"github.com/vyolayer/vyolayer/pkg/tenant"
)

type TenantProvisioner interface {
	Provision(ctx context.Context, schema string) error
}

type TenantProvisionerImpl struct {
	provisioner Provisioner
	migrator    Migrator
}

func NewTenantProvisioner(p Provisioner, m Migrator) TenantProvisioner {
	return &TenantProvisionerImpl{
		provisioner: p,
		migrator:    m,
	}
}

func (t *TenantProvisionerImpl) Provision(ctx context.Context, schema string) error {
	// validate
	if err := tenant.ValidateSchemaName(schema); err != nil {
		return err
	}

	// create schema
	if err := t.provisioner.EnsureSchema(ctx, schema); err != nil {
		return err
	}

	// run migrations
	if err := t.migrator.Run(ctx, schema); err != nil {
		return err
	}

	return nil
}
