package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	"github.com/vyolayer/vyolayer/pkg/interceptor"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Engine *grpc.Server
	port   string
	logger *logger.AppLogger
}

func NewGRPCServer(port string, logger *logger.AppLogger) *GRPCServer {
	grpcLogger := logger.WithContext("GRPC")

	v, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to initialize validator: %v", err)
	}

	// panicHandler := func(p any) (err error) {
	// 	// Log the stack trace so you can fix the bug later!
	// 	log.Printf("PANIC RECOVERED: %v\n%s", p, string(debug.Stack()))
	// 	return status.Errorf(codes.Internal, "an unexpected internal server error occurred")
	// }

	srv := grpc.NewServer(
		grpc.ConnectionTimeout(120*time.Second),
		grpc.ChainUnaryInterceptor(
			interceptor.GRPCValidationInterceptor(v),
		),
	)

	// Register health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for grpcurl
	reflection.Register(srv)

	return &GRPCServer{
		Engine: srv,
		port:   port,
		logger: grpcLogger,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	// Graceful shutdown handling
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		s.logger.Warn("Shutting down gRPC server...", "")
		s.Engine.GracefulStop()
	}()

	s.logger.Info("Console service started successfully", map[string]any{"Port": s.port})
	if err := s.Engine.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
