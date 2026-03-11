package models

import (
	"time"

	"github.com/google/uuid"
)

// Category represents an item category in data.categories.
type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName overrides the default table name.
func (Category) TableName() string {
	return "data.categories"
}
