package ctxutil

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	iamUserIDLegacyKey = "iam_user_id"
)

// InjectIAMUserID injects an IAM user ID string into the context (used for testing / in-process calls).
func InjectIAMUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKey(iamUserIDLegacyKey), userID)
}

// ExtractIAMUserID extracts the authenticated user ID for use with IAM gRPC services.
// Checks context values first (in-process), then falls back to gRPC incoming metadata.
func ExtractIAMUserID(ctx context.Context) (string, error) {
	var raw string

	// Try context value first (in-process injection)
	if val, ok := ctx.Value(contextKey(iamUserIDLegacyKey)).(string); ok && val != "" {
		raw = val
	}

	// Fall back to gRPC incoming metadata
	if raw == "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(iamUserIDLegacyKey); len(vals) > 0 {
				raw = vals[0]
			}
		}
	}

	if raw == "" {
		return "", status.Error(codes.Unauthenticated, "missing user ID in context or metadata")
	}

	return raw, nil
}
