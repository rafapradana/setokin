// Package middleware provides HTTP middleware for the Setokin API.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// AuthMiddleware returns a middleware that validates JWT access tokens.
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, apperrors.ErrTokenInvalid)
		}

		claims, err := utils.ValidateAccessToken(parts[1], jwtSecret)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				return utils.ErrorResponse(c, apperrors.ErrTokenExpired)
			}
			return utils.ErrorResponse(c, apperrors.ErrTokenInvalid)
		}

		// Store user info in context locals
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RoleMiddleware returns a middleware that checks if the user has one of the required roles.
func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, apperrors.ErrInsufficientPermissions)
	}
}
