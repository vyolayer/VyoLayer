package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	appLogger "github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
)

// ErrorHandler enforces a uniform API error response shape and structured error logs.
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())
				appLogger.LogWarning("Panic recovered", map[string]interface{}{
					"panic":     fmt.Sprintf("%v", r),
					"stackTrace": stackTrace,
					"path":      c.Path(),
					"method":    c.Method(),
					"requestID": response.GetRequestID(c),
				})

				err := errors.Internal("An unexpected error occurred")
				err.WithMetadata("panic", fmt.Sprintf("%v", r))
				appLogger.LogError(err, response.GetRequestID(c))
				_ = response.Error(c, err)
			}
		}()

		err := c.Next()
		if err == nil {
			return nil
		}

		if appErr, ok := errors.As(err); ok {
			appLogger.LogError(appErr, response.GetRequestID(c))
			return response.Error(c, appErr)
		}

		if fiberErr, ok := err.(*fiber.Error); ok {
			appErr := convertFiberError(fiberErr)
			appLogger.LogError(appErr, response.GetRequestID(c))
			return response.Error(c, appErr)
		}

		appErr := errors.InternalWrap(err, "An unexpected error occurred")
		appLogger.LogError(appErr, response.GetRequestID(c))
		return response.Error(c, appErr)
	}
}

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
