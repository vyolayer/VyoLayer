package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/gateway/config"
	"github.com/vyolayer/vyolayer/internal/gateway/server"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/internal/gateway/wire"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
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
	appLogger := logger.NewAppLogger("GATEWAY")

	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		appLogger.Warn("Note: No .env file found; relying on system environment variables", "")
	}

	// Load Configuration
	cfg := config.Load()

	// Initialize gRPC Connections
	clients, err := wire.NewClients(appLogger, cfg, grpcTimeout)
	if err != nil {
		appLogger.Error("Failed to initialize gRPC clients", err.Error())
	}
	defer clients.Close()

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

	// IAM
	iamCookieSrv := service.NewIAMCookieService(service.IAMCookieConfig{
		Atcc: fiber.Cookie{
			Name:     string(service.IAMCookieAccessToken),
			Expires:  time.Now().Add(5 * time.Minute), // TODO: Get from config
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
		},
		Rtcc: fiber.Cookie{
			Name:     string(service.IAMCookieRefreshToken),
			Expires:  time.Now().Add(2 * time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteStrictMode,
		},
	})

	iamSession := jwt.NewIamJWT(
		cfg.IAMJWT.AccessTokenSecret,
		cfg.IAMJWT.AccessTokenExpiry,
		cfg.IAMJWT.RefreshTokenExpiry,
	)

	// Initialize Handlers
	registrars := wire.NewRegistrars(
		appLogger,
		clients,
		cookieSrv,
		accountJWT,
		iamCookieSrv,
		iamSession,
	)

	// Setup and Start Server
	srv := server.New(cfg.HTTPPort)
	srv.RegisterRegistrars(registrars...)

	return srv.Start()
}
