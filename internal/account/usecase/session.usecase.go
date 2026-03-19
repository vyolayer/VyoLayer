package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

func NewSessionUsecase(
	sessionRepo repository.SessionRepository,
	accountJWT jwt.AccountJWT,
) SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
		accountJWT:  accountJWT,
	}
}

func (uc *sessionUsecase) RefreshToken(
	ctx context.Context,
	projectID uuid.UUID,
	refreshToken string,
) (*accountV1.RefreshSessionResponse, *errors.AppError) {
	refreshTokenHash := domain.SessionTokenHash(refreshToken)
	s, err := uc.sessionRepo.FindByTokenHash(ctx, projectID, refreshTokenHash)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, errors.BadRequest("Invalid refresh token")
	}
	if s.IsExpired() {
		return nil, errors.BadRequest("Refresh token expired")
	}

	deviceInfo, dErr := ctxutil.ExtractDeviceInfo(ctx)
	if dErr != nil {
		return nil, errors.Internal("Failed to extract device info", dErr)
	}

	if !s.VerifySameDevice(deviceInfo.IP, deviceInfo.UserAgent) {
		return nil, errors.BadRequest("Invalid refresh token")
	}

	accessToken, tokenErr := uc.accountJWT.GenerateAccessToken(s.UserID, s.ProjectID)
	if tokenErr != nil {
		return nil, errors.Internal("Failed to generate access token", tokenErr)
	}

	refreshToken, tokenErr = uc.accountJWT.GenerateRefreshToken()
	if tokenErr != nil {
		return nil, errors.Internal("Failed to generate refresh token", tokenErr)
	}

	s.RotateToken(refreshToken)
	updateErr := uc.sessionRepo.RotateToken(ctx, s)
	if updateErr != nil {
		return nil, errors.Internal("Failed to update session", updateErr)
	}

	return &accountV1.RefreshSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *sessionUsecase) ListSessions(
	ctx context.Context,
	projectID, userID uuid.UUID,
) (*accountV1.AllSessionsResponse, *errors.AppError) {
	sessions, err := uc.sessionRepo.FindByUserID(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	ss := make([]*accountV1.SessionData, 0, len(sessions))
	for _, s := range sessions {
		ss = append(ss, &accountV1.SessionData{
			SessionId: s.ID.String(),
			UserId:    s.UserID.String(),
			IpAddress: s.IPAddress,
			UserAgent: s.UserAgent,
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
			UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &accountV1.AllSessionsResponse{Sessions: ss}, nil
}

func (uc *sessionUsecase) RevokeSession(
	ctx context.Context,
	projectID, userID, sessionID uuid.UUID,
) *errors.AppError {
	s, err := uc.sessionRepo.FindByID(ctx, projectID, sessionID)
	if err != nil {
		return err
	}
	if s == nil {
		return errors.BadRequest("Invalid session")
	}
	if s.UserID.String() != userID.String() {
		return errors.BadRequest("Invalid session")
	}

	// delete session
	deleteErr := uc.sessionRepo.Delete(ctx, projectID, sessionID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (uc *sessionUsecase) RevokeAllSessions(
	ctx context.Context,
	projectID, userID uuid.UUID,
) *errors.AppError {
	// delete all sessions
	deleteErr := uc.sessionRepo.DeleteAll(ctx, projectID, userID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
