package middleware

import (
	"fmt"
	"runtime/debug"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/logger"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a global error handler middleware
// It catches panics and handles errors uniformly
func ErrorHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				stackTrace := string(debug.Stack())
				logger.LogWarning("Panic recovered", map[string]interface{}{
					"panic":      fmt.Sprintf("%v", r),
					"stackTrace": stackTrace,
					"path":       ctx.Path(),
					"method":     ctx.Method(),
					"requestID":  response.GetRequestID(ctx),
				})

				// Create an internal error
				err := errors.Internal("An unexpected error occurred")
				err.WithMetadata("panic", fmt.Sprintf("%v", r))

				// Log the error
				logger.LogError(err, response.GetRequestID(ctx))

				// Send error response
				response.Error(ctx, err)
			}
		}()

		// Continue with the request
		err := ctx.Next()

		// If there's an error, log and handle it
		if err != nil {
			// Check if it's already an AppError
			if appErr, ok := errors.As(err); ok {
				logger.LogError(appErr, response.GetRequestID(ctx))
				return response.Error(ctx, appErr)
			}

			// Check if it's a Fiber error
			if fiberErr, ok := err.(*fiber.Error); ok {
				// Convert Fiber error to AppError
				appErr := convertFiberError(fiberErr)
				logger.LogError(appErr, response.GetRequestID(ctx))
				return response.Error(ctx, appErr)
			}

			// For any other error, wrap as internal error
			appErr := errors.InternalWrap(err, "An unexpected error occurred")
			logger.LogError(appErr, response.GetRequestID(ctx))
			return response.Error(ctx, appErr)
		}

		return nil
	}
}

// convertFiberError converts a Fiber error to an AppError
func convertFiberError(fiberErr *fiber.Error) *errors.AppError {
	var code errors.ErrorCode

	switch fiberErr.Code {
	case fiber.StatusBadRequest:
		code = errors.ErrRequestInvalidBody
	case fiber.StatusUnauthorized:
		code = errors.ErrAuthUnauthorized
	case fiber.StatusForbidden:
		code = errors.ErrAuthForbidden
	case fiber.StatusNotFound:
		code = errors.ErrResourceNotFound
	case fiber.StatusMethodNotAllowed:
		code = errors.ErrBusinessOperationNotAllowed
	case fiber.StatusConflict:
		code = errors.ErrResourceConflict
	case fiber.StatusUnprocessableEntity:
		code = errors.ErrValidationFailed
	case fiber.StatusTooManyRequests:
		code = errors.ErrRequestRateLimited
	case fiber.StatusInternalServerError:
		code = errors.ErrInternalUnexpected
	case fiber.StatusNotImplemented:
		code = errors.ErrInternalNotImplemented
	case fiber.StatusServiceUnavailable:
		code = errors.ErrExternalServiceUnavailable
	case fiber.StatusGatewayTimeout:
		code = errors.ErrExternalServiceTimeout
	default:
		code = errors.ErrInternalUnexpected
	}

	return errors.NewWithMessage(code, fiberErr.Message)
}

// Custom404Handler handles 404 Not Found errors
func Custom404Handler(ctx *fiber.Ctx) error {
	err := errors.NotFound("Route '%s' not found", ctx.Path())
	err.WithMetadata("method", ctx.Method())
	err.WithMetadata("path", ctx.Path())

	logger.LogError(err, response.GetRequestID(ctx))
	return response.Error(ctx, err)
}
