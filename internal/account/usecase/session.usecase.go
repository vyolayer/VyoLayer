package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
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
) (*accountV1.RefreshSessionResponse, error) {
	refreshTokenHash := domain.SessionTokenHash(refreshToken)
	s, err := uc.sessionRepo.FindByTokenHash(ctx, projectID, refreshTokenHash)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, ErrSessionNotFound
	}
	if s.IsExpired() {
		return nil, ErrSessionExpired
	}

	deviceInfo, dErr := ctxutil.ExtractDeviceInfo(ctx)
	if dErr != nil {
		return nil, dErr
	}

	if !s.VerifySameDevice(deviceInfo.UserAgent) {
		return nil, ErrInvalidSession
	}

	accessToken, tokenErr := uc.accountJWT.GenerateAccessToken(s.UserID, s.ProjectID)
	if tokenErr != nil {
		return nil, ErrJwtTokenGeneration
	}

	refreshToken, tokenErr = uc.accountJWT.GenerateRefreshToken()
	if tokenErr != nil {
		return nil, ErrJwtTokenGeneration
	}

	s.RotateToken(refreshToken)
	updateErr := uc.sessionRepo.RotateToken(ctx, s)
	if updateErr != nil {
		return nil, updateErr
	}

	return &accountV1.RefreshSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *sessionUsecase) ListSessions(
	ctx context.Context,
	projectID, userID uuid.UUID,
) (*accountV1.AllSessionsResponse, error) {
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
) error {
	s, err := uc.sessionRepo.FindByID(ctx, projectID, sessionID)
	if err != nil {
		return err
	}
	if s == nil {
		return ErrSessionNotFound
	}
	if s.UserID.String() != userID.String() {
		return ErrInvalidRefreshToken
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
) error {
	// delete all sessions
	deleteErr := uc.sessionRepo.DeleteAll(ctx, projectID, userID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
