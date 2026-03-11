package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// InventoryHandler handles inventory endpoints.
type InventoryHandler struct {
	repo repositories.InventoryRepository
}

// NewInventoryHandler creates a new InventoryHandler.
func NewInventoryHandler(repo repositories.InventoryRepository) *InventoryHandler {
	return &InventoryHandler{repo: repo}
}

// GetCurrentInventory handles GET /inventory.
func (h *InventoryHandler) GetCurrentInventory(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c)

	var categoryID *uuid.UUID
	if catID := c.Query("category_id"); catID != "" {
		id, err := uuid.Parse(catID)
		if err == nil {
			categoryID = &id
		}
	}
	lowStockOnly := c.Query("low_stock_only") == "true"
	search := c.Query("search")

	rows, total, err := h.repo.GetCurrentInventory(c.Context(), categoryID, lowStockOnly, search, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to fetch inventory"))
	}

	hasMore := len(rows) > pagination.Limit
	if hasMore {
		rows = rows[:pagination.Limit]
	}
	var nextCursor *string
	if hasMore && len(rows) > 0 {
		cursor := utils.EncodeCursor(map[string]string{"id": rows[len(rows)-1].ItemID.String()})
		nextCursor = &cursor
	}

	// Add stock_status computed field
	type inventoryItem struct {
		repositories.InventoryRow
		StockStatus string `json:"stock_status"`
	}
	var result []inventoryItem
	var lowStockCount, outOfStockCount int64
	for _, row := range rows {
		status := "adequate"
		if row.TotalStock == 0 {
			status = "out_of_stock"
			outOfStockCount++
		} else if row.IsLowStock {
			status = "low"
			lowStockCount++
		}
		result = append(result, inventoryItem{InventoryRow: row, StockStatus: status})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       result,
		"pagination": utils.Pagination{NextCursor: nextCursor, HasMore: hasMore, Total: total},
		"summary": fiber.Map{
			"total_items":       total,
			"low_stock_items":   lowStockCount,
			"out_of_stock_items": outOfStockCount,
		},
	})
}

// GetExpiring handles GET /inventory/expiring.
func (h *InventoryHandler) GetExpiring(c *fiber.Ctx) error {
	days := 3
	if d := c.Query("days"); d != "" {
		parsed, err := strconv.Atoi(d)
		if err == nil && parsed > 0 && parsed <= 30 {
			days = parsed
		}
	}

	rows, err := h.repo.GetExpiringItems(c.Context(), days)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to fetch expiring items"))
	}

	// Add urgency field
	type expiringItem struct {
		repositories.ExpiringRow
		Urgency string `json:"urgency"`
	}
	var result []expiringItem
	var criticalCount, highCount int64
	for _, row := range rows {
		urgency := "medium"
		if row.DaysUntilExpiry <= 1 {
			urgency = "critical"
			criticalCount++
		} else if row.DaysUntilExpiry <= 3 {
			urgency = "high"
			highCount++
		}
		result = append(result, expiringItem{ExpiringRow: row, Urgency: urgency})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
		"summary": fiber.Map{
			"total_expiring_batches": len(result),
			"critical_count":         criticalCount,
			"high_count":             highCount,
		},
	})
}
