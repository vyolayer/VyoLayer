package controller

import (
	"vyolayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check the status of the API.
// @Tags health
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /health [get]
func (h *HealthController) HealthCheck(c *fiber.Ctx) error {
	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Welcome to VyoLayer",
		map[string]string{
			"version": "1.0.0",
			"status":  "healthy",
		},
	)
}
