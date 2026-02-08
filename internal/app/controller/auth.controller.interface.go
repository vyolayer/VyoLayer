package controller

import (
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	RegisterUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	RefreshSession(ctx *fiber.Ctx) error
	LogoutUser(ctx *fiber.Ctx) error
	ValidateSession(ctx *fiber.Ctx) error
}
