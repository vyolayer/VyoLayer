package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func grpcCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	return ctx, cancel
}

func grpcCtxMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}
