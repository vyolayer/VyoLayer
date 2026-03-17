package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return response.Error(
		c,
		errors.NotFound(fmt.Sprintf("Route %s not found", c.Path())),
	)
}
