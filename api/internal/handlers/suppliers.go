package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// SupplierHandler handles supplier endpoints.
type SupplierHandler struct {
	service services.SupplierService
}

// NewSupplierHandler creates a new SupplierHandler.
func NewSupplierHandler(service services.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

// List handles GET /suppliers.
func (h *SupplierHandler) List(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c)
	params := repositories.SupplierQueryParams{
		PaginationParams: pagination,
		Search:           c.Query("search"),
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		params.IsActive = &active
	}

	suppliers, pag, err := h.service.List(c.Context(), params)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.PaginatedResponse(c, suppliers, pag)
}

// Create handles POST /suppliers.
func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	var input services.SupplierInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	supplier, err := h.service.Create(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, supplier)
}

// Update handles PUT /suppliers/:id.
func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid supplier ID"))
	}

	var input services.SupplierInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	supplier, err := h.service.Update(c.Context(), id, input)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, supplier)
}

// Delete handles DELETE /suppliers/:id.
func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid supplier ID"))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
