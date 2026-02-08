package response

import (
	"worklayer/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
)

const (
	BAD_REQUEST           = "BAD_REQUEST"           // 400
	UNAUTHORIZED          = "UNAUTHORIZED"          // 401
	FORBIDDEN             = "FORBIDDEN"             // 403
	NOT_FOUND             = "NOT_FOUND"             // 404
	METHOD_NOT_ALLOWED    = "METHOD_NOT_ALLOWED"    // 405
	CONFLICT              = "CONFLICT"              // 409
	VALIDATION_ERROR      = "VALIDATION_ERROR"      // 422
	TOO_MANY_REQUESTS     = "TOO_MANY_REQUESTS"     // 429
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR" // 500
	NOT_IMPLEMENTED       = "NOT_IMPLEMENTED"       // 501
	BAD_GATEWAY           = "BAD_GATEWAY"           // 502
	SERVICE_UNAVAILABLE   = "SERVICE_UNAVAILABLE"   // 503
	GATEWAY_TIMEOUT       = "GATEWAY_TIMEOUT"       // 504
)

type ErrorResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Errors     any    `json:"errors,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(statusCode int, code string, message string, errors any) *ErrorResponse {
	switch statusCode {
	case fiber.StatusBadRequest:
		code = BAD_REQUEST
	case fiber.StatusUnauthorized:
		code = UNAUTHORIZED
	case fiber.StatusForbidden:
		code = FORBIDDEN
	case fiber.StatusNotFound:
		code = NOT_FOUND
	case fiber.StatusMethodNotAllowed:
		code = METHOD_NOT_ALLOWED
	case fiber.StatusConflict:
		code = CONFLICT
	case fiber.StatusUnprocessableEntity:
		code = VALIDATION_ERROR
	case fiber.StatusTooManyRequests:
		code = TOO_MANY_REQUESTS
	case fiber.StatusInternalServerError:
		code = INTERNAL_SERVER_ERROR
	case fiber.StatusNotImplemented:
		code = NOT_IMPLEMENTED
	case fiber.StatusBadGateway:
		code = BAD_GATEWAY
	case fiber.StatusServiceUnavailable:
		code = SERVICE_UNAVAILABLE
	case fiber.StatusGatewayTimeout:
		code = GATEWAY_TIMEOUT
	}

	return &ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Errors:     errors,
	}
}

func NewErrorMessage(statusCode int, message string) *ErrorResponse {
	code := ""
	switch statusCode {
	case fiber.StatusBadRequest:
		code = BAD_REQUEST
	case fiber.StatusUnauthorized:
		code = UNAUTHORIZED
	case fiber.StatusForbidden:
		code = FORBIDDEN
	case fiber.StatusNotFound:
		code = NOT_FOUND
	case fiber.StatusMethodNotAllowed:
		code = METHOD_NOT_ALLOWED
	case fiber.StatusConflict:
		code = CONFLICT
	case fiber.StatusUnprocessableEntity:
		code = VALIDATION_ERROR
	case fiber.StatusTooManyRequests:
		code = TOO_MANY_REQUESTS
	case fiber.StatusInternalServerError:
		code = INTERNAL_SERVER_ERROR
	case fiber.StatusNotImplemented:
		code = NOT_IMPLEMENTED
	case fiber.StatusBadGateway:
		code = BAD_GATEWAY
	case fiber.StatusServiceUnavailable:
		code = SERVICE_UNAVAILABLE
	case fiber.StatusGatewayTimeout:
		code = GATEWAY_TIMEOUT
	}

	return NewErrorResponse(statusCode, code, message, nil)
}

func BadRequestError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusBadRequest, BAD_REQUEST, message, nil)
}

func UnauthorizedError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusUnauthorized, UNAUTHORIZED, message, nil)
}

func ForbiddenError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusForbidden, FORBIDDEN, message, nil)
}

func NotFoundError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusNotFound, NOT_FOUND, message, nil)
}

func MethodNotAllowedError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusMethodNotAllowed, METHOD_NOT_ALLOWED, message, nil)
}

func ConflictError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusConflict, CONFLICT, message, nil)
}

func ValidationError(message string, validationErrors []*validation.ErrorResponse) *ErrorResponse {
	errs := make([]any, 0)
	for _, details := range validationErrors {
		errs = append(errs, map[string]any{
			"field":   details.Field,
			"tag":     details.Tag,
			"message": details.Message,
		})
	}

	return NewErrorResponse(fiber.StatusUnprocessableEntity, VALIDATION_ERROR, message, errs)
}

func TooManyRequestsError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusTooManyRequests, TOO_MANY_REQUESTS, message, nil)
}

func InternalServerError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusInternalServerError, INTERNAL_SERVER_ERROR, message, nil)
}

func NotImplementedError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusNotImplemented, NOT_IMPLEMENTED, message, nil)
}

func BadGatewayError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusBadGateway, BAD_GATEWAY, message, nil)
}

func ServiceUnavailableError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusServiceUnavailable, SERVICE_UNAVAILABLE, message, nil)
}

func GatewayTimeoutError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusGatewayTimeout, GATEWAY_TIMEOUT, message, nil)
}

func Error(c *fiber.Ctx, err *ErrorResponse) error {
	return c.Status(err.StatusCode).JSON(err)
}
