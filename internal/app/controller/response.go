package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/response"
)

// Success sends a successful response with data
func Success(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return response.SuccessWithMessage(ctx, statusCode, message, data)
}

// SuccessMessage sends a successful response with only a message
func SuccessMessage(ctx *fiber.Ctx, statusCode int, message string) error {
	return response.SuccessWithMessage(ctx, statusCode, message, nil)
}

// Error sends an error response
func Error(ctx *fiber.Ctx, err error) error {
	return response.Error(ctx, err)
}
