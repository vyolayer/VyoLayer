package routes

import (
	"worklayer/internal/app/controller"
	"worklayer/internal/app/middleware"
	"worklayer/internal/config"
	"worklayer/internal/repository"
	"worklayer/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Routes interface {
	SetupRoutes()
}

type routesV1 struct {
	router fiber.Router
	cfg    *config.Config
	db     *gorm.DB
}

func NewV1Routes(router fiber.Router, config *config.Config, db *gorm.DB) Routes {
	return &routesV1{
		router: router,
		cfg:    config,
		db:     db,
	}
}

func (router *routesV1) SetupRoutes() {
	// Initialize repositories
	repo := repository.NewRegistry(router.db)

	// Core services
	authService := service.NewAuthService(repo.User)
	sessionService := service.NewSessionService(repo.User, repo.Session)
	userService := service.NewUserService(repo.User)

	// Utility services
	tokenService := service.NewTokenService(router.cfg.Auth)

	// Controller
	authController := controller.NewAuthController(authService, tokenService, sessionService)
	userController := controller.NewUserController(userService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	// Register routes
	//
	// Auth routes
	authRouter := router.router.Group("/auth")
	authRouter.Post("/register", authController.RegisterUser)
	authRouter.Post("/login", authController.LoginUser)
	authRouter.Post("/refresh", authController.RefreshSession)

	authRouter.Use(authMiddleware.JwtValidated())
	authRouter.Post("/validate", authController.ValidateSession)
	authRouter.Post("/logout", authController.LogoutUser)

	// User routes
	userRouter := router.router.Group("/users")
	userRouter.Use(authMiddleware.JwtValidated())
	userRouter.Get("/me", userController.GetMe)
}
