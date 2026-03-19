package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/config"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/mail"
)

func NewAccountRecoverUsecase(
	cfg *config.Config,
	ur repository.UserRepository,
	sr repository.SessionRepository,
	tr repository.VerificationTokenRepository,
	mailer mail.Mailer,
	jwt jwt.AccountJWT,
) AccountRecoverUsecase {
	return &accountRecoverUsecase{
		cfg:    cfg,
		ur:     ur,
		sr:     sr,
		tr:     tr,
		mailer: mailer,
		jwt:    jwt,
	}
}

func (uc *accountRecoverUsecase) ChangePassword(
	ctx context.Context,
	projectID, userID uuid.UUID,
	oldPassword, newPassword string,
) error {
	u, err := uc.ur.FindByID(ctx, projectID, userID)
	if err != nil {
		return err
	}

	if !u.VerifyPassword(oldPassword) {
		return ErrInvalidPassword
	}

	if u.IsSamePassword(newPassword) {
		return ErrSamePassword
	}

	e := u.ChangePassword(newPassword)
	if e != nil {
		return e
	}

	err = uc.ur.Update(ctx, projectID, u)
	if err != nil {
		return err
	}

	uc.sr.DeleteAll(ctx, projectID, userID)

	return nil
}

func (uc *accountRecoverUsecase) ForgotPassword(
	ctx context.Context,
	projectID uuid.UUID,
	email string,
) error {
	// Find user
	u, err := uc.ur.FindByEmail(ctx, projectID, email)
	if err != nil {
		return err
	}
	rawToken, hashToken, e := domain.GenerateVerificationToken()
	if e != nil {
		return ErrJwtTokenGeneration
	}
	s := domain.NewVerificationToken(
		projectID,
		u.ID,
		hashToken,
		domain.TokenTypePasswordReset,
	)

	resetLink := strings.Join([]string{
		uc.cfg.AppURL,
		"/reset-password?token=",
		rawToken,
	}, "")
	mailerErr := uc.mailer.Send(&mail.Message{
		To:      []string{u.Email},
		Subject: "Password Reset",
		Body:    "Please click the link below to reset your password: " + resetLink,
	})
	if mailerErr != nil {
		return ErrFailedToSendEmail
	}

	err = uc.tr.Create(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func (uc *accountRecoverUsecase) ResetPassword(
	ctx context.Context,
	projectID uuid.UUID,
	token, newPassword string,
) error {
	// find verification token
	t, err := uc.tr.FindByTokenHash(
		ctx,
		projectID,
		domain.HashToken(token),
	)
	if err != nil {
		return ErrInvalidResetPasswordToken
	}
	if t.IsExpired() || !t.IsPasswordResetToken() || t.IsUsed() {
		return ErrInvalidVerificationToken
	}

	// find user
	u, err := uc.ur.FindByID(ctx, projectID, t.UserID)
	if err != nil {
		return err
	}

	if u.IsSamePassword(newPassword) {
		return ErrSamePassword
	}

	// change password
	e := u.ChangePassword(newPassword)
	if e != nil {
		return e
	}

	err = uc.ur.Update(ctx, projectID, u)
	if err != nil {
		return err
	}

	uc.tr.Delete(ctx, projectID, u.ID)
	uc.sr.DeleteAll(ctx, projectID, u.ID)
	return nil

}
