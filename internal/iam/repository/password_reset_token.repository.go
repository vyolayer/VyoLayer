package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	model "github.com/vyolayer/vyolayer/internal/iam/models/v1"
	"gorm.io/gorm"
)

// PasswordResetTokenRepository manages one-time password-reset tokens.
type PasswordResetTokenRepository interface {
	// Create generates a secure random token for the user, persists its hash,
	// and returns the raw (unhashed) token to be emailed to the user.
	Create(ctx context.Context, userID uuid.UUID) (rawToken string, err error)
	// Consume validates the token (checks hash, expiry, single-use), marks it used,
	// and returns the associated user ID.
	Consume(ctx context.Context, rawToken string) (uuid.UUID, error)
	// DeleteByUserID removes all unused reset tokens for a user.
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type passwordResetTokenRepository struct {
	client *gorm.DB
	expiry time.Duration
}

func NewPasswordResetTokenRepository(client *gorm.DB, expiry time.Duration) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{client: client, expiry: expiry}
}

// hashResetToken hashes a raw token with SHA-256. Reuses the same helper as
// the verification token repository (both are in the same package).
func hashResetToken(raw string) string {
	return hashVerificationToken(raw) // shared SHA-256 helper
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, userID uuid.UUID) (string, error) {
	raw, err := generateRawToken()
	if err != nil {
		return "", err
	}

	// Remove any existing unused tokens before creating a fresh one.
	_ = r.DeleteByUserID(ctx, userID)

	rec := model.PasswordResetToken{
		UserID:    userID,
		TokenHash: hashResetToken(raw),
		ExpiresAt: time.Now().Add(r.expiry),
	}
	if err := r.client.WithContext(ctx).Create(&rec).Error; err != nil {
		return "", ConvertDBError(err, "create password reset token")
	}
	return raw, nil
}

func (r *passwordResetTokenRepository) Consume(ctx context.Context, rawToken string) (uuid.UUID, error) {
	var rec model.PasswordResetToken
	err := r.client.WithContext(ctx).
		Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", hashResetToken(rawToken), time.Now()).
		First(&rec).Error
	if err != nil {
		return uuid.Nil, ConvertDBError(err, "password reset token not found or expired")
	}

	now := time.Now()
	if err := r.client.WithContext(ctx).
		Model(&rec).
		Update("used_at", now).Error; err != nil {
		return uuid.Nil, ConvertDBError(err, "mark password reset token used")
	}

	return rec.UserID, nil
}

func (r *passwordResetTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.client.WithContext(ctx).
		Where("user_id = ? AND used_at IS NULL", userID).
		Delete(&model.PasswordResetToken{}).Error
}
