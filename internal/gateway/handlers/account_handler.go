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
	r := router.Group("/account")
	log.Println("Account routes registered")

	r.Post("/sign-up", h.register)
	r.Post("/verify-email", h.verifyEmail)
	r.Post("/resend-verification-email", h.resendVerificationEmail)
	r.Post("/sign-in", h.login)

	r.Post("/sessions/refresh", h.refreshToken)

	ra := r.Group("/")
	ra.Use(middleware.AccountJWTVerify(h.accountJWT))

	ra.Post("/sign-out", h.logout)
	ra.Post("/validate", h.validateSession)

	ra.Get("/sessions", h.listSessions)
	ra.Post("/sessions/revoke", h.revokeSession)
	ra.Post("/sessions/revoke-all", h.revokeAllSessions)
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

func (h *AccountHandler) refreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	t, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || t == "" {
		return response.Error(c, errors.TokenInvalid("Not found"))
	}

	log.Println("Refreshing session")
	resp, e := h.client.RefreshSession(
		ctx,
		&accountV1.RefreshSessionRequest{RefreshToken: t},
	)
	if e != nil {
		log.Printf("Failed to refresh session: %v", e)
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Set(
		c,
		resp.AccessToken,
		resp.RefreshToken,
	); err != nil {
		return response.Error(c, errors.Internal("Failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Session refreshed successfully",
		resp,
	)
}

func (h *AccountHandler) listSessions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	rt, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || rt == "" {
		return response.Error(c, errors.TokenInvalid("Not found"))
	}

	req := &accountV1.AllSessionsRequest{RefreshToken: rt}

	resp, e := h.client.AllSessions(ctx, req)
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"Sessions listed successfully",
		resp,
	)
}

func (h *AccountHandler) revokeSession(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	// Session id from params
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return response.Error(c, errors.BadRequest("Session id is required"))
	}

	refreshToken, err := h.cookieSv.GetRefreshToken(c)
	if err != nil {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.RevokeSession(ctx, &accountV1.RevokeSessionRequest{
		RefreshToken: refreshToken,
		SessionId:    sessionID,
	})
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	return response.SuccessMessage(
		c,
		"Session revoked successfully",
	)
}

func (h *AccountHandler) revokeAllSessions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	refreshToken, err := h.cookieSv.GetRefreshToken(c)
	if err != nil || refreshToken == "" {
		return response.Error(c, errors.BadRequest("Refresh token not found"))
	}

	_, e := h.client.RevokeAllSessions(ctx, &accountV1.RevokeAllSessionsRequest{
		RefreshToken: refreshToken,
	})
	if e != nil {
		return response.Error(c, errors.FromGRPC(e))
	}

	if err := h.cookieSv.Clear(c); err != nil {
		log.Printf("Failed to delete refresh token cookie: %v", err)
	}

	return response.SuccessMessage(
		c,
		"All sessions revoked successfully",
	)
}
