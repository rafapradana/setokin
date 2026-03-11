package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware returns a rate limiting middleware.
func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":        "rate_limit_exceeded",
					"message":     "Rate limit exceeded. Try again later.",
					"retry_after": 60,
				},
			})
		},
	})
}
