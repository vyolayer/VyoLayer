package console

import "github.com/gofiber/fiber/v2"

func getProjectIDFromLocals(c *fiber.Ctx) string {
	id, _ := c.Locals("project_id").(string)
	return id
}
