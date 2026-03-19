package usecase

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// User
	ErrInvalidPassword = status.Error(codes.InvalidArgument, "invalid password")
	ErrInvalidEmail    = status.Error(codes.InvalidArgument, "invalid email")
	ErrInvalidUsername = status.Error(codes.InvalidArgument, "invalid username")
	ErrSamePassword    = status.Error(codes.InvalidArgument, "password is the same as before")

	ErrUserAlreadyVerified = status.Error(codes.FailedPrecondition, "user already verified")
	ErrUserNotVerified     = status.Error(codes.FailedPrecondition, "user not verified")
	ErrUserInactive        = status.Error(codes.FailedPrecondition, "user is inactive")
	ErrUserNotFound        = status.Error(codes.NotFound, "user not found")

	ErrEmailAlreadyExists    = status.Error(codes.AlreadyExists, "email already exists")
	ErrUsernameAlreadyExists = status.Error(codes.AlreadyExists, "username already exists")

	// Jwt
	ErrInvalidAccessToken  = status.Error(codes.InvalidArgument, "invalid access token")
	ErrInvalidRefreshToken = status.Error(codes.InvalidArgument, "invalid refresh token")
	ErrJwtTokenGeneration  = status.Error(codes.Internal, "failed to generate jwt token")

	// Session
	ErrSessionNotFound = status.Error(codes.NotFound, "session not found")
	ErrInvalidSession  = status.Error(codes.InvalidArgument, "invalid session")
	ErrSessionExpired  = status.Error(codes.FailedPrecondition, "session expired")

	// Verification
	ErrInvalidVerificationToken = status.Error(codes.InvalidArgument, "invalid verification token")
	ErrVerificationTokenExpired = status.Error(codes.FailedPrecondition, "verification token expired")
	ErrVerificationTokenUsed    = status.Error(codes.FailedPrecondition, "verification token already used")

	ErrInvalidResetPasswordToken = status.Error(codes.InvalidArgument, "invalid reset password token")
	ErrResetPasswordTokenExpired = status.Error(codes.FailedPrecondition, "reset password token expired")
	ErrResetPasswordTokenUsed    = status.Error(codes.FailedPrecondition, "reset password token already used")

	// Mail
	ErrFailedToSendEmail = status.Error(codes.Internal, "failed to send email")
)
