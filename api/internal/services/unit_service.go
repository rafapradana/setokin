package services

import (
	"context"

	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	apperrors "github.com/setokin/api/pkg/errors"
)

// UnitService defines the interface for unit business logic.
type UnitService interface {
	List(ctx context.Context) ([]models.Unit, error)
}

type unitService struct {
	repo repositories.UnitRepository
}

// NewUnitService creates a new UnitService.
func NewUnitService(repo repositories.UnitRepository) UnitService {
	return &unitService{repo: repo}
}

func (s *unitService) List(ctx context.Context) ([]models.Unit, error) {
	units, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to fetch units")
	}
	return units, nil
}
