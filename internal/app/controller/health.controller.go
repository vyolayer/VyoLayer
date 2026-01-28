package controller

import (
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) HealthCheck(c *fiber.Ctx) error {
	return response.Success(
		c,
		response.NewSuccessResponse(
			fiber.StatusOK,
			"Welcome to WorkLayer",
			fiber.Map{
				"version": "1.0.0",
			},
		),
	)
}
