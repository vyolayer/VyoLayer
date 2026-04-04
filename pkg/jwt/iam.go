package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type IamJWT interface {
	GenerateAccessToken(user *IAMUserJWTDto) (string, time.Time, error)
	VerifyAccessToken(token string) (*IAMUserJWTDto, error)
	GenerateRefreshToken() (string, time.Time, error)
	GetRefreshTokenExpiry() time.Time
	GetAccessTokenExpiry() time.Time
}

type iamJWT struct {
	accessTokenSecret  string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type IAMUserJWTDto struct {
	UserID          uuid.UUID `json:"user_id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Status          string    `json:"status"`
	IsEmailVerified bool      `json:"is_email_verified"`
	JoinedAt        time.Time `json:"joined_at"`
}

type iamClaims struct {
	IAMUserJWTDto
	jwt.RegisteredClaims
}

func NewIamJWT(accessTokenSecret string, accessTokenExpiry, refreshTokenExpiry time.Duration) IamJWT {
	return &iamJWT{
		accessTokenSecret:  accessTokenSecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (a *iamJWT) GenerateAccessToken(user *IAMUserJWTDto) (string, time.Time, error) {
	claims := iamClaims{
		IAMUserJWTDto: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "VyoLayer_IAM",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.accessTokenSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return token, claims.ExpiresAt.Time, nil
}

func (a *iamJWT) VerifyAccessToken(tokenStr string) (*IAMUserJWTDto, error) {
	if tokenStr == "" {
		return nil, ErrInvalidToken
	}

	t, err := jwt.ParseWithClaims(tokenStr, &iamClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.accessTokenSecret), nil
	})

	if err != nil {
		log.Printf("[JWT] Error verifying access token: %v", err)
		return nil, err
	}

	if claims, ok := t.Claims.(*iamClaims); ok && t.Valid {
		return &claims.IAMUserJWTDto, nil
	}

	log.Printf("[JWT] Token parsed but invalid or claims mismatch")
	return nil, ErrInvalidToken
}

func (a *iamJWT) GenerateRefreshToken() (string, time.Time, error) {
	var (
		b      []byte
		err    error
		token  string
		expiry time.Time
	)

	b = make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", time.Time{}, err
	}

	token = base64.URLEncoding.EncodeToString(b)
	expiry = time.Now().Add(a.refreshTokenExpiry)

	return token, expiry, nil
}

func (a *iamJWT) GetRefreshTokenExpiry() time.Time {
	return time.Now().Add(a.refreshTokenExpiry)
}

func (a *iamJWT) GetAccessTokenExpiry() time.Time {
	return time.Now().Add(a.accessTokenExpiry)
}
