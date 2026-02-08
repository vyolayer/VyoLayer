package service

import "worklayer/pkg/errors"

// Service-specific error helpers
// These provide convenient constructors for service layer errors

// WrapRepositoryError wraps a repository error for the service layer
func WrapRepositoryError(err error, context string) *errors.AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, just return it
	if appErr, ok := errors.As(err); ok {
		return appErr
	}

	// Otherwise wrap it
	return errors.InternalWrap(err, "Service error: %s", context)
}

// BusinessRuleViolation creates a business rule violation error
func BusinessRuleViolation(message string, args ...interface{}) *errors.AppError {
	return errors.NewWithMessage(errors.ErrBusinessRuleViolation, message, args...)
}

// InvalidStateError creates an invalid state error
func InvalidStateError(message string, args ...interface{}) *errors.AppError {
	return errors.NewWithMessage(errors.ErrBusinessInvalidState, message, args...)
}

// OperationNotAllowedError creates an operation not allowed error
func OperationNotAllowedError(message string, args ...interface{}) *errors.AppError {
	return errors.NewWithMessage(errors.ErrBusinessOperationNotAllowed, message, args...)
}

// ExternalServiceError creates an external service error
func ExternalServiceError(service string, err error) *errors.AppError {
	return errors.Wrap(err, errors.ErrExternalServiceError, "External service '%s' failed", service).
		WithMetadata("service", service)
}

// ConfigurationError creates a configuration error
func ConfigurationError(message string, args ...interface{}) *errors.AppError {
	return errors.NewWithMessage(errors.ErrInternalConfiguration, message, args...)
}
