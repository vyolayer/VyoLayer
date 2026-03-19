package repository

import (
	"database/sql"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ConvertDBError converts database errors to AppError
func ConvertDBError(err error, context string) error {
	// Check for GORM specific errors
	if err == gorm.ErrRecordNotFound {
		return status.Error(codes.NotFound, "Record not found: "+context)
	}

	// Check for SQL errors
	if err == sql.ErrNoRows {
		return status.Error(codes.NotFound, "No rows found: "+context)
	}

	// Check for duplicate key errors
	errMsg := err.Error()
	if strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "UNIQUE constraint") ||
		strings.Contains(errMsg, "Duplicate entry") {
		// return errors.DBDuplicateKey("Duplicate record: %s", context)
		return status.Error(codes.AlreadyExists, "Duplicate record: "+context)
	}

	// Check for constraint violations
	if strings.Contains(errMsg, "constraint") ||
		strings.Contains(errMsg, "foreign key") {
		// return errors.NewWithMessage(errors.ErrDBConstraintViolation, "Database constraint violation: %s", context).
		// WithMetadata("context", context)
		return status.Error(codes.InvalidArgument, "Database constraint violation: "+context)
	}

	// Default to query failed
	// return errors.Wrap(err, errors.ErrDBQueryFailed, "Database operation failed: %s", context)
	return status.Error(codes.Internal, "Database operation failed: "+context)
}

// NotFoundError creates a repository not found error
func NotFoundError(resourceType, identifier string) error {
	// return errors.NewWithMessage(errors.ErrDBRecordNotFound, "%s with identifier '%s' not found", resourceType, identifier).
	// 	WithMetadata("resource_type", resourceType).
	// 	WithMetadata("identifier", identifier)
	return status.Error(codes.NotFound, resourceType+" with identifier "+identifier+" not found")
}

// DuplicateError creates a repository duplicate error
func DuplicateError(resourceType, field, value string) error {
	// return errors.NewWithMessage(errors.ErrDBDuplicateKey, "%s with %s '%s' already exists", resourceType, field, value).
	// 	WithMetadata("resource_type", resourceType).
	// 	WithMetadata("field", field).
	// 	WithMetadata("value", value)
	return status.Error(codes.AlreadyExists, resourceType+" with "+field+" "+value+" already exists")
}

// QueryError wraps a query error
func QueryError(err error, query string) error {
	// return errors.DBQueryFailed(err, query)
	return status.Error(codes.Internal, "Database operation failed: "+query)
}

// TransactionError creates a transaction error
func TransactionError(err error, operation string) error {
	// return errors.Wrap(err, errors.ErrDBTransactionFailed, "Transaction failed: %s", operation).
	// 	WithMetadata("operation", operation)
	return status.Error(codes.Internal, "Transaction failed: "+operation)
}
