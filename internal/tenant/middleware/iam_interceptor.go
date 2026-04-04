package middleware

import (
	"context"

	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func IamInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("iam_user_id"); len(vals) > 0 && vals[0] != "" {
				ctx = ctxutil.InjectIAMUserID(ctx, vals[0])
			}
		}
		return handler(ctx, req)
	}
}
