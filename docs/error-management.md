# Error Management System - Usage Guide

## Overview

This error management system provides a comprehensive, production-grade solution for handling errors in the WorkLayer application. It includes:

- **Centralized error codes** with consistent naming
- **Rich error context** including stack traces, metadata, and error chaining
- **Standardized HTTP responses** with proper status codes
- **Structured logging** with severity levels
- **Easy-to-use builders and helpers**

## Table of Contents

1. [Quick Start](#quick-start)
2. [Error Codes](#error-codes)
3. [Creating Errors](#creating-errors)
4. [HTTP Responses](#http-responses)
5. [Middleware](#middleware)
6. [Logging](#logging)
7. [Migration Guide](#migration-guide)
8. [Best Practices](#best-practices)

## Quick Start

### Creating Simple Errors

```go
import "worklayer/pkg/errors"

// Use helper functions for common errors
err := errors.NotFound("User with ID %s not found", userID)
err := errors.Unauthorized("Invalid credentials")
err := errors.BadRequest("Invalid email format")
err := errors.Internal("Failed to process request")
```

### Returning Errors from Controllers

```go
import (
    "worklayer/pkg/errors"
    "worklayer/pkg/response"
)

func (c *Controller) GetUser(ctx *fiber.Ctx) error {
    userID := ctx.Params("id")

    user, err := c.service.GetUser(userID)
    if err != nil {
        // response.Error automatically handles AppError and sends proper HTTP response
        return response.Error(ctx, err)
    }

    return response.Success(ctx, user)
}
```

### Wrapping Errors

```go
// Wrap database errors
user, err := repo.FindByID(userID)
if err != nil {
    return errors.Wrap(err, errors.ErrDBQueryFailed, "Failed to find user")
}

// Wrap with context
err = service.ProcessPayment(payment)
if err != nil {
    return errors.WrapWithContext(
        err,
        errors.ErrExternalServiceError,
        "Payment processing failed",
        map[string]interface{}{
            "payment_id": payment.ID,
            "amount": payment.Amount,
        },
    )
}
```

## Error Codes

All error codes follow the format: `ERR_<CATEGORY>_<SPECIFIC_ERROR>`

### Categories

- **AUTH**: Authentication and authorization errors
- **VALIDATION**: Input validation errors
- **DB**: Database and repository errors
- **RESOURCE**: Resource management errors
- **USER**: User-specific errors
- **REQUEST**: HTTP request errors
- **EXTERNAL**: External service errors
- **INTERNAL**: Internal application errors
- **BUSINESS**: Business logic errors

### Common Error Codes

```go
// Authentication
errors.ErrAuthInvalidCredentials     // Invalid credentials provided
errors.ErrAuthUnauthorized           // Unauthorized access
errors.ErrAuthTokenExpired           // Token has expired
errors.ErrAuthSessionExpired         // Session has expired

// Validation
errors.ErrValidationFailed           // Validation failed
errors.ErrValidationRequiredField    // Required field missing
errors.ErrValidationInvalidEmail     // Invalid email format

// Database
errors.ErrDBRecordNotFound           // Record not found
errors.ErrDBDuplicateKey             // Duplicate key violation
errors.ErrDBQueryFailed              // Query execution failed

// User
errors.ErrUserNotFound               // User not found
errors.ErrUserAlreadyExists          // User already exists

// Request
errors.ErrRequestInvalidBody         // Invalid request body
errors.ErrRequestInvalidParams       // Invalid request parameters
errors.ErrRequestRateLimited         // Rate limit exceeded
```

## Creating Errors

### Method 1: Helper Functions (Recommended)

```go
// Simple errors
err := errors.NotFound("Resource not found")
err := errors.Unauthorized("Access denied")
err := errors.BadRequest("Invalid input")

// Errors with formatting
err := errors.NotFound("User with ID %s not found", userID)
err := errors.Conflict("Email %s is already registered", email)

// Domain-specific helpers
err := errors.UserNotFound(userID)
err := errors.InvalidCredentials()
err := errors.TokenExpired()
err := errors.SessionNotFound()

// Database helpers
err := errors.DBNotFound("User not found")
err := errors.DBDuplicateKey("Email already exists")

// Validation helpers
err := errors.ValidationFailed("Invalid input data")
err := errors.RequiredField("email")
err := errors.InvalidFormat("phone", "E.164 format")
```

### Method 2: Error Builder (Advanced)

```go
err := errors.NewBuilder(errors.ErrUserNotFound).
    WithMessage("User with ID %s not found", userID).
    WithMetadata("user_id", userID).
    WithMetadata("search_criteria", criteria).
    WithSeverity(errors.SeverityWarning).
    Build()
```

### Method 3: Direct Construction

```go
// Create error with code
err := errors.New(errors.ErrAuthUnauthorized)

// Create with custom message
err := errors.NewWithMessage(errors.ErrUserNotFound, "User %s not found", userID)
```

### Adding Metadata

```go
err := errors.NotFound("User not found")
err.WithMetadata("user_id", userID)
err.WithMetadata("search_type", "email")

// Add multiple metadata at once
err.WithMetadataMap(map[string]interface{}{
    "user_id": userID,
    "email": email,
    "ip_address": ipAddr,
})
```

## HTTP Responses

### Success Responses

```go
import "worklayer/pkg/response"

// Simple success
return response.Success(ctx, data)

// Success with custom message
return response.SuccessWithMessage(ctx, fiber.StatusOK, "User created", user)

// Success message only (no data)
return response.SuccessMessage(ctx, "Operation completed")

// Created (201)
return response.Created(ctx, newUser)

// No Content (204)
return response.NoContent(ctx)

// Paginated response
return response.Paginated(ctx, users, response.PaginationMeta{
    Page: 1,
    Limit: 10,
    Total: 100,
    TotalPages: 10,
})
```

### Error Responses

```go
// Automatic error handling (recommended)
if err != nil {
    return response.Error(ctx, err)
}

// Specific error responses
return response.BadRequestError(ctx, "Invalid input")
return response.UnauthorizedError(ctx, "Invalid token")
return response.NotFoundError(ctx, "Resource not found")
return response.ConflictError(ctx, "Resource already exists")
return response.InternalError(ctx, "Server error")

// Validation error with details
validationErrors := []ValidationError{...}
return response.ValidationError(ctx, "Validation failed", validationErrors)
```

### Response Format

**Success Response:**

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Operation successful",
  "data": {...},
  "meta": {
    "requestId": "req-123456",
    "timestamp": "2026-02-08T18:30:00Z"
  }
}
```

**Error Response:**

```json
{
  "success": false,
  "statusCode": 404,
  "error": {
    "code": "ERR_USER_NOT_FOUND",
    "message": "User with ID abc123 not found",
    "requestId": "req-123456",
    "timestamp": "2026-02-08T18:30:00Z",
    "metadata": {
      "user_id": "abc123"
    }
  }
}
```

## Middleware

### Error Handler Middleware

Add to your server initialization:

```go
import "worklayer/internal/app/middleware"

app := fiber.New()

// Add request context (must come first)
app.Use(middleware.RequestContext())

// Add error handler
app.Use(middleware.ErrorHandler())

// Your routes...
```

The error handler middleware:

- Catches panics and converts them to 500 errors
- Logs all errors with request context
- Converts Fiber errors to AppError
- Sends standardized error responses

### Request Context Middleware

```go
app.Use(middleware.RequestContext())
```

This middleware:

- Generates unique request ID for each request
- Adds request ID to response headers
- Makes request ID available in error responses and logs

## Logging

### Initialize Logger

```go
import "worklayer/pkg/logger"

// In main.go or initialization
isDevelopment := os.Getenv("ENV") != "production"
logger.InitLogger(isDevelopment)
```

### Log Errors

```go
// Errors are automatically logged by the error handler middleware
// But you can also log manually:

logger.LogError(appErr, requestID)
logger.LogInfo("User logged in", map[string]interface{}{
    "user_id": userID,
    "email": email,
})
logger.LogWarning("Rate limit approaching", map[string]interface{}{
    "user_id": userID,
    "requests": count,
})
```

### Log Format

**Development Mode (Pretty):**

```json
{
  "timestamp": "2026-02-08T18:30:00Z",
  "level": "ERROR",
  "code": "ERR_USER_NOT_FOUND",
  "message": "User with ID abc123 not found",
  "httpStatus": 404,
  "requestID": "req-123456",
  "metadata": {
    "user_id": "abc123"
  }
}
```

**Production Mode (Single Line):**

```json
{
  "timestamp": "2026-02-08T18:30:00Z",
  "level": "ERROR",
  "code": "ERR_USER_NOT_FOUND",
  "message": "User with ID abc123 not found",
  "httpStatus": 404,
  "requestID": "req-123456",
  "metadata": { "user_id": "abc123" }
}
```

## Migration Guide

### Repository Layer

**Before:**

```go
func (r *UserRepository) FindByID(id string) (*User, error) {
    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, domain.NewDomainError(404, "User not found")
        }
        return nil, err
    }
    return &user, nil
}
```

**After:**

```go
import (
    "worklayer/pkg/errors"
    "worklayer/internal/repository"
)

func (r *UserRepository) FindByID(id string) (*User, error) {
    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, repository.ConvertDBError(err, "finding user by ID")
    }
    return &user, nil
}
```

### Service Layer

**Before:**

```go
func (s *AuthService) LoginUser(email, password string) (*User, error) {
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        return nil, domain.NewServiceError(401, "Invalid credentials")
    }
    // ...
}
```

**After:**

```go
import "worklayer/internal/domain"

func (s *AuthService) LoginUser(email, password string) (*User, error) {
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        // If it's a not found error, return invalid credentials
        if errors.Is(err, errors.ErrDBRecordNotFound) {
            return nil, domain.InvalidCredentialsError()
        }
        // Otherwise wrap and return
        return nil, service.WrapRepositoryError(err, "login user")
    }
    // ...
}
```

### Controller Layer

**Before:**

```go
func (c *AuthController) Login(ctx *fiber.Ctx) error {
    user, err := c.service.LoginUser(email, password)
    if err != nil {
        return api.ErrorResponse(ctx, api.UNAUTHORIZED, err.Error(), nil)
    }
    return api.SuccessResponse(ctx, 200, "Login successful", user)
}
```

**After:**

```go
import "worklayer/pkg/response"

func (c *AuthController) Login(ctx *fiber.Ctx) error {
    user, err := c.service.LoginUser(email, password)
    if err != nil {
        return response.Error(ctx, err)
    }
    return response.SuccessWithMessage(ctx, fiber.StatusOK, "Login successful", user)
}
```

## Best Practices

### 1. Use Appropriate Error Codes

Choose the most specific error code:

```go
// ✓ Good
errors.UserNotFound(userID)

// ✗ Avoid
errors.NotFound("User not found")
```

### 2. Add Contextual Metadata

```go
// ✓ Good
err := errors.NotFound("Order not found")
err.WithMetadata("order_id", orderID)
err.WithMetadata("user_id", userID)

// ✗ Avoid
err := errors.NotFound("Order not found")
```

### 3. Wrap Errors at Boundaries

```go
// ✓ Good - Wrap at repository boundary
func (r *Repo) FindUser(id string) (*User, error) {
    err := r.db.Find(&user).Error
    if err != nil {
        return nil, repository.ConvertDBError(err, "finding user")
    }
    return &user, nil
}

// ✗ Avoid - Returning raw database errors
func (r *Repo) FindUser(id string) (*User, error) {
    err := r.db.Find(&user).Error
    return &user, err // Raw GORM error
}
```

### 4. Use Helper Functions

```go
// ✓ Good
return errors.InvalidCredentials()
return errors.TokenExpired()
return errors.UserNotFound(userID)

// ✗ Avoid
return errors.New(errors.ErrAuthInvalidCredentials)
return errors.NewWithMessage(errors.ErrAuthTokenExpired, "Token expired")
```

### 5. Handle Errors in Controllers

```go
// ✓ Good - Let middleware and response package handle errors
if err != nil {
    return response.Error(ctx, err)
}

// ✗ Avoid - Manual error handling
if err != nil {
    return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
}
```

### 6. Log at the Right Level

Errors are automatically logged by the middleware, so you don't need to log in most cases:

```go
// ✓ Good - Just return the error
if err != nil {
    return response.Error(ctx, err) // Automatically logged
}

// ✗ Avoid - Double logging
if err != nil {
    log.Printf("Error: %v", err) // Don't do this
    return response.Error(ctx, err)
}
```

### 7. Security Considerations

Never expose sensitive information in error messages:

```go
// ✓ Good
errors.InvalidCredentials() // Generic message

// ✗ Avoid
errors.Unauthorized("Password for '%s' is incorrect", email) // Leaks info
```

## Examples

### Complete Controller Example

```go
package controller

import (
    "worklayer/pkg/response"
    "worklayer/pkg/errors"
    "worklayer/internal/domain"
    "github.com/gofiber/fiber/v2"
)

type UserController struct {
    service *UserService
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
    userID := ctx.Params("id")
    if userID == "" {
        return response.Error(ctx, errors.InvalidParams("User ID is required"))
    }

    user, err := c.service.GetUser(userID)
    if err != nil {
        return response.Error(ctx, err)
    }

    return response.Success(ctx, user)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
    var req CreateUserRequest
    if err := ctx.BodyParser(&req); err != nil {
        return response.Error(ctx, errors.BadRequest("Invalid request body"))
    }

    user, err := c.service.CreateUser(req)
    if err != nil {
        return response.Error(ctx, err)
    }

    return response.Created(ctx, user)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
    userID := ctx.Params("id")

    var req UpdateUserRequest
    if err := ctx.BodyParser(&req); err != nil {
        return response.Error(ctx, errors.BadRequest("Invalid request body"))
    }

    user, err := c.service.UpdateUser(userID, req)
    if err != nil {
        return response.Error(ctx, err)
    }

    return response.SuccessWithMessage(ctx, fiber.StatusOK, "User updated successfully", user)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
    userID := ctx.Params("id")

    if err := c.service.DeleteUser(userID); err != nil {
        return response.Error(ctx, err)
    }

    return response.NoContent(ctx)
}
```

### Complete Service Example

```go
package service

import (
    "worklayer/internal/domain"
    "worklayer/internal/repository"
    "worklayer/pkg/errors"
)

type UserService struct {
    repo repository.UserRepository
}

func (s *UserService) GetUser(userID string) (*domain.User, error) {
    user, err := s.repo.FindByID(userID)
    if err != nil {
        if errors.Is(err, errors.ErrDBRecordNotFound) {
            return nil, domain.UserNotFoundError(userID)
        }
        return nil, WrapRepositoryError(err, "get user")
    }
    return user, nil
}

func (s *UserService) CreateUser(email, password string) (*domain.User, error) {
    // Check if user exists
    existing, err := s.repo.FindByEmail(email)
    if err == nil && existing != nil {
        return nil, domain.UserAlreadyExistsError(email)
    }

    // Create user
    user := &domain.User{Email: email}
    if err := s.repo.Create(user); err != nil {
        return nil, WrapRepositoryError(err, "create user")
    }

    return user, nil
}
```

## Testing

### Testing Error Creation

```go
func TestUserNotFound(t *testing.T) {
    err := errors.UserNotFound("user123")

    assert.Equal(t, errors.ErrUserNotFound, err.Code)
    assert.Equal(t, 404, err.HTTPStatus)
    assert.Contains(t, err.Message, "user123")
}
```

### Testing Error Responses

```go
func TestErrorResponse(t *testing.T) {
    app := fiber.New()
    app.Use(middleware.RequestContext())

    app.Get("/test", func(ctx *fiber.Ctx) error {
        return response.Error(ctx, errors.NotFound("Resource not found"))
    })

    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 404, resp.StatusCode)
    // Assert response body...
}
```

## Summary

The error management system provides:

✅ **Centralized error codes** - Easy to maintain and extend  
✅ **Rich context** - Stack traces, metadata, error chaining  
✅ **Standardized responses** - Consistent API responses  
✅ **Easy to use** - Helper functions and builders  
✅ **Production-ready** - Logging, monitoring, security  
✅ **Type-safe** - Full Go type system support

For questions or issues, please refer to the implementation files or contact the development team.
