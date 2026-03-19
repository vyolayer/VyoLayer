package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TokenTypeEmailVerify   = "email_verify"
	TokenTypePasswordReset = "password_reset"

	VerificationTokenTTL = time.Minute * 15
)

var (
	ErrVerificationTokenInvalid  = status.Error(codes.InvalidArgument, "Invalid verification token")
	ErrVerificationTokenNotFound = status.Error(codes.NotFound, "Verification token not found")

	ErrVerificationTokenExpired = status.Error(codes.FailedPrecondition, "Verification token expired")
	ErrVerificationTokenUsed    = status.Error(codes.FailedPrecondition, "Verification token already used")
)

type VerificationToken struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	Type      string
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    *time.Time
}

func NewVerificationToken(
	projectID, userID uuid.UUID,
	tokenHash, tokenType string,
) *VerificationToken {
	return &VerificationToken{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    userID,
		TokenHash: tokenHash,
		Type:      tokenType,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(VerificationTokenTTL),
	}
}

func (t *VerificationToken) Validate() error {
	if t == nil {
		return ErrVerificationTokenInvalid
	}
	if t.IsExpired() {
		return ErrVerificationTokenExpired
	}
	if t.IsUsed() {
		return ErrVerificationTokenUsed
	}
	return nil
}

func (t *VerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *VerificationToken) IsUsed() bool {
	return t.UsedAt != nil
}

func (t *VerificationToken) Use() {
	now := time.Now()
	t.UsedAt = &now
}

func (t *VerificationToken) IsEmailVerificationToken() bool {
	return t.Type == TokenTypeEmailVerify
}

func (t *VerificationToken) IsPasswordResetToken() bool {
	return t.Type == TokenTypePasswordReset
}

// GenerateRandomToken generates a cryptographically secure random token string
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates a SHA-256 hash of the given token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GenerateVerificationToken generates a new random token and its corresponding hash.
// The raw token should be sent to the user (e.g., via email link),
// and the hash should be used to create the VerificationToken record.
func GenerateVerificationToken() (string, string, error) {
	// Generate a 32-byte random token (becomes a 64-character hex string)
	rawToken, err := GenerateRandomToken(32)
	if err != nil {
		return "", "", err
	}

	tokenHash := HashToken(rawToken)
	return rawToken, tokenHash, nil
}
