package models

import (
	"time"

	"github.com/google/uuid"
)

// Unit represents a unit of measurement in data.units.
type Unit struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"uniqueIndex;not null" json:"name"`
	Abbreviation string    `gorm:"uniqueIndex;not null" json:"abbreviation"`
	Type         string    `gorm:"not null" json:"type"` // weight, volume, count
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// TableName overrides the default table name.
func (Unit) TableName() string {
	return "data.units"
}
