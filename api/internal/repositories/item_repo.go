package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/utils"
	"gorm.io/gorm"
)

// ItemRepository defines the interface for item data access.
type ItemRepository interface {
	FindAll(ctx context.Context, params ItemQueryParams) ([]models.Item, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	Create(ctx context.Context, item *models.Item) error
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uuid.UUID) error
	HasActiveBatches(ctx context.Context, itemID uuid.UUID) (bool, error)
}

// ItemQueryParams holds query parameters for listing items.
type ItemQueryParams struct {
	utils.PaginationParams
	CategoryID *uuid.UUID
	IsActive   *bool
	Search     string
}

type itemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new ItemRepository.
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) FindAll(ctx context.Context, params ItemQueryParams) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Item{})

	if params.CategoryID != nil {
		query = query.Where("category_id = ?", *params.CategoryID)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.Search != "" {
		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count items: %w", err)
	}

	result := query.
		Preload("Category").
		Preload("Unit").
		Order("name ASC").
		Limit(params.Limit + 1). // fetch one extra to check has_more
		Find(&items)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find items: %w", result.Error)
	}

	return items, total, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	var item models.Item
	result := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Unit").
		First(&item, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("item not found: %w", result.Error)
	}
	return &item, nil
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	result := r.db.WithContext(ctx).Create(item)
	if result.Error != nil {
		return fmt.Errorf("failed to create item: %w", result.Error)
	}
	// Reload with preloads
	return r.db.WithContext(ctx).Preload("Category").Preload("Unit").First(item, "id = ?", item.ID).Error
}

func (r *itemRepository) Update(ctx context.Context, item *models.Item) error {
	result := r.db.WithContext(ctx).Save(item)
	if result.Error != nil {
		return fmt.Errorf("failed to update item: %w", result.Error)
	}
	// Reload with preloads
	return r.db.WithContext(ctx).Preload("Category").Preload("Unit").First(item, "id = ?", item.ID).Error
}

func (r *itemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.Item{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete item: %w", result.Error)
	}
	return nil
}

func (r *itemRepository) HasActiveBatches(ctx context.Context, itemID uuid.UUID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.Batch{}).
		Where("item_id = ? AND is_depleted = false", itemID).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check active batches: %w", result.Error)
	}
	return count > 0, nil
}
