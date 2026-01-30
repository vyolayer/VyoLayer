package middleware

import (
	"strings"
	"worklayer/internal/app/controller"
	"worklayer/internal/service"
	"worklayer/internal/utils/constant"
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

var Error = controller.Error

type AuthMiddleware struct {
	tokenService service.TokenService
}

func NewAuthMiddleware(tokenService service.TokenService) AuthMiddleware {
	return AuthMiddleware{
		tokenService: tokenService,
	}
}

func (am *AuthMiddleware) JwtValidated() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var accessToken string

		// get access token from cookie
		accessToken = ctx.Cookies(constant.AccessTokenCookieName)

		if accessToken == "" {
			// get access token from header
			authHeader := ctx.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				accessToken = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if accessToken == "" {
			return Error(ctx, response.UnauthorizedError("Unauthorized"))
		}

		userClaims, err := am.tokenService.ValidateAccessToken(accessToken)
		if err != nil {
			return Error(ctx, response.UnauthorizedError("Unauthorized"))
		}

		ctx.Locals("user_id", userClaims.UserID)
		ctx.Locals("user_email", userClaims.Email)

		return ctx.Next()
	}
}
