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

// StockInHandler handles stock in endpoints.
type StockInHandler struct {
	service services.StockService
}

// NewStockInHandler creates a new StockInHandler.
func NewStockInHandler(service services.StockService) *StockInHandler {
	return &StockInHandler{service: service}
}

// Create handles POST /stock-in.
func (h *StockInHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
	}

	var input services.StockInInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	stockIn, batch, err := h.service.CreateStockIn(c.Context(), input, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"stock_in": stockIn,
		"batch":    batch,
	})
}

// List handles GET /stock-in.
func (h *StockInHandler) List(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c)
	params := repositories.StockInQueryParams{
		PaginationParams: pagination,
	}

	if itemID := c.Query("item_id"); itemID != "" {
		id, err := uuid.Parse(itemID)
		if err == nil {
			params.ItemID = &id
		}
	}
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		id, err := uuid.Parse(supplierID)
		if err == nil {
			params.SupplierID = &id
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

	items, pag, err := h.service.ListStockIn(c.Context(), params)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.PaginatedResponse(c, items, pag)
}
