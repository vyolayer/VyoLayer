package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/iam/domain"
	repo "github.com/vyolayer/vyolayer/internal/iam/repository"
	"github.com/vyolayer/vyolayer/internal/shared/session"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/mail"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrUserNotFound      = status.Error(codes.NotFound, "User not found")
	ErrUserAlreadyExists = status.Error(codes.AlreadyExists, "User already exists")
)

type AuthUsecase struct {
	log    *logger.AppLogger
	ur     repo.IAMUserRepository
	sr     repo.SessionRepository
	vtr    repo.VerificationTokenRepository
	prtr   repo.PasswordResetTokenRepository
	ss     *session.IAMSession
	mailer mail.Mailer
	appURL string
}

func NewAuthUsecase(
	log *logger.AppLogger,
	userRepo repo.IAMUserRepository,
	sessionRepo repo.SessionRepository,
	verificationTokenRepo repo.VerificationTokenRepository,
	passwordResetTokenRepo repo.PasswordResetTokenRepository,
	sessionService *session.IAMSession,
	mailer mail.Mailer,
	appURL string,
) *AuthUsecase {
	return &AuthUsecase{
		log:    log,
		ur:     userRepo,
		sr:     sessionRepo,
		vtr:    verificationTokenRepo,
		prtr:   passwordResetTokenRepo,
		ss:     sessionService,
		mailer: mailer,
		appURL: appURL,
	}
}

// ── Registration flow ──────────────────────────────────────────────────────────

func (uc *AuthUsecase) RegisterUser(ctx context.Context, req *iAMV1.RegisterRequest) (string, error) {
	if existing, _ := uc.ur.FindByEmail(ctx, req.GetEmail()); existing != nil {
		return "", ErrUserAlreadyExists
	}

	fallbackChar := ""
	if len(req.GetFullName()) > 0 {
		fallbackChar = req.GetFullName()[0:1]
	}
	avatar := domain.NewIAMUserAvatar("", fallbackChar)
	user := domain.NewIAMUser(req.GetEmail(), req.GetPassword(), req.GetFullName())
	user.InitAvatar(avatar)

	if err := uc.ur.Create(ctx, user); err != nil {
		return "", err
	}

	// Best-effort: log failures but don't roll back the registration.
	if err := uc.issueAndSendVerificationEmail(ctx, user.ID, user.GetEmail()); err != nil {
		uc.log.Error("(AuthUsecase.RegisterUser): send verification email: ", err.Error())
	}

	return user.ID.String(), nil
}

// VerifyEmail marks the user as email-verified using a secure single-use token.
func (uc *AuthUsecase) VerifyEmail(ctx context.Context, rawToken string) error {
	userID, err := uc.vtr.Consume(ctx, rawToken)
	if err != nil {
		return status.Error(codes.NotFound, "invalid or expired verification token")
	}

	user, err := uc.ur.FindByID(ctx, userID)
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	user.VerifyEmail()
	return uc.ur.Update(ctx, user)
}

// ResendVerificationEmail issues a new verification token and re-sends the email.
func (uc *AuthUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := uc.ur.FindByEmail(ctx, email)
	if err != nil {
		// Do NOT leak whether the email exists.
		return nil
	}

	if user.IsEmailVerified {
		return status.Error(codes.AlreadyExists, "email already verified")
	}

	if err := uc.issueAndSendVerificationEmail(ctx, user.ID, user.GetEmail()); err != nil {
		uc.log.Error("(AuthUsecase.ResendVerificationEmail): ", err.Error())
		// Surface gracefully — do not leak internal details to the caller.
	}

	return nil
}

// issueAndSendVerificationEmail generates a secure token, stores it, and emails the link.
func (uc *AuthUsecase) issueAndSendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error {
	rawToken, err := uc.vtr.Create(ctx, userID)
	if err != nil {
		return fmt.Errorf("create verification token: %w", err)
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", uc.appURL, rawToken)
	uc.log.Debug("Verification URL: ", map[string]any{"url": verifyURL})
	return uc.mailer.Send(&mail.Message{
		To:      []string{email},
		Subject: "Verify your email address",
		Body:    "Click the link to verify your account: " + verifyURL,
		IsHTML:  false,
	})
}

// ── Session flow ───────────────────────────────────────────────────────────────

func (uc *AuthUsecase) Login(ctx context.Context, req *iAMV1.LoginRequest) (*iAMV1.SessionTokenResponse, error) {
	user, err := uc.ur.FindByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if !user.Password.VerifyPassword(req.GetPassword()) {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	sess, err := uc.ss.CreateSession(&jwt.IAMUserJWTDto{
		UserID:          user.GetID(),
		FullName:        user.GetFullName(),
		Email:           user.GetEmail(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.IsEmailVerified,
		JoinedAt:        user.Timestamps.CreatedAt,
	})
	if err != nil {
		return nil, err
	}

	ip, ua := "-", "-"
	if dif, _ := ctxutil.ExtractDeviceInfo(ctx); dif != nil {
		ip = dif.IP
		ua = dif.UserAgent
	}

	if err := uc.sr.Save(ctx, user.ID, sess.SessionToken, ip, ua); err != nil {
		return nil, err
	}

	return &iAMV1.SessionTokenResponse{
		AccessToken:           sess.AccessToken,
		SessionToken:          sess.SessionToken,
		AccessTokenExpiresAt:  timestamppb.New(sess.AccessTokenExpiresAt),
		SessionTokenExpiresAt: timestamppb.New(sess.SessionTokenExpiresAt),
	}, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, req *iAMV1.RefreshSessionRequest) (*iAMV1.SessionTokenResponse, error) {
	sessionModel, err := uc.sr.FindByTokenHash(ctx, req.GetSessionToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired session token")
	}

	if time.Now().After(sessionModel.ExpiresAt) || sessionModel.RevokedAt != nil {
		return nil, status.Error(codes.Unauthenticated, "session expired or revoked")
	}
	user, err := uc.ur.FindByID(ctx, sessionModel.UserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	sess, err := uc.ss.CreateSession(&jwt.IAMUserJWTDto{
		UserID:          user.GetID(),
		FullName:        user.GetFullName(),
		Email:           user.GetEmail(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.IsEmailVerified,
		JoinedAt:        user.Timestamps.CreatedAt,
	})
	if err != nil {
		return nil, err
	}

	ip, ua := "-", "-"
	if dif, _ := ctxutil.ExtractDeviceInfo(ctx); dif != nil {
		ip = dif.IP
		ua = dif.UserAgent
	}

	if err := uc.sr.Save(ctx, sessionModel.UserID, sess.SessionToken, ip, ua); err != nil {
		return nil, err
	}

	return &iAMV1.SessionTokenResponse{
		AccessToken:           sess.AccessToken,
		SessionToken:          sess.SessionToken,
		AccessTokenExpiresAt:  timestamppb.New(sess.AccessTokenExpiresAt),
		SessionTokenExpiresAt: timestamppb.New(sess.SessionTokenExpiresAt),
	}, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, req *iAMV1.LogoutRequest) error {
	sessionModel, err := uc.sr.FindByTokenHash(ctx, req.GetSessionToken())
	if err != nil {
		return status.Error(codes.NotFound, "session not found")
	}

	if err := uc.sr.Delete(ctx, sessionModel.ID); err != nil {
		return status.Error(codes.Internal, "failed to delete session")
	}

	return nil
}

// ── Password flow ──────────────────────────────────────────────────────────────

// ChangePassword validates the old password and applies the new one.
func (uc *AuthUsecase) ChangePassword(ctx context.Context, req *iAMV1.ChangePasswordRequest) error {
	uid, err := extractCallerID(ctx)
	if err != nil {
		return err
	}

	user, err := uc.ur.FindByID(ctx, uid)
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	if !user.Password.VerifyPassword(req.GetOldPassword()) {
		return status.Error(codes.Unauthenticated, "incorrect current password")
	}

	if req.GetNewPassword() != req.GetConfirmPassword() {
		return status.Error(codes.InvalidArgument, "passwords do not match")
	}

	user.Password = domain.NewPassword(req.GetNewPassword())
	return uc.ur.Update(ctx, user)
}

// ForgotPassword generates a secure reset token, persists it, and emails a reset link.
func (uc *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := uc.ur.FindByEmail(ctx, email)
	if err != nil {
		// Do NOT leak whether the email exists.
		return nil
	}

	rawToken, err := uc.prtr.Create(ctx, user.ID)
	if err != nil {
		uc.log.Error("(AuthUsecase.ForgotPassword): create reset token: ", err)
		return nil
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", uc.appURL, rawToken)
	if err := uc.mailer.Send(&mail.Message{
		To:      []string{email},
		Subject: "Reset your password",
		Body:    "Click the link to reset your password (valid for 1 hour): " + resetURL,
		IsHTML:  false,
	}); err != nil {
		uc.log.Error("(AuthUsecase.ForgotPassword): send email: ", err)
	}

	return nil
}

// ResetPassword consumes the reset token and sets a new password.
func (uc *AuthUsecase) ResetPassword(ctx context.Context, req *iAMV1.ResetPasswordRequest) error {
	if req.GetNewPassword() != req.GetConfirmPassword() {
		return status.Error(codes.InvalidArgument, "passwords do not match")
	}

	userID, err := uc.prtr.Consume(ctx, req.GetToken())
	if err != nil {
		return status.Error(codes.NotFound, "invalid or expired reset token")
	}

	user, err := uc.ur.FindByID(ctx, userID)
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	user.Password = domain.NewPassword(req.GetNewPassword())
	return uc.ur.Update(ctx, user)
}
