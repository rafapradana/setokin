package models

import (
	"time"

	"github.com/google/uuid"
)

// Supplier represents a supplier in data.suppliers.
type Supplier struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	ContactPerson *string   `json:"contact_person,omitempty"`
	Phone         *string   `json:"phone,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Address       *string   `json:"address,omitempty"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// TableName overrides the default table name.
func (Supplier) TableName() string {
	return "data.suppliers"
}
