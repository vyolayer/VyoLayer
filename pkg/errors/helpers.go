package errors

import "fmt"

// Common error creation helpers for frequently used error types

// NotFound creates a 404 Not Found error
func NotFound(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrResourceNotFound, message)
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrAuthUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrAuthForbidden, message)
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrRequestInvalidBody, message)
}

// InvalidParams creates a 400 Bad Request error for invalid parameters
func InvalidParams(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrRequestInvalidParams, message)
}

// Conflict creates a 409 Conflict error
func Conflict(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrResourceConflict, message)
}

// Internal creates a 500 Internal Server Error
func Internal(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrInternalUnexpected, message)
}

// InternalWrap wraps an unexpected error as an Internal Server Error
func InternalWrap(err error, message string, args ...interface{}) *AppError {
	if err == nil {
		return nil
	}
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return Wrap(err, ErrInternalUnexpected, message)
}

// ValidationFailed creates a 422 Validation Error
func ValidationFailed(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrValidationFailed, message)
}

// ValidationFailedWithDetails creates a validation error with field details
func ValidationFailedWithDetails(message string, details interface{}) *AppError {
	err := NewWithMessage(ErrValidationFailed, message)
	err.WithMetadata("validation_errors", details)
	return err
}

// RequiredField creates a validation error for a required field
func RequiredField(field string) *AppError {
	return NewWithMessage(ErrValidationRequiredField, "Field '%s' is required", field).
		WithMetadata("field", field)
}

// InvalidFormat creates a validation error for invalid format
func InvalidFormat(field, expected string) *AppError {
	return NewWithMessage(ErrValidationInvalidFormat, "Field '%s' has invalid format, expected: %s", field, expected).
		WithMetadata("field", field).
		WithMetadata("expected_format", expected)
}

// TooManyRequests creates a 429 Rate Limit error
func TooManyRequests(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrRequestRateLimited, message)
}

// NotImplemented creates a 501 Not Implemented error
func NotImplemented(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrInternalNotImplemented, message)
}

// ServiceUnavailable creates a 503 Service Unavailable error
func ServiceUnavailable(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrExternalServiceUnavailable, message)
}

// Database error helpers

// DBNotFound creates a database not found error
func DBNotFound(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrDBRecordNotFound, message)
}

// DBDuplicateKey creates a duplicate key error
func DBDuplicateKey(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return NewWithMessage(ErrDBDuplicateKey, message)
}

// DBQueryFailed wraps a database query error
func DBQueryFailed(err error, query string) *AppError {
	if err == nil {
		return nil
	}
	return Wrap(err, ErrDBQueryFailed, "Database query failed").
		WithMetadata("query", query)
}

// User-specific error helpers

// UserNotFound creates a user not found error
func UserNotFound(userID string) *AppError {
	return NewWithMessage(ErrUserNotFound, "User with ID '%s' not found", userID).
		WithMetadata("user_id", userID)
}

// UserAlreadyExists creates a user already exists error
func UserAlreadyExists(identifier string) *AppError {
	return NewWithMessage(ErrUserAlreadyExists, "User with identifier '%s' already exists", identifier).
		WithMetadata("identifier", identifier)
}

// Auth-specific error helpers

// InvalidCredentials creates an invalid credentials error
func InvalidCredentials() *AppError {
	return New(ErrAuthInvalidCredentials)
}

// TokenExpired creates a token expired error
func TokenExpired() *AppError {
	return New(ErrAuthTokenExpired)
}

// TokenInvalid creates an invalid token error
func TokenInvalid(reason string) *AppError {
	return NewWithMessage(ErrAuthTokenInvalid, "Invalid token: %s", reason).
		WithMetadata("reason", reason)
}

// SessionExpired creates a session expired error
func SessionExpired() *AppError {
	return New(ErrAuthSessionExpired)
}

// SessionNotFound creates a session not found error
func SessionNotFound() *AppError {
	return New(ErrAuthSessionNotFound)
}

// AccountNotVerified creates an account not verified error
func AccountNotVerified() *AppError {
	return New(ErrAuthAccountNotVerified)
}
