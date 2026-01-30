package controller

import (
	"log"
	"worklayer/internal/service"
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	userId := ctx.Locals("user_id").(uuid.UUID)
	user, err := uc.userService.GetUser(userId)

	if err != nil {
		log.Printf("USER CONTROLLER :: GetMe : %v", err.Error())
		return Error(ctx, response.InternalServerError("Failed to get user"))
	}

	return SuccessData(
		ctx,
		fiber.StatusOK,
		user,
	)
}
