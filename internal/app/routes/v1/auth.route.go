package routes

import (
	"log"
	"worklayer/internal/app/controller"
	"worklayer/internal/app/middleware"
	"worklayer/internal/config"
	"worklayer/internal/repository"
	"worklayer/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthRoute interface {
	SetupRoutes()
}

type AuthRouteDependencies struct {
	AuthConfig  config.AuthConfig
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
}

type authRoute struct {
	router       fiber.Router
	dependencies AuthRouteDependencies
}

func NewAuthRouter(router fiber.Router, deps AuthRouteDependencies) AuthRoute {
	return &authRoute{
		router:       router,
		dependencies: deps,
	}
}

func (ar *authRoute) SetupRoutes() {
	log.Println("Setting up auth routes")

	// Initialize services
	// Core services
	authService := service.NewAuthService(ar.dependencies.UserRepo)
	sessionService := service.NewSessionService(ar.dependencies.UserRepo, ar.dependencies.SessionRepo)

	// Utility services
	tokenService := service.NewTokenService(ar.dependencies.AuthConfig)

	// Initialize controller
	authController := controller.NewAuthController(authService, tokenService, sessionService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Register routes
	authRouter := ar.router.Group("/auth")
	authRouter.Post("/register", authController.RegisterUser)
	authRouter.Post("/login", authController.LoginUser)
	authRouter.Post("/refresh", authController.RefreshSession)

	authRouter.Use(authMiddleware.JwtValidated())
	authRouter.Post("/validate", authController.ValidateSession)
	authRouter.Post("/logout", authController.LogoutUser)
}
