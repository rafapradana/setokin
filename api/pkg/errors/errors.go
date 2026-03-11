// Package errors provides custom error types for the Setokin API.
// Error codes match the API documentation.
package errors

import "fmt"

// AppError represents a structured application error.
type AppError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Status  int               `json:"-"`
	Details []ValidationError `json:"details,omitempty"`
}

// ValidationError represents a single field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithMessage returns a copy of the error with a custom message.
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: msg,
		Status:  e.Status,
		Details: e.Details,
	}
}

// WithDetails returns a copy of the error with validation details.
func (e *AppError) WithDetails(details []ValidationError) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Status:  e.Status,
		Details: details,
	}
}

// Predefined errors matching API documentation error codes.
var (
	ErrValidation = &AppError{
		Code:    "validation_error",
		Message: "Validation failed",
		Status:  422,
	}

	ErrAuthenticationRequired = &AppError{
		Code:    "authentication_required",
		Message: "Authentication required",
		Status:  401,
	}

	ErrTokenExpired = &AppError{
		Code:    "token_expired",
		Message: "Access token expired",
		Status:  401,
	}

	ErrTokenInvalid = &AppError{
		Code:    "token_invalid",
		Message: "Invalid token",
		Status:  401,
	}

	ErrInsufficientPermissions = &AppError{
		Code:    "insufficient_permissions",
		Message: "Insufficient permissions",
		Status:  403,
	}

	ErrNotFound = &AppError{
		Code:    "resource_not_found",
		Message: "Resource not found",
		Status:  404,
	}

	ErrDuplicateResource = &AppError{
		Code:    "duplicate_resource",
		Message: "Resource already exists",
		Status:  409,
	}

	ErrResourceConflict = &AppError{
		Code:    "resource_conflict",
		Message: "Resource conflict",
		Status:  409,
	}

	ErrInsufficientStock = &AppError{
		Code:    "insufficient_stock",
		Message: "Not enough stock for operation",
		Status:  422,
	}

	ErrBatchDepleted = &AppError{
		Code:    "batch_depleted",
		Message: "Batch is already depleted",
		Status:  422,
	}

	ErrInvalidQuantity = &AppError{
		Code:    "invalid_quantity",
		Message: "Invalid quantity value",
		Status:  422,
	}

	ErrRateLimitExceeded = &AppError{
		Code:    "rate_limit_exceeded",
		Message: "Too many requests",
		Status:  429,
	}

	ErrInternal = &AppError{
		Code:    "internal_error",
		Message: "Internal server error",
		Status:  500,
	}
)
