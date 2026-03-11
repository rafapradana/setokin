package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/utils"
	"gorm.io/gorm"
)

// SupplierRepository defines the interface for supplier data access.
type SupplierRepository interface {
	FindAll(ctx context.Context, params SupplierQueryParams) ([]models.Supplier, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error)
	Create(ctx context.Context, supplier *models.Supplier) error
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SupplierQueryParams holds query parameters for listing suppliers.
type SupplierQueryParams struct {
	utils.PaginationParams
	IsActive *bool
	Search   string
}

type supplierRepository struct {
	db *gorm.DB
}

// NewSupplierRepository creates a new SupplierRepository.
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) FindAll(ctx context.Context, params SupplierQueryParams) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Supplier{})

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.Search != "" {
		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	result := query.Order("name ASC").Limit(params.Limit + 1).Find(&suppliers)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find suppliers: %w", result.Error)
	}

	return suppliers, total, nil
}

func (r *supplierRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Supplier, error) {
	var supplier models.Supplier
	result := r.db.WithContext(ctx).First(&supplier, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("supplier not found: %w", result.Error)
	}
	return &supplier, nil
}

func (r *supplierRepository) Create(ctx context.Context, supplier *models.Supplier) error {
	result := r.db.WithContext(ctx).Create(supplier)
	if result.Error != nil {
		return fmt.Errorf("failed to create supplier: %w", result.Error)
	}
	return nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	result := r.db.WithContext(ctx).Save(supplier)
	if result.Error != nil {
		return fmt.Errorf("failed to update supplier: %w", result.Error)
	}
	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.Supplier{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete supplier: %w", result.Error)
	}
	return nil
}
