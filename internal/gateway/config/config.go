package config

import (
	"os"
	"time"
)

// Config holds the gateway configuration parameters
type Config struct {
	HTTPPort           string
	AccountServiceAddr string
	AccountJWT         AccountJWTConfig
}

type AccountJWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// Load reads environment variables and populates the Config struct
func Load() *Config {
	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		AccountServiceAddr: getEnv("ACCOUNT_SERVICE_ADDR", "localhost:50051"),
		AccountJWT: AccountJWTConfig{
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
