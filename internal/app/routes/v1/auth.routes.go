package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerAuthRoutes(router fiber.Router, d *dependencies) {
	auth := router.Group("/auth")

	auth.Post("/register", d.AuthCtrl.RegisterUser)
	auth.Post("/login", d.AuthCtrl.LoginUser)
	auth.Post("/refresh", d.AuthCtrl.RefreshSession)

	auth.Use(d.AuthMiddleware.JwtValidated())

	auth.Post("/validate", d.AuthCtrl.ValidateSession)
	auth.Post("/logout", d.AuthCtrl.LogoutUser)
}
