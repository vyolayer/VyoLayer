package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// ==========================================
// 1. PASSWORD HASHING (Slow & Secure)
// ==========================================

// HashPassword generates a bcrypt hash of the password.
// Use this before saving to the database.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a raw password with a stored hash.
// Returns true if they match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ==========================================
// 2. TOKEN HASHING (Fast & Efficient)
// ==========================================

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
