package domain

import "log"

// ErrorType is a type that represents the type of error.
type ErrorType string

const (
	ErrorTypeDomain     ErrorType = "domain"
	ErrorTypeRepository ErrorType = "repository"
	ErrorTypeService    ErrorType = "service"
	ErrorTypeController ErrorType = "controller"
	ErrorTypeUnknown    ErrorType = "unknown"
)

// DomainError is a struct that represents an error.
type AppError struct {
	Code      int
	Message   string
	ErrorType *ErrorType
}

// Error returns the error message.
func (e *AppError) Error() string {
	return e.Message
}

// IsErrorType checks if the error is of the given type.
func (e *AppError) IsErrorType(errorType ErrorType) bool {
	return *e.ErrorType == errorType
}

// NewError creates a new error.
func NewAppError(code int, message string, errorType ErrorType) *AppError {
	err := &AppError{Code: code, Message: message, ErrorType: &errorType}
	log.Printf("[%s] %d - %s", *err.ErrorType, err.Code, err.Message)
	return err
}

// NewDomainError creates a new domain error.
func NewDomainError(code int, message string) *AppError {
	return NewAppError(code, message, ErrorTypeDomain)
}

// NewRepositoryError creates a new repository error.
func NewRepositoryError(code int, message string) *AppError {
	return NewAppError(code, message, ErrorTypeRepository)
}

// NewServiceError creates a new service error.
func NewServiceError(code int, message string) *AppError {
	return NewAppError(code, message, ErrorTypeService)
}

// NewControllerError creates a new controller error.
func NewControllerError(code int, message string) *AppError {
	return NewAppError(code, message, ErrorTypeController)
}

// NewUnknownError creates a new unknown error.
func NewUnknownError(code int, message string) *AppError {
	return NewAppError(code, message, ErrorTypeUnknown)
}

// Common errors for domain error
var (
	ErrUserNotFound       *AppError = NewDomainError(404, "user not found")
	ErrUserNotVerified    *AppError = NewDomainError(401, "user is not verified")
	ErrInvalidCredentials *AppError = NewDomainError(401, "invalid credentials")
	ErrSessionNotFound    *AppError = NewDomainError(404, "session not found")
	ErrSessionExpired     *AppError = NewDomainError(401, "session has expired")
	ErrTokenExpired       *AppError = NewDomainError(401, "token has expired")
	ErrTokenInvalid       *AppError = NewDomainError(401, "token is invalid")
	ErrUserAlreadyExists  *AppError = NewDomainError(409, "user already exists")
	ErrInternal           *AppError = NewDomainError(500, "internal error")
	ErrPasswordHashFailed *AppError = NewDomainError(500, "password hash failed")
)

type DomainError *AppError
