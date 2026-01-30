package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

func NewSuccessResponse(statusCode int, message string, data any) *Response {
	return &Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewMessageResponse(message string) *Response {
	return &Response{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    message,
		Data:       nil,
	}
}

func NewDataResponse(statusCode int, data any) *Response {
	return &Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    "Success",
		Data:       data,
	}
}
