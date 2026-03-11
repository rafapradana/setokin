package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	FindAll(ctx context.Context) ([]models.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)
	Create(ctx context.Context, category *models.Category) error
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	result := r.db.WithContext(ctx).Order("name ASC").Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find categories: %w", result.Error)
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).First(&category, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("category not found: %w", result.Error)
	}
	return &category, nil
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&category)
	if result.Error != nil {
		return nil, fmt.Errorf("category not found: %w", result.Error)
	}
	return &category, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Create(category)
	if result.Error != nil {
		return fmt.Errorf("failed to create category: %w", result.Error)
	}
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Save(category)
	if result.Error != nil {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	return nil
}
