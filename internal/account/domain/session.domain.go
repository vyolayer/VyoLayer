package domain

import (
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	// TODO: 15 days
	SessionTTL = time.Hour * 24 * 15
)

type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ProjectID     uuid.UUID
	TokenHash     string
	IPAddress     string
	UserAgent     string
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	RevokedReason string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSession(
	userID uuid.UUID,
	projectID uuid.UUID,
	token, ipAddress, userAgent string,
) *Session {
	tokenHash := sessionTokenHash(token)
	log.Println("Session created for user: ", userID, tokenHash)

	return &Session{
		ID:            uuid.New(),
		UserID:        userID,
		ProjectID:     projectID,
		TokenHash:     tokenHash,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		RevokedReason: "",
		RevokedAt:     nil,
		ExpiresAt:     time.Now().Add(SessionTTL),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// verify same device
func (s *Session) VerifySameDevice(userAgent string) bool {
	return s.UserAgent == userAgent
}

// Rotate token
func (s *Session) RotateToken(newToken string) {
	s.TokenHash = sessionTokenHash(newToken)
	s.ExpiresAt = time.Now().Add(SessionTTL)
	s.UpdatedAt = time.Now()
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired()
}

func (s *Session) IsRevoked() bool {
	return false
}

func (s *Session) Revoke(reason string) {
	now := time.Now()
	s.RevokedAt = &now
	s.RevokedReason = reason
}

func sessionTokenHash(token string) string {
	return HashToken(token)
}

func SessionTokenHash(token string) string {
	return sessionTokenHash(token)
}
