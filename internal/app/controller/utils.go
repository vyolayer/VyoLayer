package controller

import (
	"errors"
	"worklayer/internal/platform/database/types"

	"github.com/gofiber/fiber/v2"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
	UserIDContextKey       = "user_id"
	EmailContextKey        = "user_email"
)

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx *fiber.Ctx) (types.UserID, error) {
	localUserIDVal := ctx.Locals(UserIDContextKey) // expect to be uuid.UUID
	if localUserIDVal == nil {
		return types.UserID{}, errors.New("user ID not found in context")
	}

	localUserID, err := types.ReconstructUserID(localUserIDVal.(string))
	if err != nil {
		return types.UserID{}, err
	}

	if localUserID.IsNil() || localUserID == nil || localUserID.InternalID().IsNil() {
		return types.UserID{}, errors.New("invalid user ID in context")
	}

	return *localUserID, nil
}
