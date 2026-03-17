package ctxutil

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	ProjectIDKey      contextKey = "project_id"
	OrganizationIDKey contextKey = "organization_id"
	APIKeyHashKey     contextKey = "api_key_hash"
)

func InjectProjectID(ctx context.Context, projectID uuid.UUID) context.Context {
	return context.WithValue(ctx, ProjectIDKey, projectID)
}

func InjectOrganizationID(ctx context.Context, organizationID uuid.UUID) context.Context {
	return context.WithValue(ctx, OrganizationIDKey, organizationID)
}

func InjectAPIKeyHash(ctx context.Context, apiKeyHash string) context.Context {
	return context.WithValue(ctx, APIKeyHashKey, apiKeyHash)
}

func ExtractProjectID(ctx context.Context) (uuid.UUID, error) {
	val, ok := ctx.Value(ProjectIDKey).(uuid.UUID)
	if !ok || val == uuid.Nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing project ID")
	}
	return val, nil
}

func ExtractOrganizationID(ctx context.Context) (uuid.UUID, error) {
	val, ok := ctx.Value(OrganizationIDKey).(uuid.UUID)
	if !ok || val == uuid.Nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing organization ID")
	}
	return val, nil
}

func ExtractAPIKeyHash(ctx context.Context) (string, error) {
	val, ok := ctx.Value(APIKeyHashKey).(string)
	if !ok || val == "" {
		return "", status.Error(codes.Unauthenticated, "missing project API key")
	}
	return val, nil
}
