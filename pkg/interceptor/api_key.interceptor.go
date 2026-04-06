package interceptor

import (
	"context"

	"github.com/google/uuid"
	gtmw "github.com/vyolayer/vyolayer/internal/gateway/middleware"

	// apikey "github.com/vyolayer/vyolayer/internal/shared/api-key"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type APIKeyVerifier interface {
	Verify(apiKey string, projectID uuid.UUID) error
}

// Get API Key from request and validate it
// Fixme: verify project ID
func APIKeyInterceptor(verifier APIKeyVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Skip authentication for Health Check and Reflection services
		if info.FullMethod == "/grpc.health.v1.Health/Check" ||
			info.FullMethod == "/grpc.health.v1.Health/Watch" ||
			info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
			return handler(ctx, req)
		}

		// Extract API Key from request
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		apiKey := md.Get(gtmw.ContextKeyAPIKey)
		if len(apiKey) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing API key")
		}

		projectID := md.Get(gtmw.ContextKeyProjectID)
		if len(projectID) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing project ID")
		}

		// projectUUID, err := uuid.Parse(projectID[0])
		// if err != nil {
		// 	return nil, status.Error(codes.Unauthenticated, "invalid project ID")
		// }

		// apiKeyInfo, err := verifier.Verify(apiKey[0], projectUUID)
		// if err != nil {
		// 	return nil, status.Error(codes.Unauthenticated, "invalid API key")
		// }

		// ctx = ctxutil.InjectAPIKeyInfo(ctx, apiKeyInfo)

		return handler(ctx, req)
	}
}
