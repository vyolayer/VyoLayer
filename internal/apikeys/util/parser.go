package util

import (
	"errors"
	"strings"
)

var (
	ErrInvalidAPIKeyFormat = errors.New("invalid api key format")
)

func ExtractPrefix(secret string) string {
	parts := strings.Split(secret, "_")
	if len(parts) < 4 {
		return ""
	}

	return strings.Join(parts[:3], "_")
}

func ExtractEnvironment(secret string) string {
	parts := strings.Split(secret, "_")
	if len(parts) < 3 {
		return ""
	}

	switch parts[1] {
	case EnvDev:
		return EnvDev
	case EnvLive:
		return EnvLive
	default:
		return ""
	}
}

func ValidateAPIKeyFormat(secret string) error {
	parts := strings.Split(secret, "_")

	// vyo_live_xxxxx_secret
	if len(parts) < 4 {
		return ErrInvalidAPIKeyFormat
	}

	if parts[0] != "vyo" {
		return ErrInvalidAPIKeyFormat
	}

	if parts[1] != EnvDev && parts[1] != EnvLive {
		return ErrInvalidAPIKeyFormat
	}

	if parts[2] == "" || parts[3] == "" {
		return ErrInvalidAPIKeyFormat
	}

	return nil
}
