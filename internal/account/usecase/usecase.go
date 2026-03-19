package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/config"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/mail"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

// --- Account Usecase
type AccountUsecase interface {
	Register(ctx context.Context, projectID uuid.UUID, email, username, password, firstName, lastName string) (string, *errors.AppError)
	VerifyEmail(ctx context.Context, projectID uuid.UUID, token string) *errors.AppError
	ResendVerificationEmail(ctx context.Context, projectID uuid.UUID, email string) *errors.AppError
	Login(ctx context.Context, projectID uuid.UUID, email, password string) (*accountV1.LoginResponse, *errors.AppError)
	Logout(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, refreshToken string) *errors.AppError
}

type accountUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tokenRepo   repository.VerificationTokenRepository
	mailer      mail.Mailer
	accountJWT  jwt.AccountJWT
}

// --- Session Usecase
type SessionUsecase interface {
	RefreshToken(
		ctx context.Context,
		projectID uuid.UUID,
		refreshToken string,
	) (*accountV1.RefreshSessionResponse, *errors.AppError)

	ListSessions(
		ctx context.Context,
		projectID, userID uuid.UUID,
	) (*accountV1.AllSessionsResponse, *errors.AppError)

	RevokeSession(
		ctx context.Context,
		projectID, userID, sessionID uuid.UUID,
	) *errors.AppError

	RevokeAllSessions(
		ctx context.Context,
		projectID, userID uuid.UUID,
	) *errors.AppError
}

type sessionUsecase struct {
	sessionRepo repository.SessionRepository
	accountJWT  jwt.AccountJWT
}

// --- Account recover Usecase
type AccountRecoverUsecase interface {
	// ChangePassword updates a password for an already authenticated user.
	// The userID should be extracted from the context (via your interceptor).
	ChangePassword(
		ctx context.Context,
		projectID uuid.UUID,
		userID uuid.UUID,
		oldPassword, newPassword string,
	) *errors.AppError

	// ForgotPassword initiates the recovery flow by generating a token
	// and potentially sending an email via a provider.
	ForgotPassword(
		ctx context.Context,
		projectID uuid.UUID,
		email string,
	) *errors.AppError

	// ResetPassword completes the recovery flow using the token
	// received in the forgot password step.
	ResetPassword(
		ctx context.Context,
		projectID uuid.UUID,
		token, newPassword string,
	) *errors.AppError
}

type accountRecoverUsecase struct {
	cfg    *config.Config
	ur     repository.UserRepository
	sr     repository.SessionRepository
	tr     repository.VerificationTokenRepository
	mailer mail.Mailer
	jwt    jwt.AccountJWT
}
