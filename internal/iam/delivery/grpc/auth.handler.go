package grpc

import (
	"context"

	"github.com/vyolayer/vyolayer/internal/iam/usecase"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

func NewIAMAuthHandler(authUsecase *usecase.AuthUsecase) *IAMAuthHandler {
	return &IAMAuthHandler{au: authUsecase}
}

// ── Registration flow ──────────────────────────────────────────────────────────

func (h *IAMAuthHandler) Register(ctx context.Context, req *iAMV1.RegisterRequest) (*iAMV1.RegisterResponse, error) {
	userID, err := h.au.RegisterUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &iAMV1.RegisterResponse{UserId: userID}, nil
}

func (h *IAMAuthHandler) VerifyEmail(ctx context.Context, req *iAMV1.VerifyEmailRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.VerifyEmail(ctx, req.GetToken()); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "email verified successfully"}, nil
}

func (h *IAMAuthHandler) ResendVerificationEmail(ctx context.Context, req *iAMV1.ResendVerificationEmailRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.ResendVerificationEmail(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "verification email sent"}, nil
}

// ── Session flow ───────────────────────────────────────────────────────────────

func (h *IAMAuthHandler) Login(ctx context.Context, req *iAMV1.LoginRequest) (*iAMV1.SessionTokenResponse, error) {
	resp, err := h.au.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *IAMAuthHandler) RefreshSession(ctx context.Context, req *iAMV1.RefreshSessionRequest) (*iAMV1.SessionTokenResponse, error) {
	return h.au.RefreshToken(ctx, req)
}

func (h *IAMAuthHandler) Logout(ctx context.Context, req *iAMV1.LogoutRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.Logout(ctx, req); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "logged out successfully"}, nil
}

// ── Password flow ──────────────────────────────────────────────────────────────

func (h *IAMAuthHandler) ChangePassword(ctx context.Context, req *iAMV1.ChangePasswordRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.ChangePassword(ctx, req); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "password changed successfully"}, nil
}

func (h *IAMAuthHandler) ForgotPassword(ctx context.Context, req *iAMV1.ForgotPasswordRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.ForgotPassword(ctx, req.GetEmail()); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "password reset link sent"}, nil
}

func (h *IAMAuthHandler) ResetPassword(ctx context.Context, req *iAMV1.ResetPasswordRequest) (*iAMV1.IAMSuccessResponse, error) {
	if err := h.au.ResetPassword(ctx, req); err != nil {
		return nil, err
	}

	return &iAMV1.IAMSuccessResponse{Message: "password reset successfully"}, nil
}
