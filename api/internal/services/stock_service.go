package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
	"gorm.io/gorm"
)

// StockInInput holds data for creating a stock in transaction.
type StockInInput struct {
	ItemID        uuid.UUID  `json:"item_id"`
	Quantity      float64    `json:"quantity"`
	PurchaseDate  string     `json:"purchase_date"`  // YYYY-MM-DD
	ExpiryDate    string     `json:"expiry_date"`    // YYYY-MM-DD
	SupplierID    *uuid.UUID `json:"supplier_id"`
	PurchasePrice *float64   `json:"purchase_price"`
	Notes         *string    `json:"notes"`
}

// StockOutInput holds data for creating a stock out transaction.
type StockOutInput struct {
	ItemID    uuid.UUID `json:"item_id"`
	Quantity  float64   `json:"quantity"`
	UsageDate string    `json:"usage_date"` // YYYY-MM-DD
	Notes     *string   `json:"notes"`
}

// DeductionResult represents a single FEFO batch deduction.
type DeductionResult struct {
	Batch          models.Batch `json:"batch"`
	QuantityUsed   float64      `json:"quantity_used"`
	RemainingAfter float64      `json:"remaining_after"`
}

// StockService defines the interface for stock business logic.
type StockService interface {
	// Stock In
	CreateStockIn(ctx context.Context, input StockInInput, userID uuid.UUID) (*models.StockIn, *models.Batch, error)
	ListStockIn(ctx context.Context, params repositories.StockInQueryParams) ([]models.StockIn, *utils.Pagination, error)

	// Stock Out (FEFO)
	CreateStockOut(ctx context.Context, input StockOutInput, userID uuid.UUID) (*models.StockOut, []DeductionResult, float64, error)
	ListStockOut(ctx context.Context, params repositories.StockOutQueryParams) ([]models.StockOut, *utils.Pagination, error)
	GetStockOutDetails(ctx context.Context, id uuid.UUID) (*models.StockOut, []models.StockOutDetail, error)

	// Batches
	ListBatchesByItem(ctx context.Context, itemID uuid.UUID, isDepleted *bool) ([]models.Batch, error)
	GetBatch(ctx context.Context, id uuid.UUID) (*models.Batch, []models.StockOutDetail, error)
}

type stockService struct {
	stockRepo repositories.StockRepository
	itemRepo  repositories.ItemRepository
	db        *gorm.DB
}

// NewStockService creates a new StockService.
func NewStockService(stockRepo repositories.StockRepository, itemRepo repositories.ItemRepository, db *gorm.DB) StockService {
	return &stockService{stockRepo: stockRepo, itemRepo: itemRepo, db: db}
}

// --- Stock In ---

func (s *stockService) CreateStockIn(ctx context.Context, input StockInInput, userID uuid.UUID) (*models.StockIn, *models.Batch, error) {
	// Validate
	var validationErrors []apperrors.ValidationError
	if input.ItemID == uuid.Nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "item_id", Message: "Item is required"})
	}
	if input.Quantity <= 0 {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "quantity", Message: "Quantity must be greater than 0"})
	}

	purchaseDate, err := time.Parse("2006-01-02", input.PurchaseDate)
	if err != nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "purchase_date", Message: "Invalid date format (YYYY-MM-DD)"})
	}

	expiryDate, err := time.Parse("2006-01-02", input.ExpiryDate)
	if err != nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "expiry_date", Message: "Invalid date format (YYYY-MM-DD)"})
	} else if !expiryDate.After(time.Now().Truncate(24*time.Hour)) {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "expiry_date", Message: "Expiry date must be in the future"})
	}

	if input.PurchasePrice != nil && *input.PurchasePrice < 0 {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "purchase_price", Message: "Purchase price must be >= 0"})
	}
	if len(validationErrors) > 0 {
		return nil, nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Verify item exists
	if _, err := s.itemRepo.FindByID(ctx, input.ItemID); err != nil {
		return nil, nil, apperrors.ErrNotFound.WithMessage("Item not found")
	}

	// Generate batch number
	batchNumber, err := s.stockRepo.GenerateBatchNumber(ctx, input.ItemID)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to generate batch number")
	}

	// Create batch
	batch := &models.Batch{
		ItemID:            input.ItemID,
		BatchNumber:       batchNumber,
		InitialQuantity:   input.Quantity,
		RemainingQuantity: input.Quantity,
		ExpiryDate:        expiryDate,
		IsDepleted:        false,
	}
	if err := s.stockRepo.CreateBatch(ctx, batch); err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to create batch")
	}

	// Create stock in record
	stockIn := &models.StockIn{
		ItemID:        input.ItemID,
		BatchID:       batch.ID,
		SupplierID:    input.SupplierID,
		Quantity:      input.Quantity,
		PurchaseDate:  purchaseDate,
		PurchasePrice: input.PurchasePrice,
		Notes:         input.Notes,
		CreatedBy:     userID,
	}
	if err := s.stockRepo.CreateStockIn(ctx, stockIn); err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to create stock in")
	}

	return stockIn, batch, nil
}

func (s *stockService) ListStockIn(ctx context.Context, params repositories.StockInQueryParams) ([]models.StockIn, *utils.Pagination, error) {
	items, total, err := s.stockRepo.FindAllStockIn(ctx, params)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch stock in records")
	}

	hasMore := len(items) > params.Limit
	if hasMore {
		items = items[:params.Limit]
	}
	var nextCursor *string
	if hasMore && len(items) > 0 {
		c := utils.EncodeCursor(map[string]string{"id": items[len(items)-1].ID.String()})
		nextCursor = &c
	}
	return items, &utils.Pagination{NextCursor: nextCursor, HasMore: hasMore, Total: total}, nil
}

// --- Stock Out (FEFO) ---

func (s *stockService) CreateStockOut(ctx context.Context, input StockOutInput, userID uuid.UUID) (*models.StockOut, []DeductionResult, float64, error) {
	// Validate
	var validationErrors []apperrors.ValidationError
	if input.ItemID == uuid.Nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "item_id", Message: "Item is required"})
	}
	if input.Quantity <= 0 {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "quantity", Message: "Quantity must be greater than 0"})
	}
	usageDate, err := time.Parse("2006-01-02", input.UsageDate)
	if err != nil {
		validationErrors = append(validationErrors, apperrors.ValidationError{Field: "usage_date", Message: "Invalid date format (YYYY-MM-DD)"})
	}
	if len(validationErrors) > 0 {
		return nil, nil, 0, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Verify item exists and get unit
	item, err := s.itemRepo.FindByID(ctx, input.ItemID)
	if err != nil {
		return nil, nil, 0, apperrors.ErrNotFound.WithMessage("Item not found")
	}

	// Get available batches ordered by FEFO (earliest expiry first)
	batches, err := s.stockRepo.FindAvailableBatchesFEFO(ctx, input.ItemID)
	if err != nil {
		return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to fetch batches")
	}

	// Calculate total available stock
	var totalAvailable float64
	for _, b := range batches {
		totalAvailable += b.RemainingQuantity
	}

	if totalAvailable < input.Quantity {
		return nil, nil, 0, apperrors.ErrInsufficientStock.WithMessage(
			fmt.Sprintf("Insufficient stock. Available: %.3f %s, Requested: %.3f %s",
				totalAvailable, item.Unit.Abbreviation, input.Quantity, item.Unit.Abbreviation))
	}

	// Begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to begin transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create stock out record
	stockOut := &models.StockOut{
		ItemID:    input.ItemID,
		Quantity:  input.Quantity,
		UsageDate: usageDate,
		Notes:     input.Notes,
		CreatedBy: userID,
	}
	if err := tx.Create(stockOut).Error; err != nil {
		tx.Rollback()
		return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to create stock out")
	}

	// FEFO deduction
	remaining := input.Quantity
	var deductions []DeductionResult

	for i := range batches {
		if remaining <= 0 {
			break
		}

		batch := &batches[i]
		deductQty := remaining
		if deductQty > batch.RemainingQuantity {
			deductQty = batch.RemainingQuantity
		}

		batch.RemainingQuantity -= deductQty
		if batch.RemainingQuantity == 0 {
			batch.IsDepleted = true
		}

		if err := tx.Save(batch).Error; err != nil {
			tx.Rollback()
			return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to update batch")
		}

		detail := &models.StockOutDetail{
			StockOutID:   stockOut.ID,
			BatchID:      batch.ID,
			QuantityUsed: deductQty,
		}
		if err := tx.Create(detail).Error; err != nil {
			tx.Rollback()
			return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to create deduction detail")
		}

		deductions = append(deductions, DeductionResult{
			Batch:          *batch,
			QuantityUsed:   deductQty,
			RemainingAfter: batch.RemainingQuantity,
		})

		remaining -= deductQty
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, 0, apperrors.ErrInternal.WithMessage("Failed to commit transaction")
	}

	remainingStock := totalAvailable - input.Quantity
	return stockOut, deductions, remainingStock, nil
}

func (s *stockService) ListStockOut(ctx context.Context, params repositories.StockOutQueryParams) ([]models.StockOut, *utils.Pagination, error) {
	items, total, err := s.stockRepo.FindAllStockOut(ctx, params)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch stock out records")
	}

	hasMore := len(items) > params.Limit
	if hasMore {
		items = items[:params.Limit]
	}
	var nextCursor *string
	if hasMore && len(items) > 0 {
		c := utils.EncodeCursor(map[string]string{"id": items[len(items)-1].ID.String()})
		nextCursor = &c
	}
	return items, &utils.Pagination{NextCursor: nextCursor, HasMore: hasMore, Total: total}, nil
}

func (s *stockService) GetStockOutDetails(ctx context.Context, id uuid.UUID) (*models.StockOut, []models.StockOutDetail, error) {
	stockOut, err := s.stockRepo.FindStockOutByID(ctx, id)
	if err != nil {
		return nil, nil, apperrors.ErrNotFound.WithMessage("Stock out not found")
	}
	details, err := s.stockRepo.FindStockOutDetails(ctx, id)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch deduction details")
	}
	return stockOut, details, nil
}

// --- Batches ---

func (s *stockService) ListBatchesByItem(ctx context.Context, itemID uuid.UUID, isDepleted *bool) ([]models.Batch, error) {
	batches, err := s.stockRepo.FindBatchesByItem(ctx, itemID, isDepleted)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to fetch batches")
	}
	return batches, nil
}

func (s *stockService) GetBatch(ctx context.Context, id uuid.UUID) (*models.Batch, []models.StockOutDetail, error) {
	batch, err := s.stockRepo.FindBatchByID(ctx, id)
	if err != nil {
		return nil, nil, apperrors.ErrNotFound.WithMessage("Batch not found")
	}

	// Get usage history for this batch
	var details []models.StockOutDetail
	if err := s.db.WithContext(ctx).
		Preload("Batch").
		Where("batch_id = ?", id).
		Order("created_at DESC").
		Find(&details).Error; err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch batch usage history")
	}

	return batch, details, nil
}
