package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// BatchHandler handles batch endpoints.
type BatchHandler struct {
	service services.StockService
}

// NewBatchHandler creates a new BatchHandler.
func NewBatchHandler(service services.StockService) *BatchHandler {
	return &BatchHandler{service: service}
}

// ListByItem handles GET /items/:item_id/batches.
func (h *BatchHandler) ListByItem(c *fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("item_id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid item ID"))
	}

	var isDepleted *bool
	if d := c.Query("is_depleted"); d != "" {
		val := d == "true"
		isDepleted = &val
	}

	batches, err := h.service.ListBatchesByItem(c.Context(), itemID, isDepleted)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	// Build enriched response with computed fields
	type batchItem struct {
		ID                uuid.UUID `json:"id"`
		BatchNumber       string    `json:"batch_number"`
		InitialQuantity   float64   `json:"initial_quantity"`
		RemainingQuantity float64   `json:"remaining_quantity"`
		ExpiryDate        string    `json:"expiry_date"`
		IsDepleted        bool      `json:"is_depleted"`
		Status            string    `json:"status"`
		DaysUntilExpiry   int       `json:"days_until_expiry"`
	}

	var result []batchItem
	for _, b := range batches {
		result = append(result, batchItem{
			ID:                b.ID,
			BatchNumber:       b.BatchNumber,
			InitialQuantity:   b.InitialQuantity,
			RemainingQuantity: b.RemainingQuantity,
			ExpiryDate:        b.ExpiryDate.Format("2006-01-02"),
			IsDepleted:        b.IsDepleted,
			Status:            b.Status(),
			DaysUntilExpiry:   b.DaysUntilExpiry(),
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, result)
}

// Get handles GET /batches/:id.
func (h *BatchHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid batch ID"))
	}

	batch, usageHistory, err := h.service.GetBatch(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"id":                 batch.ID,
		"batch_number":       batch.BatchNumber,
		"item":               batch.Item,
		"initial_quantity":   batch.InitialQuantity,
		"remaining_quantity": batch.RemainingQuantity,
		"expiry_date":        batch.ExpiryDate.Format("2006-01-02"),
		"is_depleted":        batch.IsDepleted,
		"status":             batch.Status(),
		"days_until_expiry":  batch.DaysUntilExpiry(),
		"usage_history":      usageHistory,
		"created_at":         batch.CreatedAt,
		"updated_at":         batch.UpdatedAt,
	})
}
