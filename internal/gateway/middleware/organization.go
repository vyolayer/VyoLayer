package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
)

func ValidateOrganizationID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Params("organizationID")
		if orgID == "" {
			return response.Error(c, errors.BadRequest("organization id is required"))
		}
		c.Locals("organization_id", orgID)
		return c.Next()
	}
}
