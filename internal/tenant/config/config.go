package config

import (
	"github.com/vyolayer/vyolayer/pkg/config"
	"github.com/vyolayer/vyolayer/pkg/postgres"
)

type Config struct {
	GRPC     config.GRPCConfig
	Database postgres.Config
}

func Load() *Config {
	return &Config{
		GRPC: config.GRPCConfig{
			GRPCPort: config.GetEnv("TENANT_SERVICE_GRPC_PORT", "50053"),
			AppURL:   config.GetEnv("APP_URL", "http://localhost:3000"),
		},
		Database: *postgres.NewConfig(postgres.DefaultConfig),
	}
}
