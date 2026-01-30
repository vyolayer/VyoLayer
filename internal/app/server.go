package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"worklayer/internal/app/middleware"
	"worklayer/internal/app/routes/v1"
	"worklayer/internal/config"
	"worklayer/internal/platform/database"
	"worklayer/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var (
	isConfigLoaded      = false
	isAppInitialized    = false
	isDatabaseConnected = false
	isRoutesSetup       = false
)

type App struct {
	app *fiber.App
	cfg *config.Config
	db  *gorm.DB
}

func New() *App {
	app := fiber.New()

	if !isAppInitialized {
		isAppInitialized = true
	}

	return &App{
		app: app,
	}
}

// Setup middleware
func (a *App) SetupMiddleware() {
	a.app.Use(recover.New())
	a.app.Use(cors.New(cors.Config{}))
	a.app.Use(middleware.NewRequestIDMiddleware().RequestIDMiddleware)
	a.app.Use(logger.New())
	a.app.Use(middleware.ErrorMiddleware)
}

// Setup routes
func (a *App) SetupRoutes() {
	userRepo := repository.NewUserRepository(a.db)
	sessionRepo := repository.NewSessionRepository(a.db)

	api := a.app.Group("/api")
	apiV1 := api.Group("/v1")

	// Health routes
	routes.NewHealthRouter(apiV1).SetupRoutes()
	// Auth routes
	routes.NewAuthRouter(
		apiV1,
		routes.AuthRouteDependencies{
			AuthConfig:  a.cfg.Auth,
			UserRepo:    userRepo,
			SessionRepo: sessionRepo,
		},
	).SetupRoutes()

	// Not found middleware
	a.app.Use(middleware.NotFoundMiddleware)

	if !isRoutesSetup {
		isRoutesSetup = true
	}
}

// Run app
func (a *App) Run() error {
	if !isConfigLoaded {
		return errors.New("config not loaded")
	}

	if !isAppInitialized {
		return errors.New("app not initialized")
	}

	if !isDatabaseConnected {
		return errors.New("database not connected")
	}

	if !isRoutesSetup {
		return errors.New("routes not setup")
	}

	go func() {
		log.Println("Server starting on :", a.cfg.Server.Port)
		if err := a.app.Listen(":" + strconv.Itoa(a.cfg.Server.Port)); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	return nil
}

func (a *App) LoadConfig() {
	if !isConfigLoaded {
		isConfigLoaded = true
	}

	cfg, err := config.Load("config/config.dev.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	a.cfg = cfg
	log.Println("Loading config...")
}

func (app *App) ConnectToDatabase() {
	db, err := database.Init(&app.cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	app.db = db

	if !isDatabaseConnected {
		isDatabaseConnected = true
	}
	log.Println("Database connected...")
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := a.app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}

	log.Println("Cleanup complete. Goodbye!")
}

func (a *App) ListenShutdownEvent() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Shutdown signal received...")
	a.Shutdown()
}
