package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/utils"
	"gorm.io/gorm"
)

// StockRepository defines the interface for stock and batch data access.
type StockRepository interface {
	// Batch operations
	CreateBatch(ctx context.Context, batch *models.Batch) error
	FindBatchByID(ctx context.Context, id uuid.UUID) (*models.Batch, error)
	FindBatchesByItem(ctx context.Context, itemID uuid.UUID, isDepleted *bool) ([]models.Batch, error)
	FindAvailableBatchesFEFO(ctx context.Context, itemID uuid.UUID) ([]models.Batch, error)
	UpdateBatch(ctx context.Context, batch *models.Batch) error
	GenerateBatchNumber(ctx context.Context, itemID uuid.UUID) (string, error)

	// Stock In operations
	CreateStockIn(ctx context.Context, stockIn *models.StockIn) error
	FindAllStockIn(ctx context.Context, params StockInQueryParams) ([]models.StockIn, int64, error)

	// Stock Out operations
	CreateStockOut(ctx context.Context, stockOut *models.StockOut) error
	CreateStockOutDetail(ctx context.Context, detail *models.StockOutDetail) error
	FindAllStockOut(ctx context.Context, params StockOutQueryParams) ([]models.StockOut, int64, error)
	FindStockOutByID(ctx context.Context, id uuid.UUID) (*models.StockOut, error)
	FindStockOutDetails(ctx context.Context, stockOutID uuid.UUID) ([]models.StockOutDetail, error)
}

// StockInQueryParams holds query parameters for listing stock in.
type StockInQueryParams struct {
	utils.PaginationParams
	ItemID     *uuid.UUID
	SupplierID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
}

// StockOutQueryParams holds query parameters for listing stock out.
type StockOutQueryParams struct {
	utils.PaginationParams
	ItemID    *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
}

type stockRepository struct {
	db *gorm.DB
}

// NewStockRepository creates a new StockRepository.
func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

// --- Batch ---

func (r *stockRepository) CreateBatch(ctx context.Context, batch *models.Batch) error {
	if err := r.db.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("failed to create batch: %w", err)
	}
	return nil
}

func (r *stockRepository) FindBatchByID(ctx context.Context, id uuid.UUID) (*models.Batch, error) {
	var batch models.Batch
	if err := r.db.WithContext(ctx).Preload("Item").Preload("Item.Unit").First(&batch, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("batch not found: %w", err)
	}
	return &batch, nil
}

func (r *stockRepository) FindBatchesByItem(ctx context.Context, itemID uuid.UUID, isDepleted *bool) ([]models.Batch, error) {
	var batches []models.Batch
	query := r.db.WithContext(ctx).Where("item_id = ?", itemID)
	if isDepleted != nil {
		query = query.Where("is_depleted = ?", *isDepleted)
	}
	if err := query.Order("expiry_date ASC").Find(&batches).Error; err != nil {
		return nil, fmt.Errorf("failed to find batches: %w", err)
	}
	return batches, nil
}

// FindAvailableBatchesFEFO returns non-depleted batches ordered by expiry date (FEFO).
func (r *stockRepository) FindAvailableBatchesFEFO(ctx context.Context, itemID uuid.UUID) ([]models.Batch, error) {
	var batches []models.Batch
	if err := r.db.WithContext(ctx).
		Where("item_id = ? AND is_depleted = false AND remaining_quantity > 0", itemID).
		Order("expiry_date ASC"). // FEFO: earliest expiry first
		Find(&batches).Error; err != nil {
		return nil, fmt.Errorf("failed to find FEFO batches: %w", err)
	}
	return batches, nil
}

func (r *stockRepository) UpdateBatch(ctx context.Context, batch *models.Batch) error {
	if err := r.db.WithContext(ctx).Save(batch).Error; err != nil {
		return fmt.Errorf("failed to update batch: %w", err)
	}
	return nil
}

func (r *stockRepository) GenerateBatchNumber(ctx context.Context, itemID uuid.UUID) (string, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Batch{}).Where("item_id = ?", itemID).Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to count batches: %w", err)
	}
	batchNumber := fmt.Sprintf("BATCH-%s-%04d", time.Now().Format("20060102"), count+1)
	return batchNumber, nil
}

// --- Stock In ---

func (r *stockRepository) CreateStockIn(ctx context.Context, stockIn *models.StockIn) error {
	if err := r.db.WithContext(ctx).Create(stockIn).Error; err != nil {
		return fmt.Errorf("failed to create stock in: %w", err)
	}
	return nil
}

func (r *stockRepository) FindAllStockIn(ctx context.Context, params StockInQueryParams) ([]models.StockIn, int64, error) {
	var items []models.StockIn
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StockIn{})
	if params.ItemID != nil {
		query = query.Where("item_id = ?", *params.ItemID)
	}
	if params.SupplierID != nil {
		query = query.Where("supplier_id = ?", *params.SupplierID)
	}
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count stock in: %w", err)
	}

	result := query.
		Preload("Item").
		Preload("Batch").
		Preload("Supplier").
		Preload("Creator").
		Order("created_at DESC").
		Limit(params.Limit + 1).
		Find(&items)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find stock in: %w", result.Error)
	}
	return items, total, nil
}

// --- Stock Out ---

func (r *stockRepository) CreateStockOut(ctx context.Context, stockOut *models.StockOut) error {
	if err := r.db.WithContext(ctx).Create(stockOut).Error; err != nil {
		return fmt.Errorf("failed to create stock out: %w", err)
	}
	return nil
}

func (r *stockRepository) CreateStockOutDetail(ctx context.Context, detail *models.StockOutDetail) error {
	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return fmt.Errorf("failed to create stock out detail: %w", err)
	}
	return nil
}

func (r *stockRepository) FindAllStockOut(ctx context.Context, params StockOutQueryParams) ([]models.StockOut, int64, error) {
	var items []models.StockOut
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StockOut{})
	if params.ItemID != nil {
		query = query.Where("item_id = ?", *params.ItemID)
	}
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count stock out: %w", err)
	}

	result := query.
		Preload("Item").
		Preload("Creator").
		Order("created_at DESC").
		Limit(params.Limit + 1).
		Find(&items)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find stock out: %w", result.Error)
	}
	return items, total, nil
}

func (r *stockRepository) FindStockOutByID(ctx context.Context, id uuid.UUID) (*models.StockOut, error) {
	var stockOut models.StockOut
	if err := r.db.WithContext(ctx).
		Preload("Item").
		Preload("Creator").
		First(&stockOut, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("stock out not found: %w", err)
	}
	return &stockOut, nil
}

func (r *stockRepository) FindStockOutDetails(ctx context.Context, stockOutID uuid.UUID) ([]models.StockOutDetail, error) {
	var details []models.StockOutDetail
	if err := r.db.WithContext(ctx).
		Preload("Batch").
		Where("stock_out_id = ?", stockOutID).
		Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to find stock out details: %w", err)
	}
	return details, nil
}
