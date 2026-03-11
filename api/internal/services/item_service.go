package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// ItemInput holds data for creating/updating an item.
type ItemInput struct {
	Name         string    `json:"name"`
	CategoryID   uuid.UUID `json:"category_id"`
	UnitID       uuid.UUID `json:"unit_id"`
	MinimumStock float64   `json:"minimum_stock"`
	Description  *string   `json:"description"`
	IsActive     *bool     `json:"is_active"`
}

// ItemService defines the interface for item business logic.
type ItemService interface {
	List(ctx context.Context, params repositories.ItemQueryParams) ([]models.Item, *utils.Pagination, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	Create(ctx context.Context, input ItemInput) (*models.Item, error)
	Update(ctx context.Context, id uuid.UUID, input ItemInput) (*models.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type itemService struct {
	itemRepo     repositories.ItemRepository
	categoryRepo repositories.CategoryRepository
	unitRepo     repositories.UnitRepository
}

// NewItemService creates a new ItemService.
func NewItemService(
	itemRepo repositories.ItemRepository,
	categoryRepo repositories.CategoryRepository,
	unitRepo repositories.UnitRepository,
) ItemService {
	return &itemService{
		itemRepo:     itemRepo,
		categoryRepo: categoryRepo,
		unitRepo:     unitRepo,
	}
}

func (s *itemService) List(ctx context.Context, params repositories.ItemQueryParams) ([]models.Item, *utils.Pagination, error) {
	items, total, err := s.itemRepo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch items")
	}

	hasMore := len(items) > params.Limit
	if hasMore {
		items = items[:params.Limit]
	}

	var nextCursor *string
	if hasMore && len(items) > 0 {
		cursor := utils.EncodeCursor(map[string]string{"id": items[len(items)-1].ID.String()})
		nextCursor = &cursor
	}

	pagination := &utils.Pagination{
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Total:      total,
	}

	return items, pagination, nil
}

func (s *itemService) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Item not found")
	}
	return item, nil
}

func (s *itemService) Create(ctx context.Context, input ItemInput) (*models.Item, error) {
	validationErrors := utils.CollectErrors(
		utils.ValidateRequired("name", input.Name),
	)
	if input.CategoryID == uuid.Nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "category_id", Message: "Category is required"})
	}
	if input.UnitID == uuid.Nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "unit_id", Message: "Unit is required"})
	}
	if input.MinimumStock < 0 {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "minimum_stock", Message: "Minimum stock must be >= 0"})
	}
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Verify category exists
	if _, err := s.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
		return nil, apperrors.ErrValidation.WithMessage("Category not found")
	}

	item := &models.Item{
		Name:         input.Name,
		CategoryID:   input.CategoryID,
		UnitID:       input.UnitID,
		MinimumStock: input.MinimumStock,
		Description:  input.Description,
		IsActive:     true,
	}

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to create item")
	}
	return item, nil
}

func (s *itemService) Update(ctx context.Context, id uuid.UUID, input ItemInput) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Item not found")
	}

	validationErrors := utils.CollectErrors(
		utils.ValidateRequired("name", input.Name),
	)
	if input.MinimumStock < 0 {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "minimum_stock", Message: "Minimum stock must be >= 0"})
	}
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	if input.CategoryID != uuid.Nil {
		if _, err := s.categoryRepo.FindByID(ctx, input.CategoryID); err != nil {
			return nil, apperrors.ErrValidation.WithMessage("Category not found")
		}
		item.CategoryID = input.CategoryID
	}
	if input.UnitID != uuid.Nil {
		item.UnitID = input.UnitID
	}

	item.Name = input.Name
	item.MinimumStock = input.MinimumStock
	item.Description = input.Description
	if input.IsActive != nil {
		item.IsActive = *input.IsActive
	}

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to update item")
	}
	return item, nil
}

func (s *itemService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.itemRepo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound.WithMessage("Item not found")
	}

	// Check for active batches
	hasActive, err := s.itemRepo.HasActiveBatches(ctx, id)
	if err != nil {
		return apperrors.ErrInternal.WithMessage("Failed to check active batches")
	}
	if hasActive {
		return apperrors.ErrResourceConflict.WithMessage("Cannot delete item with active batches")
	}

	if err := s.itemRepo.Delete(ctx, id); err != nil {
		return apperrors.ErrInternal.WithMessage("Failed to delete item")
	}
	return nil
}
