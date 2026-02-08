package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestContext enriches the request context with metadata
func RequestContext() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Generate a unique request ID
		requestID := uuid.New().String()

		// Store in context for later use
		ctx.Locals("requestID", requestID)

		// Add to response headers
		ctx.Set("X-Request-ID", requestID)

		// Continue with the request
		return ctx.Next()
	}
}
