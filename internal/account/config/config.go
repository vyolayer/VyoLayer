package config

import (
	"os"
	"strconv"
	"time"

	"github.com/vyolayer/vyolayer/pkg/postgres"
)

// Config holds the account service configuration parameters
type Config struct {
	GRPCPort string
	Database postgres.Config
	AppURL   string
	Mail     MailConfig
	JWT      JWTConfig
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseMock  bool
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
		Mail: MailConfig{
			Host:     getEnv("MAIL_HOST", "localhost"),
			Port:     getEnvInt("MAIL_PORT", "1025"),
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", "noreply@vyolayer.local"),
			UseMock:  getEnvBool("MAIL_USE_MOCK", "true"),
		},
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

func getEnvInt(key, fallback string) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	i, _ := strconv.Atoi(fallback)
	return i
}

func getEnvBool(key, fallback string) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
			return true
		case "0", "false", "FALSE", "False", "no", "NO", "No":
			return false
		}
	}
	switch fallback {
	case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
		return true
	default:
		return false
	}
}
