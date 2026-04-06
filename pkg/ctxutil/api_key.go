package ctxutil

import (
	"context"

	// apikey "github.com/vyolayer/vyolayer/internal/shared/api-key"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const APIKeyInfo contextKey = "api_key_info"

// Fixme
type TAPIKeyInfo struct {
	ProjectID uuid.UUID
}

func InjectAPIKeyInfo(ctx context.Context, apiKeyInfo TAPIKeyInfo) context.Context {
	return context.WithValue(ctx, APIKeyInfo, apiKeyInfo)
}

func ExtractAPIKeyInfo(ctx context.Context) (*TAPIKeyInfo, error) {
	val, ok := ctx.Value(APIKeyInfo).(*TAPIKeyInfo)
	if !ok || val == nil {
		return nil, status.Error(codes.Unauthenticated, "missing API key info")
	}
	return val, nil
}
