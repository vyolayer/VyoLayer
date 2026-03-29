package grpc

import (
	"context"

	"github.com/vyolayer/vyolayer/internal/iam/usecase"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

func NewIAMUserHandler(uu *usecase.UserUsecase) *IAMUserHandler {
	return &IAMUserHandler{uu: uu}
}

// GetMe returns the authenticated caller's profile.
func (h *IAMUserHandler) GetMe(ctx context.Context, _ *iAMV1.GetMeRequest) (*iAMV1.GetMeResponse, error) {
	user, err := h.uu.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	return &iAMV1.GetMeResponse{User: user}, nil
}

// UpdateMe updates the authenticated caller's profile.
func (h *IAMUserHandler) UpdateMe(ctx context.Context, req *iAMV1.UpdateMeRequest) (*iAMV1.UpdateMeResponse, error) {
	user, err := h.uu.UpdateMe(ctx, req)
	if err != nil {
		return nil, err
	}

	return &iAMV1.UpdateMeResponse{User: user}, nil
}
