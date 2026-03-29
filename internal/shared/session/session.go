package session

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken creates a SHA256 hash of a string.
// Use this for Refresh Tokens or API Keys.
func HashToken(text string) string {
	hash := sha256.Sum256([]byte(text))
	// Critical Fix: Encode to Hex string, otherwise you get garbage characters
	return hex.EncodeToString(hash[:])
}

// CompareTokenHash checks if a raw token matches its SHA256 hash.
func CompareTokenHash(rawToken, hashedToken string) bool {
	return HashToken(rawToken) == hashedToken
}
