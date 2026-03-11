package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// CategoryHandler handles category endpoints.
type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// List handles GET /categories.
func (h *CategoryHandler) List(c *fiber.Ctx) error {
	categories, err := h.service.List(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, categories)
}

// Get handles GET /categories/:id.
func (h *CategoryHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid category ID"))
	}

	category, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

// Create handles POST /categories.
func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var input services.CategoryInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	category, err := h.service.Create(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, category)
}

// Update handles PUT /categories/:id.
func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid category ID"))
	}

	var input services.CategoryInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	category, err := h.service.Update(c.Context(), id, input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, category)
}

// Delete handles DELETE /categories/:id.
func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid category ID"))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
