package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/response"
)

// ValidateProjectID extracts `:projectID` from the URL path, stores it in
// Fiber locals, and aborts with 400 if it is missing.
func ValidateProjectID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		projectID := c.Params("projectID")
		if projectID == "" {
			return response.Error(c, errors.BadRequest("project id is required"))
		}
		c.Locals("project_id", projectID)
		return c.Next()
	}
}
