package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// CategoryInput holds data for creating/updating a category.
type CategoryInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// CategoryService defines the interface for category business logic.
type CategoryService interface {
	List(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	Create(ctx context.Context, input CategoryInput) (*models.Category, error)
	Update(ctx context.Context, id uuid.UUID, input CategoryInput) (*models.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) List(ctx context.Context) ([]models.Category, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to fetch categories")
	}
	return categories, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Category not found")
	}
	return category, nil
}

func (s *categoryService) Create(ctx context.Context, input CategoryInput) (*models.Category, error) {
	// Validate
	validationErrors := utils.CollectErrors(
		utils.ValidateRequired("name", input.Name),
	)
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Check duplicate
	existing, err := s.repo.FindByName(ctx, input.Name)
	if err == nil && existing != nil {
		return nil, apperrors.ErrDuplicateResource.WithMessage("Category with this name already exists")
	}

	category := &models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to create category")
	}

	return category, nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, input CategoryInput) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Category not found")
	}

	// Validate
	validationErrors := utils.CollectErrors(
		utils.ValidateRequired("name", input.Name),
	)
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Check duplicate (only if name changed)
	if category.Name != input.Name {
		existing, err := s.repo.FindByName(ctx, input.Name)
		if err == nil && existing != nil {
			return nil, apperrors.ErrDuplicateResource.WithMessage("Category with this name already exists")
		}
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to update category")
	}

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound.WithMessage("Category not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.ErrInternal.WithMessage("Failed to delete category")
	}
	return nil
}
