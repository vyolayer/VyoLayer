package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// UserJwtDTO is the input data
type UserJwtDTO struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

// AccessClaims: Rich data for the application (Avoids DB lookups)
type AccessClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RefreshClaims: Minimal data (Just enough to identify user)
type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenManager interface
type TokenManager interface {
	GenerateAccessToken(user UserJwtDTO) (string, error)
	GenerateRefreshToken(user UserJwtDTO) (string, error)

	// Returns specific structs for each type
	ValidateAccessToken(tokenStr string) (*AccessClaims, error)
	ValidateRefreshToken(tokenStr string) (*RefreshClaims, error)
}

type tokenConfig struct {
	secret []byte
	expiry time.Duration
}

type tokenManager struct {
	access  tokenConfig
	refresh tokenConfig
}
