package ctxutil

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UserIDKey contextKey = "user_id"
)

func InjectUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func ExtractUserID(ctx context.Context) (string, error) {
	val, ok := ctx.Value(UserIDKey).(string)
	if !ok || val == "" {
		return "", status.Error(codes.Unauthenticated, "missing user ID")
	}
	return val, nil
}
