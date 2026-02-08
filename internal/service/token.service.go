package service

import (
	"log"
	"time"
	"worklayer/internal/config"
	"worklayer/internal/domain"
	"worklayer/internal/utils/token"
	"worklayer/pkg/errors"
)

type TokenService interface {
	GenerateAccessToken(user domain.User) (string, *errors.AppError)
	GenerateRefreshToken(userId string) (string, *errors.AppError)
	ValidateAccessToken(accessToken string) (*token.UserJwtDTO, *errors.AppError)
	ValidateRefreshToken(refreshToken string) (string, *errors.AppError)
	GetAccessTokenExpiry() time.Duration
	GetRefreshTokenExpiry() time.Duration
}

type tokenService struct {
	authConfig   config.AuthConfig
	tokenManager token.TokenManager
}

func NewTokenService(authConfig config.AuthConfig) TokenService {
	tokenManager := token.NewTokenManager(authConfig)
	return &tokenService{
		authConfig:   authConfig,
		tokenManager: tokenManager,
	}
}

func (ts *tokenService) GenerateAccessToken(user domain.User) (string, *errors.AppError) {
	accessToken, err := ts.tokenManager.GenerateAccessToken(token.UserJwtDTO{
		UserID: user.ID.InternalID().String(),
		Email:  user.Email,
	})
	if err != nil {
		log.Printf("TOKEN SERVICE :: GenerateAccessToken :: err : %v", err)
		return "", errors.Wrap(err, errors.ErrInternalUnexpected, "Failed to generate access token")
	}
	return accessToken, nil
}

func (ts *tokenService) GenerateRefreshToken(userID string) (string, *errors.AppError) {
	refreshToken, err := ts.tokenManager.GenerateRefreshToken(token.UserJwtDTO{
		UserID: userID,
	})
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternalUnexpected, "Failed to generate refresh token")
	}
	return refreshToken, nil
}

func (ts *tokenService) ValidateAccessToken(accessToken string) (*token.UserJwtDTO, *errors.AppError) {
	user, err := ts.tokenManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, domain.TokenInvalidError(err.Error())
	}
	return &token.UserJwtDTO{
		UserID: user.UserID,
		Email:  user.Email,
	}, nil
}

func (ts *tokenService) ValidateRefreshToken(refreshToken string) (string, *errors.AppError) {
	user, err := ts.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", domain.TokenInvalidError(err.Error())
	}

	return user.UserID, nil
}

func (ts *tokenService) GetAccessTokenExpiry() time.Duration {
	return ts.authConfig.AccessTokenTTL
}

func (ts *tokenService) GetRefreshTokenExpiry() time.Duration {
	return ts.authConfig.RefreshTokenTTL
}
