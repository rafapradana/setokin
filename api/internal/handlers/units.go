package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// UnitHandler handles unit endpoints.
type UnitHandler struct {
	service services.UnitService
}

// NewUnitHandler creates a new UnitHandler.
func NewUnitHandler(service services.UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

// List handles GET /units.
func (h *UnitHandler) List(c *fiber.Ctx) error {
	units, err := h.service.List(c.Context())
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, units)
}
