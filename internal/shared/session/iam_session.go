package session

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/pkg/jwt"
)

type IAMSessionResponse struct {
	AccessToken           string
	SessionToken          string
	AccessTokenExpiresAt  time.Time
	SessionTokenExpiresAt time.Time
}

// JWTSession for access token
type IAMSession struct {
	jwt jwt.IamJWT
}

func NewIAMSession(
	accessTokenSecret string,
	accessTokenExpiry, refreshTokenExpiry time.Duration,
) *IAMSession {
	jwt := jwt.NewIamJWT(accessTokenSecret, accessTokenExpiry, refreshTokenExpiry)

	log.Println("[IAM] Session created with JWT: ", jwt)

	return &IAMSession{
		jwt: jwt,
	}
}

func (s *IAMSession) CreateSession(userID uuid.UUID) (*IAMSessionResponse, error) {
	var (
		accessToken        string
		accessTokenExpiry  time.Time
		refreshToken       string
		refreshTokenExpiry time.Time
		err                error
	)

	accessToken, accessTokenExpiry, err = s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiry, err = s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	res := &IAMSessionResponse{
		AccessToken:           accessToken,
		SessionToken:          refreshToken,
		AccessTokenExpiresAt:  accessTokenExpiry,
		SessionTokenExpiresAt: refreshTokenExpiry,
	}

	return res, nil

}

func (s *IAMSession) VerifyAccessToken(token string) (uuid.UUID, error) {
	return s.jwt.VerifyAccessToken(token)
}
