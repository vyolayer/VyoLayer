package controller

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/pkg/errors"
	"worklayer/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	GetMe(ctx *fiber.Ctx) error
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

// GetMe godoc
// @Summary Get current user
// @Description Get the profile information of the currently authenticated user.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=dto.MeResponseDTO}
// @Failure 401 {object} response.ErrorResponse
// @Router /users/me [get]
func (uc *userController) GetMe(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return response.Error(ctx, errors.Unauthorized("Invalid or missing user context"))
	}

	user, err := uc.userService.GetUser(localUserID)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(
		ctx,
		dto.MeResponseDTO{UserDTO: dto.FromDomainUser(user)},
	)
}
