package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	model "github.com/vyolayer/vyolayer/internal/iam/models/v1"
	"gorm.io/gorm"
)

type SessionRepository interface {
	// FindByTokenHash looks up a session by the raw session token (hashed internally).
	FindByTokenHash(ctx context.Context, rawToken string) (model.Session, error)
	// Save persists a new session for the given user.
	Save(ctx context.Context, userID uuid.UUID, rawToken, ip, userAgent string) error
	// Delete removes a single session row identified by its primary key.
	Delete(ctx context.Context, sessionID int64) error
	// DeleteAllByUserID revokes all sessions belonging to a user (e.g. global sign-out).
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	client        *gorm.DB
	SessionExpiry time.Duration
}

func NewSessionRepository(client *gorm.DB, sessionExpiry time.Duration) SessionRepository {
	return &sessionRepository{
		client:        client,
		SessionExpiry: sessionExpiry,
	}
}

// hashToken produces a SHA-256 hex digest of the raw token for safe storage.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// FindByTokenHash looks up a session using the SHA-256 hash of the provided raw token.
func (r *sessionRepository) FindByTokenHash(ctx context.Context, rawToken string) (model.Session, error) {
	var s model.Session
	err := r.client.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL", hashToken(rawToken)).
		First(&s).Error
	return s, err
}

// Save stores a new session with a hashed token.
func (r *sessionRepository) Save(ctx context.Context, userID uuid.UUID, rawToken, ip, userAgent string) error {
	s := model.Session{
		UserID:    userID,
		TokenHash: hashToken(rawToken),
		IpAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(r.SessionExpiry),
	}
	return r.client.WithContext(ctx).Create(&s).Error
}

// Delete removes a single session row by its primary key.
func (r *sessionRepository) Delete(ctx context.Context, sessionID int64) error {
	return r.client.WithContext(ctx).
		Where("id = ?", sessionID).
		Delete(&model.Session{}).Error
}

// DeleteAllByUserID hard-deletes (via soft-delete) all sessions for a user.
func (r *sessionRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.client.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Session{}).Error
}
