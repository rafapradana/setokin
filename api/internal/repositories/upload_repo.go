package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"gorm.io/gorm"
)

// UploadRepository defines the interface for upload data access.
type UploadRepository interface {
	Create(ctx context.Context, upload *models.Upload) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Upload, error)
	Update(ctx context.Context, upload *models.Upload) error
}

type uploadRepository struct {
	db *gorm.DB
}

// NewUploadRepository creates a new UploadRepository.
func NewUploadRepository(db *gorm.DB) UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) Create(ctx context.Context, upload *models.Upload) error {
	if err := r.db.WithContext(ctx).Create(upload).Error; err != nil {
		return fmt.Errorf("failed to create upload: %w", err)
	}
	return nil
}

func (r *uploadRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Upload, error) {
	var upload models.Upload
	if err := r.db.WithContext(ctx).First(&upload, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("upload not found: %w", err)
	}
	return &upload, nil
}

func (r *uploadRepository) Update(ctx context.Context, upload *models.Upload) error {
	if err := r.db.WithContext(ctx).Save(upload).Error; err != nil {
		return fmt.Errorf("failed to update upload: %w", err)
	}
	return nil
}
