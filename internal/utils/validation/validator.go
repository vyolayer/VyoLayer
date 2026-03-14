package validation

import (
	"strings"
	"vyolayer/pkg/errors"
	"vyolayer/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Global instance to cache regex compilation
var validate = validator.New()

type ErrorResponse struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// ValidateStruct checks for struct tags and returns formatted errors
func ValidateStruct(payload interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = strings.ToLower(err.Field())
			element.Tag = err.Tag()
			element.Value = err.Param() // e.g. "8" for min=8

			// Custom messages
			switch err.Tag() {
			case "required":
				element.Message = element.Field + " is required"
			case "email":
				element.Message = "Invalid email format"
			case "min":
				element.Message = element.Field + " must be at least " + err.Param() + " characters"
			case "max":
				element.Message = element.Field + " must be at most " + err.Param() + " characters"
			default:
				element.Message = "Invalid " + element.Field
			}

			errors = append(errors, &element)
		}
		return errors
	}
	return nil
}

func ValidationErrorsToResponse(ctx *fiber.Ctx, errs []*ErrorResponse) error {
	return response.Error(ctx,
		errors.ValidationFailed("Validation failed").
			WithMetadata("validation_errors", errs),
	)
}
