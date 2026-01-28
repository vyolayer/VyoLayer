package routes

import (
	"worklayer/internal/app/controller"

	"github.com/gofiber/fiber/v2"
)

func HealthRoutes(app *fiber.App) {
	healthController := controller.NewHealthController()

	app.Get("/health", healthController.HealthCheck)
}
