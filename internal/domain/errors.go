package domain

import "worklayer/pkg/errors"

// Domain-specific error helpers
// These provide convenient constructors for domain layer errors

// User errors
var (
	ErrUserNotFound       = errors.ErrUserNotFound
	ErrUserAlreadyExists  = errors.ErrUserAlreadyExists
	ErrUserInactive       = errors.ErrUserInactive
	ErrInvalidCredentials = errors.ErrAuthInvalidCredentials
	ErrUserNotVerified    = errors.ErrAuthAccountNotVerified
	ErrSessionNotFound    = errors.ErrAuthSessionNotFound
	ErrSessionExpired     = errors.ErrAuthSessionExpired
	ErrTokenExpired       = errors.ErrAuthTokenExpired
	ErrTokenInvalid       = errors.ErrAuthTokenInvalid
	ErrPasswordHashFailed = errors.ErrInternalHashing
)

// UserNotFoundError creates a user not found error
func UserNotFoundError(userID string) *errors.AppError {
	return errors.UserNotFound(userID)
}

// UserAlreadyExistsError creates a user already exists error
func UserAlreadyExistsError(email string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrUserAlreadyExists, "User with email '%s' already exists", email).
		WithMetadata("email", email)
}

// InvalidCredentialsError creates an invalid credentials error
func InvalidCredentialsError() *errors.AppError {
	return errors.InvalidCredentials()
}

// UserNotVerifiedError creates a user not verified error
func UserNotVerifiedError() *errors.AppError {
	return errors.AccountNotVerified()
}

// SessionNotFoundError creates a session not found error
func SessionNotFoundError() *errors.AppError {
	return errors.SessionNotFound()
}

// SessionExpiredError creates a session expired error
func SessionExpiredError() *errors.AppError {
	return errors.SessionExpired()
}

// TokenExpiredError creates a token expired error
func TokenExpiredError() *errors.AppError {
	return errors.TokenExpired()
}

// TokenInvalidError creates an invalid token error
func TokenInvalidError(reason string) *errors.AppError {
	return errors.TokenInvalid(reason)
}

// PasswordHashFailedError creates a password hash failed error
func PasswordHashFailedError(err error) *errors.AppError {
	return errors.Wrap(err, errors.ErrInternalHashing, "Failed to hash password")
}

// ValidationError creates a validation error
func ValidationError(message string) *errors.AppError {
	return errors.ValidationFailed(message)
}

// InvalidEmailError creates an invalid email error
func InvalidEmailError(email string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrValidationInvalidEmail, "Invalid email: %s", email).
		WithMetadata("email", email)
}

// InvalidPasswordError creates an invalid password error
func InvalidPasswordError(reason string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrValidationInvalidPassword, "Invalid password: %s", reason).
		WithMetadata("reason", reason)
}
