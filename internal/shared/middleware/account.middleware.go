package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"google.golang.org/grpc/metadata"
)

type accountMiddleware struct {
	accountJWT jwt.AccountJWT
}

func NewAccountMiddleware(accountJWT jwt.AccountJWT) *accountMiddleware {
	return &accountMiddleware{
		accountJWT: accountJWT,
	}
}

func AccountJWTVerify(aJWT jwt.AccountJWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cookie or Header
		token := c.Cookies("vyo_session")
		if token == "" {
			token = c.Get("Authorization")
		}

		if token == "" {
			return errors.Unauthorized("Authorization token is required")
		}

		userID, projectID, err := aJWT.VerifyAccessToken(token)
		if err != nil {
			return err
		}

		ctx := ctxutil.InjectVyoServiceAccountDetails(c.UserContext(), userID, projectID)

		// Ensure user ID is passed to gRPC backends by appending to outgoing context metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "vyo_user_id", userID.String())

		c.SetUserContext(ctx)

		return c.Next()
	}
}
