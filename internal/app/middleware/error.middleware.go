package middleware

import (
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return response.Error(
			c,
			response.InternalServerError(err.Error()),
		)
	}
	return nil
}
