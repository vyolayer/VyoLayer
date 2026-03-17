package middleware

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"
)

const (
	headerProjectID = "x-vyo-project-id"
	headerAPIKey    = "x-vyo-key"

	ContextKeyProjectID = "vyo_project_id"
	ContextKeyAPIKey    = "vyo_api_key"
	ContextKeyIPAddress = "ip_address"
	ContextKeyUserAgent = "user_agent"
)

// GRPCMetadataMiddleware extracts headers from the HTTP request and attaches them to the Fiber
// UserContext as gRPC metadata, making it automatically available to downstream gRPC handlers.
func GRPCMetadataMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		projectID := c.Get(headerProjectID)
		if projectID != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyProjectID, projectID)
		}

		// Extract API Key from Fiber headers
		apiKey := c.Get(headerAPIKey)
		if apiKey != "" {
			// Append it to the gRPC outgoing context
			ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyAPIKey, apiKey)
		}

		// ip address and user agent
		ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyIPAddress, c.IP())
		ctx = metadata.AppendToOutgoingContext(ctx, ContextKeyUserAgent, c.Get("User-Agent"))

		// Save the updated context back to Fiber
		c.SetUserContext(ctx)

		return c.Next()
	}
}
