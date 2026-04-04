package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tenantrepo "github.com/vyolayer/vyolayer/internal/tenant/repo"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

// OrgBoundRequest matches any generated Protobuf struct that has an OrganizationId
type OrgBoundRequest interface {
	GetOrganizationId() string
}

// RequirePermissionInterceptor evaluates schema-defined PBAC rules dynamically
func RequirePermissionInterceptor(logger *logger.AppLogger, checker tenantrepo.PermissionChecker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// Read Protobuf Schema config (Uses the reflection cache we built earlier)
		authConfig := getRouteAuthConfig(info.FullMethod)

		// Immediate Bypass
		if authConfig.SkipOrgCheck {
			return handler(ctx, req)
		}

		// Extract the User ID (Ensured by the AuthMetadataInterceptor)
		userID, err := ctxutil.ExtractIAMUserUUID(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required for this route")
		}

		// Extract Organization ID from the Protobuf Request
		orgReq, ok := req.(OrgBoundRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "route %s requires organization_id but request struct does not contain it", info.FullMethod)
		}

		orgIDStr := orgReq.GetOrganizationId()
		if orgIDStr == "" {
			return nil, status.Error(codes.InvalidArgument, "organization_id is required")
		}

		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid organization_id format")
		}

		// Evaluate Specific Granular Permission
		if authConfig.RequiredOrgPermission != "" {
			hasPerm, err := checker.HasPermission(ctx, orgID, userID, authConfig.RequiredOrgPermission)
			if err != nil {
				// Don't leak database errors to the client
				return nil, status.Error(codes.Internal, "failed to verify permissions")
			}
			if !hasPerm {
				return nil, status.Errorf(codes.PermissionDenied, "missing required permission: %s", authConfig.RequiredOrgPermission)
			}

			return handler(ctx, req)
		}

		// Fallback: Basic Active Membership Check
		isMember, err := checker.IsMember(ctx, orgID, userID)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to verify organization membership")
		}
		if !isMember {
			return nil, status.Error(codes.PermissionDenied, "you are not an active member of this organization")
		}

		// Passed all checks!
		return handler(ctx, req)
	}
}
