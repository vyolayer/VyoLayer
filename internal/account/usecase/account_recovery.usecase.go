package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/config"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/pkg/errors"
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
) *errors.AppError {
	u, err := uc.ur.FindByID(ctx, projectID, userID)
	if err != nil {
		return err
	}

	if !u.VerifyPassword(oldPassword) {
		return errors.New(errors.ErrAuthPasswordMismatch)
	}

	if u.IsSamePassword(newPassword) {
		return errors.New(errors.ErrAuthSamePassword)
	}

	e := u.ChangePassword(newPassword)
	if e != nil {
		return errors.Wrap(e, errors.ErrInternalHashing, "Failed to hash password")
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
) *errors.AppError {
	// Find user
	u, err := uc.ur.FindByEmail(ctx, projectID, email)
	if err != nil {
		return err
	}
	rawToken, hashToken, e := domain.GenerateVerificationToken()
	if e != nil {
		return errors.Internal("Failed to generate verification token", e)
	}
	s := domain.NewVerificationToken(
		projectID,
		u.ID,
		hashToken,
		domain.TokenTypePasswordReset,
	)

	resetLink := strings.Join([]string{
		uc.cfg.AppURL,
		"reset-password?token=",
		rawToken,
	}, "")
	mailerErr := uc.mailer.Send(&mail.Message{
		To:      []string{u.Email},
		Subject: "Password Reset",
		Body:    "Please click the link below to reset your password: " + resetLink,
	})
	if mailerErr != nil {
		return errors.Internal("Failed to send password reset email", mailerErr)
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
) *errors.AppError {
	// find verification token
	t, err := uc.tr.FindByTokenHash(
		ctx,
		projectID,
		domain.HashToken(token),
	)
	if err != nil {
		return err
	}
	if t.IsExpired() || !t.IsPasswordResetToken() {
		return errors.BadRequest("Token is expired or invalid")
	}

	// find user
	u, err := uc.ur.FindByID(ctx, projectID, t.UserID)
	if err != nil {
		return err
	}

	// change password
	e := u.ChangePassword(newPassword)
	if e != nil {
		return errors.Wrap(e, errors.ErrInternalHashing, "Failed to hash password")
	}

	err = uc.ur.Update(ctx, projectID, u)
	if err != nil {
		return err
	}

	uc.tr.Delete(ctx, projectID, u.ID)
	uc.sr.DeleteAll(ctx, projectID, u.ID)
	return nil

}
