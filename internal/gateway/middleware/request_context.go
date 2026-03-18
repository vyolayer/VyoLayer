package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

// RequestContext enriches each request with a request ID used in responses and logs.
func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := ulid.MustNew(ulid.Now(), ulid.DefaultEntropy()).String()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}
