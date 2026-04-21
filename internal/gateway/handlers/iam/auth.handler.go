package iam

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"

	iamV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

const (
	grpcTimeout = 10 * time.Second
)

type IAMAuthGatewayHandler struct {
	auth   iamV1.AuthServiceClient
	user   iamV1.UserServiceClient
	cookie *service.IAMCookieService
	iamJWT jwt.IamJWT
	logger *logger.AppLogger
}

func NewIAMAuthGatewayHandler(
	auth iamV1.AuthServiceClient,
	user iamV1.UserServiceClient,
	cookie *service.IAMCookieService,
	iamJWT jwt.IamJWT,
	logger *logger.AppLogger,
) *IAMAuthGatewayHandler {
	return &IAMAuthGatewayHandler{
		auth:   auth,
		user:   user,
		cookie: cookie,
		iamJWT: iamJWT,
		logger: logger.WithContext("IAMAuthGatewayHandler"),
	}
}

// RegisterRoutes mounts all IAM routes under /iam.
func (h *IAMAuthGatewayHandler) RegisterRoutes(router fiber.Router) {
	// grpc ctx timeout
	grpcCtxMiddleware := middleware.NewGrpcCtxMiddleware(grpcTimeout)

	strictLimiter := middleware.NewRateLimiter(10, 1*time.Minute).Handler()
	standardLimiter := middleware.NewRateLimiter(100, 1*time.Minute).Handler()

	iam := router.Group("/iam")
	iam.Use(grpcCtxMiddleware.Handler())

	// ── Public auth endpoints ────────────────────────────────────────────────
	iam.Post("/register", strictLimiter, h.register)
	iam.Post("/verify-email", standardLimiter, h.verifyEmail)
	iam.Post("/resend-verification-email", strictLimiter, h.resendVerificationEmail)

	iam.Post("/login", strictLimiter, h.login)
	iam.Post("/logout", standardLimiter, h.logout)
	iam.Post("/refresh-session", standardLimiter, h.refreshSession)

	// iam.Post("/forgot-password", strictLimiter, h.forgotPassword)
	// iam.Post("/reset-password", strictLimiter, h.resetPassword)

	// ── Authenticated profile endpoints (/me) ───────────────────────────────
	me := iam.Group("/me", standardLimiter)
	me.Use(middleware.IamJWTVerify(h.iamJWT))
	me.Get("/", h.getMe)
	// me.Patch("/", h.updateMe)
	// me.Post("/change-password", h.changePassword)

	h.logger.Info("IAM routes registered", "")
}

// ── Registration ────────────────────────────────────────────────────────────────

func (h *IAMAuthGatewayHandler) register(c *fiber.Ctx) error {

	var req iamV1.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.Register(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusCreated, "user registered successfully", nil)
}

func (h *IAMAuthGatewayHandler) verifyEmail(c *fiber.Ctx) error {

	token := c.Query("token")
	if token == "" {
		return response.Error(c, errors.BadRequest("token is required"))
	}

	if _, err := h.auth.VerifyEmail(c.UserContext(), &iamV1.VerifyEmailRequest{Token: token}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "email verified successfully", nil)
}

func (h *IAMAuthGatewayHandler) resendVerificationEmail(c *fiber.Ctx) error {

	var req iamV1.ResendVerificationEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ResendVerificationEmail(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "verification email resent", nil)
}

// ── Session ─────────────────────────────────────────────────────────────────────

func (h *IAMAuthGatewayHandler) login(c *fiber.Ctx) error {

	var req iamV1.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	sess, err := h.auth.Login(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Set(
		c,
		sess.AccessToken,
		sess.SessionToken,
		sess.AccessTokenExpiresAt.AsTime(),
		sess.SessionTokenExpiresAt.AsTime(),
	); err != nil {
		return response.Error(c, errors.Internal("failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"login successful",
		fiber.Map{"access_token": sess.AccessToken},
	)
}

func (h *IAMAuthGatewayHandler) logout(c *fiber.Ctx) error {

	st := h.cookie.GetSessionToken(c)
	if st == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	if _, err := h.auth.Logout(c.UserContext(), &iamV1.LogoutRequest{SessionToken: st}); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Clear(c); err != nil {
		log.Printf("[IAM] failed to clear cookies: %v", err)
		return response.Error(c, errors.Internal("failed to clear cookies"))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "logged out successfully", nil)
}

func (h *IAMAuthGatewayHandler) refreshSession(c *fiber.Ctx) error {

	st := h.cookie.GetSessionToken(c)
	log.Println(st)
	if st == "" {
		return response.Error(c, errors.Unauthorized("unauthorized"))
	}

	sess, err := h.auth.RefreshSession(c.UserContext(), &iamV1.RefreshSessionRequest{SessionToken: st})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	if err := h.cookie.Set(
		c,
		sess.AccessToken,
		sess.SessionToken,
		sess.AccessTokenExpiresAt.AsTime(),
		sess.SessionTokenExpiresAt.AsTime(),
	); err != nil {
		return response.Error(c, errors.Internal("failed to set cookies"))
	}

	return response.SuccessWithMessage(
		c,
		fiber.StatusOK,
		"session refreshed",
		fiber.Map{"access_token": sess.AccessToken},
	)
}

// ── Password ─────────────────────────────────────────────────────────────────────

func (h *IAMAuthGatewayHandler) forgotPassword(c *fiber.Ctx) error {

	var req iamV1.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ForgotPassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	// Always return 200 to avoid leaking whether the email exists.
	return response.SuccessWithMessage(c, fiber.StatusOK, "if this email is registered, a reset link has been sent", nil)
}

func (h *IAMAuthGatewayHandler) resetPassword(c *fiber.Ctx) error {

	var req iamV1.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ResetPassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "password reset successfully", nil)
}

// ── Profile (/me) ─────────────────────────────────────────────────────────────────

type UserDTO struct {
	ID              string    `json:"id,omitempty"`
	Email           string    `json:"email,omitempty"`
	FullName        string    `json:"full_name,omitempty"`
	Status          string    `json:"status,omitempty"`
	IsEmailVerified bool      `json:"is_email_verified,omitempty"`
	JoinedAt        string    `json:"joined_at,omitempty"`
	Avatar          AvatarDTO `json:"avatar,omitzero"`
}

type AvatarDTO struct {
	ID            int64  `json:"id,omitempty"`
	Url           string `json:"url,omitempty"`
	FallbackChar  string `json:"fallback_char,omitempty"`
	FallbackColor string `json:"fallback_color,omitempty"`
}

type GetMeResponse struct {
	User *UserDTO `json:"user,omitempty"`
}

// getMe returns the authenticated user's profile by forwarding to the IAM UserService.
func (h *IAMAuthGatewayHandler) getMe(c *fiber.Ctx) error {

	resp, err := h.user.GetMe(c.UserContext(), &iamV1.GetMeRequest{})
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	user := resp.GetUser()
	avatar := user.GetAvatar()

	avatarDTO := &AvatarDTO{
		ID:            avatar.GetId(),
		Url:           avatar.GetUrl(),
		FallbackChar:  avatar.GetFallbackChar(),
		FallbackColor: avatar.GetFallbackColor(),
	}

	userDTO := &UserDTO{
		ID:              user.GetId(),
		Email:           user.GetEmail(),
		FullName:        user.GetFullName(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.GetIsEmailVerified(),
		JoinedAt:        user.GetJoinedAt(),
		Avatar:          *avatarDTO,
	}

	respDTO := &GetMeResponse{
		User: userDTO,
	}

	log.Print("Get user :: ", respDTO)

	return response.Success(c, respDTO)
}

// updateMe updates the authenticated user's profile.
func (h *IAMAuthGatewayHandler) updateMe(c *fiber.Ctx) error {

	var req iamV1.UpdateMeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	resp, err := h.user.UpdateMe(c.UserContext(), &req)
	if err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "profile updated", resp.User)
}

// changePassword changes the password for the authenticated user.
func (h *IAMAuthGatewayHandler) changePassword(c *fiber.Ctx) error {

	var req iamV1.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, errors.BadRequest("invalid request body"))
	}

	if _, err := h.auth.ChangePassword(c.UserContext(), &req); err != nil {
		return response.Error(c, errors.FromGRPC(err))
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "password changed successfully", nil)
}
