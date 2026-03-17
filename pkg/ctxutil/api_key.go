package ctxutil

import (
	"context"

	apikey "github.com/vyolayer/vyolayer/internal/shared/api-key"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const APIKeyInfo contextKey = "api_key_info"

func InjectAPIKeyInfo(ctx context.Context, apiKeyInfo *apikey.APIKeyInfo) context.Context {
	return context.WithValue(ctx, APIKeyInfo, apiKeyInfo)
}

func ExtractAPIKeyInfo(ctx context.Context) (*apikey.APIKeyInfo, error) {
	val, ok := ctx.Value(APIKeyInfo).(*apikey.APIKeyInfo)
	if !ok || val == nil {
		return nil, status.Error(codes.Unauthenticated, "missing API key info")
	}
	return val, nil
}
