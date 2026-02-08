package controller

import (
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
// @Success 200 {object} response.Response "Welcome to WorkLayer"
// @Router /health [get]
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
