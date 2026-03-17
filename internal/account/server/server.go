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
	apikey "github.com/vyolayer/vyolayer/internal/shared/api-key"
	"github.com/vyolayer/vyolayer/pkg/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCServer represents the gRPC Server wrapped instance
type GRPCServer struct {
	Engine *grpc.Server
	port   string
}

// NewGRPCServer creates and configures a new Server instance
func NewGRPCServer(port string, apiKeyVerifier apikey.APIKeyVerifier) *GRPCServer {
	v, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to initialize validator: %v", err)
	}

	srv := grpc.NewServer(
		grpc.ConnectionTimeout(120*time.Second),
		grpc.ChainUnaryInterceptor(
			interceptor.APIKeyInterceptor(apiKeyVerifier),
			interceptor.GRPCValidationInterceptor(v),
			interceptor.DeviceInterceptor(),
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
	}
}

// Start runs the GRPC server with graceful shutdown handling
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

		log.Println("Shutting down gRPC server...")
		s.Engine.GracefulStop()
	}()

	log.Printf("Account service listening on :%s", s.port)
	if err := s.Engine.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
