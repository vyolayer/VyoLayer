package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vyolayer/vyolayer/pkg/ctxutil"
)

// AuthMetadataInterceptor reads the user ID passed by the API Gateway
func AuthMetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// 1. Extract metadata from the incoming gRPC request
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// If no metadata exists, we let it pass. The PBAC interceptor will
			// catch unauthorized access later if the route requires a user.
			return handler(ctx, req)
		}

		// 2. Look for the specific header your Gateway sends (e.g., "x-user-id")
		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 || userIDs[0] == "" {
			return handler(ctx, req)
		}

		// 3. Parse and validate the UUID
		userID, err := uuid.Parse(userIDs[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid user identity format provided by gateway")
		}

		// 4. Inject it into the Go context for the rest of the lifecycle
		newCtx := ctxutil.InjectIAMUserID(ctx, userID.String())

		// 5. Pass the NEW context down the chain
		return handler(newCtx, req)
	}
}
