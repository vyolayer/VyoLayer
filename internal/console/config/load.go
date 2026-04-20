package config

import (
	"github.com/vyolayer/vyolayer/pkg/config"
	"github.com/vyolayer/vyolayer/pkg/postgres"
)

func ConsoleConfigLoad() *ConsoleConfig {
	return &ConsoleConfig{
		GRPC: config.GRPCConfig{
			GRPCPort: config.GetEnv("CONSOLE_SERVICE_GRPC_PORT", "50054"),
			AppURL:   config.GetEnv("APP_URL", "http://localhost:3000"),
		},
		Database: postgres.Config{
			Host:     config.GetEnv("DB_HOST", "localhost"),
			Port:     config.GetEnv("DB_PORT", "4444"),
			User:     config.GetEnv("DB_USER", "vyolayer_user"),
			Password: config.GetEnv("DB_PASSWORD", "vyolayer_password"),
			DBName:   config.GetEnv("DB_NAME", "vyolayer_db"),
			SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}
