package grpcutil

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

func NewClient(config ClientConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", config.Address, err)
	}

	log.Printf("gRPC client connected to %s", config.Address)
	return conn, nil
}

// UnaryClientInterceptor adds metadata to outgoing requests
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Can inject metadata here if needed
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
