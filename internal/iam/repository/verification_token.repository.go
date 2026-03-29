package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	model "github.com/vyolayer/vyolayer/internal/iam/models/v1"
	"gorm.io/gorm"
)

// VerificationTokenRepository manages one-time email-verification tokens.
type VerificationTokenRepository interface {
	// Create generates a secure random token for the user, persists its hash,
	// and returns the raw (unhashed) token to be sent via email.
	Create(ctx context.Context, userID uuid.UUID) (rawToken string, err error)
	// Consume finds the token by its hash, validates expiry, marks it used,
	// and returns the associated user ID. Returns an error if invalid/expired.
	Consume(ctx context.Context, rawToken string) (uuid.UUID, error)
	// DeleteByUserID removes all unused tokens for a user (used before re-issuing).
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type verificationTokenRepository struct {
	client  *gorm.DB
	expiry  time.Duration
}

func NewVerificationTokenRepository(client *gorm.DB, expiry time.Duration) VerificationTokenRepository {
	return &verificationTokenRepository{client: client, expiry: expiry}
}

// generateRawToken returns a cryptographically random 32-byte hex string.
func generateRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashVerificationToken hashes a raw token with SHA-256 for storage.
func hashVerificationToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (r *verificationTokenRepository) Create(ctx context.Context, userID uuid.UUID) (string, error) {
	raw, err := generateRawToken()
	if err != nil {
		return "", err
	}

	// Remove any prior unused tokens for this user before creating a new one.
	_ = r.DeleteByUserID(ctx, userID)

	rec := model.VerificationToken{
		UserID:    userID,
		TokenHash: hashVerificationToken(raw),
		ExpiresAt: time.Now().Add(r.expiry),
	}
	if err := r.client.WithContext(ctx).Create(&rec).Error; err != nil {
		return "", ConvertDBError(err, "create verification token")
	}
	return raw, nil
}

func (r *verificationTokenRepository) Consume(ctx context.Context, rawToken string) (uuid.UUID, error) {
	var rec model.VerificationToken
	err := r.client.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", hashVerificationToken(rawToken), time.Now()).
		First(&rec).Error
	if err != nil {
		return uuid.Nil, ConvertDBError(err, "verification token not found or expired")
	}

	now := time.Now()
	if err := r.client.WithContext(ctx).
		Model(&rec).
		Update("used_at", now).Error; err != nil {
		return uuid.Nil, ConvertDBError(err, "mark verification token used")
	}

	return rec.UserID, nil
}

func (r *verificationTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.client.WithContext(ctx).
		Where("user_id = ? AND used_at IS NULL", userID).
		Delete(&model.VerificationToken{}).Error
}
