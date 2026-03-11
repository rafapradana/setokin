package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// StockOutHandler handles stock out endpoints.
type StockOutHandler struct {
	service services.StockService
}

// NewStockOutHandler creates a new StockOutHandler.
func NewStockOutHandler(service services.StockService) *StockOutHandler {
	return &StockOutHandler{service: service}
}

// Create handles POST /stock-out (with FEFO logic).
func (h *StockOutHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
	}

	var input services.StockOutInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	stockOut, deductions, remainingStock, err := h.service.CreateStockOut(c.Context(), input, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"stock_out":       stockOut,
		"deductions":      deductions,
		"remaining_stock": remainingStock,
	})
}

// List handles GET /stock-out.
func (h *StockOutHandler) List(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c)
	params := repositories.StockOutQueryParams{
		PaginationParams: pagination,
	}

	if itemID := c.Query("item_id"); itemID != "" {
		id, err := uuid.Parse(itemID)
		if err == nil {
			params.ItemID = &id
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			params.StartDate = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			endOfDay := t.Add(24*time.Hour - time.Second)
			params.EndDate = &endOfDay
		}
	}

	items, pag, err := h.service.ListStockOut(c.Context(), params)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.PaginatedResponse(c, items, pag)
}

// GetDetails handles GET /stock-out/:id/details.
func (h *StockOutHandler) GetDetails(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid stock out ID"))
	}

	stockOut, details, err := h.service.GetStockOutDetails(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"stock_out":  stockOut,
		"deductions": details,
	})
}
