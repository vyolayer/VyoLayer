package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/iam/config"
	"github.com/vyolayer/vyolayer/internal/iam/repository"
	"github.com/vyolayer/vyolayer/internal/iam/server"
	"github.com/vyolayer/vyolayer/internal/iam/usecase"
	"github.com/vyolayer/vyolayer/internal/shared/session"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/mail"
	"github.com/vyolayer/vyolayer/pkg/postgres"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"

	iamGrpc "github.com/vyolayer/vyolayer/internal/iam/delivery/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("iam service error: %v", err)
	}
}

func run() error {
	appLogger := logger.NewAppLogger("IAM-SERVICE")

	// Load Environment Variables
	if err := godotenv.Load(); err != nil {
		appLogger.Error("Note: No .env file found; relying on system environment variables", err.Error())
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

	// ── Repositories ─────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db, appLogger.WithContext("USER REPO"))
	sessionRepo := repository.NewSessionRepository(db, cfg.JWT.RefreshTokenExpiry)
	verificationTokenRepo := repository.NewVerificationTokenRepository(db, cfg.Token.VerificationTokenExpiry)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db, cfg.Token.PasswordResetTokenExpiry)

	// ── Session service ───────────────────────────────────────────────────────
	sessionService := session.NewIAMSession(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// ── Mailer ────────────────────────────────────────────────────────────────
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

	// ── Use-cases ─────────────────────────────────────────────────────────────
	authUsecase := usecase.NewAuthUsecase(
		appLogger.WithContext("AUTH UC"),
		userRepo,
		sessionRepo,
		verificationTokenRepo,
		passwordResetTokenRepo,
		sessionService,
		mailer,
		cfg.GRPC.AppURL,
	)

	userUsecase := usecase.NewUserUsecase(
		appLogger.WithContext("USER UC"),
		userRepo,
	)

	// ── gRPC handlers ─────────────────────────────────────────────────────────
	iamAuthH := iamGrpc.NewIAMAuthHandler(authUsecase)
	iamUserH := iamGrpc.NewIAMUserHandler(userUsecase)

	// ── Server ────────────────────────────────────────────────────────────────
	grpcSrv := server.NewGRPCServer(cfg.GRPC.GRPCPort, appLogger)
	iAMV1.RegisterAuthServiceServer(grpcSrv.Engine, iamAuthH)
	iAMV1.RegisterUserServiceServer(grpcSrv.Engine, iamUserH)

	return grpcSrv.Start()
}
