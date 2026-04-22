package util

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

const (
	EnvDev  = "dev"
	EnvLive = "live"
)

type GeneratedKey struct {
	Secret string
	Prefix string
}

func GenerateAPIKey(environment string) (*GeneratedKey, error) {
	env := normalizeEnvironment(environment)

	prefixPart, err := randomToken(5)
	if err != nil {
		return nil, err
	}

	secretPart, err := randomToken(24)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("vyo_%s_%s", env, prefixPart)
	secret := fmt.Sprintf("%s_%s", prefix, secretPart)

	return &GeneratedKey{
		Secret: secret,
		Prefix: prefix,
	}, nil
}

func normalizeEnvironment(env string) string {
	switch strings.ToLower(env) {
	case EnvLive:
		return EnvLive
	default:
		return EnvDev
	}
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base32.StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(b)

	return strings.ToLower(token), nil
}
