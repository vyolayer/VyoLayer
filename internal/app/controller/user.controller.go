package controller

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/service"
	"worklayer/internal/utils/response"

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

func (uc *userController) GetMe(ctx *fiber.Ctx) error {
	localUserIDVal := ctx.Locals("user_id")
	localUserID, ok := localUserIDVal.(types.UserID)
	if !ok || localUserID.IsNil() {
		return Error(ctx, response.NewErrorMessage(fiber.StatusUnauthorized, "Invalid or missing user context"))
	}
	user, err := uc.userService.GetUser(localUserID)

	if err != nil {
		return Error(ctx, response.NewErrorMessage(err.Code, err.Message))
	}

	return SuccessData(
		ctx,
		fiber.StatusOK,
		dto.MeResponseDTO{UserDTO: dto.FromDomainUser(user)},
	)
}
