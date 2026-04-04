package session

import (
	"log"
	"time"

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

func (s *IAMSession) CreateSession(dto *jwt.IAMUserJWTDto) (*IAMSessionResponse, error) {
	var (
		accessToken        string
		accessTokenExpiry  time.Time
		refreshToken       string
		refreshTokenExpiry time.Time
		err                error
	)

	accessToken, accessTokenExpiry, err = s.jwt.GenerateAccessToken(dto)
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

func (s *IAMSession) VerifyAccessToken(token string) (*jwt.IAMUserJWTDto, error) {
	return s.jwt.VerifyAccessToken(token)
}
