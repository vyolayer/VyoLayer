package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	repo "github.com/vyolayer/vyolayer/internal/iam/repository"
	"github.com/vyolayer/vyolayer/pkg/ctxutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserUsecase handles user profile operations.
type UserUsecase struct {
	log *logger.AppLogger
	ur  repo.IAMUserRepository
}

func NewUserUsecase(log *logger.AppLogger, userRepo repo.IAMUserRepository) *UserUsecase {
	return &UserUsecase{ur: userRepo, log: log}
}

// GetMe returns the authenticated user's profile.
func (uc *UserUsecase) GetMe(ctx context.Context) (*iAMV1.User, error) {
	uid, err := extractCallerID(ctx)
	if err != nil {
		uc.log.Error("(UserUsecase.GetMe): ", err)
		return nil, err
	}

	user, err := uc.ur.FindByID(ctx, uid)
	if err != nil {
		uc.log.Error("(UserUsecase.GetMe): ", err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	avatar := &iAMV1.Avatar{
		Id:            user.Avatar.ID,
		Url:           user.Avatar.URL,
		FallbackChar:  user.Avatar.FallbackChar,
		FallbackColor: user.Avatar.FallbackColor,
	}

	u := &iAMV1.User{
		Id:              user.ID.String(),
		Email:           user.GetEmail(),
		FullName:        user.GetFullName(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.IsEmailVerified,
		JoinedAt:        user.Timestamps.CreatedAt.UTC().Format(time.RFC3339),
	}

	if user.Avatar != nil {
		u.Avatar = avatar
	}

	uc.log.Info("(UserUsecase.GetMe): ", u)
	return u, nil
}

// UpdateMe mutates the authenticated user's profile.
func (uc *UserUsecase) UpdateMe(ctx context.Context, req *iAMV1.UpdateMeRequest) (*iAMV1.User, error) {
	uid, err := extractCallerID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.ur.FindByID(ctx, uid)
	if err != nil {
		uc.log.Error("UserUsecase.UpdateMe: ", err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if n := req.GetFullName(); n != "" {
		user.SetFullName(n)
	}

	if err := uc.ur.Update(ctx, user); err != nil {
		return nil, err
	}

	return &iAMV1.User{
		Id:              user.ID.String(),
		Email:           user.GetEmail(),
		FullName:        user.GetFullName(),
		Status:          user.GetStatus(),
		IsEmailVerified: user.IsEmailVerified,
		JoinedAt:        user.Timestamps.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// extractCallerID parses the user UUID from the gRPC context.
func extractCallerID(ctx context.Context) (uuid.UUID, error) {
	raw, err := ctxutil.ExtractIAMUserID(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	uid, parseErr := uuid.Parse(raw)
	if parseErr != nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user id in context")
	}
	if uid == uuid.Nil {
		return uuid.Nil, status.Error(codes.Unauthenticated, "invalid user id in context")
	}

	return uid, nil
}
