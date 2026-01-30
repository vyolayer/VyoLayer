package controller

import (
	"log"
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/service"
	"worklayer/internal/utils/constant"
	"worklayer/internal/utils/response"
	"worklayer/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	RegisterUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	RefreshSession(ctx *fiber.Ctx) error
	LogoutUser(ctx *fiber.Ctx) error
	ValidateSession(ctx *fiber.Ctx) error
}

type authController struct {
	authService    service.AuthService
	tokenService   service.TokenService
	sessionService service.SessionService
}

func NewAuthController(authService service.AuthService, tokenService service.TokenService, sessionService service.SessionService) *authController {
	return &authController{
		authService:    authService,
		tokenService:   tokenService,
		sessionService: sessionService,
	}
}

func (ac *authController) RegisterUser(ctx *fiber.Ctx) error {
	cmd := &dto.RegisterUserDTO{}
	if err := ctx.BodyParser(cmd); err != nil {
		return Error(
			ctx,
			response.BadRequestError("Invalid request body"),
		)
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		return Error(
			ctx,
			response.ValidationError("Validation failed", errs),
		)
	}

	if err := ac.authService.RegisterUser(ctx, *cmd); err != nil {
		return Error(ctx, err.Error())
	}

	return SuccessMessage(
		ctx,
		fiber.StatusOK,
		"User registered successfully",
	)
}

func (ac *authController) LoginUser(ctx *fiber.Ctx) error {
	cmd := &dto.LoginUserDTO{}
	if err := ctx.BodyParser(cmd); err != nil {
		return Error(
			ctx,
			response.BadRequestError("Invalid request body"),
		)
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		return Error(
			ctx,
			response.ValidationError("Validation failed", errs),
		)
	}

	user, err := ac.authService.LoginUser(ctx, *cmd)
	if err != nil {
		return Error(ctx, err.Error())
	}

	accessToken, err := ac.tokenService.GenerateAccessToken(user)
	if err != nil {
		return Error(ctx, response.InternalServerError("Failed to generate access token"))
	}

	refreshToken, err := ac.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return Error(ctx, response.InternalServerError("Failed to generate refresh token"))
	}

	if err := ac.sessionService.SaveSession(
		ctx,
		user.ID,
		refreshToken,
		ac.tokenService.GetRefreshTokenExpiry(),
	); err != nil {
		return Error(ctx, response.InternalServerError("Failed to save session"))
	}

	ac.setAuthCookies(ctx, accessToken, refreshToken)
	response := dto.LoginUserResponseDTO{
		Tokens: dto.TokenResponseDTO{AccessToken: accessToken, RefreshToken: refreshToken},
		User:   *user,
	}
	return Success(
		ctx,
		fiber.StatusOK,
		"User logged in successfully",
		response,
	)
}

func (ac *authController) LogoutUser(ctx *fiber.Ctx) error {
	localUserID := ctx.Locals("user_id")
	if localUserID == nil || localUserID.(uint) == 0 {
		return Error(ctx, response.UnauthorizedError("Authorization failed"))
	}

	// get refresh token from cookie
	refreshToken := ctx.Cookies(constant.RefreshTokenCookieName)
	if refreshToken == "" {
		return Error(ctx, response.UnauthorizedError("Refresh token not found"))
	}

	// delete session
	if err := ac.sessionService.DeleteSessionByToken(ctx, refreshToken); err != nil {
		return Error(ctx, response.InternalServerError("Failed to delete session"))
	}

	// clear cookies
	ac.clearAuthCookies(ctx)
	return SuccessMessage(
		ctx,
		fiber.StatusNoContent,
		"User logged out successfully",
	)
}

func (a *authController) RefreshSession(ctx *fiber.Ctx) error {
	// get refresh token from cookie
	oldRefreshToken := ctx.Cookies(constant.RefreshTokenCookieName)
	if oldRefreshToken == "" {
		return Error(ctx, response.UnauthorizedError("Refresh token not found"))
	}

	// validate refresh token
	userID, err := a.tokenService.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return Error(ctx, response.UnauthorizedError("Refresh token is invalid"))
	}

	// generate new access and refresh tokens
	newRefreshToken, tokenErr := a.tokenService.GenerateRefreshToken(&dto.UserDTO{ID: userID})
	if tokenErr != nil {
		return Error(ctx, response.InternalServerError("Failed to generate refresh token"))
	}
	user, err := a.sessionService.RotateSession(ctx, userID, oldRefreshToken, newRefreshToken, a.tokenService.GetRefreshTokenExpiry())
	if err != nil {
		log.Printf("REFRESH TOKEN CONTROLLER :: RotateSession : %v", err.Error())
		return Error(ctx, response.InternalServerError("Failed to rotate session"))
	}
	accessToken, tokenErr := a.tokenService.GenerateAccessToken(&dto.UserDTO{ID: user.ID, Email: user.Email})
	if tokenErr != nil {
		log.Printf("REFRESH TOKEN CONTROLLER :: GenerateAccessToken : %v", tokenErr.Error())
		return Error(ctx, response.InternalServerError("Failed to generate access token"))
	}
	// set new access and refresh tokens as cookies
	a.setAuthCookies(ctx, accessToken, newRefreshToken)

	return Success(
		ctx,
		fiber.StatusOK,
		"Session refreshed successfully",
		dto.TokenResponseDTO{AccessToken: accessToken, RefreshToken: newRefreshToken},
	)
}

func (a *authController) ValidateSession(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")
	if userID == nil {
		return Error(ctx, response.UnauthorizedError("Unauthorized"))
	}

	return SuccessMessage(
		ctx,
		fiber.StatusOK,
		"Session validated successfully",
	)
}

// setAuthCookies sets the access and refresh tokens as HTTP cookies
func (a *authController) setAuthCookies(ctx *fiber.Ctx, accessToken, refreshToken string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     constant.AccessTokenCookieName,
		Value:    accessToken,
		Expires:  time.Now().Add(a.tokenService.GetAccessTokenExpiry()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     constant.RefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(a.tokenService.GetRefreshTokenExpiry()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func (a *authController) clearAuthCookies(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     constant.AccessTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     constant.RefreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}
