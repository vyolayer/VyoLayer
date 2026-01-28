package middleware

import (
	"fmt"
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return response.Error(
		c,
		response.NotFoundError(fmt.Sprintf("%s not found", c.Path())),
	)
}
