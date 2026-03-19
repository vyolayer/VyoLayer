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

	var mailer mail.Mailer
	if cfg.Mail.UseMock {
		mailer = mail.NewMockMailer()
	} else {
		mailer = mail.NewSMTPMailer(mail.SMTPConfig{
			Host:     cfg.Mail.Host,
			Port:     cfg.Mail.Port,
			Username: cfg.Mail.Username,
			Password: cfg.Mail.Password,
			From:     cfg.Mail.From,
		})
	}

	//   Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tokenRepo := repository.NewVerificationTokenRepository(db)

	//   Usecases
	accountUsecase := usecase.NewAccountUsecase(
		cfg,
		userRepo,
		sessionRepo,
		tokenRepo,
		mailer,
		accountJWT,
	)

	sessionUsecase := usecase.NewSessionUsecase(
		sessionRepo,
		accountJWT,
	)

	accountRecoverUsecase := usecase.NewAccountRecoverUsecase(
		cfg,
		userRepo,
		sessionRepo,
		tokenRepo,
		mailer,
		accountJWT,
	)

	// Handlers
	accountHandler := accountGrpc.NewAccountHandler(
		accountUsecase,
		sessionUsecase,
		accountRecoverUsecase,
	)

	// Setup and Start Server
	grpcSrv := server.NewGRPCServer(cfg.GRPCPort, apiKeyVerifier)
	accountV1.RegisterAccountServiceServer(grpcSrv.Engine, accountHandler)

	return grpcSrv.Start()
}
