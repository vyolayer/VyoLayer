package middleware

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	applogger "github.com/vyolayer/vyolayer/pkg/logger"
)

// RequestLogger logs one structured entry per request following API response semantics.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latencyMs := time.Since(start).Milliseconds()

		status := c.Response().StatusCode()
		if status == 0 {
			status = fiber.StatusOK
		}

		requestID, _ := c.Locals("requestID").(string)
		message, _ := c.Locals("response_message").(string)
		errorCode, _ := c.Locals("response_error_code").(string)

		success := status < fiber.StatusBadRequest
		if localSuccess, ok := c.Locals("response_success").(bool); ok {
			success = localSuccess
		}

		if message == "" {
			message = http.StatusText(status)
		}
		logFields := map[string]interface{}{
			"requestId":  requestID,
			"statusCode": status,
			"success":    success,
			"message":    message,
			"method":     c.Method(),
			"path":       c.OriginalURL(),
			"ip":         c.IP(),
			"latencyMs":  latencyMs,
		}
		if errorCode != "" {
			logFields["errorCode"] = errorCode
		}

		if success {
			applogger.LogInfo("http_success", logFields)

		} else {
			applogger.LogWarning("http_error", logFields)
		}

		return err
	}
}
