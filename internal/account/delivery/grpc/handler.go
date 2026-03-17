package grpc

import (
	"context"
	"log"

	"github.com/vyolayer/vyolayer/internal/account/usecase"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

type AccountHandler struct {
	accountV1.UnimplementedAccountServiceServer
	usecase usecase.AccountUsecase
}

func NewAccountHandler(usecase usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		usecase: usecase,
	}
}

func (h *AccountHandler) Register(
	ctx context.Context,
	req *accountV1.RegisterRequest,
) (*accountV1.RegisterResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	userId, appErr := h.usecase.Register(
		ctx,
		apiKeyInfo.ProjectID,
		req.Email,
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.RegisterResponse{UserId: userId}, nil
}

func (h *AccountHandler) VerifyEmail(
	ctx context.Context,
	req *accountV1.VerifyEmailRequest,
) (*accountV1.VerifyEmailResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	appErr := h.usecase.VerifyEmail(
		ctx,
		apiKeyInfo.ProjectID,
		req.Token,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.VerifyEmailResponse{}, nil
}

func (h *AccountHandler) ResendVerificationEmail(
	ctx context.Context,
	req *accountV1.ResendVerificationEmailRequest,
) (*accountV1.ResendVerificationEmailResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	appErr := h.usecase.ResendVerificationEmail(
		ctx,
		apiKeyInfo.ProjectID,
		req.Email,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.ResendVerificationEmailResponse{}, nil
}

func (h *AccountHandler) Login(
	ctx context.Context,
	req *accountV1.LoginRequest,
) (*accountV1.LoginResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	resp, appErr := h.usecase.Login(
		ctx,
		apiKeyInfo.ProjectID,
		req.Email,
		req.Password,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (h *AccountHandler) Logout(
	ctx context.Context,
	req *accountV1.LogoutRequest,
) (*accountV1.LogoutResponse, error) {
	aki, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	_, userID, err := ctxutil.ExtractVyoServiceAccountDetails(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Logout request for user: ", userID, aki.ProjectID, req.RefreshToken)
	appErr := h.usecase.Logout(
		ctx,
		aki.ProjectID,
		userID,
		req.RefreshToken,
	)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.LogoutResponse{}, nil
}
