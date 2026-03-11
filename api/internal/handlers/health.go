package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthCheck handles GET /health.
func HealthCheck(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dbStatus := "healthy"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "unhealthy"
		}

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now(),
			"services": fiber.Map{
				"database": dbStatus,
			},
		})
	}
}
