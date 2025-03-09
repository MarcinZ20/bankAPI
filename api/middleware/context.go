package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// WithTimeout adds a timeout context to the request
func WithTimeout(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.Locals("ctx", ctx)
		return c.Next()
	}
}

// GetRequestContext retrieves the timeout context from fiber context
func GetRequestContext(c *fiber.Ctx) (context.Context, bool) {
	ctx, ok := c.Locals("ctx").(context.Context)
	return ctx, ok
}
