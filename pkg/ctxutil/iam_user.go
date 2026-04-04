package ctxutil

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	iamUserIDLegacyKey = "iam_user_id"
	iamUserEmailKey    = "iam_user_email"
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

func ExtractIAMUserUUID(ctx context.Context) (uuid.UUID, error) {
	raw, err := ExtractIAMUserID(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	uid, parseErr := uuid.Parse(raw)
	if parseErr != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user id in context")
	}
	if uid == uuid.Nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user id in context")
	}

	return uid, nil
}

// Inject Email
func InjectIAMUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, contextKey(iamUserEmailKey), email)
}

// Extract email
func ExtractIAMUserEmail(ctx context.Context) (string, error) {
	var raw string

	// Try context value first (in-process injection)
	if val, ok := ctx.Value(contextKey(iamUserEmailKey)).(string); ok && val != "" {
		raw = val
	}

	// Fall back to gRPC incoming metadata
	if raw == "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(iamUserEmailKey); len(vals) > 0 {
				raw = vals[0]
			}
		}
	}

	if raw == "" {
		return "", status.Error(codes.Unauthenticated, "missing user email in context or metadata")
	}

	return raw, nil
}
