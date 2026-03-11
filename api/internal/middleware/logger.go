package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/setokin/api/pkg/logger"
	"go.uber.org/zap"
)

// LoggerMiddleware returns a middleware that logs HTTP requests.
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		logger.Info("HTTP request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
