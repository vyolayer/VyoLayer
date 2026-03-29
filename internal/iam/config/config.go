package config

import (
	"time"

	"github.com/vyolayer/vyolayer/pkg/config"
	"github.com/vyolayer/vyolayer/pkg/postgres"
)

type Config struct {
	GRPC     config.GRPCConfig
	Database postgres.Config
	Mail     config.MailConfig
	JWT      config.JWTConfig
	Token    TokenConfig
}

// TokenConfig holds expiry durations for one-time tokens sent via email.
type TokenConfig struct {
	VerificationTokenExpiry  time.Duration
	PasswordResetTokenExpiry time.Duration
}

func Load() *Config {
	return &Config{
		GRPC: config.GRPCConfig{
			GRPCPort: config.GetEnv("IAM_SERVICE_GRPC_PORT", "50052"),
			AppURL:   config.GetEnv("APP_URL", "http://localhost:3000"),
		},
		Database: *postgres.NewConfig(postgres.DefaultConfig),
		Mail:     *config.NewMailConfig(config.DefaultMailConfig),
		JWT: config.JWTConfig{
			AccessTokenSecret:  config.GetEnv("IAM_SERVICE_ACCESS_SECRET", "i_am_access_secret"),
			RefreshTokenSecret: config.GetEnv("IAM_SERVICE_REFRESH_SECRET", "i_am_refresh_secret"),
			AccessTokenExpiry:  config.GetEnvDuration("IAM_SERVICE_ACCESS_EXPIRY", "15m"),
			RefreshTokenExpiry: config.GetEnvDuration("IAM_SERVICE_REFRESH_EXPIRY", "2h"),
		},
		Token: TokenConfig{
			VerificationTokenExpiry:  config.GetEnvDuration("IAM_VERIFICATION_TOKEN_EXPIRY", "24h"),
			PasswordResetTokenExpiry: config.GetEnvDuration("IAM_PASSWORD_RESET_TOKEN_EXPIRY", "1h"),
		},
	}
}

