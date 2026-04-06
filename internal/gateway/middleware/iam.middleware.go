package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"google.golang.org/grpc/metadata"
)

func IamJWTVerify(jwt jwt.IamJWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := extractIamJWTFormCookie(c)
		if t == "" {
			t = extractIamJWTFormHeader(c)
		}

		if t == "" {
			return errors.Unauthorized("Auth token is required")
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

func extractIamJWTFormCookie(c *fiber.Ctx) string {
	return c.Cookies("__vyo_iam_auth")
}

func extractIamJWTFormHeader(c *fiber.Ctx) string {
	str := c.Get("Authorization")
	if str == "" {
		return ""
	}

	parts := strings.Split(str, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
