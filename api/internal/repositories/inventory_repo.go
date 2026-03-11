package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InventoryRow represents a row from the v_current_inventory view.
type InventoryRow struct {
	ItemID       uuid.UUID `gorm:"column:item_id" json:"item_id"`
	ItemName     string    `gorm:"column:item_name" json:"item_name"`
	CategoryName string    `gorm:"column:category_name" json:"category_name"`
	Unit         string    `gorm:"column:unit" json:"unit"`
	TotalStock   float64   `gorm:"column:total_stock" json:"total_stock"`
	MinimumStock float64   `gorm:"column:minimum_stock" json:"minimum_stock"`
	IsLowStock   bool      `gorm:"column:is_low_stock" json:"is_low_stock"`
	ActiveBatches int64    `gorm:"column:active_batches" json:"active_batches"`
}

// ExpiringRow represents a row from the v_expiring_soon view.
type ExpiringRow struct {
	BatchID           uuid.UUID `gorm:"column:batch_id" json:"batch_id"`
	ItemID            uuid.UUID `gorm:"column:item_id" json:"item_id"`
	ItemName          string    `gorm:"column:item_name" json:"item_name"`
	CategoryName      string    `gorm:"column:category_name" json:"category_name"`
	BatchNumber       string    `gorm:"column:batch_number" json:"batch_number"`
	RemainingQuantity float64   `gorm:"column:remaining_quantity" json:"remaining_quantity"`
	Unit              string    `gorm:"column:unit" json:"unit"`
	ExpiryDate        time.Time `gorm:"column:expiry_date" json:"expiry_date"`
	DaysUntilExpiry   int       `gorm:"column:days_until_expiry" json:"days_until_expiry"`
}

// InventoryRepository defines the interface for inventory data access.
type InventoryRepository interface {
	GetCurrentInventory(ctx context.Context, categoryID *uuid.UUID, lowStockOnly bool, search string, limit int) ([]InventoryRow, int64, error)
	GetExpiringItems(ctx context.Context, days int) ([]ExpiringRow, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new InventoryRepository.
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetCurrentInventory(ctx context.Context, categoryID *uuid.UUID, lowStockOnly bool, search string, limit int) ([]InventoryRow, int64, error) {
	var rows []InventoryRow
	var total int64

	query := r.db.WithContext(ctx).Table("api.v_current_inventory")

	if categoryID != nil {
		// Need to join to get category_id filter — the view doesn't expose it directly
		// We'll filter by category_name indirectly or use a subquery
		query = query.Where("item_id IN (SELECT id FROM data.items WHERE category_id = ?)", *categoryID)
	}
	if lowStockOnly {
		query = query.Where("is_low_stock = true")
	}
	if search != "" {
		query = query.Where("item_name ILIKE ?", "%"+search+"%")
	}

	countQuery := r.db.WithContext(ctx).Table("api.v_current_inventory")
	if categoryID != nil {
		countQuery = countQuery.Where("item_id IN (SELECT id FROM data.items WHERE category_id = ?)", *categoryID)
	}
	if lowStockOnly {
		countQuery = countQuery.Where("is_low_stock = true")
	}
	if search != "" {
		countQuery = countQuery.Where("item_name ILIKE ?", "%"+search+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory: %w", err)
	}

	if err := query.Order("item_name ASC").Limit(limit + 1).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch inventory: %w", err)
	}

	return rows, total, nil
}

func (r *inventoryRepository) GetExpiringItems(ctx context.Context, days int) ([]ExpiringRow, error) {
	var rows []ExpiringRow
	query := r.db.WithContext(ctx).
		Table("data.batches b").
		Select(`b.id AS batch_id, i.id AS item_id, i.name AS item_name, c.name AS category_name,
			b.batch_number, b.remaining_quantity, u.abbreviation AS unit,
			b.expiry_date, b.expiry_date - CURRENT_DATE AS days_until_expiry`).
		Joins("JOIN data.items i ON b.item_id = i.id").
		Joins("JOIN data.categories c ON i.category_id = c.id").
		Joins("JOIN data.units u ON i.unit_id = u.id").
		Where("b.is_depleted = false AND b.remaining_quantity > 0 AND b.expiry_date <= CURRENT_DATE + ? * INTERVAL '1 day' AND b.expiry_date >= CURRENT_DATE", days).
		Order("b.expiry_date ASC")

	if err := query.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch expiring items: %w", err)
	}
	return rows, nil
}
