package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err ...string) error {
	response := Response{
		Success: false,
		Message: message,
	}

	if len(err) > 0 {
		response.Error = err[0]
	}

	return c.Status(statusCode).JSON(response)
}

func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err.Error())
}
