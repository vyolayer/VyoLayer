package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success    bool        `json:"success" example:"true"`
	StatusCode int         `json:"statusCode" example:"200"`
	Message    string      `json:"message,omitempty" example:"Operation successful"`
	Data       interface{} `json:"data,omitempty"`
	Meta       *Meta       `json:"meta,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success    bool         `json:"success" example:"false"`
	StatusCode int          `json:"statusCode" example:"400"`
	Error      *ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code      string                 `json:"code" example:"ERR_VALIDATION_FAILED"`
	Message   string                 `json:"message" example:"Validation failed"`
	Details   interface{}            `json:"details,omitempty"`
	RequestID string                 `json:"requestId,omitempty" example:"req-123456"`
	Timestamp time.Time              `json:"timestamp" example:"2026-02-08T18:30:00Z"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Meta contains metadata about the response
type Meta struct {
	RequestID  string    `json:"requestId,omitempty" example:"req-123456"`
	Timestamp  time.Time `json:"timestamp" example:"2026-02-08T18:30:00Z"`
	Page       int       `json:"page,omitempty" example:"1"`
	Limit      int       `json:"limit,omitempty" example:"10"`
	Total      int64     `json:"total,omitempty" example:"100"`
	TotalPages int       `json:"totalPages,omitempty" example:"10"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"totalPages" example:"10"`
}

// GetRequestID extracts the request ID from the fiber context
func GetRequestID(ctx *fiber.Ctx) string {
	requestID := ctx.Locals("requestID")
	if requestID == nil {
		return ""
	}
	if id, ok := requestID.(string); ok {
		return id
	}
	return ""
}

// Success sends a successful response with data
func Success(ctx *fiber.Ctx, data interface{}) error {
	return SuccessWithMessage(ctx, fiber.StatusOK, "Success", data)
}

// SuccessWithMessage sends a successful response with a custom message
func SuccessWithMessage(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := SuccessResponse{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta: &Meta{
			RequestID: GetRequestID(ctx),
			Timestamp: time.Now(),
		},
	}

	return ctx.Status(statusCode).JSON(response)
}

// SuccessMessage sends a successful response with only a message (no data)
func SuccessMessage(ctx *fiber.Ctx, message string) error {
	response := SuccessResponse{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    message,
		Meta: &Meta{
			RequestID: GetRequestID(ctx),
			Timestamp: time.Now(),
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Created sends a 201 Created response
func Created(ctx *fiber.Ctx, data interface{}) error {
	return SuccessWithMessage(ctx, fiber.StatusCreated, "Resource created successfully", data)
}

// NoContent sends a 204 No Content response
func NoContent(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNoContent)
}

// Paginated sends a successful response with pagination metadata
func Paginated(ctx *fiber.Ctx, data interface{}, pagination PaginationMeta) error {
	response := SuccessResponse{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Data:       data,
		Meta: &Meta{
			RequestID:  GetRequestID(ctx),
			Timestamp:  time.Now(),
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Error sends an error response
// It automatically handles both AppError and standard errors
func Error(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return InternalError(ctx, "An unexpected error occurred")
	}

	// Check if it's an AppError
	if appErr, ok := errors.As(err); ok {
		return sendAppError(ctx, appErr)
	}

	// For standard errors, wrap as internal error
	appErr := errors.InternalWrap(err, "An unexpected error occurred")
	return sendAppError(ctx, appErr)
}

// sendAppError sends an AppError as an HTTP response
func sendAppError(ctx *fiber.Ctx, err *errors.AppError) error {
	// Determine if we should include stack trace (only in development)
	// You might want to check an environment variable here
	var metadata map[string]interface{}
	if shouldIncludeStackTrace() {
		metadata = err.Metadata
		// Optionally add stack trace to metadata
		if len(err.StackTrace) > 0 {
			metadata["stackTrace"] = err.StackTrace
		}
	} else {
		// In production, only include safe metadata
		metadata = filterSafeMetadata(err.Metadata)
	}

	response := ErrorResponse{
		Success:    false,
		StatusCode: err.HTTPStatus,
		Error: &ErrorDetail{
			Code:      string(err.Code),
			Message:   err.Message,
			RequestID: GetRequestID(ctx),
			Timestamp: err.Timestamp,
			Metadata:  metadata,
		},
	}

	// If there are validation errors in metadata, move them to details
	if validationErrors, ok := err.Metadata["validation_errors"]; ok {
		response.Error.Details = validationErrors
	}

	return ctx.Status(err.HTTPStatus).JSON(response)
}

// Helper functions for common error responses

// BadRequestError sends a 400 Bad Request error
func BadRequestError(ctx *fiber.Ctx, message string) error {
	err := errors.BadRequest(message)
	return sendAppError(ctx, err)
}

// UnauthorizedError sends a 401 Unauthorized error
func UnauthorizedError(ctx *fiber.Ctx, message string) error {
	err := errors.Unauthorized(message)
	return sendAppError(ctx, err)
}

// ForbiddenError sends a 403 Forbidden error
func ForbiddenError(ctx *fiber.Ctx, message string) error {
	err := errors.Forbidden(message)
	return sendAppError(ctx, err)
}

// NotFoundError sends a 404 Not Found error
func NotFoundError(ctx *fiber.Ctx, message string) error {
	err := errors.NotFound(message)
	return sendAppError(ctx, err)
}

// ConflictError sends a 409 Conflict error
func ConflictError(ctx *fiber.Ctx, message string) error {
	err := errors.Conflict(message)
	return sendAppError(ctx, err)
}

// ValidationError sends a 422 Validation Error with details
func ValidationError(ctx *fiber.Ctx, message string, details interface{}) error {
	err := errors.ValidationFailedWithDetails(message, details)
	return sendAppError(ctx, err)
}

// InternalError sends a 500 Internal Server Error
func InternalError(ctx *fiber.Ctx, message string) error {
	err := errors.Internal(message)
	return sendAppError(ctx, err)
}

// shouldIncludeStackTrace determines if stack traces should be included
// This should check an environment variable in production
func shouldIncludeStackTrace() bool {
	// TODO: Check environment variable
	// For now, return false (production mode)
	return false
}

// filterSafeMetadata removes sensitive metadata keys
func filterSafeMetadata(metadata map[string]interface{}) map[string]interface{} {
	if metadata == nil {
		return nil
	}

	// List of keys that should NOT be included in production responses
	sensitiveKeys := map[string]bool{
		"password":       true,
		"token":          true,
		"secret":         true,
		"api_key":        true,
		"stackTrace":     true,
		"internal_error": true,
	}

	safe := make(map[string]interface{})
	for k, v := range metadata {
		if !sensitiveKeys[k] {
			safe[k] = v
		}
	}

	return safe
}
