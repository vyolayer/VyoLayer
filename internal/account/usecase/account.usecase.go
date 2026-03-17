package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/mail"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

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

func NewAccountUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, tokenRepo repository.VerificationTokenRepository, mailer mail.Mailer, accountJWT jwt.AccountJWT) AccountUsecase {
	return &accountUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenRepo:   tokenRepo,
		mailer:      mailer,
		accountJWT:  accountJWT,
	}
}

func (uc *accountUsecase) Register(
	ctx context.Context,
	projectID uuid.UUID,
	email, username, password, firstName, lastName string,
) (string, *errors.AppError) {
	var (
		user *domain.User
		err  *errors.AppError
	)

	// Check email and username
	if username != "" {
		user, err = uc.userRepo.FindByUsername(ctx, projectID, username)
		if user != nil {
			return "", errors.BadRequest("Username already exists")
		}
	}

	user, err = uc.userRepo.FindByEmail(ctx, projectID, email)
	if user != nil {
		return "", errors.BadRequest("Email already exists")
	}

	user = domain.NewUser(projectID, email, username, password, firstName, lastName)
	user.InitAvatar()

	err = uc.userRepo.Create(ctx, projectID, user)
	if err != nil {
		return "", err
	}

	rawToken, tokenHash, tokenErr := domain.GenerateVerificationToken()
	if tokenErr != nil {
		return "", errors.Internal("Failed to generate verification token", tokenErr)
	}

	tokenRecord := domain.NewVerificationToken(
		projectID,
		user.ID,
		tokenHash,
		domain.TokenTypeEmailVerify,
	)
	err = uc.tokenRepo.Create(ctx, tokenRecord)
	if err != nil {
		return "", err
	}

	uc.mailer.Send(&mail.Message{
		To:      []string{email},
		Subject: "Verify your email address",
		Body:    "Please click on the link below to verify your email address: " + rawToken,
		IsHTML:  false,
	})

	return user.ID.String(), nil
}

func (uc *accountUsecase) VerifyEmail(
	ctx context.Context,
	projectID uuid.UUID,
	token string,
) *errors.AppError {
	// Verify token
	tokenHash := domain.HashToken(token)
	tr, err := uc.tokenRepo.FindByTokenHash(ctx, projectID, tokenHash)
	if err != nil {
		return err
	}
	if err := tr.Validate(); err != nil {
		return err
	}

	// Update user
	user, err := uc.userRepo.FindByID(ctx, projectID, tr.UserID)
	if err != nil {
		return err
	}
	if user.IsVerified() {
		return errors.BadRequest("Email already verified")
	}

	user.VerifyEmail()

	err = uc.userRepo.Update(ctx, projectID, user)
	if err != nil {
		return err
	}

	tr.Use()
	err = uc.tokenRepo.Update(ctx, projectID, tr)
	if err != nil {
		return err
	}

	return nil
}

func (uc *accountUsecase) ResendVerificationEmail(
	ctx context.Context,
	projectID uuid.UUID,
	email string,
) *errors.AppError {
	// Find user
	user, err := uc.userRepo.FindByEmail(ctx, projectID, email)
	if err != nil {
		return err
	}

	if user.IsVerified() {
		return errors.BadRequest("Email already verified")
	}

	// Generate token
	rawToken, tokenHash, tokenErr := domain.GenerateVerificationToken()
	if tokenErr != nil {
		return errors.Internal("Failed to generate verification token", tokenErr)
	}

	// Create token record
	tokenRecord := domain.NewVerificationToken(
		projectID,
		user.ID,
		tokenHash,
		domain.TokenTypeEmailVerify,
	)
	err = uc.tokenRepo.Create(ctx, tokenRecord)
	if err != nil {
		return err
	}

	// Send email
	uc.mailer.Send(&mail.Message{
		To:      []string{email},
		Subject: "Verify your email address",
		Body:    "Please click on the link below to verify your email address: " + rawToken,
		IsHTML:  false,
	})

	return nil
}

func (uc *accountUsecase) Login(
	ctx context.Context,
	projectID uuid.UUID,
	email, password string,
) (*accountV1.LoginResponse, *errors.AppError) {
	var (
		user *domain.User
		err  *errors.AppError
	)

	user, err = uc.userRepo.FindByEmail(ctx, projectID, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.BadRequest("Invalid credentials")
	}
	if !user.IsVerified() {
		return nil, errors.BadRequest("Email not verified")
	}
	if !user.IsActive() {
		return nil, errors.BadRequest("Account not active")
	}
	if !user.VerifyPassword(password) {
		return nil, errors.BadRequest("Invalid credentials")
	}

	accessToken, tokenErr := uc.accountJWT.GenerateAccessToken(user.ID, projectID)
	if tokenErr != nil {
		return nil, errors.Internal("Failed to generate access token", tokenErr)
	}

	refreshToken, tokenErr := uc.accountJWT.GenerateRefreshToken()
	if tokenErr != nil {
		return nil, errors.Internal("Failed to generate refresh token", tokenErr)
	}

	deviceInfo, _ := ctxutil.ExtractDeviceInfo(ctx)
	if deviceInfo == nil {
		deviceInfo = &ctxutil.Device{
			IP:        "-",
			UserAgent: "-",
		}
	}

	// store session
	session := domain.NewSession(
		user.ID,
		projectID,
		refreshToken,
		deviceInfo.IP,
		deviceInfo.UserAgent,
	)
	err = uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &accountV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *accountUsecase) Logout(
	ctx context.Context,
	projectID uuid.UUID,
	userID uuid.UUID,
	refreshToken string,
) *errors.AppError {
	tokenHash := domain.SessionTokenHash(refreshToken)
	session, err := uc.sessionRepo.FindByTokenHash(ctx, projectID, tokenHash)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.BadRequest("Invalid refresh token")
	}
	if session.UserID != userID {
		return errors.BadRequest("Invalid refresh token")
	}
	if session.IsRevoked() {
		return errors.BadRequest("Invalid refresh token")
	}
	if session.IsExpired() {
		return errors.BadRequest("Invalid refresh token")
	}

	session.Revoke("User logged out")
	err = uc.sessionRepo.Delete(ctx, projectID, session.ID)
	if err != nil {
		return err
	}

	return nil
}
