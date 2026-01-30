package controller

import (
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse = response.ErrorResponse

func Error(ctx *fiber.Ctx, err *ErrorResponse) error {
	return ctx.Status(err.StatusCode).JSON(err)
}

func Success(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	res := response.NewSuccessResponse(statusCode, message, data)
	return ctx.Status(statusCode).JSON(res)
}

func SuccessMessage(ctx *fiber.Ctx, statusCode int, message string) error {
	res := response.NewSuccessResponse(statusCode, message, nil)
	return ctx.Status(statusCode).JSON(res)
}

func SuccessData(ctx *fiber.Ctx, statusCode int, data interface{}) error {
	res := response.NewDataResponse(statusCode, data)
	return ctx.Status(statusCode).JSON(res)
}
