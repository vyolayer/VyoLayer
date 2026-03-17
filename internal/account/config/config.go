package config

import (
	"os"
	"time"

	"github.com/vyolayer/vyolayer/pkg/postgres"
)

// Config holds the account service configuration parameters
type Config struct {
	GRPCPort string
	Database postgres.Config
	AppURL   string
	JWT      JWTConfig
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// Load reads environment variables and populates the Config struct
func Load() *Config {
	return &Config{
		GRPCPort: getEnv("ACCOUNT_SERVICE_GRPC_PORT", "50051"),
		Database: postgres.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "4444"),
			User:     getEnv("DB_USER", "vyolayer_user"),
			Password: getEnv("DB_PASSWORD", "vyolayer_password"),
			DBName:   getEnv("DB_NAME", "vyolayer_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		AppURL: getEnv("APP_URL", "http://localhost:3000"),
		JWT: JWTConfig{
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "access_token_secret"),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "refresh_token_secret"),
			AccessTokenExpiry:  getEnvDuration("ACCESS_TOKEN_EXPIRY", "15m"),
			RefreshTokenExpiry: getEnvDuration("REFRESH_TOKEN_EXPIRY", "168h"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvDuration(key, fallback string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	d, _ := time.ParseDuration(fallback)
	return d
}
