package controller

import (
	"log"
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
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

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account with email, password, and full name.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterUserSchema true "Registration details"
// @Success 200 {object} response.Response{data=dto.UserDTO} "User registered successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation failed"
// @Router /auth/register [post]
func (ac *authController) RegisterUser(ctx *fiber.Ctx) error {
	cmd := &dto.RegisterUserSchema{}
	if err := ctx.BodyParser(cmd); err != nil {
		return Error(
			ctx,
			InvalidBodyError,
		)
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		return Error(
			ctx,
			response.ValidationError("Validation failed", errs),
		)
	}

	user, authErr := ac.authService.RegisterUser(ctx, *cmd)
	if authErr != nil {
		return Error(ctx, response.NewErrorMessage(authErr.Code, authErr.Message))
	}

	return Success(
		ctx,
		fiber.StatusOK,
		"User registered successfully",
		dto.FromDomainUser(user),
	)
}

// LoginUser godoc
// @Summary Login a user
// @Description Authenticate a user and return access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginUserSchema true "Login credentials"
// @Success 200 {object} response.Response{data=dto.LoginUserResponseDTO} "User logged in successfully"
// @Failure 401 {object} response.ErrorResponse "Invalid email or password"
// @Router /auth/login [post]
func (ac *authController) LoginUser(ctx *fiber.Ctx) error {
	cmd := &dto.LoginUserSchema{}
	if err := ctx.BodyParser(cmd); err != nil {
		return Error(
			ctx,
			InvalidBodyError,
		)
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		return Error(
			ctx,
			response.ValidationError("Validation failed", errs),
		)
	}

	user, authErr := ac.authService.LoginUser(ctx, *cmd)
	if authErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(authErr.Code, authErr.Message),
		)
	}

	accessToken, tokenErr := ac.tokenService.GenerateAccessToken(*user)
	if tokenErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(tokenErr.Code, tokenErr.Message),
		)
	}

	refreshToken, tokenErr := ac.tokenService.GenerateRefreshToken(user.ID.InternalID().String())
	if tokenErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(tokenErr.Code, tokenErr.Message),
		)
	}

	if sessionErr := ac.sessionService.SaveSession(
		ctx,
		user.ID,
		refreshToken,
		ac.tokenService.GetRefreshTokenExpiry(),
	); sessionErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(sessionErr.Code, sessionErr.Message),
		)
	}

	ac.setAuthCookies(ctx, accessToken, refreshToken)
	response := dto.LoginUserResponseDTO{
		TokenResponseDTO: dto.TokenResponseDTO{AccessToken: accessToken, RefreshToken: refreshToken},
		User:             dto.FromDomainUser(user),
	}
	return Success(
		ctx,
		fiber.StatusOK,
		"User logged in successfully",
		response,
	)
}

// LogoutUser godoc
// @Summary Logout a user
// @Description Invalidate the user's session and clear authentication cookies.
// @Tags auth
// @Produce json
// @Success 204 "User logged out successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (ac *authController) LogoutUser(ctx *fiber.Ctx) error {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return Error(ctx, response.NewErrorMessage(fiber.StatusUnauthorized, err.Error()))
	}

	// get refresh token from cookie
	refreshToken := ctx.Cookies(RefreshTokenCookieName)
	if refreshToken == "" {
		return Error(ctx, response.NewErrorMessage(fiber.StatusUnauthorized, "Refresh token not found"))
	}

	// delete session
	if err := ac.sessionService.DeleteSessionByToken(ctx, refreshToken); err != nil {
		return Error(
			ctx,
			response.NewErrorMessage(fiber.StatusInternalServerError, "Failed to delete session"),
		)
	}

	// clear cookies
	ac.clearAuthCookies(ctx)
	return SuccessMessage(
		ctx,
		fiber.StatusNoContent,
		"User logged out successfully",
	)
}

// RefreshSession godoc
// @Summary Refresh session
// @Description Get a new access token using a valid refresh token from cookies.
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response{data=dto.RefreshSessionResponseDTO} "Session refreshed successfully"
// @Failure 401 {object} response.ErrorResponse "Invalid or missing refresh token"
// @Router /auth/refresh [post]
func (a *authController) RefreshSession(ctx *fiber.Ctx) error {
	// get refresh token from cookie
	oldRefreshToken := ctx.Cookies(RefreshTokenCookieName)
	if oldRefreshToken == "" {
		return Error(
			ctx,
			response.NewErrorMessage(fiber.StatusUnauthorized, "Refresh token not found"),
		)
	}

	// validate refresh token
	jwtUserID, tokenErr := a.tokenService.ValidateRefreshToken(oldRefreshToken)
	if tokenErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(tokenErr.Code, tokenErr.Message),
		)
	}

	userID, err := types.ReconstructUserID(jwtUserID)
	log.Printf("REFRESH TOKEN CONTROLLER :: userID : %v", userID.InternalID().String())
	if err != nil {
		return Error(
			ctx,
			response.NewErrorMessage(fiber.StatusInternalServerError, "Failed to reconstruct user ID"),
		)
	}

	newRefreshToken, tokenErr := a.tokenService.GenerateRefreshToken(userID.InternalID().String())
	if tokenErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(tokenErr.Code, tokenErr.Message),
		)
	}

	user, sessionErr := a.sessionService.RotateSession(
		ctx,
		*userID,
		oldRefreshToken,
		newRefreshToken,
		a.tokenService.GetRefreshTokenExpiry(),
	)
	if sessionErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(sessionErr.Code, sessionErr.Message),
		)
	}

	accessToken, tokenErr := a.tokenService.GenerateAccessToken(*user)
	if tokenErr != nil {
		return Error(
			ctx,
			response.NewErrorMessage(tokenErr.Code, tokenErr.Message),
		)
	}

	// set new access and refresh tokens as cookies
	a.setAuthCookies(ctx, accessToken, newRefreshToken)
	responseData := dto.RefreshSessionResponseDTO{
		TokenResponseDTO: dto.TokenResponseDTO{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		},
	}
	return Success(
		ctx,
		fiber.StatusOK,
		"Session refreshed successfully",
		responseData,
	)
}

// ValidateSession godoc
// @Summary Validate session
// @Description Check if the current session (access token) is still valid.
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response "Session validated successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /auth/validate [post]
func (a *authController) ValidateSession(ctx *fiber.Ctx) error {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return Error(ctx, response.NewErrorMessage(fiber.StatusUnauthorized, err.Error()))
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
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Expires:  time.Now().Add(a.tokenService.GetAccessTokenExpiry()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(a.tokenService.GetRefreshTokenExpiry()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func (a *authController) clearAuthCookies(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})
}
