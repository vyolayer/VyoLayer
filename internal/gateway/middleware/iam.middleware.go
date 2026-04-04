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

		user, err := jwt.VerifyAccessToken(t)
		if err != nil {
			return errors.Unauthorized("invalid or expired auth token")
		}

		log.Printf("[GATEWAY - IAM] (Middleware - IAM) User ID :: %s", user.UserID.String())
		ctx := ctxutil.InjectIAMUserID(c.UserContext(), user.UserID.String())
		ctx = ctxutil.InjectIAMUserEmail(ctx, user.Email)
		ctx = metadata.AppendToOutgoingContext(ctx,
			"iam_user_id", user.UserID.String(),
			"iam_user_email", user.Email,
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
