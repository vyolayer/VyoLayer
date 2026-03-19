package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/usecase"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/errors"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountHandler struct {
	accountV1.UnimplementedAccountServiceServer
	usecase   usecase.AccountUsecase
	sessionuc usecase.SessionUsecase
	// account recover
	recoverUC usecase.AccountRecoverUsecase
}

func NewAccountHandler(
	usecase usecase.AccountUsecase,
	sessionUsecase usecase.SessionUsecase,
	recoverUC usecase.AccountRecoverUsecase,
) *AccountHandler {
	return &AccountHandler{
		usecase:   usecase,
		sessionuc: sessionUsecase,
		recoverUC: recoverUC,
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

func (h *AccountHandler) RefreshSession(
	ctx context.Context,
	req *accountV1.RefreshSessionRequest,
) (*accountV1.RefreshSessionResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	resp, appErr := h.sessionuc.RefreshToken(ctx, apiKeyInfo.ProjectID, req.RefreshToken)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.RefreshSessionResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (h *AccountHandler) AllSessions(
	ctx context.Context,
	req *accountV1.AllSessionsRequest,
) (*accountV1.AllSessionsResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}
	projectID, userID, err := ctxutil.ExtractVyoServiceAccountDetails(ctx)
	if err != nil {
		return nil, err
	}
	if apiKeyInfo.ProjectID != projectID {
		return nil, errors.NewWithMessage(errors.ErrProjectInfoNotLoaded, "invalid project id")
	}

	resp, appErr := h.sessionuc.ListSessions(ctx, projectID, userID)
	if appErr != nil {
		return nil, appErr
	}

	return resp, nil
}

func (h *AccountHandler) RevokeSession(
	ctx context.Context,
	req *accountV1.RevokeSessionRequest,
) (*accountV1.RevokeSessionResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}
	projectID, userID, err := ctxutil.ExtractVyoServiceAccountDetails(ctx)
	if err != nil {
		return nil, err
	}
	if apiKeyInfo.ProjectID != projectID {
		return nil, errors.NewWithMessage(errors.ErrProjectInfoNotLoaded, "invalid project id")
	}

	appErr := h.sessionuc.RevokeSession(ctx, projectID, userID, uuid.MustParse(req.GetSessionId()))
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.RevokeSessionResponse{}, nil
}

func (h *AccountHandler) RevokeAllSessions(
	ctx context.Context,
	req *accountV1.RevokeAllSessionsRequest,
) (*accountV1.RevokeAllSessionsResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}
	projectID, userID, err := ctxutil.ExtractVyoServiceAccountDetails(ctx)
	if err != nil {
		return nil, err
	}
	if apiKeyInfo.ProjectID != projectID {
		return nil, errors.NewWithMessage(errors.ErrProjectInfoNotLoaded, "invalid project id")
	}

	appErr := h.sessionuc.RevokeAllSessions(ctx, projectID, userID)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.RevokeAllSessionsResponse{}, nil
}

func (h *AccountHandler) ChangePassword(
	ctx context.Context,
	req *accountV1.ChangePasswordRequest,
) (*accountV1.ChangePasswordResponse, error) {
	log.Println("Change password")
	projectID, userID, err := h.validateRequest(ctx)
	if err != nil {
		return nil, err
	}

	appErr := h.recoverUC.ChangePassword(ctx, projectID, userID, req.OldPassword, req.NewPassword)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.ChangePasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

func (h *AccountHandler) ForgotPassword(
	ctx context.Context,
	req *accountV1.ForgotPasswordRequest,
) (*accountV1.ForgotPasswordResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	appErr := h.recoverUC.ForgotPassword(ctx, apiKeyInfo.ProjectID, req.Email)
	if appErr != nil {
		return nil, appErr
	}

	return &accountV1.ForgotPasswordResponse{
		Message: "Email sent successfully",
	}, nil
}

func (h *AccountHandler) ResetPassword(
	ctx context.Context,
	req *accountV1.ResetPasswordRequest,
) (*accountV1.ResetPasswordResponse, error) {
	apiKeyInfo, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return nil, err
	}

	ucErr := h.recoverUC.ResetPassword(ctx, apiKeyInfo.ProjectID, req.Token, req.NewPassword)
	if ucErr != nil {
		return nil, ucErr
	}

	return &accountV1.ResetPasswordResponse{
		Message: "Password reset successfully",
	}, nil
}

// validateRequest - gets the projectID and userID from the context
func (h *AccountHandler) validateRequest(ctx context.Context) (uuid.UUID, uuid.UUID, error) {
	aki, err := ctxutil.ExtractAPIKeyInfo(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	projectID, userID, err := ctxutil.ExtractVyoServiceAccountDetails(ctx)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	if aki.ProjectID != projectID {
		return uuid.Nil, uuid.Nil, status.Error(codes.PermissionDenied, "invalid project id")
	}
	return projectID, userID, nil
}
