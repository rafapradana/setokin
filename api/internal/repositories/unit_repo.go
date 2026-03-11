package repositories

import (
	"context"
	"fmt"

	"github.com/setokin/api/internal/models"
	"gorm.io/gorm"
)

// UnitRepository defines the interface for unit data access.
type UnitRepository interface {
	FindAll(ctx context.Context) ([]models.Unit, error)
}

type unitRepository struct {
	db *gorm.DB
}

// NewUnitRepository creates a new UnitRepository.
func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) FindAll(ctx context.Context) ([]models.Unit, error) {
	var units []models.Unit
	result := r.db.WithContext(ctx).Order("type ASC, name ASC").Find(&units)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find units: %w", result.Error)
	}
	return units, nil
}
