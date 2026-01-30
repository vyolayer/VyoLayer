package controller

import (
	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) HealthCheck(c *fiber.Ctx) error {
	return Success(
		c,
		fiber.StatusOK,
		"Welcome to WorkLayer",
		map[string]string{
			"version": "1.0.0",
		},
	)
}
