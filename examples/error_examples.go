package main

import (
	"fmt"
	"worklayer/pkg/errors"
	"worklayer/pkg/logger"
)

// This file demonstrates various error creation patterns

func demonstrateBasicErrors() {
	fmt.Println("\n=== Basic Error Creation ===")

	// Simple errors using helpers
	err1 := errors.NotFound("User not found")
	fmt.Printf("Error 1: %v (HTTP: %d)\n", err1, err1.HTTPStatus)

	err2 := errors.Unauthorized("Invalid access token")
	fmt.Printf("Error 2: %v (HTTP: %d)\n", err2, err2.HTTPStatus)

	err3 := errors.BadRequest("Invalid email format")
	fmt.Printf("Error 3: %v (HTTP: %d)\n", err3, err3.HTTPStatus)

	// Errors with formatted messages
	userID := "user123"
	err4 := errors.NotFound("User with ID %s not found", userID)
	fmt.Printf("Error 4: %v\n", err4)

	// Domain-specific errors
	err5 := errors.UserNotFound("user456")
	fmt.Printf("Error 5: %v (Code: %s)\n", err5, err5.Code)

	err6 := errors.InvalidCredentials()
	fmt.Printf("Error 6: %v (Code: %s)\n", err6, err6.Code)
}

func demonstrateErrorWithMetadata() {
	fmt.Println("\n=== Errors with Metadata ===")

	err := errors.NotFound("Order not found")
	err.WithMetadata("order_id", "ord_123456")
	err.WithMetadata("user_id", "user_789")
	err.WithMetadata("store_id", "store_001")

	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Metadata: %+v\n", err.Metadata)

	// Multiple metadata at once
	err2 := errors.Conflict("Email already registered")
	err2.WithMetadataMap(map[string]interface{}{
		"email":      "test@example.com",
		"ip_address": "192.168.1.1",
		"timestamp":  "2026-02-08T18:30:00Z",
	})
	fmt.Printf("Error 2 Metadata: %+v\n", err2.Metadata)
}

func demonstrateErrorBuilder() {
	fmt.Println("\n=== Error Builder Pattern ===")

	err := errors.NewBuilder(errors.ErrUserNotFound).
		WithMessage("User with email %s not found", "john@example.com").
		WithMetadata("email", "john@example.com").
		WithMetadata("search_type", "email").
		WithSeverity(errors.SeverityWarning).
		Build()

	fmt.Printf("Built Error: %v\n", err)
	fmt.Printf("Severity: %s\n", err.Severity)
	fmt.Printf("Metadata: %+v\n", err.Metadata)
}

func demonstrateErrorWrapping() {
	fmt.Println("\n=== Error Wrapping ===")

	// Simulate a database error
	dbErr := fmt.Errorf("connection timeout")

	// Wrap the error
	wrappedErr := errors.Wrap(dbErr, errors.ErrDBQueryFailed, "Failed to fetch user data")
	fmt.Printf("Wrapped Error: %v\n", wrappedErr)
	fmt.Printf("Underlying Error: %v\n", wrappedErr.Wrapped)

	// Wrap with context
	wrappedErr2 := errors.WrapWithContext(
		dbErr,
		errors.ErrDBConnectionFailed,
		"Database connection lost",
		map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"database": "worklayer",
		},
	)
	fmt.Printf("\nWrapped with Context: %v\n", wrappedErr2)
	fmt.Printf("Metadata: %+v\n", wrappedErr2.Metadata)
}

func demonstrateValidationErrors() {
	fmt.Println("\n=== Validation Errors ===")

	// Required field
	err1 := errors.RequiredField("email")
	fmt.Printf("Required Field Error: %v\n", err1)

	// Invalid format
	err2 := errors.InvalidFormat("phone", "E.164 format (+1234567890)")
	fmt.Printf("Invalid Format Error: %v\n", err2)

	// Validation with details
	validationDetails := []map[string]interface{}{
		{
			"field":   "email",
			"message": "Invalid email format",
		},
		{
			"field":   "password",
			"message": "Password must be at least 8 characters",
		},
	}
	err3 := errors.ValidationFailedWithDetails("Multiple validation errors", validationDetails)
	fmt.Printf("Validation Error: %v\n", err3)
	fmt.Printf("Details: %+v\n", err3.Metadata["validation_errors"])
}

func demonstrateDatabaseErrors() {
	fmt.Println("\n=== Database Errors ===")

	// Not found
	err1 := errors.DBNotFound("User with email 'test@example.com' not found")
	fmt.Printf("DB Not Found: %v (HTTP: %d)\n", err1, err1.HTTPStatus)

	// Duplicate key
	err2 := errors.DBDuplicateKey("Email 'test@example.com' already exists")
	fmt.Printf("DB Duplicate: %v (HTTP: %d)\n", err2, err2.HTTPStatus)

	// Query failed
	dbErr := fmt.Errorf("syntax error in SQL query")
	err3 := errors.DBQueryFailed(dbErr, "SELECT * FROM users WHERE id = ?")
	fmt.Printf("DB Query Failed: %v\n", err3)
}

func demonstrateErrorChecking() {
	fmt.Println("\n=== Error Checking ===")

	err := errors.UserNotFound("user123")

	// Check error code
	if errors.Is(err, errors.ErrUserNotFound) {
		fmt.Println("✓ Error is ErrUserNotFound")
	}

	// Get error code
	code := errors.GetCode(err)
	fmt.Printf("Error Code: %s\n", code)

	// Get HTTP status
	status := errors.GetHTTPStatus(err)
	fmt.Printf("HTTP Status: %d\n", status)

	// Get severity
	severity := errors.GetSeverity(err)
	fmt.Printf("Severity: %s\n", severity)

	// Convert to AppError
	if appErr, ok := errors.As(err); ok {
		fmt.Printf("✓ Converted to AppError: %s\n", appErr.Code)
	}
}

func demonstrateLogging() {
	fmt.Println("\n=== Error Logging ===")

	// Initialize logger in development mode
	logger.InitLogger(true)

	// Create and log errors
	err1 := errors.UserNotFound("user123")
	err1.WithMetadata("ip_address", "192.168.1.1")
	logger.LogError(err1, "req-123456")

	// Log info
	logger.LogInfo("User logged in successfully", map[string]interface{}{
		"user_id": "user123",
		"email":   "test@example.com",
	})

	// Log warning
	logger.LogWarning("Rate limit approaching", map[string]interface{}{
		"user_id":  "user123",
		"requests": 95,
		"limit":    100,
	})
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("ERROR MANAGEMENT SYSTEM - EXAMPLES")
	fmt.Println("==========================================")

	demonstrateBasicErrors()
	demonstrateErrorWithMetadata()
	demonstrateErrorBuilder()
	demonstrateErrorWrapping()
	demonstrateValidationErrors()
	demonstrateDatabaseErrors()
	demonstrateErrorChecking()
	demonstrateLogging()

	fmt.Println("\n==========================================")
	fmt.Println("Examples completed!")
	fmt.Println("==========================================")
}
