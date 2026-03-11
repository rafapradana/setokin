// Package handlers provides HTTP request handlers for the Setokin API.
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input services.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	user, tokens, err := h.authService.Register(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"user":   user,
		"tokens": tokens,
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input services.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	user, tokens, err := h.authService.Login(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"user":   user,
		"tokens": tokens,
	})
}

// Refresh handles POST /auth/refresh.
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("refresh_token is required"))
	}

	tokens, err := h.authService.RefreshTokens(c.Context(), body.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, tokens)
}

// Logout handles POST /auth/logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("refresh_token is required"))
	}

	if err := h.authService.Logout(c.Context(), body.RefreshToken); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetMe handles GET /auth/me.
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
	}

	user, err := h.authService.GetCurrentUser(c.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, user)
}
