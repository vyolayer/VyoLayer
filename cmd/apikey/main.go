package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vyolayer/vyolayer/internal/apikeys/config"
	"github.com/vyolayer/vyolayer/internal/apikeys/grpc"
	"github.com/vyolayer/vyolayer/internal/apikeys/repository"
	"github.com/vyolayer/vyolayer/internal/apikeys/server"
	"github.com/vyolayer/vyolayer/internal/apikeys/service"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/postgres"

	pb "github.com/vyolayer/vyolayer/proto/apikey/v1"
)

func main() {
	logger := logger.NewAppLogger("APIKEY")

	if err := godotenv.Load(); err != nil {
		logger.Warn("Note: No .env file found; relying on system environment variables", "")
	}

	cfg := config.ApikeysConfigLoad()

	// Setup DB
	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Repo
	akRepo := repository.New(db)

	// Service
	akSvc := service.New(akRepo)

	// gRPC Server
	grpcServer := server.NewGRPCServer("50055", logger)
	apikeysHandler := grpc.New(akSvc)

	pb.RegisterAPIKeyServiceServer(grpcServer.Engine, apikeysHandler)

	if err := grpcServer.Start(); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
