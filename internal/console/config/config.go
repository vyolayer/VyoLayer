package config

import (
	"github.com/vyolayer/vyolayer/pkg/config"
	"github.com/vyolayer/vyolayer/pkg/postgres"
)

type ConsoleConfig struct {
	GRPC     config.GRPCConfig
	Database postgres.Config
}
