package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/internal/shared/middleware"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/response"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

// AccountHandler manages HTTP requests related to accounts
type AccountHandler struct {
	client     accountV1.AccountServiceClient
	cookieSv   *service.AccountTokenService
	accountJWT jwt.AccountJWT
}

// NewAccountHandler creates a new AccountHandler injecting the gRPC client
func NewAccountHandler(
	client accountV1.AccountServiceClient,
	cookieSv *service.AccountTokenService,
	accountJWT jwt.AccountJWT,
) *AccountHandler {
	return &AccountHandler{
		client:     client,
		cookieSv:   cookieSv,
		accountJWT: accountJWT,
	}
}

// RegisterRoutes registers the account routes on the provided router
func (h *AccountHandler) RegisterRoutes(router fiber.Router) {
	accountMiddleware := middleware.NewAccountMiddleware(h.accountJWT)

	accountGroup := router.Group("/account")
	log.Println("Account routes registered")

	accountGroup.Post("/sign-up", h.register)
	accountGroup.Post("/verify-email", h.verifyEmail)
	accountGroup.Post("/resend-verification-email", h.resendVerificationEmail)
	accountGroup.Post("/sign-in", h.login)

	// accountGroup.Use(accountMiddleware.VerifyAccessToken())

	accountGroup.Post("/sign-out", accountMiddleware.VerifyAccessToken(), h.logout)
	accountGroup.Post("/validate", h.validateSession)
}

// register handles user registration
func (h *AccountHandler) register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	var req accountV1.RegisterRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	resp, e := h.client.Register(ctx, &req)
	if e != nil {
		appErr := errors.FromGRPC(e)
		return response.Error(c, appErr)
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusCreated,
		"User registered successfully",
		resp,
	)
}

func (h *AccountHandler) verifyEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	token := c.Query("token")
	if token == "" {
		return response.Error(c, errors.BadRequest("Token is required"))
	}

	_, e := h.client.VerifyEmail(ctx, &accountV1.VerifyEmailRequest{
		Token: token,
	})
	if e != nil {
		appErr := errors.FromGRPC(e)
		return response.Error(c, appErr)
	}

	return response.SuccessMessage(
		c,
		"Email verified successfully",
	)
}

func (h *AccountHandler) resendVerificationEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	var req accountV1.ResendVerificationEmailRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	_, e := h.client.ResendVerificationEmail(ctx, &req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessMessage(
		c,
		"Verification email resent successfully",
	)
}

func (h *AccountHandler) login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	var req accountV1.LoginRequest
	if e := c.BodyParser(&req); e != nil {
		return response.Error(c, errors.BadRequest("Invalid Request Body"))
	}

	resp, e := h.client.Login(ctx, &req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Set(c, resp.AccessToken, resp.RefreshToken); err != nil {
		return response.Error(c, errors.Internal("Failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Login successful",
		resp,
	)
}

func (h *AccountHandler) logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	t, err := h.cookieSv.GetRefreshToken(c)
	if err != nil {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.Logout(ctx, &accountV1.LogoutRequest{
		RefreshToken: t,
	})

	// Always clear the local cookies so the user isn't stuck holding invalid tokens
	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessMessage(
		c,
		"Logout successful",
	)
}

func (h *AccountHandler) validateSession(c *fiber.Ctx) error {
	// Extract the access token from the cookies
	accessToken := h.cookieSv.GetAccessToken(c)
	if accessToken == "" {
		return response.Error(c, errors.Unauthorized("No access token provided"))
	}

	// Verify the access token securely using the shared JWT verifier
	_, _, err := h.accountJWT.VerifyAccessToken(accessToken)
	if err != nil {
		log.Printf("Session validation failed: %v", err)
		return response.Error(c, errors.Unauthorized("Invalid or expired session"))
	}

	log.Println("Session validated successfully")

	return response.SuccessMessage(
		c,
		"true",
	)
}
