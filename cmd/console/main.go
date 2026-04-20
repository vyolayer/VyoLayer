package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/console/config"
	consolegrpc "github.com/vyolayer/vyolayer/internal/console/delivery/grpc"
	"github.com/vyolayer/vyolayer/internal/console/repository"
	"github.com/vyolayer/vyolayer/internal/console/server"
	"github.com/vyolayer/vyolayer/internal/console/service"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/postgres"
	consolev1 "github.com/vyolayer/vyolayer/proto/console/v1"
)

func main() {
	logger := logger.NewAppLogger("CONSOLE")

	if err := godotenv.Load(); err != nil {
		logger.Warn("Note: No .env file found; relying on system environment variables", "")
	}

	cfg := config.ConsoleConfigLoad()

	// Setup DB
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Setup Repositories
	psRepo := repository.NewProjectServiceRepository(db)
	resRepo := repository.NewResourceRepository(db)
	overRepo := repository.NewOverrideRepository(db)

	// Setup Service
	manifestSvc := service.NewManifestService(psRepo, resRepo, overRepo)

	// Setup gRPC Server
	grpcServer := server.NewGRPCServer("50054", logger)
	manifestGrpcServer := consolegrpc.NewManifestServer(manifestSvc)
	consolev1.RegisterProjectServiceManifestServer(grpcServer.Engine, manifestGrpcServer)

	if err := grpcServer.Start(); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
