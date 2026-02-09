package controller

import (
	"log"
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/internal/utils/validation"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

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
// @Success 200 {object} response.SuccessResponse{data=dto.UserDTO}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (ac *authController) RegisterUser(ctx *fiber.Ctx) error {
	cmd := &dto.RegisterUserSchema{}
	if err := ctx.BodyParser(cmd); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		validationErr := errors.ValidationFailed("Validation failed")
		validationErr.WithMetadata("validation_errors", errs)
		return response.Error(ctx, validationErr)
	}

	user, authErr := ac.authService.RegisterUser(ctx, *cmd)
	if authErr != nil {
		return response.Error(ctx, authErr)
	}

	return response.SuccessWithMessage(
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
// @Success 200 {object} response.SuccessResponse{data=dto.LoginUserResponseDTO}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/login [post]
func (ac *authController) LoginUser(ctx *fiber.Ctx) error {
	cmd := &dto.LoginUserSchema{}
	if err := ctx.BodyParser(cmd); err != nil {
		return response.Error(ctx, errors.BadRequest("Invalid request body"))
	}

	// Validate the command
	if errs := validation.ValidateStruct(cmd); len(errs) > 0 {
		validationErr := errors.ValidationFailed("Validation failed")
		validationErr.WithMetadata("validation_errors", errs)
		return response.Error(ctx, validationErr)
	}

	user, authErr := ac.authService.LoginUser(ctx, *cmd)
	if authErr != nil {
		return response.Error(ctx, authErr)
	}

	accessToken, tokenErr := ac.tokenService.GenerateAccessToken(*user)
	if tokenErr != nil {
		return response.Error(ctx, tokenErr)
	}

	refreshToken, tokenErr := ac.tokenService.GenerateRefreshToken(user.ID.InternalID().String())
	if tokenErr != nil {
		return response.Error(ctx, tokenErr)
	}

	if sessionErr := ac.sessionService.SaveSession(
		ctx,
		user.ID,
		refreshToken,
		ac.tokenService.GetRefreshTokenExpiry(),
	); sessionErr != nil {
		return response.Error(ctx, sessionErr)
	}

	ac.setAuthCookies(ctx, accessToken, refreshToken)
	responseData := dto.LoginUserResponseDTO{
		TokenResponseDTO: dto.TokenResponseDTO{AccessToken: accessToken, RefreshToken: refreshToken},
		User:             dto.FromDomainUser(user),
	}
	return response.SuccessWithMessage(
		ctx,
		fiber.StatusOK,
		"User logged in successfully",
		responseData,
	)
}

// LogoutUser godoc
// @Summary Logout a user
// @Description Invalidate the user's session and clear authentication cookies.
// @Tags auth
// @Produce json
// @Success 204 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (ac *authController) LogoutUser(ctx *fiber.Ctx) error {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return response.Error(ctx, errors.Unauthorized(err.Error()))
	}

	// get refresh token from cookie
	refreshToken := ctx.Cookies(RefreshTokenCookieName)
	if refreshToken == "" {
		return response.Error(ctx, errors.Unauthorized("Refresh token not found"))
	}

	// delete session
	if err := ac.sessionService.DeleteSessionByToken(ctx, refreshToken); err != nil {
		return response.Error(ctx, err)
	}

	// clear cookies
	ac.clearAuthCookies(ctx)
	return response.NoContent(ctx)
}

// RefreshSession godoc
// @Summary Refresh session
// @Description Get a new access token using a valid refresh token from cookies.
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=dto.RefreshSessionResponseDTO}
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (a *authController) RefreshSession(ctx *fiber.Ctx) error {
	// get refresh token from cookie
	oldRefreshToken := ctx.Cookies(RefreshTokenCookieName)
	if oldRefreshToken == "" {
		return response.Error(ctx, errors.Unauthorized("Refresh token not found"))
	}

	// validate refresh token
	jwtUserID, tokenErr := a.tokenService.ValidateRefreshToken(oldRefreshToken)
	if tokenErr != nil {
		return response.Error(ctx, tokenErr)
	}

	userID, err := types.ReconstructUserID(jwtUserID)
	log.Printf("REFRESH TOKEN CONTROLLER :: userID : %v", userID.InternalID().String())
	if err != nil {
		return response.Error(ctx, errors.Internal("Failed to reconstruct user ID"))
	}

	newRefreshToken, tokenErr := a.tokenService.GenerateRefreshToken(userID.InternalID().String())
	if tokenErr != nil {
		return response.Error(ctx, tokenErr)
	}

	user, sessionErr := a.sessionService.RotateSession(
		ctx,
		userID,
		oldRefreshToken,
		newRefreshToken,
		a.tokenService.GetRefreshTokenExpiry(),
	)
	if sessionErr != nil {
		return response.Error(ctx, sessionErr)
	}

	accessToken, tokenErr := a.tokenService.GenerateAccessToken(*user)
	if tokenErr != nil {
		return response.Error(ctx, tokenErr)
	}

	// set new access and refresh tokens as cookies
	a.setAuthCookies(ctx, accessToken, newRefreshToken)
	responseData := dto.RefreshSessionResponseDTO{
		TokenResponseDTO: dto.TokenResponseDTO{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		},
	}
	return response.SuccessWithMessage(
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
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/validate [post]
func (a *authController) ValidateSession(ctx *fiber.Ctx) error {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return response.Error(ctx, errors.Unauthorized(err.Error()))
	}

	return response.SuccessMessage(ctx, "Session validated successfully")
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

func getUserIDFromContext(ctx *fiber.Ctx) (types.UserID, error) {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user context")
	}
	return localUserID, nil
}
