package token

import (
	"time"
<<<<<<< HEAD
	"vyolayer/internal/config"
=======
>>>>>>> bc09bb9 (refactor: update module and struct names to vyolayer)

	"github.com/golang-jwt/jwt/v5"
	"github.com/vyolayer/vyolayer/internal/config"
)

func NewTokenManager(authConfig config.AuthConfig) TokenManager {
	return &tokenManager{
		access:  tokenConfig{secret: []byte(authConfig.JWTSecret), expiry: authConfig.AccessTokenTTL},
		refresh: tokenConfig{secret: []byte(authConfig.RefreshTokenSecret), expiry: authConfig.RefreshTokenTTL},
	}
}

// --- Generators ---
func (tm *tokenManager) GenerateAccessToken(user UserJwtDTO) (string, error) {
	claims := AccessClaims{
		UserID: user.UserID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.access.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "VyoLayer",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tm.access.secret)
}

func (tm *tokenManager) GenerateRefreshToken(user UserJwtDTO) (string, error) {
	claims := RefreshClaims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.refresh.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "VyoLayer",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tm.refresh.secret)
}

// --- Validators ---
func (tm *tokenManager) ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		return tm.access.secret, nil
	})

	if err != nil {
		return nil, err // Returns wrapped error (Expired/Invalid)
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func (tm *tokenManager) ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return tm.refresh.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}
