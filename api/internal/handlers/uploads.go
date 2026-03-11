package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// UploadHandler handles file upload endpoints (MinIO presigned URLs).
type UploadHandler struct {
	service services.UploadService
}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler(service services.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// RequestUpload handles POST /uploads/request.
func (h *UploadHandler) RequestUpload(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, apperrors.ErrAuthenticationRequired)
	}

	var input services.UploadRequestInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid request body"))
	}

	upload, presignedURL, err := h.service.RequestUpload(c.Context(), input, userID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"upload_id":     upload.ID,
		"upload_url":    presignedURL,
		"object_key":    upload.ObjectKey,
		"expires_in":    900, // 15 minutes in seconds
		"instructions":  "PUT the file to upload_url within 15 minutes, then call POST /uploads/:upload_id/confirm",
	})
}

// ConfirmUpload handles POST /uploads/:upload_id/confirm.
func (h *UploadHandler) ConfirmUpload(c *fiber.Ctx) error {
	uploadID, err := uuid.Parse(c.Params("upload_id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid upload ID"))
	}

	upload, err := h.service.ConfirmUpload(c.Context(), uploadID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, upload)
}

// GetDownloadURL handles GET /uploads/:upload_id/download.
func (h *UploadHandler) GetDownloadURL(c *fiber.Ctx) error {
	uploadID, err := uuid.Parse(c.Params("upload_id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid upload ID"))
	}

	downloadURL, err := h.service.GetDownloadURL(c.Context(), uploadID)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return utils.ErrorResponse(c, appErr)
		}
		return utils.ErrorResponse(c, apperrors.ErrInternal)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"download_url": downloadURL,
		"expires_in":   3600,
	})
}
