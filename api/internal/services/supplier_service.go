package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// SupplierInput holds data for creating/updating a supplier.
type SupplierInput struct {
	Name          string  `json:"name"`
	ContactPerson *string `json:"contact_person"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
}

// SupplierService defines the interface for supplier business logic.
type SupplierService interface {
	List(ctx context.Context, params repositories.SupplierQueryParams) ([]models.Supplier, *utils.Pagination, error)
	Create(ctx context.Context, input SupplierInput) (*models.Supplier, error)
	Update(ctx context.Context, id uuid.UUID, input SupplierInput) (*models.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type supplierService struct {
	repo repositories.SupplierRepository
}

// NewSupplierService creates a new SupplierService.
func NewSupplierService(repo repositories.SupplierRepository) SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) List(ctx context.Context, params repositories.SupplierQueryParams) ([]models.Supplier, *utils.Pagination, error) {
	suppliers, total, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, nil, apperrors.ErrInternal.WithMessage("Failed to fetch suppliers")
	}

	hasMore := len(suppliers) > params.Limit
	if hasMore {
		suppliers = suppliers[:params.Limit]
	}

	var nextCursor *string
	if hasMore && len(suppliers) > 0 {
		cursor := utils.EncodeCursor(map[string]string{"id": suppliers[len(suppliers)-1].ID.String()})
		nextCursor = &cursor
	}

	return suppliers, &utils.Pagination{NextCursor: nextCursor, HasMore: hasMore, Total: total}, nil
}

func (s *supplierService) Create(ctx context.Context, input SupplierInput) (*models.Supplier, error) {
	validationErrors := utils.CollectErrors(utils.ValidateRequired("name", input.Name))
	if input.Email != nil && *input.Email != "" {
		if ve := utils.ValidateEmail(*input.Email); ve != nil {
			validationErrors = append(validationErrors, *ve)
		}
	}
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	supplier := &models.Supplier{
		Name:          input.Name,
		ContactPerson: input.ContactPerson,
		Phone:         input.Phone,
		Email:         input.Email,
		Address:       input.Address,
		IsActive:      true,
	}

	if err := s.repo.Create(ctx, supplier); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to create supplier")
	}
	return supplier, nil
}

func (s *supplierService) Update(ctx context.Context, id uuid.UUID, input SupplierInput) (*models.Supplier, error) {
	supplier, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Supplier not found")
	}

	validationErrors := utils.CollectErrors(utils.ValidateRequired("name", input.Name))
	if len(validationErrors) > 0 {
		return nil, apperrors.ErrValidation.WithDetails(validationErrors)
	}

	supplier.Name = input.Name
	supplier.ContactPerson = input.ContactPerson
	supplier.Phone = input.Phone
	supplier.Email = input.Email
	supplier.Address = input.Address

	if err := s.repo.Update(ctx, supplier); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to update supplier")
	}
	return supplier, nil
}

func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return apperrors.ErrNotFound.WithMessage("Supplier not found")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.ErrInternal.WithMessage("Failed to delete supplier")
	}
	return nil
}
