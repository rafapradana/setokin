package utils

import (
	"github.com/gofiber/fiber/v2"
	apperrors "github.com/setokin/api/pkg/errors"
)

// SuccessResponse sends a JSON success response.
func SuccessResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data": data,
	})
}

// ErrorResponse sends a JSON error response from an AppError.
func ErrorResponse(c *fiber.Ctx, err *apperrors.AppError) error {
	body := fiber.Map{
		"error": fiber.Map{
			"code":    err.Code,
			"message": err.Message,
		},
	}
	if len(err.Details) > 0 {
		body["error"].(fiber.Map)["details"] = err.Details
	}
	return c.Status(err.Status).JSON(body)
}

// PaginatedResponse sends a paginated JSON response.
func PaginatedResponse(c *fiber.Ctx, data interface{}, pagination *Pagination) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       data,
		"pagination": pagination,
	})
}

// Pagination holds pagination metadata.
type Pagination struct {
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
	Total      int64   `json:"total"`
}
