package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	globalMiddleware "github.com/vyolayer/vyolayer/internal/shared/middleware"
)

// Server represents the API Gateway instance
type Server struct {
	app  *fiber.App
	port string
}

// Router interface for registering routes (ISP)
type RouteRegistrar interface {
	RegisterRoutes(router fiber.Router)
}

// New creates and configures a new Server instance
func New(port string) *Server {
	app := fiber.New(fiber.Config{AppName: "VyoLayer Gateway"})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Workspace-ID, X-Vyo-Key",
	}))

	// Inject custom middleware to propagate headers to gRPC requests
	app.Use(middleware.GRPCMetadataMiddleware())

	return &Server{
		app:  app,
		port: port,
	}
}

// RegisterRegistrars allows appending groups of routes (OCP)
func (s *Server) RegisterRegistrars(registrars ...RouteRegistrar) {
	v1 := s.app.Group("/v1")
	for _, registrar := range registrars {
		registrar.RegisterRoutes(v1)
	}

	v1.Use(globalMiddleware.NotFoundMiddleware)
}

// Start runs the HTTP server with graceful shutdown handling
func (s *Server) Start() error {
	// Graceful shutdown handling
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down HTTP server...")
		s.app.Shutdown()
	}()

	log.Printf("API Gateway listening on :%s", s.port)
	return s.app.Listen(fmt.Sprintf(":%s", s.port))
}
