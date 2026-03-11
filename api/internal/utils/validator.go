package utils

import (
	"regexp"
	"unicode"

	apperrors "github.com/setokin/api/pkg/errors"
)

var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// ValidateEmail checks if the given string is a valid email address.
func ValidateEmail(email string) *apperrors.ValidationError {
	if email == "" {
		return &apperrors.ValidationError{Field: "email", Message: "Email is required"}
	}
	if !emailRegex.MatchString(email) {
		return &apperrors.ValidationError{Field: "email", Message: "Invalid email format"}
	}
	return nil
}

// ValidatePassword checks if the password meets minimum requirements.
// Must be at least 8 characters, contain uppercase, lowercase, and a digit.
func ValidatePassword(password string) *apperrors.ValidationError {
	if len(password) < 8 {
		return &apperrors.ValidationError{Field: "password", Message: "Password must be at least 8 characters"}
	}

	var hasUpper, hasLower, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
	}

	if !hasUpper {
		return &apperrors.ValidationError{Field: "password", Message: "Password must contain at least one uppercase letter"}
	}
	if !hasLower {
		return &apperrors.ValidationError{Field: "password", Message: "Password must contain at least one lowercase letter"}
	}
	if !hasDigit {
		return &apperrors.ValidationError{Field: "password", Message: "Password must contain at least one digit"}
	}
	return nil
}

// ValidateRole checks if the role is valid.
func ValidateRole(role string) *apperrors.ValidationError {
	validRoles := map[string]bool{"owner": true, "manager": true, "staff": true}
	if role != "" && !validRoles[role] {
		return &apperrors.ValidationError{Field: "role", Message: "Role must be owner, manager, or staff"}
	}
	return nil
}

// ValidateRequired checks if a required string field is non-empty.
func ValidateRequired(field, value string) *apperrors.ValidationError {
	if value == "" {
		return &apperrors.ValidationError{Field: field, Message: field + " is required"}
	}
	return nil
}

// CollectErrors takes a list of validation results and returns non-nil ones.
func CollectErrors(checks ...*apperrors.ValidationError) []apperrors.ValidationError {
	var errs []apperrors.ValidationError
	for _, check := range checks {
		if check != nil {
			errs = append(errs, *check)
		}
	}
	return errs
}
