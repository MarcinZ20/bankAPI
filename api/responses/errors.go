package responses

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// ValidationError returns a consistent validation error response
func ValidationError(message string) error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

// DatabaseError returns a consistent database error response
func DatabaseError(err error) error {
	return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
}

// NotFoundError returns a consistent not found error response
func NotFoundError(resourceType string, identifier string) error {
	return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("%s not found: %s", resourceType, identifier))
}

// AlreadyExistsError returns a consistent conflict error response
func AlreadyExistsError(message string) error {
	return fiber.NewError(fiber.StatusConflict, message)
}

// FormattingResponseError returns a consistent response formatting error
func FormattingResponseError(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}

// InternalServerError returns a consistent internal server error response
func InternalServerError(message string) error {
	return fiber.NewError(fiber.StatusInternalServerError, message)
}
