package repository

import (
	"database/sql"
	"strings"

	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

// ConvertDBError converts database errors to AppError
func ConvertDBError(err error, context string) *errors.AppError {
	if err == nil {
		return nil
	}

	// Check for GORM specific errors
	if err == gorm.ErrRecordNotFound {
		return errors.DBNotFound("Record not found: %s", context)
	}

	// Check for SQL errors
	if err == sql.ErrNoRows {
		return errors.DBNotFound("No rows found: %s", context)
	}

	// Check for duplicate key errors
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "UNIQUE constraint") ||
		strings.Contains(errMsg, "Duplicate entry") {
		return errors.DBDuplicateKey("Duplicate record: %s", context)
	}

	// Check for constraint violations
	if strings.Contains(errMsg, "constraint") ||
		strings.Contains(errMsg, "foreign key") {
		return errors.NewWithMessage(errors.ErrDBConstraintViolation, "Database constraint violation: %s", context).
			WithMetadata("context", context)
	}

	// Default to query failed
	return errors.Wrap(err, errors.ErrDBQueryFailed, "Database operation failed: %s", context)
}

// NotFoundError creates a repository not found error
func NotFoundError(resourceType, identifier string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrDBRecordNotFound, "%s with identifier '%s' not found", resourceType, identifier).
		WithMetadata("resource_type", resourceType).
		WithMetadata("identifier", identifier)
}

// DuplicateError creates a repository duplicate error
func DuplicateError(resourceType, field, value string) *errors.AppError {
	return errors.NewWithMessage(errors.ErrDBDuplicateKey, "%s with %s '%s' already exists", resourceType, field, value).
		WithMetadata("resource_type", resourceType).
		WithMetadata("field", field).
		WithMetadata("value", value)
}

// QueryError wraps a query error
func QueryError(err error, query string) *errors.AppError {
	return errors.DBQueryFailed(err, query)
}

// TransactionError creates a transaction error
func TransactionError(err error, operation string) *errors.AppError {
	return errors.Wrap(err, errors.ErrDBTransactionFailed, "Transaction failed: %s", operation).
		WithMetadata("operation", operation)
}
