package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	grpcTimeout = 10 * time.Second
)

type GrpcCtxMiddleware struct {
	timeout time.Duration
}

func NewGrpcCtxMiddleware(timeout time.Duration) *GrpcCtxMiddleware {
	return &GrpcCtxMiddleware{timeout: timeout}
}

func (m *GrpcCtxMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), m.timeout)
		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}
