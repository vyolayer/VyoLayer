package interceptor

import (
	"context"

	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DeviceInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		ip := md.Get("ip_address")[0]
		userAgent := md.Get("user_agent")[0]

		ctx = ctxutil.InjectDeviceInfo(ctx, ip, userAgent)
		return handler(ctx, req)
	}
}
