package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerUserRoutes(router fiber.Router, d *dependencies) {
	users := router.Group("/users")
	users.Use(d.AuthMiddleware.JwtValidated())

	users.Get("/me", d.UserCtrl.GetMe)
}
