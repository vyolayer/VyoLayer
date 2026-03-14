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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/vyolayer/vyolayer/internal/app/middleware"
	v1 "github.com/vyolayer/vyolayer/internal/app/routes/v1"
	"github.com/vyolayer/vyolayer/internal/config"
	"github.com/vyolayer/vyolayer/internal/platform/database"
	"gorm.io/gorm"

	_ "github.com/vyolayer/vyolayer/docs"
)

// @title VyoLayer API
// @version 1.0
// @description This is the VyoLayer API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:6999
// @BasePath /api/v1

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
	// Request context MUST come first
	a.app.Use(middleware.RequestContext())

	// Error handler with panic recovery
	a.app.Use(middleware.ErrorHandler())

	// Other middleware
	a.app.Use(cors.New(cors.Config{
		AllowOrigins:     a.cfg.App.Cors,
		AllowCredentials: true,
	}))
	a.app.Use(logger.New())
	a.app.Use(recover.New())
	// appLogger.InitLogger(true)

	a.app.Get("/metrics", monitor.New())

	a.app.Get("/swagger/*", swagger.HandlerDefault)
}

// Setup routes
func (a *App) SetupRoutes() {
	// V1 routes
	v1.New(a.app, a.cfg, a.db).Register()

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
