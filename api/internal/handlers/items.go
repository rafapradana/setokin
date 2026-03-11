package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// ItemHandler handles item endpoints.
type ItemHandler struct {
	service services.ItemService
}

// NewItemHandler creates a new ItemHandler.
func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// List handles GET /items.
func (h *ItemHandler) List(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c)

	params := repositories.ItemQueryParams{
		PaginationParams: pagination,
		Search:           c.Query("search"),
	}

	if catID := c.Query("category_id"); catID != "" {
		id, err := uuid.Parse(catID)
		if err == nil {
			params.CategoryID = &id
		}
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		params.IsActive = &active
	}

	items, pag, err := h.service.List(c.Context(), params)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.PaginatedResponse(c, items, pag)
}

// Get handles GET /items/:id.
func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid item ID"))
	}

	item, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

// Create handles POST /items.
func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var input services.ItemInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	item, err := h.service.Create(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, item)
}

// Update handles PUT /items/:id.
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid item ID"))
	}

	var input services.ItemInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	item, err := h.service.Update(c.Context(), id, input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, item)
}

// Delete handles DELETE /items/:id.
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid item ID"))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
