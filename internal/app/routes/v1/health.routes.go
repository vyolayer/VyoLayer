package routes

import (
	"worklayer/internal/app/controller"

	"github.com/gofiber/fiber/v2"
)

type HealthRoute struct {
	router fiber.Router
}

func NewHealthRouter(router fiber.Router) *HealthRoute {
	return &HealthRoute{
		router: router,
	}
}

func (hr *HealthRoute) SetupRoutes() {
	healthController := controller.NewHealthController()

	hr.router.Get("/health", healthController.HealthCheck)
}
