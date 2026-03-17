package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

type AccountJWT interface {
	GenerateAccessToken(userID uuid.UUID, projectID uuid.UUID) (string, error)
	VerifyAccessToken(token string) (uuid.UUID, uuid.UUID, error)
	GenerateRefreshToken() (string, error)
}

type accountJWT struct {
	accessTokenSecret  string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type accountClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
	jwt.RegisteredClaims
}

func NewAccountJWT(accessTokenSecret string, accessTokenExpiry time.Duration) AccountJWT {
	return &accountJWT{
		accessTokenSecret: accessTokenSecret,
		accessTokenExpiry: accessTokenExpiry,
	}
}

func (a *accountJWT) GenerateAccessToken(userID uuid.UUID, projectID uuid.UUID) (string, error) {
	claims := accountClaims{
		UserID:    userID,
		ProjectID: projectID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "VyoLayer",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.accessTokenSecret))
}

func (a *accountJWT) VerifyAccessToken(tokenStr string) (uuid.UUID, uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accountClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(a.accessTokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	if claims, ok := token.Claims.(*accountClaims); ok && token.Valid {
		return claims.UserID, claims.ProjectID, nil
	}
	return uuid.Nil, uuid.Nil, ErrInvalidToken
}

func (a *accountJWT) GenerateRefreshToken() (string, error) {
	// Create a 32-byte array (256 bits of entropy is highly secure)
	b := make([]byte, 32)

	// Fill the array with cryptographically secure random bytes
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to a URL-safe base64 string without padding
	// This ensures it is safe to pass in HTTP headers or URLs if needed
	return base64.RawURLEncoding.EncodeToString(b), nil
}
