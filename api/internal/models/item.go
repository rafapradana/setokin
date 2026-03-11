package models

import (
	"time"

	"github.com/google/uuid"
)

// Item represents an inventory item in data.items.
type Item struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	CategoryID   uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	UnitID       uuid.UUID `gorm:"type:uuid;not null" json:"unit_id"`
	MinimumStock float64   `gorm:"type:numeric(10,3);default:0" json:"minimum_stock"`
	Description  *string   `json:"description,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;default:now()" json:"updated_at"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Unit     Unit     `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

// TableName overrides the default table name.
func (Item) TableName() string {
	return "data.items"
}
