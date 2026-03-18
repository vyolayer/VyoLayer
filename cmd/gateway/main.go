package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/gateway/config"
	"github.com/vyolayer/vyolayer/internal/gateway/handlers"
	"github.com/vyolayer/vyolayer/internal/gateway/server"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/grpcutil"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

const (
	sessionCookieName = "vyo_session"
	authCookieName    = "vyo_auth"
	grpcTimeout       = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("gateway error: %v", err)
	}
}

func run() error {
	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found; relying on system environment variables")
	}

	// Load Configuration
	cfg := config.Load()

	// Initialize gRPC Connections
	accountConn, err := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.AccountServiceAddr,
		Timeout: grpcTimeout,
	})
	if err != nil {
		return err
	}
	defer accountConn.Close()

	accountClient := accountV1.NewAccountServiceClient(accountConn)

	// Initialize Services and Utilities
	cookieSrv := service.NewAccountTokenService(service.AccountCookieConfig{
		AccessTokenCookieConfig: fiber.Cookie{
			Name:     sessionCookieName,
			Expires:  time.Now().Add(cfg.AccountJWT.AccessTokenExpiry),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
		},
		RefreshTokenCookieConfig: fiber.Cookie{
			Name:     authCookieName,
			Expires:  time.Now().Add(cfg.AccountJWT.RefreshTokenExpiry),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
		},
	})

	accountJWT := jwt.NewAccountJWT(
		cfg.AccountJWT.AccessTokenSecret,
		cfg.AccountJWT.AccessTokenExpiry,
	)

	// Initialize Handlers
	accountHandler := handlers.NewAccountHandler(
		accountClient,
		cookieSrv,
		accountJWT,
	)

	// Setup and Start Server
	srv := server.New(cfg.HTTPPort)
	srv.RegisterRegistrars(accountHandler)

	return srv.Start()
}
