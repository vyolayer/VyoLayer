package ctxutil

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	vyoUserIDKey    = "vyo_user_id"
	vyoProjectIDKey = "vyo_project_id"
)

func InjectVyoServiceAccountDetails(ctx context.Context, userID, projectID uuid.UUID) context.Context {
	ctx = context.WithValue(ctx, vyoUserIDKey, userID)
	ctx = context.WithValue(ctx, vyoProjectIDKey, projectID)
	return ctx
}

// ProjectID, UserID
func ExtractVyoServiceAccountDetails(ctx context.Context) (uuid.UUID, uuid.UUID, error) {
	projectID, _ := ctx.Value(vyoProjectIDKey).(uuid.UUID)
	userID, _ := ctx.Value(vyoUserIDKey).(uuid.UUID)

	if projectID == uuid.Nil || userID == uuid.Nil {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(vyoProjectIDKey); len(vals) > 0 && projectID == uuid.Nil {
				projectID, _ = uuid.Parse(vals[0])
			}
			if vals := md.Get(vyoUserIDKey); len(vals) > 0 && userID == uuid.Nil {
				userID, _ = uuid.Parse(vals[0])
			}
		}
	}

	if projectID == uuid.Nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.Unauthenticated, "missing project ID")
	}
	if userID == uuid.Nil {
		return uuid.Nil, uuid.Nil, status.Error(codes.Unauthenticated, "missing user ID")
	}
	
	return projectID, userID, nil
}
