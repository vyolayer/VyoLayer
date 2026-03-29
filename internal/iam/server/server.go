package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/interceptor"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Engine *grpc.Server
	port   string
	logger *logger.AppLogger
}

func NewGRPCServer(port string, logger *logger.AppLogger) *GRPCServer {
	v, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to initialize validator: %v", err)
	}

	srv := grpc.NewServer(
		grpc.ConnectionTimeout(120*time.Second),
		grpc.ChainUnaryInterceptor(
			interceptor.GRPCValidationInterceptor(v),
			interceptor.DeviceInterceptor(),
			iamInterceptor(),
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
		logger: logger,
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

	s.logger.Info("IAM service started successfully", map[string]any{"Port": s.port})
	if err := s.Engine.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// iamInterceptor propagates the caller's IAM user ID from incoming gRPC metadata
// into the request context. It is a pass-through for unauthenticated RPCs; the
// individual use-case methods are responsible for enforcing authentication.
func iamInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("iam_user_id"); len(vals) > 0 && vals[0] != "" {
				ctx = ctxutil.InjectIAMUserID(ctx, vals[0])
			}
		}
		return handler(ctx, req)
	}
}
