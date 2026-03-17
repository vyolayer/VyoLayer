package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/account/config"
	accountGrpc "github.com/vyolayer/vyolayer/internal/account/delivery/grpc"
	"github.com/vyolayer/vyolayer/internal/account/repository"
	"github.com/vyolayer/vyolayer/internal/account/server"
	"github.com/vyolayer/vyolayer/internal/account/usecase"
	apikey "github.com/vyolayer/vyolayer/internal/shared/api-key"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/mail"
	"github.com/vyolayer/vyolayer/pkg/postgres"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("account service error: %v", err)
	}
}

func run() error {
	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found; relying on system environment variables")
	}

	// Load Configuration
	cfg := config.Load()

	// Setup Database Connection
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve underlying sql database object: %w", err)
	}
	defer sqlDB.Close()

	// Initialize Dependency Injection Container
	//   Core Utilities & Security
	apiKeyVerifier := apikey.NewAPIKeyVerifier(db)
	accountJWT := jwt.NewAccountJWT(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.AccessTokenExpiry,
	)
	mailer := mail.NewMockMailer() // TODO: Read environment/config to determine if real mailer should be built

	//   Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tokenRepo := repository.NewVerificationTokenRepository(db)

	//   Usecases
	accountUsecase := usecase.NewAccountUsecase(
		userRepo,
		sessionRepo,
		tokenRepo,
		mailer,
		accountJWT,
	)

	// Handlers
	accountHandler := accountGrpc.NewAccountHandler(accountUsecase)

	// Setup and Start Server
	grpcSrv := server.NewGRPCServer(cfg.GRPCPort, apiKeyVerifier)
	accountV1.RegisterAccountServiceServer(grpcSrv.Engine, accountHandler)

	return grpcSrv.Start()
}
