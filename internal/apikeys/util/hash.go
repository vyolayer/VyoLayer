package util

import (
	"crypto/subtle"
	"encoding/base64"
	"os"

	"golang.org/x/crypto/argon2"
)

var defaultPepper = []byte(getPepper())

func HashSecret(secret string) string {
	hash := argon2.IDKey(
		[]byte(secret),
		defaultPepper,
		1,
		64*1024,
		4,
		32,
	)

	return base64.RawStdEncoding.EncodeToString(hash)
}

func VerifySecret(secret string, storedHash string) bool {
	hash := HashSecret(secret)

	return subtle.ConstantTimeCompare(
		[]byte(hash),
		[]byte(storedHash),
	) == 1
}

func getPepper() string {
	v := os.Getenv("API_KEY_PEPPER")
	if v == "" {
		return "change-this-production-pepper"
	}
	return v
}
