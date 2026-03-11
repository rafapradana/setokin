// Package models defines GORM models for the Setokin database.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the data.users table.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FullName     string    `gorm:"not null" json:"full_name"`
	Role         string    `gorm:"not null;default:staff" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "data.users"
}

// IsOwner returns true if the user has the owner role.
func (u *User) IsOwner() bool {
	return u.Role == "owner"
}

// IsManager returns true if the user has the manager role.
func (u *User) IsManager() bool {
	return u.Role == "manager"
}

// HasManagementAccess returns true if the user has owner or manager role.
func (u *User) HasManagementAccess() bool {
	return u.IsOwner() || u.IsManager()
}
