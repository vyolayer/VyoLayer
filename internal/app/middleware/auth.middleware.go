package middleware

import (
	"log"
	"strings"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/internal/utils/constant"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

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
			return response.Error(ctx, errors.Unauthorized("Missing authentication token"))
		}

		userClaims, err := am.tokenService.ValidateAccessToken(accessToken)
		if err != nil {
			return response.Error(ctx, err)
		}
		log.Printf("AUTH MIDDLEWARE :: JwtValidated : userClaims : %v", userClaims)

		// Parse user ID
		userID, parseErr := types.ReconstructUserID(userClaims.UserID)
		if parseErr != nil {
			return response.Error(ctx, errors.Unauthorized("Invalid user ID in token"))
		}

		ctx.Locals("user_id", *userID)
		ctx.Locals("user_email", userClaims.Email)

		return ctx.Next()
	}
}
