package responses

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Returns a consistent validation error response
func ValidationError(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

// Returns a consistent database error response
func DatabaseError(err error) error {
	return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
}

// Returns a consistent not found error response
func NotFoundError(resourceType string, identifier string) error {
	return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("%s not found: %s", resourceType, identifier))
}

// Returns a consistent conflict error response
func AlreadyExistsError(message string) error {
	return fiber.NewError(fiber.StatusConflict, message)
}

// Returns a consistent response formatting error
func FormattingResponseError(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

// Returns a consistent internal server error response
func InternalServerError(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
