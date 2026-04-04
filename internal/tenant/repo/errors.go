package tenantrepo

import (
	"database/sql"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func ConvertDBError(err error, context string) error {
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
		return status.Error(codes.AlreadyExists, "Duplicate record: "+context)
	}

	// Check for constraint violations
	if strings.Contains(errMsg, "constraint") ||
		strings.Contains(errMsg, "foreign key") {
		return status.Error(codes.InvalidArgument, "Database constraint violation: "+context)
	}

	// Default to query failed
	return status.Error(codes.Internal, "Database operation failed: "+context)
}

// NotFoundError creates a repository not found error
func NotFoundError(resourceType, identifier string) error {
	return status.Error(codes.NotFound, resourceType+" with identifier "+identifier+" not found")
}

// DuplicateError creates a repository duplicate error
func DuplicateError(resourceType, field, value string) error {
	return status.Error(codes.AlreadyExists, resourceType+" with "+field+" "+value+" already exists")
}

// QueryError wraps a query error
func QueryError(err error, query string) error {
	return status.Error(codes.Internal, "Database operation failed: "+query)
}

// TransactionError creates a transaction error
func TransactionError(err error, operation string) error {
	return status.Error(codes.Internal, "Transaction failed: "+operation)
}
