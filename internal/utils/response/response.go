package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func NewSuccessResponse(
	statusCode int,
	message string,
	data any,
) *Response {
	return &Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func Success(c *fiber.Ctx, response *Response) error {
	return c.Status(response.StatusCode).JSON(response)
}
