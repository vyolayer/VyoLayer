package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"google.golang.org/grpc/metadata"
)

func IamJWTVerify(jwt jwt.IamJWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t, err := extractIamJWTFormCookie(c)
		if err != nil {
			return err
		}

		id, err := jwt.VerifyAccessToken(t)
		if err != nil {
			return errors.Unauthorized("invalid or expired auth token")
		}

		log.Printf("[GATEWAY - IAM] (Middleware - IAM) User ID :: %s", id)
		ctx := ctxutil.InjectIAMUserID(c.UserContext(), id.String())
		ctx = metadata.AppendToOutgoingContext(ctx,
			"iam_user_id", id.String(),
		)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func extractIamJWTFormCookie(c *fiber.Ctx) (string, error) {
	token := c.Cookies("__vyo_iam_auth")
	if token == "" {
		return "", errors.Unauthorized("Auth token is required")
	}
	return token, nil
}
