package response

import "github.com/gofiber/fiber/v2"

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
	Success    bool              `json:"success"`
	StatusCode int               `json:"statusCode"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors,omitempty"`
}

func NewErrorResponse(statusCode int, code string, message string, errors map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Errors:     errors,
	}
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

func ValidationError(message string) *ErrorResponse {
	return NewErrorResponse(fiber.StatusUnprocessableEntity, VALIDATION_ERROR, message, nil)
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
