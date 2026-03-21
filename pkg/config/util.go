package config

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetEnvDuration(key, fallback string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	d, _ := time.ParseDuration(fallback)
	return d
}

func GetEnvInt(key, fallback string) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	i, _ := strconv.Atoi(fallback)
	return i
}

func GetEnvBool(key, fallback string) bool {
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
