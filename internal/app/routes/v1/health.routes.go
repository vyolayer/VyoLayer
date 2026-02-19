package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerHealthRoutes(router fiber.Router, d *dependencies) {
	health := router.Group("/health")

	health.Get("/", d.HealthCtrl.HealthCheck)
}
