package middleware

import (
	"fmt"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return response.Error(
		c,
		errors.NotFound(fmt.Sprintf("Route %s not found", c.Path())),
	)
}
