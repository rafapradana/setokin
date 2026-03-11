package services

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/setokin/api/internal/minio"
	"github.com/setokin/api/internal/models"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// UploadRequestInput holds data for requesting a presigned upload URL.
type UploadRequestInput struct {
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	EntityType string `json:"entity_type"` // "items", "suppliers", etc.
	EntityID   string `json:"entity_id"`
}

// UploadService defines the interface for upload business logic.
type UploadService interface {
	RequestUpload(ctx context.Context, input UploadRequestInput, userID uuid.UUID) (*models.Upload, string, error)
	ConfirmUpload(ctx context.Context, uploadID uuid.UUID) (*models.Upload, error)
	GetDownloadURL(ctx context.Context, uploadID uuid.UUID) (string, error)
}

type uploadService struct {
	repo        repositories.UploadRepository
	minioClient *minio.Client
}

// NewUploadService creates a new UploadService.
func NewUploadService(repo repositories.UploadRepository, minioClient *minio.Client) UploadService {
	return &uploadService{repo: repo, minioClient: minioClient}
}

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".pdf": true, ".doc": true, ".docx": true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true, "image/png": true, "image/gif": true,
	"image/webp": true, "application/pdf": true,
}

func (s *uploadService) RequestUpload(ctx context.Context, input UploadRequestInput, userID uuid.UUID) (*models.Upload, string, error) {
	// Validate
	validationErrors := utils.CollectErrors(
		utils.ValidateRequired("file_name", input.FileName),
		utils.ValidateRequired("file_type", input.FileType),
		utils.ValidateRequired("entity_type", input.EntityType),
	)

	ext := strings.ToLower(filepath.Ext(input.FileName))
	if !allowedExtensions[ext] {
		validationErrors = append(validationErrors, apperrors.ValidationError{
			Field: "file_name", Message: "File type not allowed. Allowed: jpg, jpeg, png, gif, webp, pdf, doc, docx",
		})
	}
	if !allowedMimeTypes[input.FileType] && !strings.HasPrefix(input.FileType, "application/") {
		validationErrors = append(validationErrors, apperrors.ValidationError{
			Field: "file_type", Message: "MIME type not allowed",
		})
	}
	if len(validationErrors) > 0 {
		return nil, "", apperrors.ErrValidation.WithDetails(validationErrors)
	}

	// Generate object key
	objectKey := minio.GenerateObjectKey(input.EntityType, input.FileName)

	// Create upload record
	upload := &models.Upload{
		FileName:   input.FileName,
		FileType:   input.FileType,
		ObjectKey:  objectKey,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Status:     "pending",
		UploadedBy: userID,
	}
	if err := s.repo.Create(ctx, upload); err != nil {
		return nil, "", apperrors.ErrInternal.WithMessage("Failed to create upload record")
	}

	// Generate presigned upload URL (15 minutes)
	presignedURL, err := s.minioClient.GeneratePresignedUploadURL(ctx, objectKey, 15*time.Minute)
	if err != nil {
		return nil, "", apperrors.ErrInternal.WithMessage("Failed to generate upload URL")
	}

	return upload, presignedURL, nil
}

func (s *uploadService) ConfirmUpload(ctx context.Context, uploadID uuid.UUID) (*models.Upload, error) {
	upload, err := s.repo.FindByID(ctx, uploadID)
	if err != nil {
		return nil, apperrors.ErrNotFound.WithMessage("Upload not found")
	}

	if upload.Status != "pending" {
		return nil, apperrors.ErrResourceConflict.WithMessage("Upload already confirmed")
	}

	// Verify file exists in MinIO
	exists, err := s.minioClient.ObjectExists(ctx, upload.ObjectKey)
	if err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to verify upload")
	}
	if !exists {
		return nil, apperrors.ErrNotFound.WithMessage("File not found in storage. Please upload the file first.")
	}

	upload.Status = "confirmed"
	if err := s.repo.Update(ctx, upload); err != nil {
		return nil, apperrors.ErrInternal.WithMessage("Failed to confirm upload")
	}

	return upload, nil
}

func (s *uploadService) GetDownloadURL(ctx context.Context, uploadID uuid.UUID) (string, error) {
	upload, err := s.repo.FindByID(ctx, uploadID)
	if err != nil {
		return "", apperrors.ErrNotFound.WithMessage("Upload not found")
	}

	if upload.Status != "confirmed" {
		return "", apperrors.ErrResourceConflict.WithMessage("Upload not yet confirmed")
	}

	// Generate presigned download URL (1 hour)
	downloadURL, err := s.minioClient.GeneratePresignedDownloadURL(ctx, upload.ObjectKey, 1*time.Hour)
	if err != nil {
		return "", apperrors.ErrInternal.WithMessage("Failed to generate download URL")
	}

	return downloadURL, nil
}
