package service

import (
	"log"
	"time"
	"worklayer/internal/config"
	"worklayer/internal/domain"
	"worklayer/internal/utils/token"
)

type TokenService interface {
	GenerateAccessToken(user domain.User) (string, ServiceError)
	GenerateRefreshToken(userId string) (string, ServiceError)
	ValidateAccessToken(accessToken string) (*token.UserJwtDTO, ServiceError)
	ValidateRefreshToken(refreshToken string) (string, ServiceError)
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

func (ts *tokenService) GenerateAccessToken(user domain.User) (string, ServiceError) {
	accessToken, err := ts.tokenManager.GenerateAccessToken(token.UserJwtDTO{
		UserID: user.ID.InternalID().String(),
		Email:  user.Email,
	})
	if err != nil {
		log.Printf("TOKEN SERVICE :: GenerateAccessToken :: err : %v", err)
		return "", NewServiceError(500, err.Error())
	}
	return accessToken, nil
}

func (ts *tokenService) GenerateRefreshToken(userID string) (string, ServiceError) {
	refreshToken, err := ts.tokenManager.GenerateRefreshToken(token.UserJwtDTO{
		UserID: userID,
	})
	if err != nil {
		return "", NewServiceError(500, err.Error())
	}
	return refreshToken, nil
}

func (ts *tokenService) ValidateAccessToken(accessToken string) (*token.UserJwtDTO, ServiceError) {
	user, err := ts.tokenManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, NewServiceError(500, err.Error())
	}
	return &token.UserJwtDTO{
		UserID: user.UserID,
		Email:  user.Email,
	}, nil
}

func (ts *tokenService) ValidateRefreshToken(refreshToken string) (string, ServiceError) {
	user, err := ts.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", NewServiceError(500, err.Error())
	}

	return user.UserID, nil
}

func (ts *tokenService) GetAccessTokenExpiry() time.Duration {
	return ts.authConfig.AccessTokenTTL
}

func (ts *tokenService) GetRefreshTokenExpiry() time.Duration {
	return ts.authConfig.RefreshTokenTTL
}
