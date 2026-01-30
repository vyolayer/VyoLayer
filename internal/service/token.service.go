package service

import (
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/config"
	"worklayer/internal/utils/response"
	"worklayer/internal/utils/token"
)

type TokenService interface {
	GenerateAccessToken(user *dto.UserDTO) (string, ServiceError)
	GenerateRefreshToken(user *dto.UserDTO) (string, ServiceError)
	ValidateAccessToken(accessToken string) (*token.UserJwtDTO, ServiceError)
	ValidateRefreshToken(refreshToken string) (uint, ServiceError)
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

func (ts *tokenService) GenerateAccessToken(user *dto.UserDTO) (string, ServiceError) {
	accessToken, err := ts.tokenManager.GenerateAccessToken(token.UserJwtDTO{
		UserID: user.ID,
		Email:  user.Email,
	})
	if err != nil {
		return "", NewServiceError(response.InternalServerError("Failed to generate access token"))
	}
	return accessToken, nil
}

func (ts *tokenService) GenerateRefreshToken(user *dto.UserDTO) (string, ServiceError) {
	refreshToken, err := ts.tokenManager.GenerateRefreshToken(token.UserJwtDTO{
		UserID: user.ID,
	})
	if err != nil {
		return "", NewServiceError(response.InternalServerError("Failed to generate refresh token"))
	}
	return refreshToken, nil
}

func (ts *tokenService) ValidateAccessToken(accessToken string) (*token.UserJwtDTO, ServiceError) {
	user, err := ts.tokenManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, NewServiceError(response.UnauthorizedError("Invalid access token"))
	}
	return &token.UserJwtDTO{
		UserID: user.UserID,
		Email:  user.Email,
	}, nil
}

func (ts *tokenService) ValidateRefreshToken(refreshToken string) (uint, ServiceError) {
	user, err := ts.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return 0, NewServiceError(response.UnauthorizedError("Invalid refresh token"))
	}
	return user.UserID, nil
}

func (ts *tokenService) GetAccessTokenExpiry() time.Duration {
	return ts.authConfig.AccessTokenTTL
}

func (ts *tokenService) GetRefreshTokenExpiry() time.Duration {
	return ts.authConfig.RefreshTokenTTL
}
