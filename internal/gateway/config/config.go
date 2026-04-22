package config

import (
	"os"
	"time"
)

// Config holds the gateway configuration parameters
type Config struct {
	HTTPPort           string
	AccountServiceAddr string
	IAMServiceAddr     string
	TenantServiceAddr  string
	ConsoleServiceAddr string
	APIKeyServiceAddr  string
	AccountJWT         AccountJWTConfig
	IAMJWT             IAMJWTConfig
}

type AccountJWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type IAMJWTConfig struct {
	AccessTokenSecret  string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// Load reads environment variables and populates the Config struct
func Load() *Config {
	return &Config{
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		AccountServiceAddr: getEnv("ACCOUNT_SERVICE_ADDR", "localhost:50051"),
		IAMServiceAddr:     getEnv("APP_SERVICE_ADDR", "localhost:50052"),
		TenantServiceAddr:  getEnv("TENANT_SERVICE_ADDR", "localhost:50053"),
		ConsoleServiceAddr: getEnv("CONSOLE_SERVICE_ADDR", "localhost:50054"),
		APIKeyServiceAddr:  getEnv("API_KEY_SERVICE_ADDR", "localhost:50055"),
		AccountJWT: AccountJWTConfig{
			AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", "access_token_secret"),
			RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "refresh_token_secret"),
			AccessTokenExpiry:  getEnvDuration("ACCESS_TOKEN_EXPIRY", "15m"),
			RefreshTokenExpiry: getEnvDuration("REFRESH_TOKEN_EXPIRY", "168h"),
		},
		IAMJWT: IAMJWTConfig{
			AccessTokenSecret:  getEnv("IAM_ACCESS_TOKEN_SECRET", "i_am_access_secret"),
			AccessTokenExpiry:  getEnvDuration("IAM_ACCESS_TOKEN_EXPIRY", "15m"),
			RefreshTokenExpiry: getEnvDuration("IAM_REFRESH_TOKEN_EXPIRY", "168h"),
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
